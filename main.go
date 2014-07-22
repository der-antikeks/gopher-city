package main

import (
	"time"
)

func main() {
	engine := NewEngine()

	terminal := NewTerminal(engine)
	engine.Subscribe(Quit, terminal)

	state := NewGameState(engine, terminal)
	engine.Subscribe(Key, state)
	engine.Subscribe(Resize, state)
	engine.Subscribe(Mouse, state)
	engine.Subscribe(Tick, state)
	engine.Subscribe(Quit, state)

	var (
		update = time.Tick(time.Duration(1000/70) * time.Millisecond)
		now    time.Time
	)

	for state.Running() {
		select {
		case now = <-update:
			engine.Publish(Message{Tick, now})
		}
	}
}
