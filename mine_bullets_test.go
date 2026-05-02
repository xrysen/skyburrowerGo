package main

import "testing"

var (
	_ Bullet = (*DrillBit)(nil)
	_ Bullet = (*FuseSpark)(nil)
)

// --- DrillBit ---

func TestDrillBit_MovesOnUpdate(t *testing.T) {
	b := NewDrillBit(100, 100, -2.0, 0.0)
	x0, y0 := b.GetPosition()
	b.Update()
	x1, y1 := b.GetPosition()
	if x1 >= x0 {
		t.Errorf("expected x to decrease (move left): was %v, now %v", x0, x1)
	}
	if y1 != y0 {
		t.Errorf("expected y unchanged: was %v, now %v", y0, y1)
	}
}

func TestDrillBit_HighDamage(t *testing.T) {
	b := NewDrillBit(100, 100, -2.0, 0.0)
	if b.GetDamage() != 3 {
		t.Errorf("want damage 3 (high), got %d", b.GetDamage())
	}
}

// --- FuseSpark ---

func TestFuseSpark_MovesOnUpdate(t *testing.T) {
	b := NewFuseSpark(200, 150, -3.5, 0.0)
	x0, _ := b.GetPosition()
	b.Update()
	x1, _ := b.GetPosition()
	if x1 >= x0 {
		t.Errorf("expected x to decrease (move left): was %v, now %v", x0, x1)
	}
}

func TestFuseSpark_MediumDamage(t *testing.T) {
	b := NewFuseSpark(200, 150, -3.5, 0.0)
	if b.GetDamage() != 2 {
		t.Errorf("want damage 2 (medium), got %d", b.GetDamage())
	}
}
