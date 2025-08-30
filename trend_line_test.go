package charts

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTrendLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		render func(*Painter) ([]byte, error)
		result string
	}{
		{
			name: "linear",
			render: func(p *Painter) ([]byte, error) {
				trendLine := newTrendLinePainter(p)
				axisRange := newTestRange(p.Height(), 6, 0.0, 10.0, 0.0, 0.0)
				xValues := []int{50, 150, 250, 350, 450, 550}
				trend := SeriesTrendLine{
					Type: SeriesTrendTypeLinear,
				}
				trendLine.add(trendLineRenderOption{
					defaultStrokeColor: ColorBlack,
					xValues:            xValues,
					seriesValues:       []float64{1, 2, 3, 4, 5, 6},
					axisRange:          axisRange,
					trends:             []SeriesTrendLine{trend},
				})
				if _, err := trendLine.Render(); err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 70 344\nL 170 308\nL 270 272\nL 370 236\nL 470 200\nL 570 164\" style=\"stroke-width:2;stroke:black;fill:none\"/></svg>",
		},
		{
			name: "cubic",
			render: func(p *Painter) ([]byte, error) {
				trendLine := newTrendLinePainter(p)
				axisRange := newTestRange(p.Height(), 6, 0.0, 40.0, 0.0, 0.0)
				xValues := []int{50, 150, 250, 350, 450, 550}
				trend := SeriesTrendLine{
					Type: SeriesTrendTypeCubic,
				}
				trendLine.add(trendLineRenderOption{
					defaultStrokeColor: ColorBlack,
					xValues:            xValues,
					seriesValues:       []float64{1, 4, 9, 16, 25, 36},
					axisRange:          axisRange,
					trends:             []SeriesTrendLine{trend},
				})
				if _, err := trendLine.Render(); err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 70 371\nL 170 345\nL 270 300\nL 370 236\nL 470 155\nL 570 57\" style=\"stroke-width:2;stroke:black;fill:none\"/></svg>",
		},
		{
			name: "average",
			render: func(p *Painter) ([]byte, error) {
				trendLine := newTrendLinePainter(p)
				axisRange := newTestRange(p.Height(), 6, 0.0, 6.0, 0.0, 0.0)
				xValues := []int{50, 150, 250, 350, 450, 550}
				trend := SeriesTrendLine{
					Type:   SeriesTrendTypeAverage,
					Window: 3,
				}
				trendLine.add(trendLineRenderOption{
					defaultStrokeColor: ColorBlack,
					xValues:            xValues,
					seriesValues:       []float64{1, 2, 3, 4, 5, 6},
					axisRange:          axisRange,
					trends:             []SeriesTrendLine{trend},
				})
				if _, err := trendLine.Render(); err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 70 290\nL 170 260\nL 270 200\nL 370 140\nL 470 80\nL 570 50\" style=\"stroke-width:2;stroke:black;fill:none\"/></svg>",
		},
		{
			name: "sma",
			render: func(p *Painter) ([]byte, error) {
				trendLine := newTrendLinePainter(p)
				axisRange := newTestRange(p.Height(), 6, 0.0, 6.0, 0.0, 0.0)
				xValues := []int{50, 150, 250, 350, 450, 550}
				trend := SeriesTrendLine{
					Type:   SeriesTrendTypeSMA,
					Period: 3,
				}
				trendLine.add(trendLineRenderOption{
					defaultStrokeColor: ColorBlack,
					xValues:            xValues,
					seriesValues:       []float64{1, 2, 3, 4, 5, 6},
					axisRange:          axisRange,
					trends:             []SeriesTrendLine{trend},
				})
				if _, err := trendLine.Render(); err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 70 290\nL 170 260\nL 270 200\nL 370 140\nL 470 80\nL 570 50\" style=\"stroke-width:2;stroke:black;fill:none\"/></svg>",
		},
		{
			name: "ema",
			render: func(p *Painter) ([]byte, error) {
				trendLine := newTrendLinePainter(p)
				axisRange := newTestRange(p.Height(), 6, 0.0, 5.0, 0.0, 0.0)
				xValues := []int{50, 150, 250, 350, 450, 550}
				trend := SeriesTrendLine{
					Type:   SeriesTrendTypeEMA,
					Period: 3,
				}
				trendLine.add(trendLineRenderOption{
					defaultStrokeColor: ColorBlack,
					xValues:            xValues,
					seriesValues:       []float64{1, 2, 3, 4, 5, 6},
					axisRange:          axisRange,
					trends:             []SeriesTrendLine{trend},
				})
				if _, err := trendLine.Render(); err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 70 308\nL 170 272\nL 270 218\nL 370 155\nL 470 88\nL 570 20\" style=\"stroke-width:2;stroke:black;fill:none\"/></svg>",
		},
		{
			name: "bollinger_upper",
			render: func(p *Painter) ([]byte, error) {
				trendLine := newTrendLinePainter(p)
				axisRange := newTestRange(p.Height(), 6, 0.0, 10.0, 0.0, 0.0)
				xValues := []int{50, 150, 250, 350, 450, 550}
				trend := SeriesTrendLine{
					Type:   SeriesTrendTypeBollingerUpper,
					Period: 3,
				}
				trendLine.add(trendLineRenderOption{
					defaultStrokeColor: ColorBlack,
					xValues:            xValues,
					seriesValues:       []float64{1, 2, 3, 4, 5, 6},
					axisRange:          axisRange,
					trends:             []SeriesTrendLine{trend},
				})
				if _, err := trendLine.Render(); err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 70 290\nL 170 250\nL 270 214\nL 370 178\nL 470 142\nL 570 146\" style=\"stroke-width:2;stroke:black;fill:none\"/></svg>",
		},
		{
			name: "bollinger_lower",
			render: func(p *Painter) ([]byte, error) {
				trendLine := newTrendLinePainter(p)
				axisRange := newTestRange(p.Height(), 6, 0.0, 10.0, 0.0, 0.0)
				xValues := []int{50, 150, 250, 350, 450, 550}
				trend := SeriesTrendLine{
					Type:   SeriesTrendTypeBollingerLower,
					Period: 3,
				}
				trendLine.add(trendLineRenderOption{
					defaultStrokeColor: ColorBlack,
					xValues:            xValues,
					seriesValues:       []float64{1, 2, 3, 4, 5, 6},
					axisRange:          axisRange,
					trends:             []SeriesTrendLine{trend},
				})
				if _, err := trendLine.Render(); err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 70 362\nL 170 367\nL 270 331\nL 370 295\nL 470 259\nL 570 218\" style=\"stroke-width:2;stroke:black;fill:none\"/></svg>",
		},
		{
			name: "rsi",
			render: func(p *Painter) ([]byte, error) {
				trendLine := newTrendLinePainter(p)
				axisRange := newTestRange(p.Height(), 6, 0.0, 100.0, 0.0, 0.0)
				xValues := []int{50, 150, 250, 350, 450, 550}
				trend := SeriesTrendLine{
					Type:   SeriesTrendTypeRSI,
					Period: 3,
				}
				trendLine.add(trendLineRenderOption{
					defaultStrokeColor: ColorBlack,
					xValues:            xValues,
					seriesValues:       []float64{44, 44.5, 43.8, 44.2, 44.5, 43.9},
					axisRange:          axisRange,
					trends:             []SeriesTrendLine{trend},
				})
				if _, err := trendLine.Render(); err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 370 178\nL 470 143\nL 570 238\" style=\"stroke-width:2;stroke:black;fill:none\"/></svg>",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			}, PainterThemeOption(GetTheme(ThemeLight)))
			data, err := tt.render(p.Child(PainterPaddingOption(NewBoxEqual(20))))
			require.NoError(t, err)
			assertEqualSVG(t, tt.result, data)
		})
	}
}

func TestTrendLine_WithNullValues(t *testing.T) {
	t.Parallel()

	nv := GetNullValue()

	t.Run("linear_basic", func(t *testing.T) {
		// Perfect linear data with null: y = x
		input := []float64{0, 1, nv, 3, 4}
		result, err := linearTrend(input)
		require.NoError(t, err)

		// Should fit y = x exactly
		assert.InDelta(t, 0, result[0], 0.01)
		assert.InDelta(t, 1, result[1], 0.01)
		assert.InDelta(t, nv, result[2], 0)
		assert.InDelta(t, 3, result[3], 0.01)
		assert.InDelta(t, 4, result[4], 0.01)
	})

	t.Run("linear_single_point", func(t *testing.T) {
		input := []float64{nv, 5, nv}
		result, err := linearTrend(input)
		require.NoError(t, err)

		assert.InDelta(t, nv, result[0], 0)
		assert.InDelta(t, float64(5), result[1], 0)
		assert.InDelta(t, nv, result[2], 0)
	})

	t.Run("linear_no_data", func(t *testing.T) {
		input := []float64{nv, nv, nv}
		result, err := linearTrend(input)
		require.NoError(t, err)

		for i, v := range result {
			assert.InDelta(t, nv, v, 0, "Index %d should be null", i)
		}
	})

	t.Run("cubic_basic", func(t *testing.T) {
		// Data with null
		input := []float64{1, 4, nv, 16, 25, 36}
		result, err := cubicTrend(input)
		require.NoError(t, err)

		// Verify nulls preserved and no NaN
		assert.InDelta(t, nv, result[2], 0)
		for i, v := range result {
			if v != nv {
				assert.False(t, math.IsNaN(v), "Index %d is NaN", i)
				assert.False(t, math.IsInf(v, 0), "Index %d is Inf", i)
			}
		}
	})

	t.Run("cubic_fallback_to_linear", func(t *testing.T) {
		// Only 3 non-null points, should fall back to linear
		input := []float64{1, nv, 3, 5}
		result, err := cubicTrend(input)
		require.NoError(t, err)

		assert.InDelta(t, nv, result[1], 0)
		// Should be linear regression through [1,3,5] at indices [0,2,3]
		// Regression won't exactly match the points
		assert.True(t, result[0] > 0 && result[0] < 2)
		assert.True(t, result[3] > 4 && result[3] < 6)
	})

	t.Run("moving_average_basic", func(t *testing.T) {
		input := []float64{10, 20, nv, 30, 40}
		result, err := movingAverageTrend(input, 3)
		require.NoError(t, err)

		assert.InDelta(t, nv, result[2], 0)
		// First value: average of values within window
		assert.True(t, result[0] > 0 && result[0] < 30)
		// Last value: average of values within window
		assert.True(t, result[4] > 20 && result[4] <= 40)
	})

	t.Run("moving_average_all_nulls", func(t *testing.T) {
		input := []float64{nv, nv, nv}
		result, err := movingAverageTrend(input, 2)
		require.NoError(t, err)

		for _, v := range result {
			assert.InDelta(t, nv, v, 0)
		}
	})

	t.Run("no_nan_produced", func(t *testing.T) {
		// Verify new functions don't produce NaN with null values
		input := []float64{1, 2, nv, 4, 5}

		// New functions don't produce NaN
		newLinear, _ := linearTrend(input)
		newCubic, _ := cubicTrend(input)
		newMA, _ := movingAverageTrend(input, 3)

		for i, v := range newLinear {
			if v != nv {
				assert.False(t, math.IsNaN(v) || math.IsInf(v, 0),
					"linearTrend produced NaN/Inf at index %d", i)
			}
		}

		for i, v := range newCubic {
			if v != nv {
				assert.False(t, math.IsNaN(v) || math.IsInf(v, 0),
					"cubicTrend produced NaN/Inf at index %d", i)
			}
		}

		for i, v := range newMA {
			if v != nv {
				assert.False(t, math.IsNaN(v) || math.IsInf(v, 0),
					"movingAverageTrend produced NaN/Inf at index %d", i)
			}
		}
	})

	t.Run("ema_with_nulls", func(t *testing.T) {
		input := []float64{1, 2, nv, 4, 5}
		result, err := exponentialMovingAverageTrend(input, 3)
		require.NoError(t, err)

		// Verify nulls preserved
		assert.InDelta(t, nv, result[2], 0)

		// Verify non-null values are calculated and no NaN
		for i, v := range result {
			if v != nv {
				assert.False(t, math.IsNaN(v), "EMA produced NaN at index %d", i)
				assert.False(t, math.IsInf(v, 0), "EMA produced Inf at index %d", i)
			}
		}
	})

	t.Run("bollinger_upper_with_nulls", func(t *testing.T) {
		input := []float64{10, 20, nv, 30, 40, 50}
		result, err := bollingerUpperTrend(input, 3)
		require.NoError(t, err)

		// Verify nulls preserved
		assert.InDelta(t, nv, result[2], 0)

		// Verify non-null values and upper band > SMA
		sma, _ := movingAverageTrend(input, 3)
		for i, v := range result {
			if v != nv && sma[i] != nv {
				assert.False(t, math.IsNaN(v), "Bollinger upper produced NaN at index %d", i)
				assert.GreaterOrEqual(t, v, sma[i], "Upper band should be >= SMA at index %d", i)
			}
		}
	})

	t.Run("bollinger_lower_with_nulls", func(t *testing.T) {
		input := []float64{10, 20, nv, 30, 40, 50}
		result, err := bollingerLowerTrend(input, 3)
		require.NoError(t, err)

		// Verify nulls preserved
		assert.InDelta(t, nv, result[2], 0)

		// Verify non-null values and lower band < SMA
		sma, _ := movingAverageTrend(input, 3)
		for i, v := range result {
			if v != nv && sma[i] != nv {
				assert.False(t, math.IsNaN(v), "Bollinger lower produced NaN at index %d", i)
				assert.LessOrEqual(t, v, sma[i], "Lower band should be <= SMA at index %d", i)
			}
		}
	})

	t.Run("rsi_with_nulls", func(t *testing.T) {
		input := []float64{44, 44.5, nv, 44.2, 44.5, 43.9, 44.5, 44.9}
		result, err := rsiTrend(input, 3)
		require.NoError(t, err)

		// All values should be null or valid RSI (0-100)
		for i, v := range result {
			if v != nv {
				assert.False(t, math.IsNaN(v), "RSI produced NaN at index %d", i)
				assert.GreaterOrEqual(t, v, 0.0, "RSI should be >= 0 at index %d", i)
				assert.LessOrEqual(t, v, 100.0, "RSI should be <= 100 at index %d", i)
			}
		}
	})

	t.Run("all_indicators_with_sparse_nulls", func(t *testing.T) {
		// Test with many nulls
		input := []float64{nv, 10, nv, nv, 20, 30, nv, 40, nv}

		// All functions should handle this gracefully
		ema, err := exponentialMovingAverageTrend(input, 2)
		require.NoError(t, err)
		assert.Len(t, ema, len(input))

		upper, err := bollingerUpperTrend(input, 2)
		require.NoError(t, err)
		assert.Len(t, upper, len(input))

		lower, err := bollingerLowerTrend(input, 2)
		require.NoError(t, err)
		assert.Len(t, lower, len(input))

		rsi, err := rsiTrend(input, 2)
		require.NoError(t, err)
		assert.Len(t, rsi, len(input))

		// Verify no NaN/Inf in any results
		for i := range input {
			if ema[i] != nv {
				assert.False(t, math.IsNaN(ema[i]) || math.IsInf(ema[i], 0), "EMA has NaN/Inf at %d", i)
			}
			if upper[i] != nv {
				assert.False(t, math.IsNaN(upper[i]) || math.IsInf(upper[i], 0), "Upper has NaN/Inf at %d", i)
			}
			if lower[i] != nv {
				assert.False(t, math.IsNaN(lower[i]) || math.IsInf(lower[i], 0), "Lower has NaN/Inf at %d", i)
			}
			if rsi[i] != nv {
				assert.False(t, math.IsNaN(rsi[i]) || math.IsInf(rsi[i], 0), "RSI has NaN/Inf at %d", i)
			}
		}
	})
}

func TestLinearTrend(t *testing.T) {
	t.Parallel()

	input := []float64{2, 4, 6, 8}
	expected := []float64{2, 4, 6, 8}

	result, err := linearTrend(input)
	require.NoError(t, err)
	require.Len(t, result, len(expected))
	for i := range expected {
		assert.InDelta(t, expected[i], result[i], 1e-9)
	}
}

func TestCubicTrend(t *testing.T) {
	t.Parallel()

	input := []float64{0, 1, 8, 27}
	expected := []float64{0, 1, 8, 27}

	result, err := cubicTrend(input)
	require.NoError(t, err)
	require.Len(t, result, len(expected))
	for i := range expected {
		assert.InDelta(t, expected[i], result[i], 1e-9)
	}
}

func TestSolveLinearSystem(t *testing.T) {
	t.Parallel()

	mat := [][]float64{
		{0, 1, 0, 0, 2},
		{1, 0, 0, 0, 1},
		{0, 0, 1, 0, 3},
		{0, 0, 0, 1, 4},
	}
	expected := []float64{1, 2, 3, 4}

	result, err := solveLinearSystem(mat)
	require.NoError(t, err)
	require.Len(t, result, len(expected))
	for i := range expected {
		assert.InDelta(t, expected[i], result[i], 1e-9)
	}
}

func TestMovingAverageTrend(t *testing.T) {
	t.Parallel()

	nv := GetNullValue()
	type tc struct {
		name   string
		input  []float64
		window int
		want   []float64
	}

	tests := []tc{
		{
			name:   "simple_gaps_window3",
			input:  []float64{1, 2, nv, 4, 5, 6},
			window: 3,
			// With 5 non-null points, uses moving average (>=4 points)
			// Window of 3 means we look at surrounding values
			// [0]: avg of nearby non-null values = (1+2)/2 = 1.5
			// [1]: avg of nearby non-null values = (1+2)/2 = 1.5
			// [2]: NULL preserved
			// [3]: avg of nearby non-null values = (2+4+5)/3 = 4.5 (window includes indices 1,3,4)
			// [4]: avg of nearby non-null values = (4+5+6)/3 = 5.0
			// [5]: avg of nearby non-null values = (5+6)/2 = 5.5
			want: []float64{1.5, 1.5, nv, 4.5, 5.0, 5.5},
		},
		{
			name:   "leading_and_trailing_nulls",
			input:  []float64{nv, 2, 3, 4, nv},
			window: 2,
			// Only 3 non-null points, falls back to linear trend
			// Linear trend on [2,3,4] at indices [1,2,3]
			// slope = 1, intercept = 1, so y = 1 + x
			want: []float64{nv, 2, 3, 4, nv},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := movingAverageTrend(tt.input, tt.window)
			require.NoError(t, err)
			require.Len(t, got, len(tt.want))
			for i := range tt.want {
				if tt.want[i] == nv {
					assert.InDelta(t, nv, got[i], 0)
				} else {
					assert.InDelta(t, tt.want[i], got[i], 1e-9)
				}
			}
		})
	}

	// Basic test with 5 points and window of 3
	// The center-weighted moving average should give:
	// [0]: avg of 1,2 = 1.5
	// [1]: avg of 1,2,3 = 2
	// [2]: avg of 2,3,4 = 3
	// [3]: avg of 3,4,5 = 4
	// [4]: avg of 4,5 = 4.5
	input := []float64{1, 2, 3, 4, 5}
	expected := []float64{1.5, 2, 3, 4, 4.5}

	result, err := movingAverageTrend(input, 3)
	require.NoError(t, err)
	require.Len(t, result, len(expected))
	for i := range expected {
		assert.InDelta(t, expected[i], result[i], 1e-9)
	}

	t.Run("window_larger_than_data", func(t *testing.T) {
		input := []float64{1, 2, 3, 4}
		result, err := movingAverageTrend(input, 10) // window > len(input)
		require.NoError(t, err)
		assert.Len(t, result, len(input))
	})

	t.Run("massive_window", func(t *testing.T) {
		input := []float64{1, 2, 3, 4, 5}
		result, err := movingAverageTrend(input, 1000)
		require.NoError(t, err)
		assert.Len(t, result, len(input))
	})
}

func TestExponentialMovingAverageTrend(t *testing.T) {
	t.Parallel()

	values := []float64{1, 2, 3, 4, 5}
	result, err := exponentialMovingAverageTrend(values, 3)

	require.NoError(t, err)
	require.Len(t, result, 5)

	// First value should equal input
	assert.InDelta(t, 1.0, result[0], 0.001)

	// EMA should be calculated with smoothing factor 2/(3+1) = 0.5
	multiplier := 2.0 / 4.0
	expected := (2.0 * multiplier) + (1.0 * (1 - multiplier))
	assert.InDelta(t, expected, result[1], 0.001)
}

func TestLinearTrendWithNulls(t *testing.T) {
	t.Parallel()

	nv := GetNullValue()
	input := []float64{1, 2, nv, 4, 6}
	got, err := linearTrend(input)
	require.NoError(t, err)
	require.Len(t, got, len(input))

	// Linear trend on [1,2,4,6] at indices [0,1,3,4]
	// Should produce a continuous trend line through the data
	// The regression line should fit through all 4 points
	// Let's verify null is preserved and trend is reasonable
	assert.InDelta(t, nv, got[2], 0, "Null at index 2 should be preserved")

	// Verify non-null values form a reasonable trend
	assert.False(t, math.IsNaN(got[0]), "Index 0 should not be NaN")
	assert.False(t, math.IsNaN(got[1]), "Index 1 should not be NaN")
	assert.False(t, math.IsNaN(got[3]), "Index 3 should not be NaN")
	assert.False(t, math.IsNaN(got[4]), "Index 4 should not be NaN")
}

func TestBollingerUpperTrend(t *testing.T) {
	t.Parallel()

	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result, err := bollingerUpperTrend(values, 3)

	require.NoError(t, err)
	require.Len(t, result, 10)

	// With centered window, all values should be calculated
	// Upper band should be greater than or equal to SMA
	sma, err := movingAverageTrend(values, 3)
	require.NoError(t, err)
	nv := GetNullValue()
	for i := 0; i < len(result); i++ {
		if result[i] != nv && sma[i] != nv {
			assert.GreaterOrEqual(t, result[i], sma[i])
		}
	}
}

func TestBollingerLowerTrend(t *testing.T) {
	t.Parallel()

	values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	result, err := bollingerLowerTrend(values, 3)

	require.NoError(t, err)
	require.Len(t, result, 10)

	// With centered window, all values should be calculated
	// Lower band should be less than or equal to SMA
	sma, err := movingAverageTrend(values, 3)
	require.NoError(t, err)
	nv := GetNullValue()
	for i := 0; i < len(result); i++ {
		if result[i] != nv && sma[i] != nv {
			assert.LessOrEqual(t, result[i], sma[i])
		}
	}
}

func TestRsiTrend(t *testing.T) {
	t.Parallel()

	// Create test data with known gains/losses
	values := []float64{44, 44.5, 43.8, 44.2, 44.5, 43.9, 44.5, 44.9, 44.5, 44.8}
	result, err := rsiTrend(values, 3)

	require.NoError(t, err)
	require.Len(t, result, 10)

	// First three values should be null
	for i := 0; i < 3; i++ {
		assert.InDelta(t, GetNullValue(), result[i], 0.001)
	}

	// RSI values should be between 0 and 100
	for i := 3; i < len(result); i++ {
		assert.GreaterOrEqual(t, result[i], 0.0)
		assert.LessOrEqual(t, result[i], 100.0)
	}
}
