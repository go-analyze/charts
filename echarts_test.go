package charts

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

func TestConvertToArray(t *testing.T) {
	assert.Equal(t, []byte(`[1]`), convertToArray([]byte("1")))
	assert.Equal(t, []byte(`[1]`), convertToArray([]byte("[1]")))
}

func TestEChartsPosition(t *testing.T) {
	var p EChartsPosition
	require.NoError(t, p.UnmarshalJSON([]byte("1")))
	assert.Equal(t, EChartsPosition("1"), p)
	require.NoError(t, p.UnmarshalJSON([]byte(`"left"`)))
	assert.Equal(t, EChartsPosition("left"), p)
}

func TestEChartsSeriesDataValue(t *testing.T) {
	es := EChartsSeriesDataValue{}
	require.NoError(t, es.UnmarshalJSON([]byte(`[1, 2]`)))
	assert.Equal(t, EChartsSeriesDataValue{
		values: []float64{
			1,
			2,
		},
	}, es)
	assert.Equal(t, NewEChartsSeriesDataValue(1, 2), es)
	assert.Equal(t, 1.0, es.First())
}

func TestEChartsSeriesData(t *testing.T) {
	es := EChartsSeriesData{}
	require.NoError(t, es.UnmarshalJSON([]byte("1.1")))
	assert.Equal(t, EChartsSeriesDataValue{
		values: []float64{
			1.1,
		},
	}, es.Value)

	require.NoError(t, es.UnmarshalJSON([]byte(`{"value":200,"itemStyle":{"color":"#a90000"}}`)))
	assert.Equal(t, EChartsSeriesData{
		Value: EChartsSeriesDataValue{
			values: []float64{
				200.0,
			},
		},
		ItemStyle: EChartStyle{
			Color: "#a90000",
		},
	}, es)
}

func TestEChartsXAxis(t *testing.T) {
	ex := EChartsXAxis{}
	require.NoError(t, ex.UnmarshalJSON([]byte(`{"boundaryGap": true, "splitNumber": 5, "data": ["a", "b"], "type": "value"}`)))

	assert.Equal(t, EChartsXAxis{
		Data: []EChartsXAxisData{
			{
				BoundaryGap: TrueFlag(),
				SplitNumber: 5,
				Data: []string{
					"a",
					"b",
				},
				Type: "value",
			},
		},
	}, ex)
}

func TestEChartStyle(t *testing.T) {
	es := EChartStyle{
		Color: "#999",
	}
	color := drawing.Color{
		R: 153,
		G: 153,
		B: 153,
		A: 255,
	}
	assert.Equal(t, Style{
		FillColor:   color,
		FontColor:   color,
		StrokeColor: color,
	}, es.ToStyle())
}

func TestEChartsPadding(t *testing.T) {
	eb := EChartsPadding{}

	require.NoError(t, eb.UnmarshalJSON([]byte(`1`)))
	assert.Equal(t, Box{
		Left:   1,
		Top:    1,
		Right:  1,
		Bottom: 1,
	}, eb.Box)

	require.NoError(t, eb.UnmarshalJSON([]byte(`[2, 3]`)))
	assert.Equal(t, Box{
		Left:   3,
		Top:    2,
		Right:  3,
		Bottom: 2,
	}, eb.Box)

	require.NoError(t, eb.UnmarshalJSON([]byte(`[4, 5, 6]`)))
	assert.Equal(t, Box{
		Left:   5,
		Top:    4,
		Right:  5,
		Bottom: 6,
	}, eb.Box)

	require.NoError(t, eb.UnmarshalJSON([]byte(`[4, 5, 6, 7]`)))
	assert.Equal(t, Box{
		Left:   7,
		Top:    4,
		Right:  5,
		Bottom: 6,
	}, eb.Box)
}

func TestEChartsMarkPoint(t *testing.T) {
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
		Data: []SeriesMarkData{
			{
				Type: "test",
			},
		},
	}, emp.ToSeriesMarkPoint())
}

func TestEChartsMarkLine(t *testing.T) {
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
		Data: []SeriesMarkData{
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
	tests := []struct {
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
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			opt := EChartsOption{}
			require.NoError(t, json.Unmarshal([]byte(tt.option), &opt))
			assert.NotEmpty(t, opt.Series)
			assert.NotEmpty(t, opt.ToOption().SeriesList)
		})
	}
}

func TestRenderEChartsToSVG(t *testing.T) {
	data, err := RenderEChartsToSVG(`{
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
	}`)
	require.NoError(t, err)
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" width=\"600\" height=\"400\">\\n<path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke-width:0;stroke:none;fill:rgba(255,255,255,1.0)\"/><path  d=\"M 182 19\nL 212 19\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><circle cx=\"197\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><text x=\"214\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Rainfall</text><path  d=\"M 286 19\nL 316 19\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><circle cx=\"301\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path  d=\"\" style=\"stroke-width:3;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><text x=\"318\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Evaporation</text><text x=\"10\" y=\"25\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Rainfall vs Evaporation</text><text x=\"54\" y=\"40\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fake Data</text><text x=\"31\" y=\"67\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">200</text><text x=\"10\" y=\"117\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">166.39</text><text x=\"10\" y=\"167\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">132.79</text><text x=\"18\" y=\"217\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">99.19</text><text x=\"18\" y=\"267\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">65.59</text><text x=\"18\" y=\"317\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">31.99</text><text x=\"22\" y=\"367\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">-1.60</text><path  d=\"M 68 60\nL 570 60\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 110\nL 570 110\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 160\nL 570 160\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 210\nL 570 210\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 260\nL 570 260\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 310\nL 570 310\" style=\"stroke-width:1;stroke:rgba(224,230,242,1.0);fill:none\"/><path  d=\"M 68 365\nL 68 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 109 365\nL 109 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 151 365\nL 151 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 235 365\nL 235 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 277 365\nL 277 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 319 365\nL 319 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 360 365\nL 360 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 402 365\nL 402 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 486 365\nL 486 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 528 365\nL 528 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 570 365\nL 570 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><path  d=\"M 68 360\nL 570 360\" style=\"stroke-width:1;stroke:rgba(110,112,121,1.0);fill:none\"/><text x=\"68\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jan</text><text x=\"117\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Feb</text><text x=\"158\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mar</text><text x=\"241\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">May</text><text x=\"285\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jun</text><text x=\"329\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jul</text><text x=\"367\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Aug</text><text x=\"410\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sep</text><text x=\"493\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Nov</text><text x=\"543\" y=\"385\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Dec</text><path  d=\"M 73 355\nL 87 355\nL 87 359\nL 73 359\nL 73 355\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 114 351\nL 128 351\nL 128 359\nL 114 359\nL 114 351\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 156 348\nL 170 348\nL 170 359\nL 156 359\nL 156 348\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 198 324\nL 212 324\nL 212 359\nL 198 359\nL 198 324\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 240 320\nL 254 320\nL 254 359\nL 240 359\nL 240 320\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 282 244\nL 296 244\nL 296 359\nL 282 359\nL 282 244\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 324 156\nL 338 156\nL 338 359\nL 324 359\nL 324 156\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 365 117\nL 379 117\nL 379 359\nL 365 359\nL 365 117\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 407 310\nL 421 310\nL 421 359\nL 407 359\nL 407 310\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 449 328\nL 463 328\nL 463 359\nL 449 359\nL 449 328\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 491 349\nL 505 349\nL 505 359\nL 491 359\nL 491 349\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 533 353\nL 547 353\nL 547 359\nL 533 359\nL 533 353\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 90 354\nL 104 354\nL 104 359\nL 90 359\nL 90 354\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 131 349\nL 145 349\nL 145 359\nL 131 359\nL 131 349\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 173 345\nL 187 345\nL 187 359\nL 173 359\nL 173 345\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 215 319\nL 229 319\nL 229 359\nL 215 359\nL 215 319\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 257 315\nL 271 315\nL 271 359\nL 257 359\nL 257 315\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 299 253\nL 313 253\nL 313 359\nL 299 359\nL 299 253\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 341 97\nL 355 97\nL 355 359\nL 341 359\nL 341 97\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 382 87\nL 396 87\nL 396 359\nL 382 359\nL 382 87\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 424 286\nL 438 286\nL 438 359\nL 424 359\nL 424 286\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 466 330\nL 480 330\nL 480 359\nL 466 359\nL 466 330\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 508 349\nL 522 349\nL 522 359\nL 508 359\nL 508 349\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 550 355\nL 564 355\nL 564 359\nL 550 359\nL 550 355\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 369 109\nA 15 15 330.00 1 1 375 109\nL 372 95\nZ\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 357 95\nQ372,132 387,95\nZ\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><text x=\"359\" y=\"100\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:10.2px;font-family:'Roboto Medium',sans-serif\">162.2</text><path  d=\"M 77 347\nA 15 15 330.00 1 1 83 347\nL 80 333\nZ\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><path  d=\"M 65 333\nQ80,370 95,333\nZ\" style=\"stroke-width:0;stroke:none;fill:rgba(84,112,198,1.0)\"/><text x=\"76\" y=\"338\" style=\"stroke-width:0;stroke:none;fill:rgba(238,238,238,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2</text><path  d=\"M 386 79\nA 15 15 330.00 1 1 392 79\nL 389 65\nZ\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 374 65\nQ389,102 404,65\nZ\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><text x=\"376\" y=\"70\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:10.2px;font-family:'Roboto Medium',sans-serif\">182.2</text><path  d=\"M 554 347\nA 15 15 330.00 1 1 560 347\nL 557 333\nZ\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><path  d=\"M 542 333\nQ557,370 572,333\nZ\" style=\"stroke-width:0;stroke:none;fill:rgba(145,204,117,1.0)\"/><text x=\"548\" y=\"338\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2.3</text><circle cx=\"71\" cy=\"296\" r=\"3\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 77 296\nL 552 296\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 552 291\nL 568 296\nL 552 301\nL 557 296\nL 552 291\" style=\"stroke-width:1;stroke:rgba(84,112,198,1.0);fill:rgba(84,112,198,1.0)\"/><text x=\"570\" y=\"300\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">41.62</text><circle cx=\"71\" cy=\"287\" r=\"3\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 77 287\nL 552 287\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 552 282\nL 568 287\nL 552 292\nL 557 287\nL 552 282\" style=\"stroke-width:1;stroke:rgba(145,204,117,1.0);fill:rgba(145,204,117,1.0)\"/><text x=\"570\" y=\"291\" style=\"stroke-width:0;stroke:none;fill:rgba(70,70,70,1.0);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">48.07</text></svg>", string(data))
}
