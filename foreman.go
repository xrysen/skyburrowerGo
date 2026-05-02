package main

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const ForemanType EnemyType = "foreman"

type Foreman struct {
	x, y         float64
	baseY        float64
	targetBaseY  float64
	repoTimer    int
	img          *ebiten.Image
	health       int
	maxHealth    int
	frameCounter int
	shootTimer   int
	stalTimer    int
	isDying        bool
	deathTimer     int
	hitFlashTimer  int
	healthBarTimer int
	beetleTimer    int
}

func NewForeman(x, y float64, health int, img *ebiten.Image) *Foreman {
	if health <= 0 {
		health = 200
	}
	return &Foreman{
		x:           x,
		y:           y,
		baseY:       y,
		targetBaseY: y,
		repoTimer:   180,
		img:         img,
		health:      health,
		maxHealth:   health,
	}
}

func (f *Foreman) Update(px, py float64, game *Game) {
	if f.isDying {
		f.deathTimer++
		return
	}
	f.frameCounter++

	// Periodic vertical repositioning
	f.repoTimer--
	if f.repoTimer <= 0 {
		f.repoTimer = 180 + rand.Intn(120)
		f.targetBaseY = 30 + float64(rand.Intn(180))
	}
	f.baseY += (f.targetBaseY - f.baseY) * 0.02

	// Sine bob: 8px amplitude
	f.y = f.baseY + 8*math.Sin(float64(f.frameCounter)*2*math.Pi/120)

	if f.hitFlashTimer > 0 {
		f.hitFlashTimer--
	}
	if f.healthBarTimer > 0 {
		f.healthBarTimer--
	}

	if f.IsPhase2() {
		f.shootTimer++
		if f.shootTimer >= 120 {
			f.shootTimer = 0
			f.fireDrillSalvo(px, py, game, 5)
		}
		f.stalTimer++
		if f.stalTimer >= 150 {
			f.stalTimer = 0
			f.dropStalactites(game, 5)
		}
		f.beetleTimer++
		if f.beetleTimer >= 300 {
			f.beetleTimer = 0
			game.enemies = append(game.enemies, NewDynamiteBeetle(f.x-40, f.y, nil))
		}
	} else {
		f.shootTimer++
		if f.shootTimer >= 180 {
			f.shootTimer = 0
			f.fireDrillSalvo(px, py, game, 3)
		}
		f.stalTimer++
		if f.stalTimer >= 240 {
			f.stalTimer = 0
			f.dropStalactites(game, 3)
		}
	}
}

func (f *Foreman) fireDrillSalvo(px, py float64, game *Game, count int) {
	cx := f.x + 32
	cy := f.y + 32
	dx := px - cx
	dy := py - cy
	baseAngle := math.Atan2(dy, dx)
	speed := 2.5
	spread := math.Pi / 12 // 15 degrees

	half := float64(count-1) / 2
	for i := 0; i < count; i++ {
		angle := baseAngle + spread*(float64(i)-half)
		vx := math.Cos(angle) * speed
		vy := math.Sin(angle) * speed
		game.enemyBullets = append(game.enemyBullets, NewDrillBit(cx, cy, vx, vy))
	}
}

func (f *Foreman) dropStalactites(game *Game, count int) {
	for i := 0; i < count; i++ {
		x := float64(80 + i*160)
		game.enemyBullets = append(game.enemyBullets, NewStalactite(x, 0, game.stalactiteImg))
	}
}

const foremanFrameCount = 5
const foremanFrameSize = 64
const foremanFrameDelay = 8

func (f *Foreman) Draw(screen *ebiten.Image) {
	if f.img != nil {
		if f.isDying {
			alpha := max(0, 255-f.deathTimer*2)
			if alpha <= 0 {
				return
			}
			frame := (f.frameCounter / foremanFrameDelay) % foremanFrameCount
			sx := frame * foremanFrameSize
			sub := f.img.SubImage(image.Rect(sx, 0, sx+foremanFrameSize, foremanFrameSize)).(*ebiten.Image)
			op := &ebiten.DrawImageOptions{}
			op.ColorScale.ScaleWithColor(color.RGBA{255, 255, 255, uint8(alpha)})
			op.GeoM.Translate(f.x, f.y)
			screen.DrawImage(sub, op)
			return
		}
		frame := (f.frameCounter / foremanFrameDelay) % foremanFrameCount
		sx := frame * foremanFrameSize
		sub := f.img.SubImage(image.Rect(sx, 0, sx+foremanFrameSize, foremanFrameSize)).(*ebiten.Image)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(f.x, f.y)
		if f.hitFlashTimer > 0 {
			op.ColorScale.Scale(1, 0.3, 0.3, 1)
		}
		screen.DrawImage(sub, op)
	}
	if f.healthBarTimer > 0 {
		f.drawHealthBar(screen)
	}
}

func (f *Foreman) drawHealthBar(screen *ebiten.Image) {
	barWidth := 64.0
	barHeight := 4.0
	barX := f.x
	barY := f.y - 8
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth), float32(barHeight), color.RGBA{255, 255, 255, 255}, false)
	healthPct := float64(f.health) / float64(f.maxHealth)
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth*healthPct), float32(barHeight), color.RGBA{255, 0, 0, 255}, false)
}

func (f *Foreman) GetPosition() (float64, float64) { return f.x, f.y }
func (f *Foreman) GetBounds() (float64, float64)   { return 64, 64 }
func (f *Foreman) TakeDamage(amount int) {
	if amount < f.health {
		f.healthBarTimer = 30
	}
	f.health -= amount
	f.hitFlashTimer = 10
	if f.health <= 0 && !f.isDying {
		f.isDying = true
	}
}
func (f *Foreman) IsDead() bool    { return f.health <= 0 }
func (f *Foreman) IsPhase2() bool  { return f.health <= f.maxHealth/2 }
func (f *Foreman) OnDeath(game *Game) {
	if !f.isDying {
		f.isDying = true
		f.deathTimer = 0
	}
}

func (f *Foreman) UpdateDeath(game *Game) {
	if f.isDying && f.deathTimer == 120 {
		f.createDeathRewards(game)
	}
	if f.isDying && f.deathTimer >= 480 {
		game.bossKilled = true
	}
}

func (f *Foreman) IsDeathComplete() bool {
	return f.isDying && f.deathTimer >= 480
}

func (f *Foreman) createDeathRewards(game *Game) {
	numCoins := 150 + rand.Intn(21)
	cx := f.x + 32
	cy := f.y + 32

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

// Stalactite is a falling hazard dropped by The Foreman. Implements Bullet.
type Stalactite struct {
	x, y  float64
	speed float64
	img   *ebiten.Image
}

func NewStalactite(x, y float64, img *ebiten.Image) *Stalactite {
	return &Stalactite{x: x, y: y, speed: 3.0, img: img}
}

func (s *Stalactite) Update() {
	s.y += s.speed
}

func (s *Stalactite) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(s.x-24, s.y)
	screen.DrawImage(s.img, op)
}

func (s *Stalactite) GetPosition() (float64, float64) { return s.x, s.y }
func (s *Stalactite) GetDamage() int                  { return 4 }
func (s *Stalactite) GetBounds() (float64, float64)   { return 48, 48 }
