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
	w, h   = 512, 256
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

	last := time.Now()

	var total float64

	angle := 90.0

	white := color.RGBA{255, 255, 255, 255}

	fall := func(fx, fy float64) {
		speed := 32.0 + (32.0 * rand.Float64())
		life := 16.0 + (1.6 * rand.Float64())

		particles = append(particles,
			newParticle(fx, fy, angle, speed, life, white),
		)
	}

	go func() {
		for {
			time.Sleep(40 * time.Millisecond)
			fall(rand.Float64()*fw, fh-(rand.Float64()*fh)/2+(fh/2))
		}
	}()

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		total += dt

		buffer := image.NewRGBA(image.Rect(0, 0, w, h))

		win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

		if win.JustPressed(pixelgl.KeyC) {
			particles = nil
			draw.Draw(buffer, buffer.Bounds(), image.Transparent, image.ZP, draw.Src)
		}

		if win.Pressed(pixelgl.KeyS) {
			fall(rand.Float64()*fw, fh-(rand.Float64()*fh)/2+(fh/2))
		}

		if win.JustPressed(pixelgl.KeyLeft) {
			angle += 4
		}

		if win.JustPressed(pixelgl.KeyRight) {
			angle -= 4
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
	p.life -= dt * 3

	if p.life > 0 {
		p.position.X += p.velocity.X * dt
		p.position.Y += p.velocity.Y * dt

		if p.life < 10 && p.color.R > 0 {
			if rand.Float64() < 0.3 {
				p.color.R -= 1
				p.color.G -= 1
				p.color.B -= 1
			}
		}
	}
}
