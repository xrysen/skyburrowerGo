package main

type SpawnConfig struct {
	EnemyType EnemyType
	SpawnRate int
	RandomY   bool
}

type LevelConfig struct {
	Name            string
	BackgroundPaths [4]string
	SpawnConfigs    []SpawnConfig
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
			{EnemyType: FlutternatType, SpawnRate: 120, RandomY: true},
		},
	}
}
