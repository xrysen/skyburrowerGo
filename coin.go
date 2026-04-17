package main

import (
	"image"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

type CoinSize int

const (
	SmallCoin CoinSize = iota
	BigCoin
)

type Coin struct {
	x, y         float64
	img          *ebiten.Image
	frameWidth   int
	frameHeight  int
	frameCounter int
	size         CoinSize
	value        int
	scale        float64
	collected    bool
	vx, vy       float64
}

func NewCoin(x, y float64, size CoinSize, img *ebiten.Image) *Coin {
	value := 1
	scale := 1.0
	if size == BigCoin {
		value = 5
		scale = 1.8
	}
	return &Coin{
		x:           x,
		y:           y,
		img:         img,
		frameWidth:  16,
		frameHeight: 16,
		size:        size,
		value:       value,
		scale:       scale,
	}
}

func (c *Coin) Update(g *Game) {
	c.frameCounter++

	// Apply velocity if set
	if c.vx != 0 || c.vy != 0 {
		c.x += c.vx
		c.y += c.vy

		// Bounce off screen boundaries
		coinWidth := float64(c.frameWidth) * c.scale
		coinHeight := float64(c.frameHeight) * c.scale

		// Bounce off left/right walls
		if c.x <= 0 {
			c.x = 0
			c.vx = math.Abs(c.vx) * 0.8 // Bounce right with energy loss
		} else if c.x >= ScreenWidth-coinWidth {
			c.x = ScreenWidth - coinWidth
			c.vx = -math.Abs(c.vx) * 0.8 // Bounce left with energy loss
		}

		// Bounce off top/bottom walls
		if c.y <= 0 {
			c.y = 0
			c.vy = math.Abs(c.vy) * 0.8 // Bounce down with energy loss
		} else if c.y >= ScreenHeight-coinHeight {
			c.y = ScreenHeight - coinHeight
			c.vy = -math.Abs(c.vy) * 0.8 // Bounce up with energy loss
		}

		// Apply friction
		c.vx *= 0.95
		c.vy *= 0.95

		// Stop if velocity is very small
		if math.Abs(c.vx) < 0.1 {
			c.vx = 0
		}
		if math.Abs(c.vy) < 0.1 {
			c.vy = 0
		}
	}

	// Check magnetic attraction
	player := g.player
	dx := player.x - c.x
	dy := player.y - c.y
	distance := (dx*dx + dy*dy) // Squared distance for performance

	// If within magnet range, attract coin to player
	if distance < player.magnetRange*player.magnetRange {
		// Calculate attraction force (stronger when closer)
		attractionSpeed := 3.0 + (player.magnetRange-math.Sqrt(distance))*0.1

		// Normalize direction and apply attraction
		if distance > 0 {
			c.x += (dx / math.Sqrt(distance)) * attractionSpeed
			c.y += (dy / math.Sqrt(distance)) * attractionSpeed
		}
	} else {
		// Normal movement when not in range
		c.x -= 1.5
	}
}

func (c *Coin) Draw(screen *ebiten.Image) {
	frame := (c.frameCounter / 8) % 7
	sx := frame * c.frameWidth
	rect := image.Rect(sx, 0, sx+c.frameWidth, c.frameHeight)
	subImg := c.img.SubImage(rect).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(c.scale, c.scale)
	op.GeoM.Translate(c.x, c.y)
	screen.DrawImage(subImg, op)
}

func ShouldSpawnBigCoin(playerLuck int, enemyDifficulty int) bool {
	chance := 2 + (playerLuck * 3) + (enemyDifficulty * 5)
	return rand.IntN(100) < int(chance)
}

func (c *Coin) SetVelocity(vx, vy float64) {
	c.vx = vx
	c.vy = vy
}
