package gui

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/juanpablocruz/shogo/clientr/internal/theme"
)

type GUI struct {
	Screen   *tcell.Screen
	boxStyle tcell.Style
	defStyle tcell.Style

	HandleEventMouse func(tcell.Event)
	HandleEventKey   func(tcell.Event)

	Theme theme.Theme
}

func NewGUI() *GUI {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	boxStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	return &GUI{
		Screen:   &s,
		defStyle: defStyle,
		boxStyle: boxStyle,
	}
}

func (gui GUI) drawText(x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range string(text) {
		(*gui.Screen).SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func (gui GUI) drawBox(x1, y1, x2, y2 int, style tcell.Style, text string) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			(*gui.Screen).SetContent(col, row, ' ', nil, style)
		}
	}
	// Draw borders
	for col := x1; col <= x2; col++ {
		(*gui.Screen).SetContent(col, y1, tcell.RuneHLine, nil, style)
		(*gui.Screen).SetContent(col, y2, tcell.RuneHLine, nil, style)
	}

	for row := y1; row <= y2; row++ {
		(*gui.Screen).SetContent(x1, row, tcell.RuneVLine, nil, style)
		(*gui.Screen).SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		(*gui.Screen).SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		(*gui.Screen).SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		(*gui.Screen).SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		(*gui.Screen).SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}

	gui.drawText(x1+1, y1+1, x2-1, y2-1, style, text)
}

func (gui GUI) Draw() {
	gui.drawBox(1, 1, 42, 7, gui.boxStyle, "Click and drag to draw a box")
	gui.drawBox(5, 9, 32, 14, gui.boxStyle, "Press C to reset")
}

func (gui GUI) Quit() {
	quit := func() {
		maybePanic := recover()
		(*gui.Screen).Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()
}

func (gui GUI) Update() {
	(*gui.Screen).Show()
	ev := (*gui.Screen).PollEvent()

	gui.ProcessEvent(ev)
}

func (gui GUI) ProcessEvent(ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventResize:
		(*gui.Screen).Sync()
	case *tcell.EventKey:
		if gui.HandleEventKey != nil {
			gui.HandleEventKey(ev)
		}
		if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
			os.Exit(0)
			return
		} else if ev.Key() == tcell.KeyCtrlL {
			(*gui.Screen).Sync()
		} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
			(*gui.Screen).Clear()
		}

	case *tcell.EventMouse:
		if gui.HandleEventMouse != nil {
			gui.HandleEventMouse(ev)
		}
	}
}
