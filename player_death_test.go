package main

import "testing"

func TestStartPlayerDeathSequence_SetsStatePlayerDying(t *testing.T) {
	g := &Game{state: StateLevel}
	g.startPlayerDeathSequence()
	if g.state != StatePlayerDying {
		t.Errorf("expected StatePlayerDying, got %v", g.state)
	}
}

func TestUpdatePlayerDying_SpawnsSmokeEvery3Frames(t *testing.T) {
	g := &Game{}
	g.startPlayerDeathSequence()
	g.player = &Player{x: 100, y: 100, frameWidth: 80, frameHeight: 80}

	for i := 0; i < 9; i++ {
		g.updatePlayerDying()
	}

	// 9 frames → 3 spawns (at frames 3, 6, 9)
	if len(g.smokeParticles) != 3 {
		t.Errorf("expected 3 smoke particles after 9 frames, got %d", len(g.smokeParticles))
	}
}

func TestUpdatePlayerDying_PlayerDivesDownAfterRotationPhase(t *testing.T) {
	g := &Game{}
	g.startPlayerDeathSequence()
	g.player = &Player{x: 100, y: 100, frameWidth: 80, frameHeight: 80}

	startY := g.player.y
	for i := 0; i < crashRotationFrames+10; i++ {
		g.updatePlayerDying()
	}

	if g.player.y <= startY {
		t.Errorf("expected player Y to increase (dive down) after rotation phase, startY=%.1f, currentY=%.1f", startY, g.player.y)
	}
}

func TestUpdatePlayerDying_TransitionsToGameOverAfterDuration(t *testing.T) {
	g := &Game{}
	g.startPlayerDeathSequence()
	g.player = &Player{x: 100, y: 100, frameWidth: 80, frameHeight: 80}

	for i := 0; i < PlayerDyingDuration; i++ {
		if g.state == StateGameOver {
			t.Fatalf("transitioned to StateGameOver too early at frame %d", i)
		}
		g.updatePlayerDying()
	}

	if g.state != StateGameOver {
		t.Errorf("expected StateGameOver after %d frames, got %v", PlayerDyingDuration, g.state)
	}
	if !g.isFading {
		t.Error("expected isFading=true after death sequence completes")
	}
}
