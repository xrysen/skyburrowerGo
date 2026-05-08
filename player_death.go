package main

import (
	"image"
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	crashRotationFrames = 45
	crashTargetAngle    = 70.0 * math.Pi / 180.0
	smokeSpawnInterval  = 3
)

type smokeParticle struct {
	x, y   float64
	vx, vy float64
	life   int
	maxLife int
	radius float64
}

func (g *Game) updatePlayerDying() {
	g.playerDyingTimer++

	// Spawn smoke every 3 frames from the tail (left-center of ship)
	if g.playerDyingTimer%smokeSpawnInterval == 0 {
		tailX := g.player.x + 8
		tailY := g.player.y + float64(g.player.frameHeight)/2
		g.smokeParticles = append(g.smokeParticles, smokeParticle{
			x:       tailX,
			y:       tailY,
			vx:      -0.5 - rand.Float64()*0.5,
			vy:      -0.3 - rand.Float64()*0.4,
			life:    20,
			maxLife: 20,
			radius:  3 + rand.Float64()*2,
		})
	}

	// Update existing smoke particles
	var alive []smokeParticle
	for _, p := range g.smokeParticles {
		p.x += p.vx
		p.y += p.vy
		p.life--
		if p.life > 0 {
			alive = append(alive, p)
		}
	}
	g.smokeParticles = alive

	// Rotate nose-down gradually over crashRotationFrames, descent begins immediately
	if g.playerDyingTimer <= crashRotationFrames {
		g.playerCrashAngle = crashTargetAngle * float64(g.playerDyingTimer) / float64(crashRotationFrames)
	}

	g.playerCrashVY += 0.12
	g.player.y += g.playerCrashVY

	if g.playerDyingTimer >= PlayerDyingDuration {
		g.gameOver()
	}
}

func (g *Game) drawSmokeParticles(screen *ebiten.Image) {
	for _, p := range g.smokeParticles {
		alpha := float64(p.life) / float64(p.maxLife)
		radius := p.radius * (1 + (1-alpha)*0.8)
		grey := uint8(80 + rand.IntN(60))
		c := color.RGBA{grey, grey, grey, uint8(alpha * 180)}
		vector.DrawFilledCircle(screen, float32(p.x), float32(p.y), float32(radius), c, false)
	}
}

func (g *Game) drawPlayerDying(screen *ebiten.Image) {
	g.drawSmokeParticles(screen)

	op := &ebiten.DrawImageOptions{}
	// Rotate around center of sprite
	cx := float64(g.player.frameWidth) / 2
	cy := float64(g.player.frameHeight) / 2
	op.GeoM.Translate(-cx, -cy)
	op.GeoM.Rotate(g.playerCrashAngle)
	op.GeoM.Translate(cx, cy)
	op.GeoM.Translate(g.player.x, g.player.y)

	// Use frame 0 of the sprite sheet
	frame := (g.player.frameCounter / 10) % 2
	subImg := g.player.img.SubImage(
		image.Rect(frame*g.player.frameWidth, 0, (frame+1)*g.player.frameWidth, g.player.frameHeight),
	).(*ebiten.Image)
	screen.DrawImage(subImg, op)
}
