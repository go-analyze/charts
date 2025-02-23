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
		result      string
	}{
		{
			name:        "basic_themed",
			themed:      true,
			makeOptions: makeFullScatterChartOption,
			result:      "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:rgb(40,40,40)\"/><path  d=\"M 21 19\nL 51 19\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:none\"/><circle cx=\"36\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><text x=\"53\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Email</text><path  d=\"M 112 19\nL 142 19\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:none\"/><circle cx=\"127\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><text x=\"144\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Union Ads</text><path  d=\"M 235 19\nL 265 19\" style=\"stroke-width:3;stroke:rgb(100,180,210);fill:none\"/><circle cx=\"250\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(100,180,210);fill:rgb(100,180,210)\"/><text x=\"267\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Video Ads</text><path  d=\"M 358 19\nL 388 19\" style=\"stroke-width:3;stroke:rgb(64,160,110);fill:none\"/><circle cx=\"373\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(64,160,110);fill:rgb(64,160,110)\"/><text x=\"390\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Direct</text><path  d=\"M 451 19\nL 481 19\" style=\"stroke-width:3;stroke:rgb(154,100,180);fill:none\"/><circle cx=\"466\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(154,100,180);fill:rgb(154,100,180)\"/><text x=\"483\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Search Engine</text><text x=\"10\" y=\"25\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Scatter</text><text x=\"10\" y=\"52\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.44k</text><text x=\"10\" y=\"87\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.28k</text><text x=\"10\" y=\"122\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.12k</text><text x=\"22\" y=\"157\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">960</text><text x=\"22\" y=\"192\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"22\" y=\"227\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">640</text><text x=\"22\" y=\"262\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">480</text><text x=\"22\" y=\"297\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">320</text><text x=\"22\" y=\"332\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"40\" y=\"367\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 59 45\nL 590 45\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 59 80\nL 590 80\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 59 115\nL 590 115\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 59 150\nL 590 150\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 59 185\nL 590 185\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 59 220\nL 590 220\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 59 255\nL 590 255\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 59 290\nL 590 290\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 59 325\nL 590 325\" style=\"stroke-width:1;stroke:rgb(72,71,83);fill:none\"/><path  d=\"M 59 365\nL 59 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 147 365\nL 147 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 236 365\nL 236 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 324 365\nL 324 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 413 365\nL 413 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 501 365\nL 501 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 590 365\nL 590 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><path  d=\"M 59 360\nL 590 360\" style=\"stroke-width:1;stroke:rgb(185,184,206);fill:none\"/><text x=\"58\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mon</text><text x=\"146\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Tue</text><text x=\"235\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Wed</text><text x=\"323\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Thu</text><text x=\"412\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fri</text><text x=\"500\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sat</text><text x=\"563\" y=\"385\" style=\"stroke:none;fill:rgb(238,238,238);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sun</text><circle cx=\"59\" cy=\"334\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"147\" cy=\"332\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"236\" cy=\"338\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"324\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"413\" cy=\"341\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"501\" cy=\"310\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"590\" cy=\"315\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,100,100);fill:rgb(255,100,100)\"/><circle cx=\"59\" cy=\"312\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"147\" cy=\"321\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"236\" cy=\"319\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"324\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"413\" cy=\"297\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"501\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"590\" cy=\"293\" r=\"2\" style=\"stroke-width:1;stroke:rgb(255,210,100);fill:rgb(255,210,100)\"/><circle cx=\"59\" cy=\"328\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:rgb(100,180,210)\"/><circle cx=\"147\" cy=\"310\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:rgb(100,180,210)\"/><circle cx=\"236\" cy=\"317\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:rgb(100,180,210)\"/><circle cx=\"324\" cy=\"327\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:rgb(100,180,210)\"/><circle cx=\"413\" cy=\"319\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:rgb(100,180,210)\"/><circle cx=\"501\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:rgb(100,180,210)\"/><circle cx=\"590\" cy=\"271\" r=\"2\" style=\"stroke-width:1;stroke:rgb(100,180,210);fill:rgb(100,180,210)\"/><circle cx=\"59\" cy=\"290\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:rgb(64,160,110)\"/><circle cx=\"147\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:rgb(64,160,110)\"/><circle cx=\"236\" cy=\"295\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:rgb(64,160,110)\"/><circle cx=\"324\" cy=\"287\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:rgb(64,160,110)\"/><circle cx=\"413\" cy=\"275\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:rgb(64,160,110)\"/><circle cx=\"501\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:rgb(64,160,110)\"/><circle cx=\"590\" cy=\"290\" r=\"2\" style=\"stroke-width:1;stroke:rgb(64,160,110);fill:rgb(64,160,110)\"/><circle cx=\"59\" cy=\"181\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:rgb(154,100,180)\"/><circle cx=\"147\" cy=\"157\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:rgb(154,100,180)\"/><circle cx=\"236\" cy=\"163\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:rgb(154,100,180)\"/><circle cx=\"324\" cy=\"156\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:rgb(154,100,180)\"/><circle cx=\"413\" cy=\"78\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:rgb(154,100,180)\"/><circle cx=\"501\" cy=\"70\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:rgb(154,100,180)\"/><circle cx=\"590\" cy=\"72\" r=\"2\" style=\"stroke-width:1;stroke:rgb(154,100,180);fill:rgb(154,100,180)\"/></svg>",
		},
		{
			name: "boundary_gap_enable",
			makeOptions: func() ScatterChartOption {
				opt := makeMinimalScatterChartOption()
				opt.XAxis.Show = Ptr(true)
				opt.XAxis.BoundaryGap = Ptr(true)
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 10 365\nL 10 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 92 365\nL 92 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 175 365\nL 175 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 258 365\nL 258 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 341 365\nL 341 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 424 365\nL 424 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 507 365\nL 507 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 590 365\nL 590 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 10 360\nL 590 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"47\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1</text><text x=\"129\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2</text><text x=\"212\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">3</text><text x=\"295\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">4</text><text x=\"378\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">5</text><text x=\"461\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">6</text><text x=\"544\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">7</text><circle cx=\"51\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"133\" cy=\"328\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"216\" cy=\"336\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"299\" cy=\"328\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"382\" cy=\"339\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"465\" cy=\"305\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"548\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"51\" cy=\"161\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"133\" cy=\"134\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"216\" cy=\"142\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"299\" cy=\"133\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"382\" cy=\"47\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"465\" cy=\"37\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"548\" cy=\"40\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/></svg>",
		},
		{
			name: "double_yaxis",
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"551\" y=\"17\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.43k</text><text x=\"551\" y=\"55\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.36k</text><text x=\"551\" y=\"94\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.29k</text><text x=\"551\" y=\"133\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.22k</text><text x=\"551\" y=\"172\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.15k</text><text x=\"551\" y=\"211\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.08k</text><text x=\"551\" y=\"250\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.01k</text><text x=\"551\" y=\"289\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">940</text><text x=\"551\" y=\"328\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">870</text><text x=\"551\" y=\"367\" style=\"stroke:none;fill:rgb(145,204,117);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">800</text><text x=\"10\" y=\"17\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">234</text><text x=\"10\" y=\"55\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">218</text><text x=\"10\" y=\"94\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">202</text><text x=\"10\" y=\"133\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">186</text><text x=\"10\" y=\"172\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">170</text><text x=\"10\" y=\"211\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">154</text><text x=\"10\" y=\"250\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">138</text><text x=\"10\" y=\"289\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">122</text><text x=\"10\" y=\"328\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">106</text><text x=\"19\" y=\"367\" style=\"stroke:none;fill:rgb(84,112,198);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">90</text><path  d=\"M 47 10\nL 541 10\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 47 48\nL 541 48\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 47 87\nL 541 87\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 47 126\nL 541 126\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 47 165\nL 541 165\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 47 204\nL 541 204\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 47 243\nL 541 243\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 47 282\nL 541 282\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 47 321\nL 541 321\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><circle cx=\"47\" cy=\"288\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"129\" cy=\"258\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"211\" cy=\"334\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"294\" cy=\"254\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"376\" cy=\"360\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"458\" cy=\"20\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"541\" cy=\"69\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"47\" cy=\"349\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"129\" cy=\"287\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"211\" cy=\"304\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"294\" cy=\"286\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"376\" cy=\"88\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"458\" cy=\"66\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"541\" cy=\"72\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/></svg>",
		},
		{
			name: "data_gap",
			makeOptions: func() ScatterChartOption {
				opt := makeMinimalScatterChartOption()
				opt.SeriesList[0].Values[4] = []float64{GetNullValue()}
				opt.SeriesList[1].Values[2] = []float64{GetNullValue()}
				return opt
			},
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><circle cx=\"10\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"106\" cy=\"328\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"203\" cy=\"336\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"300\" cy=\"328\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"396\" cy=\"2147483657\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"493\" cy=\"305\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"590\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"10\" cy=\"161\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"106\" cy=\"134\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"203\" cy=\"2147483657\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"300\" cy=\"133\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"396\" cy=\"47\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"493\" cy=\"37\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"590\" cy=\"40\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/></svg>",
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><circle cx=\"40\" cy=\"306\" r=\"4\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"126\" cy=\"304\" r=\"4\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"213\" cy=\"310\" r=\"4\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"213\" cy=\"326\" r=\"4\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"300\" cy=\"304\" r=\"4\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"386\" cy=\"312\" r=\"4\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"386\" cy=\"325\" r=\"4\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"473\" cy=\"284\" r=\"4\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"560\" cy=\"288\" r=\"4\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"40\" cy=\"165\" r=\"4\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"126\" cy=\"143\" r=\"4\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"213\" cy=\"149\" r=\"4\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"213\" cy=\"210\" r=\"4\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"300\" cy=\"142\" r=\"4\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"386\" cy=\"71\" r=\"4\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"473\" cy=\"63\" r=\"4\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"560\" cy=\"65\" r=\"4\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"43\" cy=\"326\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 49 326\nL 542 326\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 542 321\nL 558 326\nL 542 331\nL 547 326\nL 542 321\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"560\" y=\"330\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">20</text><circle cx=\"43\" cy=\"284\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 49 284\nL 542 284\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 542 279\nL 558 284\nL 542 289\nL 547 284\nL 542 279\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"560\" y=\"288\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">230</text><circle cx=\"43\" cy=\"307\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 49 307\nL 542 307\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 542 302\nL 558 307\nL 542 312\nL 547 307\nL 542 302\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"560\" y=\"311\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">118</text><circle cx=\"43\" cy=\"210\" r=\"3\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 49 210\nL 542 210\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 542 205\nL 558 210\nL 542 215\nL 547 210\nL 542 205\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"560\" y=\"214\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">600</text><circle cx=\"43\" cy=\"63\" r=\"3\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 49 63\nL 542 63\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 542 58\nL 558 63\nL 542 68\nL 547 63\nL 542 58\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"560\" y=\"67\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">1k</text><circle cx=\"43\" cy=\"126\" r=\"3\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 49 126\nL 542 126\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 542 121\nL 558 126\nL 542 131\nL 547 126\nL 542 121\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"560\" y=\"130\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">1k</text></svg>",
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><circle cx=\"10\" cy=\"331\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"106\" cy=\"328\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"203\" cy=\"336\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"203\" cy=\"356\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"300\" cy=\"328\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"396\" cy=\"339\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"396\" cy=\"354\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"493\" cy=\"305\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"590\" cy=\"309\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"10\" cy=\"161\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"106\" cy=\"134\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"203\" cy=\"142\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"203\" cy=\"215\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"300\" cy=\"133\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"396\" cy=\"47\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"493\" cy=\"37\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"590\" cy=\"40\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"15\" y=\"336\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">120</text><text x=\"111\" y=\"333\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">132</text><text x=\"208\" y=\"341\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">101</text><text x=\"208\" y=\"361\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">20</text><text x=\"305\" y=\"333\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">134</text><text x=\"401\" y=\"344\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">90</text><text x=\"401\" y=\"359\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">28</text><text x=\"498\" y=\"310\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">230</text><text x=\"595\" y=\"314\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">210</text><text x=\"15\" y=\"166\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">820</text><text x=\"111\" y=\"139\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">932</text><text x=\"208\" y=\"147\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">901</text><text x=\"208\" y=\"220\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">600</text><text x=\"305\" y=\"138\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">934</text><text x=\"401\" y=\"52\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.29k</text><text x=\"498\" y=\"42\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.33k</text><text x=\"595\" y=\"45\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.32k</text></svg>",
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 250 19\nL 280 19\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"265\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"282\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1</text><path  d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2</text><circle cx=\"10\" cy=\"334\" r=\"4\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"155\" cy=\"332\" r=\"4\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"300\" cy=\"338\" r=\"4\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"445\" cy=\"331\" r=\"4\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"590\" cy=\"341\" r=\"4\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"10\" cy=\"181\" r=\"4\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"155\" cy=\"157\" r=\"4\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"300\" cy=\"163\" r=\"4\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"445\" cy=\"156\" r=\"4\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"590\" cy=\"78\" r=\"4\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/></svg>",
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 250 19\nL 280 19\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"265\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><circle cx=\"265\" cy=\"19\" r=\"2\" style=\"stroke-width:3;stroke:white;fill:white\"/><text x=\"282\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1</text><path  d=\"M 311 19\nL 341 19\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"326\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><circle cx=\"326\" cy=\"19\" r=\"2\" style=\"stroke-width:3;stroke:white;fill:white\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2</text><circle cx=\"10\" cy=\"334\" r=\"4\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"155\" cy=\"332\" r=\"4\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"300\" cy=\"338\" r=\"4\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"445\" cy=\"331\" r=\"4\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"590\" cy=\"341\" r=\"4\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"10\" cy=\"181\" r=\"4\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"155\" cy=\"157\" r=\"4\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"300\" cy=\"163\" r=\"4\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"445\" cy=\"156\" r=\"4\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"590\" cy=\"78\" r=\"4\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/></svg>",
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 250 13\nL 280 13\nL 280 26\nL 250 26\nL 250 13\" style=\"stroke:none;fill:rgb(84,112,198)\"/><text x=\"282\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1</text><path  d=\"M 311 13\nL 341 13\nL 341 26\nL 311 26\nL 311 13\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"343\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2</text><path  d=\"M 5 329\nL 14 329\nL 14 338\nL 5 338\nL 5 329\nM 150 327\nL 159 327\nL 159 336\nL 150 336\nL 150 327\nM 295 333\nL 304 333\nL 304 342\nL 295 342\nL 295 333\nM 440 326\nL 449 326\nL 449 335\nL 440 335\nL 440 326\nM 585 336\nL 594 336\nL 594 345\nL 585 345\nL 585 336\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path  d=\"M 5 176\nL 14 176\nL 14 185\nL 5 185\nL 5 176\nM 150 152\nL 159 152\nL 159 161\nL 150 161\nL 150 152\nM 295 158\nL 304 158\nL 304 167\nL 295 167\nL 295 158\nM 440 151\nL 449 151\nL 449 160\nL 440 160\nL 440 151\nM 585 73\nL 594 73\nL 594 82\nL 585 82\nL 585 73\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/></svg>",
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
			result: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 255 10\nL 262 20\nL 255 30\nL 248 20\nL 255 10\" style=\"stroke:none;fill:rgb(84,112,198)\"/><text x=\"272\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1</text><path  d=\"M 306 10\nL 313 20\nL 306 30\nL 299 20\nL 306 10\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"323\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2</text><path  d=\"M 10 328\nL 16 334\nL 10 340\nL 4 334\nL 10 328\nM 155 326\nL 161 332\nL 155 338\nL 149 332\nL 155 326\nM 300 332\nL 306 338\nL 300 344\nL 294 338\nL 300 332\nM 445 325\nL 451 331\nL 445 337\nL 439 331\nL 445 325\nM 590 335\nL 596 341\nL 590 347\nL 584 341\nL 590 335\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path  d=\"M 10 175\nL 16 181\nL 10 187\nL 4 181\nL 10 175\nM 155 151\nL 161 157\nL 155 163\nL 149 157\nL 155 151\nM 300 157\nL 306 163\nL 300 169\nL 294 163\nL 300 157\nM 445 150\nL 451 156\nL 445 162\nL 439 156\nL 445 150\nM 590 72\nL 596 78\nL 590 84\nL 584 78\nL 590 72\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/></svg>",
		},
		{
			name:   "symbol_mixed",
			ignore: "size", // result is too big to commit
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
			result: "",
		},
		{
			name:   "dense_trends",
			ignore: "size", // result is too big to commit
			makeOptions: func() ScatterChartOption {
				opt := makeDenseScatterChartOption()
				for i := range opt.SeriesList {
					opt.SeriesList[i].TrendLine[0].StrokeSmoothingTension = 0.9 // smooth average line
					opt.SeriesList[i].TrendLine[0].Window = 5
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
			result: "",
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
		if !tt.themed {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				p := NewPainter(painterOptions)

				validateScatterChartRender(t, p, tt.makeOptions(), tt.result)
			})
		} else {
			theme := GetTheme(ThemeVividDark)
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-theme_painter", func(t *testing.T) {
				p := NewPainter(painterOptions, PainterThemeOption(theme))

				validateScatterChartRender(t, p, tt.makeOptions(), tt.result)
			})
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-theme_opt", func(t *testing.T) {
				p := NewPainter(painterOptions)
				opt := tt.makeOptions()
				opt.Theme = theme

				validateScatterChartRender(t, p, opt, tt.result)
			})
		}
	}
}

func validateScatterChartRender(t *testing.T, p *Painter, opt ScatterChartOption, expectedResult string) {
	t.Helper()

	err := p.ScatterChart(opt)
	require.NoError(t, err)
	data, err := p.Bytes()
	require.NoError(t, err)
	assertEqualSVG(t, expectedResult, data)
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
