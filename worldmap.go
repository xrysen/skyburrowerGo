package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

const (
	gearsBgY      = 0
	upgradeLayerY = 10

	// Upgrade slot positions (based on UpgradeLayer.png layout)
	upgradeSlotStartX = 140 // X position of plus signs (right of black bars)
	upgradeSlotStartY = 23  // Y position of first upgrade slot (coin slot) - 41-10 to compensate for upgradeLayerY
	upgradeSlotStepY  = 29  // Vertical distance between upgrade slots
	upgradeSlotCount  = 8   // Number of upgrade slots total

	// Upgrade bar positions (left of plus signs)
	upgradeBarStartX = 80 // X position of upgrade bars
	upgradeBarStartY = 25 // Y position of first upgrade bar

	levelGridOriginX = 353 // screen: top-left of first cell (row 0, col 0)
	levelGridOriginY = 57
	levelGridRowStep = 52 // screen Y for row 1 starts at 109 → 109 - 57
	// Horizontal distance from one column’s left edge to the next (your editor grid),
	// not necessarily PNG width + gap — tune if x looks off vs wm.png.
	levelGridColPitch = 39
	levelGridCols     = 5
	levelGridRows     = 4

	levelNumberDigitGap = 1  // horizontal space between digits when drawing 10–20
	levelNumberOffsetY  = -3 // shift level digits (and carrot anchor) up in the cell

	carrotSideGap      = 5 // horizontal space from number to left/right carrots
	carrotBelowGap     = 3 // vertical gap from number bottom to the row of 3
	carrotBelowSpacing = 2 // horizontal gap between the three bottom carrots
)

type WorldMap struct {
	gearsBg           *ebiten.Image
	upgradeLayer      *ebiten.Image
	levelSelectLayer  *ebiten.Image
	levelSelectButton *ebiten.Image
	lockIcon          *ebiten.Image
	plusSign          *ebiten.Image
	upgradeEmpty      *ebiten.Image
	upgradeFilled     *ebiten.Image
	levelDigits       [10]*ebiten.Image
	carrotEmpty       *ebiten.Image
	carrotFull        *ebiten.Image
	font              *BitmapFont
	tooltipFont       font.Face
}

func NewWorldMap(assets *Assets, font *BitmapFont) *WorldMap {
	tooltipFont, err := loadSystemFont()
	if err != nil {
		// Fallback to bitmap font if system font fails
		tooltipFont = nil
	}

	return &WorldMap{
		gearsBg:           assets.GearsBg,
		upgradeLayer:      assets.UpgradeLayer,
		levelSelectLayer:  assets.LevelSelectLayer,
		levelSelectButton: assets.LevelSelectButton,
		lockIcon:          assets.LockIcon,
		plusSign:          assets.PlusSign,
		upgradeEmpty:      assets.UpgradeEmpty,
		upgradeFilled:     assets.UpgradeFilled,
		levelDigits:       assets.LevelDigits,
		carrotEmpty:       assets.LsCarrotEmpty,
		carrotFull:        assets.LsCarrotFull,
		font:              font,
		tooltipFont:       tooltipFont,
	}
}

func (wm *WorldMap) Update(g *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && g.highestUnlockedLevel >= 1 {
		if cfg := GetLevelForWorldSlot(1); cfg != nil {
			g.startLevel(cfg)
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()

		// Check for shift-click for testing purposes
		isShiftPressed := ebiten.IsKeyPressed(ebiten.KeyShift) || ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShiftRight)

		if level, ok := wm.levelSlotAtScreen(mx, my); ok {
			// Allow shift-click to open any level for testing
			if isShiftPressed || level <= g.highestUnlockedLevel {
				if cfg := GetLevelForWorldSlot(level); cfg != nil {
					g.startLevel(cfg)
				}
			}
		} else if upgradeType, ok := wm.plusSignAtScreen(mx, my); ok {
			wm.handleUpgradeClick(g, upgradeType)
		}
	}
}

// levelSlotAtScreen returns the 1-based world level index if (x, y) is inside that slot’s button.
func (wm *WorldMap) levelSlotAtScreen(x, y int) (level int, ok bool) {
	bw := wm.levelSelectButton.Bounds().Dx()
	bh := wm.levelSelectButton.Bounds().Dy()
	for row := 0; row < levelGridRows; row++ {
		for col := 0; col < levelGridCols; col++ {
			x0 := levelGridOriginX + col*levelGridColPitch
			y0 := levelGridOriginY + row*levelGridRowStep
			if x >= x0 && y >= y0 && x < x0+bw && y < y0+bh {
				return row*levelGridCols + col + 1, true
			}
		}
	}
	return 0, false
}

// plusSignAtScreen returns the upgrade type if (x, y) is inside a plus sign button.
func (wm *WorldMap) plusSignAtScreen(x, y int) (upgradeType UpgradeType, ok bool) {
	if wm.plusSign == nil {
		return 0, false
	}

	pw := wm.plusSign.Bounds().Dx()
	ph := wm.plusSign.Bounds().Dy()

	for i := 0; i < upgradeSlotCount; i++ {
		// Skip the coin slot (index 0)
		if i == 0 {
			continue
		}

		px := upgradeSlotStartX
		py := upgradeSlotStartY + i*upgradeSlotStepY

		if x >= px && y >= py && x < px+pw && y < py+ph {
			// Map upgrade slot to upgrade type (skip coin slot)
			upgradeType := UpgradeType(i - 1)
			if upgradeType >= UpgradeCount {
				return 0, false
			}
			return upgradeType, true
		}
	}
	return 0, false
}

// handleUpgradeClick processes an upgrade purchase when a plus sign is clicked
func (wm *WorldMap) handleUpgradeClick(g *Game, upgradeType UpgradeType) {
	cost := wm.getUpgradeCost(upgradeType, g.upgrades[upgradeType].Level)
	maxLevel := 7
	if upgradeType == UpgradeBulletCount {
		maxLevel = 4
	}

	if g.player.coins >= cost && g.upgrades[upgradeType].Level < maxLevel {
		g.player.coins -= cost
		g.upgrades[upgradeType].Level++
		wm.applyUpgradeEffects(g, upgradeType)
	}
}

// getUpgradeCost returns the cost for upgrading a specific type to the next level
func (wm *WorldMap) getUpgradeCost(upgradeType UpgradeType, currentLevel int) int {
	// Base cost for first level of each upgrade type
	baseCosts := map[UpgradeType]int{
		UpgradeHealth:         50,
		UpgradeBulletStrength: 75,
		UpgradeBulletSpeed:    60,
		UpgradeBulletCount:    400,
		UpgradeSpeed:          80,
		UpgradeMagnetism:      90,
		UpgradeLuck:           65,
	}

	maxLevel := 7
	if upgradeType == UpgradeBulletCount {
		maxLevel = 4
	}

	if currentLevel >= maxLevel {
		return 999999 // Max level reached
	}

	baseCost := baseCosts[upgradeType]
	// Cost increases exponentially: baseCost * (level + 1)
	return baseCost * (currentLevel + 1)
}

// applyUpgradeEffects applies the effects of an upgrade to the player
func (wm *WorldMap) applyUpgradeEffects(g *Game, upgradeType UpgradeType) {
	switch upgradeType {
	case UpgradeHealth:
		// Each level adds 1 max health
		g.player.maxHealth = 3 + g.upgrades[upgradeType].Level
		g.player.health = g.player.maxHealth // Heal to full max health
	case UpgradeBulletStrength:
		// Each level increases bullet damage by 1
		g.player.bulletDamage = 1 + g.upgrades[upgradeType].Level
	case UpgradeBulletSpeed:
		// Each level reduces fire interval for faster shooting
		g.player.fireInterval = 30 - g.upgrades[upgradeType].Level*3
	case UpgradeBulletCount:
		// Each level adds 1 more bullet (max 3 additional)
		g.player.bulletCount = 1 + g.upgrades[upgradeType].Level
	case UpgradeSpeed:
		// Speed levels are handled in player.getMovementSpeed()
	case UpgradeMagnetism:
		// Each level increases magnet range by 20 pixels
		g.player.magnetRange = 50.0 + float64(g.upgrades[upgradeType].Level)*20.0
	case UpgradeLuck:
		// Luck affects coin spawn rate and big coin chance
		g.player.luck = 1 + g.upgrades[upgradeType].Level
	}
}

func (wm *WorldMap) Draw(screen *ebiten.Image, g *Game) {
	gearsOp := &ebiten.DrawImageOptions{}
	gearsOp.GeoM.Translate(0, float64(gearsBgY))
	screen.DrawImage(wm.gearsBg, gearsOp)
	upgradeOp := &ebiten.DrawImageOptions{}
	upgradeOp.GeoM.Translate(0, float64(upgradeLayerY))
	screen.DrawImage(wm.upgradeLayer, upgradeOp)
	levelSelectOp := &ebiten.DrawImageOptions{}
	levelSelectOp.GeoM.Translate(0, 0)
	screen.DrawImage(wm.levelSelectLayer, levelSelectOp)
	wm.drawLevelSelectButtons(screen, g)
	wm.drawUpgradeBars(screen, g)
	wm.drawUpgradePlusSigns(screen)
	wm.drawUpgradeHoverText(screen, g)
	wm.drawCoinCount(screen, g)
	wm.drawTestingHint(screen)
}

func (wm *WorldMap) drawLevelSelectButtons(screen *ebiten.Image, g *Game) {
	btn := wm.levelSelectButton

	for row := 0; row < levelGridRows; row++ {
		for col := 0; col < levelGridCols; col++ {
			op := &ebiten.DrawImageOptions{}
			x := float64(levelGridOriginX + col*levelGridColPitch)
			y := float64(levelGridOriginY + row*levelGridRowStep)
			op.GeoM.Translate(x, y)
			screen.DrawImage(btn, op)
			level := row*levelGridCols + col + 1
			if level <= g.highestUnlockedLevel {
				wm.drawLevelNumberOnButton(screen, x, y, level)
				var mask uint8
				if level >= 1 && level <= WorldLevelCount {
					mask = g.levelCarrotMask[level-1]
				}
				wm.drawLevelCarrots(screen, x, y, level, mask)
			} else {
				wm.drawLockOnButton(screen, x, y)
			}
		}
	}
}

func (wm *WorldMap) drawLockOnButton(screen *ebiten.Image, cellX, cellY float64) {
	lock := wm.lockIcon
	lw := float64(lock.Bounds().Dx())
	lh := float64(lock.Bounds().Dy())
	btnW := float64(wm.levelSelectButton.Bounds().Dx())
	btnH := float64(wm.levelSelectButton.Bounds().Dy())
	lx := cellX + (btnW-lw)/2
	ly := cellY + (btnH-lh)/2
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(lx, ly)
	screen.DrawImage(lock, op)
}

func levelNumberDigitIndices(n int) []int {
	switch {
	case n >= 1 && n <= 9:
		return []int{n}
	case n >= 10 && n <= 20:
		return []int{n / 10, n % 10}
	default:
		return nil
	}
}

// levelNumberRect returns the digit bounds for a level label (for carrot placement).
func (wm *WorldMap) levelNumberRect(cellX, cellY float64, level int) (left, top, right, bottom float64, ok bool) {
	ds := levelNumberDigitIndices(level)
	if len(ds) == 0 {
		return 0, 0, 0, 0, false
	}
	btnW := float64(wm.levelSelectButton.Bounds().Dx())
	btnH := float64(wm.levelSelectButton.Bounds().Dy())

	var totalW float64
	maxDigitH := 0
	for i, di := range ds {
		if i > 0 {
			totalW += float64(levelNumberDigitGap)
		}
		img := wm.levelDigits[di]
		totalW += float64(img.Bounds().Dx())
		h := img.Bounds().Dy()
		if h > maxDigitH {
			maxDigitH = h
		}
	}

	left = cellX + (btnW-totalW)/2
	top = cellY + (btnH-float64(maxDigitH))/2 + float64(levelNumberOffsetY)
	right = left + totalW
	bottom = top + float64(maxDigitH)
	return left, top, right, bottom, true
}

// drawLevelCarrots draws 5 carrots around the level number: bit 0 left, 1–3 below (L→R), 4 right.
func (wm *WorldMap) drawLevelCarrots(screen *ebiten.Image, cellX, cellY float64, level int, mask uint8) {
	nl, nt, nr, nb, ok := wm.levelNumberRect(cellX, cellY, level)
	if !ok {
		return
	}
	cw := float64(wm.carrotEmpty.Bounds().Dx())
	ch := float64(wm.carrotEmpty.Bounds().Dy())
	sideGap := float64(carrotSideGap)
	belowGap := float64(carrotBelowGap)
	belowSp := float64(carrotBelowSpacing)
	vCarrotY := nt + (nb-nt-ch)/2

	type pos struct{ x, y float64 }
	positions := [CarrotsPerLevel]pos{}

	// 0: left of number (extra horizontal space vs below row)
	positions[0] = pos{nl - sideGap - cw, vCarrotY}

	// 1–3: below number, evenly spaced
	belowY := nb + belowGap
	rowW := 3*cw + 2*belowSp
	startBelowX := nl + (nr-nl-rowW)/2
	for i := 0; i < 3; i++ {
		positions[1+i] = pos{startBelowX + float64(i)*(cw+belowSp), belowY}
	}

	// 4: right of number
	positions[4] = pos{nr + sideGap, vCarrotY}

	for slot := 0; slot < CarrotsPerLevel; slot++ {
		img := wm.carrotEmpty
		if (mask>>slot)&1 != 0 {
			img = wm.carrotFull
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(positions[slot].x, positions[slot].y)
		screen.DrawImage(img, op)
	}
}

func (wm *WorldMap) drawLevelNumberOnButton(screen *ebiten.Image, cellX, cellY float64, level int) {
	ds := levelNumberDigitIndices(level)
	if len(ds) == 0 {
		return
	}
	btnW := float64(wm.levelSelectButton.Bounds().Dx())
	btnH := float64(wm.levelSelectButton.Bounds().Dy())

	var totalW float64
	maxDigitH := 0
	for i, di := range ds {
		if i > 0 {
			totalW += float64(levelNumberDigitGap)
		}
		img := wm.levelDigits[di]
		totalW += float64(img.Bounds().Dx())
		h := img.Bounds().Dy()
		if h > maxDigitH {
			maxDigitH = h
		}
	}

	startX := cellX + (btnW-totalW)/2
	startY := cellY + (btnH-float64(maxDigitH))/2 + float64(levelNumberOffsetY)

	x := startX
	for i, di := range ds {
		if i > 0 {
			x += float64(levelNumberDigitGap)
		}
		img := wm.levelDigits[di]
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, startY)
		screen.DrawImage(img, op)
		x += float64(img.Bounds().Dx())
	}
}

func (wm *WorldMap) drawCoinCount(screen *ebiten.Image, g *Game) {
	if wm.font == nil {
		return
	}
	coinText := fmt.Sprintf("%d", g.player.coins)
	wm.font.DrawText(screen, coinText, 80, 27, 1.5)
}

// drawUpgradeBars draws upgrade bars (empty or filled) for each upgrade slot
func (wm *WorldMap) drawUpgradeBars(screen *ebiten.Image, g *Game) {
	for i := 0; i < upgradeSlotCount; i++ {
		// Skip the coin slot (index 0)
		if i == 0 {
			continue
		}

		// Map upgrade slot to upgrade type (skip coin slot)
		upgradeType := UpgradeType(i - 1)
		if upgradeType >= UpgradeCount {
			continue
		}

		// Determine max bars for this upgrade type
		maxBars := 7
		if upgradeType == UpgradeBulletCount {
			maxBars = 4
		}

		// Draw bars for this upgrade (left to right)
		for bar := 0; bar < maxBars; bar++ {
			x := float64(upgradeBarStartX + bar*8) // 8 pixels spacing between bars
			y := float64(upgradeBarStartY + i*upgradeSlotStepY)

			// Choose empty or filled based on upgrade level
			var img *ebiten.Image
			if bar < g.upgrades[upgradeType].Level {
				img = wm.upgradeFilled
			} else {
				img = wm.upgradeEmpty
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(x, y)
			screen.DrawImage(img, op)
		}
	}
}

// drawUpgradePlusSigns draws plus signs on all upgrade slots except the coin slot (index 0)
func (wm *WorldMap) drawUpgradePlusSigns(screen *ebiten.Image) {
	plusImg := wm.plusSign
	if plusImg == nil {
		return
	}

	for i := 0; i < upgradeSlotCount; i++ {
		// Skip the coin slot (index 0)
		if i == 0 {
			continue
		}

		x := float64(upgradeSlotStartX)
		y := float64(upgradeSlotStartY + i*upgradeSlotStepY)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(x, y)
		screen.DrawImage(plusImg, op)
	}
}

// getUpgradeInfo returns the name and description for an upgrade type
func (wm *WorldMap) getUpgradeInfo(upgradeType UpgradeType) (string, string) {
	switch upgradeType {
	case UpgradeHealth:
		return "HEALTH", "Increases max health by 1"
	case UpgradeBulletStrength:
		return "DAMAGE", "Increases bullet damage"
	case UpgradeBulletSpeed:
		return "BULLET SPEED", "Makes bullets fly faster"
	case UpgradeBulletCount:
		return "MULTISHOT", "Fire more bullets at once"
	case UpgradeSpeed:
		return "MOVEMENT", "Increases player speed"
	case UpgradeMagnetism:
		return "MAGNET", "Attracts coins from farther"
	case UpgradeLuck:
		return "LUCK", "Better coin spawn rates"
	default:
		return "UNKNOWN", "No description available"
	}
}

// wrapText wraps text to specified line length
func (wm *WorldMap) wrapText(text string, maxLen int) []string {
	if len(text) <= maxLen {
		return []string{text}
	}

	var lines []string
	for len(text) > maxLen {
		space := maxLen
		for space > 0 && text[space] != ' ' {
			space--
		}
		if space == 0 {
			space = maxLen
		}
		lines = append(lines, text[:space])
		text = text[space+1:]
	}
	lines = append(lines, text)

	return lines
}

// loadSystemFont loads a retro-style system font for tooltips
func loadSystemFont() (font.Face, error) {
	// Use basicfont which provides a clean, readable monospace font
	face := basicfont.Face7x13
	return face, nil
}

// drawUpgradeHoverText shows upgrade tooltip when hovering over plus signs
func (wm *WorldMap) drawUpgradeHoverText(screen *ebiten.Image, g *Game) {
	mx, my := ebiten.CursorPosition()
	if upgradeType, ok := wm.plusSignAtScreen(mx, my); ok {
		currentLevel := g.upgrades[upgradeType].Level

		// Get upgrade info
		name, description := wm.getUpgradeInfo(upgradeType)

		// Position tooltip above the plus sign
		tooltipX := float64(upgradeSlotStartX) + 40
		tooltipY := float64(upgradeSlotStartY + (int(upgradeType)+1)*upgradeSlotStepY)

		// Calculate tooltip dimensions
		padding := 6
		lineHeight := 14
		descLines := wm.wrapText(description, 20) // Wrap to 30 chars per line

		// Calculate background size
		bgWidth := 180
		bgHeight := padding*2 + lineHeight + len(descLines)*lineHeight + lineHeight // name + desc lines + cost

		// Draw semi-transparent black background
		bgColor := color.RGBA{0, 0, 0, 200} // Semi-transparent black
		vector.FillRect(screen, float32(tooltipX-3), float32(tooltipY-3), float32(bgWidth), float32(bgHeight), bgColor, false)

		// Use system font if available, otherwise fallback to bitmap font
		if wm.tooltipFont != nil {
			// Draw with system font
			textColor := color.RGBA{255, 255, 255, 255}

			// Draw upgrade name
			text.Draw(screen, name, wm.tooltipFont, int(tooltipX+float64(padding)), int(tooltipY+float64(padding+lineHeight-5)), textColor)

			// Draw description
			for i, line := range descLines {
				text.Draw(screen, line, wm.tooltipFont, int(tooltipX+float64(padding)), int(tooltipY+float64(padding+lineHeight-5+(i+1)*lineHeight)), textColor)
			}

			// Draw cost or MAX level
			costY := tooltipY + float64(padding+lineHeight-5+(len(descLines)+1)*lineHeight)
			maxLevel := 7
			if upgradeType == UpgradeBulletCount {
				maxLevel = 4
			}
			if currentLevel >= maxLevel {
				text.Draw(screen, "MAX LEVEL", wm.tooltipFont, int(tooltipX+float64(padding)), int(costY), textColor)
			} else {
				cost := wm.getUpgradeCost(upgradeType, currentLevel)
				costText := fmt.Sprintf("Cost: %d coins", cost)
				text.Draw(screen, costText, wm.tooltipFont, int(tooltipX+float64(padding)), int(costY), textColor)
			}
		} else {
			// Fallback to bitmap font
			if wm.font == nil {
				return
			}
			// Draw upgrade name (larger font with shadow)
			wm.font.DrawTextWithShadow(screen, name, tooltipX+float64(padding), tooltipY+float64(padding), 1.2)

			// Draw description (medium font with shadow)
			for i, line := range descLines {
				wm.font.DrawTextWithShadow(screen, line, tooltipX+float64(padding), tooltipY+float64(padding+lineHeight+i*lineHeight), 1.0)
			}

			// Draw cost or MAX level (medium font with shadow)
			costY := tooltipY + float64(padding+lineHeight+len(descLines)*lineHeight)
			maxLevel := 7
			if upgradeType == UpgradeBulletCount {
				maxLevel = 4
			}
			if currentLevel >= maxLevel {
				wm.font.DrawTextWithShadow(screen, "MAX LEVEL", tooltipX+float64(padding), costY, 1.0)
			} else {
				cost := wm.getUpgradeCost(upgradeType, currentLevel)
				costText := fmt.Sprintf("Cost: %d coins", cost)
				wm.font.DrawTextWithShadow(screen, costText, tooltipX+float64(padding), costY, 1.0)
			}
		}
	}
}

func (wm *WorldMap) drawTestingHint(screen *ebiten.Image) {
	// Show testing hint when shift is pressed
	if ebiten.IsKeyPressed(ebiten.KeyShift) || ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShiftRight) {
		if wm.font != nil {
			wm.font.DrawTextWithShadow(screen, "Shift+Click: Test any level", 10, 580, 1.0)
		}
	}
}
