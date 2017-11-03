package main

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

var (
	w, h, scale = float64(640), float64(360), float64(6)

	p = newPalette(PinkColors)

	bg = color.RGBA{7, 9, 9, 255}

	lc, pc color.RGBA

	v = rand.Float64() * scale

	balls = []*ball{
		newRandomBall(4 + v),
		newRandomBall(6 + v),
		newRandomBall(10 + v),
		newRandomBall(14 + v),
		newRandomBall(18 + v),
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

	imd.Precision = 3

	go func() {
		var step int

		for range time.Tick(256 * time.Millisecond) {
			switch imd.Precision {
			case 3:
				step = 1
			case 9:
				step = -1
			}

			imd.Precision += step
		}
	}()

	for !win.Closed() {
		win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

		var positions = []pixel.Vec{}

		for _, ball := range balls {
			positions = append(positions, ball.pos)
		}

		imd.Clear()

		imd.Color = lc
		imd.Push(positions...)
		imd.Push(positions[0])
		imd.Line(scale)

		imd.Color = pc
		imd.Push(positions...)
		imd.Polygon(0)

		for _, ball := range balls {
			imd.Color = ball.color
			imd.Push(ball.pos)
			imd.Circle(ball.radius, 0)
		}

		win.Clear(bg)

		imd.Draw(win)

		win.Update()
	}
}

func main() {
	rand.Seed(42)

	go func() {
		for range time.Tick(16 * time.Millisecond) {
			for _, ball := range balls {
				ball.update()
			}
		}
	}()

	updatePolygonColor(p.next())

	pixelgl.Run(run)
}

func updatePolygonColor(c color.RGBA) {
	lc = color.RGBA{c.R / 4, c.G / 4, c.B / 4, 192}
	pc = color.RGBA{c.R / 6, c.G / 6, c.B / 6, 64}
}

func newRandomBall(radius float64) *ball {
	pos := pixel.V(w/4+(w/4*3)*rand.Float64(), h/4+(h/4*3)*rand.Float64())
	dir := pixel.V((rand.Float64()*2)-1, (rand.Float64()*2)-1).Scaled(32 / radius)

	return &ball{pos, dir, radius, p.color(), p.clone()}
}

type ball struct {
	pos     pixel.Vec
	dir     pixel.Vec
	radius  float64
	color   color.RGBA
	palette *palette
}

func (b *ball) update() {
	b.pos.X += b.dir.X
	b.pos.Y += b.dir.Y

	if b.pos.Y <= b.radius || b.pos.Y >= h-b.radius {
		b.dir.Y *= -1
		b.color = b.palette.next()
		updatePolygonColor(p.next())
	}

	if b.pos.X <= b.radius || b.pos.X >= w-b.radius {
		b.dir.X *= -1
		b.color = b.palette.next()
		updatePolygonColor(p.next())
	}
}

func newPalette(colors []color.RGBA) *palette {
	return &palette{colors, len(colors), 0}
}

type palette struct {
	colors []color.RGBA
	size   int
	index  int
}

func (p *palette) clone() *palette {
	return &palette{p.colors, len(p.colors), p.index}
}

func (p *palette) next() color.RGBA {
	p.index++

	if p.index+1 >= p.size {
		p.index = 0
	}

	return p.colors[p.index]
}

func (p *palette) color() color.RGBA {
	return p.colors[p.index]
}

func (p *palette) random() color.RGBA {
	p.index = rand.Intn(p.size)

	return p.colors[p.index]
}

var (
	Pink            = color.RGBA{255, 192, 203, 255}
	LightPink       = color.RGBA{255, 182, 193, 255}
	HotPink         = color.RGBA{255, 105, 180, 255}
	DeepPink        = color.RGBA{255, 20, 147, 255}
	PaleVioletRed   = color.RGBA{219, 112, 147, 255}
	MediumVioletRed = color.RGBA{199, 21, 133, 255}
)

var PinkColors = []color.RGBA{
	Pink,
	LightPink,
	HotPink,
	DeepPink,
	PaleVioletRed,
	MediumVioletRed,
}
