package main

import "testing"

var _ Enemy = (*Blightmoth)(nil)

func TestBlightmoth_MovesLeft(t *testing.T) {
	b := NewBlightmoth(300, 160, nil)
	g := &Game{}
	x0, _ := b.GetPosition()
	b.Update(0, 0, g)
	x1, _ := b.GetPosition()
	if x1 >= x0 {
		t.Errorf("expected x to decrease; was %v, now %v", x0, x1)
	}
}

func TestBlightmoth_FiresThreeBulletsAfterInterval(t *testing.T) {
	b := NewBlightmoth(300, 160, nil)
	g := &Game{}
	// shootTimer fires at 120 frames; run 121 to guarantee at least one fire
	for i := 0; i < 121; i++ {
		b.Update(200, 160, g)
	}
	if len(g.enemyBullets) != 3 {
		t.Errorf("expected 3 enemy bullets after shoot interval, got %d", len(g.enemyBullets))
	}
}

func TestBlightmoth_YOscillates(t *testing.T) {
	b := NewBlightmoth(300, 160, nil)
	g := &Game{}
	_, y0 := b.GetPosition()
	yChanged := false
	for i := 0; i < 60; i++ {
		b.Update(0, 0, g)
		_, y := b.GetPosition()
		if y != y0 {
			yChanged = true
			break
		}
	}
	if !yChanged {
		t.Error("expected Y position to change over 60 updates (sine wave oscillation)")
	}
}
