package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GameState int

const (
	StateLevel GameState = iota
	StateLevelComplete
	StateGameOver
)

const (
	FadeSpeed = 0.02
)

type Game struct {
	player         *Player
	background     *Background
	bullets        []*Bullet
	bulletImg      *ebiten.Image
	enemies        []Enemy
	enemyImage     map[EnemyType]*ebiten.Image
	spawnTimers    map[EnemyType]int
	spawnCounts    map[EnemyType]int
	currentLevel   *LevelConfig
	levelTimer     int
	hud            *HUD
	coins          []*Coin
	coinImg        *ebiten.Image
	coinSpawnTimer int
	state          GameState
	bossKilled     bool

	fadeAlpha float64
	fadeSpeed float64
	isFading  bool
	fadeIn    bool
}

func (g *Game) Update() error {

	if g.isFading {
		if g.fadeIn {
			// Fade in
			g.fadeAlpha -= g.fadeSpeed
			if g.fadeAlpha <= 0 {
				g.fadeAlpha = 0
				g.isFading = false
			}
		} else {
			// Fade out
			g.fadeAlpha += g.fadeSpeed
			if g.fadeAlpha >= 1.0 {
				g.fadeAlpha = 1.0
				g.isFading = false
				if g.state == StateLevelComplete {
					g.transitionToNextLevel()
					return nil
				}
			}
		}
	}

	// Always update visuals (even during fade-out)
	g.background.Update()
	g.player.Update(g)

	var activeBullets []*Bullet
	for _, b := range g.bullets {
		b.Update()
		if b.x < 640 {
			activeBullets = append(activeBullets, b)
		}
	}
	g.bullets = activeBullets

	var activeEnemies []Enemy
	for _, e := range g.enemies {
		e.Update(g.player.x, g.player.y, g)
		ex, _ := e.GetPosition()
		if ex > -100 && !e.IsDead() {
			activeEnemies = append(activeEnemies, e)
		}
	}
	g.enemies = activeEnemies

	// Only update gameplay logic if level is active
	if g.state != StateLevelComplete {

		g.levelTimer++

		switch g.currentLevel.EndCondition {
		case EndOnTimer:
			fadeStartTime := g.currentLevel.Duration - FadeOutDuration
			if g.levelTimer >= fadeStartTime && g.state == StateLevel {
				g.completeLevel()
			}
		case EndOnBossDeath:
			if g.bossKilled && g.state == StateLevel {
				g.completeLevel()
			}
		}

		for _, spawnCfg := range g.currentLevel.SpawnConfigs {
			if g.levelTimer < spawnCfg.StartFrame {
				continue
			}

			if spawnCfg.EndFrame > 0 && g.levelTimer > spawnCfg.EndFrame {
				continue
			}
			g.spawnTimers[spawnCfg.EnemyType]++
			if g.spawnTimers[spawnCfg.EnemyType] >= spawnCfg.SpawnRate {
				g.spawnTimers[spawnCfg.EnemyType] = 0
				count := spawnCfg.MinSpawns
				if spawnCfg.MaxSpawns > spawnCfg.MinSpawns {
					count += rand.IntN(spawnCfg.MaxSpawns - spawnCfg.MinSpawns + 1)
				}

				for i := 0; i < count; i++ {
					g.spawnEnemy(spawnCfg)
				}
			}
		}

		g.coinSpawnTimer++
		spawnRate := 120 - (g.player.luck * 10)
		if spawnRate < 50 {
			spawnRate = 50
		}
		if g.coinSpawnTimer >= spawnRate {
			g.coinSpawnTimer = 0

			maxCoins := 2 + (g.player.luck / 2)
			if maxCoins > 5 {
				maxCoins = 5 // Cap at 5 coins
			}

			numCoins := 1 + rand.IntN(maxCoins)
			for i := 0; i < numCoins; i++ {
				y := rand.Float64() * 340
				size := SmallCoin
				if ShouldSpawnBigCoin(g.player.luck, 0) {
					size = BigCoin
				}
				coin := NewCoin(650+rand.Float64()*100, y, size, g.coinImg)
				g.coins = append(g.coins, coin)
			}
		}

		var activeCoins []*Coin
		for _, c := range g.coins {
			c.Update()
			if !c.collected && c.x > -50 {
				activeCoins = append(activeCoins, c)
			}
		}
		g.coins = activeCoins

		for _, c := range g.coins {
			if c.collected {
				continue
			}
			coinW := float64(c.frameWidth) * c.scale
			coinH := float64(c.frameHeight) * c.scale
			playerW := float64(g.player.frameWidth)
			playerH := float64(g.player.frameHeight)

			if checkCollision(g.player.x, g.player.y, playerW, playerH, c.x, c.y, coinW, coinH) {
				c.collected = true
				g.player.coins += c.value
			}
		}

		for _, b := range g.bullets {
			for _, e := range g.enemies {
				ex, ey := e.GetPosition()
				ew, eh := e.GetBounds()
				if checkCollision(b.x, b.y, 8, 8, ex, ey, ew, eh) {
					e.TakeDamage(1)
					b.x = 1000
					b.y = 1000
					if e.IsDead() {
						e.OnDeath(g)
					}
				}
			}

			for _, e := range g.enemies {
				ex, ey := e.GetPosition()
				ew, eh := e.GetBounds()

				hitboxWidth := float64(g.player.frameWidth) * 0.7
				hitboxHeight := float64(g.player.frameHeight) * 0.7
				offsetX := float64(g.player.frameWidth) * 0.15
				offsetY := float64(g.player.frameHeight) * 0.15

				if checkCollision(g.player.x+offsetX, g.player.y+offsetY, hitboxWidth, hitboxHeight, ex, ey, ew, eh) {
					g.player.TakeDamage(1)

					if g.player.IsDead() {
						g.gameOver()
					}
				}
			}
		}
	}

	return nil
}

func (g *Game) spawnEnemy(cfg SpawnConfig) {
	y := 180.0
	if cfg.RandomY {
		y = rand.Float64() * 360
	}

	enemy := CreateEnemy(cfg.EnemyType, 650+rand.Float64()*100, y, g.enemyImage)
	g.enemies = append(g.enemies, enemy)
	g.spawnCounts[cfg.EnemyType]++
}

func (g *Game) completeLevel() {
	g.state = StateLevelComplete

	g.fadeAlpha = 0.0
	g.fadeSpeed = FadeSpeed
	g.isFading = true
	g.fadeIn = false
}

func (g *Game) transitionToNextLevel() {
	if g.currentLevel.NextLevel != nil {
		nextLevel := g.currentLevel.NextLevel()
		g.loadLevel(nextLevel)
		g.state = StateLevel
	} else {
		g.state = StateGameOver
		fmt.Println("Game Complete!")
	}
}

func (g *Game) loadLevel(level *LevelConfig) {
	g.currentLevel = level
	g.levelTimer = 0
	g.enemies = []Enemy{}
	g.spawnTimers = make(map[EnemyType]int)
	g.spawnCounts = make(map[EnemyType]int)
	g.bossKilled = false

	// Start fade in
	g.fadeAlpha = 1.0
	g.fadeSpeed = FadeSpeed
	g.isFading = true
	g.fadeIn = true

	g.background = &Background{
		layers: []*Layer{
			{img: loadImage(level.BackgroundPaths[0]), speed: 0.5},
			{img: loadImage(level.BackgroundPaths[1]), speed: 1.0},
			{img: loadImage(level.BackgroundPaths[2]), speed: 1.5},
			{img: loadImage(level.BackgroundPaths[3]), speed: 4.0},
		},
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := 0; i < 3; i++ {
		g.background.Draw(screen, i)
	}

	for _, c := range g.coins {
		c.Draw(screen)
	}

	for _, b := range g.bullets {
		b.Draw(screen)
	}

	for _, e := range g.enemies {
		e.Draw(screen)
	}

	g.player.Draw(screen)

	g.background.Draw(screen, 3)

	if g.fadeAlpha > 0 {
		vector.FillRect(screen, 0, 0, 640, 360, color.RGBA{0, 0, 0, uint8(g.fadeAlpha * 255)}, false)
	}

	g.hud.Draw(screen, g.player.health, g.player.maxHealth, g.player.coins)

}

func (g *Game) Layout(w, h int) (int, int) {
	return 640, 360
}

func loadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		panic(fmt.Sprintf("failed to load image: %v", err))
	}
	return img
}

func main() {

	fontImg := loadImage("Assets/UI/saikyoFonto.png")
	font := NewBitmapFont(fontImg, 18, 18)

	coinImg := loadImage("Assets/Items/coin.png")

	level := GetLevel1()

	pImg := loadImage("Assets/Player/MeadowSprite-sheet.png")
	bg0 := loadImage(level.BackgroundPaths[0])
	bg1 := loadImage(level.BackgroundPaths[1])
	bg2 := loadImage(level.BackgroundPaths[2])
	bg3 := loadImage(level.BackgroundPaths[3])
	bImg := loadImage("Assets/Bullets/seedShot.png")
	heartImg := loadImage("Assets/UI/heart.png")
	backgroundImg := loadImage("Assets/UI/ui.png")

	enemyImages := map[EnemyType]*ebiten.Image{
		FlutternatType: loadImage("Assets/Enemies/Flutternat/flutterNat.png"),
	}

	game := &Game{
		player: NewPlayer(pImg),
		hud:    NewHUD(backgroundImg, heartImg, font),
		background: &Background{
			layers: []*Layer{
				{img: bg0, speed: 0.5},
				{img: bg1, speed: 1.0},
				{img: bg2, speed: 1.5},
				{img: bg3, speed: 4.0},
			},
		},
		bulletImg:    bImg,
		enemyImage:   enemyImages,
		spawnTimers:  make(map[EnemyType]int),
		spawnCounts:  make(map[EnemyType]int),
		currentLevel: level,
		coinImg:      coinImg,

		fadeAlpha: 1.0,
		fadeSpeed: FadeSpeed,
		isFading:  true,
		fadeIn:    true,
	}

	ebiten.SetWindowSize(640, 360)
	ebiten.SetWindowTitle("Skyburrower")

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}

func (g *Game) gameOver() {
	g.state = StateGameOver
	fmt.Println("Game Over!")
}
