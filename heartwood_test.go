package main

import "testing"

var _ Enemy = (*Heartwood)(nil)

func TestHeartwood_Phase1AboveThreshold(t *testing.T) {
	h := NewHeartwood(580, 0, 300, nil)
	if h.Phase() != 1 {
		t.Errorf("expected phase 1 at full health, got %d", h.Phase())
	}
}

func TestHeartwood_Phase2At50Percent(t *testing.T) {
	h := NewHeartwood(580, 0, 300, nil)
	h.TakeDamage(150) // 50% HP — within ≤66% window
	if h.Phase() != 2 {
		t.Errorf("expected phase 2 at 50%% HP, got %d", h.Phase())
	}
}

func TestHeartwood_PositionFixed(t *testing.T) {
	h := NewHeartwood(580, 0, 300, nil)
	g := &Game{}
	x0, y0 := h.GetPosition()
	for i := 0; i < 400; i++ {
		h.Update(100, 160, g)
	}
	x1, y1 := h.GetPosition()
	if x1 != x0 || y1 != y0 {
		t.Errorf("position should be fixed; was (%.0f,%.0f), now (%.0f,%.0f)", x0, y0, x1, y1)
	}
}

func TestHeartwood_Phase3At20Percent(t *testing.T) {
	h := NewHeartwood(580, 0, 300, nil)
	h.TakeDamage(240) // 20% HP — within ≤33% window
	if h.Phase() != 3 {
		t.Errorf("expected phase 3 at 20%% HP, got %d", h.Phase())
	}
}

func TestHeartwood_IsDeadFalseMidAnimation(t *testing.T) {
	h := NewHeartwood(580, 0, 5, nil)
	g := &Game{}
	h.TakeDamage(5)
	// Only 10 UpdateDeath calls — well before the 300-frame threshold
	for i := 0; i < 10; i++ {
		h.UpdateDeath(g)
	}
	if h.IsDead() {
		t.Error("IsDead should be false mid death animation")
	}
}

func TestHeartwood_IsDeathCompleteAfterAnimation(t *testing.T) {
	h := NewHeartwood(580, 0, 5, nil)
	g := &Game{}
	h.TakeDamage(5)
	for i := 0; i < 301; i++ {
		h.UpdateDeath(g)
	}
	if !h.IsDeathComplete() {
		t.Error("IsDeathComplete should be true after full animation")
	}
}
