package chartdraw

import (
	"bytes"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBarChartRender(t *testing.T) {
	t.Parallel()

	bc := BarChart{
		Width: 1024,
		Title: "Test Title",
		Bars: []Value{
			{Value: 1.0, Label: "One"},
			{Value: 2.0, Label: "Two"},
			{Value: 3.0, Label: "Three"},
			{Value: 4.0, Label: "Four"},
			{Value: 5.0, Label: "Five"},
		},
	}

	buf := bytes.NewBuffer([]byte{})
	require.NoError(t, bc.Render(PNG, buf))
	assert.NotZero(t, buf.Len())
}

func TestBarChartRenderZero(t *testing.T) {
	t.Parallel()

	bc := BarChart{
		Width: 1024,
		Title: "Test Title",
		Bars: []Value{
			{Value: 0.0, Label: "One"},
			{Value: 0.0, Label: "Two"},
		},
	}

	buf := bytes.NewBuffer([]byte{})
	require.Error(t, bc.Render(PNG, buf))
}

func TestBarChartProps(t *testing.T) {
	t.Parallel()

	bc := BarChart{}

	assert.Equal(t, DefaultDPI, bc.GetDPI())
	bc.DPI = 100
	assert.Equal(t, float64(100), bc.GetDPI())

	assert.Nil(t, bc.GetFont())
	bc.Font = GetDefaultFont()
	assert.NotNil(t, bc.GetFont())

	assert.Equal(t, DefaultChartWidth, bc.GetWidth())
	bc.Width = DefaultChartWidth - 1
	assert.Equal(t, DefaultChartWidth-1, bc.GetWidth())

	assert.Equal(t, DefaultChartHeight, bc.GetHeight())
	bc.Height = DefaultChartHeight - 1
	assert.Equal(t, DefaultChartHeight-1, bc.GetHeight())

	assert.Equal(t, DefaultBarSpacing, bc.GetBarSpacing())
	bc.BarSpacing = 150
	assert.Equal(t, 150, bc.GetBarSpacing())

	assert.Equal(t, DefaultBarWidth, bc.GetBarWidth())
	bc.BarWidth = 75
	assert.Equal(t, 75, bc.GetBarWidth())
}

func TestBarChartRenderNoBars(t *testing.T) {
	t.Parallel()

	bc := BarChart{}
	require.Error(t, bc.Render(PNG, bytes.NewBuffer([]byte{})))
}

func TestBarChartGetRanges(t *testing.T) {
	t.Parallel()

	bc := BarChart{}

	yr := bc.getRanges()
	assert.NotNil(t, yr)
	assert.False(t, yr.IsZero())

	assert.Equal(t, -math.MaxFloat64, yr.GetMax())
	assert.Equal(t, math.MaxFloat64, yr.GetMin())
}

func TestBarChartGetRangesBarsMinMax(t *testing.T) {
	t.Parallel()

	bc := BarChart{
		Bars: []Value{
			{Value: 1.0},
			{Value: 10.0},
		},
	}

	yr := bc.getRanges()
	assert.NotNil(t, yr)
	assert.False(t, yr.IsZero())

	assert.Equal(t, float64(10), yr.GetMax())
	assert.Equal(t, float64(1), yr.GetMin())
}

func TestBarChartGetRangesMinMax(t *testing.T) {
	t.Parallel()

	bc := BarChart{
		YAxis: YAxis{
			Range: &ContinuousRange{
				Min: 5.0,
				Max: 15.0,
			},
			Ticks: []Tick{
				{Value: 7.0, Label: "Foo"},
				{Value: 11.0, Label: "Foo2"},
			},
		},
		Bars: []Value{
			{Value: 1.0},
			{Value: 10.0},
		},
	}

	yr := bc.getRanges()
	assert.NotNil(t, yr)
	assert.False(t, yr.IsZero())

	assert.Equal(t, float64(15), yr.GetMax())
	assert.Equal(t, float64(5), yr.GetMin())
}

func TestBarChartGetRangesTicksMinMax(t *testing.T) {
	t.Parallel()

	bc := BarChart{
		YAxis: YAxis{
			Ticks: []Tick{
				{Value: 7.0, Label: "Foo"},
				{Value: 11.0, Label: "Foo2"},
			},
		},
		Bars: []Value{
			{Value: 1.0},
			{Value: 10.0},
		},
	}

	yr := bc.getRanges()
	assert.NotNil(t, yr)
	assert.False(t, yr.IsZero())

	assert.Equal(t, float64(11), yr.GetMax())
	assert.Equal(t, float64(7), yr.GetMin())
}

func TestBarChartHasAxes(t *testing.T) {
	t.Parallel()

	bc := BarChart{}
	assert.True(t, bc.hasAxes())
	bc.YAxis = YAxis{
		Style: Hidden(),
	}
	assert.False(t, bc.hasAxes())
}

func TestBarChartGetDefaultCanvasBox(t *testing.T) {
	t.Parallel()

	bc := BarChart{}
	b := bc.getDefaultCanvasBox()
	assert.False(t, b.IsZero())
}

func TestBarChartSetRangeDomains(t *testing.T) {
	t.Parallel()

	bc := BarChart{}
	cb := bc.box()
	yr := bc.getRanges()
	yr2 := bc.setRangeDomains(cb, yr)
	assert.NotZero(t, yr2.GetDomain())
}

func TestBarChartGetValueFormatters(t *testing.T) {
	t.Parallel()

	bc := BarChart{}
	vf := bc.getValueFormatters()
	assert.NotNil(t, vf)
	assert.Equal(t, "1234.00", vf(1234.0))

	bc.YAxis.ValueFormatter = func(_ interface{}) string { return "test" }
	assert.Equal(t, "test", bc.getValueFormatters()(1234))
}

func TestBarChartGetAxesTicks(t *testing.T) {
	t.Parallel()

	bc := BarChart{
		Bars: []Value{
			{Value: 1.0},
			{Value: 2.0},
			{Value: 3.0},
		},
	}

	r := PNG(128, 128)
	yr := bc.getRanges()
	yf := bc.getValueFormatters()

	bc.YAxis.Style.Hidden = true
	ticks := bc.getAxesTicks(r, yr, yf)
	assert.Empty(t, ticks)

	bc.YAxis.Style.Hidden = false
	ticks = bc.getAxesTicks(r, yr, yf)
	assert.Len(t, ticks, 2)
}

func TestBarChartCalculateEffectiveBarSpacing(t *testing.T) {
	t.Parallel()

	bc := BarChart{
		Width:    1024,
		BarWidth: 10,
		Bars: []Value{
			{Value: 1.0, Label: "One"},
			{Value: 2.0, Label: "Two"},
			{Value: 3.0, Label: "Three"},
			{Value: 4.0, Label: "Four"},
			{Value: 5.0, Label: "Five"},
		},
	}

	spacing := bc.calculateEffectiveBarSpacing(bc.box())
	assert.NotZero(t, spacing)

	bc.BarWidth = 250
	spacing = bc.calculateEffectiveBarSpacing(bc.box())
	assert.Zero(t, spacing)
}

func TestBarChartCalculateEffectiveBarWidth(t *testing.T) {
	t.Parallel()

	bc := BarChart{
		Width:    1024,
		BarWidth: 10,
		Bars: []Value{
			{Value: 1.0, Label: "One"},
			{Value: 2.0, Label: "Two"},
			{Value: 3.0, Label: "Three"},
			{Value: 4.0, Label: "Four"},
			{Value: 5.0, Label: "Five"},
		},
	}

	cb := bc.box()

	spacing := bc.calculateEffectiveBarSpacing(bc.box())
	assert.NotZero(t, spacing)

	barWidth := bc.calculateEffectiveBarWidth(bc.box(), spacing)
	assert.Equal(t, 10, barWidth)

	bc.BarWidth = 250
	spacing = bc.calculateEffectiveBarSpacing(bc.box())
	assert.Zero(t, spacing)
	barWidth = bc.calculateEffectiveBarWidth(bc.box(), spacing)
	assert.Equal(t, 199, barWidth)

	assert.Equal(t, cb.Width()+1, bc.calculateTotalBarWidth(barWidth, spacing))

	bw, bs, total := bc.calculateScaledTotalWidth(cb)
	assert.Equal(t, spacing, bs)
	assert.Equal(t, barWidth, bw)
	assert.Equal(t, cb.Width()+1, total)
}

func TestBarChatGetTitleFontSize(t *testing.T) {
	t.Parallel()

	size := BarChart{Width: 2049, Height: 2049}.getTitleFontSize()
	assert.Equal(t, float64(48), size)
	size = BarChart{Width: 1025, Height: 1025}.getTitleFontSize()
	assert.Equal(t, float64(24), size)
	size = BarChart{Width: 513, Height: 513}.getTitleFontSize()
	assert.Equal(t, float64(18), size)
	size = BarChart{Width: 257, Height: 257}.getTitleFontSize()
	assert.Equal(t, float64(12), size)
	size = BarChart{Width: 128, Height: 128}.getTitleFontSize()
	assert.Equal(t, float64(10), size)
}
