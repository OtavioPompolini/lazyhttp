package vim

import (
	"testing"
)

func TestLoadConfig_BasicBinding(t *testing.T) {
	engine := NewEngine()
	err := LoadConfigString(`
vim.keymap.set("n", "H", "move_line_start")
`, &engine)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// "H" should now map to ActionMoveLineStart in normal mode
	act, _, passThru := engine.ProcessKey("H")
	if act != ActionMoveLineStart {
		t.Errorf("expected ActionMoveLineStart, got %v", act)
	}
	if passThru {
		t.Errorf("expected passThru=false")
	}
}

func TestLoadConfig_Leader(t *testing.T) {
	engine := NewEngine()
	err := LoadConfigString(`
vim.g.mapleader = ","
vim.keymap.set("n", "<leader>s", "focus_out")
`, &engine)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// leader "," expanded at bind time → keymap has ",s"
	// First "," is a pending prefix
	_, _, passThru := engine.ProcessKey(",")
	_ = passThru

	// "s" resolves the sequence
	act, _, _ := engine.ProcessKey("s")
	if act != ActionFocusOut {
		t.Errorf("expected ActionFocusOut, got %v", act)
	}
}

func TestLoadConfig_MultiMode(t *testing.T) {
	engine := NewEngine()
	err := LoadConfigString(`
vim.keymap.set({"n", "v"}, "X", "delete_line")
`, &engine)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should be bound in normal
	act, _, _ := engine.ProcessKey("X")
	if act != ActionDeleteLine {
		t.Errorf("normal: expected ActionDeleteLine, got %v", act)
	}

	// And in visual
	engine.EnterMode(ModeVisual, 0, 0)
	act, _, _ = engine.ProcessKey("X")
	if act != ActionDeleteLine {
		t.Errorf("visual: expected ActionDeleteLine, got %v", act)
	}
}

func TestLoadConfig_UnknownAction(t *testing.T) {
	engine := NewEngine()
	// Unknown action should not error, just warn
	err := LoadConfigString(`
vim.keymap.set("n", "Z", "nonexistent_action")
`, &engine)
	if err != nil {
		t.Fatalf("unexpected error for unknown action: %v", err)
	}
}

func TestLoadConfig_NonexistentFile(t *testing.T) {
	engine := NewEngine()
	err := LoadConfig("/tmp/nonexistent_lazyhttp_config.lua", &engine)
	if err != nil {
		t.Errorf("expected nil for nonexistent file, got %v", err)
	}
}
