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
			result:      "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:rgb(40,40,40)\"/><path  d=\"M 144 3\nL 174 3\nL 174 16\nL 144 16\nL 144 3\" style=\"stroke:none;fill:rgb(255,100,100)\"/><text x=\"176\" y=\"15\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Allocated Budget</text><path  d=\"M 314 3\nL 344 3\nL 344 16\nL 314 16\nL 314 3\" style=\"stroke:none;fill:rgb(255,210,100)\"/><text x=\"346\" y=\"15\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Actual Spending</text><text x=\"0\" y=\"15\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Basic Radar Chart</text><path  d=\"M 300 188\nL 325 203\nL 325 231\nL 300 246\nL 275 231\nL 275 203\nL 300 188\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 159\nL 350 188\nL 350 245\nL 300 275\nL 250 246\nL 250 189\nL 300 159\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 130\nL 375 174\nL 375 260\nL 300 304\nL 225 260\nL 225 174\nL 300 130\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 101\nL 400 159\nL 400 274\nL 300 333\nL 200 275\nL 200 160\nL 300 101\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 72\nL 425 145\nL 425 289\nL 300 362\nL 175 289\nL 175 145\nL 300 72\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 217\nL 300 72\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 217\nL 425 145\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 217\nL 425 289\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 217\nL 300 362\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 217\nL 175 289\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 217\nL 175 145\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><text x=\"284\" y=\"65\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Sales</text><text x=\"430\" y=\"150\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Administration</text><text x=\"430\" y=\"294\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Information Technology</text><text x=\"248\" y=\"379\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Customer Support</text><text x=\"94\" y=\"294\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Development</text><text x=\"111\" y=\"150\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Marketing</text><path  d=\"M 300 124\nL 323 204\nL 383 265\nL 300 350\nL 180 286\nL 210 165\nL 300 124\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:none\"/><path  d=\"M 300 124\nL 323 204\nL 383 265\nL 300 350\nL 180 286\nL 210 165\nL 300 124\" style=\"stroke:none;fill:rgba(255,100,100,0.1)\"/><circle cx=\"300\" cy=\"124\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"323\" cy=\"204\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"383\" cy=\"265\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"300\" cy=\"350\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"180\" cy=\"286\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"210\" cy=\"165\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"300\" cy=\"124\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><path  d=\"M 300 106\nL 409 154\nL 417 284\nL 300 316\nL 199 275\nL 195 157\nL 300 106\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:none\"/><path  d=\"M 300 106\nL 409 154\nL 417 284\nL 300 316\nL 199 275\nL 195 157\nL 300 106\" style=\"stroke:none;fill:rgba(255,210,100,0.1)\"/><circle cx=\"300\" cy=\"106\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"409\" cy=\"154\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"417\" cy=\"284\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"300\" cy=\"316\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"199\" cy=\"275\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"195\" cy=\"157\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"300\" cy=\"106\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/></svg>",
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
