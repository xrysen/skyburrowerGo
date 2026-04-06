package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type HUD struct {
	background  *ebiten.Image
	heartImg    *ebiten.Image
	font        *BitmapFont
	carrotEmpty *ebiten.Image
	carrotFull  *ebiten.Image
	carrotScale float64
	carrotGap   float64
	carrotTopY  float64
}

func NewHUD(background, heartImg, carrotEmpty, carrotFull *ebiten.Image, font *BitmapFont) *HUD {
	return &HUD{
		heartImg:    heartImg,
		background:  background,
		font:        font,
		carrotEmpty: carrotEmpty,
		carrotFull:  carrotFull,
		carrotScale: 2,
		carrotGap:   4,
		carrotTopY:  6,
	}
}

func (h *HUD) Draw(screen *ebiten.Image, health, maxHealth, coins int, runCarrotMask uint8) {

	bgOp := &ebiten.DrawImageOptions{}
	bgOp.GeoM.Translate(10, 10)
	screen.DrawImage(h.background, bgOp)

	h.drawCarrotProgress(screen, runCarrotMask)

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
	h.font.DrawText(screen, coinText, 60, 49, 1.0)
}

func (h *HUD) drawCarrotProgress(screen *ebiten.Image, runCarrotMask uint8) {
	if h.carrotEmpty == nil || h.carrotFull == nil {
		return
	}
	cw := float64(h.carrotEmpty.Bounds().Dx()) * h.carrotScale
	gap := h.carrotGap
	rowW := float64(CarrotsPerLevel)*cw + float64(CarrotsPerLevel-1)*gap
	startX := (float64(ScreenWidth) - rowW) / 2
	y := h.carrotTopY

	for i := 0; i < CarrotsPerLevel; i++ {
		img := h.carrotEmpty
		if (runCarrotMask>>i)&1 != 0 {
			img = h.carrotFull
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(h.carrotScale, h.carrotScale)
		op.GeoM.Translate(startX+float64(i)*(cw+gap), y)
		screen.DrawImage(img, op)
	}
}
