package main

import (
	"encoding/json"
	"flag"
	"image/color"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

const w, h = 1024, 600

var (
	s       = newState()
	enc     = json.NewEncoder(os.Stdout)
	flipY   = pixel.IM.ScaledXY(pixel.ZV, pixel.V(1, -1))
	pressed = make(chan pixelgl.Button, 5)
)

type state struct {
	Depth   int        `json:"depth"`
	Theta   int        `json:"theta"`
	Angle   float64    `json:"angle"`
	Frac    float64    `json:"frac"`
	Length  float64    `json:"length"`
	Mask    color.RGBA `json:"mask"`
	Circles bool       `json:"circles"`
	updated bool
}

func newState() *state {
	return &state{
		Depth:   4,
		Angle:   29,
		Length:  100,
		Frac:    0.805,
		Mask:    color.RGBA{255, 255, 255, 255},
		updated: true,
	}
}

func run() {
	var delay time.Duration

	flag.DurationVar(&delay, "delay", 100*time.Millisecond, "delay between frames")

	flag.Parse()

	go func(dec *json.Decoder) {
		for {
			if err := dec.Decode(s); err != nil {
				break
			}

			s.updated = true
			time.Sleep(delay)
		}
	}(json.NewDecoder(os.Stdin))

	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, float64(w), float64(h)),
		VSync:       true,
		Undecorated: true,
	})
	if err != nil {
		panic(err)
	}

	go update()

	imd := imdraw.New(nil)
	imd.SetMatrix(flipY)

	for !win.Closed() {
		input(win)

		win.Clear(color.RGBA{55, 56, 84, 255})

		if s.updated {
			win.SetColorMask(s.Mask)

			imd.Clear()
			branch(imd, w/2, 0, s.Length, 0, s.Depth)
			s.updated = false
		}

		imd.Draw(win)

		win.Update()
	}
}

func branch(imd *imdraw.IMDraw, x, y, distance, direction float64, depth int) {
	direction = direction - float64(s.Theta)/11

	x2 := x + distance*math.Sin(direction*math.Pi/180)
	y2 := y - distance*math.Cos(direction*math.Pi/180)

	start, end := pixel.V(x, y), pixel.V(x2, y2)

	if depth > 2 {
		imd.Color = color.RGBA{158, 55, 159, 255}
		imd.Push(start, end)
		imd.Line((float64(depth) + 1) * 3)
	}

	if depth > 1 {
		imd.Color = color.RGBA{232, 106, 240, 55}
		imd.Push(start, end)
		imd.Line((float64(depth) + 1) * 2)
	}

	imd.Color = color.RGBA{255, 0, 0, 255}
	imd.Push(start, end)
	imd.Line(float64(depth) + 1)

	if depth < 1 {
		next := pixel.V(
			x2+distance*math.Sin(direction*math.Pi/180),
			y2-distance*math.Cos(direction*math.Pi/180),
		)

		imd.Color = color.RGBA{55, 56, 84, 55}
		imd.Push(next, pixel.V(x2, y2), end, pixel.V(w/2, -h))
		imd.Polygon(2)
	} else {
		if s.Circles {
			if depth > 0 {
				imd.Color = color.RGBA{232, 106, 240, 55}
				imd.Push(start, end)
				imd.Circle((s.Length/10)*float64(depth), s.Length/60)
			} else {
				imd.Color = color.RGBA{232, 106, 240, 55}
				imd.Push(start, end)
				imd.Line((float64(depth) + 1) * 2)
			}
		}

		branch(imd, x2, y2, distance*s.Frac, direction-s.Angle, depth-1)
		branch(imd, x2, y2, distance*s.Frac, direction+s.Angle, depth-1)
	}
}

func update() {
	for key := range pressed {
		switch key {
		case pixelgl.KeyA:
			s.Frac += 0.005
		case pixelgl.KeyZ:
			s.Frac -= 0.005
		case pixelgl.KeyC:
			s.Circles = !s.Circles
		case pixelgl.KeyR:
			s.Mask = color.RGBA{255, 128, 128, 255}
		case pixelgl.KeyG:
			s.Mask = color.RGBA{128, 255, 128, 255}
		case pixelgl.KeyB:
			s.Mask = color.RGBA{128, 128, 255, 255}
		case pixelgl.KeyW:
			s.Mask = color.RGBA{255, 255, 255, 255}
		case pixelgl.KeyUp:
			s.Length += 0.5
		case pixelgl.KeyDown:
			s.Length -= 0.5
		case pixelgl.KeyLeft:
			s.Angle += 0.5
			s.Theta += rand.Intn(5)
		case pixelgl.KeyRight:
			s.Angle -= 0.5
			s.Theta -= rand.Intn(5)
		case pixelgl.Key1:
			s.Depth = 1
		case pixelgl.Key2:
			s.Depth = 2
		case pixelgl.Key3:
			s.Depth = 3
		case pixelgl.Key4:
			s.Depth = 4
		case pixelgl.Key5:
			s.Depth = 5
		case pixelgl.Key6:
			s.Depth = 6
		case pixelgl.Key7:
			s.Depth = 7
		case pixelgl.Key8:
			s.Depth = 8
		case pixelgl.Key9:
			s.Depth = 9
		}

		enc.Encode(s)

		s.updated = true
	}
}

func input(win *pixelgl.Window) {
	win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

	if win.Pressed(pixelgl.KeyA) {
		pressed <- pixelgl.KeyA
	}

	if win.Pressed(pixelgl.KeyZ) {
		pressed <- pixelgl.KeyZ
	}

	if win.JustPressed(pixelgl.KeyC) {
		pressed <- pixelgl.KeyC
	}

	if win.JustPressed(pixelgl.KeyR) {
		pressed <- pixelgl.KeyR
	}

	if win.JustPressed(pixelgl.KeyG) {
		pressed <- pixelgl.KeyG
	}

	if win.JustPressed(pixelgl.KeyB) {
		pressed <- pixelgl.KeyB
	}

	if win.JustPressed(pixelgl.KeyW) {
		pressed <- pixelgl.KeyW
	}

	if win.Pressed(pixelgl.KeyUp) {
		pressed <- pixelgl.KeyUp
	}

	if win.Pressed(pixelgl.KeyDown) {
		pressed <- pixelgl.KeyDown
	}

	if win.Pressed(pixelgl.KeyLeft) {
		pressed <- pixelgl.KeyLeft
	}

	if win.Pressed(pixelgl.KeyRight) {
		pressed <- pixelgl.KeyRight
	}

	if win.JustPressed(pixelgl.Key1) {
		pressed <- pixelgl.Key1
	}

	if win.JustPressed(pixelgl.Key2) {
		pressed <- pixelgl.Key2
	}

	if win.JustPressed(pixelgl.Key3) {
		pressed <- pixelgl.Key3
	}

	if win.JustPressed(pixelgl.Key4) {
		pressed <- pixelgl.Key4
	}

	if win.JustPressed(pixelgl.Key5) {
		pressed <- pixelgl.Key5
	}

	if win.JustPressed(pixelgl.Key6) {
		pressed <- pixelgl.Key6
	}

	if win.JustPressed(pixelgl.Key7) {
		pressed <- pixelgl.Key7
	}

	if win.JustPressed(pixelgl.Key8) {
		pressed <- pixelgl.Key8
	}

	if win.JustPressed(pixelgl.Key9) {
		pressed <- pixelgl.Key9
	}

	if win.JustPressed(pixelgl.Key0) {
		pressed <- pixelgl.Key0
	}

}

func main() {
	pixelgl.Run(run)
}
