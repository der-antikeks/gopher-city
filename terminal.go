package main

import (
	"log"

	"github.com/nsf/termbox-go"
)

type Terminal struct {
	running bool
	engine  *Engine
}

func NewTerminal(e *Engine) *Terminal {
	if err := termbox.Init(); err != nil {
		log.Fatal(err)
		return nil
	}
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	t := &Terminal{
		running: true,
		engine:  e,
	}

	go func() {
		for t.running {
			t.HandleEvent(termbox.PollEvent())
		}
	}()

	return t
}

func (t *Terminal) Handle(m Message) {
	if m.Kind(Quit) {
		t.running = false
		termbox.Close()
	}
}

func (t *Terminal) HandleEvent(ev termbox.Event) {
	if !t.running {
		return
	}

	switch ev.Type {
	case termbox.EventKey:
		t.engine.Publish(Message{Key, ev.Key})
	case termbox.EventResize:
		t.engine.Publish(Message{Resize, Point{ev.Width, ev.Height}})
	case termbox.EventMouse:
		t.engine.Publish(Message{Mouse, MouseEvent{ev.Key, ev.MouseX, ev.MouseY}})
	case termbox.EventError:
		//t.engine.Publish(Message{Error, ev.Err})
		log.Fatal(ev.Err)
	}
}

func (t *Terminal) Size() (width, height int) {
	return termbox.Size()
}

type Point struct {
	X, Y int
}

type MouseEvent struct {
	Key  termbox.Key
	X, Y int
}
