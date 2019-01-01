package main

import (
	"image"
	"image/color"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	width  = 512
	height = 512
	size   = width * (height + 2) * 2
	delay  = 100 * time.Millisecond

	halfWidth  = width >> 1
	halfHeight = height >> 1
)

var (
	oldIdx    = width
	newIdx    = width * (height + 3)
	rippleRad = 3

	texture = xorImage(width, height)
	ripple  = xorImage(width, height)

	rippleMap = make([]int, size)
	lastMap   = make([]int, size)
)

func run() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, float64(width), float64(height)),
		VSync:       true,
		Undecorated: true,
		Resizable:   false,
	})
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(delay)

	go func() {
		for range ticker.C {
			if rand.Intn(2) == 1 {
				dropAt(rand.Intn(width), rand.Intn(height))
			}
		}
	}()

	c := win.Bounds().Center()

	for !win.Closed() {
		newRippleFrame()

		p := pixel.PictureDataFromImage(ripple)
		s := pixel.NewSprite(p, p.Bounds())

		s.Draw(win, pixel.IM.Moved(c))

		mouse := win.MousePosition()

		dropAt(int(mouse.X), int(height-mouse.Y))

		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			return
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

func dropAt(dx, dy int) {
	if dx > 1 && dy > 1 && dx < width && dy < height-rippleRad {
		for j := dy - rippleRad; j < dy+rippleRad; j++ {
			for k := dx - rippleRad; k < dx+rippleRad; k++ {
				rippleMap[oldIdx+(j*width)+k] += 512
			}
		}
	}
}

func newRippleFrame() {
	i := oldIdx

	oldIdx = newIdx
	newIdx = i

	i = 0
	mapIdx := oldIdx

	var data, oldData int

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			data = (rippleMap[mapIdx-width] +
				rippleMap[mapIdx+width] +
				rippleMap[mapIdx-1] +
				rippleMap[mapIdx+1]) >> 1

			data -= rippleMap[newIdx+i]
			data -= data >> 5

			rippleMap[newIdx+i] = data

			data = 1024 - data

			oldData = lastMap[i]
			lastMap[i] = data

			if oldData != data {
				a := (((x - halfWidth) * data / 1024) << 0) + halfWidth
				b := (((y - halfHeight) * data / 1024) << 0) + halfHeight

				if a >= width {
					a = width - 1
				}

				if a < 0 {
					a = 0
				}

				if b >= height {
					b = height - 1
				}

				if b < 0 {
					b = 0
				}

				newPixel := (a + (b * width)) * 4
				curPixel := i * 4

				ripple.Pix[curPixel] = texture.Pix[newPixel]
				ripple.Pix[curPixel+1] = texture.Pix[newPixel+1]
				ripple.Pix[curPixel+2] = texture.Pix[newPixel+2]
			}

			mapIdx++
			i++
		}
	}
}

func xorImage(w, h int) *image.RGBA {
	m := image.NewRGBA(image.Rect(0, 0, w, h))

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			c := uint8(x ^ y)

			m.Set(x, y, color.RGBA{c, c % 192, c, 255})
		}
	}

	return m
}
