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
	return noise.Eval3(p.X, p.Y, mouse.Cross(p)/float64(w))
}

func update(ticker *time.Ticker) {
	var expand bool

	for range ticker.C {
		if expand {
			scale += 0.005

			if scale > 100 {
				expand = false
			}
		} else {
			scale -= 0.005

			if scale < -100 {
				expand = true
			}
		}

		frame := image.NewRGBA(source.Bounds())

		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {

				v := pattern(pixel.V(
					float64(x)*0.0012,
					float64(y)*0.002,
				))

				if v < 0 {
					v = -v * 2
				}

				c := source.At(x, y).(color.RGBA)

				var r, g, b uint8

				switch mode {
				case 1:
					uv := uint8(int(v*255) % 255)

					if uv < c.R {
						r = c.R - uv
					}

					if uv < c.G {
						g = c.G - uv
					}

					if uv < c.B {
						b = c.B - uv
					}
				case 2:
					uv := uint8(int(v*2) % 255)

					if nr := (255 - c.R) + uv; nr < 255 {
						r = nr + c.R
					}

					if ng := (255 - c.G) + uv; ng < 255 {
						g = ng + c.G
					}

					if nb := (255 - c.B) + uv; nb < 255 {
						b = nb + c.B
					}
				default:
					r, g, b = c.R, c.G, c.B
				}

				frame.SetRGBA(x, y, color.RGBA{r, g, b, 255})
			}
		}

		target.Pix = frame.Pix
	}
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	var (
		fn   string
		seed int64
		d    time.Duration
	)

	flag.StringVar(&fn, "i", "../smokestack.jpg", "image")
	flag.Int64Var(&seed, "s", 1024, "seed")
	flag.DurationVar(&d, "d", 8*time.Millisecond, "delay")

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

	source = m
	target = image.NewRGBA(source.Bounds())
	bounds = pixel.R(0, 0, float64(w), float64(w))
	matrix = flipY.Moved(bounds.Center())

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

	canvas.SetComposeMethod(pixel.ComposeXor)

	ticker := time.NewTicker(delay)

	defer ticker.Stop()

	go update(ticker)

	for !win.Closed() {
		win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

		mouse = win.MousePosition()

		if win.JustPressed(pixelgl.Key0) {
			scale = 0
		}

		if win.JustPressed(pixelgl.Key1) {
			mode = 1
		}

		if win.JustPressed(pixelgl.Key2) {
			mode = 2
		}

		if win.JustPressed(pixelgl.Key3) {
			mode = 3
		}

		if win.JustPressed(pixelgl.KeyRight) {
			mode += 1

			if mode > 2 {
				mode = 0
			}
		}

		if win.JustPressed(pixelgl.KeyLeft) {
			mode -= 1

			if mode < 0 {
				mode = 2
			}
		}

		if win.Pressed(pixelgl.KeyUp) {
			scale += 0.01
		}

		if win.Pressed(pixelgl.KeyDown) {
			scale -= 0.01
		}

		if win.JustPressed(pixelgl.KeyL) {
			log.Info().
				Int("mode", mode).
				Float64("scale", scale).
				Msg("State")
		}

		render(win, canvas)

		win.Update()
	}
}

func render(win *pixelgl.Window, canvas *pixelgl.Canvas) {
	canvas.SetPixels(target.Pix)

	canvas.Draw(win, matrix)
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
	w, h int

	mouse pixel.Vec

	scale float64

	delay time.Duration

	noise *opensimplex.Noise

	source *image.RGBA
	target *image.RGBA

	bounds pixel.Rect
	matrix pixel.Matrix

	flipY = pixel.IM.ScaledXY(pixel.ZV, pixel.V(1, -1))

	mode = 1
)
