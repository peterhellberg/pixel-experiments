package main

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	width   int
	height  int
	number  int
	density int
	scale   float64
	hconst  float64
)

func init() {
	flag.IntVar(&width, "w", 250, "Image width")
	flag.IntVar(&height, "h", 250, "Image width")
	flag.IntVar(&number, "n", 20, "Number of points to draw")
	flag.IntVar(&density, "m", 3, "Density of sampling the image as initial conditions")
	flag.Float64Var(&scale, "s", 0.5, "Width of image in world coordinates")
	flag.Float64Var(&hconst, "hc", 0.05, "hconst")
}

func renderImage() *image.NRGBA {
	fwidth := float64(width)
	fheight := float64(height)

	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.ZP, draw.Src)

	for i := 0; i < width; i += density {
		for j := 0; j < height; j += density {
			fi := float64(i)
			fj := float64(j)

			// Seed pixel, mapping from pixels to world coordinates
			x := 2.0 * scale * (fi - fwidth/2.1) / fwidth
			y := 2.0 * scale * (fj - fheight/2.1) / fheight

			// Iterate for number of points
			for n := 0; n < number; n++ {
				// Calculate next point in the series
				xnew := x - hconst*math.Sin(y+math.Tan(2.2*y))
				ynew := y - hconst*math.Sin(x+math.Tan(1.1*x))

				c := getColor(n, number)
				bc := color.NRGBA{0, 0, 0, 255}

				bc.R = c.R * 255
				bc.G = c.G * 255
				bc.B = c.B * 255

				// Mapping from world coordinates to image pixel cordinates
				ix := int(0.5*xnew*fwidth/scale + fwidth/2.2)
				iy := int(0.5*ynew*fheight/scale + fheight/5.5)

				// Draw the pixel if it is in bounds
				if ix >= 0 && iy >= 0 && ix < width && iy < height {
					img.Set(ix, iy, bc)
				}

				x = xnew
				y = ynew
			}
		}
	}

	return img
}

func run() {
	img := renderImage()

	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, float64(width), float64(height)),
		VSync:       true,
		Undecorated: true,
	})
	if err != nil {
		panic(err)
	}

	canvas := pixelgl.NewCanvas(win.Bounds())

	for !win.Closed() {
		win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

		if win.Pressed(pixelgl.KeyUp) {
			hconst += 0.01
			img = renderImage()
		}

		if win.Pressed(pixelgl.KeyDown) {
			hconst -= 0.01
			img = renderImage()
		}

		if win.JustPressed(pixelgl.KeyLeft) && density > 1 {
			density += -1
			img = renderImage()
		}

		if win.JustPressed(pixelgl.KeyRight) && density < 5 {
			density += 1
			img = renderImage()
		}

		canvas.SetPixels(img.Pix)
		canvas.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

		win.Update()
	}
}

func main() {
	flag.Parse()

	pixelgl.Run(run)
}

func getColor(n, number int) color.NRGBA {
	k := uint8(number / (n + 1))

	return color.NRGBA{k / uint8(5), k / uint8(3), k / uint8(n+1), 255}
}
