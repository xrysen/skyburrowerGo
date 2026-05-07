package main

import "testing"

var _ Enemy = (*WraithWhisp)(nil)

func TestWraithWhisp_DriftsLeft(t *testing.T) {
	w := NewWraithWhisp(300, 160, nil)
	g := &Game{}
	x0, _ := w.GetPosition()
	w.Update(0, 0, g)
	x1, _ := w.GetPosition()
	if x1 >= x0 {
		t.Errorf("expected x to decrease; was %v, now %v", x0, x1)
	}
}

func TestWraithWhisp_ChildTracksTowardTarget(t *testing.T) {
	child := newChildWraithWhisp(300, 160, nil, 100, 160)
	g := &Game{}
	x0, _ := child.GetPosition()
	for i := 0; i < 10; i++ {
		child.Update(0, 0, g)
	}
	x1, _ := child.GetPosition()
	if x1 >= x0 {
		t.Errorf("child should move toward target x=100 from x=300; was %v, now %v", x0, x1)
	}
}

func TestWraithWhisp_ChildOnDeathSpawnsNothing(t *testing.T) {
	child := newChildWraithWhisp(300, 160, nil, 100, 100)
	g := &Game{}
	child.OnDeath(g)
	if len(g.enemies) != 0 {
		t.Errorf("expected 0 enemies after child death, got %d", len(g.enemies))
	}
}

func TestWraithWhisp_OnDeathSpawnsTwoChildren(t *testing.T) {
	w := NewWraithWhisp(300, 160, nil)
	g := &Game{}
	w.OnDeath(g)
	if len(g.enemies) != 2 {
		t.Errorf("expected 2 enemies after parent death, got %d", len(g.enemies))
	}
	for i, e := range g.enemies {
		if e == nil {
			t.Errorf("child enemy %d is nil", i)
		}
	}
}
