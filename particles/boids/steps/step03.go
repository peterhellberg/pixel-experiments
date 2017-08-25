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
	w, h   = 768, 432
	fw, fh = float64(w), float64(h)
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func run() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, fw, fh),
		Undecorated: true,
	})
	if err != nil {
		panic(err)
	}

	canvas := pixelgl.NewCanvas(win.Bounds())

	boids := []*boid{}
	avoids := []*avoid{}

	last := time.Now()

	white := color.RGBA{255, 255, 255, 255}

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		buffer := image.NewRGBA(image.Rect(0, 0, w, h))

		win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

		if win.JustPressed(pixelgl.KeyC) {
			boids = nil
			avoids = nil

			draw.Draw(buffer, buffer.Bounds(), image.Transparent, image.ZP, draw.Src)
		}

		fx, fy := win.MousePosition().XY()

		if win.JustPressed(pixelgl.KeyO) {
			avoids = append(avoids, newAvoid(fx, fy, 10, white))
		}

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			angle := 90.0 + rand.Float64()*180.0*flip()
			speed := 40.0 + dt + (20.0 * rand.Float64())

			boids = append(boids, newBoid(fx, fy, angle, speed,
				color.RGBA{255, uint8(int(fy) % 255), uint8(int(fx) % 255), 255},
			))
		}

		for _, p := range boids {
			p.update(dt)
			buffer.Set(int(p.position.X), int(p.position.Y), p.color)
		}

		for _, a := range avoids {
			buffer.Set(int(a.position.X), int(a.position.Y), a.color)
		}

		win.Clear(color.RGBA{0, 0, 0, 255})

		canvas.SetPixels(buffer.Pix)

		canvas.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

func newAvoid(x, y, s float64, c color.RGBA) *avoid {
	p := pixel.Vec{x, y}

	return &avoid{position: p, size: s, color: c}
}

type avoid struct {
	position pixel.Vec
	size     float64
	color    color.RGBA
}

type boid struct {
	position pixel.Vec
	velocity pixel.Vec
	color    color.RGBA
}

func newBoid(x, y, angle, speed float64, c color.RGBA) *boid {
	angleInRadians := angle * math.Pi / 180

	return &boid{
		position: pixel.Vec{x, y},
		velocity: pixel.Vec{
			X: speed * math.Cos(angleInRadians),
			Y: -speed * math.Sin(angleInRadians),
		},
		color: c,
	}
}

func (p *boid) update(dt float64) {
	p.position.X += p.velocity.X * dt
	p.position.Y += p.velocity.Y * dt

	if rand.Float64() < 0.4+dt {
		p.color.R -= 1
		p.color.G -= 1
		p.color.B -= 1
	}
}

func flip() float64 {
	if rand.Float64() > 0.5 {
		return 1.0
	}

	return -1.0
}
