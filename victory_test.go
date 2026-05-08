package main

import "testing"

func TestCarrotRevealScale_ZeroBeforeReveal(t *testing.T) {
	if s := carrotRevealScale(0, 1); s != 0 {
		t.Errorf("expected 0 before carrot 1 start frame, got %.2f", s)
	}
}

func TestCarrotRevealScale_OneAfterPopComplete(t *testing.T) {
	frame := 1*CarrotRevealInterval + CarrotPopDuration + 10
	if s := carrotRevealScale(frame, 1); s != 1.0 {
		t.Errorf("expected 1.0 after pop complete, got %.2f", s)
	}
}

func TestCarrotRevealScale_OvershootDuringPop(t *testing.T) {
	// Frame just after peak overshoot (elapsed=10 for carrot 0)
	if s := carrotRevealScale(10, 0); s <= 1.0 {
		t.Errorf("expected scale > 1.0 during overshoot, got %.2f", s)
	}
}

func TestEnterVictoryScreen_SetsScreenVictory(t *testing.T) {
	g := &Game{currentScreen: ScreenPlaying, runLevelCarrotMask: 0b10101}
	g.enterVictoryScreen()
	if g.currentScreen != ScreenVictory {
		t.Errorf("expected ScreenVictory, got %v", g.currentScreen)
	}
}

func TestUpdateVictory_FadeOutTransitionsToWorldMap(t *testing.T) {
	g := &Game{}
	g.enterVictoryScreen()
	// Simulate: fade-in done, all carrots shown, click-ready
	g.fadeAlpha = 0
	g.isFading = false
	g.victoryClickReady = true
	// Trigger fade-out manually (simulating a click)
	g.isFading = true
	g.fadeIn = false
	g.fadeSpeed = FadeSpeed

	// Advance until fade completes
	maxFrames := int(1.0/FadeSpeed) + 10
	for i := 0; i < maxFrames; i++ {
		g.updateVictory()
		if g.currentScreen == ScreenWorldMap {
			return
		}
	}
	t.Errorf("expected ScreenWorldMap after fade-out, still on %v after %d frames", g.currentScreen, maxFrames)
}

func TestUpdateVictory_CommitsCarrotProgressOnExit(t *testing.T) {
	g := &Game{
		runLevelCarrotMask: 0b00111,
		currentLevel:       &LevelConfig{WorldLevel: 1},
	}
	g.enterVictoryScreen()
	g.fadeAlpha = 0
	g.isFading = true
	g.fadeIn = false
	g.fadeSpeed = FadeSpeed

	maxFrames := int(1.0/FadeSpeed) + 10
	for i := 0; i < maxFrames; i++ {
		g.updateVictory()
		if g.currentScreen == ScreenWorldMap {
			break
		}
	}

	if g.levelCarrotMask[0]&0b00111 != 0b00111 {
		t.Errorf("expected carrot progress committed, got mask %08b", g.levelCarrotMask[0])
	}
}

func TestUpdateVictory_ClickReadyAfterAllCarrotsAnimate(t *testing.T) {
	g := &Game{}
	g.enterVictoryScreen()
	// Fast-forward past fade-in
	g.fadeAlpha = 0
	g.isFading = false

	lastCarrotDone := (CarrotsPerLevel-1)*CarrotRevealInterval + CarrotPopDuration
	for i := 0; i < lastCarrotDone; i++ {
		if g.victoryClickReady {
			t.Fatalf("victoryClickReady set too early at frame %d (expected %d)", i, lastCarrotDone)
		}
		g.updateVictory()
	}
	if !g.victoryClickReady {
		t.Errorf("expected victoryClickReady=true after %d frames", lastCarrotDone)
	}
}

func TestEnterVictoryScreen_CapturesCarrotMaskAndResetState(t *testing.T) {
	g := &Game{runLevelCarrotMask: 0b10101, victoryFrame: 99, victoryClickReady: true}
	g.enterVictoryScreen()
	if g.victoryCarrotMask != 0b10101 {
		t.Errorf("expected victoryCarrotMask=0b10101, got %08b", g.victoryCarrotMask)
	}
	if g.victoryFrame != 0 {
		t.Errorf("expected victoryFrame=0, got %d", g.victoryFrame)
	}
	if g.victoryClickReady {
		t.Error("expected victoryClickReady=false after reset")
	}
	if !g.isFading || !g.fadeIn || g.fadeAlpha != 1.0 {
		t.Errorf("expected fade-in started: isFading=%v fadeIn=%v fadeAlpha=%.1f", g.isFading, g.fadeIn, g.fadeAlpha)
	}
}
