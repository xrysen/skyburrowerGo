package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const LightningBugType EnemyType = "lightningbug"

type LightningBug struct {
	x, y           float64
	img            *ebiten.Image
	speed          float64
	bulletSpeed    float64
	health         int
	maxHealth      int
	frameCounter   int
	shootTimer     int
	chainCount     int
	chainTimer     int
	hitFlashTimer  int
	healthBarTimer int
}

func NewLightningBug(x, y float64, img *ebiten.Image, bulletSpeed float64) *LightningBug {
	return &LightningBug{
		x:           x,
		y:           y,
		img:         img,
		speed:       3.5,
		bulletSpeed: bulletSpeed,
		health:      3,
		maxHealth:   3,
	}
}

func (lb *LightningBug) Update(px, py float64, game *Game) {
	lb.x -= lb.speed
	lb.frameCounter++

	if lb.hitFlashTimer > 0 {
		lb.hitFlashTimer--
	}
	if lb.healthBarTimer > 0 {
		lb.healthBarTimer--
	}

	lb.shootTimer++
	if lb.shootTimer >= 120 {
		lb.shootTimer = 0
		lb.chainCount = 0
		lb.chainTimer = 0
		lb.fireChainBolt(px, py, game)
		lb.chainCount = 1
	}

	if lb.chainCount > 0 && lb.chainCount < 4 {
		lb.chainTimer++
		if lb.chainTimer >= 8 {
			lb.chainTimer = 0
			lb.fireChainBolt(px, py, game)
			lb.chainCount++
		}
	}
}

func (lb *LightningBug) fireChainBolt(px, py float64, game *Game) {
	cx := lb.x + 32
	cy := lb.y + 32
	dx := px - cx
	dy := py - cy
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist == 0 {
		dist = 1
	}
	vx := (dx / dist) * lb.bulletSpeed
	vy := (dy / dist) * lb.bulletSpeed
	bolt := NewChainLightningBolt(cx, cy, vx, vy, game.boltImg, 0.5)
	game.enemyBullets = append(game.enemyBullets, bolt)
}

func (lb *LightningBug) Draw(screen *ebiten.Image) {
	if lb.img != nil {
		frame := (lb.frameCounter / 8) % 6
		sx := frame * 64
		rect := image.Rect(sx, 0, sx+64, 64)
		subImg := lb.img.SubImage(rect).(*ebiten.Image)

		op := &ebiten.DrawImageOptions{}
		if lb.hitFlashTimer > 0 {
			op.ColorScale.ScaleWithColor(color.RGBA{255, 100, 100, 255})
		}
		op.GeoM.Translate(lb.x, lb.y)
		screen.DrawImage(subImg, op)
	}

	if lb.healthBarTimer > 0 {
		lb.drawHealthBar(screen)
	}
}

func (lb *LightningBug) drawHealthBar(screen *ebiten.Image) {
	barWidth := 64.0
	barHeight := 4.0
	barX := lb.x
	barY := lb.y - 8
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth), float32(barHeight), color.RGBA{255, 255, 255, 255}, false)
	healthPct := float64(lb.health) / float64(lb.maxHealth)
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth*healthPct), float32(barHeight), color.RGBA{255, 0, 0, 255}, false)
}

func (lb *LightningBug) GetPosition() (float64, float64) { return lb.x, lb.y }
func (lb *LightningBug) GetBounds() (float64, float64)   { return 64, 64 }
func (lb *LightningBug) TakeDamage(amount int) {
	if amount < lb.health {
		lb.healthBarTimer = 30
	}
	lb.health -= amount
	lb.hitFlashTimer = 10
}
func (lb *LightningBug) IsDead() bool { return lb.health <= 0 }
func (lb *LightningBug) OnDeath(game *Game) {
	size := SmallCoin
	if ShouldSpawnBigCoin(game.player.luck, 1) {
		size = BigCoin
	}
	coin := NewCoin(lb.x, lb.y, size, game.coinImg)
	game.coins = append(game.coins, coin)
}
