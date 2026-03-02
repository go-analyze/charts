package charts

import (
	"errors"
	"math"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/go-analyze/charts/chartdraw/matrix"
)

const (
	// TODO - v0.6 - Move these constants to types, or builder pattern on structs similar to CandlestickPatternConfig

	// SeriesTrendTypeLinear represents a linear regression trend line that fits a straight line through the data points.
	SeriesTrendTypeLinear = "linear"
	// SeriesTrendTypeCubic represents a cubic polynomial (degree 3) regression trend line that fits a curved line through the data points.
	SeriesTrendTypeCubic = "cubic"
	// Deprecated: SeriesTrendTypeAverage is deprecated, use SeriesTrendTypeSMA instead.
	SeriesTrendTypeAverage = "average"
	// SeriesTrendTypeSMA represents a Simple Moving Average trend line that smooths data using a sliding window average.
	SeriesTrendTypeSMA = "sma"
	// SeriesTrendTypeEMA represents an Exponential Moving Average trend line that gives more weight to recent data points.
	SeriesTrendTypeEMA = "ema"
	// SeriesTrendTypeBollingerUpper represents the upper Bollinger Band (SMA + 2 * standard deviation).
	// Designed for financial time-series analysis to identify volatility boundaries around price movements.
	SeriesTrendTypeBollingerUpper = "bollinger_upper"
	// SeriesTrendTypeBollingerLower represents the lower Bollinger Band (SMA - 2 * standard deviation).
	// Designed for financial time-series analysis to identify volatility boundaries around price movements.
	SeriesTrendTypeBollingerLower = "bollinger_lower"
	// SeriesTrendTypeRSI represents the Relative Strength Index momentum oscillator (0-100 scale).
	// Measures momentum by analyzing sequential price changes, designed for financial time-series analysis.
	SeriesTrendTypeRSI = "rsi"
)

// SeriesTrendLine describes the rendered trend line style.
type SeriesTrendLine struct {
	// LineStrokeWidth is the width of the rendered line.
	LineStrokeWidth float64
	// StrokeSmoothingTension should be between 0 and 1. At 0 lines are sharp and precise, 1 provides smoother lines.
	StrokeSmoothingTension float64
	// LineColor overrides the theme color for this trend line.
	LineColor Color
	// DashedLine indicates if the trend line will be a dashed line. Default depends on chart type.
	DashedLine *bool
	// Type specifies the trend line type: "linear", "cubic", "sma", "ema", "rsi".
	Type string
	// Deprecated: Window is deprecated, use Period instead.
	Window int
	// Period specifies the number of data points to consider for trend calculations.
	// Used by moving averages (SMA, EMA), Bollinger Bands, RSI, and other indicators.
	// For example, Period=20 calculates a 20-period moving average.
	Period int
}

// NewTrendLine returns a trend line for the provided type. Set on a specific Series instance.
func NewTrendLine(trendType string) []SeriesTrendLine {
	return []SeriesTrendLine{
		{
			Type: trendType,
		},
	}
}

// trendLinePainter is responsible for rendering trend lines on the chart.
type trendLinePainter struct {
	p       *Painter
	options []trendLineRenderOption
}

// newTrendLinePainter returns a new trend line renderer.
func newTrendLinePainter(p *Painter) *trendLinePainter {
	return &trendLinePainter{
		p: p,
	}
}

// add appends a trend line render option.
func (t *trendLinePainter) add(opt trendLineRenderOption) {
	t.options = append(t.options, opt)
}

// trendLineRenderOption holds configuration for rendering trend lines.
type trendLineRenderOption struct {
	defaultStrokeColor Color
	// xValues are the x-coordinates for each data sample.
	xValues []int
	// seriesValues are the raw data values.
	seriesValues []float64
	// axisRange is used to transform a raw data value into a screen y-coordinate.
	axisRange axisRange
	// trends are the list of trend lines to render for this series.
	trends []SeriesTrendLine
	// dashed indicates if the trend line will be a dashed line.
	dashed bool
}

// Render computes and draws all configured trend lines.
func (t *trendLinePainter) Render() (Box, error) {
	painter := t.p
	for _, opt := range t.options {
		if len(opt.trends) == 0 || len(opt.seriesValues) == 0 || len(opt.xValues) == 0 {
			continue
		}

		for _, trend := range opt.trends {
			if trend.Window != 0 && trend.Period == 0 {
				trend.Period = trend.Window
			}
			var fitted []float64
			var err error
			switch trend.Type {
			case SeriesTrendTypeLinear:
				fitted, err = linearTrend(opt.seriesValues)
			case SeriesTrendTypeCubic:
				fitted, err = cubicTrend(opt.seriesValues)
			case SeriesTrendTypeSMA, "average" /* long term backwards compatibility */ :
				fitted, err = movingAverageTrend(opt.seriesValues, trend.Period)
			case SeriesTrendTypeEMA:
				fitted, err = exponentialMovingAverageTrend(opt.seriesValues, trend.Period)
			case SeriesTrendTypeBollingerUpper:
				fitted, err = bollingerUpperTrend(opt.seriesValues, trend.Period)
			case SeriesTrendTypeBollingerLower:
				fitted, err = bollingerLowerTrend(opt.seriesValues, trend.Period)
			case SeriesTrendTypeRSI:
				fitted, err = rsiTrend(opt.seriesValues, trend.Period)
			default:
				err = errors.New("unknown trend type: " + trend.Type)
			}
			if err != nil {
				return BoxZero, err
			} else if len(fitted) != len(opt.xValues) {
				return BoxZero, errors.New("mismatched data length in trend line computation")
			}

			color := trend.LineColor
			if color.IsTransparent() {
				color = opt.defaultStrokeColor
			}
			strokeWidth := trend.LineStrokeWidth
			if strokeWidth == 0 {
				strokeWidth = defaultStrokeWidth
			}

			// Convert fitted data to screen points, break where fitted is null.
			points := make([]Point, len(fitted))
			for i, val := range fitted {
				if isValidExtent(val) {
					points[i] = Point{X: opt.xValues[i], Y: opt.axisRange.getRestHeight(val)}
				} else {
					points[i] = Point{X: opt.xValues[i], Y: math.MaxInt32}
				}
			}

			// Determine if this trend line should be dashed
			isDashed := opt.dashed // start with chart default
			if trend.DashedLine != nil {
				isDashed = *trend.DashedLine
			}

			if isDashed {
				// Calculate dash size based on painter dimensions for better visibility
				avgDimension := float64(t.p.box.Width()+t.p.box.Height()) / 2
				dashLength := math.Max(avgDimension*0.02, 4.0) // Minimum 4px, scale with size
				gapLength := dashLength * 0.8
				dashArray := []float64{dashLength, gapLength}
				if trend.StrokeSmoothingTension > 0 {
					painter.SmoothDashedLineStroke(points, trend.StrokeSmoothingTension, color, strokeWidth, dashArray)
				} else {
					painter.DashedLineStroke(points, color, strokeWidth, dashArray)
				}
			} else {
				if trend.StrokeSmoothingTension > 0 {
					painter.SmoothLineStroke(points, trend.StrokeSmoothingTension, color, strokeWidth)
				} else {
					painter.LineStroke(points, color, strokeWidth)
				}
			}
		}
	}
	return BoxZero, nil
}

// extractNonNullData extracts non-null values and their indices from the input.
func extractNonNullData(y []float64) ([]float64, []int) {
	cleanData := make([]float64, 0, len(y))
	cleanIndices := make([]int, 0, len(y))
	for i, v := range y {
		if isValidExtent(v) {
			cleanData = append(cleanData, v)
			cleanIndices = append(cleanIndices, i)
		}
	}
	return cleanData, cleanIndices
}

// initResultWithNulls creates a result array preserving null positions from the input.
func initResultWithNulls(y []float64) []float64 {
	result := make([]float64, len(y))
	for i, v := range y {
		if !isValidExtent(v) {
			result[i] = GetNullValue()
		}
	}
	return result
}

// linearTrend computes a linear trend over the data, preserving null positions.
func linearTrend(y []float64) ([]float64, error) {
	cleanData, cleanIndices := extractNonNullData(y)
	result := initResultWithNulls(y)

	if len(cleanData) == 0 {
		return result, nil // All nulls
	} else if len(cleanData) == 1 {
		result[cleanIndices[0]] = cleanData[0] // Single point - just preserve it
		return result, nil
	}

	return computeLinearTrend(result, y, cleanData, cleanIndices)
}

func computeLinearTrend(result, data, cleanData []float64, cleanIndices []int) ([]float64, error) {
	n := float64(len(cleanData))
	var sumX, sumY, sumXY, sumXX float64
	for i, v := range cleanData {
		x := float64(cleanIndices[i])
		sumX += x
		sumY += v
		sumXY += x * v
		sumXX += x * x
	}

	denom := n*sumXX - sumX*sumX
	if math.Abs(denom) < 1e-10 {
		return nil, errors.New("singular matrix in linear regression")
	}
	slope := (n*sumXY - sumX*sumY) / denom
	intercept := (sumY - slope*sumX) / n

	for i, v := range data {
		if isValidExtent(v) { // Apply trend to non-null positions only
			result[i] = intercept + slope*float64(i)
		}
	}

	return result, nil
}

// cubicTrend computes a cubic polynomial trend over the data, preserving null positions.
func cubicTrend(y []float64) ([]float64, error) {
	cleanData, cleanIndices := extractNonNullData(y)
	n := len(cleanData)
	result := initResultWithNulls(y)

	if n == 0 {
		return result, nil // All nulls
	} else if n == 1 {
		result[cleanIndices[0]] = cleanData[0] // Single point - just preserve it
		return result, nil
	} else if n < 4 {
		return computeLinearTrend(result, y, cleanData, cleanIndices) // Fall back to linear for less than 4 points
	}

	// Compute sums of powers of x using original indices
	var S [7]float64
	for i := 0; i < n; i++ {
		x := float64(cleanIndices[i])
		xp := 1.0
		for k := 0; k <= 6; k++ {
			S[k] += xp
			xp *= x
		}
	}

	// Compute the right-hand side vector B
	var B [4]float64
	for i := 0; i < n; i++ {
		x := float64(cleanIndices[i])
		xp := 1.0
		for j := 0; j < 4; j++ {
			B[j] += cleanData[i] * xp
			xp *= x
		}
	}

	// Build the augmented matrix
	M := make([][]float64, 4)
	for j := 0; j < 4; j++ {
		M[j] = make([]float64, 5)
		for k := 0; k < 4; k++ {
			M[j][k] = S[j+k]
		}
		M[j][4] = B[j]
	}

	coeffs, err := solveLinearSystem(M)
	if err != nil {
		return linearTrend(y) // Fall back to linear
	}

	// Apply cubic polynomial to non-null positions only
	for i, v := range y {
		if isValidExtent(v) {
			x := float64(i)
			result[i] = coeffs[0] + coeffs[1]*x + coeffs[2]*x*x + coeffs[3]*x*x*x
		}
	}

	return result, nil
}

// exponentialMovingAverageTrend computes an exponential moving average over the data, preserving null positions.
// If window is <= 0, a default based on the data size is used.
func exponentialMovingAverageTrend(y []float64, window int) ([]float64, error) {
	cleanData, cleanIndices := extractNonNullData(y)
	nonNullCount := len(cleanData)
	result := initResultWithNulls(y)

	if nonNullCount == 0 {
		return result, nil // All nulls
	} else if nonNullCount == 1 {
		result[cleanIndices[0]] = cleanData[0] // Single point - just preserve it
		return result, nil
	} else if nonNullCount < 4 {
		return computeLinearTrend(result, y, cleanData, cleanIndices) // Fall back to linear for less than 4 points
	}

	if window <= 0 {
		window = chartdraw.MaxInt(2, nonNullCount/5)
	}

	multiplier := 2.0 / (float64(window) + 1.0)

	// Calculate EMA only for non-null positions
	var ema float64
	isFirst := true
	for i, v := range y {
		if !isValidExtent(v) {
			continue
		}

		if isFirst {
			// First non-null value initializes EMA
			ema = v
			result[i] = ema
			isFirst = false
		} else {
			// Update EMA with current value
			ema = (v * multiplier) + (ema * (1 - multiplier))
			result[i] = ema
		}
	}

	return result, nil
}

// solveLinearSystem solves a 4x4 linear system represented as an augmented matrix.
// The input matrix has 4 rows and 5 columns (last column is the constants vector).
func solveLinearSystem(mat [][]float64) ([]float64, error) {
	n := len(mat)
	// Forward elimination
	for i := 0; i < n; i++ {
		// Find the pivot row
		maxRow := i
		for j := i + 1; j < n; j++ {
			if math.Abs(mat[j][i]) > math.Abs(mat[maxRow][i]) {
				maxRow = j
			}
		}
		mat[i], mat[maxRow] = mat[maxRow], mat[i]
		if math.Abs(mat[i][i]) < matrix.DefaultEpsilon {
			return nil, errors.New("singular matrix in cubic regression")
		}
		// Eliminate below
		for j := i + 1; j < n; j++ {
			factor := mat[j][i] / mat[i][i]
			for k := i; k <= n; k++ {
				mat[j][k] -= factor * mat[i][k]
			}
		}
	}
	// Back substitution
	sol := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		sol[i] = mat[i][n]
		for j := i + 1; j < n; j++ {
			sol[i] -= mat[i][j] * sol[j]
		}
		sol[i] /= mat[i][i]
	}
	return sol, nil
}

// movingAverageTrend computes a moving average over the data, preserving null positions.
func movingAverageTrend(y []float64, window int) ([]float64, error) {
	cleanData, cleanIndices := extractNonNullData(y)
	nonNullCount := len(cleanData)
	result := initResultWithNulls(y)

	if nonNullCount == 0 {
		return result, nil // All nulls
	} else if nonNullCount == 1 {
		result[cleanIndices[0]] = cleanData[0] // Single point - just preserve it
		return result, nil
	} else if nonNullCount < 4 {
		return computeLinearTrend(result, y, cleanData, cleanIndices) // Fall back to linear for less than 4 points
	}

	if window <= 0 {
		window = chartdraw.MaxInt(2, nonNullCount/5)
	}

	// Compute moving average for non-null positions
	halfWindow := window / 2
	for i, v := range y {
		if !isValidExtent(v) {
			continue
		}

		// Calculate average of surrounding non-null values within window
		var sum float64
		var count int
		start := chartdraw.MaxInt(0, i-halfWindow)
		end := chartdraw.MinInt(len(y)-1, i+halfWindow)
		for j := start; j <= end; j++ {
			if isValidExtent(y[j]) {
				sum += y[j]
				count++
			}
		}

		if count > 0 {
			result[i] = sum / float64(count)
		} // else, shouldn't happen
	}

	return result, nil
}

// bollingerBand computes a Bollinger Band (SMA ± multiplier * standard deviation).
func bollingerBand(y []float64, period int, multiplier float64) ([]float64, error) {
	cleanData, _ := extractNonNullData(y)
	nonNullCount := len(cleanData)
	result := initResultWithNulls(y)

	if nonNullCount < 2 {
		return result, nil // Not enough data
	}
	if period <= 0 {
		period = chartdraw.MaxInt(2, nonNullCount/5)
	}
	if period > nonNullCount {
		return result, nil // Period too large
	}

	// Calculate SMA first (already handles nulls)
	sma, err := movingAverageTrend(y, period)
	if err != nil {
		return nil, err
	}

	// Compute Bollinger bands with centered window
	halfWindow := period / 2
	for i, v := range y {
		if !isValidExtent(v) || !isValidExtent(sma[i]) {
			continue
		}

		// Calculate standard deviation for centered window
		mean := sma[i]
		var variance float64
		var count int
		start := chartdraw.MaxInt(0, i-halfWindow)
		end := chartdraw.MinInt(len(y)-1, i+halfWindow)

		for j := start; j <= end; j++ {
			if isValidExtent(y[j]) {
				diff := y[j] - mean
				variance += diff * diff
				count++
			}
		}

		if count > 0 {
			stddev := math.Sqrt(variance / float64(count))
			result[i] = mean + (stddev * multiplier)
		}
	}

	return result, nil
}

// bollingerUpperTrend computes the upper Bollinger Band (SMA + 2 * standard deviation), preserving null positions.
func bollingerUpperTrend(y []float64, period int) ([]float64, error) {
	return bollingerBand(y, period, 2.0)
}

// bollingerLowerTrend computes the lower Bollinger Band (SMA - 2 * standard deviation), preserving null positions.
func bollingerLowerTrend(y []float64, period int) ([]float64, error) {
	return bollingerBand(y, period, -2.0)
}

// rsiTrend computes the Relative Strength Index momentum oscillator, preserving null positions.
func rsiTrend(y []float64, period int) ([]float64, error) {
	cleanData, cleanIndices := extractNonNullData(y)
	result := initResultWithNulls(y)
	for i := 0; i < period && i < len(result); i++ {
		result[i] = GetNullValue() // set start up to period as null since it can't be calculated
	}

	if len(cleanData) < 2 {
		return result, nil // Not enough non-null data
	}
	if period <= 0 {
		period = chartdraw.MaxInt(2, len(cleanData)/5)
	}
	if len(cleanData) < period+1 {
		return result, nil // Insufficient data for RSI
	}

	// Calculate price changes between consecutive non-null values
	gains := make([]float64, len(cleanData)-1)
	losses := make([]float64, len(cleanData)-1)
	for i := 1; i < len(cleanData); i++ {
		change := cleanData[i] - cleanData[i-1]
		if change > 0 {
			gains[i-1] = change
			losses[i-1] = 0
		} else {
			gains[i-1] = 0
			losses[i-1] = -change
		}
	}

	// Calculate initial averages
	var avgGain, avgLoss float64
	for i := 0; i < period && i < len(gains); i++ {
		avgGain += gains[i]
		avgLoss += losses[i]
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	// Calculate RSI for non-null positions
	for i := period; i < len(cleanData); i++ {
		idx := cleanIndices[i]
		if avgLoss == 0 {
			result[idx] = 100
		} else {
			rs := avgGain / avgLoss
			result[idx] = 100 - (100 / (1 + rs))
		}

		// Update averages for next iteration
		if i < len(gains) {
			avgGain = ((avgGain * float64(period-1)) + gains[i]) / float64(period)
			avgLoss = ((avgLoss * float64(period-1)) + losses[i]) / float64(period)
		}
	}

	return result, nil
}
