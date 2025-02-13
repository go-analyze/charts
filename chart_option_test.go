package charts

import (
	"fmt"
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
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 21 29\nL 51 29\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"36\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"53\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path  d=\"M 112 29\nL 142 29\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"127\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"144\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path  d=\"M 235 29\nL 265 29\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"250\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"267\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><path  d=\"M 358 29\nL 388 29\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"373\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"390\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path  d=\"M 451 29\nL 481 29\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:none\"/><circle cx=\"466\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:rgb(115,192,222)\"/><text x=\"483\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><text x=\"20\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line</text><text x=\"20\" y=\"62\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1440</text><text x=\"20\" y=\"94\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1280</text><text x=\"20\" y=\"127\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1120</text><text x=\"28\" y=\"160\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"28\" y=\"193\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"28\" y=\"225\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"28\" y=\"258\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"28\" y=\"291\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"28\" y=\"324\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"46\" y=\"357\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 65 55\nL 580 55\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 65 87\nL 580 87\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 65 120\nL 580 120\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 65 153\nL 580 153\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 65 186\nL 580 186\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 65 218\nL 580 218\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 65 251\nL 580 251\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 65 284\nL 580 284\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 65 317\nL 580 317\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 65 355\nL 65 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 138 355\nL 138 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 212 355\nL 212 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 285 355\nL 285 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 359 355\nL 359 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 432 355\nL 432 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 506 355\nL 506 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 580 355\nL 580 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 65 350\nL 580 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"86\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"162\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"233\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"309\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"386\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"458\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"530\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path  d=\"M 101 326\nL 175 323\nL 248 330\nL 322 323\nL 395 332\nL 469 303\nL 543 307\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"101\" cy=\"326\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"175\" cy=\"323\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"248\" cy=\"330\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"322\" cy=\"323\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"395\" cy=\"332\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"469\" cy=\"303\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"543\" cy=\"307\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><path  d=\"M 101 305\nL 175 313\nL 248 311\nL 322 303\nL 395 291\nL 469 283\nL 543 287\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"101\" cy=\"305\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"175\" cy=\"313\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"248\" cy=\"311\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"322\" cy=\"303\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"395\" cy=\"291\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"469\" cy=\"283\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"543\" cy=\"287\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><path  d=\"M 101 320\nL 175 303\nL 248 309\nL 322 319\nL 395 312\nL 469 283\nL 543 267\" style=\"stroke-width:2;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"101\" cy=\"320\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"175\" cy=\"303\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"248\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"322\" cy=\"319\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"395\" cy=\"312\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"469\" cy=\"283\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"543\" cy=\"267\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><path  d=\"M 101 285\nL 175 282\nL 248 289\nL 322 282\nL 395 271\nL 469 283\nL 543 285\" style=\"stroke-width:2;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"101\" cy=\"285\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:white\"/><circle cx=\"175\" cy=\"282\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:white\"/><circle cx=\"248\" cy=\"289\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:white\"/><circle cx=\"322\" cy=\"282\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:white\"/><circle cx=\"395\" cy=\"271\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:white\"/><circle cx=\"469\" cy=\"283\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:white\"/><circle cx=\"543\" cy=\"285\" r=\"2\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:white\"/><path  d=\"M 101 183\nL 175 160\nL 248 166\nL 322 159\nL 395 86\nL 469 78\nL 543 80\" style=\"stroke-width:2;stroke:rgb(115,192,222);fill:none\"/><circle cx=\"101\" cy=\"183\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:white\"/><circle cx=\"175\" cy=\"160\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:white\"/><circle cx=\"248\" cy=\"166\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:white\"/><circle cx=\"322\" cy=\"159\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:white\"/><circle cx=\"395\" cy=\"86\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:white\"/><circle cx=\"469\" cy=\"78\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:white\"/><circle cx=\"543\" cy=\"80\" r=\"2\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:white\"/></svg>", data)
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
		MarkLineOptionFunc(0, SeriesMarkDataTypeAverage),
		MarkPointOptionFunc(0, SeriesMarkDataTypeMax, SeriesMarkDataTypeMin),
		// custom option func
		func(opt *ChartOption) {
			opt.Legend.Offset = OffsetRight
			opt.Legend.OverlayChart = Ptr(true)
			opt.SeriesList[1].MarkPoint = NewMarkPoint(
				SeriesMarkDataTypeMax,
				SeriesMarkDataTypeMin,
			)
			opt.SeriesList[1].MarkLine = NewMarkLine(
				SeriesMarkDataTypeAverage,
			)
		},
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 364 29\nL 394 29\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"379\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"396\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Rainfall</text><path  d=\"M 468 29\nL 498 29\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"483\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"500\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Evaporation</text><text x=\"20\" y=\"27\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">270</text><text x=\"20\" y=\"63\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">240</text><text x=\"20\" y=\"100\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">210</text><text x=\"20\" y=\"137\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">180</text><text x=\"20\" y=\"173\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">150</text><text x=\"20\" y=\"210\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">120</text><text x=\"29\" y=\"247\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">90</text><text x=\"29\" y=\"283\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">60</text><text x=\"29\" y=\"320\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">30</text><text x=\"38\" y=\"357\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 57 20\nL 580 20\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 57 56\nL 580 56\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 57 93\nL 580 93\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 57 130\nL 580 130\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 57 166\nL 580 166\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 57 203\nL 580 203\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 57 240\nL 580 240\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 57 276\nL 580 276\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 57 313\nL 580 313\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 57 355\nL 57 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 100 355\nL 100 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 144 355\nL 144 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 187 355\nL 187 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 231 355\nL 231 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 274 355\nL 274 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 318 355\nL 318 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 362 355\nL 362 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 405 355\nL 405 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 449 355\nL 449 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 492 355\nL 492 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 536 355\nL 536 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 580 355\nL 580 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 57 350\nL 580 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"65\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jan</text><text x=\"109\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Feb</text><text x=\"151\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mar</text><text x=\"197\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Apr</text><text x=\"237\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">May</text><text x=\"283\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jun</text><text x=\"330\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jul</text><text x=\"369\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Aug</text><text x=\"414\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sep</text><text x=\"458\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Oct</text><text x=\"500\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Nov</text><text x=\"545\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Dec</text><path  d=\"M 62 348\nL 77 348\nL 77 349\nL 62 349\nL 62 348\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 105 345\nL 120 345\nL 120 349\nL 105 349\nL 105 345\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 149 342\nL 164 342\nL 164 349\nL 149 349\nL 149 342\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 192 322\nL 207 322\nL 207 349\nL 192 349\nL 192 322\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 236 319\nL 251 319\nL 251 349\nL 236 349\nL 236 319\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 279 257\nL 294 257\nL 294 349\nL 279 349\nL 279 257\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 323 185\nL 338 185\nL 338 349\nL 323 349\nL 323 185\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 367 152\nL 382 152\nL 382 349\nL 367 349\nL 367 152\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 410 311\nL 425 311\nL 425 349\nL 410 349\nL 410 311\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 454 326\nL 469 326\nL 469 349\nL 454 349\nL 454 326\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 497 343\nL 512 343\nL 512 349\nL 497 349\nL 497 343\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 541 346\nL 556 346\nL 556 349\nL 541 349\nL 541 346\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 80 347\nL 95 347\nL 95 349\nL 80 349\nL 80 347\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 123 343\nL 138 343\nL 138 349\nL 123 349\nL 123 343\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 167 339\nL 182 339\nL 182 349\nL 167 349\nL 167 339\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 210 318\nL 225 318\nL 225 349\nL 210 349\nL 210 318\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 254 315\nL 269 315\nL 269 349\nL 254 349\nL 254 315\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 297 264\nL 312 264\nL 312 349\nL 297 349\nL 297 264\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 341 136\nL 356 136\nL 356 349\nL 341 349\nL 341 136\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 385 128\nL 400 128\nL 400 349\nL 385 349\nL 385 128\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 428 291\nL 443 291\nL 443 349\nL 428 349\nL 428 291\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 472 328\nL 487 328\nL 487 349\nL 472 349\nL 472 328\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 515 343\nL 530 343\nL 530 349\nL 515 349\nL 515 343\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 559 348\nL 574 348\nL 574 349\nL 559 349\nL 559 348\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 371 144\nA 14 14 330.00 1 1 377 144\nL 374 131\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 360 131\nQ374,166 388,131\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><text x=\"361\" y=\"136\" style=\"stroke:none;fill:rgb(238,238,238);font-size:10.2px;font-family:'Roboto Medium',sans-serif\">162.2</text><path  d=\"M 66 340\nA 14 14 330.00 1 1 72 340\nL 69 327\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 55 327\nQ69,362 83,327\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><text x=\"65\" y=\"332\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2</text><path  d=\"M 389 120\nA 14 14 330.00 1 1 395 120\nL 392 107\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 378 107\nQ392,142 406,107\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"379\" y=\"112\" style=\"stroke:none;fill:rgb(70,70,70);font-size:10.2px;font-family:'Roboto Medium',sans-serif\">182.2</text><path  d=\"M 563 340\nA 14 14 330.00 1 1 569 340\nL 566 327\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 552 327\nQ566,362 580,327\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"557\" y=\"332\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2.3</text><circle cx=\"60\" cy=\"300\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 66 300\nL 562 300\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 562 295\nL 578 300\nL 562 305\nL 567 300\nL 562 295\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"580\" y=\"304\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">41.63</text><circle cx=\"60\" cy=\"292\" r=\"3\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 66 292\nL 562 292\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 562 287\nL 578 292\nL 562 297\nL 567 292\nL 562 287\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"580\" y=\"296\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">48.07</text></svg>", data)
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
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 214 29\nL 244 29\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"229\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"246\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2011</text><path  d=\"M 301 29\nL 331 29\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"316\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"333\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2012</text><text x=\"20\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World Population</text><path  d=\"M 93 55\nL 98 55\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 93 104\nL 98 104\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 93 153\nL 98 153\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 93 202\nL 98 202\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 93 251\nL 98 251\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 93 300\nL 98 300\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 93 350\nL 98 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 98 55\nL 98 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"47\" y=\"86\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World</text><text x=\"48\" y=\"135\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">China</text><text x=\"54\" y=\"184\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">India</text><text x=\"58\" y=\"233\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">USA</text><text x=\"20\" y=\"282\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Indonesia</text><text x=\"49\" y=\"332\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Brazil</text><text x=\"97\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"189\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">143k</text><text x=\"281\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">286k</text><text x=\"374\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">428.99k</text><text x=\"466\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">571.99k</text><text x=\"504\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">714.99k</text><path  d=\"M 190 55\nL 190 350\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 282 55\nL 282 350\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 375 55\nL 375 350\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 467 55\nL 467 350\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 560 55\nL 560 350\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 98 305\nL 109 305\nL 109 323\nL 98 323\nL 98 305\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 98 256\nL 112 256\nL 112 274\nL 98 274\nL 98 256\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 98 207\nL 115 207\nL 115 225\nL 98 225\nL 98 207\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 98 158\nL 162 158\nL 162 176\nL 98 176\nL 98 158\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 98 109\nL 179 109\nL 179 127\nL 98 127\nL 98 109\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 98 60\nL 486 60\nL 486 78\nL 98 78\nL 98 60\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 98 326\nL 109 326\nL 109 344\nL 98 344\nL 98 326\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 98 277\nL 112 277\nL 112 295\nL 98 295\nL 98 277\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 98 228\nL 117 228\nL 117 246\nL 98 246\nL 98 228\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 98 179\nL 172 179\nL 172 197\nL 98 197\nL 98 179\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 98 130\nL 180 130\nL 180 148\nL 98 148\nL 98 130\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 98 81\nL 517 81\nL 517 99\nL 98 99\nL 98 81\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>", data)
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
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 20 29\nL 50 29\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"35\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"52\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><path  d=\"M 20 49\nL 50 49\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"35\" cy=\"49\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"52\" y=\"55\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path  d=\"M 20 69\nL 50 69\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"35\" cy=\"69\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"52\" y=\"75\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path  d=\"M 20 89\nL 50 89\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"35\" cy=\"89\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"52\" y=\"95\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path  d=\"M 20 109\nL 50 109\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:none\"/><circle cx=\"35\" cy=\"109\" r=\"5\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:rgb(115,192,222)\"/><text x=\"52\" y=\"115\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><text x=\"222\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Rainfall vs Evaporation</text><text x=\"266\" y=\"50\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fake Data</text><path  d=\"M 300 207\nL 300 93\nA 114 114 119.89 0 1 398 263\nL 300 207\nZ\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path  d=\"M 398 150\nL 411 143\nM 411 143\nL 426 143\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:none\"/><text x=\"429\" y=\"148\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">: 33.3%</text><path  d=\"M 300 207\nL 398 263\nA 114 114 84.08 0 1 254 311\nL 300 207\nZ\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path  d=\"M 335 315\nL 340 329\nM 340 329\nL 355 329\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:none\"/><text x=\"358\" y=\"334\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">: 23.35%</text><path  d=\"M 300 207\nL 254 311\nA 114 114 66.35 0 1 187 207\nL 300 207\nZ\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><path  d=\"M 205 268\nL 192 276\nM 192 276\nL 177 276\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:none\"/><text x=\"125\" y=\"281\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">: 18.43%</text><path  d=\"M 300 207\nL 187 207\nA 114 114 55.37 0 1 236 113\nL 300 207\nZ\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><path  d=\"M 200 154\nL 187 147\nM 187 147\nL 172 147\" style=\"stroke-width:1;stroke:rgb(238,102,102);fill:none\"/><text x=\"120\" y=\"152\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">: 15.37%</text><path  d=\"M 300 207\nL 236 113\nA 114 114 34.32 0 1 300 93\nL 300 207\nZ\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:rgb(115,192,222)\"/><path  d=\"M 267 99\nL 262 84\nM 262 84\nL 247 84\" style=\"stroke-width:1;stroke:rgb(115,192,222);fill:none\"/><text x=\"202\" y=\"89\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">: 9.53%</text></svg>", data)
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
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 144 29\nL 174 29\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"159\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"176\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Allocated Budget</text><path  d=\"M 314 29\nL 344 29\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"329\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"346\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Actual Spending</text><text x=\"20\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Basic Radar Chart</text><path  d=\"M 300 179\nL 319 191\nL 319 213\nL 300 225\nL 281 213\nL 281 191\nL 300 179\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 300 156\nL 339 179\nL 339 224\nL 300 248\nL 261 225\nL 261 180\nL 300 156\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 300 133\nL 359 168\nL 359 236\nL 300 271\nL 241 236\nL 241 168\nL 300 133\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 300 110\nL 379 156\nL 379 247\nL 300 294\nL 221 248\nL 221 157\nL 300 110\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 300 87\nL 399 145\nL 399 259\nL 300 317\nL 201 259\nL 201 145\nL 300 87\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 300 202\nL 300 87\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 300 202\nL 399 145\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 300 202\nL 399 259\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 300 202\nL 300 317\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 300 202\nL 201 259\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 300 202\nL 201 145\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><text x=\"284\" y=\"80\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Sales</text><text x=\"404\" y=\"150\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Administration</text><text x=\"404\" y=\"264\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Information Technology</text><text x=\"248\" y=\"334\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Customer Support</text><text x=\"120\" y=\"264\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Development</text><text x=\"137\" y=\"150\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Marketing</text><path  d=\"M 300 128\nL 318 192\nL 366 240\nL 300 307\nL 205 257\nL 229 161\nL 300 128\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><path  d=\"M 300 128\nL 318 192\nL 366 240\nL 300 307\nL 205 257\nL 229 161\nL 300 128\" style=\"stroke:none;fill:rgba(84,112,198,0.1)\"/><circle cx=\"300\" cy=\"128\" r=\"2\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"318\" cy=\"192\" r=\"2\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"366\" cy=\"240\" r=\"2\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"300\" cy=\"307\" r=\"2\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"205\" cy=\"257\" r=\"2\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"229\" cy=\"161\" r=\"2\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"300\" cy=\"128\" r=\"2\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:white\"/><path  d=\"M 300 114\nL 387 152\nL 392 255\nL 300 280\nL 220 248\nL 217 154\nL 300 114\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/><path  d=\"M 300 114\nL 387 152\nL 392 255\nL 300 280\nL 220 248\nL 217 154\nL 300 114\" style=\"stroke:none;fill:rgba(145,204,117,0.1)\"/><circle cx=\"300\" cy=\"114\" r=\"2\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"387\" cy=\"152\" r=\"2\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"392\" cy=\"255\" r=\"2\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"300\" cy=\"280\" r=\"2\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"220\" cy=\"248\" r=\"2\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"217\" cy=\"154\" r=\"2\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"300\" cy=\"114\" r=\"2\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:white\"/></svg>", data)
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
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 87 29\nL 117 29\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"102\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"119\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Show</text><path  d=\"M 177 29\nL 207 29\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"192\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"209\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Click</text><path  d=\"M 263 29\nL 293 29\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"278\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(250,200,88);fill:rgb(250,200,88)\"/><text x=\"295\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Visit</text><path  d=\"M 346 29\nL 376 29\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:none\"/><circle cx=\"361\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(238,102,102);fill:rgb(238,102,102)\"/><text x=\"378\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Inquiry</text><path  d=\"M 445 29\nL 475 29\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:none\"/><circle cx=\"460\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(115,192,222);fill:rgb(115,192,222)\"/><text x=\"477\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Order</text><text x=\"20\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Funnel</text><path  d=\"M 20 55\nL 580 55\nL 524 112\nL 76 112\nL 20 55\" style=\"stroke:none;fill:rgb(84,112,198)\"/><text x=\"280\" y=\"83\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">(100%)</text><path  d=\"M 76 114\nL 524 114\nL 468 171\nL 132 171\nL 76 114\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"284\" y=\"142\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">(80%)</text><path  d=\"M 132 173\nL 468 173\nL 412 230\nL 188 230\nL 132 173\" style=\"stroke:none;fill:rgb(250,200,88)\"/><text x=\"284\" y=\"201\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">(60%)</text><path  d=\"M 188 232\nL 412 232\nL 356 289\nL 244 289\nL 188 232\" style=\"stroke:none;fill:rgb(238,102,102)\"/><text x=\"284\" y=\"260\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">(40%)</text><path  d=\"M 244 291\nL 356 291\nL 300 348\nL 300 348\nL 244 291\" style=\"stroke:none;fill:rgb(115,192,222)\"/><text x=\"284\" y=\"319\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">(20%)</text></svg>", data)
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
			Box: NewBox(10, 200, 200, 500),
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
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"20\" y=\"27\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">450</text><text x=\"20\" y=\"63\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">410</text><text x=\"20\" y=\"100\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">370</text><text x=\"20\" y=\"137\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">330</text><text x=\"20\" y=\"173\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">290</text><text x=\"20\" y=\"210\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">250</text><text x=\"20\" y=\"247\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">210</text><text x=\"20\" y=\"283\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">170</text><text x=\"20\" y=\"320\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">130</text><text x=\"29\" y=\"357\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">90</text><path  d=\"M 57 20\nL 580 20\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 57 56\nL 580 56\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 57 93\nL 580 93\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 57 130\nL 580 130\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 57 166\nL 580 166\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 57 203\nL 580 203\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 57 240\nL 580 240\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 57 276\nL 580 276\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 57 313\nL 580 313\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 57 355\nL 57 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 131 355\nL 131 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 206 355\nL 206 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 281 355\nL 281 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 355 355\nL 355 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 430 355\nL 430 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 505 355\nL 505 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 580 355\nL 580 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 57 350\nL 580 350\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"79\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"155\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"228\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"305\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"383\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"456\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"529\" y=\"375\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path  d=\"M 94 323\nL 168 312\nL 243 340\nL 318 310\nL 392 350\nL 467 222\nL 542 240\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"94\" cy=\"323\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"168\" cy=\"312\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"243\" cy=\"340\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"318\" cy=\"310\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"392\" cy=\"350\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"467\" cy=\"222\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"542\" cy=\"240\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><path  d=\"M 94 295\nL 168 220\nL 243 249\nL 318 292\nL 392 259\nL 467 130\nL 542 57\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"94\" cy=\"295\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"168\" cy=\"220\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"243\" cy=\"249\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"318\" cy=\"292\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"392\" cy=\"259\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"467\" cy=\"130\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"542\" cy=\"57\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><path  d=\"M 94 140\nL 168 129\nL 243 157\nL 318 127\nL 392 75\nL 467 130\nL 542 140\" style=\"stroke-width:2;stroke:rgb(250,200,88);fill:none\"/><circle cx=\"94\" cy=\"140\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"168\" cy=\"129\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"243\" cy=\"157\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"318\" cy=\"127\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"392\" cy=\"75\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"467\" cy=\"130\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><circle cx=\"542\" cy=\"140\" r=\"2\" style=\"stroke-width:1;stroke:rgb(250,200,88);fill:white\"/><path  d=\"M 274 39\nL 304 39\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"289\" cy=\"39\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"306\" y=\"45\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2011</text><path  d=\"M 361 39\nL 391 39\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"376\" cy=\"39\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"393\" y=\"45\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2012</text><path  d=\"M 266 65\nL 271 65\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 266 86\nL 271 86\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 266 107\nL 271 107\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 266 128\nL 271 128\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 266 150\nL 271 150\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 271 65\nL 271 150\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"220\" y=\"82\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World</text><text x=\"221\" y=\"103\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">China</text><text x=\"227\" y=\"124\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">India</text><text x=\"231\" y=\"146\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">USA</text><text x=\"270\" y=\"175\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">70</text><text x=\"453\" y=\"175\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">144</text><path  d=\"M 340 65\nL 340 150\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 410 65\nL 410 150\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 480 65\nL 480 150\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 271 133\nL 271 133\nL 271 137\nL 271 137\nL 271 133\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 271 112\nL 329 112\nL 329 116\nL 271 116\nL 271 112\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 271 91\nL 387 91\nL 387 95\nL 271 95\nL 271 91\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 271 70\nL 445 70\nL 445 74\nL 271 74\nL 271 70\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 271 140\nL 300 140\nL 300 144\nL 271 144\nL 271 140\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 271 119\nL 358 119\nL 358 123\nL 271 123\nL 271 119\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 271 98\nL 416 98\nL 416 102\nL 271 102\nL 271 98\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 271 77\nL 474 77\nL 474 81\nL 271 81\nL 271 77\" style=\"stroke:none;fill:rgb(145,204,117)\"/></svg>", data)
}
