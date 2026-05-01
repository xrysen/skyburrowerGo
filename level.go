package main

const (
	WorldLevelCount = 20
	// CarrotsPerLevel is how many bonus carrots exist per level (world-map layout uses 5 slots).
	CarrotsPerLevel = 5
)

type LevelEndCondition int

const (
	EndOnTimer LevelEndCondition = iota
	EndOnBossDeath
)

const (
	FPS = 60

	FadeOutDuration = 4 * FPS

	Seconds30 = 30 * FPS
	Seconds45 = 45 * FPS
	Seconds60 = 60 * FPS
	Seconds90 = 90 * FPS
	Minutes2  = 120 * FPS
)

type SpawnConfig struct {
	EnemyType  EnemyType
	SpawnRate  int
	RandomY    bool
	MinSpawns  int
	MaxSpawns  int
	StartFrame int
	EndFrame   int
}

type LevelConfig struct {
	WorldLevel      int // 1-based slot on the world map; used for unlock progression
	Name            string
	BackgroundPaths [4]string
	SpawnConfigs    []SpawnConfig
	EndCondition    LevelEndCondition
	Duration        int
	BossType        EnemyType
	BossHealth      int
	NextLevel       func() *LevelConfig
	CoinSpawnConfig CoinSpawnConfig
	BulletSpeed     float64 // base bullet speed for enemies in this level
}

type CoinSpawnConfig struct {
	SpawnRate int
	RandomY   bool
}

func GetLevel1() *LevelConfig {
	return &LevelConfig{
		WorldLevel: 1,
		Name:       "Forest Level",
		BackgroundPaths: [4]string{
			"Levels/Level1/lvl1-1.png",
			"Levels/Level1/lvl1-2.png",
			"Levels/Level1/lvl1-3.png",
			"Levels/Level1/lvl1-4.png",
		},
		SpawnConfigs: []SpawnConfig{
			{EnemyType: FlutternatType, SpawnRate: 180, RandomY: true, MinSpawns: 1, MaxSpawns: 2, StartFrame: 120},
		},

		EndCondition: EndOnTimer,
		Duration:     Seconds30,
		NextLevel:    GetLevel2,
		CoinSpawnConfig: CoinSpawnConfig{
			SpawnRate: 180,
			RandomY:   true,
		},
	}
}

// GetLevelForWorldSlot returns the level for a world-map slot (1–20), or nil if not implemented yet.
func GetLevelForWorldSlot(slot int) *LevelConfig {
	switch slot {
	case 1:
		return GetLevel1()
	case 2:
		return GetLevel2()
	case 3:
		return GetLevel3()
	case 4:
		return GetLevel4()
	case 5:
		return GetLevel5()
	case 6:
		return GetLevel6()
	case 7:
		return GetLevel7()
	case 8:
		return GetLevel8()
	case 9:
		return GetLevel9()
	case 10:
		return GetLevel10()
	default:
		return nil
	}
}

func GetLevel6() *LevelConfig {
	return &LevelConfig{
		WorldLevel: 6,
		Name:       "ThunderCloud Ravine",
		BackgroundPaths: [4]string{
			"Levels/Level2/lvl2-1.png",
			"Levels/Level2/lvl2-2.png",
			"Levels/Level2/lvl2-3.png",
			"Levels/Level2/lvl2-4.png",
		},
		SpawnConfigs: []SpawnConfig{
			// CloudDrifter and LightningBug introduced together — CloudDrifter first
			{EnemyType: CloudDrifterType, SpawnRate: 240, RandomY: true, MinSpawns: 1, MaxSpawns: 1, StartFrame: 10},
			{EnemyType: LightningBugType, SpawnRate: 240, RandomY: true, MinSpawns: 1, MaxSpawns: 1, StartFrame: 120},
			// Flutternat joins as familiar background and fades out near the end
			{EnemyType: FlutternatType, SpawnRate: 260, RandomY: true, MinSpawns: 1, MaxSpawns: 1, StartFrame: 300, EndFrame: 2000},
			// Sporespinner brief mid-level cameo
			{EnemyType: SporespinnerType, SpawnRate: 320, RandomY: true, MinSpawns: 1, MaxSpawns: 1, StartFrame: 600, EndFrame: 2200},
		},
		EndCondition:    EndOnTimer,
		Duration:        Seconds45,
		BulletSpeed:     6.0,
		NextLevel:       GetLevel7,
		CoinSpawnConfig: CoinSpawnConfig{SpawnRate: 180, RandomY: true},
	}
}

func GetLevel7() *LevelConfig {
	return &LevelConfig{
		WorldLevel: 7,
		Name:       "ThunderCloud Ravine",
		BackgroundPaths: [4]string{
			"Levels/Level2/lvl2-1.png",
			"Levels/Level2/lvl2-2.png",
			"Levels/Level2/lvl2-3.png",
			"Levels/Level2/lvl2-4.png",
		},
		SpawnConfigs: []SpawnConfig{
			// StormSprite is the star — appears immediately
			{EnemyType: StormSpriteType, SpawnRate: 260, RandomY: true, MinSpawns: 1, MaxSpawns: 1, StartFrame: 10},
			// CloudDrifter and LightningBug carry over from L6 as background
			{EnemyType: CloudDrifterType, SpawnRate: 300, RandomY: true, MinSpawns: 1, MaxSpawns: 1, StartFrame: 300},
			{EnemyType: LightningBugType, SpawnRate: 280, RandomY: true, MinSpawns: 1, MaxSpawns: 1, StartFrame: 300},
			// Flutternat — early cameo, retires at the halfway point
			{EnemyType: FlutternatType, SpawnRate: 300, RandomY: true, MinSpawns: 1, MaxSpawns: 1, StartFrame: 300, EndFrame: 1350},
		},
		EndCondition:    EndOnTimer,
		Duration:        Seconds45,
		BulletSpeed:     6.5,
		NextLevel:       GetLevel8,
		CoinSpawnConfig: CoinSpawnConfig{SpawnRate: 180, RandomY: true},
	}
}

func GetLevel8() *LevelConfig {
	return &LevelConfig{
		WorldLevel: 8,
		Name:       "ThunderCloud Ravine",
		BackgroundPaths: [4]string{
			"Levels/Level2/lvl2-1.png",
			"Levels/Level2/lvl2-2.png",
			"Levels/Level2/lvl2-3.png",
			"Levels/Level2/lvl2-4.png",
		},
		SpawnConfigs: []SpawnConfig{
			// All three W2 enemies as equals — W1 fully retired
			{EnemyType: CloudDrifterType, SpawnRate: 260, RandomY: true, MinSpawns: 1, MaxSpawns: 1, StartFrame: 10},
			{EnemyType: LightningBugType, SpawnRate: 260, RandomY: true, MinSpawns: 1, MaxSpawns: 1, StartFrame: 10},
			{EnemyType: StormSpriteType, SpawnRate: 280, RandomY: true, MinSpawns: 1, MaxSpawns: 1, StartFrame: 10},
			// ThistleTurret mid-level cameo for texture
			{EnemyType: ThistleTurretType, SpawnRate: 380, RandomY: false, MinSpawns: 1, MaxSpawns: 1, StartFrame: 900, EndFrame: 3300},
		},
		EndCondition:    EndOnTimer,
		Duration:        Seconds60,
		BulletSpeed:     7.0,
		NextLevel:       GetLevel9,
		CoinSpawnConfig: CoinSpawnConfig{SpawnRate: 160, RandomY: true},
	}
}

func GetLevel9() *LevelConfig {
	return &LevelConfig{
		WorldLevel: 9,
		Name:       "Tempest Gauntlet",
		BackgroundPaths: [4]string{
			"Levels/Level2/lvl2-1.png",
			"Levels/Level2/lvl2-2.png",
			"Levels/Level2/lvl2-3.png",
			"Levels/Level2/lvl2-4.png",
		},
		SpawnConfigs: []SpawnConfig{
			// All three W2 enemies as equals — W1 fully retired
			{EnemyType: CloudDrifterType, SpawnRate: 220, RandomY: true, MinSpawns: 1, MaxSpawns: 2, StartFrame: 60},
			{EnemyType: LightningBugType, SpawnRate: 200, RandomY: true, MinSpawns: 1, MaxSpawns: 2, StartFrame: 60},
			{EnemyType: StormSpriteType, SpawnRate: 250, RandomY: true, MinSpawns: 1, MaxSpawns: 1, StartFrame: 60},
		},
		EndCondition:    EndOnTimer,
		Duration:        Seconds60,
		BulletSpeed:     7.5,
		NextLevel:       GetLevel10,
		CoinSpawnConfig: CoinSpawnConfig{SpawnRate: 140, RandomY: true},
	}
}

func GetLevel10() *LevelConfig {
	return &LevelConfig{
		WorldLevel: 10,
		Name:       "ThunderCloud Ravine: Summit",
		BackgroundPaths: [4]string{
			"Levels/Level2/lvl2-1.png",
			"Levels/Level2/lvl2-2.png",
			"Levels/Level2/lvl2-3.png",
			"Levels/Level2/lvl2-4.png",
		},
		SpawnConfigs: []SpawnConfig{
			// Light W2 pressure during the boss fight — W1 fully retired
			{EnemyType: CloudDrifterType, SpawnRate: 380, RandomY: true, MinSpawns: 1, MaxSpawns: 1, StartFrame: 120},
			{EnemyType: LightningBugType, SpawnRate: 340, RandomY: true, MinSpawns: 1, MaxSpawns: 1, StartFrame: 120},
		},
		EndCondition:    EndOnBossDeath,
		BossType:        ThunderCrabType,
		BossHealth:      180,
		BulletSpeed:     8.0,
		NextLevel:       nil,
		CoinSpawnConfig: CoinSpawnConfig{SpawnRate: 120, RandomY: true},
	}
}

func GetLevel2() *LevelConfig {
	return &LevelConfig{
		WorldLevel: 2,
		Name:       "Deep Forest",
		BackgroundPaths: [4]string{
			"Levels/Level1/lvl1-1.png",
			"Levels/Level1/lvl1-2.png",
			"Levels/Level1/lvl1-3.png",
			"Levels/Level1/lvl1-4.png",
		},
		SpawnConfigs: []SpawnConfig{
			{EnemyType: FlutternatType, SpawnRate: 150, RandomY: true, MinSpawns: 1, MaxSpawns: 3, StartFrame: 100},
			{EnemyType: SporespinnerType, SpawnRate: 160, RandomY: true, MinSpawns: 1, MaxSpawns: 2, StartFrame: 80},
		},

		EndCondition: EndOnTimer,
		Duration:     Seconds45,
		NextLevel:    GetLevel3,
		CoinSpawnConfig: CoinSpawnConfig{
			SpawnRate: 160,
			RandomY:   true,
		},
	}
}

func GetLevel3() *LevelConfig {
	return &LevelConfig{
		WorldLevel: 3,
		Name:       "Thistle Grove",
		BackgroundPaths: [4]string{
			"Levels/Level1/lvl1-1.png",
			"Levels/Level1/lvl1-2.png",
			"Levels/Level1/lvl1-3.png",
			"Levels/Level1/lvl1-4.png",
		},
		SpawnConfigs: []SpawnConfig{
			{EnemyType: FlutternatType, SpawnRate: 130, RandomY: true, MinSpawns: 2, MaxSpawns: 3, StartFrame: 80},
			{EnemyType: SporespinnerType, SpawnRate: 140, RandomY: true, MinSpawns: 1, MaxSpawns: 2, StartFrame: 60},
			{EnemyType: ThistleTurretType, SpawnRate: 200, RandomY: false, MinSpawns: 1, MaxSpawns: 1, StartFrame: 120},
		},

		EndCondition: EndOnTimer,
		Duration:     Seconds45,
		NextLevel:    GetLevel4,
		CoinSpawnConfig: CoinSpawnConfig{
			SpawnRate: 140,
			RandomY:   true,
		},
	}
}

func GetLevel4() *LevelConfig {
	return &LevelConfig{
		WorldLevel: 4,
		Name:       "Forest Chaos",
		BackgroundPaths: [4]string{
			"Levels/Level1/lvl1-1.png",
			"Levels/Level1/lvl1-2.png",
			"Levels/Level1/lvl1-3.png",
			"Levels/Level1/lvl1-4.png",
		},
		SpawnConfigs: []SpawnConfig{
			{EnemyType: FlutternatType, SpawnRate: 100, RandomY: true, MinSpawns: 2, MaxSpawns: 4, StartFrame: 60},
			{EnemyType: SporespinnerType, SpawnRate: 120, RandomY: true, MinSpawns: 2, MaxSpawns: 3, StartFrame: 40},
			{EnemyType: ThistleTurretType, SpawnRate: 180, RandomY: false, MinSpawns: 1, MaxSpawns: 2, StartFrame: 100},
		},

		EndCondition: EndOnTimer,
		Duration:     Seconds60,
		NextLevel:    GetLevel5,
		CoinSpawnConfig: CoinSpawnConfig{
			SpawnRate: 120,
			RandomY:   true,
		},
	}
}

func GetLevel5() *LevelConfig {
	return &LevelConfig{
		WorldLevel: 5,
		Name:       "Owlbert's Lair",
		BackgroundPaths: [4]string{
			"Levels/Level1/lvl1-1.png",
			"Levels/Level1/lvl1-2.png",
			"Levels/Level1/lvl1-3.png",
			"Levels/Level1/lvl1-4.png",
		},
		SpawnConfigs: []SpawnConfig{
			// Increased Flutternat spawning for boss level
			{EnemyType: FlutternatType, SpawnRate: 180, RandomY: true, MinSpawns: 2, MaxSpawns: 3, StartFrame: 120},
		},

		EndCondition: EndOnBossDeath,
		BossType:     OwlbertType,
		NextLevel:    nil,
		CoinSpawnConfig: CoinSpawnConfig{
			SpawnRate: 120, // Still spawn coins during boss fight
			RandomY:   true,
		},
	}
}
