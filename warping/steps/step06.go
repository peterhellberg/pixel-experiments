// Inspired by this article
// http://www.iquilezles.org/www/articles/warp/warp.htm

package main

import (
	"flag"
	"image"
	"image/color"
	"os"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	pixel "github.com/faiface/pixel"
	pixelgl "github.com/faiface/pixel/pixelgl"
	opensimplex "github.com/ojrac/opensimplex-go"
	zerolog "github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
)

func pattern3(p pixel.Vec) (float64, pixel.Vec, pixel.Vec) {
	q := pixel.V(
		fbm(p.Add(pixel.V(0.0, 0.0))),
		fbm(p.Add(pixel.V(15.2, 1.3))),
	)

	r := pixel.V(
		fbm(p.Add(q.Scaled(scale).Add(pixel.V(1.7, 9.2)))),
		fbm(p.Add(q.Scaled(scale).Add(pixel.V(8.3, 2.8)))),
	)

	return fbm(p.Add(r.Scaled(scale))), q, r
}

func pattern(p pixel.Vec) float64 {
	q := pixel.V(
		fbm(p.Add(pixel.V(0.0, 0.0))),
		fbm(p.Add(pixel.V(5.2, 1.3))),
	)

	r := pixel.V(
		fbm(p.Add(q.Scaled(scale).Add(pixel.V(1.7, 9.2)))),
		fbm(p.Add(q.Scaled(scale).Add(pixel.V(8.3, 2.8)))),
	)

	return fbm(p.Add(r.Scaled(scale)))
}

func fbm(p pixel.Vec) float64 {
	return noise.Eval2(p.X, p.Y)
}

func drawFrame() {
	frame := image.NewRGBA(source.Bounds())

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			fx, fy := float64(x), float64(y)

			v, q, k := pattern3(pixel.V(
				float64(x)*0.0012,
				float64(y)*0.002,
			))

			if v < 0 {
				v = -v
			}

			wx, wy := int(fx*k.X)%w, int(fy*k.Y)%h

			if wx < 0 || wx > w {
				wx = int(fx - k.X)
			}

			if wy < 0 || wy > h {
				wy = int(fy - k.Y)
			}

			c := source.At(wx, wy).(color.RGBA)

			var r, g, b uint8

			switch mode {
			case 1:
				uv := uint8(int(v*255) % 255)

				if uv > 100 && uv < c.R {
					r = c.R - uv + uint8(q.X*fw)
				} else {
					r = (c.R * 2) % 255
				}

				if uv > 100 && uv < c.G {
					g = c.G - uv + uint8(q.X*fw)
				} else {
					g = (c.G * 2) % 255
				}

				if uv > 100 && uv < c.B {
					b = c.B - uv + uint8(q.X*fw)
				} else {
					b = (c.B * 2) % 255
				}
			case 2:
				uv := uint8(int(v*250) % 254)

				if nr := (255 - c.R) + uv; nr < 155 {
					r = uv
				}

				if ng := (255 - c.G) + uv; ng < 255 {
					g = c.G
				}
			case 3:
				uv := uint8(int(v*250) % 254)

				if nb := (255 - c.B) + uv; nb < 255 {
					b = uv
				}
			default:
				r, g, b = c.R, c.G, c.B
			}

			frame.SetRGBA(x, y, color.RGBA{r, g, b, 255})
		}
	}

	target.Pix = frame.Pix
}

func update(ticker *time.Ticker) {
	drawFrame()

	for range ticker.C {
		if expand {
			scale += 0.0025

			if scale > 100 {
				expand = false
			}
		} else {
			scale -= 0.0025

			if scale < -100 {
				expand = true
			}
		}

		drawFrame()
	}
}

func render(win *pixelgl.Window, canvas *pixelgl.Canvas) {
	canvas.SetPixels(target.Pix)

	canvas.Draw(win, matrix)
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var (
		fn   string
		seed int64
		d    time.Duration
	)

	flag.StringVar(&fn, "image", "../flower.png", "image")
	flag.Float64Var(&scale, "scale", 7.847599999999528, "scale")
	flag.Int64Var(&seed, "seed", 1, "seed")
	flag.IntVar(&mode, "mode", 0, "mode")
	flag.DurationVar(&d, "delay", 8*time.Millisecond, "delay")

	flag.Parse()

	if d <= 0 {
		d = 1 * time.Millisecond
	}

	if err := setup(fn, seed, d); err != nil {
		log.Fatal().Err(err).Msg("setup")
	}

	pixelgl.Run(run)
}

func setup(fn string, seed int64, d time.Duration) error {
	delay = d
	noise = opensimplex.NewWithSeed(seed)

	m, err := loadImage(fn)
	if err != nil {
		return err
	}

	w, h = m.Bounds().Dx(), m.Bounds().Dy()
	fw, fh = float64(w), float64(h)

	source = m
	target = image.NewRGBA(source.Bounds())
	bounds = pixel.R(0, 0, float64(w), float64(h))
	matrix = flipY.Moved(bounds.Center())

	input = make(chan pixelgl.Button, 1)

	return nil
}

func run() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:       "Domain Warping",
		Bounds:      bounds,
		Resizable:   false,
		VSync:       true,
		Undecorated: true,
	})
	if err != nil {
		panic(err)
	}

	canvas := pixelgl.NewCanvas(bounds)

	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	go update(ticker)
	go handleInput()

	for !win.Closed() {
		mouse = win.MousePosition()

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

		render(win, canvas)

		win.Update()
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
			mode += 1

			if mode > 9 {
				mode = 0
			}
		case pixelgl.KeyDown:
			mode -= 1

			if mode < 0 {
				mode = 9
			}
		case pixelgl.KeyLeft:
			scale += 0.01
			expand = true
		case pixelgl.KeyRight:
			scale -= 0.01
			expand = false
		}

		logState()
	}
}

var pressedButtons = []pixelgl.Button{
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
	pixelgl.KeyUp,
	pixelgl.KeyDown,
	pixelgl.KeyL,
}

func logState() {
	log.Info().
		Int("mode", mode).
		Float64("scale", scale).
		Msg("State")
}

func loadImage(fn string) (*image.RGBA, error) {
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

var (
	mode   int
	scale  float64
	expand bool

	w, h   int
	fw, fh float64

	mouse pixel.Vec

	input chan pixelgl.Button

	delay time.Duration

	noise *opensimplex.Noise

	source *image.RGBA
	target *image.RGBA

	bounds pixel.Rect
	matrix pixel.Matrix

	flipY = pixel.IM.ScaledXY(pixel.ZV, pixel.V(1, -1))
)
