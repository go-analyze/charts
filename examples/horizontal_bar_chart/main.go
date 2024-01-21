package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	err := os.MkdirAll(tmpPath, 0700)
	if err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "horizontal-bar-chart.png")
	err = ioutil.WriteFile(file, buf, 0600)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	values := [][]float64{
		{
			10,
			30,
			50,
			70,
			90,
			110,
			130,
		},
		{
			20,
			40,
			60,
			80,
			100,
			120,
			140,
		},
	}
	p, err := charts.HorizontalBarRender(
		values,
		charts.TitleTextOptionFunc("World Population"),
		charts.PaddingOptionFunc(charts.Box{
			Top:    20,
			Right:  40,
			Bottom: 20,
			Left:   20,
		}),
		charts.LegendLabelsOptionFunc([]string{
			"2011",
			"2012",
		}),
		charts.YAxisDataOptionFunc([]string{
			"UN",
			"Brazil",
			"Indonesia",
			"USA",
			"India",
			"China",
			"World",
		}),
	)
	if err != nil {
		panic(err)
	}

	buf, err := p.Bytes()
	if err != nil {
		panic(err)
	}
	err = writeFile(buf)
	if err != nil {
		panic(err)
	}
}
