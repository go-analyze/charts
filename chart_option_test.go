package charts

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChartOption(t *testing.T) {
	t.Parallel()

	fns := []OptionFunc{
		SVGOutputOptionFunc(),
		FontOptionFunc(GetDefaultFont()),
		ThemeNameOptionFunc(ThemeVividDark),
		TitleTextOptionFunc("title"),
		LegendLabelsOptionFunc([]string{"label"}),
		XAxisLabelsOptionFunc([]string{"xaxis"}),
		YAxisLabelsOptionFunc([]string{"yaxis"}),
		DimensionsOptionFunc(800, 600),
		PaddingOptionFunc(NewBoxEqual(10)),
	}
	opt := ChartOption{}
	for _, fn := range fns {
		fn(&opt)
	}
	require.Equal(t, ChartOption{
		OutputFormat: ChartOutputSVG,
		Font:         GetDefaultFont(),
		Theme:        GetTheme(ThemeVividDark),
		Title: TitleOption{
			Text: "title",
		},
		Legend: LegendOption{
			SeriesNames: []string{"label"},
		},
		XAxis: XAxisOption{
			Labels: []string{"xaxis"},
		},
		YAxis: []YAxisOption{
			{
				Labels: []string{"yaxis"},
			},
		},
		Width:   800,
		Height:  600,
		Padding: NewBoxEqual(10),
	}, opt)
}

func TestChartOptionSeriesShowLabel(t *testing.T) {
	t.Parallel()

	opt := ChartOption{
		SeriesList: NewSeriesListPie([]float64{1, 2}).ToGenericSeriesList(),
	}
	SeriesShowLabel(true)(&opt)
	assert.True(t, flagIs(true, opt.SeriesList[0].Label.Show))

	SeriesShowLabel(false)(&opt)
	assert.True(t, flagIs(false, opt.SeriesList[0].Label.Show))
}

func newNoTypeSeriesListFromValues(values [][]float64) GenericSeriesList {
	return NewSeriesListGeneric(values, "")
}

func TestChartOptionMarkLine(t *testing.T) {
	t.Parallel()

	opt := ChartOption{
		SeriesList: newNoTypeSeriesListFromValues([][]float64{{1, 2}}),
	}
	MarkLineOptionFunc(0, "min", "max")(&opt)
	assert.Equal(t, NewMarkLine("min", "max"), opt.SeriesList[0].MarkLine)
}

func TestChartOptionMarkPoint(t *testing.T) {
	t.Parallel()

	opt := ChartOption{
		SeriesList: newNoTypeSeriesListFromValues([][]float64{{1, 2}}),
	}
	MarkPointOptionFunc(0, "min", "max")(&opt)
	assert.Equal(t, NewMarkPoint("min", "max"), opt.SeriesList[0].MarkPoint)
}

func TestLineRender(t *testing.T) {
	t.Parallel()

	values := [][]float64{
		{120, 132, 101, 134, 90, 230, 210},
		{220, 182, 191, 234, 290, 330, 310},
		{150, 232, 201, 154, 190, 330, 410},
		{320, 332, 301, 334, 390, 330, 320},
		{820, 932, 901, 934, 1290, 1330, 1320},
	}
	p, err := LineRender(
		values,
		SVGOutputOptionFunc(),
		TitleTextOptionFunc("Line"),
		XAxisLabelsOptionFunc([]string{
			"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun",
		}),
		LegendLabelsOptionFunc([]string{
			"Email", "Union Ads", "Video Ads", "Direct", "Search Engine",
		}),
		func(opt *ChartOption) {
			opt.ValueFormatter = func(f float64) string {
				return fmt.Sprintf("%.0f", f)
			}
		},
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"20\" y=\"36\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line</text><path d=\"M 21 45\nL 51 45\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"36\" cy=\"45\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"53\" y=\"51\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path d=\"M 112 45\nL 142 45\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"127\" cy=\"45\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"144\" y=\"51\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path d=\"M 235 45\nL 265 45\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"250\" cy=\"45\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"267\" y=\"51\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><path d=\"M 357 45\nL 387 45\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"372\" cy=\"45\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"389\" y=\"51\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path d=\"M 450 45\nL 480 45\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:none\"/><circle cx=\"465\" cy=\"45\" r=\"5\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:rgb(115,192,222)\"/><text x=\"482\" y=\"51\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><text x=\"19\" y=\"73\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1440</text><text x=\"19\" y=\"104\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1280</text><text x=\"19\" y=\"136\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1120</text><text x=\"27\" y=\"168\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"27\" y=\"199\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"27\" y=\"231\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"27\" y=\"263\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"27\" y=\"294\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"27\" y=\"326\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"45\" y=\"358\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path d=\"M 60 67\nL 580 67\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 60 98\nL 580 98\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 60 130\nL 580 130\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 60 162\nL 580 162\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 60 194\nL 580 194\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 60 226\nL 580 226\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 60 258\nL 580 258\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 60 290\nL 580 290\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 60 322\nL 580 322\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 64 354\nL 580 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 64 359\nL 64 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 137 359\nL 137 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 211 359\nL 211 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 285 359\nL 285 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 358 359\nL 358 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 432 359\nL 432 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 506 359\nL 506 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 580 359\nL 580 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"85\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"161\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"233\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"308\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"386\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"458\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"530\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path d=\"M 100 331\nL 174 328\nL 248 334\nL 321 328\nL 395 337\nL 469 309\nL 543 313\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"100\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"174\" cy=\"328\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"248\" cy=\"334\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"321\" cy=\"328\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"395\" cy=\"337\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"469\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"543\" cy=\"313\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><path d=\"M 100 311\nL 174 318\nL 248 316\nL 321 308\nL 395 297\nL 469 289\nL 543 293\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"100\" cy=\"311\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"174\" cy=\"318\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"248\" cy=\"316\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"321\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"395\" cy=\"297\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"469\" cy=\"289\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"543\" cy=\"293\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><path d=\"M 100 325\nL 174 308\nL 248 314\nL 321 324\nL 395 317\nL 469 289\nL 543 273\" style=\"stroke-width:2;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"100\" cy=\"325\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"174\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"248\" cy=\"314\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"321\" cy=\"324\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"395\" cy=\"317\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"469\" cy=\"289\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"543\" cy=\"273\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><path d=\"M 100 291\nL 174 288\nL 248 295\nL 321 288\nL 395 277\nL 469 289\nL 543 291\" style=\"stroke-width:2;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"100\" cy=\"291\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:white\"/><circle cx=\"174\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:white\"/><circle cx=\"248\" cy=\"295\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:white\"/><circle cx=\"321\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:white\"/><circle cx=\"395\" cy=\"277\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:white\"/><circle cx=\"469\" cy=\"289\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:white\"/><circle cx=\"543\" cy=\"291\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:white\"/><path d=\"M 100 191\nL 174 169\nL 248 175\nL 321 168\nL 395 97\nL 469 89\nL 543 91\" style=\"stroke-width:2;stroke:rgb(115,192,222);fill:none\"/><circle cx=\"100\" cy=\"191\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:white\"/><circle cx=\"174\" cy=\"169\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:white\"/><circle cx=\"248\" cy=\"175\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:white\"/><circle cx=\"321\" cy=\"168\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:white\"/><circle cx=\"395\" cy=\"97\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:white\"/><circle cx=\"469\" cy=\"89\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:white\"/><circle cx=\"543\" cy=\"91\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:white\"/></svg>", data)
}

func TestScatterRender(t *testing.T) {
	t.Parallel()

	values := [][]float64{
		{120, 132, 101, 134, 90, 230, 210},
		{220, 182, 191, 234, 290, 330, 310},
		{150, 232, 201, 154, 190, 330, 410},
		{320, 332, 301, 334, 390, 330, 320},
		{820, 932, 901, 934, 1290, 1330, 1320},
	}
	p, err := ScatterRender(
		values,
		SVGOutputOptionFunc(),
		TitleTextOptionFunc("Scatter"),
		XAxisLabelsOptionFunc([]string{
			"1", "2", "3", "4", "5", "6", "7",
		}),
		LegendLabelsOptionFunc([]string{
			"A", "B", "C", "D", "E",
		}),
		func(opt *ChartOption) {
			opt.Symbol = SymbolSquare
			opt.Legend.Symbol = SymbolSquare
		},
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"20\" y=\"36\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Scatter</text><path d=\"M 156 23\nL 186 23\nL 186 36\nL 156 36\nL 156 23\" style=\"stroke:none;fill:rgb(84,112,198)\"/><text x=\"188\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">A</text><path d=\"M 219 23\nL 249 23\nL 249 36\nL 219 36\nL 219 23\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"251\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">B</text><path d=\"M 281 23\nL 311 23\nL 311 36\nL 281 36\nL 281 23\" style=\"stroke:none;fill:rgb(250,200,88)\"/><text x=\"313\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">C</text><path d=\"M 343 23\nL 373 23\nL 373 36\nL 343 36\nL 343 23\" style=\"stroke:none;fill:rgb(238,102,102)\"/><text x=\"375\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">D</text><path d=\"M 405 23\nL 435 23\nL 435 36\nL 405 36\nL 405 23\" style=\"stroke:none;fill:rgb(115,192,222)\"/><text x=\"437\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"19\" y=\"62\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"19\" y=\"94\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.28k</text><text x=\"19\" y=\"127\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.12k</text><text x=\"31\" y=\"160\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"31\" y=\"193\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"31\" y=\"226\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"31\" y=\"259\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"31\" y=\"292\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"31\" y=\"325\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"49\" y=\"358\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path d=\"M 64 56\nL 580 56\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 64 89\nL 580 89\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 64 122\nL 580 122\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 64 155\nL 580 155\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 64 188\nL 580 188\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 64 221\nL 580 221\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 64 254\nL 580 254\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 64 287\nL 580 287\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 64 320\nL 580 320\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 68 354\nL 580 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 68 359\nL 68 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 141 359\nL 141 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 214 359\nL 214 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 287 359\nL 287 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 360 359\nL 360 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 433 359\nL 433 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 506 359\nL 506 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 580 359\nL 580 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"100\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1</text><text x=\"173\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2</text><text x=\"246\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">3</text><text x=\"319\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">4</text><text x=\"392\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">5</text><text x=\"465\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">6</text><text x=\"539\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">7</text><path d=\"M 66 328\nL 70 328\nL 70 332\nL 66 332\nL 66 328\nM 151 325\nL 155 325\nL 155 329\nL 151 329\nL 151 325\nM 236 332\nL 240 332\nL 240 336\nL 236 336\nL 236 332\nM 322 325\nL 326 325\nL 326 329\nL 322 329\nL 322 325\nM 407 334\nL 411 334\nL 411 338\nL 407 338\nL 407 334\nM 492 305\nL 496 305\nL 496 309\nL 492 309\nL 492 305\nM 578 309\nL 582 309\nL 582 313\nL 578 313\nL 578 309\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path d=\"M 66 307\nL 70 307\nL 70 311\nL 66 311\nL 66 307\nM 151 315\nL 155 315\nL 155 319\nL 151 319\nL 151 315\nM 236 313\nL 240 313\nL 240 317\nL 236 317\nL 236 313\nM 322 304\nL 326 304\nL 326 308\nL 322 308\nL 322 304\nM 407 292\nL 411 292\nL 411 296\nL 407 296\nL 407 292\nM 492 284\nL 496 284\nL 496 288\nL 492 288\nL 492 284\nM 578 288\nL 582 288\nL 582 292\nL 578 292\nL 578 288\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path d=\"M 66 321\nL 70 321\nL 70 325\nL 66 325\nL 66 321\nM 151 304\nL 155 304\nL 155 308\nL 151 308\nL 151 304\nM 236 311\nL 240 311\nL 240 315\nL 236 315\nL 236 311\nM 322 321\nL 326 321\nL 326 325\nL 322 325\nL 322 321\nM 407 313\nL 411 313\nL 411 317\nL 407 317\nL 407 313\nM 492 284\nL 496 284\nL 496 288\nL 492 288\nL 492 284\nM 578 268\nL 582 268\nL 582 272\nL 578 272\nL 578 268\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><path d=\"M 66 286\nL 70 286\nL 70 290\nL 66 290\nL 66 286\nM 151 284\nL 155 284\nL 155 288\nL 151 288\nL 151 284\nM 236 290\nL 240 290\nL 240 294\nL 236 294\nL 236 290\nM 322 283\nL 326 283\nL 326 287\nL 322 287\nL 322 283\nM 407 272\nL 411 272\nL 411 276\nL 407 276\nL 407 272\nM 492 284\nL 496 284\nL 496 288\nL 492 288\nL 492 284\nM 578 286\nL 582 286\nL 582 290\nL 578 290\nL 578 286\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><path d=\"M 66 183\nL 70 183\nL 70 187\nL 66 187\nL 66 183\nM 151 160\nL 155 160\nL 155 164\nL 151 164\nL 151 160\nM 236 166\nL 240 166\nL 240 170\nL 236 170\nL 236 166\nM 322 159\nL 326 159\nL 326 163\nL 322 163\nL 322 159\nM 407 86\nL 411 86\nL 411 90\nL 407 90\nL 407 86\nM 492 77\nL 496 77\nL 496 81\nL 492 81\nL 492 77\nM 578 79\nL 582 79\nL 582 83\nL 578 83\nL 578 79\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:rgb(115,192,222)\"/></svg>", data)
}

func TestBarRender(t *testing.T) {
	t.Parallel()

	values := [][]float64{
		{2.0, 4.9, 7.0, 23.2, 25.6, 76.7, 135.6, 162.2, 32.6, 20.0, 6.4, 3.3},
		{2.6, 5.9, 9.0, 26.4, 28.7, 70.7, 175.6, 182.2, 48.7, 18.8, 6.0, 2.3},
	}
	p, err := BarRender(
		values,
		SVGOutputOptionFunc(),
		XAxisLabelsOptionFunc([]string{
			"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
		}),
		LegendLabelsOptionFunc([]string{
			"Rainfall", "Evaporation",
		}),
		MarkLineOptionFunc(0, SeriesMarkTypeAverage),
		MarkPointOptionFunc(0, SeriesMarkTypeMax, SeriesMarkTypeMin),
		// custom option func
		func(opt *ChartOption) {
			opt.Legend.Offset = OffsetRight
			opt.Legend.OverlayChart = Ptr(true)
			opt.SeriesList[1].MarkPoint = NewMarkPoint(SeriesMarkTypeMax, SeriesMarkTypeMin)
			opt.SeriesList[1].MarkLine = NewMarkLine(SeriesMarkTypeAverage)
		},
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 364 29\nL 394 29\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"379\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"396\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Rainfall</text><path d=\"M 468 29\nL 498 29\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"483\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"500\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Evaporation</text><text x=\"19\" y=\"26\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">270</text><text x=\"19\" y=\"62\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">240</text><text x=\"19\" y=\"99\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">210</text><text x=\"19\" y=\"136\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">180</text><text x=\"19\" y=\"173\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">150</text><text x=\"19\" y=\"210\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">120</text><text x=\"28\" y=\"247\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">90</text><text x=\"28\" y=\"284\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">60</text><text x=\"28\" y=\"321\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">30</text><text x=\"37\" y=\"358\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path d=\"M 52 20\nL 580 20\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 57\nL 580 57\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 94\nL 580 94\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 131\nL 580 131\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 168\nL 580 168\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 205\nL 580 205\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 242\nL 580 242\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 279\nL 580 279\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 316\nL 580 316\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 56 354\nL 580 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 56 359\nL 56 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 99 359\nL 99 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 143 359\nL 143 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 187 359\nL 187 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 230 359\nL 230 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 274 359\nL 274 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 318 359\nL 318 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 361 359\nL 361 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 405 359\nL 405 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 449 359\nL 449 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 492 359\nL 492 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 536 359\nL 536 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 580 359\nL 580 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"64\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jan</text><text x=\"108\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Feb</text><text x=\"151\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mar</text><text x=\"196\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Apr</text><text x=\"237\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">May</text><text x=\"283\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jun</text><text x=\"329\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jul</text><text x=\"369\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Aug</text><text x=\"414\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sep</text><text x=\"458\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Oct</text><text x=\"500\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Nov</text><text x=\"545\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Dec</text><path d=\"M 61 352\nL 76 352\nL 76 353\nL 61 353\nL 61 352\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 104 348\nL 119 348\nL 119 353\nL 104 353\nL 104 348\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 148 346\nL 163 346\nL 163 353\nL 148 353\nL 148 346\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 192 326\nL 207 326\nL 207 353\nL 192 353\nL 192 326\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 235 323\nL 250 323\nL 250 353\nL 235 353\nL 235 323\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 279 260\nL 294 260\nL 294 353\nL 279 353\nL 279 260\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 323 187\nL 338 187\nL 338 353\nL 323 353\nL 323 187\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 366 154\nL 381 154\nL 381 353\nL 366 353\nL 366 154\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 410 314\nL 425 314\nL 425 353\nL 410 353\nL 410 314\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 454 330\nL 469 330\nL 469 353\nL 454 353\nL 454 330\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 497 347\nL 512 347\nL 512 353\nL 497 353\nL 497 347\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 541 350\nL 556 350\nL 556 353\nL 541 353\nL 541 350\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 79 351\nL 94 351\nL 94 353\nL 79 353\nL 79 351\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 122 347\nL 137 347\nL 137 353\nL 122 353\nL 122 347\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 166 343\nL 181 343\nL 181 353\nL 166 353\nL 166 343\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 210 322\nL 225 322\nL 225 353\nL 210 353\nL 210 322\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 253 319\nL 268 319\nL 268 353\nL 253 353\nL 253 319\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 297 267\nL 312 267\nL 312 353\nL 297 353\nL 297 267\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 341 137\nL 356 137\nL 356 353\nL 341 353\nL 341 137\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 384 129\nL 399 129\nL 399 353\nL 384 353\nL 384 129\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 428 294\nL 443 294\nL 443 353\nL 428 353\nL 428 294\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 472 331\nL 487 331\nL 487 353\nL 472 353\nL 472 331\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 515 347\nL 530 347\nL 530 353\nL 515 353\nL 515 347\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 559 352\nL 574 352\nL 574 353\nL 559 353\nL 559 352\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 369 147\nA 14 14 330.00 1 1 377 147\nL 373 133\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 359 133\nQ373,168 387,133\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><text x=\"360\" y=\"138\" style=\"stroke:none;fill:rgb(238,238,238);font-size:10.2px;font-family:'Roboto Medium',sans-serif\">162.2</text><path d=\"M 64 345\nA 14 14 330.00 1 1 72 345\nL 68 331\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 54 331\nQ68,366 82,331\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><text x=\"64\" y=\"336\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2</text><path d=\"M 387 122\nA 14 14 330.00 1 1 395 122\nL 391 108\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 377 108\nQ391,143 405,108\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"378\" y=\"113\" style=\"stroke:none;fill:rgb(70,70,70);font-size:10.2px;font-family:'Roboto Medium',sans-serif\">182.2</text><path d=\"M 562 345\nA 14 14 330.00 1 1 570 345\nL 566 331\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 552 331\nQ566,366 580,331\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"557\" y=\"336\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2.3</text><circle cx=\"59\" cy=\"303\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 65 303\nL 562 303\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 562 298\nL 578 303\nL 562 308\nL 567 303\nL 562 298\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"580\" y=\"307\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">41.63</text><circle cx=\"59\" cy=\"295\" r=\"3\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 65 295\nL 562 295\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 562 290\nL 578 295\nL 562 300\nL 567 295\nL 562 290\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"580\" y=\"299\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">48.07</text></svg>", data)
}

func TestHorizontalBarRender(t *testing.T) {
	t.Parallel()

	values := [][]float64{
		{18203, 23489, 29034, 104970, 131744, 630230},
		{19325, 23438, 31000, 121594, 134141, 681807},
	}
	p, err := HorizontalBarRender(
		values,
		SVGOutputOptionFunc(),
		TitleTextOptionFunc("World Population"),
		PaddingOptionFunc(Box{
			Top:    20,
			Right:  40,
			Bottom: 20,
			Left:   20,
		}),
		LegendLabelsOptionFunc([]string{"2011", "2012"}),
		YAxisLabelsOptionFunc([]string{"Brazil", "Indonesia", "USA", "India", "China", "World"}),
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"20\" y=\"36\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World Population</text><path d=\"M 214 29\nL 244 29\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"229\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"246\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2011</text><path d=\"M 301 29\nL 331 29\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"316\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"333\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2012</text><path d=\"M 97 56\nL 97 356\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 92 56\nL 97 56\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 92 106\nL 97 106\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 92 156\nL 97 156\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 92 206\nL 97 206\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 92 256\nL 97 256\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 92 306\nL 97 306\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 92 356\nL 97 356\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"46\" y=\"86\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World</text><text x=\"47\" y=\"136\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">China</text><text x=\"53\" y=\"186\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">India</text><text x=\"57\" y=\"235\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">USA</text><text x=\"19\" y=\"285\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Indonesia</text><text x=\"48\" y=\"335\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Brazil</text><text x=\"97\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"174\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">120k</text><text x=\"251\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">240k</text><text x=\"328\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">360k</text><text x=\"405\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480k</text><text x=\"482\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">600k</text><text x=\"525\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">720k</text><path d=\"M 175 56\nL 175 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 252 56\nL 252 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 329 56\nL 329 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 406 56\nL 406 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 483 56\nL 483 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 560 56\nL 560 352\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 98 316\nL 109 316\nL 109 328\nL 98 328\nL 98 316\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 98 266\nL 113 266\nL 113 278\nL 98 278\nL 98 266\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 98 216\nL 116 216\nL 116 228\nL 98 228\nL 98 216\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 98 166\nL 165 166\nL 165 178\nL 98 178\nL 98 166\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 98 116\nL 182 116\nL 182 128\nL 98 128\nL 98 116\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 98 66\nL 502 66\nL 502 78\nL 98 78\nL 98 66\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 98 333\nL 110 333\nL 110 345\nL 98 345\nL 98 333\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 98 283\nL 113 283\nL 113 295\nL 98 295\nL 98 283\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 98 233\nL 117 233\nL 117 245\nL 98 245\nL 98 233\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 98 183\nL 176 183\nL 176 195\nL 98 195\nL 98 183\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 98 133\nL 184 133\nL 184 145\nL 98 145\nL 98 133\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 98 83\nL 535 83\nL 535 95\nL 98 95\nL 98 83\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>", data)
}

func TestPieRender(t *testing.T) {
	t.Parallel()

	values := []float64{1048, 735, 580, 484, 300}
	p, err := PieRender(
		values,
		SVGOutputOptionFunc(),
		TitleOptionFunc(TitleOption{
			Text:    "Rainfall vs Evaporation",
			Subtext: "Fake Data",
			Offset:  OffsetCenter,
		}),
		PaddingOptionFunc(NewBoxEqual(20)),
		LegendOptionFunc(LegendOption{
			Vertical:    Ptr(true),
			SeriesNames: []string{"Search Engine", "Direct", "Email", "Union Ads", "Video Ads"},
			Offset:      OffsetLeft,
		}),
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"222\" y=\"36\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Rainfall vs Evaporation</text><text x=\"266\" y=\"52\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fake Data</text><path d=\"M 20 29\nL 50 29\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"35\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"52\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><path d=\"M 20 49\nL 50 49\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"35\" cy=\"49\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"52\" y=\"55\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path d=\"M 20 69\nL 50 69\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"35\" cy=\"69\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"52\" y=\"75\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path d=\"M 20 89\nL 50 89\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"35\" cy=\"89\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"52\" y=\"95\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path d=\"M 20 109\nL 50 109\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:none\"/><circle cx=\"35\" cy=\"109\" r=\"5\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:rgb(115,192,222)\"/><text x=\"52\" y=\"115\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><path d=\"M 300 223\nL 300 98\nA 125 125 119.89 0 1 409 285\nL 300 223\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 408 161\nL 421 153\nM 421 153\nL 436 153\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:none\"/><text x=\"439\" y=\"158\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">: 33.3%</text><path d=\"M 300 223\nL 409 285\nA 125 125 84.08 0 1 249 337\nL 300 223\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 338 342\nL 343 356\nM 343 356\nL 358 356\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:none\"/><text x=\"361\" y=\"361\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">: 23.35%</text><path d=\"M 300 223\nL 249 337\nA 125 125 66.35 0 1 175 222\nL 300 223\nZ\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path d=\"M 195 290\nL 183 299\nM 183 299\nL 168 299\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:none\"/><text x=\"116\" y=\"304\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">: 18.43%</text><path d=\"M 300 223\nL 175 222\nA 125 125 55.37 0 1 229 120\nL 300 223\nZ\" style=\"stroke:none;fill:rgb(238,102,102)\"/><path d=\"M 190 165\nL 177 158\nM 177 158\nL 162 158\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:none\"/><text x=\"110\" y=\"163\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">: 15.37%</text><path d=\"M 300 223\nL 229 120\nA 125 125 34.32 0 1 300 98\nL 300 223\nZ\" style=\"stroke:none;fill:rgb(115,192,222)\"/><path d=\"M 264 104\nL 259 90\nM 259 90\nL 244 90\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:none\"/><text x=\"199\" y=\"95\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">: 9.53%</text></svg>", data)
}

func TestDoughnutRender(t *testing.T) {
	t.Parallel()

	values := []float64{1048, 735, 580, 484, 300}
	p, err := DoughnutRender(
		values,
		SVGOutputOptionFunc(),
		TitleOptionFunc(TitleOption{
			Text:    "Title",
			Subtext: "Fake Data",
			Offset:  OffsetCenter,
		}),
		PaddingOptionFunc(NewBoxEqual(20)),
		LegendOptionFunc(LegendOption{
			Vertical:    Ptr(true),
			SeriesNames: []string{"A", "B", "C", "D", "E"},
			Offset:      OffsetLeft,
		}),
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"285\" y=\"36\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Title</text><text x=\"266\" y=\"52\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fake Data</text><path d=\"M 20 29\nL 50 29\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"35\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"52\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">A</text><path d=\"M 20 49\nL 50 49\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"35\" cy=\"49\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"52\" y=\"55\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">B</text><path d=\"M 20 69\nL 50 69\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"35\" cy=\"69\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"52\" y=\"75\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">C</text><path d=\"M 20 89\nL 50 89\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"35\" cy=\"89\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"52\" y=\"95\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">D</text><path d=\"M 20 109\nL 50 109\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:none\"/><circle cx=\"35\" cy=\"109\" r=\"5\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:rgb(115,192,222)\"/><text x=\"52\" y=\"115\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">E</text><path d=\"M 300 223\nL 300 98\nA 125 125 119.89 0 1 409 285\nL 300 223\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 408 161\nL 421 153\nM 421 153\nL 436 153\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:none\"/><text x=\"439\" y=\"158\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">: 33.3%</text><path d=\"M 300 223\nL 409 285\nA 125 125 84.08 0 1 249 337\nL 300 223\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 338 342\nL 343 356\nM 343 356\nL 358 356\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:none\"/><text x=\"361\" y=\"361\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">: 23.35%</text><path d=\"M 300 223\nL 249 337\nA 125 125 66.35 0 1 175 222\nL 300 223\nZ\" style=\"stroke:none;fill:rgb(250,200,88)\"/><path d=\"M 195 290\nL 183 299\nM 183 299\nL 168 299\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:none\"/><text x=\"116\" y=\"304\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">: 18.43%</text><path d=\"M 300 223\nL 175 222\nA 125 125 55.37 0 1 229 120\nL 300 223\nZ\" style=\"stroke:none;fill:rgb(238,102,102)\"/><path d=\"M 190 165\nL 177 158\nM 177 158\nL 162 158\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:none\"/><text x=\"110\" y=\"163\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">: 15.37%</text><path d=\"M 300 223\nL 229 120\nA 125 125 34.32 0 1 300 98\nL 300 223\nZ\" style=\"stroke:none;fill:rgb(115,192,222)\"/><path d=\"M 264 104\nL 259 90\nM 259 90\nL 244 90\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:none\"/><text x=\"199\" y=\"95\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">: 9.53%</text><circle cx=\"300\" cy=\"223\" r=\"75\" style=\"stroke:none;fill:white\"/></svg>", data)
}

func TestRadarRender(t *testing.T) {
	t.Parallel()

	values := [][]float64{
		{4200, 3000, 20000, 35000, 50000, 18000},
		{5000, 14000, 28000, 26000, 42000, 21000},
	}
	p, err := RadarRender(
		values,
		SVGOutputOptionFunc(),
		TitleTextOptionFunc("Basic Radar Chart"),
		LegendLabelsOptionFunc([]string{
			"Allocated Budget", "Actual Spending",
		}),
		RadarIndicatorOptionFunc([]string{
			"Sales",
			"Administration",
			"Information Technology",
			"Customer Support",
			"Development",
			"Marketing",
		}, []float64{
			6500, 16000, 30000, 38000, 52000, 25000,
		}),
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"20\" y=\"36\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Basic Radar Chart</text><path d=\"M 143 29\nL 173 29\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"158\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"175\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Allocated Budget</text><path d=\"M 313 29\nL 343 29\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"328\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"345\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Actual Spending</text><path d=\"M 300 193\nL 321 206\nL 321 230\nL 300 243\nL 279 230\nL 279 206\nL 300 193\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 300 168\nL 343 193\nL 343 242\nL 300 268\nL 257 243\nL 257 194\nL 300 168\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 300 143\nL 364 181\nL 364 255\nL 300 293\nL 236 255\nL 236 181\nL 300 143\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 300 118\nL 386 168\nL 386 267\nL 300 318\nL 214 268\nL 214 169\nL 300 118\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 300 93\nL 408 156\nL 408 280\nL 300 343\nL 192 280\nL 192 156\nL 300 93\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 300 218\nL 300 93\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 300 218\nL 408 156\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 300 218\nL 408 280\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 300 218\nL 300 343\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 300 218\nL 192 280\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 300 218\nL 192 156\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><text x=\"284\" y=\"85\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Sales</text><text x=\"413\" y=\"161\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Administration</text><text x=\"413\" y=\"285\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Information Technology</text><text x=\"248\" y=\"361\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Customer Support</text><text x=\"111\" y=\"285\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Development</text><text x=\"129\" y=\"161\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Marketing</text><path d=\"M 300 138\nL 320 207\nL 372 259\nL 300 333\nL 196 278\nL 223 174\nL 300 138\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><path d=\"M 300 138\nL 320 207\nL 372 259\nL 300 333\nL 196 278\nL 223 174\nL 300 138\" style=\"stroke:none;fill:rgba(84,112,198,0.1)\"/><circle cx=\"300\" cy=\"138\" r=\"2\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"320\" cy=\"207\" r=\"2\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"372\" cy=\"259\" r=\"2\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"300\" cy=\"333\" r=\"2\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"196\" cy=\"278\" r=\"2\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"223\" cy=\"174\" r=\"2\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"300\" cy=\"138\" r=\"2\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:white\"/><path d=\"M 300 122\nL 394 164\nL 401 276\nL 300 303\nL 213 268\nL 210 166\nL 300 122\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/><path d=\"M 300 122\nL 394 164\nL 401 276\nL 300 303\nL 213 268\nL 210 166\nL 300 122\" style=\"stroke:none;fill:rgba(145,204,117,0.1)\"/><circle cx=\"300\" cy=\"122\" r=\"2\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"394\" cy=\"164\" r=\"2\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"401\" cy=\"276\" r=\"2\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"300\" cy=\"303\" r=\"2\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"213\" cy=\"268\" r=\"2\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"210\" cy=\"166\" r=\"2\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"300\" cy=\"122\" r=\"2\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:white\"/></svg>", data)
}

func TestFunnelRender(t *testing.T) {
	t.Parallel()

	values := []float64{
		100, 80, 60, 40, 20,
	}
	p, err := FunnelRender(
		values,
		SVGOutputOptionFunc(),
		TitleTextOptionFunc("Funnel"),
		LegendLabelsOptionFunc([]string{
			"Show", "Click", "Visit", "Inquiry", "Order",
		}),
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"20\" y=\"36\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Funnel</text><path d=\"M 86 29\nL 116 29\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"101\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"118\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Show</text><path d=\"M 176 29\nL 206 29\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"191\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"208\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Click</text><path d=\"M 262 29\nL 292 29\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"277\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"294\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Visit</text><path d=\"M 345 29\nL 375 29\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"360\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"377\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Inquiry</text><path d=\"M 444 29\nL 474 29\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:none\"/><circle cx=\"459\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:rgb(115,192,222)\"/><text x=\"476\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Order</text><path d=\"M 20 56\nL 580 56\nL 524 119\nL 76 119\nL 20 56\" style=\"stroke:none;fill:rgb(84,112,198)\"/><text x=\"280\" y=\"87\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">(100%)</text><path d=\"M 76 121\nL 524 121\nL 468 184\nL 132 184\nL 76 121\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"284\" y=\"152\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">(80%)</text><path d=\"M 132 186\nL 468 186\nL 412 249\nL 188 249\nL 132 186\" style=\"stroke:none;fill:rgb(250,200,88)\"/><text x=\"284\" y=\"217\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">(60%)</text><path d=\"M 188 251\nL 412 251\nL 356 314\nL 244 314\nL 188 251\" style=\"stroke:none;fill:rgb(238,102,102)\"/><text x=\"284\" y=\"282\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">(40%)</text><path d=\"M 244 316\nL 356 316\nL 300 379\nL 300 379\nL 244 316\" style=\"stroke:none;fill:rgb(115,192,222)\"/><text x=\"284\" y=\"347\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">(20%)</text></svg>", data)
}

func TestChildRender(t *testing.T) {
	p, err := LineRender(
		[][]float64{
			{120, 132, 101, 134, 90, 230, 210},
			{150, 232, 201, 154, 190, 330, 410},
			{320, 332, 301, 334, 390, 330, 320},
		},
		SVGOutputOptionFunc(),
		XAxisLabelsOptionFunc([]string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}),
		ChildOptionFunc(ChartOption{
			Box: NewBox(200, 10, 500, 200),
			SeriesList: NewSeriesListHorizontalBar([][]float64{
				{70, 90, 110, 130},
				{80, 100, 120, 140},
			}).ToGenericSeriesList(),
			Legend: LegendOption{
				SeriesNames: []string{"2011", "2012"},
			},
			YAxis: []YAxisOption{
				{
					Labels: []string{"USA", "India", "China", "World"},
				},
			},
		}),
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"19\" y=\"26\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">450</text><text x=\"19\" y=\"62\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">410</text><text x=\"19\" y=\"99\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">370</text><text x=\"19\" y=\"136\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">330</text><text x=\"19\" y=\"173\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">290</text><text x=\"19\" y=\"210\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">250</text><text x=\"19\" y=\"247\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">210</text><text x=\"19\" y=\"284\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">170</text><text x=\"19\" y=\"321\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">130</text><text x=\"28\" y=\"358\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">90</text><path d=\"M 52 20\nL 580 20\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 57\nL 580 57\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 94\nL 580 94\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 131\nL 580 131\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 168\nL 580 168\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 205\nL 580 205\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 242\nL 580 242\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 279\nL 580 279\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 316\nL 580 316\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 56 354\nL 580 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 56 359\nL 56 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 130 359\nL 130 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 205 359\nL 205 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 280 359\nL 280 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 355 359\nL 355 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 430 359\nL 430 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 505 359\nL 505 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 580 359\nL 580 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"78\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"154\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"227\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"304\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"383\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"456\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"529\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path d=\"M 93 327\nL 167 316\nL 242 344\nL 317 314\nL 392 354\nL 467 225\nL 542 243\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"93\" cy=\"327\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"167\" cy=\"316\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"242\" cy=\"344\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"317\" cy=\"314\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"392\" cy=\"354\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"467\" cy=\"225\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"542\" cy=\"243\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><path d=\"M 93 299\nL 167 223\nL 242 252\nL 317 295\nL 392 262\nL 467 132\nL 542 58\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"93\" cy=\"299\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"167\" cy=\"223\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"242\" cy=\"252\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"317\" cy=\"295\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"392\" cy=\"262\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"467\" cy=\"132\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"542\" cy=\"58\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><path d=\"M 93 141\nL 167 130\nL 242 159\nL 317 128\nL 392 76\nL 467 132\nL 542 141\" style=\"stroke-width:2;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"93\" cy=\"141\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"167\" cy=\"130\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"242\" cy=\"159\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"317\" cy=\"128\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"392\" cy=\"76\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"467\" cy=\"132\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"542\" cy=\"141\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><path d=\"M 274 39\nL 304 39\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"289\" cy=\"39\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"306\" y=\"45\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2011</text><path d=\"M 361 39\nL 391 39\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"376\" cy=\"39\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"393\" y=\"45\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2012</text><path d=\"M 270 66\nL 270 156\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 265 66\nL 270 66\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 265 88\nL 270 88\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 265 111\nL 270 111\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 265 133\nL 270 133\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 265 156\nL 270 156\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"219\" y=\"83\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World</text><text x=\"220\" y=\"105\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">China</text><text x=\"226\" y=\"127\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">India</text><text x=\"230\" y=\"149\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">USA</text><text x=\"270\" y=\"175\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">70</text><text x=\"374\" y=\"175\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">110</text><text x=\"453\" y=\"175\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">150</text><path d=\"M 375 66\nL 375 152\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 480 66\nL 480 152\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 271 138\nL 271 138\nL 271 142\nL 271 142\nL 271 138\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 271 116\nL 323 116\nL 323 120\nL 271 120\nL 271 116\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 271 93\nL 375 93\nL 375 97\nL 271 97\nL 271 93\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 271 71\nL 427 71\nL 427 75\nL 271 75\nL 271 71\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 271 145\nL 297 145\nL 297 149\nL 271 149\nL 271 145\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 271 123\nL 349 123\nL 349 127\nL 271 127\nL 271 123\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 271 100\nL 401 100\nL 401 104\nL 271 104\nL 271 100\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 271 78\nL 453 78\nL 453 82\nL 271 82\nL 271 78\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>", data)
}

func TestChartCombos(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		optFunc     func() ChartOption
		expectedSVG string
	}{
		{
			name: "line+scatter",
			optFunc: func() ChartOption {
				lineValues := [][]float64{
					{120, 132, 101, 134, 90, 230, 210},
				}
				scatterValues := [][]float64{
					{180, 200, 150, 180, 160, 280, 250},
				}

				lineSeries := NewSeriesListLine(lineValues, LineSeriesOption{
					Names: []string{"Line Data"},
				}).ToGenericSeriesList()
				scatterSeries := NewSeriesListScatter(scatterValues, ScatterSeriesOption{
					Names: []string{"Scatter Data"},
				}).ToGenericSeriesList()

				return ChartOption{
					Width:  600,
					Height: 400,
					XAxis: XAxisOption{
						Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
					},
					YAxis: []YAxisOption{
						{
							Max: Ptr(500.0),
						},
					},
					Legend: LegendOption{
						SeriesNames: []string{"Line Data", "Scatter Data"},
					},
					SeriesList: append(lineSeries, scatterSeries...),
				}
			},
			expectedSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 184 29\nL 214 29\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"199\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"216\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line Data</text><path d=\"M 301 29\nL 331 29\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"316\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"333\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Scatter Data</text><text x=\"40\" y=\"62\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">500</text><text x=\"19\" y=\"94\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">454.44</text><text x=\"19\" y=\"127\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">408.89</text><text x=\"19\" y=\"160\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">363.33</text><text x=\"19\" y=\"193\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">317.78</text><text x=\"19\" y=\"226\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">272.22</text><text x=\"19\" y=\"259\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">226.67</text><text x=\"19\" y=\"292\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">181.11</text><text x=\"19\" y=\"325\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">135.56</text><text x=\"49\" y=\"358\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">90</text><path d=\"M 73 56\nL 580 56\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 89\nL 580 89\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 122\nL 580 122\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 155\nL 580 155\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 188\nL 580 188\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 221\nL 580 221\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 254\nL 580 254\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 287\nL 580 287\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 320\nL 580 320\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 77 354\nL 580 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 77 359\nL 77 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 148 359\nL 148 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 220 359\nL 220 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 292 359\nL 292 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 364 359\nL 364 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 436 359\nL 436 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 508 359\nL 508 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 580 359\nL 580 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"97\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"171\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"241\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"315\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"391\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"461\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"531\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path d=\"M 112 333\nL 184 324\nL 256 347\nL 328 323\nL 400 354\nL 472 253\nL 544 267\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"112\" cy=\"333\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"184\" cy=\"324\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"256\" cy=\"347\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"328\" cy=\"323\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"400\" cy=\"354\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"472\" cy=\"253\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"544\" cy=\"267\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"77\" cy=\"289\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"160\" cy=\"275\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"244\" cy=\"311\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"328\" cy=\"289\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"412\" cy=\"304\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"496\" cy=\"216\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"580\" cy=\"238\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/></svg>",
		},
		{
			name: "line+bar",
			optFunc: func() ChartOption {
				lineValues := [][]float64{
					{120, 132, 101, 134, 90, 230, 210},
				}
				barValues := [][]float64{
					{70, 90, 110, 130, 80, 100, 120},
				}

				lineSeries := NewSeriesListLine(lineValues, LineSeriesOption{
					Names: []string{"Line Data"},
				}).ToGenericSeriesList()
				barSeries := NewSeriesListBar(barValues, BarSeriesOption{
					Names: []string{"Bar Data"},
				}).ToGenericSeriesList()

				return ChartOption{
					Width:  600,
					Height: 400,
					XAxis: XAxisOption{
						Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
					},
					YAxis: []YAxisOption{
						{
							Max: Ptr(400.0),
						},
					},
					Legend: LegendOption{
						SeriesNames: []string{"Line Data", "Bar Data"},
					},
					SeriesList: append(lineSeries, barSeries...),
				}
			},
			expectedSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 197 29\nL 227 29\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"212\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"229\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line Data</text><path d=\"M 314 29\nL 344 29\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"329\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"346\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Bar Data</text><text x=\"40\" y=\"62\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">400</text><text x=\"19\" y=\"94\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">356.11</text><text x=\"19\" y=\"127\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">312.22</text><text x=\"19\" y=\"160\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">268.33</text><text x=\"19\" y=\"193\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">224.44</text><text x=\"19\" y=\"226\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">180.56</text><text x=\"19\" y=\"259\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">136.67</text><text x=\"27\" y=\"292\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">92.78</text><text x=\"27\" y=\"325\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">48.89</text><text x=\"58\" y=\"358\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">5</text><path d=\"M 73 56\nL 580 56\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 89\nL 580 89\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 122\nL 580 122\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 155\nL 580 155\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 188\nL 580 188\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 221\nL 580 221\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 254\nL 580 254\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 287\nL 580 287\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 320\nL 580 320\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 77 354\nL 580 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 77 359\nL 77 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 148 359\nL 148 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 220 359\nL 220 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 292 359\nL 292 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 364 359\nL 364 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 436 359\nL 436 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 508 359\nL 508 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 580 359\nL 580 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"97\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"171\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"241\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"315\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"391\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"461\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"531\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path d=\"M 87 305\nL 138 305\nL 138 353\nL 87 353\nL 87 305\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 158 290\nL 209 290\nL 209 353\nL 158 353\nL 158 290\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 230 275\nL 281 275\nL 281 353\nL 230 353\nL 230 275\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 302 260\nL 353 260\nL 353 353\nL 302 353\nL 302 260\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 374 298\nL 425 298\nL 425 353\nL 374 353\nL 374 298\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 446 283\nL 497 283\nL 497 353\nL 446 353\nL 446 283\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 518 268\nL 569 268\nL 569 353\nL 518 353\nL 518 268\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 112 268\nL 184 259\nL 256 282\nL 328 257\nL 400 290\nL 472 185\nL 544 200\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"112\" cy=\"268\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"184\" cy=\"259\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"256\" cy=\"282\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"328\" cy=\"257\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"400\" cy=\"290\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"472\" cy=\"185\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"544\" cy=\"200\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/></svg>",
		},
		{
			name: "line+candlestick",
			optFunc: func() ChartOption {
				lineValues := [][]float64{
					{25, 28, 23, 30, 27, 32, 29},
				}
				candlestickData := [][]OHLCData{
					{
						{Open: 24, High: 30, Low: 20, Close: 25},
						{Open: 25, High: 32, Low: 22, Close: 28},
						{Open: 28, High: 30, Low: 18, Close: 23},
						{Open: 23, High: 35, Low: 25, Close: 30},
						{Open: 30, High: 33, Low: 24, Close: 27},
						{Open: 27, High: 36, Low: 28, Close: 32},
						{Open: 32, High: 34, Low: 26, Close: 29},
					},
				}

				lineSeries := NewSeriesListLine(lineValues, LineSeriesOption{
					Names: []string{"Line Data"},
				}).ToGenericSeriesList()
				candlestickSeries := NewSeriesListCandlestick(candlestickData, CandlestickSeriesOption{
					Names: []string{"OHLC"},
				}).ToGenericSeriesList()

				return ChartOption{
					Width:  600,
					Height: 400,
					XAxis: XAxisOption{
						Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
					},
					YAxis: []YAxisOption{
						{
							Max: Ptr(40.0),
						},
					},
					Legend: LegendOption{
						SeriesNames: []string{"Line Data", "OHLC"},
					},
					SeriesList: append(lineSeries, candlestickSeries...),
				}
			},
			expectedSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 206 29\nL 236 29\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"221\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"238\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line Data</text><path d=\"M 323 36\nL 338 36\nL 330 23\nL 323 36\" style=\"stroke:none;fill:rgb(115,192,222)\"/><path d=\"M 338 23\nL 353 23\nL 345 36\nL 338 23\" style=\"stroke:none;fill:rgb(252,132,82)\"/><text x=\"355\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">OHLC</text><text x=\"41\" y=\"62\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">40</text><text x=\"19\" y=\"94\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">37.56</text><text x=\"19\" y=\"127\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">35.11</text><text x=\"19\" y=\"160\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">32.67</text><text x=\"19\" y=\"193\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">30.22</text><text x=\"19\" y=\"226\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">27.78</text><text x=\"19\" y=\"259\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">25.33</text><text x=\"19\" y=\"292\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">22.89</text><text x=\"19\" y=\"325\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">20.44</text><text x=\"41\" y=\"358\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">18</text><path d=\"M 65 56\nL 580 56\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 89\nL 580 89\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 122\nL 580 122\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 155\nL 580 155\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 188\nL 580 188\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 221\nL 580 221\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 254\nL 580 254\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 287\nL 580 287\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 320\nL 580 320\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 69 354\nL 580 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 69 359\nL 69 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 142 359\nL 142 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 215 359\nL 215 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 288 359\nL 288 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 361 359\nL 361 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 434 359\nL 434 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 507 359\nL 507 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 580 359\nL 580 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"90\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"165\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"236\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"311\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"388\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"459\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"530\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path d=\"M 105 192\nL 105 260\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:none\"/><path d=\"M 105 273\nL 105 327\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:none\"/><path d=\"M 91 192\nL 119 192\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:none\"/><path d=\"M 91 327\nL 119 327\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:none\"/><path d=\"M 76 260\nL 134 260\nL 134 273\nL 76 273\nL 76 260\" style=\"stroke:none;fill:rgb(115,192,222)\"/><path d=\"M 178 165\nL 178 219\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:none\"/><path d=\"M 178 260\nL 178 300\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:none\"/><path d=\"M 164 165\nL 192 165\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:none\"/><path d=\"M 164 300\nL 192 300\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:none\"/><path d=\"M 149 219\nL 207 219\nL 207 260\nL 149 260\nL 149 219\" style=\"stroke:none;fill:rgb(115,192,222)\"/><path d=\"M 251 192\nL 251 219\" style=\"stroke-width:1;stroke:rgb(252,132,82);fill:none\"/><path d=\"M 251 287\nL 251 354\" style=\"stroke-width:1;stroke:rgb(252,132,82);fill:none\"/><path d=\"M 237 192\nL 265 192\" style=\"stroke-width:1;stroke:rgb(252,132,82);fill:none\"/><path d=\"M 237 354\nL 265 354\" style=\"stroke-width:1;stroke:rgb(252,132,82);fill:none\"/><path d=\"M 222 219\nL 280 219\nL 280 287\nL 222 287\nL 222 219\" style=\"stroke:none;fill:rgb(252,132,82)\"/><path d=\"M 397 151\nL 397 192\" style=\"stroke-width:1;stroke:rgb(252,132,82);fill:none\"/><path d=\"M 397 233\nL 397 273\" style=\"stroke-width:1;stroke:rgb(252,132,82);fill:none\"/><path d=\"M 383 151\nL 411 151\" style=\"stroke-width:1;stroke:rgb(252,132,82);fill:none\"/><path d=\"M 383 273\nL 411 273\" style=\"stroke-width:1;stroke:rgb(252,132,82);fill:none\"/><path d=\"M 368 192\nL 426 192\nL 426 233\nL 368 233\nL 368 192\" style=\"stroke:none;fill:rgb(252,132,82)\"/><path d=\"M 543 138\nL 543 165\" style=\"stroke-width:1;stroke:rgb(252,132,82);fill:none\"/><path d=\"M 543 205\nL 543 246\" style=\"stroke-width:1;stroke:rgb(252,132,82);fill:none\"/><path d=\"M 529 138\nL 557 138\" style=\"stroke-width:1;stroke:rgb(252,132,82);fill:none\"/><path d=\"M 529 246\nL 557 246\" style=\"stroke-width:1;stroke:rgb(252,132,82);fill:none\"/><path d=\"M 514 165\nL 572 165\nL 572 205\nL 514 205\nL 514 165\" style=\"stroke:none;fill:rgb(252,132,82)\"/><path d=\"M 105 260\nL 178 219\nL 251 287\nL 324 192\nL 397 233\nL 470 165\nL 543 205\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"105\" cy=\"260\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"178\" cy=\"219\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"251\" cy=\"287\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"324\" cy=\"192\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"397\" cy=\"233\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"470\" cy=\"165\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"543\" cy=\"205\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/></svg>",
		},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			opt := tc.optFunc()
			opt.OutputFormat = ChartOutputSVG

			p, err := Render(opt)
			require.NoError(t, err)
			data, err := p.Bytes()
			require.NoError(t, err)
			assertEqualSVG(t, tc.expectedSVG, data)
		})
	}
}

func TestDualAxisChartCombos(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		optFunc     func() ChartOption
		expectedSVG string
	}{
		{
			name: "scatter+line",
			optFunc: func() ChartOption {
				scatterValues := [][]float64{
					{180, 200, 150, 180, 160, 280, 250},
				}
				lineValues := [][]float64{
					{1.2, 1.5, 1.1, 1.8, 1.4, 2.1, 1.7},
				}

				scatterSeries := NewSeriesListScatter(scatterValues, ScatterSeriesOption{
					Names: []string{"Scatter Data"},
				}).ToGenericSeriesList()
				lineSeries := NewSeriesListLine(lineValues, LineSeriesOption{
					Names: []string{"Line Data"},
				}).ToGenericSeriesList()
				lineSeries[0].YAxisIndex = 1

				return ChartOption{
					Width:  600,
					Height: 400,
					XAxis: XAxisOption{
						Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
					},
					YAxis: []YAxisOption{
						{
							Max: Ptr(300.0),
						},
						{
							Max: Ptr(3.0),
						},
					},
					Legend: LegendOption{
						SeriesNames: []string{"Scatter Data", "Line Data"},
					},
					SeriesList: append(scatterSeries, lineSeries...),
				}
			},
			expectedSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 184 29\nL 214 29\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"199\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"216\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Scatter Data</text><path d=\"M 321 29\nL 351 29\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"336\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"353\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line Data</text><text x=\"550\" y=\"62\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">3</text><text x=\"550\" y=\"94\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2.78</text><text x=\"550\" y=\"127\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2.56</text><text x=\"550\" y=\"160\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2.33</text><text x=\"550\" y=\"193\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2.11</text><text x=\"550\" y=\"226\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.89</text><text x=\"550\" y=\"259\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.67</text><text x=\"550\" y=\"292\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44</text><text x=\"550\" y=\"325\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.22</text><text x=\"550\" y=\"358\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1</text><text x=\"40\" y=\"62\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">300</text><text x=\"19\" y=\"94\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">283.33</text><text x=\"19\" y=\"127\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">266.67</text><text x=\"40\" y=\"160\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">250</text><text x=\"19\" y=\"193\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">233.33</text><text x=\"19\" y=\"226\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">216.67</text><text x=\"40\" y=\"259\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">200</text><text x=\"19\" y=\"292\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">183.33</text><text x=\"19\" y=\"325\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">166.67</text><text x=\"40\" y=\"358\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">150</text><path d=\"M 73 56\nL 540 56\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 89\nL 540 89\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 122\nL 540 122\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 155\nL 540 155\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 188\nL 540 188\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 221\nL 540 221\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 254\nL 540 254\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 287\nL 540 287\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 320\nL 540 320\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 77 354\nL 540 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 77 359\nL 77 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 143 359\nL 143 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 209 359\nL 209 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 275 359\nL 275 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 341 359\nL 341 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 407 359\nL 407 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 473 359\nL 473 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 540 359\nL 540 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"95\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"163\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"227\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"295\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"365\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"429\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"493\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path d=\"M 110 325\nL 176 280\nL 242 340\nL 308 235\nL 374 295\nL 440 191\nL 506 250\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"110\" cy=\"325\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"176\" cy=\"280\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"242\" cy=\"340\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"308\" cy=\"235\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"374\" cy=\"295\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"440\" cy=\"191\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"506\" cy=\"250\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"77\" cy=\"295\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"154\" cy=\"255\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"231\" cy=\"354\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"308\" cy=\"295\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"385\" cy=\"335\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"462\" cy=\"96\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"540\" cy=\"156\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/></svg>",
		},
		{
			name: "bar+line",
			optFunc: func() ChartOption {
				barValues := [][]float64{
					{70, 90, 110, 130, 80, 100, 120},
				}
				lineValues := [][]float64{
					{1.2, 1.8, 2.3, 1.9, 2.5, 1.6, 2.1},
				}

				barSeries := NewSeriesListBar(barValues, BarSeriesOption{
					Names: []string{"Bar Data"},
				}).ToGenericSeriesList()
				lineSeries := NewSeriesListLine(lineValues, LineSeriesOption{
					Names: []string{"Line Data"},
				}).ToGenericSeriesList()
				lineSeries[0].YAxisIndex = 1

				return ChartOption{
					Width:  600,
					Height: 400,
					XAxis: XAxisOption{
						Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
					},
					YAxis: []YAxisOption{
						{
							Max: Ptr(150.0),
						},
						{
							Max: Ptr(3.0),
						},
					},
					Legend: LegendOption{
						SeriesNames: []string{"Bar Data", "Line Data"},
					},
					SeriesList: append(barSeries, lineSeries...),
				}
			},
			expectedSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 197 29\nL 227 29\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"212\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"229\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Bar Data</text><path d=\"M 308 29\nL 338 29\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"323\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"340\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line Data</text><text x=\"550\" y=\"62\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">3</text><text x=\"550\" y=\"94\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2.78</text><text x=\"550\" y=\"127\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2.56</text><text x=\"550\" y=\"160\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2.33</text><text x=\"550\" y=\"193\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2.11</text><text x=\"550\" y=\"226\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.89</text><text x=\"550\" y=\"259\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.67</text><text x=\"550\" y=\"292\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44</text><text x=\"550\" y=\"325\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.22</text><text x=\"550\" y=\"358\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1</text><text x=\"40\" y=\"62\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">150</text><text x=\"19\" y=\"94\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">141.11</text><text x=\"19\" y=\"127\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">132.22</text><text x=\"19\" y=\"160\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">123.33</text><text x=\"19\" y=\"193\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">114.44</text><text x=\"19\" y=\"226\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">105.56</text><text x=\"27\" y=\"259\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">96.67</text><text x=\"27\" y=\"292\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">87.78</text><text x=\"27\" y=\"325\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">78.89</text><text x=\"49\" y=\"358\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">70</text><path d=\"M 73 56\nL 540 56\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 89\nL 540 89\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 122\nL 540 122\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 155\nL 540 155\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 188\nL 540 188\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 221\nL 540 221\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 254\nL 540 254\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 287\nL 540 287\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 73 320\nL 540 320\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 77 354\nL 540 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 77 359\nL 77 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 143 359\nL 143 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 209 359\nL 209 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 275 359\nL 275 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 341 359\nL 341 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 407 359\nL 407 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 473 359\nL 473 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 540 359\nL 540 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"95\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"163\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"227\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"295\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"365\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"429\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"493\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path d=\"M 87 354\nL 133 354\nL 133 353\nL 87 353\nL 87 354\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 153 280\nL 199 280\nL 199 353\nL 153 353\nL 153 280\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 219 205\nL 265 205\nL 265 353\nL 219 353\nL 219 205\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 285 131\nL 331 131\nL 331 353\nL 285 353\nL 285 131\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 351 317\nL 397 317\nL 397 353\nL 351 353\nL 351 317\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 417 243\nL 463 243\nL 463 353\nL 417 353\nL 417 243\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 483 168\nL 529 168\nL 529 353\nL 483 353\nL 483 168\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 110 325\nL 176 235\nL 242 161\nL 308 220\nL 374 131\nL 440 265\nL 506 191\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"110\" cy=\"325\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"176\" cy=\"235\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"242\" cy=\"161\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"308\" cy=\"220\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"374\" cy=\"131\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"440\" cy=\"265\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"506\" cy=\"191\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/></svg>",
		},
		{
			name: "candlestick+line",
			optFunc: func() ChartOption {
				candlestickData := [][]OHLCData{
					{
						{Open: 24, High: 30, Low: 20, Close: 25},
						{Open: 25, High: 32, Low: 22, Close: 28},
						{Open: 28, High: 30, Low: 18, Close: 23},
						{Open: 23, High: 35, Low: 25, Close: 30},
						{Open: 30, High: 33, Low: 24, Close: 27},
						{Open: 27, High: 36, Low: 28, Close: 32},
						{Open: 32, High: 34, Low: 26, Close: 29},
					},
				}
				lineValues := [][]float64{
					{1200, 1400, 1100, 1600, 1300, 1800, 1500},
				}

				candlestickSeries := NewSeriesListCandlestick(candlestickData, CandlestickSeriesOption{
					Names: []string{"OHLC"},
				}).ToGenericSeriesList()
				lineSeries := NewSeriesListLine(lineValues, LineSeriesOption{
					Names: []string{"Volume"},
				}).ToGenericSeriesList()
				lineSeries[0].YAxisIndex = 1

				return ChartOption{
					Width:  600,
					Height: 400,
					XAxis: XAxisOption{
						Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
					},
					YAxis: []YAxisOption{
						{
							Max: Ptr(40.0),
						},
						{
							Max: Ptr(2000.0),
						},
					},
					Legend: LegendOption{
						SeriesNames: []string{"OHLC", "Volume"},
					},
					SeriesList: append(candlestickSeries, lineSeries...),
				}
			},
			expectedSVG: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 212 36\nL 227 36\nL 219 23\nL 212 36\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 227 23\nL 242 23\nL 234 36\nL 227 23\" style=\"stroke:none;fill:rgb(238,102,102)\"/><text x=\"244\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">OHLC</text><path d=\"M 304 29\nL 334 29\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"319\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"336\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Volume</text><text x=\"542\" y=\"62\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2k</text><text x=\"542\" y=\"94\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.89k</text><text x=\"542\" y=\"127\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.78k</text><text x=\"542\" y=\"160\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.67k</text><text x=\"542\" y=\"193\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.56k</text><text x=\"542\" y=\"226\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"542\" y=\"259\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.33k</text><text x=\"542\" y=\"292\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.22k</text><text x=\"542\" y=\"325\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.11k</text><text x=\"542\" y=\"358\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1k</text><text x=\"41\" y=\"62\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">40</text><text x=\"19\" y=\"94\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">37.56</text><text x=\"19\" y=\"127\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">35.11</text><text x=\"19\" y=\"160\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">32.67</text><text x=\"19\" y=\"193\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">30.22</text><text x=\"19\" y=\"226\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">27.78</text><text x=\"19\" y=\"259\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">25.33</text><text x=\"19\" y=\"292\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">22.89</text><text x=\"19\" y=\"325\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">20.44</text><text x=\"41\" y=\"358\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">18</text><path d=\"M 65 56\nL 532 56\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 89\nL 532 89\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 122\nL 532 122\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 155\nL 532 155\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 188\nL 532 188\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 221\nL 532 221\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 254\nL 532 254\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 287\nL 532 287\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 320\nL 532 320\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 69 354\nL 532 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 69 359\nL 69 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 135 359\nL 135 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 201 359\nL 201 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 267 359\nL 267 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 333 359\nL 333 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 399 359\nL 399 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 465 359\nL 465 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 532 359\nL 532 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"87\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"155\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"219\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"287\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"357\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"421\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"485\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path d=\"M 102 192\nL 102 260\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:none\"/><path d=\"M 102 273\nL 102 327\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:none\"/><path d=\"M 89 192\nL 115 192\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:none\"/><path d=\"M 89 327\nL 115 327\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:none\"/><path d=\"M 76 260\nL 128 260\nL 128 273\nL 76 273\nL 76 260\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 168 165\nL 168 219\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:none\"/><path d=\"M 168 260\nL 168 300\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:none\"/><path d=\"M 155 165\nL 181 165\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:none\"/><path d=\"M 155 300\nL 181 300\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:none\"/><path d=\"M 142 219\nL 194 219\nL 194 260\nL 142 260\nL 142 219\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 234 192\nL 234 219\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:none\"/><path d=\"M 234 287\nL 234 354\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:none\"/><path d=\"M 221 192\nL 247 192\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:none\"/><path d=\"M 221 354\nL 247 354\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:none\"/><path d=\"M 208 219\nL 260 219\nL 260 287\nL 208 287\nL 208 219\" style=\"stroke:none;fill:rgb(238,102,102)\"/><path d=\"M 366 151\nL 366 192\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:none\"/><path d=\"M 366 233\nL 366 273\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:none\"/><path d=\"M 353 151\nL 379 151\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:none\"/><path d=\"M 353 273\nL 379 273\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:none\"/><path d=\"M 340 192\nL 392 192\nL 392 233\nL 340 233\nL 340 192\" style=\"stroke:none;fill:rgb(238,102,102)\"/><path d=\"M 498 138\nL 498 165\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:none\"/><path d=\"M 498 205\nL 498 246\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:none\"/><path d=\"M 485 138\nL 511 138\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:none\"/><path d=\"M 485 246\nL 511 246\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:none\"/><path d=\"M 472 165\nL 524 165\nL 524 205\nL 472 205\nL 472 165\" style=\"stroke:none;fill:rgb(238,102,102)\"/><path d=\"M 102 295\nL 168 235\nL 234 325\nL 300 176\nL 366 265\nL 432 116\nL 498 205\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"102\" cy=\"295\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"168\" cy=\"235\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"234\" cy=\"325\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"300\" cy=\"176\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"366\" cy=\"265\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"432\" cy=\"116\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"498\" cy=\"205\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/></svg>",
		},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			opt := tc.optFunc()
			opt.OutputFormat = ChartOutputSVG

			p, err := Render(opt)
			require.NoError(t, err)
			data, err := p.Bytes()
			require.NoError(t, err)
			assertEqualSVG(t, tc.expectedSVG, data)
		})
	}
}
