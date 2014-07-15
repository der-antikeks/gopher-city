package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nsf/termbox-go"
)

func main() {
	if err := termbox.Init(); err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	var (
		update = time.Tick(time.Duration(1000/70) * time.Millisecond)
		event  = make(chan termbox.Event)
		quit   = make(chan struct{})
	)

	go func() {
		for {
			event <- termbox.PollEvent()
		}
	}()

	w, h := termbox.Size()
	resize(w, h)
	console = "initalized"

	for {
		select {
		case <-update:
			simulate()
			draw()

		case ev := <-event:
			switch ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc:
					close(quit)
					return
				case termbox.KeyF1:
					mode = ModeIdle
					console = "idle mode"
				case termbox.KeyF2:
					mode = ModeResidential
					console = "residential mode"
				case termbox.KeyF3:
					mode = ModeCommercial
					console = "commercial mode"
				case termbox.KeyF4:
					mode = ModeIndustrial
					console = "industrial mode"
				case termbox.KeyF5:
					mode = ModeDelete
					console = "delete mode"
				}
			case termbox.EventResize:
				resize(ev.Width, ev.Height)
			case termbox.EventMouse:
				switch ev.Key {
				case termbox.MouseLeft:
					paint(ev.MouseX, ev.MouseY)
				}
			case termbox.EventError:
				log.Fatal(ev.Err)
			}

		case <-quit:
			return
		}
	}
}

type Cell struct { // termbox.Cell
	Ch     rune
	Fg, Bg termbox.Attribute
	Start  time.Time
}

var (
	width, height int
	data          []Cell
)

func resize(w, h int) {
	h -= 1
	if width == w && height == h {
		return
	}

	oldw, oldh := width, height
	olddata := data

	width, height = w, h
	data = make([]Cell, width*height)
	for i := range data {
		data[i].Ch = ' '
		data[i].Fg = termbox.ColorDefault
		data[i].Bg = termbox.ColorDefault
	}

	minw, minh := oldw, oldh
	if w < minw {
		minw = w
	}
	if h < minh {
		minh = h
	}

	for i := 0; i < minh; i++ {
		srco, dsto := i*oldw, i*width
		src := olddata[srco : srco+minw]
		dst := data[dsto : dsto+minw]
		copy(dst, src)
	}

}

var ascii = map[string][]rune{
	"quality": []rune{'.', 'o', 'O'},
	"density": []rune{'░', '▒', '▓', '█'},
	"thin":    []rune{'┌', '─', '┐', '│', '└', '┘'},
	"thick":   []rune{'╔', '═', '╗', '║', '╚', '╝'},
}

type ClickMode int

const (
	ModeIdle ClickMode = iota
	ModeResidential
	ModeCommercial
	ModeIndustrial
	ModeDelete
)

var mode ClickMode

func paint(x, y int) {
	console = fmt.Sprintf("mouse at %v:%v", x, y)

	p := y*width + x
	if p >= len(data) {
		console += " out of bounds"
		return
	}

	color := termbox.ColorDefault // ModeDelete
	switch mode {
	case ModeIdle:
		return
	case ModeResidential:
		color = termbox.ColorGreen
	case ModeCommercial:
		color = termbox.ColorCyan
	case ModeIndustrial:
		color = termbox.ColorYellow
	}

	if data[p].Bg == color {
		return
	}

	if data[p].Bg != termbox.ColorDefault && color != termbox.ColorDefault {
		return
	}

	data[p].Ch = ' '
	data[p].Fg = termbox.ColorDefault
	data[p].Bg = color
	data[p].Start = time.Now()
}

func simulate() {
	for i, c := range data {
		if c.Bg != termbox.ColorDefault {
			delta := time.Since(c.Start)
			switch {
			case 10*time.Second < delta:
				data[i].Ch = ascii["quality"][2]
			case 5*time.Second < delta:
				data[i].Ch = ascii["quality"][1]
			case 1*time.Second < delta:
				data[i].Ch = ascii["quality"][0]
			}
		}
	}
}

var console string

func draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// data
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := data[y*width+x]
			termbox.SetCell(x, y, c.Ch, c.Fg, c.Bg)
		}
	}

	// console
	for p, c := range console {
		termbox.SetCell(p, height, c, termbox.ColorDefault, termbox.ColorDefault)
	}

	// menu
	for y := 1; y < height-1; y++ {
		termbox.SetCell(width-15, y, ascii["thin"][3], termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(width-1, y, ascii["thin"][3], termbox.ColorDefault, termbox.ColorDefault)
	}
	for x := width - 14; x < width-1; x++ {
		termbox.SetCell(x, 0, ascii["thin"][1], termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(x, height-1, ascii["thin"][1], termbox.ColorDefault, termbox.ColorDefault)
	}

	termbox.SetCell(width-15, 0, ascii["thin"][0], termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(width-1, 0, ascii["thin"][2], termbox.ColorDefault, termbox.ColorDefault)

	termbox.SetCell(width-15, height-1, ascii["thin"][4], termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(width-1, height-1, ascii["thin"][5], termbox.ColorDefault, termbox.ColorDefault)

	// flush
	if err := termbox.Flush(); err != nil {
		log.Fatal(err)
	}
}
