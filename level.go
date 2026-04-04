package main

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
	NextLevel       func() *LevelConfig
	CoinSpawnConfig CoinSpawnConfig
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
			{EnemyType: FlutternatType, SpawnRate: 120, RandomY: true, MinSpawns: 1, MaxSpawns: 3, StartFrame: 100},
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
	default:
		return nil
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
			{EnemyType: FlutternatType, SpawnRate: 90, RandomY: true, MinSpawns: 2, MaxSpawns: 4, StartFrame: 60},
		},

		EndCondition: EndOnTimer,
		Duration:     Seconds45,
		NextLevel:    nil,
	}
}
