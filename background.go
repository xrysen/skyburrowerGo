package main

import "github.com/hajimehoshi/ebiten/v2"

type Layer struct {
	img   *ebiten.Image
	speed float64
	x     float64
}

type Background struct {
	layers []*Layer
}

func (b *Background) Update() {
	for _, l := range b.layers {
		l.x -= l.speed
		w := float64(l.img.Bounds().Dx())
		if l.x < -w {
			l.x = 0
		}
	}
}

func (b *Background) Draw(screen *ebiten.Image, layerIndex int) {
	l := b.layers[layerIndex]
	w := float64(l.img.Bounds().Dx())

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(l.x, 0)
	screen.DrawImage(l.img, op)

	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(l.x+w, 0)
	screen.DrawImage(l.img, op2)
}
