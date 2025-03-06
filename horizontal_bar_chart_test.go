package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeBasicHorizontalBarChartOption() HorizontalBarChartOption {
	return HorizontalBarChartOption{
		Padding: NewBoxEqual(10),
		SeriesList: NewSeriesListHorizontalBar([][]float64{
			{18203, 23489, 29034, 104970, 131744, 630230},
			{19325, 23438, 31000, 121594, 134141, 681807},
		}),
		Title: TitleOption{
			Text: "World Population",
		},
		Legend: LegendOption{
			SeriesNames: []string{"2011", "2012"},
			Symbol:      SymbolDot,
		},
		YAxis: YAxisOption{
			Labels: []string{"Brazil", "Indonesia", "USA", "India", "China", "World"},
		},
	}
}

func makeMinimalHorizontalBarChartOption() HorizontalBarChartOption {
	opt := NewHorizontalBarChartOptionWithData([][]float64{
		{12, 24},
		{24, 48},
	})
	opt.YAxis = YAxisOption{
		Show:   Ptr(false),
		Labels: []string{"A", "B"},
	}
	opt.XAxis.Show = Ptr(false)
	return opt
}

func makeFullHorizontalBarChartStackedOption() HorizontalBarChartOption {
	seriesList := NewSeriesListHorizontalBar([][]float64{
		{4.9, 23.2, 25.6, 102.6, 142.2, 32.6, 20.0, 3.3},
		{19.0, 26.4, 28.7, 144.6, 122.2, 48.7, 28.8, 22.3},
		{80.0, 40.4, 28.4, 28.8, 24.4, 24.2, 40.8, 80.8},
	}, BarSeriesOption{
		Label: SeriesLabel{
			Show: Ptr(true),
			ValueFormatter: func(f float64) string {
				return strconv.Itoa(int(f))
			},
		},
	})
	dataLabels := []string{"A", "B", "C"}
	return HorizontalBarChartOption{
		Padding:     NewBoxEqual(20),
		SeriesList:  seriesList,
		StackSeries: Ptr(true),
		Legend: LegendOption{
			SeriesNames: dataLabels,
			Symbol:      SymbolDot,
		},
		YAxis: YAxisOption{
			Labels: []string{"1", "2", "3", "4", "5", "6", "7", "8"},
		},
	}
}

func TestNewHorizontalBarChartOptionWithData(t *testing.T) {
	t.Parallel()

	opt := NewHorizontalBarChartOptionWithData([][]float64{
		{12, 24},
		{24, 48},
	})

	assert.Len(t, opt.SeriesList, 2)
	assert.Equal(t, ChartTypeHorizontalBar, opt.SeriesList[0].getType())
	assert.Equal(t, defaultPadding, opt.Padding)

	p := NewPainter(PainterOptions{})
	assert.NoError(t, p.HorizontalBarChart(opt))
}

func TestHorizontalBarChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		themed      bool
		makeOptions func() HorizontalBarChartOption
		result      string
	}{
		{
			name:        "basic_themed",
			themed:      true,
			makeOptions: makeBasicHorizontalBarChartOption,
			result:      "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:rgb(40,40,40)\"/><path  d=\"M 224 19\nL 254 19\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:none\"/><circle cx=\"239\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><text x=\"256\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2011</text><path  d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:none\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2012</text><text x=\"10\" y=\"26\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World Population</text><path  d=\"M 87 46\nL 87 365\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 82 46\nL 87 46\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 82 99\nL 87 99\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 82 152\nL 87 152\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 82 205\nL 87 205\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 82 258\nL 87 258\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 82 311\nL 87 311\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 82 365\nL 87 365\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><text x=\"36\" y=\"78\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World</text><text x=\"37\" y=\"130\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">China</text><text x=\"43\" y=\"183\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">India</text><text x=\"47\" y=\"236\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">USA</text><text x=\"9\" y=\"289\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Indonesia</text><text x=\"38\" y=\"342\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Brazil</text><text x=\"87\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"187\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">144k</text><text x=\"287\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">288k</text><text x=\"388\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">432k</text><text x=\"488\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">576k</text><text x=\"555\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">720k</text><path  d=\"M 188 46\nL 188 362\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 288 46\nL 288 362\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 389 46\nL 389 362\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 489 46\nL 489 362\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 590 46\nL 590 362\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 88 321\nL 100 321\nL 100 335\nL 88 335\nL 88 321\" style=\"stroke:none;fill:rgb(255,100,100)\"/><path  d=\"M 88 268\nL 103 268\nL 103 282\nL 88 282\nL 88 268\" style=\"stroke:none;fill:rgb(255,100,100)\"/><path  d=\"M 88 215\nL 107 215\nL 107 229\nL 88 229\nL 88 215\" style=\"stroke:none;fill:rgb(255,100,100)\"/><path  d=\"M 88 162\nL 158 162\nL 158 176\nL 88 176\nL 88 162\" style=\"stroke:none;fill:rgb(255,100,100)\"/><path  d=\"M 88 109\nL 176 109\nL 176 123\nL 88 123\nL 88 109\" style=\"stroke:none;fill:rgb(255,100,100)\"/><path  d=\"M 88 56\nL 509 56\nL 509 70\nL 88 70\nL 88 56\" style=\"stroke:none;fill:rgb(255,100,100)\"/><path  d=\"M 88 340\nL 100 340\nL 100 354\nL 88 354\nL 88 340\" style=\"stroke:none;fill:rgb(255,210,100)\"/><path  d=\"M 88 287\nL 103 287\nL 103 301\nL 88 301\nL 88 287\" style=\"stroke:none;fill:rgb(255,210,100)\"/><path  d=\"M 88 234\nL 108 234\nL 108 248\nL 88 248\nL 88 234\" style=\"stroke:none;fill:rgb(255,210,100)\"/><path  d=\"M 88 181\nL 169 181\nL 169 195\nL 88 195\nL 88 181\" style=\"stroke:none;fill:rgb(255,210,100)\"/><path  d=\"M 88 128\nL 177 128\nL 177 142\nL 88 142\nL 88 128\" style=\"stroke:none;fill:rgb(255,210,100)\"/><path  d=\"M 88 75\nL 544 75\nL 544 89\nL 88 89\nL 88 75\" style=\"stroke:none;fill:rgb(255,210,100)\"/></svg>",
		},
		{
			name: "custom_fonts",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeBasicHorizontalBarChartOption()
				customFont := NewFontStyleWithSize(4.0).WithColor(ColorBlue)
				opt.Legend.FontStyle = customFont
				opt.XAxis.FontStyle = customFont
				opt.Title.FontStyle = customFont
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 247 19\nL 277 19\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"262\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"279\" y=\"25\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">2011</text><path  d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">2012</text><text x=\"10\" y=\"16\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">World Population</text><path  d=\"M 87 36\nL 87 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 82 36\nL 87 36\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 82 90\nL 87 90\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 82 145\nL 87 145\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 82 200\nL 87 200\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 82 255\nL 87 255\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 82 310\nL 87 310\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 82 365\nL 87 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"36\" y=\"69\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World</text><text x=\"37\" y=\"123\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">China</text><text x=\"43\" y=\"178\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">India</text><text x=\"47\" y=\"232\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">USA</text><text x=\"9\" y=\"287\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Indonesia</text><text x=\"38\" y=\"341\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Brazil</text><text x=\"87\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"187\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">144k</text><text x=\"287\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">288k</text><text x=\"388\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">432k</text><text x=\"488\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">576k</text><text x=\"578\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">720k</text><path  d=\"M 188 36\nL 188 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 288 36\nL 288 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 389 36\nL 389 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 489 36\nL 489 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 590 36\nL 590 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 88 320\nL 100 320\nL 100 334\nL 88 334\nL 88 320\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 265\nL 103 265\nL 103 279\nL 88 279\nL 88 265\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 210\nL 107 210\nL 107 224\nL 88 224\nL 88 210\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 155\nL 158 155\nL 158 169\nL 88 169\nL 88 155\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 100\nL 176 100\nL 176 114\nL 88 114\nL 88 100\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 46\nL 509 46\nL 509 60\nL 88 60\nL 88 46\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 339\nL 100 339\nL 100 353\nL 88 353\nL 88 339\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 284\nL 103 284\nL 103 298\nL 88 298\nL 88 284\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 229\nL 108 229\nL 108 243\nL 88 243\nL 88 229\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 174\nL 169 174\nL 169 188\nL 88 188\nL 88 174\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 119\nL 177 119\nL 177 133\nL 88 133\nL 88 119\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 65\nL 544 65\nL 544 79\nL 88 79\nL 88 65\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
		},
		{
			name: "value_labels",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeBasicHorizontalBarChartOption()
				series := opt.SeriesList
				for i := range series {
					series[i].Label.Show = Ptr(true)
				}
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 224 19\nL 254 19\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"239\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"256\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2011</text><path  d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2012</text><text x=\"10\" y=\"26\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World Population</text><path  d=\"M 87 46\nL 87 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 82 46\nL 87 46\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 82 99\nL 87 99\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 82 152\nL 87 152\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 82 205\nL 87 205\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 82 258\nL 87 258\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 82 311\nL 87 311\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 82 365\nL 87 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"36\" y=\"78\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World</text><text x=\"37\" y=\"130\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">China</text><text x=\"43\" y=\"183\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">India</text><text x=\"47\" y=\"236\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">USA</text><text x=\"9\" y=\"289\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Indonesia</text><text x=\"38\" y=\"342\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Brazil</text><text x=\"87\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"187\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">144k</text><text x=\"287\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">288k</text><text x=\"388\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">432k</text><text x=\"488\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">576k</text><text x=\"555\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">720k</text><path  d=\"M 188 46\nL 188 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 288 46\nL 288 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 389 46\nL 389 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 489 46\nL 489 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 590 46\nL 590 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 88 321\nL 100 321\nL 100 335\nL 88 335\nL 88 321\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 268\nL 103 268\nL 103 282\nL 88 282\nL 88 268\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 215\nL 107 215\nL 107 229\nL 88 229\nL 88 215\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 162\nL 158 162\nL 158 176\nL 88 176\nL 88 162\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 109\nL 176 109\nL 176 123\nL 88 123\nL 88 109\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 56\nL 509 56\nL 509 70\nL 88 70\nL 88 56\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 340\nL 100 340\nL 100 354\nL 88 354\nL 88 340\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 287\nL 103 287\nL 103 301\nL 88 301\nL 88 287\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 234\nL 108 234\nL 108 248\nL 88 248\nL 88 234\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 181\nL 169 181\nL 169 195\nL 88 195\nL 88 181\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 128\nL 177 128\nL 177 142\nL 88 142\nL 88 128\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 75\nL 544 75\nL 544 89\nL 88 89\nL 88 75\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"105\" y=\"332\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">18203</text><text x=\"108\" y=\"279\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">23489</text><text x=\"112\" y=\"226\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">29034</text><text x=\"163\" y=\"173\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">104970</text><text x=\"181\" y=\"120\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">131744</text><text x=\"514\" y=\"67\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">630230</text><text x=\"105\" y=\"351\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">19325</text><text x=\"108\" y=\"298\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">23438</text><text x=\"113\" y=\"245\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">31000</text><text x=\"174\" y=\"192\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">121594</text><text x=\"182\" y=\"139\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">134141</text><text x=\"549\" y=\"86\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">681807</text></svg>",
		},
		{
			name: "value_formatter",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeBasicHorizontalBarChartOption()
				series := opt.SeriesList
				for i := range series {
					series[i].Label.Show = Ptr(true)
				}
				opt.ValueFormatter = func(f float64) string {
					return "f"
				}
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 224 19\nL 254 19\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"239\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"256\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2011</text><path  d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2012</text><text x=\"10\" y=\"26\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World Population</text><path  d=\"M 87 46\nL 87 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 82 46\nL 87 46\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 82 99\nL 87 99\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 82 152\nL 87 152\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 82 205\nL 87 205\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 82 258\nL 87 258\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 82 311\nL 87 311\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 82 365\nL 87 365\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"36\" y=\"78\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World</text><text x=\"37\" y=\"130\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">China</text><text x=\"43\" y=\"183\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">India</text><text x=\"47\" y=\"236\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">USA</text><text x=\"9\" y=\"289\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Indonesia</text><text x=\"38\" y=\"342\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Brazil</text><text x=\"87\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"187\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"287\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"388\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"488\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"584\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><path  d=\"M 188 46\nL 188 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 288 46\nL 288 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 389 46\nL 389 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 489 46\nL 489 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 590 46\nL 590 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 88 321\nL 100 321\nL 100 335\nL 88 335\nL 88 321\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 268\nL 103 268\nL 103 282\nL 88 282\nL 88 268\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 215\nL 107 215\nL 107 229\nL 88 229\nL 88 215\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 162\nL 158 162\nL 158 176\nL 88 176\nL 88 162\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 109\nL 176 109\nL 176 123\nL 88 123\nL 88 109\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 56\nL 509 56\nL 509 70\nL 88 70\nL 88 56\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 88 340\nL 100 340\nL 100 354\nL 88 354\nL 88 340\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 287\nL 103 287\nL 103 301\nL 88 301\nL 88 287\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 234\nL 108 234\nL 108 248\nL 88 248\nL 88 234\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 181\nL 169 181\nL 169 195\nL 88 195\nL 88 181\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 128\nL 177 128\nL 177 142\nL 88 142\nL 88 128\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 88 75\nL 544 75\nL 544 89\nL 88 89\nL 88 75\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"105\" y=\"332\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">18203</text><text x=\"108\" y=\"279\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">23489</text><text x=\"112\" y=\"226\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">29034</text><text x=\"163\" y=\"173\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">104970</text><text x=\"181\" y=\"120\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">131744</text><text x=\"514\" y=\"67\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">630230</text><text x=\"105\" y=\"351\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">19325</text><text x=\"108\" y=\"298\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">23438</text><text x=\"113\" y=\"245\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">31000</text><text x=\"174\" y=\"192\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">121594</text><text x=\"182\" y=\"139\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">134141</text><text x=\"549\" y=\"86\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">681807</text></svg>",
		},
		{
			name: "bar_height_truncate",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeBasicHorizontalBarChartOption()
				opt.Title.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				opt.YAxis.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				opt.BarHeight = 1000
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 10 336\nL 24 336\nL 24 355\nL 10 355\nL 10 336\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 10 273\nL 28 273\nL 28 292\nL 10 292\nL 10 273\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 10 210\nL 32 210\nL 32 229\nL 10 229\nL 10 210\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 10 146\nL 91 146\nL 91 165\nL 10 165\nL 10 146\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 10 83\nL 111 83\nL 111 102\nL 10 102\nL 10 83\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 10 20\nL 497 20\nL 497 39\nL 10 39\nL 10 20\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 10 360\nL 24 360\nL 24 379\nL 10 379\nL 10 360\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 10 297\nL 28 297\nL 28 316\nL 10 316\nL 10 297\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 10 234\nL 33 234\nL 33 253\nL 10 253\nL 10 234\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 10 170\nL 104 170\nL 104 189\nL 10 189\nL 10 170\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 10 107\nL 113 107\nL 113 126\nL 10 126\nL 10 107\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 10 44\nL 537 44\nL 537 63\nL 10 63\nL 10 44\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
		},
		{
			name: "mark_line",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeBasicHorizontalBarChartOption()
				opt.SeriesList[0].MarkLine = NewMarkLine(SeriesMarkTypeMax, SeriesMarkTypeAverage)
				opt.YAxis.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"10\" y=\"26\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World Population</text><text x=\"9\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"125\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">144k</text><text x=\"241\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">288k</text><text x=\"357\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">432k</text><text x=\"473\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">576k</text><text x=\"555\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">720k</text><path  d=\"M 126 41\nL 126 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 242 41\nL 242 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 358 41\nL 358 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 474 41\nL 474 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 590 41\nL 590 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 10 321\nL 24 321\nL 24 335\nL 10 335\nL 10 321\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 10 267\nL 28 267\nL 28 281\nL 10 281\nL 10 267\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 10 213\nL 32 213\nL 32 227\nL 10 227\nL 10 213\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 10 159\nL 91 159\nL 91 173\nL 10 173\nL 10 159\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 10 105\nL 111 105\nL 111 119\nL 10 119\nL 10 105\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 10 51\nL 497 51\nL 497 65\nL 10 65\nL 10 51\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 10 340\nL 24 340\nL 24 354\nL 10 354\nL 10 340\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 10 286\nL 28 286\nL 28 300\nL 10 300\nL 10 286\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 10 232\nL 33 232\nL 33 246\nL 10 246\nL 10 232\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 10 178\nL 104 178\nL 104 192\nL 10 192\nL 10 178\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 10 124\nL 113 124\nL 113 138\nL 10 138\nL 10 124\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 10 70\nL 537 70\nL 537 84\nL 10 84\nL 10 70\" style=\"stroke:none;fill:rgb(145,204,117)\"/><circle cx=\"497\" cy=\"362\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 497 43\nL 497 365\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 492 59\nL 497 43\nL 502 59\nL 497 54\nL 492 59\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"473\" y=\"41\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">630.23k</text><circle cx=\"130\" cy=\"362\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 130 43\nL 130 365\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 125 59\nL 130 43\nL 135 59\nL 130 54\nL 125 59\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"106\" y=\"41\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">156.28k</text></svg>",
		},
		{
			name: "bar_height_thin",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeMinimalHorizontalBarChartOption()
				opt.BarHeight = 2
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 20 285\nL 98 285\nL 98 287\nL 20 287\nL 20 285\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 105\nL 232 105\nL 232 107\nL 20 107\nL 20 105\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 292\nL 232 292\nL 232 294\nL 20 294\nL 20 292\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 20 112\nL 501 112\nL 501 114\nL 20 114\nL 20 112\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
		},
		{
			name: "bar_margin_narrow",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeMinimalHorizontalBarChartOption()
				opt.BarMargin = Ptr(0.0)
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 20 210\nL 98 210\nL 98 290\nL 20 290\nL 20 210\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 30\nL 232 30\nL 232 110\nL 20 110\nL 20 30\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 290\nL 232 290\nL 232 370\nL 20 370\nL 20 290\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 20 110\nL 501 110\nL 501 190\nL 20 190\nL 20 110\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
		},
		{
			name: "bar_margin_wide",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeMinimalHorizontalBarChartOption()
				opt.BarMargin = Ptr(1000.0) // will be limited to fit graph
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 20 210\nL 98 210\nL 98 245\nL 20 245\nL 20 210\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 30\nL 232 30\nL 232 65\nL 20 65\nL 20 30\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 335\nL 232 335\nL 232 370\nL 20 370\nL 20 335\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 20 155\nL 501 155\nL 501 190\nL 20 190\nL 20 155\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
		},
		{
			name: "bar_height_and_narrow_margin",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeMinimalHorizontalBarChartOption()
				opt.BarHeight = 10
				opt.BarMargin = Ptr(0.0)
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 20 280\nL 98 280\nL 98 290\nL 20 290\nL 20 280\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 100\nL 232 100\nL 232 110\nL 20 110\nL 20 100\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 290\nL 232 290\nL 232 300\nL 20 300\nL 20 290\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 20 110\nL 501 110\nL 501 120\nL 20 120\nL 20 110\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
		},
		{
			name: "bar_height_and_wide_margin",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeMinimalHorizontalBarChartOption()
				opt.BarHeight = 10
				opt.BarMargin = Ptr(1000.0) // will be limited for readability
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 20 240\nL 98 240\nL 98 250\nL 20 250\nL 20 240\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 60\nL 232 60\nL 232 70\nL 20 70\nL 20 60\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 330\nL 232 330\nL 232 340\nL 20 340\nL 20 330\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 20 150\nL 501 150\nL 501 160\nL 20 160\nL 20 150\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
		},
		{
			name:        "stack_series",
			makeOptions: makeFullHorizontalBarChartStackedOption,
			result:      "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 217 29\nL 247 29\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"232\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"249\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">A</text><path  d=\"M 280 29\nL 310 29\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"295\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"312\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">B</text><path  d=\"M 342 29\nL 372 29\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"357\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"374\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">C</text><path  d=\"M 38 56\nL 38 355\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 33 56\nL 38 56\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 33 93\nL 38 93\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 33 130\nL 38 130\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 33 168\nL 38 168\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 33 205\nL 38 205\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 33 242\nL 38 242\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 33 280\nL 38 280\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 33 317\nL 38 317\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 33 355\nL 38 355\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"19\" y=\"80\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">8</text><text x=\"19\" y=\"117\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">7</text><text x=\"19\" y=\"154\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">6</text><text x=\"19\" y=\"191\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">5</text><text x=\"19\" y=\"228\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">4</text><text x=\"19\" y=\"265\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">3</text><text x=\"19\" y=\"302\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2</text><text x=\"19\" y=\"340\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1</text><text x=\"38\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"115\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">42.43</text><text x=\"192\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">84.86</text><text x=\"269\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">127.29</text><text x=\"347\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">169.71</text><text x=\"424\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">212.14</text><text x=\"501\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">254.57</text><text x=\"553\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">297</text><path  d=\"M 116 56\nL 116 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 193 56\nL 193 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 270 56\nL 270 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 348 56\nL 348 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 425 56\nL 425 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 502 56\nL 502 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 580 56\nL 580 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 39 322\nL 48 322\nL 48 349\nL 39 349\nL 39 322\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 39 285\nL 81 285\nL 81 312\nL 39 312\nL 39 285\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 39 247\nL 86 247\nL 86 274\nL 39 274\nL 39 247\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 39 210\nL 227 210\nL 227 237\nL 39 237\nL 39 210\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 39 173\nL 300 173\nL 300 200\nL 39 200\nL 39 173\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 39 135\nL 98 135\nL 98 162\nL 39 162\nL 39 135\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 39 98\nL 75 98\nL 75 125\nL 39 125\nL 39 98\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 39 61\nL 45 61\nL 45 88\nL 39 88\nL 39 61\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 48 322\nL 82 322\nL 82 349\nL 48 349\nL 48 322\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 81 285\nL 129 285\nL 129 312\nL 81 312\nL 81 285\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 86 247\nL 138 247\nL 138 274\nL 86 274\nL 86 247\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 227 210\nL 493 210\nL 493 237\nL 227 237\nL 227 210\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 300 173\nL 524 173\nL 524 200\nL 300 200\nL 300 173\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 98 135\nL 187 135\nL 187 162\nL 98 162\nL 98 135\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 75 98\nL 127 98\nL 127 125\nL 75 125\nL 75 98\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 45 61\nL 86 61\nL 86 88\nL 45 88\nL 45 61\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 82 322\nL 229 322\nL 229 349\nL 82 349\nL 82 322\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 129 285\nL 203 285\nL 203 312\nL 129 312\nL 129 285\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 138 247\nL 190 247\nL 190 274\nL 138 274\nL 138 247\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 493 210\nL 545 210\nL 545 237\nL 493 237\nL 493 210\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 524 173\nL 568 173\nL 568 200\nL 524 200\nL 524 173\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 187 135\nL 231 135\nL 231 162\nL 187 162\nL 187 135\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 127 98\nL 202 98\nL 202 125\nL 127 125\nL 127 98\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 86 61\nL 234 61\nL 234 88\nL 86 88\nL 86 61\" style=\"stroke:none;fill:rgb(250,200,88)\"/><text x=\"53\" y=\"339\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">4</text><text x=\"86\" y=\"302\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">23</text><text x=\"91\" y=\"264\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">25</text><text x=\"232\" y=\"227\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">102</text><text x=\"305\" y=\"190\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">142</text><text x=\"103\" y=\"152\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">32</text><text x=\"80\" y=\"115\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">20</text><text x=\"50\" y=\"78\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">3</text><text x=\"87\" y=\"339\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">19</text><text x=\"134\" y=\"302\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">26</text><text x=\"143\" y=\"264\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">28</text><text x=\"498\" y=\"227\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">144</text><text x=\"529\" y=\"190\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">122</text><text x=\"192\" y=\"152\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">48</text><text x=\"132\" y=\"115\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">28</text><text x=\"91\" y=\"78\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">22</text><text x=\"234\" y=\"339\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">80</text><text x=\"208\" y=\"302\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">40</text><text x=\"195\" y=\"264\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">28</text><text x=\"550\" y=\"227\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">28</text><text x=\"573\" y=\"190\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">24</text><text x=\"236\" y=\"152\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">24</text><text x=\"207\" y=\"115\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">40</text><text x=\"239\" y=\"78\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">80</text></svg>",
		},
		{
			name: "stack_series_with_mark",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeFullHorizontalBarChartStackedOption()
				opt.SeriesList[0].MarkLine = NewMarkLine(SeriesMarkTypeMax, SeriesMarkTypeAverage)
				opt.SeriesList[len(opt.SeriesList)-1].MarkLine = NewMarkLine(SeriesMarkTypeMax)
				opt.SeriesList[len(opt.SeriesList)-1].MarkLine.Lines[0].Global = true
				opt.YAxis.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"19\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"99\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">42.43</text><text x=\"179\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">84.86</text><text x=\"259\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">127.29</text><text x=\"339\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">169.71</text><text x=\"419\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">212.14</text><text x=\"499\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">254.57</text><text x=\"553\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">297</text><path  d=\"M 100 20\nL 100 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 180 20\nL 180 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 260 20\nL 260 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 340 20\nL 340 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 420 20\nL 420 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 500 20\nL 500 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 580 20\nL 580 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 20 318\nL 29 318\nL 29 349\nL 20 349\nL 20 318\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 276\nL 64 276\nL 64 307\nL 20 307\nL 20 276\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 234\nL 68 234\nL 68 265\nL 20 265\nL 20 234\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 192\nL 215 192\nL 215 223\nL 20 223\nL 20 192\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 150\nL 290 150\nL 290 181\nL 20 181\nL 20 150\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 108\nL 82 108\nL 82 139\nL 20 139\nL 20 108\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 66\nL 58 66\nL 58 97\nL 20 97\nL 20 66\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 20 25\nL 26 25\nL 26 56\nL 20 56\nL 20 25\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 29 318\nL 65 318\nL 65 349\nL 29 349\nL 29 318\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 64 276\nL 114 276\nL 114 307\nL 64 307\nL 64 276\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 68 234\nL 122 234\nL 122 265\nL 68 265\nL 68 234\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 215 192\nL 490 192\nL 490 223\nL 215 223\nL 215 192\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 290 150\nL 522 150\nL 522 181\nL 290 181\nL 290 150\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 82 108\nL 174 108\nL 174 139\nL 82 139\nL 82 108\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 58 66\nL 112 66\nL 112 97\nL 58 97\nL 58 66\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 26 25\nL 68 25\nL 68 56\nL 26 56\nL 26 25\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 65 318\nL 217 318\nL 217 349\nL 65 349\nL 65 318\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 114 276\nL 190 276\nL 190 307\nL 114 307\nL 114 276\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 122 234\nL 176 234\nL 176 265\nL 122 265\nL 122 234\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 490 192\nL 544 192\nL 544 223\nL 490 223\nL 490 192\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 522 150\nL 568 150\nL 568 181\nL 522 181\nL 522 150\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 174 108\nL 220 108\nL 220 139\nL 174 139\nL 174 108\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 112 66\nL 189 66\nL 189 97\nL 112 97\nL 112 66\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 68 25\nL 221 25\nL 221 56\nL 68 56\nL 68 25\" style=\"stroke:none;fill:rgb(250,200,88)\"/><circle cx=\"290\" cy=\"352\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 290 22\nL 290 355\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 285 38\nL 290 22\nL 295 38\nL 290 33\nL 285 38\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"278\" y=\"20\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">142</text><circle cx=\"104\" cy=\"352\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 104 22\nL 104 355\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 99 38\nL 104 22\nL 109 38\nL 104 33\nL 99 38\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"96\" y=\"20\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">44</text><circle cx=\"570\" cy=\"352\" r=\"3\" style=\"stroke-width:1;stroke:rgb(211,211,211);fill:rgb(211,211,211)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 570 22\nL 570 355\" style=\"stroke-width:1;stroke:rgb(211,211,211);fill:rgb(211,211,211)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 565 38\nL 570 22\nL 575 38\nL 570 33\nL 565 38\" style=\"stroke-width:1;stroke:rgb(211,211,211);fill:rgb(211,211,211)\"/><text x=\"558\" y=\"20\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">288</text><text x=\"34\" y=\"337\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">4</text><text x=\"69\" y=\"295\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">23</text><text x=\"73\" y=\"253\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">25</text><text x=\"220\" y=\"211\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">102</text><text x=\"295\" y=\"169\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">142</text><text x=\"87\" y=\"127\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">32</text><text x=\"63\" y=\"85\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">20</text><text x=\"31\" y=\"44\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">3</text><text x=\"70\" y=\"337\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">19</text><text x=\"119\" y=\"295\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">26</text><text x=\"127\" y=\"253\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">28</text><text x=\"495\" y=\"211\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">144</text><text x=\"527\" y=\"169\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">122</text><text x=\"179\" y=\"127\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">48</text><text x=\"117\" y=\"85\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">28</text><text x=\"73\" y=\"44\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">22</text><text x=\"222\" y=\"337\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">80</text><text x=\"195\" y=\"295\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">40</text><text x=\"181\" y=\"253\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">28</text><text x=\"549\" y=\"211\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">28</text><text x=\"573\" y=\"169\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">24</text><text x=\"225\" y=\"127\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">24</text><text x=\"194\" y=\"85\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">40</text><text x=\"226\" y=\"44\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">80</text></svg>",
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

				validateHorizontalBarChartRender(t, p, tt.makeOptions(), tt.result)
			})
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-options", func(t *testing.T) {
				p := NewPainter(painterOptions)
				opt := tt.makeOptions()
				opt.Theme = GetTheme(ThemeVividDark)

				validateHorizontalBarChartRender(t, p, opt, tt.result)
			})
		} else {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				p := NewPainter(painterOptions)

				validateHorizontalBarChartRender(t, p, tt.makeOptions(), tt.result)
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

func TestHorizontalBarChartError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		makeOptions      func() HorizontalBarChartOption
		errorMsgContains string
	}{
		{
			name: "empty_series",
			makeOptions: func() HorizontalBarChartOption {
				return NewHorizontalBarChartOptionWithData([][]float64{})
			},
			errorMsgContains: "empty series list",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			})

			err := p.HorizontalBarChart(tt.makeOptions())
			require.Error(t, err)
			require.ErrorContains(t, err, tt.errorMsgContains)
		})
	}
}
