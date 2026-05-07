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
	if cfg.NextLevel == nil {
		t.Fatal("Level15 NextLevel should point to Level16")
	}
	if cfg.NextLevel().WorldLevel != 16 {
		t.Errorf("Level15 NextLevel: want WorldLevel 16, got %d", cfg.NextLevel().WorldLevel)
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

func TestLevels1To15_NoStaticBackground(t *testing.T) {
	getters := []func() *LevelConfig{
		GetLevel1, GetLevel2, GetLevel3, GetLevel4, GetLevel5,
		GetLevel6, GetLevel7, GetLevel8, GetLevel9, GetLevel10,
		GetLevel11, GetLevel12, GetLevel13, GetLevel14, GetLevel15,
	}
	for i, get := range getters {
		cfg := get()
		if cfg.StaticBackgroundPath != "" {
			t.Errorf("level %d: expected empty StaticBackgroundPath, got %q", i+1, cfg.StaticBackgroundPath)
		}
	}
}

func TestLevelConfig_StaticBackgroundPathField(t *testing.T) {
	cfg := LevelConfig{StaticBackgroundPath: "Levels/Level4/lvl4-moon.png"}
	if cfg.StaticBackgroundPath != "Levels/Level4/lvl4-moon.png" {
		t.Fatalf("expected path set, got %q", cfg.StaticBackgroundPath)
	}
	empty := LevelConfig{}
	if empty.StaticBackgroundPath != "" {
		t.Fatalf("expected empty string zero value, got %q", empty.StaticBackgroundPath)
	}
}

const level4BgPrefix = "Levels/Level4/"

func TestGetLevel16_Config(t *testing.T) {
	cfg := GetLevel16()
	if cfg.WorldLevel != 16 {
		t.Errorf("WorldLevel: want 16, got %d", cfg.WorldLevel)
	}
	for i, p := range cfg.BackgroundPaths {
		if len(p) < len(level4BgPrefix) || p[:len(level4BgPrefix)] != level4BgPrefix {
			t.Errorf("BackgroundPaths[%d] = %q, want Level4 asset", i, p)
		}
	}
	if cfg.StaticBackgroundPath != "Levels/Level4/lvl4-moon.png" {
		t.Errorf("StaticBackgroundPath: want moon path, got %q", cfg.StaticBackgroundPath)
	}
	if cfg.EndCondition != EndOnTimer {
		t.Errorf("EndCondition: want EndOnTimer, got %v", cfg.EndCondition)
	}
	if !hasSpawnType(cfg.SpawnConfigs, BlightmothType) {
		t.Error("Level16 should spawn BlightmothType")
	}
}

func TestGetLevel15_NextLevelIsLevel16(t *testing.T) {
	cfg := GetLevel15()
	if cfg.NextLevel == nil {
		t.Fatal("Level15 NextLevel should point to Level16")
	}
	if cfg.NextLevel().WorldLevel != 16 {
		t.Errorf("Level15 NextLevel: want WorldLevel 16, got %d", cfg.NextLevel().WorldLevel)
	}
}

func TestGetLevelForWorldSlot_16(t *testing.T) {
	cfg := GetLevelForWorldSlot(16)
	if cfg == nil {
		t.Fatal("GetLevelForWorldSlot(16) returned nil")
	}
	if cfg.WorldLevel != 16 {
		t.Errorf("slot 16: WorldLevel want 16, got %d", cfg.WorldLevel)
	}
}

func TestGetLevel17_Config(t *testing.T) {
	cfg := GetLevel17()
	if cfg.WorldLevel != 17 {
		t.Errorf("WorldLevel: want 17, got %d", cfg.WorldLevel)
	}
	for i, p := range cfg.BackgroundPaths {
		if len(p) < len(level4BgPrefix) || p[:len(level4BgPrefix)] != level4BgPrefix {
			t.Errorf("BackgroundPaths[%d] = %q, want Level4 asset", i, p)
		}
	}
	if cfg.StaticBackgroundPath != "Levels/Level4/lvl4-moon.png" {
		t.Errorf("StaticBackgroundPath: want moon path, got %q", cfg.StaticBackgroundPath)
	}
	if cfg.EndCondition != EndOnTimer {
		t.Errorf("EndCondition: want EndOnTimer, got %v", cfg.EndCondition)
	}
	if !hasSpawnType(cfg.SpawnConfigs, HollowStagType) {
		t.Error("Level17 should spawn HollowStagType")
	}
}

func TestGetLevel16_NextLevelIsLevel17(t *testing.T) {
	cfg := GetLevel16()
	if cfg.NextLevel == nil {
		t.Fatal("Level16 NextLevel should point to Level17")
	}
	if cfg.NextLevel().WorldLevel != 17 {
		t.Errorf("Level16 NextLevel: want WorldLevel 17, got %d", cfg.NextLevel().WorldLevel)
	}
}

func TestGetLevelForWorldSlot_17(t *testing.T) {
	cfg := GetLevelForWorldSlot(17)
	if cfg == nil {
		t.Fatal("GetLevelForWorldSlot(17) returned nil")
	}
	if cfg.WorldLevel != 17 {
		t.Errorf("slot 17: WorldLevel want 17, got %d", cfg.WorldLevel)
	}
}

func spawnRate(cfgs []SpawnConfig, et EnemyType) int {
	for _, c := range cfgs {
		if c.EnemyType == et {
			return c.SpawnRate
		}
	}
	return 0
}

func TestGetLevel18_SpawnConfigs(t *testing.T) {
	cfg := GetLevel18()
	if !hasSpawnType(cfg.SpawnConfigs, BlightmothType) {
		t.Error("Level18 should spawn Blightmoth")
	}
	if !hasSpawnType(cfg.SpawnConfigs, HollowStagType) {
		t.Error("Level18 should spawn HollowStag")
	}
	// Higher density = lower spawn rate than intro levels (L16: 220, L17: 240)
	if r := spawnRate(cfg.SpawnConfigs, BlightmothType); r >= 220 {
		t.Errorf("Level18 Blightmoth spawn rate %d should be higher density than Level16 (rate 220)", r)
	}
	if r := spawnRate(cfg.SpawnConfigs, HollowStagType); r >= 240 {
		t.Errorf("Level18 HollowStag spawn rate %d should be higher density than Level17 (rate 240)", r)
	}
}

func TestGetLevelForWorldSlot_18And19(t *testing.T) {
	for _, slot := range []int{18, 19} {
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

func TestGetLevel19_SpawnConfigs(t *testing.T) {
	cfg := GetLevel19()
	if !hasSpawnType(cfg.SpawnConfigs, BlightmothType) {
		t.Error("Level19 should spawn Blightmoth")
	}
	if !hasSpawnType(cfg.SpawnConfigs, HollowStagType) {
		t.Error("Level19 should spawn HollowStag")
	}
	if !hasSpawnType(cfg.SpawnConfigs, WraithWhispType) {
		t.Error("Level19 should spawn WraithWhisp")
	}
	// Higher density than Level18 (Blightmoth: 180, HollowStag: 200)
	if r := spawnRate(cfg.SpawnConfigs, BlightmothType); r >= 180 {
		t.Errorf("Level19 Blightmoth spawn rate %d should be higher density than Level18 (rate 180)", r)
	}
	if r := spawnRate(cfg.SpawnConfigs, HollowStagType); r >= 200 {
		t.Errorf("Level19 HollowStag spawn rate %d should be higher density than Level18 (rate 200)", r)
	}
}

func TestGetLevel19_BasicConfig(t *testing.T) {
	cfg := GetLevel19()
	if cfg.WorldLevel != 19 {
		t.Errorf("WorldLevel: want 19, got %d", cfg.WorldLevel)
	}
	for i, p := range cfg.BackgroundPaths {
		if len(p) < len(level4BgPrefix) || p[:len(level4BgPrefix)] != level4BgPrefix {
			t.Errorf("BackgroundPaths[%d] = %q, want Level4 asset", i, p)
		}
	}
	if cfg.StaticBackgroundPath != "Levels/Level4/lvl4-moon.png" {
		t.Errorf("StaticBackgroundPath: want moon path, got %q", cfg.StaticBackgroundPath)
	}
	if cfg.EndCondition != EndOnTimer {
		t.Errorf("EndCondition: want EndOnTimer, got %v", cfg.EndCondition)
	}
}

func TestGetLevel18_NextLevelIsLevel19(t *testing.T) {
	cfg := GetLevel18()
	if cfg.NextLevel == nil {
		t.Fatal("Level18 NextLevel should point to Level19")
	}
	if cfg.NextLevel().WorldLevel != 19 {
		t.Errorf("Level18 NextLevel: want WorldLevel 19, got %d", cfg.NextLevel().WorldLevel)
	}
}

func TestGetLevel17_NextLevelIsLevel18(t *testing.T) {
	cfg := GetLevel17()
	if cfg.NextLevel == nil {
		t.Fatal("Level17 NextLevel should point to Level18")
	}
	if cfg.NextLevel().WorldLevel != 18 {
		t.Errorf("Level17 NextLevel: want WorldLevel 18, got %d", cfg.NextLevel().WorldLevel)
	}
}

func TestGetLevel18_BasicConfig(t *testing.T) {
	cfg := GetLevel18()
	if cfg.WorldLevel != 18 {
		t.Errorf("WorldLevel: want 18, got %d", cfg.WorldLevel)
	}
	for i, p := range cfg.BackgroundPaths {
		if len(p) < len(level4BgPrefix) || p[:len(level4BgPrefix)] != level4BgPrefix {
			t.Errorf("BackgroundPaths[%d] = %q, want Level4 asset", i, p)
		}
	}
	if cfg.StaticBackgroundPath != "Levels/Level4/lvl4-moon.png" {
		t.Errorf("StaticBackgroundPath: want moon path, got %q", cfg.StaticBackgroundPath)
	}
	if cfg.EndCondition != EndOnTimer {
		t.Errorf("EndCondition: want EndOnTimer, got %v", cfg.EndCondition)
	}
}

func TestGetLevelForWorldSlot_20(t *testing.T) {
	cfg := GetLevelForWorldSlot(20)
	if cfg == nil {
		t.Fatal("GetLevelForWorldSlot(20) returned nil")
	}
	if cfg.WorldLevel != 20 {
		t.Errorf("slot 20: WorldLevel want 20, got %d", cfg.WorldLevel)
	}
}

func TestGetLevel19_NextLevelIsLevel20(t *testing.T) {
	cfg := GetLevel19()
	if cfg.NextLevel == nil {
		t.Fatal("Level19 NextLevel should point to Level20")
	}
	if cfg.NextLevel().WorldLevel != 20 {
		t.Errorf("Level19 NextLevel: want WorldLevel 20, got %d", cfg.NextLevel().WorldLevel)
	}
}

func TestGetLevel20_Config(t *testing.T) {
	cfg := GetLevel20()
	if cfg.WorldLevel != 20 {
		t.Errorf("WorldLevel: want 20, got %d", cfg.WorldLevel)
	}
	if cfg.EndCondition != EndOnBossDeath {
		t.Errorf("EndCondition: want EndOnBossDeath, got %v", cfg.EndCondition)
	}
	if cfg.BossType != HeartwoodType {
		t.Errorf("BossType: want HeartwoodType, got %v", cfg.BossType)
	}
	if len(cfg.SpawnConfigs) != 0 {
		t.Errorf("SpawnConfigs: want empty, got %d entries", len(cfg.SpawnConfigs))
	}
	for i, p := range cfg.BackgroundPaths {
		if len(p) < len(level4BgPrefix) || p[:len(level4BgPrefix)] != level4BgPrefix {
			t.Errorf("BackgroundPaths[%d] = %q, want Level4 asset", i, p)
		}
	}
	if cfg.StaticBackgroundPath != "Levels/Level4/lvl4-moon.png" {
		t.Errorf("StaticBackgroundPath: want moon path, got %q", cfg.StaticBackgroundPath)
	}
	if cfg.NextLevel != nil {
		t.Error("Level20 NextLevel should be nil (final level)")
	}
}
