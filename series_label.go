package charts

import (
	"sort"
	"strings"

	"github.com/dustin/go-humanize"
)

var (
	// LabelFormatterValueShort provides a short value with at most 2 decimal places.
	LabelFormatterValueShort = func(index int, name string, val float64) (string, *LabelStyle) {
		return defaultValueFormatter(val), nil
	}

	// LabelFormatterNameShortValue puts the series name next to the value with up to 2 decimal places.
	LabelFormatterNameShortValue = func(index int, name string, val float64) (string, *LabelStyle) {
		return name + ": " + defaultValueFormatter(val), nil
	}
)

// LabelFormatterThresholdMin returns a SeriesLabelFormatter that only shows labels for values above the specified threshold.
// Values at or below the threshold will have empty labels (effectively hiding them).
func LabelFormatterThresholdMin(threshold float64) SeriesLabelFormatter {
	return func(index int, name string, val float64) (string, *LabelStyle) {
		if val >= threshold {
			return defaultValueFormatter(val), nil
		}
		return "", nil
	}
}

// LabelFormatterThresholdMax returns a SeriesLabelFormatter that only shows labels for values at or below the specified threshold.
// Values above the threshold will have empty labels (effectively hiding them).
func LabelFormatterThresholdMax(threshold float64) SeriesLabelFormatter {
	return func(index int, name string, val float64) (string, *LabelStyle) {
		if val <= threshold {
			return defaultValueFormatter(val), nil
		}
		return "", nil
	}
}

// LabelFormatterTopN returns a SeriesLabelFormatter that only shows labels for the top N highest values in the provided slice.
// This formatter requires the complete data values to determine which values are in the top N.
func LabelFormatterTopN(values []float64, n int) SeriesLabelFormatter {
	if n <= 0 {
		return func(index int, name string, val float64) (string, *LabelStyle) {
			return "", nil
		}
	} else if len(values) <= n {
		return LabelFormatterValueShort // Threshold below minimum, show all values
	}

	// Sort values in ascending order to find the nth highest
	sortedValues := make([]float64, len(values))
	copy(sortedValues, values)
	sort.Float64s(sortedValues)

	// set threshold to the nth highest (from the end of ascending sorted slice)
	threshold := sortedValues[len(sortedValues)-n]
	return LabelFormatterThresholdMin(threshold)
}

// LabelFormatterGradientGreenRed returns a SeriesLabelFormatter that colors labels from green (minimum) to red (maximum).
// This formatter requires the complete data values to determine the min/max range for color interpolation.
func LabelFormatterGradientGreenRed(values []float64) SeriesLabelFormatter {
	return LabelFormatterGradientColor(values, ColorGreen, ColorYellowAlt1, ColorRed)
}

// LabelFormatterGradientRedGreen returns a SeriesLabelFormatter that colors labels from red (minimum) to green (maximum).
// This formatter requires the complete data values to determine the min/max range for color interpolation.
func LabelFormatterGradientRedGreen(values []float64) SeriesLabelFormatter {
	return LabelFormatterGradientColor(values, ColorRed, ColorYellowAlt1, ColorGreen)
}

// LabelFormatterGradientColor returns a SeriesLabelFormatter that colors labels with a gradient between multiple colors.
// The minimum value gets the first color, maximum value gets the last color, and intermediate values are interpolated.
// For two colors: lowColor -> highColor. For three+ colors: first -> second -> ... -> last.
// This formatter requires the complete data values to determine the min/max range for color interpolation.
func LabelFormatterGradientColor(values []float64, colors ...Color) SeriesLabelFormatter {
	// Handle case where no colors are provided
	if len(colors) == 0 {
		colors = []Color{ColorBlack} // default to black
	}

	summary := summarizePopulationData(values)
	minVal, maxVal := summary.Min, summary.Max

	// Handle edge case where all values are the same
	if minVal == maxVal {
		return func(index int, name string, val float64) (string, *LabelStyle) {
			return defaultValueFormatter(val), &LabelStyle{
				FontStyle: FontStyle{FontColor: colors[0]},
			}
		}
	}

	return func(index int, name string, val float64) (string, *LabelStyle) {
		// Calculate the interpolation factor (0.0 = min, 1.0 = max)
		factor := (val - minVal) / (maxVal - minVal)
		if factor < 0 {
			factor = 0
		} else if factor > 1 {
			factor = 1
		}

		// Interpolate between the colors
		interpolatedColor := interpolateMultipleColors(colors, factor)

		return defaultValueFormatter(val), &LabelStyle{
			FontStyle: FontStyle{FontColor: interpolatedColor},
		}
	}
}

// interpolateColor linearly interpolates between two colors.
// factor should be between 0.0 (returns color1) and 1.0 (returns color2).
func interpolateColor(color1, color2 Color, factor float64) Color {
	if factor < 0 {
		factor = 0
	} else if factor > 1 {
		factor = 1
	}

	r1, g1, b1, a1 := color1.RGBA()
	r2, g2, b2, a2 := color2.RGBA()

	// Convert from 16-bit to 8-bit values
	r1, g1, b1, a1 = r1>>8, g1>>8, b1>>8, a1>>8
	r2, g2, b2, a2 = r2>>8, g2>>8, b2>>8, a2>>8

	// Interpolate each component
	r := uint8(float64(r1) + factor*float64(int(r2)-int(r1)))
	g := uint8(float64(g1) + factor*float64(int(g2)-int(g1)))
	b := uint8(float64(b1) + factor*float64(int(b2)-int(b1)))
	a := uint8(float64(a1) + factor*float64(int(a2)-int(a1)))

	return Color{R: r, G: g, B: b, A: a}
}

// interpolateMultipleColors interpolates between multiple colors based on a factor (0.0 to 1.0).
// For a single color, returns that color.
// For two colors, behaves like interpolateColor.
// For three+ colors, divides the factor range equally between color segments.
func interpolateMultipleColors(colors []Color, factor float64) Color {
	if len(colors) == 0 {
		return ColorBlack // fallback
	}
	if len(colors) == 1 {
		return colors[0]
	}
	// Clamp factor to [0, 1] range for all cases
	if factor < 0 {
		factor = 0
	} else if factor > 1 {
		factor = 1
	}

	if len(colors) == 2 {
		return interpolateColor(colors[0], colors[1], factor)
	}

	// For 3+ colors, divide the factor range into segments
	numSegments := len(colors) - 1
	segmentSize := 1.0 / float64(numSegments)

	// Find which segment the factor falls into
	segmentIndex := int(factor / segmentSize)

	// Handle edge case where factor is exactly 1.0
	if segmentIndex >= numSegments {
		segmentIndex = numSegments - 1
	}

	// Calculate the local factor within this segment (0.0 to 1.0)
	localFactor := (factor - float64(segmentIndex)*segmentSize) / segmentSize

	// Interpolate between the two colors in this segment
	return interpolateColor(colors[segmentIndex], colors[segmentIndex+1], localFactor)
}

type labelRenderValue struct {
	text            string
	fontStyle       FontStyle
	x               int
	y               int
	radians         float64
	backgroundColor Color
	cornerRadius    int
	borderColor     Color
	borderWidth     float64
}

type labelValue struct {
	index     int
	value     float64
	x         int
	y         int
	radians   float64
	fontStyle FontStyle
	vertical  bool
	offset    OffsetInt
}

type seriesLabelPainter struct {
	p           *Painter
	seriesNames []string
	label       *SeriesLabel
	theme       ColorPalette
	values      []labelRenderValue
}

func newSeriesLabelPainter(p *Painter, seriesNames []string, label SeriesLabel,
	theme ColorPalette) *seriesLabelPainter {
	return &seriesLabelPainter{
		p:           p,
		seriesNames: seriesNames,
		label:       &label,
		theme:       theme,
	}
}

func (o *seriesLabelPainter) Add(value labelValue) {
	label := o.label
	if flagIs(false, label.Show) {
		return
	}
	distance := label.Distance
	if distance == 0 {
		distance = 5
	}
	var text string
	var labelStyleOverride *LabelStyle
	if label.LabelFormatter != nil {
		var name string
		if len(o.seriesNames) > value.index {
			name = o.seriesNames[value.index]
		}
		text, labelStyleOverride = label.LabelFormatter(value.index, name, value.value)
	} else {
		if label.ValueFormatter == nil {
			label.ValueFormatter = defaultValueFormatter
		}
		if label.FormatTemplate != "" {
			text = labelFormatValue(o.seriesNames, label.FormatTemplate, label.ValueFormatter,
				value.index, value.value, -1)
		} else {
			text = label.ValueFormatter(value.value)
		}
	}
	if text == "" {
		return // nothing to render
	}

	labelFontStyle := mergeFontStyles(label.FontStyle, value.fontStyle, FontStyle{
		FontColor: o.theme.GetLabelTextColor(),
		FontSize:  defaultLabelFontSize,
		Font:      getPreferredFont(label.FontStyle.Font, value.fontStyle.Font),
	})
	if labelStyleOverride != nil { // Prefer per-point style overrides if present
		labelFontStyle = mergeFontStyles(labelStyleOverride.FontStyle, labelFontStyle)
	}

	textBox := o.p.MeasureText(text, value.radians, labelFontStyle)
	renderValue := labelRenderValue{
		text:      text,
		fontStyle: labelFontStyle,
		x:         value.x,
		y:         value.y,
		radians:   value.radians,
	}

	// Set background color, corner radius, and border styling if specified
	if labelStyleOverride != nil {
		renderValue.backgroundColor = labelStyleOverride.BackgroundColor
		renderValue.cornerRadius = labelStyleOverride.CornerRadius
		renderValue.borderColor = labelStyleOverride.BorderColor
		renderValue.borderWidth = labelStyleOverride.BorderWidth
	}
	if value.vertical {
		renderValue.x -= textBox.Width() >> 1
		renderValue.y -= distance
	} else {
		renderValue.x += distance
		renderValue.y += textBox.Height() >> 1
		renderValue.y -= 2
	}
	if value.radians != 0 {
		renderValue.x = value.x + (textBox.Width() >> 1) - 1
	}
	renderValue.x += value.offset.Left
	renderValue.y += value.offset.Top
	o.values = append(o.values, renderValue)
}

// drawLabelWithBackground draws a text label with optional background styling.
// This helper function can be used by various chart types to render labels with custom styling.
//
// Parameters:
//   - p: The painter instance to draw on
//   - text: The text to render (empty strings are ignored)
//   - x, y: Text position coordinates
//   - radians: Text rotation angle in radians
//   - fontStyle: Font styling for the text
//   - backgroundColor: Background color (transparent colors are ignored)
//   - cornerRadius: Corner radius for background rectangle (negative values are treated as 0)
//   - borderColor: Border color around the background (transparent colors are ignored)
//   - borderWidth: Width of the border around the background (values <= 0 are ignored)
func drawLabelWithBackground(p *Painter, text string, x, y int, radians float64, fontStyle FontStyle,
	backgroundColor Color, cornerRadius int, borderColor Color, borderWidth float64) {
	if text == "" {
		return
	}

	if cornerRadius < 0 {
		cornerRadius = 0
	}

	if !backgroundColor.IsTransparent() || (!borderColor.IsTransparent() && borderWidth > 0) {
		textBox := p.MeasureText(text, radians, fontStyle)

		const padding = 4
		bgBox := Box{
			Left:   x - padding,
			Top:    y - textBox.Height() - padding,
			Right:  x + textBox.Width() + padding,
			Bottom: y + padding,
		}

		if cornerRadius > 0 {
			// Clamp corner radius to prevent visual artifacts
			boxWidth := bgBox.Width() / 2
			boxHeight := bgBox.Height() / 2
			maxRadius := boxWidth
			if boxHeight < boxWidth {
				maxRadius = boxHeight
			}
			if cornerRadius > maxRadius {
				cornerRadius = maxRadius
			}
			p.roundedRect(bgBox, cornerRadius, true, true, backgroundColor, borderColor, borderWidth)
		} else {
			p.FilledRect(bgBox.Left, bgBox.Top, bgBox.Right, bgBox.Bottom, backgroundColor, borderColor, borderWidth)
		}
	}

	p.Text(text, x, y, radians, fontStyle)
}

func (o *seriesLabelPainter) Render() (Box, error) {
	for _, item := range o.values {
		if item.text != "" {
			drawLabelWithBackground(o.p, item.text, item.x, item.y, item.radians,
				item.fontStyle, item.backgroundColor, item.cornerRadius, item.borderColor, item.borderWidth)
		}
	}
	return BoxZero, nil
}

// Deprecated: labelFormatPie is deprecated.
func labelFormatPie(seriesName string, layout string, valueFormatter ValueFormatter,
	value float64, percent float64) string {
	if len(layout) == 0 {
		layout = "{b}: {d}"
	}
	return newLabelFormatter([]string{seriesName}, layout, valueFormatter)(0, value, percent)
}

// Deprecated: labelFormatFunnel is deprecated.
func labelFormatFunnel(seriesName string, layout string, valueFormatter ValueFormatter,
	value float64, percent float64) string {
	if len(layout) == 0 {
		layout = "{b}({d})"
	}
	return newLabelFormatter([]string{seriesName}, layout, valueFormatter)(0, value, percent)
}

// Deprecated: labelFormatValue is deprecated.
func labelFormatValue(seriesNames []string, layout string, valueFormatter ValueFormatter,
	index int, value float64, percent float64) string {
	if len(layout) == 0 {
		layout = "{c}"
	}
	return newLabelFormatter(seriesNames, layout, valueFormatter)(index, value, percent)
}

// Deprecated: newLabelFormatter is deprecated.
func newLabelFormatter(seriesNames []string, layout string, valueFormatter ValueFormatter) func(index int, value float64, percent float64) string {
	if valueFormatter == nil {
		valueFormatter = defaultValueFormatter
	}
	return func(index int, value, percent float64) string {
		var percentText string
		if percent >= 0 {
			percentText = humanize.FtoaWithDigits(percent*100, 2) + "%"
		}
		var name string
		if len(seriesNames) > index {
			name = seriesNames[index]
		}
		text := strings.ReplaceAll(layout, "{c}", valueFormatter(value))
		text = strings.ReplaceAll(text, "{d}", percentText)
		text = strings.ReplaceAll(text, "{b}", name)
		return text
	}
}
