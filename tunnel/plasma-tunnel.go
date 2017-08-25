package main

import (
	"image"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"

	"github.com/peterhellberg/plasma"
	"github.com/peterhellberg/plasma/palette"
)

const (
	w, h, size    = 640, 480, 128
	fw, fh, fsize = float64(w), float64(h), float64(size)
)

var (
	start         = time.Now()
	texture       = plasmaImage(size, size, 1)
	distanceTable = [h * 2][w * 2]int{}
	angleTable    = [h * 2][w * 2]int{}
	ratio         = 32.0
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

	win.SetSmooth(true)

	canvas := pixelgl.NewCanvas(win.Bounds())

	go func() {
		s := 1
		for range time.Tick(32 * time.Millisecond) {
			s++

			texture = plasmaImage(size, size, s)
		}
	}()

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

	for !win.Closed() {
		win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

		drawFrame(canvas)

		canvas.Draw(win, pixel.IM.Moved(c))

		win.Update()
	}
}

func drawFrame(canvas *pixelgl.Canvas) {
	animation := time.Since(start).Seconds()

	shiftX := int(fsize * 0.2 * animation)
	shiftY := int(fsize * 0.05 * animation)

	shiftLookX := int(fw/2 + float64(int(fw/2*math.Sin(animation))))
	shiftLookY := int(fh/2 + float64(int(fh/2*math.Sin(animation*1.6))))

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

func plasmaImage(w, h, s int) image.Image {
	return plasma.New(w, h, 4).Image(w, h, s, palette.DefaultGradient)
}

func main() {
	pixelgl.Run(run)
}
