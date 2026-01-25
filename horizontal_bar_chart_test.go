package charts

import (
	"strconv"
	"testing"

	"github.com/dustin/go-humanize"
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
	return HorizontalBarChartOption{
		Padding:     NewBoxEqual(20),
		SeriesList:  seriesList,
		StackSeries: Ptr(true),
		Legend: LegendOption{
			Symbol: SymbolDot,
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
		svg         string
		pngCRC      uint32
	}{
		{
			name:        "basic_themed",
			themed:      true,
			makeOptions: makeBasicHorizontalBarChartOption,
			svg:         "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:rgb(40,40,40)\"/><text x=\"10\" y=\"26\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World Population</text><path d=\"M 224 19\nL 254 19\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:none\"/><circle cx=\"239\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><text x=\"256\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2011</text><path d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:none\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2012</text><path d=\"M 87 46\nL 87 366\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path d=\"M 82 46\nL 87 46\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path d=\"M 82 99\nL 87 99\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path d=\"M 82 152\nL 87 152\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path d=\"M 82 206\nL 87 206\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path d=\"M 82 259\nL 87 259\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path d=\"M 82 312\nL 87 312\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path d=\"M 82 366\nL 87 366\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><text x=\"36\" y=\"78\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World</text><text x=\"37\" y=\"131\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">China</text><text x=\"43\" y=\"184\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">India</text><text x=\"47\" y=\"237\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">USA</text><text x=\"9\" y=\"290\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Indonesia</text><text x=\"38\" y=\"343\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Brazil</text><text x=\"87\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"170\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">120k</text><text x=\"254\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">240k</text><text x=\"338\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">360k</text><text x=\"421\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480k</text><text x=\"505\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">600k</text><text x=\"555\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">720k</text><path d=\"M 171 46\nL 171 362\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path d=\"M 255 46\nL 255 362\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path d=\"M 339 46\nL 339 362\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path d=\"M 422 46\nL 422 362\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path d=\"M 506 46\nL 506 362\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path d=\"M 590 46\nL 590 362\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path d=\"M 88 322\nL 100 322\nL 100 336\nL 88 336\nL 88 322\" style=\"stroke:none;fill:rgb(255,100,100)\"/><path d=\"M 88 269\nL 104 269\nL 104 283\nL 88 283\nL 88 269\" style=\"stroke:none;fill:rgb(255,100,100)\"/><path d=\"M 88 216\nL 108 216\nL 108 230\nL 88 230\nL 88 216\" style=\"stroke:none;fill:rgb(255,100,100)\"/><path d=\"M 88 162\nL 161 162\nL 161 176\nL 88 176\nL 88 162\" style=\"stroke:none;fill:rgb(255,100,100)\"/><path d=\"M 88 109\nL 179 109\nL 179 123\nL 88 123\nL 88 109\" style=\"stroke:none;fill:rgb(255,100,100)\"/><path d=\"M 88 56\nL 527 56\nL 527 70\nL 88 70\nL 88 56\" style=\"stroke:none;fill:rgb(255,100,100)\"/><path d=\"M 88 341\nL 101 341\nL 101 355\nL 88 355\nL 88 341\" style=\"stroke:none;fill:rgb(255,210,100)\"/><path d=\"M 88 288\nL 104 288\nL 104 302\nL 88 302\nL 88 288\" style=\"stroke:none;fill:rgb(255,210,100)\"/><path d=\"M 88 235\nL 109 235\nL 109 249\nL 88 249\nL 88 235\" style=\"stroke:none;fill:rgb(255,210,100)\"/><path d=\"M 88 181\nL 172 181\nL 172 195\nL 88 195\nL 88 181\" style=\"stroke:none;fill:rgb(255,210,100)\"/><path d=\"M 88 128\nL 181 128\nL 181 142\nL 88 142\nL 88 128\" style=\"stroke:none;fill:rgb(255,210,100)\"/><path d=\"M 88 75\nL 563 75\nL 563 89\nL 88 89\nL 88 75\" style=\"stroke:none;fill:rgb(255,210,100)\"/></svg>",
			pngCRC:      0xfb168012,
		},
		{
			name: "custom_fonts",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeBasicHorizontalBarChartOption()
				customFont := NewFontStyleWithSize(4.0).WithColor(ColorBlue)
				opt.Legend.FontStyle = customFont
				opt.XAxis.FontStyle = customFont
				opt.YAxis.FontStyle = customFont
				opt.Title.FontStyle = customFont
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"10\" y=\"16\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">World Population</text><path d=\"M 247 19\nL 277 19\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"262\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"279\" y=\"25\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">2011</text><path d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">2012</text><path d=\"M 42 36\nL 42 366\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 37 36\nL 42 36\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 37 91\nL 42 91\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 37 146\nL 42 146\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 37 201\nL 42 201\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 37 256\nL 42 256\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 37 311\nL 42 311\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 37 366\nL 42 366\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"18\" y=\"64\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">World</text><text x=\"18\" y=\"118\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">China</text><text x=\"20\" y=\"173\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">India</text><text x=\"22\" y=\"228\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">USA</text><text x=\"9\" y=\"282\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Indonesia</text><text x=\"19\" y=\"337\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Brazil</text><text x=\"42\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"102\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">80k</text><text x=\"163\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">160k</text><text x=\"224\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">240k</text><text x=\"285\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">320k</text><text x=\"345\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">400k</text><text x=\"406\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">480k</text><text x=\"467\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">560k</text><text x=\"528\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">640k</text><text x=\"578\" y=\"375\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">720k</text><path d=\"M 103 36\nL 103 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 164 36\nL 164 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 225 36\nL 225 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 286 36\nL 286 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 346 36\nL 346 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 407 36\nL 407 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 468 36\nL 468 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 529 36\nL 529 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 590 36\nL 590 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 43 321\nL 56 321\nL 56 336\nL 43 336\nL 43 321\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 43 266\nL 60 266\nL 60 281\nL 43 281\nL 43 266\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 43 211\nL 65 211\nL 65 226\nL 43 226\nL 43 211\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 43 156\nL 122 156\nL 122 171\nL 43 171\nL 43 156\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 43 101\nL 143 101\nL 143 116\nL 43 116\nL 43 101\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 43 46\nL 521 46\nL 521 61\nL 43 61\nL 43 46\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 43 341\nL 57 341\nL 57 356\nL 43 356\nL 43 341\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 43 286\nL 60 286\nL 60 301\nL 43 301\nL 43 286\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 43 231\nL 66 231\nL 66 246\nL 43 246\nL 43 231\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 43 176\nL 135 176\nL 135 191\nL 43 191\nL 43 176\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 43 121\nL 144 121\nL 144 136\nL 43 136\nL 43 121\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 43 66\nL 560 66\nL 560 81\nL 43 81\nL 43 66\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
			pngCRC: 0x7d773c06,
		},
		{
			name: "value_labels",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeBasicHorizontalBarChartOption()
				series := opt.SeriesList
				for i := range series {
					series[i].Label.Show = Ptr(true)
					series[i].Label.ValueFormatter = func(f float64) string {
						return humanize.FtoaWithDigits(f, 2)
					}
				}
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"10\" y=\"26\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World Population</text><path d=\"M 224 19\nL 254 19\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"239\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"256\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2011</text><path d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2012</text><path d=\"M 87 46\nL 87 366\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 82 46\nL 87 46\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 82 99\nL 87 99\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 82 152\nL 87 152\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 82 206\nL 87 206\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 82 259\nL 87 259\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 82 312\nL 87 312\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 82 366\nL 87 366\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"36\" y=\"78\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World</text><text x=\"37\" y=\"131\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">China</text><text x=\"43\" y=\"184\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">India</text><text x=\"47\" y=\"237\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">USA</text><text x=\"9\" y=\"290\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Indonesia</text><text x=\"38\" y=\"343\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Brazil</text><text x=\"87\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"170\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">120k</text><text x=\"254\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">240k</text><text x=\"338\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">360k</text><text x=\"421\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480k</text><text x=\"505\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">600k</text><text x=\"555\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">720k</text><path d=\"M 171 46\nL 171 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 255 46\nL 255 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 339 46\nL 339 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 422 46\nL 422 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 506 46\nL 506 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 590 46\nL 590 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 88 322\nL 100 322\nL 100 336\nL 88 336\nL 88 322\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 88 269\nL 104 269\nL 104 283\nL 88 283\nL 88 269\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 88 216\nL 108 216\nL 108 230\nL 88 230\nL 88 216\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 88 162\nL 161 162\nL 161 176\nL 88 176\nL 88 162\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 88 109\nL 179 109\nL 179 123\nL 88 123\nL 88 109\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 88 56\nL 527 56\nL 527 70\nL 88 70\nL 88 56\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 88 341\nL 101 341\nL 101 355\nL 88 355\nL 88 341\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 88 288\nL 104 288\nL 104 302\nL 88 302\nL 88 288\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 88 235\nL 109 235\nL 109 249\nL 88 249\nL 88 235\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 88 181\nL 172 181\nL 172 195\nL 88 195\nL 88 181\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 88 128\nL 181 128\nL 181 142\nL 88 142\nL 88 128\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 88 75\nL 563 75\nL 563 89\nL 88 89\nL 88 75\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"105\" y=\"333\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">18203</text><text x=\"109\" y=\"280\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">23489</text><text x=\"113\" y=\"227\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">29034</text><text x=\"166\" y=\"173\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">104970</text><text x=\"184\" y=\"120\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">131744</text><text x=\"532\" y=\"67\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">630230</text><text x=\"106\" y=\"352\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">19325</text><text x=\"109\" y=\"299\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">23438</text><text x=\"114\" y=\"246\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">31000</text><text x=\"177\" y=\"192\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">121594</text><text x=\"186\" y=\"139\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">134141</text><text x=\"556\" y=\"86\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">681807</text></svg>",
			pngCRC: 0xdda92f8,
		},
		{
			name: "value_formatter",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeBasicHorizontalBarChartOption()
				opt.ValueFormatter = func(f float64) string {
					return "f"
				}
				series := opt.SeriesList
				for i := range series {
					series[i].Label.Show = Ptr(true)
					series[i].Label.ValueFormatter = opt.ValueFormatter
				}
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"10\" y=\"26\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World Population</text><path d=\"M 224 19\nL 254 19\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"239\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"256\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2011</text><path d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2012</text><path d=\"M 87 46\nL 87 366\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 82 46\nL 87 46\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 82 99\nL 87 99\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 82 152\nL 87 152\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 82 206\nL 87 206\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 82 259\nL 87 259\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 82 312\nL 87 312\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 82 366\nL 87 366\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"36\" y=\"78\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World</text><text x=\"37\" y=\"131\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">China</text><text x=\"43\" y=\"184\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">India</text><text x=\"47\" y=\"237\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">USA</text><text x=\"9\" y=\"290\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Indonesia</text><text x=\"38\" y=\"343\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Brazil</text><text x=\"87\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"142\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"198\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"254\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"310\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"365\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"421\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"477\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"533\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"584\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">f</text><path d=\"M 143 46\nL 143 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 199 46\nL 199 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 255 46\nL 255 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 311 46\nL 311 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 366 46\nL 366 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 422 46\nL 422 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 478 46\nL 478 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 534 46\nL 534 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 590 46\nL 590 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 88 322\nL 100 322\nL 100 336\nL 88 336\nL 88 322\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 88 269\nL 104 269\nL 104 283\nL 88 283\nL 88 269\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 88 216\nL 108 216\nL 108 230\nL 88 230\nL 88 216\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 88 162\nL 161 162\nL 161 176\nL 88 176\nL 88 162\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 88 109\nL 179 109\nL 179 123\nL 88 123\nL 88 109\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 88 56\nL 527 56\nL 527 70\nL 88 70\nL 88 56\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 88 341\nL 101 341\nL 101 355\nL 88 355\nL 88 341\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 88 288\nL 104 288\nL 104 302\nL 88 302\nL 88 288\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 88 235\nL 109 235\nL 109 249\nL 88 249\nL 88 235\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 88 181\nL 172 181\nL 172 195\nL 88 195\nL 88 181\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 88 128\nL 181 128\nL 181 142\nL 88 142\nL 88 128\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 88 75\nL 563 75\nL 563 89\nL 88 89\nL 88 75\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"105\" y=\"333\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"109\" y=\"280\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"113\" y=\"227\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"166\" y=\"173\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"184\" y=\"120\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"532\" y=\"67\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"106\" y=\"352\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"109\" y=\"299\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"114\" y=\"246\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"177\" y=\"192\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"186\" y=\"139\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">f</text><text x=\"568\" y=\"86\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">f</text></svg>",
			pngCRC: 0xd72c917f,
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
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 10 336\nL 24 336\nL 24 355\nL 10 355\nL 10 336\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 10 273\nL 28 273\nL 28 292\nL 10 292\nL 10 273\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 10 210\nL 33 210\nL 33 229\nL 10 229\nL 10 210\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 10 146\nL 94 146\nL 94 165\nL 10 165\nL 10 146\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 10 83\nL 116 83\nL 116 102\nL 10 102\nL 10 83\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 10 20\nL 517 20\nL 517 39\nL 10 39\nL 10 20\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 10 360\nL 25 360\nL 25 379\nL 10 379\nL 10 360\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 10 297\nL 28 297\nL 28 316\nL 10 316\nL 10 297\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 10 234\nL 34 234\nL 34 253\nL 10 253\nL 10 234\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 10 170\nL 107 170\nL 107 189\nL 10 189\nL 10 170\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 10 107\nL 118 107\nL 118 126\nL 10 126\nL 10 107\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 10 44\nL 559 44\nL 559 63\nL 10 63\nL 10 44\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
			pngCRC: 0xb35224f4,
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
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"10\" y=\"26\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World Population</text><text x=\"9\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"105\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">120k</text><text x=\"202\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">240k</text><text x=\"299\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">360k</text><text x=\"395\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480k</text><text x=\"492\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">600k</text><text x=\"555\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">720k</text><path d=\"M 106 41\nL 106 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 203 41\nL 203 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 300 41\nL 300 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 396 41\nL 396 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 493 41\nL 493 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 590 41\nL 590 362\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 10 321\nL 24 321\nL 24 335\nL 10 335\nL 10 321\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 10 267\nL 28 267\nL 28 281\nL 10 281\nL 10 267\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 10 213\nL 33 213\nL 33 227\nL 10 227\nL 10 213\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 10 159\nL 94 159\nL 94 173\nL 10 173\nL 10 159\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 10 105\nL 116 105\nL 116 119\nL 10 119\nL 10 105\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 10 51\nL 517 51\nL 517 65\nL 10 65\nL 10 51\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 10 340\nL 25 340\nL 25 354\nL 10 354\nL 10 340\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 10 286\nL 28 286\nL 28 300\nL 10 300\nL 10 286\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 10 232\nL 34 232\nL 34 246\nL 10 246\nL 10 232\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 10 178\nL 107 178\nL 107 192\nL 10 192\nL 10 178\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 10 124\nL 118 124\nL 118 138\nL 10 138\nL 10 124\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 10 70\nL 559 70\nL 559 84\nL 10 84\nL 10 70\" style=\"stroke:none;fill:rgb(145,204,117)\"/><circle cx=\"517\" cy=\"363\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 517 43\nL 517 366\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 512 59\nL 517 43\nL 522 59\nL 517 54\nL 512 59\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"493\" y=\"41\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">630.23k</text><circle cx=\"135\" cy=\"363\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 135 43\nL 135 366\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 130 59\nL 135 43\nL 140 59\nL 135 54\nL 130 59\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"111\" y=\"41\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">156.28k</text></svg>",
			pngCRC: 0x56386ea2,
		},
		{
			name: "bar_height_thin",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeMinimalHorizontalBarChartOption()
				opt.BarHeight = 2
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 20 285\nL 98 285\nL 98 287\nL 20 287\nL 20 285\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 20 105\nL 232 105\nL 232 107\nL 20 107\nL 20 105\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 20 292\nL 232 292\nL 232 294\nL 20 294\nL 20 292\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 20 112\nL 501 112\nL 501 114\nL 20 114\nL 20 112\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
			pngCRC: 0xedf8c602,
		},
		{
			name: "bar_margin_narrow",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeMinimalHorizontalBarChartOption()
				opt.BarMargin = Ptr(0.0)
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 20 210\nL 98 210\nL 98 290\nL 20 290\nL 20 210\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 20 30\nL 232 30\nL 232 110\nL 20 110\nL 20 30\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 20 290\nL 232 290\nL 232 370\nL 20 370\nL 20 290\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 20 110\nL 501 110\nL 501 190\nL 20 190\nL 20 110\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
			pngCRC: 0x62e361ae,
		},
		{
			name: "bar_margin_wide",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeMinimalHorizontalBarChartOption()
				opt.BarMargin = Ptr(1000.0) // will be limited to fit graph
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 20 210\nL 98 210\nL 98 245\nL 20 245\nL 20 210\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 20 30\nL 232 30\nL 232 65\nL 20 65\nL 20 30\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 20 335\nL 232 335\nL 232 370\nL 20 370\nL 20 335\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 20 155\nL 501 155\nL 501 190\nL 20 190\nL 20 155\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
			pngCRC: 0xe22b0ccd,
		},
		{
			name: "bar_height_and_narrow_margin",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeMinimalHorizontalBarChartOption()
				opt.BarHeight = 10
				opt.BarMargin = Ptr(0.0)
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 20 280\nL 98 280\nL 98 290\nL 20 290\nL 20 280\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 20 100\nL 232 100\nL 232 110\nL 20 110\nL 20 100\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 20 290\nL 232 290\nL 232 300\nL 20 300\nL 20 290\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 20 110\nL 501 110\nL 501 120\nL 20 120\nL 20 110\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
			pngCRC: 0x8a6043ab,
		},
		{
			name: "bar_height_and_wide_margin",
			makeOptions: func() HorizontalBarChartOption {
				opt := makeMinimalHorizontalBarChartOption()
				opt.BarHeight = 10
				opt.BarMargin = Ptr(1000.0) // will be limited for readability
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 20 240\nL 98 240\nL 98 250\nL 20 250\nL 20 240\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 20 60\nL 232 60\nL 232 70\nL 20 70\nL 20 60\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 20 330\nL 232 330\nL 232 340\nL 20 340\nL 20 330\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 20 150\nL 501 150\nL 501 160\nL 20 160\nL 20 150\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
			pngCRC: 0x56436af8,
		},
		{
			name:        "stack_series",
			makeOptions: makeFullHorizontalBarChartStackedOption,
			svg:         "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 38 20\nL 38 356\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 33 20\nL 38 20\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 33 62\nL 38 62\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 33 104\nL 38 104\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 33 146\nL 38 146\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 33 188\nL 38 188\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 33 230\nL 38 230\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 33 272\nL 38 272\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 33 314\nL 38 314\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 33 356\nL 38 356\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"19\" y=\"46\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">8</text><text x=\"19\" y=\"88\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">7</text><text x=\"19\" y=\"130\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">6</text><text x=\"19\" y=\"172\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">5</text><text x=\"19\" y=\"213\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">4</text><text x=\"19\" y=\"255\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">3</text><text x=\"19\" y=\"297\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2</text><text x=\"19\" y=\"339\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1</text><text x=\"38\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"115\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">42</text><text x=\"192\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">84</text><text x=\"269\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">126</text><text x=\"347\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">168</text><text x=\"424\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">210</text><text x=\"501\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">252</text><text x=\"553\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">294</text><path d=\"M 116 20\nL 116 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 193 20\nL 193 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 270 20\nL 270 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 348 20\nL 348 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 425 20\nL 425 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 502 20\nL 502 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 580 20\nL 580 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 39 319\nL 48 319\nL 48 351\nL 39 351\nL 39 319\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 39 277\nL 81 277\nL 81 309\nL 39 309\nL 39 277\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 39 235\nL 86 235\nL 86 267\nL 39 267\nL 39 235\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 39 193\nL 227 193\nL 227 225\nL 39 225\nL 39 193\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 39 151\nL 300 151\nL 300 183\nL 39 183\nL 39 151\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 39 109\nL 98 109\nL 98 141\nL 39 141\nL 39 109\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 39 67\nL 75 67\nL 75 99\nL 39 99\nL 39 67\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 39 25\nL 45 25\nL 45 57\nL 39 57\nL 39 25\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 48 319\nL 82 319\nL 82 351\nL 48 351\nL 48 319\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 81 277\nL 129 277\nL 129 309\nL 81 309\nL 81 277\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 86 235\nL 138 235\nL 138 267\nL 86 267\nL 86 235\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 227 193\nL 493 193\nL 493 225\nL 227 225\nL 227 193\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 300 151\nL 524 151\nL 524 183\nL 300 183\nL 300 151\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 98 109\nL 187 109\nL 187 141\nL 98 141\nL 98 109\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 75 67\nL 127 67\nL 127 99\nL 75 99\nL 75 67\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 45 25\nL 86 25\nL 86 57\nL 45 57\nL 45 25\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 82 319\nL 229 319\nL 229 351\nL 82 351\nL 82 319\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path d=\"M 129 277\nL 203 277\nL 203 309\nL 129 309\nL 129 277\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path d=\"M 138 235\nL 190 235\nL 190 267\nL 138 267\nL 138 235\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path d=\"M 493 193\nL 545 193\nL 545 225\nL 493 225\nL 493 193\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path d=\"M 524 151\nL 568 151\nL 568 183\nL 524 183\nL 524 151\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path d=\"M 187 109\nL 231 109\nL 231 141\nL 187 141\nL 187 109\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path d=\"M 127 67\nL 202 67\nL 202 99\nL 127 99\nL 127 67\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path d=\"M 86 25\nL 234 25\nL 234 57\nL 86 57\nL 86 25\" style=\"stroke:none;fill:rgb(250,200,88)\"/><text x=\"53\" y=\"339\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">4</text><text x=\"86\" y=\"297\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">23</text><text x=\"91\" y=\"255\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">25</text><text x=\"232\" y=\"213\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">102</text><text x=\"305\" y=\"171\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">142</text><text x=\"103\" y=\"129\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">32</text><text x=\"80\" y=\"87\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">20</text><text x=\"50\" y=\"45\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">3</text><text x=\"87\" y=\"339\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">19</text><text x=\"134\" y=\"297\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">26</text><text x=\"143\" y=\"255\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">28</text><text x=\"498\" y=\"213\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">144</text><text x=\"529\" y=\"171\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">122</text><text x=\"192\" y=\"129\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">48</text><text x=\"132\" y=\"87\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">28</text><text x=\"91\" y=\"45\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">22</text><text x=\"234\" y=\"339\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">80</text><text x=\"208\" y=\"297\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">40</text><text x=\"195\" y=\"255\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">28</text><text x=\"550\" y=\"213\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">28</text><text x=\"573\" y=\"171\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">24</text><text x=\"236\" y=\"129\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">24</text><text x=\"207\" y=\"87\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">40</text><text x=\"239\" y=\"45\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">80</text></svg>",
			pngCRC:      0xf75e263c,
		},
		{
			name: "stack_series_simple",
			makeOptions: func() HorizontalBarChartOption {
				opt := NewHorizontalBarChartOptionWithData([][]float64{{4.0}, {1.0}})
				opt.StackSeries = Ptr(true)
				opt.XAxis.Unit = 1
				// disable extra
				opt.YAxis.Show = Ptr(false)
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"19\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"205\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2</text><text x=\"392\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">4</text><text x=\"571\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">6</text><path d=\"M 206 20\nL 206 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 393 20\nL 393 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 580 20\nL 580 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 20 30\nL 393 30\nL 393 346\nL 20 346\nL 20 30\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 393 30\nL 486 30\nL 486 346\nL 393 346\nL 393 30\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>",
			pngCRC: 0x2f4f3f65,
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
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"19\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"99\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">42</text><text x=\"179\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">84</text><text x=\"259\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">126</text><text x=\"339\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">168</text><text x=\"419\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">210</text><text x=\"499\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">252</text><text x=\"553\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">294</text><path d=\"M 100 20\nL 100 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 180 20\nL 180 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 260 20\nL 260 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 340 20\nL 340 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 420 20\nL 420 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 500 20\nL 500 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 580 20\nL 580 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 20 319\nL 29 319\nL 29 351\nL 20 351\nL 20 319\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 20 277\nL 64 277\nL 64 309\nL 20 309\nL 20 277\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 20 235\nL 68 235\nL 68 267\nL 20 267\nL 20 235\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 20 193\nL 215 193\nL 215 225\nL 20 225\nL 20 193\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 20 151\nL 290 151\nL 290 183\nL 20 183\nL 20 151\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 20 109\nL 82 109\nL 82 141\nL 20 141\nL 20 109\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 20 67\nL 58 67\nL 58 99\nL 20 99\nL 20 67\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 20 25\nL 26 25\nL 26 57\nL 20 57\nL 20 25\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 29 319\nL 65 319\nL 65 351\nL 29 351\nL 29 319\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 64 277\nL 114 277\nL 114 309\nL 64 309\nL 64 277\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 68 235\nL 122 235\nL 122 267\nL 68 267\nL 68 235\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 215 193\nL 490 193\nL 490 225\nL 215 225\nL 215 193\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 290 151\nL 522 151\nL 522 183\nL 290 183\nL 290 151\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 82 109\nL 174 109\nL 174 141\nL 82 141\nL 82 109\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 58 67\nL 112 67\nL 112 99\nL 58 99\nL 58 67\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 26 25\nL 68 25\nL 68 57\nL 26 57\nL 26 25\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 65 319\nL 217 319\nL 217 351\nL 65 351\nL 65 319\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path d=\"M 114 277\nL 190 277\nL 190 309\nL 114 309\nL 114 277\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path d=\"M 122 235\nL 176 235\nL 176 267\nL 122 267\nL 122 235\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path d=\"M 490 193\nL 544 193\nL 544 225\nL 490 225\nL 490 193\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path d=\"M 522 151\nL 568 151\nL 568 183\nL 522 183\nL 522 151\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path d=\"M 174 109\nL 220 109\nL 220 141\nL 174 141\nL 174 109\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path d=\"M 112 67\nL 189 67\nL 189 99\nL 112 99\nL 112 67\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path d=\"M 68 25\nL 221 25\nL 221 57\nL 68 57\nL 68 25\" style=\"stroke:none;fill:rgb(250,200,88)\"/><circle cx=\"290\" cy=\"353\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 290 22\nL 290 356\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 285 38\nL 290 22\nL 295 38\nL 290 33\nL 285 38\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"278\" y=\"20\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">142</text><circle cx=\"104\" cy=\"353\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 104 22\nL 104 356\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 99 38\nL 104 22\nL 109 38\nL 104 33\nL 99 38\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"96\" y=\"20\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">44</text><circle cx=\"570\" cy=\"353\" r=\"3\" style=\"stroke-width:1;stroke:rgb(211,211,211);fill:rgb(211,211,211)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 570 22\nL 570 356\" style=\"stroke-width:1;stroke:rgb(211,211,211);fill:rgb(211,211,211)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 565 38\nL 570 22\nL 575 38\nL 570 33\nL 565 38\" style=\"stroke-width:1;stroke:rgb(211,211,211);fill:rgb(211,211,211)\"/><text x=\"558\" y=\"20\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">288</text><text x=\"34\" y=\"339\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">4</text><text x=\"69\" y=\"297\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">23</text><text x=\"73\" y=\"255\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">25</text><text x=\"220\" y=\"213\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">102</text><text x=\"295\" y=\"171\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">142</text><text x=\"87\" y=\"129\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">32</text><text x=\"63\" y=\"87\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">20</text><text x=\"31\" y=\"45\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">3</text><text x=\"70\" y=\"339\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">19</text><text x=\"119\" y=\"297\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">26</text><text x=\"127\" y=\"255\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">28</text><text x=\"495\" y=\"213\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">144</text><text x=\"527\" y=\"171\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">122</text><text x=\"179\" y=\"129\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">48</text><text x=\"117\" y=\"87\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">28</text><text x=\"73\" y=\"45\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">22</text><text x=\"222\" y=\"339\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">80</text><text x=\"195\" y=\"297\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">40</text><text x=\"181\" y=\"255\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">28</text><text x=\"549\" y=\"213\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">28</text><text x=\"573\" y=\"171\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">24</text><text x=\"225\" y=\"129\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">24</text><text x=\"194\" y=\"87\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">40</text><text x=\"226\" y=\"45\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">80</text></svg>",
			pngCRC: 0x2e00befd,
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

				validateHorizontalBarChartRender(t, p, rp, tt.makeOptions(), tt.svg, tt.pngCRC)
			})
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-options", func(t *testing.T) {
				p := NewPainter(painterOptions)
				rp := NewPainter(rasterOptions)
				opt := tt.makeOptions()
				opt.Theme = GetTheme(ThemeVividDark)

				validateHorizontalBarChartRender(t, p, rp, opt, tt.svg, tt.pngCRC)
			})
		} else {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				p := NewPainter(painterOptions)
				rp := NewPainter(rasterOptions)

				validateHorizontalBarChartRender(t, p, rp, tt.makeOptions(), tt.svg, tt.pngCRC)
			})
		}
	}
}

func validateHorizontalBarChartRender(t *testing.T, svgP, pngP *Painter, opt HorizontalBarChartOption, expectedSVG string, expectedCRC uint32) {
	t.Helper()

	err := svgP.HorizontalBarChart(opt)
	require.NoError(t, err)
	data, err := svgP.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, expectedSVG, data)

	err = pngP.HorizontalBarChart(opt)
	require.NoError(t, err)
	rasterData, err := pngP.Bytes()
	require.NoError(t, err)
	assertEqualPNGCRC(t, expectedCRC, rasterData)
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
