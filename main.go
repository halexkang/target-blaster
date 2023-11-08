package main

import (
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const width = 1024
const height = 768
const growRate = 1
const maxSize = 200

type target struct {
	x    float64
	y    float64
	size float64
	grow bool
}

func (t *target) resize() {
	if t.size+growRate >= maxSize {
		t.grow = false
	}
	if t.size-growRate <= 0 {
		t.grow = true
	}
	if t.grow {
		t.size += growRate
	} else {
		t.size -= growRate
	}
}

func (t target) drawTarget(win *pixelgl.Window) {
	imd := imdraw.New(nil)
	imd.Color = colornames.Limegreen
	imd.Push(pixel.V(t.x, t.y))
	imd.Circle(t.size, 0)
	imd.Draw(win)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Target Blaster",
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	target := target{x: 300, y: 500, size: 100, grow: true}

	for !win.Closed() {
		win.Clear(colornames.Aliceblue)
		target.resize()
		target.drawTarget(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run) // gives pixelgl control of main func
}
