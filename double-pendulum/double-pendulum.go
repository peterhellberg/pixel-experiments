package main

import (
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

var (
	g      = 0.1
	r1     = 180.0
	r2     = 90.0
	m1     = 32.0
	m2     = 16.0
	a1v    = 0.0
	a2v    = 0.0
	a1, a2 = a1a2DefaultValues()
)

func run() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, 600, 310),
		VSync:       true,
		Undecorated: true,
	})
	if err != nil {
		panic(err)
	}

	win.SetSmooth(true)
	win.SetMatrix(pixel.IM.ScaledXY(pixel.ZV, pixel.V(1, -1)).Moved(pixel.V(300, 300)))

	for !win.Closed() {
		win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

		if win.JustPressed(pixelgl.KeySpace) {
			a1, a2 = a1a2DefaultValues()
		}

		win.Clear(color.NRGBA{44, 44, 84, 255})

		a, b := update()

		imd := imdraw.New(nil)

		imd.Color = color.NRGBA{64, 64, 122, 255}
		imd.Push(pixel.ZV, a, b)
		imd.Line(3)

		imd.Color = color.NRGBA{51, 217, 178, 255}
		imd.Push(a)
		imd.Circle(m1/2, 0)

		imd.Color = color.NRGBA{52, 172, 224, 255}
		imd.Push(b)
		imd.Circle(m2/2, 0)

		imd.Draw(win)
		win.Update()
	}
}

func update() (pixel.Vec, pixel.Vec) {
	a1a := a1aCalculation()
	a2a := a2aCalculation()

	a1v += a1a
	a2v += a2a

	a1 += a1v
	a2 += a2v

	a1v *= 0.9996
	a2v *= 0.9996

	a := pixel.V(r1*math.Sin(a1), r1*math.Cos(a1))
	b := pixel.V(a.X+r2*math.Sin(a2), a.Y+r2*math.Cos(a2))

	return a, b
}

func main() {
	pixelgl.Run(run)
}

func a1a2DefaultValues() (float64, float64) {
	return math.Pi / 2, math.Pi / 3
}

func a1aCalculation() float64 {
	num1 := -g * (2*m1 + m2) * math.Sin(a1)
	num2 := -m2 * g * math.Sin(a1-2*a2)
	num3 := -2 * math.Sin(a1-a2) * m2
	num4 := a2v*a2v*r2 + a1v*a1v*r1*math.Cos(a1-a2)
	den := r1 * (2*m1 + m2 - m2*math.Cos(2*a1-2*a2))

	return (num1 + num2 + num3*num4) / den
}

func a2aCalculation() float64 {
	num1 := 2 * math.Sin(a1-a2)
	num2 := (a1v * a1v * r1 * (m1 + m2))
	num3 := g * (m1 + m2) * math.Cos(a1)
	num4 := a2v * a2v * r2 * m2 * math.Cos(a1-a2)
	den := r2 * (2*m1 + m2 - m2*math.Cos(2*a2-2*a2))

	return (num1 * (num2 + num3 + num4)) / den
}
