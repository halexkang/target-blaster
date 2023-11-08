package main

import (
	"math/rand"
	"time"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const width = 1024
const height = 768
const growRate = 0.1
const maxSize = 20

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

func addTargets(targets *[]*target) {
	target := target{x: float64(rand.Intn(width)), y: float64(rand.Intn(height)), size: 0, grow: true}
	*targets = append(*targets, &target)
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

	targets := []*target{}
	last := time.Now()
	for !win.Closed() {
		win.Clear(colornames.Aliceblue)
		dt := time.Since(last).Seconds()
		if dt >= 0.5 {
			addTargets(&targets)
			last = time.Now()
		}
		for _, target := range targets {
			target.resize()
			target.drawTarget(win)
		}
		win.Update()
	}
}

func main() {
	pixelgl.Run(run) // gives pixelgl control of main func
}
