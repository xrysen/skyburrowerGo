package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Assets struct {
	PlayerImg *ebiten.Image
	BulletImg *ebiten.Image

	EnemyImages map[EnemyType]*ebiten.Image

	HeartImg *ebiten.Image
	HudBg    *ebiten.Image
	FontImg  *ebiten.Image
	CoinImg  *ebiten.Image

	GearsBg           *ebiten.Image
	UpgradeLayer      *ebiten.Image
	LevelSelectLayer  *ebiten.Image
	LevelSelectButton *ebiten.Image
	LockIcon          *ebiten.Image
	LevelDigits       [10]*ebiten.Image

	LevelBgs map[string]*ebiten.Image
}

func LoadAssets() *Assets {
	a := &Assets{}

	a.PlayerImg = loadImage("Assets/Player/MeadowSprite-sheet.png")
	a.BulletImg = loadImage("Assets/Bullets/seedShot.png")

	a.EnemyImages = map[EnemyType]*ebiten.Image{
		FlutternatType: loadImage("Assets/Enemies/Flutternat/flutterNat.png"),
	}

	a.HeartImg = loadImage("Assets/UI/heart.png")
	a.HudBg = loadImage("Assets/UI/ui.png")
	a.FontImg = loadImage("Assets/UI/saikyoFonto.png")
	a.CoinImg = loadImage("Assets/Items/coin.png")

	a.GearsBg = loadImage("Assets/WorldMap/gearsBg.png")
	a.UpgradeLayer = loadImage("Assets/WorldMap/UpgradeLayer.png")
	a.LevelSelectLayer = loadImage("Assets/WorldMap/wm.png")
	a.LevelSelectButton = loadImage("Assets/WorldMap/levelSelectButton.png")
	a.LockIcon = loadImage("Assets/WorldMap/lockIcon.png")
	for d := 0; d < 10; d++ {
		a.LevelDigits[d] = loadImage(fmt.Sprintf("Assets/WorldMap/%d.png", d))
	}

	a.LevelBgs = map[string]*ebiten.Image{
		"Level1": loadImage("Levels/Level1/lvl1-1.png"),
		"Level2": loadImage("Levels/Level1/lvl1-2.png"),
		"Level3": loadImage("Levels/Level1/lvl1-3.png"),
		"Level4": loadImage("Levels/Level1/lvl1-4.png"),
	}

	return a
}
