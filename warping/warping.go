// Inspired by this article
// http://www.iquilezles.org/www/articles/warp/warp.htm

package main

import (
	"flag"
	"image"
	"image/color"
	"os"
	"time"

	_ "image/jpeg"
	_ "image/png"

	pixel "github.com/faiface/pixel"
	pixelgl "github.com/faiface/pixel/pixelgl"
	opensimplex "github.com/ojrac/opensimplex-go"
	zerolog "github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
)

var (
	mode   int
	scale  float64
	expand bool

	w, h   int
	fw, fh float64

	pos pixel.Vec

	input chan pixelgl.Button

	delay time.Duration
	start time.Time

	noise *opensimplex.Noise

	source *image.RGBA
	target *image.RGBA

	bounds pixel.Rect
	matrix pixel.Matrix

	flipY = pixel.IM.ScaledXY(pixel.ZV, pixel.V(1, -1))
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var (
		fn   string
		seed int64
		posX float64
		posY float64
		d    time.Duration
	)

	flag.StringVar(&fn, "image", "", "image")
	flag.Float64Var(&scale, "scale", 1, "scale")
	flag.Float64Var(&posX, "x", -2.07, "pos.X")
	flag.Float64Var(&posY, "y", 0.257, "pos.Y")
	flag.Int64Var(&seed, "seed", 1, "seed")
	flag.IntVar(&mode, "mode", 1, "mode")
	flag.DurationVar(&delay, "delay", 100*time.Millisecond, "delay")

	flag.Parse()

	if d <= 0 {
		d = 1 * time.Millisecond
	}

	if err := setup(fn, pixel.V(posX, posY), seed); err != nil {
		log.Fatal().Err(err).Msg("setup")
	}

	pixelgl.Run(run)
}

func setup(fn string, p pixel.Vec, seed int64) error {
	m, err := loadImage(fn)
	if err != nil {
		m = xorImage(256, 256)
	}

	w, h = m.Bounds().Dx(), m.Bounds().Dy()

	fw, fh = float64(w), float64(h)

	source = m
	target = image.NewRGBA(source.Bounds())
	bounds = pixel.R(0, 0, fw, fh)
	matrix = flipY.Moved(bounds.Center()).Scaled(pixel.ZV, scale)

	input = make(chan pixelgl.Button, 1024)

	pos = p

	noise = opensimplex.NewWithSeed(seed)

	return nil
}

func background(ticker *time.Ticker) {
	for range ticker.C {
		logState()
	}
}

func run() {
	start = time.Now()

	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, fw*scale, fh*scale),
		VSync:       true,
		Undecorated: true,
	})
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	go background(ticker)
	go handleInput()

	canvas := pixelgl.NewCanvas(bounds)

	tickRate := 1.0 / 30
	adt := 0.0

	last := time.Now()
	for !win.Closed() {
		adt += time.Since(last).Seconds()
		last = time.Now()

		processInput(win)

		if adt >= tickRate {
			adt -= tickRate
			update()
		}

		canvas.SetPixels(target.Pix)
		canvas.Draw(win, matrix)

		win.Update()
	}
}

func update() {
	if expand {
		pos.Y += 0.0025
		pos.X += 0.0025

		if pos.X > 4 {
			expand = false
		}
	} else {
		pos.X -= 0.0025
		pos.Y -= 0.0025

		if pos.X < -4 {
			expand = true
		}
	}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			target.SetRGBA(x, y, pixelColor(x, y))
		}
	}
}

func pattern(p pixel.Vec) (float64, pixel.Vec, pixel.Vec) {
	q := pixel.V(
		fbm(p.Add(pixel.V(0.0, 0.0))),
		fbm(p.Add(pixel.V(5.2, 1.3))),
	)

	k := pixel.V(
		fbm(p.Add(q.Scaled(pos.Y).Add(pixel.V(1.7, 9.2)))),
		fbm(p.Add(q.Scaled(pos.X).Add(pixel.V(8.3, 2.8)))),
	)

	v := fbm(p.Add(k.Scaled(pos.Y * pos.X)))

	return v, q, k
}

func fbm(p pixel.Vec) float64 {
	return noise.Eval2(p.X, p.Y)
}

func warp(fx, fy float64, k pixel.Vec) (int, int) {
	wx, wy := int(fx*k.X)%w, int(fy*k.Y)%h

	if wx < 0 || wx > w {
		wx = int(fx - k.X)
	}

	if wy < 0 || wy > h {
		wy = int(fy - k.Y)
	}

	return wx, wy
}

func pixelColor(x, y int) color.RGBA {
	fx, fy := float64(x), float64(y)

	v, q, k := pattern(pixel.V(fx*0.00123, fy*0.00321))

	wx, wy := warp(fx, fy, k)

	c := source.At(wx, wy).(color.RGBA)

	switch mode {
	case 1: // Source warped by pattern
		c = source.At(wx, wy).(color.RGBA)
	case 2: // Source warped by pattern, glowing edges
		if v > 0.1 {
			c = source.At(x, y).(color.RGBA)
		} else if v > -0.1 {
			c = color.RGBA{c.R, c.G / 2, c.B / 4, 25}
		}
	case 3: // Grayscale
		// Weighted luminosity average for human perception as per the GIMP docs:
		// https://docs.gimp.org/2.8/en/gimp-tool-desaturate.html
		a := uint8(float64(c.R)*0.21 + float64(c.G)*0.72 + float64(c.B)*0.07)

		c = color.RGBA{a, a, a, 255}
	case 4: // Black and white
		uv := uint8(int(v*255) % 255)

		switch {
		case uv > 127:
			c = color.RGBA{255, 255, 255, 255}
		default:
			c = color.RGBA{0, 0, 0, 255}
		}
	case 5: // See through black and white
		uv := uint8(int(v*255) % 255)

		switch {
		case uv > 32 && uv < 128:
			a := uint8(float64(c.R)*0.21 + float64(c.G)*0.72 + float64(c.B)*0.07)

			c = color.RGBA{a, a, a, 255}
		case uv > 64:
			c = color.RGBA{255, 255, 255, 255}
		default:
			c = color.RGBA{0, 0, 0, 255}
		}
	case 6: // Black for now
		c = color.RGBA{0, 0, 0, 255}
	case 7: // Black for now
		c = color.RGBA{0, 0, 0, 255}
	case 8: // Black for now
		c = color.RGBA{0, 0, 0, 255}
	case 9: // Experimental
		c = source.At(x, y).(color.RGBA)
		a := uint8(float64(c.R)*0.21 + float64(c.G)*0.72 + float64(c.B)*0.07)

		switch {
		case q.Y > 0 && v > 0:
			c = color.RGBA{255, a, a, 255}
		case q.Y < 0 && v < 0:
			c = color.RGBA{0, a, 255, 255}
		case q.X > 0 && k.X < 0:
			c = color.RGBA{255, 255, a, 255}
		case q.X > 0 && k.Y > 0:
			c = color.RGBA{a, 255, 255, 255}
		default:
			c = color.RGBA{a, a, a, 255}
		}
	case 0: // Source image without any warping
		c = source.At(x, y).(color.RGBA)
	}

	return c
}

func processInput(win *pixelgl.Window) {
	win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

	for _, button := range pressedButtons {
		if win.Pressed(button) {
			input <- button
		}
	}

	for _, button := range justPressedButtons {
		if win.JustPressed(button) {
			input <- button
		}
	}
}

func handleInput() {
	for button := range input {
		switch button {
		case pixelgl.Key0:
			mode = 0
		case pixelgl.Key1:
			mode = 1
		case pixelgl.Key2:
			mode = 2
		case pixelgl.Key3:
			mode = 3
		case pixelgl.Key4:
			mode = 4
		case pixelgl.Key5:
			mode = 5
		case pixelgl.Key6:
			mode = 6
		case pixelgl.Key7:
			mode = 7
		case pixelgl.Key8:
			mode = 8
		case pixelgl.Key9:
			mode = 9
		case pixelgl.KeyUp:
			pos.Y += 0.05
			expand = true
		case pixelgl.KeyDown:
			pos.Y -= 0.05
			expand = false
		case pixelgl.KeyLeft:
			pos.X += 0.05
			expand = true
		case pixelgl.KeyRight:
			pos.X -= 0.05
			expand = false
		case pixelgl.KeyS:
			pos.X = 0
			pos.Y = 0
		}

		logState()
	}
}

var pressedButtons = []pixelgl.Button{
	pixelgl.KeyUp,
	pixelgl.KeyDown,
	pixelgl.KeyLeft,
	pixelgl.KeyRight,
}

var justPressedButtons = []pixelgl.Button{
	pixelgl.Key0,
	pixelgl.Key1,
	pixelgl.Key2,
	pixelgl.Key3,
	pixelgl.Key4,
	pixelgl.Key5,
	pixelgl.Key6,
	pixelgl.Key7,
	pixelgl.Key8,
	pixelgl.Key9,
	pixelgl.KeyL,
	pixelgl.KeyS,
}

func logState() {
	log.Info().
		Int("mode", mode).
		Interface("pos", pos).
		Msg("State")
}

func loadImage(fn string) (*image.RGBA, error) {
	if fn == "" {
		return xorImage(400, 300), nil
	}

	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}

	m, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	if rgba, ok := m.(*image.RGBA); ok {
		return rgba, nil
	}

	b := m.Bounds()

	rgba := image.NewRGBA(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			rgba.Set(x, y, m.At(x, y))
		}
	}

	return rgba, nil
}

func xorImage(w, h int) *image.RGBA {
	m := image.NewRGBA(image.Rect(0, 0, w, h))

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			v := uint8(x ^ y)

			c := color.RGBA{v, v % 128, v, 255}

			if x >= w/2 {
				c = color.RGBA{v, v, v % 128, 255}
			}

			m.Set(x, y, c)
		}
	}

	return m
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
