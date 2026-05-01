package main

import "testing"

// Compile-time check: ThunderCrab must implement Enemy.
var _ Enemy = (*ThunderCrab)(nil)

// --- Level 10 config ---

func TestLevel10_HasThunderCrabBoss(t *testing.T) {
	cfg := GetLevel10()
	if cfg.BossType != ThunderCrabType {
		t.Errorf("Level 10 BossType: want %q, got %q", ThunderCrabType, cfg.BossType)
	}
	if cfg.EndCondition != EndOnBossDeath {
		t.Error("Level 10 should use EndOnBossDeath")
	}
	if cfg.BossHealth != 180 {
		t.Errorf("Level 10 BossHealth: want 180, got %d", cfg.BossHealth)
	}
}

// --- Health and phases ---

func TestThunderCrab_StartsAliveWithGivenHealth(t *testing.T) {
	tc := NewThunderCrab(400, 200, 180, 8.0, nil)
	if tc.IsDead() {
		t.Fatal("new ThunderCrab should not be dead")
	}
	tc.TakeDamage(179)
	if tc.IsDead() {
		t.Error("should survive 179 damage (has 180 HP)")
	}
	tc.TakeDamage(1)
	if !tc.IsDead() {
		t.Error("should die after 180 total damage")
	}
}

// --- Attacks ---

func TestThunderCrab_Phase1FiresLightningBoltWithTell(t *testing.T) {
	tc := NewThunderCrab(400, 200, 180, 8.0, nil)
	g := &Game{}

	// Advance to just before the lightning tell fires (interval-1 frames)
	for i := 0; i < lightningInterval-1; i++ {
		tc.Update(0, 0, g)
	}
	if len(g.enemyBullets) != 0 {
		t.Fatalf("no bullets expected before tell fires, got %d", len(g.enemyBullets))
	}

	// Frame that crosses the threshold → tell begins (30 frames), no bullet yet
	tc.Update(0, 0, g)
	if tc.attackTellTimer != 30 {
		t.Errorf("expected 30-frame tell, got %d", tc.attackTellTimer)
	}
	if tc.pendingAttack != tellLightning {
		t.Errorf("expected pending attack = tellLightning (%d), got %d", tellLightning, tc.pendingAttack)
	}
	if len(g.enemyBullets) != 0 {
		t.Fatalf("bullet should not fire during tell, got %d", len(g.enemyBullets))
	}

	// Advance through the tell (29 more frames)
	for i := 0; i < 29; i++ {
		tc.Update(0, 0, g)
	}
	if len(g.enemyBullets) != 0 {
		t.Fatalf("bullet should not fire before tell expires, got %d", len(g.enemyBullets))
	}

	// Final tell frame → bolt fires
	tc.Update(0, 0, g)
	if len(g.enemyBullets) != 1 {
		t.Fatalf("expected 1 bullet after tell, got %d", len(g.enemyBullets))
	}
	if _, ok := g.enemyBullets[0].(*LightningBolt); !ok {
		t.Error("expected a *LightningBolt")
	}
}

func TestThunderCrab_Phase1FiresCloudProjectile(t *testing.T) {
	tc := NewThunderCrab(400, 200, 180, 8.0, nil)
	g := &Game{}

	// Exhaust the lightning timer first so it doesn't interfere
	for i := 0; i < lightningInterval+30; i++ {
		tc.Update(0, 0, g)
	}
	g.enemyBullets = nil // reset

	// Advance until a CloudProjectile appears
	found := false
	for i := 0; i < cloudInterval+30+1; i++ {
		tc.Update(0, 0, g)
		for _, b := range g.enemyBullets {
			if _, ok := b.(*CloudProjectile); ok {
				found = true
			}
		}
		if found {
			break
		}
	}
	if !found {
		t.Error("expected a CloudProjectile from phase 1")
	}
}

func TestThunderCrab_Phase2AddsElectricalRing(t *testing.T) {
	tc := NewThunderCrab(400, 200, 180, 8.0, nil)
	g := &Game{}

	// Transition to phase 2
	tc.TakeDamage(62) // ~65.5% HP → phase 2
	tc.Update(0, 0, g)
	if tc.phase != 2 {
		t.Fatalf("setup failed: expected phase 2, got %d", tc.phase)
	}
	g.enemyBullets = nil

	// Use a large window: tells from lightning/cloud pause the ring timer,
	// so we need headroom beyond ringInterval alone.
	found := false
	for i := 0; i < 600; i++ {
		tc.Update(0, 0, g)
		for _, b := range g.enemyBullets {
			if _, ok := b.(*ElectricalRing); ok {
				found = true
			}
		}
		if found {
			break
		}
	}
	if !found {
		t.Error("expected an ElectricalRing in phase 2")
	}
}

func TestThunderCrab_Phase2DoesNotFireRingInPhase1(t *testing.T) {
	tc := NewThunderCrab(400, 200, 180, 8.0, nil)
	g := &Game{}

	// Stay in phase 1 — run 600 frames
	for i := 0; i < 600; i++ {
		tc.Update(0, 0, g)
	}
	for _, b := range g.enemyBullets {
		if _, ok := b.(*ElectricalRing); ok {
			t.Error("ElectricalRing should not fire in phase 1")
		}
	}
}

func TestThunderCrab_Phase3AddsChainLightningBurst(t *testing.T) {
	tc := NewThunderCrab(400, 200, 180, 8.0, nil)
	g := &Game{}

	tc.TakeDamage(125) // 55 HP ≈ 30.5% → phase 3
	tc.Update(0, 0, g)
	if tc.phase != 3 {
		t.Fatalf("setup failed: expected phase 3, got %d", tc.phase)
	}
	g.enemyBullets = nil

	found := false
	for i := 0; i < 800; i++ {
		tc.Update(0, 0, g)
		for _, b := range g.enemyBullets {
			if _, ok := b.(*ChainLightningBolt); ok {
				found = true
			}
		}
		if found {
			break
		}
	}
	if !found {
		t.Error("expected a ChainLightningBolt in phase 3")
	}
}

func TestThunderCrab_Phase3AddsShockwave(t *testing.T) {
	tc := NewThunderCrab(400, 200, 180, 8.0, nil)
	g := &Game{}

	tc.TakeDamage(125) // phase 3
	tc.Update(0, 0, g)
	g.enemyBullets = nil

	found := false
	for i := 0; i < 1200; i++ {
		tc.Update(0, 0, g)
		for _, b := range g.enemyBullets {
			if _, ok := b.(*Shockwave); ok {
				found = true
			}
		}
		if found {
			break
		}
	}
	if !found {
		t.Error("expected a Shockwave in phase 3")
	}
}

func TestThunderCrab_Phase3LightningFiresFasterThanPhase1(t *testing.T) {
	// Phase 1: count frames until first LightningBolt fires
	tc1 := NewThunderCrab(400, 200, 180, 8.0, nil)
	g1 := &Game{}
	framesP1 := 0
	for i := 0; i < 500; i++ {
		tc1.Update(0, 0, g1)
		framesP1++
		if len(g1.enemyBullets) > 0 {
			break
		}
	}

	// Phase 3: count frames until first LightningBolt fires
	tc3 := NewThunderCrab(400, 200, 180, 8.0, nil)
	tc3.TakeDamage(125) // phase 3
	tc3.Update(0, 0, &Game{})
	g3 := &Game{}
	framesP3 := 0
	for i := 0; i < 500; i++ {
		tc3.Update(0, 0, g3)
		framesP3++
		if len(g3.enemyBullets) > 0 {
			break
		}
	}

	if framesP3 >= framesP1 {
		t.Errorf("phase 3 first attack should fire faster than phase 1: phase1=%d frames, phase3=%d frames", framesP1, framesP3)
	}
}

func TestThunderCrab_AttackTellIs30Frames(t *testing.T) {
	tc := NewThunderCrab(400, 200, 180, 8.0, nil)
	g := &Game{}

	// Drive until a tell starts
	for i := 0; i < 300; i++ {
		tc.Update(0, 0, g)
		if tc.attackTellTimer == 30 {
			// Verify it counts down without firing for 29 frames
			beforeBullets := len(g.enemyBullets)
			for j := 0; j < 29; j++ {
				tc.Update(0, 0, g)
				if len(g.enemyBullets) > beforeBullets {
					t.Errorf("bullet fired before 30-frame tell expired (frame %d of tell)", j+1)
				}
			}
			return
		}
	}
	t.Error("no attack tell observed in 300 frames")
}

// --- Death ---

func TestThunderCrab_DeathAnimationSetsBossKilled(t *testing.T) {
	tc := NewThunderCrab(400, 200, 1, 8.0, nil)
	g := &Game{player: &Player{}}
	tc.TakeDamage(1)
	tc.OnDeath(g)

	if !tc.isDying {
		t.Fatal("expected isDying after OnDeath")
	}

	// Simulate game loop: Update increments deathTimer, then UpdateDeath checks it.
	for i := 0; i < 480; i++ {
		tc.Update(0, 0, g)  // isDying → deathTimer++
		tc.UpdateDeath(g)
	}

	if !g.bossKilled {
		t.Error("bossKilled should be true after death animation completes (480 frames)")
	}
	if !tc.IsDeathComplete() {
		t.Error("IsDeathComplete should return true after 480 frames")
	}
}

func TestThunderCrab_DeathCoinsSpawnedAt120Frames(t *testing.T) {
	tc := NewThunderCrab(400, 200, 1, 8.0, nil)
	g := &Game{player: &Player{}}
	tc.TakeDamage(1)
	tc.OnDeath(g)

	// Simulate game loop for 120 frames
	for i := 0; i < 120; i++ {
		tc.Update(0, 0, g)
		tc.UpdateDeath(g)
	}

	if len(g.coins) == 0 {
		t.Error("expected coins to spawn at death frame 120")
	}
}

func TestThunderCrab_PhaseTransitions(t *testing.T) {
	tc := NewThunderCrab(400, 200, 180, 8.0, nil)
	g := &Game{}

	if tc.phase != 1 {
		t.Errorf("new crab should be phase 1, got %d", tc.phase)
	}

	// Drop to exactly 66% HP → should still be phase 1 (>66%)
	tc.TakeDamage(59) // 59 damage → 121 HP ≈ 67%
	tc.Update(0, 0, g)
	if tc.phase != 1 {
		t.Errorf("at ~67%% HP expected phase 1, got %d", tc.phase)
	}

	// Drop to ≤66% HP → phase 2
	tc.TakeDamage(3) // 118 HP ≈ 65.5%
	tc.Update(0, 0, g)
	if tc.phase != 2 {
		t.Errorf("at ≤66%% HP expected phase 2, got %d", tc.phase)
	}

	// Drop to ≤33% HP → phase 3
	tc.TakeDamage(62) // 56 HP ≈ 31%
	tc.Update(0, 0, g)
	if tc.phase != 3 {
		t.Errorf("at ≤33%% HP expected phase 3, got %d", tc.phase)
	}
}
