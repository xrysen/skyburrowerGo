package main

import (
	"image"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type BitmapFont struct {
	img        *ebiten.Image
	charWidth  int
	charHeight int
}

func NewBitmapFont(img *ebiten.Image, charWidth, charHeight int) *BitmapFont {
	return &BitmapFont{
		img:        img,
		charWidth:  charWidth,
		charHeight: charHeight,
	}
}

func (f *BitmapFont) DrawText(screen *ebiten.Image, text string, x, y float64, scale float64) {
	charMap := " !\"\"$%*'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_~abcdefghijklmnopqrstuvwxyz"

	charsPerRow := 20
	baseScale := 0.5 * scale
	charSpacing := 9 * scale

	for i, char := range text {
		charIndex := strings.IndexRune(charMap, char)
		if charIndex == -1 {
			//charIndex = 0
			continue
		}

		col := charIndex % charsPerRow
		row := charIndex / charsPerRow

		srcX := 1 + col*19
		srcY := 1 + row*19

		rect := image.Rect(srcX, srcY, srcX+18, srcY+18)
		charImg := f.img.SubImage(rect).(*ebiten.Image)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(baseScale, baseScale)
		op.GeoM.Translate(x+float64(i)*charSpacing, y)
		screen.DrawImage(charImg, op)
	}
}
