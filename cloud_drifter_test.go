package main

import (
	"math"
	"testing"
)

// Compile-time check: CloudDrifter must implement Enemy.
var _ Enemy = (*CloudDrifter)(nil)

func TestCloudDrifter_MovesLeftAt1_5(t *testing.T) {
	cd := NewCloudDrifter(300, 160, nil, 6.0)
	g := &Game{}
	x0, _ := cd.GetPosition()
	cd.Update(0, 0, g)
	x1, _ := cd.GetPosition()
	if x1 != x0-1.5 {
		t.Errorf("expected x to decrease by 1.5: was %v, now %v", x0, x1)
	}
}

func TestCloudDrifter_YDriftsSinusoidally(t *testing.T) {
	cd := NewCloudDrifter(300, 160, nil, 6.0)
	g := &Game{}
	var ys []float64
	for i := 0; i < 30; i++ {
		cd.Update(0, 0, g)
		_, y := cd.GetPosition()
		ys = append(ys, y)
	}
	// y must not be constant — sinusoidal drift means values vary
	allSame := true
	for _, y := range ys[1:] {
		if y != ys[0] {
			allSame = false
			break
		}
	}
	if allSame {
		t.Error("y should drift sinusoidally, but all values were identical")
	}
	// y must also return near its base after a full period (180 frames)
	cd2 := NewCloudDrifter(300, 160, nil, 6.0)
	for i := 0; i < 180; i++ {
		cd2.Update(0, 0, g)
	}
	_, yAfterPeriod := cd2.GetPosition()
	if math.Abs(yAfterPeriod-160) > 1.0 {
		t.Errorf("y should return near base after 180 frames; got %v, base 160", yAfterPeriod)
	}
}

func TestCloudDrifter_FiresEightBoltsEvery90Frames(t *testing.T) {
	cd := NewCloudDrifter(300, 160, nil, 6.0)
	g := &Game{}
	for i := 0; i < 89; i++ {
		cd.Update(0, 0, g)
	}
	if len(g.enemyBullets) != 0 {
		t.Errorf("expected 0 bullets before frame 90, got %d", len(g.enemyBullets))
	}
	cd.Update(0, 0, g)
	if len(g.enemyBullets) != 8 {
		t.Errorf("expected 8 bullets at frame 90, got %d", len(g.enemyBullets))
	}
	// Verify all 8 are LightningBolts
	for i, b := range g.enemyBullets {
		if _, ok := b.(*LightningBolt); !ok {
			t.Errorf("bullet %d is not a *LightningBolt", i)
		}
	}
}

func TestCloudDrifter_BulletSpeedForwardedToBolts(t *testing.T) {
	const wantSpeed = 7.0
	cd := NewCloudDrifter(300, 160, nil, wantSpeed)
	g := &Game{}
	for i := 0; i < 90; i++ {
		cd.Update(0, 0, g)
	}
	if len(g.enemyBullets) != 8 {
		t.Fatalf("expected 8 bolts, got %d", len(g.enemyBullets))
	}
	// Each bolt travels at wantSpeed per frame. Check max displacement across all 8.
	maxDist := 0.0
	for _, b := range g.enemyBullets {
		x0, y0 := b.GetPosition()
		b.Update()
		x1, y1 := b.GetPosition()
		dist := math.Sqrt((x1-x0)*(x1-x0) + (y1-y0)*(y1-y0))
		if dist > maxDist {
			maxDist = dist
		}
	}
	if math.Abs(maxDist-wantSpeed) > 0.001 {
		t.Errorf("bolt speed: want %v, got %v", wantSpeed, maxDist)
	}
}

func TestLevel6_HasCloudDrifterSpawns(t *testing.T) {
	cfg := GetLevel6()
	found := false
	for _, sc := range cfg.SpawnConfigs {
		if sc.EnemyType == CloudDrifterType {
			found = true
			break
		}
	}
	if !found {
		t.Error("Level 6 should have at least one CloudDrifter spawn config")
	}
}

func TestCloudDrifter_StartsWithSixHP(t *testing.T) {
	cd := NewCloudDrifter(100, 100, nil, 6.0)
	if cd.IsDead() {
		t.Fatal("new CloudDrifter should not be dead")
	}
	cd.TakeDamage(5)
	if cd.IsDead() {
		t.Error("should survive 5 damage (has 6 HP)")
	}
	cd.TakeDamage(1)
	if !cd.IsDead() {
		t.Error("should die after 6 total damage")
	}
}
