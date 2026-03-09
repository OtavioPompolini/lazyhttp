package vim

import (
	"testing"
)

func TestParser_SimpleMotion(t *testing.T) {
	km := DefaultKeyMap(ModeNormal)
	p := &Parser{}

	action, count, pending, reset := p.Feed("j", km)
	if action != ActionMoveDown {
		t.Errorf("expected ActionMoveDown, got %v", action)
	}
	if count != 0 {
		t.Errorf("expected count 0, got %d", count)
	}
	if pending || reset {
		t.Errorf("unexpected pending=%v reset=%v", pending, reset)
	}
}

func TestParser_CountMotion(t *testing.T) {
	km := DefaultKeyMap(ModeNormal)
	p := &Parser{}

	// Feed "5"
	action, _, pending, _ := p.Feed("5", km)
	if action != ActionNone || !pending {
		t.Errorf("after '5': expected ActionNone pending=true, got action=%v pending=%v", action, pending)
	}

	// Feed "j"
	action, count, pending, reset := p.Feed("j", km)
	if action != ActionMoveDown {
		t.Errorf("expected ActionMoveDown after '5j', got %v", action)
	}
	if count != 5 {
		t.Errorf("expected count=5, got %d", count)
	}
	if pending || reset {
		t.Errorf("unexpected pending=%v reset=%v", pending, reset)
	}
}

func TestParser_ZeroIsLineStart(t *testing.T) {
	km := DefaultKeyMap(ModeNormal)
	p := &Parser{}

	// "0" with no prior count = line start motion
	action, _, pending, reset := p.Feed("0", km)
	if action != ActionMoveLineStart {
		t.Errorf("expected ActionMoveLineStart for bare '0', got %v", action)
	}
	if pending || reset {
		t.Errorf("unexpected pending=%v reset=%v", pending, reset)
	}
}

func TestParser_ZeroAfterCountIsDigit(t *testing.T) {
	km := DefaultKeyMap(ModeNormal)
	p := &Parser{}

	// "1" starts count
	p.Feed("1", km)
	// "0" should accumulate as digit (count = 10)
	action, _, pending, _ := p.Feed("0", km)
	if action != ActionNone || !pending {
		t.Errorf("after '10': expected ActionNone pending=true, got action=%v pending=%v", action, pending)
	}

	// "j" resolves
	action, count, _, _ := p.Feed("j", km)
	if action != ActionMoveDown {
		t.Errorf("expected ActionMoveDown after '10j'")
	}
	if count != 10 {
		t.Errorf("expected count=10, got %d", count)
	}
}

func TestParser_TwoCharSequence(t *testing.T) {
	km := DefaultKeyMap(ModeNormal)
	p := &Parser{}

	// First "d" — pending prefix
	action, _, pending, reset := p.Feed("d", km)
	if action != ActionNone || !pending {
		t.Errorf("after 'd': expected pending, got action=%v pending=%v reset=%v", action, pending, reset)
	}

	// Second "d" → dd → delete line
	action, _, pending, reset = p.Feed("d", km)
	if action != ActionDeleteLine {
		t.Errorf("expected ActionDeleteLine after 'dd', got %v", action)
	}
	if pending || reset {
		t.Errorf("unexpected pending=%v reset=%v", pending, reset)
	}
}

func TestParser_ThreeCharSequence(t *testing.T) {
	km := DefaultKeyMap(ModeNormal)
	p := &Parser{}

	p.Feed("c", km)
	p.Feed("i", km)
	action, _, _, _ := p.Feed("w", km)
	if action != ActionChangeInnerWord {
		t.Errorf("expected ActionChangeInnerWord after 'ciw', got %v", action)
	}
}

func TestParser_InvalidSequenceReset(t *testing.T) {
	km := DefaultKeyMap(ModeNormal)
	p := &Parser{}

	// "d" then "z" (dz is not a binding)
	p.Feed("d", km)
	action, _, pending, reset := p.Feed("z", km)
	if action != ActionNone || pending || !reset {
		t.Errorf("expected reset after invalid sequence 'dz', got action=%v pending=%v reset=%v", action, pending, reset)
	}
	// After reset, buf should be empty
	if p.HasPending() {
		t.Errorf("expected buf to be empty after reset")
	}
}

func TestParser_Leader(t *testing.T) {
	// Leader key is expanded at bind time; keymap has the raw leader char + key.
	// e.g. leader=" ", binding "<leader>s" → " s" in keymap.
	km := KeyMap{
		" s": ActionNone, // space+s for <leader>s with leader=" "
	}
	p := &Parser{}
	p.SetLeader(" ")

	// Feed space (leader key) — pending prefix " s"
	action, _, pending, _ := p.Feed(" ", km)
	if action != ActionNone || !pending {
		t.Errorf("after leader: expected pending, got action=%v pending=%v", action, pending)
	}

	// Feed "s" — resolves to ActionNone
	action, _, pending, reset := p.Feed("s", km)
	if action != ActionNone {
		t.Errorf("expected ActionNone (mapped action), got %v", action)
	}
	if pending || reset {
		t.Errorf("unexpected pending=%v reset=%v", pending, reset)
	}
}

func TestParser_InsertMode(t *testing.T) {
	km := DefaultKeyMap(ModeInsert)
	p := &Parser{}

	// In insert mode, only esc is mapped
	action, _, _, _ := p.Feed("esc", km)
	if action != ActionEscape {
		t.Errorf("expected ActionEscape for 'esc' in insert mode, got %v", action)
	}

	// Other keys reset (not a valid prefix)
	action, _, _, reset := p.Feed("j", km)
	if !reset {
		t.Errorf("expected reset for 'j' in insert mode (should pass through)")
	}
	_ = action
}
