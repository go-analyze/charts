package main

import (
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example violin chart as a population pyramid comparing US and Japan age demographics.
Uses approximate 2020 absolute population counts (people) per single-year age.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "violin-chart-1-basic.png")
	return os.WriteFile(file, buf, 0600)
}

type controlPoint struct {
	age  int
	male float64
	fem  float64
}

// interpolateYearly expands 5-year control points into single-year bands (0 through maxAge-1).
func interpolateYearly(points []controlPoint, maxAge int) [][2]float64 {
	result := make([][2]float64, maxAge)
	for age := 0; age < maxAge; age++ {
		// Find surrounding control points
		var lo, hi int
		for hi = 1; hi < len(points); hi++ {
			if points[hi].age >= age {
				break
			}
		}
		lo = hi - 1
		if lo < 0 {
			lo = 0
		}
		if hi >= len(points) {
			hi = len(points) - 1
		}

		if lo == hi {
			result[age] = [2]float64{points[lo].male, points[lo].fem}
			continue
		}

		// Linear interpolation between lo and hi
		t := float64(age-points[lo].age) / float64(points[hi].age-points[lo].age)
		result[age] = [2]float64{
			points[lo].male + t*(points[hi].male-points[lo].male),
			points[lo].fem + t*(points[hi].fem-points[lo].fem),
		}
	}
	return result
}

func main() {
	// Approximate 2020 control points (5-year midpoints, people per single-year age).
	usPoints := []controlPoint{
		{2, 1_980_000, 1_900_000}, {7, 2_090_000, 2_000_000}, {12, 2_140_000, 2_040_000},
		{17, 2_130_000, 2_040_000}, {22, 2_180_000, 2_090_000}, {27, 2_310_000, 2_220_000},
		{32, 2_250_000, 2_180_000}, {37, 2_130_000, 2_090_000}, {42, 2_000_000, 1_970_000},
		{47, 2_040_000, 2_020_000}, {52, 2_100_000, 2_090_000}, {57, 2_040_000, 2_120_000},
		{62, 1_900_000, 2_000_000}, {67, 1_660_000, 1_810_000}, {72, 1_330_000, 1_500_000},
		{77, 970_000, 1_150_000}, {82, 640_000, 830_000}, {87, 380_000, 570_000},
		{92, 170_000, 290_000}, {97, 50_000, 95_000},
	}
	jpPoints := []controlPoint{
		{2, 980_000, 930_000}, {7, 1_070_000, 1_010_000}, {12, 1_170_000, 1_110_000},
		{17, 1_220_000, 1_170_000}, {22, 1_270_000, 1_220_000}, {27, 1_370_000, 1_320_000},
		{32, 1_470_000, 1_420_000}, {37, 1_570_000, 1_520_000}, {42, 1_680_000, 1_630_000},
		{47, 1_870_000, 1_810_000}, {52, 1_910_000, 1_860_000}, {57, 1_620_000, 1_620_000},
		{62, 1_470_000, 1_520_000}, {67, 1_410_000, 1_510_000}, {72, 1_600_000, 1_780_000},
		{77, 1_300_000, 1_560_000}, {82, 900_000, 1_200_000}, {87, 550_000, 850_000},
		{92, 190_000, 450_000}, {97, 45_000, 150_000},
	}

	const maxAge = 100
	usData := interpolateYearly(usPoints, maxAge)
	jpData := interpolateYearly(jpPoints, maxAge)

	opt := charts.NewViolinChartOptionWithData([][][2]float64{usData, jpData})
	opt.Title.Text = "Population Pyramid"
	opt.Title.Subtext = "Age Distribution by Gender"
	opt.Title.FontStyle = charts.NewFontStyleWithSize(14)
	opt.Title.SubtextFontStyle = charts.NewFontStyleWithSize(10)
	opt.Legend.SeriesNames = []string{"United States", "Japan"}
	opt.Legend.Offset = charts.OffsetRight
	opt.ValueAxis.Limit = charts.Ptr(2_400_000.0)
	opt.ValueAxis.LabelCount = 9

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        800,
		Height:       600,
	})
	if err := p.ViolinChart(opt); err != nil {
		panic(err)
	} else if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}
