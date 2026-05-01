package main

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const StormSpriteType EnemyType = "stormsprite"

type StormSprite struct {
	x, y           float64
	img            *ebiten.Image
	bulletSpeed    float64
	health         int
	maxHealth      int
	frameCounter   int
	shootTimer     int
	hitFlashTimer  int
	healthBarTimer int
}

func NewStormSprite(x, y float64, img *ebiten.Image, bulletSpeed float64) *StormSprite {
	return &StormSprite{
		x:           x,
		y:           y,
		img:         img,
		bulletSpeed: bulletSpeed,
		health:      15,
		maxHealth:   15,
	}
}

func (s *StormSprite) Update(px, py float64, game *Game) {
	s.x -= 1.0
	s.frameCounter++

	if s.hitFlashTimer > 0 {
		s.hitFlashTimer--
	}
	if s.healthBarTimer > 0 {
		s.healthBarTimer--
	}

	s.shootTimer++
	if s.shootTimer >= 150 {
		s.shootTimer = 0
		ring := NewElectricalRing(s.x+32, s.y+32, -s.bulletSpeed)
		game.enemyBullets = append(game.enemyBullets, ring)
	}
}

func (s *StormSprite) Draw(screen *ebiten.Image) {
	if s.img != nil {
		frame := (s.frameCounter / 8) % 6
		sx := frame * 64
		rect := image.Rect(sx, 0, sx+64, 64)
		subImg := s.img.SubImage(rect).(*ebiten.Image)

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

func (s *StormSprite) drawHealthBar(screen *ebiten.Image) {
	barWidth := 64.0
	barHeight := 4.0
	barX := s.x
	barY := s.y - 8
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth), float32(barHeight), color.RGBA{255, 255, 255, 255}, false)
	healthPct := float64(s.health) / float64(s.maxHealth)
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth*healthPct), float32(barHeight), color.RGBA{255, 0, 0, 255}, false)
}

func (s *StormSprite) GetPosition() (float64, float64) { return s.x, s.y }
func (s *StormSprite) GetBounds() (float64, float64)   { return 64, 64 }
func (s *StormSprite) TakeDamage(amount int) {
	if amount < s.health {
		s.healthBarTimer = 30
	}
	s.health -= amount
	s.hitFlashTimer = 10
}
func (s *StormSprite) IsDead() bool { return s.health <= 0 }
func (s *StormSprite) OnDeath(game *Game) {
	size := SmallCoin
	if ShouldSpawnBigCoin(game.player.luck, 1) {
		size = BigCoin
	}
	coin := NewCoin(s.x, s.y, size, game.coinImg)
	game.coins = append(game.coins, coin)
}
