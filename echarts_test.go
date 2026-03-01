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
		name   string
		option string
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
			assert.NotEmpty(t, opt.ToOption().SeriesList)

			if len(opt.XAxis.Data) > 0 {
				assert.NotEmpty(t, opt.XAxis.Data[0].Data)
				assert.NotEmpty(t, opt.XAxis.Data[0].Type)
			}
		})
	}
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := RenderEChartsToSVG(tt.jsonData)
			require.NoError(t, err)
			assertTestdataSVG(t, data)
		})
	}
}
