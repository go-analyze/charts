package charts

import (
	"strconv"
	"testing"

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
			Data: []string{
				"Allocated Budget", "Actual Spending",
			},
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

func TestRadarChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		defaultTheme bool
		makeOptions  func() RadarChartOption
		result       string
	}{
		{
			name:         "default",
			defaultTheme: true,
			makeOptions:  makeBasicRadarChartOption,
			result:       "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 560 0\nL 560 360\nL 0 360\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 144 29\nL 174 29\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><circle cx=\"159\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><text x=\"176\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Allocated Budget</text><path  d=\"M 314 29\nL 344 29\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><circle cx=\"329\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><text x=\"346\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Actual Spending</text><text x=\"20\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Basic Radar Chart</text><path  d=\"M 300 179\nL 319 191\nL 319 213\nL 300 225\nL 281 213\nL 281 191\nL 300 179\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 300 156\nL 339 179\nL 339 224\nL 300 248\nL 261 225\nL 261 180\nL 300 156\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 300 133\nL 359 168\nL 359 236\nL 300 271\nL 241 236\nL 241 168\nL 300 133\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 300 110\nL 379 156\nL 379 247\nL 300 294\nL 221 248\nL 221 157\nL 300 110\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 300 87\nL 399 145\nL 399 259\nL 300 317\nL 201 259\nL 201 145\nL 300 87\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 300 202\nL 300 87\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 300 202\nL 399 145\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 300 202\nL 399 259\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 300 202\nL 300 317\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 300 202\nL 201 259\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 300 202\nL 201 145\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><text x=\"284\" y=\"80\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Sales</text><text x=\"404\" y=\"150\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Administration</text><text x=\"404\" y=\"264\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Information Technology</text><text x=\"248\" y=\"334\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Customer Support</text><text x=\"120\" y=\"264\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Development</text><text x=\"137\" y=\"150\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Marketing</text><path  d=\"M 300 128\nL 318 192\nL 366 240\nL 300 307\nL 205 257\nL 229 161\nL 300 128\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,0.1)\"/><path  d=\"M 300 128\nL 318 192\nL 366 240\nL 300 307\nL 205 257\nL 229 161\nL 300 128\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,0.1)\"/><circle cx=\"300\" cy=\"128\" r=\"2\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"318\" cy=\"192\" r=\"2\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"366\" cy=\"240\" r=\"2\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"300\" cy=\"307\" r=\"2\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"205\" cy=\"257\" r=\"2\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"229\" cy=\"161\" r=\"2\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"300\" cy=\"128\" r=\"2\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 300 114\nL 387 152\nL 392 255\nL 300 280\nL 220 248\nL 217 154\nL 300 114\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,0.1)\"/><path  d=\"M 300 114\nL 387 152\nL 392 255\nL 300 280\nL 220 248\nL 217 154\nL 300 114\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,0.1)\"/><circle cx=\"300\" cy=\"114\" r=\"2\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"387\" cy=\"152\" r=\"2\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"392\" cy=\"255\" r=\"2\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"300\" cy=\"280\" r=\"2\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"220\" cy=\"248\" r=\"2\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"217\" cy=\"154\" r=\"2\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"300\" cy=\"114\" r=\"2\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/></svg>",
		},
		{
			name:         "themed",
			defaultTheme: false,
			makeOptions:  makeBasicRadarChartOption,
			result:       "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(40,40,40,1.0)\"/><path  d=\"M 144 9\nL 174 9\" style=\"stroke-width:3;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><circle cx=\"159\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><text x=\"176\" y=\"15\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Allocated Budget</text><path  d=\"M 314 9\nL 344 9\" style=\"stroke-width:3;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><circle cx=\"329\" cy=\"9\" r=\"5\" style=\"stroke-width:3;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><text x=\"346\" y=\"15\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Actual Spending</text><text x=\"0\" y=\"15\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Basic Radar Chart</text><path  d=\"M 300 176\nL 322 189\nL 322 214\nL 300 228\nL 278 215\nL 278 190\nL 300 176\" style=\"stroke-width:1;stroke:rgba(72,71,83,1.0);fill:none\"/><path  d=\"M 300 150\nL 345 176\nL 345 227\nL 300 254\nL 255 228\nL 255 177\nL 300 150\" style=\"stroke-width:1;stroke:rgba(72,71,83,1.0);fill:none\"/><path  d=\"M 300 124\nL 367 163\nL 367 240\nL 300 280\nL 233 241\nL 233 164\nL 300 124\" style=\"stroke-width:1;stroke:rgba(72,71,83,1.0);fill:none\"/><path  d=\"M 300 98\nL 390 150\nL 390 253\nL 300 306\nL 210 254\nL 210 151\nL 300 98\" style=\"stroke-width:1;stroke:rgba(72,71,83,1.0);fill:none\"/><path  d=\"M 300 72\nL 412 137\nL 412 266\nL 300 332\nL 188 267\nL 188 138\nL 300 72\" style=\"stroke-width:1;stroke:rgba(72,71,83,1.0);fill:none\"/><path  d=\"M 300 202\nL 300 72\" style=\"stroke-width:1;stroke:rgba(72,71,83,1.0);fill:none\"/><path  d=\"M 300 202\nL 412 137\" style=\"stroke-width:1;stroke:rgba(72,71,83,1.0);fill:none\"/><path  d=\"M 300 202\nL 412 266\" style=\"stroke-width:1;stroke:rgba(72,71,83,1.0);fill:none\"/><path  d=\"M 300 202\nL 300 332\" style=\"stroke-width:1;stroke:rgba(72,71,83,1.0);fill:none\"/><path  d=\"M 300 202\nL 188 267\" style=\"stroke-width:1;stroke:rgba(72,71,83,1.0);fill:none\"/><path  d=\"M 300 202\nL 188 138\" style=\"stroke-width:1;stroke:rgba(72,71,83,1.0);fill:none\"/><text x=\"284\" y=\"65\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Sales</text><text x=\"417\" y=\"142\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Administration</text><text x=\"417\" y=\"271\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Information Technology</text><text x=\"248\" y=\"349\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Customer Support</text><text x=\"107\" y=\"272\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Development</text><text x=\"124\" y=\"143\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Marketing</text><path  d=\"M 300 118\nL 321 190\nL 375 245\nL 300 321\nL 192 264\nL 219 156\nL 300 118\" style=\"stroke-width:2;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,0.1)\"/><path  d=\"M 300 118\nL 321 190\nL 375 245\nL 300 321\nL 192 264\nL 219 156\nL 300 118\" style=\"stroke-width:2;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,0.1)\"/><circle cx=\"300\" cy=\"118\" r=\"2\" style=\"stroke-width:2;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><circle cx=\"321\" cy=\"190\" r=\"2\" style=\"stroke-width:2;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><circle cx=\"375\" cy=\"245\" r=\"2\" style=\"stroke-width:2;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><circle cx=\"300\" cy=\"321\" r=\"2\" style=\"stroke-width:2;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><circle cx=\"192\" cy=\"264\" r=\"2\" style=\"stroke-width:2;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><circle cx=\"219\" cy=\"156\" r=\"2\" style=\"stroke-width:2;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><circle cx=\"300\" cy=\"118\" r=\"2\" style=\"stroke-width:2;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><path  d=\"M 300 102\nL 398 146\nL 405 262\nL 300 290\nL 210 254\nL 206 148\nL 300 102\" style=\"stroke-width:2;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,0.1)\"/><path  d=\"M 300 102\nL 398 146\nL 405 262\nL 300 290\nL 210 254\nL 206 148\nL 300 102\" style=\"stroke-width:2;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,0.1)\"/><circle cx=\"300\" cy=\"102\" r=\"2\" style=\"stroke-width:2;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><circle cx=\"398\" cy=\"146\" r=\"2\" style=\"stroke-width:2;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><circle cx=\"405\" cy=\"262\" r=\"2\" style=\"stroke-width:2;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><circle cx=\"300\" cy=\"290\" r=\"2\" style=\"stroke-width:2;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><circle cx=\"210\" cy=\"254\" r=\"2\" style=\"stroke-width:2;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><circle cx=\"206\" cy=\"148\" r=\"2\" style=\"stroke-width:2;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><circle cx=\"300\" cy=\"102\" r=\"2\" style=\"stroke-width:2;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/></svg>",
		},
	}

	for i, tt := range tests {
		painterOptions := PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        600,
			Height:       400,
		}
		if tt.defaultTheme {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				p, err := NewPainter(painterOptions)
				require.NoError(t, err)

				validateRadarChartRender(t, p.Child(PainterPaddingOption(Box{
					Left:   20,
					Top:    20,
					Right:  20,
					Bottom: 20,
				})), tt.makeOptions(), tt.result)
			})
		} else {
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-painter", func(t *testing.T) {
				p, err := NewPainter(painterOptions, PainterThemeOption(GetTheme(ThemeVividDark)))
				require.NoError(t, err)

				validateRadarChartRender(t, p, tt.makeOptions(), tt.result)
			})
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-options", func(t *testing.T) {
				p, err := NewPainter(painterOptions)
				require.NoError(t, err)
				opt := tt.makeOptions()
				opt.Theme = GetTheme(ThemeVividDark)

				validateRadarChartRender(t, p, opt, tt.result)
			})
		}
	}
}

func validateRadarChartRender(t *testing.T, p *Painter, opt RadarChartOption, expectedResult string) {
	t.Helper()

	_, err := NewRadarChart(p, opt).Render()
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, expectedResult, data)
}
