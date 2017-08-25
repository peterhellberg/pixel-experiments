package main

import (
	"fmt"
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
	w, h = 768, 432
	//w, h   = 192, 108
	fw, fh = float64(w), float64(h)

	globalScale = 1.0
)

var (
	maxSpeed     = 2.1 * globalScale
	friendRadius = 60 * globalScale
	crowdRadius  = friendRadius / 1.3
	avoidRadius  = 90 * globalScale
	coheseRadius = friendRadius

	boids  = Boids{}
	avoids = Avoids{}

	gray = color.RGBA{55, 55, 55, 255}
)

func init() {
	rand.Seed(time.Now().UnixNano())

	setup()
}

func setup() {
	for x := 0; x < w; x += 10 {
		avoids = append(avoids, newAvoid(pixel.V(float64(x), 10), 10, gray))
		avoids = append(avoids, newAvoid(pixel.V(float64(x), fh-10), 10, gray))
	}
}

func run() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, fw, fh),
		Undecorated: true,
		VSync:       true,
	})
	if err != nil {
		panic(err)
	}

	canvas := pixelgl.NewCanvas(win.Bounds())

	last := time.Now()

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

		if win.JustPressed(pixelgl.KeyC) {
			boids, avoids = nil, nil
			setup()
		}

		pos := win.MousePosition()

		if win.JustPressed(pixelgl.KeyO) {
			avoids = append(avoids, newAvoid(pos, 10, gray))
		}

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			boids = append(boids, randomColorBoidAt(pos, dt))
		}

		win.Clear(color.RGBA{0, 0, 0, 255})

		drawFrame(canvas, dt)

		canvas.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

		win.Update()
	}
}

func drawFrame(canvas *pixelgl.Canvas, dt float64) {
	buffer := image.NewRGBA(image.Rect(0, 0, w, h))

	for _, b := range boids {
		b.move(dt)
		b.draw(buffer)
	}

	for _, a := range avoids {
		a.draw(buffer)
	}

	canvas.SetPixels(buffer.Pix)
}

func randomColorBoidAt(p pixel.Vec, dt float64) *boid {
	angle := 90.0 + rand.Float64()*180.0*flip()
	speed := 20.0 + dt + (10.0 * rand.Float64())

	return newBoid(p.X, p.Y, angle, speed,
		color.RGBA{
			uint8(rand.Intn(200)),
			uint8(rand.Intn(200) + 55),
			uint8(rand.Intn(200) + 55),
			255,
		},
	)
}

func main() {
	pixelgl.Run(run)
}

type Avoids []*avoid

func newAvoid(p pixel.Vec, s float64, c color.RGBA) *avoid {
	return &avoid{position: p, size: s, color: c}
}

type avoid struct {
	position pixel.Vec
	size     float64
	color    color.RGBA
}

func (a *avoid) draw(m *image.RGBA) {
	x, y := int(a.position.X), int(a.position.Y)

	r := image.Rect(x-3, y-3, x+3, y+3)

	draw.Draw(m, r, &image.Uniform{a.color}, image.ZP, draw.Src)
}

type boid struct {
	angle         float64
	speed         float64
	position      pixel.Vec
	color         color.RGBA
	originalColor color.RGBA
	friends       []*boid
}

func newBoid(x, y, angle, speed float64, c color.RGBA) *boid {
	return &boid{
		angle:         angle,
		speed:         speed,
		position:      pixel.Vec{x, y},
		color:         c,
		originalColor: c,
		friends:       nil,
	}
}

func (b *boid) velocity() pixel.Vec {
	angleInRadians := b.angle * math.Pi / 180

	return pixel.Vec{
		X: b.speed * math.Cos(angleInRadians),
		Y: -b.speed * math.Sin(angleInRadians),
	}
}

func (b *boid) move(dt float64) {
	b.updatePosition(dt)
	b.updateFriends()
}

func (b *boid) updatePosition(dt float64) {
	v := b.velocity()

	b.position.X += v.X * dt
	b.position.Y += v.Y * dt

	if b.position.X < 0 {
		b.position.X = fw
	}

	if b.position.X > fw {
		b.position.X = 0
	}

	if b.position.Y < 0 {
		b.position.Y = fh
	}

	if b.position.Y > fh {
		b.position.Y = 0
	}
}

type Boids []*boid

func (b *boid) updateFriends() {
	var nearby []*boid

	for _, t := range boids {
		if t != b {
			if math.Abs(t.position.X-b.position.X) < friendRadius &&
				math.Abs(t.position.Y-b.position.Y) < friendRadius {

				// change speed
				//if t.speed > 0 { t.speed -= 0.1 }

				// flip angle
				//t.angle = -t.angle

				nearby = append(nearby, t)
			}
		}
	}

	if len(nearby) > 0 {
		fmt.Println(b, "found", len(nearby), "nearby friends")

		for _, n := range nearby {
			n.color = color.RGBA{255, 0, 0, 255}
		}

	} else {
		b.color = b.originalColor
	}

	b.friends = nearby
}

func (b *boid) draw(m *image.RGBA) {
	x, y := int(b.position.X), int(b.position.Y)

	r := image.Rect(x-3, y-3, x+3, y+3)

	draw.Draw(m, r, &image.Uniform{b.color}, image.ZP, draw.Src)
}

func flip() float64 {
	if rand.Float64() > 0.5 {
		return 1.0
	}

	return -1.0
}
