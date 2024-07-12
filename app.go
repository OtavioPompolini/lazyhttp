package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jroimartin/gocui"
)

type RequestWindow struct {
	Name          string
	X, Y          int
	Width, Height int
	SelectKey     gocui.Key
}

type Window struct {
	Title         string
	Name          string
	X, Y          int
	Width, Height int
	SelectKey     gocui.Key
	Content       string
}

type Request struct {
	Name string
}

var (
	requests = []Request{
		{"pudim"},
		{"pudim1"},
		{"pudim2"},
		{"pudim3"},
		{"pudim4"},
		{"pudim5"},
	}
)

func newWindow(name string, x, y, w, h int, content string, key gocui.Key) *Window {
	return &Window{Name: name, X: x, Y: y, Width: w, Height: h, Content: content, SelectKey: key, Title: "Requests"}

}

func (w *Window) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.Name, w.X, w.Y, w.X+w.Width, w.Y+w.Height)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	v.Clear()

	for _, request := range requests {
		fmt.Fprintln(v, request.Name)
	}
	v.Title = w.Title

	if err := g.SetKeybinding("", w.SelectKey, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err = g.SetCurrentView(w.Name)
		return err
	}); err != nil {
		log.Panicln(err)
	}
	return nil
}

type AddRequestWindow struct {
	Name          string
	Title         string
	X, Y          int
	Width, Height int
}

func newAddRequestWindow(a, b int) *AddRequestWindow {

	return &AddRequestWindow{
		Name:   "addNewRequest",
		Title:  "puidm",
		X:      a / 2,
		Y:      b / 2,
		Width:  20,
		Height: 2,
	}
}

func (arw *AddRequestWindow) Layout(g *gocui.Gui) error {
	v, err := g.SetView(arw.Name, arw.X, arw.Y, arw.X+arw.Width, arw.Y+arw.Height)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Title = arw.Title

		if err := g.SetKeybinding(arw.Name, gocui.KeyF6, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			test := v.Buffer()
			requests = append(requests, Request{test})
			return nil
		}); err != nil {
			log.Panicln(err)
		}
	}
	return nil
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	g.Highlight = true
	g.SelFgColor = gocui.ColorBlack

	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	view1 := newWindow("requests", 0, 0, 14, 20, strconv.FormatBool(g.Highlight), gocui.KeyF1)
	// view2 := newWindow("env", 0, 3, 14, 2, "marcos viado", gocui.KeyF2)
	addRequestWindow := newAddRequestWindow(g.Size())

	g.SetManager(view1, addRequestWindow)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", 'a', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err = g.SetCurrentView(addRequestWindow.Name)
		g.SetViewOnTop(addRequestWindow.Name)
		return err
	}); err != nil {
		log.Panicln(err)
	}

	g.SetCurrentView("requests")

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
