package main

import "fmt"

func main() {
	e := NewEngine()
	e.SubscribeFunc(Tick, func(m Message) {
		fmt.Println("tick:", m)
	})
	e.SubscribeFunc(Add, func(m Message) {
		fmt.Println("add:", m)
	})
	e.SubscribeFunc(Add|Position, func(m Message) {
		fmt.Println("add position:", m)
	})
	e.SubscribeFunc(Update|Position, func(m Message) {
		fmt.Println("update position:", m)
	})
	e.SubscribeFunc(Update|Position|Velocity, func(m Message) {
		fmt.Println("update position, velocity:", m)
	})

	e.Publish(Message{
		Add | Position | Velocity | Geometry,
		"add, position, velocity, geometry",
	})

	e.Publish(Message{
		Update | Position,
		"update, position",
	})

	e.Publish(Message{
		Update | Position | Velocity | Geometry,
		"update, position, velocity, geometry",
	})

	e.Publish(Message{
		Tick,
		"tick",
	})
}

type Entity uint

type Kind uint16

func (k Kind) String() string {
	return fmt.Sprintf("%016b", k)
}

func (a Kind) Contains(b Kind) bool {
	return b&^a == 0
}

func (a Kind) Intersects(b Kind) bool {
	return b&a > 0
}

const (
	// process
	Tick Kind = 1 << iota
	Quit

	// entity
	Add
	Update
	Remove

	// components
	Position
	Velocity
	Geometry
)

type Message struct {
	Flags   Kind
	Payload interface{}
}

type System interface {
	Handle(Message)
}

type SystemFunc func(Message)

func (f SystemFunc) Handle(m Message) { f(m) }

type Engine struct {
	// sync.RWMutex
	subscribers map[Kind][]System // replace with tree
}

func NewEngine() *Engine {
	return &Engine{
		subscribers: make(map[Kind][]System),
	}
}

var DefaultEngine = NewEngine()

func (e *Engine) Subscribe(k Kind, s System) {
	e.subscribers[k] = append(e.subscribers[k], s)
}

func (e *Engine) SubscribeFunc(k Kind, f func(Message)) {
	e.Subscribe(k, SystemFunc(f))
}

func (e *Engine) Publish(m Message) {
	for f, kinds := range e.subscribers {
		if m.Flags.Contains(f) {
			for _, h := range kinds {
				h.Handle(m)
			}
		}
	}
}
