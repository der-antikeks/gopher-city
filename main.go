package main

import (
	"fmt"
	"log"

	"github.com/nsf/termbox-go"
)

func main() {
	if err := termbox.Init(); err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	width, height := termbox.Size()
	data := [][]termbox.Cell{}

	print := func(v string) {
		tmp := make([]termbox.Cell, len(v))
		for pos, ch := range v {
			tmp[pos] = termbox.Cell{Ch: ch, Fg: termbox.ColorDefault, Bg: termbox.ColorDefault}
		}

		for len(tmp) > width {
			data = append(data, tmp[0:width])
			tmp = tmp[width:len(tmp)]
		}
		data = append(data, tmp)

		if o := len(data) - height; o > 0 {
			data = data[o:len(data)]
		}
	}

	print("init")
	print("Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.")
	draw(data)

	for {
		ev := termbox.PollEvent()
		switch ev.Type {

		case termbox.EventKey:
			switch ev.Key {
			default:
				print(fmt.Sprint("key: ", ev.Key, ev.Ch))
			case termbox.KeyEsc:
				return
			}
		case termbox.EventResize:
			width, height = ev.Width, ev.Height
			print(fmt.Sprint("resize to: ", width, height))
		case termbox.EventMouse:
			print(fmt.Sprint("mouse: ", ev.Key, ev.MouseX, ev.MouseY))
		case termbox.EventError:
			print(fmt.Sprint("error: ", ev.Err))
		default:
			print(fmt.Sprint("unknown event: ", ev))
		}
		draw(data)
	}
}

func draw(data [][]termbox.Cell) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for y, xc := range data {
		for x, c := range xc {
			termbox.SetCell(x, y, c.Ch, c.Fg, c.Bg)
		}
	}

	if err := termbox.Flush(); err != nil {
		log.Fatal(err)
	}
}
