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
)

func main() {
	var (
		fn    string
		alpha float64
		beta  float64
		n     int
		seed  int64
		d     time.Duration
	)

	flag.StringVar(&fn, "i", "../smokestack.jpg", "image")
	flag.Float64Var(&alpha, "a", 0.105, "alpha")
	flag.Float64Var(&beta, "b", 0.085, "beta")
	flag.IntVar(&n, "n", 2, "n")
	flag.Int64Var(&seed, "s", 123, "seed")
	flag.DurationVar(&d, "d", 16*time.Millisecond, "delay")
	flag.Parse()

	if d <= 0 {
		d = 1 * time.Millisecond
	}

	if err := setup(fn, alpha, beta, n, seed, d); err != nil {
		log.Fatal().Err(err).Msg("setup")
	}

	pixelgl.Run(run)
}

func setup(fn string, alpha, beta float64, n int, seed int64, d time.Duration) error {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

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

	for !win.Closed() {
		win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

		mouse = win.MousePosition()

		if win.JustPressed(pixelgl.Key0) {
			scale = 0
		}

		if win.JustPressed(pixelgl.Key1) {
			scale = 1
		}

		if win.JustPressed(pixelgl.Key2) {
			scale = 2
		}

		if win.JustPressed(pixelgl.Key3) {
			scale = 3
		}

		if win.Pressed(pixelgl.KeyUp) {
			scale += 0.01
		}

		if win.Pressed(pixelgl.KeyDown) {
			scale -= 0.01
		}

		render(win, canvas)

		win.Update()
	}
}

func update(ticker *time.Ticker) {
	sampled := log.Sample(&zerolog.BasicSampler{N: 1000000})

	var expand bool

	//start := time.Now()

	for range ticker.C {
		//d := time.Since(start).Seconds()

		if expand {
			scale += 0.005

			if scale > 200 {
				expand = false
			}
		} else {
			scale -= 0.005

			if scale < -200 {
				expand = true
			}
		}

		frame := image.NewRGBA(source.Bounds())

		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {

				v := pattern(pixel.V(
					float64(x)*0.002,
					float64(y)*0.002,
				))

				if v < 0 {
					v = -v
				}

				if false {
					sampled.Print(v)
				}

				c := source.At(x, y).(color.RGBA)

				var r, g, b uint8

				switch "darken" {
				default:
					r, g, b = c.R, c.G, c.B
				case "black-white":
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
				case "darken":
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
				}

				frame.Set(x, y, color.RGBA{r, g, b, 255})
			}
		}

		target.Pix = frame.Pix
	}
}

// Inspired by this article
// http://www.iquilezles.org/www/articles/warp/warp.htm
func render(win *pixelgl.Window, canvas *pixelgl.Canvas) {
	canvas.SetPixels(target.Pix)

	canvas.Draw(win, matrix)
}

func pattern1(p pixel.Vec) float64 {
	return fbm(p)
}

func pattern2(p pixel.Vec) float64 {
	q := pixel.V(
		fbm(p.Add(pixel.V(0.0, 0.0))),
		fbm(p.Add(pixel.V(5.2, 1.3))),
	)

	return fbm(p.Add(q.Scaled(scale)))
}

func pattern3(p pixel.Vec) float64 {
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

func pattern(p pixel.Vec) float64 {
	return pattern3(p)
}

func fbm(p pixel.Vec) float64 {
	return noise.Eval2(p.X, p.Y)
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
