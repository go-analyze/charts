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
		expected string // Placeholder for expected SVG, can be empty for now
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
			expected: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"10\" y=\"26\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Rainfall vs Evaporation</text><text x=\"54\" y=\"42\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fake Data</text><path d=\"M 182 19\nL 212 19\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"197\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"214\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Rainfall</text><path d=\"M 286 19\nL 316 19\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"301\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"318\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Evaporation</text><text x=\"9\" y=\"63\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">270</text><text x=\"9\" y=\"96\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">240</text><text x=\"9\" y=\"130\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">210</text><text x=\"9\" y=\"164\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">180</text><text x=\"9\" y=\"198\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">150</text><text x=\"9\" y=\"232\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">120</text><text x=\"18\" y=\"266\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">90</text><text x=\"18\" y=\"300\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">60</text><text x=\"18\" y=\"334\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">30</text><text x=\"27\" y=\"368\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path d=\"M 42 57\nL 570 57\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 91\nL 570 91\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 125\nL 570 125\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 159\nL 570 159\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 193\nL 570 193\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 227\nL 570 227\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 261\nL 570 261\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 295\nL 570 295\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 42 329\nL 570 329\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 46 364\nL 570 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 46 369\nL 46 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 89 369\nL 89 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 133 369\nL 133 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 177 369\nL 177 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 220 369\nL 220 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 264 369\nL 264 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 308 369\nL 308 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 351 369\nL 351 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 395 369\nL 395 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 439 369\nL 439 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 482 369\nL 482 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 526 369\nL 526 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 570 369\nL 570 364\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"54\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jan</text><text x=\"98\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Feb</text><text x=\"141\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mar</text><text x=\"186\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Apr</text><text x=\"227\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">May</text><text x=\"273\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jun</text><text x=\"319\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jul</text><text x=\"359\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Aug</text><text x=\"404\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sep</text><text x=\"448\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Oct</text><text x=\"490\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Nov</text><text x=\"535\" y=\"390\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Dec</text><path d=\"M 51 362\nL 66 362\nL 66 363\nL 51 363\nL 51 362\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 94 359\nL 109 359\nL 109 363\nL 94 363\nL 94 359\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 138 357\nL 153 357\nL 153 363\nL 138 363\nL 138 357\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 182 338\nL 197 338\nL 197 363\nL 182 363\nL 182 338\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 225 335\nL 240 335\nL 240 363\nL 225 363\nL 225 335\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 269 277\nL 284 277\nL 284 363\nL 269 363\nL 269 277\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 313 210\nL 328 210\nL 328 363\nL 313 363\nL 313 210\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 356 180\nL 371 180\nL 371 363\nL 356 363\nL 356 180\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 400 327\nL 415 327\nL 415 363\nL 400 363\nL 400 327\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 444 342\nL 459 342\nL 459 363\nL 444 363\nL 444 342\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 487 357\nL 502 357\nL 502 363\nL 487 363\nL 487 357\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 531 361\nL 546 361\nL 546 363\nL 531 363\nL 531 361\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 69 362\nL 84 362\nL 84 363\nL 69 363\nL 69 362\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 112 358\nL 127 358\nL 127 363\nL 112 363\nL 112 358\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 156 354\nL 171 354\nL 171 363\nL 156 363\nL 156 354\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 200 334\nL 215 334\nL 215 363\nL 200 363\nL 200 334\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 243 332\nL 258 332\nL 258 363\nL 243 363\nL 243 332\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 287 284\nL 302 284\nL 302 363\nL 287 363\nL 287 284\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 331 165\nL 346 165\nL 346 363\nL 331 363\nL 331 165\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 374 157\nL 389 157\nL 389 363\nL 374 363\nL 374 157\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 418 309\nL 433 309\nL 433 363\nL 418 363\nL 418 309\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 462 343\nL 477 343\nL 477 363\nL 462 363\nL 462 343\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 505 358\nL 520 358\nL 520 363\nL 505 363\nL 505 358\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 549 362\nL 564 362\nL 564 363\nL 549 363\nL 549 362\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 359 173\nA 14 14 330.00 1 1 367 173\nL 363 159\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 349 159\nQ363,194 377,159\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><text x=\"350\" y=\"164\" style=\"stroke:none;fill:rgb(238,238,238);font-size:10.2px;font-family:'Roboto Medium',sans-serif\">162.2</text><path d=\"M 54 355\nA 14 14 330.00 1 1 62 355\nL 58 341\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 44 341\nQ58,376 72,341\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><text x=\"54\" y=\"346\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2</text><path d=\"M 377 150\nA 14 14 330.00 1 1 385 150\nL 381 136\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 367 136\nQ381,171 395,136\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"368\" y=\"141\" style=\"stroke:none;fill:rgb(70,70,70);font-size:10.2px;font-family:'Roboto Medium',sans-serif\">182.2</text><path d=\"M 552 355\nA 14 14 330.00 1 1 560 355\nL 556 341\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path d=\"M 542 341\nQ556,376 570,341\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"547\" y=\"346\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2.3</text><circle cx=\"49\" cy=\"317\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 55 317\nL 552 317\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 552 312\nL 568 317\nL 552 322\nL 557 317\nL 552 312\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"570\" y=\"321\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">41.63</text><circle cx=\"49\" cy=\"310\" r=\"3\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 55 310\nL 552 310\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 552 305\nL 568 310\nL 552 315\nL 557 310\nL 552 305\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"570\" y=\"314\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">48.07</text></svg>",
		},
		{
			name: "basic_bar_chart",
			jsonData: `{
				"title": { "text": "Sales" },
				"xAxis": { "type": "category", "data": ["Jan", "Feb"] },
				"yAxis": { "type": "value" },
				"series": [{ "data": [100, 200], "type": "bar" }]
			}`,
			expected: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"20\" y=\"36\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sales</text><text x=\"19\" y=\"57\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">208</text><text x=\"19\" y=\"90\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">196</text><text x=\"19\" y=\"123\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">184</text><text x=\"19\" y=\"157\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">172</text><text x=\"19\" y=\"190\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">160</text><text x=\"19\" y=\"224\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">148</text><text x=\"19\" y=\"257\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">136</text><text x=\"19\" y=\"291\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">124</text><text x=\"19\" y=\"324\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">112</text><text x=\"19\" y=\"358\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">100</text><path d=\"M 52 51\nL 580 51\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 84\nL 580 84\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 118\nL 580 118\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 152\nL 580 152\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 185\nL 580 185\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 219\nL 580 219\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 253\nL 580 253\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 286\nL 580 286\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 52 320\nL 580 320\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 56 354\nL 580 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 56 359\nL 56 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 318 359\nL 318 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 580 359\nL 580 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"174\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jan</text><text x=\"436\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Feb</text><path d=\"M 66 354\nL 308 354\nL 308 353\nL 66 353\nL 66 354\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 328 74\nL 570 74\nL 570 353\nL 328 353\nL 328 74\" style=\"stroke:none;fill:rgb(84,112,198)\"/></svg>",
		},
		{
			name: "axis_styling",
			jsonData: `{
				"xAxis": { "axisLabel": { "color": "#ff0000", "fontSize": 14 } },
				"yAxis": { "axisLabel": { "color": "#00ff00", "fontSize": 12 } },
				"series": [{ "data": [10, 20], "type": "bar" }]
			}`,
			expected: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"28\" y=\"26\" style=\"stroke:none;fill:lime;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">20.5</text><text x=\"19\" y=\"62\" style=\"stroke:none;fill:lime;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">19.33</text><text x=\"19\" y=\"99\" style=\"stroke:none;fill:lime;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">18.17</text><text x=\"41\" y=\"136\" style=\"stroke:none;fill:lime;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">17</text><text x=\"19\" y=\"172\" style=\"stroke:none;fill:lime;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">15.83</text><text x=\"19\" y=\"209\" style=\"stroke:none;fill:lime;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">14.67</text><text x=\"28\" y=\"246\" style=\"stroke:none;fill:lime;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">13.5</text><text x=\"19\" y=\"282\" style=\"stroke:none;fill:lime;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">12.33</text><text x=\"19\" y=\"319\" style=\"stroke:none;fill:lime;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">11.17</text><text x=\"41\" y=\"356\" style=\"stroke:none;fill:lime;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">10</text><path d=\"M 65 20\nL 580 20\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 56\nL 580 56\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 93\nL 580 93\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 130\nL 580 130\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 167\nL 580 167\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 204\nL 580 204\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 241\nL 580 241\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 278\nL 580 278\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 315\nL 580 315\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 69 352\nL 580 352\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 69 357\nL 69 352\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 324 357\nL 324 352\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 580 357\nL 580 352\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"191\" y=\"380\" style=\"stroke:none;fill:red;font-size:17.9px;font-family:'Roboto Medium',sans-serif\">1</text><text x=\"447\" y=\"380\" style=\"stroke:none;fill:red;font-size:17.9px;font-family:'Roboto Medium',sans-serif\">2</text><path d=\"M 79 352\nL 314 352\nL 314 351\nL 79 351\nL 79 352\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 334 36\nL 569 36\nL 569 351\nL 334 351\nL 334 36\" style=\"stroke:none;fill:rgb(84,112,198)\"/></svg>",
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
			expected: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"29\" y=\"18\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">15.5</text><text x=\"29\" y=\"55\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">14.33</text><text x=\"29\" y=\"92\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">13.17</text><text x=\"29\" y=\"129\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">12</text><text x=\"29\" y=\"166\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">10.83</text><text x=\"29\" y=\"203\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">9.67</text><text x=\"29\" y=\"240\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">8.5</text><text x=\"29\" y=\"277\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">7.33</text><text x=\"29\" y=\"314\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">6.17</text><text x=\"29\" y=\"351\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">5</text><path d=\"M 25 20\nL 580 20\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 57\nL 580 57\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 94\nL 580 94\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 131\nL 580 131\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 168\nL 580 168\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 206\nL 580 206\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 243\nL 580 243\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 280\nL 580 280\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 317\nL 580 317\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 29 355\nL 580 355\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 29 360\nL 29 355\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 304 360\nL 304 355\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 580 360\nL 580 355\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"166\" y=\"365\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">X1</text><text x=\"442\" y=\"365\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">X2</text><path d=\"M 39 355\nL 294 355\nL 294 354\nL 39 354\nL 39 355\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 314 36\nL 569 36\nL 569 354\nL 314 354\nL 314 36\" style=\"stroke:none;fill:rgb(84,112,198)\"/></svg>",
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
			expected: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 259 29\nL 289 29\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"274\" cy=\"29\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"291\" y=\"35\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Series1</text><path d=\"M 249 51\nL 249 10\nL 351 10\nL 351 51\nL 249 51\" style=\"stroke-width:2;stroke:lime;fill:none\"/><text x=\"29\" y=\"54\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">30.5</text><text x=\"29\" y=\"87\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">29.33</text><text x=\"29\" y=\"120\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">28.17</text><text x=\"29\" y=\"153\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">27</text><text x=\"29\" y=\"186\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">25.83</text><text x=\"29\" y=\"219\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">24.67</text><text x=\"29\" y=\"252\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">23.5</text><text x=\"29\" y=\"285\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">22.33</text><text x=\"29\" y=\"318\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">21.17</text><text x=\"29\" y=\"351\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">20</text><path d=\"M 25 56\nL 580 56\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 89\nL 580 89\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 122\nL 580 122\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 155\nL 580 155\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 188\nL 580 188\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 222\nL 580 222\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 255\nL 580 255\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 288\nL 580 288\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 321\nL 580 321\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 29 355\nL 580 355\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 29 360\nL 29 355\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 304 360\nL 304 355\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 580 360\nL 580 355\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"166\" y=\"365\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"442\" y=\"365\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">B</text><path d=\"M 166 355\nL 442 71\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"166\" cy=\"355\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"442\" cy=\"71\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/></svg>",
		},
		{
			name: "yaxis_line_show",
			jsonData: `{
				"yAxis": {
					"axisLine": { "show": true, "lineStyle": { "color": "#ff0000", "opacity": 0.8 } }
				},
				"series": [{ "data": [5, 15], "type": "bar" }]
			}`,
			expected: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path d=\"M 69 20\nL 69 354\" style=\"stroke-width:1;stroke:rgba(255,0,0,0.8);fill:none\"/><path d=\"M 64 20\nL 69 20\" style=\"stroke-width:1;stroke:rgba(255,0,0,0.8);fill:none\"/><path d=\"M 64 57\nL 69 57\" style=\"stroke-width:1;stroke:rgba(255,0,0,0.8);fill:none\"/><path d=\"M 64 94\nL 69 94\" style=\"stroke-width:1;stroke:rgba(255,0,0,0.8);fill:none\"/><path d=\"M 64 131\nL 69 131\" style=\"stroke-width:1;stroke:rgba(255,0,0,0.8);fill:none\"/><path d=\"M 64 168\nL 69 168\" style=\"stroke-width:1;stroke:rgba(255,0,0,0.8);fill:none\"/><path d=\"M 64 205\nL 69 205\" style=\"stroke-width:1;stroke:rgba(255,0,0,0.8);fill:none\"/><path d=\"M 64 242\nL 69 242\" style=\"stroke-width:1;stroke:rgba(255,0,0,0.8);fill:none\"/><path d=\"M 64 279\nL 69 279\" style=\"stroke-width:1;stroke:rgba(255,0,0,0.8);fill:none\"/><path d=\"M 64 316\nL 69 316\" style=\"stroke-width:1;stroke:rgba(255,0,0,0.8);fill:none\"/><path d=\"M 64 354\nL 69 354\" style=\"stroke-width:1;stroke:rgba(255,0,0,0.8);fill:none\"/><text x=\"28\" y=\"26\" style=\"stroke:none;fill:rgba(255,0,0,0.8);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">15.5</text><text x=\"19\" y=\"62\" style=\"stroke:none;fill:rgba(255,0,0,0.8);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">14.33</text><text x=\"19\" y=\"99\" style=\"stroke:none;fill:rgba(255,0,0,0.8);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">13.17</text><text x=\"41\" y=\"136\" style=\"stroke:none;fill:rgba(255,0,0,0.8);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">12</text><text x=\"19\" y=\"173\" style=\"stroke:none;fill:rgba(255,0,0,0.8);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">10.83</text><text x=\"28\" y=\"210\" style=\"stroke:none;fill:rgba(255,0,0,0.8);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">9.67</text><text x=\"37\" y=\"247\" style=\"stroke:none;fill:rgba(255,0,0,0.8);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">8.5</text><text x=\"28\" y=\"284\" style=\"stroke:none;fill:rgba(255,0,0,0.8);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">7.33</text><text x=\"28\" y=\"321\" style=\"stroke:none;fill:rgba(255,0,0,0.8);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">6.17</text><text x=\"50\" y=\"358\" style=\"stroke:none;fill:rgba(255,0,0,0.8);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">5</text><path d=\"M 65 20\nL 580 20\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 57\nL 580 57\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 94\nL 580 94\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 131\nL 580 131\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 168\nL 580 168\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 205\nL 580 205\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 242\nL 580 242\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 279\nL 580 279\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 65 316\nL 580 316\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 70 354\nL 580 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 70 359\nL 70 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 325 359\nL 325 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 580 359\nL 580 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"193\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1</text><text x=\"448\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2</text><path d=\"M 80 354\nL 315 354\nL 315 353\nL 80 353\nL 80 354\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 335 36\nL 570 36\nL 570 353\nL 335 353\nL 335 36\" style=\"stroke:none;fill:rgb(84,112,198)\"/></svg>",
		},
		{
			name: "two_yaxis",
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
			expected: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"559\" y=\"26\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">5</text><text x=\"559\" y=\"92\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">4.2</text><text x=\"559\" y=\"158\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">3.4</text><text x=\"559\" y=\"225\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">2.6</text><text x=\"559\" y=\"291\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1.8</text><text x=\"559\" y=\"358\" style=\"stroke:none;fill:blue;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">1</text><text x=\"19\" y=\"26\" style=\"stroke:none;fill:red;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">61.5</text><text x=\"32\" y=\"62\" style=\"stroke:none;fill:red;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">58</text><text x=\"19\" y=\"99\" style=\"stroke:none;fill:red;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">54.5</text><text x=\"32\" y=\"136\" style=\"stroke:none;fill:red;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">51</text><text x=\"19\" y=\"173\" style=\"stroke:none;fill:red;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">47.5</text><text x=\"32\" y=\"210\" style=\"stroke:none;fill:red;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">44</text><text x=\"19\" y=\"247\" style=\"stroke:none;fill:red;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">40.5</text><text x=\"32\" y=\"284\" style=\"stroke:none;fill:red;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">37</text><text x=\"19\" y=\"321\" style=\"stroke:none;fill:red;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">33.5</text><text x=\"32\" y=\"358\" style=\"stroke:none;fill:red;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">30</text><path d=\"M 56 20\nL 549 20\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 56 57\nL 549 57\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 56 94\nL 549 94\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 56 131\nL 549 131\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 56 168\nL 549 168\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 56 205\nL 549 205\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 56 242\nL 549 242\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 56 279\nL 549 279\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 56 316\nL 549 316\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 60 354\nL 549 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 60 359\nL 60 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 304 359\nL 304 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 549 359\nL 549 354\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"169\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jan</text><text x=\"413\" y=\"380\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Feb</text><path d=\"M 70 354\nL 294 354\nL 294 353\nL 70 353\nL 70 354\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 314 36\nL 538 36\nL 538 353\nL 314 353\nL 314 36\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path d=\"M 182 313\nL 426 171\" style=\"stroke-width:2;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"182\" cy=\"313\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/><circle cx=\"426\" cy=\"171\" r=\"2\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:white\"/></svg>",
		},
		{
			name: "background_color",
			jsonData: `{
				"backgroundColor": "#e0e0e0",
				"xAxis": { "axisLabel": { "show": false }, "type": "category", "data": ["A", "B"] },
				"yAxis": { "axisLabel": { "show": false }, "type": "value" },
				"series": [{ "data": [40, 70], "type": "line" }]
			}`,
			expected: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:rgb(224,224,224)\"/><text x=\"29\" y=\"18\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">71.5</text><text x=\"29\" y=\"55\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">68</text><text x=\"29\" y=\"92\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">64.5</text><text x=\"29\" y=\"129\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">61</text><text x=\"29\" y=\"166\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">57.5</text><text x=\"29\" y=\"203\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">54</text><text x=\"29\" y=\"240\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">50.5</text><text x=\"29\" y=\"277\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">47</text><text x=\"29\" y=\"314\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">43.5</text><text x=\"29\" y=\"351\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">40</text><path d=\"M 25 20\nL 580 20\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 57\nL 580 57\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 94\nL 580 94\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 131\nL 580 131\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 168\nL 580 168\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 206\nL 580 206\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 243\nL 580 243\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 280\nL 580 280\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 317\nL 580 317\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 29 355\nL 580 355\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 29 360\nL 29 355\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 304 360\nL 304 355\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 580 360\nL 580 355\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"166\" y=\"365\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"442\" y=\"365\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">B</text><path d=\"M 166 355\nL 442 36\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"166\" cy=\"355\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(224,224,224)\"/><circle cx=\"442\" cy=\"36\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(224,224,224)\"/></svg>",
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
			expected: "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><text x=\"20\" y=\"36\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Title</text><path d=\"M 10 46\nL 10 10\nL 61 10\nL 61 46\nL 10 46\" style=\"stroke-width:2;stroke:lime;fill:none\"/><text x=\"29\" y=\"49\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">71.5</text><text x=\"29\" y=\"82\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">68</text><text x=\"29\" y=\"116\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">64.5</text><text x=\"29\" y=\"149\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">61</text><text x=\"29\" y=\"183\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">57.5</text><text x=\"29\" y=\"216\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">54</text><text x=\"29\" y=\"250\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">50.5</text><text x=\"29\" y=\"283\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">47</text><text x=\"29\" y=\"317\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">43.5</text><text x=\"29\" y=\"351\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">40</text><path d=\"M 25 51\nL 580 51\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 84\nL 580 84\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 118\nL 580 118\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 152\nL 580 152\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 186\nL 580 186\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 219\nL 580 219\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 253\nL 580 253\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 287\nL 580 287\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 25 321\nL 580 321\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path d=\"M 29 355\nL 580 355\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 29 360\nL 29 355\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 304 360\nL 304 355\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path d=\"M 580 360\nL 580 355\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"166\" y=\"365\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">A</text><text x=\"442\" y=\"365\" style=\"stroke:none;fill:none;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">B</text><path d=\"M 166 355\nL 442 66\" style=\"stroke-width:2;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"166\" cy=\"355\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/><circle cx=\"442\" cy=\"66\" r=\"2\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:white\"/></svg>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := RenderEChartsToSVG(tt.jsonData)
			require.NoError(t, err)
			assertEqualSVG(t, tt.expected, data)
		})
	}
}
