package vim

// Action represents a semantic vim operation.
type Action int

const (
	ActionNone Action = iota

	// Motions
	ActionMoveLeft
	ActionMoveRight
	ActionMoveUp
	ActionMoveDown
	ActionMoveWordForward
	ActionMoveWordEnd
	ActionMoveWordBack
	ActionMoveWORDForward
	ActionMoveWORDEnd
	ActionMoveWORDBack
	ActionMoveLineStart  // 0
	ActionMoveLineFirst  // ^
	ActionMoveLineEnd    // $
	ActionMoveDocStart   // gg
	ActionMoveDocEnd     // G

	// Insert-entry actions
	ActionInsertBefore    // i
	ActionInsertAfter     // a
	ActionInsertLineStart // I
	ActionInsertLineEnd   // A
	ActionNewLineBelow    // o
	ActionNewLineAbove    // O

	// Pre-composed delete ops
	ActionDeleteLine        // dd
	ActionDeleteWord        // dw
	ActionDeleteToLineEnd   // D
	ActionDeleteInnerWord   // diw
	ActionDeleteInnerWORD   // diW
	ActionDeleteInnerQuoteS // di'
	ActionDeleteInnerQuoteD // di"
	ActionDeleteInnerQuoteB // di`
	ActionDeleteInnerParen  // di(
	ActionDeleteInnerBrace  // di{
	ActionDeleteInnerBracket // di[

	// Pre-composed change ops
	ActionChangeLine        // cc
	ActionChangeWord        // cw
	ActionChangeToLineEnd   // C
	ActionChangeInnerWord   // ciw
	ActionChangeInnerWORD   // ciW
	ActionChangeInnerQuoteS // ci'
	ActionChangeInnerQuoteD // ci"
	ActionChangeInnerQuoteB // ci`
	ActionChangeInnerParen  // ci(
	ActionChangeInnerBrace  // ci{
	ActionChangeInnerBracket // ci[

	// Pre-composed yank ops
	ActionYankLine      // yy
	ActionYankWord      // yw
	ActionYankInnerWord // yiw

	// Clipboard
	ActionPasteAfter  // p
	ActionPasteBefore // P

	// Undo/Redo
	ActionUndo // u
	ActionRedo // ctrl+r

	// Visual mode entry
	ActionVisualChar  // v
	ActionVisualLine  // V
	ActionVisualBlock // ctrl+v

	// Visual mode operations
	ActionVisualDelete
	ActionVisualYank
	ActionVisualChange

	// Visual mode text object selections (iw, iW, i', i", i`, i(, i{, i[)
	ActionVisualSelectInnerWord
	ActionVisualSelectInnerWORD
	ActionVisualSelectInnerQuoteS
	ActionVisualSelectInnerQuoteD
	ActionVisualSelectInnerQuoteB
	ActionVisualSelectInnerParen
	ActionVisualSelectInnerBrace
	ActionVisualSelectInnerBracket

	// Control
	ActionEscape   // esc in insert/visual → Normal
	ActionFocusOut // esc in Normal → exit pane
)

// ActionName maps Lua rhs strings to Action values.
var ActionName = map[string]Action{
	"none": ActionNone,

	"move_left":         ActionMoveLeft,
	"move_right":        ActionMoveRight,
	"move_up":           ActionMoveUp,
	"move_down":         ActionMoveDown,
	"move_word_forward": ActionMoveWordForward,
	"move_word_end":     ActionMoveWordEnd,
	"move_word_back":    ActionMoveWordBack,
	"move_WORD_forward": ActionMoveWORDForward,
	"move_WORD_end":     ActionMoveWORDEnd,
	"move_WORD_back":    ActionMoveWORDBack,
	"move_line_start":   ActionMoveLineStart,
	"move_line_first":   ActionMoveLineFirst,
	"move_line_end":     ActionMoveLineEnd,
	"move_doc_start":    ActionMoveDocStart,
	"move_doc_end":      ActionMoveDocEnd,

	"insert_before":     ActionInsertBefore,
	"insert_after":      ActionInsertAfter,
	"insert_line_start": ActionInsertLineStart,
	"insert_line_end":   ActionInsertLineEnd,
	"new_line_below":    ActionNewLineBelow,
	"new_line_above":    ActionNewLineAbove,

	"delete_line":          ActionDeleteLine,
	"delete_word":          ActionDeleteWord,
	"delete_to_line_end":   ActionDeleteToLineEnd,
	"delete_inner_word":    ActionDeleteInnerWord,
	"delete_inner_WORD":    ActionDeleteInnerWORD,
	"delete_inner_quote_s": ActionDeleteInnerQuoteS,
	"delete_inner_quote_d": ActionDeleteInnerQuoteD,
	"delete_inner_quote_b": ActionDeleteInnerQuoteB,
	"delete_inner_paren":   ActionDeleteInnerParen,
	"delete_inner_brace":   ActionDeleteInnerBrace,
	"delete_inner_bracket": ActionDeleteInnerBracket,

	"change_line":          ActionChangeLine,
	"change_word":          ActionChangeWord,
	"change_to_line_end":   ActionChangeToLineEnd,
	"change_inner_word":    ActionChangeInnerWord,
	"change_inner_WORD":    ActionChangeInnerWORD,
	"change_inner_quote_s": ActionChangeInnerQuoteS,
	"change_inner_quote_d": ActionChangeInnerQuoteD,
	"change_inner_quote_b": ActionChangeInnerQuoteB,
	"change_inner_paren":   ActionChangeInnerParen,
	"change_inner_brace":   ActionChangeInnerBrace,
	"change_inner_bracket": ActionChangeInnerBracket,

	"yank_line":       ActionYankLine,
	"yank_word":       ActionYankWord,
	"yank_inner_word": ActionYankInnerWord,

	"paste_after":  ActionPasteAfter,
	"paste_before": ActionPasteBefore,

	"undo": ActionUndo,
	"redo": ActionRedo,

	"visual_char":  ActionVisualChar,
	"visual_line":  ActionVisualLine,
	"visual_block": ActionVisualBlock,

	"visual_delete": ActionVisualDelete,
	"visual_yank":   ActionVisualYank,
	"visual_change": ActionVisualChange,

	"visual_select_inner_word":    ActionVisualSelectInnerWord,
	"visual_select_inner_WORD":    ActionVisualSelectInnerWORD,
	"visual_select_inner_quote_s": ActionVisualSelectInnerQuoteS,
	"visual_select_inner_quote_d": ActionVisualSelectInnerQuoteD,
	"visual_select_inner_quote_b": ActionVisualSelectInnerQuoteB,
	"visual_select_inner_paren":   ActionVisualSelectInnerParen,
	"visual_select_inner_brace":   ActionVisualSelectInnerBrace,
	"visual_select_inner_bracket": ActionVisualSelectInnerBracket,

	"escape":    ActionEscape,
	"focus_out": ActionFocusOut,
}
