package charts

import (
	"strconv"
	"strings"
	"testing"

	"github.com/go-analyze/bulk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDojiPattern(t *testing.T) {
	t.Parallel()

	// Valid doji: open ≈ close
	doji := OHLCData{Open: 100, High: 105, Low: 95, Close: 100.1}
	data := []OHLCData{doji}
	for _, tt := range []struct {
		name      string
		threshold float64
		expected  bool
	}{
		{"low", 0.009, false},
		{"default", 0.01, true},
		{"high", 0.011, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, detectDojiAt(data, 0, CandlestickPatternConfig{DojiThreshold: tt.threshold}))
		})
	}

	// Invalid: body too large
	notDoji := OHLCData{Open: 100, High: 105, Low: 95, Close: 103}
	data = []OHLCData{notDoji}
	assert.False(t, detectDojiAt(data, 0, CandlestickPatternConfig{DojiThreshold: 0.01}))

	// Invalid: invalid OHLC
	invalidOHLC := OHLCData{Open: 100, High: 95, Low: 105, Close: 98}
	data = []OHLCData{invalidOHLC}
	assert.False(t, detectDojiAt(data, 0, CandlestickPatternConfig{DojiThreshold: 0.01}))
}

func TestHammerPattern(t *testing.T) {
	t.Parallel()

	// Valid hammer: long lower shadow, small body at top
	hammer := OHLCData{Open: 105, High: 107, Low: 95, Close: 106}
	data := []OHLCData{hammer}
	for _, tt := range []struct {
		name     string
		ratio    float64
		expected bool
	}{
		{"low", 1.0, true},
		{"default", 2.0, true},
		{"high", 11.1, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, detectHammerAt(data, 0, CandlestickPatternConfig{ShadowRatio: tt.ratio}))
		})
	}

	// Invalid: short lower shadow
	notHammer := OHLCData{Open: 105, High: 107, Low: 104, Close: 106}
	data = []OHLCData{notHammer}
	assert.False(t, detectHammerAt(data, 0, CandlestickPatternConfig{ShadowRatio: 2.0}))

	// Invalid: long upper shadow
	notHammer2 := OHLCData{Open: 95, High: 107, Low: 94, Close: 96}
	data = []OHLCData{notHammer2}
	assert.False(t, detectHammerAt(data, 0, CandlestickPatternConfig{ShadowRatio: 2.0}))
}

func TestInvertedHammerPattern(t *testing.T) {
	t.Parallel()

	// Valid inverted hammer: long upper shadow, small body at bottom
	invertedHammer := OHLCData{Open: 95, High: 107, Low: 94, Close: 96}
	for _, tt := range []struct {
		name     string
		ratio    float64
		expected bool
	}{
		{"low", 1.0, true},
		{"default", 2.0, true},
		{"high", 11.1, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, detectInvertedHammerAt([]OHLCData{invertedHammer}, 0, CandlestickPatternConfig{ShadowRatio: tt.ratio}))
		})
	}

	// Invalid: short upper shadow
	notInvertedHammer := OHLCData{Open: 95, High: 97, Low: 94, Close: 96}
	assert.False(t, detectInvertedHammerAt([]OHLCData{notInvertedHammer}, 0, CandlestickPatternConfig{ShadowRatio: 2.0}))
}

func TestEngulfingPattern(t *testing.T) {
	t.Parallel()

	prevBearish := OHLCData{Open: 110, High: 112, Low: 105, Close: 106}
	currentBullish := OHLCData{Open: 104, High: 115, Low: 103, Close: 114}
	for _, tt := range []struct {
		name     string
		size     float64
		expected bool
	}{
		{"low", 0.5, true},
		{"default", 0.8, true},
		{"high", 2.6, false},
	} {
		t.Run("bullish_"+tt.name, func(t *testing.T) {
			detected := detectBullishEngulfingAt([]OHLCData{prevBearish, currentBullish}, 1, CandlestickPatternConfig{EngulfingMinSize: tt.size})
			assert.Equal(t, tt.expected, detected)
		})
	}
	assert.False(t, detectBearishEngulfingAt([]OHLCData{prevBearish, currentBullish}, 1, CandlestickPatternConfig{EngulfingMinSize: 0.8}))

	prevBullish := OHLCData{Open: 106, High: 112, Low: 105, Close: 110}
	currentBearish := OHLCData{Open: 114, High: 115, Low: 103, Close: 104}

	for _, tt := range []struct {
		name     string
		size     float64
		expected bool
	}{
		{"low", 0.5, true},
		{"default", 0.8, true},
		{"high", 2.6, false},
	} {
		t.Run("bearish_"+tt.name, func(t *testing.T) {
			detected := detectBearishEngulfingAt([]OHLCData{prevBullish, currentBearish}, 1, CandlestickPatternConfig{EngulfingMinSize: tt.size})
			assert.Equal(t, tt.expected, detected)
		})
	}
	assert.False(t, detectBullishEngulfingAt([]OHLCData{prevBullish, currentBearish}, 1, CandlestickPatternConfig{EngulfingMinSize: 0.8}))

	// Test non-engulfing
	nonEngulfing := OHLCData{Open: 107, High: 109, Low: 106, Close: 108}
	assert.False(t, detectBullishEngulfingAt([]OHLCData{prevBullish, nonEngulfing}, 1, CandlestickPatternConfig{EngulfingMinSize: 0.8}))
	assert.False(t, detectBearishEngulfingAt([]OHLCData{prevBullish, nonEngulfing}, 1, CandlestickPatternConfig{EngulfingMinSize: 0.8}))
}

func TestShootingStarPattern(t *testing.T) {
	t.Parallel()

	// Valid shooting star: small body at bottom, long upper shadow
	shootingStar := OHLCData{Open: 106, High: 125, Low: 105, Close: 107}
	for _, tt := range []struct {
		name     string
		ratio    float64
		expected bool
	}{
		{"low", 1.0, true},
		{"default", 2.0, true},
		{"high", 18.1, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, detectShootingStarAt([]OHLCData{shootingStar}, 0, CandlestickPatternConfig{ShadowRatio: tt.ratio}))
		})
	}

	// Invalid: body not near bottom
	notShootingStar := OHLCData{Open: 115, High: 125, Low: 105, Close: 117}
	assert.False(t, detectShootingStarAt([]OHLCData{notShootingStar}, 0, CandlestickPatternConfig{ShadowRatio: 2.0}))

	// Invalid: upper shadow too short
	shortShadow := OHLCData{Open: 106, High: 110, Low: 105, Close: 107}
	assert.False(t, detectShootingStarAt([]OHLCData{shortShadow}, 0, CandlestickPatternConfig{ShadowRatio: 2.0}))
}

func TestGravestoneDojiPattern(t *testing.T) {
	t.Parallel()

	gravestoneDoji := OHLCData{Open: 108, High: 120, Low: 107, Close: 108.1}
	for _, tt := range []struct {
		name      string
		threshold float64
		shadow    float64
		expected  bool
	}{
		{"low_threshold", 0.004, 2.0, false},
		{"default", 0.01, 2.0, true},
		{"high_shadow", 0.01, 200, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			opt := CandlestickPatternConfig{DojiThreshold: tt.threshold, ShadowRatio: tt.shadow}
			assert.Equal(t, tt.expected, detectGravestoneDojiAt([]OHLCData{gravestoneDoji}, 0, opt))
		})
	}

	// Invalid: not a doji (body too large)
	notDoji := OHLCData{Open: 108, High: 120, Low: 107, Close: 115}
	assert.False(t, detectGravestoneDojiAt([]OHLCData{notDoji}, 0, CandlestickPatternConfig{DojiThreshold: 0.01, ShadowRatio: 2.0}))

	// Invalid: doji but no long upper shadow
	dojiNoShadow := OHLCData{Open: 108, High: 109, Low: 107, Close: 108.1}
	assert.False(t, detectGravestoneDojiAt([]OHLCData{dojiNoShadow}, 0, CandlestickPatternConfig{DojiThreshold: 0.01, ShadowRatio: 2.0}))
}

func TestDragonflyDojiPattern(t *testing.T) {
	t.Parallel()

	dragonflyDoji := OHLCData{Open: 109, High: 110, Low: 90, Close: 108.9}
	for _, tt := range []struct {
		name      string
		threshold float64
		shadow    float64
		expected  bool
	}{
		{"low_threshold", 0.004, 2.0, false},
		{"default", 0.01, 2.0, true},
		{"high_shadow", 0.01, 200, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			opt := CandlestickPatternConfig{DojiThreshold: tt.threshold, ShadowRatio: tt.shadow}
			assert.Equal(t, tt.expected, detectDragonflyDojiAt([]OHLCData{dragonflyDoji}, 0, opt))
		})
	}

	// Invalid: not a doji
	notDoji := OHLCData{Open: 109, High: 110, Low: 90, Close: 102}
	assert.False(t, detectDragonflyDojiAt([]OHLCData{notDoji}, 0, CandlestickPatternConfig{DojiThreshold: 0.01, ShadowRatio: 2.0}))

	// Invalid: doji but no long lower shadow
	dojiNoShadow := OHLCData{Open: 109, High: 110, Low: 108, Close: 108.9}
	assert.False(t, detectDragonflyDojiAt([]OHLCData{dojiNoShadow}, 0, CandlestickPatternConfig{DojiThreshold: 0.01, ShadowRatio: 2.0}))
}

func TestMorningStarPattern(t *testing.T) {
	t.Parallel()

	opt := CandlestickPatternConfig{}

	// Valid morning star pattern
	first := OHLCData{Open: 120, High: 125, Low: 105, Close: 108}  // Large bearish
	second := OHLCData{Open: 102, High: 104, Low: 100, Close: 103} // Small body, gap down
	third := OHLCData{Open: 108, High: 125, Low: 106, Close: 122}  // Large bullish, gap up

	assert.True(t, detectMorningStarAt([]OHLCData{first, second, third}, 2, opt))

	// Invalid: first candle not bearish
	invalidFirst := OHLCData{Open: 108, High: 125, Low: 105, Close: 120} // Bullish
	assert.False(t, detectMorningStarAt([]OHLCData{invalidFirst, second, third}, 2, opt))

	// Invalid: no gap down between first and second
	noGapSecond := OHLCData{Open: 109, High: 111, Low: 107, Close: 110} // No gap
	assert.False(t, detectMorningStarAt([]OHLCData{first, noGapSecond, third}, 2, opt))

	// Invalid: third candle not bullish
	invalidThird := OHLCData{Open: 108, High: 110, Low: 105, Close: 107} // Bearish
	assert.False(t, detectMorningStarAt([]OHLCData{first, second, invalidThird}, 2, opt))
}

func TestEveningStarPattern(t *testing.T) {
	t.Parallel()

	opt := CandlestickPatternConfig{}

	// Valid evening star pattern
	first := OHLCData{Open: 122, High: 140, Low: 120, Close: 138}  // Large bullish
	second := OHLCData{Open: 142, High: 144, Low: 140, Close: 143} // Small body, gap up
	third := OHLCData{Open: 138, High: 140, Low: 115, Close: 118}  // Large bearish, gap down

	assert.True(t, detectEveningStarAt([]OHLCData{first, second, third}, 2, opt))

	// Invalid: first candle not bullish
	invalidFirst := OHLCData{Open: 138, High: 140, Low: 120, Close: 122} // Bearish
	assert.False(t, detectEveningStarAt([]OHLCData{invalidFirst, second, third}, 2, opt))

	// Invalid: no gap up between first and second
	noGapSecond := OHLCData{Open: 136, High: 140, Low: 134, Close: 139} // No gap
	assert.False(t, detectEveningStarAt([]OHLCData{first, noGapSecond, third}, 2, opt))

	// Invalid: third candle not bearish
	invalidThird := OHLCData{Open: 138, High: 145, Low: 135, Close: 142} // Bullish
	assert.False(t, detectEveningStarAt([]OHLCData{first, second, invalidThird}, 2, opt))
}

func newCandlestickWithPatterns(data []OHLCData, options ...CandlestickPatternConfig) CandlestickSeries {
	// Start with defaults and override with provided options
	config := &CandlestickPatternConfig{
		PreferPatternLabels: true,
		EnabledPatterns:     (&CandlestickPatternConfig{}).WithPatternsAll().EnabledPatterns,
		DojiThreshold:       0.001,
		ShadowTolerance:     0.01,
		ShadowRatio:         2.0,
		EngulfingMinSize:    0.8,
	}
	if len(options) > 0 {
		// Merge provided options with defaults
		opt := options[0]
		if opt.DojiThreshold > 0 {
			config.DojiThreshold = opt.DojiThreshold
		}
		if opt.ShadowRatio > 0 {
			config.ShadowRatio = opt.ShadowRatio
		}
		if opt.EngulfingMinSize > 0 {
			config.EngulfingMinSize = opt.EngulfingMinSize
		}
		if opt.ShadowTolerance > 0 {
			config.ShadowTolerance = opt.ShadowTolerance
		}
	}

	return CandlestickSeries{
		Data:          data,
		PatternConfig: config,
	}
}

func makePatternChartOption(data []OHLCData, config CandlestickPatternConfig) CandlestickChartOption {
	series := newCandlestickWithPatterns(data, config)
	labels := make([]string, len(data))
	for i := range labels {
		labels[i] = strconv.Itoa(i + 1)
	}
	return CandlestickChartOption{
		XAxis:      XAxisOption{Labels: labels},
		YAxis:      make([]YAxisOption, 1),
		SeriesList: CandlestickSeriesList{series},
		Padding:    NewBoxEqual(10),
	}
}

func TestMarubozuPattern(t *testing.T) {
	t.Parallel()

	// Bullish Marubozu - no shadows
	bullishMarubozu := OHLCData{Open: 100, High: 120, Low: 100, Close: 120}
	for _, tt := range []struct {
		tol      float64
		expected bool
	}{
		{0.005, true},
		{0.01, true},
		{0.02, true},
	} {
		t.Run("bullish_tol_"+strconv.FormatFloat(tt.tol, 'f', 3, 64), func(t *testing.T) {
			detected := detectBullishMarubozuAt([]OHLCData{bullishMarubozu}, 0, CandlestickPatternConfig{ShadowTolerance: tt.tol})
			assert.Equal(t, tt.expected, detected)
		})
	}
	assert.False(t, detectBearishMarubozuAt([]OHLCData{bullishMarubozu}, 0, CandlestickPatternConfig{ShadowTolerance: 0.01}))

	// Bearish Marubozu - no shadows
	bearishMarubozu := OHLCData{Open: 120, High: 120, Low: 100, Close: 100}
	for _, tt := range []struct {
		tol      float64
		expected bool
	}{
		{0.005, true},
		{0.01, true},
		{0.02, true},
	} {
		t.Run("bearish_tol_"+strconv.FormatFloat(tt.tol, 'f', 3, 64), func(t *testing.T) {
			detected := detectBearishMarubozuAt([]OHLCData{bearishMarubozu}, 0, CandlestickPatternConfig{ShadowTolerance: tt.tol})
			assert.Equal(t, tt.expected, detected)
		})
	}
	assert.False(t, detectBullishMarubozuAt([]OHLCData{bearishMarubozu}, 0, CandlestickPatternConfig{ShadowTolerance: 0.01}))

	// Not a marubozu - has significant shadows
	notMarubozu := OHLCData{Open: 105, High: 125, Low: 95, Close: 115}
	assert.False(t, detectBullishMarubozuAt([]OHLCData{notMarubozu}, 0, CandlestickPatternConfig{ShadowTolerance: 0.01}))
	assert.False(t, detectBearishMarubozuAt([]OHLCData{notMarubozu}, 0, CandlestickPatternConfig{ShadowTolerance: 0.01}))
	assert.True(t, detectBullishMarubozuAt([]OHLCData{notMarubozu}, 0, CandlestickPatternConfig{ShadowTolerance: 0.7}))
}

func TestPiercingLinePattern(t *testing.T) {
	t.Parallel()

	// Classic piercing line - bearish then bullish with gap down and close above midpoint
	prev := OHLCData{Open: 120, High: 120, Low: 110, Close: 110}    // Bearish
	current := OHLCData{Open: 108, High: 118, Low: 108, Close: 116} // Bullish, opens below prev low, closes above midpoint (115)
	detected := detectPiercingLineAt([]OHLCData{prev, current}, 1, CandlestickPatternConfig{})
	assert.True(t, detected)

	// Not piercing line - current closes below midpoint
	current = OHLCData{Open: 108, High: 114, Low: 108, Close: 112}
	detected = detectPiercingLineAt([]OHLCData{prev, current}, 1, CandlestickPatternConfig{})
	assert.False(t, detected)
}

func TestDarkCloudCoverPattern(t *testing.T) {
	t.Parallel()

	// Classic dark cloud cover - bullish then bearish with gap up and close below midpoint
	prev := OHLCData{Open: 110, High: 120, Low: 110, Close: 120}    // Bullish
	current := OHLCData{Open: 122, High: 122, Low: 112, Close: 114} // Bearish, opens above prev high, closes below midpoint (115)
	detected := detectDarkCloudCoverAt([]OHLCData{prev, current}, 1, CandlestickPatternConfig{})
	assert.True(t, detected)

	// Not dark cloud cover - current closes above midpoint
	current = OHLCData{Open: 122, High: 122, Low: 118, Close: 118}
	detected = detectDarkCloudCoverAt([]OHLCData{prev, current}, 1, CandlestickPatternConfig{})
	assert.False(t, detected)
}

func TestPatternValidation(t *testing.T) {
	t.Parallel()

	// Test with invalid OHLC data
	invalidOHLC := OHLCData{Open: 100, High: 95, Low: 105, Close: 98} // High < Low

	assert.False(t, detectDojiAt([]OHLCData{invalidOHLC}, 0, CandlestickPatternConfig{DojiThreshold: 0.01}))
	assert.False(t, detectHammerAt([]OHLCData{invalidOHLC}, 0, CandlestickPatternConfig{ShadowRatio: 2.0}))
	assert.False(t, detectShootingStarAt([]OHLCData{invalidOHLC}, 0, CandlestickPatternConfig{ShadowRatio: 2.0}))

	// Test three-candle patterns with invalid data
	validOHLC := OHLCData{Open: 100, High: 110, Low: 95, Close: 105}
	opt := CandlestickPatternConfig{}

	assert.False(t, detectMorningStarAt([]OHLCData{invalidOHLC, validOHLC, validOHLC}, 2, opt))
	assert.False(t, detectMorningStarAt([]OHLCData{validOHLC, invalidOHLC, validOHLC}, 2, opt))
	assert.False(t, detectMorningStarAt([]OHLCData{validOHLC, validOHLC, invalidOHLC}, 2, opt))

	assert.False(t, detectEveningStarAt([]OHLCData{invalidOHLC, validOHLC, validOHLC}, 2, opt))
	assert.False(t, detectEveningStarAt([]OHLCData{validOHLC, invalidOHLC, validOHLC}, 2, opt))
	assert.False(t, detectEveningStarAt([]OHLCData{validOHLC, validOHLC, invalidOHLC}, 2, opt))
}

func TestPatternScanningComprehensive(t *testing.T) {
	t.Parallel()

	data := []OHLCData{
		// Index 0: Normal candle
		{Open: 100, High: 110, Low: 95, Close: 105},
		// Index 1: Doji
		{Open: 105, High: 108, Low: 102, Close: 105.05},
		// Index 2: Hammer
		{Open: 108, High: 109, Low: 98, Close: 107},
		// Index 3: Shooting Star
		{Open: 106, High: 125, Low: 105, Close: 107},
		// Index 4: Gravestone Doji
		{Open: 108, High: 120, Low: 107, Close: 108.1},
		// Index 5: Dragonfly Doji
		{Open: 109, High: 110, Low: 90, Close: 108.9},
		// Index 6-8: Morning Star sequence
		{Open: 120, High: 125, Low: 105, Close: 108}, // 6: Large bearish
		{Open: 102, High: 104, Low: 100, Close: 103}, // 7: Small body, gap down
		{Open: 108, High: 125, Low: 106, Close: 122}, // 8: Large bullish, gap up
		// Index 9-11: Evening Star sequence
		{Open: 122, High: 140, Low: 120, Close: 138}, // 9: Large bullish
		{Open: 142, High: 144, Low: 140, Close: 143}, // 10: Small body, gap up
		{Open: 138, High: 140, Low: 115, Close: 118}, // 11: Large bearish, gap down
		// Index 12: Bullish Marubozu (no shadows)
		{Open: 120, High: 135, Low: 120, Close: 135},
		// Index 13: Bearish Marubozu (no shadows)
		{Open: 135, High: 135, Low: 115, Close: 115},
		// Index 14: Spinning Top (small body, long shadows)
		{Open: 118, High: 125, Low: 110, Close: 119},
		// Index 15: Setup for Piercing Line - bearish candle
		{Open: 120, High: 121, Low: 115, Close: 115},
		// Index 16: Piercing Line - bullish candle opening below prev low, closing above midpoint
		{Open: 112, High: 119, Low: 112, Close: 118}, // Opens below 115, closes above midpoint (117.5)
		// Index 17: Setup for Dark Cloud Cover - bullish candle
		{Open: 118, High: 125, Low: 118, Close: 125},
		// Index 18: Dark Cloud Cover - bearish candle opening above prev high, closing below midpoint
		{Open: 127, High: 127, Low: 120, Close: 121}, // Opens above 125, closes below midpoint (121.5)
		// Index 19: Setup for Tweezer Bottom - bearish with low at 100
		{Open: 125, High: 126, Low: 100, Close: 102},
		// Index 20: Tweezer Bottom - bullish with same low at 100
		{Open: 102, High: 108, Low: 100, Close: 107},
		// Index 21-23: Three White Soldiers sequence
		{Open: 110, High: 115, Low: 109, Close: 114}, // 21: First soldier
		{Open: 113, High: 118, Low: 112, Close: 117}, // 22: Second soldier
		{Open: 116, High: 121, Low: 115, Close: 120}, // 23: Third soldier
		// Index 24-26: Three Black Crows sequence
		{Open: 120, High: 121, Low: 115, Close: 116}, // 24: First crow
		{Open: 117, High: 118, Low: 112, Close: 113}, // 25: Second crow
		{Open: 114, High: 115, Low: 108, Close: 109}, // 26: Third crow
	}

	opt := (&CandlestickPatternConfig{}).WithPatternsAll()
	opt.DojiThreshold = 0.01
	opt.ShadowRatio = 2.0
	opt.EngulfingMinSize = 0.8
	indexPatterns := scanForCandlestickPatterns(data, *opt)

	// Verify specific patterns were detected
	patternsByIndex := make(map[int][]string)
	uniquePatterns := make(map[string]bool)
	for index, patterns := range indexPatterns {
		for _, pattern := range patterns {
			patternsByIndex[index] = append(patternsByIndex[index], pattern.PatternType)
			uniquePatterns[pattern.PatternType] = true
		}
	}

	// Check expected patterns
	assert.Len(t, uniquePatterns, 13)
	assert.Contains(t, patternsByIndex[1], "doji")
	assert.Contains(t, patternsByIndex[2], "hammer")
	assert.Contains(t, patternsByIndex[3], "shooting_star")
	assert.Contains(t, patternsByIndex[4], "gravestone_doji")
	assert.Contains(t, patternsByIndex[5], "dragonfly_doji")
	assert.Contains(t, patternsByIndex[8], "morning_star")
	assert.Contains(t, patternsByIndex[11], "evening_star")
	assert.Contains(t, patternsByIndex[12], "marubozu_bull")
	assert.Contains(t, patternsByIndex[13], "marubozu_bear")
	assert.Contains(t, patternsByIndex[16], "piercing_line")
	assert.Contains(t, patternsByIndex[18], "dark_cloud_cover")
}

func TestCandlestickPatternSets(t *testing.T) {
	t.Parallel()

	t.Run("all", func(t *testing.T) {
		config := (&CandlestickPatternConfig{}).WithPatternsAll()

		assert.Contains(t, config.EnabledPatterns, "doji")
		assert.Contains(t, config.EnabledPatterns, "hammer")
		assert.Len(t, config.EnabledPatterns, 14)
	})

	t.Run("core", func(t *testing.T) {
		config := (&CandlestickPatternConfig{}).WithPatternsCore()

		assert.Contains(t, config.EnabledPatterns, "engulfing_bull")
		assert.Contains(t, config.EnabledPatterns, "hammer")
		assert.Len(t, config.EnabledPatterns, 6)
	})

	t.Run("bullish", func(t *testing.T) {
		config := (&CandlestickPatternConfig{}).WithPatternsBullish()

		assert.Contains(t, config.EnabledPatterns, "hammer")
		assert.NotContains(t, config.EnabledPatterns, "shooting_star")
		assert.Len(t, config.EnabledPatterns, 7)
	})

	t.Run("bearish", func(t *testing.T) {
		config := (&CandlestickPatternConfig{}).WithPatternsBearish()

		assert.Contains(t, config.EnabledPatterns, "shooting_star")
		assert.NotContains(t, config.EnabledPatterns, "hammer")
		assert.Len(t, config.EnabledPatterns, 6)
	})

	t.Run("reversal", func(t *testing.T) {
		config := (&CandlestickPatternConfig{}).WithPatternsReversal()

		assert.Contains(t, config.EnabledPatterns, "hammer")
		assert.NotContains(t, config.EnabledPatterns, "marubozu_bull")
		assert.Len(t, config.EnabledPatterns, 10)
	})

	t.Run("trend", func(t *testing.T) {
		config := (&CandlestickPatternConfig{}).WithPatternsTrend()

		assert.Contains(t, config.EnabledPatterns, "marubozu_bull")
		assert.NotContains(t, config.EnabledPatterns, "hammer")
		assert.Len(t, config.EnabledPatterns, 2)
	})
}

func TestPatternFormatterCustom(t *testing.T) {
	t.Parallel()

	// Data with a clear Doji at index 1
	data := []OHLCData{
		{Open: 100, High: 110, Low: 95, Close: 105},
		{Open: 105, High: 107, Low: 103, Close: 105}, // Doji
		{Open: 105, High: 112, Low: 98, Close: 108},
	}

	// Custom formatter that prefixes with PF: and joins all detected pattern types
	customFormatter := func(patterns []PatternDetectionResult, seriesName string, value float64) (string, *LabelStyle) {
		if len(patterns) == 0 {
			return "", nil
		}
		names := make([]string, len(patterns))
		for i, p := range patterns {
			names[i] = p.PatternType
		}
		return "PF:" + strings.Join(names, "+"), &LabelStyle{FontStyle: FontStyle{FontColor: ColorGray}}
	}

	mkOpt := func(cfg CandlestickPatternConfig, userLabel bool) CandlestickChartOption {
		series := CandlestickSeries{
			Data: data,
			Label: SeriesLabel{
				Show: Ptr(userLabel),
				LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
					return "UserLabel", nil
				},
			},
			PatternConfig: &cfg,
		}
		labels := []string{"1", "2", "3"}
		return CandlestickChartOption{
			XAxis:      XAxisOption{Labels: labels},
			YAxis:      make([]YAxisOption, 1),
			SeriesList: CandlestickSeriesList{series},
			Padding:    NewBoxEqual(10),
		}
	}

	t.Run("pattern_priority_mode", func(t *testing.T) {
		cfg := CandlestickPatternConfig{
			PreferPatternLabels: true,
			EnabledPatterns:     []string{"doji"},
			DojiThreshold:       0.001,
			PatternFormatter:    customFormatter,
		}
		opt := mkOpt(cfg, true)
		p := NewPainter(PainterOptions{OutputFormat: ChartOutputSVG, Width: 800, Height: 600})
		require.NoError(t, p.CandlestickChart(opt))
		svg, err := p.Bytes()
		require.NoError(t, err)
		s := string(svg)
		// With pattern priority, pattern label should be shown at index 1 where Doji is detected
		assert.Contains(t, s, "PF:"+"doji")
		// User labels should still appear at indices 0 and 2 where no patterns are detected
		assert.Contains(t, s, "UserLabel")
	})

	t.Run("user_priority_mode", func(t *testing.T) {
		cfg := CandlestickPatternConfig{
			PreferPatternLabels: false,
			EnabledPatterns:     []string{"doji"},
			DojiThreshold:       0.001,
			PatternFormatter:    customFormatter,
		}
		opt := mkOpt(cfg, true)
		p := NewPainter(PainterOptions{OutputFormat: ChartOutputSVG, Width: 800, Height: 600})
		require.NoError(t, p.CandlestickChart(opt))
		svg, err := p.Bytes()
		require.NoError(t, err)
		s := string(svg)
		// With user priority, user labels take precedence everywhere they're provided
		assert.Contains(t, s, "UserLabel")
		assert.NotContains(t, s, "PF:")
	})

	t.Run("no_user_labels_shows_patterns", func(t *testing.T) {
		cfg := CandlestickPatternConfig{
			PreferPatternLabels: false,
			EnabledPatterns:     []string{"doji"},
			DojiThreshold:       0.001,
			PatternFormatter:    customFormatter,
		}
		opt := mkOpt(cfg, false) // user label disabled
		p := NewPainter(PainterOptions{OutputFormat: ChartOutputSVG, Width: 800, Height: 600})
		require.NoError(t, p.CandlestickChart(opt))
		svg, err := p.Bytes()
		require.NoError(t, err)
		s := string(svg)
		// When no user labels are provided, patterns should be shown
		assert.Contains(t, s, "PF:"+"doji")
		assert.NotContains(t, s, "UserLabel")
	})
}

func TestCandlestickChartPatterns(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		optGen func() CandlestickChartOption
		pngCRC uint32
	}{
		{
			name: "doji",
			optGen: func() CandlestickChartOption {
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105},  // Normal candle
					{Open: 105, High: 107, Low: 103, Close: 105}, // Pure Doji pattern - minimal body and minimal shadows
					{Open: 105, High: 112, Low: 98, Close: 108},  // Normal candle
				}
				return makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold: 0.01,
					ShadowRatio:   2.0,
				})
			},
			pngCRC: 0xd14dc283,
		},
		{
			name: "hammer",
			optGen: func() CandlestickChartOption {
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105},  // Normal candle
					{Open: 108, High: 109, Low: 98, Close: 107},  // Hammer pattern
					{Open: 107, High: 112, Low: 102, Close: 110}, // Normal candle
				}
				return makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold: 0.01,
					ShadowRatio:   2.0,
				})
			},
			pngCRC: 0x48755ffa,
		},
		{
			name: "inverted_hammer",
			optGen: func() CandlestickChartOption {
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105}, // Normal candle
					{Open: 95, High: 107, Low: 94, Close: 96},   // Inverted hammer
					{Open: 96, High: 102, Low: 91, Close: 98},   // Normal candle
				}
				return makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold: 0.01,
					ShadowRatio:   2.0,
				})
			},
			pngCRC: 0x2db8b33b,
		},
		{
			name: "shooting_star",
			optGen: func() CandlestickChartOption {
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105},  // Normal candle
					{Open: 107, High: 125, Low: 106, Close: 108}, // Shooting star - small body at bottom, long upper shadow
					{Open: 107, High: 112, Low: 102, Close: 109}, // Normal candle
				}
				return makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold: 0.01,
					ShadowRatio:   2.0,
				})
			},
			pngCRC: 0xf9c1a5ec,
		},
		{
			name: "gravestone_doji",
			optGen: func() CandlestickChartOption {
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105},  // Normal candle
					{Open: 108, High: 125, Low: 108, Close: 108}, // Gravestone doji - minimal body at bottom, long upper shadow only
					{Open: 108, High: 115, Low: 103, Close: 110}, // Normal candle
				}
				return makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold: 0.01,
					ShadowRatio:   2.0,
				})
			},
			pngCRC: 0x8159f53b,
		},
		{
			name: "dragonfly_doji",
			optGen: func() CandlestickChartOption {
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105},  // Normal candle
					{Open: 109, High: 110, Low: 90, Close: 109},  // Dragonfly doji
					{Open: 109, High: 115, Low: 104, Close: 112}, // Normal candle
				}
				return makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold: 0.01,
					ShadowRatio:   2.0,
				})
			},
			pngCRC: 0xa6e7aaef,
		},
		{
			name: "bullish_marubozu",
			optGen: func() CandlestickChartOption {
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105},  // Normal candle
					{Open: 100, High: 120, Low: 100, Close: 120}, // Bullish marubozu
					{Open: 120, High: 125, Low: 115, Close: 122}, // Normal candle
				}
				return makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold: 0.01,
					ShadowRatio:   2.0,
				})
			},
			pngCRC: 0xf59a696f,
		},
		{
			name: "bearish_marubozu",
			optGen: func() CandlestickChartOption {
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105},  // Normal candle
					{Open: 120, High: 120, Low: 100, Close: 100}, // Bearish marubozu
					{Open: 100, High: 105, Low: 95, Close: 102},  // Normal candle
				}
				return makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold: 0.01,
					ShadowRatio:   2.0,
				})
			},
			pngCRC: 0x78ba7419,
		},
		{
			name: "bullish_engulfing",
			optGen: func() CandlestickChartOption {
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105},  // Normal candle
					{Open: 110, High: 112, Low: 105, Close: 106}, // Small bearish candle
					{Open: 104, High: 115, Low: 103, Close: 114}, // Bullish engulfing
					{Open: 114, High: 120, Low: 112, Close: 118}, // Normal candle
				}
				return makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold:    0.01,
					ShadowRatio:      2.0,
					EngulfingMinSize: 0.8,
				})
			},
			pngCRC: 0xd40689c8,
		},
		{
			name: "bearish_engulfing",
			optGen: func() CandlestickChartOption {
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105},  // Normal candle
					{Open: 106, High: 112, Low: 105, Close: 110}, // Small bullish candle
					{Open: 114, High: 115, Low: 103, Close: 104}, // Bearish engulfing
					{Open: 104, High: 108, Low: 100, Close: 102}, // Normal candle
				}
				return makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold:    0.01,
					ShadowRatio:      2.0,
					EngulfingMinSize: 0.8,
				})
			},
			pngCRC: 0xa1cca1bc,
		},
		{
			name: "morning_star",
			optGen: func() CandlestickChartOption {
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105},  // Normal candle
					{Open: 120, High: 125, Low: 105, Close: 108}, // Large bearish
					{Open: 102, High: 104, Low: 100, Close: 103}, // Small body, gap down - overlaps are expected
					{Open: 108, High: 125, Low: 106, Close: 122}, // Large bullish, gap up
					{Open: 122, High: 128, Low: 120, Close: 125}, // Normal candle
				}
				return makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold: 0.01,
					ShadowRatio:   2.0,
				})
			},
			pngCRC: 0xe848b023,
		},
		{
			name: "evening_star",
			optGen: func() CandlestickChartOption {
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105},  // Normal candle
					{Open: 122, High: 140, Low: 120, Close: 138}, // Large bullish
					{Open: 142, High: 144, Low: 140, Close: 143}, // Small body, gap up - overlaps are expected
					{Open: 138, High: 140, Low: 115, Close: 118}, // Large bearish, gap down
					{Open: 118, High: 122, Low: 115, Close: 120}, // Normal candle - harami overlap is expected
				}
				return makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold: 0.01,
					ShadowRatio:   2.0,
				})
			},
			pngCRC: 0x1683fe48,
		},
		{
			name: "piercing_line",
			optGen: func() CandlestickChartOption {
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105},  // Normal candle
					{Open: 120, High: 121, Low: 115, Close: 115}, // Bearish candle
					{Open: 112, High: 119, Low: 111, Close: 118}, // Piercing line (opens below prev low, closes above midpoint)
					{Open: 118, High: 125, Low: 116, Close: 122}, // Normal candle
				}
				return makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold: 0.01,
					ShadowRatio:   2.0,
				})
			},
			pngCRC: 0x957f0bb4,
		},
		{
			name: "dark_cloud_cover",
			optGen: func() CandlestickChartOption {
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105},  // Normal candle
					{Open: 118, High: 125, Low: 117, Close: 125}, // Bullish candle
					{Open: 127, High: 128, Low: 120, Close: 121}, // Dark cloud cover (opens above prev high, closes below midpoint)
					{Open: 121, High: 124, Low: 118, Close: 120}, // Normal candle
				}
				return makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold: 0.01,
					ShadowRatio:   2.0,
				})
			},
			pngCRC: 0x1f2e9bc3,
		},
		{
			name: "engulfing_and_stars",
			optGen: func() CandlestickChartOption {
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105},  // Normal candle
					{Open: 110, High: 112, Low: 105, Close: 106}, // Small bearish candle
					{Open: 104, High: 115, Low: 103, Close: 114}, // Bullish engulfing
					{Open: 120, High: 125, Low: 105, Close: 108}, // Large bearish (morning star setup)
					{Open: 102, High: 104, Low: 100, Close: 103}, // Small body, gap down
					{Open: 108, High: 125, Low: 106, Close: 122}, // Large bullish, gap up (morning star completion)
					{Open: 122, High: 128, Low: 120, Close: 125}, // Normal candle
				}
				opt := makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold:    0.01,
					ShadowRatio:      2.0,
					EngulfingMinSize: 0.8,
				})
				opt.XAxis = XAxisOption{Show: Ptr(false)}
				return opt
			},
			pngCRC: 0x84037a1d,
		},
		{
			name: "combination_mixed",
			optGen: func() CandlestickChartOption {
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105},     // Normal candle
					{Open: 105, High: 108, Low: 102, Close: 105.05}, // Doji pattern
					{Open: 105, High: 107, Low: 95, Close: 106},     // Hammer pattern
					{Open: 110, High: 125, Low: 95, Close: 112},     // Spinning top pattern
					{Open: 100, High: 120, Low: 100, Close: 120},    // Bullish marubozu pattern
					{Open: 120, High: 120, Low: 100, Close: 100},    // Bearish marubozu pattern
					{Open: 110, High: 112, Low: 105, Close: 106},    // Small bearish candle
					{Open: 104, High: 115, Low: 103, Close: 114},    // Bullish engulfing
					{Open: 106, High: 125, Low: 105, Close: 107},    // Shooting star pattern
					{Open: 109, High: 110, Low: 90, Close: 108.9},   // Dragonfly doji pattern
					{Open: 108, High: 120, Low: 107, Close: 108.1},  // Gravestone doji pattern
					{Open: 108, High: 115, Low: 103, Close: 110},    // Normal candle
				}
				opt := makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold:    0.01,
					ShadowRatio:      2.0,
					EngulfingMinSize: 0.8,
				})
				opt.SeriesList[0].PatternConfig.EnabledPatterns = bulk.SliceFilterInPlace(func(pattern string) bool {
					// remove high volume patterns
					if pattern == "doji" {
						return false
					}
					return true
				}, opt.SeriesList[0].PatternConfig.EnabledPatterns)
				opt.XAxis = XAxisOption{Show: Ptr(false)}
				return opt
			},
			pngCRC: 0x0aa5328a,
		},
		{
			name: "combination_three_candle_patterns",
			optGen: func() CandlestickChartOption {
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105}, // Normal candle
					// Morning star sequence
					{Open: 120, High: 125, Low: 105, Close: 108}, // Large bearish
					{Open: 102, High: 104, Low: 100, Close: 103}, // Small body, gap down
					{Open: 108, High: 125, Low: 106, Close: 122}, // Large bullish, gap up
					// Three white soldiers sequence
					{Open: 110, High: 115, Low: 109, Close: 114}, // First soldier
					{Open: 113, High: 118, Low: 112, Close: 117}, // Second soldier
					{Open: 116, High: 121, Low: 115, Close: 120}, // Third soldier
					// Evening star sequence
					{Open: 122, High: 140, Low: 120, Close: 138}, // Large bullish
					{Open: 142, High: 144, Low: 140, Close: 143}, // Small body, gap up
					{Open: 138, High: 140, Low: 115, Close: 118}, // Large bearish, gap down
					// Three black crows sequence
					{Open: 120, High: 121, Low: 115, Close: 116}, // Second crow
					{Open: 117, High: 118, Low: 112, Close: 113}, // Third crow
					{Open: 113, High: 132, Low: 106, Close: 128}, // Normal candle
				}
				opt := makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold: 0.01,
					ShadowRatio:   2.0,
				})
				opt.XAxis = XAxisOption{Show: Ptr(false)}
				opt.SeriesList[0].PatternConfig.EnabledPatterns = []string{
					"morning_star",
					"evening_star",
				}
				return opt
			},
			pngCRC: 0x7a2d1492,
		},
		{
			name: "bullish_patterns",
			optGen: func() CandlestickChartOption {
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105},  // Normal candle
					{Open: 110, High: 112, Low: 105, Close: 106}, // Small bearish candle
					{Open: 104, High: 115, Low: 103, Close: 114}, // Bullish engulfing
					{Open: 108, High: 109, Low: 98, Close: 107},  // Hammer pattern
					{Open: 100, High: 120, Low: 100, Close: 120}, // Bullish belt hold / marubozu
					{Open: 120, High: 140, Low: 118, Close: 138}, // Large bullish
					{Open: 110, High: 119, Low: 110, Close: 118}, // Piercing line
					{Open: 118, High: 125, Low: 115, Close: 122}, // Normal candle
				}
				opt := makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold:    0.01,
					ShadowRatio:      2.0,
					EngulfingMinSize: 0.8,
				})
				opt.SeriesList[0].PatternConfig = (&CandlestickPatternConfig{}).WithPatternsBullish()
				opt.XAxis = XAxisOption{Show: Ptr(false)}
				return opt
			},
			pngCRC: 0x9b02ec55,
		},
		{
			name: "bearish_patterns",
			optGen: func() CandlestickChartOption {
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105},  // Normal candle
					{Open: 106, High: 112, Low: 105, Close: 110}, // Small bullish candle
					{Open: 114, High: 115, Low: 103, Close: 104}, // Bearish engulfing
					{Open: 106, High: 125, Low: 105, Close: 107}, // Shooting star pattern
					{Open: 120, High: 120, Low: 100, Close: 100}, // Bearish belt hold / marubozu
					{Open: 118, High: 125, Low: 117, Close: 125}, // Bullish candle
					{Open: 127, High: 128, Low: 120, Close: 121}, // Dark cloud cover
					{Open: 121, High: 124, Low: 118, Close: 120}, // Normal candle
				}
				opt := makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold:    0.01,
					ShadowRatio:      2.0,
					EngulfingMinSize: 0.8,
				})
				opt.SeriesList[0].PatternConfig = (&CandlestickPatternConfig{}).WithPatternsBearish()
				opt.XAxis = XAxisOption{Show: Ptr(false)}
				return opt
			},
			pngCRC: 0x1e8a4f88,
		},
		{
			name: "reversal_patterns",
			optGen: func() CandlestickChartOption {
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105},  // Normal candle
					{Open: 120, High: 121, Low: 115, Close: 115}, // Bearish candle
					{Open: 112, High: 119, Low: 112, Close: 118}, // Piercing line (bullish reversal)
					{Open: 118, High: 125, Low: 118, Close: 125}, // Bullish candle
					{Open: 127, High: 127, Low: 120, Close: 121}, // Dark cloud cover (bearish reversal)
					{Open: 125, High: 126, Low: 100, Close: 102}, // Bearish with low at 100
					{Open: 102, High: 108, Low: 100, Close: 107}, // Tweezer bottom (bullish reversal)
					{Open: 107, High: 112, Low: 102, Close: 110}, // Normal candle
					// Additional reversal patterns
					{Open: 115, High: 117, Low: 95, Close: 114},    // Hammer pattern (bullish reversal)
					{Open: 112, High: 130, Low: 111, Close: 113},   // Shooting star pattern (bearish reversal)
					{Open: 108, High: 110, Low: 85, Close: 108.1},  // Dragonfly doji (bullish reversal)
					{Open: 105, High: 125, Low: 104, Close: 105.1}, // Gravestone doji (bearish reversal)
					{Open: 130, High: 135, Low: 110, Close: 115},   // Large bearish for engulfing setup
					{Open: 110, High: 140, Low: 108, Close: 138},   // Bullish engulfing (reversal)
					{Open: 140, High: 145, Low: 105, Close: 110},   // Bearish engulfing (reversal)
					// Three candle reversal patterns
					{Open: 125, High: 130, Low: 105, Close: 110}, // Large bearish for morning star
					{Open: 105, High: 108, Low: 102, Close: 106}, // Small body (morning star middle)
					{Open: 110, High: 135, Low: 108, Close: 130}, // Large bullish (morning star completion)
					{Open: 115, High: 140, Low: 113, Close: 135}, // Large bullish for evening star
					{Open: 138, High: 145, Low: 136, Close: 140}, // Small body (evening star middle)
					{Open: 135, High: 136, Low: 110, Close: 115}, // Large bearish (evening star completion)
				}
				opt := makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold: 0.01,
					ShadowRatio:   2.0,
				})
				opt.SeriesList[0].PatternConfig = (&CandlestickPatternConfig{}).WithPatternsReversal()
				opt.XAxis = XAxisOption{Show: Ptr(false)}
				return opt
			},
			pngCRC: 0x4abd9948,
		},
		{
			name: "trend_patterns",
			optGen: func() CandlestickChartOption {
				data := []OHLCData{
					{Open: 100, High: 110, Low: 95, Close: 105},  // Normal candle
					{Open: 110, High: 120, Low: 100, Close: 120}, // Marubozu bullish - trend continuation
					{Open: 125, High: 125, Low: 115, Close: 115}, // Marubozu bearish - trend continuation
					{Open: 120, High: 130, Low: 115, Close: 125}, // Large bullish for belt hold setup
					{Open: 120, High: 140, Low: 120, Close: 140}, // Belt hold bullish - trend continuation
					{Open: 135, High: 135, Low: 115, Close: 115}, // Belt hold bearish - trend continuation
					{Open: 118, High: 125, Low: 117, Close: 122}, // Normal candle
					{Open: 122, High: 130, Low: 120, Close: 128}, // Trend continuation candle
				}
				opt := makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold:    0.01,
					ShadowRatio:      2.0,
					EngulfingMinSize: 0.8,
				})
				opt.SeriesList[0].PatternConfig = (&CandlestickPatternConfig{}).WithPatternsTrend()
				opt.XAxis = XAxisOption{Show: Ptr(false)}
				return opt
			},
			pngCRC: 0x410c3669,
		},
		{
			name: "all_patterns_showcase",
			optGen: func() CandlestickChartOption {
				// Comprehensive dataset showcasing all supported candlestick patterns
				data := []OHLCData{
					// 0: Setup - Normal candle
					{Open: 100, High: 110, Low: 95, Close: 105},
					// 1: Regular candle (reduce spinning top frequency)
					{Open: 105, High: 108, Low: 102, Close: 107},
					// 2: Hammer pattern
					{Open: 108, High: 109, Low: 98, Close: 107},
					// 3: Regular candle (was inverted hammer, reduce shooting star frequency)
					{Open: 95, High: 102, Low: 94, Close: 100},
					// 4: Regular candle (was shooting star, reduce frequency)
					{Open: 106, High: 115, Low: 105, Close: 112},
					// 5: Gravestone Doji pattern
					{Open: 108, High: 120, Low: 107, Close: 108.1},
					// 6: Hammer-like pattern (preserve dragonfly, reduce doji frequency)
					{Open: 109, High: 111, Low: 90, Close: 108},
					// 7: Bullish Marubozu pattern
					{Open: 100, High: 120, Low: 100, Close: 120},
					// 8: Bearish Marubozu pattern
					{Open: 120, High: 120, Low: 100, Close: 100},
					// 9: Regular candle (break harami pattern, reduce spinning top)
					{Open: 110, High: 120, Low: 107, Close: 118},
					// Setup for two-candle patterns - Large bearish candle
					{Open: 130, High: 135, Low: 110, Close: 115},
					// 11: Bullish Engulfing pattern
					{Open: 110, High: 140, Low: 108, Close: 138},
					// Setup for Bearish Engulfing - Large bullish candle
					{Open: 110, High: 140, Low: 108, Close: 138},
					// 13: Bearish Engulfing pattern (fixed to properly engulf)
					{Open: 140, High: 142, Low: 105, Close: 107},
					// Setup for Harami - Large bearish candle
					{Open: 130, High: 135, Low: 100, Close: 105},
					// 15: Regular candle (break harami by extending body)
					{Open: 110, High: 125, Low: 95, Close: 120},
					// Setup for Bearish Harami - Large bullish candle
					{Open: 100, High: 135, Low: 98, Close: 130},
					// 17: Bearish Harami pattern
					{Open: 125, High: 128, Low: 120, Close: 122},
					// Setup for Piercing Line - Bearish candle
					{Open: 120, High: 125, Low: 110, Close: 112},
					// 19: Piercing Line pattern
					{Open: 108, High: 125, Low: 107, Close: 118},
					// Setup for Dark Cloud Cover - Bullish candle
					{Open: 110, High: 125, Low: 108, Close: 123},
					// 21: Dark Cloud Cover pattern (fixed to gap up and close below midpoint)
					{Open: 128, High: 130, Low: 112, Close: 115},
					// Setup for Tweezer Top - Two candles with same high
					{Open: 110, High: 130, Low: 108, Close: 125},
					// 23: Tweezer Top pattern
					{Open: 123, High: 130, Low: 115, Close: 118},
					// Setup for Tweezer Bottom - Two candles with same low
					{Open: 120, High: 125, Low: 100, Close: 105},
					// 25: Tweezer Bottom pattern
					{Open: 108, High: 115, Low: 100, Close: 112},
					// Setup for Morning Star - Large bearish candle
					{Open: 130, High: 135, Low: 110, Close: 115},
					// 27: Morning Star middle - Small body with gap down (reduce spinning top)
					{Open: 108, High: 112, Low: 107, Close: 110},
					// 28: Morning Star completion - Large bullish candle
					{Open: 115, High: 140, Low: 113, Close: 135},
					// Setup for Evening Star - Large bullish candle
					{Open: 110, High: 140, Low: 108, Close: 135},
					// 30: Evening Star middle - Small body with proper gap up (fixed)
					{Open: 137, High: 145, Low: 136, Close: 140},
					// 31: Evening Star completion - Large bearish candle (fixed)
					{Open: 135, High: 136, Low: 115, Close: 120},
					// Setup for Three White Soldiers - Start with bearish sentiment
					{Open: 120, High: 125, Low: 110, Close: 115},
					// 33: Three White Soldiers - First soldier
					{Open: 118, High: 128, Low: 116, Close: 125},
					// 34: Three White Soldiers - Second soldier
					{Open: 127, High: 135, Low: 125, Close: 132},
					// 35: Three White Soldiers - Third soldier
					{Open: 134, High: 142, Low: 132, Close: 140},
					// Setup for Three Black Crows - Start with bullish sentiment
					{Open: 130, High: 145, Low: 128, Close: 142},
					// 37: Three Black Crows - First crow (fixed to open within previous body)
					{Open: 138, High: 140, Low: 128, Close: 132},
					// 38: Three Black Crows - Second crow (fixed to open within previous body)
					{Open: 130, High: 132, Low: 120, Close: 125},
					// 39: Three Black Crows - Third crow (fixed to open within previous body)
					{Open: 124, High: 127, Low: 115, Close: 118},
					// 40: Regular candle (reduce spinning top frequency)
					{Open: 115, High: 120, Low: 114, Close: 118},
					// 41: Regular candle (was spinning top, reduce frequency)
					{Open: 118, High: 125, Low: 115, Close: 122},
					// 42: Setup for Shooting Star - rising trend
					{Open: 120, High: 125, Low: 118, Close: 124},
					// 43: Shooting Star pattern - long upper shadow, small body near low
					{Open: 123, High: 140, Low: 122, Close: 125},
					// 44: Setup for Gravestone Doji - uptrend
					{Open: 125, High: 130, Low: 123, Close: 128},
					// 45: Gravestone Doji pattern - doji with long upper shadow
					{Open: 128, High: 145, Low: 127, Close: 128.05},
					// 46: Setup for Dragonfly Doji - downtrend
					{Open: 128, High: 130, Low: 125, Close: 126},
					// 47: Dragonfly Doji pattern - doji with long lower shadow
					{Open: 125, High: 126, Low: 110, Close: 125.05},
					// 48: Setup for Tweezer Bottom - bearish candle
					{Open: 125, High: 127, Low: 115, Close: 118},
					// 49: Tweezer Bottom pattern - same low as previous, bullish reversal
					{Open: 120, High: 125, Low: 115, Close: 123},
					// 50: Setup for Three Black Crows - high bullish candle
					{Open: 120, High: 135, Low: 118, Close: 133},
					// 51: Three Black Crows - First crow (bearish, substantial body)
					{Open: 132, High: 133, Low: 125, Close: 126},
					// 52: Three Black Crows - Second crow (bearish, opens within prev body, closes lower)
					{Open: 130, High: 131, Low: 121, Close: 122},
					// 53: Three Black Crows - Third crow (bearish, opens within prev body, closes lower)
					{Open: 125, High: 126, Low: 115, Close: 116},
					// 54: Long-Legged Doji pattern - very long shadows on both sides, small body
					{Open: 118, High: 135, Low: 95, Close: 118.1},
				}

				opt := makePatternChartOption(data, CandlestickPatternConfig{
					DojiThreshold:    0.01,
					ShadowRatio:      2.0,
					EngulfingMinSize: 0.8,
				})
				opt.XAxis = XAxisOption{Show: Ptr(false)}
				return opt
			},
			pngCRC: 0x2272c6bb,
		},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        800,
				Height:       600,
			})
			r := NewPainter(PainterOptions{
				OutputFormat: ChartOutputPNG,
				Width:        800,
				Height:       600,
			})

			opt := tc.optGen()
			opt.Theme = GetTheme(ThemeVividLight)

			validateCandlestickChartRender(t, p, r, opt, tc.pngCRC)
		})
	}
}

func TestCandlestickPatternConfigMergePatterns(t *testing.T) {
	t.Parallel()

	t.Run("merge_two_configs", func(t *testing.T) {
		config1 := &CandlestickPatternConfig{
			PreferPatternLabels: true,
			EnabledPatterns:     []string{"doji", "hammer"},
			DojiThreshold:       0.01,
		}
		config2 := &CandlestickPatternConfig{
			PreferPatternLabels: false,
			EnabledPatterns:     []string{"shooting_star", "doji"}, // Doji is duplicate
			DojiThreshold:       0.02,
		}

		merged := config1.MergePatterns(config2)

		// Should preserve config1's settings
		assert.True(t, merged.PreferPatternLabels)
		assert.InDelta(t, 0.01, merged.DojiThreshold, 0)

		// Should have union of patterns without duplicates, preserving order
		assert.Len(t, merged.EnabledPatterns, 3)
		assert.Equal(t, "doji", merged.EnabledPatterns[0])
		assert.Equal(t, "hammer", merged.EnabledPatterns[1])
		assert.Equal(t, "shooting_star", merged.EnabledPatterns[2])
	})

	t.Run("merge_with_nil", func(t *testing.T) {
		config := &CandlestickPatternConfig{
			PreferPatternLabels: true,
			EnabledPatterns:     []string{"doji", "hammer"},
		}

		// Merge nil with config
		var nilConfig *CandlestickPatternConfig
		merged1 := nilConfig.MergePatterns(config)
		assert.NotNil(t, merged1)
		assert.True(t, merged1.PreferPatternLabels)
		assert.Len(t, merged1.EnabledPatterns, 2)

		// Merge config with nil
		merged2 := config.MergePatterns(nil)
		assert.NotNil(t, merged2)
		assert.True(t, merged2.PreferPatternLabels)
		assert.Len(t, merged2.EnabledPatterns, 2)

		// Merge nil with nil
		merged3 := nilConfig.MergePatterns(nil)
		assert.Nil(t, merged3)
	})

	t.Run("merge_identical_patterns", func(t *testing.T) {
		config1 := &CandlestickPatternConfig{
			EnabledPatterns: []string{"doji", "hammer", "shooting_star"},
		}
		config2 := &CandlestickPatternConfig{
			EnabledPatterns: []string{"doji", "hammer", "shooting_star"},
		}

		merged := config1.MergePatterns(config2)
		assert.Len(t, merged.EnabledPatterns, 3) // No duplicates
		assert.Equal(t, "doji", merged.EnabledPatterns[0])
		assert.Equal(t, "hammer", merged.EnabledPatterns[1])
		assert.Equal(t, "shooting_star", merged.EnabledPatterns[2])
	})

	t.Run("merge_empty_patterns", func(t *testing.T) {
		config1 := &CandlestickPatternConfig{
			PreferPatternLabels: true,
			EnabledPatterns:     []string{},
		}
		config2 := &CandlestickPatternConfig{
			EnabledPatterns: []string{"doji", "hammer"},
		}

		merged := config1.MergePatterns(config2)
		assert.True(t, merged.PreferPatternLabels)
		assert.Len(t, merged.EnabledPatterns, 2)
		assert.Equal(t, "doji", merged.EnabledPatterns[0])
		assert.Equal(t, "hammer", merged.EnabledPatterns[1])
	})

	t.Run("merge_predefined_configs", func(t *testing.T) {
		core := (&CandlestickPatternConfig{}).WithPatternsCore()
		trend := (&CandlestickPatternConfig{}).WithPatternsTrend()

		merged := core.MergePatterns(trend)

		// Should have all patterns from both configs
		assert.Len(t, merged.EnabledPatterns, len(core.EnabledPatterns)+len(trend.EnabledPatterns))

		// Should preserve core config's settings
		assert.Equal(t, core.PreferPatternLabels, merged.PreferPatternLabels)

		// Should contain patterns from both
		assert.Contains(t, merged.EnabledPatterns, "engulfing_bull") // From core
		assert.Contains(t, merged.EnabledPatterns, "marubozu_bull")  // From Trend
	})
}
