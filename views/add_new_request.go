package views

import (
	"github.com/OtavioPompolini/project-postman/types"
	"github.com/jroimartin/gocui"
)

const (
	AddNewRequestName = "addNewRequest"
)

type AddRequestWindow struct {
	Name          string
	Title         string
	X, Y          int
	Width, Height int
	OnAddRequest  func(types.Request)
}

func NewAddNewRequestView(g *gocui.Gui) (*AddRequestWindow, error) {
	a, b := g.Size()
	arw := &AddRequestWindow{
		Name:   AddNewRequestName,
		Title:  "puidm",
		X:      a / 2,
		Y:      b / 2,
		Width:  20,
		Height: 2,
	}

	v, err := g.SetView(arw.Name, arw.X, arw.Y, arw.X+arw.Width, arw.Y+arw.Height)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}

		if _, err := g.SetCurrentView(arw.Name); err != nil {
			return nil, err
		}

		v.Editable = true
		v.Title = arw.Title
	}
	return arw, nil
}
