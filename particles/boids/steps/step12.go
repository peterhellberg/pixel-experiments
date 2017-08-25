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
	w, h   = 192 * 2, 108 * 2
	fw, fh = float64(w), float64(h)

	globalScale = 0.98
)

var (
	maxSpeed     = 2.1 * globalScale
	desireAmount = 0.5 * globalScale
	friendRadius = 120 * globalScale
	crowdRadius  = friendRadius / 1.3
	avoidRadius  = 30 * globalScale
	coheseRadius = friendRadius / 4.1

	boids  = Boids{}
	avoids = Avoids{}

	gray = color.RGBA{55, 55, 55, 255}
)

func init() {
	rand.Seed(time.Now().UnixNano())

	setup()
}

func setup() {
	for x := 0; x < w; x += 15 {
		avoids = append(avoids, newAvoid(pixel.V(float64(x+5), 10), 0, gray))
		avoids = append(avoids, newAvoid(pixel.V(float64(x+5), fh-10), 0, gray))
	}

	for y := 0; y < h; y += 15 {
		avoids = append(avoids, newAvoid(pixel.V(10, float64(y+5)), 0, gray))
		avoids = append(avoids, newAvoid(pixel.V(fw-10, float64(y+5)), 0, gray))
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

		if win.Pressed(pixelgl.Key1) {
			desireAmount = 1
		}

		if win.Pressed(pixelgl.Key2) {
			desireAmount = 2
		}

		if win.Pressed(pixelgl.Key3) {
			desireAmount = 10
		}

		if win.Pressed(pixelgl.Key4) {
			desireAmount = 20
		}

		if win.Pressed(pixelgl.KeyUp) {
			desireAmount += 0.5
		}

		if win.Pressed(pixelgl.KeyDown) {
			desireAmount -= 0.5
		}

		pos := win.MousePosition()

		if win.Pressed(pixelgl.KeyO) {
			avoids = append(avoids, newAvoid(pos, 10, gray))
		}

		if win.Pressed(pixelgl.MouseButtonLeft) {
			boids = append(boids, randomColorBoidAt(pos, dt))
		}

		win.Clear(color.RGBA{0, 0, 0, 255})

		drawFrame(canvas)

		canvas.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

		win.Update()
	}
}

func drawFrame(canvas *pixelgl.Canvas) {
	buffer := image.NewRGBA(image.Rect(0, 0, w, h))

	for _, a := range avoids {
		a.draw(buffer)
	}

	for _, b := range boids {
		b.increment()
		b.wrap()

		if b.think == 0 {
			b.updateFriends()
		}

		b.flock()

		b.updatePosition()

		b.draw(buffer)
	}

	canvas.SetPixels(buffer.Pix)
}

type Boids []*boid

type boid struct {
	think         int
	position      pixel.Vec
	velocity      pixel.Vec
	color         color.RGBA
	originalColor color.RGBA
	friends       []*boid
}

func newBoid(x, y, angle, speed float64, c color.RGBA) *boid {
	angleInRadians := angle * math.Pi / 180

	return &boid{
		think:    rand.Intn(100),
		position: pixel.Vec{x, y},
		velocity: pixel.Vec{
			X: speed * math.Cos(angleInRadians),
			Y: -speed * math.Sin(angleInRadians),
		},
		color:         c,
		originalColor: c,
		friends:       nil,
	}
}

func randomColorBoidAt(p pixel.Vec, dt float64) *boid {
	angle := (90.0 + rand.Float64()) * 180.0 * flip()
	speed := maxSpeed * rand.Float64()

	return newBoid(p.X, p.Y, angle, speed, randomColor())
}

func (b *boid) increment() {
	b.think = (b.think + 1) % 10
}

func (b *boid) wrap() {
	b.position.X = float64(int(b.position.X+fw) % w)
	b.position.Y = float64(int(b.position.Y+fh) % h)
}

func (b *boid) updatePosition() {
	b.position = b.position.Add(b.velocity)
}

func (b *boid) updateFriends() {
	var nearby []*boid

	for _, t := range boids {
		if t != b {
			if math.Abs(t.position.X-b.position.X) < friendRadius &&
				math.Abs(t.position.Y-b.position.Y) < friendRadius {
				nearby = append(nearby, t)
			}
		}
	}

	b.friends = nearby
}

func (b *boid) getAverageColor() color.RGBA {
	if false {
		return color.RGBA{255, 0, 0, 255}
	}

	c := len(b.friends)

	tr, tg, tb := 0, 0, 0
	br, bg, bb := int(b.color.R), int(b.color.G), int(b.color.B)

	for _, f := range b.friends {
		fr, fg, fb := int(f.originalColor.R), int(f.originalColor.G), int(f.originalColor.B)

		if fr-br < -128 {
			tr += fr + 255 - br
		} else if fr-br > 128 {
			tr += fr - 255 - br
		} else {
			tr += fr - br
		}

		if fg-bg < -128 {
			tg += fg + 255 - bg
		} else if fg-bg > 128 {
			tg += fg - 255 - bg
		} else {
			tg += fg - bg
		}

		if fb-bb < -128 {
			tb += fb + 255 - bb
		} else if fb-bb > 128 {
			tb += fb - 255 - bb
		} else {
			tb += fb - bb
		}
	}

	return color.RGBA{
		uint8(float64(tr) / float64(c)),
		uint8(float64(tg) / float64(c)),
		uint8(float64(tb) / float64(c)),
		255,
	}
}

func (b *boid) getAverageDir() pixel.Vec {
	sum := pixel.V(0, 0)

	for _, f := range b.friends {
		d := dist(b.position, f.position)

		if d > 0 && d < friendRadius {
			sum = sum.Add(div(f.velocity.Unit(), d))
		}
	}

	return sum
}

func (b *boid) getAvoidDir() pixel.Vec {
	steer := pixel.V(0, 0)

	for _, f := range b.friends {
		d := dist(b.position, f.position)

		if d > 0 && d < crowdRadius {
			diff := div(b.position.Sub(f.position).Unit(), d)
			steer = steer.Add(diff)
		}
	}

	return steer
}

func (b *boid) getAvoidObjects() pixel.Vec {
	steer := pixel.V(0, 0)

	for _, f := range avoids {
		d := dist(b.position, f.position)

		if d > 0 && d < avoidRadius {
			diff := div(b.position.Sub(f.position).Unit(), d)
			steer = steer.Add(diff)
		}
	}

	return steer
}

func (b *boid) getCohesion() pixel.Vec {
	sum := pixel.V(0, 0)

	count := 0

	for _, other := range b.friends {
		d := dist(b.position, other.position)

		if d > 0 && d < coheseRadius {
			sum = sum.Add(other.position)
			count++
		}
	}

	if count > 0 {
		desired := div(sum, float64(count)).Sub(b.position)

		return desired.Unit().Scaled(desireAmount)
	}

	return pixel.V(0, 0)
}

func (b *boid) move(v pixel.Vec) {
	b.velocity = b.velocity.Add(v)
}

func (b *boid) limitSpeed(s float64) {
	b.velocity = b.velocity.Unit().Scaled(s)
}

func (b *boid) updateColor() {
	if len(b.friends) > 0 {
		b.color = b.getAverageColor()
	}
}

func (b *boid) flock() {
	var (
		align        = b.getAverageDir().Scaled(1)
		cohesion     = b.getCohesion().Scaled(1)
		avoidDir     = b.getAvoidDir().Scaled(1)
		avoidObjects = b.getAvoidObjects().Scaled(1)

		noise = pixel.V(rand.Float64()*2-1, rand.Float64()*2-1).Scaled(0.05)
	)

	b.move(align)
	b.move(avoidDir)
	b.move(avoidObjects)
	b.move(noise)
	b.move(cohesion)

	b.limitSpeed(maxSpeed)

	b.updateColor()
}

func (b *boid) draw(m *image.RGBA) {
	x, y := int(b.position.X), int(b.position.Y)

	r := image.Rect(x-3, y-3, x+3, y+3)

	draw.Draw(m, r, &image.Uniform{b.color}, image.ZP, draw.Src)

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

	r := image.Rect(x-2, y-3, x+2, y+3)

	draw.Draw(m, r, &image.Uniform{a.color}, image.ZP, draw.Src)
}

func randomColor() color.RGBA {
	if true {
		return color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255}
	}

	var r, g, b uint8

	i := rand.Intn(3)

	switch i {
	case 0:
		r = 255
		g = uint8(rand.Intn(100) + 50)
		b = uint8(rand.Intn(100) + 50)
	case 1:
		r = uint8(rand.Intn(100) + 50)
		g = 255
		b = uint8(rand.Intn(100) + 50)
	default:
		r = uint8(rand.Intn(200) + 50)
		g = uint8(rand.Intn(100) + 50)
		b = 255
	}

	return color.RGBA{r, g, b, 255}
}

func flip() float64 {
	if rand.Float64() > 0.5 {
		return 1.0
	}

	return -1.0
}

func dist(a, b pixel.Vec) float64 {
	return (a.X - b.X) + (a.Y - b.Y)
	//return math.Sqrt(a.Sub(b).Dot(a))
}

func div(v pixel.Vec, d float64) pixel.Vec {
	v.X /= d
	v.Y /= d

	return v
}

func main() {
	pixelgl.Run(run)
}
