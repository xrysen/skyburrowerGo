package main

import "testing"

var _ Enemy = (*HollowStag)(nil)

func TestHollowStag_TelegraphNoMovement(t *testing.T) {
	s := NewHollowStag(640, 160, nil)
	g := &Game{}
	x0, _ := s.GetPosition()
	for i := 0; i < 30; i++ {
		s.Update(0, 0, g)
	}
	x1, _ := s.GetPosition()
	if x1 != x0 {
		t.Errorf("expected x unchanged during telegraph; was %v, now %v", x0, x1)
	}
}

func TestHollowStag_ChargeMovesLeft(t *testing.T) {
	s := NewHollowStag(640, 160, nil)
	g := &Game{}
	// exhaust the telegraph (60 frames)
	for i := 0; i < 61; i++ {
		s.Update(0, 0, g)
	}
	x0, _ := s.GetPosition()
	s.Update(0, 0, g)
	x1, _ := s.GetPosition()
	if x1 >= x0 {
		t.Errorf("expected x to decrease after charge starts; was %v, now %v", x0, x1)
	}
}

func TestHollowStag_NoEnemyBullets(t *testing.T) {
	s := NewHollowStag(640, 160, nil)
	g := &Game{}
	for i := 0; i < 200; i++ {
		s.Update(0, 0, g)
	}
	if len(g.enemyBullets) != 0 {
		t.Errorf("expected no enemy bullets, got %d", len(g.enemyBullets))
	}
}
