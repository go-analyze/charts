package charts

import (
	"errors"
	"math"

	"github.com/go-analyze/charts/chartdraw"
	"github.com/go-analyze/charts/chartdraw/matrix"
)

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
}

// Render computes and draws all configured trend lines.
func (t *trendLinePainter) Render() (Box, error) {
	painter := t.p
	for _, opt := range t.options {
		if len(opt.trends) == 0 || len(opt.seriesValues) == 0 || len(opt.xValues) == 0 {
			continue
		}

		for _, trend := range opt.trends {
			var fitted []float64
			var err error
			switch trend.Type {
			case SeriesTrendTypeLinear:
				fitted, err = linearTrend(opt.seriesValues)
			case SeriesTrendTypeCubic:
				fitted, err = cubicTrend(opt.seriesValues)
			case SeriesTrendTypeAverage:
				fitted, err = movingAverageTrend(opt.seriesValues, trend.Window)
			default:
				// Unknown trend type; skip.
				continue
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

			// Convert fitted data to screen points.
			points := make([]Point, len(fitted))
			for i, val := range fitted {
				points[i] = Point{
					X: opt.xValues[i],
					Y: opt.axisRange.getRestHeight(val),
				}
			}

			if trend.StrokeSmoothingTension > 0 {
				painter.SmoothLineStroke(points, trend.StrokeSmoothingTension, color, strokeWidth)
			} else {
				painter.LineStroke(points, color, strokeWidth)
			}
		}
	}
	return BoxZero, nil
}

// linearTrend computes a linear regression over the provided data.
func linearTrend(y []float64) ([]float64, error) {
	n := float64(len(y))
	if n < 2 {
		return nil, errors.New("not enough data points for linear trend")
	}

	var sumX, sumY, sumXY, sumXX float64
	for i, v := range y {
		x := float64(i)
		sumX += x
		sumY += v
		sumXY += x * v
		sumXX += x * x
	}

	denom := n*sumXX - sumX*sumX
	if math.Abs(denom) < matrix.DefaultEpsilon {
		return nil, errors.New("degenerate x values for linear regression")
	}
	slope := (n*sumXY - sumX*sumY) / denom
	intercept := (sumY - slope*sumX) / n

	fitted := make([]float64, len(y))
	for i := range y {
		fitted[i] = intercept + slope*float64(i)
	}
	return fitted, nil
}

// cubicTrend computes a cubic (degree 3) polynomial regression over the data.
// If there are fewer than 4 points, it falls back to a linear trend.
func cubicTrend(y []float64) ([]float64, error) {
	n := len(y)
	if n < 2 {
		return nil, errors.New("not enough data points for cubic trend")
	} else if n < 4 {
		return linearTrend(y)
	}

	// Compute sums of powers of x.
	var S [7]float64 // S[k] = Σ x^k for k = 0..6.
	for i := 0; i < n; i++ {
		x := float64(i)
		xp := 1.0
		for k := 0; k <= 6; k++ {
			S[k] += xp
			xp *= x
		}
	}

	// Compute the right-hand side vector B.
	var B [4]float64
	for i := 0; i < n; i++ {
		x := float64(i)
		xp := 1.0
		for j := 0; j < 4; j++ {
			B[j] += y[i] * xp
			xp *= x
		}
	}

	// Build the augmented matrix for the normal equations.
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
		// fallback to linear
		return linearTrend(y)
	}

	fitted := make([]float64, n)
	for i := 0; i < n; i++ {
		x := float64(i)
		fitted[i] = coeffs[0] + coeffs[1]*x + coeffs[2]*x*x + coeffs[3]*x*x*x
	}
	return fitted, nil
}

// movingAverageTrend computes a moving average over the data using the given window size.
// If window is <= 0, a default based on the data size is used.
func movingAverageTrend(y []float64, window int) ([]float64, error) {
	n := len(y)
	if n < 2 {
		return nil, errors.New("not enough data points for average trend")
	} else if n < 4 {
		return linearTrend(y)
	}
	if window <= 0 {
		window = chartdraw.MaxInt(2, n/5)
	}

	fitted := make([]float64, n)
	var sum float64
	for i := 0; i < n; i++ {
		sum += y[i]
		if i >= window {
			sum -= y[i-window]
			fitted[i] = sum / float64(window)
		} else {
			fitted[i] = sum / float64(i+1)
		}
	}
	return fitted, nil
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
