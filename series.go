package charts

import (
	"math"
	"sort"
	"strings"

	"github.com/dustin/go-humanize"
)

// newSeriesListFromValues returns a series list for the given values and chart type.
func newSeriesListFromValues(values [][]float64, chartType string, label SeriesLabel, names []string,
	radius string, markPoint SeriesMarkPoint, markLine SeriesMarkLine) SeriesList {
	seriesList := make(SeriesList, len(values))
	for index, value := range values {
		s := Series{
			Data:      value,
			Type:      chartType,
			Label:     label,
			Radius:    radius,
			MarkPoint: markPoint,
			MarkLine:  markLine,
		}
		if index < len(names) {
			s.Name = names[index]
		}
		seriesList[index] = s
	}
	return seriesList
}

type SeriesLabel struct {
	// Deprecated: Formatter is deprecated, use FormatTemplate as a direct replacement.
	Formatter string
	// FormatTemplate is a string template for formatting the data label.
	// {b}: the name of a data item.
	// {c}: the value of a data item.
	// {d}: the percent of a data item(pie chart).
	FormatTemplate string
	// ValueFormatter is an alternative method of providing a format for the data label.
	ValueFormatter ValueFormatter
	// FontStyle specifies the font and style for the label.
	FontStyle FontStyle
	// Show flag for label, if unset the behavior will be defaulted based on the chart.
	Show *bool
	// Distance to the host graphic element.
	Distance int // TODO - do we want to replace with just Offset?
	// Deprecated: Position is deprecated, this value was only used on bar and horizontal bar charts. Instead use
	// SeriesLabelPosition on those chart options directly.
	Position string
	// Offset specifies an offset from the position.
	Offset OffsetInt
}

const (
	// Deprecated: SeriesMarkDataTypeMax is deprecated, use SeriesMarkTypeMax.
	SeriesMarkDataTypeMax = SeriesMarkTypeMax
	// Deprecated: SeriesMarkDataTypeMin is deprecated, use SeriesMarkTypeMin.
	SeriesMarkDataTypeMin = SeriesMarkTypeMin
	// Deprecated: SeriesMarkDataTypeAverage is deprecated, use SeriesMarkTypeAverage.
	SeriesMarkDataTypeAverage = SeriesMarkTypeAverage
	SeriesMarkTypeMax         = "max"
	SeriesMarkTypeMin         = "min"
	SeriesMarkTypeAverage     = "average"
)

type SeriesMarkData struct {
	// Type is the mark data type, it can be "max", "min", "average". "average" is only for mark line.
	Type string
}

type SeriesMarkPoint struct {
	// SymbolSize is the width of symbol, default value is 28.
	SymbolSize int
	// ValueFormatter is used to produce the label for the Mark Point.
	ValueFormatter ValueFormatter
	// GlobalPoint specifies optionally that the point should reference the sum of series. This option is only
	// used when the Series is "Stacked" and the point is on the LAST Series of the SeriesList.
	GlobalPoint bool
	// Data is the mark data for the series mark point.
	Data []SeriesMarkData
}

type SeriesMarkLine struct {
	// ValueFormatter is used to produce the label for the Mark Line.
	ValueFormatter ValueFormatter
	// GlobalLine specifies optionally that the line should reference the sum of series. This option is only
	// used when the Series is "Stacked" and the line is on the LAST Series of the SeriesList.
	GlobalLine bool
	// Data is the mark data for the series mark line.
	Data []SeriesMarkData
}

// Series references a population of data.
type Series struct {
	// Type is the type of series, it can be "line", "bar" or "pie". Default value is "line".
	Type string
	// Data provides the series data list.
	Data []float64
	// YAxisIndex is the index for the axis, it must be 0 or 1.
	YAxisIndex int
	// Label provides the series labels.
	Label SeriesLabel
	// Name specifies a name for the series.
	Name string
	// Radius for Pie chart, e.g.: 40%, default is "40%"
	Radius string
	// MarkPoint provides a series for mark points. If Label is also enabled, the MarkPoint will replace the label
	// where rendered.
	MarkPoint SeriesMarkPoint
	// MarkLine provides a series for mark lines. When using a MarkLine, you will want to configure padding to the
	// chart on the right for the values.
	MarkLine SeriesMarkLine
}

// SeriesList is a list of series to be rendered on the chart, typically constructed using NewSeriesListLine,
// NewSeriesListBar, NewSeriesListHorizontalBar, NewSeriesListPie, NewSeriesListRadar, or NewSeriesListFunnel.
// These Series can be appended to each other if multiple chart types should be rendered to the same axis.
type SeriesList []Series

// Deprecated: Filter is deprecated, this function is not expected to be used outside the internal chart
// implementation. If you make use of this function open a GitHub issue to mention its use.
func (sl SeriesList) Filter(chartType string) SeriesList {
	arr := make(SeriesList, 0, len(sl))
	for index, item := range sl {
		if chartTypeMatch(chartType, item.Type) {
			arr = append(arr, sl[index])
		}
	}
	return arr
}

func chartTypeMatch(expected, actual string) bool {
	return expected == "" || expected == actual || (expected == ChartTypeLine && actual == "")
}

func (sl SeriesList) getYAxisCount() int {
	for _, series := range sl {
		if series.YAxisIndex == 1 {
			return 2
		} else if series.YAxisIndex != 0 {
			return -1
		}
	}
	return 1
}

// Deprecated: GetMinMax is deprecated, instead use Series.Summary().  For example seriesList[0].Summary().
func (sl SeriesList) GetMinMax(yaxisIndex int) (float64, float64) {
	min, max, _ := sl.getMinMaxSumMax(yaxisIndex, false)
	return min, max
}

// getMinMaxSumMax returns the min, max, and maximum sum of the series for a given y-axis index (either 0 or 1).
// This is a higher performance option for internal use. calcSum provides an optimization to
// only calculate the sumMax if it will be used.
func (sl SeriesList) getMinMaxSumMax(yaxisIndex int, calcSum bool) (float64, float64, float64) {
	min := math.MaxFloat64
	max := -math.MaxFloat64
	var sums []float64
	if calcSum {
		sums = make([]float64, sl.getMaxDataCount(""))
	}
	for _, series := range sl {
		if series.YAxisIndex != yaxisIndex {
			continue
		}
		for i, item := range series.Data {
			if item == GetNullValue() {
				continue
			}
			if item > max {
				max = item
			}
			if item < min {
				min = item
			}
			if calcSum {
				sums[i] += item
			}
		}
	}
	maxSum := max
	if calcSum {
		for _, val := range sums {
			if val > maxSum {
				maxSum = val
			}
		}
	}
	return min, max, maxSum
}

// LineSeriesOption provides series customization for NewSeriesListLine.
type LineSeriesOption struct {
	Label     SeriesLabel
	Names     []string
	MarkPoint SeriesMarkPoint
	MarkLine  SeriesMarkLine
}

// NewSeriesListLine builds a SeriesList for a line chart. The first dimension of the values indicates the population
// of the data, while the second dimension provides the samples for the population.
func NewSeriesListLine(values [][]float64, opts ...LineSeriesOption) SeriesList {
	var opt LineSeriesOption
	if len(opts) != 0 {
		opt = opts[0]
	}
	return newSeriesListFromValues(values, ChartTypeLine,
		opt.Label, opt.Names, "", opt.MarkPoint, opt.MarkLine)
}

// BarSeriesOption provides series customization for NewSeriesListBar or NewSeriesListHorizontalBar.
type BarSeriesOption struct {
	Label     SeriesLabel
	Names     []string
	MarkPoint SeriesMarkPoint
	MarkLine  SeriesMarkLine
}

// NewSeriesListBar builds a SeriesList for a bar chart. The first dimension of the values indicates the population
// of the data, while the second dimension provides the samples for the population (on the X-Axis).
func NewSeriesListBar(values [][]float64, opts ...BarSeriesOption) SeriesList {
	var opt BarSeriesOption
	if len(opts) != 0 {
		opt = opts[0]
	}
	return newSeriesListFromValues(values, ChartTypeBar,
		opt.Label, opt.Names, "", opt.MarkPoint, opt.MarkLine)
}

// NewSeriesListHorizontalBar builds a SeriesList for a horizontal bar chart. Horizontal bar charts are unique in that
// these Series can not be combined with any other chart type.
func NewSeriesListHorizontalBar(values [][]float64, opts ...BarSeriesOption) SeriesList {
	var opt BarSeriesOption
	if len(opts) != 0 {
		opt = opts[0]
	}
	return newSeriesListFromValues(values, ChartTypeHorizontalBar,
		opt.Label, opt.Names, "", opt.MarkPoint, opt.MarkLine)
}

// PieSeriesOption provides series customization for NewSeriesListPie.
type PieSeriesOption struct {
	// Deprecated: Radius is deprecated, instead set the Radius in PieChartOption.
	Radius string
	Label  SeriesLabel
	Names  []string
}

// NewSeriesListPie builds a SeriesList for a pie chart.
func NewSeriesListPie(values []float64, opts ...PieSeriesOption) SeriesList {
	result := make([]Series, len(values))
	var opt PieSeriesOption
	if len(opts) != 0 {
		opt = opts[0]
	}
	for index, v := range values {
		var name string
		if index < len(opt.Names) {
			name = opt.Names[index]
		}
		s := Series{
			Type:   ChartTypePie,
			Data:   []float64{v},
			Radius: opt.Radius,
			Label:  opt.Label,
			Name:   name,
		}
		result[index] = s
	}
	return result
}

// RadarSeriesOption provides series customization for NewSeriesListRadar.
type RadarSeriesOption struct {
	Label SeriesLabel
	Names []string
}

// NewSeriesListRadar builds a SeriesList for a Radar chart.
func NewSeriesListRadar(values [][]float64, opts ...RadarSeriesOption) SeriesList {
	var opt RadarSeriesOption
	if len(opts) != 0 {
		opt = opts[0]
	}
	return newSeriesListFromValues(values, ChartTypeRadar,
		opt.Label, opt.Names, "", SeriesMarkPoint{}, SeriesMarkLine{})
}

// FunnelSeriesOption provides series customization for NewSeriesListFunnel.
type FunnelSeriesOption struct {
	Label SeriesLabel
	Names []string
}

// NewSeriesListFunnel builds a series list for funnel charts.
func NewSeriesListFunnel(values []float64, opts ...FunnelSeriesOption) SeriesList {
	var opt FunnelSeriesOption
	if len(opts) != 0 {
		opt = opts[0]
	}
	seriesList := make(SeriesList, len(values))
	for index, value := range values {
		var name string
		if index < len(opt.Names) {
			name = opt.Names[index]
		}
		seriesList[index] = Series{
			Data:  []float64{value},
			Type:  ChartTypeFunnel,
			Label: opt.Label,
			Name:  name,
		}
	}
	return seriesList
}

type populationSummary struct {
	// Max is the maximum value in the series.
	Max float64
	// MaxIndex is the index of the maximum value in the series. If the series is empty this value will be -1.
	MaxIndex int
	// Min is the minimum value in the series.
	Min float64
	// MinIndex is the index of the minimum value in the series. If the series is empty this value will be -1.
	MinIndex int
	// Average is the mean of all values in the series.
	Average float64
	// Median is the middle value of the series when it is sorted in ascending order.
	Median float64
	// StandardDeviation is a measure of the amount of variation or dispersion of a set of values. A low standard
	// deviation indicates that the values tend to be close to the mean of the set, while a high standard deviation
	// indicates that the values are spread out over a wider range.
	StandardDeviation float64
	// Skewness measures the asymmetry of the distribution of values in the series around the mean. If skewness is zero,
	// the data are perfectly symmetrical, although not necessarily normal. If skewness is positive, the data is skewed
	// right, meaning that the right tail is longer or fatter than the left. If skewness is negative, the data is skewed
	// left, meaning that the left tail is longer or fatter than the right.
	Skewness float64
	// Kurtosis is a measure of the "tailedness" of the probability distribution of a real-valued random variable.
	// High kurtosis in a data set is an indicator of substantial outliers. A negative kurtosis indicates a relatively flat distribution.
	Kurtosis float64
}

// Summary returns numeric summary of series values (population statistics).
func (s *Series) Summary() populationSummary {
	return summarizePopulationData(s.Data)
}

// summarizePopulationData returns numeric summary of series values (population statistics).
func summarizePopulationData(data []float64) populationSummary {
	var minIndex, maxIndex int
	minValue := math.MaxFloat64
	maxValue := -math.MaxFloat64
	var sum, sumSq, sumCu, sumQd float64
	sortedData := make([]float64, 0, len(data))
	for i, x := range data {
		if x == GetNullValue() {
			continue
		}
		sortedData = append(sortedData, x)

		if x < minValue {
			minValue = x
			minIndex = i
		}
		if x > maxValue {
			maxValue = x
			maxIndex = i
		}

		sum += x
		sumSq += x * x
		sumCu += x * x * x
		sumQd += x * x * x * x
	}
	sort.Float64s(sortedData) // sort non-null values for median and other computations
	ni := len(sortedData)
	if ni == 0 {
		return populationSummary{
			MinIndex: -1,
			MaxIndex: -1,
		}
	}
	nf := float64(ni)

	// Compute average (mean)
	mean := sum / nf
	// Compute median: copy the data and sort
	var median float64
	mid := ni / 2
	if ni%2 == 0 {
		median = (sortedData[mid-1] + sortedData[mid]) / 2.0
	} else {
		median = sortedData[mid]
	}

	// Compute population variance = E[X^2] - (E[X])^2
	variance := sumSq/nf - mean*mean
	stdDev := math.Sqrt(variance)
	// Compute population skewness:
	// thirdCentral = Σ x^3 - 3μΣ x^2 + 3μ^2Σ x - nμ^3
	// skewness = thirdCentral / (n * σ^3)
	var skewness float64
	if stdDev != 0 { // zero stdDev will result in a divide by zero
		thirdCentral := sumCu - 3*mean*sumSq + 3*mean*mean*sum - nf*mean*mean*mean
		skewness = thirdCentral / (nf * stdDev * stdDev * stdDev)
	}

	// Compute population excess kurtosis:
	// fourthCentral = Σ x^4
	//                 - 4μΣ x^3
	//                 + 6μ^2Σ x^2
	//                 - 4μ^3Σ x
	//                 + nμ^4
	// kurtosis = (fourthCentral / (n * σ^4))
	// We don't subtract 3 (excess kurtosis) in our implementation.
	var kurtosis float64
	if variance != 0 {
		fourthCentral := sumQd -
			4*mean*sumCu +
			6*mean*mean*sumSq -
			4*mean*mean*mean*sum +
			nf*mean*mean*mean*mean

		kurtosis = fourthCentral / (nf * variance * variance)
	} // else, all points might be the same => kurtosis is undefined

	return populationSummary{
		Max:               maxValue,
		MaxIndex:          maxIndex,
		Min:               minValue,
		MinIndex:          minIndex,
		Average:           mean,
		Median:            median,
		StandardDeviation: stdDev,
		Skewness:          skewness,
		Kurtosis:          kurtosis,
	}
}

// Deprecated: Names is deprecated, this is expected to be used internally only, if you use this function please open
// a GitHub issue to let us know it's useful to you.
func (sl SeriesList) Names() []string {
	names := make([]string, len(sl))
	for index, s := range sl {
		names[index] = s.Name
	}
	return names
}

// SumSeries will return a single Series which represents the sum of the entire SeriesList. This is useful for
// providing global statistics through Series.Summary().
func (sl SeriesList) SumSeries() Series {
	return sl.makeSumSeries("", -1)
}

func (sl SeriesList) makeSumSeries(chartType string, yaxisIndex int) Series {
	result := Series{
		Type: chartType,
	}
	// check for fast path result
	switch len(sl) {
	case 0:
		return result
	case 1:
		if chartTypeMatch(chartType, sl[0].Type) && (yaxisIndex < 0 || sl[0].YAxisIndex == yaxisIndex) {
			return sl[0]
		} else {
			return result
		}
	}

	sumValues := make([]float64, sl.getMaxDataCount(chartType))
	for _, s := range sl {
		if chartTypeMatch(chartType, s.Type) && (yaxisIndex < 0 || s.YAxisIndex == yaxisIndex) {
			result = s // ensure other series values are set into the result
			for i, f := range s.Data {
				if f != GetNullValue() {
					sumValues[i] += f
				}
			}
		}
	}
	result.Data = sumValues
	return result
}

func (sl SeriesList) getMaxDataCount(chartType string) int {
	result := 0
	for _, s := range sl {
		if chartTypeMatch(chartType, s.Type) {
			count := len(s.Data)
			if count > result {
				result = count
			}
		}
	}
	return result
}

// labelFormatPie formats the value for a pie chart label.
func labelFormatPie(seriesNames []string, layout string, index int, value float64, percent float64) string {
	if len(layout) == 0 {
		layout = "{b}: {d}"
	}
	return newLabelFormatter(seriesNames, layout)(index, value, percent)
}

// labelFormatFunnel formats the value for a funnel chart label.
func labelFormatFunnel(seriesNames []string, layout string, index int, value float64, percent float64) string {
	if len(layout) == 0 {
		layout = "{b}({d})"
	}
	return newLabelFormatter(seriesNames, layout)(index, value, percent)
}

// labelFormatValue returns a formatted value.
func labelFormatValue(seriesNames []string, layout string, index int, value float64, percent float64) string {
	if len(layout) == 0 {
		layout = "{c}"
	}
	return newLabelFormatter(seriesNames, layout)(index, value, percent)
}

// newLabelFormatter returns a label formatter.
func newLabelFormatter(seriesNames []string, layout string) func(index int, value float64, percent float64) string {
	return func(index int, value, percent float64) string {
		var percentText string
		if percent >= 0 {
			percentText = humanize.FtoaWithDigits(percent*100, 2) + "%"
		}
		valueText := humanize.FtoaWithDigits(value, 2)
		var name string
		if len(seriesNames) > index {
			name = seriesNames[index]
		}
		text := strings.ReplaceAll(layout, "{c}", valueText)
		text = strings.ReplaceAll(text, "{d}", percentText)
		text = strings.ReplaceAll(text, "{b}", name)
		return text
	}
}
