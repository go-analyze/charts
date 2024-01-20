package charts

import (
	"errors"
	"testing"

	"github.com/wcharczuk/go-chart/v2"
)

func BenchmarkMultiChartPNGRender(b *testing.B) {
	for i := 0; i < b.N; i++ {
		opt := ChartOption{
			Type: ChartOutputPNG,
			Legend: LegendOption{
				Top: "-90",
				Data: []string{
					"Milk Tea",
					"Matcha Latte",
					"Cheese Cocoa",
					"Walnut Brownie",
				},
			},
			Padding: chart.Box{
				Top:    100,
				Right:  10,
				Bottom: 10,
				Left:   10,
			},
			XAxis: NewXAxisOption([]string{
				"2012",
				"2013",
				"2014",
				"2015",
				"2016",
				"2017",
			}),
			YAxisOptions: []YAxisOption{
				{

					Min: NewFloatPoint(0),
					Max: NewFloatPoint(90),
				},
			},
			SeriesList: []Series{
				NewSeriesFromValues([]float64{
					56.5,
					82.1,
					88.7,
					70.1,
					53.4,
					85.1,
				}),
				NewSeriesFromValues([]float64{
					51.1,
					51.4,
					55.1,
					53.3,
					73.8,
					68.7,
				}),
				NewSeriesFromValues([]float64{
					40.1,
					62.2,
					69.5,
					36.4,
					45.2,
					32.5,
				}, ChartTypeBar),
				NewSeriesFromValues([]float64{
					25.2,
					37.1,
					41.2,
					18,
					33.9,
					49.1,
				}, ChartTypeBar),
			},
			Children: []ChartOption{
				{
					Legend: LegendOption{
						Show: FalseFlag(),
						Data: []string{
							"Milk Tea",
							"Matcha Latte",
							"Cheese Cocoa",
							"Walnut Brownie",
						},
					},
					Box: chart.Box{
						Top:    20,
						Left:   400,
						Right:  500,
						Bottom: 120,
					},
					SeriesList: NewPieSeriesList([]float64{
						435.9,
						354.3,
						285.9,
						204.5,
					}, PieSeriesOption{
						Label: SeriesLabel{
							Show: true,
						},
						Radius: "35%",
					}),
				},
			},
		}
		d, err := Render(opt)
		if err != nil {
			panic(err)
		}
		buf, err := d.Bytes()
		if err != nil {
			panic(err)
		}
		if len(buf) == 0 {
			panic(errors.New("data is nil"))
		}
	}
}

func BenchmarkMultiChartSVGRender(b *testing.B) {
	for i := 0; i < b.N; i++ {
		opt := ChartOption{
			Legend: LegendOption{
				Top: "-90",
				Data: []string{
					"Milk Tea",
					"Matcha Latte",
					"Cheese Cocoa",
					"Walnut Brownie",
				},
			},
			Padding: chart.Box{
				Top:    100,
				Right:  10,
				Bottom: 10,
				Left:   10,
			},
			XAxis: NewXAxisOption([]string{
				"2012",
				"2013",
				"2014",
				"2015",
				"2016",
				"2017",
			}),
			YAxisOptions: []YAxisOption{
				{

					Min: NewFloatPoint(0),
					Max: NewFloatPoint(90),
				},
			},
			SeriesList: []Series{
				NewSeriesFromValues([]float64{
					56.5,
					82.1,
					88.7,
					70.1,
					53.4,
					85.1,
				}),
				NewSeriesFromValues([]float64{
					51.1,
					51.4,
					55.1,
					53.3,
					73.8,
					68.7,
				}),
				NewSeriesFromValues([]float64{
					40.1,
					62.2,
					69.5,
					36.4,
					45.2,
					32.5,
				}, ChartTypeBar),
				NewSeriesFromValues([]float64{
					25.2,
					37.1,
					41.2,
					18,
					33.9,
					49.1,
				}, ChartTypeBar),
			},
			Children: []ChartOption{
				{
					Legend: LegendOption{
						Show: FalseFlag(),
						Data: []string{
							"Milk Tea",
							"Matcha Latte",
							"Cheese Cocoa",
							"Walnut Brownie",
						},
					},
					Box: chart.Box{
						Top:    20,
						Left:   400,
						Right:  500,
						Bottom: 120,
					},
					SeriesList: NewPieSeriesList([]float64{
						435.9,
						354.3,
						285.9,
						204.5,
					}, PieSeriesOption{
						Label: SeriesLabel{
							Show: true,
						},
						Radius: "35%",
					}),
				},
			},
		}
		d, err := Render(opt)
		if err != nil {
			panic(err)
		}
		buf, err := d.Bytes()
		if err != nil {
			panic(err)
		}
		if len(buf) == 0 {
			panic(errors.New("data is nil"))
		}
	}
}
