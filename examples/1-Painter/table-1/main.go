package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-analyze/charts"
)

/* Example table chart with a variety of basic configuration options shown using the Painter API. */
func writeFile(buf []byte, filename string) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}
	file := filepath.Join(tmpPath, filename)
	return os.WriteFile(file, buf, 0600)
}

func main() {
	// First table chart
	header := []string{
		"Name",
		"Age",
		"Address",
		"Tag",
		"Action",
	}
	data := [][]string{
		{
			"John Brown",
			"32",
			"New York No. 1 Lake Park",
			"nice, developer",
			"Send Mail",
		},
		{
			"Jim Green ",
			"42",
			"London No. 1 Lake Park",
			"wow",
			"Send Mail",
		},
		{
			"Joe Black ",
			"32",
			"Sidney No. 1 Lake Park",
			"cool, teacher",
			"Send Mail",
		},
	}

	// Create painter for the first table
	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        810,
		Height:       200,
	})
	p.FilledRect(0, 0, 810, 300, charts.ColorWhite, charts.ColorWhite, 0.0)

	// Define options for the first table
	tableOpt1 := charts.TableChartOption{
		Header: header,
		Data:   data,
		Spans: []int{
			2, 1, 3, 2, 2,
		},
		Padding: charts.Box{
			Top:    15,
			Right:  10,
			Bottom: 15,
			Left:   10,
		},
	}

	if err := p.TableChart(tableOpt1); err != nil {
		panic(err)
	} else if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err := writeFile(buf, "table-1.png"); err != nil {
		panic(err)
	}

	// Second table chart with styling
	bgColor := charts.ColorRGB(28, 28, 32)
	tableOpt2 := charts.TableChartOption{
		Header:                []string{"Name", "Price", "Change"},
		BackgroundColor:       bgColor,
		HeaderBackgroundColor: charts.ColorRGB(80, 80, 80),
		RowBackgroundColors:   []charts.Color{bgColor},
		HeaderFontColor:       charts.ColorWhite,
		FontStyle: charts.FontStyle{
			FontColor: charts.ColorWhite,
		},
		Padding: charts.Box{
			Top:    15,
			Right:  10,
			Bottom: 15,
			Left:   10,
		},
		Data: [][]string{
			{
				"Datadog Inc",
				"97.32",
				"-7.49%",
			},
			{
				"Hashicorp Inc",
				"28.66",
				"-9.25%",
			},
			{
				"Gitlab Inc",
				"51.63",
				"+4.32%",
			},
		},
		TextAligns: []string{
			"",
			charts.AlignRight,
			charts.AlignRight,
		},
		CellModifier: func(tc charts.TableCell) charts.TableCell {
			tc.FillColor = charts.ColorTransparent
			column := tc.Column
			if column != 2 {
				return tc
			}
			value, _ := strconv.ParseFloat(strings.Replace(tc.Text, "%", "", 1), 64)
			if value == 0 {
				return tc
			}
			if value > 0 {
				tc.FillColor = charts.ColorRGB(179, 53, 20)
			} else {
				tc.FillColor = charts.ColorRGB(33, 124, 50)
			}
			return tc
		},
	}

	// Create painter for the second table
	p2 := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        810,
		Height:       300,
	})
	if err := p2.TableChart(tableOpt2); err != nil {
		panic(err)
	} else if buf, err := p2.Bytes(); err != nil {
		panic(err)
	} else if err := writeFile(buf, "table-1-color.png"); err != nil {
		panic(err)
	}
}
