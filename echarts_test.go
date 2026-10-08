package charts

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertToArray(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []byte(`[1]`), convertToArray([]byte("1")))
	assert.Equal(t, []byte(`[1]`), convertToArray([]byte("[1]")))
}

func TestEChartsPosition(t *testing.T) {
	t.Parallel()

	var p EChartsPosition
	require.NoError(t, p.UnmarshalJSON([]byte("1")))
	assert.Equal(t, EChartsPosition("1"), p)
	require.NoError(t, p.UnmarshalJSON([]byte(`"left"`)))
	assert.Equal(t, EChartsPosition("left"), p)
	require.NoError(t, p.UnmarshalJSON([]byte("-10")))
	assert.Equal(t, EChartsPosition("-10"), p)
	require.NoError(t, p.UnmarshalJSON([]byte("1.5")))
	assert.Equal(t, EChartsPosition("1.5"), p)
	require.NoError(t, p.UnmarshalJSON([]byte("1e3")))
	assert.Equal(t, EChartsPosition("1e3"), p)

	p = EChartsPosition("top")
	require.NoError(t, p.UnmarshalJSON([]byte("null")))
	assert.Equal(t, EChartsPosition("top"), p) // null leaves the value untouched
}

func TestEChartsSeriesDataValue(t *testing.T) {
	t.Parallel()

	es := EChartsSeriesDataValue{}
	require.NoError(t, es.UnmarshalJSON([]byte(`[1, 2]`)))
	assert.Equal(t, EChartsSeriesDataValue{
		values: []float64{1, 2},
	}, es)
	assert.Equal(t, EChartsSeriesDataValue{values: []float64{1, 2}}, es)
	assert.InDelta(t, 1.0, es.First(), 0)

	require.NoError(t, es.UnmarshalJSON([]byte("1e3")))
	assert.Equal(t, EChartsSeriesDataValue{values: []float64{1000}}, es)

	require.NoError(t, es.UnmarshalJSON([]byte("null")))
	assert.Equal(t, EChartsSeriesDataValue{values: []float64{GetNullValue()}}, es)

	require.NoError(t, es.UnmarshalJSON([]byte(`"-"`)))
	assert.Equal(t, EChartsSeriesDataValue{values: []float64{GetNullValue()}}, es)

	require.NoError(t, es.UnmarshalJSON([]byte(`[1, null, "-", 2]`)))
	assert.Equal(t, EChartsSeriesDataValue{
		values: []float64{1, GetNullValue(), GetNullValue(), 2},
	}, es)

	require.Error(t, es.UnmarshalJSON([]byte(`"foo"`)))

	var empty EChartsSeriesDataValue
	assert.InDelta(t, GetNullValue(), empty.First(), 0)
}

func TestEChartsSeriesData(t *testing.T) {
	t.Parallel()

	es := EChartsSeriesData{}
	require.NoError(t, es.UnmarshalJSON([]byte("1.1")))
	assert.Equal(t, EChartsSeriesDataValue{
		values: []float64{1.1},
	}, es.Value)

	require.NoError(t, es.UnmarshalJSON([]byte(`{"value":200,"itemStyle":{"color":"#a90000"}}`)))
	assert.Equal(t, EChartsSeriesData{
		Value: EChartsSeriesDataValue{
			values: []float64{200.0},
		},
		ItemStyle: EChartStyle{
			Color: "#a90000",
		},
	}, es)

	es = EChartsSeriesData{}
	require.NoError(t, es.UnmarshalJSON([]byte("1e3")))
	assert.Equal(t, EChartsSeriesDataValue{values: []float64{1000}}, es.Value)

	require.NoError(t, es.UnmarshalJSON([]byte("null")))
	assert.Equal(t, EChartsSeriesDataValue{values: []float64{GetNullValue()}}, es.Value)

	require.NoError(t, es.UnmarshalJSON([]byte(`"-"`)))
	assert.Equal(t, EChartsSeriesDataValue{values: []float64{GetNullValue()}}, es.Value)

	require.NoError(t, es.UnmarshalJSON([]byte(`[1, 2]`)))
	assert.Equal(t, EChartsSeriesDataValue{values: []float64{1, 2}}, es.Value)

	require.NoError(t, es.UnmarshalJSON([]byte(`{"value":null,"name":"foo"}`)))
	assert.Equal(t, EChartsSeriesData{
		Name:  "foo",
		Value: EChartsSeriesDataValue{values: []float64{GetNullValue()}},
	}, es)
}

func TestEChartsXAxis(t *testing.T) {
	t.Parallel()

	ex := EChartsXAxis{}
	require.NoError(t, ex.UnmarshalJSON([]byte(`{"boundaryGap": true, "splitNumber": 5, "data": ["a", "b"], "type": "value"}`)))

	assert.Equal(t, EChartsXAxis{
		Data: []EChartsXAxisData{
			{
				BoundaryGap: Ptr(true),
				SplitNumber: 5,
				Data:        []string{"a", "b"},
				Type:        "value",
			},
		},
	}, ex)
}

func TestEChartsPadding(t *testing.T) {
	t.Parallel()

	eb := EChartsPadding{}

	require.NoError(t, eb.UnmarshalJSON([]byte(`1`)))
	assert.Equal(t, NewBoxEqual(1), eb.Box)

	require.NoError(t, eb.UnmarshalJSON([]byte(`[2, 3]`)))
	assert.Equal(t, Box{
		Left:   3,
		Top:    2,
		Right:  3,
		Bottom: 2,
		IsSet:  true,
	}, eb.Box)

	require.NoError(t, eb.UnmarshalJSON([]byte(`[4, 5, 6]`)))
	assert.Equal(t, Box{
		Left:   5,
		Top:    4,
		Right:  5,
		Bottom: 6,
		IsSet:  true,
	}, eb.Box)

	require.NoError(t, eb.UnmarshalJSON([]byte(`[4, 5, 6, 7]`)))
	assert.Equal(t, Box{
		Left:   7,
		Top:    4,
		Right:  5,
		Bottom: 6,
		IsSet:  true,
	}, eb.Box)
}

func TestEChartsMarkPoint(t *testing.T) {
	t.Parallel()

	emp := EChartsMarkPoint{
		SymbolSize: 30,
		Data: []EChartsMarkData{
			{
				Type: "test",
			},
		},
	}
	assert.Equal(t, SeriesMarkPoint{
		SymbolSize: 30,
		Points: []SeriesMark{
			{
				Type: "test",
			},
		},
	}, emp.ToSeriesMarkPoint())
}

func TestEChartsMarkLine(t *testing.T) {
	t.Parallel()

	eml := EChartsMarkLine{
		Data: []EChartsMarkData{
			{
				Type: "min",
			},
			{
				Type: "max",
			},
		},
	}
	assert.Equal(t, SeriesMarkLine{
		Lines: []SeriesMark{
			{
				Type: "min",
			},
			{
				Type: "max",
			},
		},
	}, eml.ToSeriesMarkLine())
}

func TestEChartsOption(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		option         string
		expectedValues []float64 // values of the first series, verified when set
	}{
		{
			option: `{
				"xAxis": {
					"type": "category",
					"data": [
						"Mon",
						"Tue",
						"Wed",
						"Thu",
						"Fri",
						"Sat",
						"Sun"
					]
				},
				"yAxis": {
					"type": "value"
				},
				"series": [
					{
						"data": [
							120,
							{
								"value": 200,
								"itemStyle": {
									"color": "#a90000"
								}
							},
							150,
							80,
							70,
							110,
							130
						],
						"type": "bar"
					}
				]
			}`,
		},
		{
			option: `{
				"title": {
					"text": "Referer of a Website",
					"subtext": "Fake Data",
					"left": "center"
				},
				"tooltip": {
					"trigger": "item"
				},
				"legend": {
					"orient": "vertical",
					"left": "left"
				},
				"series": [
					{
						"name": "Access From",
						"type": "pie",
						"radius": "50%",
						"data": [
							{
								"value": 1048,
								"name": "Search Engine"
							},
							{
								"value": 735,
								"name": "Direct"
							},
							{
								"value": 580,
								"name": "Email"
							},
							{
								"value": 484,
								"name": "Union Ads"
							},
							{
								"value": 300,
								"name": "Video Ads"
							}
						],
						"emphasis": {
							"itemStyle": {
								"shadowBlur": 10,
								"shadowOffsetX": 0,
								"shadowColor": "rgba(0, 0, 0, 0.5)"
							}
						}
					}
				]
			}`,
		},
		{
			option: `{
				"title": {
					"text": "Rainfall vs Evaporation",
					"subtext": "Fake Data"
				},
				"tooltip": {
					"trigger": "axis"
				},
				"legend": {
					"data": [
						"Rainfall",
						"Evaporation"
					]
				},
				"toolbox": {
					"show": true,
					"feature": {
						"dataView": {
							"show": true,
							"readOnly": false
						},
						"magicType": {
							"show": true,
							"type": [
								"line",
								"bar"
							]
						},
						"restore": {
							"show": true
						},
						"saveAsImage": {
							"show": true
						}
					}
				},
				"calculable": true,
				"xAxis": [
					{
						"type": "category",
						"data": [
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
							"Dec"
						]
					}
				],
				"yAxis": [
					{
						"type": "value"
					}
				],
				"series": [
					{
						"name": "Rainfall",
						"type": "bar",
						"data": [
							2,
							4.9,
							7,
							23.2,
							25.6,
							76.7,
							135.6,
							162.2,
							32.6,
							20,
							6.4,
							3.3
						],
						"markPoint": {
							"data": [
								{
									"type": "max",
									"name": "Max"
								},
								{
									"type": "min",
									"name": "Min"
								}
							]
						},
						"markLine": {
							"data": [
								{
									"type": "average",
									"name": "Avg"
								}
							]
						}
					},
					{
						"name": "Evaporation",
						"type": "bar",
						"data": [
							2.6,
							5.9,
							9,
							26.4,
							28.7,
							70.7,
							175.6,
							182.2,
							48.7,
							18.8,
							6,
							2.3
						],
						"markPoint": {
							"data": [
								{
									"name": "Max",
									"value": 182.2,
									"xAxis": 7,
									"yAxis": 183
								},
								{
									"name": "Min",
									"value": 2.3,
									"xAxis": 11,
									"yAxis": 3
								}
							]
						},
						"markLine": {
							"data": [
								{
									"type": "average",
									"name": "Avg"
								}
							]
						}
					}
				]
			}`,
		},
		{
			name: "basic_bar_Chart",
			option: `{
				"xAxis": { "type": "category", "data": ["Mon", "Tue", "Wed"] },
				"yAxis": { "type": "value" },
				"series": [{ "data": [120, 200, 150], "type": "bar" }]
			}`,
		},
		{
			name: "basic_pie_chart",
			option: `{
				"title": { "text": "Website Traffic", "left": "center" },
				"series": [{ "name": "Source", "type": "pie", "data": [{ "value": 100, "name": "Google" }] }]
			}`,
		},
		{
			name: "numeric_title_position",
			option: `{
				"title": { "text": "Offset Title", "left": -10, "top": 20 },
				"xAxis": { "type": "category", "data": ["Mon", "Tue"] },
				"yAxis": { "type": "value" },
				"series": [{ "data": [120, 200], "type": "bar" }]
			}`,
		},
		{
			name: "null_and_exponent_data",
			option: `{
				"xAxis": { "type": "category", "data": ["Mon", "Tue", "Wed", "Thu"] },
				"yAxis": { "type": "value" },
				"series": [{ "data": [150, null, "-", 1e3], "type": "line" }]
			}`,
			expectedValues: []float64{150, GetNullValue(), GetNullValue(), 1000},
		},
	}

	for i, tt := range tests {
		name := strconv.Itoa(i)
		if tt.name != "" {
			name += tt.name
		}
		t.Run(name, func(t *testing.T) {
			opt := EChartsOption{}
			require.NoError(t, json.Unmarshal([]byte(tt.option), &opt))
			assert.NotEmpty(t, opt.Series)
			chartOpt, err := opt.ToOption()
			require.NoError(t, err)
			seriesList := chartOpt.SeriesList
			assert.NotEmpty(t, seriesList)
			if tt.expectedValues != nil {
				assert.Equal(t, tt.expectedValues, seriesList[0].Values)
			}

			if len(opt.XAxis.Data) > 0 {
				assert.NotEmpty(t, opt.XAxis.Data[0].Data)
				assert.NotEmpty(t, opt.XAxis.Data[0].Type)
			}
		})
	}
}

// mustParseEChartsOption unmarshals the option JSON and converts it, failing the test on error.
func mustParseEChartsOption(t *testing.T, option string) (EChartsOption, ChartOption) {
	t.Helper()
	eChartsOpt := EChartsOption{}
	require.NoError(t, json.Unmarshal([]byte(option), &eChartsOpt))
	chartOpt, err := eChartsOpt.ToOption()
	require.NoError(t, err)
	return eChartsOpt, chartOpt
}

func TestEChartsSeriesStack(t *testing.T) {
	t.Parallel()

	t.Run("stacked_sets_chart_flag", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t, `{"series": [
			{"type": "bar", "stack": "total", "data": [1, 2]},
			{"type": "bar", "stack": "total", "data": [3, 4]}
		]}`)
		assert.True(t, flagIs(true, chartOpt.StackSeries))
		assert.Len(t, chartOpt.SeriesList, 2)
	})

	t.Run("mixed_stack_errors", func(t *testing.T) {
		eChartsOpt := EChartsOption{}
		require.NoError(t, json.Unmarshal([]byte(`{"series": [
			{"type": "bar", "stack": "total", "data": [1, 2]},
			{"type": "bar", "data": [3, 4]}
		]}`), &eChartsOpt))
		_, err := eChartsOpt.ToOption()
		assert.Error(t, err)
	})

	t.Run("multiple_groups_error", func(t *testing.T) {
		eChartsOpt := EChartsOption{}
		require.NoError(t, json.Unmarshal([]byte(`{"series": [
			{"type": "bar", "stack": "a", "data": [1, 2]},
			{"type": "bar", "stack": "b", "data": [3, 4]}
		]}`), &eChartsOpt))
		_, err := eChartsOpt.ToOption()
		assert.Error(t, err)
	})

	t.Run("non_stackable_types_ignored", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t, `{"series": [
			{"type": "bar", "stack": "total", "data": [1, 2]},
			{"type": "radar", "data": [{"value": [1, 2, 3], "name": "r"}]}
		]}`)
		assert.True(t, flagIs(true, chartOpt.StackSeries))
	})
}

func TestEChartsSeriesSymbols(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		option        string
		expectedShape SymbolShape
		expectedSize  float64
	}{
		{
			name: "first_defined_wins",
			option: `{"series": [
				{"type": "line", "data": [1, 2], "symbol": "circle"},
				{"type": "line", "data": [3, 4], "symbol": "rect", "symbolSize": 6}
			]}`,
			expectedShape: SymbolCircle,
			expectedSize:  6,
		},
		{
			name:          "rect_maps_to_square",
			option:        `{"series": [{"type": "line", "data": [1, 2], "symbol": "roundRect"}]}`,
			expectedShape: SymbolSquare,
		},
		{
			name:          "unsupported_falls_to_square",
			option:        `{"series": [{"type": "line", "data": [1, 2], "symbol": "triangle"}]}`,
			expectedShape: SymbolSquare,
		},
		{
			name:          "none_disables_symbols",
			option:        `{"series": [{"type": "line", "data": [1, 2], "symbol": "none"}]}`,
			expectedShape: SymbolNone,
		},
		{
			name:          "unset_keeps_default",
			option:        `{"series": [{"type": "line", "data": [1, 2]}]}`,
			expectedShape: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, chartOpt := mustParseEChartsOption(t, tt.option)
			assert.Equal(t, tt.expectedShape, chartOpt.Symbol.Shape)
			assert.InDelta(t, tt.expectedSize, chartOpt.Symbol.Size, 0)
		})
	}
}

func TestEChartsSeriesStrokeAndFill(t *testing.T) {
	t.Parallel()

	t.Run("line_width_first_wins", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t, `{"series": [
			{"type": "line", "data": [1, 2], "lineStyle": {"width": 3}},
			{"type": "line", "data": [3, 4], "lineStyle": {"width": 5}}
		]}`)
		assert.InDelta(t, 3.0, chartOpt.LineStrokeWidth, 0)
	})

	t.Run("area_style_enables_fill", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t,
			`{"series": [{"type": "line", "data": [1, 2], "areaStyle": {}}]}`)
		require.NotNil(t, chartOpt.FillArea)
		assert.True(t, *chartOpt.FillArea)
		assert.Zero(t, chartOpt.FillOpacity)
	})

	t.Run("area_opacity_scales_alpha", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t,
			`{"series": [{"type": "line", "data": [1, 2], "areaStyle": {"opacity": 0.5}}]}`)
		assert.Equal(t, uint8(128), chartOpt.FillOpacity)
	})

	t.Run("fill_area_first_wins", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t, `{"series": [
			{"type": "line", "data": [1, 2]},
			{"type": "line", "data": [3, 4], "areaStyle": {"opacity": 0.25}}
		]}`)
		require.NotNil(t, chartOpt.FillArea)
		assert.True(t, *chartOpt.FillArea)
		assert.Equal(t, uint8(64), chartOpt.FillOpacity)
	})
}

func TestEChartsSeriesBarSizing(t *testing.T) {
	t.Parallel()

	t.Run("bar_width_percent", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t,
			`{"series": [{"type": "bar", "data": [1, 2], "barWidth": "35%"}]}`)
		assert.InDelta(t, 0.35, chartOpt.BarSize, 0)
		assert.Nil(t, chartOpt.BarMargin)
	})

	t.Run("bar_gap_percent", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t,
			`{"series": [{"type": "bar", "data": [1, 2], "barGap": "30%"}]}`)
		require.NotNil(t, chartOpt.BarMargin)
		assert.InDelta(t, 0.30, *chartOpt.BarMargin, 0)
	})

	t.Run("absolute_sizes_ignored", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t,
			`{"series": [{"type": "bar", "data": [1, 2], "barWidth": "30", "barGap": "8", "barCategoryGap": "20%"}]}`)
		assert.Zero(t, chartOpt.BarSize)
		assert.Nil(t, chartOpt.BarMargin)
	})

	t.Run("bar_width_sets_candle_width", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t,
			`{"series": [{"type": "candlestick", "data": [[100, 110, 95, 105]], "barWidth": "40%"}]}`)
		assert.InDelta(t, 0.40, chartOpt.BarSize, 0)
	})
}

func TestEChartsSeriesSmooth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		option          string
		expectedTension float64
	}{
		{
			name:            "boolean_selects_default",
			option:          `{"series": [{"type": "line", "data": [1, 2], "smooth": true}]}`,
			expectedTension: 0.5,
		},
		{
			name:            "numeric_maps_directly",
			option:          `{"series": [{"type": "line", "data": [1, 2], "smooth": 0.8}]}`,
			expectedTension: 0.8,
		},
		{
			name:            "false_stays_disabled",
			option:          `{"series": [{"type": "line", "data": [1, 2], "smooth": false}]}`,
			expectedTension: 0,
		},
		{
			name:            "first_defined_wins",
			option:          `{"series": [{"type": "line", "data": [1, 2], "smooth": false}, {"type": "line", "data": [3, 4], "smooth": true}]}`,
			expectedTension: 0,
		},
		{
			name:            "values_clamp_to_one",
			option:          `{"series": [{"type": "line", "data": [1, 2], "smooth": 1.5}]}`,
			expectedTension: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, chartOpt := mustParseEChartsOption(t, tt.option)
			assert.InDelta(t, tt.expectedTension, chartOpt.StrokeSmoothingTension, 0)
		})
	}
}

func TestEChartsSeriesDoughnut(t *testing.T) {
	t.Parallel()

	t.Run("expands_data_items", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t,
			`{"series": [{"type": "doughnut", "data": [
				{"value": 10, "name": "a"},
				{"value": 20, "name": "b"}
			]}]}`)
		require.Len(t, chartOpt.SeriesList, 2)
		assert.Equal(t, ChartTypeDoughnut, chartOpt.SeriesList[0].Type)
		assert.Equal(t, "a", chartOpt.SeriesList[0].Name)
		assert.Equal(t, []float64{10}, chartOpt.SeriesList[0].Values)
		assert.Equal(t, "b", chartOpt.SeriesList[1].Name)
		assert.True(t, flagIs(true, chartOpt.SeriesList[0].Label.Show))
	})

	t.Run("radius_single_value", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t,
			`{"series": [{"type": "doughnut", "radius": "55%", "data": [{"value": 10, "name": "a"}]}]}`)
		assert.Equal(t, "55%", chartOpt.Radius)
		assert.Empty(t, chartOpt.RadiusCenter)
	})

	t.Run("radius_numeric_value", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t,
			`{"series": [{"type": "doughnut", "radius": 60, "data": [{"value": 10, "name": "a"}]}]}`)
		assert.Equal(t, "60", chartOpt.Radius)
	})

	t.Run("radius_inner_outer_pair", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t,
			`{"series": [{"type": "doughnut", "radius": ["40%", "70%"], "data": [{"value": 10, "name": "a"}]}]}`)
		assert.Equal(t, "70%", chartOpt.Radius)
		assert.Equal(t, "40%", chartOpt.RadiusCenter)
		assert.Equal(t, "70%", chartOpt.SeriesList[0].Radius)
	})
}

func TestEChartsAxisEnrichment(t *testing.T) {
	t.Parallel()

	t.Run("axis_names_and_styles", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t, `{"xAxis": {
				"type": "category", "data": ["a", "b"],
				"name": "XTitle", "nameTextStyle": {"color": "#ff0000", "fontSize": 14}
			},
			"yAxis": [{"name": "YTitle", "nameTextStyle": {"color": "#00ff00", "fontSize": 12}}],
			"series": [{"type": "bar", "data": [1, 2]}]}`)
		assert.Equal(t, "XTitle", chartOpt.XAxis.Title)
		assert.Equal(t, ParseColor("#ff0000"), chartOpt.XAxis.TitleFontStyle.FontColor)
		assert.InDelta(t, 14.0, chartOpt.XAxis.TitleFontStyle.FontSize, 0)
		require.NotEmpty(t, chartOpt.YAxis)
		assert.Equal(t, "YTitle", chartOpt.YAxis[0].Title)
		assert.Equal(t, ParseColor("#00ff00"), chartOpt.YAxis[0].TitleFontStyle.FontColor)
		assert.InDelta(t, 12.0, chartOpt.YAxis[0].TitleFontStyle.FontSize, 0)
	})

	t.Run("label_rotation_and_offsets", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t, `{"xAxis": {
				"type": "category", "data": ["a", "b"],
				"axisLabel": {"rotate": 30, "margin": 4}
			},
			"yAxis": [{"axisLabel": {"rotate": 45, "margin": 8, "interval": 2}, "splitNumber": 6}],
			"series": [{"type": "bar", "data": [1, 2]}]}`)
		assert.InDelta(t, DegreesToRadians(30), chartOpt.XAxis.LabelRotation, 0)
		assert.Equal(t, OffsetInt{Top: 4}, chartOpt.XAxis.LabelOffset)
		require.NotEmpty(t, chartOpt.YAxis)
		assert.InDelta(t, DegreesToRadians(45), chartOpt.YAxis[0].LabelRotation, 0)
		assert.Equal(t, OffsetInt{Left: -8}, chartOpt.YAxis[0].LabelOffset)
		assert.Equal(t, 2, chartOpt.YAxis[0].LabelSkipCount)
		assert.Equal(t, 6, chartOpt.YAxis[0].LabelCount)
	})

	t.Run("split_line_default_shown", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t,
			`{"yAxis": [{"type": "value"}], "series": [{"type": "bar", "data": [1, 2]}]}`)
		require.NotEmpty(t, chartOpt.YAxis)
		require.NotNil(t, chartOpt.YAxis[0].SplitLineShow)
		assert.True(t, *chartOpt.YAxis[0].SplitLineShow)
	})

	t.Run("split_line_explicit_hidden", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t,
			`{"yAxis": [{"splitLine": {"show": false}}], "series": [{"type": "bar", "data": [1, 2]}]}`)
		require.NotEmpty(t, chartOpt.YAxis)
		require.NotNil(t, chartOpt.YAxis[0].SplitLineShow)
		assert.False(t, *chartOpt.YAxis[0].SplitLineShow)
	})

	t.Run("yaxis_position", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t,
			`{"yAxis": [{"position": "right"}, {"position": "diagonal"}], "series": [{"type": "bar", "data": [1, 2]}]}`)
		require.Len(t, chartOpt.YAxis, 2)
		assert.Equal(t, PositionRight, chartOpt.YAxis[0].Position)
		assert.Empty(t, chartOpt.YAxis[1].Position)
	})

	t.Run("yaxis_margin_right_position", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t,
			`{"yAxis": [{"position": "right", "axisLabel": {"margin": 8}}], "series": [{"type": "bar", "data": [1, 2]}]}`)
		require.NotEmpty(t, chartOpt.YAxis)
		assert.Equal(t, OffsetInt{Left: 8}, chartOpt.YAxis[0].LabelOffset)
	})
}

func TestEChartsLabelFormatter(t *testing.T) {
	t.Parallel()

	t.Run("template_tokens", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t, `{"series": [{
				"type": "bar", "name": "S", "data": [{"value": 10, "name": "a"}, {"value": 30, "name": "b"}],
				"label": {"show": true, "formatter": "{a}/{b}: {c}"}
			}]}`)
		require.Len(t, chartOpt.SeriesList, 1)
		require.NotNil(t, chartOpt.SeriesList[0].Label.LabelFormatter)
		text, _ := chartOpt.SeriesList[0].Label.LabelFormatter(0, "S", 10)
		assert.Equal(t, "S/S: 10", text)
	})

	t.Run("null_value_hidden", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t, `{"series": [{
				"type": "bar", "name": "S", "data": [10], "label": {"formatter": "{c}"}
			}]}`)
		require.Len(t, chartOpt.SeriesList, 1)
		require.NotNil(t, chartOpt.SeriesList[0].Label.LabelFormatter)
		text, _ := chartOpt.SeriesList[0].Label.LabelFormatter(0, "S", GetNullValue())
		assert.Empty(t, text)
	})

	t.Run("percent_token", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t, `{"series": [{
				"type": "pie", "data": [{"value": 10, "name": "a"}, {"value": 30, "name": "b"}],
				"label": {"formatter": "{d}%"}
			}]}`)
		require.Len(t, chartOpt.SeriesList, 2)
		require.NotNil(t, chartOpt.SeriesList[0].Label.LabelFormatter)
		text, _ := chartOpt.SeriesList[0].Label.LabelFormatter(0, "a", 10)
		assert.Equal(t, "25%", text)
	})

	t.Run("position_bottom_offset", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t, `{"series": [{
				"type": "bar", "data": [10], "label": {"position": "bottom", "distance": 8}
			}]}`)
		require.Len(t, chartOpt.SeriesList, 1)
		assert.Equal(t, OffsetInt{Top: 16}, chartOpt.SeriesList[0].Label.Offset)
		assert.Equal(t, 8, chartOpt.SeriesList[0].Label.Distance)
	})

	t.Run("non_string_ignored", func(t *testing.T) {
		_, chartOpt := mustParseEChartsOption(t, `{"series": [{
				"type": "bar", "data": [10], "label": {"formatter": {"a": 1}}
			}]}`)
		require.Len(t, chartOpt.SeriesList, 1)
		assert.Nil(t, chartOpt.SeriesList[0].Label.LabelFormatter)
	})
}

func TestRenderEChartsToSVG(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		jsonData string
	}{
		{
			name: "detailed",
			jsonData: `{
		"title": {
			"text": "Rainfall vs Evaporation",
			"subtext": "Fake Data"
		},
		"legend": {
			"data": [
				"Rainfall",
				"Evaporation"
			]
		},
		"padding": [10, 30, 10, 10],
		"xAxis": [
			{
				"type": "category",
				"data": [
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
					"Dec"
				]
			}
		],
		"series": [
			{
				"name": "Rainfall",
				"type": "bar",
				"data": [
					2,
					4.9,
					7,
					23.2,
					25.6,
					76.7,
					135.6,
					162.2,
					32.6,
					20,
					6.4,
					3.3
				],
				"markPoint": {
					"data": [
						{
							"type": "max"
						},
						{
							"type": "min"
						}
					]
				},
				"markLine": {
					"data": [
						{
							"type": "average"
						}
					]
				}
			},
			{
				"name": "Evaporation",
				"type": "bar",
				"data": [
					2.6,
					5.9,
					9,
					26.4,
					28.7,
					70.7,
					175.6,
					182.2,
					48.7,
					18.8,
					6,
					2.3
				],
				"markPoint": {
					"data": [
						{
							"type": "max"
						},
						{
							"type": "min"
						}
					]
				},
				"markLine": {
					"data": [
						{
							"type": "average"
						}
					]
				}
			}
		]
	}`,
		},
		{
			name: "basic_bar_chart",
			jsonData: `{
				"title": { "text": "Sales" },
				"xAxis": { "type": "category", "data": ["Jan", "Feb"] },
				"yAxis": { "type": "value" },
				"series": [{ "data": [100, 200], "type": "bar" }]
			}`,
		},
		{
			name: "axis_styling",
			jsonData: `{
				"xAxis": { "axisLabel": { "color": "#ff0000", "fontSize": 14 } },
				"yAxis": { "axisLabel": { "color": "#00ff00", "fontSize": 12 } },
				"series": [{ "data": [10, 20], "type": "bar" }]
			}`,
		},
		{
			name: "title_and_axis_labels_hidden",
			jsonData: `{
				"title": {
					"show": false,
					"text": "Hidden Title"
				},
				"xAxis": { "axisLabel": { "show": false }, "type": "category", "data": ["X1", "X2"] },
				"yAxis": { "axisLabel": { "show": false }, "type": "value" },
				"series": [{ "data": [5, 15], "type": "bar" }]
			}`,
		},
		{
			name: "legend_border_color",
			jsonData: `{
				"legend": {
					"borderColor": "#00ff00",
					"data": ["Series1"]
				},
				"xAxis": { "axisLabel": { "show": false }, "type": "category", "data": ["A", "B"] },
				"yAxis": { "axisLabel": { "show": false }, "type": "value" },
				"series": [{ "data": [20, 30], "type": "line" }]
			}`,
		},
		{
			name: "yaxis_line_show",
			jsonData: `{
				"yAxis": {
					"axisLine": { "show": true, "lineStyle": { "color": "#ff0000", "opacity": 0.8 } }
				},
				"series": [{ "data": [5, 15], "type": "bar" }]
			}`,
		},
		{
			name: "dual_yaxis",
			jsonData: `{
				"xAxis": { "type": "category", "data": ["Jan", "Feb"] },
				"yAxis": [
					{ "type": "value", "axisLabel": { "color": "#ff0000" } },
					{ "type": "value", "axisLabel": { "color": "#0000ff" } }
				],
				"series": [
					{ "data": [30, 60], "type": "bar", "yAxisIndex": 0 },
					{ "data": [1.5, 3.2], "type": "line", "yAxisIndex": 1 }
				]
			}`,
		},
		{
			name: "background_color",
			jsonData: `{
				"backgroundColor": "#e0e0e0",
				"xAxis": { "axisLabel": { "show": false }, "type": "category", "data": ["A", "B"] },
				"yAxis": { "axisLabel": { "show": false }, "type": "value" },
				"series": [{ "data": [40, 70], "type": "line" }]
			}`,
		},
		{
			name: "title_border",
			jsonData: `{
				"title": {
					"text": "Title",
					"borderColor": "#00ff00"
				},
				"xAxis": { "axisLabel": { "show": false }, "type": "category", "data": ["A", "B"] },
				"yAxis": { "axisLabel": { "show": false }, "type": "value" },
				"series": [{ "data": [40, 70], "type": "line" }]
			}`,
		},
		{
			name: "stacked_bars",
			jsonData: `{
				"xAxis": { "type": "category", "data": ["Q1", "Q2", "Q3"] },
				"series": [
					{ "name": "Online", "type": "bar", "stack": "total", "data": [120, 200, 150] },
					{ "name": "Retail", "type": "bar", "stack": "total", "data": [80, 70, 110] }
				]
			}`,
		},
		{
			name: "line_symbols",
			jsonData: `{
				"xAxis": { "type": "category", "data": ["A", "B", "C", "D"] },
				"series": [
					{ "name": "Circle", "type": "line", "symbol": "circle", "symbolSize": 10, "data": [10, 30, 20, 40] },
					{ "name": "Square", "type": "line", "symbol": "rect", "symbolSize": 8, "data": [25, 15, 35, 30] },
					{ "name": "None", "type": "line", "symbol": "none", "data": [40, 25, 45, 15] }
				]
			}`,
		},
		{
			name: "line_smooth_area",
			jsonData: `{
				"xAxis": { "type": "category", "data": ["Mon", "Tue", "Wed", "Thu", "Fri"] },
				"series": [
					{ "name": "Smooth", "type": "line", "smooth": true, "lineStyle": { "width": 3 }, "data": [15, 35, 25, 45, 30] },
					{ "name": "Area", "type": "line", "areaStyle": { "opacity": 0.3 }, "data": [30, 20, 40, 35, 50] }
				]
			}`,
		},
		{
			name: "bar_width_and_gap",
			jsonData: `{
				"xAxis": { "type": "category", "data": ["A", "B", "C"] },
				"series": [
					{ "name": "S1", "type": "bar", "barWidth": "25%", "barGap": "20%", "data": [40, 70, 55] },
					{ "name": "S2", "type": "bar", "barWidth": "25%", "data": [30, 50, 65] }
				]
			}`,
		},
		{
			name: "doughnut_radius_formatter",
			jsonData: `{
				"series": [{
					"name": "Share", "type": "doughnut", "radius": ["40%", "70%"],
					"label": { "formatter": "{b}: {d}%" },
					"data": [
						{ "value": 40, "name": "Alpha" },
						{ "value": 30, "name": "Beta" },
						{ "value": 30, "name": "Gamma" }
					]
				}]
			}`,
		},
		{
			name: "axis_titles_and_rotation",
			jsonData: `{
				"xAxis": {
					"type": "category", "data": ["Jan", "Feb", "Mar", "Apr"],
					"name": "Month", "nameTextStyle": { "color": "#ff0000", "fontSize": 12 },
					"axisLabel": { "rotate": 30 }
				},
				"yAxis": {
					"name": "Value", "nameTextStyle": { "color": "#0000ff", "fontSize": 12 },
					"axisLabel": { "rotate": 45 }
				},
				"series": [{ "name": "S", "type": "bar", "data": [10, 25, 15, 30] }]
			}`,
		},
		{
			name: "label_formatter_position",
			jsonData: `{
				"xAxis": { "type": "category", "data": ["A", "B", "C"] },
				"series": [{
					"name": "Sales", "type": "bar",
					"label": { "show": true, "position": "top", "formatter": "{a}: {c}" },
					"data": [10, 25, 15]
				}]
			}`,
		},
		{
			name: "label_position_bottom",
			jsonData: `{
				"xAxis": { "type": "category", "data": ["A", "B", "C"] },
				"series": [{
					"name": "Sales", "type": "bar",
					"label": { "show": true, "position": "bottom", "distance": 8 },
					"data": [10, 25, 15]
				}]
			}`,
		},
		{
			name: "dual_yaxis_right",
			jsonData: `{
				"xAxis": { "type": "category", "data": ["Jan", "Feb", "Mar"] },
				"yAxis": [
					{ "type": "value", "name": "Left" },
					{ "type": "value", "name": "Right", "position": "right", "axisLabel": { "margin": 12 } }
				],
				"series": [
					{ "name": "Bars", "type": "bar", "yAxisIndex": 0, "data": [30, 60, 45] },
					{ "name": "Line", "type": "line", "yAxisIndex": 1, "data": [1.5, 3.2, 2.4] }
				]
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := RenderEChartsToSVG(tt.jsonData)
			require.NoError(t, err)
			assertTestdataSVG(t, data)
		})
	}
}
