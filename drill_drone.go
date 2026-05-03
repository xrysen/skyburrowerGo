package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const DrillDroneType EnemyType = "drilldrone"

type DrillDrone struct {
	x, y           float64
	img            *ebiten.Image
	speed          float64
	health         int
	maxHealth      int
	frameCounter    int
	shootTimer      int
	dashTimer       int
	dashRemaining   float64
	dashDirection   float64
	hitFlashTimer   int
	healthBarTimer  int
}

func NewDrillDrone(x, y float64, img *ebiten.Image) *DrillDrone {
	return &DrillDrone{
		x:         x,
		y:         y,
		img:           img,
		speed:         2.0,
		health:        20,
		maxHealth:     20,
		dashDirection: 1,
	}
}

func (d *DrillDrone) Update(px, py float64, game *Game) {
	d.x -= d.speed
	d.frameCounter++

	if d.hitFlashTimer > 0 {
		d.hitFlashTimer--
	}
	if d.healthBarTimer > 0 {
		d.healthBarTimer--
	}

	// Vertical dash
	if d.dashRemaining > 0 {
		step := math.Min(4.0, d.dashRemaining)
		d.y += d.dashDirection * step
		d.dashRemaining -= step
		if d.dashRemaining <= 0 {
			d.dashDirection = -d.dashDirection
		}
	} else {
		d.dashTimer++
		if d.dashTimer >= 120 {
			d.dashTimer = 0
			d.dashRemaining = 40
		}
	}

	d.shootTimer++
	if d.shootTimer >= 100 {
		d.shootTimer = 0
		dx := px - (d.x + 16)
		dy := py - (d.y + 16)
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist == 0 {
			dist = 1
		}
		speed := 2.5
		vx := dx / dist * speed
		vy := dy / dist * speed
		game.enemyBullets = append(game.enemyBullets, NewDrillBit(d.x+16, d.y+16, vx, vy, game.drillBitImg))
	}
}

func (d *DrillDrone) Draw(screen *ebiten.Image) {
	if d.img != nil {
		frame := (d.frameCounter / 6) % 5
		sx := frame * 32
		rect := image.Rect(sx, 0, sx+32, 32)
		subImg := d.img.SubImage(rect).(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		if d.hitFlashTimer > 0 {
			op.ColorScale.ScaleWithColor(color.RGBA{255, 100, 100, 255})
		}
		op.GeoM.Translate(d.x, d.y)
		screen.DrawImage(subImg, op)
	}
	if d.healthBarTimer > 0 {
		d.drawHealthBar(screen)
	}
}

func (d *DrillDrone) drawHealthBar(screen *ebiten.Image) {
	barWidth := 32.0
	barHeight := 4.0
	barX := d.x
	barY := d.y - 8
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth), float32(barHeight), color.RGBA{255, 255, 255, 255}, false)
	healthPct := float64(d.health) / float64(d.maxHealth)
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth*healthPct), float32(barHeight), color.RGBA{255, 0, 0, 255}, false)
}

func (d *DrillDrone) GetPosition() (float64, float64) { return d.x, d.y }
func (d *DrillDrone) GetBounds() (float64, float64)   { return 32, 32 }
func (d *DrillDrone) TakeDamage(amount int) {
	if amount < d.health {
		d.healthBarTimer = 30
	}
	d.health -= amount
	d.hitFlashTimer = 10
}
func (d *DrillDrone) IsDead() bool { return d.health <= 0 }
func (d *DrillDrone) OnDeath(game *Game) {
	size := BigCoin
	if game.player != nil && !ShouldSpawnBigCoin(game.player.luck, 1) {
		size = SmallCoin
	}
	coin := NewCoin(d.x, d.y, size, game.coinImg)
	game.coins = append(game.coins, coin)
}
