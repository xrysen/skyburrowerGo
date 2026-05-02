package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const DarkWingType EnemyType = "darkwing"

type DarkWing struct {
	x, y           float64
	img            *ebiten.Image
	speed          float64
	trackSpeed     float64
	health         int
	maxHealth      int
	frameCounter   int
	hitFlashTimer  int
	healthBarTimer int
}

func NewDarkWing(x, y float64, img *ebiten.Image) *DarkWing {
	return &DarkWing{
		x:          x,
		y:          y,
		img:        img,
		speed:      2.0,
		trackSpeed: 1.5,
		health:     3,
		maxHealth:  3,
	}
}

func (d *DarkWing) Update(px, py float64, game *Game) {
	d.x -= d.speed
	d.frameCounter++

	// Track player Y position, capped per frame
	if py > d.y {
		d.y += d.trackSpeed
		if d.y > py {
			d.y = py
		}
	} else if py < d.y {
		d.y -= d.trackSpeed
		if d.y < py {
			d.y = py
		}
	}

	// Clamp to screen bounds
	if d.y < 0 {
		d.y = 0
	}
	if d.y > float64(ScreenHeight)-32 {
		d.y = float64(ScreenHeight) - 32
	}

	if d.hitFlashTimer > 0 {
		d.hitFlashTimer--
	}
	if d.healthBarTimer > 0 {
		d.healthBarTimer--
	}
}

func (d *DarkWing) animFrame() int {
	// Ping-pong over 8 frames: 0,1,2,3,4,5,6,7,6,5,4,3,2,1,...
	pos := (d.frameCounter / 6) % 14
	if pos < 8 {
		return pos
	}
	return 14 - pos
}

func (d *DarkWing) Draw(screen *ebiten.Image) {
	if d.img != nil {
		frame := d.animFrame()
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

func (d *DarkWing) drawHealthBar(screen *ebiten.Image) {
	barWidth := 32.0
	barHeight := 4.0
	barX := d.x
	barY := d.y - 8
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth), float32(barHeight), color.RGBA{255, 255, 255, 255}, false)
	healthPct := float64(d.health) / float64(d.maxHealth)
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth*healthPct), float32(barHeight), color.RGBA{255, 0, 0, 255}, false)
}

func (d *DarkWing) GetPosition() (float64, float64) { return d.x, d.y }
func (d *DarkWing) GetBounds() (float64, float64)   { return 32, 32 }
func (d *DarkWing) TakeDamage(amount int) {
	if amount < d.health {
		d.healthBarTimer = 30
	}
	d.health -= amount
	d.hitFlashTimer = 10
}
func (d *DarkWing) IsDead() bool { return d.health <= 0 }
func (d *DarkWing) OnDeath(game *Game) {
	size := SmallCoin
	if ShouldSpawnBigCoin(game.player.luck, 1) {
		size = BigCoin
	}
	coin := NewCoin(d.x, d.y, size, game.coinImg)
	game.coins = append(game.coins, coin)
}
