package charts

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeBasicDoughnutChartOption() DoughnutChartOption {
	values := []float64{
		1048, 735, 580, 484, 300,
	}
	return DoughnutChartOption{
		SeriesList: NewSeriesListDoughnut(values),
		Title: TitleOption{
			Text:    "Title",
			Subtext: "Sub",
			Offset:  OffsetCenter,
		},
		Padding: NewBoxEqual(20),
		Legend: LegendOption{
			Vertical:    Ptr(true),
			SeriesNames: []string{"Series-A", "Series-B", "Series-C", "Series-D", "Series-E"},
			Offset:      OffsetLeft,
			Symbol:      SymbolDot,
		},
	}
}

func makeMinimalDoughnutChartOption() DoughnutChartOption {
	opt := makeBasicDoughnutChartOption()
	// disable extras
	for i := range opt.SeriesList {
		opt.SeriesList[i].Label.Show = Ptr(false)
	}
	opt.Title.Show = Ptr(false)
	opt.Legend.Show = Ptr(false)
	return opt
}

func TestNewDoughnutChartOptionWithData(t *testing.T) {
	t.Parallel()

	opt := NewDoughnutChartOptionWithData([]float64{12, 24, 48})

	assert.Len(t, opt.SeriesList, 3)
	assert.Equal(t, ChartTypeDoughnut, opt.SeriesList[0].getType())
	assert.Equal(t, defaultPadding, opt.Padding)

	p := NewPainter(PainterOptions{})
	assert.NoError(t, p.DoughnutChart(opt))
}

func TestDoughnutChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		themed      bool
		makeOptions func() DoughnutChartOption
		result      string
	}{
		{
			name: "defaults",
			makeOptions: func() DoughnutChartOption {
				opt := makeBasicDoughnutChartOption()
				opt.Title.Show = Ptr(false)
				opt.Legend.Offset = OffsetStr{}
				opt.Legend.Symbol = ""
				opt.Legend.Vertical = nil
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 20 20\nL 580 20\nL 580 380\nL 20 380\nL 20 20\" style=\"stroke:none;fill:white\"/><path  d=\"M 40 43\nL 70 43\nL 70 56\nL 40 56\nL 40 43\" style=\"stroke:none;fill:rgb(84,112,198)\"/><text x=\"72\" y=\"55\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Series-A</text><path  d=\"M 151 43\nL 181 43\nL 181 56\nL 151 56\nL 151 43\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"183\" y=\"55\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Series-B</text><path  d=\"M 261 43\nL 291 43\nL 291 56\nL 261 56\nL 261 43\" style=\"stroke:none;fill:rgb(250,200,88)\"/><text x=\"293\" y=\"55\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Series-C</text><path  d=\"M 371 43\nL 401 43\nL 401 56\nL 371 56\nL 371 43\" style=\"stroke:none;fill:rgb(238,102,102)\"/><text x=\"403\" y=\"55\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Series-D</text><path  d=\"M 40 59\nL 70 59\nL 70 72\nL 40 72\nL 40 59\" style=\"stroke:none;fill:rgb(115,192,222)\"/><text x=\"72\" y=\"71\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Series-E</text><path  d=\"M 300 226\nL 300 119\nA 107 107 119.89 0 1 393 279\nL 300 226\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 392 173\nL 405 165\nM 405 165\nL 420 165\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:none\"/><text x=\"423\" y=\"170\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-A: 33.3%</text><path  d=\"M 300 226\nL 393 279\nA 107 107 84.08 0 1 256 324\nL 300 226\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 333 327\nL 337 342\nM 337 342\nL 352 342\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:none\"/><text x=\"355\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-B: 23.35%</text><path  d=\"M 300 226\nL 256 324\nA 107 107 66.35 0 1 193 225\nL 300 226\nZ\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 210 284\nL 198 292\nM 198 292\nL 183 292\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:none\"/><text x=\"83\" y=\"297\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-C: 18.43%</text><path  d=\"M 300 226\nL 193 225\nA 107 107 55.37 0 1 240 137\nL 300 226\nZ\" style=\"stroke:none;fill:rgb(238,102,102)\"/><path  d=\"M 206 176\nL 193 169\nM 193 169\nL 178 169\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:none\"/><text x=\"78\" y=\"174\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-D: 15.37%</text><path  d=\"M 300 226\nL 240 137\nA 107 107 34.32 0 1 300 119\nL 300 226\nZ\" style=\"stroke:none;fill:rgb(115,192,222)\"/><path  d=\"M 269 124\nL 264 110\nM 264 110\nL 249 110\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:none\"/><text x=\"157\" y=\"115\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-E: 9.53%</text><circle cx=\"300\" cy=\"226\" r=\"64\" style=\"stroke:none;fill:white\"/></svg>",
		},
		{
			name:        "basic_themed",
			themed:      true,
			makeOptions: makeBasicDoughnutChartOption,
			result:      "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:rgb(40,40,40)\"/><path  d=\"M 20 29\nL 50 29\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:none\"/><circle cx=\"35\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><text x=\"52\" y=\"35\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Series-A</text><path  d=\"M 20 49\nL 50 49\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:none\"/><circle cx=\"35\" cy=\"49\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><text x=\"52\" y=\"55\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Series-B</text><path  d=\"M 20 69\nL 50 69\" style=\"stroke-width:3;stroke:rgb(100,180,210);fill:none\"/><circle cx=\"35\" cy=\"69\" r=\"5\" style=\"stroke-width:3;stroke:rgb(100,180,210);fill:rgb(100,180,210)\"/><text x=\"52\" y=\"75\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Series-C</text><path  d=\"M 20 89\nL 50 89\" style=\"stroke-width:3;stroke:rgb(64,160,110);fill:none\"/><circle cx=\"35\" cy=\"89\" r=\"5\" style=\"stroke-width:3;stroke:rgb(64,160,110);fill:rgb(64,160,110)\"/><text x=\"52\" y=\"95\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Series-D</text><path  d=\"M 20 109\nL 50 109\" style=\"stroke-width:3;stroke:rgb(154,100,180);fill:none\"/><circle cx=\"35\" cy=\"109\" r=\"5\" style=\"stroke-width:3;stroke:rgb(154,100,180);fill:rgb(154,100,180)\"/><text x=\"52\" y=\"115\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Series-E</text><text x=\"285\" y=\"36\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Title</text><text x=\"287\" y=\"52\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sub</text><path  d=\"M 300 223\nL 300 98\nA 125 125 119.89 0 1 409 285\nL 300 223\nZ\" style=\"stroke:none;fill:rgb(255,100,100)\"/><path  d=\"M 408 161\nL 421 153\nM 421 153\nL 436 153\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:none\"/><text x=\"439\" y=\"158\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-A: 33.3%</text><path  d=\"M 300 223\nL 409 285\nA 125 125 84.08 0 1 249 337\nL 300 223\nZ\" style=\"stroke:none;fill:rgb(255,210,100)\"/><path  d=\"M 338 342\nL 343 356\nM 343 356\nL 358 356\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:none\"/><text x=\"361\" y=\"361\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-B: 23.35%</text><path  d=\"M 300 223\nL 249 337\nA 125 125 66.35 0 1 175 222\nL 300 223\nZ\" style=\"stroke:none;fill:rgb(100,180,210)\"/><path  d=\"M 195 290\nL 183 299\nM 183 299\nL 168 299\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:none\"/><text x=\"68\" y=\"304\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-C: 18.43%</text><path  d=\"M 300 223\nL 175 222\nA 125 125 55.37 0 1 229 120\nL 300 223\nZ\" style=\"stroke:none;fill:rgb(64,160,110)\"/><path  d=\"M 190 165\nL 177 158\nM 177 158\nL 162 158\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:none\"/><text x=\"62\" y=\"163\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-D: 15.37%</text><path  d=\"M 300 223\nL 229 120\nA 125 125 34.32 0 1 300 98\nL 300 223\nZ\" style=\"stroke:none;fill:rgb(154,100,180)\"/><path  d=\"M 264 104\nL 259 90\nM 259 90\nL 244 90\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:none\"/><text x=\"152\" y=\"95\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-E: 9.53%</text><circle cx=\"300\" cy=\"223\" r=\"75\" style=\"stroke:none;fill:rgb(40,40,40)\"/></svg>",
		},
		{
			name: "custom_fonts",
			makeOptions: func() DoughnutChartOption {
				opt := makeBasicDoughnutChartOption()
				customFont := NewFontStyleWithSize(4.0).WithColor(ColorBlue)
				opt.SeriesList[0].Label.FontStyle = customFont
				opt.Legend.FontStyle = customFont
				opt.Title.FontStyle = customFont
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 20 20\nL 580 20\nL 580 380\nL 20 380\nL 20 20\" style=\"stroke:none;fill:white\"/><path  d=\"M 40 49\nL 70 49\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"55\" cy=\"49\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"72\" y=\"55\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Series-A</text><path  d=\"M 40 69\nL 70 69\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"55\" cy=\"69\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"72\" y=\"75\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Series-B</text><path  d=\"M 40 89\nL 70 89\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"55\" cy=\"89\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"72\" y=\"95\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Series-C</text><path  d=\"M 40 109\nL 70 109\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"55\" cy=\"109\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"72\" y=\"115\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Series-D</text><path  d=\"M 40 129\nL 70 129\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:none\"/><circle cx=\"55\" cy=\"129\" r=\"5\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:rgb(115,192,222)\"/><text x=\"72\" y=\"135\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Series-E</text><text x=\"295\" y=\"46\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Title</text><text x=\"296\" y=\"52\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Sub</text><path  d=\"M 300 213\nL 300 96\nA 117 117 119.89 0 1 402 271\nL 300 213\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 401 155\nL 414 147\nM 414 147\nL 429 147\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:none\"/><text x=\"432\" y=\"149\" style=\"stroke:none;fill:blue;font-size:5.1px;font-family:'Roboto Medium',sans-serif\">Series-A: 33.3%</text><path  d=\"M 300 213\nL 402 271\nA 117 117 84.08 0 1 252 320\nL 300 213\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 336 324\nL 341 338\nM 341 338\nL 356 338\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:none\"/><text x=\"359\" y=\"343\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-B: 23.35%</text><path  d=\"M 300 213\nL 252 320\nA 117 117 66.35 0 1 183 212\nL 300 213\nZ\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 202 276\nL 189 284\nM 189 284\nL 174 284\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:none\"/><text x=\"74\" y=\"289\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-C: 18.43%</text><path  d=\"M 300 213\nL 183 212\nA 117 117 55.37 0 1 234 116\nL 300 213\nZ\" style=\"stroke:none;fill:rgb(238,102,102)\"/><path  d=\"M 197 158\nL 184 151\nM 184 151\nL 169 151\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:none\"/><text x=\"69\" y=\"156\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-D: 15.37%</text><path  d=\"M 300 213\nL 234 116\nA 117 117 34.32 0 1 300 96\nL 300 213\nZ\" style=\"stroke:none;fill:rgb(115,192,222)\"/><path  d=\"M 266 102\nL 261 87\nM 261 87\nL 246 87\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:none\"/><text x=\"154\" y=\"92\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-E: 9.53%</text><circle cx=\"300\" cy=\"213\" r=\"70\" style=\"stroke:none;fill:white\"/></svg>",
		},
		{
			name: "variable_series_radius",
			makeOptions: func() DoughnutChartOption {
				opt := makeMinimalDoughnutChartOption()
				for i := range opt.SeriesList {
					opt.SeriesList[i].Radius = strconv.Itoa((i+1)*10) + "%"
				}
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 20 20\nL 580 20\nL 580 380\nL 20 380\nL 20 20\" style=\"stroke:none;fill:white\"/><path  d=\"M 300 200\nL 300 168\nA 32 32 119.89 0 1 328 216\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 300 200\nL 355 232\nA 64 64 84.08 0 1 274 258\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 300 200\nL 261 288\nA 96 96 66.35 0 1 204 199\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 300 200\nL 172 199\nA 128 128 55.37 0 1 228 94\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(238,102,102)\"/><path  d=\"M 300 200\nL 210 68\nA 160 160 34.32 0 1 300 40\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(115,192,222)\"/><circle cx=\"300\" cy=\"200\" r=\"19\" style=\"stroke:none;fill:white\"/></svg>",
		},
		{
			name: "center_radius_small",
			makeOptions: func() DoughnutChartOption {
				opt := makeMinimalDoughnutChartOption()
				opt.RadiusCenter = "20"
				for i := range opt.SeriesList {
					opt.SeriesList[i].Label.Show = Ptr(true)
				}
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 20 20\nL 580 20\nL 580 380\nL 20 380\nL 20 20\" style=\"stroke:none;fill:white\"/><path  d=\"M 300 200\nL 300 72\nA 128 128 119.89 0 1 411 264\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 410 136\nL 423 129\nM 423 129\nL 438 129\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:none\"/><text x=\"441\" y=\"134\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-A: 33.3%</text><path  d=\"M 300 200\nL 411 264\nA 128 128 84.08 0 1 248 317\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 339 321\nL 344 335\nM 344 335\nL 359 335\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:none\"/><text x=\"362\" y=\"340\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-B: 23.35%</text><path  d=\"M 300 200\nL 248 317\nA 128 128 66.35 0 1 172 199\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 193 269\nL 180 277\nM 180 277\nL 165 277\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:none\"/><text x=\"65\" y=\"282\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-C: 18.43%</text><path  d=\"M 300 200\nL 172 199\nA 128 128 55.37 0 1 228 94\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(238,102,102)\"/><path  d=\"M 187 140\nL 174 133\nM 174 133\nL 159 133\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:none\"/><text x=\"59\" y=\"138\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-D: 15.37%</text><path  d=\"M 300 200\nL 228 94\nA 128 128 34.32 0 1 300 72\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(115,192,222)\"/><path  d=\"M 263 78\nL 258 64\nM 258 64\nL 243 64\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:none\"/><text x=\"151\" y=\"69\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-E: 9.53%</text><circle cx=\"300\" cy=\"200\" r=\"20\" style=\"stroke:none;fill:white\"/></svg>",
		},
		{
			name: "center_radius_large",
			makeOptions: func() DoughnutChartOption {
				opt := makeMinimalDoughnutChartOption()
				opt.RadiusRing = "42%"
				opt.RadiusCenter = "40%"
				for i := range opt.SeriesList {
					opt.SeriesList[i].Label.Show = Ptr(true)
				}
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 20 20\nL 580 20\nL 580 380\nL 20 380\nL 20 20\" style=\"stroke:none;fill:white\"/><path  d=\"M 300 200\nL 300 66\nA 134 134 119.89 0 1 417 267\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 416 133\nL 429 126\nM 429 126\nL 444 126\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:none\"/><text x=\"447\" y=\"131\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-A: 33.3%</text><path  d=\"M 300 200\nL 417 267\nA 134 134 84.08 0 1 245 323\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 341 327\nL 346 342\nM 346 342\nL 361 342\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:none\"/><text x=\"364\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-B: 23.35%</text><path  d=\"M 300 200\nL 245 323\nA 134 134 66.35 0 1 166 199\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 188 272\nL 175 281\nM 175 281\nL 160 281\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:none\"/><text x=\"60\" y=\"286\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-C: 18.43%</text><path  d=\"M 300 200\nL 166 199\nA 134 134 55.37 0 1 224 89\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(238,102,102)\"/><path  d=\"M 182 137\nL 169 130\nM 169 130\nL 154 130\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:none\"/><text x=\"54\" y=\"135\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-D: 15.37%</text><path  d=\"M 300 200\nL 224 89\nA 134 134 34.32 0 1 300 66\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(115,192,222)\"/><path  d=\"M 261 72\nL 256 58\nM 256 58\nL 241 58\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:none\"/><text x=\"149\" y=\"63\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-E: 9.53%</text><circle cx=\"300\" cy=\"200\" r=\"124\" style=\"stroke:none;fill:white\"/></svg>",
		},
		{
			name: "segment_gap",
			makeOptions: func() DoughnutChartOption {
				opt := makeMinimalDoughnutChartOption()
				opt.SegmentGap = 20.0
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 20 20\nL 580 20\nL 580 380\nL 20 380\nL 20 20\" style=\"stroke:none;fill:white\"/><path  d=\"M 300 200\nL 300 72\nA 128 128 119.89 0 1 411 264\nL 300 200\nZ\" style=\"stroke-width:20;stroke:white;fill:rgb(84,112,198)\"/><path  d=\"M 300 200\nL 411 264\nA 128 128 84.08 0 1 248 317\nL 300 200\nZ\" style=\"stroke-width:20;stroke:white;fill:rgb(145,204,117)\"/><path  d=\"M 300 200\nL 248 317\nA 128 128 66.35 0 1 172 199\nL 300 200\nZ\" style=\"stroke-width:20;stroke:white;fill:rgb(250,200,88)\"/><path  d=\"M 300 200\nL 172 199\nA 128 128 55.37 0 1 228 94\nL 300 200\nZ\" style=\"stroke-width:20;stroke:white;fill:rgb(238,102,102)\"/><path  d=\"M 300 200\nL 228 94\nA 128 128 34.32 0 1 300 72\nL 300 200\nZ\" style=\"stroke-width:20;stroke:white;fill:rgb(115,192,222)\"/><circle cx=\"300\" cy=\"200\" r=\"77\" style=\"stroke:none;fill:white\"/></svg>",
		},
		{
			name: "center_sum",
			makeOptions: func() DoughnutChartOption {
				opt := makeMinimalDoughnutChartOption()
				opt.CenterValues = "sum"
				opt.CenterValuesFontStyle.FontSize = 24.0
				opt.CenterValuesFontStyle.FontColor = ColorNavy
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 20 20\nL 580 20\nL 580 380\nL 20 380\nL 20 20\" style=\"stroke:none;fill:white\"/><path  d=\"M 300 200\nL 300 72\nA 128 128 119.89 0 1 411 264\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 300 200\nL 411 264\nA 128 128 84.08 0 1 248 317\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 300 200\nL 248 317\nA 128 128 66.35 0 1 172 199\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 300 200\nL 172 199\nA 128 128 55.37 0 1 228 94\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(238,102,102)\"/><path  d=\"M 300 200\nL 228 94\nA 128 128 34.32 0 1 300 72\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(115,192,222)\"/><circle cx=\"300\" cy=\"200\" r=\"77\" style=\"stroke:none;fill:white\"/><text x=\"262\" y=\"215\" style=\"stroke:none;fill:navy;font-size:30.7px;font-family:'Roboto Medium',sans-serif\">3.15k</text></svg>",
		},
		{
			name: "center_labels",
			makeOptions: func() DoughnutChartOption {
				opt := makeMinimalDoughnutChartOption()
				opt.CenterValues = "labels"
				for i := range opt.SeriesList {
					opt.SeriesList[i].Label.Show = Ptr(true)
				}
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 20 20\nL 580 20\nL 580 380\nL 20 380\nL 20 20\" style=\"stroke:none;fill:white\"/><path  d=\"M 300 200\nL 300 40\nA 160 160 119.89 0 1 439 280\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 300 200\nL 439 280\nA 160 160 84.08 0 1 235 346\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 300 200\nL 235 346\nA 160 160 66.35 0 1 140 199\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 300 200\nL 140 199\nA 160 160 55.37 0 1 210 68\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(238,102,102)\"/><path  d=\"M 300 200\nL 210 68\nA 160 160 34.32 0 1 300 40\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(115,192,222)\"/><circle cx=\"300\" cy=\"200\" r=\"96\" style=\"stroke:none;fill:white\"/><path  d=\"M 378 146\nL 366 154\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><text x=\"276\" y=\"161\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-A: 33.3%</text><path  d=\"M 356 277\nL 350 269\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/><text x=\"253\" y=\"269\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-B: 23.35%</text><path  d=\"M 215 244\nL 235 234\" style=\"stroke-width:2;stroke:rgb(250,200,88);fill:none\"/><text x=\"235\" y=\"241\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-C: 18.43%</text><path  d=\"M 213 162\nL 227 168\" style=\"stroke-width:2;stroke:rgb(238,102,102);fill:none\"/><text x=\"227\" y=\"175\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-D: 15.37%</text><path  d=\"M 253 117\nL 255 120\" style=\"stroke-width:2;stroke:rgb(115,192,222);fill:none\"/><text x=\"255\" y=\"133\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Series-E: 9.53%</text></svg>",
		},
		{
			name: "center_lots_labels",
			makeOptions: func() DoughnutChartOption {
				values := []float64{
					9104772, 11754004, 10827529, 10394055, 9597085,
					17947406, 36753736, 10467366, 19051562, 10521556,
				}

				return DoughnutChartOption{
					SeriesList:   NewSeriesListDoughnut(values),
					CenterValues: "labels",
					Legend: LegendOption{
						SeriesNames: []string{
							"Cyprus", "Denmark", "Estonia", "Finland", "France",
							"Germany", "Greece", "Hungary", "Ireland", "Italy",
						},
						Show: Ptr(false),
					},
				}
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 20 20\nL 580 20\nL 580 380\nL 20 380\nL 20 20\" style=\"stroke:none;fill:white\"/><path  d=\"M 300 200\nL 440 87\nA 180 180 26.62 0 1 476 162\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path  d=\"M 300 200\nL 369 34\nA 180 180 28.90 0 1 440 87\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 300 200\nL 300 20\nA 180 180 22.39 0 1 369 34\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 300 200\nL 476 162\nA 180 180 25.56 0 1 475 242\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(238,102,102)\"/><path  d=\"M 300 200\nL 475 242\nA 180 180 23.60 0 1 444 308\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(115,192,222)\"/><path  d=\"M 300 200\nL 444 308\nA 180 180 44.13 0 1 328 378\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(59,162,114)\"/><path  d=\"M 300 200\nL 328 378\nA 180 180 90.37 0 1 122 226\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(252,132,82)\"/><path  d=\"M 300 200\nL 122 226\nA 180 180 25.74 0 1 128 147\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(154,96,180)\"/><path  d=\"M 300 200\nL 128 147\nA 180 180 46.84 0 1 221 38\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(234,124,204)\"/><path  d=\"M 300 200\nL 221 38\nA 180 180 25.87 0 1 300 20\nL 300 200\nZ\" style=\"stroke:none;fill:rgb(123,142,198)\"/><circle cx=\"300\" cy=\"200\" r=\"108\" style=\"stroke:none;fill:white\"/><path  d=\"M 341 101\nL 342 106\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><text x=\"260\" y=\"119\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Cyprus: 6.21%</text><path  d=\"M 359 111\nL 354 119\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/><text x=\"261\" y=\"132\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Denmark: 8.02%</text><path  d=\"M 393 147\nL 378 155\" style=\"stroke-width:2;stroke:rgb(250,200,88);fill:none\"/><text x=\"294\" y=\"162\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Estonia: 7.39%</text><path  d=\"M 407 193\nL 386 194\" style=\"stroke-width:2;stroke:rgb(238,102,102);fill:none\"/><text x=\"302\" y=\"201\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Finland: 7.09%</text><path  d=\"M 401 237\nL 378 229\" style=\"stroke-width:2;stroke:rgb(115,192,222);fill:none\"/><text x=\"297\" y=\"236\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">France: 6.55%</text><path  d=\"M 362 288\nL 352 274\" style=\"stroke-width:2;stroke:rgb(59,162,114);fill:none\"/><text x=\"251\" y=\"274\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Germany: 12.25%</text><path  d=\"M 224 276\nL 239 261\" style=\"stroke-width:2;stroke:rgb(252,132,82);fill:none\"/><text x=\"239\" y=\"261\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Greece: 25.1%</text><path  d=\"M 195 177\nL 214 181\" style=\"stroke-width:2;stroke:rgb(154,96,180);fill:none\"/><text x=\"214\" y=\"188\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Hungary: 7.14%</text><path  d=\"M 224 125\nL 232 133\" style=\"stroke-width:2;stroke:rgb(234,124,204);fill:none\"/><text x=\"232\" y=\"146\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Ireland: 13.01%</text><path  d=\"M 266 98\nL 266 98\" style=\"stroke-width:2;stroke:rgb(123,142,198);fill:none\"/><text x=\"266\" y=\"111\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Italy: 7.18%</text></svg>",
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

				validateDoughnutChartRender(t, p, tt.makeOptions(), tt.result)
			})
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-options", func(t *testing.T) {
				p := NewPainter(painterOptions)
				opt := tt.makeOptions()
				opt.Theme = GetTheme(ThemeVividDark)

				validateDoughnutChartRender(t, p, opt, tt.result)
			})
		} else {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				p := NewPainter(painterOptions)

				validateDoughnutChartRender(t, p.Child(PainterPaddingOption(NewBoxEqual(20))),
					tt.makeOptions(), tt.result)
			})
		}
	}
}

func validateDoughnutChartRender(t *testing.T, p *Painter, opt DoughnutChartOption, expectedResult string) {
	t.Helper()

	err := p.DoughnutChart(opt)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, expectedResult, data)
}

func TestDoughnutChartError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		makeOptions      func() DoughnutChartOption
		errorMsgContains string
	}{
		{
			name: "empty_series",
			makeOptions: func() DoughnutChartOption {
				return NewDoughnutChartOptionWithData([]float64{})
			},
			errorMsgContains: "empty series list",
		},
		{
			name: "zero_sum",
			makeOptions: func() DoughnutChartOption {
				return NewDoughnutChartOptionWithData([]float64{0.0, 0.0})
			},
			errorMsgContains: "greater than 0",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			})

			err := p.DoughnutChart(tt.makeOptions())
			require.Error(t, err)
			require.ErrorContains(t, err, tt.errorMsgContains)
		})
	}
}

func TestClampAngleToSector(t *testing.T) {
	t.Parallel()

	cases := []struct {
		angle, start, end float64
		want              float64
	}{
		{0.5, 0, math.Pi / 2, 0.5},
		{-0.1, 0, math.Pi / 2, math.Pi / 2},
		{math.Pi, 0, math.Pi / 2, math.Pi / 2},
		{math.Pi, 3 * math.Pi / 2, math.Pi / 2, 3 * math.Pi / 2},
		{0, 3 * math.Pi / 2, math.Pi / 2, 0},
	}
	for _, tt := range cases {
		got := clampAngleToSector(tt.angle, tt.start, tt.end)
		assert.InDelta(t, tt.want, got, 1e-6)
	}
}

func TestConnectionPoint(t *testing.T) {
	t.Parallel()

	s := sector{startAngle: 0, delta: math.Pi / 2}
	x, y := s.connectionPoint(0, 0, 5, 0, 10)
	assert.Equal(t, 0, x)
	assert.Equal(t, 5, y)

	x, y = s.connectionPoint(0, 0, 5, 5, 5)
	assert.Equal(t, 3, x)
	assert.Equal(t, 3, y)

	x, y = s.connectionPoint(0, 0, 5, -10, 0)
	assert.Equal(t, 0, x)
	assert.Equal(t, 5, y)
}

func TestIsInsideCircle(t *testing.T) {
	t.Parallel()

	c := isInsideCircle(NewBox(-1, -1, 1, 1), 0, 0, 5)
	assert.True(t, c)
	c = isInsideCircle(NewBox(4, 4, 6, 6), 0, 0, 5)
	assert.False(t, c)
}

func TestClampInsideCircle(t *testing.T) {
	t.Parallel()

	lp := &labelPlacement{box: NewBox(6, -1, 8, 1)}
	clampInsideCircle(lp, 0, 0, 5)
	assert.True(t, isInsideCircle(lp.box, 0, 0, 5))
}

func TestMinimalRadialPush(t *testing.T) {
	t.Parallel()

	p := &labelPlacement{box: NewBox(1, 0, 3, 2)}
	q := &labelPlacement{box: NewBox(0, 0, 2, 2)}
	dx, dy := minimalRadialPush(p, q, 0, 0)
	assert.Equal(t, 3, dx)
	assert.Equal(t, 2, dy)
}

func TestProjectBoxRadially(t *testing.T) {
	t.Parallel()

	b := NewBox(0, 0, 2, 2)
	min, max := projectBoxRadially(b, 1, 0)
	assert.InDelta(t, 0.0, min, 1e-6)
	assert.InDelta(t, 2.0, max, 1e-6)

	min, max = projectBoxRadially(b, 0, 1)
	assert.InDelta(t, 0.0, min, 1e-6)
	assert.InDelta(t, 2.0, max, 1e-6)
}

func TestShiftLabelHorizontallyTowardSector(t *testing.T) {
	t.Parallel()

	lp := &labelPlacement{box: NewBox(-2, -1, 0, 1)}
	shiftLabelHorizontallyTowardSector(lp, 0, 0)
	assert.Equal(t, NewBox(0, -1, 2, 1), lp.box)

	lp = &labelPlacement{box: NewBox(2, -1, 4, 1)}
	shiftLabelHorizontallyTowardSector(lp, 0, math.Pi)
	assert.Equal(t, NewBox(-2, -1, 0, 1), lp.box)
}

func TestAnyLabelCollision(t *testing.T) {
	t.Parallel()

	placed := []*labelPlacement{{box: NewBox(1, 1, 3, 3)}}
	c := anyLabelCollision(NewBox(0, 0, 2, 2), placed)
	assert.True(t, c)
	c = anyLabelCollision(NewBox(3, 3, 5, 5), placed)
	assert.False(t, c)
}

func TestComputeLabelBox(t *testing.T) {
	t.Parallel()

	b := computeLabelBox(0, 0, 5, 0, NewBox(0, 0, 2, 2))
	assert.Equal(t, NewBox(3, -2, 5, 0), b)

	b = computeLabelBox(0, 0, 5, math.Pi, NewBox(0, 0, 2, 2))
	assert.Equal(t, NewBox(-5, -2, -3, 0), b)
}
