package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/wcharczuk/go-chart/v2/drawing"

	"github.com/go-analyze/charts"
)

/*
Example table chart with a variety of basic configuration options shown.
*/

func writeFile(buf []byte, filename string) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, filename)
	return os.WriteFile(file, buf, 0600)
}

func main() {
	charts.SetDefaultWidth(810)
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
	}
	spans := map[int]int{
		0: 2,
		1: 1,
		// set the span of the third column
		2: 3,
		3: 2,
		4: 2,
	}
	p, err := charts.TableRender(
		header,
		data,
		spans,
	)
	if err != nil {
		panic(err)
	}

	if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf, "table.png"); err != nil {
		panic(err)
	}

	bgColor := charts.Color{R: 16, G: 22, B: 30, A: 255}
	p, err = charts.TableOptionRender(charts.TableChartOption{
		Header: []string{
			"Name",
			"Price",
			"Change",
		},
		BackgroundColor:       bgColor,
		HeaderBackgroundColor: bgColor,
		RowBackgroundColors: []charts.Color{
			bgColor,
		},
		HeaderFontColor: drawing.ColorWhite,
		FontColor:       drawing.ColorWhite,
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
		CellStyle: func(tc charts.TableCell) *charts.Style {
			column := tc.Column
			if column != 2 {
				return nil
			}
			value, _ := strconv.ParseFloat(strings.Replace(tc.Text, "%", "", 1), 64)
			if value == 0 {
				return nil
			}
			style := charts.Style{
				Padding: charts.Box{
					Bottom: 5,
				},
			}
			if value > 0 {
				style.FillColor = charts.Color{R: 179, G: 53, B: 20, A: 255}
			} else if value < 0 {
				style.FillColor = charts.Color{R: 33, G: 124, B: 50, A: 255}
			}
			return &style
		},
	})
	if err != nil {
		panic(err)
	}

	if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf, "table-color.png"); err != nil {
		panic(err)
	}
}
