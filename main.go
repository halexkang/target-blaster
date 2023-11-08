package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"
	"github.com/gopxl/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

const width = 1024
const height = 768
const growRate = 0.1
const maxSize = 20
const lives = 50

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

func (t target) drawTargetHelper(win *pixelgl.Window, scale float64, color color.RGBA) {
	imd := imdraw.New(nil)
	imd.Color = color
	imd.Push(pixel.V(t.x, t.y))
	imd.Circle(t.size*scale, 0)
	imd.Draw(win)
}

func (t target) drawTarget(win *pixelgl.Window) {
	t.drawTargetHelper(win, 1, colornames.Orange)
	t.drawTargetHelper(win, 0.8, colornames.White)
	t.drawTargetHelper(win, 0.6, colornames.Orange)
	t.drawTargetHelper(win, 0.4, colornames.White)
}

func (t target) collide(x, y float64) bool {
	r := math.Sqrt(math.Pow(t.x-x, 2) + math.Pow(t.y-y, 2))
	return r <= t.size
}

func menu(win *pixelgl.Window, hits int, misses int, lives int) {
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	txt := text.New(pixel.V(20, height-30), atlas)
	fmt.Fprintln(txt, "Hits: ", hits, "Misses: ", misses, "Lives: ", lives)
	txt.Draw(win, pixel.IM.Scaled(txt.Orig, 1.5))
}

func addTargets(targets *[]*target) {
	randX := float64(rand.Intn(width-maxSize*4) + maxSize*2)
	randY := float64(rand.Intn(height-maxSize*4) + maxSize*2)
	target := target{x: randX, y: randY, size: 0, grow: true}
	*targets = append(*targets, &target)
}

func run() {
	// create window
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
	clicked := false
	mousePos := win.MousePosition()
	misses := 0
	hits := 0
	gameOver := false
	for !win.Closed() {
		if !gameOver {
			if win.JustPressed(pixelgl.MouseButtonLeft) || win.JustPressed(pixelgl.MouseButtonRight) {
				clicked = true
				mousePos = win.MousePosition()
			}
			dt := time.Since(last).Seconds()
			if dt >= 1 { // time between creating new targets
				last = time.Now()
				addTargets(&targets)
			}
			win.Clear(colornames.Gray)
			menu(win, hits, misses, lives-misses)
			for i := len(targets) - 1; i >= 0; i-- {
				targets[i].resize()
				if targets[i].size <= 0 { // if target missed
					targets = append(targets[:i], targets[i+1:]...)
					misses += 1
				}
				if clicked && targets[i].collide(mousePos.X, mousePos.Y) { // if target hit
					targets = append(targets[:i], targets[i+1:]...)
					hits += 1
				}
				targets[i].drawTarget(win)
			}
			if misses > lives {
				gameOver = true
			}
			win.Update()
		} else {
			break
		}
	}

}

func main() {
	pixelgl.Run(run) // gives pixelgl control of main func
}
