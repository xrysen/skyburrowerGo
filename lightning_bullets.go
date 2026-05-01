package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// LightningBolt is a velocity-based projectile drawn using the bolt sprite.
type LightningBolt struct {
	x, y   float64
	vx, vy float64
	img    *ebiten.Image
	angle  float64
	scale  float64
}

func NewLightningBolt(x, y, vx, vy float64, img *ebiten.Image, scale float64) *LightningBolt {
	return &LightningBolt{x: x, y: y, vx: vx, vy: vy, img: img, angle: math.Atan2(vy, vx), scale: scale}
}

func (b *LightningBolt) Update() {
	b.x += b.vx
	b.y += b.vy
}

func (b *LightningBolt) Draw(screen *ebiten.Image) {
	if b.img != nil {
		w := float64(b.img.Bounds().Dx())
		h := float64(b.img.Bounds().Dy())
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-w/2, -h/2)
		op.GeoM.Rotate(b.angle - math.Pi/2) // bolt.png faces down
		op.GeoM.Scale(b.scale, b.scale)
		op.GeoM.Translate(b.x, b.y)
		screen.DrawImage(b.img, op)
	} else {
		vector.FillRect(screen, float32(b.x), float32(b.y), 12, 4, color.RGBA{255, 255, 0, 255}, false)
	}
}

func (b *LightningBolt) hitboxSize() float64 { return b.scale * 24 }

// GetPosition returns the top-left of the hitbox (centered on the bullet).
func (b *LightningBolt) GetPosition() (float64, float64) {
	h := b.hitboxSize()
	return b.x - h/2, b.y - h/2
}
func (b *LightningBolt) GetDamage() int                { return 1 }
func (b *LightningBolt) GetBounds() (float64, float64) { h := b.hitboxSize(); return h, h }

// CloudProjectile is a velocity-based projectile drawn as a white puff.
type CloudProjectile struct {
	x, y   float64
	vx, vy float64
}

func NewCloudProjectile(x, y, vx, vy float64) *CloudProjectile {
	return &CloudProjectile{x: x, y: y, vx: vx, vy: vy}
}

func (b *CloudProjectile) Update() {
	b.x += b.vx
	b.y += b.vy
}

func (b *CloudProjectile) Draw(screen *ebiten.Image) {
	vector.FillCircle(screen, float32(b.x), float32(b.y), 6, color.RGBA{200, 220, 255, 255}, false)
}

func (b *CloudProjectile) GetPosition() (float64, float64) { return b.x, b.y }
func (b *CloudProjectile) GetDamage() int                  { return 1 }

// ChainLightningBolt is a velocity-based aimed bolt drawn using the bolt sprite.
type ChainLightningBolt struct {
	x, y   float64
	vx, vy float64
	img    *ebiten.Image
	angle  float64
	scale  float64
}

func NewChainLightningBolt(x, y, vx, vy float64, img *ebiten.Image, scale float64) *ChainLightningBolt {
	return &ChainLightningBolt{x: x, y: y, vx: vx, vy: vy, img: img, angle: math.Atan2(vy, vx), scale: scale}
}

func (b *ChainLightningBolt) Update() {
	b.x += b.vx
	b.y += b.vy
}

func (b *ChainLightningBolt) Draw(screen *ebiten.Image) {
	if b.img != nil {
		w := float64(b.img.Bounds().Dx())
		h := float64(b.img.Bounds().Dy())
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-w/2, -h/2)
		op.GeoM.Rotate(b.angle - math.Pi/2) // bolt.png faces down
		op.GeoM.Scale(b.scale, b.scale)
		op.GeoM.Translate(b.x, b.y)
		screen.DrawImage(b.img, op)
	} else {
		// Jagged three-segment line fallback
		vector.StrokeLine(screen, float32(b.x), float32(b.y), float32(b.x-5), float32(b.y+3), 2, color.RGBA{255, 240, 0, 255}, false)
		vector.StrokeLine(screen, float32(b.x-5), float32(b.y+3), float32(b.x-10), float32(b.y-2), 2, color.RGBA{255, 240, 0, 255}, false)
		vector.StrokeLine(screen, float32(b.x-10), float32(b.y-2), float32(b.x-15), float32(b.y), 2, color.RGBA{255, 240, 0, 255}, false)
	}
}

func (b *ChainLightningBolt) hitboxSize() float64 { return b.scale * 24 }

// GetPosition returns the top-left of the hitbox (centered on the bullet).
func (b *ChainLightningBolt) GetPosition() (float64, float64) {
	h := b.hitboxSize()
	return b.x - h/2, b.y - h/2
}
func (b *ChainLightningBolt) GetDamage() int                { return 1 }
func (b *ChainLightningBolt) GetBounds() (float64, float64) { h := b.hitboxSize(); return h, h }

// ElectricalRing expands outward from its spawn point.
// Collision is active while radius is in [8, 20].
// At radius >= 20 the ring moves off-screen so the game loop removes it.
const (
	electricalRingMinRadius = 8.0
	electricalRingMaxRadius = 20.0
	electricalRingGrowRate  = 2.0
)

type ElectricalRing struct {
	cx, cy float64
	vx     float64
	radius float64
}

func NewElectricalRing(cx, cy, vx float64) *ElectricalRing {
	return &ElectricalRing{cx: cx, cy: cy, vx: vx}
}

func (b *ElectricalRing) Update() {
	b.cx += b.vx
	if b.radius < electricalRingMaxRadius {
		b.radius += electricalRingGrowRate
	}
}

func (b *ElectricalRing) Draw(screen *ebiten.Image) {
	if b.radius <= 0 {
		return
	}
	vector.StrokeCircle(screen, float32(b.cx), float32(b.cy), float32(b.radius), 2, color.RGBA{100, 180, 255, 220}, false)
}

func (b *ElectricalRing) GetPosition() (float64, float64) { return b.cx, b.cy }
func (b *ElectricalRing) GetDamage() int                  { return 1 }
func (b *ElectricalRing) GetRadius() float64              { return b.radius }

// IsActive returns true while the ring is in its damage window (radius 8–20).
func (b *ElectricalRing) IsActive() bool {
	return b.radius >= electricalRingMinRadius && b.radius < electricalRingMaxRadius
}

// Shockwave expands from its spawn point.
// Collision is active while radius is in [10, 40].
// At radius >= 40 the shockwave moves off-screen.
const (
	shockwaveMinRadius = 10.0
	shockwaveMaxRadius = 40.0
	shockwaveGrowRate  = 0.8
)

type Shockwave struct {
	cx, cy float64
	x, y   float64
	radius float64
}

func NewShockwave(cx, cy float64) *Shockwave {
	return &Shockwave{cx: cx, cy: cy, x: cx, y: cy}
}

func (b *Shockwave) Update() {
	if b.radius < shockwaveMaxRadius {
		b.radius += shockwaveGrowRate
	}
	if b.radius >= shockwaveMaxRadius {
		b.x = -9999
		b.y = -9999
	}
}

func (b *Shockwave) Draw(screen *ebiten.Image) {
	if b.x < -100 {
		return
	}
	vector.FillCircle(screen, float32(b.cx), float32(b.cy), float32(b.radius), color.RGBA{255, 200, 50, 80}, false)
	vector.StrokeCircle(screen, float32(b.cx), float32(b.cy), float32(b.radius), 3, color.RGBA{255, 160, 0, 200}, false)
}

func (b *Shockwave) GetPosition() (float64, float64) { return b.x, b.y }
func (b *Shockwave) GetDamage() int                  { return 2 }
func (b *Shockwave) GetRadius() float64              { return b.radius }

// IsActive returns true while the shockwave is in its damage window (radius 10–40).
func (b *Shockwave) IsActive() bool {
	return b.radius >= shockwaveMinRadius && b.radius < shockwaveMaxRadius
}
