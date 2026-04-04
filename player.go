package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	x, y         float64
	img          *ebiten.Image
	frameWidth   int
	frameHeight  int
	frameCounter int
	speedLevel   int
	fireCounter  int
	fireInterval int

	health          int
	maxHealth       int
	hitFlashTimer   int
	invincibleTimer int
	luck            int
	coins           int
}

func NewPlayer(img *ebiten.Image) *Player {
	return &Player{
		x:            50,
		y:            100,
		img:          img,
		frameWidth:   80,
		frameHeight:  80,
		frameCounter: 0,
		speedLevel:   1,
		fireInterval: 30,
		health:       3,
		maxHealth:    3,
		luck:         1,
		coins:        0,
	}
}

func (p *Player) getMovementSpeed() float64 {
	const slow = 0.03
	const fast = 1.0

	return slow + (fast-slow)*(float64(p.speedLevel-1)/6.0)
}

func (p *Player) Update(g *Game) {
	mx, my := ebiten.CursorPosition()

	targetX := float64(mx) - (float64(p.frameWidth) / 2)
	targetY := float64(my) - (float64(p.frameHeight) / 2)

	lerpFactor := p.getMovementSpeed()
	if lerpFactor > 1 {
		lerpFactor = 1
	}

	p.x += (targetX - p.x) * lerpFactor
	p.y += (targetY - p.y) * lerpFactor

	if p.x < 0 {
		p.x = 0
	} else if p.x > float64(ScreenWidth)-float64(p.frameWidth) {
		p.x = float64(ScreenWidth) - float64(p.frameWidth)
	}

	if p.y < 0 {
		p.y = 0
	} else if p.y > float64(ScreenHeight)-float64(p.frameHeight) {
		p.y = float64(ScreenHeight) - float64(p.frameHeight)
	}

	p.fireCounter++

	if p.fireCounter >= p.fireInterval {
		p.fireCounter = 0

		bx := p.x + float64(p.frameWidth)/2
		by := p.y + float64(p.frameHeight)/2 - 4

		newBullet := NewBullet(bx, by, g.bulletImg)
		g.bullets = append(g.bullets, newBullet)

	}
	p.frameCounter++

	if p.invincibleTimer > 0 {
		p.invincibleTimer--
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	if p.invincibleTimer > 0 {
		if (p.invincibleTimer/8)%2 == 1 {
			return
		}
	}
	// Shows from 0 or 1
	frame := (p.frameCounter / 10) % 2

	sx := frame * p.frameWidth
	rect := image.Rect(sx, 0, sx+p.frameWidth, p.frameHeight)
	subImg := p.img.SubImage(rect).(*ebiten.Image)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(p.x, p.y)
	screen.DrawImage(subImg, op)
}

func (p *Player) TakeDamage(amount int) {
	if p.invincibleTimer > 0 {
		return
	}
	p.health -= amount
	p.hitFlashTimer = 100
	p.invincibleTimer = 100
}

func (p *Player) IsDead() bool {
	return p.health <= 0
}
