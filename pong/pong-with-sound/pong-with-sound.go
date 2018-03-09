package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/mattetti/filebuffer"
)

var (
	beeep    beep.StreamSeekCloser
	plop     beep.StreamSeekCloser
	peeeeeep beep.StreamSeekCloser
)

func init() {
	var (
		err    error
		format beep.Format
	)

	plop, _, err = mp3.Decode(plopBuffer)
	if err != nil {
		panic(err)
	}

	peeeeeep, _, err = mp3.Decode(peeeeeepBuffer)
	if err != nil {
		panic(err)
	}

	beeep, format, err = mp3.Decode(beeepBuffer)
	if err != nil {
		panic(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/35))
}

func main() {
	pixelgl.Run(func() {
		p := newPong(pixel.R(0, 0, 858, 525))

		win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
			Bounds:      p.Rect,
			Undecorated: true,
			VSync:       true,
		})

		if err == nil {
			win.SetSmooth(false)

			last := time.Now()

			for !win.Closed() {
				dt := time.Since(last).Seconds()
				last = time.Now()

				p.update(win, dt)
				p.input(win)
				p.draw(win)
			}
		}
	})
}

type pong struct {
	pixel.Rect
	ball     pixel.Vec
	velocity pixel.Vec
	left     *player
	right    *player
}

func newPong(r pixel.Rect) *pong {
	return &pong{
		r, r.Center(), pixel.V(5, (rand.Float64()*2)-1*3),
		&player{
			pos: pixel.V(16, r.Max.Y/2),
			txt: text.New(pixel.V(r.Max.X/2-320, r.Max.Y-112), text.Atlas7x13),
		},
		&player{
			pos: pixel.V(r.Max.X-16, r.Max.Y/2),
			txt: text.New(pixel.V(r.Max.X/2+128, r.Max.Y-112), text.Atlas7x13),
		},
	}
}

func (p *pong) input(win *pixelgl.Window) {
	win.SetClosed(win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ))

	if win.Pressed(pixelgl.KeyW) && p.left.pos.Y < p.Max.Y-50 {
		p.left.pos.Y += 10
	}
	if win.Pressed(pixelgl.KeyS) && p.left.pos.Y > 50 {
		p.left.pos.Y -= 10
	}
	if win.Pressed(pixelgl.KeyUp) && p.right.pos.Y < p.Max.Y-50 {
		p.right.pos.Y += 10
	}
	if win.Pressed(pixelgl.KeyDown) && p.right.pos.Y > 50 {
		p.right.pos.Y -= 10
	}
	if win.Pressed(pixelgl.KeyR) {
		p.ball, p.left.score, p.right.score = p.Center(), 0, 0
	}
}

func (p *pong) draw(win *pixelgl.Window) {
	win.Clear(color.NRGBA{16, 16, 16, 255})

	imd := imdraw.New(nil)

	imd.Color = color.NRGBA{18, 18, 18, 255}
	imd.Push(pixel.V(p.Max.X/2, 0), pixel.V(p.Max.X/2, p.Max.Y))
	imd.Line(16)

	imd.Color = color.White
	imd.Push(p.ball.Add(pixel.V(-8, -8)), p.ball.Add(pixel.V(8, 8)))
	imd.Rectangle(0)

	p.right.draw(imd, win)
	p.left.draw(imd, win)

	imd.Draw(win)
}

func (p *pong) update(win *pixelgl.Window, dt float64) {
	switch {
	case p.left.Contains(p.ball):
		go func() {
			speaker.Play(beeep)
			beeep.Seek(0)
		}()

		p.velocity.X = -p.velocity.X
		p.ball.X = 24
	case p.right.Contains(p.ball):
		go func() {
			speaker.Play(beeep)
			beeep.Seek(0)
		}()

		p.velocity.X = -p.velocity.X
		p.ball.X = p.Max.X - 24
	case p.ball.Y < 16, p.ball.Y > p.Max.Y-16:
		go func() {
			speaker.Play(plop)
			plop.Seek(0)
		}()

		p.velocity.Y = -p.velocity.Y

		if p.ball.Y < 16 {
			p.ball.Y = 16
		} else {
			p.ball.Y = p.Max.Y - 16
		}
	case p.ball.X < 8 || p.ball.X > p.Max.X+8:
		go func() {
			speaker.Play(peeeeeep)
			peeeeeep.Seek(0)
		}()

		if p.ball.X < 8 {
			p.right.score++
		} else {
			p.left.score++
		}

		p.velocity = pixel.V(-8, (rand.Float64()-0.5)*16)

		if rand.Float64() > 0.5 {
			p.velocity.X = 8
		}

		p.ball = p.Center()
	}

	p.ball = p.ball.Add(p.velocity.Scaled(dt * 60))

	p.left.update()
	p.right.update()

	win.Update()
}

type player struct {
	pixel.Rect
	pos   pixel.Vec
	txt   *text.Text
	score int
}

func (p *player) draw(imd *imdraw.IMDraw, t pixel.Target) {
	imd.Color = color.White
	imd.Push(p.Min, p.Max)
	imd.Rectangle(0)
	p.txt.Draw(t, pixel.IM.Scaled(p.txt.Orig, 8))
}

func (p *player) update() {
	p.txt.Clear()
	fmt.Fprintf(p.txt, "% 3d", p.score)
	p.Rect = pixel.Rect{p.pos.Add(pixel.V(-8, -50)), p.pos.Add(pixel.V(8, 50))}
}

var (
	beeepBuffer    = filebuffer.New([]byte{255, 251, 24, 196, 0, 0, 5, 124, 29, 95, 148, 97, 128, 1, 22, 165, 236, 195, 32, 208, 0, 128, 0, 0, 0, 8, 70, 4, 103, 179, 201, 147, 76, 192, 24, 12, 153, 52, 193, 240, 124, 16, 4, 3, 16, 112, 16, 12, 97, 250, 131, 31, 131, 225, 255, 193, 7, 127, 206, 126, 15, 255, 254, 80, 48, 137, 118, 97, 88, 95, 1, 63, 137, 216, 79, 192, 215, 19, 226, 70, 90, 53, 161, 42, 184, 112, 46, 69, 104, 164, 174, 78, 40, 228, 153, 236, 105, 79, 157, 38, 37, 212, 65, 245, 18, 253, 20, 219, 242, 191, 232, 159, 235, 45, 243, 135, 245, 26, 122, 103, 247, 53, 243, 173, 214, 105, 234, 111, 169, 31, 58, 223, 155, 121, 213, 104, 0, 133, 65, 34, 50, 2, 222, 244, 255, 251, 24, 196, 4, 128, 7, 221, 31, 99, 156, 211, 128, 0, 243, 35, 174, 240, 20, 168, 150, 75, 161, 37, 0, 20, 22, 12, 9, 171, 122, 116, 73, 163, 160, 26, 5, 145, 74, 175, 250, 168, 132, 44, 93, 255, 252, 192, 88, 19, 28, 159, 254, 62, 37, 127, 254, 120, 64, 79, 255, 242, 0, 137, 6, 251, 191, 212, 78, 81, 63, 254, 105, 2, 110, 43, 66, 90, 13, 194, 160, 12, 106, 166, 10, 3, 58, 142, 14, 136, 147, 190, 130, 54, 113, 166, 171, 87, 234, 61, 28, 255, 254, 68, 223, 253, 178, 16, 188, 71, 221, 247, 67, 146, 114, 136, 32, 225, 172, 105, 191, 85, 52, 213, 64, 174, 199, 33, 201, 244, 244, 27, 255, 255, 30, 186, 87, 80, 4, 89, 96, 86, 43, 177, 255, 251, 24, 196, 4, 0, 7, 201, 27, 103, 224, 189, 164, 128, 236, 155, 177, 62, 146, 112, 6, 0, 3, 55, 44, 158, 192, 211, 170, 58, 43, 181, 9, 186, 148, 143, 24, 197, 7, 89, 130, 159, 253, 101, 68, 159, 255, 169, 76, 48, 193, 26, 223, 254, 161, 157, 90, 247, 82, 234, 213, 51, 20, 126, 255, 189, 98, 76, 175, 255, 206, 22, 255, 255, 36, 154, 16, 137, 90, 21, 26, 182, 219, 0, 12, 153, 1, 223, 125, 113, 55, 81, 136, 101, 154, 78, 186, 186, 129, 50, 40, 98, 43, 59, 169, 139, 65, 247, 11, 134, 159, 255, 241, 89, 181, 255, 245, 36, 237, 186, 127, 69, 11, 142, 167, 253, 21, 40, 45, 46, 199, 254, 175, 243, 21, 196, 150, 83, 126, 14, 80, 157, 32, 208, 255, 251, 24, 196, 4, 128, 7, 116, 253, 98, 24, 245, 0, 0, 0, 0, 52, 131, 128, 0, 4, 191, 100, 136, 204, 251, 242, 112, 163, 81, 87, 2, 194, 120, 246, 190, 129, 114, 55, 52, 226, 103, 243, 17, 134, 99, 207, 53, 147, 254, 61, 39, 152, 100, 211, 83, 255, 204, 41, 231, 151, 206, 245, 55, 255, 232, 97, 27, 115, 139, 124, 170, 76, 65, 77, 69, 51, 46, 57, 57, 46, 53, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 170, 174, 170})
	peeeeeepBuffer = filebuffer.New([]byte{255, 251, 24, 196, 0, 0, 4, 212, 27, 97, 52, 97, 128, 1, 49, 149, 110, 103, 30, 112, 0, 128, 1, 0, 54, 176, 47, 123, 65, 2, 24, 120, 0, 0, 66, 29, 7, 3, 2, 0, 65, 210, 128, 129, 207, 40, 15, 191, 203, 131, 242, 127, 224, 248, 127, 235, 127, 252, 16, 71, 250, 127, 250, 127, 127, 250, 64, 0, 5, 152, 177, 245, 144, 255, 25, 166, 105, 168, 218, 125, 43, 98, 235, 63, 10, 135, 82, 35, 153, 20, 150, 132, 3, 83, 212, 228, 23, 171, 176, 67, 160, 80, 161, 215, 52, 145, 29, 226, 41, 61, 226, 49, 201, 117, 74, 26, 151, 65, 37, 80, 230, 92, 101, 140, 18, 92, 138, 52, 18, 17, 173, 194, 29, 31, 187, 232, 223, 67, 45, 171, 182, 243, 96, 1, 255, 251, 24, 196, 4, 0, 7, 208, 231, 137, 188, 51, 128, 48, 242, 29, 45, 188, 55, 168, 200, 5, 176, 184, 228, 12, 132, 128, 151, 62, 231, 237, 228, 132, 65, 161, 1, 179, 186, 237, 244, 21, 11, 78, 255, 254, 41, 143, 147, 142, 50, 19, 52, 57, 135, 152, 229, 81, 56, 241, 219, 45, 205, 87, 247, 133, 18, 190, 211, 189, 40, 55, 49, 253, 190, 213, 134, 48, 102, 33, 24, 32, 232, 128, 2, 138, 161, 94, 80, 66, 113, 42, 169, 197, 212, 203, 215, 245, 5, 0, 222, 74, 199, 143, 171, 175, 229, 11, 127, 171, 254, 50, 149, 45, 61, 99, 178, 106, 52, 164, 71, 19, 39, 239, 54, 167, 83, 104, 221, 175, 250, 254, 164, 214, 119, 235, 158, 67, 171, 171, 54, 171, 100, 0, 255, 251, 24, 196, 4, 0, 7, 188, 233, 123, 164, 152, 232, 80, 245, 157, 44, 240, 87, 156, 144, 175, 100, 155, 212, 196, 0, 176, 186, 129, 38, 38, 209, 163, 163, 244, 5, 228, 66, 115, 74, 138, 217, 119, 244, 148, 6, 35, 126, 191, 196, 241, 210, 108, 123, 15, 141, 156, 37, 186, 41, 134, 195, 173, 233, 95, 225, 87, 239, 189, 191, 139, 45, 103, 215, 30, 162, 178, 6, 19, 176, 6, 12, 169, 243, 48, 131, 236, 226, 66, 99, 49, 205, 122, 126, 84, 5, 8, 131, 15, 48, 245, 219, 248, 49, 62, 245, 254, 40, 174, 210, 173, 30, 218, 177, 64, 35, 249, 217, 106, 111, 17, 11, 166, 109, 90, 38, 218, 170, 80, 69, 181, 42, 187, 191, 66, 149, 85, 193, 97, 162, 160, 34, 32, 255, 251, 24, 196, 3, 128, 7, 128, 231, 105, 160, 189, 66, 192, 253, 33, 109, 116, 22, 40, 144, 0, 73, 105, 121, 17, 96, 201, 64, 77, 26, 198, 210, 78, 198, 159, 186, 1, 104, 139, 74, 205, 87, 87, 253, 69, 31, 235, 252, 87, 177, 50, 34, 203, 71, 168, 171, 83, 197, 198, 219, 215, 45, 237, 25, 50, 125, 103, 255, 29, 87, 242, 41, 236, 146, 193, 108, 2, 144, 107, 128, 0, 67, 154, 27, 143, 136, 130, 211, 235, 156, 148, 54, 191, 195, 114, 67, 139, 77, 174, 159, 149, 28, 79, 235, 252, 88, 185, 247, 99, 9, 98, 179, 171, 202, 40, 184, 79, 190, 211, 104, 212, 241, 45, 172, 156, 188, 133, 230, 214, 84, 84, 174, 254, 255, 148, 250, 21, 114, 32, 218, 130, 18, 162, 255, 251, 24, 196, 3, 0, 7, 84, 233, 115, 161, 188, 228, 144, 249, 157, 47, 52, 132, 156, 210, 64, 1, 142, 226, 16, 239, 33, 101, 120, 213, 57, 183, 145, 119, 110, 163, 66, 66, 41, 232, 113, 189, 191, 19, 150, 79, 173, 91, 179, 66, 121, 131, 202, 86, 138, 130, 41, 198, 155, 66, 46, 31, 111, 74, 127, 34, 246, 244, 167, 232, 45, 255, 84, 214, 29, 82, 53, 37, 91, 32, 9, 105, 121, 71, 5, 1, 96, 232, 18, 40, 196, 178, 212, 253, 5, 35, 114, 229, 208, 250, 239, 248, 132, 183, 233, 183, 226, 201, 66, 179, 101, 6, 78, 3, 158, 78, 62, 133, 67, 83, 253, 232, 244, 218, 32, 218, 148, 107, 178, 122, 11, 236, 79, 238, 235, 150, 65, 35, 102, 194, 164, 100, 1, 255, 251, 24, 196, 4, 0, 7, 168, 231, 121, 164, 148, 235, 208, 249, 157, 46, 244, 103, 156, 146, 110, 166, 114, 236, 168, 63, 34, 90, 105, 103, 154, 159, 3, 0, 129, 7, 67, 84, 121, 18, 200, 184, 214, 253, 127, 224, 210, 137, 69, 114, 232, 85, 204, 238, 10, 143, 89, 141, 169, 89, 131, 40, 172, 235, 7, 236, 191, 120, 239, 221, 70, 223, 228, 166, 144, 247, 16, 172, 168, 216, 0, 64, 216, 37, 7, 200, 127, 147, 20, 163, 98, 38, 149, 50, 191, 198, 1, 217, 18, 202, 93, 85, 230, 211, 196, 255, 218, 141, 248, 147, 67, 103, 74, 187, 137, 169, 153, 19, 134, 39, 239, 108, 189, 247, 133, 79, 181, 42, 210, 153, 251, 168, 115, 33, 63, 90, 186, 83, 35, 134, 164, 163, 132, 255, 251, 24, 196, 3, 128, 7, 188, 231, 121, 164, 148, 233, 80, 238, 156, 236, 244, 156, 28, 200, 0, 164, 224, 174, 89, 129, 112, 16, 116, 13, 126, 88, 141, 77, 208, 62, 40, 7, 41, 202, 247, 191, 226, 23, 255, 79, 227, 105, 198, 77, 99, 90, 63, 49, 204, 179, 3, 45, 165, 42, 246, 218, 40, 46, 167, 109, 54, 75, 58, 113, 227, 194, 177, 223, 218, 237, 98, 32, 69, 32, 68, 64, 3, 191, 82, 186, 108, 96, 148, 144, 60, 187, 58, 111, 79, 199, 192, 80, 225, 7, 60, 197, 188, 207, 196, 37, 237, 237, 79, 226, 46, 109, 114, 200, 40, 163, 79, 137, 131, 51, 253, 235, 252, 70, 107, 211, 121, 95, 74, 139, 151, 202, 242, 126, 202, 154, 17, 35, 130, 146, 107, 96, 0, 255, 251, 24, 196, 4, 0, 7, 176, 231, 109, 180, 179, 128, 0, 239, 150, 172, 163, 32, 112, 0, 247, 114, 202, 121, 88, 116, 38, 111, 241, 93, 99, 213, 55, 152, 18, 150, 56, 110, 76, 179, 93, 127, 138, 75, 254, 244, 254, 35, 79, 45, 237, 31, 116, 241, 50, 187, 58, 107, 86, 75, 94, 23, 103, 116, 125, 17, 171, 189, 72, 4, 92, 234, 121, 30, 0, 209, 64, 40, 128, 0, 10, 227, 24, 138, 88, 52, 198, 209, 137, 142, 32, 130, 71, 140, 177, 96, 156, 231, 192, 24, 29, 22, 53, 124, 71, 97, 28, 6, 61, 83, 200, 169, 179, 130, 31, 255, 84, 49, 144, 74, 115, 62, 191, 250, 202, 153, 113, 19, 191, 244, 14, 80, 242, 181, 76, 65, 77, 69, 51, 46, 57, 57, 46, 255, 251, 24, 196, 4, 131, 192, 0, 1, 164, 28, 0, 0, 32, 0, 0, 52, 128, 0, 0, 4, 53, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85})
	plopBuffer     = filebuffer.New([]byte{255, 251, 24, 196, 0, 0, 6, 177, 24, 228, 20, 18, 128, 8, 244, 38, 38, 11, 2, 32, 0, 255, 255, 140, 132, 109, 78, 119, 35, 127, 255, 157, 255, 249, 206, 115, 255, 255, 243, 156, 255, 252, 231, 57, 200, 66, 55, 83, 135, 3, 130, 132, 212, 225, 240, 0, 56, 194, 97, 240, 248, 124, 239, 144, 132, 35, 40, 152, 112, 56, 28, 20, 32, 0, 129, 70, 20, 0, 0, 240, 251, 23, 196, 253, 31, 245, 69, 251, 108, 165, 79, 235, 62, 86, 52, 199, 255, 116, 98, 18, 72, 97, 65, 76, 12, 4, 40, 15, 254, 236, 119, 57, 226, 209, 136, 20, 96, 197, 12, 40, 48, 16, 127, 255, 185, 24, 132, 146, 254, 25, 194, 152, 162, 74, 2, 69, 240, 2, 248, 177, 188, 56, 191, 42, 255, 251, 24, 196, 4, 0, 7, 177, 49, 36, 24, 35, 128, 0, 234, 36, 165, 151, 2, 32, 1, 79, 244, 38, 159, 228, 11, 141, 8, 127, 248, 225, 3, 200, 21, 27, 255, 255, 145, 27, 184, 220, 144, 144, 50, 34, 127, 255, 248, 144, 40, 22, 13, 194, 65, 28, 39, 20, 131, 180, 255, 255, 255, 193, 216, 152, 28, 9, 64, 224, 176, 56, 26, 131, 177, 4, 6, 16, 0, 139, 89, 117, 60, 149, 115, 156, 182, 238, 70, 230, 253, 26, 134, 67, 47, 231, 59, 144, 156, 178, 213, 191, 232, 192, 96, 103, 0, 0, 52, 198, 67, 22, 165, 255, 216, 231, 156, 231, 0, 0, 34, 149, 21, 3, 33, 159, 255, 254, 6, 6, 112, 13, 64, 241, 64, 129, 159, 224, 73, 131, 112, 120, 88, 96, 255, 251, 24, 196, 5, 128, 7, 136, 227, 80, 25, 22, 128, 1, 6, 140, 173, 151, 30, 144, 0, 156, 60, 3, 88, 72, 255, 30, 97, 204, 52, 255, 64, 192, 148, 52, 255, 240, 189, 142, 114, 92, 57, 131, 208, 148, 255, 255, 37, 203, 230, 230, 244, 201, 66, 231, 255, 255, 154, 23, 203, 230, 229, 194, 224, 13, 255, 255, 57, 135, 228, 155, 73, 36, 1, 132, 19, 141, 197, 27, 33, 10, 124, 189, 218, 163, 200, 161, 172, 126, 229, 168, 26, 24, 146, 211, 210, 251, 68, 27, 229, 150, 127, 213, 66, 53, 53, 106, 21, 241, 220, 250, 158, 1, 29, 120, 179, 220, 196, 92, 29, 21, 185, 165, 103, 201, 169, 55, 181, 11, 247, 48, 135, 106, 255, 214, 4, 22, 11, 6, 137, 173, 12, 8, 255, 251, 24, 196, 4, 1, 7, 136, 99, 70, 60, 147, 0, 16, 212, 36, 31, 84, 17, 155, 193, 129, 38, 86, 69, 43, 140, 106, 89, 71, 37, 90, 104, 40, 6, 165, 28, 106, 167, 157, 121, 54, 113, 129, 128, 72, 209, 42, 42, 26, 193, 161, 224, 168, 240, 87, 196, 176, 104, 96, 119, 255, 42, 116, 74, 226, 192, 215, 81, 238, 10, 245, 134, 214, 16, 35, 35, 115, 173, 68, 121, 190, 190, 231, 189, 220, 211, 199, 190, 17, 250, 247, 149, 204, 185, 127, 254, 101, 34, 147, 68, 53, 205, 97, 176, 77, 166, 108, 155, 103, 233, 206, 74, 242, 208, 170, 44, 8, 73, 19, 80, 57, 109, 154, 69, 41, 234, 156, 0, 64, 1, 59, 166, 103, 55, 108, 171, 109, 210, 200, 160, 229, 174, 100, 255, 251, 24, 196, 8, 128, 71, 109, 32, 244, 192, 140, 222, 9, 12, 36, 93, 148, 17, 39, 128, 243, 107, 122, 243, 223, 223, 206, 124, 238, 185, 121, 60, 136, 98, 141, 247, 65, 94, 137, 2, 24, 120, 194, 34, 115, 157, 52, 106, 142, 74, 65, 18, 7, 6, 4, 30, 72, 28, 210, 6, 16, 145, 128, 16, 18, 19, 147, 9, 4, 60, 132, 194, 40, 133, 33, 51, 113, 10, 66, 127, 254, 38, 16, 197, 248, 98, 124, 133, 87, 133, 127, 82, 229, 46, 105, 64, 74, 213, 37, 154, 244, 180, 163, 31, 82, 216, 199, 220, 82, 33, 103, 245, 136, 154, 85, 159, 72, 183, 8, 131, 75, 146, 238, 22, 21, 61, 84, 56, 138, 53, 4, 12, 140, 252, 140, 140, 200, 254, 102, 95, 255, 255, 254, 255, 251, 24, 196, 6, 131, 196, 201, 20, 158, 32, 4, 120, 168, 0, 0, 52, 128, 0, 0, 4, 204, 141, 90, 201, 122, 202, 8, 24, 32, 86, 86, 10, 8, 24, 32, 64, 193, 58, 60, 176, 48, 80, 81, 3, 3, 21, 76, 65, 77, 69, 51, 46, 57, 57, 46, 53, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85, 85})
)
