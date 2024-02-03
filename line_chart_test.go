package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLineChart(t *testing.T) {
	tests := []struct {
		render func(*Painter) ([]byte, error)
		result string
	}{
		{
			render: func(p *Painter) ([]byte, error) {
				values := [][]float64{
					{
						120,
						132,
						101,
						134,
						90,
						230,
						210,
					},
					{
						220,
						182,
						191,
						234,
						290,
						330,
						310,
					},
					{
						150,
						232,
						201,
						154,
						190,
						330,
						410,
					},
					{
						320,
						332,
						301,
						334,
						390,
						330,
						320,
					},
					{
						820,
						932,
						901,
						934,
						1290,
						1330,
						1320,
					},
				}
				_, err := NewLineChart(p, LineChartOption{
					Title: TitleOption{
						Text: "Line",
					},
					Padding: Box{
						Top:    10,
						Right:  10,
						Bottom: 10,
						Left:   10,
					},
					XAxis: NewXAxisOption([]string{
						"Mon",
						"Tue",
						"Wed",
						"Thu",
						"Fri",
						"Sat",
						"Sun",
					}),
					Legend: NewLegendOption([]string{
						"Email",
						"Union Ads",
						"Video Ads",
						"Direct",
						"Search Engine",
					}, PositionCenter),
					SeriesList: NewSeriesListDataFromValues(values),
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 20 19\nL 50 19\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><circle cx=\"35\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><text x=\"52\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path  d=\"M 111 19\nL 141 19\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><circle cx=\"126\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><text x=\"143\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path  d=\"M 234 19\nL 264 19\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><circle cx=\"249\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><text x=\"266\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><path  d=\"M 357 19\nL 387 19\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><circle cx=\"372\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><text x=\"389\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path  d=\"M 450 19\nL 480 19\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><circle cx=\"465\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><text x=\"482\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><text x=\"10\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line</text><text x=\"28\" y=\"52\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.4k</text><text x=\"19\" y=\"104\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.16k</text><text x=\"10\" y=\"157\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">933.33</text><text x=\"31\" y=\"209\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">700</text><text x=\"10\" y=\"262\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">466.66</text><text x=\"10\" y=\"314\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">233.33</text><text x=\"49\" y=\"367\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 68 45\nL 590 45\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 97\nL 590 97\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 150\nL 590 150\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 202\nL 590 202\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 255\nL 590 255\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 307\nL 590 307\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 365\nL 68 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 142 365\nL 142 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 217 365\nL 217 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 291 365\nL 291 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 366 365\nL 366 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 440 365\nL 440 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 515 365\nL 515 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 590 365\nL 590 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 68 360\nL 590 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"90\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"166\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"239\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"315\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"394\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"466\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"539\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path  d=\"M 105 333\nL 179 331\nL 254 338\nL 328 330\nL 403 340\nL 477 309\nL 552 313\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:none\"/><circle cx=\"105\" cy=\"333\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"179\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"254\" cy=\"338\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"328\" cy=\"330\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"403\" cy=\"340\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"477\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"313\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 105 311\nL 179 320\nL 254 318\nL 328 308\nL 403 295\nL 477 286\nL 552 291\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:none\"/><circle cx=\"105\" cy=\"311\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"179\" cy=\"320\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"254\" cy=\"318\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"328\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"403\" cy=\"295\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"477\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"291\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 105 327\nL 179 308\nL 254 315\nL 328 326\nL 403 318\nL 477 286\nL 552 268\" style=\"stroke-width:2;stroke:rgba(250,200,88,1.0);fill:none\"/><circle cx=\"105\" cy=\"327\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"179\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"254\" cy=\"315\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"328\" cy=\"326\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"403\" cy=\"318\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"477\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"268\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 105 288\nL 179 286\nL 254 293\nL 328 285\nL 403 273\nL 477 286\nL 552 288\" style=\"stroke-width:2;stroke:rgba(238,102,102,1.0);fill:none\"/><circle cx=\"105\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"179\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"254\" cy=\"293\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"328\" cy=\"285\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"403\" cy=\"273\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"477\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 105 176\nL 179 151\nL 254 158\nL 328 150\nL 403 70\nL 477 61\nL 552 63\" style=\"stroke-width:2;stroke:rgba(115,192,222,1.0);fill:none\"/><circle cx=\"105\" cy=\"176\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"179\" cy=\"151\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"254\" cy=\"158\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"328\" cy=\"150\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"403\" cy=\"70\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"477\" cy=\"61\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"63\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/></svg>",
		},
		{
			render: func(p *Painter) ([]byte, error) {
				values := [][]float64{
					{
						120,
						132,
						101,
						134,
						90,
						230,
						210,
					},
					{
						220,
						182,
						191,
						234,
						290,
						330,
						310,
					},
					{
						150,
						232,
						201,
						154,
						190,
						330,
						410,
					},
					{
						320,
						332,
						301,
						334,
						390,
						330,
						320,
					},
					{
						820,
						932,
						901,
						934,
						1290,
						1330,
						1320,
					},
				}
				_, err := NewLineChart(p, LineChartOption{
					Title: TitleOption{
						Text: "Line",
					},
					Padding: Box{
						Top:    10,
						Right:  10,
						Bottom: 10,
						Left:   10,
					},
					XAxis: NewXAxisOption([]string{
						"Mon",
						"Tue",
						"Wed",
						"Thu",
						"Fri",
						"Sat",
						"Sun",
					}, FalseFlag()),
					Legend: NewLegendOption([]string{
						"Email",
						"Union Ads",
						"Video Ads",
						"Direct",
						"Search Engine",
					}, PositionCenter),
					SeriesList: NewSeriesListDataFromValues(values),
				}).Render()
				if err != nil {
					return nil, err
				}
				return p.Bytes()
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 20 19\nL 50 19\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><circle cx=\"35\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><text x=\"52\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path  d=\"M 111 19\nL 141 19\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><circle cx=\"126\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><text x=\"143\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path  d=\"M 234 19\nL 264 19\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><circle cx=\"249\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><text x=\"266\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><path  d=\"M 357 19\nL 387 19\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><circle cx=\"372\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><text x=\"389\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path  d=\"M 450 19\nL 480 19\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><circle cx=\"465\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><text x=\"482\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><text x=\"10\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line</text><text x=\"28\" y=\"52\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.4k</text><text x=\"19\" y=\"104\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.16k</text><text x=\"10\" y=\"157\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">933.33</text><text x=\"31\" y=\"209\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">700</text><text x=\"10\" y=\"262\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">466.66</text><text x=\"10\" y=\"314\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">233.33</text><text x=\"49\" y=\"367\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 68 45\nL 590 45\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 97\nL 590 97\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 150\nL 590 150\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 202\nL 590 202\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 255\nL 590 255\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 307\nL 590 307\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 365\nL 68 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 155 365\nL 155 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 242 365\nL 242 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 329 365\nL 329 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 416 365\nL 416 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 503 365\nL 503 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 590 365\nL 590 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 68 360\nL 590 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"68\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"155\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"242\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"329\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"416\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"503\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"563\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path  d=\"M 68 333\nL 155 331\nL 242 338\nL 329 330\nL 416 340\nL 503 309\nL 590 313\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:none\"/><circle cx=\"68\" cy=\"333\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"155\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"242\" cy=\"338\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"329\" cy=\"330\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"416\" cy=\"340\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"503\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"590\" cy=\"313\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 68 311\nL 155 320\nL 242 318\nL 329 308\nL 416 295\nL 503 286\nL 590 291\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:none\"/><circle cx=\"68\" cy=\"311\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"155\" cy=\"320\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"242\" cy=\"318\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"329\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"416\" cy=\"295\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"503\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"590\" cy=\"291\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 68 327\nL 155 308\nL 242 315\nL 329 326\nL 416 318\nL 503 286\nL 590 268\" style=\"stroke-width:2;stroke:rgba(250,200,88,1.0);fill:none\"/><circle cx=\"68\" cy=\"327\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"155\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"242\" cy=\"315\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"329\" cy=\"326\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"416\" cy=\"318\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"503\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"590\" cy=\"268\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 68 288\nL 155 286\nL 242 293\nL 329 285\nL 416 273\nL 503 286\nL 590 288\" style=\"stroke-width:2;stroke:rgba(238,102,102,1.0);fill:none\"/><circle cx=\"68\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"155\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"242\" cy=\"293\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"329\" cy=\"285\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"416\" cy=\"273\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"503\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"590\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 68 176\nL 155 151\nL 242 158\nL 329 150\nL 416 70\nL 503 61\nL 590 63\" style=\"stroke-width:2;stroke:rgba(115,192,222,1.0);fill:none\"/><circle cx=\"68\" cy=\"176\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"155\" cy=\"151\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"242\" cy=\"158\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"329\" cy=\"150\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"416\" cy=\"70\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"503\" cy=\"61\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"590\" cy=\"63\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/></svg>",
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
