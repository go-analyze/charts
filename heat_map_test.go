package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHeatMapOptionWithData(t *testing.T) {
	data := [][]float64{
		{10, 20},
		{30, 40},
	}
	opt := NewHeatMapOptionWithData(data)
	require.Equal(t, data, opt.Values)

	p := NewPainter(PainterOptions{
		OutputFormat: ChartOutputSVG,
		Width:        600,
		Height:       400,
	})
	err := p.HeatMapChart(opt)
	require.NoError(t, err)
}

func makeBasicHeatMapOption() HeatMapOption {
	return HeatMapOption{
		Title: TitleOption{Text: "Heat Map"},
		Values: [][]float64{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		},
		XAxis: HeatMapAxis{
			Title:  "X-Axis",
			Labels: []string{"A", "B", "C"},
		},
		YAxis: HeatMapAxis{
			Title:  "Y-Axis",
			Labels: []string{"Row1", "Row2", "Row3"},
		},
	}
}

func makeMinimalHeatMapOption() HeatMapOption {
	return NewHeatMapOptionWithData([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})
}

func makeDenseHeatMapOption() HeatMapOption {
	const size = 24
	values := make([][]float64, size)
	for i := range values {
		values[i] = make([]float64, size)
	}
	// create a grid of varying intensity
	for _, index := range []int{6, 12, 18} {
		for i := 0; i < size; i++ {
			inc := 1.0
			if i%2 == 0 {
				inc *= 2
			}
			if index == 12 { // make center line more intense
				inc *= 2
			}
			values[index][i] += inc
			values[i][index] += inc
		}
	}

	opt := NewHeatMapOptionWithData(values)
	opt.BaseColorIndex = 1
	opt.XAxis.LabelCount = size
	opt.YAxis.LabelCount = size
	return opt
}

func TestHeatMapChart(t *testing.T) {
	tests := []struct {
		name        string
		themed      bool
		makeOptions func() HeatMapOption
		pngCRC      uint32
	}{
		{
			name:        "basic_themed",
			themed:      true,
			makeOptions: makeBasicHeatMapOption,
			pngCRC:      0xce5f0709,
		},
		{
			name: "scale_override",
			makeOptions: func() HeatMapOption {
				opt := makeMinimalHeatMapOption()
				minVal, maxVal := 0.0, 20.0
				opt.ScaleMinValue = &minVal
				opt.ScaleMaxValue = &maxVal
				return opt
			},
			pngCRC: 0x7abb75db,
		},
		{
			name: "values_label",
			makeOptions: func() HeatMapOption {
				opt := makeMinimalHeatMapOption()
				opt.ValuesLabel = SeriesLabel{
					Show: Ptr(true),
					FontStyle: FontStyle{
						FontSize:  14,
						FontColor: ColorBlue,
					},
					ValueFormatter: func(f float64) string {
						return strconv.FormatFloat(f, 'f', 0, 64)
					},
				}
				return opt
			},
			pngCRC: 0x3c6cead7,
		},
		{
			name: "varying_row_lengths",
			makeOptions: func() HeatMapOption {
				return NewHeatMapOptionWithData([][]float64{
					{1, 2, 3, 4},
					{5, 6},
					{7, 8, 9},
					nil,
				})
			},
			pngCRC: 0x6bc4178,
		},
		{
			name:        "dense_data",
			makeOptions: makeDenseHeatMapOption,
			pngCRC:      0x314e2a7a,
		},
	}

	for i, tt := range tests {
		painterOptions := PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        600,
			Height:       400,
		}
		rasterOptions := PainterOptions{
			OutputFormat: ChartOutputPNG,
			Width:        600,
			Height:       400,
		}
		if !tt.themed {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				p := NewPainter(painterOptions)
				rp := NewPainter(rasterOptions)
				validateHeatMapChartRender(t, p, rp, tt.makeOptions(), tt.pngCRC)
			})
		} else {
			theme := GetTheme(ThemeVividDark)
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-theme_painter", func(t *testing.T) {
				p := NewPainter(painterOptions, PainterThemeOption(theme))
				rp := NewPainter(rasterOptions, PainterThemeOption(theme))
				validateHeatMapChartRender(t, p, rp, tt.makeOptions(), tt.pngCRC)
			})
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-theme_opt", func(t *testing.T) {
				p := NewPainter(painterOptions)
				rp := NewPainter(rasterOptions)
				opt := tt.makeOptions()
				opt.Theme = theme
				validateHeatMapChartRender(t, p, rp, opt, tt.pngCRC)
			})
		}
	}
}

func validateHeatMapChartRender(t *testing.T, svgP, pngP *Painter, opt HeatMapOption, expectedCRC uint32) {
	t.Helper()

	err := svgP.HeatMapChart(opt)
	require.NoError(t, err)
	data, err := svgP.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, data)

	err = pngP.HeatMapChart(opt)
	require.NoError(t, err)
	rasterData, err := pngP.Bytes()
	require.NoError(t, err)
	assertEqualPNGCRC(t, expectedCRC, rasterData)
}

func TestHeatMapChartError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		makeOptions      func() HeatMapOption
		errorMsgContains string
	}{
		{
			name: "empty_values",
			makeOptions: func() HeatMapOption {
				return HeatMapOption{
					Values: [][]float64{},
				}
			},
			errorMsgContains: "empty values",
		},
		{
			name: "no_columns",
			makeOptions: func() HeatMapOption {
				return HeatMapOption{
					Values: [][]float64{{}, {}},
				}
			},
			errorMsgContains: "heat map has no columns",
		},
		{
			name: "insufficient_space",
			makeOptions: func() HeatMapOption {
				return HeatMapOption{
					Padding: NewBoxEqual(2),
					Values: [][]float64{
						{1, 2, 3},
						{4, 5, 6},
						{7, 8, 9},
					},
					XAxis: HeatMapAxis{
						Title:          "X-Axis",
						Labels:         []string{"A", "B", "C"},
						LabelFontStyle: FontStyle{FontSize: 10, Font: GetDefaultFont(), FontColor: ColorBlack},
					},
					YAxis: HeatMapAxis{
						Title:          "Y-Axis",
						Labels:         []string{"Row1", "Row2", "Row3"},
						LabelFontStyle: FontStyle{FontSize: 10, Font: GetDefaultFont(), FontColor: ColorBlack},
					},
				}
			},
			errorMsgContains: "insufficient space for heat map cells",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			var p *Painter
			if tt.name == "insufficient_space" {
				p = NewPainter(PainterOptions{
					Width:  40,
					Height: 40,
				})
			} else {
				p = NewPainter(PainterOptions{
					Width:  600,
					Height: 400,
				})
			}
			err := p.HeatMapChart(tt.makeOptions())
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.errorMsgContains)
		})
	}
}

func TestComputeMinMax(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		values      [][]float64
		numCol      int
		expectedMin float64
		expectedMax float64
	}{
		{
			name:        "empty",
			values:      [][]float64{},
			numCol:      0,
			expectedMin: 0,
			expectedMax: 0,
		},
		{
			name: "uneven_rows",
			values: [][]float64{
				{1, 2, 3},
				{4},
			},
			numCol:      3,
			expectedMin: 1,
			expectedMax: 4,
		},
		{
			name: "negative_values",
			values: [][]float64{
				{-5, -2},
				{-3, -1},
			},
			numCol:      2,
			expectedMin: -5,
			expectedMax: -1,
		},
		{
			name: "default_column_padding",
			values: [][]float64{
				{-1},
				{},
			},
			numCol:      1,
			expectedMin: 0,
			expectedMax: 0,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			min, max := computeMinMax(tt.values, tt.numCol)

			assert.InDelta(t, tt.expectedMin, min, 0.0)
			assert.InDelta(t, tt.expectedMax, max, 0.0)
		})
	}
}
