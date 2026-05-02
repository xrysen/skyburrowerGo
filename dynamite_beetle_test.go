package main

import "testing"

var _ Enemy = (*DynamiteBeetle)(nil)

func TestDynamiteBeetle_MovesLeft(t *testing.T) {
	db := NewDynamiteBeetle(300, 160, nil)
	g := &Game{}
	x0, _ := db.GetPosition()
	db.Update(0, 0, g)
	x1, _ := db.GetPosition()
	if x1 >= x0 {
		t.Errorf("expected x to decrease; was %v, now %v", x0, x1)
	}
}

func TestDynamiteBeetle_FiresFuseSparkAtInterval(t *testing.T) {
	db := NewDynamiteBeetle(300, 160, nil)
	g := &Game{}

	// Should not fire before 150 frames
	for i := 0; i < 149; i++ {
		db.Update(0, 0, g)
	}
	if len(g.enemyBullets) != 0 {
		t.Errorf("expected no bullets before interval, got %d", len(g.enemyBullets))
	}

	// Frame 150 should fire one FuseSpark
	db.Update(0, 0, g)
	if len(g.enemyBullets) != 1 {
		t.Errorf("expected 1 bullet at frame 150, got %d", len(g.enemyBullets))
	}

	// Frame 300 should fire a second
	for i := 0; i < 150; i++ {
		db.Update(0, 0, g)
	}
	if len(g.enemyBullets) != 2 {
		t.Errorf("expected 2 bullets at frame 300, got %d", len(g.enemyBullets))
	}
}

func TestDynamiteBeetle_OnDeathFires4Cardinals(t *testing.T) {
	db := NewDynamiteBeetle(300, 160, nil)
	g := &Game{}
	db.OnDeath(g)

	if len(g.enemyBullets) != 4 {
		t.Fatalf("expected exactly 4 death bullets, got %d", len(g.enemyBullets))
	}

	// Verify all 4 cardinal directions are represented (one per axis pair)
	type vec struct{ vx, vy float64 }
	dirs := map[vec]bool{}
	for _, b := range g.enemyBullets {
		fs, ok := b.(*FuseSpark)
		if !ok {
			t.Fatal("death bullet is not a FuseSpark")
		}
		dirs[vec{fs.vx, fs.vy}] = true
	}
	cardinals := []vec{{-3, 0}, {3, 0}, {0, -3}, {0, 3}}
	for _, c := range cardinals {
		if !dirs[c] {
			t.Errorf("missing cardinal direction %v", c)
		}
	}
}

func TestDynamiteBeetle_MediumHealth(t *testing.T) {
	db := NewDynamiteBeetle(300, 160, nil)
	if db.IsDead() {
		t.Fatal("new DynamiteBeetle should not be dead")
	}
	db.TakeDamage(3)
	if db.IsDead() {
		t.Error("DynamiteBeetle should survive 3 damage (has 6 HP)")
	}
	db.TakeDamage(3)
	if !db.IsDead() {
		t.Error("DynamiteBeetle should die after 6 total damage")
	}
}
