package main

import "github.com/hajimehoshi/ebiten/v2"

type Bullet interface {
	Update()
	Draw(screen *ebiten.Image)
	GetPosition() (float64, float64)
	GetDamage() int
}

type BaseBullet struct {
	x, y   float64
	img    *ebiten.Image
	damage int
}

type BulletImpl struct {
	BaseBullet
	speed float64
}

func NewBullet(x, y float64, img *ebiten.Image, damage int) *BulletImpl {
	return &BulletImpl{
		BaseBullet: BaseBullet{
			x:      x,
			y:      y,
			img:    img,
			damage: damage,
		},
		speed: 7.0,
	}
}

func (b *BulletImpl) Update() {
	b.x += b.speed
}

func (b *BulletImpl) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.x, b.y)
	screen.DrawImage(b.img, op)
}

func (b *BulletImpl) GetPosition() (float64, float64) {
	return b.x, b.y
}

func (b *BulletImpl) GetDamage() int {
	return b.damage
}

type SpreadBullet struct {
	BulletImpl
	vx, vy float64
}

func NewSpreadBullet(x, y float64, img *ebiten.Image, damage int, vx, vy float64) *SpreadBullet {
	return &SpreadBullet{
		BulletImpl: BulletImpl{
			BaseBullet: BaseBullet{
				x:      x,
				y:      y,
				img:    img,
				damage: damage,
			},
			speed: 7.0,
		},
		vx: vx,
		vy: vy,
	}
}

func (b *SpreadBullet) Update() {
	b.x += b.vx
	b.y += b.vy
}

func (b *SpreadBullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.x, b.y)
	screen.DrawImage(b.img, op)
}

func (b *SpreadBullet) GetPosition() (float64, float64) {
	return b.x, b.y
}

func (b *SpreadBullet) GetDamage() int {
	return b.damage
}

type PodBullet struct {
	BaseBullet
	vx, vy float64
}

func NewPodBullet(x, y float64, img *ebiten.Image, vx, vy float64) *PodBullet {
	return &PodBullet{
		BaseBullet: BaseBullet{
			x:      x,
			y:      y,
			img:    img,
			damage: 1,
		},
		vx: vx,
		vy: vy,
	}
}

func (b *PodBullet) Update() {
	b.x += b.vx
	b.y += b.vy
}

func (b *PodBullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.x, b.y)
	screen.DrawImage(b.img, op)
}

func (b *PodBullet) GetPosition() (float64, float64) {
	return b.x, b.y
}

func (b *PodBullet) GetDamage() int {
	return b.damage
}

type SporeBullet struct {
	BaseBullet
	vx, vy float64
}

func NewSporeBullet(x, y float64, img *ebiten.Image) *SporeBullet {
	return &SporeBullet{
		BaseBullet: BaseBullet{
			x:      x,
			y:      y,
			img:    img,
			damage: 1,
		},
		vx: -3.0, // Default speed to the left
		vy: 0.0,
	}
}

func (b *SporeBullet) SetVelocity(vx, vy float64) {
	b.vx = vx
	b.vy = vy
}

func (b *SporeBullet) Update() {
	b.x += b.vx
	b.y += b.vy
}

func (b *SporeBullet) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(b.x, b.y)
	screen.DrawImage(b.img, op)
}

func (b *SporeBullet) GetPosition() (float64, float64) {
	return b.x, b.y
}

func (b *SporeBullet) GetDamage() int {
	return b.damage
}
