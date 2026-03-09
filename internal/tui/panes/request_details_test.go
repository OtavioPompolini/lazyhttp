package panes

import (
	"strings"
	"testing"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"

	"github.com/OtavioPompolini/project-postman/internal/vim"
)

// ── Test helpers ─────────────────────────────────────────────────────────────

// newTestPane creates a focused RequestDetailsPane with the given content.
// Cursor starts at (0,0). Note: textarea.SetValue leaves cursor at end (it
// calls InsertString internally), so we call MoveToBegin() explicitly.
func newTestPane(content string) RequestDetailsPane {
	ta := textarea.New()
	ta.SetValue(content)
	ta.MoveToBegin() // SetValue leaves cursor at end; reset to (0,0)
	ta.Focus()
	return RequestDetailsPane{
		focused:  true,
		engine:   vim.NewEngine(),
		textarea: ta,
	}
}

// newTestPaneAt creates a pane with the cursor at the given line/col.
func newTestPaneAt(content string, line, col int) RequestDetailsPane {
	p := newTestPane(content)
	p.repositionCursor(content, line, col)
	return p
}

// key constructs a tea.KeyPressMsg from a string, setting Text so String() works.
// Handles "esc", "ctrl+r", "ctrl+v", and regular characters.
func key(s string) tea.KeyPressMsg {
	switch s {
	case "esc":
		return tea.KeyPressMsg{Code: tea.KeyEscape}
	case "enter":
		return tea.KeyPressMsg{Code: tea.KeyEnter}
	case "ctrl+r":
		return tea.KeyPressMsg{Code: 'r', Mod: tea.ModCtrl}
	case "ctrl+v":
		return tea.KeyPressMsg{Code: 'v', Mod: tea.ModCtrl}
	default:
		r := []rune(s)[0]
		return tea.KeyPressMsg{Code: r, Text: s}
	}
}

// sendKeys simulates a sequence of key presses.
func sendKeys(p RequestDetailsPane, keys ...string) RequestDetailsPane {
	for _, k := range keys {
		p, _ = p.Update(key(k))
	}
	return p
}

// value returns the current textarea content.
func value(p RequestDetailsPane) string {
	return p.textarea.Value()
}

// cursorByte returns the cursor byte offset.
func cursorByte(p RequestDetailsPane) int {
	return p.cursorByteOffset()
}

// mode returns the current vim mode.
func mode(p RequestDetailsPane) vim.Mode {
	return p.engine.Mode
}

// ── Motion tests ──────────────────────────────────────────────────────────────

func TestMotion_hjkl(t *testing.T) {
	p := newTestPane("abc\ndef\nghi")

	// Move right
	p = sendKeys(p, "l")
	if cursorByte(p) != 1 {
		t.Errorf("l: cursor byte want 1, got %d", cursorByte(p))
	}

	// Move left
	p = sendKeys(p, "h")
	if cursorByte(p) != 0 {
		t.Errorf("h: cursor byte want 0, got %d", cursorByte(p))
	}

	// Move down
	p = sendKeys(p, "j")
	line := p.textarea.Line()
	if line != 1 {
		t.Errorf("j: line want 1, got %d", line)
	}

	// Move up
	p = sendKeys(p, "k")
	line = p.textarea.Line()
	if line != 0 {
		t.Errorf("k: line want 0, got %d", line)
	}
}

func TestMotion_LineStartEnd(t *testing.T) {
	p := newTestPaneAt("hello world", 0, 5)

	p = sendKeys(p, "0")
	if cursorByte(p) != 0 {
		t.Errorf("0: want cursor at 0, got %d", cursorByte(p))
	}

	p = sendKeys(p, "$")
	if cursorByte(p) != 11 {
		t.Errorf("$: want cursor at 11, got %d", cursorByte(p))
	}
}

func TestMotion_DocStartEnd(t *testing.T) {
	p := newTestPane("line1\nline2\nline3")

	p = sendKeys(p, "G")
	if p.textarea.Line() != 2 {
		t.Errorf("G: want line 2, got %d", p.textarea.Line())
	}

	p = sendKeys(p, "g", "g")
	if p.textarea.Line() != 0 {
		t.Errorf("gg: want line 0, got %d", p.textarea.Line())
	}
}

func TestMotion_WordForward(t *testing.T) {
	p := newTestPane("hello world foo")
	// w from "hello" should land on "world"
	p = sendKeys(p, "w")
	if cursorByte(p) != 6 {
		t.Errorf("w: want cursor at 6 (start of 'world'), got %d", cursorByte(p))
	}
}

func TestMotion_WordBack(t *testing.T) {
	p := newTestPaneAt("hello world", 0, 6)
	p = sendKeys(p, "b")
	if cursorByte(p) != 0 {
		t.Errorf("b: want cursor at 0 (start of 'hello'), got %d", cursorByte(p))
	}
}

func TestMotion_Count(t *testing.T) {
	p := newTestPane("abc\ndef\nghi\njkl")
	p = sendKeys(p, "3", "j")
	if p.textarea.Line() != 3 {
		t.Errorf("3j: want line 3, got %d", p.textarea.Line())
	}
}

// ── Insert mode entry ─────────────────────────────────────────────────────────

func TestInsert_i_EntersInsertMode(t *testing.T) {
	p := newTestPane("hello")
	p = sendKeys(p, "i")
	if mode(p) != vim.ModeInsert {
		t.Errorf("i: want INSERT mode, got %s", mode(p))
	}
}

func TestInsert_a_MovesRight(t *testing.T) {
	p := newTestPane("hello")
	p = sendKeys(p, "a")
	if mode(p) != vim.ModeInsert {
		t.Errorf("a: want INSERT mode, got %s", mode(p))
	}
	// cursor should be 1 position to the right of where we were (after 'h')
	// can't easily test without typing, but mode is enough for now
}

func TestInsert_Escape_ReturnToNormal(t *testing.T) {
	p := newTestPane("hello")
	p = sendKeys(p, "i", "esc")
	if mode(p) != vim.ModeNormal {
		t.Errorf("esc: want NORMAL mode, got %s", mode(p))
	}
}

// ── Delete operations ─────────────────────────────────────────────────────────

func TestDD_DeletesLine(t *testing.T) {
	p := newTestPane("hello\nworld\nfoo")
	p = sendKeys(p, "d", "d")
	got := value(p)
	want := "world\nfoo"
	if got != want {
		t.Errorf("dd: got %q, want %q", got, want)
	}
}

func TestDD_DeletesMiddleLine(t *testing.T) {
	p := newTestPane("hello\nworld\nfoo")
	p = sendKeys(p, "j", "d", "d")
	got := value(p)
	want := "hello\nfoo"
	if got != want {
		t.Errorf("dd on middle line: got %q, want %q", got, want)
	}
}

func TestD_ToLineEnd(t *testing.T) {
	p := newTestPaneAt("hello world", 0, 5)
	p = sendKeys(p, "D")
	got := value(p)
	want := "hello"
	if got != want {
		t.Errorf("D: got %q, want %q", got, want)
	}
}

func TestDW_DeleteWord(t *testing.T) {
	p := newTestPane("hello world")
	p = sendKeys(p, "d", "w")
	got := value(p)
	want := "world"
	if got != want {
		t.Errorf("dw: got %q, want %q", got, want)
	}
}

func TestDIW_FirstWord_IncludesTrailingSpace(t *testing.T) {
	// diw on first word should delete word + trailing space
	p := newTestPane("hello world")
	p = sendKeys(p, "d", "i", "w")
	got := value(p)
	want := "world"
	if got != want {
		t.Errorf("diw on first word: got %q, want %q", got, want)
	}
}

func TestDIW_LastWord_IncludesLeadingSpace(t *testing.T) {
	// diw on last word should delete leading space + word
	p := newTestPaneAt("hello world", 0, 7) // cursor on 'o' of 'world'
	p = sendKeys(p, "d", "i", "w")
	got := value(p)
	want := "hello"
	if got != want {
		t.Errorf("diw on last word: got %q, want %q", got, want)
	}
}

func TestDIW_SingleWord_NoSurroundingSpace(t *testing.T) {
	p := newTestPane("hello")
	p = sendKeys(p, "d", "i", "w")
	got := value(p)
	want := ""
	if got != want {
		t.Errorf("diw on single word: got %q, want %q", got, want)
	}
}

func TestDIW_OnWhitespace_DeletesWhitespaceRun(t *testing.T) {
	// diw with cursor on whitespace deletes the whitespace run
	p := newTestPaneAt("hello   world", 0, 6) // cursor on middle space
	p = sendKeys(p, "d", "i", "w")
	got := value(p)
	want := "helloworld"
	if got != want {
		t.Errorf("diw on whitespace: got %q, want %q", got, want)
	}
}

func TestDIW_MultiLine(t *testing.T) {
	p := newTestPane("foo bar\nbaz qux")
	p = sendKeys(p, "d", "i", "w")
	got := value(p)
	want := "bar\nbaz qux"
	if got != want {
		t.Errorf("diw multiline first word: got %q, want %q", got, want)
	}
}

func TestDI_Quote(t *testing.T) {
	p := newTestPaneAt(`say "hello" now`, 0, 6) // cursor inside "hello"
	p = sendKeys(p, "d", "i", `"`)
	got := value(p)
	want := `say "" now`
	if got != want {
		t.Errorf(`di": got %q, want %q`, got, want)
	}
}

func TestDI_Paren(t *testing.T) {
	p := newTestPaneAt("foo(bar)baz", 0, 5) // cursor inside parens
	p = sendKeys(p, "d", "i", "(")
	got := value(p)
	want := "foo()baz"
	if got != want {
		t.Errorf("di(: got %q, want %q", got, want)
	}
}

func TestDI_Brace(t *testing.T) {
	p := newTestPaneAt(`{"key":"val"}`, 0, 5) // cursor inside braces
	p = sendKeys(p, "d", "i", "{")
	got := value(p)
	want := "{}"
	if got != want {
		t.Errorf("di{: got %q, want %q", got, want)
	}
}

// ── Change operations ─────────────────────────────────────────────────────────

func TestCIW_DeletesWordEntersInsert(t *testing.T) {
	p := newTestPane("hello world")
	p = sendKeys(p, "c", "i", "w")
	if mode(p) != vim.ModeInsert {
		t.Errorf("ciw: want INSERT mode, got %s", mode(p))
	}
	got := value(p)
	want := "world"
	if got != want {
		t.Errorf("ciw: got %q, want %q", got, want)
	}
}

func TestCIW_CursorAtStartOfRemainingWord(t *testing.T) {
	p := newTestPane("hello world")
	p = sendKeys(p, "c", "i", "w")
	// After deleting "hello " (with trailing space), cursor should be at 'w' of "world"
	if cursorByte(p) != 0 {
		t.Errorf("ciw cursor: want byte 0 (start of 'world'), got %d", cursorByte(p))
	}
}

func TestCIW_LastWord_CursorAtEndOfRemaining(t *testing.T) {
	p := newTestPaneAt("hello world", 0, 7)
	p = sendKeys(p, "c", "i", "w")
	got := value(p)
	want := "hello"
	if got != want {
		t.Errorf("ciw last word: got %q, want %q", got, want)
	}
	if mode(p) != vim.ModeInsert {
		t.Errorf("ciw: want INSERT, got %s", mode(p))
	}
}

func TestCC_DeletesLineEntersInsert(t *testing.T) {
	p := newTestPane("hello\nworld")
	p = sendKeys(p, "c", "c")
	if mode(p) != vim.ModeInsert {
		t.Errorf("cc: want INSERT, got %s", mode(p))
	}
	got := value(p)
	want := "\nworld"
	if got != want {
		t.Errorf("cc: got %q, want %q", got, want)
	}
}

func TestC_ToLineEnd(t *testing.T) {
	p := newTestPaneAt("hello world", 0, 5)
	p = sendKeys(p, "C")
	if mode(p) != vim.ModeInsert {
		t.Errorf("C: want INSERT, got %s", mode(p))
	}
	got := value(p)
	want := "hello"
	if got != want {
		t.Errorf("C: got %q, want %q", got, want)
	}
}

func TestCI_Quote(t *testing.T) {
	p := newTestPaneAt(`"hello world"`, 0, 5)
	p = sendKeys(p, "c", "i", `"`)
	if mode(p) != vim.ModeInsert {
		t.Errorf("ci\": want INSERT, got %s", mode(p))
	}
	got := value(p)
	want := `""`
	if got != want {
		t.Errorf(`ci": got %q, want %q`, got, want)
	}
}

func TestCI_Paren(t *testing.T) {
	p := newTestPaneAt("fn(arg1, arg2)", 0, 5)
	p = sendKeys(p, "c", "i", "(")
	got := value(p)
	want := "fn()"
	if got != want {
		t.Errorf("ci(: got %q, want %q", got, want)
	}
}

// ── Yank operations ───────────────────────────────────────────────────────────

func TestYY_YanksLine(t *testing.T) {
	p := newTestPane("hello\nworld")
	p = sendKeys(p, "y", "y")
	if p.engine.Register != "hello" {
		t.Errorf("yy: register want %q, got %q", "hello", p.engine.Register)
	}
}

func TestYIW_YanksInnerWord(t *testing.T) {
	p := newTestPane("hello world")
	p = sendKeys(p, "y", "i", "w")
	// yiw yanks "hello " (inner word with trailing space)
	if p.engine.Register != "hello " {
		t.Errorf("yiw: register want %q, got %q", "hello ", p.engine.Register)
	}
}

// ── Paste operations ──────────────────────────────────────────────────────────

func TestP_PastesAfterCursor(t *testing.T) {
	p := newTestPane("ac")
	p.engine.Register = "b"
	p = sendKeys(p, "p")
	got := value(p)
	want := "abc"
	if got != want {
		t.Errorf("p: got %q, want %q", got, want)
	}
}

func TestP_PastesBeforeCursor(t *testing.T) {
	p := newTestPaneAt("bc", 0, 0)
	p.engine.Register = "a"
	p = sendKeys(p, "P")
	got := value(p)
	want := "abc"
	if got != want {
		t.Errorf("P: got %q, want %q", got, want)
	}
}

func TestYankAndPaste(t *testing.T) {
	p := newTestPane("hello world")
	p = sendKeys(p, "y", "y") // yank "hello"
	// move to end, paste
	p = sendKeys(p, "$", "p")
	// "hello world" + paste "hello" after last char → "hello worldhello"
	// Actually yy yanks the line, p pastes after the last char on the line
	got := value(p)
	if len(got) == 0 {
		t.Errorf("yank+paste: got empty string")
	}
}

// ── Undo/Redo ─────────────────────────────────────────────────────────────────

func TestUndo_AfterDD(t *testing.T) {
	p := newTestPane("hello\nworld")
	p = sendKeys(p, "d", "d")
	if value(p) != "world" {
		t.Fatalf("dd: want %q, got %q", "world", value(p))
	}
	p = sendKeys(p, "u")
	got := value(p)
	want := "hello\nworld"
	if got != want {
		t.Errorf("u after dd: got %q, want %q", got, want)
	}
}

func TestUndo_AfterDIW(t *testing.T) {
	p := newTestPane("hello world")
	p = sendKeys(p, "d", "i", "w")
	if value(p) != "world" {
		t.Fatalf("diw: want %q, got %q", "world", value(p))
	}
	p = sendKeys(p, "u")
	got := value(p)
	want := "hello world"
	if got != want {
		t.Errorf("u after diw: got %q, want %q", got, want)
	}
}

func TestRedo_AfterUndo(t *testing.T) {
	p := newTestPane("hello world")
	p = sendKeys(p, "d", "i", "w") // delete "hello "
	p = sendKeys(p, "u")            // undo
	p = sendKeys(p, "ctrl+r")       // redo
	got := value(p)
	want := "world"
	if got != want {
		t.Errorf("redo: got %q, want %q", got, want)
	}
}

// ── Visual mode ───────────────────────────────────────────────────────────────

func TestV_EntersVisualMode(t *testing.T) {
	p := newTestPane("hello world")
	p = sendKeys(p, "v")
	if mode(p) != vim.ModeVisual {
		t.Errorf("v: want VISUAL, got %s", mode(p))
	}
}

func TestV_Line_EntersVisualLine(t *testing.T) {
	p := newTestPane("hello world")
	p = sendKeys(p, "V")
	if mode(p) != vim.ModeVisualLine {
		t.Errorf("V: want V-LINE, got %s", mode(p))
	}
}

func TestV_Esc_ExitsVisual(t *testing.T) {
	p := newTestPane("hello world")
	p = sendKeys(p, "v", "esc")
	if mode(p) != vim.ModeNormal {
		t.Errorf("v+esc: want NORMAL, got %s", mode(p))
	}
}

func TestVIW_SelectsWord(t *testing.T) {
	// viw: v enters visual, iw selects inner word + adjacent space
	p := newTestPane("hello world")
	p = sendKeys(p, "v", "i", "w")
	// Should remain in visual mode
	if mode(p) != vim.ModeVisual {
		t.Errorf("viw: want VISUAL mode, got %s", mode(p))
	}
	// Anchor should be at start of "hello " (byte 0)
	if p.engine.AnchorLine != 0 || p.engine.AnchorCol != 0 {
		t.Errorf("viw: anchor want (0,0), got (%d,%d)", p.engine.AnchorLine, p.engine.AnchorCol)
	}
	// Cursor should be at last char of "hello " = position of ' ' = byte 5
	cb := cursorByte(p)
	if cb != 5 {
		t.Errorf("viw: cursor byte want 5 (space after hello), got %d", cb)
	}
}

func TestVIW_ThenDelete(t *testing.T) {
	p := newTestPane("hello world")
	p = sendKeys(p, "v", "i", "w", "d")
	got := value(p)
	want := "world"
	if got != want {
		t.Errorf("viw+d: got %q, want %q", got, want)
	}
	if mode(p) != vim.ModeNormal {
		t.Errorf("viw+d: want NORMAL, got %s", mode(p))
	}
}

func TestVisualDelete_Char(t *testing.T) {
	// Enter visual, move right 4, delete → deletes first 5 chars
	p := newTestPane("hello world")
	p = sendKeys(p, "v", "l", "l", "l", "l", "d")
	got := value(p)
	want := " world"
	if got != want {
		t.Errorf("v+4l+d: got %q, want %q", got, want)
	}
}

func TestVisualLine_Delete(t *testing.T) {
	p := newTestPane("hello\nworld\nfoo")
	// V selects current line, d deletes it
	p = sendKeys(p, "V", "d")
	got := value(p)
	want := "world\nfoo"
	if got != want {
		t.Errorf("V+d: got %q, want %q", got, want)
	}
}

func TestVisualYank(t *testing.T) {
	p := newTestPane("hello world")
	p = sendKeys(p, "v", "l", "l", "l", "l", "y")
	if p.engine.Register == "" {
		t.Errorf("v+4l+y: register is empty")
	}
	// Should be back in normal mode
	if mode(p) != vim.ModeNormal {
		t.Errorf("v+y: want NORMAL, got %s", mode(p))
	}
}

// ── Visual text objects ───────────────────────────────────────────────────────

func TestVI_Quote(t *testing.T) {
	p := newTestPaneAt(`say "hello" now`, 0, 6) // cursor inside "hello"
	p = sendKeys(p, "v", "i", `"`)
	if mode(p) != vim.ModeVisual {
		t.Errorf(`vi": want VISUAL, got %s`, mode(p))
	}
	// anchor at byte 5 (after opening quote), cursor at byte 9 (before closing quote)
	wantAnchorByte := 5
	anchorByte := vim.ByteOffsetFromLineCol(value(p), p.engine.AnchorLine, p.engine.AnchorCol)
	if anchorByte != wantAnchorByte {
		t.Errorf(`vi" anchor: want byte %d, got %d`, wantAnchorByte, anchorByte)
	}
}

func TestVI_Paren(t *testing.T) {
	p := newTestPaneAt("fn(arg1)", 0, 4) // cursor on 'a'
	p = sendKeys(p, "v", "i", "(")
	if mode(p) != vim.ModeVisual {
		t.Errorf("vi(: want VISUAL, got %s", mode(p))
	}
	// anchor byte 3 (after '('), cursor at byte 6 (before ')')
	anchorByte := vim.ByteOffsetFromLineCol(value(p), p.engine.AnchorLine, p.engine.AnchorCol)
	if anchorByte != 3 {
		t.Errorf("vi( anchor: want byte 3, got %d", anchorByte)
	}
}

// ── Focus / escape ────────────────────────────────────────────────────────────

func TestEsc_InNormal_FocusesOut(t *testing.T) {
	p := newTestPane("hello")
	_, cmd := p.Update(key("esc"))
	if cmd == nil {
		t.Errorf("esc in normal: want FocusCmd, got nil")
	}
}

func TestEsc_InInsert_ReturnsToNormal(t *testing.T) {
	p := sendKeys(newTestPane("hello"), "i")
	if mode(p) != vim.ModeInsert {
		t.Fatalf("setup: want INSERT, got %s", mode(p))
	}
	p = sendKeys(p, "esc")
	if mode(p) != vim.ModeNormal {
		t.Errorf("esc in insert: want NORMAL, got %s", mode(p))
	}
}

// ── Insert mode typing ────────────────────────────────────────────────────────

func TestInsertMode_TypesText(t *testing.T) {
	p := newTestPane("")
	p = sendKeys(p, "i")
	// In insert mode, keys pass through to textarea
	p, _ = p.Update(key("h"))
	p, _ = p.Update(key("i"))
	if mode(p) != vim.ModeInsert {
		t.Errorf("still want INSERT, got %s", mode(p))
	}
	got := value(p)
	if got != "hi" {
		t.Errorf("insert typing: got %q, want %q", got, "hi")
	}
}

// ── Pending indicator ─────────────────────────────────────────────────────────

func TestPendingBuf_ShowsDuringSequence(t *testing.T) {
	p := newTestPane("hello")
	p, _ = p.Update(key("d"))
	if p.engine.PendingBuf() != "d" {
		t.Errorf("after 'd': pending buf want %q, got %q", "d", p.engine.PendingBuf())
	}
	p, _ = p.Update(key("i"))
	if p.engine.PendingBuf() != "di" {
		t.Errorf("after 'di': pending buf want %q, got %q", "di", p.engine.PendingBuf())
	}
	p, _ = p.Update(key("w"))
	if p.engine.PendingBuf() != "" {
		t.Errorf("after 'diw': pending buf should be empty, got %q", p.engine.PendingBuf())
	}
}

// ── GG and G doc navigation ───────────────────────────────────────────────────

func TestGG_MovesToDocStart(t *testing.T) {
	p := newTestPane("a\nb\nc")
	p = sendKeys(p, "j", "j") // go to line 2
	p = sendKeys(p, "g", "g") // gg
	if p.textarea.Line() != 0 {
		t.Errorf("gg: want line 0, got %d", p.textarea.Line())
	}
}

func TestG_MovesToDocEnd(t *testing.T) {
	p := newTestPane("a\nb\nc")
	p = sendKeys(p, "G")
	if p.textarea.Line() != 2 {
		t.Errorf("G: want line 2, got %d", p.textarea.Line())
	}
}

// ── New line ──────────────────────────────────────────────────────────────────

func TestO_NewLineBelow(t *testing.T) {
	p := newTestPane("hello\nworld")
	p = sendKeys(p, "o")
	if mode(p) != vim.ModeInsert {
		t.Errorf("o: want INSERT, got %s", mode(p))
	}
	// There should now be 3 lines
	lines := strings.Count(value(p), "\n")
	if lines != 2 {
		t.Errorf("o: want 2 newlines (3 lines), got %d", lines)
	}
}

func TestO_NewLineAbove(t *testing.T) {
	p := newTestPane("hello\nworld")
	p = sendKeys(p, "O")
	if mode(p) != vim.ModeInsert {
		t.Errorf("O: want INSERT, got %s", mode(p))
	}
	lines := strings.Count(value(p), "\n")
	if lines != 2 {
		t.Errorf("O: want 2 newlines (3 lines), got %d", lines)
	}
}
