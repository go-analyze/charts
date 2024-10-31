package charts

import (
	"math"
	"strings"

	"github.com/dustin/go-humanize"
)

type SeriesData struct {
	// Value is the retained value for the data.
	Value float64
}

// NewSeriesListDataFromValues returns a series list
func NewSeriesListDataFromValues(values [][]float64, chartType ...string) SeriesList {
	seriesList := make(SeriesList, len(values))
	for index, value := range values {
		seriesList[index] = NewSeriesFromValues(value, chartType...)
	}
	return seriesList
}

// NewSeriesFromValues returns a series
func NewSeriesFromValues(values []float64, chartType ...string) Series {
	s := Series{
		Data: NewSeriesDataFromValues(values),
	}
	if len(chartType) != 0 {
		s.Type = chartType[0]
	}
	return s
}

// NewSeriesDataFromValues return a series data
func NewSeriesDataFromValues(values []float64) []SeriesData {
	data := make([]SeriesData, len(values))
	for index, value := range values {
		data[index] = SeriesData{
			Value: value,
		}
	}
	return data
}

type SeriesLabel struct {
	// Data label formatter, which supports string template.
	// {b}: the name of a data item.
	// {c}: the value of a data item.
	// {d}: the percent of a data item(pie chart).
	Formatter string
	// Color specifies the color for the label.
	Color Color
	// Show flag for label
	Show bool
	// Distance to the host graphic element.
	Distance int
	// Position defines the label position.
	Position string
	// Offset specifies an offset from the position.
	Offset Offset
	// FontSize is the size for text rendering.
	FontSize float64
}

const (
	SeriesMarkDataTypeMax     = "max"
	SeriesMarkDataTypeMin     = "min"
	SeriesMarkDataTypeAverage = "average"
)

type SeriesMarkData struct {
	// Type is the mark data type, it can be "max", "min", "average". "average" is only for mark line.
	Type string
}

type SeriesMarkPoint struct {
	// SymbolSize is the width of symbol, default value is 30.
	SymbolSize int
	// Data is the mark data for the series mark point.
	Data []SeriesMarkData
}

type SeriesMarkLine struct {
	// Data is the mark data for the series mark line.
	Data []SeriesMarkData
}

type Series struct {
	index int
	// Type is the type of series, it can be "line", "bar" or "pie". Default value is "line".
	Type string
	// Data provides the series data list.
	Data []SeriesData
	// YAxisIndex is the index for the axis, it must be 0 or 1.
	YAxisIndex int
	// Label provides the series labels.
	Label SeriesLabel
	// Name specifies a name for the series.
	Name string
	// Radius for Pie chart, e.g.: 40%, default is "40%"
	Radius string
	// MarkPoint provides a series for mark points.
	MarkPoint SeriesMarkPoint
	// MarkLine provides a series for mark lines.
	MarkLine SeriesMarkLine
	// Max value of series
	Min *float64
	// Min value of series
	Max *float64
}
type SeriesList []Series

func (sl SeriesList) init() {
	if len(sl) == 0 || sl[len(sl)-1].index != 0 {
		return // already initialized
	}
	for i := 0; i < len(sl); i++ {
		if sl[i].Type == "" {
			sl[i].Type = ChartTypeLine
		}
		sl[i].index = i
	}
}

func (sl SeriesList) Filter(chartType string) SeriesList {
	arr := make(SeriesList, 0)
	for index, item := range sl {
		if item.Type == chartType {
			arr = append(arr, sl[index])
		}
	}
	return arr
}

// GetMinMax get max and min value of series list
func (sl SeriesList) GetMinMax(axisIndex int) (float64, float64) {
	min := math.MaxFloat64
	max := -math.MaxFloat64
	for _, series := range sl {
		if series.YAxisIndex != axisIndex {
			continue
		}
		for _, item := range series.Data {
			if item.Value == nullValue {
				continue
			}
			if item.Value > max {
				max = item.Value
			}
			if item.Value < min {
				min = item.Value
			}
		}
	}
	return min, max
}

type PieSeriesOption struct {
	Radius string
	Label  SeriesLabel
	Names  []string
}

func NewPieSeriesList(values []float64, opts ...PieSeriesOption) SeriesList {
	result := make([]Series, len(values))
	var opt PieSeriesOption
	if len(opts) != 0 {
		opt = opts[0]
	}
	for index, v := range values {
		name := ""
		if index < len(opt.Names) {
			name = opt.Names[index]
		}
		s := Series{
			Type: ChartTypePie,
			Data: []SeriesData{
				{
					Value: v,
				},
			},
			Radius: opt.Radius,
			Label:  opt.Label,
			Name:   name,
		}
		result[index] = s
	}
	return result
}

type seriesSummary struct {
	// The index of max value
	MaxIndex int
	// The max value
	MaxValue float64
	// The index of min value
	MinIndex int
	// The min value
	MinValue float64
	// THe average value
	AverageValue float64
}

// Summary get summary of series
func (s *Series) Summary() seriesSummary {
	minIndex := -1
	maxIndex := -1
	minValue := math.MaxFloat64
	maxValue := -math.MaxFloat64
	sum := float64(0)
	for j, item := range s.Data {
		if item.Value < minValue {
			minIndex = j
			minValue = item.Value
		}
		if item.Value > maxValue {
			maxIndex = j
			maxValue = item.Value
		}
		sum += item.Value
	}
	return seriesSummary{
		MaxIndex:     maxIndex,
		MaxValue:     maxValue,
		MinIndex:     minIndex,
		MinValue:     minValue,
		AverageValue: sum / float64(len(s.Data)),
	}
}

// Names returns the names of series list
func (sl SeriesList) Names() []string {
	names := make([]string, len(sl))
	for index, s := range sl {
		names[index] = s.Name
	}
	return names
}

// LabelFormatter label formatter
type LabelFormatter func(index int, value float64, percent float64) string

// NewPieLabelFormatter returns a pie label formatter
func NewPieLabelFormatter(seriesNames []string, layout string) LabelFormatter {
	if len(layout) == 0 {
		layout = "{b}: {d}"
	}
	return NewLabelFormatter(seriesNames, layout)
}

// NewFunnelLabelFormatter returns a funner label formatter
func NewFunnelLabelFormatter(seriesNames []string, layout string) LabelFormatter {
	if len(layout) == 0 {
		layout = "{b}({d})"
	}
	return NewLabelFormatter(seriesNames, layout)
}

// NewValueLabelFormatter returns a value formatter
func NewValueLabelFormatter(seriesNames []string, layout string) LabelFormatter {
	if len(layout) == 0 {
		layout = "{c}"
	}
	return NewLabelFormatter(seriesNames, layout)
}

// NewLabelFormatter returns a label formatter
func NewLabelFormatter(seriesNames []string, layout string) LabelFormatter {
	return func(index int, value, percent float64) string {
		percentText := ""
		if percent >= 0 {
			percentText = humanize.FtoaWithDigits(percent*100, 2) + "%"
		}
		valueText := humanize.FtoaWithDigits(value, 2)
		name := ""
		if len(seriesNames) > index {
			name = seriesNames[index]
		}
		text := strings.ReplaceAll(layout, "{c}", valueText)
		text = strings.ReplaceAll(text, "{d}", percentText)
		text = strings.ReplaceAll(text, "{b}", name)
		return text
	}
}
