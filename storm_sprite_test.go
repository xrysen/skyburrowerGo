package main

import "testing"

// Compile-time check: StormSprite must implement Enemy.
var _ Enemy = (*StormSprite)(nil)

func TestStormSprite_StartsWithFifteenHP(t *testing.T) {
	ss := NewStormSprite(100, 100, nil, 7.0)
	if ss.IsDead() {
		t.Fatal("new StormSprite should not be dead")
	}
	ss.TakeDamage(14)
	if ss.IsDead() {
		t.Error("should survive 14 damage (has 15 HP)")
	}
	ss.TakeDamage(1)
	if !ss.IsDead() {
		t.Error("should die after 15 total damage")
	}
}

func TestStormSprite_FiresRingEvery150Frames(t *testing.T) {
	ss := NewStormSprite(300, 160, nil, 7.0)
	g := &Game{}
	for i := 0; i < 149; i++ {
		ss.Update(0, 0, g)
	}
	if len(g.enemyBullets) != 0 {
		t.Errorf("expected 0 bullets before frame 150, got %d", len(g.enemyBullets))
	}
	ss.Update(0, 0, g)
	if len(g.enemyBullets) != 1 {
		t.Errorf("expected 1 ring at frame 150, got %d", len(g.enemyBullets))
	}
	if _, ok := g.enemyBullets[0].(*ElectricalRing); !ok {
		t.Error("bullet at frame 150 should be an *ElectricalRing")
	}
}

func TestStormSprite_RingSpawnsAtSpriteCenter(t *testing.T) {
	ss := NewStormSprite(300, 160, nil, 7.0)
	g := &Game{}
	for i := 0; i < 150; i++ {
		ss.Update(0, 0, g)
	}
	if len(g.enemyBullets) == 0 {
		t.Fatal("expected a ring after 150 frames")
	}
	ring := g.enemyBullets[0].(*ElectricalRing)
	// Ring center is (cx, cy) — accessible via GetPosition before it moves off-screen.
	// At birth radius is 0, position == center.
	rx, ry := ring.GetPosition()
	sx, sy := ss.GetPosition()
	wantX := sx + 32
	wantY := sy + 32
	if rx != wantX || ry != wantY {
		t.Errorf("ring center: want (%.1f, %.1f), got (%.1f, %.1f)", wantX, wantY, rx, ry)
	}
}

func TestLevel8_HasStormSpriteSpawns(t *testing.T) {
	cfg := GetLevel8()
	found := false
	for _, sc := range cfg.SpawnConfigs {
		if sc.EnemyType == StormSpriteType {
			found = true
			break
		}
	}
	if !found {
		t.Error("Level 8 should have at least one StormSprite spawn config")
	}
}

func TestStormSprite_MovesLeftAt1_0(t *testing.T) {
	ss := NewStormSprite(300, 160, nil, 7.0)
	g := &Game{}
	x0, _ := ss.GetPosition()
	ss.Update(0, 0, g)
	x1, _ := ss.GetPosition()
	if x1 != x0-1.0 {
		t.Errorf("expected x to decrease by 1.0: was %v, now %v", x0, x1)
	}
}
