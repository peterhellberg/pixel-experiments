package main

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	width      = 600
	height     = 600
	seed       = 6502
	numCircles = 4
)

var (
	threshold = 0.6
	pxSize    = 3
)

type circle struct {
	x     int
	y     int
	r     int
	vx    int
	vy    int
	color color.RGBA
}

var circles = []*circle{}

func init() {
	rand.Seed(seed)

	for i := 0; i < numCircles; i++ {
		c := color.RGBA{
			uint8(random(100, 255)),
			uint8(random(20, 255)),
			uint8(random(50, 255)),
			255,
		}

		circles = append(circles, &circle{
			x:     random(width/3, width),
			y:     random(height/3, height),
			r:     random(64, 128),
			vx:    random(-3, 3),
			vy:    random(-3, 3),
			color: c,
		})
	}
}

func drawPicture() *pixel.PictureData {
	m := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.Draw(m, m.Bounds(), &image.Uniform{color.Black}, image.ZP, draw.Src)

	for x := 0; x < width; x += pxSize {
		for y := 0; y < height; y += pxSize {
			sum, closestD2 := 0.0, math.MaxFloat64

			var closestColor color.RGBA

			for _, c := range circles {
				dx := float64(x - c.x)
				dy := float64(y - c.y)
				d2 := dx*dx + dy*dy

				if d2 >= 0 {
					sum += float64(c.r) * float64(c.r) / d2
				}

				if d2 < closestD2 {
					closestD2 = d2
					closestColor = c.color
				}
			}

			if sum > threshold {
				rect := image.Rect(x, y, x+pxSize, y+pxSize)
				draw.Draw(m, rect, &image.Uniform{closestColor}, image.ZP, draw.Src)
			} else {
				m.Set(x, y, color.RGBA{150, 150, 250, 255})
			}
		}
	}

	return pixel.PictureDataFromImage(m)
}

func update() {
	for _, c := range circles {
		c.x += c.vx
		c.y += c.vy

		if c.x-c.r < 0+c.r {
			c.vx = +abs(c.vx)
		}

		if c.x+c.r > width-c.r {
			c.vx = -abs(c.vx)
		}

		if c.y-c.r < 0+c.r {
			c.vy = +abs(c.vy)
		}

		if c.y+c.r > height-c.r {
			c.vy = -abs(c.vy)
		}
	}
}

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

	c := win.Bounds().Center()

	go func() {
		for {
			update()
			time.Sleep(60 * time.Millisecond)
		}
	}()

	for !win.Closed() {
		win.Update()

		p := drawPicture()

		s := pixel.NewSprite(p, p.Bounds())
		s.Draw(win, pixel.IM.Moved(c).Scaled(c, 1))

		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			return
		}

		if win.Pressed(pixelgl.KeyRight) {
			threshold += 0.01
		}

		if win.Pressed(pixelgl.KeyLeft) {
			threshold -= 0.01
		}

		if win.Pressed(pixelgl.KeyUp) {
			pxSize++
		}

		if win.Pressed(pixelgl.KeyDown) {
			if pxSize > 1 {
				pxSize--
			}
		}

		if win.JustPressed(pixelgl.KeyR) {
			for _, circle := range circles {
				circle.color = color.RGBA{
					uint8(random(100, 255)),
					uint8(random(20, 255)),
					uint8(random(50, 255)),
					155,
				}
			}
		}

		if win.JustPressed(pixelgl.KeyS) {
			c := color.RGBA{
				uint8(random(100, 255)),
				uint8(random(20, 255)),
				uint8(random(50, 255)),
				155,
			}

			for _, circle := range circles {
				circle.color = c
			}
		}
	}
}

func main() {
	pixelgl.Run(run)
}

func random(min, max int) int {
	return min - rand.Intn(max-min)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}

	return n
}
