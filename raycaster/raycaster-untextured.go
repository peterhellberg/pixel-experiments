package main

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

var (
	fullscreen = false
	width      = 320
	height     = 200
	scale      = 3.0

	pos   = pixel.V(18.0, 9.5)
	dir   = pixel.V(-1.0, 0.0)
	plane = pixel.V(0.0, 0.66)
)

var world = [24][24]int{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 2, 2, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, 0, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 2, 0, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 0, 0, 0, 5, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
}

func getColor(x, y int) color.RGBA {
	switch world[x][y] {
	case 0:
		return color.RGBA{64, 64, 64, 255}
	case 1:
		return color.RGBA{244, 115, 33, 255}
	case 2:
		return color.RGBA{54, 124, 43, 255}
	case 3:
		return color.RGBA{0, 125, 198, 255}
	case 4:
		return color.RGBA{255, 255, 255, 255}
	default:
		return color.RGBA{255, 194, 32, 255}
	}
}

func frame() *image.RGBA {
	m := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.Draw(m, image.Rect(0, 0, width, height/2), &image.Uniform{color.RGBA{192, 192, 192, 255}}, image.ZP, draw.Src)
	draw.Draw(m, image.Rect(0, height/2, width, height), &image.Uniform{color.RGBA{64, 64, 64, 255}}, image.ZP, draw.Src)

	for x := 0; x < width; x++ {
		var (
			sideDist     pixel.Vec
			perpWallDist float64
			stepX, stepY int
			hit, side    int

			rayPos, worldX, worldY = pos, int(pos.X), int(pos.Y)

			cameraX = 2*float64(x)/float64(width) - 1

			rayDir = pixel.V(
				dir.X+plane.X*cameraX,
				dir.Y+plane.Y*cameraX,
			)

			deltaDist = pixel.V(
				math.Sqrt(1.0+(rayDir.Y*rayDir.Y)/(rayDir.X*rayDir.X)),
				math.Sqrt(1.0+(rayDir.X*rayDir.X)/(rayDir.Y*rayDir.Y)),
			)
		)

		if rayDir.X < 0 {
			stepX = -1
			sideDist.X = (rayPos.X - float64(worldX)) * deltaDist.X
		} else {
			stepX = 1
			sideDist.X = (float64(worldX) + 1.0 - rayPos.X) * deltaDist.X
		}

		if rayDir.Y < 0 {
			stepY = -1
			sideDist.Y = (rayPos.Y - float64(worldY)) * deltaDist.Y
		} else {
			stepY = 1
			sideDist.Y = (float64(worldY) + 1.0 - rayPos.Y) * deltaDist.Y
		}

		for hit == 0 {
			if sideDist.X < sideDist.Y {
				sideDist.X += deltaDist.X
				worldX += stepX
				side = 0
			} else {
				sideDist.Y += deltaDist.Y
				worldY += stepY
				side = 1
			}

			if world[worldX][worldY] > 0 {
				hit = 1
			}
		}

		if side == 0 {
			perpWallDist = (float64(worldX) - rayPos.X + (1-float64(stepX))/2) / rayDir.X
		} else {
			perpWallDist = (float64(worldY) - rayPos.Y + (1-float64(stepY))/2) / rayDir.Y
		}

		lineHeight := int(float64(height) / perpWallDist)

		drawStart := -lineHeight/2 + height/2
		if drawStart < 0 {
			drawStart = 0
		}

		drawEnd := lineHeight/2 + height/2
		if drawEnd >= height {
			drawEnd = height - 1
		}

		c := getColor(worldX, worldY)

		if side == 1 {
			c.R = c.R / 2
			c.G = c.G / 2
			c.B = c.B / 2
		}

		for y := drawStart; y < drawEnd; y++ {
			if y+1 == drawEnd {
				m.Set(x, y, color.RGBA{32, 32, 32, 255})
			} else {
				m.Set(x, y, c)
			}
		}
	}

	for x, row := range world {
		for y, _ := range row {
			m.Set(x, y, getColor(x, y))
		}
	}

	m.Set(int(pos.X), int(pos.Y), color.RGBA{255, 0, 0, 255})

	return m
}

func run() {
	cfg := pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, float64(width)*scale, float64(height)*scale),
		VSync:       true,
		Undecorated: true,
	}

	if fullscreen {
		cfg.Monitor = pixelgl.PrimaryMonitor()
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	c := win.Bounds().Center()

	last := time.Now()

	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			return
		}

		dt := time.Since(last).Seconds()
		last = time.Now()

		if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW) {
			moveForward(5 * dt)
		}

		if win.Pressed(pixelgl.KeyA) {
			moveLeft(5 * dt)
		}

		if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS) {
			moveBackwards(5 * dt)
		}

		if win.Pressed(pixelgl.KeyD) {
			moveRight(5 * dt)
		}

		if win.Pressed(pixelgl.KeyRight) {
			turnRight(2 * dt)
		}

		if win.Pressed(pixelgl.KeyLeft) {
			turnLeft(2 * dt)
		}

		p := pixel.PictureDataFromImage(frame())
		s := pixel.NewSprite(p, p.Bounds())

		s.Draw(win, pixel.IM.Moved(c).Scaled(c, scale))

		win.Update()
	}
}

func moveForward(s float64) {
	if world[int(pos.X+dir.X*s)][int(pos.Y)] == 0 {
		pos.X += dir.X * s
	}

	if world[int(pos.X)][int(pos.Y+dir.Y*s)] == 0 {
		pos.Y += dir.Y * s
	}
}

func moveLeft(s float64) {
	if world[int(pos.X-plane.X*s)][int(pos.Y)] == 0 {
		pos.X -= plane.X * s
	}

	if world[int(pos.X)][int(pos.Y-plane.Y*s)] == 0 {
		pos.Y -= plane.Y * s
	}
}

func moveBackwards(s float64) {
	if world[int(pos.X-dir.X*s)][int(pos.Y)] == 0 {
		pos.X -= dir.X * s
	}

	if world[int(pos.X)][int(pos.Y-dir.Y*s)] == 0 {
		pos.Y -= dir.Y * s
	}
}

func moveRight(s float64) {
	if world[int(pos.X+plane.X*s)][int(pos.Y)] == 0 {
		pos.X += plane.X * s
	}

	if world[int(pos.X)][int(pos.Y+plane.Y*s)] == 0 {
		pos.Y += plane.Y * s
	}
}

func turnRight(s float64) {
	dir.Y = dir.X*math.Sin(-s) + dir.Y*math.Cos(-s)
	dir.X = dir.X*math.Cos(-s) - dir.Y*math.Sin(-s)

	plane.Y = plane.X*math.Sin(-s) + plane.Y*math.Cos(-s)
	plane.X = plane.X*math.Cos(-s) - plane.Y*math.Sin(-s)
}

func turnLeft(s float64) {
	dir.Y = dir.X*math.Sin(s) + dir.Y*math.Cos(s)
	dir.X = dir.X*math.Cos(s) - dir.Y*math.Sin(s)

	plane.Y = plane.X*math.Sin(s) + plane.Y*math.Cos(s)
	plane.X = plane.X*math.Cos(s) - plane.Y*math.Sin(s)
}

func main() {
	flag.BoolVar(&fullscreen, "f", fullscreen, "fullscreen")
	flag.IntVar(&width, "w", width, "width")
	flag.IntVar(&height, "h", height, "height")
	flag.Float64Var(&scale, "s", scale, "scale")
	flag.Parse()

	pixelgl.Run(run)
}
