package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-analyze/charts"
)

/*
Hosts an example http server which will render the charts to the web page.
*/

var html = `<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8" />
		<link rel="icon" href="/favicon.ico" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<style>
			body {
				background-color: #e0e0e0;
			}
			* {
				margin: 0;
				padding: 0;
			}
			.charts {
				margin: 10px auto;
				overflow: hidden;
			}
			.grid {
				float: left;
				margin-right: 10px;
			}
			.grid:nth-child(even) {
				margin-right: 0px;
			}
			.grid svg {
				margin-bottom: 10px;
			}
			h1 {
				text-align: center;
			}
			pre {
				width: 100%;
				margin: auto auto 20px auto;
				max-height: 300px;
				overflow: auto;
				display: block;
			}
			svg{
				margin: auto auto 10px auto;
				display: block;
			}
		</style>
		<title>go-charts</title>
	</head>
	<body>
		<div class="charts">{{body}}</div>
	</body>
</html>
`

func handler(w http.ResponseWriter, req *http.Request, chartOptions []charts.ChartOption, echartsOptions []string) {
	if req.URL.Path != "/" &&
		req.URL.Path != "/echarts" {
		return
	}
	query := req.URL.Query()
	theme := query.Get("theme")
	width, _ := strconv.Atoi(query.Get("width"))
	height, _ := strconv.Atoi(query.Get("height"))
	charts.SetDefaultWidth(width)
	charts.SetDefaultWidth(height)
	bytesList := make([][]byte, 0)
	for _, opt := range chartOptions {
		opt.Theme = charts.GetTheme(theme)
		opt.OutputFormat = charts.ChartOutputSVG
		d, err := charts.Render(opt)
		if err != nil {
			panic(err)
		}
		buf, err := d.Bytes()
		if err != nil {
			panic(err)
		}
		bytesList = append(bytesList, buf)
	}
	for _, opt := range echartsOptions {
		buf, err := charts.RenderEChartsToSVG(opt)
		if err != nil {
			panic(err)
		}
		bytesList = append(bytesList, buf)
	}

	p, err := charts.TableOptionRender(charts.TableChartOption{
		OutputFormat: charts.ChartOutputSVG,
		Header: []string{
			"Name",
			"Age",
			"Address",
			"Tag",
			"Action",
		},
		Data: [][]string{
			{
				"John Brown",
				"32",
				"New York No. 1 Lake Park",
				"nice, developer",
				"Send Mail",
			},
			{
				"Jim Green	",
				"42",
				"London No. 1 Lake Park",
				"wow",
				"Send Mail",
			},
			{
				"Joe Black	",
				"32",
				"Sidney No. 1 Lake Park",
				"cool, teacher",
				"Send Mail",
			},
		},
	})
	if err != nil {
		panic(err)
	}
	buf, err := p.Bytes()
	if err != nil {
		panic(err)
	}
	bytesList = append(bytesList, buf)

	data := bytes.ReplaceAll([]byte(html), []byte("{{body}}"), bytes.Join(bytesList, []byte("")))
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	chartOptions := []charts.ChartOption{
		{
			Title: charts.TitleOption{
				Text: "Line",
				FontStyle: charts.FontStyle{
					FontSize: 18,
				},
			},
			Legend: charts.LegendOption{
				Data: []string{
					"Email",
					"Union Ads",
					"Video Ads",
					"Direct",
					"Search Engine",
				},
				Padding: charts.Box{Left: 100},
			},
			XAxis: charts.XAxisOption{
				Data: []string{
					"Mon",
					"Tue",
					"Wed",
					"Thu",
					"Fri",
					"Sat",
					"Sun",
				},
			},
			SeriesList: []charts.Series{
				charts.NewSeriesFromValues([]float64{
					120,
					132,
					101,
					134,
					90,
					230,
					210,
				}),
				charts.NewSeriesFromValues([]float64{
					220,
					182,
					191,
					234,
					290,
					330,
					310,
				}),
				charts.NewSeriesFromValues([]float64{
					150,
					232,
					201,
					154,
					190,
					330,
					410,
				}),
				charts.NewSeriesFromValues([]float64{
					320,
					332,
					301,
					334,
					390,
					330,
					320,
				}),
				charts.NewSeriesFromValues([]float64{
					820,
					932,
					901,
					934,
					1290,
					1330,
					1320,
				}),
			},
		},
		// temperature line chart
		{
			Title: charts.TitleOption{
				Text: "Temperature Change in the Coming Week",
			},
			Padding: charts.Box{
				Top:    20,
				Left:   20,
				Right:  30,
				Bottom: 20,
			},
			Legend: charts.LegendOption{
				Data: []string{
					"Highest",
					"Lowest",
				},
			},
			XAxis: charts.XAxisOption{
				Data: []string{
					"Mon",
					"Tue",
					"Wed",
					"Thu",
					"Fri",
					"Sat",
					"Sun",
				},
				BoundaryGap: charts.False(),
			},
			SeriesList: []charts.Series{
				{
					Data: charts.NewSeriesDataFromValues([]float64{
						14,
						11,
						13,
						11,
						12,
						12,
						7,
					}),
					MarkPoint: charts.NewMarkPoint(charts.SeriesMarkDataTypeMax, charts.SeriesMarkDataTypeMin),
					MarkLine:  charts.NewMarkLine(charts.SeriesMarkDataTypeAverage),
				},
				{
					Data: charts.NewSeriesDataFromValues([]float64{
						1,
						-2,
						2,
						5,
						3,
						2,
						0,
					}),
					MarkLine: charts.NewMarkLine(charts.SeriesMarkDataTypeAverage),
				},
			},
		},
		{
			Title: charts.TitleOption{
				Text: "Line Area",
				FontStyle: charts.FontStyle{
					FontSize: 18,
				},
			},
			Legend: charts.LegendOption{
				Data: []string{"Email"},
			},
			XAxis: charts.XAxisOption{
				Data: []string{
					"Mon",
					"Tue",
					"Wed",
					"Thu",
					"Fri",
					"Sat",
					"Sun",
				},
			},
			YAxis: []charts.YAxisOption{{
				Min: charts.FloatPointer(0.0), // ensure y-axis starts at 0
			}},
			SeriesList: []charts.Series{
				charts.NewSeriesFromValues([]float64{
					120,
					132,
					101,
					134,
					90,
					230,
					210,
				}),
			},
			FillArea: true,
		},
		// histogram
		{
			Title: charts.TitleOption{
				Text: "Bar",
				FontStyle: charts.FontStyle{
					FontSize: 18,
				},
			},
			XAxis: charts.XAxisOption{
				Data: []string{
					"Mon",
					"Tue",
					"Wed",
					"Thu",
					"Fri",
					"Sat",
					"Sun",
				},
			},
			Legend: charts.LegendOption{
				Data: []string{
					"Rainfall",
					"Evaporation",
				},
				Icon: charts.IconRect,
			},
			SeriesList: []charts.Series{
				charts.NewSeriesFromValues([]float64{
					120,
					200,
					150,
					80,
					70,
					110,
					130,
				}, charts.ChartTypeBar),
				{
					Type: charts.ChartTypeBar,
					Data: []charts.SeriesData{
						{
							Value: 100,
						},
						{
							Value: 190,
						},
						{
							Value: 230,
						},
						{
							Value: 140,
						},
						{
							Value: 100,
						},
						{
							Value: 200,
						},
						{
							Value: 180,
						},
					},
					Label: charts.SeriesLabel{
						Show:     true,
						Position: charts.PositionBottom,
					},
				},
			},
		},
		// horizontal bar chart
		{
			Title: charts.TitleOption{
				Text: "World Population",
				FontStyle: charts.FontStyle{
					FontSize: 18,
				},
			},
			Padding: charts.Box{
				Top:    20,
				Right:  40,
				Bottom: 20,
				Left:   20,
			},
			Legend: charts.LegendOption{
				Data: []string{
					"2011",
					"2012",
				},
			},
			YAxis: []charts.YAxisOption{
				{
					Data: []string{
						"Brazil",
						"Indonesia",
						"USA",
						"India",
						"China",
						"World",
					},
				},
			},
			SeriesList: []charts.Series{
				{
					Type: charts.ChartTypeHorizontalBar,
					Data: charts.NewSeriesDataFromValues([]float64{
						18203,
						23489,
						29034,
						104970,
						131744,
						630230,
					}),
				},
				{
					Type: charts.ChartTypeHorizontalBar,
					Data: charts.NewSeriesDataFromValues([]float64{
						19325,
						23438,
						31000,
						121594,
						134141,
						681807,
					}),
				},
			},
		},
		// histogram+marker
		{
			Title: charts.TitleOption{
				Text:    "Rainfall vs Evaporation",
				Subtext: "(Fake Data)",
				FontStyle: charts.FontStyle{
					FontSize: 16,
				},
				SubtextFontStyle: charts.FontStyle{
					FontSize: 10,
				},
			},
			Padding: charts.Box{
				Top:    20,
				Right:  20,
				Bottom: 20,
				Left:   20,
			},
			XAxis: charts.XAxisOption{
				Data: []string{
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
					"Dec",
				},
			},
			Legend: charts.LegendOption{
				Data: []string{
					"Rainfall",
					"Evaporation",
				},
			},
			SeriesList: []charts.Series{
				{
					Type: charts.ChartTypeBar,
					Data: charts.NewSeriesDataFromValues([]float64{
						2.0,
						4.9,
						7.0,
						23.2,
						25.6,
						76.7,
						135.6,
						162.2,
						32.6,
						20.0,
						6.4,
						3.3,
					}),
					MarkPoint: charts.NewMarkPoint(
						charts.SeriesMarkDataTypeMax,
						charts.SeriesMarkDataTypeMin,
					),
					MarkLine: charts.NewMarkLine(
						charts.SeriesMarkDataTypeAverage,
					),
				},
				{
					Type: charts.ChartTypeBar,
					Data: charts.NewSeriesDataFromValues([]float64{
						2.6,
						5.9,
						9.0,
						26.4,
						28.7,
						70.7,
						175.6,
						182.2,
						48.7,
						18.8,
						6.0,
						2.3,
					}),
					MarkPoint: charts.NewMarkPoint(
						charts.SeriesMarkDataTypeMax,
						charts.SeriesMarkDataTypeMin,
					),
					MarkLine: charts.NewMarkLine(
						charts.SeriesMarkDataTypeAverage,
					),
				},
			},
		},
		// dual y-axis
		{
			Title: charts.TitleOption{
				Text: "Temperature",
				FontStyle: charts.FontStyle{
					FontSize: 16,
				},
			},
			XAxis: charts.XAxisOption{
				Data: []string{
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
					"Dec",
				},
			},
			Legend: charts.LegendOption{
				Data: []string{
					"Evaporation",
					"Precipitation",
					"Temperature",
				},
				Padding: charts.Box{Left: 100},
			},
			YAxis: []charts.YAxisOption{
				{
					Formatter: "{value}ml",
					AxisColor: charts.Color{R: 84, G: 112, B: 198, A: 255},
				},
				{
					Formatter: "{value}°C",
					AxisColor: charts.Color{R: 250, G: 200, B: 88, A: 255},
				},
			},
			SeriesList: []charts.Series{
				{
					Type: charts.ChartTypeBar,
					Data: charts.NewSeriesDataFromValues([]float64{
						2.0,
						4.9,
						7.0,
						23.2,
						25.6,
						76.7,
						135.6,
						162.2,
						32.6,
						20.0,
						6.4,
						3.3,
					}),
				},
				{
					Type: charts.ChartTypeBar,
					Data: charts.NewSeriesDataFromValues([]float64{
						2.6,
						5.9,
						9.0,
						26.4,
						28.7,
						70.7,
						175.6,
						182.2,
						48.7,
						18.8,
						6.0,
						2.3,
					}),
				},
				{
					Data: charts.NewSeriesDataFromValues([]float64{
						2.0,
						2.2,
						3.3,
						4.5,
						6.3,
						10.2,
						20.3,
						23.4,
						23.0,
						16.5,
						12.0,
						6.2,
					}),
					YAxisIndex: 1,
				},
			},
		},
		// pie chart
		{
			Title: charts.TitleOption{
				Text:    "Referer of a Website",
				Subtext: "(Fake Data)",
				Offset:  charts.OffsetCenter,
				FontStyle: charts.FontStyle{
					FontSize: 16,
				},
				SubtextFontStyle: charts.FontStyle{
					FontSize: 10,
				},
			},
			Legend: charts.LegendOption{
				Vertical: true,
				Data: []string{
					"Search Engine",
					"Direct",
					"Email",
					"Union Ads",
					"Video Ads",
				},
				Offset: charts.OffsetLeft,
			},
			SeriesList: charts.NewPieSeriesList([]float64{
				1048,
				735,
				580,
				484,
				300,
			}, charts.PieSeriesOption{
				Label: charts.SeriesLabel{
					Show: true,
				},
				Radius: "35%",
			}),
		},
		// radar chart
		{
			Title: charts.TitleOption{
				Text: "Basic Radar Chart",
				FontStyle: charts.FontStyle{
					FontSize: 16,
				},
			},
			Legend: charts.LegendOption{
				Data: []string{
					"Allocated Budget",
					"Actual Spending",
				},
				Padding: charts.Box{Left: 100},
			},
			RadarIndicators: []charts.RadarIndicator{
				{
					Name: "Sales",
					Max:  6500,
				},
				{
					Name: "Administration",
					Max:  16000,
				},
				{
					Name: "Information Technology",
					Max:  30000,
				},
				{
					Name: "Customer Support",
					Max:  38000,
				},
				{
					Name: "Development",
					Max:  52000,
				},
				{
					Name: "Marketing",
					Max:  25000,
				},
			},
			SeriesList: charts.SeriesList{
				{
					Type: charts.ChartTypeRadar,
					Data: charts.NewSeriesDataFromValues([]float64{
						4200,
						3000,
						20000,
						35000,
						50000,
						18000,
					}),
				},
				{
					Type: charts.ChartTypeRadar,
					Data: charts.NewSeriesDataFromValues([]float64{
						5000,
						14000,
						28000,
						26000,
						42000,
						21000,
					}),
				},
			},
		},
		// funnel chart
		{
			Title: charts.TitleOption{
				Text: "Funnel",
			},
			Legend: charts.LegendOption{
				Data: []string{
					"Show",
					"Click",
					"Visit",
					"Inquiry",
					"Order",
				},
			},
			SeriesList: []charts.Series{
				{
					Type: charts.ChartTypeFunnel,
					Name: "Show",
					Data: charts.NewSeriesDataFromValues([]float64{
						100,
					}),
				},
				{
					Type: charts.ChartTypeFunnel,
					Name: "Click",
					Data: charts.NewSeriesDataFromValues([]float64{
						80,
					}),
				},
				{
					Type: charts.ChartTypeFunnel,
					Name: "Visit",
					Data: charts.NewSeriesDataFromValues([]float64{
						60,
					}),
				},
				{
					Type: charts.ChartTypeFunnel,
					Name: "Inquiry",
					Data: charts.NewSeriesDataFromValues([]float64{
						40,
					}),
				},
				{
					Type: charts.ChartTypeFunnel,
					Name: "Order",
					Data: charts.NewSeriesDataFromValues([]float64{
						20,
					}),
				},
			},
		},
		// multiple picture
		{
			Legend: charts.LegendOption{
				Offset: charts.OffsetStr{
					Top: "-90",
				},
				Data: []string{
					"Milk Tea",
					"Matcha Latte",
					"Cheese Cocoa",
					"Walnut Brownie",
				},
			},
			Padding: charts.Box{
				Top:    100,
				Right:  10,
				Bottom: 10,
				Left:   10,
			},
			XAxis: charts.XAxisOption{
				Data: []string{
					"2012",
					"2013",
					"2014",
					"2015",
					"2016",
					"2017",
				},
			},
			YAxis: []charts.YAxisOption{
				{

					Min: charts.FloatPointer(0),
					Max: charts.FloatPointer(90),
				},
			},
			SeriesList: []charts.Series{
				charts.NewSeriesFromValues([]float64{
					56.5,
					82.1,
					88.7,
					70.1,
					53.4,
					85.1,
				}),
				charts.NewSeriesFromValues([]float64{
					51.1,
					51.4,
					55.1,
					53.3,
					73.8,
					68.7,
				}),
				charts.NewSeriesFromValues([]float64{
					40.1,
					62.2,
					69.5,
					36.4,
					45.2,
					32.5,
				}, charts.ChartTypeBar),
				charts.NewSeriesFromValues([]float64{
					25.2,
					37.1,
					41.2,
					18,
					33.9,
					49.1,
				}, charts.ChartTypeBar),
			},
			Children: []charts.ChartOption{
				{
					Legend: charts.LegendOption{
						Show: charts.False(),
						Data: []string{
							"Milk Tea",
							"Matcha Latte",
							"Cheese Cocoa",
							"Walnut Brownie",
						},
					},
					Box: charts.Box{
						Top:    20,
						Left:   400,
						Right:  500,
						Bottom: 120,
					},
					SeriesList: charts.NewPieSeriesList([]float64{
						435.9,
						354.3,
						285.9,
						204.5,
					}, charts.PieSeriesOption{
						Label: charts.SeriesLabel{
							Show: true,
						},
						Radius: "35%",
					}),
				},
			},
		},
	}
	handler(w, req, chartOptions, nil)
}

func echartsHandler(w http.ResponseWriter, req *http.Request) {
	echartsOptions := []string{
		`{
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
						150,
						230,
						224,
						218,
						135,
						147,
						260
					],
					"type": "line"
				}
			]
		}`,
		`{
			"title": {
				"text": "Multiple Line"
			},
			"tooltip": {
				"trigger": "axis"
			},
			"legend": {
				"left": "right",
				"data": [
					"Email",
					"Union Ads",
					"Video Ads",
					"Direct",
					"Search Engine"
				]
			},
			"grid": {
				"left": "3%",
				"right": "4%",
				"bottom": "3%",
				"containLabel": true
			},
			"toolbox": {
				"feature": {
					"saveAsImage": {}
				}
			},
			"xAxis": {
				"type": "category",
				"boundaryGap": false,
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
					"name": "Email",
					"type": "line",
					"data": [
						120,
						132,
						101,
						134,
						90,
						230,
						210
					]
				},
				{
					"name": "Union Ads",
					"type": "line",
					"data": [
						220,
						182,
						191,
						234,
						290,
						330,
						310
					]
				},
				{
					"name": "Video Ads",
					"type": "line",
					"data": [
						150,
						232,
						201,
						154,
						190,
						330,
						410
					]
				},
				{
					"name": "Direct",
					"type": "line",
					"data": [
						320,
						332,
						301,
						334,
						390,
						330,
						320
					]
				},
				{
					"name": "Search Engine",
					"type": "line",
					"data": [
						820,
						932,
						901,
						934,
						1290,
						1330,
						1320
					]
				}
			]
		}`,
		`{
			"title": {
				"text": "Temperature Change in the Coming Week"
			},
			"legend": {
				"left": "right"
			},
			"padding": [10, 30, 10, 10],
			"xAxis": {
				"type": "category",
				"boundaryGap": false,
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
				"axisLabel": {
					"formatter": "{value} °C"
				}
			},
			"series": [
				{
					"name": "Highest",
					"type": "line",
					"data": [
						10,
						11,
						13,
						11,
						12,
						12,
						9
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
					"name": "Lowest",
					"type": "line",
					"data": [
						1,
						-2,
						2,
						5,
						3,
						2,
						0
					],
					"markPoint": {
						"data": [
							{
								"type": "min"
							}
						]
					},
					"markLine": {
						"data": [
							{
								"type": "average"
							},
							{
								"type": "max"
							}
						]
					}
				}
			]
		}`,
		`{
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
						200,
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
		`{
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
		`{
			"title": {
				"text": "World Population"
			},
			"tooltip": {
				"trigger": "axis",
				"axisPointer": {
					"type": "shadow"
				}
			},
			"legend": {},
			"grid": {
				"left": "3%",
				"right": "4%",
				"bottom": "3%",
				"containLabel": true
			},
			"xAxis": {
				"type": "value"
			},
			"yAxis": {
				"type": "category",
				"data": [
					"Brazil",
					"Indonesia",
					"USA",
					"India",
					"China",
					"World"
				]
			},
			"series": [
				{
					"name": "2011",
					"type": "bar",
					"data": [
						18203,
						23489,
						29034,
						104970,
						131744,
						630230
					]
				},
				{
					"name": "2012",
					"type": "bar",
					"data": [
						19325,
						23438,
						31000,
						121594,
						134141,
						681807
					]
				}
			]
		}`,
		`{
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
		`{
			"legend": {
				"data": [
					"Evaporation",
					"Precipitation",
					"Temperature"
				]
			},
			"xAxis": [
				{
					"type": "category",
					"data": [
						"Mon",
						"Tue",
						"Wed",
						"Thu",
						"Fri",
						"Sat",
						"Sun"
					],
					"axisPointer": {
						"type": "shadow"
					}
				}
			],
			"yAxis": [
				{
					"type": "value",
					"name": "Precipitation",
					"min": 0,
					"max": 240,
					"axisLabel": {
						"formatter": "{value} ml"
					}
				},
				{
					"type": "value",
					"name": "Temperature",
					"min": 0,
					"max": 24,
					"axisLabel": {
						"formatter": "{value} °C"
					}
				}
			],
			"series": [
				{
					"name": "Evaporation",
					"type": "bar",
					"tooltip": {},
					"data": [
						2,
						4.9,
						7,
						23.2,
						25.6,
						76.7,
						135.6
					]
				},
				{
					"name": "Precipitation",
					"type": "bar",
					"tooltip": {},
					"data": [
						2.6,
						5.9,
						9,
						26.4,
						28.7,
						70.7,
						175.6
					]
				},
				{
					"name": "Temperature",
					"type": "line",
					"yAxisIndex": 1,
					"tooltip": {},
					"data": [
						2,
						2.2,
						3.3,
						4.5,
						6.3,
						10.2,
						20.3
					]
				}
			]
		}`,
		`{
			"tooltip": {
				"trigger": "axis",
				"axisPointer": {
					"type": "cross"
				}
			},
			"grid": {
				"right": "20%"
			},
			"toolbox": {
				"feature": {
					"dataView": {
						"show": true,
						"readOnly": false
					},
					"restore": {
						"show": true
					},
					"saveAsImage": {
						"show": true
					}
				}
			},
			"legend": {
				"data": [
					"Evaporation",
					"Precipitation",
					"Temperature"
				]
			},
			"xAxis": [
				{
					"type": "category",
					"axisTick": {
						"alignWithLabel": true
					},
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
					"type": "value",
					"name": "温度",
					"position": "left",
					"alignTicks": true,
					"axisLine": {
						"show": true,
						"lineStyle": {
							"color": "#EE6666"
						}
					},
					"axisLabel": {
						"formatter": "{value} °C"
					}
				},
				{
					"type": "value",
					"name": "Evaporation",
					"position": "right",
					"alignTicks": true,
					"axisLine": {
						"show": true,
						"lineStyle": {
							"color": "#5470C6"
						}
					},
					"axisLabel": {
						"formatter": "{value} ml"
					}
				}
			],
			"series": [
				{
					"name": "Evaporation",
					"type": "bar",
					"yAxisIndex": 1,
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
					]
				},
				{
					"name": "Precipitation",
					"type": "bar",
					"yAxisIndex": 1,
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
					]
				},
				{
					"name": "Temperature",
					"type": "line",
					"data": [
						2,
						2.2,
						3.3,
						4.5,
						6.3,
						10.2,
						20.3,
						23.4,
						23,
						16.5,
						12,
						6.2
					]
				}
			]
		}`,
		`{
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
					]
				}
			]
		}`,
		`{
			"title": {
				"text": "Rainfall"
			},
			"padding": [10, 10, 10, 30],
			"legend": {
				"data": [
					"GZ",
					"SH"
				]
			},
			"xAxis": {
				"type": "category",
				"splitNumber": 6,
				"data": [
					"01-01",
					"01-02",
					"01-03",
					"01-04",
					"01-05",
					"01-06",
					"01-07",
					"01-08",
					"01-09",
					"01-10",
					"01-11",
					"01-12",
					"01-13",
					"01-14",
					"01-15",
					"01-16",
					"01-17",
					"01-18",
					"01-19",
					"01-20",
					"01-21",
					"01-22",
					"01-23",
					"01-24",
					"01-25",
					"01-26",
					"01-27",
					"01-28",
					"01-29",
					"01-30",
					"01-31"
				]
			},
			"yAxis": {
				"axisLabel": {
					"formatter": "{value} mm"
				}
			},
			"series": [
				{
					"type": "bar",
					"data": [
						928,
						821,
						889,
						600,
						547,
						783,
						197,
						853,
						430,
						346,
						63,
						465,
						309,
						334,
						141,
						538,
						792,
						58,
						922,
						807,
						298,
						243,
						744,
						885,
						812,
						231,
						330,
						220,
						984,
						221,
						429
					]
				},
				{
					"type": "bar",
					"data": [
						749,
						201,
						296,
						579,
						255,
						159,
						902,
						246,
						149,
						158,
						507,
						776,
						186,
						79,
						390,
						222,
						601,
						367,
						221,
						411,
						714,
						620,
						966,
						73,
						203,
						631,
						833,
						610,
						487,
						677,
						596
					]
				}
			]
		}`,
		`{
			"title": {
				"text": "Basic Radar Chart"
			},
			"legend": {
				"data": [
					"Allocated Budget",
					"Actual Spending"
				]
			},
			"radar": {
				"indicator": [
					{
						"name": "Sales",
						"max": 6500
					},
					{
						"name": "Administration",
						"max": 16000
					},
					{
						"name": "Information Technology",
						"max": 30000
					},
					{
						"name": "Customer Support",
						"max": 38000
					},
					{
						"name": "Development",
						"max": 52000
					},
					{
						"name": "Marketing",
						"max": 25000
					}
				]
			},
			"series": [
				{
					"name": "Budget vs spending",
					"type": "radar",
					"data": [
						{
							"value": [
								4200,
								3000,
								20000,
								35000,
								50000,
								18000
							],
							"name": "Allocated Budget"
						},
						{
							"value": [
								5000,
								14000,
								28000,
								26000,
								42000,
								21000
							],
							"name": "Actual Spending"
						}
					]
				}
			]
		}`,
		`{
			"title": {
				"text": "Funnel"
			},
			"tooltip": {
				"trigger": "item",
				"formatter": "{a} <br/>{b} : {c}%"
			},
			"toolbox": {
				"feature": {
					"dataView": {
						"readOnly": false
					},
					"restore": {},
					"saveAsImage": {}
				}
			},
			"legend": {
				"data": [
					"Show",
					"Click",
					"Visit",
					"Inquiry",
					"Order"
				]
			},
			"series": [
				{
					"name": "Funnel",
					"type": "funnel",
					"left": "10%",
					"top": 60,
					"bottom": 60,
					"width": "80%",
					"min": 0,
					"max": 100,
					"minSize": "0%",
					"maxSize": "100%",
					"sort": "descending",
					"gap": 2,
					"label": {
						"show": true,
						"position": "inside"
					},
					"labelLine": {
						"length": 10,
						"lineStyle": {
							"width": 1,
							"type": "solid"
						}
					},
					"itemStyle": {
						"borderColor": "#fff",
						"borderWidth": 1
					},
					"emphasis": {
						"label": {
							"fontSize": 20
						}
					},
					"data": [
						{
							"value": 60,
							"name": "Visit"
						},
						{
							"value": 40,
							"name": "Inquiry"
						},
						{
							"value": 20,
							"name": "Order"
						},
						{
							"value": 80,
							"name": "Click"
						},
						{
							"value": 100,
							"name": "Show"
						}
					]
				}
			]
		}`,
		`{
			"legend": {
				"top": "-140",
				"data": [
					"Milk Tea",
					"Matcha Latte",
					"Cheese Cocoa",
					"Walnut Brownie"
				]
			},
			"padding": [
				150,
				10,
				10,
				10
			],
			"xAxis": [
				{
					"data": [
						"2012",
						"2013",
						"2014",
						"2015",
						"2016",
						"2017"
					]
				}
			],
			"series": [
				{
					"data": [
						56.5,
						82.1,
						88.7,
						70.1,
						53.4,
						85.1
					]
				},
				{
					"data": [
						51.1,
						51.4,
						55.1,
						53.3,
						73.8,
						68.7
					]
				},
				{
					"data": [
						40.1,
						62.2,
						69.5,
						36.4,
						45.2,
						32.5
					]
				},
				{
					"data": [
						25.2,
						37.1,
						41.2,
						18,
						33.9,
						49.1
					]
				}
			],
			"children": [
				{
					"box": {
						"left": 0,
						"top": 30,
						"right": 600,
						"bottom": 150
					},
					"legend": {
						"show": false		
					},
					"series": [
						{
							"type": "pie",
							"radius": "50%",
							"data": [
								{
									"value": 435.9,
									"name": "Milk Tea"
								},
								{
									"value": 354.3,
									"name": "Matcha Latte"
								},
								{
									"value": 285.9,
									"name": "Cheese Cocoa"
								},
								{
									"value": 204.5,
									"name": "Walnut Brownie"
								}
							]
						}
					]
				}
			]
		}`,
	}
	handler(w, req, nil, echartsOptions)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/echarts", echartsHandler)
	fmt.Println("http://127.0.0.1:3012/")
	http.ListenAndServe(":3012", nil)
}
