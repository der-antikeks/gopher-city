package main

type Message struct {
	Flags   Kind
	Payload interface{}
}

func (m Message) Kind(f Kind) bool {
	return m.Flags.Contains(f)
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
