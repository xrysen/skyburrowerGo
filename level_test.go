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

func TestWeatherType_LevelConfigField(t *testing.T) {
	cfg := LevelConfig{Weather: WeatherRain}
	if cfg.Weather != WeatherRain {
		t.Fatalf("expected WeatherRain, got %v", cfg.Weather)
	}
	cfg2 := LevelConfig{}
	if cfg2.Weather != WeatherNone {
		t.Fatalf("expected WeatherNone zero value, got %v", cfg2.Weather)
	}
}

func TestLevels1To5_WeatherNone(t *testing.T) {
	getters := []func() *LevelConfig{GetLevel1, GetLevel2, GetLevel3, GetLevel4, GetLevel5}
	for i, get := range getters {
		cfg := get()
		if cfg.Weather != WeatherNone {
			t.Errorf("level %d: want WeatherNone, got %v", i+1, cfg.Weather)
		}
	}
}

func TestLevels6To10_WeatherRain(t *testing.T) {
	getters := []func() *LevelConfig{GetLevel6, GetLevel7, GetLevel8, GetLevel9, GetLevel10}
	for i, get := range getters {
		cfg := get()
		if cfg.Weather != WeatherRain {
			t.Errorf("level %d: want WeatherRain, got %v", i+6, cfg.Weather)
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

// Level3 background path prefix used by World 3 levels.
const level3BgPrefix = "Levels/Level3/"

func hasSpawnType(cfgs []SpawnConfig, et EnemyType) bool {
	for _, c := range cfgs {
		if c.EnemyType == et {
			return true
		}
	}
	return false
}

func w2CameoHasEndFrame(cfgs []SpawnConfig, et EnemyType) bool {
	for _, c := range cfgs {
		if c.EnemyType == et {
			return c.EndFrame > 0
		}
	}
	return false
}

func TestGetLevel11_BasicConfig(t *testing.T) {
	cfg := GetLevel11()
	if cfg.WorldLevel != 11 {
		t.Errorf("WorldLevel: want 11, got %d", cfg.WorldLevel)
	}
	if cfg.BulletSpeed != 9.0 {
		t.Errorf("BulletSpeed: want 9.0, got %v", cfg.BulletSpeed)
	}
	if cfg.EndCondition != EndOnTimer {
		t.Errorf("EndCondition: want EndOnTimer, got %v", cfg.EndCondition)
	}
	if cfg.Duration != Seconds45 {
		t.Errorf("Duration: want Seconds45 (%d), got %d", Seconds45, cfg.Duration)
	}
	for i, p := range cfg.BackgroundPaths {
		if len(p) < len(level3BgPrefix) || p[:len(level3BgPrefix)] != level3BgPrefix {
			t.Errorf("BackgroundPaths[%d] = %q, want Level3 asset", i, p)
		}
	}
}

func TestGetLevel11_SpawnConfigs(t *testing.T) {
	cfg := GetLevel11()
	if !hasSpawnType(cfg.SpawnConfigs, DarkWingType) {
		t.Error("Level11 should have DarkWing spawn config")
	}
	if !hasSpawnType(cfg.SpawnConfigs, CloudDrifterType) {
		t.Error("Level11 should have CloudDrifter W2 cameo")
	}
	if !hasSpawnType(cfg.SpawnConfigs, LightningBugType) {
		t.Error("Level11 should have LightningBug W2 cameo")
	}
	if !w2CameoHasEndFrame(cfg.SpawnConfigs, CloudDrifterType) {
		t.Error("Level11 CloudDrifter cameo must have EndFrame set to phase out mid-level")
	}
	if !w2CameoHasEndFrame(cfg.SpawnConfigs, LightningBugType) {
		t.Error("Level11 LightningBug cameo must have EndFrame set to phase out mid-level")
	}
}

func TestGetLevel12_Config(t *testing.T) {
	cfg := GetLevel12()
	if cfg.WorldLevel != 12 {
		t.Errorf("WorldLevel: want 12, got %d", cfg.WorldLevel)
	}
	if cfg.BulletSpeed != 9.5 {
		t.Errorf("BulletSpeed: want 9.5, got %v", cfg.BulletSpeed)
	}
	if cfg.Duration != Seconds50 {
		t.Errorf("Duration: want Seconds50 (%d), got %d", Seconds50, cfg.Duration)
	}
	if !hasSpawnType(cfg.SpawnConfigs, DarkWingType) {
		t.Error("Level12 should have DarkWing spawn config")
	}
	if !hasSpawnType(cfg.SpawnConfigs, DynamiteBeetleType) {
		t.Error("Level12 should introduce DynamiteBeetle")
	}
}

func TestGetLevel13_Config(t *testing.T) {
	cfg := GetLevel13()
	if cfg.WorldLevel != 13 {
		t.Errorf("WorldLevel: want 13, got %d", cfg.WorldLevel)
	}
	if cfg.BulletSpeed != 10.0 {
		t.Errorf("BulletSpeed: want 10.0, got %v", cfg.BulletSpeed)
	}
	if cfg.Duration != Seconds55 {
		t.Errorf("Duration: want Seconds55 (%d), got %d", Seconds55, cfg.Duration)
	}
	if !hasSpawnType(cfg.SpawnConfigs, DarkWingType) {
		t.Error("Level13 should have DarkWing")
	}
	if !hasSpawnType(cfg.SpawnConfigs, DynamiteBeetleType) {
		t.Error("Level13 should have DynamiteBeetle")
	}
	if !hasSpawnType(cfg.SpawnConfigs, DrillDroneType) {
		t.Error("Level13 should introduce DrillDrone")
	}
	if hasSpawnType(cfg.SpawnConfigs, CloudDrifterType) || hasSpawnType(cfg.SpawnConfigs, LightningBugType) {
		t.Error("Level13 should have no W2 enemies")
	}
}

func TestGetLevel14_Config(t *testing.T) {
	cfg := GetLevel14()
	if cfg.WorldLevel != 14 {
		t.Errorf("WorldLevel: want 14, got %d", cfg.WorldLevel)
	}
	if cfg.BulletSpeed != 10.5 {
		t.Errorf("BulletSpeed: want 10.5, got %v", cfg.BulletSpeed)
	}
	if cfg.Duration != Seconds60 {
		t.Errorf("Duration: want Seconds60 (%d), got %d", Seconds60, cfg.Duration)
	}
	if !hasSpawnType(cfg.SpawnConfigs, DarkWingType) {
		t.Error("Level14 should have DarkWing")
	}
	if !hasSpawnType(cfg.SpawnConfigs, DynamiteBeetleType) {
		t.Error("Level14 should have DynamiteBeetle")
	}
	if !hasSpawnType(cfg.SpawnConfigs, DrillDroneType) {
		t.Error("Level14 should have DrillDrone")
	}
}

func TestGetLevel15_Config(t *testing.T) {
	cfg := GetLevel15()
	if cfg.WorldLevel != 15 {
		t.Errorf("WorldLevel: want 15, got %d", cfg.WorldLevel)
	}
	if cfg.BulletSpeed != 11.0 {
		t.Errorf("BulletSpeed: want 11.0, got %v", cfg.BulletSpeed)
	}
	if cfg.EndCondition != EndOnBossDeath {
		t.Errorf("EndCondition: want EndOnBossDeath, got %v", cfg.EndCondition)
	}
	if cfg.BossType != ForemanType {
		t.Errorf("BossType: want ForemanType, got %v", cfg.BossType)
	}
	if len(cfg.SpawnConfigs) == 0 {
		t.Error("Level15 should have light mine enemy pressure during boss fight")
	}
	if cfg.NextLevel != nil {
		t.Error("Level15 NextLevel should be nil (end of world)")
	}
}

func TestLevels11To15_Level3Backgrounds(t *testing.T) {
	getters := []func() *LevelConfig{GetLevel11, GetLevel12, GetLevel13, GetLevel14, GetLevel15}
	for i, get := range getters {
		cfg := get()
		for j, p := range cfg.BackgroundPaths {
			if len(p) < len(level3BgPrefix) || p[:len(level3BgPrefix)] != level3BgPrefix {
				t.Errorf("level %d bg[%d] = %q, want Level3 asset", i+11, j, p)
			}
		}
	}
}

func TestNextLevel_Chain11To15(t *testing.T) {
	getters := []func() *LevelConfig{GetLevel11, GetLevel12, GetLevel13, GetLevel14}
	expectedNext := []int{12, 13, 14, 15}
	for i, get := range getters {
		cfg := get()
		if cfg.NextLevel == nil {
			t.Errorf("level %d: NextLevel should not be nil", i+11)
			continue
		}
		next := cfg.NextLevel()
		if next.WorldLevel != expectedNext[i] {
			t.Errorf("level %d NextLevel: want %d, got %d", i+11, expectedNext[i], next.WorldLevel)
		}
	}
}

func TestGetLevel10_NextLevelIsLevel11(t *testing.T) {
	cfg := GetLevel10()
	if cfg.NextLevel == nil {
		t.Fatal("Level10 NextLevel should point to Level11")
	}
	if cfg.NextLevel().WorldLevel != 11 {
		t.Errorf("Level10 NextLevel: want WorldLevel 11, got %d", cfg.NextLevel().WorldLevel)
	}
}

func TestGetLevelForWorldSlot_Slots11To15(t *testing.T) {
	for slot := 11; slot <= 15; slot++ {
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
