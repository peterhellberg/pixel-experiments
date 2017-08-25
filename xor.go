package main

import (
	"image"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	width  = 768
	height = 768
	size   = 256
)

func run() {
	scale := float64(height) / float64(size)

	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, float64(width), float64(height)),
		VSync:       true,
		Undecorated: true,
		Resizable:   false,
	})
	if err != nil {
		panic(err)
	}

	win.SetSmooth(false)

	c := win.Bounds().Center()
	p := xorPicture(size, size)
	s := pixel.NewSprite(p, p.Bounds())

	for !win.Closed() {
		win.Update()

		s.Draw(win, pixel.IM.Moved(c).Scaled(c, scale))

		if win.Pressed(pixelgl.KeyUp) {
			scale += 0.1
		}

		if win.Pressed(pixelgl.KeyDown) {
			scale -= 0.1
		}

		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			return
		}
	}
}

func xorPicture(w, h int) *pixel.PictureData {
	m := image.NewRGBA(image.Rect(0, 0, w, h))

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			c := uint8(x ^ y)

			m.Set(x, y, color.RGBA{c, c % 192, c, 255})
		}
	}

	return pixel.PictureDataFromImage(m)
}

func main() {
	pixelgl.Run(run)
}
