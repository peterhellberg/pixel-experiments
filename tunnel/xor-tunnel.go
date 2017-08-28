package main

import (
	"image"
	"image/color"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	w, h, size    = 640, 480, 128
	fw, fh, fsize = float64(w), float64(h), float64(size)
)

var (
	start         = time.Now()
	texture       = xorImage(size, size)
	distanceTable = [h * 2][w * 2]int{}
	angleTable    = [h * 2][w * 2]int{}
	ratio         = 16.0
)

func run() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, float64(w), float64(h)),
		VSync:       true,
		Undecorated: true,
	})
	if err != nil {
		panic(err)
	}

	win.SetSmooth(false)

	canvas := pixelgl.NewCanvas(win.Bounds())

	for y := 0; y < h*2; y++ {
		for x := 0; x < w*2; x++ {
			fx, fy := float64(x), float64(y)

			distance := int(ratio*size/math.Sqrt((fx-fw)*(fx-fw)+(fy-fh)*(fy-fh))) % size
			angle := int(0.5 * size * math.Atan2(fy-fh, fx-fw) / math.Pi)

			distanceTable[y][x] = distance
			angleTable[y][x] = angle
		}
	}

	c := win.Bounds().Center()

	m := pixel.IM.Moved(c)

	for !win.Closed() {
		win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

		drawFrame(canvas)

		canvas.Draw(win, m)

		win.Update()
	}
}

func drawFrame(canvas *pixelgl.Canvas) {
	animation := time.Since(start).Seconds()

	shiftX := int(fsize * 0.16 * animation)
	shiftY := int(0)

	shiftLookX := int(fw/3 + float64(int(fw/3*math.Sin(animation))))
	shiftLookY := int(fh/2 + float64(int(fh/3*math.Sin(animation*2.0))))

	buffer := image.NewRGBA(image.Rect(0, 0, w, h))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			buffer.Set(x, y, texture.At(
				int(uint(distanceTable[y+shiftLookY][x+shiftLookX]+shiftX)%size),
				int(uint(angleTable[y+shiftLookY][x+shiftLookX]+shiftY)%size),
			))
		}
	}

	canvas.SetPixels(buffer.Pix)
}

func xorImage(w, h int) image.Image {
	m := image.NewRGBA(image.Rect(0, 0, w, h))

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			u := uint8(x ^ y)

			c := color.RGBA{
				(u) % 192,
				24,
				(u & uint8(y)) % 128,
				200,
			}

			m.Set(x, y, c)
		}
	}

	return m
}

func main() {
	pixelgl.Run(run)
}
