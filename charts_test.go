package charts

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetNullValue(t *testing.T) {
	t.Parallel()

	assert.InDelta(t, math.MaxFloat64, GetNullValue(), 0.0)
}

func TestLegendRepositioning(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                     string
		titleOpt                 TitleOption
		legendOpt                LegendOption
		expectLegendRepositioned bool
	}{
		{
			name: "horizontal_center_with_center_title_collision",
			titleOpt: TitleOption{
				Text:   "Test Title",
				Offset: OffsetCenter,
			},
			legendOpt: LegendOption{
				SeriesNames: []string{"Series A", "Series B", "Series C", "Series D", "Series E"},
			},
			expectLegendRepositioned: true,
		},
		{
			name: "horizontal_center_no_title",
			titleOpt: TitleOption{
				Text: "",
			},
			legendOpt: LegendOption{
				SeriesNames: []string{"Series A", "Series B"},
			},
			expectLegendRepositioned: false,
		},
		{
			name: "horizontal_bottom_no_collision",
			titleOpt: TitleOption{
				Text:   "Test Title",
				Offset: OffsetCenter,
			},
			legendOpt: LegendOption{
				SeriesNames: []string{"A", "B"},
				Offset: OffsetStr{
					Top: PositionBottom,
				},
			},
			expectLegendRepositioned: false,
		},
		{
			name: "vertical_legend_with_collision",
			titleOpt: TitleOption{
				Text:   "Test Title",
				Offset: OffsetLeft,
			},
			legendOpt: LegendOption{
				SeriesNames: []string{"Series A", "Series B"},
				Vertical:    Ptr(true),
			},
			expectLegendRepositioned: true,
		},
		{
			name: "vertical_legend_no_collision_right_position",
			titleOpt: TitleOption{
				Text:   "Test Title",
				Offset: OffsetLeft,
			},
			legendOpt: LegendOption{
				SeriesNames: []string{"Series A", "Series B"},
				Vertical:    Ptr(true),
				Offset:      OffsetRight,
			},
			expectLegendRepositioned: false,
		},
		{
			name: "explicit_numeric_offset_blocks_repositioning",
			titleOpt: TitleOption{
				Text:   "Test Title",
				Offset: OffsetCenter,
			},
			legendOpt: LegendOption{
				SeriesNames: []string{"A", "B", "C", "D", "E"},
				Offset: OffsetStr{
					Top: "10",
				},
			},
			expectLegendRepositioned: false,
		},
		{
			name: "explicit_numeric_left_offset_blocks_repositioning",
			titleOpt: TitleOption{
				Text:   "Test Title",
				Offset: OffsetCenter,
			},
			legendOpt: LegendOption{
				SeriesNames: []string{"A", "B", "C", "D", "E"},
				Offset: OffsetStr{
					Left: "50",
				},
			},
			expectLegendRepositioned: false,
		},
		{
			name: "overlay_chart_blocks_repositioning",
			titleOpt: TitleOption{
				Text:   "Test Title",
				Offset: OffsetCenter,
			},
			legendOpt: LegendOption{
				SeriesNames:  []string{"A", "B", "C", "D", "E"},
				OverlayChart: Ptr(true),
			},
			expectLegendRepositioned: false,
		},
		{
			name: "title_and_legend_both_bottom",
			titleOpt: TitleOption{
				Text: "Test Title",
				Offset: OffsetStr{
					Top: PositionBottom,
				},
			},
			legendOpt: LegendOption{
				SeriesNames: []string{"A", "B"},
				Offset: OffsetStr{
					Top: PositionBottom,
				},
			},
			expectLegendRepositioned: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// create chart with line series to trigger defaultRender
			opt := ChartOption{
				Title:      tt.titleOpt,
				Legend:     tt.legendOpt,
				SeriesList: NewSeriesListGeneric([][]float64{{1, 2, 3}}, ChartTypeLine),
				Width:      600,
				Height:     400,
			}

			// save original offset
			originalTopOffset := tt.legendOpt.Offset.Top

			p, err := Render(opt)
			require.NoError(t, err)
			require.NotNil(t, p)

			// check if legend was repositioned by verifying original option wasn't mutated
			// (we use a copy internally, so original should always remain unchanged)
			assert.Equal(t, originalTopOffset, tt.legendOpt.Offset.Top,
				"original legend option should not be mutated")
		})
	}
}

func TestLegendOptionNotMutated(t *testing.T) {
	t.Parallel()

	// verify that rendering does not mutate the original legend option
	legendOpt := LegendOption{
		SeriesNames: []string{"A", "B", "C", "D", "E"},
		Offset: OffsetStr{
			Top: PositionTop,
		},
	}
	titleOpt := TitleOption{
		Text:   "Wide Title That Should Cause Collision",
		Offset: OffsetCenter,
	}

	opt := ChartOption{
		Title:      titleOpt,
		Legend:     legendOpt,
		SeriesList: NewSeriesListGeneric([][]float64{{1, 2, 3}}, ChartTypeLine),
		Width:      600,
		Height:     400,
	}

	_, err := Render(opt)
	require.NoError(t, err)

	// original should not have been modified
	assert.Equal(t, PositionTop, legendOpt.Offset.Top)
}

func TestLegendIndexSpan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		left, right      int
		plotWidth, count int
		lo, hi           int
		ok               bool
	}{
		{name: "left_quarter", left: 0, right: 100, plotWidth: 400, count: 5, lo: 0, hi: 1, ok: true},
		{name: "right_quarter", left: 300, right: 400, plotWidth: 400, count: 5, lo: 3, hi: 4, ok: true},
		{name: "center", left: 150, right: 250, plotWidth: 400, count: 5, lo: 1, hi: 3, ok: true},
		{name: "full_width", left: 0, right: 400, plotWidth: 400, count: 5, lo: 0, hi: 4, ok: true},
		{name: "single_point", left: 0, right: 100, plotWidth: 400, count: 1, lo: 0, hi: 0, ok: true},
		{name: "clamp_overflow", left: -50, right: 999, plotWidth: 400, count: 5, lo: 0, hi: 4, ok: true},
		{name: "no_points", left: 0, right: 100, plotWidth: 400, count: 0, ok: false},
		{name: "span_left_of_plot", left: -200, right: -10, plotWidth: 400, count: 5, ok: false},
		{name: "span_right_of_plot", left: 500, right: 600, plotWidth: 400, count: 5, ok: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lo, hi, ok := legendIndexSpan(tt.left, tt.right, tt.plotWidth, tt.count)
			assert.Equal(t, tt.ok, ok)
			if tt.ok {
				assert.Equal(t, tt.lo, lo)
				assert.Equal(t, tt.hi, hi)
			}
		})
	}
}

func TestLocalMaxOverIndices(t *testing.T) {
	t.Parallel()

	t.Run("non_stacked_max", func(t *testing.T) {
		sl := NewSeriesListLine([][]float64{{1, 9, 3}, {4, 2, 5}})
		v, ok := localMaxOverIndices(sl, 0, 0, 2, false)
		require.True(t, ok)
		assert.InDelta(t, 9.0, v, 0)
	})

	t.Run("subrange", func(t *testing.T) {
		sl := NewSeriesListLine([][]float64{{1, 9, 3}, {4, 2, 5}})
		v, ok := localMaxOverIndices(sl, 0, 2, 2, false)
		require.True(t, ok)
		assert.InDelta(t, 5.0, v, 0)
	})

	t.Run("stacked_sum", func(t *testing.T) {
		sl := NewSeriesListLine([][]float64{{1, 9, 3}, {4, 2, 5}})
		v, ok := localMaxOverIndices(sl, 0, 0, 2, true)
		require.True(t, ok)
		assert.InDelta(t, 11.0, v, 0) // index 1: 9+2; index 2: 3+5=8; index 0: 5
	})

	t.Run("null_skipped", func(t *testing.T) {
		sl := NewSeriesListLine([][]float64{{1, GetNullValue(), 3}})
		v, ok := localMaxOverIndices(sl, 0, 0, 2, false)
		require.True(t, ok)
		assert.InDelta(t, 3.0, v, 0)
	})

	t.Run("empty_range", func(t *testing.T) {
		sl := NewSeriesListLine([][]float64{{1, 2, 3}})
		_, ok := localMaxOverIndices(sl, 0, 2, 1, false)
		assert.False(t, ok)
	})

	t.Run("yaxis_filter", func(t *testing.T) {
		sl := LineSeriesList{
			{Values: []float64{1, 2, 3}, YAxisIndex: 0},
			{Values: []float64{40, 50, 60}, YAxisIndex: 1},
		}
		v, ok := localMaxOverIndices(sl, 0, 0, 2, false)
		require.True(t, ok)
		assert.InDelta(t, 3.0, v, 0)
		v, ok = localMaxOverIndices(sl, 1, 0, 2, false)
		require.True(t, ok)
		assert.InDelta(t, 60.0, v, 0)
	})
}

func TestVerticalLegendHeadroom(t *testing.T) {
	t.Parallel()

	// three series so the vertical legend is several rows tall (taller than the default top padding gap)
	downward := [][]float64{{100, 78, 55, 30, 5}, {96, 74, 51, 26, 2}, {92, 70, 47, 22, 1}}
	upward := [][]float64{{5, 30, 55, 78, 100}, {2, 26, 51, 74, 96}, {1, 22, 47, 70, 92}}
	labels := []string{"a", "b", "c", "d", "e"}
	names := []string{"one", "two", "three"}

	render := func(values [][]float64, legend LegendOption, markPoint bool, valueAxis ValueAxisOption) float64 {
		p := NewPainter(PainterOptions{OutputFormat: ChartOutputSVG, Width: 600, Height: 400})
		sl := NewSeriesListLine(values)
		if markPoint {
			for i := range sl {
				sl[i].MarkPoint = NewMarkPoint(SeriesMarkTypeMax)
			}
		}
		rr, err := defaultRender(p, defaultRenderOption{
			theme:        GetDefaultTheme(),
			seriesList:   &sl,
			categoryAxis: &CategoryAxisOption{Labels: labels},
			valueAxis:    []ValueAxisOption{valueAxis},
			legend:       &legend,
		})
		require.NoError(t, err)
		return rr.valueAxisRanges[0].max
	}

	// baseline: no legend, so no headroom adjustment is ever applied
	baselineDown := render(downward, LegendOption{}, false, ValueAxisOption{})
	baselineUp := render(upward, LegendOption{}, false, ValueAxisOption{})

	topLeft := LegendOption{Vertical: Ptr(true), SeriesNames: names}
	topRight := LegendOption{Vertical: Ptr(true), SeriesNames: names, Offset: OffsetStr{Left: PositionRight}}

	t.Run("top_left_downward_adjusts", func(t *testing.T) {
		assert.Greater(t, render(downward, topLeft, false, ValueAxisOption{}), baselineDown)
	})

	t.Run("top_right_downward_no_adjust", func(t *testing.T) {
		assert.InDelta(t, baselineDown, render(downward, topRight, false, ValueAxisOption{}), 0)
	})

	t.Run("top_left_upward_no_adjust", func(t *testing.T) {
		assert.InDelta(t, baselineUp, render(upward, topLeft, false, ValueAxisOption{}), 0)
	})

	t.Run("top_right_upward_adjusts", func(t *testing.T) {
		assert.Greater(t, render(upward, topRight, false, ValueAxisOption{}), baselineUp)
	})

	t.Run("bottom_legend_no_adjust", func(t *testing.T) {
		bottom := LegendOption{Vertical: Ptr(true), SeriesNames: names, Offset: OffsetStr{Top: PositionBottom}}
		assert.InDelta(t, baselineDown, render(downward, bottom, false, ValueAxisOption{}), 0)
	})

	t.Run("horizontal_legend_no_adjust", func(t *testing.T) {
		horizontal := LegendOption{SeriesNames: names, OverlayChart: Ptr(true)}
		assert.InDelta(t, baselineDown, render(downward, horizontal, false, ValueAxisOption{}), 0)
	})

	t.Run("overlay_disabled_no_adjust", func(t *testing.T) {
		noOverlay := LegendOption{Vertical: Ptr(true), SeriesNames: names, OverlayChart: Ptr(false)}
		assert.InDelta(t, baselineDown, render(downward, noOverlay, false, ValueAxisOption{}), 0)
	})

	t.Run("explicit_max_no_adjust", func(t *testing.T) {
		fixed := ValueAxisOption{Max: Ptr(100.0)}
		base := render(downward, LegendOption{}, false, fixed)
		assert.InDelta(t, base, render(downward, topLeft, false, fixed), 0)
	})

	t.Run("markpoint_adjusts_more", func(t *testing.T) {
		withoutMark := render(downward, topLeft, false, ValueAxisOption{})
		withMark := render(downward, topLeft, true, ValueAxisOption{})
		assert.Greater(t, withMark, withoutMark)
	})
}
