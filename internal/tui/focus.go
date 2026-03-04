package tui

import "github.com/OtavioPompolini/project-postman/internal/tui/msgs"

type FocusManager struct {
	current msgs.FocusTarget
	prev    msgs.FocusTarget
}

func newFocusManager() FocusManager {
	return FocusManager{current: msgs.FocusCollections, prev: msgs.FocusCollections}
}

func (f *FocusManager) Current() msgs.FocusTarget {
	return f.current
}

func (f *FocusManager) MoveTo(target msgs.FocusTarget) {
	f.prev = f.current
	f.current = target
}

func (f *FocusManager) RestorePrev() {
	f.current, f.prev = f.prev, f.current
}
