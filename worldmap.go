package main

import "github.com/hajimehoshi/ebiten/v2"

type WorldMap struct {
	gearsBg          *ebiten.Image
	upgradeLayer     *ebiten.Image
	levelSelectLayer *ebiten.Image
	lockIcon         *ebiten.Image
}

func NewWorldMap(assets *Assets) *WorldMap {
	return &WorldMap{
		gearsBg:          assets.GearsBg,
		upgradeLayer:     assets.UpgradeLayer,
		levelSelectLayer: assets.LevelSelectLayer,
		lockIcon:         assets.LockIcon,
	}
}

func (wm *WorldMap) Update(g *Game) {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		g.startLevel(GetLevel1())
	}
}

func (wm *WorldMap) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(wm.gearsBg, op)
}
