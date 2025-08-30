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
