package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	colorful "github.com/lucasb-eyer/go-colorful"
)

func main() {
	keypoints := GradientTable{
		{colorful.MakeColor(color.NRGBA{243, 119, 54, 255}), 0.0},
		{colorful.MakeColor(color.NRGBA{123, 192, 67, 255}), 0.5},
		{colorful.MakeColor(color.NRGBA{255, 255, 255, 255}), 1.0},
	}

	for i := 0; i < 100; i++ {
		var (
			c       = keypoints.GetInterpolatedColorFor(float64(i+1) * 0.01)
			r, g, b = c.RGB255()
			d       = 1 + i/16
			t       = (-i * 2) + 50
			a       = i / 3
			f       = 0.34 + (float64(i) * 0.00455)
			l       = float64(i) * 1.45
		)

		fmt.Printf(`{`+
			`"depth":%d,`+
			`"theta":%d,`+
			`"angle":%d,`+
			`"frac":%f,`+
			`"length":%f,`+
			`"mask":{`+
			`"R":%d,`+
			`"G":%d,`+
			`"B":%d,`+
			`"A":%d`+
			`},`+
			`"circles":false,`+
			`"light":false}`+"\n",
			d, t, a, f, l, r, g, b, 255)
	}
}

func saveImage(m image.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, m)
}

type GradientTable []struct {
	Col colorful.Color
	Pos float64
}

func (self GradientTable) GetInterpolatedColorFor(t float64) colorful.Color {
	for i := 0; i < len(self)-1; i++ {
		c1 := self[i]
		c2 := self[i+1]
		if c1.Pos <= t && t <= c2.Pos {
			t := (t - c1.Pos) / (c2.Pos - c1.Pos)
			return c1.Col.BlendHcl(c2.Col, t).Clamped()
		}
	}

	return self[len(self)-1].Col
}
