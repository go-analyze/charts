package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-analyze/charts/chartdraw/drawing"
)

func makeBasicHorizontalBarChartOption() HorizontalBarChartOption {
	return HorizontalBarChartOption{
		Padding: Box{
			Top:    10,
			Right:  10,
			Bottom: 10,
			Left:   10,
		},
		SeriesList: NewSeriesListHorizontalBar([][]float64{
			{18203, 23489, 29034, 104970, 131744, 630230},
			{19325, 23438, 31000, 121594, 134141, 681807},
		}),
		Title: TitleOption{
			Text: "World Population",
		},
		Legend: LegendOption{
			Data: []string{
				"2011", "2012",
			},
		},
		YAxis: []YAxisOption{
			{
				Data: []string{
					"Brazil", "Indonesia", "USA", "India", "China", "World",
				},
			},
		},
	}
}

func makeMinimalHorizontalBarChartOption() HorizontalBarChartOption {
	opt := NewHorizontalBarChartOptionWithData([][]float64{
		{12, 24},
		{24, 48},
	})
	opt.YAxis = []YAxisOption{
		{
			Show: False(),
			Data: []string{
				"A", "B",
			},
		},
	}
	opt.XAxis.Show = False()
	return opt
}

func TestNewHorizontalBarChartOptionWithData(t *testing.T) {
	t.Parallel()

	opt := NewHorizontalBarChartOptionWithData([][]float64{
		{12, 24},
		{24, 48},
	})

	assert.Len(t, opt.SeriesList, 2)
	assert.Equal(t, ChartTypeHorizontalBar, opt.SeriesList[0].Type)
	assert.Len(t, opt.YAxis, 1)
	assert.Equal(t, defaultPadding, opt.Padding)

	p := NewPainter(PainterOptions{})
	assert.NoError(t, p.HorizontalBarChart(opt))
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
			result:       "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 224 19\nL 254 19\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"239\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"256\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2011</text><path  d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2012</text><text x=\"10\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World Population</text><path  d=\"M 83 45\nL 88 45\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 97\nL 88 97\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 150\nL 88 150\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 202\nL 88 202\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 255\nL 88 255\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 307\nL 88 307\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 360\nL 88 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 88 45\nL 88 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"37\" y=\"78\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World</text><text x=\"38\" y=\"130\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">China</text><text x=\"44\" y=\"183\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">India</text><text x=\"48\" y=\"235\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">USA</text><text x=\"10\" y=\"288\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Indonesia</text><text x=\"39\" y=\"340\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Brazil</text><text x=\"87\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"187\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">144k</text><text x=\"287\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">288k</text><text x=\"388\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">432k</text><text x=\"488\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">576k</text><text x=\"555\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">720k</text><path  d=\"M 188 45\nL 188 360\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 288 45\nL 288 360\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 389 45\nL 389 360\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 489 45\nL 489 360\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 590 45\nL 590 360\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 88 317\nL 100 317\nL 100 330\nL 88 330\nL 88 317\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 265\nL 103 265\nL 103 278\nL 88 278\nL 88 265\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 212\nL 107 212\nL 107 225\nL 88 225\nL 88 212\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 160\nL 158 160\nL 158 173\nL 88 173\nL 88 160\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 107\nL 176 107\nL 176 120\nL 88 120\nL 88 107\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 55\nL 509 55\nL 509 68\nL 88 68\nL 88 55\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 335\nL 100 335\nL 100 348\nL 88 348\nL 88 335\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 283\nL 103 283\nL 103 296\nL 88 296\nL 88 283\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 230\nL 108 230\nL 108 243\nL 88 243\nL 88 230\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 178\nL 169 178\nL 169 191\nL 88 191\nL 88 178\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 125\nL 177 125\nL 177 138\nL 88 138\nL 88 125\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 73\nL 544 73\nL 544 86\nL 88 86\nL 88 73\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
		},
		{
			name:         "themed",
			defaultTheme: false,
			makeOptions:  makeBasicHorizontalBarChartOption,
			result:       "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:rgb(40,40,40)\"/><path  d=\"M 224 19\nL 254 19\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:none\"/><circle cx=\"239\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><text x=\"256\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2011</text><path  d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:none\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2012</text><text x=\"10\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World Population</text><path  d=\"M 83 45\nL 88 45\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 83 97\nL 88 97\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 83 150\nL 88 150\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 83 202\nL 88 202\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 83 255\nL 88 255\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 83 307\nL 88 307\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 83 360\nL 88 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 88 45\nL 88 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><text x=\"37\" y=\"78\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World</text><text x=\"38\" y=\"130\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">China</text><text x=\"44\" y=\"183\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">India</text><text x=\"48\" y=\"235\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">USA</text><text x=\"10\" y=\"288\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Indonesia</text><text x=\"39\" y=\"340\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Brazil</text><text x=\"87\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"187\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">144k</text><text x=\"287\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">288k</text><text x=\"388\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">432k</text><text x=\"488\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">576k</text><text x=\"555\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">720k</text><path  d=\"M 188 45\nL 188 360\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 288 45\nL 288 360\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 389 45\nL 389 360\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 489 45\nL 489 360\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 590 45\nL 590 360\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 88 317\nL 100 317\nL 100 330\nL 88 330\nL 88 317\" style=\"stroke:none;fill:rgb(255,100,100)\"/><path  d=\"M 88 265\nL 103 265\nL 103 278\nL 88 278\nL 88 265\" style=\"stroke:none;fill:rgb(255,100,100)\"/><path  d=\"M 88 212\nL 107 212\nL 107 225\nL 88 225\nL 88 212\" style=\"stroke:none;fill:rgb(255,100,100)\"/><path  d=\"M 88 160\nL 158 160\nL 158 173\nL 88 173\nL 88 160\" style=\"stroke:none;fill:rgb(255,100,100)\"/><path  d=\"M 88 107\nL 176 107\nL 176 120\nL 88 120\nL 88 107\" style=\"stroke:none;fill:rgb(255,100,100)\"/><path  d=\"M 88 55\nL 509 55\nL 509 68\nL 88 68\nL 88 55\" style=\"stroke:none;fill:rgb(255,100,100)\"/><path  d=\"M 88 335\nL 100 335\nL 100 348\nL 88 348\nL 88 335\" style=\"stroke:none;fill:rgb(255,210,100)\"/><path  d=\"M 88 283\nL 103 283\nL 103 296\nL 88 296\nL 88 283\" style=\"stroke:none;fill:rgb(255,210,100)\"/><path  d=\"M 88 230\nL 108 230\nL 108 243\nL 88 243\nL 88 230\" style=\"stroke:none;fill:rgb(255,210,100)\"/><path  d=\"M 88 178\nL 169 178\nL 169 191\nL 88 191\nL 88 178\" style=\"stroke:none;fill:rgb(255,210,100)\"/><path  d=\"M 88 125\nL 177 125\nL 177 138\nL 88 138\nL 88 125\" style=\"stroke:none;fill:rgb(255,210,100)\"/><path  d=\"M 88 73\nL 544 73\nL 544 86\nL 88 86\nL 88 73\" style=\"stroke:none;fill:rgb(255,210,100)\"/></svg>",
		},
		{
			name:         "custom_fonts",
			defaultTheme: true,
			makeOptions: func() HorizontalBarChartOption {
				opt := makeBasicHorizontalBarChartOption()
				customFont := FontStyle{
					FontSize:  4.0,
					FontColor: drawing.ColorBlue,
				}
				opt.Legend.FontStyle = customFont
				opt.XAxis.FontStyle = customFont
				opt.Title.FontStyle = customFont
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 247 19\nL 277 19\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"262\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"279\" y=\"25\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">2011</text><path  d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">2012</text><text x=\"10\" y=\"15\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">World Population</text><path  d=\"M 83 35\nL 88 35\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 89\nL 88 89\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 143\nL 88 143\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 197\nL 88 197\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 251\nL 88 251\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 305\nL 88 305\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 360\nL 88 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 88 35\nL 88 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"37\" y=\"69\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World</text><text x=\"38\" y=\"123\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">China</text><text x=\"44\" y=\"177\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">India</text><text x=\"48\" y=\"231\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">USA</text><text x=\"10\" y=\"285\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Indonesia</text><text x=\"39\" y=\"339\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Brazil</text><text x=\"87\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"187\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">144k</text><text x=\"287\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">288k</text><text x=\"388\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">432k</text><text x=\"488\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">576k</text><text x=\"578\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">720k</text><path  d=\"M 188 35\nL 188 360\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 288 35\nL 288 360\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 389 35\nL 389 360\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 489 35\nL 489 360\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 590 35\nL 590 360\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 88 315\nL 100 315\nL 100 329\nL 88 329\nL 88 315\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 261\nL 103 261\nL 103 275\nL 88 275\nL 88 261\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 207\nL 107 207\nL 107 221\nL 88 221\nL 88 207\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 153\nL 158 153\nL 158 167\nL 88 167\nL 88 153\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 99\nL 176 99\nL 176 113\nL 88 113\nL 88 99\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 45\nL 509 45\nL 509 59\nL 88 59\nL 88 45\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 334\nL 100 334\nL 100 348\nL 88 348\nL 88 334\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 280\nL 103 280\nL 103 294\nL 88 294\nL 88 280\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 226\nL 108 226\nL 108 240\nL 88 240\nL 88 226\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 172\nL 169 172\nL 169 186\nL 88 186\nL 88 172\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 118\nL 177 118\nL 177 132\nL 88 132\nL 88 118\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 64\nL 544 64\nL 544 78\nL 88 78\nL 88 64\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
		},
		{
			name:         "value_labels",
			defaultTheme: true,
			makeOptions: func() HorizontalBarChartOption {
				opt := makeBasicHorizontalBarChartOption()
				series := opt.SeriesList
				for i := range series {
					series[i].Label.Show = True()
				}
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 224 19\nL 254 19\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"239\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"256\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2011</text><path  d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2012</text><text x=\"10\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World Population</text><path  d=\"M 83 45\nL 88 45\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 97\nL 88 97\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 150\nL 88 150\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 202\nL 88 202\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 255\nL 88 255\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 307\nL 88 307\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 360\nL 88 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 88 45\nL 88 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"37\" y=\"78\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World</text><text x=\"38\" y=\"130\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">China</text><text x=\"44\" y=\"183\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">India</text><text x=\"48\" y=\"235\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">USA</text><text x=\"10\" y=\"288\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Indonesia</text><text x=\"39\" y=\"340\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Brazil</text><text x=\"87\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"187\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">144k</text><text x=\"287\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">288k</text><text x=\"388\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">432k</text><text x=\"488\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">576k</text><text x=\"555\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">720k</text><path  d=\"M 188 45\nL 188 360\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 288 45\nL 288 360\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 389 45\nL 389 360\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 489 45\nL 489 360\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 590 45\nL 590 360\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 88 317\nL 100 317\nL 100 330\nL 88 330\nL 88 317\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 265\nL 103 265\nL 103 278\nL 88 278\nL 88 265\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 212\nL 107 212\nL 107 225\nL 88 225\nL 88 212\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 160\nL 158 160\nL 158 173\nL 88 173\nL 88 160\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 107\nL 176 107\nL 176 120\nL 88 120\nL 88 107\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 55\nL 509 55\nL 509 68\nL 88 68\nL 88 55\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 335\nL 100 335\nL 100 348\nL 88 348\nL 88 335\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 283\nL 103 283\nL 103 296\nL 88 296\nL 88 283\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 230\nL 108 230\nL 108 243\nL 88 243\nL 88 230\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 178\nL 169 178\nL 169 191\nL 88 191\nL 88 178\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 125\nL 177 125\nL 177 138\nL 88 138\nL 88 125\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 73\nL 544 73\nL 544 86\nL 88 86\nL 88 73\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"105\" y=\"327\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">18203</text><text x=\"108\" y=\"275\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">23489</text><text x=\"112\" y=\"222\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">29034</text><text x=\"163\" y=\"170\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">104970</text><text x=\"181\" y=\"117\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">131744</text><text x=\"514\" y=\"65\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">630230</text><text x=\"105\" y=\"345\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">19325</text><text x=\"108\" y=\"293\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">23438</text><text x=\"113\" y=\"240\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">31000</text><text x=\"174\" y=\"188\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">121594</text><text x=\"182\" y=\"135\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">134141</text><text x=\"549\" y=\"83\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">681807</text></svg>",
		},
		{
			name:         "value_formatter",
			defaultTheme: true,
			makeOptions: func() HorizontalBarChartOption {
				opt := makeBasicHorizontalBarChartOption()
				series := opt.SeriesList
				for i := range series {
					series[i].Label.Show = True()
				}
				opt.ValueFormatter = func(f float64) string {
					return "f"
				}
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 224 19\nL 254 19\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"239\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"256\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2011</text><path  d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2012</text><text x=\"10\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World Population</text><path  d=\"M 83 45\nL 88 45\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 97\nL 88 97\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 150\nL 88 150\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 202\nL 88 202\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 255\nL 88 255\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 307\nL 88 307\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 83 360\nL 88 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 88 45\nL 88 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"37\" y=\"78\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World</text><text x=\"38\" y=\"130\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">China</text><text x=\"44\" y=\"183\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">India</text><text x=\"48\" y=\"235\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">USA</text><text x=\"10\" y=\"288\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Indonesia</text><text x=\"39\" y=\"340\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Brazil</text><text x=\"87\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"187\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"287\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"388\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"488\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"584\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><path  d=\"M 188 45\nL 188 360\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 288 45\nL 288 360\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 389 45\nL 389 360\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 489 45\nL 489 360\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 590 45\nL 590 360\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 88 317\nL 100 317\nL 100 330\nL 88 330\nL 88 317\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 265\nL 103 265\nL 103 278\nL 88 278\nL 88 265\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 212\nL 107 212\nL 107 225\nL 88 225\nL 88 212\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 160\nL 158 160\nL 158 173\nL 88 173\nL 88 160\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 107\nL 176 107\nL 176 120\nL 88 120\nL 88 107\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 55\nL 509 55\nL 509 68\nL 88 68\nL 88 55\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 335\nL 100 335\nL 100 348\nL 88 348\nL 88 335\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 283\nL 103 283\nL 103 296\nL 88 296\nL 88 283\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 230\nL 108 230\nL 108 243\nL 88 243\nL 88 230\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 178\nL 169 178\nL 169 191\nL 88 191\nL 88 178\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 125\nL 177 125\nL 177 138\nL 88 138\nL 88 125\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 73\nL 544 73\nL 544 86\nL 88 86\nL 88 73\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"105\" y=\"327\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">18203</text><text x=\"108\" y=\"275\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">23489</text><text x=\"112\" y=\"222\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">29034</text><text x=\"163\" y=\"170\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">104970</text><text x=\"181\" y=\"117\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">131744</text><text x=\"514\" y=\"65\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">630230</text><text x=\"105\" y=\"345\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">19325</text><text x=\"108\" y=\"293\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">23438</text><text x=\"113\" y=\"240\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">31000</text><text x=\"174\" y=\"188\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">121594</text><text x=\"182\" y=\"135\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">134141</text><text x=\"549\" y=\"83\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">681807</text></svg>",
		},
		{
			name:         "bar_height_truncate",
			defaultTheme: true,
			makeOptions: func() HorizontalBarChartOption {
				opt := makeBasicHorizontalBarChartOption()
				opt.Title.Show = False()
				opt.XAxis.Show = False()
				opt.YAxis[0].Show = False()
				opt.Legend.Show = False()
				opt.BarHeight = 1000
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 10 311\nL 24 311\nL 24 327\nL 10 327\nL 10 311\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 10 253\nL 28 253\nL 28 269\nL 10 269\nL 10 253\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 10 195\nL 32 195\nL 32 211\nL 10 211\nL 10 195\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 10 136\nL 91 136\nL 91 152\nL 10 152\nL 10 136\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 10 78\nL 111 78\nL 111 94\nL 10 94\nL 10 78\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 10 20\nL 497 20\nL 497 36\nL 10 36\nL 10 20\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 10 332\nL 24 332\nL 24 348\nL 10 348\nL 10 332\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 10 274\nL 28 274\nL 28 290\nL 10 290\nL 10 274\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 10 216\nL 33 216\nL 33 232\nL 10 232\nL 10 216\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 10 157\nL 104 157\nL 104 173\nL 10 173\nL 10 157\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 10 99\nL 113 99\nL 113 115\nL 10 115\nL 10 99\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 10 41\nL 537 41\nL 537 57\nL 10 57\nL 10 41\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
		},
		{
			name:         "bar_height_thin",
			defaultTheme: true,
			makeOptions: func() HorizontalBarChartOption {
				opt := makeMinimalHorizontalBarChartOption()
				opt.BarHeight = 2
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 20 263\nL 98 263\nL 98 265\nL 20 265\nL 20 263\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 98\nL 232 98\nL 232 100\nL 20 100\nL 20 98\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 270\nL 232 270\nL 232 272\nL 20 272\nL 20 270\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 20 105\nL 501 105\nL 501 107\nL 20 107\nL 20 105\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
		},
		{
			name:         "bar_margin_narrow",
			defaultTheme: true,
			makeOptions: func() HorizontalBarChartOption {
				opt := makeMinimalHorizontalBarChartOption()
				opt.BarMargin = FloatPointer(0)
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 20 195\nL 98 195\nL 98 267\nL 20 267\nL 20 195\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 30\nL 232 30\nL 232 102\nL 20 102\nL 20 30\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 267\nL 232 267\nL 232 339\nL 20 339\nL 20 267\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 20 102\nL 501 102\nL 501 174\nL 20 174\nL 20 102\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
		},
		{
			name:         "bar_margin_wide",
			defaultTheme: true,
			makeOptions: func() HorizontalBarChartOption {
				opt := makeMinimalHorizontalBarChartOption()
				opt.BarMargin = FloatPointer(1000) // will be limited to fit graph
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 20 195\nL 98 195\nL 98 226\nL 20 226\nL 20 195\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 30\nL 232 30\nL 232 61\nL 20 61\nL 20 30\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 308\nL 232 308\nL 232 339\nL 20 339\nL 20 308\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 20 143\nL 501 143\nL 501 174\nL 20 174\nL 20 143\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
		},
		{
			name:         "bar_height_and_narrow_margin",
			defaultTheme: true,
			makeOptions: func() HorizontalBarChartOption {
				opt := makeMinimalHorizontalBarChartOption()
				opt.BarHeight = 10
				opt.BarMargin = FloatPointer(0)
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 20 257\nL 98 257\nL 98 267\nL 20 267\nL 20 257\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 92\nL 232 92\nL 232 102\nL 20 102\nL 20 92\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 267\nL 232 267\nL 232 277\nL 20 277\nL 20 267\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 20 102\nL 501 102\nL 501 112\nL 20 112\nL 20 102\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
		},
		{
			name:         "bar_height_and_wide_margin",
			defaultTheme: true,
			makeOptions: func() HorizontalBarChartOption {
				opt := makeMinimalHorizontalBarChartOption()
				opt.BarHeight = 10
				opt.BarMargin = FloatPointer(1000) // will be limited for readability
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 20 221\nL 98 221\nL 98 231\nL 20 231\nL 20 221\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 56\nL 232 56\nL 232 66\nL 20 66\nL 20 56\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 303\nL 232 303\nL 232 313\nL 20 313\nL 20 303\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 20 138\nL 501 138\nL 501 148\nL 20 148\nL 20 138\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
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
				p := NewPainter(painterOptions)

				validateHorizontalBarChartRender(t, p, tt.makeOptions(), tt.result)
			})
		} else {
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-painter", func(t *testing.T) {
				p := NewPainter(painterOptions, PainterThemeOption(GetTheme(ThemeVividDark)))

				validateHorizontalBarChartRender(t, p, tt.makeOptions(), tt.result)
			})
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-options", func(t *testing.T) {
				p := NewPainter(painterOptions)
				opt := tt.makeOptions()
				opt.Theme = GetTheme(ThemeVividDark)

				validateHorizontalBarChartRender(t, p, opt, tt.result)
			})
		}
	}
}

func validateHorizontalBarChartRender(t *testing.T, p *Painter, opt HorizontalBarChartOption, expectedResult string) {
	t.Helper()

	err := p.HorizontalBarChart(opt)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, expectedResult, data)
}
