package main

import (
	"image"
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	w, h   = 512, 512
	fw, fh = float64(w), float64(h)
)

func run() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, fw, fh),
		Undecorated: true,
	})
	if err != nil {
		panic(err)
	}

	canvas := pixelgl.NewCanvas(win.Bounds())

	for !win.Closed() {
		win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

		drawFrame(win, canvas)

		win.Update()
	}
}

func drawFrame(win *pixelgl.Window, canvas *pixelgl.Canvas) {
	fx, fy := win.MousePosition().XY()
	x, y := int(fx), int(fy)

	buffer := image.NewRGBA(image.Rect(0, 0, w, h))

	buffer.Set(x, y, color.RGBA{255, uint8(y % 255), uint8(x % 255), 255})

	canvas.SetPixels(buffer.Pix)

	canvas.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
}

func main() {
	pixelgl.Run(run)
}
