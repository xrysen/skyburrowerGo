package main

import "testing"

// Level2 background path prefix used by World 2 levels.
const level2BgPrefix = "Levels/Level2/"

func TestLevelConfig_HasBulletSpeedField(t *testing.T) {
	cfg := LevelConfig{BulletSpeed: 6.0}
	if cfg.BulletSpeed != 6.0 {
		t.Fatalf("expected BulletSpeed 6.0, got %v", cfg.BulletSpeed)
	}
}

func TestGetLevel6_Config(t *testing.T) {
	cfg := GetLevel6()
	if cfg.WorldLevel != 6 {
		t.Errorf("WorldLevel: want 6, got %d", cfg.WorldLevel)
	}
	if cfg.BulletSpeed != 6.0 {
		t.Errorf("BulletSpeed: want 6.0, got %v", cfg.BulletSpeed)
	}
	for i, p := range cfg.BackgroundPaths {
		if len(p) < len(level2BgPrefix) || p[:len(level2BgPrefix)] != level2BgPrefix {
			t.Errorf("BackgroundPaths[%d] = %q, want Level2 asset", i, p)
		}
	}
}

func TestGetLevel7_Config(t *testing.T) {
	cfg := GetLevel7()
	if cfg.WorldLevel != 7 {
		t.Errorf("WorldLevel: want 7, got %d", cfg.WorldLevel)
	}
	if cfg.BulletSpeed != 6.5 {
		t.Errorf("BulletSpeed: want 6.5, got %v", cfg.BulletSpeed)
	}
}

func TestGetLevel8_Config(t *testing.T) {
	cfg := GetLevel8()
	if cfg.WorldLevel != 8 {
		t.Errorf("WorldLevel: want 8, got %d", cfg.WorldLevel)
	}
	if cfg.BulletSpeed != 7.0 {
		t.Errorf("BulletSpeed: want 7.0, got %v", cfg.BulletSpeed)
	}
}

func TestGetLevel9_Config(t *testing.T) {
	cfg := GetLevel9()
	if cfg.WorldLevel != 9 {
		t.Errorf("WorldLevel: want 9, got %d", cfg.WorldLevel)
	}
	if cfg.BulletSpeed != 7.5 {
		t.Errorf("BulletSpeed: want 7.5, got %v", cfg.BulletSpeed)
	}
}

func TestGetLevel10_Config(t *testing.T) {
	cfg := GetLevel10()
	if cfg.WorldLevel != 10 {
		t.Errorf("WorldLevel: want 10, got %d", cfg.WorldLevel)
	}
	if cfg.BulletSpeed != 8.0 {
		t.Errorf("BulletSpeed: want 8.0, got %v", cfg.BulletSpeed)
	}
	if cfg.EndCondition != EndOnBossDeath {
		t.Errorf("EndCondition: want EndOnBossDeath, got %v", cfg.EndCondition)
	}
	for i, p := range cfg.BackgroundPaths {
		if len(p) < len(level2BgPrefix) || p[:len(level2BgPrefix)] != level2BgPrefix {
			t.Errorf("BackgroundPaths[%d] = %q, want Level2 asset", i, p)
		}
	}
}

func TestGetLevelForWorldSlot_Slots6To10(t *testing.T) {
	for slot := 6; slot <= 10; slot++ {
		cfg := GetLevelForWorldSlot(slot)
		if cfg == nil {
			t.Errorf("GetLevelForWorldSlot(%d) returned nil", slot)
			continue
		}
		if cfg.WorldLevel != slot {
			t.Errorf("slot %d: WorldLevel want %d, got %d", slot, slot, cfg.WorldLevel)
		}
	}
}
