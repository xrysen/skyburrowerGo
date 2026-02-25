package main

type LevelEndCondition int

const (
	EndOnTimer LevelEndCondition = iota
	EndOnBossDeath
)

const (
	FPS = 60

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
	Name            string
	BackgroundPaths [4]string
	SpawnConfigs    []SpawnConfig
	EndCondition    LevelEndCondition
	Duration        int
	BossType        EnemyType
}

func GetLevel1() *LevelConfig {
	return &LevelConfig{
		Name: "Forest Level",
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
	}
}
