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

type UpgradeType int

const (
	UpgradeHealth UpgradeType = iota
	UpgradeBulletStrength
	UpgradeBulletSpeed
	UpgradeBulletCount
	UpgradeSpeed
	UpgradeMagnetism
	UpgradeLuck
	UpgradeCount // Total number of upgrades
)

type Upgrade struct {
	Level int // 0 = empty, 1+ = filled levels
}

const (
	StateLevel GameState = iota
	StateLevelComplete
	StateGameOver
)

const (
	FadeSpeed = 0.02
)

type Game struct {
	currentScreen     Screen
	worldMap          *WorldMap
	player            *Player
	background        *Background
	bullets           []*Bullet
	bulletImg         *ebiten.Image
	enemies           []Enemy
	enemyImage        map[EnemyType]*ebiten.Image
	spawnTimers       map[EnemyType]int
	spawnCounts       map[EnemyType]int
	currentLevel      *LevelConfig
	levelTimer        int
	hud               *HUD
	coins             []*Coin
	coinImg           *ebiten.Image
	coinSpawnTimer    int
	levelCarrots      []*LevelCarrot
	carrotImg         *ebiten.Image
	carrotSpawnFrames [CarrotsPerLevel]int
	carrotSpawned     [CarrotsPerLevel]bool
	state             GameState
	bossKilled        bool
	upgrades          [UpgradeCount]Upgrade

	fadeAlpha float64
	fadeSpeed float64
	isFading  bool
	fadeIn    bool

	// Levels 1..highestUnlockedLevel are playable on the world map (start: only 1).
	highestUnlockedLevel int

	// Per world level (index 0 = level 1): bits 0–4 = which of the 5 bonus carrots
	// have ever been collected (merged across replays). See CarrotsPerLevel.
	levelCarrotMask [WorldLevelCount]uint8
	// Bits collected during the current run; merged into levelCarrotMask when the
	// level ends and returns to the map. Gameplay sets bits via CollectLevelCarrot.
	runLevelCarrotMask uint8
}

func (g *Game) Update() error {

	switch g.currentScreen {
	case ScreenWorldMap, ScreenGameOver:
		ebiten.SetCursorMode(ebiten.CursorModeVisible)
	case ScreenPlaying:
		ebiten.SetCursorMode(ebiten.CursorModeHidden)
	}

	switch g.currentScreen {
	case ScreenWorldMap:
		g.worldMap.Update(g)
	case ScreenPlaying:
		g.updatePlaying()
	case ScreenGameOver:
		g.updateGameOver()
	}
	return nil
}

func (g *Game) updatePlaying() {
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
				} else if g.state == StateGameOver {
					g.transitionToWorldMapFromGameOver()
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
		if b.x < ScreenWidth {
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
	if g.state != StateLevelComplete && g.state != StateGameOver {

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
				y := rand.Float64() * float64(ScreenHeight)
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

		for i := 0; i < CarrotsPerLevel; i++ {
			if g.carrotSpawned[i] {
				continue
			}
			if g.levelTimer < g.carrotSpawnFrames[i] {
				continue
			}
			cy := rand.Float64() * float64(max(1, ScreenHeight-g.carrotImg.Bounds().Dy()))
			carrot := NewLevelCarrot(650+rand.Float64()*100, cy, i, g.carrotImg)
			g.levelCarrots = append(g.levelCarrots, carrot)
			g.carrotSpawned[i] = true
		}

		var activeCarrots []*LevelCarrot
		for _, c := range g.levelCarrots {
			c.Update()
			if !c.collected && c.x > -50 {
				activeCarrots = append(activeCarrots, c)
			}
		}
		g.levelCarrots = activeCarrots

		for _, c := range g.levelCarrots {
			if c.collected {
				continue
			}
			cw := float64(c.frameWidth)
			ch := float64(c.frameHeight)
			playerW := float64(g.player.frameWidth)
			playerH := float64(g.player.frameHeight)
			if checkCollision(g.player.x, g.player.y, playerW, playerH, c.x, c.y, cw, ch) {
				c.collected = true
				g.CollectLevelCarrot(c.slotIndex)
			}
		}

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
}

func (g *Game) startLevel(level *LevelConfig) {
	g.loadLevel(level)
	g.state = StateLevel
	g.currentScreen = ScreenPlaying
}

func (g *Game) updateGameOver() {
	// Game over fade happens in updatePlaying()
	// This function remains empty since we never actually switch to ScreenGameOver
}

func (g *Game) drawGameOver(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Game Over")
}

func (g *Game) spawnEnemy(cfg SpawnConfig) {
	y := float64(ScreenHeight) / 2
	if cfg.RandomY {
		y = rand.Float64() * ScreenHeight
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
	if g.currentLevel != nil && g.currentLevel.WorldLevel > 0 {
		next := min(g.currentLevel.WorldLevel+1, 20)
		g.highestUnlockedLevel = max(g.highestUnlockedLevel, next)
	}
	g.commitRunCarrotProgress()
	g.currentScreen = ScreenWorldMap
}

func (g *Game) transitionToWorldMapFromGameOver() {
	// Reset player state (but keep coins as per requirement)
	g.player.health = g.player.maxHealth
	g.player.x = 50
	g.player.y = 100
	g.player.invincibleTimer = 0
	g.player.hitFlashTimer = 0

	// Clear all active game objects
	g.enemies = []Enemy{}
	g.bullets = []*Bullet{}
	g.coins = []*Coin{}
	g.levelCarrots = []*LevelCarrot{}

	// Reset timers and counters
	g.levelTimer = 0
	g.spawnTimers = make(map[EnemyType]int)
	g.spawnCounts = make(map[EnemyType]int)
	g.coinSpawnTimer = 0
	g.bossKilled = false

	// Reset carrot spawn tracking
	for i := range g.carrotSpawned {
		g.carrotSpawned[i] = false
	}

	// Do NOT commit carrot progress (player failed the level)
	g.runLevelCarrotMask = 0

	// Do NOT unlock next level (player didn't complete)

	// Return to worldmap
	g.currentScreen = ScreenWorldMap
	g.state = StateLevel // Reset state for next level attempt
}

func (g *Game) loadLevel(level *LevelConfig) {
	g.currentLevel = level
	g.runLevelCarrotMask = 0
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

	g.levelCarrots = nil
	g.carrotSpawnFrames = planCarrotSpawnFrames(level)
	for i := range g.carrotSpawned {
		g.carrotSpawned[i] = false
	}
}

// commitRunCarrotProgress ORs this run’s carrot bits into the persistent per-level mask.
func (g *Game) commitRunCarrotProgress() {
	if g.currentLevel == nil {
		return
	}
	wl := g.currentLevel.WorldLevel
	if wl < 1 || wl > WorldLevelCount {
		return
	}
	const carrotBitsMask = (1 << CarrotsPerLevel) - 1
	g.levelCarrotMask[wl-1] |= g.runLevelCarrotMask & uint8(carrotBitsMask)
}

// CollectLevelCarrot records bonus carrot `index` (0–4) for the current level this run.
// Bit/layout order: 0 = left of number, 1–3 = below (L→R), 4 = right of number.
func (g *Game) CollectLevelCarrot(index int) {
	if g.currentScreen != ScreenPlaying || g.currentLevel == nil {
		return
	}
	if index < 0 || index >= CarrotsPerLevel {
		return
	}
	g.runLevelCarrotMask |= uint8(1) << index
}

func (g *Game) Draw(screen *ebiten.Image) {

	switch g.currentScreen {
	case ScreenWorldMap:
		g.worldMap.Draw(screen, g)
	case ScreenPlaying:
		g.drawPlaying(screen)
	case ScreenGameOver:
		g.drawGameOver(screen)
	}
}

func (g *Game) drawPlaying(screen *ebiten.Image) {
	for i := 0; i < 3; i++ {
		g.background.Draw(screen, i)
	}

	for _, c := range g.coins {
		c.Draw(screen)
	}

	for _, c := range g.levelCarrots {
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
		vector.FillRect(screen, 0, 0, float32(ScreenWidth), float32(ScreenHeight), color.RGBA{0, 0, 0, uint8(g.fadeAlpha * 255)}, false)
	}

	g.hud.Draw(screen, g.player.health, g.player.maxHealth, g.player.coins, g.runLevelCarrotMask)
}

func (g *Game) Layout(w, h int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func loadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		panic(fmt.Sprintf("failed to load image: %v", err))
	}
	return img
}

func main() {

	assets := LoadAssets()
	font := NewBitmapFont(assets.FontImg, 18, 18)

	game := &Game{
		currentScreen:        ScreenWorldMap,
		highestUnlockedLevel: 1,
		worldMap:             NewWorldMap(assets, font),
		player:               NewPlayer(assets.PlayerImg),
		hud:                  NewHUD(assets.HudBg, assets.HeartImg, assets.LsCarrotEmpty, assets.LsCarrotFull, font),
		bulletImg:            assets.BulletImg,
		enemyImage:           assets.EnemyImages,
		spawnTimers:          make(map[EnemyType]int),
		spawnCounts:          make(map[EnemyType]int),
		coinImg:              assets.CoinImg,
		carrotImg:            assets.CarrotImg,
	}

	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Skyburrower")

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}

func (g *Game) gameOver() {
	g.state = StateGameOver

	// Start fadeout (same as completeLevel)
	g.fadeAlpha = 0.0
	g.fadeSpeed = FadeSpeed
	g.isFading = true
	g.fadeIn = false
}
