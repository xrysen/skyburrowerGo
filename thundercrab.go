package main

import (
	"image"
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	ThunderCrabWidth  = 64
	ThunderCrabHeight = 64
)

type ThunderCrab struct {
	x, y           float64
	img            *ebiten.Image
	bulletSpeed    float64
	health         int
	maxHealth      int
	frameCounter   int
	hitFlashTimer  int
	healthBarTimer int

	// Movement
	baseY         float64
	movementTimer float64

	// Phase (1, 2, or 3)
	phase int

	// Attack timers (count up; fire when >= threshold)
	lightningTimer  int
	cloudTimer      int
	ringTimer       int
	chainTimer      int
	shockwaveTimer  int
	attackTellTimer int
	pendingAttack   int // 1=lightning, 2=cloud, 3=ring, 4=chain, 5=shockwave

	// Death animation
	deathTimer int
	isDying    bool
}

// Attack tell type constants (also used for tell color)
const (
	tellNone      = 0
	tellLightning = 1 // white
	tellCloud     = 2 // blue
	tellRing      = 3 // yellow
	tellChain     = 4 // purple
	tellShockwave = 5 // red
)

// Base attack intervals (frames). Phase 3 runs at 75% of these.
const (
	lightningInterval = 90
	cloudInterval     = 120
	ringInterval      = 150
	chainInterval     = 180
	shockwaveInterval = 200
)

func NewThunderCrab(x, y float64, health int, bulletSpeed float64, img *ebiten.Image) *ThunderCrab {
	return &ThunderCrab{
		x:           x,
		y:           y,
		baseY:       y,
		img:         img,
		bulletSpeed: bulletSpeed,
		health:      health,
		maxHealth:   health,
		phase:       1,
	}
}

func (tc *ThunderCrab) Update(px, py float64, game *Game) {
	if tc.isDying {
		tc.deathTimer++
		return
	}

	tc.frameCounter++
	tc.movementTimer += 0.025

	if tc.hitFlashTimer > 0 {
		tc.hitFlashTimer--
	}
	if tc.healthBarTimer > 0 {
		tc.healthBarTimer--
	}

	tc.updatePhase()
	tc.updateMovement()
	tc.updateAttacks(px, py, game)
}

func (tc *ThunderCrab) updatePhase() {
	hp := float64(tc.health) / float64(tc.maxHealth)
	switch {
	case hp <= 0.33:
		tc.phase = 3
	case hp <= 0.66:
		tc.phase = 2
	default:
		tc.phase = 1
	}
}

func (tc *ThunderCrab) updateMovement() {
	// Oscillate around vertical center of the playfield
	scaledH := float64(ThunderCrabHeight) * 1.5
	centerY := (float64(ScreenHeight) - scaledH) / 2
	tc.y = centerY + math.Sin(tc.movementTimer)*70

	// Horizontal drift: stay on the right side of the screen
	targetX := 500.0
	tc.x += (targetX - tc.x) * 0.015

	if tc.x < 400 {
		tc.x = 400
	}
	if tc.x > 600 {
		tc.x = 600
	}
}

func (tc *ThunderCrab) intervalFor(base int) int {
	if tc.phase == 3 {
		return int(float64(base) * 0.75)
	}
	return base
}

func (tc *ThunderCrab) updateAttacks(px, py float64, game *Game) {
	if tc.attackTellTimer > 0 {
		tc.attackTellTimer--
		if tc.attackTellTimer == 0 {
			tc.fireAttack(tc.pendingAttack, px, py, game)
			tc.pendingAttack = tellNone
		}
		return
	}

	// Increment all applicable timers and queue the first that fires
	tc.lightningTimer++
	tc.cloudTimer++
	if tc.phase >= 2 {
		tc.ringTimer++
	}
	if tc.phase >= 3 {
		tc.chainTimer++
		tc.shockwaveTimer++
	}

	switch {
	case tc.lightningTimer >= tc.intervalFor(lightningInterval):
		tc.lightningTimer = 0
		tc.queueAttack(tellLightning)
	case tc.cloudTimer >= tc.intervalFor(cloudInterval):
		tc.cloudTimer = 0
		tc.queueAttack(tellCloud)
	case tc.phase >= 2 && tc.ringTimer >= tc.intervalFor(ringInterval):
		tc.ringTimer = 0
		tc.queueAttack(tellRing)
	case tc.phase >= 3 && tc.chainTimer >= tc.intervalFor(chainInterval):
		tc.chainTimer = 0
		tc.queueAttack(tellChain)
	case tc.phase >= 3 && tc.shockwaveTimer >= tc.intervalFor(shockwaveInterval):
		tc.shockwaveTimer = 0
		tc.queueAttack(tellShockwave)
	}
}

func (tc *ThunderCrab) queueAttack(tellType int) {
	tc.attackTellTimer = 30
	tc.pendingAttack = tellType
}

func (tc *ThunderCrab) fireAttack(attackType int, px, py float64, game *Game) {
	cx := tc.x + ThunderCrabWidth/2
	cy := tc.y + ThunderCrabHeight/2

	switch attackType {
	case tellLightning:
		angle := math.Atan2(py-cy, px-cx)
		speed := tc.bulletSpeed
		bolt := NewLightningBolt(cx, cy, math.Cos(angle)*speed, math.Sin(angle)*speed, game.boltImg, 1.0)
		game.enemyBullets = append(game.enemyBullets, bolt)

	case tellCloud:
		angle := math.Atan2(py-cy, px-cx)
		speed := tc.bulletSpeed * 0.7
		cloud := NewCloudProjectile(cx, cy, math.Cos(angle)*speed, math.Sin(angle)*speed)
		game.enemyBullets = append(game.enemyBullets, cloud)

	case tellRing:
		ring := NewElectricalRing(cx, cy, -tc.bulletSpeed)
		game.enemyBullets = append(game.enemyBullets, ring)

	case tellChain:
		// 3-way spread
		baseAngle := math.Atan2(py-cy, px-cx)
		for i := -1; i <= 1; i++ {
			a := baseAngle + float64(i)*math.Pi/8
			bolt := NewChainLightningBolt(cx, cy, math.Cos(a)*tc.bulletSpeed, math.Sin(a)*tc.bulletSpeed, game.boltImg, 1.0)
			game.enemyBullets = append(game.enemyBullets, bolt)
		}

	case tellShockwave:
		sw := NewShockwave(cx, cy)
		game.enemyBullets = append(game.enemyBullets, sw)
	}
}

func (tc *ThunderCrab) Draw(screen *ebiten.Image) {
	if tc.isDying {
		alpha := max(0, 255-tc.deathTimer*2)
		if alpha <= 0 {
			return
		}
		if tc.img != nil {
			frame := (tc.frameCounter / 8) % 5
			sx := frame * ThunderCrabWidth
			rect := image.Rect(sx, 0, sx+ThunderCrabWidth, ThunderCrabHeight)
			subImg := tc.img.SubImage(rect).(*ebiten.Image)
			op := &ebiten.DrawImageOptions{}
			op.ColorScale.ScaleWithColor(color.RGBA{255, 255, 255, uint8(alpha)})
			op.GeoM.Scale(1.5, 1.5)
			op.GeoM.Translate(tc.x, tc.y)
			screen.DrawImage(subImg, op)
		}
		return
	}

	if tc.img != nil {
		frame := (tc.frameCounter / 8) % 5
		sx := frame * ThunderCrabWidth
		rect := image.Rect(sx, 0, sx+ThunderCrabWidth, ThunderCrabHeight)
		subImg := tc.img.SubImage(rect).(*ebiten.Image)

		op := &ebiten.DrawImageOptions{}
		tc.applyTellColor(op)
		op.GeoM.Scale(1.5, 1.5)
		op.GeoM.Translate(tc.x, tc.y)
		screen.DrawImage(subImg, op)
	}

	tc.drawHealthBar(screen)
}

func (tc *ThunderCrab) applyTellColor(op *ebiten.DrawImageOptions) {
	if tc.attackTellTimer <= 0 {
		if tc.hitFlashTimer > 0 {
			op.ColorScale.ScaleWithColor(color.RGBA{255, 100, 100, 255})
		}
		return
	}
	switch tc.pendingAttack {
	case tellLightning:
		op.ColorScale.ScaleWithColor(color.RGBA{255, 255, 255, 255}) // white
	case tellCloud:
		op.ColorScale.ScaleWithColor(color.RGBA{100, 100, 255, 255}) // blue
	case tellRing:
		op.ColorScale.ScaleWithColor(color.RGBA{255, 255, 0, 255}) // yellow
	case tellChain:
		op.ColorScale.ScaleWithColor(color.RGBA{180, 0, 255, 255}) // purple
	case tellShockwave:
		op.ColorScale.ScaleWithColor(color.RGBA{255, 0, 0, 255}) // red
	}
}

func (tc *ThunderCrab) drawHealthBar(screen *ebiten.Image) {
	scaledWidth := float64(ThunderCrabWidth) * 1.5
	barHeight := 6.0
	barX := tc.x
	barY := tc.y - 18

	vector.FillRect(screen, float32(barX), float32(barY), float32(scaledWidth), float32(barHeight), color.RGBA{50, 50, 50, 255}, false)
	healthPct := float64(tc.health) / float64(tc.maxHealth)
	vector.FillRect(screen, float32(barX), float32(barY), float32(scaledWidth*healthPct), float32(barHeight), color.RGBA{255, 0, 0, 255}, false)
	vector.StrokeRect(screen, float32(barX), float32(barY), float32(scaledWidth), float32(barHeight), 1, color.RGBA{255, 255, 255, 255}, false)
}

func (tc *ThunderCrab) GetPosition() (float64, float64) { return tc.x, tc.y }
func (tc *ThunderCrab) GetBounds() (float64, float64) {
	return float64(ThunderCrabWidth) * 1.5, float64(ThunderCrabHeight) * 1.5
}
func (tc *ThunderCrab) TakeDamage(amount int) {
	tc.health -= amount
	tc.hitFlashTimer = 10
	tc.healthBarTimer = 60
}
func (tc *ThunderCrab) IsDead() bool { return tc.health <= 0 || tc.isDying }
func (tc *ThunderCrab) OnDeath(game *Game) {
	if !tc.isDying {
		tc.isDying = true
		tc.deathTimer = 0
	}
}

func (tc *ThunderCrab) UpdateDeath(game *Game) {
	if tc.isDying && tc.deathTimer == 120 {
		tc.createDeathRewards(game)
	}
	if tc.isDying && tc.deathTimer >= 480 {
		game.bossKilled = true
	}
}

func (tc *ThunderCrab) IsDeathComplete() bool {
	return tc.isDying && tc.deathTimer >= 480
}

func (tc *ThunderCrab) createDeathRewards(game *Game) {
	numCoins := 150 + rand.IntN(21)
	cx := tc.x + float64(ThunderCrabWidth)/2
	cy := tc.y + float64(ThunderCrabHeight)/2

	for i := 0; i < numCoins; i++ {
		angle := 2 * math.Pi * float64(i) / float64(numCoins)
		speed := 1.5 + rand.Float64()*2.0
		vx := math.Cos(angle) * speed
		vy := math.Sin(angle) * speed

		size := SmallCoin
		if ShouldSpawnBigCoin(game.player.luck, 1) {
			size = BigCoin
		}
		coin := NewCoin(cx, cy, size, game.coinImg)
		coin.SetVelocity(vx, vy)
		game.coins = append(game.coins, coin)
	}
}
