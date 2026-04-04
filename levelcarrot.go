package main

import (
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

// LevelCarrot is a scrolling pickup; slotIndex 0–4 maps to world-map carrot bits.
type LevelCarrot struct {
	x, y        float64
	img         *ebiten.Image
	frameWidth  int
	frameHeight int
	slotIndex   int
	collected   bool
}

func NewLevelCarrot(x, y float64, slot int, img *ebiten.Image) *LevelCarrot {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	return &LevelCarrot{
		x:           x,
		y:           y,
		img:         img,
		frameWidth:  w,
		frameHeight: h,
		slotIndex:   slot,
	}
}

func (c *LevelCarrot) Update() {
	c.x -= 1.5
}

func (c *LevelCarrot) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(c.x, c.y)
	screen.DrawImage(c.img, op)
}

// planCarrotSpawnFrames picks 5 monotonically increasing spawn frames spread across
// playable time so the last carrot appears with time left before the level-end fade.
func planCarrotSpawnFrames(level *LevelConfig) [CarrotsPerLevel]int {
	var times [CarrotsPerLevel]int

	playEnd := level.Duration - FadeOutDuration - 3*FPS
	if playEnd <= 0 {
		playEnd = Minutes2 - FadeOutDuration
	}

	const startMin = 2 * FPS
	if playEnd <= startMin+CarrotsPerLevel*30 {
		playEnd = startMin + CarrotsPerLevel*60
	}

	span := playEnd - startMin
	seg := span / CarrotsPerLevel
	if seg < 30 {
		seg = 30
	}

	const windowTailMargin = 10 // keep spawn pick inside window, before next wave / fade
	for i := 0; i < CarrotsPerLevel; i++ {
		windowStart := startMin + i*seg
		windowEnd := startMin + (i+1)*seg
		if windowEnd > playEnd {
			windowEnd = playEnd
		}
		hi := windowEnd - windowTailMargin
		if hi <= windowStart {
			hi = windowStart
		}
		times[i] = windowStart + rand.IntN(hi-windowStart+1)
	}

	return times
}
