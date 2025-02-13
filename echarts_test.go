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
				BoundaryGap: True(),
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
		Data: []SeriesMarkData{
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
	t.Parallel()

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
	t.Parallel()

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
	assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 600 400\"><path  d=\"M 0 0\nL 600 0\nL 600 400\nL 0 400\nL 0 0\" style=\"stroke:none;fill:white\"/><path  d=\"M 182 19\nL 212 19\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:none\"/><circle cx=\"197\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"214\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Rainfall</text><path  d=\"M 286 19\nL 316 19\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:none\"/><circle cx=\"301\" cy=\"19\" r=\"5\" style=\"stroke-width:3;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"318\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Evaporation</text><text x=\"10\" y=\"25\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Rainfall vs Evaporation</text><text x=\"54\" y=\"40\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Fake Data</text><text x=\"10\" y=\"62\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">270</text><text x=\"10\" y=\"95\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">240</text><text x=\"10\" y=\"129\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">210</text><text x=\"10\" y=\"163\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">180</text><text x=\"10\" y=\"197\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">150</text><text x=\"10\" y=\"231\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">120</text><text x=\"19\" y=\"265\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">90</text><text x=\"19\" y=\"299\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">60</text><text x=\"19\" y=\"333\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">30</text><text x=\"28\" y=\"367\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">0</text><path  d=\"M 47 55\nL 570 55\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 47 88\nL 570 88\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 47 122\nL 570 122\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 47 156\nL 570 156\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 47 190\nL 570 190\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 47 224\nL 570 224\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 47 258\nL 570 258\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 47 292\nL 570 292\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 47 326\nL 570 326\" style=\"stroke-width:1;stroke:rgb(224,230,242);fill:none\"/><path  d=\"M 47 365\nL 47 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 90 365\nL 90 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 134 365\nL 134 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 177 365\nL 177 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 221 365\nL 221 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 264 365\nL 264 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 308 365\nL 308 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 352 365\nL 352 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 395 365\nL 395 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 439 365\nL 439 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 482 365\nL 482 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 526 365\nL 526 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 570 365\nL 570 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><path  d=\"M 47 360\nL 570 360\" style=\"stroke-width:1;stroke:rgb(110,112,121);fill:none\"/><text x=\"55\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jan</text><text x=\"99\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Feb</text><text x=\"141\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Mar</text><text x=\"187\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Apr</text><text x=\"227\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">May</text><text x=\"273\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jun</text><text x=\"320\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Jul</text><text x=\"359\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Aug</text><text x=\"404\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Sep</text><text x=\"448\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Oct</text><text x=\"490\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Nov</text><text x=\"535\" y=\"385\" style=\"stroke:none;fill:rgb(70,70,70);font-size:15.3px;font-family:'Roboto Medium',sans-serif\">Dec</text><path  d=\"M 52 358\nL 67 358\nL 67 359\nL 52 359\nL 52 358\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 95 355\nL 110 355\nL 110 359\nL 95 359\nL 95 355\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 139 353\nL 154 353\nL 154 359\nL 139 359\nL 139 353\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 182 334\nL 197 334\nL 197 359\nL 182 359\nL 182 334\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 226 332\nL 241 332\nL 241 359\nL 226 359\nL 226 332\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 269 274\nL 284 274\nL 284 359\nL 269 359\nL 269 274\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 313 207\nL 328 207\nL 328 359\nL 313 359\nL 313 207\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 357 177\nL 372 177\nL 372 359\nL 357 359\nL 357 177\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 400 324\nL 415 324\nL 415 359\nL 400 359\nL 400 324\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 444 338\nL 459 338\nL 459 359\nL 444 359\nL 444 338\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 487 353\nL 502 353\nL 502 359\nL 487 359\nL 487 353\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 531 357\nL 546 357\nL 546 359\nL 531 359\nL 531 357\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 70 358\nL 85 358\nL 85 359\nL 70 359\nL 70 358\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 113 354\nL 128 354\nL 128 359\nL 113 359\nL 113 354\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 157 350\nL 172 350\nL 172 359\nL 157 359\nL 157 350\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 200 331\nL 215 331\nL 215 359\nL 200 359\nL 200 331\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 244 328\nL 259 328\nL 259 359\nL 244 359\nL 244 328\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 287 281\nL 302 281\nL 302 359\nL 287 359\nL 287 281\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 331 162\nL 346 162\nL 346 359\nL 331 359\nL 331 162\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 375 155\nL 390 155\nL 390 359\nL 375 359\nL 375 155\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 418 305\nL 433 305\nL 433 359\nL 418 359\nL 418 305\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 462 339\nL 477 339\nL 477 359\nL 462 359\nL 462 339\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 505 354\nL 520 354\nL 520 359\nL 505 359\nL 505 354\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 549 358\nL 564 358\nL 564 359\nL 549 359\nL 549 358\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 361 169\nA 14 14 330.00 1 1 367 169\nL 364 156\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 350 156\nQ364,191 378,156\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><text x=\"351\" y=\"161\" style=\"stroke:none;fill:rgb(238,238,238);font-size:10.2px;font-family:'Roboto Medium',sans-serif\">162.2</text><path  d=\"M 56 350\nA 14 14 330.00 1 1 62 350\nL 59 337\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><path  d=\"M 45 337\nQ59,372 73,337\nZ\" style=\"stroke:none;fill:rgb(84,112,198)\"/><text x=\"55\" y=\"342\" style=\"stroke:none;fill:rgb(238,238,238);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2</text><path  d=\"M 379 147\nA 14 14 330.00 1 1 385 147\nL 382 134\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 368 134\nQ382,169 396,134\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"369\" y=\"139\" style=\"stroke:none;fill:rgb(70,70,70);font-size:10.2px;font-family:'Roboto Medium',sans-serif\">182.2</text><path  d=\"M 553 350\nA 14 14 330.00 1 1 559 350\nL 556 337\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><path  d=\"M 542 337\nQ556,372 570,337\nZ\" style=\"stroke:none;fill:rgb(145,204,117)\"/><text x=\"547\" y=\"342\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">2.3</text><circle cx=\"50\" cy=\"313\" r=\"3\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 56 313\nL 552 313\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 552 308\nL 568 313\nL 552 318\nL 557 313\nL 552 308\" style=\"stroke-width:1;stroke:rgb(84,112,198);fill:rgb(84,112,198)\"/><text x=\"570\" y=\"317\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">41.63</text><circle cx=\"50\" cy=\"306\" r=\"3\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 56 306\nL 552 306\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><path stroke-dasharray=\"4.0, 2.0\" d=\"M 552 301\nL 568 306\nL 552 311\nL 557 306\nL 552 301\" style=\"stroke-width:1;stroke:rgb(145,204,117);fill:rgb(145,204,117)\"/><text x=\"570\" y=\"310\" style=\"stroke:none;fill:rgb(70,70,70);font-size:12.8px;font-family:'Roboto Medium',sans-serif\">48.07</text></svg>", data)
}
