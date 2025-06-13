package charts

import (
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 70 320\nL 170 290\nL 270 260\nL 370 200\nL 470 140\nL 570 80\" style=\"stroke-width:2;stroke:black;fill:none\"/></svg>",
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

func TestMovingAverageTrend(t *testing.T) {
	t.Parallel()

	input := []float64{1, 2, 3, 4, 5}
	expected := []float64{1, 1.5, 2, 3, 4}

	result, err := movingAverageTrend(input, 3)
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
