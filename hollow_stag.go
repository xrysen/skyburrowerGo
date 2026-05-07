package main

import (
	"image"
	"image/color"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const HollowStagType EnemyType = "hollowstag"

const (
	stagTelegraphDuration = 60
	stagChargeSpeed       = 8.0
	stagResetX            = 640.0
	stagOffscreenLeft     = -64.0
)

type stagState int

const (
	stagStateTelegraph stagState = iota
	stagStateCharge
)

type HollowStag struct {
	x, y           float64
	img            *ebiten.Image
	state          stagState
	frameCounter   int
	telegraphTimer int
	health         int
	maxHealth      int
	hitFlashTimer  int
	healthBarTimer int
}

func NewHollowStag(x, y float64, img *ebiten.Image) *HollowStag {
	return &HollowStag{
		x:              x,
		y:              y,
		img:            img,
		state:          stagStateTelegraph,
		telegraphTimer: stagTelegraphDuration,
		health:         5,
		maxHealth:      5,
	}
}

func (s *HollowStag) Update(px, py float64, game *Game) {
	s.frameCounter++
	switch s.state {
	case stagStateTelegraph:
		s.telegraphTimer--
		if s.telegraphTimer <= 0 {
			s.state = stagStateCharge
		}
	case stagStateCharge:
		s.x -= stagChargeSpeed
		if s.x < stagOffscreenLeft {
			s.x = stagResetX
			s.state = stagStateTelegraph
			s.telegraphTimer = stagTelegraphDuration
		}
	}

	if s.hitFlashTimer > 0 {
		s.hitFlashTimer--
	}
	if s.healthBarTimer > 0 {
		s.healthBarTimer--
	}
}

func (s *HollowStag) Draw(screen *ebiten.Image) {
	if s.img != nil {
		frame := (s.frameCounter / 8) % 5
		sx := frame * 64
		subImg := s.img.SubImage(image.Rect(sx, 0, sx+64, 64)).(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		if s.hitFlashTimer > 0 {
			op.ColorScale.ScaleWithColor(color.RGBA{255, 100, 100, 255})
		}
		op.GeoM.Translate(s.x, s.y)
		screen.DrawImage(subImg, op)
	}
	if s.healthBarTimer > 0 {
		s.drawHealthBar(screen)
	}
}

func (s *HollowStag) drawHealthBar(screen *ebiten.Image) {
	barWidth := 32.0
	barHeight := 4.0
	barX := s.x
	barY := s.y - 8
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth), float32(barHeight), color.RGBA{255, 255, 255, 255}, false)
	healthPct := float64(s.health) / float64(s.maxHealth)
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth*healthPct), float32(barHeight), color.RGBA{255, 0, 0, 255}, false)
}

func (s *HollowStag) GetPosition() (float64, float64) { return s.x, s.y }
func (s *HollowStag) GetBounds() (float64, float64)   { return 32, 32 }
func (s *HollowStag) TakeDamage(amount int) {
	if amount < s.health {
		s.healthBarTimer = 30
	}
	s.health -= amount
	s.hitFlashTimer = 10
}
func (s *HollowStag) IsDead() bool          { return s.health <= 0 }
func (s *HollowStag) OnDeath(game *Game) {
	numCoins := 3 + rand.IntN(2)
	for i := 0; i < numCoins; i++ {
		size := SmallCoin
		if ShouldSpawnBigCoin(game.player.luck, 1) {
			size = BigCoin
		}
		coin := NewCoin(s.x+float64(i*8), s.y, size, game.coinImg)
		game.coins = append(game.coins, coin)
	}
}
