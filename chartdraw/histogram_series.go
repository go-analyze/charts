package chartdraw

import (
	"errors"
)

// HistogramSeries is a special type of series that draws as a histogram.
// Some peculiarities; it will always be lower bounded at 0 (at the very least).
// This may alter ranges a bit and generally you want to put a histogram series on its own y-axis.
type HistogramSeries struct {
	Name        string
	Style       Style
	YAxis       YAxisType
	InnerSeries ValuesProvider
}

// GetName returns the series name (for Series interface).
func (hs HistogramSeries) GetName() string {
	return hs.Name
}

// GetStyle returns the style used for rendering (for Series interface).
func (hs HistogramSeries) GetStyle() Style {
	return hs.Style
}

// GetYAxis returns which yaxis the series is mapped to.
func (hs HistogramSeries) GetYAxis() YAxisType {
	return hs.YAxis
}

// Len returns the number of values in the inner series (for BoundedValuesProvider interface).
func (hs HistogramSeries) Len() int {
	return hs.InnerSeries.Len()
}

// GetValues proxies value access to the inner series (for ValuesProvider interface).
func (hs HistogramSeries) GetValues(index int) (x, y float64) {
	return hs.InnerSeries.GetValues(index)
}

// GetBoundedValues returns the bounded values used for drawing bars (for BoundedValuesProvider interface).
func (hs HistogramSeries) GetBoundedValues(index int) (x, y1, y2 float64) {
	vx, vy := hs.InnerSeries.GetValues(index)

	x = vx

	if vy > 0 {
		y1 = vy
		return
	}

	y2 = vy
	return
}

// Render draws the histogram series using the given renderer (for Series interface).
func (hs HistogramSeries) Render(r Renderer, canvasBox Box, xrange, yrange Range, defaults Style) {
	style := hs.Style.InheritFrom(defaults)
	Draw.HistogramSeries(r, canvasBox, xrange, yrange, style, hs)
}

// Validate validates the series.
func (hs HistogramSeries) Validate() error {
	if hs.InnerSeries == nil {
		return errors.New("histogram series requires InnerSeries to be set")
	}
	return nil
}
