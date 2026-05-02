package main

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Layer struct {
	img   *ebiten.Image
	speed float64
	x     float64
}

type rainDrop struct {
	x, y   float64
	dy     float64
	length float64
}

type Background struct {
	layers          []*Layer
	weather         WeatherType
	drops           []rainDrop
	lightningTimer  int
	lightningAlpha  float64
	flashDecay      float64
	flashOriginX    float64
	pendingFlickers int
	rainIntensity   float64
}

func NewBackground(layers []*Layer, weather WeatherType) *Background {
	b := &Background{layers: layers, weather: weather}
	if weather == WeatherRain {
		b.drops = make([]rainDrop, 100)
		for i := range b.drops {
			b.drops[i] = rainDrop{
				x:      rand.Float64() * ScreenWidth,
				y:      rand.Float64() * ScreenHeight,
				dy:     3.0 + rand.Float64()*2.0,
				length: 6.0 + rand.Float64()*4.0,
			}
		}
		b.lightningTimer = 180 + rand.Intn(301)
		b.rainIntensity = 1.0
	}
	return b
}

func (b *Background) Update() {
	for _, l := range b.layers {
		l.x -= l.speed
		w := float64(l.img.Bounds().Dx())
		if l.x < -w {
			l.x = 0
		}
	}
	for i := range b.drops {
		b.drops[i].y += b.drops[i].dy * b.rainIntensity
		b.drops[i].x -= 0.8 * b.rainIntensity
		if b.drops[i].y > ScreenHeight || b.drops[i].x < 0 {
			b.drops[i].x = rand.Float64() * ScreenWidth
			b.drops[i].y = 0
		}
	}
	if b.weather == WeatherRain {
		b.rainIntensity -= 0.03
		if b.rainIntensity < 1.0 {
			b.rainIntensity = 1.0
		}
		b.lightningTimer--
		if b.lightningTimer <= 0 {
			b.flashOriginX = rand.Float64() * ScreenWidth
			if b.pendingFlickers > 0 {
				// pre-flicker: brief dim blip
				b.lightningAlpha = 0.05
				b.flashDecay = 0.55
				b.pendingFlickers--
				b.lightningTimer = 3 + rand.Intn(4)
			} else {
				// main flash: mostly quick, occasionally a slow linger
				if rand.Float64() < 0.25 {
					b.lightningAlpha = 0.12
					b.flashDecay = 0.84 + rand.Float64()*0.04 // long (18–25 frames)
				} else {
					b.lightningAlpha = 0.15
					b.flashDecay = 0.75 // quick (~8 frames)
				}
				b.rainIntensity = 2.5
				b.lightningTimer = 180 + rand.Intn(301)
				if rand.Float64() < 0.4 {
					b.pendingFlickers = 1 + rand.Intn(2)
				}
			}
		}
		b.lightningAlpha *= b.flashDecay
		if b.lightningAlpha < 0.005 {
			b.lightningAlpha = 0
		}
	}
}

func (b *Background) DrawLightningFlash(screen *ebiten.Image) {
	if b.lightningAlpha == 0 {
		return
	}
	const strips = 10
	stripW := float32(ScreenWidth) / strips
	for i := 0; i < strips; i++ {
		cx := float32(i)*stripW + stripW/2
		dist := cx - float32(b.flashOriginX)
		if dist < 0 {
			dist = -dist
		}
		a := b.lightningAlpha * (1.0 - 0.65*float64(dist)/ScreenWidth)
		if a <= 0 {
			continue
		}
		vector.FillRect(screen, float32(i)*stripW, 0, stripW, float32(ScreenHeight),
			color.RGBA{R: 255, G: 255, B: 255, A: uint8(a * 255)}, false)
	}
}

func (b *Background) DrawRain(screen *ebiten.Image) {
	rainColor := color.RGBA{R: 180, G: 210, B: 255, A: 120}
	for _, d := range b.drops {
		x1 := float32(d.x)
		y1 := float32(d.y)
		x2 := float32(d.x - d.length*0.2)
		y2 := float32(d.y + d.length)
		vector.StrokeLine(screen, x1, y1, x2, y2, 1, rainColor, false)
	}
}

func (b *Background) Draw(screen *ebiten.Image, layerIndex int) {
	l := b.layers[layerIndex]
	w := float64(l.img.Bounds().Dx())

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(l.x, 0)
	screen.DrawImage(l.img, op)

	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(l.x+w, 0)
	screen.DrawImage(l.img, op2)
}
