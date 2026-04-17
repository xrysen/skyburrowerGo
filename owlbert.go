package main

import (
	"image"
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Owlbert struct {
	x, y           float64
	img            *ebiten.Image
	health         int
	maxHealth      int
	frameWidth     int
	frameHeight    int
	frameCounter   int
	hitFlashTimer  int
	healthBarTimer int

	// Movement
	baseY         float64
	movementTimer float64
	movementSpeed float64

	// Attacks
	featherTimer    int
	talonTimer      int
	minionTimer     int
	attackTellTimer int
	attackTellType  int // 0 = none, 1 = feather, 2 = talon, 3 = minion

	// Death animation
	deathTimer int
	isDying    bool
}

const (
	OwlbertWidth  = 65
	OwlbertHeight = 65
	OwlbertHealth = 100
)

func NewOwlbert(x, y float64, img *ebiten.Image) *Owlbert {
	return &Owlbert{
		x:             x,
		y:             y,
		baseY:         y,
		img:           img,
		health:        OwlbertHealth,
		maxHealth:     OwlbertHealth,
		frameWidth:    OwlbertWidth,
		frameHeight:   OwlbertHeight,
		movementSpeed: 2.0,
		featherTimer:  60,  // 1 second at 60 FPS (more frequent)
		talonTimer:    300, // 5 seconds at 60 FPS (less frequent)
		minionTimer:   300, // 5 seconds at 60 FPS
	}
}

func (o *Owlbert) Update(px, py float64, game *Game) {
	if o.isDying {
		o.deathTimer++
		return
	}

	o.frameCounter++
	o.movementTimer += 0.03

	// Update timers
	if o.featherTimer > 0 {
		o.featherTimer--
	}
	if o.talonTimer > 0 {
		o.talonTimer--
	}
	if o.minionTimer > 0 {
		o.minionTimer--
	}
	if o.attackTellTimer > 0 {
		o.attackTellTimer--
	}

	if o.hitFlashTimer > 0 {
		o.hitFlashTimer--
	}

	if o.healthBarTimer > 0 {
		o.healthBarTimer--
	}

	// Movement pattern
	o.updateMovement(px, py)

	// Attack logic
	o.updateAttacks(px, py, game)
}

func (o *Owlbert) updateMovement(px, py float64) {
	// Sinusoidal vertical movement - biased upward (go up more, down less)
	o.y = o.baseY - 60 + math.Sin(o.movementTimer)*80

	// Horizontal drift towards player but stay more centered
	targetX := px - 150 // Stay to the right of player but not too far
	if targetX > 300 {
		targetX = 300 // Don't go too far left
	}
	if targetX < 450 {
		targetX = 450 // Stay on the right side but more centered
	}

	// Smooth movement towards target position
	o.x += (targetX - o.x) * 0.02

	// Boundary checking - keep more centered on screen
	screenWidth := 800.0
	minX := screenWidth * 0.25 // Minimum 25% from left (200px)
	maxX := screenWidth * 0.65 // Maximum 65% from left (520px)

	if o.x > maxX {
		o.x = maxX
	}
	if o.x < minX {
		o.x = minX
	}
}

func (o *Owlbert) updateAttacks(px, py float64, game *Game) {
	// Feather attack
	if o.featherTimer <= 0 && o.attackTellTimer == 0 {
		o.attackTellTimer = 30 // 0.5 second tell
		o.attackTellType = 1
		o.featherTimer = 90 // Reset timer (more frequent)
	}

	// Talon strike
	if o.talonTimer <= 0 && o.attackTellTimer == 0 {
		o.attackTellTimer = 30 // 0.5 second tell
		o.attackTellType = 2
		o.talonTimer = 420 // Reset timer (less frequent)
	}

	// Minion spawn
	if o.minionTimer <= 0 && o.attackTellTimer == 0 {
		o.attackTellTimer = 30 // 0.5 second tell
		o.attackTellType = 3
		o.minionTimer = 300 // Reset timer
	}

	// Execute attack after tell
	if o.attackTellTimer == 1 {
		switch o.attackTellType {
		case 1:
			o.shootFeather(game)
		case 2:
			o.talonStrike(px, py, game)
		case 3:
			o.spawnMinions(game)
		}
		o.attackTellType = 0
	}
}

func (o *Owlbert) shootFeather(game *Game) {
	feather := NewFeatherBullet(o.x+float64(o.frameWidth), o.y+float64(o.frameHeight)/2, game.featherImg)
	game.enemyBullets = append(game.enemyBullets, feather)
}

func (o *Owlbert) talonStrike(px, py float64, game *Game) {
	// Shoot arc of 3 feathers towards player
	centerX := o.x + float64(o.frameWidth)/2
	centerY := o.y + float64(o.frameHeight)/2

	// Calculate angles for 3-feather arc
	baseAngle := math.Atan2(py-centerY, px-centerX) // Angle to player
	spreadAngle := math.Pi / 6                      // 30 degree spread

	for i := 0; i < 3; i++ {
		angle := baseAngle - spreadAngle + spreadAngle*float64(i) // -30°, 0°, +30°
		speed := 4.0

		feather := NewFeatherBullet(centerX, centerY, game.featherImg)
		feather.SetVelocity(math.Cos(angle)*speed, math.Sin(angle)*speed)
		game.enemyBullets = append(game.enemyBullets, feather)
	}
}

func (o *Owlbert) spawnMinions(game *Game) {
	// Spawn 1-2 Flutternats
	numMinions := 1 + rand.IntN(2)
	for i := 0; i < numMinions; i++ {
		x := float64(100 + i*50)
		y := 100.0
		minion := CreateEnemy(FlutternatType, x, y, game.enemyImage, game.podImg)
		game.enemies = append(game.enemies, minion)
	}
}

func (o *Owlbert) Draw(screen *ebiten.Image) {
	if o.isDying {
		// Death animation - fade out
		alpha := max(0, 255-o.deathTimer*2)
		if alpha <= 0 {
			return
		}

		// Use sprite animation even in death
		frame := (o.frameCounter / 8) % 4
		sx := frame * o.frameWidth
		rect := image.Rect(sx, 0, sx+o.frameWidth, o.frameHeight)
		subImg := o.img.SubImage(rect).(*ebiten.Image)

		op := &ebiten.DrawImageOptions{}
		op.ColorScale.ScaleWithColor(color.RGBA{255, 255, 255, uint8(alpha)})
		op.GeoM.Scale(1.5, 1.5) // Scale to 1.5x size
		op.GeoM.Translate(o.x, o.y)
		screen.DrawImage(subImg, op)
		return
	}

	// Sprite animation
	frame := (o.frameCounter / 8) % 4
	sx := frame * o.frameWidth
	rect := image.Rect(sx, 0, sx+o.frameWidth, o.frameHeight)
	subImg := o.img.SubImage(rect).(*ebiten.Image)

	// Attack tell coloring
	op := &ebiten.DrawImageOptions{}

	if o.attackTellTimer > 0 {
		switch o.attackTellType {
		case 1: // Feather attack - white flash
			op.ColorScale.ScaleWithColor(color.RGBA{255, 255, 255, 255})
		case 2: // Talon strike - red flash
			op.ColorScale.ScaleWithColor(color.RGBA{255, 100, 100, 255})
		case 3: // Minion spawn - blue flash
			op.ColorScale.ScaleWithColor(color.RGBA{100, 100, 255, 255})
		}
	} else if o.hitFlashTimer > 0 {
		op.ColorScale.ScaleWithColor(color.RGBA{255, 100, 100, 255})
	}

	op.GeoM.Scale(1.5, 1.5) // Scale to 1.5x size
	op.GeoM.Translate(o.x, o.y)
	screen.DrawImage(subImg, op)

	// Always show health bar for boss
	o.drawHealthBar(screen)
}

func (o *Owlbert) drawHealthBar(screen *ebiten.Image) {
	scaledWidth := float64(o.frameWidth) * 1.5
	barHeight := 6.0
	barX := o.x
	barY := o.y - 18 // Adjusted for scaling

	// Background
	vector.FillRect(screen, float32(barX), float32(barY), float32(scaledWidth), float32(barHeight), color.RGBA{50, 50, 50, 255}, false)

	// Health
	healthPercent := float64(o.health) / float64(o.maxHealth)
	redWidth := scaledWidth * healthPercent
	vector.FillRect(screen, float32(barX), float32(barY), float32(redWidth), float32(barHeight), color.RGBA{255, 0, 0, 255}, false)

	// Border
	vector.StrokeRect(screen, float32(barX), float32(barY), float32(scaledWidth), float32(barHeight), 1, color.RGBA{255, 255, 255, 255}, false)
}

func (o *Owlbert) GetPosition() (x, y float64) {
	return o.x, o.y
}

func (o *Owlbert) GetBounds() (width, height float64) {
	return float64(o.frameWidth) * 1.5, float64(o.frameHeight) * 1.5
}

func (o *Owlbert) TakeDamage(amount int) {
	o.health -= amount
	o.hitFlashTimer = 10
	o.healthBarTimer = 60 // Show health bar longer for boss
}

func (o *Owlbert) IsDead() bool {
	return o.health <= 0 || o.isDying
}

func (o *Owlbert) OnDeath(game *Game) {
	if !o.isDying {
		o.isDying = true
		o.deathTimer = 0
	}
}

func (o *Owlbert) UpdateDeath(game *Game) {
	if o.isDying && o.deathTimer == 120 { // 2 seconds - create rewards
		o.CreateDeathRewards(game)
	}
	if o.isDying && o.deathTimer >= 480 { // 8 seconds - then allow level completion
		game.bossKilled = true
	}
}

func (o *Owlbert) IsDeathComplete() bool {
	return o.isDying && o.deathTimer >= 480 // 8 seconds
}

func (o *Owlbert) CreateDeathRewards(game *Game) {

	numCoins := 150 + rand.IntN(21)
	centerX := o.x + float64(o.frameWidth)/2
	centerY := o.y + float64(o.frameHeight)/2

	for i := 0; i < numCoins; i++ {
		angle := 2 * math.Pi * float64(i) / float64(numCoins)
		speed := 1.5 + rand.Float64()*2.0 // Reduced speed to prevent coins from going off screen
		vx := math.Cos(angle) * speed
		vy := math.Sin(angle) * speed

		size := SmallCoin
		if ShouldSpawnBigCoin(game.player.luck, 1) {
			size = BigCoin
		}

		coin := NewCoin(centerX, centerY, size, game.coinImg)
		coin.SetVelocity(vx, vy)
		game.coins = append(game.coins, coin)
	}
}
