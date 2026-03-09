package vim

// Mode represents the current vim editing mode.
type Mode int

const (
	ModeNormal      Mode = iota // Normal mode
	ModeInsert                  // Insert mode
	ModeVisual                  // v — character-wise visual
	ModeVisualLine              // V — line-wise visual
	ModeVisualBlock             // ctrl+v — column block visual
)

func (m Mode) String() string {
	switch m {
	case ModeNormal:
		return "NORMAL"
	case ModeInsert:
		return "INSERT"
	case ModeVisual:
		return "VISUAL"
	case ModeVisualLine:
		return "V-LINE"
	case ModeVisualBlock:
		return "V-BLOCK"
	default:
		return "NORMAL"
	}
}
