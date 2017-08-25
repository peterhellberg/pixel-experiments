package main

import (
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/peterhellberg/plasma"
	"github.com/peterhellberg/plasma/palette"
)

const (
	width  = 768
	height = 768
	size   = 256
)

func run() {
	scale := float64(height) / float64(size)

	cfg := pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, 1024, 512),
		VSync:       true,
		Resizable:   false,
		Undecorated: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.SetSmooth(false)

	var s *pixel.Sprite

	size := 12.0

	p := plasmaPicture(256, 128, size, 0)
	s = pixel.NewSprite(p, p.Bounds())

	go func() {
		c := time.Tick(32 * time.Millisecond)

		var i int

		for range c {
			i++

			p := plasmaPicture(256, 128, size, i)

			s.Set(p, p.Bounds())
		}
	}()

	win.Clear(color.Black)

	c := win.Bounds().Center()

	for !win.Closed() {
		win.Update()

		s.Draw(win, pixel.IM.Moved(c).Scaled(c, scale))

		if win.Pressed(pixelgl.KeyUp) {
			size += 0.2
		}

		if win.Pressed(pixelgl.KeyDown) {
			size -= 0.2
		}

		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			return
		}
	}
}

func plasmaPicture(w, h int, s float64, i int) *pixel.PictureData {
	return pixel.PictureDataFromImage(plasma.New(w, h, s).
		Image(w, h, i, palette.DefaultGradient))
}

func main() {
	pixelgl.Run(run)
}
