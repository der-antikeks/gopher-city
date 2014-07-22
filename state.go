package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nsf/termbox-go"
)

type ClickMode int

const (
	ModeIdle ClickMode = iota
	ModeResidential
	ModeCommercial
	ModeIndustrial
	ModeDelete
)

type GameState struct {
	running  bool
	engine   *Engine
	terminal *Terminal

	mode    ClickMode
	console string

	width, height int
	data          []Cell
}

func NewGameState(e *Engine, t *Terminal) *GameState {
	gs := &GameState{
		running:  true,
		engine:   e,
		terminal: t,

		console: "initalized",
	}

	w, h := t.Size()
	gs.resize(w, h)

	return gs
}

func (gs *GameState) Handle(m Message) {
	if !gs.running {
		return
	}

	switch {
	case m.Kind(Key):
		switch m.Payload.(termbox.Key) {
		case termbox.KeyEsc:
			gs.engine.Publish(Message{Flags: Quit})
		case termbox.KeyF1:
			gs.mode = ModeIdle
			gs.console = "idle mode"
		case termbox.KeyF2:
			gs.mode = ModeResidential
			gs.console = "residential mode"
		case termbox.KeyF3:
			gs.mode = ModeCommercial
			gs.console = "commercial mode"
		case termbox.KeyF4:
			gs.mode = ModeIndustrial
			gs.console = "industrial mode"
		case termbox.KeyF5:
			gs.mode = ModeDelete
			gs.console = "delete mode"
		}

	case m.Kind(Resize):
		re := m.Payload.(Point)
		gs.resize(re.X, re.Y)

	case m.Kind(Mouse):
		me := m.Payload.(MouseEvent)
		switch me.Key {
		case termbox.MouseLeft:
			gs.paint(me.X, me.Y)
		}

	case m.Kind(Tick):
		//now := m.Payload.(time.Time)
		gs.simulate()
		gs.draw()

	case m.Kind(Quit):
		gs.running = false

	}
}

func (gs *GameState) Running() bool {
	return gs.running
}

type Cell struct { // termbox.Cell
	Ch     rune
	Fg, Bg termbox.Attribute
	Start  time.Time
}

func (gs *GameState) resize(w, h int) {
	h -= 1
	if gs.width == w && gs.height == h {
		return
	}

	oldw, oldh := gs.width, gs.height
	olddata := gs.data

	gs.width, gs.height = w, h
	gs.data = make([]Cell, gs.width*gs.height)
	for i := range gs.data {
		gs.data[i].Ch = ' '
		gs.data[i].Fg = termbox.ColorDefault
		gs.data[i].Bg = termbox.ColorDefault
	}

	minw, minh := oldw, oldh
	if w < minw {
		minw = w
	}
	if h < minh {
		minh = h
	}

	for i := 0; i < minh; i++ {
		srco, dsto := i*oldw, i*gs.width
		src := olddata[srco : srco+minw]
		dst := gs.data[dsto : dsto+minw]
		copy(dst, src)
	}

}

var ascii = map[string][]rune{
	"quality": []rune{'.', 'o', 'O'},
	"density": []rune{'░', '▒', '▓', '█'},
	"thin":    []rune{'┌', '─', '┐', '│', '└', '┘'},
	"thick":   []rune{'╔', '═', '╗', '║', '╚', '╝'},
}

func (gs *GameState) paint(x, y int) {
	gs.console = fmt.Sprintf("mouse at %v:%v", x, y)

	p := y*gs.width + x
	if p >= len(gs.data) {
		gs.console += " out of bounds"
		return
	}

	color := termbox.ColorDefault // ModeDelete
	switch gs.mode {
	case ModeIdle:
		return
	case ModeResidential:
		color = termbox.ColorGreen
	case ModeCommercial:
		color = termbox.ColorCyan
	case ModeIndustrial:
		color = termbox.ColorYellow
	}

	if gs.data[p].Bg == color {
		return
	}

	if gs.data[p].Bg != termbox.ColorDefault && color != termbox.ColorDefault {
		return
	}

	gs.data[p].Ch = ' '
	gs.data[p].Fg = termbox.ColorDefault
	gs.data[p].Bg = color
	gs.data[p].Start = time.Now()
}

func (gs *GameState) simulate() {
	for i, c := range gs.data {
		if c.Bg != termbox.ColorDefault {
			delta := time.Since(c.Start)
			switch {
			case 10*time.Second < delta:
				gs.data[i].Ch = ascii["quality"][2]
			case 5*time.Second < delta:
				gs.data[i].Ch = ascii["quality"][1]
			case 1*time.Second < delta:
				gs.data[i].Ch = ascii["quality"][0]
			}
		}
	}
}

func (gs *GameState) draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// data
	for y := 0; y < gs.height; y++ {
		for x := 0; x < gs.width; x++ {
			c := gs.data[y*gs.width+x]
			termbox.SetCell(x, y, c.Ch, c.Fg, c.Bg)
		}
	}

	// console
	for p, c := range gs.console {
		termbox.SetCell(p, gs.height, c, termbox.ColorDefault, termbox.ColorDefault)
	}

	// menu
	for y := 1; y < gs.height-1; y++ {
		termbox.SetCell(gs.width-15, y, ascii["thin"][3], termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(gs.width-1, y, ascii["thin"][3], termbox.ColorDefault, termbox.ColorDefault)
	}
	for x := gs.width - 14; x < gs.width-1; x++ {
		termbox.SetCell(x, 0, ascii["thin"][1], termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(x, gs.height-1, ascii["thin"][1], termbox.ColorDefault, termbox.ColorDefault)
	}

	termbox.SetCell(gs.width-15, 0, ascii["thin"][0], termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(gs.width-1, 0, ascii["thin"][2], termbox.ColorDefault, termbox.ColorDefault)

	termbox.SetCell(gs.width-15, gs.height-1, ascii["thin"][4], termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(gs.width-1, gs.height-1, ascii["thin"][5], termbox.ColorDefault, termbox.ColorDefault)

	// flush
	if err := termbox.Flush(); err != nil {
		log.Fatal(err)
	}
}
