package main

import (
	"flag"
	"image"
	"math/rand"
	"os"

	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	perlin "github.com/aquilax/go-perlin"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	alpha, beta float64
	seed        int64

	w, h int

	source *image.NRGBA

	noise *perlin.Perlin

	bounds pixel.Rect
	matrix pixel.Matrix

	flipY = pixel.IM.ScaledXY(pixel.ZV, pixel.V(1, -1))
)

func render(win *pixelgl.Window, canvas *pixelgl.Canvas) {
	frame := image.NewRGBA(image.Rect(0, 0, w, h))

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			p := pixel.V(float64(x)*0.012, float64(y)*0.04)
			v := pattern(p)

			frame.Set(x, y, color.NRGBA{
				uint8(v * 255),
				uint8(v * 255),
				uint8(v * 255),
				255,
			})
		}
	}

	canvas.Draw(win, matrix)

	canvas.SetPixels(frame.Pix)
}

func pattern(p pixel.Vec) float64 {
	return pattern1(p)
}

func pattern1(p pixel.Vec) float64 {
	return fbm(p)
}

func fbm(p pixel.Vec) float64 {
	return noise.Noise2D(p.X, p.Y)
}

func main() {
	var fn string

	flag.StringVar(&fn, "i", "../smokestack.jpg", "image")
	flag.Float64Var(&alpha, "a", 7.99, "alpha")
	flag.Float64Var(&beta, "b", 1.42, "beta")
	flag.Int64Var(&seed, "s", 123, "seed")

	flag.Parse()

	if err := setup(fn); err == nil {
		pixelgl.Run(run)
	}
}

func setup(fn string) error {
	noise = perlin.NewPerlinRandSource(alpha, beta, 1, rand.NewSource(seed))

	m, err := loadImage(fn)
	if err != nil {
		return err
	}

	w, h = m.Bounds().Dx(), m.Bounds().Dy()

	source = m
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

	for !win.Closed() {
		win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

		render(win, canvas)

		win.Update()
	}
}

func loadImage(fn string) (*image.NRGBA, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}

	m, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	if rgba, ok := m.(*image.NRGBA); ok {
		return rgba, nil
	}

	b := m.Bounds()

	rgba := image.NewNRGBA(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			rgba.Set(x, y, m.At(x, y))
		}
	}

	return rgba, nil
}
