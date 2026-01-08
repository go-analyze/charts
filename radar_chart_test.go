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
		svg         string
		pngCRC      uint32
	}{
		{
			name:        "basic_themed",
			themed:      true,
			makeOptions: makeBasicRadarChartOption,
			svg:         "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:rgb(40,40,40)\"/><path  d=\"M 143 3\nL 173 3\nL 173 16\nL 143 16\nL 143 3\" style=\"stroke:none;fill:rgb(255,100,100)\"/><text x=\"175\" y=\"15\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Allocated Budget</text><path  d=\"M 313 3\nL 343 3\nL 343 16\nL 313 16\nL 313 3\" style=\"stroke:none;fill:rgb(255,210,100)\"/><text x=\"345\" y=\"15\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Actual Spending</text><text x=\"0\" y=\"16\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Basic Radar Chart</text><path  d=\"M 300 189\nL 325 204\nL 325 232\nL 300 247\nL 275 232\nL 275 204\nL 300 189\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 160\nL 350 189\nL 350 246\nL 300 276\nL 250 247\nL 250 190\nL 300 160\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 131\nL 375 175\nL 375 261\nL 300 305\nL 225 261\nL 225 175\nL 300 131\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 102\nL 400 160\nL 400 275\nL 300 334\nL 200 276\nL 200 161\nL 300 102\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 73\nL 425 146\nL 425 290\nL 300 363\nL 175 290\nL 175 146\nL 300 73\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 218\nL 300 73\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 218\nL 425 146\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 218\nL 425 290\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 218\nL 300 363\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 218\nL 175 290\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 300 218\nL 175 146\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><text x=\"284\" y=\"65\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Sales</text><text x=\"430\" y=\"151\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Administration</text><text x=\"430\" y=\"295\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Information Technology</text><text x=\"248\" y=\"381\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Customer Support</text><text x=\"94\" y=\"295\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Development</text><text x=\"112\" y=\"151\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Marketing</text><path  d=\"M 300 125\nL 323 205\nL 383 266\nL 300 351\nL 180 287\nL 210 166\nL 300 125\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:none\"/><path  d=\"M 300 125\nL 323 205\nL 383 266\nL 300 351\nL 180 287\nL 210 166\nL 300 125\" style=\"stroke:none;fill:rgba(255,100,100,0.1)\"/><circle cx=\"300\" cy=\"125\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"323\" cy=\"205\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"383\" cy=\"266\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"300\" cy=\"351\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"180\" cy=\"287\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"210\" cy=\"166\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"300\" cy=\"125\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><path  d=\"M 300 107\nL 409 155\nL 417 285\nL 300 317\nL 199 276\nL 195 158\nL 300 107\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:none\"/><path  d=\"M 300 107\nL 409 155\nL 417 285\nL 300 317\nL 199 276\nL 195 158\nL 300 107\" style=\"stroke:none;fill:rgba(255,210,100,0.1)\"/><circle cx=\"300\" cy=\"107\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"409\" cy=\"155\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"417\" cy=\"285\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"300\" cy=\"317\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"199\" cy=\"276\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"195\" cy=\"158\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"300\" cy=\"107\" r=\"2\" style=\"stroke-width:2;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/></svg>",
			pngCRC:      0x20b7cb3e,
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

				validateRadarChartRender(t, p, rp, tt.makeOptions(), tt.svg, tt.pngCRC)
			})
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-options", func(t *testing.T) {
				p := NewPainter(painterOptions)
				rp := NewPainter(rasterOptions)
				opt := tt.makeOptions()
				opt.Theme = GetTheme(ThemeVividDark)

				validateRadarChartRender(t, p, rp, opt, tt.svg, tt.pngCRC)
			})
		} else {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				p := NewPainter(painterOptions)
				rp := NewPainter(rasterOptions)

				validateRadarChartRender(t, p.Child(PainterPaddingOption(NewBoxEqual(20))),
					rp.Child(PainterPaddingOption(NewBoxEqual(20))), tt.makeOptions(), tt.svg, tt.pngCRC)
			})
		}
	}
}

func validateRadarChartRender(t *testing.T, svgP, pngP *Painter, opt RadarChartOption, expectedSVG string, expectedCRC uint32) {
	t.Helper()

	err := svgP.RadarChart(opt)
	require.NoError(t, err)
	data, err := svgP.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, expectedSVG, data)

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
