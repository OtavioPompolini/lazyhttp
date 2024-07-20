package ui

import "github.com/jroimartin/gocui"

type View struct {
	v gocui.View
	name string
}

func NewView(name string) *View {
	return &View{
		name: name,
	}
}

func (v *View) Name() string {
	return v.name
}


