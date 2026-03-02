package main

import (
	"math"
	"os"
	"path/filepath"

	"github.com/go-analyze/charts"
)

/*
Example violin chart from sample data using KDE, with median and average mark lines.
Uses NewViolinChartOptionWithSamples with deterministic pseudo-random data generation.
*/

func writeFile(buf []byte) error {
	tmpPath := "./tmp"
	if err := os.MkdirAll(tmpPath, 0700); err != nil {
		return err
	}

	file := filepath.Join(tmpPath, "violin-chart-2-samples.png")
	return os.WriteFile(file, buf, 0600)
}

// lcg is a simple linear congruential generator for deterministic pseudo-random numbers.
type lcg struct {
	state uint64
}

func (l *lcg) next() float64 {
	// Parameters from Numerical Recipes.
	l.state = l.state*6364136223846793005 + 1442695040888963407
	return float64(l.state>>33) / float64(1<<31)
}

// boxMuller generates a standard normal variate from two uniform variates.
func (l *lcg) boxMuller() float64 {
	u1 := l.next()
	u2 := l.next()
	if u1 < 1e-10 {
		u1 = 1e-10
	}
	return math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
}

// normal returns a sample from N(mean, stddev).
func (l *lcg) normal(mean, stddev float64) float64 {
	return mean + stddev*l.boxMuller()
}

func main() {
	const sampleCount = 200
	rng := &lcg{state: 42}

	// Generate four distinct distribution shapes
	normalSamples := make([]float64, sampleCount)
	rightSkewedSamples := make([]float64, sampleCount)
	bimodalSamples := make([]float64, sampleCount)
	tightSamples := make([]float64, sampleCount)
	for i := 0; i < sampleCount; i++ {
		// Normal: symmetric bell curve
		normalSamples[i] = rng.normal(50, 10)

		// Right skewed: lognormal-like
		rightSkewedSamples[i] = 30 + math.Exp(rng.normal(0, 1))*10

		// Bimodal: 45/55% mixture of two clusters
		if rng.next() < 0.45 {
			bimodalSamples[i] = rng.normal(35, 6)
		} else {
			bimodalSamples[i] = rng.normal(65, 6)
		}

		// Tight: low variance
		tightSamples[i] = rng.normal(50, 3)
	}

	samples := [][]float64{normalSamples, rightSkewedSamples, bimodalSamples, tightSamples}
	opt, err := charts.NewViolinChartOptionWithSamples(samples, 80)
	if err != nil {
		panic(err)
	}
	opt.Padding = charts.NewBox(5, 5, 50, 5)
	opt.Title.Text = "Distribution Shapes"
	opt.Legend.SeriesNames = []string{"Normal", "Right Skewed", "Bimodal", "Tight"}
	opt.Legend.Offset = charts.OffsetRight

	// Add median and average mark lines to every series
	for i := range opt.SeriesList {
		opt.SeriesList[i].MarkLine.AddLines(charts.SeriesMarkTypeAverage)
	}

	p := charts.NewPainter(charts.PainterOptions{
		OutputFormat: charts.ChartOutputPNG,
		Width:        1200,
		Height:       800,
	})
	if err := p.ViolinChart(opt); err != nil {
		panic(err)
	} else if buf, err := p.Bytes(); err != nil {
		panic(err)
	} else if err = writeFile(buf); err != nil {
		panic(err)
	}
}
