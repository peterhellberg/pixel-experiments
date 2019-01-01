package main

import (
	"image"
	"image/color"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	width  = 160
	height = 96
	scale  = 6
	delay  = 32 * time.Millisecond
)

func run() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, float64(width)*scale, float64(height)*scale),
		VSync:       true,
		Undecorated: true,
		Resizable:   false,
	})
	if err != nil {
		panic(err)
	}

	i := NewInferno(width, height)
	c := win.Bounds().Center()

	ticker := time.NewTicker(delay)

	go func() {
		for range ticker.C {
			i.Spread()
			i.Render()
		}
	}()

	for !win.Closed() {
		p := pixel.PictureDataFromImage(i.buffer)

		pixel.NewSprite(p, p.Bounds()).
			Draw(win, pixel.IM.Moved(c).Scaled(c, scale*1.1))

		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			return
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

type Inferno struct {
	width  int
	height int
	grid   []int8
	buffer *image.RGBA
}

func NewInferno(width, height int) *Inferno {
	i := &Inferno{width: width, height: height}

	i.init()

	return i
}

func (i *Inferno) init() {
	i.grid = make([]int8, i.width*i.height)

	for j := 0; j < i.width; j++ {
		i.grid[((i.height-1)*i.width)+j] = 42
	}

	i.buffer = image.NewRGBA(image.Rect(0, 0, i.width, i.height))
}

func (i *Inferno) Render() {
	for y := 0; y < i.height; y += 2 {
		for x := 0; x < i.width; x++ {
			if pos := (y * i.width) + x; pos > 0 && i.grid[pos] != i.grid[pos-1] {
				i.buffer.Set(x, y, mapColor(i.grid[pos]))
			}

			if pos := ((y + 1) * i.width) + x; i.grid[pos] != i.grid[pos-1] {
				i.buffer.Set(x, y+1, mapColor(i.grid[pos]))
			}
		}
	}
}

func (i *Inferno) Spread() {
	for y := i.height - 1; y > 0; y-- {
		for x := 0; x < i.width; x++ {
			src := (y * i.width) + x
			dst := (src - i.width) + rand.Intn(5) - 2

			if dst < 0 {
				dst = 0
			}

			if end := (i.width * i.height) - 1; dst > end {
				dst = end
			}

			i.grid[dst] = i.grid[src] - int8(rand.Intn(6)-1)

			if i.grid[dst] > 32 {
				i.grid[dst] = 32
			}

			if i.grid[dst] < 0 {
				i.grid[dst] = 0
			}
		}
	}
}

var cmap = []color.RGBA{
	{0x07, 0x07, 0x07, 0xdc}, {0x1f, 0x07, 0x07, 0xdc},
	{0x2f, 0x0f, 0x07, 0xdc}, {0x47, 0x0f, 0x07, 0xdc},
	{0x57, 0x17, 0x07, 0xdc}, {0x67, 0x1f, 0x07, 0xdc},
	{0x77, 0x1f, 0x07, 0xdc}, {0x8f, 0x27, 0x07, 0xdc},
	{0x9f, 0x2f, 0x07, 0xdc}, {0xaf, 0x3f, 0x07, 0xdc},
	{0xbf, 0x47, 0x07, 0xdc}, {0xc7, 0x47, 0x07, 0xdc},
	{0xdf, 0x4f, 0x07, 0xdc}, {0xdf, 0x57, 0x07, 0xdc},
	{0xdf, 0x57, 0x07, 0xdc}, {0xd7, 0x5f, 0x07, 0xdc},
	{0xd7, 0x67, 0x0f, 0xdc}, {0xcf, 0x6f, 0x0f, 0xdc},
	{0xcf, 0x77, 0x0f, 0xdc}, {0xcf, 0x7f, 0x0f, 0xdc},
	{0xcf, 0x87, 0x17, 0xdc}, {0xc7, 0x87, 0x17, 0xdc},
	{0xc7, 0x8f, 0x17, 0xdc}, {0xc7, 0x97, 0x1f, 0xdc},
	{0xbf, 0x9f, 0x1f, 0xdc}, {0xbf, 0x9f, 0x1f, 0xdc},
	{0xbf, 0xa7, 0x27, 0xdc}, {0xbf, 0xa7, 0x27, 0xdc},
	{0xbf, 0xaf, 0x2f, 0xdc}, {0xb7, 0xaf, 0x2f, 0xdc},
	{0xb7, 0xb7, 0x2f, 0xdc}, {0xb7, 0xb7, 0x37, 0xdc},
	{0xcf, 0xcf, 0x6f, 0xdc}, {0xdf, 0xdf, 0x9f, 0xdc},
	{0xef, 0xef, 0xc7, 0xdc}, {0xff, 0xff, 0xff, 0xdc},
}

func mapColor(v int8) color.RGBA {
	if v < 0 || int(v) >= len(cmap) {
		return color.RGBA{0, 0, 0, 255}
	}

	return cmap[v]
}
