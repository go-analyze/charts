package charts

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeFullScatterChartOption() ScatterChartOption {
	values := [][]float64{
		{120, 132, 101, 134, 90, 230, 210},
		{220, 182, 191, 234, 290, 330, 310},
		{150, 232, 201, 154, 190, 330, 410},
		{320, 332, 301, 334, 390, 330, 320},
		{820, 932, 901, 934, 1290, 1330, 1320},
	}
	return ScatterChartOption{
		Title: TitleOption{
			Text: "Scatter",
		},
		Padding: NewBoxEqual(10),
		XAxis: XAxisOption{
			Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
		},
		YAxis: make([]YAxisOption, 1),
		Legend: LegendOption{
			SeriesNames: []string{"Email", "Union Ads", "Video Ads", "Direct", "Search Engine"},
		},
		SeriesList: NewSeriesListScatter(values),
	}
}

func makeBasicScatterChartOption() ScatterChartOption {
	values := [][]float64{
		{120, 132, 101, 134, 90, 230, 210},
		{820, 932, 901, 934, 1290, 1330, 1320},
	}
	return ScatterChartOption{
		Padding: NewBoxEqual(10),
		XAxis: XAxisOption{
			Labels: []string{"A", "B", "C", "D", "E", "F", "G"},
		},
		YAxis: make([]YAxisOption, 1),
		Legend: LegendOption{
			SeriesNames: []string{"1", "2"},
		},
		SeriesList: NewSeriesListScatter(values),
	}
}

func makeMinimalScatterChartOption() ScatterChartOption {
	values := [][]float64{
		{120, 132, 101, 134, 90, 230, 210},
		{820, 932, 901, 934, 1290, 1330, 1320},
	}
	return ScatterChartOption{
		Padding: NewBoxEqual(10),
		XAxis: XAxisOption{
			Labels: []string{"1", "2", "3", "4", "5", "6", "7"},
			Show:   Ptr(false),
		},
		YAxis: []YAxisOption{
			{
				Show: Ptr(false),
			},
		},
		SeriesList: NewSeriesListScatter(values),
	}
}

func makeMinimalMultiValueScatterChartOption() ScatterChartOption {
	values := [][][]float64{
		{{120, GetNullValue()}, {132}, {101, 20}, {134}, {90, 28}, {230}, {210}},
		{{820, GetNullValue()}, {932}, {901, 600}, {934}, {1290}, {1330}, {1320}},
	}
	return ScatterChartOption{
		Padding: NewBoxEqual(10),
		XAxis: XAxisOption{
			Labels: []string{"1", "2", "3", "4", "5", "6", "7"},
			Show:   Ptr(false),
		},
		YAxis: []YAxisOption{
			{
				Show: Ptr(false),
			},
		},
		SeriesList: NewSeriesListScatterMultiValue(values),
	}
}

func generateRandomScatterData(seriesCount int, dataPointCount int, maxVariationPercentage float64) [][][]float64 {
	data := make([][][]float64, seriesCount)
	for i := 0; i < seriesCount; i++ {
		data[i] = make([][]float64, dataPointCount)
	}
	r := rand.New(rand.NewSource(1))

	for i := 0; i < seriesCount; i++ {
		for j := 0; j < dataPointCount; j++ {
			if j == 0 {
				// Set the initial value for the line
				data[i][j] = []float64{r.Float64() * 100}
			} else {
				// Calculate the allowed variation range
				variationRange := data[i][j-1][0] * maxVariationPercentage / 100
				min := data[i][j-1][0] - variationRange
				max := data[i][j-1][0] + variationRange

				// Generate a random value within the allowed range
				values := []float64{min + r.Float64()*(max-min)}
				if j%2 == 0 {
					values = append(values, min+r.Float64()*(max-min))
				}
				if j%10 == 0 {
					values = append(values, min+r.Float64()*(max-min))
				}
				data[i][j] = values
			}
		}
	}

	return data
}

func makeDenseScatterChartOption() ScatterChartOption {
	const dataPointCount = 100
	values := generateRandomScatterData(3, dataPointCount, 10)

	xAxisLabels := make([]string, dataPointCount)
	for i := 0; i < dataPointCount; i++ {
		xAxisLabels[i] = strconv.Itoa(i)
	}

	return ScatterChartOption{
		SeriesList: NewSeriesListScatterMultiValue(values, ScatterSeriesOption{
			TrendLine: NewTrendLine(SeriesMarkTypeAverage),
			Label: SeriesLabel{
				ValueFormatter: func(f float64) string {
					return FormatValueHumanizeShort(f, 0, false)
				},
			},
		}),
		Padding: NewBoxEqual(20),
		Theme:   GetTheme(ThemeLight),
		YAxis: []YAxisOption{
			{
				Min:            Ptr(0.0), // force min to be zero
				Max:            Ptr(200.0),
				Unit:           10,
				LabelSkipCount: 1,
			},
		},
		XAxis: XAxisOption{
			Labels:        xAxisLabels,
			BoundaryGap:   Ptr(false),
			LabelCount:    10,
			LabelRotation: DegreesToRadians(45),
		},
	}
}

func TestNewScatterChartOptionWithData(t *testing.T) {
	t.Parallel()

	opt := NewScatterChartOptionWithData([][]float64{
		{12, 24},
		{24, 48},
	})

	assert.Len(t, opt.SeriesList, 2)
	assert.Equal(t, ChartTypeScatter, opt.SeriesList[0].getType())
	assert.Len(t, opt.YAxis, 1)
	assert.Equal(t, defaultPadding, opt.Padding)

	p := NewPainter(PainterOptions{})
	assert.NoError(t, p.ScatterChart(opt))
}

func TestScatterChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		ignore      string // specified if the test is ignored
		themed      bool
		makeOptions func() ScatterChartOption
		svg         string
		pngCRC      uint32
	}{
		{
			name:        "basic_themed",
			themed:      true,
			makeOptions: makeFullScatterChartOption,
			svg:         "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:rgb(40,40,40)\"/><text x=\"10\" y=\"26\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Scatter</text><path d=\"M 21 35\nL 51 35\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:none\"/><circle cx=\"36\" cy=\"35\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><text x=\"53\" y=\"41\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path d=\"M 112 35\nL 142 35\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:none\"/><circle cx=\"127\" cy=\"35\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><text x=\"144\" y=\"41\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path d=\"M 235 35\nL 265 35\" style=\"stroke-width:3;stroke:rgb(100,180,210);fill:none\"/><circle cx=\"250\" cy=\"35\" r=\"5\" style=\"stroke-width:3;stroke:rgb(100,180,210);fill:rgb(100,180,210)\"/><text x=\"267\" y=\"41\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><path d=\"M 357 35\nL 387 35\" style=\"stroke-width:3;stroke:rgb(64,160,110);fill:none\"/><circle cx=\"372\" cy=\"35\" r=\"5\" style=\"stroke-width:3;stroke:rgb(64,160,110);fill:rgb(64,160,110)\"/><text x=\"389\" y=\"41\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path d=\"M 450 35\nL 480 35\" style=\"stroke-width:3;stroke:rgb(154,96,180);fill:none\"/><circle cx=\"465\" cy=\"35\" r=\"5\" style=\"stroke-width:3;stroke:rgb(154,96,180);fill:rgb(154,96,180)\"/><text x=\"482\" y=\"41\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><text x=\"9\" y=\"63\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"9\" y=\"96\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.28k</text><text x=\"9\" y=\"130\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.12k</text><text x=\"21\" y=\"164\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"21\" y=\"198\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"21\" y=\"232\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"21\" y=\"266\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"21\" y=\"300\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"21\" y=\"334\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"39\" y=\"368\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path d=\"M 54 57\nL 590 57\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path d=\"M 54 91\nL 590 91\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path d=\"M 54 125\nL 590 125\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path d=\"M 54 159\nL 590 159\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path d=\"M 54 193\nL 590 193\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path d=\"M 54 227\nL 590 227\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path d=\"M 54 261\nL 590 261\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path d=\"M 54 295\nL 590 295\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path d=\"M 54 329\nL 590 329\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path d=\"M 58 364\nL 590 364\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path d=\"M 58 369\nL 58 364\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path d=\"M 146 369\nL 146 364\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path d=\"M 235 369\nL 235 364\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path d=\"M 324 369\nL 324 364\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path d=\"M 412 369\nL 412 364\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path d=\"M 501 369\nL 501 364\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path d=\"M 590 369\nL 590 364\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><text x=\"57\" y=\"390\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"145\" y=\"390\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"234\" y=\"390\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"323\" y=\"390\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"411\" y=\"390\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"500\" y=\"390\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"563\" y=\"390\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><circle cx=\"58\" cy=\"339\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"146\" cy=\"336\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"235\" cy=\"343\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"324\" cy=\"336\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"412\" cy=\"345\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"501\" cy=\"315\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"590\" cy=\"320\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"58\" cy=\"318\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"146\" cy=\"326\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"235\" cy=\"324\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"324\" cy=\"315\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"412\" cy=\"303\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"501\" cy=\"294\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"590\" cy=\"298\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"58\" cy=\"333\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:rgb(100,180,210)\"/><circle cx=\"146\" cy=\"315\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:rgb(100,180,210)\"/><circle cx=\"235\" cy=\"322\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:rgb(100,180,210)\"/><circle cx=\"324\" cy=\"332\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:rgb(100,180,210)\"/><circle cx=\"412\" cy=\"324\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:rgb(100,180,210)\"/><circle cx=\"501\" cy=\"294\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:rgb(100,180,210)\"/><circle cx=\"590\" cy=\"277\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:rgb(100,180,210)\"/><circle cx=\"58\" cy=\"296\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:rgb(64,160,110)\"/><circle cx=\"146\" cy=\"294\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:rgb(64,160,110)\"/><circle cx=\"235\" cy=\"300\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:rgb(64,160,110)\"/><circle cx=\"324\" cy=\"293\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:rgb(64,160,110)\"/><circle cx=\"412\" cy=\"281\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:rgb(64,160,110)\"/><circle cx=\"501\" cy=\"294\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:rgb(64,160,110)\"/><circle cx=\"590\" cy=\"296\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:rgb(64,160,110)\"/><circle cx=\"58\" cy=\"190\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,96,180);fill:rgb(154,96,180)\"/><circle cx=\"146\" cy=\"166\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,96,180);fill:rgb(154,96,180)\"/><circle cx=\"235\" cy=\"172\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,96,180);fill:rgb(154,96,180)\"/><circle cx=\"324\" cy=\"165\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,96,180);fill:rgb(154,96,180)\"/><circle cx=\"412\" cy=\"89\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,96,180);fill:rgb(154,96,180)\"/><circle cx=\"501\" cy=\"81\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,96,180);fill:rgb(154,96,180)\"/><circle cx=\"590\" cy=\"83\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,96,180);fill:rgb(154,96,180)\"/></svg>",
			pngCRC:      0xe4ea791c,
		},
		{
			name: "boundary_gap_enable",
			makeOptions: func() ScatterChartOption {
				opt := makeMinimalScatterChartOption()
				opt.XAxis.Show = Ptr(true)
				opt.XAxis.BoundaryGap = Ptr(true)
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 10 364\nL 590 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 10 369\nL 10 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 92 369\nL 92 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 175 369\nL 175 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 258 369\nL 258 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 341 369\nL 341 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 424 369\nL 424 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 507 369\nL 507 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 590 369\nL 590 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"47\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1</text><text x=\"129\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2</text><text x=\"212\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">3</text><text x=\"295\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">4</text><text x=\"378\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">5</text><text x=\"461\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">6</text><text x=\"544\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">7</text><circle cx=\"51\" cy=\"335\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"133\" cy=\"332\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"216\" cy=\"340\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"299\" cy=\"332\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"382\" cy=\"342\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"465\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"548\" cy=\"313\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"51\" cy=\"163\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"133\" cy=\"135\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"216\" cy=\"143\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"299\" cy=\"135\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"382\" cy=\"47\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"465\" cy=\"38\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"548\" cy=\"40\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/></svg>",
			pngCRC: 0xd9149cb1,
		},
		{
			name: "dual_yaxis",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.Theme = GetTheme(ThemeLight)
				opt.SeriesList[1].YAxisIndex = 1
				opt.YAxis = append(opt.YAxis, opt.YAxis[0])
				opt.YAxis[0].Theme = opt.Theme.WithYAxisSeriesColor(0)
				opt.YAxis[1].Theme = opt.Theme.WithYAxisSeriesColor(1)
				opt.XAxis.Show = Ptr(false)
				opt.Title.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"552\" y=\"16\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.43k</text><text x=\"552\" y=\"58\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.36k</text><text x=\"552\" y=\"100\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.29k</text><text x=\"552\" y=\"142\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.22k</text><text x=\"552\" y=\"184\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.15k</text><text x=\"552\" y=\"226\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.08k</text><text x=\"552\" y=\"268\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.01k</text><text x=\"552\" y=\"310\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">940</text><text x=\"552\" y=\"352\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">870</text><text x=\"552\" y=\"394\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"9\" y=\"16\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">234</text><text x=\"9\" y=\"58\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">218</text><text x=\"9\" y=\"100\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">202</text><text x=\"9\" y=\"142\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">186</text><text x=\"9\" y=\"184\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">170</text><text x=\"9\" y=\"226\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">154</text><text x=\"9\" y=\"268\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">138</text><text x=\"9\" y=\"310\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">122</text><text x=\"9\" y=\"352\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">106</text><text x=\"18\" y=\"394\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">90</text><path d=\"M 42 10\nL 542 10\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 52\nL 542 52\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 94\nL 542 94\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 136\nL 542 136\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 178\nL 542 178\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 221\nL 542 221\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 263\nL 542 263\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 305\nL 542 305\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 347\nL 542 347\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><circle cx=\"46\" cy=\"311\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"128\" cy=\"280\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"211\" cy=\"361\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"294\" cy=\"274\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"376\" cy=\"390\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"459\" cy=\"21\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"542\" cy=\"74\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"46\" cy=\"378\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"128\" cy=\"311\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"211\" cy=\"330\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"294\" cy=\"310\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"376\" cy=\"95\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"459\" cy=\"71\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"542\" cy=\"77\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/></svg>",
			pngCRC: 0x295ebe5b,
		},
		{
			name: "no_nice_interval",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.Theme = GetTheme(ThemeLight)
				opt.SeriesList[1].YAxisIndex = 1
				opt.YAxis = append(opt.YAxis, opt.YAxis[0])
				opt.YAxis[0].PreferNiceIntervals = Ptr(false)
				opt.YAxis[1].PreferNiceIntervals = Ptr(false)
				opt.YAxis[0].Theme = opt.Theme.WithYAxisSeriesColor(0)
				opt.YAxis[1].Theme = opt.Theme.WithYAxisSeriesColor(1)
				opt.XAxis.Show = Ptr(false)
				opt.Title.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"552\" y=\"16\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.43k</text><text x=\"552\" y=\"58\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.36k</text><text x=\"552\" y=\"100\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.29k</text><text x=\"552\" y=\"142\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.22k</text><text x=\"552\" y=\"184\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.15k</text><text x=\"552\" y=\"226\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.08k</text><text x=\"552\" y=\"268\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.01k</text><text x=\"552\" y=\"310\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">940</text><text x=\"552\" y=\"352\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">870</text><text x=\"552\" y=\"394\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"9\" y=\"16\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">234</text><text x=\"9\" y=\"58\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">218</text><text x=\"9\" y=\"100\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">202</text><text x=\"9\" y=\"142\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">186</text><text x=\"9\" y=\"184\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">170</text><text x=\"9\" y=\"226\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">154</text><text x=\"9\" y=\"268\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">138</text><text x=\"9\" y=\"310\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">122</text><text x=\"9\" y=\"352\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">106</text><text x=\"18\" y=\"394\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">90</text><path d=\"M 42 10\nL 542 10\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 52\nL 542 52\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 94\nL 542 94\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 136\nL 542 136\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 178\nL 542 178\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 221\nL 542 221\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 263\nL 542 263\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 305\nL 542 305\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 347\nL 542 347\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><circle cx=\"46\" cy=\"311\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"128\" cy=\"280\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"211\" cy=\"361\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"294\" cy=\"274\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"376\" cy=\"390\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"459\" cy=\"21\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"542\" cy=\"74\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"46\" cy=\"378\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"128\" cy=\"311\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"211\" cy=\"330\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"294\" cy=\"310\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"376\" cy=\"95\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"459\" cy=\"71\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"542\" cy=\"77\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/></svg>",
			pngCRC: 0x295ebe5b,
		},
		{
			name: "left_nice_interval",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.Theme = GetTheme(ThemeLight)
				opt.SeriesList[1].YAxisIndex = 1
				opt.YAxis = append(opt.YAxis, opt.YAxis[0])
				opt.YAxis[0].PreferNiceIntervals = Ptr(true)
				opt.YAxis[1].PreferNiceIntervals = Ptr(false)
				opt.YAxis[0].Theme = opt.Theme.WithYAxisSeriesColor(0)
				opt.YAxis[1].Theme = opt.Theme.WithYAxisSeriesColor(1)
				opt.XAxis.Show = Ptr(false)
				opt.Title.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"552\" y=\"16\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.36k</text><text x=\"552\" y=\"63\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.29k</text><text x=\"552\" y=\"110\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.22k</text><text x=\"552\" y=\"157\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.15k</text><text x=\"552\" y=\"205\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.08k</text><text x=\"552\" y=\"252\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.01k</text><text x=\"552\" y=\"299\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">940</text><text x=\"552\" y=\"346\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">870</text><text x=\"552\" y=\"394\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"9\" y=\"16\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">250</text><text x=\"9\" y=\"63\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">230</text><text x=\"9\" y=\"110\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">210</text><text x=\"9\" y=\"157\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">190</text><text x=\"9\" y=\"205\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">170</text><text x=\"9\" y=\"252\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">150</text><text x=\"9\" y=\"299\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">130</text><text x=\"9\" y=\"346\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">110</text><text x=\"18\" y=\"394\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">90</text><path d=\"M 42 10\nL 542 10\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 57\nL 542 57\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 105\nL 542 105\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 152\nL 542 152\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 200\nL 542 200\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 247\nL 542 247\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 295\nL 542 295\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 342\nL 542 342\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><circle cx=\"46\" cy=\"319\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"128\" cy=\"291\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"211\" cy=\"364\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"294\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"376\" cy=\"390\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"459\" cy=\"58\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"542\" cy=\"105\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"46\" cy=\"377\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"128\" cy=\"301\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"211\" cy=\"322\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"294\" cy=\"300\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"376\" cy=\"58\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"459\" cy=\"31\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"542\" cy=\"38\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/></svg>",
			pngCRC: 0x2f46b989,
		},
		{
			name: "right_nice_interval",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.Theme = GetTheme(ThemeLight)
				opt.SeriesList[1].YAxisIndex = 1
				opt.YAxis = append(opt.YAxis, opt.YAxis[0])
				opt.YAxis[0].PreferNiceIntervals = Ptr(false)
				opt.YAxis[1].PreferNiceIntervals = Ptr(true)
				opt.YAxis[0].Theme = opt.Theme.WithYAxisSeriesColor(0)
				opt.YAxis[1].Theme = opt.Theme.WithYAxisSeriesColor(1)
				opt.XAxis.Show = Ptr(false)
				opt.Title.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"561\" y=\"16\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.4k</text><text x=\"561\" y=\"142\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.2k</text><text x=\"561\" y=\"268\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1k</text><text x=\"561\" y=\"394\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"9\" y=\"16\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">240</text><text x=\"9\" y=\"142\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">190</text><text x=\"9\" y=\"268\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">140</text><text x=\"18\" y=\"394\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">90</text><path d=\"M 42 10\nL 551 10\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 136\nL 551 136\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 263\nL 551 263\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><circle cx=\"46\" cy=\"314\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"130\" cy=\"284\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"214\" cy=\"363\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"298\" cy=\"279\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"382\" cy=\"390\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"466\" cy=\"36\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"551\" cy=\"86\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"46\" cy=\"378\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"130\" cy=\"307\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"214\" cy=\"327\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"298\" cy=\"306\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"382\" cy=\"80\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"466\" cy=\"55\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"551\" cy=\"61\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/></svg>",
			pngCRC: 0x5d2cb52f,
		},
		{
			name: "data_gap",
			makeOptions: func() ScatterChartOption {
				opt := makeMinimalScatterChartOption()
				opt.SeriesList[0].Values[4] = []float64{GetNullValue()}
				opt.SeriesList[1].Values[2] = []float64{GetNullValue()}
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><circle cx=\"10\" cy=\"359\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"106\" cy=\"356\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"203\" cy=\"364\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"300\" cy=\"355\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"396\" cy=\"2147483657\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"493\" cy=\"330\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"590\" cy=\"335\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"10\" cy=\"174\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"106\" cy=\"145\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"203\" cy=\"2147483657\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"300\" cy=\"144\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"396\" cy=\"50\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"493\" cy=\"40\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"590\" cy=\"42\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/></svg>",
			pngCRC: 0x2ec0cf1c,
		},
		{
			name: "mark_line",
			makeOptions: func() ScatterChartOption {
				opt := makeMinimalMultiValueScatterChartOption()
				opt.Padding = NewBoxEqual(40)
				opt.SymbolSize = 4.5
				for i := range opt.SeriesList {
					markLine := NewMarkLine("min", "max", "average")
					markLine.ValueFormatter = func(f float64) string {
						return FormatValueHumanizeShort(f, 0, false)
					}
					opt.SeriesList[i].MarkLine = markLine
				}
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><circle cx=\"40\" cy=\"334\" r=\"5\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"126\" cy=\"331\" r=\"5\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"213\" cy=\"338\" r=\"5\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"213\" cy=\"356\" r=\"5\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"300\" cy=\"331\" r=\"5\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"386\" cy=\"340\" r=\"5\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"386\" cy=\"354\" r=\"5\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"473\" cy=\"309\" r=\"5\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"560\" cy=\"314\" r=\"5\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"40\" cy=\"178\" r=\"5\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"126\" cy=\"153\" r=\"5\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"213\" cy=\"160\" r=\"5\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"213\" cy=\"227\" r=\"5\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"300\" cy=\"153\" r=\"5\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"386\" cy=\"74\" r=\"5\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"473\" cy=\"65\" r=\"5\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"560\" cy=\"67\" r=\"5\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"43\" cy=\"356\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 49 356\nL 542 356\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 542 351\nL 558 356\nL 542 361\nL 547 356\nL 542 351\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"560\" y=\"360\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">20</text><circle cx=\"43\" cy=\"309\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 49 309\nL 542 309\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 542 304\nL 558 309\nL 542 314\nL 547 309\nL 542 304\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"560\" y=\"313\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">230</text><circle cx=\"43\" cy=\"334\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 49 334\nL 542 334\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 542 329\nL 558 334\nL 542 339\nL 547 334\nL 542 329\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"560\" y=\"338\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">118</text><circle cx=\"43\" cy=\"227\" r=\"3\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 49 227\nL 542 227\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 542 222\nL 558 227\nL 542 232\nL 547 227\nL 542 222\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"560\" y=\"231\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">600</text><circle cx=\"43\" cy=\"65\" r=\"3\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 49 65\nL 542 65\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 542 60\nL 558 65\nL 542 70\nL 547 65\nL 542 60\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"560\" y=\"69\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">1k</text><circle cx=\"43\" cy=\"135\" r=\"3\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 49 135\nL 542 135\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 542 130\nL 558 135\nL 542 140\nL 547 135\nL 542 130\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"560\" y=\"139\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">1k</text></svg>",
			pngCRC: 0x77aa06d5,
		},
		{
			name: "series_label",
			makeOptions: func() ScatterChartOption {
				opt := makeMinimalMultiValueScatterChartOption()
				opt.YAxis[0].Show = Ptr(false)
				for i := range opt.SeriesList {
					opt.SeriesList[i].Label.Show = Ptr(true)
					opt.SeriesList[i].Label.FontStyle = FontStyle{
						FontSize:  12.0,
						Font:      GetDefaultFont(),
						FontColor: ColorBlue,
					}
					opt.SeriesList[i].Label.ValueFormatter = func(f float64) string {
						return FormatValueHumanizeShort(f, 2, false)
					}
				}
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><circle cx=\"10\" cy=\"359\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"106\" cy=\"356\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"203\" cy=\"364\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"203\" cy=\"385\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"300\" cy=\"355\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"396\" cy=\"367\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"396\" cy=\"383\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"493\" cy=\"330\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"590\" cy=\"335\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"10\" cy=\"174\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"106\" cy=\"145\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"203\" cy=\"153\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"203\" cy=\"232\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"300\" cy=\"144\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"396\" cy=\"50\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"493\" cy=\"40\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"590\" cy=\"42\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"15\" y=\"365\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">120</text><text x=\"111\" y=\"362\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">132</text><text x=\"208\" y=\"370\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">101</text><text x=\"208\" y=\"391\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">20</text><text x=\"305\" y=\"361\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">134</text><text x=\"401\" y=\"373\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">90</text><text x=\"401\" y=\"389\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">28</text><text x=\"498\" y=\"336\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">230</text><text x=\"573\" y=\"341\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">210</text><text x=\"15\" y=\"180\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">820</text><text x=\"111\" y=\"151\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">932</text><text x=\"208\" y=\"159\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">901</text><text x=\"208\" y=\"238\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">600</text><text x=\"305\" y=\"150\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">934</text><text x=\"401\" y=\"56\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.29k</text><text x=\"498\" y=\"46\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.33k</text><text x=\"561\" y=\"48\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.32k</text></svg>",
			pngCRC: 0x2dd0ba8d,
		},
		{
			name: "symbol_dot",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.XAxis.Labels = opt.XAxis.Labels[:5]
				for i := range opt.SeriesList {
					opt.SeriesList[i].Values = opt.SeriesList[i].Values[:5]
				}
				opt.SymbolSize = 4.5
				opt.Symbol = SymbolDot
				opt.Legend.Symbol = SymbolDot
				opt.Title.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 250 19\nL 280 19\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"265\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"282\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1</text><path d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2</text><circle cx=\"10\" cy=\"362\" r=\"5\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"155\" cy=\"359\" r=\"5\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"300\" cy=\"366\" r=\"5\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"445\" cy=\"358\" r=\"5\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"590\" cy=\"369\" r=\"5\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"10\" cy=\"195\" r=\"5\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"155\" cy=\"168\" r=\"5\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"300\" cy=\"175\" r=\"5\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"445\" cy=\"167\" r=\"5\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"590\" cy=\"82\" r=\"5\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/></svg>",
			pngCRC: 0x65aa700f,
		},
		{
			name: "symbol_circle",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.XAxis.Labels = opt.XAxis.Labels[:5]
				for i := range opt.SeriesList {
					opt.SeriesList[i].Values = opt.SeriesList[i].Values[:5]
				}
				opt.SymbolSize = 4.5
				opt.Symbol = SymbolCircle
				opt.Legend.Symbol = SymbolCircle
				opt.Title.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 250 19\nL 280 19\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"265\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"265\" cy=\"19\" r=\"2\" style=\"stroke-width:3;stroke:white;fill:white\"/><text x=\"282\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1</text><path d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"326\" cy=\"19\" r=\"2\" style=\"stroke-width:3;stroke:white;fill:white\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2</text><circle cx=\"10\" cy=\"362\" r=\"5\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"155\" cy=\"359\" r=\"5\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"300\" cy=\"366\" r=\"5\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"445\" cy=\"358\" r=\"5\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"590\" cy=\"369\" r=\"5\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"10\" cy=\"195\" r=\"5\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"155\" cy=\"168\" r=\"5\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"300\" cy=\"175\" r=\"5\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"445\" cy=\"167\" r=\"5\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"590\" cy=\"82\" r=\"5\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/></svg>",
			pngCRC: 0xdd29c17b,
		},
		{
			name: "symbol_square",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.XAxis.Labels = opt.XAxis.Labels[:5]
				for i := range opt.SeriesList {
					opt.SeriesList[i].Values = opt.SeriesList[i].Values[:5]
				}
				opt.SymbolSize = 4.5
				opt.Symbol = SymbolSquare
				opt.Legend.Symbol = SymbolSquare
				opt.Title.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 250 13\nL 280 13\nL 280 26\nL 250 26\nL 250 13\" style=\"stroke:none;fill:rgb(84,112,198)\"/><text x=\"282\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1</text><path d=\"M 311 13\nL 341 13\nL 341 26\nL 311 26\nL 311 13\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2</text><path d=\"M 5 357\nL 14 357\nL 14 366\nL 5 366\nL 5 357\nM 150 354\nL 159 354\nL 159 363\nL 150 363\nL 150 354\nM 295 361\nL 304 361\nL 304 370\nL 295 370\nL 295 361\nM 440 353\nL 449 353\nL 449 362\nL 440 362\nL 440 353\nM 585 364\nL 594 364\nL 594 373\nL 585 373\nL 585 364\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path d=\"M 5 190\nL 14 190\nL 14 199\nL 5 199\nL 5 190\nM 150 163\nL 159 163\nL 159 172\nL 150 172\nL 150 163\nM 295 170\nL 304 170\nL 304 179\nL 295 179\nL 295 170\nM 440 162\nL 449 162\nL 449 171\nL 440 171\nL 440 162\nM 585 77\nL 594 77\nL 594 86\nL 585 86\nL 585 77\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/></svg>",
			pngCRC: 0x11891d83,
		},
		{
			name: "symbol_diamond",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.XAxis.Labels = opt.XAxis.Labels[:5]
				for i := range opt.SeriesList {
					opt.SeriesList[i].Values = opt.SeriesList[i].Values[:5]
				}
				opt.SymbolSize = 4.5
				opt.Symbol = SymbolDiamond
				opt.Legend.Symbol = SymbolDiamond
				opt.Title.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 265 10\nL 272 20\nL 265 30\nL 258 20\nL 265 10\" style=\"stroke:none;fill:rgb(84,112,198)\"/><text x=\"282\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1</text><path d=\"M 316 10\nL 323 20\nL 316 30\nL 309 20\nL 316 10\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"333\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2</text><path d=\"M 10 356\nL 16 362\nL 10 368\nL 4 362\nL 10 356\nM 155 353\nL 161 359\nL 155 365\nL 149 359\nL 155 353\nM 300 360\nL 306 366\nL 300 372\nL 294 366\nL 300 360\nM 445 352\nL 451 358\nL 445 364\nL 439 358\nL 445 352\nM 590 363\nL 596 369\nL 590 375\nL 584 369\nL 590 363\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path d=\"M 10 189\nL 16 195\nL 10 201\nL 4 195\nL 10 189\nM 155 162\nL 161 168\nL 155 174\nL 149 168\nL 155 162\nM 300 169\nL 306 175\nL 300 181\nL 294 175\nL 300 169\nM 445 161\nL 451 167\nL 445 173\nL 439 167\nL 445 161\nM 590 76\nL 596 82\nL 590 88\nL 584 82\nL 590 76\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/></svg>",
			pngCRC: 0x5caf6eb3,
		},
		{
			name:   "symbol_mixed",
			ignore: "size", // svg is too big to commit
			makeOptions: func() ScatterChartOption {
				opt := makeFullScatterChartOption()
				opt.XAxis.Labels = opt.XAxis.Labels[:5]
				for i := range opt.SeriesList {
					opt.SeriesList[i].Values = opt.SeriesList[i].Values[:5]
				}
				opt.SymbolSize = 4.0
				opt.SeriesList[0].Symbol = SymbolCircle
				opt.SeriesList[1].Symbol = SymbolSquare
				opt.SeriesList[2].Symbol = SymbolDiamond
				opt.SeriesList[3].Symbol = SymbolDot
				opt.Title.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
			svg:    "",
			pngCRC: 0,
		},
		{
			name:   "dense_trends",
			ignore: "size", // svg is too big to commit
			makeOptions: func() ScatterChartOption {
				opt := makeDenseScatterChartOption()
				for i := range opt.SeriesList {
					opt.SeriesList[i].TrendLine[0].StrokeSmoothingTension = 0.9 // smooth average line
					opt.SeriesList[i].TrendLine[0].Period = 5
					c1 := Color{
						R: uint8(80 + (20 * i)),
						G: uint8(80 + (20 * i)),
						B: uint8(80 + (20 * i)),
						A: 255,
					}
					c2 := c1
					if i%2 == 0 {
						c2.R = 200
					} else {
						c2.B = 200
					}
					trendLine1 := SeriesTrendLine{
						Type:      SeriesTrendTypeCubic,
						LineColor: c1,
					}
					trendLine2 := SeriesTrendLine{
						Type:      SeriesTrendTypeLinear,
						LineColor: c2,
					}
					opt.SeriesList[i].TrendLine = append(opt.SeriesList[i].TrendLine, trendLine1, trendLine2)
				}
				// disable extras
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
			svg:    "",
			pngCRC: 0,
		},
		{
			name: "trend_line_dashed",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.Theme = GetDefaultTheme()
				opt.Title.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				// Add dashed trend lines to each series
				for i := range opt.SeriesList {
					opt.SeriesList[i].TrendLine = []SeriesTrendLine{
						{
							StrokeSmoothingTension: 0.8,
							Type:                   SeriesTrendTypeSMA,
							DashedLine:             Ptr(true), // Explicitly set to dashed
							LineColor:              opt.Theme.GetSeriesTrendColor(i).WithAdjustHSL(0, .2, -.2),
						},
						{
							Type:       SeriesTrendTypeCubic,
							DashedLine: Ptr(true), // Explicitly set to dashed
							LineColor:  opt.Theme.GetSeriesTrendColor(i).WithAdjustHSL(0, .4, -.4),
						},
					}
				}
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><circle cx=\"10\" cy=\"359\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"106\" cy=\"356\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"203\" cy=\"364\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"300\" cy=\"355\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"396\" cy=\"367\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"493\" cy=\"330\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"590\" cy=\"335\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"10\" cy=\"174\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"106\" cy=\"145\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"203\" cy=\"153\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"300\" cy=\"144\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"396\" cy=\"50\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"493\" cy=\"40\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"590\" cy=\"42\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path stroke-dasharray=\"9.6, 7.7\" d=\"M 10 357\nQ106,359 144,358\nQ203,358 241,359\nQ300,362 338,357\nQ396,351 434,348\nQ493,344 531,339\nQ493,344 590,332\" style=\"stroke-width:2;stroke:rgb(12,38,115);fill:none\"/><path stroke-dasharray=\"9.6, 7.7\" d=\"M 10 357\nL 106 361\nL 203 361\nL 300 359\nL 396 353\nL 493 344\nL 590 331\" style=\"stroke-width:2;stroke:rgb(0,6,25);fill:none\"/><path stroke-dasharray=\"9.6, 7.7\" d=\"M 10 159\nQ106,157 144,153\nQ203,147 241,134\nQ300,116 338,100\nQ396,78 434,64\nQ493,44 531,42\nQ493,44 590,41\" style=\"stroke-width:2;stroke:rgb(61,146,20);fill:none\"/><path stroke-dasharray=\"9.6, 7.7\" d=\"M 10 164\nL 106 170\nL 203 148\nL 300 112\nL 396 73\nL 493 44\nL 590 37\" style=\"stroke-width:2;stroke:rgb(21,63,1);fill:none\"/></svg>",
			pngCRC: 0x8faf9659,
		},
		{
			name: "with_conditional_labels",
			makeOptions: func() ScatterChartOption {
				return ScatterChartOption{
					Padding: NewBoxEqual(10),
					XAxis: XAxisOption{
						Labels: []string{"A", "B", "C", "D", "E"},
					},
					YAxis: []YAxisOption{{}},
					SeriesList: NewSeriesListScatter([][]float64{
						{50, 150, 100, 200, 175},
						{75, 125, 90, 160, 140},
					}, ScatterSeriesOption{
						Names: []string{"Dataset1", "Dataset2"},
						Label: SeriesLabel{
							Show: Ptr(true),
							LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
								// Show labels only for values above 120
								if val > 120 {
									switch {
									case val >= 180: // High values - gold styling
										return "⭐ " + strconv.FormatFloat(val, 'f', 0, 64), &LabelStyle{
											FontStyle:       FontStyle{FontColor: ColorBlack, FontSize: 14},
											BackgroundColor: ColorFromHex("#FFD700"), // Gold
											CornerRadius:    6,
										}
									case val >= 150: // Medium-high values - silver styling
										return "📊 " + strconv.FormatFloat(val, 'f', 0, 64), &LabelStyle{
											FontStyle:       FontStyle{FontColor: ColorBlack, FontSize: 12},
											BackgroundColor: ColorFromHex("#C0C0C0"), // Silver
											CornerRadius:    4,
										}
									default: // Values above 120 but below 150 - simple styling
										return strconv.FormatFloat(val, 'f', 0, 64), &LabelStyle{
											FontStyle: FontStyle{FontColor: ColorBlue, FontSize: 10},
										}
									}
								}
								// Hide labels for values <= 120
								return "", nil
							},
						},
					}),
					Title: TitleOption{
						Show: Ptr(false),
					},
					Legend: LegendOption{
						Show: Ptr(false),
					},
				}
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"9\" y=\"16\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">230</text><text x=\"9\" y=\"55\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">210</text><text x=\"9\" y=\"94\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">190</text><text x=\"9\" y=\"133\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">170</text><text x=\"9\" y=\"172\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">150</text><text x=\"9\" y=\"211\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">130</text><text x=\"9\" y=\"250\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">110</text><text x=\"18\" y=\"289\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">90</text><text x=\"18\" y=\"328\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">70</text><text x=\"18\" y=\"368\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">50</text><path d=\"M 42 10\nL 590 10\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 49\nL 590 49\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 88\nL 590 88\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 128\nL 590 128\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 167\nL 590 167\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 206\nL 590 206\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 246\nL 590 246\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 285\nL 590 285\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 324\nL 590 324\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 46 364\nL 590 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 46 369\nL 46 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 182 369\nL 182 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 318 369\nL 318 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 454 369\nL 454 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 590 369\nL 590 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"45\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"181\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"317\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"453\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">D</text><text x=\"581\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">E</text><circle cx=\"46\" cy=\"364\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"182\" cy=\"168\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"318\" cy=\"266\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"454\" cy=\"69\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"590\" cy=\"119\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"46\" cy=\"315\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"182\" cy=\"217\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"318\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"454\" cy=\"148\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"590\" cy=\"187\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path d=\"M 187 154\nL 229 154\nL 229 154\nA 4 4 90.00 0 1 233 158\nL 233 174\nL 233 174\nA 4 4 90.00 0 1 229 178\nL 187 178\nL 187 178\nA 4 4 90.00 0 1 183 174\nL 183 158\nL 183 158\nA 4 4 90.00 0 1 187 154\nZ\" style=\"stroke:none;fill:silver\"/><text x=\"187\" y=\"174\" style=\"stroke:none;fill:black;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">📊 150</text><path d=\"M 461 54\nL 506 54\nL 506 54\nA 6 6 90.00 0 1 512 60\nL 512 74\nL 512 74\nA 6 6 90.00 0 1 506 80\nL 461 80\nL 461 80\nA 6 6 90.00 0 1 455 74\nL 455 60\nL 455 60\nA 6 6 90.00 0 1 461 54\nZ\" style=\"stroke:none;fill:rgb(255,215,0)\"/><text x=\"459\" y=\"76\" style=\"stroke:none;fill:black;font-size:17.9px;font-family:'Roboto Medium',sans-serif\">⭐ 200</text><path d=\"M 558 105\nL 600 105\nL 600 105\nA 4 4 90.00 0 1 604 109\nL 604 125\nL 604 125\nA 4 4 90.00 0 1 600 129\nL 558 129\nL 558 129\nA 4 4 90.00 0 1 554 125\nL 554 109\nL 554 109\nA 4 4 90.00 0 1 558 105\nZ\" style=\"stroke:none;fill:silver\"/><text x=\"558\" y=\"125\" style=\"stroke:none;fill:black;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">📊 175</text><text x=\"187\" y=\"221\" style=\"stroke:none;fill:blue;font-size:12.8px;font-family:'Roboto Medium',sans-serif\">125</text><path d=\"M 459 134\nL 501 134\nL 501 134\nA 4 4 90.00 0 1 505 138\nL 505 154\nL 505 154\nA 4 4 90.00 0 1 501 158\nL 459 158\nL 459 158\nA 4 4 90.00 0 1 455 154\nL 455 138\nL 455 138\nA 4 4 90.00 0 1 459 134\nZ\" style=\"stroke:none;fill:silver\"/><text x=\"459\" y=\"154\" style=\"stroke:none;fill:black;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">📊 160</text><text x=\"578\" y=\"191\" style=\"stroke:none;fill:blue;font-size:12.8px;font-family:'Roboto Medium',sans-serif\">140</text></svg>",
			pngCRC: 0xce19f6cc,
		},
		{
			name: "bollinger",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.SeriesList[0].TrendLine = []SeriesTrendLine{
					{Type: SeriesTrendTypeBollingerLower, Period: 3},
				}
				opt.SeriesList[1].TrendLine = []SeriesTrendLine{
					{Type: SeriesTrendTypeBollingerUpper, Period: 3},
				}
				opt.Legend.Show = Ptr(false)
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"9\" y=\"16\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"9\" y=\"55\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.28k</text><text x=\"9\" y=\"94\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.12k</text><text x=\"21\" y=\"133\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"21\" y=\"172\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"21\" y=\"211\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"21\" y=\"250\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"21\" y=\"289\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"21\" y=\"328\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"39\" y=\"368\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path d=\"M 54 10\nL 590 10\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 54 49\nL 590 49\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 54 88\nL 590 88\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 54 128\nL 590 128\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 54 167\nL 590 167\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 54 206\nL 590 206\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 54 246\nL 590 246\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 54 285\nL 590 285\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 54 324\nL 590 324\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 58 364\nL 590 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 58 369\nL 58 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 146 369\nL 146 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 235 369\nL 235 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 324 369\nL 324 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 412 369\nL 412 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 501 369\nL 501 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 590 369\nL 590 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"57\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"145\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"234\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"323\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">D</text><text x=\"411\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"500\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">F</text><text x=\"579\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">G</text><circle cx=\"58\" cy=\"335\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"146\" cy=\"332\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"235\" cy=\"340\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"324\" cy=\"332\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"412\" cy=\"342\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"501\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"590\" cy=\"313\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"58\" cy=\"163\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"146\" cy=\"135\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"235\" cy=\"143\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"324\" cy=\"135\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"412\" cy=\"47\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"501\" cy=\"38\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"590\" cy=\"40\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path d=\"M 58 336\nL 146 342\nL 235 342\nL 324 347\nL 412 356\nL 501 351\nL 590 315\" style=\"stroke-width:2;stroke:rgb(46,80,184);fill:none\"/><path d=\"M 58 122\nL 146 124\nL 235 130\nL 324 22\nL 412 10\nL 501 33\nL 590 36\" style=\"stroke-width:2;stroke:rgb(111,202,67);fill:none\"/></svg>",
			pngCRC: 0x2731f499,
		},
		{
			name: "rsi",
			makeOptions: func() ScatterChartOption {
				opt := makeBasicScatterChartOption()
				opt.SeriesList[0].TrendLine = []SeriesTrendLine{
					{Type: SeriesTrendTypeRSI, Period: 3},
				}
				opt.Legend.Show = Ptr(false)
				return opt
			},
			svg:    "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"9\" y=\"16\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"9\" y=\"55\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.28k</text><text x=\"9\" y=\"94\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.12k</text><text x=\"21\" y=\"133\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"21\" y=\"172\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"21\" y=\"211\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"21\" y=\"250\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"21\" y=\"289\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"21\" y=\"328\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"39\" y=\"368\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path d=\"M 54 10\nL 590 10\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 54 49\nL 590 49\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 54 88\nL 590 88\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 54 128\nL 590 128\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 54 167\nL 590 167\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 54 206\nL 590 206\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 54 246\nL 590 246\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 54 285\nL 590 285\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 54 324\nL 590 324\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 58 364\nL 590 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 58 369\nL 58 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 146 369\nL 146 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 235 369\nL 235 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 324 369\nL 324 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 412 369\nL 412 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 501 369\nL 501 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 590 369\nL 590 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"57\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"145\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">B</text><text x=\"234\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">C</text><text x=\"323\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">D</text><text x=\"411\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">E</text><text x=\"500\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">F</text><text x=\"579\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">G</text><circle cx=\"58\" cy=\"335\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"146\" cy=\"332\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"235\" cy=\"340\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"324\" cy=\"332\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"412\" cy=\"342\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"501\" cy=\"308\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"590\" cy=\"313\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"58\" cy=\"163\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"146\" cy=\"135\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"235\" cy=\"143\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"324\" cy=\"135\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"412\" cy=\"47\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"501\" cy=\"38\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"590\" cy=\"40\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path d=\"M 324 350\nL 412 357\nL 501 345\nL 590 348\" style=\"stroke-width:2;stroke:rgb(46,80,184);fill:none\"/></svg>",
			pngCRC: 0x6a3d4fd4,
		},
	}

	for i, tt := range tests {
		if tt.ignore != "" {
			continue
		}
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
		if !tt.themed {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				p := NewPainter(painterOptions)
				rp := NewPainter(rasterOptions)

				validateScatterChartRender(t, p, rp, tt.makeOptions(), tt.svg, tt.pngCRC)
			})
		} else {
			theme := GetTheme(ThemeVividDark)
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-theme_painter", func(t *testing.T) {
				p := NewPainter(painterOptions, PainterThemeOption(theme))
				rp := NewPainter(rasterOptions, PainterThemeOption(theme))

				validateScatterChartRender(t, p, rp, tt.makeOptions(), tt.svg, tt.pngCRC)
			})
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-theme_opt", func(t *testing.T) {
				p := NewPainter(painterOptions)
				rp := NewPainter(rasterOptions)
				opt := tt.makeOptions()
				opt.Theme = theme

				validateScatterChartRender(t, p, rp, opt, tt.svg, tt.pngCRC)
			})
		}
	}
}

func validateScatterChartRender(t *testing.T, svgP, pngP *Painter, opt ScatterChartOption, expectedSVG string, expectedCRC uint32) {
	t.Helper()

	err := svgP.ScatterChart(opt)
	require.NoError(t, err)
	data, err := svgP.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, expectedSVG, data)

	err = pngP.ScatterChart(opt)
	require.NoError(t, err)
	rasterData, err := pngP.Bytes()
	require.NoError(t, err)
	assertEqualPNGCRC(t, expectedCRC, rasterData)
}

func TestScatterChartError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		makeOptions      func() ScatterChartOption
		errorMsgContains string
	}{
		{
			name: "empty_series",
			makeOptions: func() ScatterChartOption {
				return NewScatterChartOptionWithData([][]float64{})
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

			err := p.ScatterChart(tt.makeOptions())
			require.Error(t, err)
			require.ErrorContains(t, err, tt.errorMsgContains)
		})
	}
}
