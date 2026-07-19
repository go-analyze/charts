package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeFullLineChartOption() LineChartOption {
	values := [][]float64{
		{120, 132, 101, 134, 90, 230, 210},
		{220, 182, 191, 234, 290, 330, 310},
		{150, 232, 201, 154, 190, 330, 410},
		{320, 332, 301, 334, 390, 330, 320},
		{820, 932, 901, 934, 1290, 1330, 1320},
	}
	return LineChartOption{
		Title: TitleOption{
			Text: "Line",
		},
		Padding: NewBoxEqual(10),
		XAxis: XAxisOption{
			Labels: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
		},
		YAxis: make([]YAxisOption, 1),
		Legend: LegendOption{
			SeriesNames: []string{"Email", "Union Ads", "Video Ads", "Direct", "Search Engine"},
			Symbol:      SymbolDot,
		},
		SeriesList: NewSeriesListLine(values),
	}
}

func makeBasicLineChartOption() LineChartOption {
	values := [][]float64{
		{120, 132, 101, 134, 90, 230, 210},
		{820, 932, 901, 934, 1290, 1330, 1320},
	}
	return LineChartOption{
		Title: TitleOption{
			Text: "Line",
		},
		Padding: NewBoxEqual(10),
		XAxis: XAxisOption{
			Labels: []string{"A", "B", "C", "D", "E", "F", "G"},
		},
		YAxis: make([]YAxisOption, 1),
		Legend: LegendOption{
			SeriesNames: []string{"1", "2"},
			Symbol:      SymbolDot,
		},
		SeriesList: NewSeriesListLine(values),
	}
}

func makeMinimalLineChartLegendOption() LineChartOption {
	opt := makeMinimalLineChartOption()
	opt.Legend = LegendOption{
		SeriesNames: []string{"1", "2"},
	}
	return opt
}

func makeMinimalLineChartOption() LineChartOption {
	values := [][]float64{
		{120, 132, 101, 134, 90, 230, 210},
		{820, 932, 901, 934, 1290, 1330, 1320},
	}
	return LineChartOption{
		Padding: NewBoxEqual(10),
		XAxis: XAxisOption{
			Labels: []string{"1", "2", "3", "4", "5", "6", "7"},
			Show:   Ptr(false),
		},
		Legend: LegendOption{
			Symbol: SymbolDot,
		},
		YAxis:      make([]YAxisOption, 1),
		Symbol:     Symbol{Shape: SymbolNone},
		SeriesList: NewSeriesListLine(values),
	}
}

func makeFullLineChartStackedOption() LineChartOption {
	seriesList := NewSeriesListLine([][]float64{
		{4.9, 23.2, 25.6, 102.6, 142.2, 32.6, 20.0, 3.3},
		{9.0, 26.4, 28.7, 144.6, 122.2, 48.7, 18.8, 2.3},
		{80.0, 40.4, 28.4, 28.8, 24.4, 24.2, 40.8, 80.8},
	}, LineSeriesOption{
		Label: SeriesLabel{
			Show: Ptr(true),
		},
		MarkPoint: NewMarkPoint("min", "max"),
	})
	return LineChartOption{
		Padding:     NewBoxEqual(20),
		SeriesList:  seriesList,
		StackSeries: Ptr(true),
		XAxis: XAxisOption{
			Labels:      []string{"1", "2", "3", "4", "5", "6", "7", "8"},
			BoundaryGap: Ptr(true),
		},
		Legend: LegendOption{
			SeriesNames: []string{"A", "B", "C"},
			Symbol:      SymbolDot,
		},
		YAxis: []YAxisOption{
			{
				RangeValuePaddingScale: Ptr(1.0),
				PreferNiceIntervals:    Ptr(true),
			},
		},
	}
}

func TestNewLineChartOptionWithData(t *testing.T) {
	t.Parallel()

	opt := NewLineChartOptionWithData([][]float64{
		{12, 24},
		{24, 48},
	})

	assert.Len(t, opt.SeriesList, 2)
	assert.Equal(t, ChartTypeLine, opt.SeriesList[0].getType())
	assert.Len(t, opt.YAxis, 1)
	assert.Equal(t, defaultPadding, opt.Padding)

	p := NewPainter(PainterOptions{})
	assert.NoError(t, p.LineChart(opt))
}

func TestLineChart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		themed      bool
		makeOptions func() LineChartOption
		pngCRC      uint32
	}{
		{
			name:        "basic_themed",
			themed:      true,
			makeOptions: makeFullLineChartOption,
			pngCRC:      0x3f0900e9,
		},
		{
			name: "boundary_gap_disable",
			makeOptions: func() LineChartOption {
				opt := makeFullLineChartOption()
				opt.XAxis.BoundaryGap = Ptr(false)
				// disable extras
				opt.Title.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x60869a42,
		},
		{
			name: "boundary_gap_enable",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.XAxis.BoundaryGap = Ptr(true)
				return opt
			},
			pngCRC: 0x269d492f,
		},
		{
			name: "08Y_skip1",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxis = []YAxisOption{
					{
						LabelCount:     8,
						LabelSkipCount: 1,
					},
				}
				return opt
			},
			pngCRC: 0xc4b6c46e,
		},
		{
			name: "09Y_skip1",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxis = []YAxisOption{
					{
						LabelCount:     9,
						LabelSkipCount: 1,
					},
				}
				return opt
			},
			pngCRC: 0x4be7d06b,
		},
		{
			name: "08Y_skip2",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxis = []YAxisOption{
					{
						LabelCount:     8,
						LabelSkipCount: 2,
					},
				}
				return opt
			},
			pngCRC: 0x11756954,
		},
		{
			name: "09Y_skip2",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxis = []YAxisOption{
					{
						LabelCount:     9,
						LabelSkipCount: 2,
					},
				}
				return opt
			},
			pngCRC: 0x40214911,
		},
		{
			name: "10Y_skip2",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxis = []YAxisOption{
					{
						LabelCount:     10,
						LabelSkipCount: 2,
					},
				}
				return opt
			},
			pngCRC: 0x3d920ed,
		},
		{
			name: "08Y_skip3",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxis = []YAxisOption{
					{
						LabelCount:     8,
						LabelSkipCount: 3,
					},
				}
				return opt
			},
			pngCRC: 0x1ac94fa1,
		},
		{
			name: "09Y_skip3",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxis = []YAxisOption{
					{
						LabelCount:     9,
						LabelSkipCount: 3,
					},
				}
				return opt
			},
			pngCRC: 0x2424201c,
		},
		{
			name: "10Y_skip3",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxis = []YAxisOption{
					{
						LabelCount:     10,
						LabelSkipCount: 3,
					},
				}
				return opt
			},
			pngCRC: 0x377f49de,
		},
		{
			name: "11Y_skip3",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxis = []YAxisOption{
					{
						LabelCount:     11,
						LabelSkipCount: 3,
					},
				}
				return opt
			},
			pngCRC: 0x12fda160,
		},
		{
			name: "no_yaxis_split_line",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxis = []YAxisOption{
					{
						LabelCount:    4,
						SplitLineShow: Ptr(false),
					},
				}
				return opt
			},
			pngCRC: 0x3bcddb42,
		},
		{
			name: "yaxis_spine_line_show",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxis = []YAxisOption{
					{
						LabelCount:    4,
						SpineLineShow: Ptr(true),
					},
				}
				return opt
			},
			pngCRC: 0x54b19a85,
		},
		{
			name: "dual_yaxis",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				opt.Theme = GetTheme(ThemeLight)
				opt.SeriesList[1].YAxisIndex = 1
				opt.YAxis = append(opt.YAxis, opt.YAxis[0])
				opt.YAxis[0].Theme = opt.Theme.WithYAxisSeriesColor(0)
				opt.YAxis[1].Theme = opt.Theme.WithYAxisSeriesColor(1)
				// disable extras
				opt.Title.Text = ""
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0xa9a691a1,
		},
		{
			name: "no_nice_interval",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				opt.Theme = GetTheme(ThemeLight)
				opt.SeriesList[1].YAxisIndex = 1
				opt.YAxis = append(opt.YAxis, opt.YAxis[0])
				opt.YAxis[0].PreferNiceIntervals = Ptr(false)
				opt.YAxis[1].PreferNiceIntervals = Ptr(false)
				opt.YAxis[0].Theme = opt.Theme.WithYAxisSeriesColor(0)
				opt.YAxis[1].Theme = opt.Theme.WithYAxisSeriesColor(1)
				// disable extras
				opt.Title.Text = ""
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0xc1b14f19,
		},
		{
			name: "left_nice_interval",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				opt.Theme = GetTheme(ThemeLight)
				opt.SeriesList[1].YAxisIndex = 1
				opt.YAxis = append(opt.YAxis, opt.YAxis[0])
				opt.YAxis[0].PreferNiceIntervals = Ptr(true)
				opt.YAxis[1].PreferNiceIntervals = Ptr(false)
				opt.YAxis[0].Theme = opt.Theme.WithYAxisSeriesColor(0)
				opt.YAxis[1].Theme = opt.Theme.WithYAxisSeriesColor(1)
				// disable extras
				opt.Title.Text = ""
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x8f0883de,
		},
		{
			name: "right_nice_interval",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				opt.Theme = GetTheme(ThemeLight)
				opt.Title.Text = "T"
				opt.SeriesList[1].YAxisIndex = 1
				opt.YAxis = append(opt.YAxis, opt.YAxis[0])
				opt.YAxis[0].PreferNiceIntervals = Ptr(false)
				opt.YAxis[1].PreferNiceIntervals = Ptr(true)
				opt.YAxis[0].Theme = opt.Theme.WithYAxisSeriesColor(0)
				opt.YAxis[1].Theme = opt.Theme.WithYAxisSeriesColor(1)
				// disable extras
				opt.Title.Text = ""
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0xa9a691a1,
		},
		{
			name: "right_yaxis",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				opt.Theme = GetTheme(ThemeLight)
				opt.YAxis[0].Position = PositionRight
				// disable extras
				opt.Title.Text = ""
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0xae898031,
		},
		{
			name: "zero_data",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				values := [][]float64{
					{0, 0, 0, 0, 0, 0, 0},
					{0, 0, 0, 0, 0, 0, 0},
				}
				opt.SeriesList = NewSeriesListLine(values)
				return opt
			},
			pngCRC: 0x58ea4a38,
		},
		{
			name: "tiny_range",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				values := [][]float64{
					{0.1, 0.2, 0.1, 0.2, 0.4, 0.2, 0.1},
					{0.2, 0.4, 0.8, 0.4, 0.2, 0.1, 0.2},
				}
				opt.SeriesList = NewSeriesListLine(values)
				return opt
			},
			pngCRC: 0x75aa23cd,
		},
		{
			name: "custom_font",
			makeOptions: func() LineChartOption {
				opt := makeFullLineChartOption()
				customFont := NewFontStyleWithSize(4.0).WithColor(ColorBlue)
				opt.Legend.FontStyle = customFont
				opt.XAxis.LabelFontStyle = customFont
				opt.Title.FontStyle = customFont
				return opt
			},
			pngCRC: 0xb32a74ec,
		},
		{
			name: "title_offset_center_legend_right",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				opt.Title.Offset = OffsetCenter
				opt.Legend.Offset = OffsetRight
				return opt
			},
			pngCRC: 0x590ab0e5,
		},
		{
			name: "title_offset_right",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				opt.Title.Offset = OffsetRight
				return opt
			},
			pngCRC: 0x7725cb09,
		},
		{
			name: "title_offset_bottom_center",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				opt.Title.Offset = OffsetStr{
					Top:  PositionBottom,
					Left: PositionCenter,
				}
				return opt
			},
			pngCRC: 0x6a854a98,
		},
		{
			name: "legend_offset_bottom",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				opt.Legend.Offset = OffsetStr{
					Top: PositionBottom,
				}
				return opt
			},
			pngCRC: 0x35acfac5,
		},
		{
			name: "legend_padding_top",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				opt.Legend.Padding = NewBoxEqual(50)
				opt.Legend.BorderWidth = 1.0
				// disable extra stuff
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
			pngCRC: 0x5ce8f982,
		},
		{
			name: "legend_padding_bottom",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				opt.Legend.Padding = NewBoxEqual(50)
				opt.Legend.Offset = OffsetStr{
					Top: PositionBottom,
				}
				opt.Legend.BorderWidth = 1.0
				// disable extra stuff
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
			pngCRC: 0x31f4a06f,
		},
		{
			name: "title_and_legend_offset_bottom",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				bottomOffset := OffsetStr{
					Top:  PositionBottom,
					Left: PositionCenter,
				}
				opt.Title.Offset = bottomOffset
				opt.Legend.Offset = bottomOffset
				return opt
			},
			pngCRC: 0xe78f2365,
		},
		{
			name: "vertical_legend_offset_right",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				opt.Legend.Vertical = Ptr(true)
				opt.Legend.Offset = OffsetRight
				return opt
			},
			pngCRC: 0x7a05c1db,
		},
		{
			name: "legend_overlap_chart",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				opt.Title.Show = Ptr(false)
				opt.Legend.Offset = OffsetRight
				opt.Legend.OverlayChart = Ptr(true)
				// disable extra
				opt.Title.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x94c6bb83,
		},
		{
			name: "legend_boxed_offset_bottom",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartLegendOption()
				opt.Legend.Offset = OffsetStr{
					Top: PositionBottom,
				}
				opt.Legend.BorderWidth = 2.0
				return opt
			},
			pngCRC: 0x13dffa63,
		},
		{
			name: "vertical_legend_boxed_offset_right",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartLegendOption()
				opt.Legend.Vertical = Ptr(true)
				opt.Legend.Offset = OffsetRight
				opt.Legend.BorderWidth = 2.0
				return opt
			},
			pngCRC: 0x783d4ef2,
		},
		{
			name: "legend_boxed_overlap_chart",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartLegendOption()
				opt.Title.Show = Ptr(false)
				opt.Legend.Offset = OffsetRight
				opt.Legend.OverlayChart = Ptr(true)
				opt.Legend.BorderWidth = 2.0
				return opt
			},
			pngCRC: 0xccd0fe5b,
		},
		{
			name: "curved_line",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.StrokeSmoothingTension = 0.8
				return opt
			},
			pngCRC: 0x9b90df89,
		},
		{
			name: "line_gap",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.SeriesList[0].Values[3] = GetNullValue()
				return opt
			},
			pngCRC: 0x52ee1b05,
		},
		{
			name: "line_gap_dot",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.Symbol.Shape = SymbolDot
				opt.SeriesList[0].Values[3] = GetNullValue()
				opt.SeriesList[0].Values[5] = GetNullValue()
				return opt
			},
			pngCRC: 0x97150462,
		},
		{
			name: "line_gap_fill_area",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.SeriesList[0].Values[3] = GetNullValue()
				opt.FillArea = Ptr(true)
				return opt
			},
			pngCRC: 0xa5d75beb,
		},
		{
			name: "line_gap_start_fill_area",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.SeriesList[0].Values[0] = GetNullValue()
				opt.SeriesList[0].Values[1] = GetNullValue()
				opt.SeriesList[1].Values[0] = GetNullValue()
				opt.FillArea = Ptr(true)
				return opt
			},
			pngCRC: 0x587cdb17,
		},
		{
			name: "curved_line_gap",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.StrokeSmoothingTension = 0.8
				opt.SeriesList[0].Values[3] = GetNullValue()
				return opt
			},
			pngCRC: 0xce51085e,
		},
		{
			name: "curved_line_gap_fill_area",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.StrokeSmoothingTension = 0.8
				opt.SeriesList[0].Values[3] = GetNullValue()
				opt.FillArea = Ptr(true)
				return opt
			},
			pngCRC: 0xa19ba7e0,
		},
		{
			name: "fill_area",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.FillArea = Ptr(true)
				opt.FillOpacity = 100
				return opt
			},
			pngCRC: 0xa74354ee,
		},
		{
			name: "fill_area_boundary_gap",
			makeOptions: func() LineChartOption {
				opt := makeFullLineChartOption()
				opt.FillArea = Ptr(true)
				opt.FillOpacity = 100
				opt.XAxis.BoundaryGap = Ptr(true)
				// disable extra
				opt.Symbol.Shape = SymbolNone
				opt.Title.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
			pngCRC: 0x3c0b1162,
		},
		{
			name: "fill_area_curved_boundary_gap",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.FillArea = Ptr(true)
				opt.StrokeSmoothingTension = 0.8
				opt.XAxis.BoundaryGap = Ptr(true)
				return opt
			},
			pngCRC: 0x8ec4cbaf,
		},
		{
			name: "fill_area_curved_no_gap",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.FillArea = Ptr(true)
				opt.StrokeSmoothingTension = 0.8
				opt.XAxis.BoundaryGap = Ptr(false)
				return opt
			},
			pngCRC: 0x7007e29a,
		},
		{
			name: "value_formatter",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxis[0].ValueFormatter = func(f float64) string {
					return "f"
				}
				return opt
			},
			pngCRC: 0xa8064a0e,
		},
		{
			name: "mark_line",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.Padding = NewBoxEqual(40)
				opt.YAxis[0].Show = Ptr(false)
				for i := range opt.SeriesList {
					markLine := NewMarkLine("min", "max", "average")
					markLine.ValueFormatter = func(f float64) string {
						return FormatValueHumanizeShort(f, 0, false)
					}
					opt.SeriesList[i].MarkLine = markLine
				}
				return opt
			},
			pngCRC: 0x150eb948,
		},
		{
			name: "mark_point",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxis[0].Show = Ptr(false)
				for i := range opt.SeriesList {
					markPoint := NewMarkPoint("min", "max")
					markPoint.ValueFormatter = func(f float64) string {
						return FormatValueHumanizeShort(f, 0, false)
					}
					opt.SeriesList[i].MarkPoint = markPoint
				}
				return opt
			},
			pngCRC: 0x7ce0855c,
		},
		{
			name: "series_label",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.YAxis[0].Show = Ptr(false)
				for i := range opt.SeriesList {
					opt.SeriesList[i].Label.Show = Ptr(true)
					opt.SeriesList[i].Label.FontStyle = FontStyle{
						FontSize:  12.0,
						Font:      GetDefaultFont(),
						FontColor: ColorBlue,
					}
					opt.SeriesList[i].MarkPoint = NewMarkPoint("min", "max")
					opt.SeriesList[i].Label.ValueFormatter = func(f float64) string {
						return FormatValueHumanizeShort(f, 2, false)
					}
				}
				return opt
			},
			pngCRC: 0xdff58aa8,
		},
		{
			name:        "stack_series",
			makeOptions: makeFullLineChartStackedOption,
			pngCRC:      0x79fb5751,
		},
		{
			name: "stack_series_global_mark_point",
			makeOptions: func() LineChartOption {
				opt := makeFullLineChartStackedOption()
				// add global point configurations
				opt.SeriesList[len(opt.SeriesList)-1].MarkPoint.AddGlobalPoints(SeriesMarkTypeMin, SeriesMarkTypeMax)
				// disable extra stuff
				opt.Padding = NewBox(20, 40, 20, 20)
				opt.Legend.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				return opt
			},
			pngCRC: 0xea99641e,
		},
		{
			name: "stack_series_global_mark_line",
			makeOptions: func() LineChartOption {
				opt := makeFullLineChartStackedOption()
				// add global line configurations
				opt.SeriesList[len(opt.SeriesList)-1].MarkLine.AddGlobalLines(SeriesMarkTypeAverage,
					SeriesMarkTypeMin, SeriesMarkTypeMax)
				// disable extra stuff
				for i := range opt.SeriesList {
					opt.SeriesList[i].MarkPoint = SeriesMarkPoint{}
					opt.SeriesList[i].Label.Show = Ptr(false)
				}
				opt.Padding = NewBox(20, 40, 20, 20)
				opt.Legend.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				return opt
			},
			pngCRC: 0xd052fb55,
		},
		{
			name: "stack_series_dual_yaxis",
			makeOptions: func() LineChartOption {
				opt := makeFullLineChartStackedOption()
				opt.SeriesList[len(opt.SeriesList)-1].YAxisIndex = 1
				// disable extra stuff
				for i := range opt.SeriesList {
					opt.SeriesList[i].MarkLine = SeriesMarkLine{}
					opt.SeriesList[i].MarkPoint = SeriesMarkPoint{}
					opt.SeriesList[i].Label.Show = Ptr(false)
				}
				opt.Legend.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x4cfc45a1,
		},
		{
			name: "series_legend_order_sync",
			makeOptions: func() LineChartOption {
				opt := makeFullLineChartOption()
				// set names in weird order, series should shuffle to match legend order
				opt.SeriesList[0].Name = opt.Legend.SeriesNames[3]
				opt.SeriesList[1].Name = opt.Legend.SeriesNames[4]
				opt.SeriesList[2].Name = opt.Legend.SeriesNames[0]
				opt.SeriesList[3].Name = opt.Legend.SeriesNames[1]
				opt.SeriesList[4].Name = opt.Legend.SeriesNames[2]
				opt.SeriesList[4].MarkLine = NewMarkLine("min")
				// disable extra stuff
				opt.YAxis[0].Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				return opt
			},
			pngCRC: 0xed840231,
		},
		{
			name: "symbol_dot",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				opt.Symbol.Shape = SymbolDot
				opt.Legend.Symbol = SymbolDot
				opt.Title.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
			pngCRC: 0xbbe1d8ee,
		},
		{
			name: "symbol_circle",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				opt.Symbol.Shape = SymbolCircle
				opt.Legend.Symbol = SymbolCircle
				opt.Title.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
			pngCRC: 0x465c5745,
		},
		{
			name: "symbol_square",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				opt.Symbol.Shape = SymbolSquare
				opt.Legend.Symbol = SymbolSquare
				opt.Title.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
			pngCRC: 0xaab4f046,
		},
		{
			name: "symbol_diamond",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				opt.Symbol.Shape = SymbolDiamond
				opt.Legend.Symbol = SymbolDiamond
				opt.Title.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
			pngCRC: 0x8bcacb5d,
		},
		{
			name: "symbol_mixed",
			makeOptions: func() LineChartOption {
				opt := makeFullLineChartOption()
				opt.XAxis.Labels = opt.XAxis.Labels[:5]
				for i := range opt.SeriesList {
					opt.SeriesList[i].Values = opt.SeriesList[i].Values[:5]
				}
				opt.Symbol.Shape = SymbolNone
				opt.SeriesList[0].Symbol.Shape = SymbolCircle
				opt.SeriesList[1].Symbol.Shape = SymbolSquare
				opt.SeriesList[2].Symbol.Shape = SymbolDiamond
				opt.SeriesList[3].Symbol.Shape = SymbolDot
				opt.Title.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
			pngCRC: 0x7a516f45,
		},
		{
			name: "text_color_themes",
			makeOptions: func() LineChartOption {
				opt := makeFullLineChartOption()
				opt.Padding = NewBoxEqual(40)
				for i := range opt.SeriesList {
					markLine := NewMarkLine(SeriesMarkTypeAverage)
					markLine.ValueFormatter = func(f float64) string {
						return FormatValueHumanizeShort(f, 0, false)
					}
					opt.SeriesList[i].MarkLine = markLine
				}
				opt.Theme = GetTheme(ThemeAnt).WithTitleTextColor(ColorRed).WithLegendTextColor(ColorBlue).
					WithYAxisTextColor(ColorGreen).WithXAxisTextColor(ColorAqua).WithMarkTextColor(ColorPurple).
					WithLabelTextColor(ColorGold)
				opt.SeriesList = opt.SeriesList[:2]
				opt.Legend.SeriesNames = opt.Legend.SeriesNames[:2]
				return opt
			},
			pngCRC: 0xf1a7f880,
		},
		{
			name: "axis_titles",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				opt.Title.Text = ""
				opt.Legend.Show = Ptr(false)
				opt.XAxis.Title = "x-axis"
				opt.XAxis.TitleFontStyle.FontColor = ColorBlue
				opt.XAxis.TitleFontStyle.FontSize = 18
				opt.SeriesList[1].YAxisIndex = 1
				opt.YAxis = append(opt.YAxis, opt.YAxis[0])
				opt.YAxis[0].PreferNiceIntervals = Ptr(true)
				opt.YAxis[1].PreferNiceIntervals = Ptr(true)
				opt.YAxis[0].Title = "left y-axis"
				opt.YAxis[0].TitleFontStyle.FontColor = ColorBlue
				opt.YAxis[0].TitleFontStyle.FontSize = 18
				opt.YAxis[0].SpineLineShow = Ptr(true)
				opt.YAxis[1].Title = "right y-axis"
				opt.YAxis[1].TitleFontStyle.FontColor = ColorBlue
				opt.YAxis[1].TitleFontStyle.FontSize = 18
				opt.YAxis[1].SpineLineShow = Ptr(true)
				return opt
			},
			pngCRC: 0xec49ed33,
		},
		{
			name: "trend_line_linear",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.SeriesList = NewSeriesListLine([][]float64{
					{120, 132, 101, 134, 190, 230, 210},
				}, LineSeriesOption{
					TrendLine: []SeriesTrendLine{
						{Type: SeriesTrendTypeLinear, LineColor: ColorRed, DashedLine: Ptr(false)},
					},
				})
				return opt
			},
			pngCRC: 0xd9c39cf,
		},
		{
			name: "trend_line_cubic",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.SeriesList = NewSeriesListLine([][]float64{
					{120, 132, 101, 134, 190, 230, 210},
				}, LineSeriesOption{
					TrendLine: []SeriesTrendLine{
						{Type: SeriesTrendTypeCubic, LineColor: ColorRed, DashedLine: Ptr(false)},
					},
				})
				return opt
			},
			pngCRC: 0xc571e061,
		},
		{
			name: "trend_line_average",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.SeriesList = NewSeriesListLine([][]float64{
					{120, 132, 101, 134, 190, 230, 210},
				}, LineSeriesOption{
					TrendLine: []SeriesTrendLine{
						{Type: SeriesTrendTypeSMA, Period: 3, LineColor: ColorRed, DashedLine: Ptr(false)},
					},
				})
				return opt
			},
			pngCRC: 0x164293aa,
		},
		{
			name: "trend_line_sma",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.SeriesList = NewSeriesListLine([][]float64{
					{120, 132, 101, 134, 190, 230, 210},
				}, LineSeriesOption{
					TrendLine: []SeriesTrendLine{
						{Type: SeriesTrendTypeSMA, Period: 3, LineColor: ColorRed, DashedLine: Ptr(false)},
					},
				})
				return opt
			},
			pngCRC: 0x164293aa,
		},
		{
			name: "trend_line_multiple",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.SeriesList = NewSeriesListLine([][]float64{
					{120, 132, 101, 134, 190, 230, 210},
					{220, 232, 201, 234, 290, 330, 310},
				}, LineSeriesOption{
					TrendLine: []SeriesTrendLine{
						{Type: SeriesTrendTypeLinear, LineColor: ColorRed, DashedLine: Ptr(false)},
						{Type: SeriesTrendTypeSMA, LineColor: ColorBrown},
						{Type: SeriesTrendTypeCubic, LineColor: ColorRedAlt1, StrokeSmoothingTension: 0.8},
					},
				})
				return opt
			},
			pngCRC: 0xdfa24892,
		},
		{
			name: "bollinger",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				opt.SeriesList[0].TrendLine = []SeriesTrendLine{
					{Type: SeriesTrendTypeBollingerLower, Period: 3},
				}
				opt.SeriesList[1].TrendLine = []SeriesTrendLine{
					{Type: SeriesTrendTypeBollingerUpper, Period: 3},
				}
				opt.Title.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x4600b14,
		},
		{
			name: "rsi",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				opt.SeriesList[0].TrendLine = []SeriesTrendLine{
					{Type: SeriesTrendTypeRSI, Period: 3},
				}
				opt.Title.Show = Ptr(false)
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x82e7475a,
		},
		{
			name: "empty_series",
			makeOptions: func() LineChartOption {
				opt := NewLineChartOptionWithData([][]float64{})
				opt.Padding = NewBoxEqual(10)
				opt.Legend = LegendOption{
					Show:        Ptr(true),
					SeriesNames: []string{"Series A", "Series B"},
				}
				opt.XAxis = XAxisOption{Labels: []string{"Jan", "Feb", "Mar"}}
				opt.YAxis = []YAxisOption{{Show: Ptr(true)}, {Show: Ptr(true)}}
				return opt
			},
			pngCRC: 0xe7e75e6c,
		},
		{
			name: "symbol_large_size",
			makeOptions: func() LineChartOption {
				opt := makeBasicLineChartOption()
				opt.Symbol = Symbol{Shape: SymbolCircle, Size: 12}
				opt.Legend.Symbol = SymbolCircle
				opt.Title.Show = Ptr(false)
				opt.XAxis.Show = Ptr(false)
				opt.YAxis[0].Show = Ptr(false)
				return opt
			},
			pngCRC: 0xf3df345b,
		},
		{
			name: "legend_collision",
			makeOptions: func() LineChartOption {
				return LineChartOption{
					Padding: NewBoxEqual(10),
					XAxis: XAxisOption{
						Labels: []string{"A", "B", "C", "D", "E", "F", "G", "H"},
					},
					YAxis: []YAxisOption{
						{
							RangeValuePaddingScale: Ptr(0.1), // force collision
						},
					},
					Legend: LegendOption{
						SeriesNames: []string{"First Series", "Second Series"},
						Vertical:    Ptr(true),
						Offset:      OffsetRight,
					},
					SeriesList: NewSeriesListLine([][]float64{
						{1, 2, 4, 8, 16, 32, 64, 128},
						{1, 2, 3, 4, 5, 6, 7, 8},
					}),
				}
			},
			pngCRC: 0xfce4a185,
		},
		{
			name: "legend_collision_stacked",
			makeOptions: func() LineChartOption {
				return LineChartOption{
					Padding: NewBoxEqual(10),
					XAxis: XAxisOption{
						Labels: []string{"A", "B", "C", "D", "E", "F", "G", "H"},
					},
					YAxis: []YAxisOption{
						{
							RangeValuePaddingScale: Ptr(0.1), // force collision
						},
					},
					Legend: LegendOption{
						SeriesNames: []string{"First Series", "Second Series"},
						Vertical:    Ptr(true),
						Offset:      OffsetRight,
					},
					StackSeries: Ptr(true),
					SeriesList: NewSeriesListLine([][]float64{
						{1, 2, 3, 4, 5, 6, 7, 8},
						{1, 2, 3, 4, 5, 6, 7, 8},
					}),
				}
			},
			pngCRC: 0xfa2e7856,
		},
		{
			name: "legend_collision_markpoint",
			makeOptions: func() LineChartOption {
				sl := NewSeriesListLine([][]float64{
					{2, 4, 6, 8, 10, 12, 14, 16},
					{1, 2, 3, 4, 5, 6, 7, 8},
				})
				sl[0].MarkPoint.AddPoints(SeriesMarkTypeMax)
				return LineChartOption{
					Padding: NewBoxEqual(10),
					XAxis: XAxisOption{
						Labels: []string{"A", "B", "C", "D", "E", "F", "G", "H"},
					},
					Legend: LegendOption{
						SeriesNames: []string{"First Series", "Second Series"},
						Vertical:    Ptr(true),
						Offset:      OffsetRight,
					},
					SeriesList: sl,
				}
			},
			pngCRC: 0x1a686a7b,
		},
		{
			// Threshold formatter hides some labels, shortening the label slice below the data length
			name: "mark_point_threshold_labels",
			makeOptions: func() LineChartOption {
				return LineChartOption{
					Padding: NewBoxEqual(20),
					SeriesList: NewSeriesListLine([][]float64{{5, 900, 15, 300, 200}}, LineSeriesOption{
						Label: SeriesLabel{
							Show:           Ptr(true),
							LabelFormatter: LabelFormatterThresholdMin(100),
						},
						MarkPoint: NewMarkPoint(SeriesMarkTypeMax),
					}),
					XAxis:  XAxisOption{Labels: []string{"a", "b", "c", "d", "e"}},
					YAxis:  []YAxisOption{{Show: Ptr(false)}},
					Legend: LegendOption{Show: Ptr(false)},
				}
			},
			pngCRC: 0xa131d49d,
		},
		{
			name: "line_gap_label",
			makeOptions: func() LineChartOption {
				opt := makeMinimalLineChartOption()
				opt.SeriesList[0].Label.Show = Ptr(true)
				opt.SeriesList[0].Values[3] = GetNullValue()
				return opt
			},
			pngCRC: 0xae12cc16,
		},
		{
			name: "stack_series_interleaved_marks",
			makeOptions: func() LineChartOption {
				opt := NewLineChartOptionWithData([][]float64{
					{50, 500, 50, 500}, {1, 2, 3, 4}, {4, 3, 2, 1}, {500, 50, 500, 50},
				})
				opt.SeriesList[0].YAxisIndex = 1
				opt.SeriesList[3].YAxisIndex = 1
				opt.YAxis = make([]YAxisOption, 2)
				opt.StackSeries = Ptr(true)
				// marks are gated on the first and last stacked series, not the first and last series
				opt.SeriesList[1].MarkLine = NewMarkLine(SeriesMarkTypeMax)
				opt.SeriesList[2].MarkLine.AddGlobalLines(SeriesMarkTypeMax)
				opt.XAxis.Labels = []string{"a", "b", "c", "d"}
				opt.Legend.Show = Ptr(false)
				return opt
			},
			pngCRC: 0x7d2f36b3,
		},
	}

	for i, tt := range tests {
		painterOptions := PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        600,
			Height:       400,
		}
		if !tt.themed {
			t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
				p := NewPainter(painterOptions)
				r := NewPainter(PainterOptions{
					OutputFormat: ChartOutputPNG,
					Width:        600,
					Height:       400,
				})

				validateLineChartRender(t, p, r, tt.makeOptions(), tt.pngCRC)
			})
		} else {
			theme := GetTheme(ThemeVividDark)
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-theme_painter", func(t *testing.T) {
				p := NewPainter(painterOptions, PainterThemeOption(theme))
				r := NewPainter(PainterOptions{
					OutputFormat: ChartOutputPNG,
					Width:        600,
					Height:       400,
				}, PainterThemeOption(theme))

				validateLineChartRender(t, p, r, tt.makeOptions(), tt.pngCRC)
			})
			t.Run(strconv.Itoa(i)+"-"+tt.name+"-theme_opt", func(t *testing.T) {
				p := NewPainter(painterOptions)
				r := NewPainter(PainterOptions{
					OutputFormat: ChartOutputPNG,
					Width:        600,
					Height:       400,
				})
				opt := tt.makeOptions()
				opt.Theme = theme

				validateLineChartRender(t, p, r, opt, tt.pngCRC)
			})
		}
	}
}

func validateLineChartRender(t *testing.T, svgP, pngP *Painter, opt LineChartOption, expectedCRC uint32) {
	t.Helper()

	err := svgP.LineChart(opt)
	require.NoError(t, err)
	data, err := svgP.Bytes()
	require.NoError(t, err)
	assertTestdataSVG(t, data)

	err = pngP.LineChart(opt)
	require.NoError(t, err)
	rdata, err := pngP.Bytes()
	require.NoError(t, err)
	assertEqualPNGCRC(t, expectedCRC, rdata)
}

func TestLineChartError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		makeOptions      func() LineChartOption
		errorMsgContains string
	}{
		{
			name: "invalid_yaxis_index",
			makeOptions: func() LineChartOption {
				opt := NewLineChartOptionWithData([][]float64{{1, 2, 3}})
				opt.SeriesList[0].YAxisIndex = 2
				return opt
			},
			errorMsgContains: "invalid y-axis index",
		},
		{
			name: "negative_yaxis_index",
			makeOptions: func() LineChartOption {
				opt := NewLineChartOptionWithData([][]float64{{1, 2, 3}})
				opt.SeriesList[0].YAxisIndex = -1
				return opt
			},
			errorMsgContains: "invalid y-axis index",
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			})

			err := p.LineChart(tt.makeOptions())
			require.Error(t, err)
			require.ErrorContains(t, err, tt.errorMsgContains)
		})
	}
}

func TestBoundaryGapAxisPositions(t *testing.T) {
	t.Parallel()

	got := boundaryGapAxisPositions(10, false, 3)
	assert.Equal(t, []int{0, 5, 10}, got)
	assert.Equal(t, 0, got[0])
	assert.Equal(t, 10, got[len(got)-1])

	got = boundaryGapAxisPositions(10, true, 3)
	assert.Equal(t, []int{1, 4, 8}, got)
	assert.Equal(t, 1, got[0])
	assert.Equal(t, 8, got[len(got)-1])
}

func TestLineChartOptionNotMutated(t *testing.T) {
	t.Parallel()

	opt := NewLineChartOptionWithData([][]float64{{1, 2, 3}})
	opt.YAxis = []YAxisOption{{}, {}}

	p := NewPainter(PainterOptions{Width: 600, Height: 400}, PainterThemeOption(GetTheme(ThemeLight)))
	require.NoError(t, p.LineChart(opt))

	// theme must stay unset so a later render can apply its own painter theme
	assert.Nil(t, opt.YAxis[0].Theme)
	assert.Nil(t, opt.YAxis[1].Theme)
}
