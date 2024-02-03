package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHorizontalBarChart(t *testing.T) {
	tests := []struct {
		render func(*Painter) ([]byte, error)
		result string
	}{
		{
			render: func(p *Painter) ([]byte, error) {
				_, err := NewHorizontalBarChart(p, HorizontalBarChartOption{
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
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 224 19\nL 254 19\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><circle cx=\"239\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><text x=\"256\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2011</text><path  d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><text x=\"343\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2012</text><text x=\"10\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World Population</text><path  d=\"M 83 45\nL 88 45\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 83 97\nL 88 97\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 83 150\nL 88 150\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 83 202\nL 88 202\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 83 255\nL 88 255\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 83 307\nL 88 307\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 83 360\nL 88 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 88 45\nL 88 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"37\" y=\"78\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World</text><text x=\"38\" y=\"130\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">China</text><text x=\"44\" y=\"183\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">India</text><text x=\"48\" y=\"235\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">USA</text><text x=\"10\" y=\"288\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Indonesia</text><text x=\"39\" y=\"340\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Brazil</text><text x=\"88\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"171\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">116.66k</text><text x=\"255\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">233.33k</text><text x=\"339\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">350k</text><text x=\"422\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">466.66k</text><text x=\"506\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">583.33k</text><text x=\"555\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">700k</text><path  d=\"M 171 45\nL 171 360\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 255 45\nL 255 360\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 339 45\nL 339 360\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 422 45\nL 422 360\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 506 45\nL 506 360\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 590 45\nL 590 360\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 88 317\nL 88 317\nL 88 330\nL 88 330\nL 88 317\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 88 265\nL 91 265\nL 91 278\nL 88 278\nL 88 265\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 88 212\nL 96 212\nL 96 225\nL 88 225\nL 88 212\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 88 160\nL 153 160\nL 153 173\nL 88 173\nL 88 160\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 88 107\nL 173 107\nL 173 120\nL 88 120\nL 88 107\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 88 55\nL 550 55\nL 550 68\nL 88 68\nL 88 55\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 88 335\nL 88 335\nL 88 348\nL 88 348\nL 88 335\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 88 283\nL 91 283\nL 91 296\nL 88 296\nL 88 283\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 88 230\nL 97 230\nL 97 243\nL 88 243\nL 88 230\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 88 178\nL 166 178\nL 166 191\nL 88 191\nL 88 178\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 88 125\nL 175 125\nL 175 138\nL 88 138\nL 88 125\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 88 73\nL 590 73\nL 590 86\nL 88 86\nL 88 73\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/></svg>",
		},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			p, err := NewPainter(PainterOptions{
				Type:   ChartOutputSVG,
				Width:  600,
				Height: 400,
			}, PainterThemeOption(defaultTheme))
			require.NoError(t, err)
			data, err := tt.render(p)
			require.NoError(t, err)
			assertEqualSVG(t, tt.result, string(data))
		})
	}
}
