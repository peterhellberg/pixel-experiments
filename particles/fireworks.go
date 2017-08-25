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

	var total float64

	explode := func(fx, fy float64) {
		if len(particles) > 500 {
			return
		}

		x := int(fx)
		y := int(fy)

		n := 30 + (6 * rand.Intn(10))

		c := color.RGBA{255, uint8(int(total) + y%255), uint8((int(total) ^ x) % 255), 0}

		for i := 0; i < n; i++ {
			angle := 90.0 + rand.Float64()*180.0*flip()
			speed := (2.0 * float64(n)) + (float64(i) * rand.Float64())
			life := 0.1 + (1.3 * rand.Float64())

			particles = append(particles,
				newParticle(fx, fy, angle, speed, life, c),
			)
		}
	}

	go func() {
		for {
			time.Sleep(2700 * time.Millisecond)
			explode(rand.Float64()*fw, (rand.Float64()*fh)/2+(fh/2))
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

		if win.Pressed(pixelgl.MouseButtonLeft) {
			explode(win.MousePosition().XY())
		}

		if win.JustPressed(pixelgl.KeyS) {
			explode(rand.Float64()*fw, (rand.Float64()*fh)/2+(fh/2))
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
	angle    float64
	speed    float64
	position pixel.Vec
	velocity pixel.Vec
	color    color.RGBA
	life     float64
}

func newParticle(x, y, angle, speed, life float64, c color.RGBA) *particle {
	angleInRadians := angle * math.Pi / 180

	return &particle{
		angle:    angle,
		speed:    speed,
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

		if p.life > 1.0 {
			if rand.Float64() < 0.6 {
				p.color.G -= 1
				p.color.B -= 1
			}
		} else {
			p.color.R -= 2
			p.color.G -= 1
			p.color.B -= 1
		}
	}
}
