package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeBasicRadarChartOption() RadarChartOption {
	values := [][]float64{
		{4200, 3000, 20000, 35000, 50000, 18000},
		{5000, 14000, 28000, 26000, 42000, 21000},
	}
	return RadarChartOption{
		SeriesList: NewSeriesListRadar(values),
		Title: TitleOption{
			Text: "Basic Radar Chart",
		},
		Legend: LegendOption{
			SeriesNames: []string{"Allocated Budget", "Actual Spending"},
		},
		RadarIndicators: NewRadarIndicators([]string{
			"Sales",
			"Administration",
			"Information Technology",
			"Customer Support",
			"Development",
			"Marketing",
		}, []float64{
			6500, 16000, 30000, 38000, 52000, 25000,
		}),
	}
}

func TestNewRadarChartOptionWithData(t *testing.T) {
	t.Parallel()

	opt := NewRadarChartOptionWithData([][]float64{
		{4200, 3000, 20000, 35000, 50000, 18000},
		{5000, 14000, 28000, 26000, 42000, 21000},
	}, []string{
		"Sales",
		"Administration",
		"Information Technology",
		"Customer Support",
		"Development",
		"Marketing",
	}, []float64{
		6500, 16000, 30000, 38000, 52000, 25000,
	})

	assert.Len(t, opt.SeriesList, 2)
	assert.Equal(t, ChartTypeRadar, opt.SeriesList[0].getType())
	assert.Equal(t, defaultPadding, opt.Padding)

	p := NewPainter(PainterOptions{})
	assert.NoError(t, p.RadarChart(opt))
}

func TestRadarChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		themed      bool
		makeOptions func() RadarChartOption
		pngCRC      uint32
	}{
		{
			name:        "basic_themed",
			themed:      true,
			makeOptions: makeBasicRadarChartOption,
			pngCRC:      0xaa7f7dcd,
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
		if tt.themed {
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-painter", func(t *testing.T) {
				p := NewPainter(painterOptions, PainterThemeOption(GetTheme(ThemeVividDark)))
				rp := NewPainter(rasterOptions, PainterThemeOption(GetTheme(ThemeVividDark)))

				validateRadarChartRender(t, p, rp, tt.makeOptions(), tt.pngCRC)
			})
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-options", func(t *testing.T) {
				p := NewPainter(painterOptions)
				rp := NewPainter(rasterOptions)
				opt := tt.makeOptions()
				opt.Theme = GetTheme(ThemeVividDark)

				validateRadarChartRender(t, p, rp, opt, tt.pngCRC)
			})
		} else {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				p := NewPainter(painterOptions)
				rp := NewPainter(rasterOptions)

				validateRadarChartRender(t, p.Child(PainterPaddingOption(NewBoxEqual(20))),
					rp.Child(PainterPaddingOption(NewBoxEqual(20))), tt.makeOptions(), tt.pngCRC)
			})
		}
	}
}

func validateRadarChartRender(t *testing.T, svgP, pngP *Painter, opt RadarChartOption, expectedCRC uint32) {
	t.Helper()

	err := svgP.RadarChart(opt)
	require.NoError(t, err)
	data, err := svgP.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, data)

	err = pngP.RadarChart(opt)
	require.NoError(t, err)
	rasterData, err := pngP.Bytes()
	require.NoError(t, err)
	assertEqualPNGCRC(t, expectedCRC, rasterData)
}

func TestRadarChartError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		makeOptions      func() RadarChartOption
		errorMsgContains string
	}{
		{
			name: "empty_series",
			makeOptions: func() RadarChartOption {
				return NewRadarChartOptionWithData([][]float64{}, []string{"foo", "bar", "foobar"}, []float64{1, 2, 3})
			},
			errorMsgContains: "empty series list",
		},
		{
			name: "too_few_indicators",
			makeOptions: func() RadarChartOption {
				return NewRadarChartOptionWithData([][]float64{{0.0}}, []string{"foo", "bar"}, []float64{1, 2})
			},
			errorMsgContains: "indicator count",
		},
		{
			name: "indicator_name_value_mismatch",
			makeOptions: func() RadarChartOption {
				return NewRadarChartOptionWithData([][]float64{{1, 2, 3}}, []string{"foo", "bar"}, []float64{1, 2, 3})
			},
			errorMsgContains: "indicator count",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			})

			err := p.RadarChart(tt.makeOptions())
			require.Error(t, err)
			require.ErrorContains(t, err, tt.errorMsgContains)
		})
	}
}
