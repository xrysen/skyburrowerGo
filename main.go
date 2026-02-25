package main

import (
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GameState int

const (
	StateLevel GameState = iota
	StateLevelComplete
	StateGameOver
)

type Game struct {
	player       *Player
	background   *Background
	bullets      []*Bullet
	bulletImg    *ebiten.Image
	enemies      []Enemy
	enemyImage   map[EnemyType]*ebiten.Image
	spawnTimers  map[EnemyType]int
	spawnCounts  map[EnemyType]int
	currentLevel *LevelConfig
	levelTimer   int
	state        GameState
	bossKilled   bool
}

func (g *Game) Update() error {
	g.levelTimer++
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

	switch g.currentLevel.EndCondition {
	case EndOnTimer:
		if g.levelTimer >= g.currentLevel.Duration {
			g.completeLevel()
		}
	case EndOnBossDeath:
		if g.bossKilled {
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

	var activeEnemies []Enemy
	for _, e := range g.enemies {
		e.Update(g.player.x, g.player.y, g)
		ex, _ := e.GetPosition()
		if ex > -100 && !e.IsDead() {
			activeEnemies = append(activeEnemies, e)
		}
	}

	g.enemies = activeEnemies

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
}

func (g *Game) loadLevel(level *LevelConfig) {
	g.currentLevel = level
	g.levelTimer = 0
	g.enemies = []Enemy{}
	g.spawnTimers = make(map[EnemyType]int)
	g.spawnCounts = make(map[EnemyType]int)
	g.bossKilled = false

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

	for _, b := range g.bullets {
		b.Draw(screen)
	}

	for _, e := range g.enemies {
		e.Draw(screen)
	}

	g.player.Draw(screen)

	g.background.Draw(screen, 3)
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
	level := GetLevel1()

	pImg := loadImage("Assets/Player/MeadowSprite-sheet.png")
	bg0 := loadImage(level.BackgroundPaths[0])
	bg1 := loadImage(level.BackgroundPaths[1])
	bg2 := loadImage(level.BackgroundPaths[2])
	bg3 := loadImage(level.BackgroundPaths[3])
	bImg := loadImage("Assets/Bullets/seedShot.png")

	enemyImages := map[EnemyType]*ebiten.Image{
		FlutternatType: loadImage("Assets/Enemies/Flutternat/flutterNat.png"),
	}

	game := &Game{
		player: NewPlayer(pImg),
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
	}

	ebiten.SetWindowSize(640, 360)
	ebiten.SetWindowTitle("Skyburrower")

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
