package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Flutternat struct {
	x, y           float64
	img            *ebiten.Image
	speed          float64
	health         int
	maxHealth      int
	frameWidth     int
	frameHeight    int
	frameCounter   int
	hitFlashTimer  int
	healthBarTimer int
}

func (f *Flutternat) Update(px, py float64, game *Game) {
	f.x -= f.speed
	f.frameCounter++

	if f.hitFlashTimer > 0 {
		f.hitFlashTimer--
	}

	if f.healthBarTimer > 0 {
		f.healthBarTimer--
	}
}

func (f *Flutternat) Draw(screen *ebiten.Image) {
	frame := (f.frameCounter / 8) % 4
	sx := frame * f.frameWidth
	rect := image.Rect(sx, 0, sx+f.frameWidth, f.frameHeight)
	subImg := f.img.SubImage(rect).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}

	if f.hitFlashTimer > 0 {
		op.ColorScale.ScaleWithColor(color.RGBA{255, 100, 100, 255})
	}

	op.GeoM.Translate(f.x, f.y)
	screen.DrawImage(subImg, op)

	if f.healthBarTimer > 0 {
		f.drawHealthBar(screen)
	}
}

func (f *Flutternat) drawHealthBar(screen *ebiten.Image) {
	barWidth := float64(f.frameWidth)
	barHeight := 4.0
	barX := f.x
	barY := f.y - 8

	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth), float32(barHeight), color.RGBA{255, 255, 255, 255}, false)

	healthPercent := float64(f.health) / float64(f.maxHealth)
	redWidth := barWidth * healthPercent
	vector.FillRect(screen, float32(barX), float32(barY), float32(redWidth), float32(barHeight), color.RGBA{255, 0, 0, 255}, false)
}

func (f *Flutternat) GetPosition() (x, y float64) {
	return f.x, f.y
}

func (f *Flutternat) GetBounds() (width, height float64) {
	return float64(f.frameWidth), float64(f.frameHeight)
}

func (f *Flutternat) TakeDamage(amount int) {
	f.health -= amount
	f.hitFlashTimer = 10
	f.healthBarTimer = 30
}

func (f *Flutternat) IsDead() bool {
	return f.health <= 0
}

func (f *Flutternat) OnDeath(game *Game) {
	size := SmallCoin
	if ShouldSpawnBigCoin(game.player.luck, 1) {
		size = BigCoin
	}
	coin := NewCoin(f.x, f.y, size, game.coinImg)
	game.coins = append(game.coins, coin)
	// TODO: if boss, set bossKilled to true
}
