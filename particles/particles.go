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
	w, h   = 512, 512
	fw, fh = float64(w), float64(h)
)

func flip() float64 {
	if rand.Float64() > 0.5 {
		return 1.0
	}

	return -1.0
}

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

	last := time.Now()

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		buffer := image.NewRGBA(image.Rect(0, 0, w, h))

		win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

		if win.JustPressed(pixelgl.KeyC) {
			particles = nil
			draw.Draw(buffer, buffer.Bounds(), image.Transparent, image.ZP, draw.Src)
		}

		fx, fy := win.MousePosition().XY()
		x := int(fx)
		y := int(fy)

		if win.Pressed(pixelgl.MouseButtonLeft) {
			c := color.RGBA{255, uint8(y % 255), uint8(x % 255), 255}

			for i := 0; i < 10; i++ {
				angle := 90.0 + rand.Float64()*180.0*flip()
				speed := 40.0 + dt + (20.0 * rand.Float64())
				life := 0.2 + (2.0 * rand.Float64())

				particles = append(particles,
					newParticle(fx, fy, angle, speed, life, c),
				)
			}
		}

		for _, p := range particles {
			p.update(dt)

			x := int(p.position.X)
			y := int(p.position.Y)

			buffer.Set(x, y, p.color)
			buffer.Set(x-1, y, p.color)
			buffer.Set(x+1, y, p.color)
			buffer.Set(x, y-1, p.color)
			buffer.Set(x, y+1, p.color)
		}

		aliveParticles := []*particle{}

		for _, p := range particles {
			if p.life > 0 {
				aliveParticles = append(aliveParticles, p)
			}
		}

		particles = aliveParticles

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
