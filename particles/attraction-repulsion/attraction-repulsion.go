package main

import (
	"flag"
	"image/color"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

var w, h int

var G float64

func run() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:       "Attraction Repulsion",
		Bounds:      pixel.R(0, 0, float64(w), float64(h)),
		VSync:       true,
		Undecorated: true,
	})
	if err != nil {
		panic(err)
	}

	var (
		last       = time.Now()
		particles  = []*particle{}
		attractors = []pixel.Vec{}
	)

	for !win.Closed() {
		win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

		dt := time.Since(last).Seconds()
		last = time.Now()

		imd := imdraw.New(nil)
		imd.Precision = 6

		for _, p := range particles {
			p.update(dt, attractors)

			if win.Bounds().Contains(p.pos) {
				imd.Color = color.NRGBA{45, 47, 49, 200}
				imd.Push(p.pos)
				imd.Circle(p.mass, 0)

				imd.Color = color.NRGBA{45, 47, 49, 40}
				imd.Push(p.pos)
				imd.Circle(4, 1.2)

				imd.Color = color.NRGBA{0, 0, 0, 16}
				imd.Push(p.pos, p.pos.Add(p.vel.Unit().Scaled(16)))
				imd.Line(3)
			}
		}

		for _, a := range attractors {
			for _, p := range particles {
				if l := a.Sub(p.pos).Len(); l < 112 {
					imd.Color = color.NRGBA{54, 57, 59, 0}
					imd.Push(p.pos)
					imd.Color = color.NRGBA{119, 187, 17, 112 - uint8(l)}
					imd.Push(a)
					imd.Line(3)

					switch {
					case l < 80:
						imd.Push(p.pos)
						imd.Circle(3, 0)
					case l < 96:
						imd.Push(p.pos)
						imd.Circle(2.5, 0)
					case l < 112:
						imd.Push(p.pos)
						imd.Circle(2, 0)
					}
				}
			}
		}

		for _, a := range attractors {
			imd.Color = color.NRGBA{119, 187, 17, 255}
			imd.Push(a)
			imd.Circle(12, 0)

			imd.Color = color.NRGBA{199, 244, 100, 255}
			imd.Push(a)
			imd.Circle(6, 0)
		}

		win.Clear(color.NRGBA{54, 57, 59, 255})
		imd.Draw(win)

		if win.JustPressed(pixelgl.KeyC) {
			particles = []*particle{}
			attractors = []pixel.Vec{}
		}

		if win.JustPressed(pixelgl.MouseButtonRight) {
			attractors = append(attractors, win.MousePosition())
		}

		if win.Pressed(pixelgl.KeyUp) {
			G += 0.01
		}

		if win.Pressed(pixelgl.KeyDown) {
			G -= 0.01
		}

		if win.Pressed(pixelgl.MouseButtonLeft) {
			particles = append(particles, &particle{
				pos:  win.MousePosition(),
				vel:  pixel.V(rand.Float64()-0.5, rand.Float64()-0.5).Scaled(0.15),
				mass: 4 + rand.Float64()*6,
			})
		}

		win.Update()
	}
}

func main() {
	flag.IntVar(&w, "w", 1024, "width")
	flag.IntVar(&h, "h", 576, "width")
	flag.Float64Var(&G, "G", 0.6673, "gravity")
	flag.Parse()

	rand.Seed(42)

	pixelgl.Run(run)
}

type particle struct {
	vel pixel.Vec
	pos pixel.Vec
	acc pixel.Vec

	mass float64
}

func (p *particle) update(dt float64, attractors []pixel.Vec) {
	for _, a := range attractors {
		f := a.Sub(p.pos)
		d := f.Len()

		if d > 128 {
			d = 128
		}

		s := ((G * ((p.mass + 1) * 3)) / (d * d))

		f = f.Unit().Scaled(s)

		if d < 24 {
			p.acc = a.To(p.pos).Unit()
		}

		p.acc = p.acc.Add(f)
	}

	p.pos = p.pos.Add(p.vel)
	p.vel = p.vel.Add(p.acc).Scaled(1 / p.vel.Len() * (p.mass / 12)) // Dampening by speed
	p.acc = p.acc.Scaled(0)
}
