package main

import "github.com/hajimehoshi/ebiten/v2"

type Bullet struct {
	x, y  float64
	img   *ebiten.Image
	speed float64
}

func NewBullet(x, y float64, img *ebiten.Image) *Bullet {
	return &Bullet{
		x:     x,
		y:     y,
		img:   img,
		speed: 7.0,
	}
}

func (b *Bullet) Update() {
	b.x += b.speed
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.x, b.y)
	screen.DrawImage(b.img, op)
}
