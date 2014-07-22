package main

type Kind uint32

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
	Error

	// input
	Resize
	Key
	Mouse

	// entity
	Add
	Update
	Remove

	// components
	Position
	Velocity
	Geometry
)
