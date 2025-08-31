package charts

import (
	"math"
	"strings"
)

const (
	/** Single candle patterns **/

	// candlestickPatternDoji represents a doji candle where open and close prices are nearly equal, indicating market indecision.
	candlestickPatternDoji = "doji"
	// candlestickPatternHammer represents a hammer candle with a small body and long lower shadow, signaling potential bullish reversal.
	candlestickPatternHammer = "hammer"
	// candlestickPatternInvertedHammer represents an inverted hammer with a small body and long upper shadow, signaling potential bullish reversal.
	candlestickPatternInvertedHammer = "inverted_hammer"
	// candlestickPatternShootingStar represents a shooting star with a small body and long upper shadow, signaling potential bearish reversal.
	candlestickPatternShootingStar = "shooting_star"
	// candlestickPatternGravestone represents a gravestone doji with long upper shadow and no lower shadow, indicating bearish sentiment.
	candlestickPatternGravestone = "gravestone_doji"
	// candlestickPatternDragonfly represents a dragonfly doji with long lower shadow and no upper shadow, indicating bullish sentiment.
	candlestickPatternDragonfly = "dragonfly_doji"
	// candlestickPatternMarubozuBull represents a bullish marubozu with no shadows and closing at the high, showing strong buying pressure.
	candlestickPatternMarubozuBull = "marubozu_bull"
	// candlestickPatternMarubozuBear represents a bearish marubozu with no shadows and closing at the low, showing strong selling pressure.
	candlestickPatternMarubozuBear = "marubozu_bear"

	/** Two candle patterns **/

	// candlestickPatternEngulfingBull represents a bullish engulfing pattern where a large bullish candle engulfs the previous bearish candle.
	candlestickPatternEngulfingBull = "engulfing_bull"
	// candlestickPatternEngulfingBear represents a bearish engulfing pattern where a large bearish candle engulfs the previous bullish candle.
	candlestickPatternEngulfingBear = "engulfing_bear"
	// candlestickPatternPiercingLine represents a piercing line where a bullish candle closes above the midpoint of the previous bearish candle.
	candlestickPatternPiercingLine = "piercing_line"
	// candlestickPatternDarkCloudCover represents a dark cloud cover where a bearish candle closes below the midpoint of the previous bullish candle.
	candlestickPatternDarkCloudCover = "dark_cloud_cover"

	/** Three candle patterns **/

	// candlestickPatternMorningStar represents a bullish morning star pattern with a doji or small candle between two opposite-colored candles.
	candlestickPatternMorningStar = "morning_star"
	// candlestickPatternEveningStar represents a bearish evening star pattern with a doji or small candle between two opposite-colored candles.
	candlestickPatternEveningStar = "evening_star"
)

// PatternFormatter allows custom formatting of detected patterns.
type PatternFormatter func(patterns []PatternDetectionResult, seriesName string, value float64) (string, *LabelStyle)

// CandlestickPatternConfig configures automatic pattern detection.
// EXPERIMENTAL: Pattern detection logic is under active development and may change in future versions.
type CandlestickPatternConfig struct {
	// PreferPatternLabels controls pattern/user label precedence
	// true = pattern labels have priority over user labels, false = user labels have priority over pattern labels
	PreferPatternLabels bool

	// PatternFormatter allows custom formatting, if nil, uses default formatting with theme colors.
	PatternFormatter PatternFormatter

	// EnabledPatterns lists specific patterns to detect
	// nil or empty = no patterns detected (PatternConfig must be set to enable)
	// Use With* methods to configure patterns
	EnabledPatterns []string

	// DojiThreshold is the body-to-range ratio threshold for doji pattern detection.
	// Default: 0.05 (5% - standard textbook definition)
	// A candlestick where the body is ≤5% of the total range is considered a doji.
	DojiThreshold float64

	// ShadowTolerance is the shadow-to-range ratio threshold for patterns requiring minimal shadows.
	// Used by marubozu patterns to determine acceptable shadow size.
	// Default: 0.01 (1% of range)
	ShadowTolerance float64

	// ShadowRatio is the minimum shadow-to-body ratio for patterns requiring long shadows.
	// Used by hammer, shooting star, and similar patterns.
	// Default: 2.0 (shadow must be at least 2x the body - standard textbook definition)
	ShadowRatio float64

	// EngulfingMinSize is the minimum size ratio for engulfing patterns.
	// The engulfing candle body must be at least this percentage of the engulfed candle body.
	// Default: 1.0 (100% - must completely engulf the previous body)
	EngulfingMinSize float64
}

// MergePatterns creates a new CandlestickPatternConfig by combining the enabled patterns config with another.
// It returns a union of both pattern sets with the current config taking precedence for other settings.
func (c *CandlestickPatternConfig) MergePatterns(other *CandlestickPatternConfig) *CandlestickPatternConfig {
	if c == nil && other == nil {
		return nil
	} else if c == nil {
		result := *other // Return a copy of other
		result.EnabledPatterns = make([]string, len(other.EnabledPatterns))
		copy(result.EnabledPatterns, other.EnabledPatterns)
		return &result
	} else if other == nil {
		result := *c // Return a copy of c
		result.EnabledPatterns = make([]string, len(c.EnabledPatterns))
		copy(result.EnabledPatterns, c.EnabledPatterns)
		return &result
	}

	// Create union of patterns, preserving order and avoiding duplicates
	seen := make(map[string]bool)
	var mergedPatterns []string
	// Add patterns from current config first (in order)
	for _, pattern := range c.EnabledPatterns {
		if !seen[pattern] {
			mergedPatterns = append(mergedPatterns, pattern)
			seen[pattern] = true
		}
	}
	// Add patterns from other config that aren't already present
	for _, pattern := range other.EnabledPatterns {
		if !seen[pattern] {
			mergedPatterns = append(mergedPatterns, pattern)
			seen[pattern] = true
		}
	}

	// Merge numeric configuration fields (keeping > 0 values with priority to c if both are set)
	dojiThreshold := c.DojiThreshold
	if dojiThreshold <= 0 {
		dojiThreshold = other.DojiThreshold
	}
	shadowTolerance := c.ShadowTolerance
	if shadowTolerance <= 0 {
		shadowTolerance = other.ShadowTolerance
	}
	shadowRatio := c.ShadowRatio
	if shadowRatio <= 0 {
		shadowRatio = other.ShadowRatio
	}
	engulfingMinSize := c.EngulfingMinSize
	if engulfingMinSize <= 0 {
		engulfingMinSize = other.EngulfingMinSize
	}

	return &CandlestickPatternConfig{
		PreferPatternLabels: c.PreferPatternLabels,
		EnabledPatterns:     mergedPatterns,
		PatternFormatter:    c.PatternFormatter,
		DojiThreshold:       dojiThreshold,
		ShadowTolerance:     shadowTolerance,
		ShadowRatio:         shadowRatio,
		EngulfingMinSize:    engulfingMinSize,
	}
}

// PatternDetectionResult holds detected pattern information.
type PatternDetectionResult struct {
	// Index is the pattern's series data point position.
	Index int
	// PatternName is the display name (e.g., "Doji", "Hammer").
	PatternName string
	// PatternType is the identifier constant (e.g., CandlestickPatternDoji).
	PatternType string
}

// addPattern is a helper to add a pattern to the list if not already present.
func (c *CandlestickPatternConfig) addPattern(pattern string) {
	for _, p := range c.EnabledPatterns {
		if p == pattern {
			return
		}
	}
	c.EnabledPatterns = append(c.EnabledPatterns, pattern)
}

// addPattern is a helper to add a set of patterns to the list if not already present.
func (c *CandlestickPatternConfig) addPatterns(patterns ...string) {
	if len(c.EnabledPatterns) == 0 {
		c.EnabledPatterns = patterns
	} else {
		for _, pattern := range patterns {
			c.addPattern(pattern)
		}
	}
}

// WithPatternsAll enables all standard patterns.
func (c *CandlestickPatternConfig) WithPatternsAll() *CandlestickPatternConfig {
	c.addPatterns(
		// Strong reversal patterns
		candlestickPatternEngulfingBull, candlestickPatternEngulfingBear, candlestickPatternHammer,
		candlestickPatternMorningStar, candlestickPatternEveningStar, candlestickPatternShootingStar,
		// Moderate patterns
		candlestickPatternDarkCloudCover, candlestickPatternDragonfly, candlestickPatternGravestone,
		candlestickPatternMarubozuBear, candlestickPatternMarubozuBull, candlestickPatternPiercingLine,
		// Neutral/indecision patterns
		candlestickPatternDoji, candlestickPatternInvertedHammer,
	)
	return c
}

// WithPatternsCore enables only the most reliable patterns that work well without volume.
func (c *CandlestickPatternConfig) WithPatternsCore() *CandlestickPatternConfig {
	c.addPatterns(
		candlestickPatternEngulfingBull, candlestickPatternEngulfingBear,
		candlestickPatternHammer, candlestickPatternShootingStar,
		candlestickPatternMorningStar, candlestickPatternEveningStar,
	)
	return c
}

// WithPatternsBullish enables only bullish patterns.
func (c *CandlestickPatternConfig) WithPatternsBullish() *CandlestickPatternConfig {
	c.addPatterns(
		candlestickPatternHammer, candlestickPatternInvertedHammer, candlestickPatternDragonfly,
		candlestickPatternMarubozuBull, candlestickPatternEngulfingBull, candlestickPatternPiercingLine,
		candlestickPatternMorningStar,
	)
	return c
}

// WithPatternsBearish enables only bearish patterns.
func (c *CandlestickPatternConfig) WithPatternsBearish() *CandlestickPatternConfig {
	c.addPatterns(
		candlestickPatternShootingStar, candlestickPatternGravestone, candlestickPatternMarubozuBear,
		candlestickPatternEngulfingBear, candlestickPatternDarkCloudCover, candlestickPatternEveningStar,
	)
	return c
}

// WithPatternsReversal enables only reversal patterns.
func (c *CandlestickPatternConfig) WithPatternsReversal() *CandlestickPatternConfig {
	c.addPatterns(
		// Single candle reversals
		candlestickPatternHammer, candlestickPatternShootingStar,
		candlestickPatternDragonfly, candlestickPatternGravestone,
		// Two candle reversals
		candlestickPatternEngulfingBull, candlestickPatternEngulfingBear,
		candlestickPatternPiercingLine, candlestickPatternDarkCloudCover,
		// Three candle reversals
		candlestickPatternMorningStar, candlestickPatternEveningStar,
	)
	return c
}

// WithPatternsTrend enables only trend continuation patterns.
func (c *CandlestickPatternConfig) WithPatternsTrend() *CandlestickPatternConfig {
	c.addPatterns(
		candlestickPatternMarubozuBull, candlestickPatternMarubozuBear,
	)
	return c
}

// WithDoji adds the doji pattern.
func (c *CandlestickPatternConfig) WithDoji() *CandlestickPatternConfig {
	c.addPattern(candlestickPatternDoji)
	return c
}

// WithHammer adds the hammer pattern.
func (c *CandlestickPatternConfig) WithHammer() *CandlestickPatternConfig {
	c.addPattern(candlestickPatternHammer)
	return c
}

// WithInvertedHammer adds the inverted hammer pattern.
func (c *CandlestickPatternConfig) WithInvertedHammer() *CandlestickPatternConfig {
	c.addPattern(candlestickPatternInvertedHammer)
	return c
}

// WithShootingStar adds the shooting star pattern.
func (c *CandlestickPatternConfig) WithShootingStar() *CandlestickPatternConfig {
	c.addPattern(candlestickPatternShootingStar)
	return c
}

// WithGravestone adds the gravestone doji pattern.
func (c *CandlestickPatternConfig) WithGravestone() *CandlestickPatternConfig {
	c.addPattern(candlestickPatternGravestone)
	return c
}

// WithDragonfly adds the dragonfly doji pattern.
func (c *CandlestickPatternConfig) WithDragonfly() *CandlestickPatternConfig {
	c.addPattern(candlestickPatternDragonfly)
	return c
}

// WithMarubozuBull adds the bullish marubozu pattern.
func (c *CandlestickPatternConfig) WithMarubozuBull() *CandlestickPatternConfig {
	c.addPattern(candlestickPatternMarubozuBull)
	return c
}

// WithMarubozuBear adds the bearish marubozu pattern.
func (c *CandlestickPatternConfig) WithMarubozuBear() *CandlestickPatternConfig {
	c.addPattern(candlestickPatternMarubozuBear)
	return c
}

// WithEngulfingBull adds the bullish engulfing pattern.
func (c *CandlestickPatternConfig) WithEngulfingBull() *CandlestickPatternConfig {
	c.addPattern(candlestickPatternEngulfingBull)
	return c
}

// WithEngulfingBear adds the bearish engulfing pattern.
func (c *CandlestickPatternConfig) WithEngulfingBear() *CandlestickPatternConfig {
	c.addPattern(candlestickPatternEngulfingBear)
	return c
}

// WithPiercingLine adds the piercing line pattern.
func (c *CandlestickPatternConfig) WithPiercingLine() *CandlestickPatternConfig {
	c.addPattern(candlestickPatternPiercingLine)
	return c
}

// WithDarkCloudCover adds the dark cloud cover pattern.
func (c *CandlestickPatternConfig) WithDarkCloudCover() *CandlestickPatternConfig {
	c.addPattern(candlestickPatternDarkCloudCover)
	return c
}

// WithMorningStar adds the morning star pattern.
func (c *CandlestickPatternConfig) WithMorningStar() *CandlestickPatternConfig {
	c.addPattern(candlestickPatternMorningStar)
	return c
}

// WithEveningStar adds the evening star pattern.
func (c *CandlestickPatternConfig) WithEveningStar() *CandlestickPatternConfig {
	c.addPattern(candlestickPatternEveningStar)
	return c
}

// WithPreferPatternLabels sets whether pattern labels have priority over user labels.
func (c *CandlestickPatternConfig) WithPreferPatternLabels(prefer bool) *CandlestickPatternConfig {
	c.PreferPatternLabels = prefer
	return c
}

// WithPatternFormatter sets a custom pattern formatter.
func (c *CandlestickPatternConfig) WithPatternFormatter(formatter PatternFormatter) *CandlestickPatternConfig {
	c.PatternFormatter = formatter
	return c
}

// WithDojiThreshold sets the doji threshold (default: 0.05).
func (c *CandlestickPatternConfig) WithDojiThreshold(threshold float64) *CandlestickPatternConfig {
	c.DojiThreshold = threshold
	return c
}

// WithShadowTolerance sets the shadow tolerance (default: 0.01).
func (c *CandlestickPatternConfig) WithShadowTolerance(tolerance float64) *CandlestickPatternConfig {
	c.ShadowTolerance = tolerance
	return c
}

// WithShadowRatio sets the shadow ratio (default: 2.0).
func (c *CandlestickPatternConfig) WithShadowRatio(ratio float64) *CandlestickPatternConfig {
	c.ShadowRatio = ratio
	return c
}

// WithEngulfingMinSize sets the engulfing minimum size (default: 1.0).
func (c *CandlestickPatternConfig) WithEngulfingMinSize(size float64) *CandlestickPatternConfig {
	c.EngulfingMinSize = size
	return c
}

// scanForCandlestickPatterns scans entire series upfront for configured patterns (private)
func scanForCandlestickPatterns(data []OHLCData, config CandlestickPatternConfig) map[int][]PatternDetectionResult {
	if len(config.EnabledPatterns) == 0 {
		return nil
	}

	patternMap := make(map[int][]PatternDetectionResult)
	for _, patternType := range config.EnabledPatterns {
		detector, ok := patternDetectors[patternType]
		if !ok {
			continue
		}
		// Scan series for this specific pattern
		for i := detector.minCandles - 1; i < len(data); i++ {
			if detector.detectFunc(data, i, config) {
				patternMap[i] = append(patternMap[i], PatternDetectionResult{
					Index:       i,
					PatternName: detector.patternName,
					PatternType: patternType,
				})
			}
		}
	}

	return patternMap
}

func detectDojiAt(data []OHLCData, index int, options CandlestickPatternConfig) bool {
	ohlc := data[index]
	if !validateOHLCData(ohlc) {
		return false
	}

	// Use configured threshold or default to standard 5% (0.05)
	threshold := options.DojiThreshold
	if threshold <= 0 {
		threshold = 0.05 // Standard textbook definition: body ≤5% of range
	}

	bodySize := math.Abs(ohlc.Close - ohlc.Open)
	priceRange := ohlc.High - ohlc.Low
	if priceRange == 0 {
		return false
	}
	return (bodySize / priceRange) <= threshold
}

func detectHammerAt(data []OHLCData, index int, options CandlestickPatternConfig) bool {
	ohlc := data[index]
	if !validateOHLCData(ohlc) {
		return false
	}

	// Use configured ratio or default to standard 2:1
	shadowRatio := options.ShadowRatio
	if shadowRatio <= 0 {
		shadowRatio = 2.0 // Standard: shadow at least 2x the body
	}

	bodySize := math.Abs(ohlc.Close - ohlc.Open)
	lowerShadow := math.Min(ohlc.Open, ohlc.Close) - ohlc.Low
	upperShadow := ohlc.High - math.Max(ohlc.Open, ohlc.Close)

	// Hammer: long lower shadow, short upper shadow, small body
	return lowerShadow >= shadowRatio*bodySize && upperShadow <= lowerShadow*0.3
}

func detectInvertedHammerAt(data []OHLCData, index int, options CandlestickPatternConfig) bool {
	ohlc := data[index]
	if !validateOHLCData(ohlc) {
		return false
	}

	// Use configured ratio or default to standard 2:1
	shadowRatio := options.ShadowRatio
	if shadowRatio <= 0 {
		shadowRatio = 2.0 // Standard: shadow at least 2x the body
	}

	bodySize := math.Abs(ohlc.Close - ohlc.Open)
	lowerShadow := math.Min(ohlc.Open, ohlc.Close) - ohlc.Low
	upperShadow := ohlc.High - math.Max(ohlc.Open, ohlc.Close)

	// Inverted hammer: long upper shadow, short lower shadow, small body
	return upperShadow >= shadowRatio*bodySize && lowerShadow <= upperShadow*0.3
}

func detectShootingStarAt(data []OHLCData, index int, options CandlestickPatternConfig) bool {
	ohlc := data[index]
	if !validateOHLCData(ohlc) {
		return false
	}

	// Use configured ratio or default to standard 2:1
	shadowRatio := options.ShadowRatio
	if shadowRatio <= 0 {
		shadowRatio = 2.0 // Standard: shadow at least 2x the body
	}

	bodySize := math.Abs(ohlc.Close - ohlc.Open)
	lowerShadow := math.Min(ohlc.Open, ohlc.Close) - ohlc.Low
	upperShadow := ohlc.High - math.Max(ohlc.Open, ohlc.Close)

	// Shooting star: long upper shadow, relatively small lower shadow, small body near the low
	hasLongUpperShadow := upperShadow >= shadowRatio*bodySize
	hasShortLowerShadow := lowerShadow <= upperShadow*0.3

	// Body should be in lower third of the total range
	totalRange := ohlc.High - ohlc.Low
	if totalRange == 0 {
		return false
	}
	bodyPosition := (math.Min(ohlc.Open, ohlc.Close) - ohlc.Low) / totalRange
	isNearLow := bodyPosition <= 0.33

	return hasLongUpperShadow && hasShortLowerShadow && isNearLow
}

func detectGravestoneDojiAt(data []OHLCData, index int, options CandlestickPatternConfig) bool {
	ohlc := data[index]
	if !validateOHLCData(ohlc) {
		return false
	}

	// Must be a doji first
	threshold := options.DojiThreshold
	if threshold <= 0 {
		threshold = 0.05 // Standard: body ≤5% of range
	}

	bodySize := math.Abs(ohlc.Close - ohlc.Open)
	priceRange := ohlc.High - ohlc.Low
	if priceRange == 0 {
		return false
	}
	if (bodySize / priceRange) > threshold {
		return false
	}

	bodyMidpoint := (ohlc.Open + ohlc.Close) / 2
	upperShadow := ohlc.High - bodyMidpoint
	lowerShadow := bodyMidpoint - ohlc.Low

	shadowRatio := options.ShadowRatio
	if shadowRatio <= 0 {
		shadowRatio = 2.0 // Standard: shadow at least 2x the body
	}

	// Gravestone doji: long upper shadow, minimal lower shadow
	hasLongUpperShadow := upperShadow >= shadowRatio*math.Abs(ohlc.Close-ohlc.Open)
	hasMinimalLowerShadow := lowerShadow <= upperShadow*0.3

	return hasLongUpperShadow && hasMinimalLowerShadow
}

func detectDragonflyDojiAt(data []OHLCData, index int, options CandlestickPatternConfig) bool {
	ohlc := data[index]
	if !validateOHLCData(ohlc) {
		return false
	}

	// Must be a doji first
	threshold := options.DojiThreshold
	if threshold <= 0 {
		threshold = 0.05 // Standard: body ≤5% of range
	}

	bodySize := math.Abs(ohlc.Close - ohlc.Open)
	priceRange := ohlc.High - ohlc.Low
	if priceRange == 0 {
		return false
	}
	if (bodySize / priceRange) > threshold {
		return false
	}

	bodyMidpoint := (ohlc.Open + ohlc.Close) / 2
	upperShadow := ohlc.High - bodyMidpoint
	lowerShadow := bodyMidpoint - ohlc.Low

	shadowRatio := options.ShadowRatio
	if shadowRatio <= 0 {
		shadowRatio = 2.0 // Standard: shadow at least 2x the body
	}

	// Dragonfly doji: long lower shadow, minimal upper shadow
	hasLongLowerShadow := lowerShadow >= shadowRatio*math.Abs(ohlc.Close-ohlc.Open)
	hasMinimalUpperShadow := upperShadow <= lowerShadow*0.3

	return hasLongLowerShadow && hasMinimalUpperShadow
}

func detectBullishMarubozuAt(data []OHLCData, index int, options CandlestickPatternConfig) bool {
	ohlc := data[index]
	if !validateOHLCData(ohlc) {
		return false
	}

	// Use configured tolerance or default to 1% of range
	threshold := options.ShadowTolerance
	if threshold <= 0 {
		threshold = 0.01 // Standard: shadows ≤1% of range
	}

	// Calculate shadow sizes
	upper := ohlc.High - math.Max(ohlc.Open, ohlc.Close)
	lower := math.Min(ohlc.Open, ohlc.Close) - ohlc.Low
	body := math.Abs(ohlc.Close - ohlc.Open)
	total := ohlc.High - ohlc.Low

	if total == 0 || body == 0 {
		return false
	}

	// Shadows should be minimal compared to total range
	hasMinimalShadows := (upper+lower)/total <= threshold

	if !hasMinimalShadows {
		return false
	}
	return ohlc.Close > ohlc.Open
}

func detectBearishMarubozuAt(data []OHLCData, index int, options CandlestickPatternConfig) bool {
	ohlc := data[index]
	if !validateOHLCData(ohlc) {
		return false
	}

	// Use configured tolerance or default to 1% of range
	threshold := options.ShadowTolerance
	if threshold <= 0 {
		threshold = 0.01 // Standard: shadows ≤1% of range
	}

	// Calculate shadow sizes
	upper := ohlc.High - math.Max(ohlc.Open, ohlc.Close)
	lower := math.Min(ohlc.Open, ohlc.Close) - ohlc.Low
	body := math.Abs(ohlc.Close - ohlc.Open)
	total := ohlc.High - ohlc.Low

	if total == 0 || body == 0 {
		return false
	}

	// Shadows should be minimal compared to total range
	hasMinimalShadows := (upper+lower)/total <= threshold

	if !hasMinimalShadows {
		return false
	}
	return ohlc.Close < ohlc.Open
}

func detectBullishEngulfingAt(data []OHLCData, index int, options CandlestickPatternConfig) bool {
	if index < 1 {
		return false
	}
	prev := data[index-1]
	current := data[index]
	if !validateOHLCData(prev) || !validateOHLCData(current) {
		return false
	}

	// Use configured size or default to complete engulfing (100%)
	minSize := options.EngulfingMinSize
	if minSize <= 0 {
		minSize = 1.0 // Standard: must completely engulf previous body
	}

	prevBody := math.Abs(prev.Close - prev.Open)
	currentBody := math.Abs(current.Close - current.Open)

	// Current candle must engulf previous candle's body
	prevTop := math.Max(prev.Open, prev.Close)
	prevBottom := math.Min(prev.Open, prev.Close)
	currentTop := math.Max(current.Open, current.Close)
	currentBottom := math.Min(current.Open, current.Close)

	isEngulfing := currentTop > prevTop && currentBottom < prevBottom
	isSizeSignificant := currentBody >= minSize*prevBody

	if !isEngulfing || !isSizeSignificant {
		return false
	}

	// Determine bullish
	prevBearish := prev.Close < prev.Open
	currentBullish := current.Close > current.Open

	return prevBearish && currentBullish
}

func detectBearishEngulfingAt(data []OHLCData, index int, options CandlestickPatternConfig) bool {
	if index < 1 {
		return false
	}
	prev := data[index-1]
	current := data[index]
	if !validateOHLCData(prev) || !validateOHLCData(current) {
		return false
	}

	// Use configured size or default to complete engulfing (100%)
	minSize := options.EngulfingMinSize
	if minSize <= 0 {
		minSize = 1.0 // Standard: must completely engulf previous body
	}

	prevBody := math.Abs(prev.Close - prev.Open)
	currentBody := math.Abs(current.Close - current.Open)

	// Current candle must engulf previous candle's body
	prevTop := math.Max(prev.Open, prev.Close)
	prevBottom := math.Min(prev.Open, prev.Close)
	currentTop := math.Max(current.Open, current.Close)
	currentBottom := math.Min(current.Open, current.Close)

	isEngulfing := currentTop > prevTop && currentBottom < prevBottom
	isSizeSignificant := currentBody >= minSize*prevBody

	if !isEngulfing || !isSizeSignificant {
		return false
	}

	// Determine bearish
	prevBullish := prev.Close > prev.Open
	currentBearish := current.Close < current.Open

	return prevBullish && currentBearish
}

func detectPiercingLineAt(data []OHLCData, index int, _ CandlestickPatternConfig) bool {
	if index < 1 {
		return false
	}
	prev := data[index-1]
	current := data[index]
	if !validateOHLCData(prev) || !validateOHLCData(current) {
		return false
	} else if prev.Close >= prev.Open { // Previous candle must be bearish
		return false
	} else if current.Close <= current.Open { // Current candle must be bullish
		return false
	} else if current.Open >= prev.Close { // Current must open below previous close (gap down)
		return false
	}
	// Current must close above midpoint of previous candle's body
	prevMidpoint := (prev.Open + prev.Close) / 2
	if current.Close <= prevMidpoint {
		return false
	}
	// Current close should not exceed previous open (not engulfing)
	if current.Close >= prev.Open {
		return false
	}
	return true
}

func detectDarkCloudCoverAt(data []OHLCData, index int, _ CandlestickPatternConfig) bool {
	if index < 1 {
		return false
	}
	prev := data[index-1]
	current := data[index]
	if !validateOHLCData(prev) || !validateOHLCData(current) {
		return false
	} else if prev.Close <= prev.Open { // Previous candle must be bullish
		return false
	} else if current.Close >= current.Open { // Current candle must be bearish
		return false
	} else if current.Open <= prev.Close { // Current must open above previous close (gap up)
		return false
	}
	// Current must close below midpoint of previous candle's body
	prevMidpoint := (prev.Open + prev.Close) / 2
	if current.Close >= prevMidpoint {
		return false
	}
	// Current close should not go below previous open (not engulfing)
	if current.Close <= prev.Open {
		return false
	}
	return true
}

func detectMorningStarAt(data []OHLCData, index int, _ CandlestickPatternConfig) bool {
	if index < 2 {
		return false
	}
	first := data[index-2]
	second := data[index-1]
	third := data[index]
	if !validateOHLCData(first) || !validateOHLCData(second) || !validateOHLCData(third) {
		return false
	}

	// First candle: bearish (long red)
	if first.Close >= first.Open {
		return false
	}
	firstBody := first.Open - first.Close

	// Second candle: small body (doji-like), gaps down
	secondBody := math.Abs(second.Close - second.Open)
	if secondBody > firstBody*0.3 { // Second body should be small
		return false
	}
	// Gap down: second candle should open below first candle's body (allow some overlap)
	if second.Open >= first.Close {
		return false
	}
	// Third candle: bullish (long green), gaps up
	if third.Close <= third.Open {
		return false
	}
	thirdBody := third.Close - third.Open
	// Gap up: third candle should open above second candle's body (allow some overlap)
	if third.Open <= math.Max(second.Open, second.Close) {
		return false
	}
	// Third candle should close well into first candle's body
	if third.Close <= (first.Open+first.Close)/2 {
		return false
	}
	// Bodies should be reasonably sized
	if thirdBody < firstBody*0.5 {
		return false
	}
	return true
}

func detectEveningStarAt(data []OHLCData, index int, _ CandlestickPatternConfig) bool {
	if index < 2 {
		return false
	}
	first := data[index-2]
	second := data[index-1]
	third := data[index]
	if !validateOHLCData(first) || !validateOHLCData(second) || !validateOHLCData(third) {
		return false
	}
	// First candle: bullish (long green)
	if first.Close <= first.Open {
		return false
	}
	firstBody := first.Close - first.Open
	// Second candle: small body (doji-like), gaps up
	secondBody := math.Abs(second.Close - second.Open)
	if secondBody > firstBody*0.3 { // Second body should be small
		return false
	}
	// Gap up: second candle should open above first candle's body (allow some overlap)
	if second.Open <= first.Close {
		return false
	}
	// Third candle: bearish (long red), gaps down
	if third.Close >= third.Open {
		return false
	}
	thirdBody := third.Open - third.Close
	// Gap down: third candle should open below second candle's body (allow some overlap)
	if third.Open >= math.Min(second.Open, second.Close) {
		return false
	}
	// Third candle should close well into first candle's body
	if third.Close >= (first.Open+first.Close)/2 {
		return false
	}
	// Bodies should be reasonably sized
	if thirdBody < firstBody*0.5 {
		return false
	}
	return true
}

// patternDetector defines a single pattern detection function with metadata.
type patternDetector struct {
	patternName string
	detectFunc  func([]OHLCData, int, CandlestickPatternConfig) bool
	minCandles  int
}

// patternDetectors contains all available pattern detectors organized by type
var patternDetectors = map[string]patternDetector{
	// single candle patterns
	candlestickPatternDoji:           {"Doji", detectDojiAt, 1},
	candlestickPatternHammer:         {"Hammer", detectHammerAt, 1},
	candlestickPatternInvertedHammer: {"Inverted Hammer", detectInvertedHammerAt, 1},
	candlestickPatternShootingStar:   {"Shooting Star", detectShootingStarAt, 1},
	candlestickPatternGravestone:     {"Gravestone Doji", detectGravestoneDojiAt, 1},
	candlestickPatternDragonfly:      {"Dragonfly Doji", detectDragonflyDojiAt, 1},
	candlestickPatternMarubozuBull:   {"Bullish Marubozu", detectBullishMarubozuAt, 1},
	candlestickPatternMarubozuBear:   {"Bearish Marubozu", detectBearishMarubozuAt, 1},
	// double candle patterns
	candlestickPatternEngulfingBull:  {"Bullish Engulfing", detectBullishEngulfingAt, 2},
	candlestickPatternEngulfingBear:  {"Bearish Engulfing", detectBearishEngulfingAt, 2},
	candlestickPatternPiercingLine:   {"Piercing Line", detectPiercingLineAt, 2},
	candlestickPatternDarkCloudCover: {"Dark Cloud Cover", detectDarkCloudCoverAt, 2},
	// triple candle patterns
	candlestickPatternMorningStar: {"Morning Star", detectMorningStarAt, 3},
	candlestickPatternEveningStar: {"Evening Star", detectEveningStarAt, 3},
}

// formatPatternsDefault provides default pattern formatting (private)
func formatPatternsDefault(patterns []PatternDetectionResult, seriesIndex int, theme ColorPalette) (string, *LabelStyle) {
	if len(patterns) == 0 {
		return "", nil
	}

	// Build display names and determine color
	displayNames := make([]string, len(patterns))
	var bullishCount, bearishCount, neutralCount int
	for i, pattern := range patterns {
		displayName := getPatternDisplayName(pattern.PatternType)
		if displayName == "" {
			displayName = pattern.PatternName // fallback name without icon
		}
		displayNames[i] = displayName

		// Count pattern types to determine color
		switch pattern.PatternType {
		case candlestickPatternHammer, candlestickPatternMorningStar, candlestickPatternEngulfingBull, candlestickPatternDragonfly, candlestickPatternMarubozuBull, candlestickPatternPiercingLine:
			bullishCount++
		case candlestickPatternShootingStar, candlestickPatternEveningStar, candlestickPatternEngulfingBear, candlestickPatternGravestone, candlestickPatternMarubozuBear, candlestickPatternDarkCloudCover:
			bearishCount++
		default: // Doji and other neutral patterns
			neutralCount++
		}
	}

	// Determine color based on predominant pattern type
	upColor, downColor := theme.GetSeriesUpDownColors(seriesIndex)
	var color Color
	if bullishCount > bearishCount && bullishCount > neutralCount {
		color = upColor
	} else if bearishCount > bullishCount && bearishCount > neutralCount {
		color = downColor
	} else {
		if theme.IsDark() {
			color = Color{R: 100, G: 100, B: 100, A: 255} // dark gray
		} else {
			color = Color{R: 200, G: 200, B: 200, A: 255} // light gray
		}
	}

	// Use theme-appropriate background based on dark/light mode
	var backgroundColor, fontColor Color
	if theme.IsDark() {
		backgroundColor = ColorBlack.WithAlpha(180)
		fontColor = color.WithAdjustHSL(0, 0, 0.28) // Lighter for dark backgrounds
	} else {
		backgroundColor = ColorWhite.WithAlpha(180)
		fontColor = color.WithAdjustHSL(0, 0, -0.28) // Darker for light backgrounds
	}

	return strings.Join(displayNames, "\n"), &LabelStyle{
		FontStyle: FontStyle{
			FontColor: fontColor,
			FontSize:  10,
		},
		BackgroundColor: backgroundColor,
		CornerRadius:    4,
		BorderColor:     color,
		BorderWidth:     1.2,
	}
}

/* All symbols currently supported:
Balance: ≈
Volatility: ~
Directional: ^ v
Greek: Α Β Γ Δ Ε Ζ Η Θ Ι Κ Λ Μ Ν Ξ Ο Π Ρ Σ Τ Υ Φ Χ Ψ Ω α β γ δ ε ζ η θ ι κ λ μ ν ξ ο π ρ ς σ τ υ φ χ ψ ω ϑ ϕ ϖ ϗ Ϙ ϙ Ϛ ϛ Ϝ ϝ Ϟ ϟ Ϡ ϡ ϰ ϱ ϲ ϳ ϴ ϵ ϶ Ϸ ϸ Ϲ Ϻ ϻ ϼ Ͻ Ͼ Ͽ
Geometric: ◊ ◌
Math: ∂ ∆ ∏ ∑ √ ∞ ∫ ≠ ≤ ≥
Financial: $ ¢ £ ¤ ¥ ₡ ₢ ₣ ₤ ₥ ₦ ₧ ₨ ₩ ₪ ₫ € ₭ ₮ ₯ ₰ ₱ ₲ ₳ ₴ ₵ ₶ ₷ ₸ ₹ ₺ ₻ ₼ ₽ ₾ ₿
Star: * ※ ⁎
Special: ! " # % & ' ( ) + , - . / : ; < = > ? @ K V [ \ ] _ ` { | } ¡ ¦ § ¨ © ª « ¬ ® ¯ ° ± ² ³ ´ µ ¶ · ¸ ¹ º » ¼ ½ ¾ ¿ Å × ÷ ƒ ǀ ʘ ˆ ˇ ˘ ˙ ˚ ˛ ˜ ˝ – — † ‡ • ‣ ‰ ‱ ′ ″ ‴ ‵ ‶ ‷ ‸ ‹ › ‼ ‽ ⁂ ⁄ ⁅ ⁆ ⁇ ⁈ ⁉ ⁊ ⁋ ⁌ ⁍ ⁏ ℀ ℁ ℂ ℃ ℄ ℅ ℆ ℇ ℈ ℉ ℊ ℋ ℌ ℍ ℎ ℏ ℐ ℑ ℒ ℓ ℔ ℕ № ℗ ℘ ℙ ℚ ℛ ℜ ℝ ℞ ℟ ℠ ℡ ™ ℣ ℤ ℥ ℧ ℨ ℩ ℬ ℭ ℮ ℯ ℰ ℱ Ⅎ ℳ ℴ ℵ ℶ ℷ ℸ ℹ ℺ ℻ ℼ ℽ ℾ ℿ ⅀ ⅁ ⅂ ⅃ ⅄ ⅅ ⅆ ⅇ ⅈ ⅉ ⅊ ⅋ ⅌ ⅍ ⅎ ⅏ Ʇ
*/

// getPatternDisplayName returns the pattern name with appropriate symbol.
func getPatternDisplayName(patternType string) string {
	switch patternType {
	case candlestickPatternDoji:
		// Current: ± (plus-minus, balance symbol)
		// Alternatives: ≈ (approximately equal), ∏ (product)
		return "± Doji"
	case candlestickPatternHammer:
		// Current: Γ (Greek gamma, hammer shape)
		// Alternatives: Τ (Greek tau), τ (small tau)
		return "Γ Hammer"
	case candlestickPatternInvertedHammer:
		// Current: Ʇ (turned T, upside-down hammer)
		return "Ʇ Inv. Hammer"
	case candlestickPatternShootingStar:
		// Current: ※ (reference mark, star-like)
		// Alternatives: * (asterisk), ⁎ (low asterisk), ‣ (triangular bullet), • (bullet)
		return "※ Shooting Star"
	case candlestickPatternGravestone:
		// Current: † (dagger, cross symbol)
		// Alternatives: ‡ (double dagger)
		return "† Gravestone"
	case candlestickPatternDragonfly:
		// Current: ψ (small psi, trident-like)
		// Alternatives: Ψ (capital psi), ‡ (double dagger), ◊ (geometric diamond)
		return "ψ Dragonfly"
	case candlestickPatternMarubozuBull:
		// Current: ^ (circumflex, upward direction)
		// Alternatives: Λ (lambda), Δ (delta)
		return "^ Bull Marubozu"
	case candlestickPatternMarubozuBear:
		// Current: v (lowercase v, downward direction)
		// Alternatives: V (capital v)
		return "v Bear Marubozu"
	case candlestickPatternEngulfingBull:
		// Current: Λ (Lambda, upward V shape, engulfing)
		// Alternatives: Δ (delta), < (less than)
		return "Λ Bull Engulfing"
	case candlestickPatternEngulfingBear:
		// Current: V (capital V, downward engulfing)
		// Alternatives: v (lowercase v), > (greater than)
		return "V Bear Engulfing"
	case candlestickPatternMorningStar:
		// Current: * (asterisk, star symbol)
		// Alternatives: ※ (reference mark), ⁎ (low asterisk), ‣ (triangular bullet), • (bullet)
		return "* Morning Star"
	case candlestickPatternEveningStar:
		// Current: ⁎ (low asterisk, evening star)
		// Alternatives: ※ (reference mark), * (asterisk), ‣ (triangular bullet), • (bullet)
		return "⁎ Evening Star"
	case candlestickPatternPiercingLine:
		// Current: | (vertical bar)
		// Alternatives: ǀ (dental click), ¦ (broken bar)
		return "| Piercing Line"
	case candlestickPatternDarkCloudCover:
		// Current: Ξ (Xi, horizontal lines like cloud layers)
		// Alternatives: ≈ (approximately equal), ∞ (infinity), ~ (tilde)
		return "Ξ Dark Cloud"
	default:
		return ""
	}
}
