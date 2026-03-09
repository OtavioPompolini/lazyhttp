package vim

const maxUndoStack = 100

// Engine is the central vim state machine.
type Engine struct {
	Mode                  Mode
	AnchorLine, AnchorCol int    // visual selection anchor
	Register              string // unnamed yank register

	undoStack []string
	redoStack []string

	parser    Parser
	normalMap KeyMap
	insertMap KeyMap
	visualMap KeyMap // shared for all 3 visual sub-modes
}

// NewEngine creates an Engine with default keymaps.
func NewEngine() Engine {
	return Engine{
		Mode:      ModeNormal,
		normalMap: DefaultKeyMap(ModeNormal),
		insertMap: DefaultKeyMap(ModeInsert),
		visualMap: DefaultKeyMap(ModeVisual),
	}
}

// ProcessKey routes the key through the correct keymap and returns the
// resolved action, count, and whether the key should pass through to the textarea.
func (e *Engine) ProcessKey(key string) (action Action, count int, passThru bool) {
	var km KeyMap
	switch e.Mode {
	case ModeInsert:
		km = e.insertMap
	case ModeVisual, ModeVisualLine, ModeVisualBlock:
		km = e.visualMap
	default:
		km = e.normalMap
	}

	act, cnt, pending, reset := e.parser.Feed(key, km)

	switch {
	case pending:
		// Waiting for more keys — consume silently
		return ActionNone, 0, false
	case reset:
		// Sequence didn't match any binding
		if e.Mode == ModeInsert {
			// In insert mode, unrecognized keys pass through to textarea
			return ActionNone, 0, true
		}
		return ActionNone, 0, false
	default:
		// Resolved an action
		resolvedCount := cnt
		if resolvedCount == 0 {
			resolvedCount = 1
		}
		return act, resolvedCount, false
	}
}

// SetLeader sets the leader key on the parser.
func (e *Engine) SetLeader(key string) {
	e.parser.SetLeader(key)
}

// BindKey adds or overrides a binding in the given mode's keymap.
func (e *Engine) BindKey(mode Mode, seq string, action Action) {
	switch mode {
	case ModeNormal:
		e.normalMap[seq] = action
	case ModeInsert:
		e.insertMap[seq] = action
	case ModeVisual, ModeVisualLine, ModeVisualBlock:
		e.visualMap[seq] = action
	}
}

// EnterMode transitions to a new mode and resets the parser.
// For visual modes, sets the anchor to (curLine, curCol).
func (e *Engine) EnterMode(mode Mode, curLine, curCol int) {
	e.Mode = mode
	e.parser.Reset()
	if mode == ModeVisual || mode == ModeVisualLine || mode == ModeVisualBlock {
		e.AnchorLine = curLine
		e.AnchorCol = curCol
	}
}

// PushUndo saves the current text value onto the undo stack.
func (e *Engine) PushUndo(value string) {
	if len(e.undoStack) >= maxUndoStack {
		e.undoStack = e.undoStack[1:]
	}
	e.undoStack = append(e.undoStack, value)
	// A new edit clears the redo stack
	e.redoStack = nil
}

// PopUndo pops the most recent undo snapshot.
func (e *Engine) PopUndo() (string, bool) {
	if len(e.undoStack) == 0 {
		return "", false
	}
	val := e.undoStack[len(e.undoStack)-1]
	e.undoStack = e.undoStack[:len(e.undoStack)-1]
	return val, true
}

// PopRedo pops the most recent redo snapshot.
func (e *Engine) PopRedo() (string, bool) {
	if len(e.redoStack) == 0 {
		return "", false
	}
	val := e.redoStack[len(e.redoStack)-1]
	e.redoStack = e.redoStack[:len(e.redoStack)-1]
	return val, true
}

// PushRedo saves a value onto the redo stack (called by undo logic).
func (e *Engine) PushRedo(value string) {
	e.redoStack = append(e.redoStack, value)
}

// PendingBuf delegates to the parser.
func (e *Engine) PendingBuf() string {
	return e.parser.PendingBuf()
}
