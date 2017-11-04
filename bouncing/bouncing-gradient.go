package main

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

var (
	w, h, scale = float64(640), float64(360), float64(48)

	p, bg = newPalette(Colors), color.RGBA{255, 228, 225, 255}

	balls = []*ball{
		newRandomBall(scale / 2),
		newRandomBall(scale / 2),
		newRandomBall(scale / 2),
		newRandomBall(scale / 2),
		newRandomBall(scale / 2),
	}
)

func run() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, w, h),
		VSync:       true,
		Undecorated: true,
	})
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)
	imd.EndShape = imdraw.RoundEndShape

	for !win.Closed() {
		win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

		if win.JustPressed(pixelgl.KeySpace) {
			for _, ball := range balls {
				ball.color = ball.palette.next()
			}
		}

		if win.JustPressed(pixelgl.KeyEnter) {
			for _, ball := range balls {
				ball.pos = center()
				ball.vel = randomVelocity()
			}
		}

		imd.Clear()

		for _, ball := range balls {
			imd.Color = ball.color
			imd.Push(ball.pos)
		}

		imd.Polygon(scale)

		//for _, ball := range balls {
		//	imd.Color = ball.color
		//	imd.Push(ball.pos)
		//	imd.Circle(ball.radius, 0)
		//}

		win.Clear(bg)
		imd.Draw(win)
		win.Update()
	}
}

func main() {
	rand.Seed(4)

	go func() {
		for range time.Tick(32 * time.Millisecond) {
			for _, ball := range balls {
				ball.update()
			}
		}
	}()

	pixelgl.Run(run)
}

func newRandomBall(radius float64) *ball {
	return &ball{
		center(), randomVelocity(),
		math.Pi * (radius * radius),
		radius, p.random(), p,
	}
}

func center() pixel.Vec {
	return pixel.V(w/2, h/2)
}

func randomVelocity() pixel.Vec {
	return pixel.V((rand.Float64()*2)-1, (rand.Float64()*2)-1)
}

type ball struct {
	pos     pixel.Vec
	vel     pixel.Vec
	mass    float64
	radius  float64
	color   color.RGBA
	palette *Palette
}

func (b *ball) update() {
	b.pos = b.pos.Add(b.vel)

	var bounced bool

	switch {
	case b.pos.Y <= b.radius || b.pos.Y >= h-(b.radius):
		b.vel.Y *= -1.0
		bounced = true

		if b.pos.Y < b.radius {
			b.pos.Y = b.radius
		}
	case b.pos.X <= b.radius || b.pos.X >= w-(b.radius):
		b.vel.X *= -1.0
		bounced = true

		if b.pos.X < b.radius {
			b.pos.X = b.radius
		}
	}

	for _, a := range balls {
		if a != b {
			d := a.pos.Sub(b.pos)

			if d.Len() > a.radius+b.radius {
				continue
			}

			pen := d.Unit().Scaled(a.radius + b.radius - d.Len())

			a.pos = a.pos.Add(pen.Scaled(b.mass / (a.mass + b.mass)))
			b.pos = b.pos.Sub(pen.Scaled(a.mass / (a.mass + b.mass)))

			u := d.Unit()
			v := 2 * (a.vel.Dot(u) - b.vel.Dot(u)) / (a.mass + b.mass)

			a.vel = a.vel.Sub(u.Scaled(v * a.mass))
			b.vel = b.vel.Add(u.Scaled(v * b.mass))

			//b.color = p.next()
			//a.color = p.next()
		}
	}

	if bounced {
		b.color = p.next()
	}
}

func newPalette(cc []color.Color) *Palette {
	colors := []color.RGBA{}

	for _, v := range cc {
		if c, ok := v.(color.RGBA); ok {
			colors = append(colors, c)
		}
	}

	return &Palette{colors, len(colors), 0}
}

type Palette struct {
	colors []color.RGBA
	size   int
	index  int
}

func (p *Palette) clone() *Palette {
	return &Palette{p.colors, p.size, p.index}
}

func (p *Palette) next() color.RGBA {
	if p.index++; p.index >= p.size {
		p.index = 0
	}

	return p.colors[p.index]
}

func (p *Palette) color() color.RGBA {
	return p.colors[p.index]
}

func (p *Palette) random() color.RGBA {
	p.index = rand.Intn(p.size)

	return p.colors[p.index]
}

var Colors = []color.Color{
	color.RGBA{190, 38, 51, 255},
	color.RGBA{224, 111, 139, 255},
	color.RGBA{73, 60, 43, 255},
	color.RGBA{164, 100, 34, 255},
	color.RGBA{235, 137, 49, 255},
	color.RGBA{247, 226, 107, 255},
	color.RGBA{47, 72, 78, 255},
	color.RGBA{68, 137, 26, 255},
	color.RGBA{163, 206, 39, 255},
	color.RGBA{0, 87, 132, 255},
	color.RGBA{49, 162, 242, 255},
	color.RGBA{178, 220, 239, 255},
}
