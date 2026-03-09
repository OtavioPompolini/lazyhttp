# Plan: Full Vim Motion Engine with Lua Config

## Context

`RequestDetailsPane` currently has a hand-coded flat switch for vim keys — a `motionToKey()` helper emitting synthetic `tea.KeyPressMsg`s and a single `pendingG bool`. This cannot scale to counts (`2j`), operators+motions (`dw`, `ciw`), visual modes, or user-configurable bindings.

The goal is a proper vim engine (`internal/vim/` package) with a Lua config at `~/.config/lazyhttp/init.lua` using Neovim's API style (`vim.keymap.set`, `vim.g.mapleader`).

**New external dep**: `github.com/yuin/gopher-lua` — pure-Go Lua 5.1 VM, no CGo.

---

## Architecture

```
~/.config/lazyhttp/init.lua
        ↓  gopher-lua at startup
  vim.Engine  (internal/vim/)
        ↓  ProcessKey(key) → (Action, count, passThru)
  RequestDetailsPane.executeAction(action, count)
        ↓  textarea direct API + text manipulation
  textarea.Model  (charm.land/bubbles/v2)
```

**Key design choice**: Pre-composed operator+motion bindings (e.g., `"ciw" → ActionChangeInnerWord`) rather than dynamic operator-pending state machine. This is appropriate for a REST client editor and makes the keymap user-customisable without grammar knowledge.

---

## New Package: `internal/vim/`

### `mode.go`
```go
type Mode int
const (
    ModeNormal Mode = iota
    ModeInsert
    ModeVisual       // v — character-wise
    ModeVisualLine   // V — line-wise
    ModeVisualBlock  // ctrl+v — column block
)
func (m Mode) String() string  // "NORMAL" / "INSERT" / "VISUAL" / "V-LINE" / "V-BLOCK"
```

### `action.go`
Full enum of semantic operations:
- **Motions**: `ActionMoveLeft/Right/Up/Down`, `ActionMoveWordForward/End/Back`, `ActionMoveWORD*`, `ActionMoveLineStart/First/End`, `ActionMoveDocStart/End`
- **Insert-entry**: `ActionInsertBefore/After/LineStart/LineEnd/NewLineBelow/NewLineAbove`
- **Pre-composed text ops**: `ActionDeleteLine` (dd), `ActionDeleteWord` (dw), `ActionDeleteToLineEnd` (D), `ActionChangeLine` (cc), `ActionChangeWord` (cw), `ActionChangeToLineEnd` (C), `ActionChangeInnerWord` (ciw), `ActionChangeInnerWORD` (ciW), `ActionChangeInnerQuoteS` (ci'), `ActionChangeInnerQuoteD` (ci"), `ActionChangeInnerQuoteB` (ci`), `ActionChangeInnerParen` (ci()), `ActionChangeInnerBrace` (ci{}), `ActionChangeInnerBracket` (ci[]), `ActionDeleteInnerWord` (diw), `ActionDeleteInnerQuoteS/D/Paren/Brace/Bracket`, `ActionYankLine` (yy), `ActionYankWord` (yw), `ActionYankInnerWord` (yiw)
- **Clipboard**: `ActionPasteAfter` (p), `ActionPasteBefore` (P)
- **Undo/Redo**: `ActionUndo` (u), `ActionRedo` (ctrl+r)
- **Visual**: `ActionVisualChar/Line/Block`, `ActionVisualDelete/Yank/Change`
- **Control**: `ActionEscape` (→ Normal from insert/visual), `ActionFocusOut` (esc in Normal → exit pane)

```go
// ActionName maps Lua rhs strings to Action values, e.g. "delete_line" → ActionDeleteLine
var ActionName = map[string]Action{ ... }  // one entry per Action
```

### `keymap.go`
```go
type KeyMap map[string]Action  // "dd" → ActionDeleteLine, "ciw" → ActionChangeInnerWord
func DefaultKeyMap(mode Mode) KeyMap
```

Default normal map: all motions (h/j/k/l/w/e/b/W/E/B/0/^/$/gg/G), all insert-entries (i/a/I/A/o/O), all pre-composed delete/change/yank variants, p/P, u/ctrl+r, v/V/ctrl+v, `esc → ActionFocusOut`.

Default visual map: same motions + d/y/c (visual ops) + esc/v/V → ActionEscape.

Default insert map: only `esc → ActionEscape`. Everything else passes through to textarea.

**Note**: Default maps never bind bare `" "` (space) so it is safe as `mapleader`.

### `parser.go`
```go
type Parser struct {
    buf       string   // accumulated key sequence (post-leader substitution)
    count     int      // accumulated numeric count (0 = no count typed)
    leaderKey string
}

// Feed processes one raw key string against km.
// Returns (action, count, pending, reset):
//   action  — resolved Action or ActionNone
//   count   — resolved count (≥1 after resolution)
//   pending — buf is a valid prefix; wait for more keys
//   reset   — sequence discarded (no match and no prefix)
func (p *Parser) Feed(key string, km KeyMap) (Action, int, bool, bool)

func (p *Parser) Reset()
func (p *Parser) HasPending() bool
func (p *Parser) PendingBuf() string
```

**Count/0 ambiguity**: `0` accumulates into count only when `count > 0`. When `count==0` and `buf==""`, `0` appends to buf → matches `ActionMoveLineStart`. This is exact vim behaviour.

**Leader**: Before appending to buf, if `leaderKey != ""` and key equals the leader, append the leader string (e.g. `" "`). The keymap then has entries like `" s"` for `<leader>s`.

**Prefix check**: `strings.HasPrefix(binding, buf)` scanned over all keymap entries (~60). No trie needed.

### `textobject.go`
Pure string functions — no TUI deps. Used by `executeAction` to compute byte ranges before `SetValue`.

```go
// Word boundary functions (use github.com/rivo/uniseg for grapheme-aware boundaries)
func WordBounds(text string, cursorByte int, inner bool) (start, end int)
func WORDBounds(text string, cursorByte int, inner bool) (start, end int)
func LineBounds(text string, cursorByte int, inner bool) (start, end int)
func QuoteBounds(text string, cursorByte int, quote rune, inner bool) (start, end int, ok bool)
func BracketBounds(text string, cursorByte int, open rune, inner bool) (start, end int, ok bool)

// Word navigation
func WordForwardByte(text string, cursorByte int) int
func WordEndByte(text string, cursorByte int) int
func WordBackByte(text string, cursorByte int) int

// Byte ↔ line/col conversion (bridges textarea.Line()/LineInfo() with byte offsets)
func ByteOffsetFromLineCol(text string, line, col int) int
func LineColFromByteOffset(text string, offset int) (line, col int)
```

**`ciw` on whitespace**: If `text[cursorByte]` is whitespace, `WordBounds` returns the contiguous whitespace run — matches vim exactly.

Write comprehensive **table-driven unit tests** for this file first (cursor at BOL, EOL, empty line, Unicode, nested brackets, whitespace edge cases).

### `engine.go`
```go
type Engine struct {
    Mode                   Mode
    AnchorLine, AnchorCol  int       // visual selection anchor (set on visual mode entry)
    Register               string    // unnamed yank register
    undoStack, redoStack   []string  // full-value snapshots, capped at 100
    parser                 Parser
    normalMap, insertMap   KeyMap
    visualMap              KeyMap    // shared for all 3 visual sub-modes
}

func NewEngine() Engine   // DefaultKeyMap for all modes

// ProcessKey routes through the correct keymap. Returns:
//   passThru=true → forward original KeyPressMsg to textarea.Update (insert mode keys)
func (e *Engine) ProcessKey(key string) (action Action, count int, passThru bool)

func (e *Engine) SetLeader(key string)
func (e *Engine) BindKey(mode Mode, seq string, action Action)
// EnterMode transitions mode, resets parser, sets anchor for visual modes
func (e *Engine) EnterMode(mode Mode, curLine, curCol int)
func (e *Engine) PushUndo(value string)
func (e *Engine) PopUndo() (string, bool)
func (e *Engine) PopRedo() (string, bool)
func (e *Engine) PendingBuf() string  // delegates to parser
```

### `config.go`
```go
// LoadConfig executes path as Lua 5.1 (gopher-lua) and applies bindings to engine.
// Returns nil if path does not exist. Logs warnings for unknown action names (non-fatal).
func LoadConfig(path string, engine *Engine) error
```

**Lua VM lifecycle**: Create `lua.LState`, execute config, `Close()`. Not kept alive post-init.

**Lua API**:
```lua
-- ~/.config/lazyhttp/init.lua
vim.g.mapleader = " "

vim.keymap.set("n", "H", "move_line_start")
vim.keymap.set("n", "<leader>s", "save")
vim.keymap.set({"n", "v"}, "<leader>p", "paste_after")
```

`vim.g` uses a `__newindex` metatable: assignment to `mapleader` calls `engine.SetLeader(value)`; all other writes are silently accepted.

Mode strings: `"n"` → Normal, `"i"` → Insert, `"v"` → all three visual modes.

`<leader>` in lhs is expanded with the current leader at bind time.

---

## Modified Files

### `internal/tui/panes/request_details.go`

**Struct** (replace `mode vimMode`, `pendingG bool`):
```go
type RequestDetailsPane struct {
    focused        bool
    width, height  int
    engine         vim.Engine
    textarea       textarea.Model
    currentRequest *types.Request
    requestSystem  *state.RequestManager
}
```

**Constructor**: `NewRequestDetailsPane(st *state.State, engine vim.Engine) RequestDetailsPane`

**Update()** key handling:
```go
case tea.KeyPressMsg:
    if !p.focused { break }
    action, count, passThru := p.engine.ProcessKey(m.String())
    if passThru {
        var cmd tea.Cmd
        p.textarea, cmd = p.textarea.Update(msg)
        return p, cmd
    }
    if action != vim.ActionNone {
        return p.executeAction(action, count)
    }
    // pending sequence — consume silently
```

**`executeAction(action, count)`** — pointer receiver method. Key patterns:

| Action group | Implementation |
|---|---|
| Insert-entry | `p.engine.EnterMode(ModeInsert, ...)` + cursor reposition, return `textarea.Blink` |
| Pure motions | `CursorUp/Down()` ×count (clamped), `CursorStart/End()`, `MoveToBegin/End()` |
| Word motions | `textobject.WordForwardByte()` etc. → `LineColFromByteOffset` → `repositionCursor()` |
| Text mutations | `PushUndo(Value())`, compute byte range via textobject, `SetValue(newText)`, `repositionCursor()` |
| Visual entry | `p.engine.EnterMode(ModeVisual*, line, col)` |
| Visual ops | Read selection from anchor+cursor → mutate or copy to Register |
| Paste | Insert `engine.Register` at cursor via SetValue + reposition |
| Undo/Redo | `PopUndo/Redo()` → `SetValue(prev)` |
| ActionEscape | `p.engine.EnterMode(ModeNormal, ...)` |
| ActionFocusOut | `return p, msgs.FocusCmd(msgs.FocusRequests)` |

**Cursor reposition after SetValue** — critical helper (SetValue resets cursor to (0,0)):
```go
func repositionCursor(ta *textarea.Model, newText string, targetLine, targetCol int) {
    clamp := strings.Count(newText, "\n")
    for i := 0; i < min(targetLine, clamp); i++ {
        ta.CursorDown()
    }
    ta.SetCursorColumn(targetCol)
}
```

Always compute `targetLine, targetCol` **before** calling `SetValue`, then call `repositionCursor` immediately after.

**View()** — switch on mode:
```go
func (p RequestDetailsPane) View() string {
    // border + mode indicator (yellow=INSERT, magenta=VISUAL*, plain=NORMAL)
    // pending indicator: show p.engine.PendingBuf() if non-empty (e.g. "d" waiting)
    var body string
    switch p.engine.Mode {
    case vim.ModeVisual, vim.ModeVisualLine, vim.ModeVisualBlock:
        body = p.renderVisual()
    default:
        body = p.textarea.View()
    }
}
```

**`renderVisual()`**:
- Read `Value()`, split into lines
- Get cursor from `textarea.Line()` + `LineInfo().ColumnOffset`; anchor from `engine.AnchorLine/Col`
- For each character: apply `lipgloss.NewStyle().Reverse(true)` if in selection, `Underline(true)` on cursor
- `isInSelection(line, col, anchor, cursor, mode)` handles char/line/block geometry
- Scroll: approximate with `startRender = max(0, curLine - visibleLines/2)`

Selection geometry:
- **Char** (`ModeVisual`): contiguous range from anchor to cursor (handles reverse direction)
- **Line** (`ModeVisualLine`): all chars on lines between anchor and cursor
- **Block** (`ModeVisualBlock`): rectangle `[minLine..maxLine][minCol..maxCol]`; clamp to line length

### `internal/app/app.go`
```go
engine := vim.NewEngine()
configPath := filepath.Join(xdg.ConfigHome, "lazyhttp", "init.lua")
if err := vim.LoadConfig(configPath, &engine); err != nil {
    log.Printf("warn: vim config: %v", err)
}
rootModel := tui.NewRootModel(st, eb, engine)
```

### `internal/tui/model.go`
`NewRootModel` gains `engine vim.Engine`, forwards to `panes.NewRequestDetailsPane(st, engine)`. No other changes.

### `go.mod`
Add `github.com/yuin/gopher-lua`.

---

## Implementation Order

1. `go get github.com/yuin/gopher-lua` — verify clean build
2. `internal/vim/mode.go` — trivial, no deps
3. `internal/vim/action.go` — trivial, no deps
4. `internal/vim/textobject.go` + **unit tests** — most critical; test exhaustively before engine
5. `internal/vim/keymap.go` — data only
6. `internal/vim/parser.go` + unit tests (count, leader, 0-ambiguity, prefix matching)
7. `internal/vim/engine.go` + unit tests (ProcessKey sequences)
8. `internal/vim/config.go` + unit tests (use `L.DoString()` not file I/O in tests)
9. `internal/app/app.go` + `internal/tui/model.go` — thread engine
10. `internal/tui/panes/request_details.go` — executeAction + renderVisual (final, requires all above)

---

## Key Edge Cases

| Problem | Solution |
|---|---|
| `SetValue` resets cursor to (0,0) | Compute target before `SetValue`; call `repositionCursor` (CursorDown×N + SetCursorColumn) immediately after |
| `0` = count digit OR line-start | `0` → count only when `count > 0`; else appends to buf → matches `ActionMoveLineStart` |
| `ciw` on whitespace | `WordBounds` returns whitespace run when `text[cursor]` is space/tab |
| Count clamping (e.g. `1000j`) | `min(count, LineCount()-Line()-1)` before CursorDown loop |
| Visual block on short lines | Clamp `colIdx` to `len(runes)` in `isInSelection` for ModeVisualBlock |
| Space as leader | Default keymap never binds bare `" "`, so leader combos always wait for second key |
| Undo memory | Full snapshots capped at 100; request bodies are small (< 10 KB) |
| Visual scroll offset | `startRender = max(0, curLine - visibleLines/2)` — approximate centring |
| `repositionCursor` column trap | Always call `CursorDown` N times first, then `SetCursorColumn` — column is relative to current line |

---

## Verification

1. `go build ./...` — clean build
2. `go test ./internal/vim/...` — all unit tests pass
3. Manual smoke test:
   - Focus request details pane
   - `5j` — cursor moves down 5 lines
   - `dd` — current line deleted
   - `ciw` — word deleted, enters INSERT mode
   - `v` → `VISUAL` mode indicator; cursor movement highlights text
   - `V` → `V-LINE`; `ctrl+v` → `V-BLOCK`
   - `esc` → back to NORMAL; `esc` again → focus returns to requests pane
   - `u` — last mutation undone
4. Lua config smoke test:
   ```lua
   -- ~/.config/lazyhttp/init.lua
   vim.g.mapleader = ","
   vim.keymap.set("n", "<leader>s", "save")
   vim.keymap.set("n", "H", "move_line_start")
   ```
   Relaunch: `H` moves to line start; `,s` triggers save action
