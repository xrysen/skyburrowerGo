package main

import "testing"

// Compile-time interface checks — if any type doesn't implement Bullet the build fails.
var (
	_ Bullet = (*LightningBolt)(nil)
	_ Bullet = (*CloudProjectile)(nil)
	_ Bullet = (*ElectricalRing)(nil)
	_ Bullet = (*ChainLightningBolt)(nil)
	_ Bullet = (*Shockwave)(nil)
)

// --- LightningBolt ---

func TestLightningBolt_MovesOnUpdate(t *testing.T) {
	b := NewLightningBolt(100, 100, -5.0, 0.0)
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

func TestLightningBolt_Damage(t *testing.T) {
	b := NewLightningBolt(100, 100, -5.0, 0.0)
	if b.GetDamage() != 1 {
		t.Errorf("want damage 1, got %d", b.GetDamage())
	}
}

// --- CloudProjectile ---

func TestCloudProjectile_MovesOnUpdate(t *testing.T) {
	b := NewCloudProjectile(200, 150, -4.0, 0.0)
	x0, _ := b.GetPosition()
	b.Update()
	x1, _ := b.GetPosition()
	if x1 >= x0 {
		t.Errorf("expected x to decrease: was %v, now %v", x0, x1)
	}
}

func TestCloudProjectile_Damage(t *testing.T) {
	b := NewCloudProjectile(200, 150, -4.0, 0.0)
	if b.GetDamage() != 1 {
		t.Errorf("want damage 1, got %d", b.GetDamage())
	}
}

// --- ChainLightningBolt ---

func TestChainLightningBolt_MovesOnUpdate(t *testing.T) {
	b := NewChainLightningBolt(300, 200, -6.0, 1.0)
	x0, y0 := b.GetPosition()
	b.Update()
	x1, y1 := b.GetPosition()
	if x1 == x0 && y1 == y0 {
		t.Error("expected position to change after Update")
	}
}

func TestChainLightningBolt_Damage(t *testing.T) {
	b := NewChainLightningBolt(300, 200, -6.0, 1.0)
	if b.GetDamage() != 1 {
		t.Errorf("want damage 1, got %d", b.GetDamage())
	}
}

// --- ElectricalRing ---

func TestElectricalRing_RadiusGrowsOnUpdate(t *testing.T) {
	b := NewElectricalRing(200, 200, 0)
	r0 := b.GetRadius()
	b.Update()
	r1 := b.GetRadius()
	if r1 <= r0 {
		t.Errorf("expected radius to grow: was %v, now %v", r0, r1)
	}
}

func TestElectricalRing_ActiveWindowOnly(t *testing.T) {
	b := NewElectricalRing(200, 200, 0)
	// Advance until just inside the active window (radius 8–20)
	for b.GetRadius() < 8 {
		b.Update()
		if !b.IsActive() && b.GetRadius() >= 8 {
			t.Errorf("should be active at radius %v", b.GetRadius())
		}
	}
	if !b.IsActive() {
		t.Errorf("should be active at radius %v (in window 8–20)", b.GetRadius())
	}
	// Advance past max radius (20)
	for b.GetRadius() < 20 {
		b.Update()
	}
	b.Update() // push past 20
	if b.IsActive() {
		t.Errorf("should be inactive past radius 20, got radius %v", b.GetRadius())
	}
}

func TestElectricalRing_ContinuesMovingLeftAfterMaxRadius(t *testing.T) {
	b := NewElectricalRing(200, 200, -1.5)
	// Drive past max radius
	for i := 0; i < 1000; i++ {
		b.Update()
		if b.GetRadius() >= 20 {
			break
		}
	}
	x0, _ := b.GetPosition()
	b.Update()
	x1, _ := b.GetPosition()
	if x1 >= x0 {
		t.Errorf("ring should keep moving left after max radius: was %v, now %v", x0, x1)
	}
}

func TestElectricalRing_Damage(t *testing.T) {
	b := NewElectricalRing(200, 200, 0)
	if b.GetDamage() != 1 {
		t.Errorf("want damage 1, got %d", b.GetDamage())
	}
}

// --- Shockwave ---

func TestShockwave_RadiusGrowsOnUpdate(t *testing.T) {
	b := NewShockwave(200, 200)
	r0 := b.GetRadius()
	b.Update()
	r1 := b.GetRadius()
	if r1 <= r0 {
		t.Errorf("expected radius to grow: was %v, now %v", r0, r1)
	}
}

func TestShockwave_ActiveWindowOnly(t *testing.T) {
	b := NewShockwave(200, 200)
	for b.GetRadius() < 10 {
		b.Update()
	}
	if !b.IsActive() {
		t.Errorf("should be active at radius %v (window 10–40)", b.GetRadius())
	}
	for b.GetRadius() < 40 {
		b.Update()
	}
	b.Update()
	if b.IsActive() {
		t.Errorf("should be inactive past radius 40, got %v", b.GetRadius())
	}
}

func TestShockwave_ReportsDeadAtMaxRadius(t *testing.T) {
	b := NewShockwave(200, 200)
	for i := 0; i < 1000; i++ {
		b.Update()
		if b.GetRadius() >= 40 {
			break
		}
	}
	x, y := b.GetPosition()
	if x > -50 && x < ScreenWidth+50 && y > -50 && y < ScreenHeight+50 {
		t.Errorf("dead shockwave should be off-screen; got (%v, %v)", x, y)
	}
}

func TestShockwave_Damage(t *testing.T) {
	b := NewShockwave(200, 200)
	if b.GetDamage() != 2 {
		t.Errorf("want damage 2, got %d", b.GetDamage())
	}
}
