package main

import (
	"image"
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const CloudDrifterType EnemyType = "clouddrifter"

type CloudDrifter struct {
	x, y           float64
	baseY          float64
	img            *ebiten.Image
	speed          float64
	bulletSpeed    float64
	health         int
	maxHealth      int
	frameCounter   int
	shootTimer     int
	hitFlashTimer  int
	healthBarTimer int
}

func NewCloudDrifter(x, y float64, img *ebiten.Image, bulletSpeed float64) *CloudDrifter {
	return &CloudDrifter{
		x:           x,
		y:           y,
		baseY:       y,
		img:         img,
		speed:       1.5,
		bulletSpeed: bulletSpeed,
		health:      6,
		maxHealth:   6,
	}
}

func (c *CloudDrifter) Update(px, py float64, game *Game) {
	c.x -= c.speed
	c.frameCounter++

	// Sinusoidal vertical drift: 30px amplitude, one cycle every 180 frames
	c.y = c.baseY + 30*math.Sin(float64(c.frameCounter)*2*math.Pi/180)
	if c.y < 0 {
		c.y = 0
	}
	if c.y > float64(ScreenHeight)-64 {
		c.y = float64(ScreenHeight) - 64
	}

	if c.hitFlashTimer > 0 {
		c.hitFlashTimer--
	}
	if c.healthBarTimer > 0 {
		c.healthBarTimer--
	}

	c.shootTimer++
	if c.shootTimer >= 90 {
		c.shootTimer = 0
		c.fireCircularSpread(game)
	}
}

func (c *CloudDrifter) fireCircularSpread(game *Game) {
	cx := c.x + 32
	cy := c.y + 32
	for i := 0; i < 8; i++ {
		angle := 2 * math.Pi * float64(i) / 8
		vx := math.Cos(angle) * c.bulletSpeed
		vy := math.Sin(angle) * c.bulletSpeed
		bolt := NewLightningBolt(cx, cy, vx, vy, game.boltImg, 0.5)
		game.enemyBullets = append(game.enemyBullets, bolt)
	}
}

func (c *CloudDrifter) Draw(screen *ebiten.Image) {
	if c.img != nil {
		frame := (c.frameCounter / 8) % 5
		sx := frame * 64
		rect := image.Rect(sx, 0, sx+64, 64)
		subImg := c.img.SubImage(rect).(*ebiten.Image)

		op := &ebiten.DrawImageOptions{}
		if c.hitFlashTimer > 0 {
			op.ColorScale.ScaleWithColor(color.RGBA{255, 100, 100, 255})
		}
		op.GeoM.Translate(c.x, c.y)
		screen.DrawImage(subImg, op)
	}

	if c.healthBarTimer > 0 {
		c.drawHealthBar(screen)
	}
}

func (c *CloudDrifter) drawHealthBar(screen *ebiten.Image) {
	barWidth := 64.0
	barHeight := 4.0
	barX := c.x
	barY := c.y - 8
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth), float32(barHeight), color.RGBA{255, 255, 255, 255}, false)
	healthPct := float64(c.health) / float64(c.maxHealth)
	vector.FillRect(screen, float32(barX), float32(barY), float32(barWidth*healthPct), float32(barHeight), color.RGBA{255, 0, 0, 255}, false)
}

func (c *CloudDrifter) GetPosition() (float64, float64)  { return c.x, c.y }
func (c *CloudDrifter) GetBounds() (float64, float64)    { return 64, 64 }
func (c *CloudDrifter) TakeDamage(amount int) {
	if amount < c.health {
		c.healthBarTimer = 30
	}
	c.health -= amount
	c.hitFlashTimer = 10
}
func (c *CloudDrifter) IsDead() bool { return c.health <= 0 }
func (c *CloudDrifter) OnDeath(game *Game) {
	numCoins := 2 + rand.IntN(2)
	for i := 0; i < numCoins; i++ {
		size := SmallCoin
		if ShouldSpawnBigCoin(game.player.luck, 1) {
			size = BigCoin
		}
		coin := NewCoin(c.x+float64(i*8), c.y, size, game.coinImg)
		game.coins = append(game.coins, coin)
	}
}
