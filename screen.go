package main

const (
	ScreenWidth  = 640
	ScreenHeight = 320
)

type Screen int

const (
	ScreenWorldMap Screen = iota
	ScreenPlaying
	ScreenGameOver
)
