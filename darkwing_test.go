package main

import "testing"

var _ Enemy = (*DarkWing)(nil)

func TestDarkWing_MovesLeft(t *testing.T) {
	dw := NewDarkWing(300, 160, nil)
	g := &Game{}
	x0, _ := dw.GetPosition()
	dw.Update(0, 0, g)
	x1, _ := dw.GetPosition()
	if x1 >= x0 {
		t.Errorf("expected x to decrease; was %v, now %v", x0, x1)
	}
}

func TestDarkWing_YTracksPlayerY(t *testing.T) {
	dw := NewDarkWing(300, 100, nil)
	g := &Game{}
	playerY := 200.0
	// After several updates Y should move toward playerY
	for i := 0; i < 20; i++ {
		dw.Update(0, playerY, g)
	}
	_, y := dw.GetPosition()
	if y <= 100 {
		t.Errorf("y should have moved toward %v from 100, got %v", playerY, y)
	}
	if y > playerY {
		t.Errorf("y should not overshoot playerY %v, got %v", playerY, y)
	}
}

func TestDarkWing_YStaysWithinScreenBounds(t *testing.T) {
	g := &Game{}

	// Player way above screen — Y should not go below 0
	dw := NewDarkWing(300, 100, nil)
	for i := 0; i < 100; i++ {
		dw.Update(0, -9999, g)
	}
	_, y := dw.GetPosition()
	if y < 0 {
		t.Errorf("y went below 0: %v", y)
	}

	// Player way below screen — Y should not exceed ScreenHeight-32
	dw2 := NewDarkWing(300, 100, nil)
	for i := 0; i < 100; i++ {
		dw2.Update(0, 9999, g)
	}
	_, y2 := dw2.GetPosition()
	if y2 > float64(ScreenHeight)-32 {
		t.Errorf("y exceeded screen bottom: %v", y2)
	}
}

func TestDarkWing_PingPongAnimFrameInBounds(t *testing.T) {
	dw := NewDarkWing(300, 160, nil)
	g := &Game{}
	// Run for 3 full ping-pong cycles (14 anim steps * 6 frame ticks = 84 ticks * 3)
	for i := 0; i < 84*3; i++ {
		dw.Update(0, 160, g)
		frame := dw.animFrame()
		if frame < 0 || frame > 7 {
			t.Errorf("animFrame out of bounds at tick %d: got %d", i, frame)
		}
	}
}

func TestDarkWing_PingPongDoesNotLoopForwardOnly(t *testing.T) {
	dw := NewDarkWing(300, 160, nil)
	g := &Game{}
	// Collect frames over two full cycles; must contain descending values after peak
	peaked := false
	prev := 0
	sawDecrement := false
	for i := 0; i < 84*2; i++ {
		dw.Update(0, 160, g)
		frame := dw.animFrame()
		if frame == 7 {
			peaked = true
		}
		if peaked && frame < prev {
			sawDecrement = true
			break
		}
		prev = frame
	}
	if !sawDecrement {
		t.Error("ping-pong animation never decremented after reaching peak frame 7")
	}
}

func TestDarkWing_DoesNotFireProjectiles(t *testing.T) {
	dw := NewDarkWing(300, 160, nil)
	g := &Game{}
	for i := 0; i < 300; i++ {
		dw.Update(0, 160, g)
	}
	if len(g.enemyBullets) != 0 {
		t.Errorf("expected 0 enemy bullets, got %d", len(g.enemyBullets))
	}
}

func TestDarkWing_LowHealth(t *testing.T) {
	dw := NewDarkWing(300, 160, nil)
	if dw.IsDead() {
		t.Fatal("new DarkWing should not be dead")
	}
	dw.TakeDamage(2)
	if dw.IsDead() {
		t.Error("DarkWing should survive 2 damage (has 3 HP)")
	}
	dw.TakeDamage(1)
	if !dw.IsDead() {
		t.Error("DarkWing should die after 3 total damage")
	}
}
