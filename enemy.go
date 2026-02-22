package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Enemy interface {
	Update(playerX, playerY float64, game *Game)
	Draw(screen *ebiten.Image)
	GetPosition() (x, y float64)
	GetBounds() (width, height float64)
	TakeDamage(amount int)
	IsDead() bool
	OnDeath(game *Game)
}

type EnemyType string

const (
	FlutternatType EnemyType = "flutternat"
)

func CreateEnemy(enemyType EnemyType, x, y float64, image map[EnemyType]*ebiten.Image) Enemy {
	switch enemyType {
	case FlutternatType:
		return &Flutternat{
			x:           x,
			y:           y,
			img:         image[FlutternatType],
			speed:       2.0,
			health:      3,
			maxHealth:   3,
			frameWidth:  59,
			frameHeight: 32,
		}
	}
	return nil
}

func checkCollision(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}
