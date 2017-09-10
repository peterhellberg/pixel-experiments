// Inspired by this article
// http://www.iquilezles.org/www/articles/warp/warp.htm

package main

import (
	"flag"
	"image"
	"image/color"
	"math"
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

	picture *pixel.PictureData

	bounds pixel.Rect
	matrix pixel.Matrix

	flipY = pixel.IM.ScaledXY(pixel.ZV, pixel.V(1, -1))
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var (
		fn   string
		seed int64
		d    time.Duration
	)

	flag.StringVar(&fn, "image", "", "image")
	flag.Float64Var(&scale, "scale", 2.667599999999603, "scale")
	flag.Int64Var(&seed, "seed", 1, "seed")
	flag.IntVar(&mode, "mode", 0, "mode")
	flag.DurationVar(&delay, "delay", 32*time.Millisecond, "delay")

	flag.Parse()

	if d <= 0 {
		d = 1 * time.Millisecond
	}

	if err := setup(fn, seed); err != nil {
		log.Fatal().Err(err).Msg("setup")
	}

	pixelgl.Run(run)
}

func pattern(p pixel.Vec) (float64, pixel.Vec, pixel.Vec) {
	q := pixel.V(
		fbm(p.Add(pixel.V(0.0, 0.0))),
		fbm(p.Add(pixel.V(5.2, 1.3))),
	)

	k := pixel.V(
		fbm(p.Add(q.Scaled(scale).Add(pixel.V(1.7, 9.2)))),
		fbm(p.Add(q.Scaled(scale).Add(pixel.V(8.3, 2.8)))),
	)

	return fbm(p.Add(k.Scaled(scale * scale))), q, k
}

func fbm(p pixel.Vec) float64 {
	return noise.Eval2(p.X, p.Y)
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func renderFrame() *image.RGBA {
	frame := image.NewRGBA(source.Bounds())

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			fx, fy := float64(x), float64(y)

			v, q, k := pattern(pixel.V(
				float64(x)*0.00123,
				float64(y)*0.00321,
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
				rx, ry := q.Map(math.Abs).Rotated(k.X * 10).Unit().XY()
				r, g, b = uint8(math.Abs(rx)*255), 0, uint8(math.Abs(ry)*255)
			case 2:
				r, g, b = 255, 255, 0
			case 3:
				r, g, b = 255, 255, 0
			case 4:
				r, g, b = 255, 255, 0
			case 5:
				r, g, b = 255, 255, 0
			case 6:
				r, g, b = 0, 0, 0

				if v < 0.01 {
					r = uint8(255/1 + uint8(int(math.Abs(scale)))%254)
					g = uint8(255/1 + uint8(int(math.Abs(scale)))%254)
				}

				if v < 0.4 {
					b = uint8(255 - uint8(100*math.Abs(scale))%255)
				}
			case 7:
				rx, ry := q.Map(math.Abs).Rotated(k.X * 10).Unit().XY()
				r, g, b = uint8(math.Abs(ry)*255), uint8(math.Abs(rx)*255), 0
			case 8: // See through black and white
				uv := uint8(int(v*255) % 255)

				switch {
				case uv > 64 && uv < 128:
					r, g, b = c.B, c.B, c.B
				case uv > 64:
					r, g, b = 255, 255, 255
				default:
					r, g, b = 0, 0, 0
				}
			case 9: // Black and white
				uv := uint8(int(v*255) % 255)

				switch {
				case uv > 127:
					r, g, b = 255, 255, 255
				default:
					r, g, b = 0, 0, 0
				}
			default:
				r, g, b = c.R, c.G, c.B
			}

			frame.SetRGBA(x, y, color.RGBA{r, g, b, 255})
		}
	}

	return frame

}

func update(ticker *time.Ticker) {
	for range ticker.C {
		if expand {
			scale += 0.0025

			if scale > 4 {
				expand = false
			}
		} else {
			scale -= 0.0025

			if scale < -4 {
				expand = true
			}
		}

		logState()
	}
}

func render(win *pixelgl.Window, canvas *pixelgl.Canvas) {
	canvas.SetPixels(renderFrame().Pix)
	canvas.Draw(win, matrix)
}

func setup(fn string, seed int64) error {
	noise = opensimplex.NewWithSeed(seed)

	m, err := loadImage(fn)
	if err != nil {
		m = xorImage(500, 500)
	}

	w, h = m.Bounds().Dx(), m.Bounds().Dy()

	fw, fh = float64(w), float64(h)

	source = m
	target = image.NewRGBA(m.Bounds())
	bounds = pixel.R(0, 0, fw, fh)
	matrix = flipY.Moved(bounds.Center())

	picture = pixel.PictureDataFromImage(source)

	input = make(chan pixelgl.Button, 4)

	return nil
}

func run() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Bounds:      bounds,
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
			scale += 0.05
			expand = true
		case pixelgl.KeyRight:
			scale -= 0.05
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
	if fn == "" {
		return xorImage(500, 500), nil
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
			c := uint8(x ^ y)

			m.Set(x, y, color.RGBA{c, c % 192, c, 255})
		}
	}

	return m
}
