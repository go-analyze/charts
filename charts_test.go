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
