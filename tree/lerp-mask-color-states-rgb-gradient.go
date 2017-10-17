package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/peterhellberg/gradient"
)

func main() {
	n := 128

	hg := gradient.NewHorizontal(n, 10, []gradient.Stop{
		{0.0, color.NRGBA{243, 119, 54, 255}},
		{0.5, color.NRGBA{3, 146, 207, 255}},
		{1.0, color.NRGBA{123, 192, 67, 255}},
	})

	saveImage(hg, "/tmp/horizontal-gradient.png")

	for i := 0; i < n; i++ {
		c := hg.At(i, 0).(color.NRGBA)

		fmt.Printf(`{`+
			`"depth":7,`+
			`"theta":27,`+
			`"angle":34,`+
			`"frac":0.839,`+
			`"length":93.518,`+
			`"mask":{`+
			`"R":%d,`+
			`"G":%d,`+
			`"B":%d,`+
			`"A":%d`+
			`},`+
			`"circles":true,`+
			`"light":true}`+"\n",
			c.R, c.G, c.B, c.A)
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

//RED: {"depth":7,"theta":27,"angle":34,"frac":0.839,"length":93.518,"mask":{"R":243,"G":119,"B":54,"A":55},"circles":true,"light":true}
//GREEN:   {"depth":7,"theta":27,"angle":34,"frac":0.839,"length":93.518,"mask":{"R":123,"G":192,"B":67,"A":55},"circles":true,"light":true}
