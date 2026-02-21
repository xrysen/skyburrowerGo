package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	x, y         float64
	img          *ebiten.Image
	frameWidth   int
	frameHeight  int
	frameCounter int
}

func NewPlayer(img *ebiten.Image) *Player {
	return &Player{
		x:            50,
		y:            100,
		img:          img,
		frameWidth:   80,
		frameHeight:  80,
		frameCounter: 0,
	}
}

func (p *Player) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		p.x += 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		p.x -= 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		p.y -= 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		p.y += 3
	}

	p.frameCounter++
}

func (p *Player) Draw(screen *ebiten.Image) {
	// Shows from 0 or 1
	frame := (p.frameCounter / 10) % 2

	sx := frame * p.frameWidth
	rect := image.Rect(sx, 0, sx+p.frameWidth, p.frameHeight)
	subImg := p.img.SubImage(rect).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.x, p.y)
	screen.DrawImage(subImg, op)
}
