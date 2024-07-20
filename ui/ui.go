package ui

import (
	"github.com/jroimartin/gocui"
)


type UI struct {
	g *gocui.Gui
}

func NewUI() (*UI, error){
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}

	return &UI{
		g: g,
	}, nil
}

func (ui *UI) SetHightlight(e bool) {
	ui.g.Highlight = e
}

func (ui *UI) SetFgColor(color gocui.Attribute) {
	ui.g.FgColor = color
}

func (ui *UI) SetManager(views ...View) {
}

func (ui *UI) MainLoop() error {
	return ui.g.MainLoop()
}

func (ui *UI) Close() {
	ui.g.Close()
}

func (ui *UI) ActiveViewName() string {
	return ui.g.CurrentView().Name()
}
