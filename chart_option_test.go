package charts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChartOption(t *testing.T) {
	t.Parallel()

	fns := []OptionFunc{
		SVGOutputOption(),
		FontOptionFunc(GetDefaultFont()),
		ThemeNameOptionFunc(ThemeVividDark),
		TitleTextOptionFunc("title"),
		LegendLabelsOptionFunc([]string{"label"}),
		XAxisDataOptionFunc([]string{"xaxis"}),
		YAxisDataOptionFunc([]string{"yaxis"}),
		WidthOptionFunc(800),
		HeightOptionFunc(600),
		PaddingOptionFunc(Box{
			Left:   10,
			Top:    10,
			Right:  10,
			Bottom: 10,
		}),
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
			Data: []string{"label"},
		},
		XAxis: XAxisOption{
			Data: []string{"xaxis"},
		},
		YAxis: []YAxisOption{
			{
				Data: []string{"yaxis"},
			},
		},
		Width:  800,
		Height: 600,
		Padding: Box{
			Left:   10,
			Top:    10,
			Right:  10,
			Bottom: 10,
		},
	}, opt)
}

func TestChartOptionPieSeriesShowLabel(t *testing.T) {
	t.Parallel()

	opt := ChartOption{
		SeriesList: NewPieSeriesList([]float64{1, 2}),
	}
	PieSeriesShowLabel()(&opt)
	assert.True(t, opt.SeriesList[0].Label.Show)
}

func TestChartOptionMarkLine(t *testing.T) {
	t.Parallel()

	opt := ChartOption{
		SeriesList: NewSeriesListDataFromValues([][]float64{{1, 2}}),
	}
	MarkLineOptionFunc(0, "min", "max")(&opt)
	assert.Equal(t, NewMarkLine("min", "max"), opt.SeriesList[0].MarkLine)
}

func TestChartOptionMarkPoint(t *testing.T) {
	t.Parallel()

	opt := ChartOption{
		SeriesList: NewSeriesListDataFromValues([][]float64{{1, 2}}),
	}
	MarkPointOptionFunc(0, "min", "max")(&opt)
	assert.Equal(t, NewMarkPoint("min", "max"), opt.SeriesList[0].MarkPoint)
}

func TestLineRender(t *testing.T) {
	t.Parallel()

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
	p, err := LineRender(
		values,
		SVGOutputOption(),
		TitleTextOptionFunc("Line"),
		XAxisDataOptionFunc([]string{
			"Mon",
			"Tue",
			"Wed",
			"Thu",
			"Fri",
			"Sat",
			"Sun",
		}),
		LegendLabelsOptionFunc([]string{
			"Email",
			"Union Ads",
			"Video Ads",
			"Direct",
			"Search Engine",
		}, PositionCenter),
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 20 29\nL 50 29\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><circle cx=\"35\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><text x=\"52\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path  d=\"M 111 29\nL 141 29\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><circle cx=\"126\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><text x=\"143\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path  d=\"M 234 29\nL 264 29\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><circle cx=\"249\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><text x=\"266\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><path  d=\"M 357 29\nL 387 29\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><circle cx=\"372\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><text x=\"389\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path  d=\"M 450 29\nL 480 29\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><circle cx=\"465\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><text x=\"482\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><text x=\"20\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Line</text><text x=\"20\" y=\"62\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"20\" y=\"94\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.28k</text><text x=\"20\" y=\"127\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.12k</text><text x=\"32\" y=\"160\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"32\" y=\"193\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"32\" y=\"225\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"32\" y=\"258\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"32\" y=\"291\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"32\" y=\"324\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"50\" y=\"357\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 69 55\nL 580 55\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 69 87\nL 580 87\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 69 120\nL 580 120\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 69 153\nL 580 153\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 69 186\nL 580 186\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 69 218\nL 580 218\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 69 251\nL 580 251\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 69 284\nL 580 284\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 69 317\nL 580 317\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 69 355\nL 69 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 142 355\nL 142 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 215 355\nL 215 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 288 355\nL 288 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 361 355\nL 361 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 434 355\nL 434 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 507 355\nL 507 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 580 355\nL 580 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 69 350\nL 580 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"90\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"165\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"236\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"311\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"388\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"459\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"530\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><path  d=\"M 105 326\nL 178 323\nL 251 330\nL 324 323\nL 397 332\nL 470 303\nL 543 307\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:none\"/><circle cx=\"105\" cy=\"326\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"178\" cy=\"323\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"251\" cy=\"330\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"323\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"397\" cy=\"332\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"470\" cy=\"303\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"543\" cy=\"307\" r=\"2\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 105 305\nL 178 313\nL 251 311\nL 324 303\nL 397 291\nL 470 283\nL 543 287\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:none\"/><circle cx=\"105\" cy=\"305\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"178\" cy=\"313\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"251\" cy=\"311\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"303\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"397\" cy=\"291\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"470\" cy=\"283\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"543\" cy=\"287\" r=\"2\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 105 320\nL 178 303\nL 251 309\nL 324 319\nL 397 312\nL 470 283\nL 543 267\" style=\"stroke-width:2;stroke:rgba(250,200,88,1.0);fill:none\"/><circle cx=\"105\" cy=\"320\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"178\" cy=\"303\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"251\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"319\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"397\" cy=\"312\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"470\" cy=\"283\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"543\" cy=\"267\" r=\"2\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 105 285\nL 178 282\nL 251 289\nL 324 282\nL 397 271\nL 470 283\nL 543 285\" style=\"stroke-width:2;stroke:rgba(238,102,102,1.0);fill:none\"/><circle cx=\"105\" cy=\"285\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"178\" cy=\"282\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"251\" cy=\"289\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"282\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"397\" cy=\"271\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"470\" cy=\"283\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"543\" cy=\"285\" r=\"2\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 105 183\nL 178 160\nL 251 166\nL 324 159\nL 397 86\nL 470 78\nL 543 80\" style=\"stroke-width:2;stroke:rgba(115,192,222,1.0);fill:none\"/><circle cx=\"105\" cy=\"183\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"178\" cy=\"160\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"251\" cy=\"166\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"324\" cy=\"159\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"397\" cy=\"86\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"470\" cy=\"78\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"543\" cy=\"80\" r=\"2\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(255,255,255,1.0)\"/></svg>", string(data))
}

func TestBarRender(t *testing.T) {
	t.Parallel()

	values := [][]float64{
		{
			2.0,
			4.9,
			7.0,
			23.2,
			25.6,
			76.7,
			135.6,
			162.2,
			32.6,
			20.0,
			6.4,
			3.3,
		},
		{
			2.6,
			5.9,
			9.0,
			26.4,
			28.7,
			70.7,
			175.6,
			182.2,
			48.7,
			18.8,
			6.0,
			2.3,
		},
	}
	p, err := BarRender(
		values,
		SVGOutputOption(),
		XAxisDataOptionFunc([]string{
			"Jan",
			"Feb",
			"Mar",
			"Apr",
			"May",
			"Jun",
			"Jul",
			"Aug",
			"Sep",
			"Oct",
			"Nov",
			"Dec",
		}),
		LegendLabelsOptionFunc([]string{
			"Rainfall",
			"Evaporation",
		}, PositionRight),
		MarkLineOptionFunc(0, SeriesMarkDataTypeAverage),
		MarkPointOptionFunc(0, SeriesMarkDataTypeMax, SeriesMarkDataTypeMin),
		// custom option func
		func(opt *ChartOption) {
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
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 364 29\nL 394 29\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><circle cx=\"379\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><text x=\"396\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Rainfall</text><path  d=\"M 468 29\nL 498 29\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><circle cx=\"483\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><text x=\"500\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Evaporation</text><text x=\"20\" y=\"27\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">270</text><text x=\"20\" y=\"63\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">240</text><text x=\"20\" y=\"100\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">210</text><text x=\"20\" y=\"137\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">180</text><text x=\"20\" y=\"173\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">150</text><text x=\"20\" y=\"210\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">120</text><text x=\"29\" y=\"247\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">90</text><text x=\"29\" y=\"283\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">60</text><text x=\"29\" y=\"320\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">30</text><text x=\"38\" y=\"357\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 57 20\nL 580 20\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 57 56\nL 580 56\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 57 93\nL 580 93\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 57 130\nL 580 130\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 57 166\nL 580 166\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 57 203\nL 580 203\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 57 240\nL 580 240\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 57 276\nL 580 276\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 57 313\nL 580 313\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 57 355\nL 57 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 100 355\nL 100 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 144 355\nL 144 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 187 355\nL 187 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 231 355\nL 231 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 274 355\nL 274 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 318 355\nL 318 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 362 355\nL 362 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 405 355\nL 405 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 449 355\nL 449 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 492 355\nL 492 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 536 355\nL 536 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 580 355\nL 580 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 57 350\nL 580 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"65\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jan</text><text x=\"109\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Feb</text><text x=\"151\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mar</text><text x=\"197\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Apr</text><text x=\"237\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">May</text><text x=\"283\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jun</text><text x=\"330\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jul</text><text x=\"369\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Aug</text><text x=\"414\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sep</text><text x=\"458\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Oct</text><text x=\"500\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Nov</text><text x=\"545\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Dec</text><path  d=\"M 62 348\nL 77 348\nL 77 349\nL 62 349\nL 62 348\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 105 345\nL 120 345\nL 120 349\nL 105 349\nL 105 345\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 149 342\nL 164 342\nL 164 349\nL 149 349\nL 149 342\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 192 322\nL 207 322\nL 207 349\nL 192 349\nL 192 322\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 236 319\nL 251 319\nL 251 349\nL 236 349\nL 236 319\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 279 257\nL 294 257\nL 294 349\nL 279 349\nL 279 257\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 323 185\nL 338 185\nL 338 349\nL 323 349\nL 323 185\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 367 152\nL 382 152\nL 382 349\nL 367 349\nL 367 152\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 410 311\nL 425 311\nL 425 349\nL 410 349\nL 410 311\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 454 326\nL 469 326\nL 469 349\nL 454 349\nL 454 326\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 497 343\nL 512 343\nL 512 349\nL 497 349\nL 497 343\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 541 346\nL 556 346\nL 556 349\nL 541 349\nL 541 346\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 80 347\nL 95 347\nL 95 349\nL 80 349\nL 80 347\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 123 343\nL 138 343\nL 138 349\nL 123 349\nL 123 343\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 167 339\nL 182 339\nL 182 349\nL 167 349\nL 167 339\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 210 318\nL 225 318\nL 225 349\nL 210 349\nL 210 318\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 254 315\nL 269 315\nL 269 349\nL 254 349\nL 254 315\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 297 264\nL 312 264\nL 312 349\nL 297 349\nL 297 264\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 341 136\nL 356 136\nL 356 349\nL 341 349\nL 341 136\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 385 128\nL 400 128\nL 400 349\nL 385 349\nL 385 128\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 428 291\nL 443 291\nL 443 349\nL 428 349\nL 428 291\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 472 328\nL 487 328\nL 487 349\nL 472 349\nL 472 328\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 515 343\nL 530 343\nL 530 349\nL 515 349\nL 515 343\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 559 348\nL 574 348\nL 574 349\nL 559 349\nL 559 348\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 371 144\nA 14 14 330.00 1 1 377 144\nL 374 131\nZ\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 360 131\nQ374,166 388,131\nZ\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><text x=\"361\" y=\"136\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:10.2px;font-family:'Roboto Medium',sans-serif\">162.2</text><path  d=\"M 66 340\nA 14 14 330.00 1 1 72 340\nL 69 327\nZ\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 55 327\nQ69,362 83,327\nZ\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><text x=\"65\" y=\"332\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2</text><path  d=\"M 389 120\nA 14 14 330.00 1 1 395 120\nL 392 107\nZ\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 378 107\nQ392,142 406,107\nZ\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><text x=\"379\" y=\"112\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:10.2px;font-family:'Roboto Medium',sans-serif\">182.2</text><path  d=\"M 563 340\nA 14 14 330.00 1 1 569 340\nL 566 327\nZ\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 552 327\nQ566,362 580,327\nZ\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><text x=\"557\" y=\"332\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2.3</text><circle cx=\"60\" cy=\"300\" r=\"3\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 66 300\nL 562 300\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 562 295\nL 578 300\nL 562 305\nL 567 300\nL 562 295\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><text x=\"580\" y=\"304\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">41.62</text><circle cx=\"60\" cy=\"292\" r=\"3\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 66 292\nL 562 292\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 562 287\nL 578 292\nL 562 297\nL 567 292\nL 562 287\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><text x=\"580\" y=\"296\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">48.07</text></svg>", string(data))
}

func TestHorizontalBarRender(t *testing.T) {
	t.Parallel()

	values := [][]float64{
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
	}
	p, err := HorizontalBarRender(
		values,
		SVGOutputOption(),
		TitleTextOptionFunc("World Population"),
		PaddingOptionFunc(Box{
			Top:    20,
			Right:  40,
			Bottom: 20,
			Left:   20,
		}),
		LegendLabelsOptionFunc([]string{
			"2011",
			"2012",
		}),
		YAxisDataOptionFunc([]string{
			"Brazil",
			"Indonesia",
			"USA",
			"India",
			"China",
			"World",
		}),
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 214 29\nL 244 29\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><circle cx=\"229\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><text x=\"246\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2011</text><path  d=\"M 301 29\nL 331 29\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><circle cx=\"316\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><text x=\"333\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2012</text><text x=\"20\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World Population</text><path  d=\"M 93 55\nL 98 55\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 93 104\nL 98 104\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 93 153\nL 98 153\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 93 202\nL 98 202\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 93 251\nL 98 251\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 93 300\nL 98 300\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 93 350\nL 98 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 98 55\nL 98 350\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"47\" y=\"86\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">World</text><text x=\"48\" y=\"135\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">China</text><text x=\"54\" y=\"184\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">India</text><text x=\"58\" y=\"233\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">USA</text><text x=\"20\" y=\"282\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Indonesia</text><text x=\"49\" y=\"332\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Brazil</text><text x=\"97\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><text x=\"189\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">142.99k</text><text x=\"281\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">285.99k</text><text x=\"374\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">428.99k</text><text x=\"466\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">571.99k</text><text x=\"504\" y=\"375\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">714.98k</text><path  d=\"M 190 55\nL 190 350\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 282 55\nL 282 350\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 375 55\nL 375 350\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 467 55\nL 467 350\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 560 55\nL 560 350\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 98 305\nL 109 305\nL 109 323\nL 98 323\nL 98 305\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 98 256\nL 112 256\nL 112 274\nL 98 274\nL 98 256\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 98 207\nL 115 207\nL 115 225\nL 98 225\nL 98 207\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 98 158\nL 162 158\nL 162 176\nL 98 176\nL 98 158\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 98 109\nL 179 109\nL 179 127\nL 98 127\nL 98 109\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 98 60\nL 486 60\nL 486 78\nL 98 78\nL 98 60\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 98 326\nL 109 326\nL 109 344\nL 98 344\nL 98 326\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 98 277\nL 112 277\nL 112 295\nL 98 295\nL 98 277\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 98 228\nL 117 228\nL 117 246\nL 98 246\nL 98 228\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 98 179\nL 172 179\nL 172 197\nL 98 197\nL 98 179\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 98 130\nL 180 130\nL 180 148\nL 98 148\nL 98 130\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 98 81\nL 517 81\nL 517 99\nL 98 99\nL 98 81\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/></svg>", string(data))
}

func TestPieRender(t *testing.T) {
	t.Parallel()

	values := []float64{
		1048,
		735,
		580,
		484,
		300,
	}
	p, err := PieRender(
		values,
		SVGOutputOption(),
		TitleOptionFunc(TitleOption{
			Text:    "Rainfall vs Evaporation",
			Subtext: "Fake Data",
			Left:    PositionCenter,
		}),
		PaddingOptionFunc(Box{
			Top:    20,
			Right:  20,
			Bottom: 20,
			Left:   20,
		}),
		LegendOptionFunc(LegendOption{
			Orient: OrientVertical,
			Data: []string{
				"Search Engine",
				"Direct",
				"Email",
				"Union Ads",
				"Video Ads",
			},
			Left: PositionLeft,
		}),
		PieSeriesShowLabel(),
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 20 29\nL 50 29\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><circle cx=\"35\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><text x=\"52\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><path  d=\"M 20 49\nL 50 49\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><circle cx=\"35\" cy=\"49\" r=\"5\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><text x=\"52\" y=\"55\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path  d=\"M 20 69\nL 50 69\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><circle cx=\"35\" cy=\"69\" r=\"5\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><text x=\"52\" y=\"75\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path  d=\"M 20 89\nL 50 89\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><circle cx=\"35\" cy=\"89\" r=\"5\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><text x=\"52\" y=\"95\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path  d=\"M 20 109\nL 50 109\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><circle cx=\"35\" cy=\"109\" r=\"5\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><text x=\"52\" y=\"115\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><text x=\"222\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Rainfall vs Evaporation</text><text x=\"266\" y=\"50\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fake Data</text><path  d=\"M 300 210\nL 300 98\nA 112 112 119.89 0 1 397 265\nL 300 210\nZ\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path  d=\"M 396 154\nL 409 147\nM 409 147\nL 424 147\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><text x=\"427\" y=\"152\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Search Engine: 33.3%</text><path  d=\"M 300 210\nL 397 265\nA 112 112 84.08 0 1 255 312\nL 300 210\nZ\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path  d=\"M 334 316\nL 339 330\nM 339 330\nL 354 330\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><text x=\"357\" y=\"335\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Direct: 23.35%</text><path  d=\"M 300 210\nL 255 312\nA 112 112 66.35 0 1 189 210\nL 300 210\nZ\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><path  d=\"M 206 270\nL 194 278\nM 194 278\nL 179 278\" style=\"stroke-width:1;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><text x=\"95\" y=\"283\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Email: 18.43%</text><path  d=\"M 300 210\nL 189 210\nA 112 112 55.37 0 1 237 118\nL 300 210\nZ\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><path  d=\"M 202 158\nL 188 151\nM 188 151\nL 173 151\" style=\"stroke-width:1;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><text x=\"62\" y=\"156\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Union Ads: 15.37%</text><path  d=\"M 300 210\nL 237 118\nA 112 112 34.32 0 1 300 98\nL 300 210\nZ\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><path  d=\"M 267 103\nL 263 89\nM 263 89\nL 248 89\" style=\"stroke-width:1;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><text x=\"145\" y=\"94\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Video Ads: 9.53%</text></svg>", string(data))
}

func TestRadarRender(t *testing.T) {
	t.Parallel()

	values := [][]float64{
		{
			4200,
			3000,
			20000,
			35000,
			50000,
			18000,
		},
		{
			5000,
			14000,
			28000,
			26000,
			42000,
			21000,
		},
	}
	p, err := RadarRender(
		values,
		SVGOutputOption(),
		TitleTextOptionFunc("Basic Radar Chart"),
		LegendLabelsOptionFunc([]string{
			"Allocated Budget",
			"Actual Spending",
		}),
		RadarIndicatorOptionFunc([]string{
			"Sales",
			"Administration",
			"Information Technology",
			"Customer Support",
			"Development",
			"Marketing",
		}, []float64{
			6500,
			16000,
			30000,
			38000,
			52000,
			25000,
		}),
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 143 29\nL 173 29\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><circle cx=\"158\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><text x=\"175\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Allocated Budget</text><path  d=\"M 313 29\nL 343 29\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><circle cx=\"328\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><text x=\"345\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Actual Spending</text><text x=\"20\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Basic Radar Chart</text><path  d=\"M 300 179\nL 319 191\nL 319 213\nL 300 225\nL 281 213\nL 281 191\nL 300 179\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 300 156\nL 339 179\nL 339 224\nL 300 248\nL 261 225\nL 261 180\nL 300 156\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 300 133\nL 359 168\nL 359 236\nL 300 271\nL 241 236\nL 241 168\nL 300 133\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 300 110\nL 379 156\nL 379 247\nL 300 294\nL 221 248\nL 221 157\nL 300 110\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 300 87\nL 399 145\nL 399 259\nL 300 317\nL 201 259\nL 201 145\nL 300 87\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 300 202\nL 300 87\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 300 202\nL 399 145\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 300 202\nL 399 259\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 300 202\nL 300 317\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 300 202\nL 201 259\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 300 202\nL 201 145\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><text x=\"284\" y=\"80\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Sales</text><text x=\"404\" y=\"150\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Administration</text><text x=\"404\" y=\"264\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Information Technology</text><text x=\"248\" y=\"334\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Customer Support</text><text x=\"120\" y=\"264\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Development</text><text x=\"137\" y=\"150\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Marketing</text><path  d=\"M 300 128\nL 318 192\nL 366 240\nL 300 307\nL 205 257\nL 229 161\nL 300 128\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,0.1)\"/><path  d=\"M 300 128\nL 318 192\nL 366 240\nL 300 307\nL 205 257\nL 229 161\nL 300 128\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,0.1)\"/><circle cx=\"300\" cy=\"128\" r=\"2\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"318\" cy=\"192\" r=\"2\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"366\" cy=\"240\" r=\"2\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"300\" cy=\"307\" r=\"2\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"205\" cy=\"257\" r=\"2\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"229\" cy=\"161\" r=\"2\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"300\" cy=\"128\" r=\"2\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(84,112,198,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"M 300 114\nL 387 152\nL 392 255\nL 300 280\nL 220 248\nL 217 154\nL 300 114\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,0.1)\"/><path  d=\"M 300 114\nL 387 152\nL 392 255\nL 300 280\nL 220 248\nL 217 154\nL 300 114\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,0.1)\"/><circle cx=\"300\" cy=\"114\" r=\"2\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"387\" cy=\"152\" r=\"2\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"392\" cy=\"255\" r=\"2\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"300\" cy=\"280\" r=\"2\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"220\" cy=\"248\" r=\"2\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"217\" cy=\"154\" r=\"2\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><circle cx=\"300\" cy=\"114\" r=\"2\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/><path  d=\"\" style=\"stroke-width:2;stroke:rgba(145,204,117,1.0);fill:rgba(255,255,255,1.0)\"/></svg>", string(data))
}

func TestFunnelRender(t *testing.T) {
	t.Parallel()

	values := []float64{
		100,
		80,
		60,
		40,
		20,
	}
	p, err := FunnelRender(
		values,
		SVGOutputOption(),
		TitleTextOptionFunc("Funnel"),
		LegendLabelsOptionFunc([]string{
			"Show",
			"Click",
			"Visit",
			"Inquiry",
			"Order",
		}),
	)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 86 29\nL 116 29\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><circle cx=\"101\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><text x=\"118\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Show</text><path  d=\"M 176 29\nL 206 29\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><circle cx=\"191\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><text x=\"208\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Click</text><path  d=\"M 262 29\nL 292 29\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><circle cx=\"277\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(250,200,88,1.0);fill:rgba(250,200,88,1.0)\"/><text x=\"294\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Visit</text><path  d=\"M 345 29\nL 375 29\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><circle cx=\"360\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(238,102,102,1.0);fill:rgba(238,102,102,1.0)\"/><text x=\"377\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Inquiry</text><path  d=\"M 444 29\nL 474 29\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><circle cx=\"459\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(115,192,222,1.0);fill:rgba(115,192,222,1.0)\"/><text x=\"476\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Order</text><text x=\"20\" y=\"35\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Funnel</text><path  d=\"M 20 55\nL 580 55\nL 524 112\nL 76 112\nL 20 55\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><text x=\"264\" y=\"83\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Show(100%)</text><path  d=\"M 76 114\nL 524 114\nL 468 171\nL 132 171\nL 76 114\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><text x=\"269\" y=\"142\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Click(80%)</text><path  d=\"M 132 173\nL 468 173\nL 412 230\nL 188 230\nL 132 173\" style=\"stroke-width:0;stroke:none;fill:rgba(250,200,88,1.0)\"/><text x=\"271\" y=\"201\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Visit(60%)</text><path  d=\"M 188 232\nL 412 232\nL 356 289\nL 244 289\nL 188 232\" style=\"stroke-width:0;stroke:none;fill:rgba(238,102,102,1.0)\"/><text x=\"264\" y=\"260\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Inquiry(40%)</text><path  d=\"M 244 291\nL 356 291\nL 300 348\nL 300 348\nL 244 291\" style=\"stroke-width:0;stroke:none;fill:rgba(115,192,222,1.0)\"/><text x=\"268\" y=\"319\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">Order(20%)</text></svg>", string(data))
}
