package chartdraw

import (
	"bytes"
	"testing"

	"github.com/go-analyze/charts/chartdraw/drawing"
)

var benchmarkLineChartXValues = Seq{Sequence: NewLinearSequence().WithStart(1.0).WithEnd(100.0)}.Values()
var benchmarkLineChartYValues = Seq{Sequence: NewRandomSequence().WithLen(100).WithMin(100).WithMax(512)}.Values()

func makeBenchmarkLineChart() Chart {
	graph := Chart{
		Background: Style{
			Padding: Box{
				Top:    50,
				Left:   25,
				Right:  25,
				Bottom: 10,
			},
			FillColor: drawing.ColorFromHex("efefef"),
		},
		Title: "Test Line Chart",
		Series: []Series{
			ContinuousSeries{
				XValues: benchmarkLineChartXValues,
				YValues: benchmarkLineChartYValues,
			},
		},
	}
	graph.Elements = []Renderable{
		Legend(&graph),
	}
	return graph
}

func BenchmarkLineChartPNGRender(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buffer := bytes.Buffer{}
		if err := makeBenchmarkLineChart().Render(PNG, &buffer); err != nil {
			panic(err)
		}
	}
}

func BenchmarkLineChartSVGRender(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buffer := bytes.Buffer{}
		if err := makeBenchmarkLineChart().Render(SVG, &buffer); err != nil {
			panic(err)
		}
	}
}

func makeBenchmarkBarChart() BarChart {
	return BarChart{
		Width: 1024,
		Background: Style{
			Padding: Box{
				Top:    50,
				Left:   25,
				Right:  25,
				Bottom: 10,
			},
			FillColor: drawing.ColorFromHex("efefef"),
		},
		Title: "Test Bar Chart",
		Bars: []Value{
			{Value: 1.0, Label: "One"},
			{Value: 2.0, Label: "Two"},
			{Value: 3.0, Label: "Three"},
			{Value: 4.0, Label: "Four"},
			{Value: 5.0, Label: "Five"},
		},
	}
}

func BenchmarkBarChartPNGRender(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buffer := bytes.Buffer{}
		if err := makeBenchmarkBarChart().Render(PNG, &buffer); err != nil {
			panic(err)
		}
	}
}

func BenchmarkBarChartSVGRender(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buffer := bytes.Buffer{}
		if err := makeBenchmarkBarChart().Render(SVG, &buffer); err != nil {
			panic(err)
		}
	}
}

func makeBenchmarkPieChart() PieChart {
	return PieChart{
		Width:  512,
		Height: 512,
		Background: Style{
			Padding: Box{
				Top:    50,
				Left:   25,
				Right:  25,
				Bottom: 10,
			},
			FillColor: drawing.ColorFromHex("efefef"),
		},
		Title: "Test Pie Chart",
		Values: []Value{
			{Value: 5, Label: "Blue"},
			{Value: 5, Label: "Green"},
			{Value: 4, Label: "Gray"},
			{Value: 4, Label: "Orange"},
			{Value: 3, Label: "Deep Blue"},
			{Value: 3, Label: "??"},
			{Value: 1, Label: "!!"},
		},
	}
}

func BenchmarkPieChartPNGRender(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buffer := bytes.Buffer{}
		if err := makeBenchmarkPieChart().Render(PNG, &buffer); err != nil {
			panic(err)
		}
	}
}

func BenchmarkPieChartSVGRender(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buffer := bytes.Buffer{}
		if err := makeBenchmarkPieChart().Render(SVG, &buffer); err != nil {
			panic(err)
		}
	}
}
