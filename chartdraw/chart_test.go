package chartdraw

import (
	"bytes"
	"image"
	"image/png"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-analyze/charts/chartdraw/drawing"
)

func TestChartGetDPI(t *testing.T) {
	t.Parallel()

	unset := Chart{}
	assert.InDelta(t, DefaultDPI, unset.GetDPI(), 0)
	assert.InDelta(t, float64(192), unset.GetDPI(192), 0)

	set := Chart{DPI: 128}
	assert.InDelta(t, float64(128), set.GetDPI(), 0)
	assert.InDelta(t, float64(128), set.GetDPI(192), 0)
}

func TestChartGetFont(t *testing.T) {
	t.Parallel()

	unset := Chart{}
	require.Nil(t, unset.GetFont())

	set := Chart{Font: GetDefaultFont()}
	require.NotNil(t, set.GetFont())
}

func TestChartGetWidth(t *testing.T) {
	t.Parallel()

	unset := Chart{}
	assert.Equal(t, DefaultChartWidth, unset.GetWidth())

	set := Chart{Width: DefaultChartWidth + 10}
	assert.Equal(t, DefaultChartWidth+10, set.GetWidth())
}

func TestChartGetHeight(t *testing.T) {
	t.Parallel()

	unset := Chart{}
	assert.Equal(t, DefaultChartHeight, unset.GetHeight())

	set := Chart{Height: DefaultChartHeight + 10}
	assert.Equal(t, DefaultChartHeight+10, set.GetHeight())
}

func TestChartGetRanges(t *testing.T) {
	t.Parallel()

	c := Chart{
		Series: []Series{
			ContinuousSeries{
				XValues: []float64{-2.0, -1.0, 0, 1.0, 2.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 4.5},
			},
			ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{-2.1, -1.0, 0, 1.0, 2.0},
			},
			ContinuousSeries{
				YAxis:   YAxisSecondary,
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{10.0, 11.0, 12.0, 13.0, 14.0},
			},
		},
	}

	xrange, yrange, yrangeAlt := c.getRanges()
	assert.InDelta(t, -2.0, xrange.GetMin(), 0)
	assert.InDelta(t, 5.0, xrange.GetMax(), 0)

	assert.InDelta(t, -2.1, yrange.GetMin(), 0)
	assert.InDelta(t, 4.5, yrange.GetMax(), 0)

	assert.InDelta(t, 10.0, yrangeAlt.GetMin(), 0)
	assert.InDelta(t, 14.0, yrangeAlt.GetMax(), 0)

	cSet := Chart{
		XAxis: XAxis{
			Range: &ContinuousRange{Min: 9.8, Max: 19.8},
		},
		YAxis: YAxis{
			Range: &ContinuousRange{Min: 9.9, Max: 19.9},
		},
		YAxisSecondary: YAxis{
			Range: &ContinuousRange{Min: 9.7, Max: 19.7},
		},
		Series: []Series{
			ContinuousSeries{
				XValues: []float64{-2.0, -1.0, 0, 1.0, 2.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 4.5},
			},
			ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{-2.1, -1.0, 0, 1.0, 2.0},
			},
			ContinuousSeries{
				YAxis:   YAxisSecondary,
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{10.0, 11.0, 12.0, 13.0, 14.0},
			},
		},
	}

	xr2, yr2, yra2 := cSet.getRanges()
	assert.InDelta(t, 9.8, xr2.GetMin(), 0)
	assert.InDelta(t, 19.8, xr2.GetMax(), 0)

	assert.InDelta(t, 9.9, yr2.GetMin(), 0)
	assert.InDelta(t, 19.9, yr2.GetMax(), 0)

	assert.InDelta(t, 9.7, yra2.GetMin(), 0)
	assert.InDelta(t, 19.7, yra2.GetMax(), 0)
}

func TestChartGetRangesUseTicks(t *testing.T) {
	t.Parallel()

	// this test asserts that ticks should supercede manual ranges when generating the overall ranges.

	c := Chart{
		YAxis: YAxis{
			Ticks: []Tick{
				{0.0, "Zero"},
				{1.0, "1.0"},
				{2.0, "2.0"},
				{3.0, "3.0"},
				{4.0, "4.0"},
				{5.0, "Five"},
			},
			Range: &ContinuousRange{
				Min: -5.0,
				Max: 5.0,
			},
		},
		Series: []Series{
			ContinuousSeries{
				XValues: []float64{-2.0, -1.0, 0, 1.0, 2.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 4.5},
			},
		},
	}

	xr, yr, yar := c.getRanges()
	assert.InDelta(t, -2.0, xr.GetMin(), 0)
	assert.InDelta(t, 2.0, xr.GetMax(), 0)
	assert.InDelta(t, 0.0, yr.GetMin(), 0)
	assert.InDelta(t, 5.0, yr.GetMax(), 0)
	assert.True(t, yar.IsZero(), yar.String())
}

func TestChartGetRangesUseUserRanges(t *testing.T) {
	t.Parallel()

	c := Chart{
		YAxis: YAxis{
			Range: &ContinuousRange{
				Min: -5.0,
				Max: 5.0,
			},
		},
		Series: []Series{
			ContinuousSeries{
				XValues: []float64{-2.0, -1.0, 0, 1.0, 2.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 4.5},
			},
		},
	}

	xr, yr, yar := c.getRanges()
	assert.InDelta(t, -2.0, xr.GetMin(), 0)
	assert.InDelta(t, 2.0, xr.GetMax(), 0)
	assert.InDelta(t, -5.0, yr.GetMin(), 0)
	assert.InDelta(t, 5.0, yr.GetMax(), 0)
	assert.True(t, yar.IsZero(), yar.String())
}

func TestChartGetBackgroundStyle(t *testing.T) {
	t.Parallel()

	c := Chart{
		Background: Style{
			FillColor: drawing.ColorBlack,
		},
	}

	bs := c.getBackgroundStyle()
	assert.Equal(t, bs.FillColor.String(), drawing.ColorBlack.String())
}

func TestChartGetCanvasStyle(t *testing.T) {
	c := Chart{
		Canvas: Style{
			FillColor: drawing.ColorBlack,
		},
	}

	bs := c.getCanvasStyle()
	assert.Equal(t, bs.FillColor.String(), drawing.ColorBlack.String())
}

func TestChartGetDefaultCanvasBox(t *testing.T) {
	t.Parallel()

	c := Chart{}
	canvasBoxDefault := c.getDefaultCanvasBox()
	assert.False(t, canvasBoxDefault.IsZero())
	assert.Equal(t, DefaultBackgroundPadding.Top, canvasBoxDefault.Top)
	assert.Equal(t, DefaultBackgroundPadding.Left, canvasBoxDefault.Left)
	assert.Equal(t, c.GetWidth()-DefaultBackgroundPadding.Right, canvasBoxDefault.Right)
	assert.Equal(t, c.GetHeight()-DefaultBackgroundPadding.Bottom, canvasBoxDefault.Bottom)

	custom := Chart{
		Background: Style{
			Padding: Box{
				Top:    DefaultBackgroundPadding.Top + 1,
				Left:   DefaultBackgroundPadding.Left + 1,
				Right:  DefaultBackgroundPadding.Right + 1,
				Bottom: DefaultBackgroundPadding.Bottom + 1,
			},
		},
	}
	canvasBoxCustom := custom.getDefaultCanvasBox()
	assert.False(t, canvasBoxCustom.IsZero())
	assert.Equal(t, DefaultBackgroundPadding.Top+1, canvasBoxCustom.Top)
	assert.Equal(t, DefaultBackgroundPadding.Left+1, canvasBoxCustom.Left)
	assert.Equal(t, c.GetWidth()-(DefaultBackgroundPadding.Right+1), canvasBoxCustom.Right)
	assert.Equal(t, c.GetHeight()-(DefaultBackgroundPadding.Bottom+1), canvasBoxCustom.Bottom)
}

func TestChartGetValueFormatters(t *testing.T) {
	t.Parallel()

	c := Chart{
		Series: []Series{
			ContinuousSeries{
				XValues: []float64{-2.0, -1.0, 0, 1.0, 2.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 4.5},
			},
			ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{-2.1, -1.0, 0, 1.0, 2.0},
			},
			ContinuousSeries{
				YAxis:   YAxisSecondary,
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{10.0, 11.0, 12.0, 13.0, 14.0},
			},
		},
	}

	dxf, dyf, dyaf := c.getValueFormatters()
	assert.NotNil(t, dxf)
	assert.NotNil(t, dyf)
	assert.NotNil(t, dyaf)
}

func TestChartHasAxes(t *testing.T) {
	t.Parallel()

	assert.True(t, Chart{}.hasAxes())
	assert.False(t, Chart{XAxis: XAxis{Style: Hidden()}, YAxis: YAxis{Style: Hidden()}, YAxisSecondary: YAxis{Style: Hidden()}}.hasAxes())

	x := Chart{
		XAxis: XAxis{
			Style: Hidden(),
		},
		YAxis: YAxis{
			Style: Shown(),
		},
		YAxisSecondary: YAxis{
			Style: Hidden(),
		},
	}
	assert.True(t, x.hasAxes())

	y := Chart{
		XAxis: XAxis{
			Style: Shown(),
		},
		YAxis: YAxis{
			Style: Hidden(),
		},
		YAxisSecondary: YAxis{
			Style: Hidden(),
		},
	}
	assert.True(t, y.hasAxes())

	ya := Chart{
		XAxis: XAxis{
			Style: Hidden(),
		},
		YAxis: YAxis{
			Style: Hidden(),
		},
		YAxisSecondary: YAxis{
			Style: Shown(),
		},
	}
	assert.True(t, ya.hasAxes())
}

func TestChartGetAxesTicks(t *testing.T) {
	t.Parallel()

	r := PNG(1024, 1024)

	c := Chart{
		XAxis: XAxis{
			Range: &ContinuousRange{Min: 9.8, Max: 19.8},
		},
		YAxis: YAxis{
			Range: &ContinuousRange{Min: 9.9, Max: 19.9},
		},
		YAxisSecondary: YAxis{
			Range: &ContinuousRange{Min: 9.7, Max: 19.7},
		},
	}
	xr, yr, yar := c.getRanges()

	xt, yt, yat := c.getAxesTicks(r, xr, yr, yar, FloatValueFormatter, FloatValueFormatter, FloatValueFormatter)
	assert.NotEmpty(t, xt)
	assert.NotEmpty(t, yt)
	assert.NotEmpty(t, yat)
}

func TestChartSingleSeries(t *testing.T) {
	t.Parallel()

	now := time.Now()
	c := Chart{
		Title:  "Hello!",
		Width:  1024,
		Height: 400,
		YAxis: YAxis{
			Range: &ContinuousRange{
				Min: 0.0,
				Max: 4.0,
			},
		},
		Series: []Series{
			TimeSeries{
				Name:    "goog",
				XValues: []time.Time{now.AddDate(0, 0, -3), now.AddDate(0, 0, -2), now.AddDate(0, 0, -1)},
				YValues: []float64{1.0, 2.0, 3.0},
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	require.NoError(t, c.Render(PNG, buffer))
	assert.NotEmpty(t, buffer.Bytes())
}

func TestChartRegressionBadRanges(t *testing.T) {
	t.Parallel()

	c := Chart{
		Series: []Series{
			ContinuousSeries{
				XValues: []float64{math.Inf(1), math.Inf(1), math.Inf(1), math.Inf(1), math.Inf(1)},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 4.5},
			},
		},
	}
	buffer := bytes.NewBuffer([]byte{})
	require.Error(t, c.Render(PNG, buffer))
}

func TestChartRegressionBadRangesByUser(t *testing.T) {
	t.Parallel()

	c := Chart{
		YAxis: YAxis{
			Range: &ContinuousRange{
				Min: math.Inf(-1),
				Max: math.Inf(1), // this could really happen? eh.
			},
		},
		Series: []Series{
			ContinuousSeries{
				XValues: LinearRange(1.0, 10.0),
				YValues: LinearRange(1.0, 10.0),
			},
		},
	}
	buffer := bytes.NewBuffer([]byte{})
	require.Error(t, c.Render(PNG, buffer))
}

func TestChartValidatesSeries(t *testing.T) {
	t.Parallel()

	c := Chart{
		Series: []Series{
			ContinuousSeries{
				XValues: LinearRange(1.0, 10.0),
				YValues: LinearRange(1.0, 10.0),
			},
		},
	}

	require.NoError(t, c.validateSeries())

	c = Chart{
		Series: []Series{
			ContinuousSeries{
				XValues: LinearRange(1.0, 10.0),
			},
		},
	}

	require.Error(t, c.validateSeries())
}

func TestChartCheckRanges(t *testing.T) {
	t.Parallel()

	c := Chart{
		Series: []Series{
			ContinuousSeries{
				XValues: []float64{1.0, 2.0},
				YValues: []float64{3.10, 3.14},
			},
		},
	}

	xr, yr, yra := c.getRanges()
	require.NoError(t, c.checkRanges(xr, yr, yra))
}

func TestChartCheckRangesWithRanges(t *testing.T) {
	t.Parallel()

	c := Chart{
		XAxis: XAxis{
			Range: &ContinuousRange{
				Min: 0,
				Max: 10,
			},
		},
		YAxis: YAxis{
			Range: &ContinuousRange{
				Min: 0,
				Max: 5,
			},
		},
		Series: []Series{
			ContinuousSeries{
				XValues: []float64{1.0, 2.0},
				YValues: []float64{3.14, 3.14},
			},
		},
	}

	xr, yr, yra := c.getRanges()
	require.NoError(t, c.checkRanges(xr, yr, yra))
}

func at(i image.Image, x, y int) drawing.Color {
	return drawing.ColorFromAlphaMixedRGBA(i.At(x, y).RGBA())
}

func TestChartE2ELine(t *testing.T) {
	t.Parallel()

	c := Chart{
		Height:         50,
		Width:          50,
		TitleStyle:     Hidden(),
		XAxis:          HideXAxis(),
		YAxis:          HideYAxis(),
		YAxisSecondary: HideYAxis(),
		Canvas: Style{
			Padding: BoxZero,
		},
		Background: Style{
			Padding: BoxZero,
		},
		Series: []Series{
			ContinuousSeries{
				XValues: LinearRangeWithStep(0, 4, 1),
				YValues: LinearRangeWithStep(0, 4, 1),
			},
		},
	}

	var buffer = &bytes.Buffer{}
	require.NoError(t, c.Render(PNG, buffer))

	// do color tests ...

	i, err := png.Decode(buffer)
	require.NoError(t, err)

	// test the bottom and top of the line
	assert.Equal(t, drawing.ColorWhite, at(i, 0, 0))
	assert.Equal(t, drawing.ColorWhite, at(i, 49, 49))

	// test a line mid-point
	defaultSeriesColor := GetDefaultColor(0)
	assert.Equal(t, defaultSeriesColor, at(i, 0, 49))
	assert.Equal(t, defaultSeriesColor, at(i, 49, 0))
	assert.Equal(t, drawing.ColorFromHex("bddbf6"), at(i, 24, 24))
}

func TestChartE2ELineWithFill(t *testing.T) {
	t.Parallel()

	c := Chart{
		Height: 50,
		Width:  50,
		Canvas: Style{
			Padding: BoxZero,
		},
		Background: Style{
			Padding: BoxZero,
		},
		TitleStyle:     Hidden(),
		XAxis:          HideXAxis(),
		YAxis:          HideYAxis(),
		YAxisSecondary: HideYAxis(),
		Series: []Series{
			ContinuousSeries{
				Style: Style{
					StrokeColor: drawing.ColorBlue,
					FillColor:   drawing.ColorRed,
				},
				XValues: LinearRangeWithStep(0, 4, 1),
				YValues: LinearRangeWithStep(0, 4, 1),
			},
		},
	}

	assert.Len(t, c.Series[0].(ContinuousSeries).XValues, 5)
	assert.Len(t, c.Series[0].(ContinuousSeries).YValues, 5)

	var buffer = &bytes.Buffer{}
	require.NoError(t, c.Render(PNG, buffer))

	i, err := png.Decode(buffer)
	require.NoError(t, err)

	// test the bottom and top of the line
	assert.Equal(t, drawing.ColorWhite, at(i, 0, 0))
	assert.Equal(t, drawing.ColorRed, at(i, 49, 49))

	// test a line mid-point
	defaultSeriesColor := drawing.ColorBlue
	assert.Equal(t, defaultSeriesColor, at(i, 0, 49))
	assert.Equal(t, defaultSeriesColor, at(i, 49, 0))
}

func Test_Chart_cve(t *testing.T) {
	t.Parallel()

	poc := StackedBarChart{
		Title: "poc",
		Bars: []StackedBar{
			{
				Name: "11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
				Values: []Value{
					{Value: 1, Label: "infinite"},
					{Value: 1, Label: "loop"},
				},
			},
		},
	}

	var imgContent bytes.Buffer
	err := poc.Render(PNG, &imgContent)
	assert.Error(t, err)
}
