package main

import (
	"image"
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const WraithWhispType EnemyType = "wraithwhisp"

type WraithWhisp struct {
	x, y           float64
	img            *ebiten.Image
	speed          float64
	frameCounter   int
	driftTimer     int
	health         int
	maxHealth      int
	hitFlashTimer  int
	healthBarTimer int
	isChild        bool
	targetX        float64
	targetY        float64
	trackSpeed     float64
}

func NewWraithWhisp(x, y float64, img *ebiten.Image) *WraithWhisp {
	return &WraithWhisp{
		x:         x,
		y:         y,
		img:       img,
		speed:     0.8,
		health:    4,
		maxHealth: 4,
	}
}

func newChildWraithWhisp(x, y float64, img *ebiten.Image, targetX, targetY float64) *WraithWhisp {
	return &WraithWhisp{
		x:          x,
		y:          y,
		img:        img,
		speed:      1.5,
		trackSpeed: 1.2,
		health:     2,
		maxHealth:  2,
		isChild:    true,
		targetX:    targetX,
		targetY:    targetY,
	}
}

func (w *WraithWhisp) Update(px, py float64, game *Game) {
	w.frameCounter++
	if w.isChild {
		dx := w.targetX - w.x
		dy := w.targetY - w.y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > 0 {
			w.x += (dx / dist) * w.trackSpeed
			w.y += (dy / dist) * w.trackSpeed
		}
		w.x -= w.speed * 0.3
	} else {
		w.x -= w.speed
		w.driftTimer++
		w.y += math.Sin(float64(w.driftTimer)*0.05) * 0.4
	}

	if w.hitFlashTimer > 0 {
		w.hitFlashTimer--
	}
	if w.healthBarTimer > 0 {
		w.healthBarTimer--
	}
}

func (w *WraithWhisp) Draw(screen *ebiten.Image) {
	if w.img != nil {
		frame := (w.frameCounter / 8) % 5
		sx := frame * 64
		subImg := w.img.SubImage(image.Rect(sx, 0, sx+64, 64)).(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		if w.hitFlashTimer > 0 {
			op.ColorScale.ScaleWithColor(color.RGBA{255, 100, 100, 255})
		}
		op.GeoM.Translate(w.x, w.y)
		screen.DrawImage(subImg, op)
	}
	if w.healthBarTimer > 0 {
		w.drawHealthBar(screen)
	}
}

func (w *WraithWhisp) drawHealthBar(screen *ebiten.Image) {
	barWidth := 32.0
	barHeight := 4.0
	barX := w.x
	barY := w.y - 8
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth), float32(barHeight), color.RGBA{255, 255, 255, 255}, false)
	healthPct := float64(w.health) / float64(w.maxHealth)
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth*healthPct), float32(barHeight), color.RGBA{255, 0, 0, 255}, false)
}

func (w *WraithWhisp) GetPosition() (float64, float64) { return w.x, w.y }
func (w *WraithWhisp) GetBounds() (float64, float64)   { return 24, 24 }
func (w *WraithWhisp) TakeDamage(amount int) {
	if amount < w.health {
		w.healthBarTimer = 30
	}
	w.health -= amount
	w.hitFlashTimer = 10
}
func (w *WraithWhisp) IsDead() bool { return w.health <= 0 }
func (w *WraithWhisp) OnDeath(game *Game) {
	numCoins := 2 + rand.IntN(2)
	for i := 0; i < numCoins; i++ {
		size := SmallCoin
		if ShouldSpawnBigCoin(game.player.luck, 1) {
			size = BigCoin
		}
		coin := NewCoin(w.x+float64(i*8), w.y, size, game.coinImg)
		game.coins = append(game.coins, coin)
	}
	if w.isChild {
		return
	}
	px, py := w.x-100, w.y
	if game.player != nil {
		px = game.player.x
		py = game.player.y
	}
	game.enemies = append(game.enemies,
		newChildWraithWhisp(w.x-8, w.y-8, w.img, px, py),
		newChildWraithWhisp(w.x+8, w.y+8, w.img, px, py),
	)
}
