package main

import "testing"

// compile-time interface checks
var _ Enemy = (*Foreman)(nil)

func TestForeman_XFixed(t *testing.T) {
	f := NewForeman(560, 160, 200, nil)
	g := &Game{}
	x0, _ := f.GetPosition()
	for i := 0; i < 300; i++ {
		f.Update(100, 160, g)
	}
	x1, _ := f.GetPosition()
	if x1 != x0 {
		t.Errorf("X should be fixed: started %v, now %v", x0, x1)
	}
}

func TestForeman_YBobs(t *testing.T) {
	f := NewForeman(560, 160, 200, nil)
	g := &Game{}
	_, y0 := f.GetPosition()
	f.Update(100, 160, g)
	_, y1 := f.GetPosition()
	if y1 == y0 {
		t.Error("Y should change on first update (sine bob)")
	}
}

func TestForeman_YBobAmplitude(t *testing.T) {
	f := NewForeman(560, 160, 200, nil)
	g := &Game{}
	minY, maxY := 160.0, 160.0
	for i := 0; i < 360; i++ {
		f.Update(100, 160, g)
		_, y := f.GetPosition()
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}
	amplitude := (maxY - minY) / 2
	if amplitude < 6 || amplitude > 10 {
		t.Errorf("bob amplitude should be ~8px, got %.2f", amplitude)
	}
}

func TestForeman_FiresDrillBitSalvoOf3(t *testing.T) {
	f := NewForeman(560, 160, 200, nil)
	g := &Game{}
	for i := 0; i < 181; i++ {
		f.Update(100, 160, g)
	}
	drillCount := 0
	for _, b := range g.enemyBullets {
		if _, ok := b.(*DrillBit); ok {
			drillCount++
		}
	}
	if drillCount != 3 {
		t.Errorf("expected 3 DrillBit bullets after 181 frames, got %d", drillCount)
	}
}

func TestForeman_DrillBitsFiredOnlyOncePerInterval(t *testing.T) {
	f := NewForeman(560, 160, 200, nil)
	g := &Game{}
	for i := 0; i < 360; i++ {
		f.Update(100, 160, g)
	}
	drillCount := 0
	for _, b := range g.enemyBullets {
		if _, ok := b.(*DrillBit); ok {
			drillCount++
		}
	}
	if drillCount != 6 {
		t.Errorf("expected 6 DrillBit bullets after 360 frames (2 salvos), got %d", drillCount)
	}
}

func TestForeman_DropsStalactitesOf3(t *testing.T) {
	f := NewForeman(560, 160, 200, nil)
	g := &Game{}
	for i := 0; i < 241; i++ {
		f.Update(100, 160, g)
	}
	stalCount := 0
	for _, b := range g.enemyBullets {
		if _, ok := b.(*Stalactite); ok {
			stalCount++
		}
	}
	if stalCount != 3 {
		t.Errorf("expected 3 Stalactites after 241 frames, got %d", stalCount)
	}
}

func TestForeman_StalactiteStartsAtTopOfScreen(t *testing.T) {
	f := NewForeman(560, 160, 200, nil)
	g := &Game{}
	for i := 0; i < 241; i++ {
		f.Update(100, 160, g)
	}
	for _, b := range g.enemyBullets {
		if s, ok := b.(*Stalactite); ok {
			_, y := s.GetPosition()
			if y != 0 {
				t.Errorf("stalactite should start at Y=0, got %v", y)
			}
		}
	}
}

func TestForeman_StalactiteFallsDown(t *testing.T) {
	s := NewStalactite(320, 0, nil)
	_, y0 := s.GetPosition()
	s.Update()
	_, y1 := s.GetPosition()
	if y1 <= y0 {
		t.Errorf("stalactite should fall (Y increase): was %v, now %v", y0, y1)
	}
}

func TestForeman_StalactiteDamage(t *testing.T) {
	s := NewStalactite(320, 0, nil)
	if s.GetDamage() <= 0 {
		t.Errorf("stalactite damage should be positive, got %d", s.GetDamage())
	}
}

func TestForeman_HighHealth(t *testing.T) {
	f := NewForeman(560, 160, 200, nil)
	if f.health < 100 {
		t.Errorf("Foreman health should be >= 100 for a world boss, got %d", f.health)
	}
}

func TestForeman_TakeDamageAndDie(t *testing.T) {
	f := NewForeman(560, 160, 5, nil)
	if f.IsDead() {
		t.Fatal("new Foreman should not be dead")
	}
	f.TakeDamage(5)
	if !f.IsDead() {
		t.Error("Foreman should die after taking full health damage")
	}
}

func TestForeman_Phase2DrillSalvoCount(t *testing.T) {
	f := NewForeman(560, 160, 200, nil)
	g := &Game{}
	// Drop to phase 2 (≤50% HP)
	f.TakeDamage(100)
	// Run 121 frames — first phase 2 salvo fires at 120
	for i := 0; i < 121; i++ {
		f.Update(100, 160, g)
	}
	drillCount := 0
	for _, b := range g.enemyBullets {
		if _, ok := b.(*DrillBit); ok {
			drillCount++
		}
	}
	if drillCount != 5 {
		t.Errorf("expected 5 DrillBits in phase 2 salvo, got %d", drillCount)
	}
}

func TestForeman_Phase2StalactiteCount(t *testing.T) {
	f := NewForeman(560, 160, 200, nil)
	g := &Game{}
	f.TakeDamage(100) // enter phase 2
	for i := 0; i < 151; i++ {
		f.Update(100, 160, g)
	}
	stalCount := 0
	for _, b := range g.enemyBullets {
		if _, ok := b.(*Stalactite); ok {
			stalCount++
		}
	}
	if stalCount != 5 {
		t.Errorf("expected 5 Stalactites in phase 2 barrage, got %d", stalCount)
	}
}

func TestForeman_Phase2SpawnsBeetle(t *testing.T) {
	f := NewForeman(560, 160, 200, nil)
	g := &Game{}
	f.TakeDamage(100) // enter phase 2
	for i := 0; i < 301; i++ {
		f.Update(100, 160, g)
	}
	beetleCount := 0
	for _, e := range g.enemies {
		if _, ok := e.(*DynamiteBeetle); ok {
			beetleCount++
		}
	}
	if beetleCount != 1 {
		t.Errorf("expected 1 DynamiteBeetle spawned after 301 phase 2 frames, got %d", beetleCount)
	}
}

func TestForeman_Phase2TriggersAtHalfHP(t *testing.T) {
	f := NewForeman(560, 160, 200, nil)
	if f.IsPhase2() {
		t.Error("should not be in phase 2 at full health")
	}
	f.TakeDamage(99) // 101 HP — still phase 1
	if f.IsPhase2() {
		t.Error("should not be in phase 2 at 101/200 HP")
	}
	f.TakeDamage(1) // 100 HP — exactly 50%
	if !f.IsPhase2() {
		t.Error("should be in phase 2 at exactly 50% HP")
	}
	f.TakeDamage(1) // 99 HP — below 50%
	if !f.IsPhase2() {
		t.Error("should remain in phase 2 below 50% HP")
	}
}

func TestForeman_UpdateDeathAndIsDeathComplete(t *testing.T) {
	f := NewForeman(560, 160, 5, nil)
	g := &Game{}
	f.TakeDamage(5)
	// Death animation completes after enough UpdateDeath calls
	for i := 0; i < 500; i++ {
		f.UpdateDeath(g)
	}
	if !f.IsDeathComplete() {
		t.Error("IsDeathComplete should be true after 500 death frames")
	}
}
