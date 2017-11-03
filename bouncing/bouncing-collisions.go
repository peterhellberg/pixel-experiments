package main

import (
	"image/color"
	"image/color/palette"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

var (
	w, h, scale = float64(640), float64(360), float64(6)

	p, bg = newPalette(palette.WebSafe[128:192]), color.RGBA{255, 228, 225, 255}

	balls = []*ball{
		newRandomBall(30),
		newRandomBall(40),
		newRandomBall(30),
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

	win.SetSmooth(true)

	imd := imdraw.New(nil)
	imd.EndShape = imdraw.RoundEndShape

	for !win.Closed() {
		win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

		c0 := balls[0].color
		c1 := balls[1].color
		c2 := balls[2].color

		lc := color.RGBA{(c1.R / 3) * 2, (c1.G / 3) * 2, (c1.B / 3) * 2, 255}

		imd.Clear()

		imd.Color = lc
		imd.Push(balls[0].pos, balls[1].pos, balls[2].pos)
		imd.Line(18)

		imd.Push(balls[1].pos)
		imd.Circle(balls[1].radius, 24)

		imd.Color = color.RGBA{
			uint8((int(c0.R) + int(c2.R) + int(c1.R)) / 4),
			uint8((int(c0.G) + int(c2.G) + int(c1.G)) / 4),
			uint8((int(c0.B) + int(c2.B) + int(c1.B)) / 4),
			128,
		}

		imd.Push(balls[0].pos, balls[2].pos)
		imd.Line(balls[0].radius * 2.5)

		for _, ball := range []*ball{balls[0], balls[2]} {
			imd.Color = ball.color
			imd.Push(ball.pos)
			imd.Circle(ball.radius, 0)
		}

		imd.Color = balls[1].color
		imd.Push(balls[1].pos)
		imd.Circle(balls[1].radius, 0)

		win.Clear(bg)

		imd.Draw(win)

		win.Update()
	}
}

func main() {
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
		pixel.V(w/2, h/2),
		pixel.V((rand.Float64()*2)-1, (rand.Float64()*2)-1).Scaled(scale / 2),
		math.Pi * (radius * radius),
		radius, p.next(), p.clone(),
	}
}

type ball struct {
	pos     pixel.Vec
	dir     pixel.Vec
	mass    float64
	radius  float64
	color   color.RGBA
	palette *Palette
}

func (b *ball) update() {
	b.pos.X += b.dir.X
	b.pos.Y += b.dir.Y

	if b.pos.Y <= b.radius+6 || b.pos.Y >= h-(b.radius+6) {
		b.dir.Y *= -1.0
		b.color = b.palette.next()
	}

	if b.pos.X <= b.radius+6 || b.pos.X >= w-(b.radius+6) {
		b.dir.X *= -1.0
		b.color = b.palette.next()
	}

	for _, a := range balls {
		if b != a {
			d := a.pos.Sub(b.pos)

			if d.Len() > a.radius+b.radius {
				continue
			}

			pen := d.Unit().Scaled(a.radius + b.radius - d.Len())

			a.pos = a.pos.Add(pen.Scaled(b.mass / (a.mass + b.mass)))
			b.pos = b.pos.Sub(pen.Scaled(a.mass / (a.mass + b.mass)))

			u := d.Unit()
			v := 2 * (a.dir.Dot(u) - b.dir.Dot(u)) / (a.mass + b.mass)

			a.dir = a.dir.Sub(u.Scaled(v * a.mass))
			b.dir = b.dir.Add(u.Scaled(v * b.mass))

			a.color = a.palette.next()
			b.color = b.palette.next()
		}
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
	return &Palette{p.colors, len(p.colors), p.index}
}

func (p *Palette) next() color.RGBA {
	p.index++

	if p.index+1 >= p.size {
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
