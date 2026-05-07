package main

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const HeartwoodType EnemyType = "heartwood"

const (
	heartwoodSpikeInterval1    = 120
	heartwoodWispInterval1     = 300
	heartwoodSpikeInterval2    = 90  // 75% of phase 1
	heartwoodWispInterval2     = 225 // 75% of phase 1
	heartwoodLightningInterval = 180
	heartwoodDrillInterval     = 200
	// 3s flicker (180) + 2s fade (120) = 300 frames
	heartwoodDeathDuration = 300
	heartwoodIntroDelay    = 90
	heartwoodIntroDuration = 180
)

type debrisParticle struct {
	x, y   float64
	vx, vy float64
	life   int
	maxLife int
	size   float64
	r, g, b uint8
}

type Heartwood struct {
	x, y           float64
	img            *ebiten.Image
	health         int
	maxHealth      int
	hitFlashTimer  int
	healthBarTimer int
	shootTimer     int
	wispTimer      int
	lightningTimer int
	drillTimer     int
	frameCounter   int
	isDying        bool
	deathTimer     int
	introTimer     int
	targetX        float64
	targetY        float64
	debris         []debrisParticle
}

var debrisColors = []color.RGBA{
	{139, 90, 43, 255},
	{101, 67, 33, 255},
	{128, 128, 128, 255},
	{160, 140, 100, 255},
	{85, 55, 30, 255},
}

func NewHeartwood(x, y float64, health int, img *ebiten.Image) *Heartwood {
	if health <= 0 {
		health = 300
	}
	return &Heartwood{
		x:         x,
		y:         float64(ScreenHeight) + 20,
		targetX:   x,
		targetY:   y,
		img:       img,
		health:    health,
		maxHealth: health,
	}
}

func (h *Heartwood) Phase() int {
	pct := float64(h.health) / float64(h.maxHealth)
	if pct <= 0.33 {
		return 3
	}
	if pct <= 0.66 {
		return 2
	}
	return 1
}

func (h *Heartwood) spikeInterval() int {
	if h.Phase() >= 2 {
		return heartwoodSpikeInterval2
	}
	return heartwoodSpikeInterval1
}

func (h *Heartwood) wispInterval() int {
	if h.Phase() >= 2 {
		return heartwoodWispInterval2
	}
	return heartwoodWispInterval1
}

func (h *Heartwood) introActive() bool {
	return h.introTimer < heartwoodIntroDelay+heartwoodIntroDuration
}

func (h *Heartwood) spawnDebris() {
	c := debrisColors[rand.Intn(len(debrisColors))]
	baseX := h.x + float64(rand.Intn(80))
	baseY := h.y + 280 + float64(rand.Intn(30))
	side := 1.0
	if rand.Intn(2) == 0 {
		side = -1.0
	}
	vx := side * (1.5 + rand.Float64()*3.0)
	vy := -(1.0 + rand.Float64()*4.0)
	life := 25 + rand.Intn(35)
	h.debris = append(h.debris, debrisParticle{
		x: baseX, y: baseY,
		vx: vx, vy: vy,
		life: life, maxLife: life,
		size: 2 + rand.Float64()*5,
		r: c.R, g: c.G, b: c.B,
	})
}

func (h *Heartwood) Update(px, py float64, game *Game) {
	h.frameCounter++
	if h.isDying {
		h.deathTimer++
		return
	}

	// Update debris particles regardless of intro state
	active := h.debris[:0]
	for i := range h.debris {
		p := &h.debris[i]
		p.x += p.vx
		p.y += p.vy
		p.vy += 0.15 // gravity
		p.life--
		if p.life > 0 {
			active = append(active, *p)
		}
	}
	h.debris = active

	if h.introActive() {
		h.introTimer++
		if h.introTimer <= heartwoodIntroDelay {
			return
		}
		riseTimer := h.introTimer - heartwoodIntroDelay
		progress := float64(riseTimer) / float64(heartwoodIntroDuration)

		// Ease out: decelerate as it approaches final position
		h.y += (h.targetY - h.y) * 0.04

		// Jitter decreases as the boss settles
		jitter := math.Sin(float64(h.introTimer)*1.2) * (1.0 - progress) * 3.0
		h.x = h.targetX + jitter

		// Spawn debris bursts while rising — more frequent early on
		if h.introTimer%4 == 0 {
			count := 3
			if progress < 0.4 {
				count = 5
			}
			for i := 0; i < count; i++ {
				h.spawnDebris()
			}
		}

		if riseTimer >= heartwoodIntroDuration {
			h.x = h.targetX
			h.y = h.targetY
		}
		return
	}

	if h.hitFlashTimer > 0 {
		h.hitFlashTimer--
	}
	if h.healthBarTimer > 0 {
		h.healthBarTimer--
	}

	h.shootTimer++
	if h.shootTimer >= h.spikeInterval() {
		h.shootTimer = 0
		h.fireRootSpikes(game)
	}

	h.wispTimer++
	if h.wispTimer >= h.wispInterval() {
		h.wispTimer = 0
		h.spawnWraithWisp(game)
	}

	if h.Phase() == 3 {
		h.lightningTimer++
		if h.lightningTimer >= heartwoodLightningInterval {
			h.lightningTimer = 0
			h.fireCorruptedLightning(game)
		}

		h.drillTimer++
		if h.drillTimer >= heartwoodDrillInterval {
			h.drillTimer = 0
			h.fireDrillSpikeSpread(game)
		}
	}
}

func (h *Heartwood) fireRootSpikes(game *Game) {
	cx := h.x + 16
	cy := h.y + 160
	for i := -1; i <= 1; i++ {
		vy := float64(i) * 0.6
		game.enemyBullets = append(game.enemyBullets, NewSpreadBullet(cx, cy, game.sporeImg, 2, -2.0, vy))
	}
}

func (h *Heartwood) spawnWraithWisp(game *Game) {
	img := game.enemyImage[WraithWhispType]
	game.enemies = append(game.enemies, NewWraithWhisp(h.x-40, h.y+float64(rand.Intn(200)), img))
}

func (h *Heartwood) fireCorruptedLightning(game *Game) {
	cx := h.x + 16
	cy := h.y + 160
	for i := 0; i < 5; i++ {
		angle := math.Pi + (float64(i)-2)*math.Pi/16
		vx := math.Cos(angle) * 3.5
		vy := math.Sin(angle) * 3.5
		game.enemyBullets = append(game.enemyBullets, NewLightningBolt(cx, cy, vx, vy, game.boltImg, 0.8))
	}
}

func (h *Heartwood) fireDrillSpikeSpread(game *Game) {
	cx := h.x + 16
	cy := h.y + 160
	for i := 0; i < 5; i++ {
		angle := math.Pi + (float64(i)-2)*math.Pi/12
		vx := math.Cos(angle) * 2.5
		vy := math.Sin(angle) * 2.5
		game.enemyBullets = append(game.enemyBullets, NewDrillBit(cx, cy, vx, vy, game.drillBitImg))
	}
}

func (h *Heartwood) Draw(screen *ebiten.Image) {
	// Draw debris particles behind the boss
	for _, p := range h.debris {
		alpha := uint8(255 * float64(p.life) / float64(p.maxLife))
		vector.FillRect(screen,
			float32(p.x), float32(p.y),
			float32(p.size), float32(p.size),
			color.RGBA{p.r, p.g, p.b, alpha}, false)
	}

	if h.img != nil {
		if h.isDying {
			var alpha uint8
			if h.deathTimer < 180 {
				if (h.deathTimer/6)%2 == 0 {
					alpha = 255
				} else {
					alpha = 160
				}
			} else {
				fadeProgress := float64(h.deathTimer-180) / 120.0
				if fadeProgress > 1 {
					fadeProgress = 1
				}
				alpha = uint8(255 * (1 - fadeProgress))
			}
			if alpha > 0 {
				frame := (h.frameCounter / 8) % 7
				sx := frame * 128
				subImg := h.img.SubImage(image.Rect(sx, 0, sx+128, 128)).(*ebiten.Image)
				scale := 320.0 / 128.0
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(scale, scale)
				op.GeoM.Translate(h.x, h.y)
				op.ColorScale.ScaleWithColor(color.RGBA{180, 80, 200, alpha})
				screen.DrawImage(subImg, op)
			}
			return
		}

		frame := (h.frameCounter / 8) % 7
		sx := frame * 128
		subImg := h.img.SubImage(image.Rect(sx, 0, sx+128, 128)).(*ebiten.Image)
		scale := 320.0 / 128.0
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		op.GeoM.Translate(h.x, h.y)
		if h.hitFlashTimer > 0 {
			op.ColorScale.Scale(1, 0.3, 0.3, 1)
		}
		screen.DrawImage(subImg, op)
	}

	if !h.introActive() {
		h.drawHealthBar(screen)
	}
}

func (h *Heartwood) drawHealthBar(screen *ebiten.Image) {
	barWidth := 80.0
	barHeight := 6.0
	barX := h.x
	barY := h.y - 10
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth), float32(barHeight), color.RGBA{255, 255, 255, 255}, false)
	healthPct := float64(h.health) / float64(h.maxHealth)
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth*healthPct), float32(barHeight), color.RGBA{150, 50, 200, 255}, false)
}

func (h *Heartwood) GetPosition() (float64, float64) { return h.x, h.y }
func (h *Heartwood) GetBounds() (float64, float64)   { return 80, 320 }

func (h *Heartwood) TakeDamage(amount int) {
	if h.introActive() {
		return
	}
	if amount < h.health {
		h.healthBarTimer = 30
	}
	h.health -= amount
	h.hitFlashTimer = 10
	if h.health <= 0 && !h.isDying {
		h.isDying = true
	}
}

func (h *Heartwood) IsDead() bool { return h.isDying && h.deathTimer >= heartwoodDeathDuration }

func (h *Heartwood) OnDeath(game *Game) {
	if !h.isDying {
		h.isDying = true
		h.deathTimer = 0
	}
}

func (h *Heartwood) UpdateDeath(game *Game) {
	if h.isDying {
		h.deathTimer++
	}
	if h.isDying && h.deathTimer >= heartwoodDeathDuration {
		game.bossKilled = true
	}
}

func (h *Heartwood) IsDeathComplete() bool {
	return h.isDying && h.deathTimer >= heartwoodDeathDuration
}
