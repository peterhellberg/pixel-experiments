package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

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

			for !win.Closed() {
				p.input(win)
				p.draw(win)
				p.update(win)
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
	win.Clear(color.NRGBA{16, 16, 16, 15})

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

func (p *pong) update(win *pixelgl.Window) {
	switch {
	case p.left.Contains(p.ball):
		p.velocity.X = -p.velocity.X
		p.ball.X = 24
	case p.right.Contains(p.ball):
		p.velocity.X = -p.velocity.X
		p.ball.X = p.Max.X - 24
	case p.ball.Y < 16, p.ball.Y > p.Max.Y-16:
		p.velocity.Y = -p.velocity.Y
	case p.ball.X < 8 || p.ball.X > p.Max.X+8:
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

	p.ball = p.ball.Add(p.velocity)

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
