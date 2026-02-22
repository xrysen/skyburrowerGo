package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	player     *Player
	background *Background
	bullets    []*Bullet
	bulletImg  *ebiten.Image
}

func (g *Game) Update() error {
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

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i := 0; i < 3; i++ {
		g.background.Draw(screen, i)
	}

	for _, b := range g.bullets {
		b.Draw(screen)
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
	pImg := loadImage("Assets/Player/MeadowSprite-sheet.png")
	bg0 := loadImage("Levels/Level1/lvl1-1.png")
	bg1 := loadImage("Levels/Level1/lvl1-2.png")
	bg2 := loadImage("Levels/Level1/lvl1-3.png")
	bg3 := loadImage("Levels/Level1/lvl1-4.png")
	bImg := loadImage("Assets/Bullets/seedShot.png")

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
		bulletImg: bImg,
	}

	ebiten.SetWindowSize(640, 360)
	ebiten.SetWindowTitle("Skyburrower")

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

}
