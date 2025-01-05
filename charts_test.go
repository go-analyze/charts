package charts

import (
	"errors"
	"testing"

	"github.com/go-analyze/charts/chartdraw"
)

func BenchmarkMultiChartPNGRender(b *testing.B) {
	for i := 0; i < b.N; i++ {
		opt := ChartOption{
			OutputFormat: ChartOutputPNG,
			Legend: LegendOption{
				Offset: OffsetStr{
					Top: "-90",
				},
				Data: []string{
					"Milk Tea", "Matcha Latte", "Cheese Cocoa", "Walnut Brownie",
				},
			},
			Padding: chartdraw.Box{
				Top:    100,
				Right:  10,
				Bottom: 10,
				Left:   10,
			},
			XAxis: XAxisOption{
				Data: []string{
					"2012", "2013", "2014", "2015", "2016", "2017",
				},
			},
			YAxis: []YAxisOption{
				{

					Min: FloatPointer(0),
					Max: FloatPointer(90),
				},
			},
			SeriesList: append(
				NewSeriesListLine([][]float64{
					{56.5, 82.1, 88.7, 70.1, 53.4, 85.1},
					{51.1, 51.4, 55.1, 53.3, 73.8, 68.7},
				}),
				NewSeriesListBar([][]float64{
					{40.1, 62.2, 69.5, 36.4, 45.2, 32.5},
					{25.2, 37.1, 41.2, 18, 33.9, 49.1},
				})...),
			Children: []ChartOption{
				{
					Legend: LegendOption{
						Show: False(),
						Data: []string{
							"Milk Tea", "Matcha Latte", "Cheese Cocoa", "Walnut Brownie",
						},
					},
					Box: chartdraw.Box{
						Top:    20,
						Left:   400,
						Right:  500,
						Bottom: 120,
					},
					SeriesList: NewSeriesListPie([]float64{
						435.9, 354.3, 285.9, 204.5,
					}, PieSeriesOption{
						Radius: "35%",
					}),
				},
			},
		}

		if d, err := Render(opt); err != nil {
			panic(err)
		} else if buf, err := d.Bytes(); err != nil {
			panic(err)
		} else if len(buf) == 0 {
			panic(errors.New("data is nil"))
		}
	}
}

func BenchmarkMultiChartSVGRender(b *testing.B) {
	for i := 0; i < b.N; i++ {
		opt := ChartOption{
			Legend: LegendOption{
				Offset: OffsetStr{
					Top: "-90",
				},
				Data: []string{
					"Milk Tea", "Matcha Latte", "Cheese Cocoa", "Walnut Brownie",
				},
			},
			Padding: chartdraw.Box{
				Top:    100,
				Right:  10,
				Bottom: 10,
				Left:   10,
			},
			XAxis: XAxisOption{
				Data: []string{
					"2012",
					"2013",
					"2014",
					"2015",
					"2016",
					"2017",
				},
			},
			YAxis: []YAxisOption{
				{
					Min: FloatPointer(0),
					Max: FloatPointer(90),
				},
			},
			SeriesList: append(
				NewSeriesListLine([][]float64{
					{56.5, 82.1, 88.7, 70.1, 53.4, 85.1},
					{51.1, 51.4, 55.1, 53.3, 73.8, 68.7},
				}),
				NewSeriesListBar([][]float64{
					{40.1, 62.2, 69.5, 36.4, 45.2, 32.5},
					{25.2, 37.1, 41.2, 18, 33.9, 49.1},
				})...),
			Children: []ChartOption{
				{
					Legend: LegendOption{
						Show: False(),
						Data: []string{
							"Milk Tea", "Matcha Latte", "Cheese Cocoa", "Walnut Brownie",
						},
					},
					Box: chartdraw.Box{
						Top:    20,
						Left:   400,
						Right:  500,
						Bottom: 120,
					},
					SeriesList: NewSeriesListPie([]float64{
						435.9, 354.3, 285.9, 204.5,
					}, PieSeriesOption{
						Radius: "35%",
					}),
				},
			},
		}
		if d, err := Render(opt); err != nil {
			panic(err)
		} else if buf, err := d.Bytes(); err != nil {
			panic(err)
		} else if len(buf) == 0 {
			panic(errors.New("data is nil"))
		}
	}
}
