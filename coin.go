package main

import (
	"image"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

type CoinSize int

const (
	SmallCoin CoinSize = iota
	BigCoin
)

type Coin struct {
	x, y         float64
	img          *ebiten.Image
	frameWidth   int
	frameHeight  int
	frameCounter int
	size         CoinSize
	value        int
	scale        float64
	collected    bool
}

func NewCoin(x, y float64, size CoinSize, img *ebiten.Image) *Coin {
	value := 1
	scale := 1.0
	if size == BigCoin {
		value = 5
		scale = 1.8
	}
	return &Coin{
		x:           x,
		y:           y,
		img:         img,
		frameWidth:  16,
		frameHeight: 16,
		size:        size,
		value:       value,
		scale:       scale,
	}
}

func (c *Coin) Update() {
	c.frameCounter++
	c.x -= 1.5
}

func (c *Coin) Draw(screen *ebiten.Image) {
	frame := (c.frameCounter / 8) % 7
	sx := frame * c.frameWidth
	rect := image.Rect(sx, 0, sx+c.frameWidth, c.frameHeight)
	subImg := c.img.SubImage(rect).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(c.scale, c.scale)
	op.GeoM.Translate(c.x, c.y)
	screen.DrawImage(subImg, op)
}

func ShouldSpawnBigCoin(playerLuck int, enemyDifficulty int) bool {
	chance := 2 + (playerLuck * 3) + (enemyDifficulty * 5)
	return rand.IntN(100) < int(chance)
}
