package charts

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func makeBasicPieChartOption() PieChartOption {
	values := []float64{
		1048,
		735,
		580,
		484,
		300,
	}
	return PieChartOption{
		SeriesList: NewPieSeriesList(values, PieSeriesOption{
			Label: SeriesLabel{
				Show: true,
			},
		}),
		Title: TitleOption{
			Text:    "Rainfall vs Evaporation",
			Subtext: "Fake Data",
			Left:    PositionCenter,
		},
		Padding: Box{
			Top:    20,
			Right:  20,
			Bottom: 20,
			Left:   20,
		},
		Legend: LegendOption{
			Orient: OrientVertical,
			Data: []string{
				"Search Engine",
				"Direct",
				"Email",
				"Union Ads",
				"Video Ads",
			},
			Left: PositionLeft,
		},
	}
}

func TestPieChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		defaultTheme bool
		makeOptions  func() PieChartOption
		result       string
	}{
		{
			name:         "default",
			defaultTheme: true,
			makeOptions:  makeBasicPieChartOption,
			result:       "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 560 0\nL 560 360\nL 0 360\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 40 49\nL 70 49\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><circle cx=\"55\" cy=\"49\" r=\"5\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><text x=\"72\" y=\"55\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><path  d=\"M 40 69\nL 70 69\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><circle cx=\"55\" cy=\"69\" r=\"5\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><text x=\"72\" y=\"75\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path  d=\"M 40 89\nL 70 89\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><circle cx=\"55\" cy=\"89\" r=\"5\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><text x=\"72\" y=\"95\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path  d=\"M 40 109\nL 70 109\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><circle cx=\"55\" cy=\"109\" r=\"5\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><text x=\"72\" y=\"115\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path  d=\"M 40 129\nL 70 129\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><circle cx=\"55\" cy=\"129\" r=\"5\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><text x=\"72\" y=\"135\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><text x=\"222\" y=\"55\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Rainfall vs Evaporation</text><text x=\"266\" y=\"70\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fake Data</text><path  d=\"M 300 210\nL 300 114\nA 96 96 119.89 0 1 383 257\nL 300 210\nZ\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path  d=\"M 383 162\nL 396 155\nM 396 155\nL 411 155\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><text x=\"414\" y=\"160\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Search Engine: 33.3%</text><path  d=\"M 300 210\nL 383 257\nA 96 96 84.08 0 1 262 297\nL 300 210\nZ\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path  d=\"M 329 301\nL 334 315\nM 334 315\nL 349 315\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><text x=\"352\" y=\"320\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Direct: 23.35%</text><path  d=\"M 300 210\nL 262 297\nA 96 96 66.35 0 1 205 210\nL 300 210\nZ\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><path  d=\"M 220 262\nL 207 270\nM 207 270\nL 192 270\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><text x=\"108\" y=\"275\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Email: 18.43%</text><path  d=\"M 300 210\nL 205 210\nA 96 96 55.37 0 1 246 131\nL 300 210\nZ\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><path  d=\"M 216 165\nL 202 158\nM 202 158\nL 187 158\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><text x=\"76\" y=\"163\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Union Ads: 15.37%</text><path  d=\"M 300 210\nL 246 131\nA 96 96 34.32 0 1 300 114\nL 300 210\nZ\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><path  d=\"M 272 119\nL 268 104\nM 268 104\nL 253 104\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><text x=\"150\" y=\"109\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Video Ads: 9.53%</text></svg>",
		},
		{
			name:         "themed",
			defaultTheme: false,
			makeOptions:  makeBasicPieChartOption,
			result:       "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(40,40,40,1.0)\"/><path  d=\"M 20 29\nL 50 29\" style=\"stroke-width:3;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><circle cx=\"35\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><text x=\"52\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><path  d=\"M 20 49\nL 50 49\" style=\"stroke-width:3;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><circle cx=\"35\" cy=\"49\" r=\"5\" style=\"stroke-width:3;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><text x=\"52\" y=\"55\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path  d=\"M 20 69\nL 50 69\" style=\"stroke-width:3;stroke:rgba(100,180,210,1.0);fill:rgba(100,180,210,1.0)\"/><circle cx=\"35\" cy=\"69\" r=\"5\" style=\"stroke-width:3;stroke:rgba(100,180,210,1.0);fill:rgba(100,180,210,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(100,180,210,1.0);fill:rgba(100,180,210,1.0)\"/><text x=\"52\" y=\"75\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path  d=\"M 20 89\nL 50 89\" style=\"stroke-width:3;stroke:rgba(64,160,110,1.0);fill:rgba(64,160,110,1.0)\"/><circle cx=\"35\" cy=\"89\" r=\"5\" style=\"stroke-width:3;stroke:rgba(64,160,110,1.0);fill:rgba(64,160,110,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(64,160,110,1.0);fill:rgba(64,160,110,1.0)\"/><text x=\"52\" y=\"95\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path  d=\"M 20 109\nL 50 109\" style=\"stroke-width:3;stroke:rgba(154,100,180,1.0);fill:rgba(154,100,180,1.0)\"/><circle cx=\"35\" cy=\"109\" r=\"5\" style=\"stroke-width:3;stroke:rgba(154,100,180,1.0);fill:rgba(154,100,180,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(154,100,180,1.0);fill:rgba(154,100,180,1.0)\"/><text x=\"52\" y=\"115\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><text x=\"222\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Rainfall vs Evaporation</text><text x=\"266\" y=\"50\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fake Data</text><path  d=\"M 300 210\nL 300 98\nA 112 112 119.89 0 1 397 265\nL 300 210\nZ\" style=\"stroke-width:1;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><path  d=\"M 396 154\nL 409 147\nM 409 147\nL 424 147\" style=\"stroke-width:1;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><text x=\"427\" y=\"152\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Search Engine: 33.3%</text><path  d=\"M 300 210\nL 397 265\nA 112 112 84.08 0 1 255 312\nL 300 210\nZ\" style=\"stroke-width:1;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><path  d=\"M 334 316\nL 339 330\nM 339 330\nL 354 330\" style=\"stroke-width:1;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><text x=\"357\" y=\"335\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Direct: 23.35%</text><path  d=\"M 300 210\nL 255 312\nA 112 112 66.35 0 1 189 210\nL 300 210\nZ\" style=\"stroke-width:1;stroke:rgba(100,180,210,1.0);fill:rgba(100,180,210,1.0)\"/><path  d=\"M 206 270\nL 194 278\nM 194 278\nL 179 278\" style=\"stroke-width:1;stroke:rgba(100,180,210,1.0);fill:rgba(100,180,210,1.0)\"/><text x=\"95\" y=\"283\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Email: 18.43%</text><path  d=\"M 300 210\nL 189 210\nA 112 112 55.37 0 1 237 118\nL 300 210\nZ\" style=\"stroke-width:1;stroke:rgba(64,160,110,1.0);fill:rgba(64,160,110,1.0)\"/><path  d=\"M 202 158\nL 188 151\nM 188 151\nL 173 151\" style=\"stroke-width:1;stroke:rgba(64,160,110,1.0);fill:rgba(64,160,110,1.0)\"/><text x=\"62\" y=\"156\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Union Ads: 15.37%</text><path  d=\"M 300 210\nL 237 118\nA 112 112 34.32 0 1 300 98\nL 300 210\nZ\" style=\"stroke-width:1;stroke:rgba(154,100,180,1.0);fill:rgba(154,100,180,1.0)\"/><path  d=\"M 267 103\nL 263 89\nM 263 89\nL 248 89\" style=\"stroke-width:1;stroke:rgba(154,100,180,1.0);fill:rgba(154,100,180,1.0)\"/><text x=\"145\" y=\"94\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Video Ads: 9.53%</text></svg>",
		},
	}

	for _, tt := range tests {
		painterOptions := PainterOptions{
			Type:   ChartOutputSVG,
			Width:  600,
			Height: 400,
		}
		if tt.defaultTheme {
			t.Run(tt.name, func(t *testing.T) {
				p, err := NewPainter(painterOptions)
				require.NoError(t, err)

				validatePieChartRender(t, p.Child(PainterPaddingOption(Box{
					Left:   20,
					Top:    20,
					Right:  20,
					Bottom: 20,
				})), tt.makeOptions(), tt.result)
			})
		} else {
			t.Run(tt.name+"-painter", func(t *testing.T) {
				p, err := NewPainter(painterOptions, PainterThemeOption(GetTheme(ThemeVividDark)))
				require.NoError(t, err)

				validatePieChartRender(t, p, tt.makeOptions(), tt.result)
			})
			t.Run(tt.name+"-options", func(t *testing.T) {
				p, err := NewPainter(painterOptions)
				require.NoError(t, err)
				opt := tt.makeOptions()
				opt.Theme = GetTheme(ThemeVividDark)

				validatePieChartRender(t, p, opt, tt.result)
			})
		}
	}
}

func validatePieChartRender(t *testing.T, p *Painter, opt PieChartOption, expectedResult string) {
	t.Helper()

	_, err := NewPieChart(p, opt).Render()
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, expectedResult, string(data))
}
