package main

import (
	"math"
	"testing"
)

// Compile-time check: LightningBug must implement Enemy.
var _ Enemy = (*LightningBug)(nil)

func TestLightningBug_MovesLeftAt3_5(t *testing.T) {
	lb := NewLightningBug(300, 160, nil, 6.5)
	g := &Game{}
	x0, _ := lb.GetPosition()
	lb.Update(0, 0, g)
	x1, _ := lb.GetPosition()
	if x1 != x0-3.5 {
		t.Errorf("expected x to decrease by 3.5: was %v, now %v", x0, x1)
	}
}

func TestLightningBug_StartsWithThreeHP(t *testing.T) {
	lb := NewLightningBug(100, 100, nil, 6.5)
	if lb.IsDead() {
		t.Fatal("new LightningBug should not be dead")
	}
	lb.TakeDamage(2)
	if lb.IsDead() {
		t.Error("should survive 2 damage (has 3 HP)")
	}
	lb.TakeDamage(1)
	if !lb.IsDead() {
		t.Error("should die after 3 total damage")
	}
}

func TestLightningBug_FiresFourBoltsPerBurst(t *testing.T) {
	lb := NewLightningBug(300, 160, nil, 6.5)
	g := &Game{}
	// Advance to just before burst triggers
	for i := 0; i < 119; i++ {
		lb.Update(0, 0, g)
	}
	if len(g.enemyBullets) != 0 {
		t.Errorf("expected 0 bullets before frame 120, got %d", len(g.enemyBullets))
	}
	// Advance through the full burst window (120 + 3*8 = 144 frames total)
	for i := 0; i < 25; i++ {
		lb.Update(0, 0, g)
	}
	if len(g.enemyBullets) != 4 {
		t.Errorf("expected 4 ChainLightningBolts after full burst, got %d", len(g.enemyBullets))
	}
	for i, b := range g.enemyBullets {
		if _, ok := b.(*ChainLightningBolt); !ok {
			t.Errorf("bullet %d is not a *ChainLightningBolt", i)
		}
	}
}

func TestLightningBug_BurstFiresOneEvery8Frames(t *testing.T) {
	lb := NewLightningBug(300, 160, nil, 6.5)
	g := &Game{}
	// Advance to burst trigger
	for i := 0; i < 120; i++ {
		lb.Update(0, 0, g)
	}
	after120 := len(g.enemyBullets)
	if after120 < 1 {
		t.Fatalf("expected at least 1 bolt at frame 120, got %d", after120)
	}
	// 8 more frames: one more bolt
	for i := 0; i < 8; i++ {
		lb.Update(0, 0, g)
	}
	after128 := len(g.enemyBullets)
	if after128 != after120+1 {
		t.Errorf("expected exactly 1 more bolt after 8 frames; had %d, now %d", after120, after128)
	}
}

func TestLightningBug_BoltsAimedAtPlayer(t *testing.T) {
	// Bug starts at 800; after 120 frames it's at x=380 (center 412).
	// Player is directly to the left at same height so bolt must travel left.
	const bx, by = 800.0, 200.0
	const px, py = 0.0, 232.0 // same y as bug center (by+32)
	lb := NewLightningBug(bx, by, nil, 6.5)
	g := &Game{}
	for i := 0; i < 120; i++ {
		lb.Update(px, py, g)
	}
	if len(g.enemyBullets) == 0 {
		t.Fatal("expected at least one bolt after 120 frames")
	}
	bolt := g.enemyBullets[0]
	x0, y0 := bolt.GetPosition()
	bolt.Update()
	x1, y1 := bolt.GetPosition()
	// Bolt should move leftward (toward player) not rightward
	if x1 >= x0 {
		t.Errorf("bolt should move toward player (left); dx=%v", x1-x0)
	}
	// Vertical drift should be near zero since player is at same Y
	dy := math.Abs(y1 - y0)
	dx := math.Abs(x1 - x0)
	if dy > dx*0.5 {
		t.Errorf("bolt aimed poorly: dx=%v dy=%v (player is directly left)", dx, dy)
	}
}

func TestLightningBug_BulletSpeedForwarded(t *testing.T) {
	const wantSpeed = 6.5
	lb := NewLightningBug(300, 160, nil, wantSpeed)
	g := &Game{}
	for i := 0; i < 120; i++ {
		lb.Update(0, 0, g)
	}
	if len(g.enemyBullets) == 0 {
		t.Fatal("expected at least one bolt after burst trigger")
	}
	bolt := g.enemyBullets[0]
	x0, y0 := bolt.GetPosition()
	bolt.Update()
	x1, y1 := bolt.GetPosition()
	dist := math.Sqrt((x1-x0)*(x1-x0) + (y1-y0)*(y1-y0))
	if math.Abs(dist-wantSpeed) > 0.001 {
		t.Errorf("bolt speed: want %v, got %v", wantSpeed, dist)
	}
}

func TestLevel7_HasLightningBugSpawns(t *testing.T) {
	cfg := GetLevel7()
	found := false
	for _, sc := range cfg.SpawnConfigs {
		if sc.EnemyType == LightningBugType {
			found = true
			break
		}
	}
	if !found {
		t.Error("Level 7 should have at least one LightningBug spawn config")
	}
}
