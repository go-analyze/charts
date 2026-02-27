package charts

import (
	"math"
	"sort"

	"github.com/go-analyze/bulk"

	"github.com/go-analyze/charts/chartdraw"
)

const (
	SeriesMarkTypeMax     = "max"
	SeriesMarkTypeMin     = "min"
	SeriesMarkTypeAverage = "average"
)

// SeriesMark describes a single mark line or point type.
type SeriesMark struct {
	// Type is the mark data type: "max", "min", "average". "average" is only for mark line.
	Type string
	// Global specifies the mark references the sum of all series. Only used when
	// the Series is "Stacked" and the mark is on the LAST Series of the SeriesList.
	Global bool
}

// NewSeriesMarkList returns a SeriesMarkList initialized for the given types.
func NewSeriesMarkList(markTypes ...string) SeriesMarkList {
	return appendMarks(nil, false, markTypes)
}

// NewSeriesMarkGlobalList returns a slice of SeriesMark initialized for the given types with the global flag set.
// Global marks reference the sum of all series. Only used when the Series is "Stacked" and the mark is
// on the LAST Series of the SeriesList.
func NewSeriesMarkGlobalList(markTypes ...string) SeriesMarkList {
	return appendMarks(nil, true, markTypes)
}

func appendMarks(m SeriesMarkList, global bool, markTypes []string) SeriesMarkList {
	for _, mt := range markTypes {
		if !hasMarkType(m, global, mt) {
			m = append(m, SeriesMark{
				Type:   mt,
				Global: global,
			})
		}
	}
	return m
}

// SeriesMarkList is a slice of SeriesMark values.
type SeriesMarkList []SeriesMark

func (m SeriesMarkList) splitGlobal() (SeriesMarkList, SeriesMarkList) {
	return bulk.SliceSplit(func(v SeriesMark) bool {
		return !v.Global
	}, m)
}

func (m SeriesMarkList) filterGlobal(global bool) SeriesMarkList {
	return bulk.SliceFilter(func(v SeriesMark) bool {
		return v.Global == global
	}, m)
}

func hasMarkType(seriesMarks []SeriesMark, global bool, typeStr string) bool {
	for _, sm := range seriesMarks {
		if sm.Global == global && sm.Type == typeStr {
			return true
		}
	}
	return false
}

// GenericSeries references a population of data for any chart type. Chart-specific fields are only active
// for chart types that support them.
type GenericSeries struct {
	// Type is the series chart type. Default is "line".
	Type string
	// Values provides the series data values.
	// For ChartTypeCandlestick, the Values field must contain OHLC data encoded as groups of 4 consecutive
	// float64 values: [Open, High, Low, Close, ...]. For N candlesticks, Values must have exactly N*4 elements.
	Values []float64
	// YAxisIndex is the y-axis to apply the series to: must be 0 or 1.
	YAxisIndex int
	// Label provides the series labels.
	Label SeriesLabel
	// Name specifies a name for the series.
	Name string
	// Radius for circular charts. Default is "40%".
	Radius string
	// MarkPoint provides a mark point configuration for this series. If Label is enabled, MarkPoint
	// replaces the label where rendered.
	MarkPoint SeriesMarkPoint
	// MarkLine provides amark line configuration for this series. When using MarkLine, configure
	// padding on the chart's right side to ensure space for the values.
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

func (g GenericSeriesList) getSeriesLen(index int) int {
	if g[index].Type == ChartTypeCandlestick {
		if l := len(g[index].Values); l%4 == 0 {
			return l / 4
		} else {
			return l // invalid OHLC format, each value will get its own candle
		}
	}
	return len(g[index].Values)
}

func (g GenericSeriesList) getSeriesSymbol(i int) Symbol {
	if g[i].Type == ChartTypeCandlestick {
		return symbolCandlestick // return type here so captured for defaultRender in ChartOptions render
	}
	return ""
}

func (g GenericSeriesList) hasMarkPoint() bool {
	for _, s := range g {
		if len(s.MarkPoint.Points) > 0 {
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

// SetSeriesLabels sets the label for all elements in the series.
func (g GenericSeriesList) SetSeriesLabels(label SeriesLabel) {
	for i := range g {
		g[i].Label = label
	}
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
	// MarkPoint provides a configuration for mark points for this series. If Label is also enabled, the MarkPoint
	// will replace the label where rendered.
	MarkPoint SeriesMarkPoint
	// MarkLine provides a configuration for mark lines for this series. When using a MarkLine, you will want to
	// configure padding to the chart on the right for the values.
	MarkLine SeriesMarkLine
	// TrendLine provides configurations for trend lines for this series.
	TrendLine []SeriesTrendLine
	// Symbol specifies a custom symbol for the series.
	Symbol Symbol // TODO - v0.6 - consider combining symbol with size into a SymbolStyle struct

	// absThemeIndex represents the series index when combined with other chart types.
	absThemeIndex *int
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

// SumSeries returns a float64 slice with the sum of each series.
func (l LineSeriesList) SumSeries() []float64 {
	return sumSeries(l)
}

// SumSeriesValues returns a float64 slice with each series totaled by the value index.
func (l LineSeriesList) SumSeriesValues() []float64 {
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

func (l LineSeriesList) getSeriesLen(index int) int {
	return len(l[index].Values)
}

func (l LineSeriesList) getSeriesSymbol(index int) Symbol {
	return l[index].Symbol
}

func (l LineSeriesList) hasMarkPoint() bool {
	for _, s := range l {
		if len(s.MarkPoint.Points) > 0 {
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

// SetSeriesLabels sets the label for all elements in the series.
func (l LineSeriesList) SetSeriesLabels(label SeriesLabel) {
	for i := range l {
		l[i].Label = label
	}
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

// ScatterSeries references a population of data for scatter charts.
type ScatterSeries struct {
	// Values provides the series data values.
	Values [][]float64
	// YAxisIndex is the index for the axis, it must be 0 or 1.
	YAxisIndex int
	// Label provides the series labels.
	Label SeriesLabel
	// Name specifies a name for the series.
	Name string
	// MarkLine provides a configuration for mark lines for this series. When using a MarkLine, you will want to
	// configure padding to the chart on the right for the values.
	MarkLine SeriesMarkLine
	// TrendLine provides configurations for trend lines for this series.
	TrendLine []SeriesTrendLine
	// Symbol specifies a custom symbol for the series.
	Symbol Symbol // TODO - v0.6 - consider combining symbol with size into a SymbolStyle struct

	// absThemeIndex represents the series index when combined with other chart types.
	absThemeIndex *int
}

func (s *ScatterSeries) getYAxisIndex() int {
	return s.YAxisIndex
}

func (s *ScatterSeries) getValues() []float64 {
	if len(s.Values) == 0 {
		return nil
	}
	result := make([]float64, 0, len(s.Values)*len(s.Values[0]))
	for _, v := range s.Values {
		result = append(result, v...)
	}
	return result
}

func (s *ScatterSeries) avgValues() []float64 {
	values := make([]float64, len(s.Values))
	for i, v := range s.Values {
		values[i] = chartdraw.MeanFloat64(v...)
	}
	return values
}

func (s *ScatterSeries) getType() string {
	return ChartTypeScatter
}

func (s *ScatterSeries) Summary() populationSummary {
	return summarizePopulationData(s.getValues())
}

// ScatterSeriesList provides the data populations for scatter charts (ScatterChartOption).
type ScatterSeriesList []ScatterSeries

func (s ScatterSeriesList) names() []string {
	return seriesNames(s)
}

// SumSeries returns a float64 slice with the sum of each series.
func (s ScatterSeriesList) SumSeries() []float64 {
	return sumSeries(s)
}

func (s ScatterSeriesList) len() int {
	return len(s)
}

func (s ScatterSeriesList) getSeries(index int) series {
	return &s[index]
}

func (s ScatterSeriesList) getSeriesName(index int) string {
	return s[index].Name
}

func (s ScatterSeriesList) getSeriesValues(index int) []float64 {
	return s[index].getValues()
}

func (s ScatterSeriesList) getSeriesLen(index int) int {
	return len(s[index].Values)
}

func (s ScatterSeriesList) getSeriesSymbol(index int) Symbol {
	return s[index].Symbol
}

func (s ScatterSeriesList) hasMarkPoint() bool {
	return false
}

func (s ScatterSeriesList) setSeriesName(index int, name string) {
	s[index].Name = name
}

func (s ScatterSeriesList) sortByNameIndex(dict map[string]int) {
	sort.Slice(s, func(i, j int) bool {
		return dict[s[i].Name] < dict[s[j].Name]
	})
}

// SetSeriesLabels sets the label for all elements in the series.
func (s ScatterSeriesList) SetSeriesLabels(label SeriesLabel) {
	for i := range s {
		s[i].Label = label
	}
}

func (s ScatterSeriesList) ToGenericSeriesList() GenericSeriesList {
	result := make([]GenericSeries, len(s))
	for i, s := range s {
		result[i] = GenericSeries{
			Values:     s.avgValues(),
			YAxisIndex: s.YAxisIndex,
			Label:      s.Label,
			Name:       s.Name,
			Type:       ChartTypeScatter,
			MarkLine:   s.MarkLine,
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
	// MarkPoint provides a configuration for mark points for this series. If Label is also enabled, the MarkPoint
	// will replace the label where rendered.
	MarkPoint SeriesMarkPoint
	// MarkLine provides a configuration for mark lines for this series. When using a MarkLine, you will want to
	// configure padding to the chart on the right for the values.
	MarkLine SeriesMarkLine

	// absThemeIndex represents the series index when combined with other chart types.
	absThemeIndex *int
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

// BarSeriesList provides the data populations for bar charts (BarChartOption).
type BarSeriesList []BarSeries

func (b BarSeriesList) names() []string {
	return seriesNames(b)
}

// SumSeries returns a float64 slice with the sum of each series.
func (b BarSeriesList) SumSeries() []float64 {
	return sumSeries(b)
}

// SumSeriesValues returns a float64 slice with each series totaled by the value index.
func (b BarSeriesList) SumSeriesValues() []float64 {
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

func (b BarSeriesList) getSeriesLen(index int) int {
	return len(b[index].Values)
}

func (b BarSeriesList) getSeriesSymbol(_ int) Symbol {
	return ""
}

func (b BarSeriesList) hasMarkPoint() bool {
	for _, s := range b {
		if len(s.MarkPoint.Points) > 0 {
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

// SetSeriesLabels sets the label for all elements in the series.
func (b BarSeriesList) SetSeriesLabels(label SeriesLabel) {
	for i := range b {
		b[i].Label = label
	}
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
	// MarkLine provides a configuration for mark lines for this series.
	MarkLine SeriesMarkLine
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

// SumSeries returns a float64 slice with the sum of each series.
func (h HorizontalBarSeriesList) SumSeries() []float64 {
	return sumSeries(h)
}

// SumSeriesValues returns a float64 slice with each series totaled by the value index.
func (h HorizontalBarSeriesList) SumSeriesValues() []float64 {
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

func (h HorizontalBarSeriesList) getSeriesLen(index int) int {
	return len(h[index].Values)
}

func (h HorizontalBarSeriesList) getSeriesSymbol(_ int) Symbol {
	return ""
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

// SetSeriesLabels sets the label for all elements in the series.
func (h HorizontalBarSeriesList) SetSeriesLabels(label SeriesLabel) {
	for i := range h {
		h[i].Label = label
	}
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

func (f FunnelSeriesList) getSeriesLen(_ int) int {
	return 1
}

func (f FunnelSeriesList) getSeriesSymbol(_ int) Symbol {
	return ""
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

// SetSeriesLabels sets the label for all elements in the series.
func (f FunnelSeriesList) SetSeriesLabels(label SeriesLabel) {
	for i := range f {
		f[i].Label = label
	}
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

func (p PieSeriesList) SumSeries() float64 {
	var sum float64
	for _, s := range p {
		sum += s.Value
	}
	return sum
}

// MaxValue returns the maximum value within the series, or MinInt64 if no values.
func (p PieSeriesList) MaxValue() float64 {
	max := float64(math.MinInt64)
	for _, s := range p {
		if s.Value > max {
			max = s.Value
		}
	}
	return max
}

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

func (p PieSeriesList) getSeriesLen(_ int) int {
	return 1
}

func (p PieSeriesList) getSeriesSymbol(_ int) Symbol {
	return ""
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

// SetSeriesLabels sets the label for all elements in the series.
func (p PieSeriesList) SetSeriesLabels(label SeriesLabel) {
	for i := range p {
		p[i].Label = label
	}
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

// DoughnutSeries references a population of data for doughnut charts.
type DoughnutSeries struct {
	// Value provides the value for the Doughnut section.
	Value float64
	// Label provides the series labels.
	Label SeriesLabel
	// Name specifies a name for the series.
	Name string
	// Radius for Doughnut chart, e.g.: 40%, default is "40%"
	Radius string
}

func (d *DoughnutSeries) getYAxisIndex() int {
	return 0
}

func (d *DoughnutSeries) getValues() []float64 {
	return []float64{d.Value}
}

func (d *DoughnutSeries) getType() string {
	return ChartTypeDoughnut
}

// DoughnutSeriesList provides the data populations for Doughnut charts (DoughnutChartOption).
type DoughnutSeriesList []DoughnutSeries

func (d DoughnutSeriesList) SumSeries() float64 {
	var sum float64
	for _, s := range d {
		sum += s.Value
	}
	return sum
}

// MaxValue returns the maximum value within the series, or MinInt64 if no values.
func (d DoughnutSeriesList) MaxValue() float64 {
	max := float64(math.MinInt64)
	for _, s := range d {
		if s.Value > max {
			max = s.Value
		}
	}
	return max
}

func (d DoughnutSeriesList) names() []string {
	return seriesNames(d)
}

func (d DoughnutSeriesList) len() int {
	return len(d)
}

func (d DoughnutSeriesList) getSeries(index int) series {
	return &d[index]
}

func (d DoughnutSeriesList) getSeriesName(index int) string {
	return d[index].Name
}

func (d DoughnutSeriesList) getSeriesValues(index int) []float64 {
	return []float64{d[index].Value}
}

func (d DoughnutSeriesList) getSeriesLen(_ int) int {
	return 1
}

func (d DoughnutSeriesList) getSeriesSymbol(_ int) Symbol {
	return ""
}

func (d DoughnutSeriesList) hasMarkPoint() bool {
	return false // not supported on this chart type
}

func (d DoughnutSeriesList) setSeriesName(index int, name string) {
	d[index].Name = name
}

func (d DoughnutSeriesList) sortByNameIndex(dict map[string]int) {
	sort.Slice(d, func(i, j int) bool {
		return dict[d[i].Name] < dict[d[j].Name]
	})
}

// SetSeriesLabels sets the label for all elements in the series.
func (d DoughnutSeriesList) SetSeriesLabels(label SeriesLabel) {
	for i := range d {
		d[i].Label = label
	}
}

func (d DoughnutSeriesList) ToGenericSeriesList() GenericSeriesList {
	result := make([]GenericSeries, len(d))
	for i, s := range d {
		result[i] = GenericSeries{
			Values: []float64{s.Value},
			Label:  s.Label,
			Name:   s.Name,
			Type:   ChartTypeDoughnut,
			Radius: s.Radius,
		}
	}
	return result
}

func (d DoughnutSeriesList) toPieSeriesList() PieSeriesList {
	result := make([]PieSeries, len(d))
	for i, s := range d {
		result[i] = PieSeries(s)
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

// RadarSeriesList provides the data populations for radar charts (RadarChartOption).
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

func (r RadarSeriesList) getSeriesLen(index int) int {
	return len(r[index].Values)
}

func (r RadarSeriesList) getSeriesSymbol(_ int) Symbol {
	return ""
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

// SetSeriesLabels sets the label for all elements in the series.
func (r RadarSeriesList) SetSeriesLabels(label SeriesLabel) {
	for i := range r {
		r[i].Label = label
	}
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
	getSeriesLen(i int) int
	names() []string
	hasMarkPoint() bool
	setSeriesName(index int, name string)
	sortByNameIndex(dict map[string]int)
	getSeriesSymbol(index int) Symbol
	//SetSeriesLabels(label SeriesLabel) // informally included in interface (not used internally in interface)
}

// series interface is used to provide the raw series struct to callers of seriesList, allowing direct type checks.
type series interface {
	getType() string
	getYAxisIndex() int
	getValues() []float64
}

func expandSingleValueScatterSeries(vals []float64) [][]float64 {
	result := make([][]float64, len(vals))
	for i, v := range vals {
		result[i] = []float64{v}
	}
	return result
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
						Values:        v.Values,
						YAxisIndex:    v.YAxisIndex,
						Label:         v.Label,
						Name:          v.Name,
						MarkLine:      v.MarkLine,
						MarkPoint:     v.MarkPoint,
						absThemeIndex: Ptr(i),
					})
				}
			}
		}
		return any(result).(T)
	case ChartTypeScatter:
		result := make(ScatterSeriesList, 0, sl.len())
		for i := 0; i < sl.len(); i++ {
			s := sl.getSeries(i)
			if chartTypeMatch(chartType, s.getType()) {
				switch v := s.(type) {
				case *ScatterSeries:
					result = append(result, *v)
				case *GenericSeries:
					result = append(result, ScatterSeries{
						Values:        expandSingleValueScatterSeries(v.Values),
						YAxisIndex:    v.YAxisIndex,
						Label:         v.Label,
						Name:          v.Name,
						MarkLine:      v.MarkLine,
						absThemeIndex: Ptr(i),
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
						Values:        v.Values,
						YAxisIndex:    v.YAxisIndex,
						Label:         v.Label,
						Name:          v.Name,
						MarkLine:      v.MarkLine,
						MarkPoint:     v.MarkPoint,
						absThemeIndex: Ptr(i),
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
	case ChartTypeDoughnut:
		result := make(DoughnutSeriesList, 0, sl.len())
		for i := 0; i < sl.len(); i++ {
			s := sl.getSeries(i)
			if chartTypeMatch(chartType, s.getType()) {
				switch v := s.(type) {
				case *DoughnutSeries:
					result = append(result, *v)
				case *GenericSeries:
					result = append(result, DoughnutSeries{
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
	case ChartTypeCandlestick:
		result := make(CandlestickSeriesList, 0, sl.len())
		for i := 0; i < sl.len(); i++ {
			s := sl.getSeries(i)
			if chartTypeMatch(chartType, s.getType()) {
				switch v := s.(type) {
				case *CandlestickSeries:
					result = append(result, *v)
				case *GenericSeries:
					// Convert GenericSeries to CandlestickSeries
					// If Values length is divisible by 4, assume it contains OHLC data
					// Otherwise, create basic OHLC where O=H=L=C (flat line) from single values
					var ohlcData []OHLCData
					if len(v.Values)%4 == 0 && len(v.Values) > 0 {
						// Assume OHLC encoding: groups of 4 consecutive values
						candleCount := len(v.Values) / 4
						ohlcData = make([]OHLCData, candleCount)
						for j := 0; j < candleCount; j++ {
							baseIdx := j * 4
							ohlcData[j] = OHLCData{
								Open:  v.Values[baseIdx],
								High:  v.Values[baseIdx+1],
								Low:   v.Values[baseIdx+2],
								Close: v.Values[baseIdx+3],
							}
						}
					} else {
						// Fallback: create basic OHLC where O=H=L=C (flat line) from single values
						ohlcData = make([]OHLCData, len(v.Values))
						for j, val := range v.Values {
							ohlcData[j] = OHLCData{
								Open:  val,
								High:  val,
								Low:   val,
								Close: val,
							}
						}
					}
					result = append(result, CandlestickSeries{
						Data:           ohlcData,
						YAxisIndex:     v.YAxisIndex,
						Label:          v.Label,
						Name:           v.Name,
						CloseMarkLine:  v.MarkLine,
						CloseMarkPoint: v.MarkPoint,
						absThemeIndex:  Ptr(i),
					})
				}
			}
		}
		return any(result).(T)
	default:
		var zero T
		return zero
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
	// If min was not updated then there were no valid data points. Return
	// zeros to avoid propagating sentinel values like math.MaxFloat64 which
	// can corrupt downstream range calculations.
	if min == math.MaxFloat64 && max == -math.MaxFloat64 {
		return 0, 0, 0
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
	TrendLine []SeriesTrendLine
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
			TrendLine: opt.TrendLine,
		}
		if index < len(opt.Names) {
			s.Name = opt.Names[index]
		}
		seriesList[index] = s
	}
	return seriesList
}

// ScatterSeriesOption provides series customization for NewSeriesListScatter and NewSeriesListScatterMultiValue.
type ScatterSeriesOption struct {
	Label     SeriesLabel
	Names     []string
	MarkLine  SeriesMarkLine
	TrendLine []SeriesTrendLine
}

// NewSeriesListScatter builds a SeriesList for a scatter chart. The first dimension of the values indicates the population
// of the data, while the second dimension provides the samples for the population.
func NewSeriesListScatter(values [][]float64, opts ...ScatterSeriesOption) ScatterSeriesList {
	var opt ScatterSeriesOption
	if len(opts) != 0 {
		opt = opts[0]
	}

	seriesList := make([]ScatterSeries, len(values))
	for index, v := range values {
		s := ScatterSeries{
			Values:    expandSingleValueScatterSeries(v),
			Label:     opt.Label,
			MarkLine:  opt.MarkLine,
			TrendLine: opt.TrendLine,
		}
		if index < len(opt.Names) {
			s.Name = opt.Names[index]
		}
		seriesList[index] = s
	}
	return seriesList
}

// NewSeriesListScatterMultiValue builds a SeriesList for a scatter charts. The first dimension of the values indicates
// the population of the data, while the second dimension provides the samples for the population. Multiple values for
// a single sample can be provided using the last dimension.
func NewSeriesListScatterMultiValue(values [][][]float64, opts ...ScatterSeriesOption) ScatterSeriesList {
	var opt ScatterSeriesOption
	if len(opts) != 0 {
		opt = opts[0]
	}

	seriesList := make([]ScatterSeries, len(values))
	for index, v := range values {
		s := ScatterSeries{
			Values:    v,
			Label:     opt.Label,
			MarkLine:  opt.MarkLine,
			TrendLine: opt.TrendLine,
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
	Label SeriesLabel
	Names []string
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
			Value: v,
			Label: opt.Label,
		}
		if index < len(opt.Names) {
			s.Name = opt.Names[index]
		}
		result[index] = s
	}
	return result
}

// DoughnutSeriesOption provides series customization for NewSeriesListDoughnut.
type DoughnutSeriesOption struct {
	Label SeriesLabel
	Names []string
}

// NewSeriesListDoughnut builds a SeriesList for a doughnut chart.
func NewSeriesListDoughnut(values []float64, opts ...DoughnutSeriesOption) DoughnutSeriesList {
	var opt DoughnutSeriesOption
	if len(opts) != 0 {
		opt = opts[0]
	}

	result := make([]DoughnutSeries, len(values))
	for index, v := range values {
		s := DoughnutSeries{
			Value: v,
			Label: opt.Label,
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

// OHLCData represents Open, High, Low, Close financial data for a single time period.
// All values must satisfy: High >= Open, Close and Low <= Open, Close for valid candlesticks.
type OHLCData struct {
	// Open is the opening price for the time period.
	Open float64
	// High is the highest price during the time period.
	High float64
	// Low is the lowest price during the time period.
	Low float64
	// Close is the closing price for the time period.
	Close float64
}

const (
	// CandleStyleFilled always fills bodies.
	CandleStyleFilled = "filled"
	// CandleStyleTraditional uses hollow bullish, filled bearish.
	CandleStyleTraditional = "traditional"
	// CandleStyleOutline always outlines only.
	CandleStyleOutline = "outline"
)

// CandlestickSeries references OHLC data for candlestick charts.
type CandlestickSeries struct {
	// Data provides OHLC data for each time period.
	Data []OHLCData
	// YAxisIndex is the index for the axis, it must be 0 or 1.
	YAxisIndex int
	// Label provides the series labels.
	Label SeriesLabel
	// Name specifies a name for the series.
	Name string

	// OpenMarkPoint provides mark points for open values.
	OpenMarkPoint SeriesMarkPoint
	// OpenMarkLine provides mark lines for open values.
	OpenMarkLine SeriesMarkLine
	// OpenTrendLine provides trend lines for open values.
	OpenTrendLine []SeriesTrendLine
	// HighMarkPoint provides mark points for high values.
	HighMarkPoint SeriesMarkPoint
	// HighMarkLine provides mark lines for high values.
	HighMarkLine SeriesMarkLine
	// HighTrendLine provides trend lines for high values.
	HighTrendLine []SeriesTrendLine
	// LowMarkPoint provides mark points for low values.
	LowMarkPoint SeriesMarkPoint
	// LowMarkLine provides mark lines for low values.
	LowMarkLine SeriesMarkLine
	// LowTrendLine provides trend lines for low values.
	LowTrendLine []SeriesTrendLine
	// CloseMarkPoint provides mark points for close values.
	CloseMarkPoint SeriesMarkPoint
	// CloseMarkLine provides mark lines for close values.
	CloseMarkLine SeriesMarkLine
	// CloseTrendLine provides trend lines for close values.
	CloseTrendLine []SeriesTrendLine

	// ShowWicks hides wicks when false (body only). Overrides chart-level setting.
	ShowWicks *bool
	// CandleStyle specifies the visual style: CandleStyleFilled, CandleStyleTraditional, or CandleStyleOutline.
	CandleStyle string
	// PatternConfig configures automatic pattern detection and labeling.
	PatternConfig *CandlestickPatternConfig

	// absThemeIndex represents the series index when combined with other chart types.
	absThemeIndex *int
}

func (k *CandlestickSeries) getYAxisIndex() int {
	return k.YAxisIndex
}

// getValues returns only High and Low values for range calculations.
// Candlesticks cannot be meaningfully summed or aggregated, but the full price range
// (High/Low) is needed for axis scaling.
func (k *CandlestickSeries) getValues() []float64 {
	// candlestick can't be summed or flattened, so values are only used for range
	// for that reason we only add the high and low values here
	result := make([]float64, 0, len(k.Data)*2)
	for _, ohlc := range k.Data {
		result = append(result, ohlc.High, ohlc.Low)
	}
	return result
}

func (k *CandlestickSeries) getType() string {
	return ChartTypeCandlestick
}

func (k *CandlestickSeries) Summary() populationSummary {
	return summarizePopulationData(k.getValues())
}

// validateOHLCData ensures OHLC data follows financial market rules.
func validateOHLCData(ohlc OHLCData) bool {
	return validateOHLCOpen(ohlc) && validateOHLCClose(ohlc)
}

// validateOHLCHighLow validates that High >= Low and neither is null.
func validateOHLCHighLow(ohlc OHLCData) bool {
	if ohlc.High == GetNullValue() || ohlc.Low == GetNullValue() {
		return false
	}
	return ohlc.High >= ohlc.Low
}

// validateOHLCOpen validates that Open is within the High-Low range.
func validateOHLCOpen(ohlc OHLCData) bool {
	if ohlc.Open == GetNullValue() || !validateOHLCHighLow(ohlc) {
		return false
	}
	return ohlc.High >= ohlc.Open && ohlc.Low <= ohlc.Open
}

// validateOHLCClose validates that Close is within the High-Low range.
func validateOHLCClose(ohlc OHLCData) bool {
	if ohlc.Close == GetNullValue() || !validateOHLCHighLow(ohlc) {
		return false
	}
	return ohlc.High >= ohlc.Close && ohlc.Low <= ohlc.Close
}

// CandlestickSeriesList holds multiple CandlestickSeries values.
type CandlestickSeriesList []CandlestickSeries

func (k CandlestickSeriesList) names() []string {
	return seriesNames(k)
}

func (k CandlestickSeriesList) len() int {
	return len(k)
}

// SumSeries returns a float64 slice with the sum of each series.
// For candlestick series, this sums the High and Low values (used for range calculations).
func (k CandlestickSeriesList) SumSeries() []float64 {
	return sumSeries(k)
}

func (k CandlestickSeriesList) getSeries(index int) series {
	return &k[index]
}

func (k CandlestickSeriesList) getSeriesName(index int) string {
	return k[index].Name
}

func (k CandlestickSeriesList) getSeriesValues(index int) []float64 {
	return k[index].getValues()
}

func (k CandlestickSeriesList) getSeriesLen(index int) int {
	return len(k[index].Data)
}

func (k CandlestickSeriesList) getSeriesSymbol(_ int) Symbol {
	return "" // no need to set symbol here, configured globally in candlestick_chart.go before defaultRender
}

func (k CandlestickSeriesList) hasMarkPoint() bool {
	for _, s := range k {
		if len(s.OpenMarkPoint.Points) > 0 || len(s.HighMarkPoint.Points) > 0 || len(s.LowMarkPoint.Points) > 0 || len(s.CloseMarkPoint.Points) > 0 {
			return true
		}
	}
	return false
}

func (k CandlestickSeriesList) setSeriesName(index int, name string) {
	k[index].Name = name
}

func (k CandlestickSeriesList) sortByNameIndex(dict map[string]int) {
	sort.Slice(k, func(i, j int) bool {
		return dict[k[i].Name] < dict[k[j].Name]
	})
}

// SetSeriesLabels sets the label for all elements in the series.
func (k CandlestickSeriesList) SetSeriesLabels(label SeriesLabel) {
	for i := range k {
		k[i].Label = label
	}
}

// ToGenericSeriesList converts candlestick series to generic series format.
// Each candlestick is encoded as 4 consecutive float64 values: [Open, High, Low, Close].
// Invalid OHLC data is encoded with null values to maintain the 4-element structure.
func (k CandlestickSeriesList) ToGenericSeriesList() GenericSeriesList {
	result := make([]GenericSeries, len(k))
	for i, s := range k {
		// Encode OHLC data as four in-order float64 values per candlestick
		values := make([]float64, 0, len(s.Data)*4)
		for _, ohlc := range s.Data {
			if validateOHLCData(ohlc) {
				values = append(values, ohlc.Open, ohlc.High, ohlc.Low, ohlc.Close)
			} else if validateOHLCHighLow(ohlc) {
				values = append(values, GetNullValue(), ohlc.High, ohlc.Low, GetNullValue())
			} else { // For invalid OHLC data, use null values to maintain structure
				values = append(values, GetNullValue(), GetNullValue(), GetNullValue(), GetNullValue())
			}
		}
		result[i] = GenericSeries{
			Values:     values,
			YAxisIndex: s.YAxisIndex,
			Label:      s.Label,
			Name:       s.Name,
			Type:       ChartTypeCandlestick,
			// For generic representation, use close values as primary
			MarkLine:  s.CloseMarkLine,
			MarkPoint: s.CloseMarkPoint,
		}
	}
	return result
}

// CandlestickSeriesOption configures optional elements when building
// candlestick series.
type CandlestickSeriesOption struct {
	// Label styles the series labels.
	Label SeriesLabel
	// Names provide data names for each series.
	Names []string
	// OpenMarkPoint marks open prices.
	OpenMarkPoint SeriesMarkPoint
	// OpenMarkLine draws reference lines for open prices.
	OpenMarkLine SeriesMarkLine
	// OpenTrendLine adds trend lines based on opens.
	OpenTrendLine []SeriesTrendLine
	// HighMarkPoint marks high prices.
	HighMarkPoint SeriesMarkPoint
	// HighMarkLine draws reference lines for highs.
	HighMarkLine SeriesMarkLine
	// HighTrendLine adds trend lines based on highs.
	HighTrendLine []SeriesTrendLine
	// LowMarkPoint marks low prices.
	LowMarkPoint SeriesMarkPoint
	// LowMarkLine draws reference lines for lows.
	LowMarkLine SeriesMarkLine
	// LowTrendLine adds trend lines based on lows.
	LowTrendLine []SeriesTrendLine
	// CloseMarkPoint marks closing prices.
	CloseMarkPoint SeriesMarkPoint
	// CloseMarkLine draws reference lines for closes.
	CloseMarkLine SeriesMarkLine
	// CloseTrendLine adds trend lines based on closes.
	CloseTrendLine []SeriesTrendLine
	// CandleStyle sets the drawing style for candles.
	CandleStyle string
	// PatternConfig configures candlestick pattern detection.
	PatternConfig *CandlestickPatternConfig
}

// NewSeriesListCandlestick builds a SeriesList for candlestick charts from OHLC data.
func NewSeriesListCandlestick(data [][]OHLCData, opts ...CandlestickSeriesOption) CandlestickSeriesList {
	var opt CandlestickSeriesOption
	if len(opts) != 0 {
		opt = opts[0]
	}

	seriesList := make([]CandlestickSeries, len(data))
	for index, ohlcData := range data {
		s := CandlestickSeries{
			Data:           ohlcData,
			Label:          opt.Label,
			OpenMarkPoint:  opt.OpenMarkPoint,
			OpenMarkLine:   opt.OpenMarkLine,
			OpenTrendLine:  opt.OpenTrendLine,
			HighMarkPoint:  opt.HighMarkPoint,
			HighMarkLine:   opt.HighMarkLine,
			HighTrendLine:  opt.HighTrendLine,
			LowMarkPoint:   opt.LowMarkPoint,
			LowMarkLine:    opt.LowMarkLine,
			LowTrendLine:   opt.LowTrendLine,
			CloseMarkPoint: opt.CloseMarkPoint,
			CloseMarkLine:  opt.CloseMarkLine,
			CloseTrendLine: opt.CloseTrendLine,
			CandleStyle:    opt.CandleStyle,
			PatternConfig:  opt.PatternConfig,
		}
		if index < len(opt.Names) {
			s.Name = opt.Names[index]
		}
		seriesList[index] = s
	}
	return seriesList
}

// ExtractOpenPrices extracts open prices from OHLC data.
func (k *CandlestickSeries) ExtractOpenPrices() []float64 {
	result := make([]float64, len(k.Data))
	for i, ohlc := range k.Data {
		if validateOHLCOpen(ohlc) {
			result[i] = ohlc.Open
		} else {
			result[i] = GetNullValue()
		}
	}
	return result
}

// ExtractClosePrices extracts close prices from OHLC data for use with indicators.
func (k *CandlestickSeries) ExtractClosePrices() []float64 {
	result := make([]float64, len(k.Data))
	for i, ohlc := range k.Data {
		if validateOHLCClose(ohlc) {
			result[i] = ohlc.Close
		} else {
			result[i] = GetNullValue()
		}
	}
	return result
}

// ExtractHighPrices extracts high prices from OHLC data.
func (k *CandlestickSeries) ExtractHighPrices() []float64 {
	result := make([]float64, len(k.Data))
	for i, ohlc := range k.Data {
		if validateOHLCHighLow(ohlc) {
			result[i] = ohlc.High
		} else {
			result[i] = GetNullValue()
		}
	}
	return result
}

// ExtractLowPrices extracts low prices from OHLC data.
func (k *CandlestickSeries) ExtractLowPrices() []float64 {
	result := make([]float64, len(k.Data))
	for i, ohlc := range k.Data {
		if validateOHLCHighLow(ohlc) {
			result[i] = ohlc.Low
		} else {
			result[i] = GetNullValue()
		}
	}
	return result
}

// AggregateCandlestick aggregates OHLC data by the specified factor.
func AggregateCandlestick(data CandlestickSeries, factor int) CandlestickSeries {
	if factor <= 1 {
		return data
	}

	aggregated := make([]OHLCData, 0, len(data.Data)/factor)
	for i := 0; i < len(data.Data); i += factor {
		end := i + factor
		if end > len(data.Data) {
			end = len(data.Data)
		}

		// Aggregate OHLC for this period
		open := data.Data[i].Open       // First open
		close := data.Data[end-1].Close // Last close
		high := data.Data[i].High       // Find max high
		low := data.Data[i].Low         // Find min low

		for j := i; j < end; j++ {
			if data.Data[j].High > high {
				high = data.Data[j].High
			}
			if data.Data[j].Low < low {
				low = data.Data[j].Low
			}
		}

		aggregated = append(aggregated, OHLCData{
			Open:  open,
			High:  high,
			Low:   low,
			Close: close,
		})
	}

	return CandlestickSeries{
		Data:           aggregated,
		YAxisIndex:     data.YAxisIndex,
		Label:          data.Label,
		Name:           data.Name,
		CloseMarkPoint: data.CloseMarkPoint,
		CloseMarkLine:  data.CloseMarkLine,
		CandleStyle:    data.CandleStyle,
	}
}

type populationSummary struct {
	// Max is the maximum value in the series.
	Max float64
	// MaxFirstIndex is the first index of the maximum value in the series. If the series is empty this value will be -1.
	MaxFirstIndex int
	// MaxIndex is the index of the maximum value in the series. If the series is empty this value will be -1.
	MaxIndex int
	// Min is the minimum value in the series.
	Min float64
	// MinFirstIndex is the first index of the minimum value in the series. If the series is empty this value will be -1.
	MinFirstIndex int
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
	var minFirstIndex, minIndex, maxFirstIndex, maxIndex int
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
			minFirstIndex = i
			minIndex = i
		} else if i > 0 && x == minValue && data[i-1] != minValue {
			minIndex = i // update the index to be the first point we returned to a minimum value
		}
		if x > maxValue {
			maxValue = x
			maxFirstIndex = i
			maxIndex = i
		} else if i > 0 && x == maxValue && data[i-1] != maxValue {
			maxIndex = i // update the index to be the first point we returned to a maximum value
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
		MaxFirstIndex:     maxFirstIndex,
		MaxIndex:          maxIndex,
		Min:               minValue,
		MinFirstIndex:     minFirstIndex,
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

func sumSeries(sl seriesList) []float64 {
	seriesLen := sl.len()
	sumValues := make([]float64, seriesLen)
	for i := 0; i < seriesLen; i++ {
		var total float64
		for _, v := range sl.getSeriesValues(i) {
			if v != GetNullValue() {
				total += v
			}
		}
		sumValues[i] = total
	}
	return sumValues
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

func getSeriesMaxDataCount(sl seriesList) (result int) {
	for i := 0; i < sl.len(); i++ {
		count := sl.getSeriesLen(i)
		if count > result {
			result = count
		}
	}
	return
}
