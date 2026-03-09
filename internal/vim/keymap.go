package vim

// KeyMap maps key sequences to Actions.
type KeyMap map[string]Action

// DefaultKeyMap returns the default key bindings for the given mode.
func DefaultKeyMap(mode Mode) KeyMap {
	switch mode {
	case ModeNormal:
		return KeyMap{
			// Motions
			"h":      ActionMoveLeft,
			"l":      ActionMoveRight,
			"j":      ActionMoveDown,
			"k":      ActionMoveUp,
			"w":      ActionMoveWordForward,
			"e":      ActionMoveWordEnd,
			"b":      ActionMoveWordBack,
			"W":      ActionMoveWORDForward,
			"E":      ActionMoveWORDEnd,
			"B":      ActionMoveWORDBack,
			"0":      ActionMoveLineStart,
			"^":      ActionMoveLineFirst,
			"$":      ActionMoveLineEnd,
			"gg":     ActionMoveDocStart,
			"G":      ActionMoveDocEnd,

			// Insert-entry
			"i": ActionInsertBefore,
			"a": ActionInsertAfter,
			"I": ActionInsertLineStart,
			"A": ActionInsertLineEnd,
			"o": ActionNewLineBelow,
			"O": ActionNewLineAbove,

			// Delete ops
			"dd":  ActionDeleteLine,
			"dw":  ActionDeleteWord,
			"D":   ActionDeleteToLineEnd,
			"diw": ActionDeleteInnerWord,
			"diW": ActionDeleteInnerWORD,
			"di'": ActionDeleteInnerQuoteS,
			`di"`: ActionDeleteInnerQuoteD,
			"di`": ActionDeleteInnerQuoteB,
			"di(": ActionDeleteInnerParen,
			"di{": ActionDeleteInnerBrace,
			"di[": ActionDeleteInnerBracket,

			// Change ops
			"cc":  ActionChangeLine,
			"cw":  ActionChangeWord,
			"C":   ActionChangeToLineEnd,
			"ciw": ActionChangeInnerWord,
			"ciW": ActionChangeInnerWORD,
			"ci'": ActionChangeInnerQuoteS,
			`ci"`: ActionChangeInnerQuoteD,
			"ci`": ActionChangeInnerQuoteB,
			"ci(": ActionChangeInnerParen,
			"ci{": ActionChangeInnerBrace,
			"ci[": ActionChangeInnerBracket,

			// Yank ops
			"yy":  ActionYankLine,
			"yw":  ActionYankWord,
			"yiw": ActionYankInnerWord,

			// Clipboard
			"p": ActionPasteAfter,
			"P": ActionPasteBefore,

			// Undo/Redo
			"u":      ActionUndo,
			"ctrl+r": ActionRedo,

			// Visual
			"v":      ActionVisualChar,
			"V":      ActionVisualLine,
			"ctrl+v": ActionVisualBlock,

			// Control
			"esc": ActionFocusOut,
		}

	case ModeInsert:
		return KeyMap{
			"esc": ActionEscape,
		}

	case ModeVisual, ModeVisualLine, ModeVisualBlock:
		return KeyMap{
			// Motions (same as normal)
			"h":  ActionMoveLeft,
			"l":  ActionMoveRight,
			"j":  ActionMoveDown,
			"k":  ActionMoveUp,
			"w":  ActionMoveWordForward,
			"e":  ActionMoveWordEnd,
			"b":  ActionMoveWordBack,
			"W":  ActionMoveWORDForward,
			"E":  ActionMoveWORDEnd,
			"B":  ActionMoveWORDBack,
			"0":  ActionMoveLineStart,
			"^":  ActionMoveLineFirst,
			"$":  ActionMoveLineEnd,
			"gg": ActionMoveDocStart,
			"G":  ActionMoveDocEnd,

			// Visual ops
			"d": ActionVisualDelete,
			"y": ActionVisualYank,
			"c": ActionVisualChange,

			// Text object selections
			"iw":  ActionVisualSelectInnerWord,
			"iW":  ActionVisualSelectInnerWORD,
			"i'":  ActionVisualSelectInnerQuoteS,
			`i"`:  ActionVisualSelectInnerQuoteD,
			"i`":  ActionVisualSelectInnerQuoteB,
			"i(":  ActionVisualSelectInnerParen,
			"i{":  ActionVisualSelectInnerBrace,
			"i[":  ActionVisualSelectInnerBracket,

			// Mode changes
			"v":   ActionEscape, // v again exits visual
			"V":   ActionEscape,
			"esc": ActionEscape,
		}
	}
	return KeyMap{}
}
