package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

var (
	left, up = float64(55), float64(-65)

	texts = [6][6]*text.Text{
		{newText(0, 125), newText(25, 125), newText(50, 125), newText(75, 125), newText(100, 125), newText(125, 125)},
		{newText(0, 100), newText(25, 100), newText(50, 100), newText(75, 100), newText(100, 100), newText(125, 100)},
		{newText(0, 75), newText(25, 75), newText(50, 75), newText(75, 75), newText(100, 75), newText(125, 75)},
		{newText(0, 50), newText(25, 50), newText(50, 50), newText(75, 50), newText(100, 50), newText(125, 50)},
		{newText(0, 25), newText(25, 25), newText(50, 25), newText(75, 25), newText(100, 25), newText(125, 25)},
		{newText(0, 0), newText(25, 0), newText(50, 0), newText(75, 0), newText(100, 0), newText(125, 0)},
	}

	colors = []color.RGBA{
		color.RGBA{255, 32, 32, 200},  // RED
		color.RGBA{32, 255, 32, 200},  // GREEN
		color.RGBA{32, 32, 255, 200},  // BLUE
		color.RGBA{32, 255, 255, 200}, // CYAN
		color.RGBA{255, 32, 255, 200}, // MAGENTA
	}
)

func run() {
	start := time.Now()

	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, float64(880), float64(400)),
		VSync:       true,
		Undecorated: true,
	})
	if err != nil {
		panic(err)
	}

	win.SetSmooth(false)
	win.SetMatrix(pixel.IM.Moved(win.Bounds().Center()).Scaled(win.Bounds().Center(), 2))

	s := new(state)

	go func() {
		for range time.Tick(96 * time.Millisecond) {
			s.Update(time.Since(start))
		}
	}()

	txt := text.New(pixel.V(-250+left, 50+up), text.Atlas7x13)
	txt.WriteString("THE MAGIC HAND OF")

	for !win.Closed() {
		win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))
		win.Clear(color.RGBA{24, 24, 24, 255})

		txt.Draw(win, pixel.IM.ScaledXY(txt.Orig, pixel.V(2, 1)))

		for y := 0; y < 6; y++ {
			for x := 0; x < 6; x++ {
				t := texts[y][x]
				t.Clear()
				t.WriteString(string(s.field[y][x]))
				t.Draw(win, pixel.IM.ScaledXY(t.Orig, pixel.V(2, 1)))
			}
		}

		if win.JustPressed(pixelgl.KeySpace) {
			start = time.Now()
		}

		win.Update()
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	pixelgl.Run(run)
}

type field [6]row

type row [6]rune

func (r row) String() string {
	return fmt.Sprintf("%s %s %s %s %s %s",
		string(r[0]),
		string(r[1]),
		string(r[2]),
		string(r[3]),
		string(r[4]),
		string(r[5]),
	)
}

type state struct {
	field field
}

func (s *state) Update(d time.Duration) bool {
	var letters string

	switch {
	case d < 1*time.Second:
		letters = "#"
	case d < 2*time.Second:
		letters = "ACEHN####"
	case d < 3*time.Second:
		letters = "ACEHN###"
	default:
		letters = "ACEHN##"

		if s.field[3].String() == "C H A N C E" {
			return true
		}
	}

	for y, row := range s.field {
		for x, _ := range row {
			if d > 3*time.Second && y == 3 {
				switch {
				case
					x == 0 && row[0] == 'C',
					x == 1 && row[0] == 'C' && row[1] == 'H',
					x == 2 && row[0] == 'C' && row[1] == 'H' && row[2] == 'A',
					x == 3 && row[0] == 'C' && row[1] == 'H' && row[2] == 'A' && row[3] == 'N',
					x == 4 && row[0] == 'C' && row[1] == 'H' && row[2] == 'A' && row[3] == 'N' && row[4] == 'C':
					continue
				}
			}

			s.field[y][x] = rune(letters[rand.Intn(len(letters))])

			if rand.Intn(3) == 0 {
				texts[y][x].Color = colors[rand.Intn(len(colors))]
			}
		}
	}

	return false
}

func (s *state) String() string {
	var f string

	for _, row := range s.field {
		f += fmt.Sprintf("%s %s %s %s %s %s\n",
			string(row[0]),
			string(row[1]),
			string(row[2]),
			string(row[3]),
			string(row[4]),
			string(row[5]),
		)
	}

	return f
}

func newText(x, y float64) *text.Text {
	return text.New(pixel.V(left+x, up+y), text.Atlas7x13)
}
