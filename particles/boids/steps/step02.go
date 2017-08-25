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

func run() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, fw, fh),
		Undecorated: true,
	})
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())

	canvas := pixelgl.NewCanvas(win.Bounds())

	particles := []*particle{}
	obstacles := []*particle{}

	last := time.Now()

	white := color.RGBA{255, 255, 255, 255}

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		buffer := image.NewRGBA(image.Rect(0, 0, w, h))

		win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

		if win.JustPressed(pixelgl.KeyC) {
			particles = nil
			obstacles = nil

			draw.Draw(buffer, buffer.Bounds(), image.Transparent, image.ZP, draw.Src)
		}

		fx, fy := win.MousePosition().XY()

		if win.JustPressed(pixelgl.KeyO) {
			obstacles = append(obstacles, newParticle(fx, fy, 0, 0, 100, white))
		}

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			particles = append(particles, newParticle(fx, fy, 0, 0, 100,
				color.RGBA{255, uint8(int(fy) % 255), uint8(int(fx) % 255), 255},
			))
		}

		for _, p := range particles {
			p.update(dt)
			buffer.Set(int(p.position.X), int(p.position.Y), p.color)
		}

		for _, o := range obstacles {
			o.update(dt)
			buffer.Set(int(o.position.X), int(o.position.Y), o.color)
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

type particle struct {
	position pixel.Vec
	velocity pixel.Vec
	color    color.RGBA
	life     float64
}

func newParticle(x, y, angle, speed, life float64, c color.RGBA) *particle {
	angleInRadians := angle * math.Pi / 180

	return &particle{
		position: pixel.Vec{x, y},
		velocity: pixel.Vec{
			X: speed * math.Cos(angleInRadians),
			Y: -speed * math.Sin(angleInRadians),
		},
		life:  life,
		color: c,
	}
}

func (p *particle) update(dt float64) {
	p.life -= dt

	if p.life > 0 {
		p.position.X += p.velocity.X * dt
		p.position.Y += p.velocity.Y * dt

		if p.life > 1.5 {
			if rand.Float64() < 0.4+dt {
				p.color.R -= 1
				p.color.G -= 1
				p.color.B -= 1
			}
		} else {
			p.color.R -= 1
			p.color.G -= 1
			p.color.B -= 1
		}
	}
}
