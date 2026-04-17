package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type ThistleTurret struct {
	x, y           float64
	img            *ebiten.Image
	podImg         *ebiten.Image
	health         int
	maxHealth      int
	frameWidth     int
	frameHeight    int
	frameCounter   int
	hitFlashTimer  int
	healthBarTimer int

	// Growth behavior
	targetY      float64
	growthSpeed  float64
	isFullyGrown bool

	// Shooting behavior
	shootTimer    int
	shootCooldown int
	currentFrame  int
}

func NewThistleTurret(x, y float64, img, podImg *ebiten.Image) *ThistleTurret {
	return &ThistleTurret{
		x:             x,
		y:             y,
		img:           img,
		podImg:        podImg,
		health:        5,
		maxHealth:     5,
		frameWidth:    64,      // Assuming frame width, adjust based on actual sprite
		frameHeight:   64,      // Assuming frame height, adjust based on actual sprite
		targetY:       y - 110, // Grow up 110 pixels from spawn position
		growthSpeed:   1.0,
		shootTimer:    0,
		shootCooldown: 0,
	}
}

func (t *ThistleTurret) Update(px, py float64, game *Game) {
	// Growth behavior
	if !t.isFullyGrown {
		t.y -= t.growthSpeed
		if t.y <= t.targetY {
			t.y = t.targetY
			t.isFullyGrown = true
		}
	}

	// Update timers
	t.frameCounter++

	if t.hitFlashTimer > 0 {
		t.hitFlashTimer--
	}

	if t.healthBarTimer > 0 {
		t.healthBarTimer--
	}

	// Shoot with cooldown when fully grown
	if t.isFullyGrown {
		t.shootTimer++
		if t.shootCooldown > 0 {
			t.shootCooldown--
		}

		if t.shootTimer >= 60 && t.shootCooldown == 0 { // Every 60 frames (1 second)
			t.shootTimer = 0
			t.shootCooldown = 60 // 1 second cooldown between shots
			t.shoot(px, py, game)
		}
	}
}

func (t *ThistleTurret) shoot(px, py float64, game *Game) {
	// Determine direction and animation frame
	playerLeft := px < t.x
	playerAbove := py < t.y

	// Set animation frame based on player position
	if playerLeft && !playerAbove {
		t.currentFrame = 0 // Frames 0-1: left and level/below
	} else if !playerLeft && !playerAbove {
		t.currentFrame = 2 // Frames 2-3: right and level/below
	} else if playerLeft && playerAbove {
		t.currentFrame = 4 // Frames 4-5: left and above
	} else {
		t.currentFrame = 6 // Frames 6-7: right and above
	}

	// Calculate direction to player
	dx := px - t.x
	dy := py - t.y
	distance := math.Sqrt(dx*dx + dy*dy)

	if distance > 0 {
		// Normalize direction
		vx := (dx / distance) * 4.0 // Slower speed for better balance
		vy := (dy / distance) * 4.0

		// Create pod bullet
		bullet := NewPodBullet(t.x, t.y, t.podImg, vx, vy)
		game.enemyBullets = append(game.enemyBullets, bullet)
	}
}

func (t *ThistleTurret) Draw(screen *ebiten.Image) {
	// Constant animation at 5fps (12 frames per animation at 60fps = ~200ms per frame)
	animFrame := (t.frameCounter / 12) % 2
	frame := t.currentFrame + animFrame

	sx := frame * t.frameWidth
	rect := image.Rect(sx, 0, sx+t.frameWidth, t.frameHeight)
	subImg := t.img.SubImage(rect).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}

	if t.hitFlashTimer > 0 {
		op.ColorScale.ScaleWithColor(color.RGBA{255, 100, 100, 255})
	}

	op.GeoM.Translate(t.x, t.y)
	screen.DrawImage(subImg, op)

	if t.healthBarTimer > 0 {
		t.drawHealthBar(screen)
	}
}

func (t *ThistleTurret) drawHealthBar(screen *ebiten.Image) {
	barWidth := float64(t.frameWidth)
	barHeight := 4.0
	barX := t.x
	barY := t.y - 8

	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth), float32(barHeight), color.RGBA{255, 255, 255, 255}, false)

	healthPercent := float64(t.health) / float64(t.maxHealth)
	redWidth := barWidth * healthPercent
	vector.FillRect(screen, float32(barX), float32(barY), float32(redWidth), float32(barHeight), color.RGBA{255, 0, 0, 255}, false)
}

func (t *ThistleTurret) GetPosition() (x, y float64) {
	return t.x, t.y
}

func (t *ThistleTurret) GetBounds() (width, height float64) {
	return float64(t.frameWidth), float64(t.frameHeight)
}

func (t *ThistleTurret) TakeDamage(amount int) {
	// Only show health bar if enemy won't die in one hit
	if amount < t.health {
		t.healthBarTimer = 30
	}
	t.health -= amount
	t.hitFlashTimer = 10
}

func (t *ThistleTurret) IsDead() bool {
	return t.health <= 0
}

func (t *ThistleTurret) OnDeath(game *Game) {
	size := SmallCoin
	if ShouldSpawnBigCoin(game.player.luck, 1) {
		size = BigCoin
	}
	coin := NewCoin(t.x, t.y, size, game.coinImg)
	game.coins = append(game.coins, coin)
}
