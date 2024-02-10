package charts

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func makeFullLineChartOption() LineChartOption {
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
	return LineChartOption{
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
	}
}

func makeBasicLineChartOption() LineChartOption {
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
			820,
			932,
			901,
			934,
			1290,
			1330,
			1320,
		},
	}
	return LineChartOption{
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
			"A",
			"B",
			"C",
			"D",
			"E",
			"F",
			"G",
		}),
		Legend: NewLegendOption([]string{
			"1",
			"2",
		}, PositionCenter),
		SeriesList: NewSeriesListDataFromValues(values),
	}
}

func makeMinimalLineChartOption() LineChartOption {
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
			820,
			932,
			901,
			934,
			1290,
			1330,
			1320,
		},
	}
	return LineChartOption{
		Padding: Box{
			Top:    10,
			Right:  10,
			Bottom: 10,
			Left:   10,
		},
		XAxis: XAxisOption{
			Data: []string{
				"1",
				"2",
				"3",
				"4",
				"5",
				"6",
				"7",
			},
			Show: FalseFlag(),
		},
		SeriesList: NewSeriesListDataFromValues(values),
	}
}

func TestLineChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		makeOptions func() LineChartOption
		result      string
	}{
		{
			name:        "Basic",
			makeOptions: makeFullLineChartOption,
			result:      "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 20 19\nL 50 19\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><circle cx=\"35\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><text x=\"52\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path  d=\"M 111 19\nL 141 19\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><circle cx=\"126\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><text x=\"143\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path  d=\"M 234 19\nL 264 19\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><circle cx=\"249\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><text x=\"266\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><path  d=\"M 357 19\nL 387 19\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><circle cx=\"372\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><text x=\"389\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path  d=\"M 450 19\nL 480 19\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><circle cx=\"465\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><text x=\"482\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><text x=\"10\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line</text><text x=\"10\" y=\"52\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"10\" y=\"87\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.28k</text><text x=\"10\" y=\"122\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.12k</text><text x=\"22\" y=\"157\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"22\" y=\"192\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"22\" y=\"227\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"22\" y=\"262\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"22\" y=\"297\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"22\" y=\"332\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"40\" y=\"367\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 59 45\nL 590 45\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 80\nL 590 80\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 115\nL 590 115\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 150\nL 590 150\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 185\nL 590 185\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 220\nL 590 220\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 255\nL 590 255\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 290\nL 590 290\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 325\nL 590 325\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 365\nL 59 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 134 365\nL 134 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 210 365\nL 210 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 286 365\nL 286 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 362 365\nL 362 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 438 365\nL 438 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 514 365\nL 514 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 590 365\nL 590 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 59 360\nL 590 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"81\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"159\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"233\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"311\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"391\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"465\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"539\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path  d=\"M 96 334\nL 172 332\nL 248 338\nL 324 331\nL 400 341\nL 476 310\nL 552 315\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:none\"/><circle cx=\"96\" cy=\"334\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"172\" cy=\"332\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"248\" cy=\"338\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"400\" cy=\"341\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"476\" cy=\"310\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"315\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 96 312\nL 172 321\nL 248 319\nL 324 309\nL 400 297\nL 476 288\nL 552 293\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:none\"/><circle cx=\"96\" cy=\"312\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"172\" cy=\"321\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"248\" cy=\"319\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"400\" cy=\"297\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"476\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"293\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 96 328\nL 172 310\nL 248 317\nL 324 327\nL 400 319\nL 476 288\nL 552 271\" style=\"stroke-width:2;stroke:rgba(250,200,88,1.0);fill:none\"/><circle cx=\"96\" cy=\"328\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"172\" cy=\"310\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"248\" cy=\"317\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"327\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"400\" cy=\"319\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"476\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"271\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 96 290\nL 172 288\nL 248 295\nL 324 287\nL 400 275\nL 476 288\nL 552 290\" style=\"stroke-width:2;stroke:rgba(238,102,102,1.0);fill:none\"/><circle cx=\"96\" cy=\"290\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"172\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"248\" cy=\"295\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"287\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"400\" cy=\"275\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"476\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"290\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 96 181\nL 172 157\nL 248 163\nL 324 156\nL 400 78\nL 476 70\nL 552 72\" style=\"stroke-width:2;stroke:rgba(115,192,222,1.0);fill:none\"/><circle cx=\"96\" cy=\"181\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"172\" cy=\"157\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"248\" cy=\"163\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"156\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"400\" cy=\"78\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"476\" cy=\"70\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"72\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/></svg>",
		},
		{
			name: "BasicWithoutBoundary",
			makeOptions: func() LineChartOption {
				opt := makeFullLineChartOption()
				opt.XAxis.BoundaryGap = FalseFlag()
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 20 19\nL 50 19\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><circle cx=\"35\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><text x=\"52\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path  d=\"M 111 19\nL 141 19\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><circle cx=\"126\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><text x=\"143\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path  d=\"M 234 19\nL 264 19\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><circle cx=\"249\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><text x=\"266\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><path  d=\"M 357 19\nL 387 19\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><circle cx=\"372\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><text x=\"389\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path  d=\"M 450 19\nL 480 19\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><circle cx=\"465\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><text x=\"482\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><text x=\"10\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line</text><text x=\"10\" y=\"52\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"10\" y=\"87\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.28k</text><text x=\"10\" y=\"122\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.12k</text><text x=\"22\" y=\"157\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"22\" y=\"192\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"22\" y=\"227\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"22\" y=\"262\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"22\" y=\"297\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"22\" y=\"332\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"40\" y=\"367\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 59 45\nL 590 45\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 80\nL 590 80\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 115\nL 590 115\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 150\nL 590 150\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 185\nL 590 185\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 220\nL 590 220\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 255\nL 590 255\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 290\nL 590 290\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 325\nL 590 325\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 365\nL 59 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 147 365\nL 147 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 236 365\nL 236 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 324 365\nL 324 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 413 365\nL 413 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 501 365\nL 501 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 590 365\nL 590 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 59 360\nL 590 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"59\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"147\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"236\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"324\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"413\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"501\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"563\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path  d=\"M 59 334\nL 147 332\nL 236 338\nL 324 331\nL 413 341\nL 501 310\nL 590 315\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:none\"/><circle cx=\"59\" cy=\"334\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"147\" cy=\"332\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"236\" cy=\"338\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"413\" cy=\"341\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"501\" cy=\"310\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"590\" cy=\"315\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 59 312\nL 147 321\nL 236 319\nL 324 309\nL 413 297\nL 501 288\nL 590 293\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:none\"/><circle cx=\"59\" cy=\"312\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"147\" cy=\"321\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"236\" cy=\"319\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"413\" cy=\"297\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"501\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"590\" cy=\"293\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 59 328\nL 147 310\nL 236 317\nL 324 327\nL 413 319\nL 501 288\nL 590 271\" style=\"stroke-width:2;stroke:rgba(250,200,88,1.0);fill:none\"/><circle cx=\"59\" cy=\"328\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"147\" cy=\"310\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"236\" cy=\"317\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"327\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"413\" cy=\"319\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"501\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"590\" cy=\"271\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 59 290\nL 147 288\nL 236 295\nL 324 287\nL 413 275\nL 501 288\nL 590 290\" style=\"stroke-width:2;stroke:rgba(238,102,102,1.0);fill:none\"/><circle cx=\"59\" cy=\"290\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"147\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"236\" cy=\"295\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"287\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"413\" cy=\"275\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"501\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"590\" cy=\"290\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 59 181\nL 147 157\nL 236 163\nL 324 156\nL 413 78\nL 501 70\nL 590 72\" style=\"stroke-width:2;stroke:rgba(115,192,222,1.0);fill:none\"/><circle cx=\"59\" cy=\"181\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"147\" cy=\"157\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"236\" cy=\"163\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"156\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"413\" cy=\"78\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"501\" cy=\"70\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"590\" cy=\"72\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/></svg>",
		},
		{
			name: "08YSkip1",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxisOptions = []YAxisOption{
					{
						LabelCount:     8,
						LabelSkipCount: 1,
					},
				}
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><text x=\"10\" y=\"17\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.4k</text><text x=\"23\" y=\"117\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1k</text><text x=\"13\" y=\"217\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">600</text><text x=\"13\" y=\"317\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">200</text><text x=\"31\" y=\"367\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 50 10\nL 590 10\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 50 60\nL 590 60\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 50 110\nL 590 110\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 50 160\nL 590 160\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 50 210\nL 590 210\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 50 260\nL 590 260\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 50 310\nL 590 310\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 88 330\nL 165 327\nL 242 335\nL 319 327\nL 396 338\nL 473 303\nL 551 308\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:none\"/><circle cx=\"88\" cy=\"330\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"165\" cy=\"327\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"242\" cy=\"335\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"319\" cy=\"327\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"396\" cy=\"338\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"473\" cy=\"303\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"551\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 88 155\nL 165 127\nL 242 135\nL 319 127\nL 396 38\nL 473 28\nL 551 30\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:none\"/><circle cx=\"88\" cy=\"155\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"165\" cy=\"127\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"242\" cy=\"135\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"319\" cy=\"127\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"396\" cy=\"38\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"473\" cy=\"28\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"551\" cy=\"30\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/></svg>",
		},
		{
			name: "09YSkip1",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxisOptions = []YAxisOption{
					{
						LabelCount:     9,
						LabelSkipCount: 1,
					},
				}
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><text x=\"10\" y=\"17\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"10\" y=\"104\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.08k</text><text x=\"22\" y=\"192\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">720</text><text x=\"22\" y=\"279\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">360</text><text x=\"40\" y=\"367\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 59 10\nL 590 10\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 53\nL 590 53\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 97\nL 590 97\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 141\nL 590 141\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 185\nL 590 185\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 228\nL 590 228\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 272\nL 590 272\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 316\nL 590 316\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 96 331\nL 172 328\nL 248 336\nL 324 328\nL 400 339\nL 476 305\nL 552 309\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:none\"/><circle cx=\"96\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"172\" cy=\"328\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"248\" cy=\"336\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"328\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"400\" cy=\"339\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"476\" cy=\"305\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 96 161\nL 172 134\nL 248 142\nL 324 133\nL 400 47\nL 476 37\nL 552 40\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:none\"/><circle cx=\"96\" cy=\"161\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"172\" cy=\"134\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"248\" cy=\"142\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"133\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"400\" cy=\"47\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"476\" cy=\"37\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"40\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/></svg>",
		},
		{
			name: "08YSkip2",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxisOptions = []YAxisOption{
					{
						LabelCount:     8,
						LabelSkipCount: 2,
					},
				}
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><text x=\"10\" y=\"17\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.4k</text><text x=\"13\" y=\"167\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"13\" y=\"317\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">200</text><text x=\"31\" y=\"367\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 50 10\nL 590 10\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 50 60\nL 590 60\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 50 110\nL 590 110\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 50 160\nL 590 160\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 50 210\nL 590 210\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 50 260\nL 590 260\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 50 310\nL 590 310\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 88 330\nL 165 327\nL 242 335\nL 319 327\nL 396 338\nL 473 303\nL 551 308\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:none\"/><circle cx=\"88\" cy=\"330\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"165\" cy=\"327\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"242\" cy=\"335\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"319\" cy=\"327\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"396\" cy=\"338\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"473\" cy=\"303\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"551\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 88 155\nL 165 127\nL 242 135\nL 319 127\nL 396 38\nL 473 28\nL 551 30\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:none\"/><circle cx=\"88\" cy=\"155\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"165\" cy=\"127\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"242\" cy=\"135\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"319\" cy=\"127\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"396\" cy=\"38\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"473\" cy=\"28\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"551\" cy=\"30\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/></svg>",
		},
		{
			name: "09YSkip2",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxisOptions = []YAxisOption{
					{
						LabelCount:     9,
						LabelSkipCount: 2,
					},
				}
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><text x=\"10\" y=\"17\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"22\" y=\"148\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">900</text><text x=\"22\" y=\"279\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">360</text><text x=\"40\" y=\"367\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 59 10\nL 590 10\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 53\nL 590 53\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 97\nL 590 97\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 141\nL 590 141\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 185\nL 590 185\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 228\nL 590 228\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 272\nL 590 272\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 316\nL 590 316\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 96 331\nL 172 328\nL 248 336\nL 324 328\nL 400 339\nL 476 305\nL 552 309\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:none\"/><circle cx=\"96\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"172\" cy=\"328\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"248\" cy=\"336\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"328\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"400\" cy=\"339\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"476\" cy=\"305\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 96 161\nL 172 134\nL 248 142\nL 324 133\nL 400 47\nL 476 37\nL 552 40\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:none\"/><circle cx=\"96\" cy=\"161\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"172\" cy=\"134\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"248\" cy=\"142\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"133\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"400\" cy=\"47\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"476\" cy=\"37\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"40\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/></svg>",
		},
		{
			name: "10YSkip2",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxisOptions = []YAxisOption{
					{
						LabelCount:     10,
						LabelSkipCount: 2,
					},
				}
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><text x=\"10\" y=\"17\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"22\" y=\"133\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"22\" y=\"250\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"40\" y=\"367\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 59 10\nL 590 10\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 48\nL 590 48\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 87\nL 590 87\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 126\nL 590 126\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 165\nL 590 165\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 204\nL 590 204\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 243\nL 590 243\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 282\nL 590 282\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 321\nL 590 321\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 96 331\nL 172 328\nL 248 336\nL 324 328\nL 400 339\nL 476 305\nL 552 309\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:none\"/><circle cx=\"96\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"172\" cy=\"328\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"248\" cy=\"336\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"328\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"400\" cy=\"339\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"476\" cy=\"305\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 96 161\nL 172 134\nL 248 142\nL 324 133\nL 400 47\nL 476 37\nL 552 40\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:none\"/><circle cx=\"96\" cy=\"161\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"172\" cy=\"134\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"248\" cy=\"142\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"133\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"400\" cy=\"47\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"476\" cy=\"37\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"40\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/></svg>",
		},
		{
			name: "08YSkip3",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxisOptions = []YAxisOption{
					{
						LabelCount:     8,
						LabelSkipCount: 3,
					},
				}
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><text x=\"10\" y=\"17\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.4k</text><text x=\"13\" y=\"217\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">600</text><text x=\"31\" y=\"367\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 50 10\nL 590 10\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 50 60\nL 590 60\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 50 110\nL 590 110\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 50 160\nL 590 160\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 50 210\nL 590 210\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 50 260\nL 590 260\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 50 310\nL 590 310\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 88 330\nL 165 327\nL 242 335\nL 319 327\nL 396 338\nL 473 303\nL 551 308\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:none\"/><circle cx=\"88\" cy=\"330\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"165\" cy=\"327\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"242\" cy=\"335\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"319\" cy=\"327\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"396\" cy=\"338\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"473\" cy=\"303\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"551\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 88 155\nL 165 127\nL 242 135\nL 319 127\nL 396 38\nL 473 28\nL 551 30\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:none\"/><circle cx=\"88\" cy=\"155\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"165\" cy=\"127\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"242\" cy=\"135\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"319\" cy=\"127\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"396\" cy=\"38\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"473\" cy=\"28\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"551\" cy=\"30\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/></svg>",
		},
		{
			name: "09YSkip3",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxisOptions = []YAxisOption{
					{
						LabelCount:     9,
						LabelSkipCount: 3,
					},
				}
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><text x=\"10\" y=\"17\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"22\" y=\"192\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">720</text><text x=\"40\" y=\"367\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 59 10\nL 590 10\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 53\nL 590 53\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 97\nL 590 97\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 141\nL 590 141\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 185\nL 590 185\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 228\nL 590 228\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 272\nL 590 272\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 316\nL 590 316\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 96 331\nL 172 328\nL 248 336\nL 324 328\nL 400 339\nL 476 305\nL 552 309\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:none\"/><circle cx=\"96\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"172\" cy=\"328\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"248\" cy=\"336\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"328\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"400\" cy=\"339\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"476\" cy=\"305\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 96 161\nL 172 134\nL 248 142\nL 324 133\nL 400 47\nL 476 37\nL 552 40\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:none\"/><circle cx=\"96\" cy=\"161\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"172\" cy=\"134\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"248\" cy=\"142\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"133\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"400\" cy=\"47\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"476\" cy=\"37\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"40\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/></svg>",
		},
		{
			name: "10YSkip3",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxisOptions = []YAxisOption{
					{
						LabelCount:     10,
						LabelSkipCount: 3,
					},
				}
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><text x=\"10\" y=\"17\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"22\" y=\"172\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"22\" y=\"328\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"40\" y=\"367\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 59 10\nL 590 10\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 48\nL 590 48\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 87\nL 590 87\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 126\nL 590 126\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 165\nL 590 165\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 204\nL 590 204\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 243\nL 590 243\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 282\nL 590 282\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 321\nL 590 321\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 96 331\nL 172 328\nL 248 336\nL 324 328\nL 400 339\nL 476 305\nL 552 309\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:none\"/><circle cx=\"96\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"172\" cy=\"328\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"248\" cy=\"336\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"328\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"400\" cy=\"339\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"476\" cy=\"305\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 96 161\nL 172 134\nL 248 142\nL 324 133\nL 400 47\nL 476 37\nL 552 40\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:none\"/><circle cx=\"96\" cy=\"161\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"172\" cy=\"134\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"248\" cy=\"142\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"133\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"400\" cy=\"47\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"476\" cy=\"37\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"40\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/></svg>",
		},
		{
			name: "11YSkip3",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxisOptions = []YAxisOption{
					{
						LabelCount:     11,
						LabelSkipCount: 3,
					},
				}
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><text x=\"19\" y=\"17\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.4k</text><text x=\"22\" y=\"157\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">840</text><text x=\"22\" y=\"297\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">280</text><text x=\"40\" y=\"367\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 59 10\nL 590 10\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 45\nL 590 45\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 80\nL 590 80\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 115\nL 590 115\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 150\nL 590 150\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 185\nL 590 185\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 220\nL 590 220\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 255\nL 590 255\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 290\nL 590 290\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 59 325\nL 590 325\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 96 330\nL 172 327\nL 248 335\nL 324 327\nL 400 338\nL 476 303\nL 552 308\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:none\"/><circle cx=\"96\" cy=\"330\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"172\" cy=\"327\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"248\" cy=\"335\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"327\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"400\" cy=\"338\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"476\" cy=\"303\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 96 155\nL 172 127\nL 248 135\nL 324 127\nL 400 38\nL 476 28\nL 552 30\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:none\"/><circle cx=\"96\" cy=\"155\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"172\" cy=\"127\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"248\" cy=\"135\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"127\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"400\" cy=\"38\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"476\" cy=\"28\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"552\" cy=\"30\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/></svg>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := NewPainter(PainterOptions{
				Type:   ChartOutputSVG,
				Width:  600,
				Height: 400,
			}, PainterThemeOption(GetTheme(ThemeLight)))
			require.NoError(t, err)
			opt := tt.makeOptions()

			_, err = NewLineChart(p, opt).Render()
			require.NoError(t, err)
			data, err := p.Bytes()
			require.NoError(t, err)
			assertEqualSVG(t, tt.result, string(data))
		})
	}
}
