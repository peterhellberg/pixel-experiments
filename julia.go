package main

import (
	"image"
	"image/color"
	"math/cmplx"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/peterhellberg/plasma/palette"
)

const (
	width, height = 1280, 720
	maxIterations = 1024
	juliaConstant = complex(-0.7, 0.27015)
	//juliaConstant = complex(0.285, 0.01)
)

var z = 1.0

func main() {
	pixelgl.Run(run)
}

func run() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, float64(width), float64(height)),
		VSync:       true,
		Undecorated: true,
	})
	if err != nil {
		panic(err)
	}

	buffer := image.NewRGBA(image.Rect(0, 0, width, height))
	canvas := pixelgl.NewCanvas(win.Bounds())

	go draw(buffer)

	c := win.Bounds().Center()

	for !win.Closed() {
		win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

		if win.Pressed(pixelgl.KeyUp) {
			z += 0.01
		}

		if win.Pressed(pixelgl.KeyDown) {
			z -= 0.01
		}

		win.Clear(color.Black)

		canvas.SetPixels(buffer.Pix)
		canvas.Draw(win, pixel.IM.Moved(c))

		win.Update()
	}
}

func draw(buffer *image.RGBA) {
	for iterations := 0; iterations < maxIterations; iterations++ {
		i := &imageIterator{bounds: buffer.Bounds()}

		for i.next() {
			center := i.centered()

			value := calculateValue(
				complex(1.5*real(center)/(0.45*width*z), imag(center)/(0.45*height*z)),
				iterations,
			)

			if value >= 18 {
				r, g, b := palette.MaterialDesign500[value%255].RGB255()

				buffer.Set(i.X, i.Y, color.RGBA{r, g, b, 255})
			} else {
				buffer.Set(i.X, i.Y, color.RGBA{0, 0, 0, 255})
			}
		}
	}

}

func calculateValue(value complex128, iterations int) (i int) {
	for i = 0; i < iterations; i++ {
		value = (value * value) + juliaConstant

		if cmplx.Abs(value) > 4.0 {
			return i
		}
	}

	return i
}

type imageIterator struct {
	image.Point
	bounds image.Rectangle
}

func (i *imageIterator) centered() complex128 {
	return complex(
		float64(i.X)-float64(i.bounds.Max.X)/2.0,
		float64(i.Y)-float64(i.bounds.Max.Y)/2.0,
	)
}

func (i *imageIterator) next() bool {
	if i.X < i.bounds.Max.X {
		i.X++

		return true
	}

	if i.Y < i.bounds.Max.Y {
		i.X = 0
		i.Y++

		return true
	}

	return false
}
