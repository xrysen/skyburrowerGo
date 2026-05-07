package main

import (
	"image"
	"image/color"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const DynamiteBeetleType EnemyType = "dynamitebeetle"

type DynamiteBeetle struct {
	x, y           float64
	img            *ebiten.Image
	speed          float64
	health         int
	maxHealth      int
	frameCounter   int
	shootTimer     int
	hitFlashTimer  int
	healthBarTimer int
}

func NewDynamiteBeetle(x, y float64, img *ebiten.Image) *DynamiteBeetle {
	return &DynamiteBeetle{
		x:         x,
		y:         y,
		img:       img,
		speed:     1.0,
		health:    6,
		maxHealth: 6,
	}
}

func (d *DynamiteBeetle) Update(px, py float64, game *Game) {
	d.x -= d.speed
	d.frameCounter++

	if d.hitFlashTimer > 0 {
		d.hitFlashTimer--
	}
	if d.healthBarTimer > 0 {
		d.healthBarTimer--
	}

	d.shootTimer++
	if d.shootTimer >= 150 {
		d.shootTimer = 0
		vx := -3.5
		if px > d.x {
			vx = 3.5
		}
		game.enemyBullets = append(game.enemyBullets, NewFuseSpark(d.x+16, d.y+16, vx, 0))
	}
}

func (d *DynamiteBeetle) Draw(screen *ebiten.Image) {
	if d.img != nil {
		frame := (d.frameCounter / 10) % 6
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

func (d *DynamiteBeetle) drawHealthBar(screen *ebiten.Image) {
	barWidth := 32.0
	barHeight := 4.0
	barX := d.x
	barY := d.y - 8
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth), float32(barHeight), color.RGBA{255, 255, 255, 255}, false)
	healthPct := float64(d.health) / float64(d.maxHealth)
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth*healthPct), float32(barHeight), color.RGBA{255, 0, 0, 255}, false)
}

func (d *DynamiteBeetle) GetPosition() (float64, float64) { return d.x, d.y }
func (d *DynamiteBeetle) GetBounds() (float64, float64)   { return 32, 32 }
func (d *DynamiteBeetle) TakeDamage(amount int) {
	if amount < d.health {
		d.healthBarTimer = 30
	}
	d.health -= amount
	d.hitFlashTimer = 10
}
func (d *DynamiteBeetle) IsDead() bool { return d.health <= 0 }
func (d *DynamiteBeetle) OnDeath(game *Game) {
	cx := d.x + 16
	cy := d.y + 16
	speed := 3.0
	dirs := [][2]float64{{speed, 0}, {-speed, 0}, {0, speed}, {0, -speed}}
	for _, dir := range dirs {
		game.enemyBullets = append(game.enemyBullets, NewFuseSpark(cx, cy, dir[0], dir[1]))
	}

	numCoins := 2 + rand.IntN(2)
	for i := 0; i < numCoins; i++ {
		size := SmallCoin
		if game.player != nil && ShouldSpawnBigCoin(game.player.luck, 1) {
			size = BigCoin
		}
		coin := NewCoin(d.x+float64(i*8), d.y, size, game.coinImg)
		game.coins = append(game.coins, coin)
	}
}
