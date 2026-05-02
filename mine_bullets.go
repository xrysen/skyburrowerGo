package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// DrillBit is a slow, heavy projectile drawn as a small rotated rectangle.
type DrillBit struct {
	x, y   float64
	vx, vy float64
	angle  float64
}

func NewDrillBit(x, y, vx, vy float64) *DrillBit {
	return &DrillBit{x: x, y: y, vx: vx, vy: vy, angle: math.Atan2(vy, vx)}
}

func (b *DrillBit) Update() {
	b.x += b.vx
	b.y += b.vy
}

func (b *DrillBit) Draw(screen *ebiten.Image) {
	img := ebiten.NewImage(8, 4)
	img.Fill(color.RGBA{80, 80, 90, 255})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-4, -2)
	op.GeoM.Rotate(b.angle)
	op.GeoM.Translate(b.x, b.y)
	screen.DrawImage(img, op)
}

func (b *DrillBit) GetPosition() (float64, float64) { return b.x, b.y }
func (b *DrillBit) GetDamage() int                  { return 3 }
func (b *DrillBit) GetBounds() (float64, float64)   { return 8, 4 }

// FuseSpark is a medium-speed projectile drawn as a small orange circle.
type FuseSpark struct {
	x, y   float64
	vx, vy float64
}

func NewFuseSpark(x, y, vx, vy float64) *FuseSpark {
	return &FuseSpark{x: x, y: y, vx: vx, vy: vy}
}

func (b *FuseSpark) Update() {
	b.x += b.vx
	b.y += b.vy
}

func (b *FuseSpark) Draw(screen *ebiten.Image) {
	vector.FillCircle(screen, float32(b.x), float32(b.y), 4, color.RGBA{255, 140, 0, 255}, false)
}

func (b *FuseSpark) GetPosition() (float64, float64) { return b.x, b.y }
func (b *FuseSpark) GetDamage() int                  { return 2 }
func (b *FuseSpark) GetBounds() (float64, float64)   { return 8, 8 }
