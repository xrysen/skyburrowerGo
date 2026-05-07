package main

import (
	"image"
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const BlightmothType EnemyType = "blightmoth"

type Blightmoth struct {
	x, y           float64
	img            *ebiten.Image
	speed          float64
	frameCounter   int
	sineTimer      int
	shootTimer     int
	health         int
	maxHealth      int
	hitFlashTimer  int
	healthBarTimer int
}

func NewBlightmoth(x, y float64, img *ebiten.Image) *Blightmoth {
	return &Blightmoth{
		x:         x,
		y:         y,
		img:       img,
		speed:     1.2,
		health:    5,
		maxHealth: 5,
	}
}

func (b *Blightmoth) Update(px, py float64, game *Game) {
	b.x -= b.speed
	b.frameCounter++
	b.sineTimer++
	b.y += math.Sin(float64(b.sineTimer)*0.04) * 0.6

	b.shootTimer++
	if b.shootTimer >= 120 {
		b.shootTimer = 0
		b.fireSpread(px, py, game)
	}

	if b.hitFlashTimer > 0 {
		b.hitFlashTimer--
	}
	if b.healthBarTimer > 0 {
		b.healthBarTimer--
	}
}

func (b *Blightmoth) fireSpread(px, py float64, game *Game) {
	cx := b.x + 16
	cy := b.y + 16
	dx := px - cx
	dy := py - cy
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist == 0 {
		dist = 1
	}
	nx := dx / dist
	ny := dy / dist
	speed := 2.5
	offsets := []float64{-0.3, 0, 0.3}
	for _, off := range offsets {
		vx := (nx*math.Cos(off) - ny*math.Sin(off)) * speed
		vy := (nx*math.Sin(off) + ny*math.Cos(off)) * speed
		bullet := NewSpreadBullet(cx, cy, game.sporeImg, 1, vx, vy)
		game.enemyBullets = append(game.enemyBullets, bullet)
	}
}

func (b *Blightmoth) Draw(screen *ebiten.Image) {
	if b.img != nil {
		frame := (b.frameCounter / 8) % 8
		sx := frame * 64
		subImg := b.img.SubImage(image.Rect(sx, 0, sx+64, 64)).(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		if b.hitFlashTimer > 0 {
			op.ColorScale.ScaleWithColor(color.RGBA{255, 100, 100, 255})
		}
		op.GeoM.Translate(b.x, b.y)
		screen.DrawImage(subImg, op)
	}
	if b.healthBarTimer > 0 {
		b.drawHealthBar(screen)
	}
}

func (b *Blightmoth) drawHealthBar(screen *ebiten.Image) {
	barWidth := 32.0
	barHeight := 4.0
	barX := b.x
	barY := b.y - 8
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth), float32(barHeight), color.RGBA{255, 255, 255, 255}, false)
	pct := float64(b.health) / float64(b.maxHealth)
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth*pct), float32(barHeight), color.RGBA{255, 0, 0, 255}, false)
}

func (b *Blightmoth) GetPosition() (float64, float64) { return b.x, b.y }
func (b *Blightmoth) GetBounds() (float64, float64)   { return 32, 32 }
func (b *Blightmoth) TakeDamage(amount int) {
	if amount < b.health {
		b.healthBarTimer = 30
	}
	b.health -= amount
	b.hitFlashTimer = 10
}
func (b *Blightmoth) IsDead() bool    { return b.health <= 0 }
func (b *Blightmoth) OnDeath(game *Game) {
	numCoins := 2 + rand.IntN(2)
	for i := 0; i < numCoins; i++ {
		size := SmallCoin
		if ShouldSpawnBigCoin(game.player.luck, 1) {
			size = BigCoin
		}
		coin := NewCoin(b.x+float64(i*8), b.y, size, game.coinImg)
		game.coins = append(game.coins, coin)
	}
}
