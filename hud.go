package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type HUD struct {
	background *ebiten.Image
	heartImg   *ebiten.Image
	font       *BitmapFont
}

func NewHUD(background, heartImg *ebiten.Image, font *BitmapFont) *HUD {
	return &HUD{
		heartImg:   heartImg,
		background: background,
		font:       font,
	}
}

func (h *HUD) Draw(screen *ebiten.Image, health, maxHealth, coins int) {

	bgOp := &ebiten.DrawImageOptions{}
	bgOp.GeoM.Translate(10, 10)
	screen.DrawImage(h.background, bgOp)

	for i := 0; i < maxHealth; i++ {
		x := 0.0
		y := 0.0
		if i < 5 {
			x = float64(48 + (i * 15))
			y = 15.0
		} else {
			x = float64(48 + ((i - 5) * 15))
			y = 30.0
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)

		if i < health {
			screen.DrawImage(h.heartImg, op)
		} else {
			op.ColorScale.Scale(0.3, 0.3, 0.3, 0.5)
		}
	}
	coinText := fmt.Sprintf("%d", coins)
	h.font.DrawText(screen, coinText, 60, 49)
}
