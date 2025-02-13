package charts

import (
	"math"
	"sort"
	"strings"

	"github.com/dustin/go-humanize"

	"github.com/go-analyze/charts/chartdraw"
)

type SeriesLabel struct {
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
	// Offset specifies an offset from the position.
	Offset OffsetInt
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

// GenericSeries references a population of data for any type of charts. The chart specific fields will only be active
// for chart types which support them.
type GenericSeries struct {
	// Type is the type of series, it can be "line", "bar" or "pie". Default value is "line".
	Type string
	// Values provides the series data values.
	Values []float64
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

func (g *GenericSeries) getYAxisIndex() int {
	return g.YAxisIndex
}

func (g *GenericSeries) getValues() []float64 {
	return g.Values
}

func (g *GenericSeries) getType() string {
	return g.Type
}

// GenericSeriesList provides the data populations for any chart type configured through ChartOption.
type GenericSeriesList []GenericSeries

func (g GenericSeriesList) names() []string {
	return seriesNames(g)
}

func (g GenericSeriesList) len() int {
	return len(g)
}

func (g GenericSeriesList) getSeries(index int) series {
	return &g[index]
}

func (g GenericSeriesList) getSeriesName(index int) string {
	return g[index].Name
}

func (g GenericSeriesList) getSeriesValues(index int) []float64 {
	return g[index].Values
}

func (g GenericSeriesList) hasMarkPoint() bool {
	for _, s := range g {
		if len(s.MarkPoint.Data) > 0 {
			return true
		}
	}
	return false
}

func (g GenericSeriesList) setSeriesName(index int, name string) {
	g[index].Name = name
}

func (g GenericSeriesList) sortByNameIndex(dict map[string]int) {
	sort.Slice(g, func(i, j int) bool {
		return dict[g[i].Name] < dict[g[j].Name]
	})
}

// LineSeries references a population of data for line charts.
type LineSeries struct {
	// Values provides the series data values.
	Values []float64
	// YAxisIndex is the index for the axis, it must be 0 or 1.
	YAxisIndex int
	// Label provides the series labels.
	Label SeriesLabel
	// Name specifies a name for the series.
	Name string
	// MarkPoint provides a series for mark points. If Label is also enabled, the MarkPoint will replace the label
	// where rendered.
	MarkPoint SeriesMarkPoint
	// MarkLine provides a series for mark lines. When using a MarkLine, you will want to configure padding to the
	// chart on the right for the values.
	MarkLine SeriesMarkLine
}

func (l *LineSeries) getYAxisIndex() int {
	return l.YAxisIndex
}

func (l *LineSeries) getValues() []float64 {
	return l.Values
}

func (l *LineSeries) getType() string {
	return ChartTypeLine
}

func (l *LineSeries) Summary() populationSummary {
	return summarizePopulationData(l.Values)
}

// LineSeriesList provides the data populations for line charts (LineChartOption).
type LineSeriesList []LineSeries

func (l LineSeriesList) names() []string {
	return seriesNames(l)
}

func (l LineSeriesList) SumSeries() []float64 {
	return sumSeriesData(l, -1)
}

func (l LineSeriesList) len() int {
	return len(l)
}

func (l LineSeriesList) getSeries(index int) series {
	return &l[index]
}

func (l LineSeriesList) getSeriesName(index int) string {
	return l[index].Name
}

func (l LineSeriesList) getSeriesValues(index int) []float64 {
	return l[index].Values
}

func (l LineSeriesList) hasMarkPoint() bool {
	for _, s := range l {
		if len(s.MarkPoint.Data) > 0 {
			return true
		}
	}
	return false
}

func (l LineSeriesList) setSeriesName(index int, name string) {
	l[index].Name = name
}

func (l LineSeriesList) sortByNameIndex(dict map[string]int) {
	sort.Slice(l, func(i, j int) bool {
		return dict[l[i].Name] < dict[l[j].Name]
	})
}

func (l LineSeriesList) ToGenericSeriesList() GenericSeriesList {
	result := make([]GenericSeries, len(l))
	for i, s := range l {
		result[i] = GenericSeries{
			Values:     s.Values,
			YAxisIndex: s.YAxisIndex,
			Label:      s.Label,
			Name:       s.Name,
			Type:       ChartTypeLine,
			MarkLine:   s.MarkLine,
			MarkPoint:  s.MarkPoint,
		}
	}
	return result
}

// BarSeries references a population of data for bar charts.
type BarSeries struct {
	// Values provides the series data values.
	Values []float64
	// YAxisIndex is the index for the axis, it must be 0 or 1.
	YAxisIndex int
	// Label provides the series labels.
	Label SeriesLabel
	// Name specifies a name for the series.
	Name string
	// MarkPoint provides a series for mark points. If Label is also enabled, the MarkPoint will replace the label
	// where rendered.
	MarkPoint SeriesMarkPoint
	// MarkLine provides a series for mark lines. When using a MarkLine, you will want to configure padding to the
	// chart on the right for the values.
	MarkLine SeriesMarkLine
}

func (b *BarSeries) getYAxisIndex() int {
	return b.YAxisIndex
}

func (b *BarSeries) getValues() []float64 {
	return b.Values
}

func (b *BarSeries) getType() string {
	return ChartTypeBar
}

func (b *BarSeries) Summary() populationSummary {
	return summarizePopulationData(b.Values)
}

// BarSeriesList provides the data populations for line charts (BarChartOption).
type BarSeriesList []BarSeries

func (b BarSeriesList) names() []string {
	return seriesNames(b)
}

func (b BarSeriesList) SumSeries() []float64 {
	return sumSeriesData(b, -1)
}

func (b BarSeriesList) len() int {
	return len(b)
}

func (b BarSeriesList) getSeries(index int) series {
	return &b[index]
}

func (b BarSeriesList) getSeriesName(index int) string {
	return b[index].Name
}

func (b BarSeriesList) getSeriesValues(index int) []float64 {
	return b[index].Values
}

func (b BarSeriesList) hasMarkPoint() bool {
	for _, s := range b {
		if len(s.MarkPoint.Data) > 0 {
			return true
		}
	}
	return false
}

func (b BarSeriesList) setSeriesName(index int, name string) {
	b[index].Name = name
}

func (b BarSeriesList) sortByNameIndex(dict map[string]int) {
	sort.Slice(b, func(i, j int) bool {
		return dict[b[i].Name] < dict[b[j].Name]
	})
}

func (b BarSeriesList) ToGenericSeriesList() GenericSeriesList {
	result := make([]GenericSeries, len(b))
	for i, s := range b {
		result[i] = GenericSeries{
			Values:     s.Values,
			YAxisIndex: s.YAxisIndex,
			Label:      s.Label,
			Name:       s.Name,
			Type:       ChartTypeBar,
			MarkLine:   s.MarkLine,
			MarkPoint:  s.MarkPoint,
		}
	}
	return result
}

// HorizontalBarSeries references a population of data for horizontal bar charts.
type HorizontalBarSeries struct {
	// Values provides the series data values.
	Values []float64
	// Label provides the series labels.
	Label SeriesLabel
	// Name specifies a name for the series.
	Name string
}

func (h *HorizontalBarSeries) getYAxisIndex() int {
	return 0
}

func (h *HorizontalBarSeries) getValues() []float64 {
	return h.Values
}

func (h *HorizontalBarSeries) getType() string {
	return ChartTypeHorizontalBar
}

func (h *HorizontalBarSeries) Summary() populationSummary {
	return summarizePopulationData(h.Values)
}

// HorizontalBarSeriesList provides the data populations for horizontal bar charts (HorizontalBarChartOption).
type HorizontalBarSeriesList []HorizontalBarSeries

func (h HorizontalBarSeriesList) names() []string {
	return seriesNames(h)
}

func (h HorizontalBarSeriesList) SumSeries() []float64 {
	return sumSeriesData(h, -1)
}

func (h HorizontalBarSeriesList) len() int {
	return len(h)
}

func (h HorizontalBarSeriesList) getSeries(index int) series {
	return &h[index]
}

func (h HorizontalBarSeriesList) getSeriesName(index int) string {
	return h[index].Name
}

func (h HorizontalBarSeriesList) getSeriesValues(index int) []float64 {
	return h[index].Values
}

func (h HorizontalBarSeriesList) hasMarkPoint() bool {
	return false // not currently supported on this chart type
}

func (h HorizontalBarSeriesList) setSeriesName(index int, name string) {
	h[index].Name = name
}

func (h HorizontalBarSeriesList) sortByNameIndex(dict map[string]int) {
	sort.Slice(h, func(i, j int) bool {
		return dict[h[i].Name] < dict[h[j].Name]
	})
}

func (h HorizontalBarSeriesList) ToGenericSeriesList() GenericSeriesList {
	result := make([]GenericSeries, len(h))
	for i, s := range h {
		result[i] = GenericSeries{
			Values: s.Values,
			Label:  s.Label,
			Name:   s.Name,
			Type:   ChartTypeHorizontalBar,
		}
	}
	return result
}

// FunnelSeries references a population of data for funnel charts.
type FunnelSeries struct {
	// Value provides the value for the funnel section.
	Value float64
	// Label provides the series labels.
	Label SeriesLabel
	// Name specifies a name for the series.
	Name string
}

func (f *FunnelSeries) getYAxisIndex() int {
	return 0
}

func (f *FunnelSeries) getValues() []float64 {
	return []float64{f.Value}
}

func (f *FunnelSeries) getType() string {
	return ChartTypeFunnel
}

// FunnelSeriesList provides the data populations for funnel charts (FunnelChartOption).
type FunnelSeriesList []FunnelSeries

func (f FunnelSeriesList) names() []string {
	return seriesNames(f)
}

func (f FunnelSeriesList) len() int {
	return len(f)
}

func (f FunnelSeriesList) getSeries(index int) series {
	return &f[index]
}

func (f FunnelSeriesList) getSeriesName(index int) string {
	return f[index].Name
}

func (f FunnelSeriesList) getSeriesValues(index int) []float64 {
	return []float64{f[index].Value}
}

func (f FunnelSeriesList) hasMarkPoint() bool {
	return false // not supported on this chart type
}

func (f FunnelSeriesList) setSeriesName(index int, name string) {
	f[index].Name = name
}

func (f FunnelSeriesList) sortByNameIndex(dict map[string]int) {
	sort.Slice(f, func(i, j int) bool {
		return dict[f[i].Name] < dict[f[j].Name]
	})
}

func (f FunnelSeriesList) ToGenericSeriesList() GenericSeriesList {
	result := make([]GenericSeries, len(f))
	for i, s := range f {
		result[i] = GenericSeries{
			Values: []float64{s.Value},
			Label:  s.Label,
			Name:   s.Name,
			Type:   ChartTypeFunnel,
		}
	}
	return result
}

// PieSeries references a population of data for pie charts.
type PieSeries struct {
	// Value provides the value for the pie section.
	Value float64
	// Label provides the series labels.
	Label SeriesLabel
	// Name specifies a name for the series.
	Name string
	// Radius for Pie chart, e.g.: 40%, default is "40%"
	Radius string
}

func (p *PieSeries) getYAxisIndex() int {
	return 0
}

func (p *PieSeries) getValues() []float64 {
	return []float64{p.Value}
}

func (p *PieSeries) getType() string {
	return ChartTypePie
}

// PieSeriesList provides the data populations for pie charts (PieChartOption).
type PieSeriesList []PieSeries

func (p PieSeriesList) names() []string {
	return seriesNames(p)
}

func (p PieSeriesList) len() int {
	return len(p)
}

func (p PieSeriesList) getSeries(index int) series {
	return &p[index]
}

func (p PieSeriesList) getSeriesName(index int) string {
	return p[index].Name
}

func (p PieSeriesList) getSeriesValues(index int) []float64 {
	return []float64{p[index].Value}
}

func (p PieSeriesList) hasMarkPoint() bool {
	return false // not supported on this chart type
}

func (p PieSeriesList) setSeriesName(index int, name string) {
	p[index].Name = name
}

func (p PieSeriesList) sortByNameIndex(dict map[string]int) {
	sort.Slice(p, func(i, j int) bool {
		return dict[p[i].Name] < dict[p[j].Name]
	})
}

func (p PieSeriesList) ToGenericSeriesList() GenericSeriesList {
	result := make([]GenericSeries, len(p))
	for i, s := range p {
		result[i] = GenericSeries{
			Values: []float64{s.Value},
			Label:  s.Label,
			Name:   s.Name,
			Type:   ChartTypePie,
			Radius: s.Radius,
		}
	}
	return result
}

// RadarSeries references a population of data for radar charts.
type RadarSeries struct {
	// Values provides the series data list.
	Values []float64
	// Label provides the series labels.
	Label SeriesLabel
	// Name specifies a name for the series.
	Name string
}

func (r *RadarSeries) getYAxisIndex() int {
	return 0
}

func (r *RadarSeries) getValues() []float64 {
	return r.Values
}

func (r *RadarSeries) getType() string {
	return ChartTypeRadar
}

// RadarSeriesList provides the data populations for line charts (RadarChartOption).
type RadarSeriesList []RadarSeries

func (r RadarSeriesList) names() []string {
	return seriesNames(r)
}

func (r RadarSeriesList) len() int {
	return len(r)
}

func (r RadarSeriesList) getSeries(index int) series {
	return &r[index]
}

func (r RadarSeriesList) getSeriesName(index int) string {
	return r[index].Name
}

func (r RadarSeriesList) getSeriesValues(index int) []float64 {
	return r[index].Values
}

func (r RadarSeriesList) hasMarkPoint() bool {
	return false // not supported on this chart type
}

func (r RadarSeriesList) setSeriesName(index int, name string) {
	r[index].Name = name
}

func (r RadarSeriesList) sortByNameIndex(dict map[string]int) {
	sort.Slice(r, func(i, j int) bool {
		return dict[r[i].Name] < dict[r[j].Name]
	})
}

func (r RadarSeriesList) ToGenericSeriesList() GenericSeriesList {
	result := make([]GenericSeries, len(r))
	for i, s := range r {
		result[i] = GenericSeries{
			Values: s.Values,
			Label:  s.Label,
			Name:   s.Name,
			Type:   ChartTypeRadar,
		}
	}
	return result
}

// seriesList contains internal functions for operations that occur across chart types. Most of this interface usage
// is within `series.go` and `charts.go`.
type seriesList interface {
	len() int
	getSeries(index int) series
	getSeriesName(index int) string
	getSeriesValues(index int) []float64
	names() []string
	hasMarkPoint() bool
	setSeriesName(index int, name string)
	sortByNameIndex(dict map[string]int)
}

// series interface is used to provide the raw series struct to callers of seriesList, allowing direct type checks.
type series interface {
	getType() string
	getYAxisIndex() int
	getValues() []float64
}

func filterSeriesList[T any](sl seriesList, chartType string) T {
	switch chartType {
	case ChartTypeLine:
		result := make(LineSeriesList, 0, sl.len())
		for i := 0; i < sl.len(); i++ {
			s := sl.getSeries(i)
			if chartTypeMatch(chartType, s.getType()) {
				switch v := s.(type) {
				case *LineSeries:
					result = append(result, *v)
				case *GenericSeries:
					result = append(result, LineSeries{
						Values:     v.Values,
						YAxisIndex: v.YAxisIndex,
						Label:      v.Label,
						Name:       v.Name,
						MarkLine:   v.MarkLine,
						MarkPoint:  v.MarkPoint,
					})
				}
			}
		}
		return any(result).(T)
	case ChartTypeBar:
		result := make(BarSeriesList, 0, sl.len())
		for i := 0; i < sl.len(); i++ {
			s := sl.getSeries(i)
			if chartTypeMatch(chartType, s.getType()) {
				switch v := s.(type) {
				case *BarSeries:
					result = append(result, *v)
				case *GenericSeries:
					result = append(result, BarSeries{
						Values:     v.Values,
						YAxisIndex: v.YAxisIndex,
						Label:      v.Label,
						Name:       v.Name,
						MarkLine:   v.MarkLine,
						MarkPoint:  v.MarkPoint,
					})
				}
			}
		}
		return any(result).(T)
	case ChartTypeHorizontalBar:
		result := make(HorizontalBarSeriesList, 0, sl.len())
		for i := 0; i < sl.len(); i++ {
			s := sl.getSeries(i)
			if chartTypeMatch(chartType, s.getType()) {
				switch v := s.(type) {
				case *HorizontalBarSeries:
					result = append(result, *v)
				case *GenericSeries:
					result = append(result, HorizontalBarSeries{
						Values: v.Values,
						Label:  v.Label,
						Name:   v.Name,
					})
				}
			}
		}
		return any(result).(T)
	case ChartTypePie:
		result := make(PieSeriesList, 0, sl.len())
		for i := 0; i < sl.len(); i++ {
			s := sl.getSeries(i)
			if chartTypeMatch(chartType, s.getType()) {
				switch v := s.(type) {
				case *PieSeries:
					result = append(result, *v)
				case *GenericSeries:
					result = append(result, PieSeries{
						Value:  chartdraw.SumFloat64(v.Values...),
						Label:  v.Label,
						Name:   v.Name,
						Radius: v.Radius,
					})
				}
			}
		}
		return any(result).(T)
	case ChartTypeRadar:
		result := make(RadarSeriesList, 0, sl.len())
		for i := 0; i < sl.len(); i++ {
			s := sl.getSeries(i)
			if chartTypeMatch(chartType, s.getType()) {
				switch v := s.(type) {
				case *RadarSeries:
					result = append(result, *v)
				case *GenericSeries:
					result = append(result, RadarSeries{
						Values: v.Values,
						Label:  v.Label,
						Name:   v.Name,
					})
				}
			}
		}
		return any(result).(T)
	case ChartTypeFunnel:
		result := make(FunnelSeriesList, 0, sl.len())
		for i := 0; i < sl.len(); i++ {
			s := sl.getSeries(i)
			if chartTypeMatch(chartType, s.getType()) {
				switch v := s.(type) {
				case *FunnelSeries:
					result = append(result, *v)
				case *GenericSeries:
					result = append(result, FunnelSeries{
						Value: chartdraw.SumFloat64(v.Values...),
						Label: v.Label,
						Name:  v.Name,
					})
				}
			}
		}
		return any(result).(T)
	default:
		result := make(GenericSeriesList, 0, sl.len())
		for i := 0; i < sl.len(); i++ {
			s := sl.getSeries(i)
			if chartTypeMatch(chartType, s.getType()) {
				switch v := s.(type) {
				case *LineSeries:
					result = append(result, GenericSeries{
						Values:     v.Values,
						YAxisIndex: v.YAxisIndex,
						Label:      v.Label,
						Name:       v.Name,
						MarkLine:   v.MarkLine,
						MarkPoint:  v.MarkPoint,
					})
				case *BarSeries:
					result = append(result, GenericSeries{
						Values:     v.Values,
						YAxisIndex: v.YAxisIndex,
						Label:      v.Label,
						Name:       v.Name,
						MarkLine:   v.MarkLine,
						MarkPoint:  v.MarkPoint,
					})
				case *HorizontalBarSeries:
					result = append(result, GenericSeries{
						Values: v.Values,
						Label:  v.Label,
						Name:   v.Name,
					})
				case *PieSeries:
					result = append(result, GenericSeries{
						Values: []float64{v.Value},
						Label:  v.Label,
						Name:   v.Name,
						Radius: v.Radius,
					})
				case *RadarSeries:
					result = append(result, GenericSeries{
						Values: v.Values,
						Label:  v.Label,
						Name:   v.Name,
					})
				case *FunnelSeries:
					result = append(result, GenericSeries{
						Values: []float64{v.Value},
						Label:  v.Label,
						Name:   v.Name,
					})
				case *GenericSeries:
					result = append(result, *v)
				}
			}
		}
		return any(result).(T)
	}
}

func chartTypeMatch(expected, actual string) bool {
	return expected == "" || expected == actual || (expected == ChartTypeLine && actual == "")
}

func getSeriesYAxisCount(sl seriesList) int {
	for i := 0; i < sl.len(); i++ {
		axis := sl.getSeries(i).getYAxisIndex()
		if axis == 1 {
			return 2
		} else if axis != 0 {
			return -1
		}
	}
	return 1
}

// getSeriesMinMaxSumMax returns the min, max, and maximum sum of the series for a given y-axis index (either 0 or 1).
// This is a higher performance option for internal use. calcSum provides an optimization to
// only calculate the sumMax if it will be used.
func getSeriesMinMaxSumMax(sl seriesList, yaxisIndex int, calcSum bool) (float64, float64, float64) {
	min := math.MaxFloat64
	max := -math.MaxFloat64
	var sums []float64
	if calcSum {
		sums = make([]float64, getSeriesMaxDataCount(sl))
	}
	for i := 0; i < sl.len(); i++ {
		series := sl.getSeries(i)
		if series.getYAxisIndex() != yaxisIndex {
			continue
		}
		for i, item := range series.getValues() {
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

// NewSeriesListGeneric returns a Generic series list for the given values and chart type (used in ChartOption).
func NewSeriesListGeneric(values [][]float64, chartType string) GenericSeriesList {
	seriesList := make([]GenericSeries, len(values))
	for index, v := range values {
		seriesList[index] = GenericSeries{
			Values: v,
			Type:   chartType,
		}
	}
	return seriesList
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
func NewSeriesListLine(values [][]float64, opts ...LineSeriesOption) LineSeriesList {
	var opt LineSeriesOption
	if len(opts) != 0 {
		opt = opts[0]
	}

	seriesList := make([]LineSeries, len(values))
	for index, v := range values {
		s := LineSeries{
			Values:    v,
			Label:     opt.Label,
			MarkPoint: opt.MarkPoint,
			MarkLine:  opt.MarkLine,
		}
		if index < len(opt.Names) {
			s.Name = opt.Names[index]
		}
		seriesList[index] = s
	}
	return seriesList
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
func NewSeriesListBar(values [][]float64, opts ...BarSeriesOption) BarSeriesList {
	var opt BarSeriesOption
	if len(opts) != 0 {
		opt = opts[0]
	}

	seriesList := make([]BarSeries, len(values))
	for index, v := range values {
		s := BarSeries{
			Values:    v,
			Label:     opt.Label,
			MarkPoint: opt.MarkPoint,
			MarkLine:  opt.MarkLine,
		}
		if index < len(opt.Names) {
			s.Name = opt.Names[index]
		}
		seriesList[index] = s
	}
	return seriesList
}

// NewSeriesListHorizontalBar builds a SeriesList for a horizontal bar chart. Horizontal bar charts are unique in that
// these Series can not be combined with any other chart type.
func NewSeriesListHorizontalBar(values [][]float64, opts ...BarSeriesOption) HorizontalBarSeriesList {
	var opt BarSeriesOption
	if len(opts) != 0 {
		opt = opts[0]
	}

	seriesList := make([]HorizontalBarSeries, len(values))
	for index, v := range values {
		s := HorizontalBarSeries{
			Values: v,
			Label:  opt.Label,
		}
		if index < len(opt.Names) {
			s.Name = opt.Names[index]
		}
		seriesList[index] = s
	}
	return seriesList
}

// PieSeriesOption provides series customization for NewSeriesListPie.
type PieSeriesOption struct {
	Radius string
	Label  SeriesLabel
	Names  []string
}

// NewSeriesListPie builds a SeriesList for a pie chart.
func NewSeriesListPie(values []float64, opts ...PieSeriesOption) PieSeriesList {
	var opt PieSeriesOption
	if len(opts) != 0 {
		opt = opts[0]
	}

	result := make([]PieSeries, len(values))
	for index, v := range values {
		s := PieSeries{
			Value:  v,
			Label:  opt.Label,
			Radius: opt.Radius,
		}
		if index < len(opt.Names) {
			s.Name = opt.Names[index]
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
func NewSeriesListRadar(values [][]float64, opts ...RadarSeriesOption) RadarSeriesList {
	var opt RadarSeriesOption
	if len(opts) != 0 {
		opt = opts[0]
	}

	result := make([]RadarSeries, len(values))
	for index, v := range values {
		s := RadarSeries{
			Values: v,
			Label:  opt.Label,
		}
		if index < len(opt.Names) {
			s.Name = opt.Names[index]
		}
		result[index] = s
	}
	return result
}

// FunnelSeriesOption provides series customization for NewSeriesListFunnel.
type FunnelSeriesOption struct {
	Label SeriesLabel
	Names []string
}

// NewSeriesListFunnel builds a series list for funnel charts.
func NewSeriesListFunnel(values []float64, opts ...FunnelSeriesOption) FunnelSeriesList {
	var opt FunnelSeriesOption
	if len(opts) != 0 {
		opt = opts[0]
	}

	seriesList := make([]FunnelSeries, len(values))
	for index, value := range values {
		s := FunnelSeries{
			Value: value,
			Label: opt.Label,
		}
		if index < len(opt.Names) {
			s.Name = opt.Names[index]
		}
		seriesList[index] = s
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

// summarizePopulationData returns numeric summary of the values (population statistics).
func summarizePopulationData(data []float64) populationSummary {
	n := float64(len(data))
	if n == 0 {
		return populationSummary{
			MinIndex: -1,
			MaxIndex: -1,
		}
	}

	// Initialize tracking variables
	var minIndex, maxIndex int
	minValue := math.MaxFloat64
	maxValue := -math.MaxFloat64
	// For sums of powers:
	var sum, sumSq, sumCu, sumQd float64

	// Single pass to gather everything we need
	for i, x := range data {
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

	// Compute average (mean)
	mean := sum / n
	// Compute population variance = E[X^2] - (E[X])^2
	variance := sumSq/n - mean*mean
	stdDev := math.Sqrt(variance)
	// Compute median: copy the data and sort
	sortedData := make([]float64, len(data))
	copy(sortedData, data)
	sort.Float64s(sortedData)
	var median float64
	mid := len(sortedData) / 2
	if len(sortedData)%2 == 0 {
		median = (sortedData[mid-1] + sortedData[mid]) / 2.0
	} else {
		median = sortedData[mid]
	}

	// Compute population skewness:
	// thirdCentral = Σ x^3 - 3μΣ x^2 + 3μ^2Σ x - nμ^3
	// skewness = thirdCentral / (n * σ^3)
	var skewness float64
	if stdDev != 0 { // zero stdDev will result in a divide by zero
		thirdCentral := sumCu - 3*mean*sumSq + 3*mean*mean*sum - n*mean*mean*mean
		skewness = thirdCentral / (n * stdDev * stdDev * stdDev)
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
			n*mean*mean*mean*mean

		kurtosis = fourthCentral / (n * variance * variance)
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

// seriesNames returns the names of series list.
func seriesNames(sl seriesList) []string {
	names := make([]string, sl.len())
	for index := range names {
		names[index] = sl.getSeriesName(index)
	}
	return names
}

func sumSeriesData(sl seriesList, yaxisIndex int) []float64 {
	seriesLen := sl.len()
	// check for fast path result
	switch seriesLen {
	case 0:
		return make([]float64, 0)
	case 1:
		s := sl.getSeries(0)
		if yaxisIndex < 0 || s.getYAxisIndex() == yaxisIndex {
			return s.getValues()
		}
	}

	sumValues := make([]float64, getSeriesMaxDataCount(sl))
	for i1 := 0; i1 < seriesLen; i1++ {
		s := sl.getSeries(i1)
		if yaxisIndex > -1 && s.getYAxisIndex() != yaxisIndex {
			continue
		}
		for i2, f := range s.getValues() {
			if f != GetNullValue() {
				sumValues[i2] += f
			}
		}
	}
	return sumValues
}

func getSeriesMaxDataCount(sl seriesList) int {
	result := 0
	for i := 0; i < sl.len(); i++ {
		count := len(sl.getSeriesValues(i))
		if count > result {
			result = count
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
