package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Sporespinner struct {
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
	shootTimer     int
}

func NewSporespinner(x, y float64, img *ebiten.Image) *Sporespinner {
	return &Sporespinner{
		x:           x,
		y:           y,
		img:         img,
		speed:       2.0,
		health:      3,
		maxHealth:   3,
		frameWidth:  64,
		frameHeight: 64,
		shootTimer:  0,
	}
}

func (s *Sporespinner) Update(px, py float64, game *Game) {
	s.x -= s.speed
	s.frameCounter++
	s.shootTimer++

	if s.hitFlashTimer > 0 {
		s.hitFlashTimer--
	}

	if s.healthBarTimer > 0 {
		s.healthBarTimer--
	}

	// Shoot spore every 90 frames (1.5 seconds at 60 FPS)
	if s.shootTimer >= 90 {
		s.shootTimer = 0
		s.shootSpore(game)
	}
}

func (s *Sporespinner) shootSpore(game *Game) {
	// Create spore bullet moving in straight line to the right
	spore := NewSporeBullet(s.x, s.y+float64(s.frameHeight)/2, game.sporeImg)
	game.enemyBullets = append(game.enemyBullets, spore)
}

func (s *Sporespinner) Draw(screen *ebiten.Image) {
	frame := (s.frameCounter / 8) % 4
	sx := frame * s.frameWidth
	rect := image.Rect(sx, 0, sx+s.frameWidth, s.frameHeight)
	subImg := s.img.SubImage(rect).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}

	if s.hitFlashTimer > 0 {
		op.ColorScale.ScaleWithColor(color.RGBA{255, 100, 100, 255})
	}

	op.GeoM.Translate(s.x, s.y)
	screen.DrawImage(subImg, op)

	if s.healthBarTimer > 0 {
		s.drawHealthBar(screen)
	}
}

func (s *Sporespinner) drawHealthBar(screen *ebiten.Image) {
	barWidth := float64(s.frameWidth)
	barHeight := 4.0
	barX := s.x
	barY := s.y - 8

	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth), float32(barHeight), color.RGBA{255, 255, 255, 255}, false)

	healthPercent := float64(s.health) / float64(s.maxHealth)
	redWidth := barWidth * healthPercent
	vector.FillRect(screen, float32(barX), float32(barY), float32(redWidth), float32(barHeight), color.RGBA{255, 0, 0, 255}, false)
}

func (s *Sporespinner) GetPosition() (x, y float64) {
	return s.x, s.y
}

func (s *Sporespinner) GetBounds() (width, height float64) {
	return float64(s.frameWidth), float64(s.frameHeight)
}

func (s *Sporespinner) TakeDamage(amount int) {
	// Only show health bar if enemy won't die in one hit
	if amount < s.health {
		s.healthBarTimer = 30
	}
	s.health -= amount
	s.hitFlashTimer = 10
}

func (s *Sporespinner) IsDead() bool {
	return s.health <= 0
}

func (s *Sporespinner) OnDeath(game *Game) {
	// Create explosion of 8 spores in circular pattern
	s.createSporeExplosion(game)
	
	// Drop coin
	size := SmallCoin
	if ShouldSpawnBigCoin(game.player.luck, 1) {
		size = BigCoin
	}
	coin := NewCoin(s.x, s.y, size, game.coinImg)
	game.coins = append(game.coins, coin)
}

func (s *Sporespinner) createSporeExplosion(game *Game) {
	centerX := s.x + float64(s.frameWidth)/2
	centerY := s.y + float64(s.frameHeight)/2
	numSpores := 8
	explosionSpeed := 3.0

	for i := 0; i < numSpores; i++ {
		angle := 2 * math.Pi * float64(i) / float64(numSpores)
		vx := math.Cos(angle) * explosionSpeed
		vy := math.Sin(angle) * explosionSpeed
		
		spore := NewSporeBullet(centerX, centerY, game.sporeImg)
		spore.SetVelocity(vx, vy)
		game.enemyBullets = append(game.enemyBullets, spore)
	}
}
