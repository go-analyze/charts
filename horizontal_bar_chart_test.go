package charts

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func makeBasicHorizontalBarChartOption() HorizontalBarChartOption {
	return HorizontalBarChartOption{
		Padding: Box{
			Top:    10,
			Right:  10,
			Bottom: 10,
			Left:   10,
		},
		SeriesList: NewSeriesListDataFromValues([][]float64{
			{
				18203,
				23489,
				29034,
				104970,
				131744,
				630230,
			},
			{
				19325,
				23438,
				31000,
				121594,
				134141,
				681807,
			},
		}, ChartTypeHorizontalBar),
		Title: TitleOption{
			Text: "World Population",
		},
		Legend: NewLegendOption([]string{
			"2011",
			"2012",
		}),
		YAxisOptions: NewYAxisOptions([]string{
			"Brazil",
			"Indonesia",
			"USA",
			"India",
			"China",
			"World",
		}),
	}
}

func TestHorizontalBarChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		defaultTheme bool
		makeOptions  func() HorizontalBarChartOption
		result       string
	}{
		{
			name:         "default",
			defaultTheme: true,
			makeOptions:  makeBasicHorizontalBarChartOption,
			result:       "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 224 19\nL 254 19\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><circle cx=\"239\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><text x=\"256\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2011</text><path  d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><text x=\"343\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2012</text><text x=\"10\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World Population</text><path  d=\"M 83 45\nL 88 45\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 83 97\nL 88 97\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 83 150\nL 88 150\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 83 202\nL 88 202\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 83 255\nL 88 255\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 83 307\nL 88 307\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 83 360\nL 88 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 88 45\nL 88 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"37\" y=\"78\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World</text><text x=\"38\" y=\"130\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">China</text><text x=\"44\" y=\"183\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">India</text><text x=\"48\" y=\"235\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">USA</text><text x=\"10\" y=\"288\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Indonesia</text><text x=\"39\" y=\"340\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Brazil</text><text x=\"88\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"188\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">144k</text><text x=\"288\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">288k</text><text x=\"389\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">432k</text><text x=\"489\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">576k</text><text x=\"555\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">720k</text><path  d=\"M 188 45\nL 188 360\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 288 45\nL 288 360\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 389 45\nL 389 360\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 489 45\nL 489 360\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 590 45\nL 590 360\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 88 317\nL 100 317\nL 100 330\nL 88 330\nL 88 317\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 88 265\nL 103 265\nL 103 278\nL 88 278\nL 88 265\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 88 212\nL 107 212\nL 107 225\nL 88 225\nL 88 212\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 88 160\nL 158 160\nL 158 173\nL 88 173\nL 88 160\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 88 107\nL 176 107\nL 176 120\nL 88 120\nL 88 107\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 88 55\nL 509 55\nL 509 68\nL 88 68\nL 88 55\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 88 335\nL 100 335\nL 100 348\nL 88 348\nL 88 335\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 88 283\nL 103 283\nL 103 296\nL 88 296\nL 88 283\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 88 230\nL 108 230\nL 108 243\nL 88 243\nL 88 230\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 88 178\nL 169 178\nL 169 191\nL 88 191\nL 88 178\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 88 125\nL 177 125\nL 177 138\nL 88 138\nL 88 125\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 88 73\nL 544 73\nL 544 86\nL 88 86\nL 88 73\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/></svg>",
		},
		{
			name:         "themed",
			defaultTheme: false,
			makeOptions:  makeBasicHorizontalBarChartOption,
			result:       "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(40,40,40,1.0)\"/><path  d=\"M 224 19\nL 254 19\" style=\"stroke-width:3;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><circle cx=\"239\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(255,100,100,1.0);fill:rgba(255,100,100,1.0)\"/><text x=\"256\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2011</text><path  d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(255,210,100,1.0);fill:rgba(255,210,100,1.0)\"/><text x=\"343\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2012</text><text x=\"10\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World Population</text><path  d=\"M 83 45\nL 88 45\" style=\"stroke-width:1;stroke:rgba(185,184,206,1.0);fill:none\"/><path  d=\"M 83 97\nL 88 97\" style=\"stroke-width:1;stroke:rgba(185,184,206,1.0);fill:none\"/><path  d=\"M 83 150\nL 88 150\" style=\"stroke-width:1;stroke:rgba(185,184,206,1.0);fill:none\"/><path  d=\"M 83 202\nL 88 202\" style=\"stroke-width:1;stroke:rgba(185,184,206,1.0);fill:none\"/><path  d=\"M 83 255\nL 88 255\" style=\"stroke-width:1;stroke:rgba(185,184,206,1.0);fill:none\"/><path  d=\"M 83 307\nL 88 307\" style=\"stroke-width:1;stroke:rgba(185,184,206,1.0);fill:none\"/><path  d=\"M 83 360\nL 88 360\" style=\"stroke-width:1;stroke:rgba(185,184,206,1.0);fill:none\"/><path  d=\"M 88 45\nL 88 360\" style=\"stroke-width:1;stroke:rgba(185,184,206,1.0);fill:none\"/><text x=\"37\" y=\"78\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World</text><text x=\"38\" y=\"130\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">China</text><text x=\"44\" y=\"183\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">India</text><text x=\"48\" y=\"235\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">USA</text><text x=\"10\" y=\"288\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Indonesia</text><text x=\"39\" y=\"340\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Brazil</text><text x=\"88\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"188\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">144k</text><text x=\"288\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">288k</text><text x=\"389\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">432k</text><text x=\"489\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">576k</text><text x=\"555\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">720k</text><path  d=\"M 188 45\nL 188 360\" style=\"stroke-width:1;stroke:rgba(72,71,83,1.0);fill:none\"/><path  d=\"M 288 45\nL 288 360\" style=\"stroke-width:1;stroke:rgba(72,71,83,1.0);fill:none\"/><path  d=\"M 389 45\nL 389 360\" style=\"stroke-width:1;stroke:rgba(72,71,83,1.0);fill:none\"/><path  d=\"M 489 45\nL 489 360\" style=\"stroke-width:1;stroke:rgba(72,71,83,1.0);fill:none\"/><path  d=\"M 590 45\nL 590 360\" style=\"stroke-width:1;stroke:rgba(72,71,83,1.0);fill:none\"/><path  d=\"M 88 317\nL 100 317\nL 100 330\nL 88 330\nL 88 317\" style=\"stroke-width:0;stroke:none;fill:rgba(255,100,100,1.0)\"/><path  d=\"M 88 265\nL 103 265\nL 103 278\nL 88 278\nL 88 265\" style=\"stroke-width:0;stroke:none;fill:rgba(255,100,100,1.0)\"/><path  d=\"M 88 212\nL 107 212\nL 107 225\nL 88 225\nL 88 212\" style=\"stroke-width:0;stroke:none;fill:rgba(255,100,100,1.0)\"/><path  d=\"M 88 160\nL 158 160\nL 158 173\nL 88 173\nL 88 160\" style=\"stroke-width:0;stroke:none;fill:rgba(255,100,100,1.0)\"/><path  d=\"M 88 107\nL 176 107\nL 176 120\nL 88 120\nL 88 107\" style=\"stroke-width:0;stroke:none;fill:rgba(255,100,100,1.0)\"/><path  d=\"M 88 55\nL 509 55\nL 509 68\nL 88 68\nL 88 55\" style=\"stroke-width:0;stroke:none;fill:rgba(255,100,100,1.0)\"/><path  d=\"M 88 335\nL 100 335\nL 100 348\nL 88 348\nL 88 335\" style=\"stroke-width:0;stroke:none;fill:rgba(255,210,100,1.0)\"/><path  d=\"M 88 283\nL 103 283\nL 103 296\nL 88 296\nL 88 283\" style=\"stroke-width:0;stroke:none;fill:rgba(255,210,100,1.0)\"/><path  d=\"M 88 230\nL 108 230\nL 108 243\nL 88 243\nL 88 230\" style=\"stroke-width:0;stroke:none;fill:rgba(255,210,100,1.0)\"/><path  d=\"M 88 178\nL 169 178\nL 169 191\nL 88 191\nL 88 178\" style=\"stroke-width:0;stroke:none;fill:rgba(255,210,100,1.0)\"/><path  d=\"M 88 125\nL 177 125\nL 177 138\nL 88 138\nL 88 125\" style=\"stroke-width:0;stroke:none;fill:rgba(255,210,100,1.0)\"/><path  d=\"M 88 73\nL 544 73\nL 544 86\nL 88 86\nL 88 73\" style=\"stroke-width:0;stroke:none;fill:rgba(255,210,100,1.0)\"/></svg>",
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

				validateHorizontalBarChartRender(t, p, tt.makeOptions(), tt.result)
			})
		} else {
			t.Run(tt.name+"-painter", func(t *testing.T) {
				p, err := NewPainter(painterOptions, PainterThemeOption(GetTheme(ThemeVividDark)))
				require.NoError(t, err)

				validateHorizontalBarChartRender(t, p, tt.makeOptions(), tt.result)
			})
			t.Run(tt.name+"-options", func(t *testing.T) {
				p, err := NewPainter(painterOptions)
				require.NoError(t, err)
				opt := tt.makeOptions()
				opt.Theme = GetTheme(ThemeVividDark)

				validateHorizontalBarChartRender(t, p, opt, tt.result)
			})
		}
	}
}

func validateHorizontalBarChartRender(t *testing.T, p *Painter, opt HorizontalBarChartOption, expectedResult string) {
	t.Helper()

	_, err := NewHorizontalBarChart(p, opt).Render()
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, expectedResult, string(data))
}
