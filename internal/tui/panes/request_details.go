package panes

import (
	"strings"
	"unicode/utf8"

	"charm.land/bubbles/v2/textarea"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/OtavioPompolini/project-postman/internal/state"
	"github.com/OtavioPompolini/project-postman/internal/tui/msgs"
	"github.com/OtavioPompolini/project-postman/internal/types"
	"github.com/OtavioPompolini/project-postman/internal/vim"
)

type RequestDetailsPane struct {
	focused        bool
	width, height  int
	engine         vim.Engine
	textarea       textarea.Model
	currentRequest *types.Request
	requestSystem  *state.RequestManager
}

func NewRequestDetailsPane(st *state.State, engine vim.Engine) RequestDetailsPane {
	ta := textarea.New()
	ta.Placeholder = "Request body..."
	ta.ShowLineNumbers = false
	ta.SetWidth(40)
	ta.SetHeight(20)

	return RequestDetailsPane{
		textarea:      ta,
		requestSystem: st.RequestManager,
		engine:        engine,
	}
}

func (p RequestDetailsPane) Init() tea.Cmd { return nil }

func (p RequestDetailsPane) Update(msg tea.Msg) (RequestDetailsPane, tea.Cmd) {
	switch m := msg.(type) {
	case msgs.RequestChangedMsg:
		if !p.focused {
			p.currentRequest = nil
			if len(m.Event.Requests) > 0 && m.Event.Cursor >= 0 && m.Event.Cursor < len(m.Event.Requests) {
				p.currentRequest = m.Event.Requests[m.Event.Cursor]
			}
			body := ""
			if p.currentRequest != nil {
				body = p.currentRequest.Body
			}
			p.textarea.SetValue(body)
			p.textarea.MoveToBegin() // SetValue leaves cursor at end; reset to (0,0)
		}

	case msgs.FocusLostMsg:
		p.focused = false
		p.engine.EnterMode(vim.ModeNormal, 0, 0)
		p.textarea.Blur()
		if p.currentRequest != nil {
			p.requestSystem.Update(&types.Request{
				Id:   p.currentRequest.Id,
				Body: p.textarea.Value(),
			})
		}

	case tea.KeyPressMsg:
		if !p.focused {
			break
		}
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
	}
	return p, nil
}

// executeAction performs the semantic vim action on the textarea.
func (p RequestDetailsPane) executeAction(action vim.Action, count int) (RequestDetailsPane, tea.Cmd) {
	switch action {

	// ── Mode entry ──────────────────────────────────────────────────────────

	case vim.ActionInsertBefore:
		p.engine.EnterMode(vim.ModeInsert, 0, 0)
		return p, textarea.Blink

	case vim.ActionInsertAfter:
		p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: tea.KeyRight})
		p.engine.EnterMode(vim.ModeInsert, 0, 0)
		return p, textarea.Blink

	case vim.ActionInsertLineStart:
		p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: 'a', Mod: tea.ModCtrl})
		p.engine.EnterMode(vim.ModeInsert, 0, 0)
		return p, textarea.Blink

	case vim.ActionInsertLineEnd:
		p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: 'e', Mod: tea.ModCtrl})
		p.engine.EnterMode(vim.ModeInsert, 0, 0)
		return p, textarea.Blink

	case vim.ActionNewLineBelow:
		p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: 'e', Mod: tea.ModCtrl})
		p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
		p.engine.EnterMode(vim.ModeInsert, 0, 0)
		return p, textarea.Blink

	case vim.ActionNewLineAbove:
		p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: 'a', Mod: tea.ModCtrl})
		p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
		p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: tea.KeyUp})
		p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: 'e', Mod: tea.ModCtrl})
		p.engine.EnterMode(vim.ModeInsert, 0, 0)
		return p, textarea.Blink

	// ── Pure motions ─────────────────────────────────────────────────────────

	case vim.ActionMoveLeft:
		for i := 0; i < count; i++ {
			p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
		}

	case vim.ActionMoveRight:
		for i := 0; i < count; i++ {
			p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: tea.KeyRight})
		}

	case vim.ActionMoveUp:
		for i := 0; i < count; i++ {
			p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: tea.KeyUp})
		}

	case vim.ActionMoveDown:
		for i := 0; i < count; i++ {
			p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: tea.KeyDown})
		}

	case vim.ActionMoveLineStart, vim.ActionMoveLineFirst:
		p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: 'a', Mod: tea.ModCtrl})

	case vim.ActionMoveLineEnd:
		p.textarea, _ = p.textarea.Update(tea.KeyPressMsg{Code: 'e', Mod: tea.ModCtrl})

	case vim.ActionMoveDocStart:
		p.textarea.MoveToBegin()

	case vim.ActionMoveDocEnd:
		p.textarea.MoveToEnd()

	case vim.ActionMoveWordForward, vim.ActionMoveWORDForward:
		text := p.textarea.Value()
		curByte := p.cursorByteOffset()
		for i := 0; i < count; i++ {
			if action == vim.ActionMoveWordForward {
				curByte = vim.WordForwardByte(text, curByte)
			} else {
				curByte = vim.WORDForwardByte(text, curByte)
			}
		}
		line, col := vim.LineColFromByteOffset(text, curByte)
		p.repositionCursor(text, line, col)

	case vim.ActionMoveWordEnd, vim.ActionMoveWORDEnd:
		text := p.textarea.Value()
		curByte := p.cursorByteOffset()
		for i := 0; i < count; i++ {
			if action == vim.ActionMoveWordEnd {
				curByte = vim.WordEndByte(text, curByte)
			} else {
				curByte = vim.WORDEndByte(text, curByte)
			}
		}
		line, col := vim.LineColFromByteOffset(text, curByte)
		p.repositionCursor(text, line, col)

	case vim.ActionMoveWordBack, vim.ActionMoveWORDBack:
		text := p.textarea.Value()
		curByte := p.cursorByteOffset()
		for i := 0; i < count; i++ {
			if action == vim.ActionMoveWordBack {
				curByte = vim.WordBackByte(text, curByte)
			} else {
				curByte = vim.WORDBackByte(text, curByte)
			}
		}
		line, col := vim.LineColFromByteOffset(text, curByte)
		p.repositionCursor(text, line, col)

	// ── Delete ops ───────────────────────────────────────────────────────────

	case vim.ActionDeleteLine:
		text := p.textarea.Value()
		curByte := p.cursorByteOffset()
		for i := 0; i < count; i++ {
			s, e := vim.LineBounds(text, curByte, false)
			p.engine.Register = text[s:e]
			text = text[:s] + text[e:]
		}
		p.engine.PushUndo(p.textarea.Value())
		targetLine, targetCol := vim.LineColFromByteOffset(text, curByte)
		p.textarea.SetValue(text)
		p.repositionCursor(text, targetLine, targetCol)

	case vim.ActionDeleteWord:
		text := p.textarea.Value()
		curByte := p.cursorByteOffset()
		p.engine.PushUndo(text)
		end := vim.WordForwardByte(text, curByte)
		p.engine.Register = text[curByte:end]
		newText := text[:curByte] + text[end:]
		targetLine, targetCol := vim.LineColFromByteOffset(newText, curByte)
		p.textarea.SetValue(newText)
		p.repositionCursor(newText, targetLine, targetCol)

	case vim.ActionDeleteToLineEnd:
		text := p.textarea.Value()
		curByte := p.cursorByteOffset()
		p.engine.PushUndo(text)
		_, e := vim.LineBounds(text, curByte, true)
		p.engine.Register = text[curByte:e]
		newText := text[:curByte] + text[e:]
		targetLine, targetCol := vim.LineColFromByteOffset(newText, curByte)
		p.textarea.SetValue(newText)
		p.repositionCursor(newText, targetLine, targetCol)

	case vim.ActionDeleteInnerWord:
		return p.deleteTextObject(func(text string, curByte int) (int, int, bool) {
			s, e := vim.WordBounds(text, curByte, true)
			return s, e, true
		})

	case vim.ActionDeleteInnerWORD:
		return p.deleteTextObject(func(text string, curByte int) (int, int, bool) {
			s, e := vim.WORDBounds(text, curByte, true)
			return s, e, true
		})

	case vim.ActionDeleteInnerQuoteS:
		return p.deleteTextObject(func(text string, curByte int) (int, int, bool) {
			s, e, ok := vim.QuoteBounds(text, curByte, '\'', true)
			return s, e, ok
		})

	case vim.ActionDeleteInnerQuoteD:
		return p.deleteTextObject(func(text string, curByte int) (int, int, bool) {
			s, e, ok := vim.QuoteBounds(text, curByte, '"', true)
			return s, e, ok
		})

	case vim.ActionDeleteInnerQuoteB:
		return p.deleteTextObject(func(text string, curByte int) (int, int, bool) {
			s, e, ok := vim.QuoteBounds(text, curByte, '`', true)
			return s, e, ok
		})

	case vim.ActionDeleteInnerParen:
		return p.deleteTextObject(func(text string, curByte int) (int, int, bool) {
			s, e, ok := vim.BracketBounds(text, curByte, '(', true)
			return s, e, ok
		})

	case vim.ActionDeleteInnerBrace:
		return p.deleteTextObject(func(text string, curByte int) (int, int, bool) {
			s, e, ok := vim.BracketBounds(text, curByte, '{', true)
			return s, e, ok
		})

	case vim.ActionDeleteInnerBracket:
		return p.deleteTextObject(func(text string, curByte int) (int, int, bool) {
			s, e, ok := vim.BracketBounds(text, curByte, '[', true)
			return s, e, ok
		})

	// ── Change ops ───────────────────────────────────────────────────────────

	case vim.ActionChangeLine:
		text := p.textarea.Value()
		curByte := p.cursorByteOffset()
		p.engine.PushUndo(text)
		s, e := vim.LineBounds(text, curByte, true)
		p.engine.Register = text[s:e]
		newText := text[:s] + text[e:]
		targetLine, targetCol := vim.LineColFromByteOffset(newText, s)
		p.textarea.SetValue(newText)
		p.repositionCursor(newText, targetLine, targetCol)
		p.engine.EnterMode(vim.ModeInsert, 0, 0)
		return p, textarea.Blink

	case vim.ActionChangeWord:
		text := p.textarea.Value()
		curByte := p.cursorByteOffset()
		p.engine.PushUndo(text)
		end := vim.WordForwardByte(text, curByte)
		p.engine.Register = text[curByte:end]
		newText := text[:curByte] + text[end:]
		targetLine, targetCol := vim.LineColFromByteOffset(newText, curByte)
		p.textarea.SetValue(newText)
		p.repositionCursor(newText, targetLine, targetCol)
		p.engine.EnterMode(vim.ModeInsert, 0, 0)
		return p, textarea.Blink

	case vim.ActionChangeToLineEnd:
		text := p.textarea.Value()
		curByte := p.cursorByteOffset()
		p.engine.PushUndo(text)
		_, e := vim.LineBounds(text, curByte, true)
		p.engine.Register = text[curByte:e]
		newText := text[:curByte] + text[e:]
		targetLine, targetCol := vim.LineColFromByteOffset(newText, curByte)
		p.textarea.SetValue(newText)
		p.repositionCursor(newText, targetLine, targetCol)
		p.engine.EnterMode(vim.ModeInsert, 0, 0)
		return p, textarea.Blink

	case vim.ActionChangeInnerWord:
		return p.changeTextObject(func(text string, curByte int) (int, int, bool) {
			s, e := vim.WordBounds(text, curByte, true)
			return s, e, true
		})

	case vim.ActionChangeInnerWORD:
		return p.changeTextObject(func(text string, curByte int) (int, int, bool) {
			s, e := vim.WORDBounds(text, curByte, true)
			return s, e, true
		})

	case vim.ActionChangeInnerQuoteS:
		return p.changeTextObject(func(text string, curByte int) (int, int, bool) {
			s, e, ok := vim.QuoteBounds(text, curByte, '\'', true)
			return s, e, ok
		})

	case vim.ActionChangeInnerQuoteD:
		return p.changeTextObject(func(text string, curByte int) (int, int, bool) {
			s, e, ok := vim.QuoteBounds(text, curByte, '"', true)
			return s, e, ok
		})

	case vim.ActionChangeInnerQuoteB:
		return p.changeTextObject(func(text string, curByte int) (int, int, bool) {
			s, e, ok := vim.QuoteBounds(text, curByte, '`', true)
			return s, e, ok
		})

	case vim.ActionChangeInnerParen:
		return p.changeTextObject(func(text string, curByte int) (int, int, bool) {
			s, e, ok := vim.BracketBounds(text, curByte, '(', true)
			return s, e, ok
		})

	case vim.ActionChangeInnerBrace:
		return p.changeTextObject(func(text string, curByte int) (int, int, bool) {
			s, e, ok := vim.BracketBounds(text, curByte, '{', true)
			return s, e, ok
		})

	case vim.ActionChangeInnerBracket:
		return p.changeTextObject(func(text string, curByte int) (int, int, bool) {
			s, e, ok := vim.BracketBounds(text, curByte, '[', true)
			return s, e, ok
		})

	// ── Yank ops ─────────────────────────────────────────────────────────────

	case vim.ActionYankLine:
		text := p.textarea.Value()
		curByte := p.cursorByteOffset()
		s, e := vim.LineBounds(text, curByte, true)
		p.engine.Register = text[s:e]

	case vim.ActionYankWord:
		text := p.textarea.Value()
		curByte := p.cursorByteOffset()
		end := vim.WordForwardByte(text, curByte)
		p.engine.Register = text[curByte:end]

	case vim.ActionYankInnerWord:
		text := p.textarea.Value()
		curByte := p.cursorByteOffset()
		s, e := vim.WordBounds(text, curByte, true)
		p.engine.Register = text[s:e]

	// ── Clipboard ────────────────────────────────────────────────────────────

	case vim.ActionPasteAfter:
		if p.engine.Register == "" {
			break
		}
		text := p.textarea.Value()
		curByte := p.cursorByteOffset()
		p.engine.PushUndo(text)
		// paste after cursor
		insertAt := curByte
		if insertAt < len(text) {
			_, size := utf8.DecodeRuneInString(text[insertAt:])
			insertAt += size
		}
		newText := text[:insertAt] + p.engine.Register + text[insertAt:]
		targetLine, targetCol := vim.LineColFromByteOffset(newText, insertAt)
		p.textarea.SetValue(newText)
		p.repositionCursor(newText, targetLine, targetCol)

	case vim.ActionPasteBefore:
		if p.engine.Register == "" {
			break
		}
		text := p.textarea.Value()
		curByte := p.cursorByteOffset()
		p.engine.PushUndo(text)
		newText := text[:curByte] + p.engine.Register + text[curByte:]
		targetLine, targetCol := vim.LineColFromByteOffset(newText, curByte)
		p.textarea.SetValue(newText)
		p.repositionCursor(newText, targetLine, targetCol)

	// ── Undo/Redo ────────────────────────────────────────────────────────────

	case vim.ActionUndo:
		if prev, ok := p.engine.PopUndo(); ok {
			p.engine.PushRedo(p.textarea.Value())
			p.textarea.SetValue(prev)
		}

	case vim.ActionRedo:
		if next, ok := p.engine.PopRedo(); ok {
			p.engine.PushUndo(p.textarea.Value())
			p.textarea.SetValue(next)
		}

	// ── Visual mode ──────────────────────────────────────────────────────────

	case vim.ActionVisualChar:
		curLine := p.textarea.Line()
		curCol := p.textarea.LineInfo().ColumnOffset
		p.engine.EnterMode(vim.ModeVisual, curLine, curCol)

	case vim.ActionVisualLine:
		curLine := p.textarea.Line()
		curCol := p.textarea.LineInfo().ColumnOffset
		p.engine.EnterMode(vim.ModeVisualLine, curLine, curCol)

	case vim.ActionVisualBlock:
		curLine := p.textarea.Line()
		curCol := p.textarea.LineInfo().ColumnOffset
		p.engine.EnterMode(vim.ModeVisualBlock, curLine, curCol)

	case vim.ActionVisualDelete:
		text := p.textarea.Value()
		s, e := p.visualByteRange(text)
		if s > e {
			s, e = e, s
		}
		p.engine.PushUndo(text)
		p.engine.Register = text[s:e]
		newText := text[:s] + text[e:]
		targetLine, targetCol := vim.LineColFromByteOffset(newText, s)
		p.textarea.SetValue(newText)
		p.repositionCursor(newText, targetLine, targetCol)
		p.engine.EnterMode(vim.ModeNormal, 0, 0)

	case vim.ActionVisualYank:
		text := p.textarea.Value()
		s, e := p.visualByteRange(text)
		if s > e {
			s, e = e, s
		}
		p.engine.Register = text[s:e]
		p.engine.EnterMode(vim.ModeNormal, 0, 0)

	case vim.ActionVisualChange:
		text := p.textarea.Value()
		s, e := p.visualByteRange(text)
		if s > e {
			s, e = e, s
		}
		p.engine.PushUndo(text)
		p.engine.Register = text[s:e]
		newText := text[:s] + text[e:]
		targetLine, targetCol := vim.LineColFromByteOffset(newText, s)
		p.textarea.SetValue(newText)
		p.repositionCursor(newText, targetLine, targetCol)
		p.engine.EnterMode(vim.ModeInsert, 0, 0)
		return p, textarea.Blink

	// ── Visual text object selections ────────────────────────────────────────

	case vim.ActionVisualSelectInnerWord:
		return p.visualSelectTextObject(func(text string, curByte int) (int, int, bool) {
			s, e := vim.WordBounds(text, curByte, true)
			return s, e, true
		})

	case vim.ActionVisualSelectInnerWORD:
		return p.visualSelectTextObject(func(text string, curByte int) (int, int, bool) {
			s, e := vim.WORDBounds(text, curByte, true)
			return s, e, true
		})

	case vim.ActionVisualSelectInnerQuoteS:
		return p.visualSelectTextObject(func(text string, curByte int) (int, int, bool) {
			s, e, ok := vim.QuoteBounds(text, curByte, '\'', true)
			return s, e, ok
		})

	case vim.ActionVisualSelectInnerQuoteD:
		return p.visualSelectTextObject(func(text string, curByte int) (int, int, bool) {
			s, e, ok := vim.QuoteBounds(text, curByte, '"', true)
			return s, e, ok
		})

	case vim.ActionVisualSelectInnerQuoteB:
		return p.visualSelectTextObject(func(text string, curByte int) (int, int, bool) {
			s, e, ok := vim.QuoteBounds(text, curByte, '`', true)
			return s, e, ok
		})

	case vim.ActionVisualSelectInnerParen:
		return p.visualSelectTextObject(func(text string, curByte int) (int, int, bool) {
			s, e, ok := vim.BracketBounds(text, curByte, '(', true)
			return s, e, ok
		})

	case vim.ActionVisualSelectInnerBrace:
		return p.visualSelectTextObject(func(text string, curByte int) (int, int, bool) {
			s, e, ok := vim.BracketBounds(text, curByte, '{', true)
			return s, e, ok
		})

	case vim.ActionVisualSelectInnerBracket:
		return p.visualSelectTextObject(func(text string, curByte int) (int, int, bool) {
			s, e, ok := vim.BracketBounds(text, curByte, '[', true)
			return s, e, ok
		})

	// ── Control ──────────────────────────────────────────────────────────────

	case vim.ActionEscape:
		p.engine.EnterMode(vim.ModeNormal, 0, 0)

	case vim.ActionFocusOut:
		return p, msgs.FocusCmd(msgs.FocusRequests)
	}

	return p, nil
}

// deleteTextObject deletes the text object identified by fn.
type textObjectFn func(text string, curByte int) (start, end int, ok bool)

func (p RequestDetailsPane) deleteTextObject(fn textObjectFn) (RequestDetailsPane, tea.Cmd) {
	text := p.textarea.Value()
	curByte := p.cursorByteOffset()
	s, e, ok := fn(text, curByte)
	if !ok {
		return p, nil
	}
	p.engine.PushUndo(text)
	p.engine.Register = text[s:e]
	newText := text[:s] + text[e:]
	targetLine, targetCol := vim.LineColFromByteOffset(newText, s)
	p.textarea.SetValue(newText)
	p.repositionCursor(newText, targetLine, targetCol)
	return p, nil
}

// changeTextObject deletes the text object and enters insert mode.
func (p RequestDetailsPane) changeTextObject(fn textObjectFn) (RequestDetailsPane, tea.Cmd) {
	text := p.textarea.Value()
	curByte := p.cursorByteOffset()
	s, e, ok := fn(text, curByte)
	if !ok {
		return p, nil
	}
	p.engine.PushUndo(text)
	p.engine.Register = text[s:e]
	newText := text[:s] + text[e:]
	targetLine, targetCol := vim.LineColFromByteOffset(newText, s)
	p.textarea.SetValue(newText)
	p.repositionCursor(newText, targetLine, targetCol)
	p.engine.EnterMode(vim.ModeInsert, 0, 0)
	return p, textarea.Blink
}

// visualSelectTextObject sets anchor to start and cursor to end of the text object,
// remaining in visual mode.
func (p RequestDetailsPane) visualSelectTextObject(fn textObjectFn) (RequestDetailsPane, tea.Cmd) {
	text := p.textarea.Value()
	curByte := p.cursorByteOffset()
	s, e, ok := fn(text, curByte)
	if !ok || e <= s {
		return p, nil
	}
	startLine, startCol := vim.LineColFromByteOffset(text, s)
	p.engine.AnchorLine = startLine
	p.engine.AnchorCol = startCol
	// Move cursor to the last character of the object (e is exclusive).
	endByte := e - 1
	if endByte < s {
		endByte = s
	}
	endLine, endCol := vim.LineColFromByteOffset(text, endByte)
	// Reset to (0,0) then reposition — SetValue with same content resets cursor.
	p.textarea.SetValue(text)
	p.repositionCursor(text, endLine, endCol)
	return p, nil
}

// cursorByteOffset returns the current cursor position as a byte offset into Value().
func (p *RequestDetailsPane) cursorByteOffset() int {
	line := p.textarea.Line()
	col := p.textarea.LineInfo().ColumnOffset
	return vim.ByteOffsetFromLineCol(p.textarea.Value(), line, col)
}

// repositionCursor moves the textarea cursor to (targetLine, targetCol).
// SetValue leaves the cursor at the END of text (it calls InsertString internally),
// so we always call MoveToBegin first, then CursorDown N times, then SetCursorColumn.
func (p *RequestDetailsPane) repositionCursor(text string, targetLine, targetCol int) {
	p.textarea.MoveToBegin()
	clamp := strings.Count(text, "\n")
	if targetLine > clamp {
		targetLine = clamp
	}
	for i := 0; i < targetLine; i++ {
		p.textarea.CursorDown()
	}
	p.textarea.SetCursorColumn(targetCol)
}

// visualByteRange returns the [start, end) byte range of the visual selection.
func (p *RequestDetailsPane) visualByteRange(text string) (start, end int) {
	curLine := p.textarea.Line()
	curCol := p.textarea.LineInfo().ColumnOffset
	anchorLine := p.engine.AnchorLine
	anchorCol := p.engine.AnchorCol

	switch p.engine.Mode {
	case vim.ModeVisual:
		s := vim.ByteOffsetFromLineCol(text, anchorLine, anchorCol)
		e := vim.ByteOffsetFromLineCol(text, curLine, curCol)
		if s > e {
			s, e = e, s
		}
		// include the char under cursor
		if e < len(text) {
			_, size := utf8.DecodeRuneInString(text[e:])
			e += size
		}
		return s, e

	case vim.ModeVisualLine:
		minLine, maxLine := anchorLine, curLine
		if minLine > maxLine {
			minLine, maxLine = maxLine, minLine
		}
		s, _ := vim.LineBounds(text, vim.ByteOffsetFromLineCol(text, minLine, 0), false)
		_, e := vim.LineBounds(text, vim.ByteOffsetFromLineCol(text, maxLine, 0), false)
		return s, e

	case vim.ModeVisualBlock:
		// For block: return contiguous range from top-left to bottom-right
		minLine, maxLine := anchorLine, curLine
		if minLine > maxLine {
			minLine, maxLine = maxLine, minLine
		}
		minCol, maxCol := anchorCol, curCol
		if minCol > maxCol {
			minCol, maxCol = maxCol, minCol
		}
		s := vim.ByteOffsetFromLineCol(text, minLine, minCol)
		e := vim.ByteOffsetFromLineCol(text, maxLine, maxCol)
		return s, e
	}
	return 0, 0
}

func (p RequestDetailsPane) View() string {
	border := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(p.width - 2).
		Height(p.height - 2)

	if p.focused {
		border = border.BorderForeground(lipgloss.Color("2"))
	} else {
		border = border.BorderForeground(lipgloss.Color("8"))
	}

	// Mode indicator
	var modeStr string
	switch p.engine.Mode {
	case vim.ModeInsert:
		modeStr = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render("INSERT")
	case vim.ModeVisual:
		modeStr = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Render("VISUAL")
	case vim.ModeVisualLine:
		modeStr = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Render("V-LINE")
	case vim.ModeVisualBlock:
		modeStr = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Render("V-BLOCK")
	default:
		modeStr = "NORMAL"
	}

	pending := p.engine.PendingBuf()
	titleBase := lipgloss.NewStyle().Bold(true).Render("Request Details")
	title := titleBase + "  [" + modeStr + "]"
	if pending != "" {
		title += "  " + lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(pending+"…")
	}

	p.textarea.SetWidth(p.width - 4)
	p.textarea.SetHeight(p.height - 5)

	var body string
	switch p.engine.Mode {
	case vim.ModeVisual, vim.ModeVisualLine, vim.ModeVisualBlock:
		body = p.renderVisual()
	default:
		body = p.textarea.View()
	}

	content := title + "\n" + body
	return border.Render(content)
}

// renderVisual renders the textarea content with visual selection highlighted.
func (p RequestDetailsPane) renderVisual() string {
	text := p.textarea.Value()
	lines := strings.Split(text, "\n")

	curLine := p.textarea.Line()
	curCol := p.textarea.LineInfo().ColumnOffset
	anchorLine := p.engine.AnchorLine
	anchorCol := p.engine.AnchorCol

	// Compute visible range
	visibleLines := p.height - 5
	if visibleLines < 1 {
		visibleLines = 1
	}
	startRender := curLine - visibleLines/2
	if startRender < 0 {
		startRender = 0
	}
	endRender := startRender + visibleLines
	if endRender > len(lines) {
		endRender = len(lines)
	}

	var sb strings.Builder
	for lineIdx := startRender; lineIdx < endRender; lineIdx++ {
		if lineIdx > startRender {
			sb.WriteByte('\n')
		}
		runes := []rune(lines[lineIdx])
		for colIdx, r := range runes {
			inSel := isInSelection(lineIdx, colIdx, anchorLine, anchorCol, curLine, curCol, len(runes), p.engine.Mode)
			isCursor := lineIdx == curLine && colIdx == curCol

			st := lipgloss.NewStyle()
			if isCursor {
				st = st.Underline(true)
			}
			if inSel {
				st = st.Reverse(true)
			}
			sb.WriteString(st.Render(string(r)))
		}
		// Cursor at end of line
		if curLine == lineIdx && curCol >= len(runes) {
			sb.WriteString(lipgloss.NewStyle().Underline(true).Render(" "))
		}
	}
	return sb.String()
}

// isInSelection returns true if (line, col) falls within the visual selection.
func isInSelection(line, col, anchorLine, anchorCol, curLine, curCol, lineLen int, mode vim.Mode) bool {
	switch mode {
	case vim.ModeVisual:
		// Determine ordered start/end
		startLine, startCol := anchorLine, anchorCol
		endLine, endCol := curLine, curCol
		if startLine > endLine || (startLine == endLine && startCol > endCol) {
			startLine, startCol, endLine, endCol = endLine, endCol, startLine, startCol
		}
		if line < startLine || line > endLine {
			return false
		}
		if line == startLine && line == endLine {
			return col >= startCol && col <= endCol
		}
		if line == startLine {
			return col >= startCol
		}
		if line == endLine {
			return col <= endCol
		}
		return true

	case vim.ModeVisualLine:
		minLine, maxLine := anchorLine, curLine
		if minLine > maxLine {
			minLine, maxLine = maxLine, minLine
		}
		return line >= minLine && line <= maxLine

	case vim.ModeVisualBlock:
		minLine, maxLine := anchorLine, curLine
		if minLine > maxLine {
			minLine, maxLine = maxLine, minLine
		}
		minCol, maxCol := anchorCol, curCol
		if minCol > maxCol {
			minCol, maxCol = maxCol, minCol
		}
		if line < minLine || line > maxLine {
			return false
		}
		clampedCol := col
		if clampedCol > lineLen-1 {
			return false
		}
		return clampedCol >= minCol && clampedCol <= maxCol
	}
	return false
}

func (p *RequestDetailsPane) SetSize(w, h int) {
	p.width = w
	p.height = h
}

func (p *RequestDetailsPane) SetFocused(f bool) {
	p.focused = f
	if f {
		p.engine.EnterMode(vim.ModeNormal, 0, 0)
		p.textarea.Focus()
	} else {
		p.engine.EnterMode(vim.ModeNormal, 0, 0)
		p.textarea.Blur()
	}
}
