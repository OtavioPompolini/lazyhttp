package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/jroimartin/gocui"
	"github.com/OtavioPompolini/project-postman/views"
)

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
	Requests = []Request{
		{"pudim"},
		{"pudim1"},
		{"pudim2"},
		{"pudim3"},
		{"pudim4"},
		{"pudim5"},
	}
)

func newWindow(name string, x, y, w, h int, content string, key gocui.Key) *Window {
	return &Window{
		Name: name,
		X: x,
		Y: y,
		Width: w,
		Height: h,
		Content: content,
		SelectKey: key,
		Title: "Requests",
	}
}

func (w *Window) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.Name, w.X, w.Y, w.X+w.Width, w.Y+w.Height)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	if err == gocui.ErrUnknownView {
		_, err := g.SetCurrentView(w.Name)
		return err
	}

	v.Clear()

	for _, request := range Requests {
		fmt.Fprintln(v, request.Name)
	}
	v.Title = w.Title
	v.SelFgColor = gocui.ColorRed

	if err := g.SetKeybinding("", w.SelectKey, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		_, err = g.SetCurrentView(w.Name)
		return err
	}); err != nil {
		log.Panicln(err)
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

	g.SetManager(view1)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("requests", 'n', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		err := views.OpenAddNewRequestView(g)
		return err
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding(views.AddNewRequestName, gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		read_line := v.Buffer()
		read_line = strings.TrimSuffix(read_line, "\n")
		Requests = append(Requests, Request{read_line})
		g.DeleteView("addNewRequest")
		g.SetCurrentView("requests")
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}


func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
