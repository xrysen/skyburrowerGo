package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	CarrotRevealInterval = 20
	CarrotPopDuration    = 15
)

func (g *Game) enterVictoryScreen() {
	g.victoryCarrotMask = g.runLevelCarrotMask
	g.victoryFrame = 0
	g.victoryClickReady = false
	g.currentScreen = ScreenVictory
	g.fadeAlpha = 1.0
	g.fadeSpeed = FadeSpeed
	g.isFading = true
	g.fadeIn = true
}

func (g *Game) updateVictory() {
	if g.isFading {
		if g.fadeIn {
			g.fadeAlpha -= g.fadeSpeed
			if g.fadeAlpha <= 0 {
				g.fadeAlpha = 0
				g.isFading = false
			}
		} else {
			g.fadeAlpha += g.fadeSpeed
			if g.fadeAlpha >= 1.0 {
				g.fadeAlpha = 1.0
				g.isFading = false
				g.transitionToNextLevel()
			}
		}
		return
	}

	g.victoryFrame++

	lastCarrotDone := (CarrotsPerLevel-1)*CarrotRevealInterval + CarrotPopDuration
	if g.victoryFrame >= lastCarrotDone {
		g.victoryClickReady = true
	}

	if g.victoryClickReady && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.isFading = true
		g.fadeIn = false
		g.fadeSpeed = FadeSpeed
	}
}

// carrotRevealScale returns the draw scale for a carrot slot during the pop-in animation.
// Returns 0 if the carrot hasn't started revealing yet.
func carrotRevealScale(victoryFrame, carrotIndex int) float64 {
	start := carrotIndex * CarrotRevealInterval
	elapsed := victoryFrame - start
	if elapsed <= 0 {
		return 0
	}
	if elapsed < 10 {
		return float64(elapsed) / 10.0 * 1.3
	}
	if elapsed < CarrotPopDuration {
		return 1.3 - (float64(elapsed-10)/float64(CarrotPopDuration-10) * 0.3)
	}
	return 1.0
}

func (g *Game) drawVictory(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{20, 20, 20, 255}, false)

	if g.font != nil {
		headerText := "LEVEL COMPLETE"
		textWidth := float64(len(headerText)) * 10.0 * 2.0
		g.font.DrawText(screen, headerText, float64(ScreenWidth)/2-textWidth/2, 60, 2.0)
	}

	if g.victoryCarrotEmpty != nil && g.victoryCarrotFull != nil {
		carrotW := float64(g.victoryCarrotEmpty.Bounds().Dx())
		carrotH := float64(g.victoryCarrotEmpty.Bounds().Dy())
		baseScale := 2.5
		scaledW := carrotW * baseScale
		scaledH := carrotH * baseScale
		spacing := 20.0
		slotPad := 8.0
		totalWidth := float64(CarrotsPerLevel)*scaledW + float64(CarrotsPerLevel-1)*spacing
		startX := (float64(ScreenWidth) - totalWidth) / 2
		cy := float64(ScreenHeight) / 2

		// Draw all slot backgrounds first so they're always visible
		for i := 0; i < CarrotsPerLevel; i++ {
			cx := startX + float64(i)*(scaledW+spacing) + scaledW/2
			vector.DrawFilledRect(screen,
				float32(cx-scaledW/2-slotPad),
				float32(cy-scaledH/2-slotPad),
				float32(scaledW+slotPad*2),
				float32(scaledH+slotPad*2),
				color.RGBA{70, 70, 70, 220}, false)
		}

		for i := 0; i < CarrotsPerLevel; i++ {
			scale := carrotRevealScale(g.victoryFrame, i)
			if scale <= 0 {
				continue
			}
			collected := (g.victoryCarrotMask>>i)&1 == 1
			var img *ebiten.Image
			if collected {
				img = g.victoryCarrotFull
			} else {
				img = g.victoryCarrotEmpty
			}

			cx := startX + float64(i)*(scaledW+spacing) + scaledW/2

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(-carrotW/2, -carrotH/2)
			op.GeoM.Scale(baseScale*scale, baseScale*scale)
			op.GeoM.Translate(cx, cy)
			screen.DrawImage(img, op)
		}
	}

	if g.victoryClickReady && g.font != nil {
		promptText := "CLICK TO CONTINUE"
		textWidth := float64(len(promptText)) * 10.0 * 1.0
		g.font.DrawText(screen, promptText, float64(ScreenWidth)/2-textWidth/2, float64(ScreenHeight)-50, 1.0)
	}

	if g.fadeAlpha > 0 {
		vector.DrawFilledRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{0, 0, 0, uint8(g.fadeAlpha * 255)}, false)
	}
}
