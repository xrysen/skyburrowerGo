package main

import "testing"

var _ Enemy = (*DrillDrone)(nil)

func TestDrillDrone_MovesLeft(t *testing.T) {
	d := NewDrillDrone(300, 160, nil)
	g := &Game{}
	x0, _ := d.GetPosition()
	d.Update(0, 160, g)
	x1, _ := d.GetPosition()
	if x1 >= x0 {
		t.Errorf("expected x to decrease; was %v, now %v", x0, x1)
	}
}

func TestDrillDrone_FiresDrillBitAtInterval(t *testing.T) {
	d := NewDrillDrone(300, 160, nil)
	g := &Game{}

	for i := 0; i < 99; i++ {
		d.Update(0, 160, g)
	}
	if len(g.enemyBullets) != 0 {
		t.Errorf("expected no bullets before interval, got %d", len(g.enemyBullets))
	}

	d.Update(0, 160, g)
	if len(g.enemyBullets) != 1 {
		t.Errorf("expected 1 DrillBit at frame 100, got %d", len(g.enemyBullets))
	}
	if _, ok := g.enemyBullets[0].(*DrillBit); !ok {
		t.Error("bullet should be a DrillBit")
	}

	for i := 0; i < 100; i++ {
		d.Update(0, 160, g)
	}
	if len(g.enemyBullets) != 2 {
		t.Errorf("expected 2 DrillBits after 200 frames, got %d", len(g.enemyBullets))
	}
}

func TestDrillDrone_VerticalDashTriggersAndCompletes(t *testing.T) {
	d := NewDrillDrone(300, 160, nil)
	g := &Game{}

	// Run to just before dash window
	for i := 0; i < 119; i++ {
		d.Update(0, 160, g)
	}
	_, yBefore := d.GetPosition()

	// Frame 120 should begin a dash — Y moves away from baseline
	var yDuringDash float64
	for i := 0; i < 20; i++ {
		d.Update(0, 160, g)
		_, yDuringDash = d.GetPosition()
	}
	if yDuringDash == yBefore {
		t.Error("expected Y to change during dash window")
	}

	// After dash completes (~40px), further updates should not keep moving Y away indefinitely
	for i := 0; i < 200; i++ {
		d.Update(0, 160, g)
	}
	_, yAfter := d.GetPosition()
	displacement := yAfter - yBefore
	if displacement > 80 || displacement < -80 {
		t.Errorf("expected Y to settle after dash, total displacement %v is too large", displacement)
	}
}

func TestDrillDrone_HighHealth(t *testing.T) {
	d := NewDrillDrone(300, 160, nil)
	if d.IsDead() {
		t.Fatal("new DrillDrone should not be dead")
	}
	// Must have more HP than StormSprite (15), the highest prior regular enemy
	d.TakeDamage(15)
	if d.IsDead() {
		t.Error("DrillDrone should survive 15 damage (must have more than 15 HP)")
	}
}
