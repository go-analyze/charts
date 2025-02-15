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
		result      string
	}{
		{
			name:        "basic_themed",
			themed:      true,
			makeOptions: makeBasicRadarChartOption,
			result:      "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:rgb(40,40,40)\"/><path  d=\"M 144 3\nL 174 3\nL 174 16\nL 144 16\nL 144 3\" style=\"stroke:none;fill:rgb(255,100,100)\"/><text x=\"176\" y=\"15\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Allocated Budget</text><path  d=\"M 314 3\nL 344 3\nL 344 16\nL 314 16\nL 314 3\" style=\"stroke:none;fill:rgb(255,210,100)\"/><text x=\"346\" y=\"15\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Actual Spending</text><text x=\"0\" y=\"15\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Basic Radar Chart</text><path  d=\"M 300 176\nL 322 189\nL 322 214\nL 300 228\nL 278 215\nL 278 190\nL 300 176\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 150\nL 345 176\nL 345 227\nL 300 254\nL 255 228\nL 255 177\nL 300 150\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 124\nL 367 163\nL 367 240\nL 300 280\nL 233 241\nL 233 164\nL 300 124\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 98\nL 390 150\nL 390 253\nL 300 306\nL 210 254\nL 210 151\nL 300 98\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 72\nL 412 137\nL 412 266\nL 300 332\nL 188 267\nL 188 138\nL 300 72\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 202\nL 300 72\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 202\nL 412 137\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 202\nL 412 266\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 202\nL 300 332\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 202\nL 188 267\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 202\nL 188 138\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><text x=\"284\" y=\"65\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Sales</text><text x=\"417\" y=\"142\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Administration</text><text x=\"417\" y=\"271\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Information Technology</text><text x=\"248\" y=\"349\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Customer Support</text><text x=\"107\" y=\"272\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Development</text><text x=\"124\" y=\"143\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Marketing</text><path  d=\"M 300 118\nL 321 190\nL 375 245\nL 300 321\nL 192 264\nL 219 156\nL 300 118\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:none\"/><path  d=\"M 300 118\nL 321 190\nL 375 245\nL 300 321\nL 192 264\nL 219 156\nL 300 118\" style=\"stroke:none;fill:rgba(255,100,100,0.1)\"/><circle cx=\"300\" cy=\"118\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"321\" cy=\"190\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"375\" cy=\"245\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"300\" cy=\"321\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"192\" cy=\"264\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"219\" cy=\"156\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"300\" cy=\"118\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><path  d=\"M 300 102\nL 398 146\nL 405 262\nL 300 290\nL 210 254\nL 206 148\nL 300 102\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:none\"/><path  d=\"M 300 102\nL 398 146\nL 405 262\nL 300 290\nL 210 254\nL 206 148\nL 300 102\" style=\"stroke:none;fill:rgba(255,210,100,0.1)\"/><circle cx=\"300\" cy=\"102\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"398\" cy=\"146\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"405\" cy=\"262\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"300\" cy=\"290\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"210\" cy=\"254\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"206\" cy=\"148\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"300\" cy=\"102\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/></svg>",
		},
	}

	for i, tt := range tests {
		painterOptions := PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        600,
			Height:       400,
		}
		if tt.themed {
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-painter", func(t *testing.T) {
				p := NewPainter(painterOptions, PainterThemeOption(GetTheme(ThemeVividDark)))

				validateRadarChartRender(t, p, tt.makeOptions(), tt.result)
			})
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-options", func(t *testing.T) {
				p := NewPainter(painterOptions)
				opt := tt.makeOptions()
				opt.Theme = GetTheme(ThemeVividDark)

				validateRadarChartRender(t, p, opt, tt.result)
			})
		} else {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				p := NewPainter(painterOptions)

				validateRadarChartRender(t, p.Child(PainterPaddingOption(NewBoxEqual(20))),
					tt.makeOptions(), tt.result)
			})
		}
	}
}

func validateRadarChartRender(t *testing.T, p *Painter, opt RadarChartOption, expectedResult string) {
	t.Helper()

	err := p.RadarChart(opt)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, expectedResult, data)
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
