package main

import "testing"

func TestNewBackground_StoresWeather(t *testing.T) {
	bg := NewBackground(nil, WeatherRain)
	if bg.weather != WeatherRain {
		t.Fatalf("expected WeatherRain, got %v", bg.weather)
	}
}

func TestNewBackground_DefaultWeatherNone(t *testing.T) {
	bg := NewBackground(nil, WeatherNone)
	if bg.weather != WeatherNone {
		t.Fatalf("expected WeatherNone, got %v", bg.weather)
	}
}

func TestNewBackground_WeatherRain_InitDrops(t *testing.T) {
	bg := NewBackground(nil, WeatherRain)
	if len(bg.drops) < 90 || len(bg.drops) > 110 {
		t.Fatalf("expected ~100 drops for WeatherRain, got %d", len(bg.drops))
	}
}

func TestNewBackground_WeatherRain_DropsStaggered(t *testing.T) {
	bg := NewBackground(nil, WeatherRain)
	atZero := 0
	for _, d := range bg.drops {
		if d.y == 0 {
			atZero++
		}
	}
	if atZero == len(bg.drops) {
		t.Fatal("all drops start at y=0; expected staggered vertical positions")
	}
}

func TestBackground_Update_MovesDropsDownAndLeft(t *testing.T) {
	bg := NewBackground(nil, WeatherRain)
	// pin the first drop to a known safe position mid-screen
	bg.drops[0] = rainDrop{x: 300, y: 100, dy: 4.0, length: 8}
	bg.Update()
	d := bg.drops[0]
	if d.y <= 100 {
		t.Errorf("expected y > 100 after update, got %v", d.y)
	}
	if d.x >= 300 {
		t.Errorf("expected x < 300 after update, got %v", d.x)
	}
}

func TestBackground_Update_WrapsDropAtBottom(t *testing.T) {
	bg := NewBackground(nil, WeatherRain)
	bg.drops[0] = rainDrop{x: 300, y: 319, dy: 4.0, length: 8}
	bg.Update()
	d := bg.drops[0]
	if d.y >= ScreenHeight {
		t.Errorf("expected drop to wrap to top, got y=%v", d.y)
	}
	if d.x < 0 || d.x >= ScreenWidth {
		t.Errorf("wrapped drop x out of range: %v", d.x)
	}
}

func TestBackground_Update_WrapsDropAtLeftEdge(t *testing.T) {
	bg := NewBackground(nil, WeatherRain)
	bg.drops[0] = rainDrop{x: 0.5, y: 100, dy: 4.0, length: 8}
	bg.Update()
	d := bg.drops[0]
	if d.x < 0 {
		t.Errorf("expected drop to wrap before going negative, got x=%v", d.x)
	}
	if d.y != 0 {
		t.Errorf("wrapped drop should reset to y=0, got y=%v", d.y)
	}
}

func TestNewBackground_WeatherRain_LightningTimerInitialised(t *testing.T) {
	bg := NewBackground(nil, WeatherRain)
	if bg.lightningTimer < 180 || bg.lightningTimer > 480 {
		t.Errorf("expected lightningTimer in [180,480], got %d", bg.lightningTimer)
	}
}

func TestBackground_Update_DecrementsLightningTimer(t *testing.T) {
	bg := NewBackground(nil, WeatherRain)
	bg.lightningTimer = 300
	bg.Update()
	if bg.lightningTimer != 299 {
		t.Errorf("expected lightningTimer 299 after one Update, got %d", bg.lightningTimer)
	}
}

func TestBackground_Update_FlashTriggerOnTimerZero(t *testing.T) {
	bg := NewBackground(nil, WeatherRain)
	bg.lightningTimer = 1
	bg.lightningAlpha = 0
	bg.Update()
	if bg.lightningAlpha <= 0 {
		t.Errorf("expected lightningAlpha > 0 after timer fires, got %v", bg.lightningAlpha)
	}
	if bg.lightningTimer < 180 || bg.lightningTimer > 480 {
		t.Errorf("expected new timer in [180,480] after flash, got %d", bg.lightningTimer)
	}
}

func TestBackground_Update_LightningAlphaDecays(t *testing.T) {
	bg := NewBackground(nil, WeatherRain)
	bg.lightningTimer = 999 // prevent trigger
	bg.lightningAlpha = 0.35
	bg.Update()
	if bg.lightningAlpha >= 0.35 {
		t.Errorf("expected alpha to decay below 0.35, got %v", bg.lightningAlpha)
	}
	if bg.lightningAlpha < 0 {
		t.Errorf("alpha went negative: %v", bg.lightningAlpha)
	}
}

func TestBackground_Update_LightningAlphaClampedAtZero(t *testing.T) {
	bg := NewBackground(nil, WeatherRain)
	bg.lightningTimer = 999
	bg.lightningAlpha = 0.001 // less than one decay step
	bg.Update()
	if bg.lightningAlpha != 0 {
		t.Errorf("expected alpha clamped to 0, got %v", bg.lightningAlpha)
	}
}

func TestBackground_WeatherNone_NoLightning(t *testing.T) {
	bg := NewBackground(nil, WeatherNone)
	bg.Update()
	if bg.lightningTimer != 0 {
		t.Errorf("expected lightningTimer 0 for WeatherNone, got %d", bg.lightningTimer)
	}
	if bg.lightningAlpha != 0 {
		t.Errorf("expected lightningAlpha 0 for WeatherNone, got %v", bg.lightningAlpha)
	}
}

func TestNewBackground_WeatherNone_NoDrops(t *testing.T) {
	bg := NewBackground(nil, WeatherNone)
	if len(bg.drops) != 0 {
		t.Fatalf("expected 0 drops for WeatherNone, got %d", len(bg.drops))
	}
}
