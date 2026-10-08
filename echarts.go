package charts

import (
	"bytes"
	"encoding/json"
	"errors"
	"slices"
	"strconv"
	"strings"

	"github.com/golang/freetype/truetype"

	"github.com/go-analyze/charts/chartdraw/drawing"
)

func convertToArray(data []byte) []byte {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return nil
	}
	if data[0] != '[' {
		data = []byte("[" + string(data) + "]")
	}
	return data
}

// EChartsPosition represents a CSS-like position value that can be either a string (like "center", "left") or a numeric value.
type EChartsPosition string

// UnmarshalJSON decodes a position JSON value that may be a string or number.
func (p *EChartsPosition) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return nil
	}
	if c := data[0]; c == '-' || c == '+' || c == '.' || (c >= '0' && c <= '9') {
		data = append(append([]byte{'"'}, data...), '"')
	}
	s := (*string)(p)
	return json.Unmarshal(data, s)
}

// EChartStyle describes color and opacity for ECharts elements.
type EChartStyle struct {
	Color   string   `json:"color"`
	Opacity *float64 `json:"opacity,omitempty"`
}

// EChartsSeriesDataValue holds numeric values from an ECharts data entry.
type EChartsSeriesDataValue struct {
	values []float64
}

// UnmarshalJSON decodes a series data value that may be a single number or array.
func (value *EChartsSeriesDataValue) UnmarshalJSON(data []byte) error {
	data = convertToArray(data)
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	values := make([]float64, len(raw))
	for i, r := range raw {
		v, err := parseSeriesValue(r)
		if err != nil {
			return err
		}
		values[i] = v
	}
	value.values = values
	return nil
}

// parseSeriesValue returns the number for a JSON value, with null or "-" returning GetNullValue.
func parseSeriesValue(data []byte) (float64, error) {
	switch string(bytes.TrimSpace(data)) {
	case "", "null", `"-"`:
		return GetNullValue(), nil
	}
	var f float64
	err := json.Unmarshal(data, &f)
	return f, err
}

// First returns the first value, or GetNullValue when empty.
func (value *EChartsSeriesDataValue) First() float64 {
	if len(value.values) == 0 {
		return GetNullValue()
	}
	return value.values[0]
}

// EChartsSeriesData describes a single data item from ECharts.
type EChartsSeriesData struct {
	Value     EChartsSeriesDataValue `json:"value"`
	Name      string                 `json:"name"`
	ItemStyle EChartStyle            `json:"itemStyle,omitempty"` // TODO - add support
}
type _EChartsSeriesData EChartsSeriesData

// UnmarshalJSON parses a series data item that may be a number or object.
func (es *EChartsSeriesData) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return nil
	}
	if data[0] != '{' { // scalar, null or array value
		return es.Value.UnmarshalJSON(data)
	}
	v := _EChartsSeriesData{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	es.Name = v.Name
	es.Value = v.Value
	es.ItemStyle = v.ItemStyle
	return nil
}

// EChartsXAxisData holds x-axis configuration extracted from ECharts JSON.
type EChartsXAxisData struct {
	BoundaryGap   *bool                `json:"boundaryGap,omitempty"`
	SplitNumber   int                  `json:"splitNumber,omitempty"`
	Name          string               `json:"name,omitempty"`
	NameTextStyle EChartsNameTextStyle `json:"nameTextStyle,omitempty"`
	AxisLabel     EChartsAxisLabel     `json:"axisLabel,omitempty"`
	AxisLine      EChartsAxisLine      `json:"axisLine,omitempty"`
	Data          []string             `json:"data"`
	Type          string               `json:"type"`
}

// EChartsAxisLine describes the line styling for an axis.
type EChartsAxisLine struct {
	Show      *bool `json:"show,omitempty"`
	LineStyle struct {
		Color   string   `json:"color,omitempty"`
		Opacity *float64 `json:"opacity,omitempty"`
		Width   *int     `json:"width,omitempty"` // TODO - add support
	} `json:"lineStyle,omitempty"`
}

// EChartsXAxis holds a list of x-axis options.
type EChartsXAxis struct {
	Data []EChartsXAxisData
}

// UnmarshalJSON decodes x-axis options that may be a single object or an array.
func (ex *EChartsXAxis) UnmarshalJSON(data []byte) error {
	data = convertToArray(data)
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, &ex.Data)
}

// EChartsAxisLabel configures axis label display for ECharts.
type EChartsAxisLabel struct {
	Formatter string   `json:"formatter,omitempty"`
	Show      *bool    `json:"show,omitempty"`
	Color     string   `json:"color,omitempty"`
	FontSize  *int     `json:"fontSize,omitempty"`
	Rotate    *float64 `json:"rotate,omitempty"`
	Interval  *int     `json:"interval,omitempty"`
	Margin    *int     `json:"margin,omitempty"`
}

func (al EChartsAxisLabel) makeFontStyle() FontStyle {
	var axisFont FontStyle
	if al.FontSize != nil {
		axisFont.FontSize = float64(*al.FontSize)
	}
	if flagIs(false, al.Show) {
		axisFont.FontColor = ColorTransparent
	} else if axisTextColor := ParseColor(al.Color); !axisTextColor.IsZero() {
		axisFont.FontColor = axisTextColor
	}
	return axisFont
}

// EChartsNameTextStyle styles an axis name (title).
type EChartsNameTextStyle struct {
	Color    string `json:"color,omitempty"`
	FontSize *int   `json:"fontSize,omitempty"`
}

// makeFontStyle builds the axis title style, unset fields fall back to axis defaults.
func (nt EChartsNameTextStyle) makeFontStyle() FontStyle {
	var titleFont FontStyle
	if nt.FontSize != nil {
		titleFont.FontSize = float64(*nt.FontSize)
	}
	if titleColor := ParseColor(nt.Color); !titleColor.IsZero() {
		titleFont.FontColor = titleColor
	}
	return titleFont
}

// EChartsSplitLine controls axis split line visibility.
type EChartsSplitLine struct {
	Show *bool `json:"show,omitempty"`
}

// EChartsAreaStyle describes the area fill styling for line series.
type EChartsAreaStyle struct {
	Opacity *float64 `json:"opacity,omitempty"`
	Color   string   `json:"color,omitempty"` // TODO - add support
}

// echartsSmoothDefaultTension selects the line tension when smooth is set to true.
const echartsSmoothDefaultTension = 0.5

// EChartsSmooth holds line smoothing as a boolean or a 0-1 tension value.
type EChartsSmooth struct {
	value float64
	set   bool
}

// UnmarshalJSON decodes smoothing provided as a boolean or numeric tension.
func (s *EChartsSmooth) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	switch string(data) {
	case "", "null":
		return nil
	case "true":
		s.value = echartsSmoothDefaultTension
		s.set = true
		return nil
	case "false":
		s.set = true
		return nil
	}
	if err := json.Unmarshal(data, &s.value); err != nil {
		return err
	}
	s.set = true
	return nil
}

// tension returns the smoothing tension clamped to 0-1, 0 when unset or disabled.
func (s EChartsSmooth) tension() float64 {
	if !s.set || s.value <= 0 {
		return 0
	}
	return min(s.value, 1)
}

// EChartsRadius holds a circular chart radius: a single value or an [inner, outer] pair.
type EChartsRadius struct {
	// Inner is the center-hole radius used by doughnuts.
	Inner string
	// Outer is the ring (or pie) radius.
	Outer string
	// IsSet is true when a radius was present in the source JSON.
	IsSet bool
}

// UnmarshalJSON decodes a radius provided as a number, percent string, or [inner, outer] pair.
func (r *EChartsRadius) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	if data[0] == '[' {
		var pair []EChartsPosition
		if err := json.Unmarshal(data, &pair); err != nil {
			return err
		}
		if len(pair) > 0 {
			r.Inner = string(pair[0])
		}
		if len(pair) > 1 {
			r.Outer = string(pair[1])
		}
		r.IsSet = true
		return nil
	}
	var single EChartsPosition
	if err := json.Unmarshal(data, &single); err != nil {
		return err
	}
	r.Outer = string(single)
	r.IsSet = true
	return nil
}

// EChartsYAxisData holds a single y-axis configuration block.
type EChartsYAxisData struct {
	Min           *float64             `json:"min,omitempty"`
	Max           *float64             `json:"max,omitempty"`
	SplitNumber   int                  `json:"splitNumber,omitempty"`
	Name          string               `json:"name,omitempty"`
	NameTextStyle EChartsNameTextStyle `json:"nameTextStyle,omitempty"`
	Position      string               `json:"position,omitempty"`
	SplitLine     EChartsSplitLine     `json:"splitLine,omitempty"`
	AxisLabel     EChartsAxisLabel     `json:"axisLabel,omitempty"`
	AxisLine      EChartsAxisLine      `json:"axisLine,omitempty"`
	Data          []string             `json:"data"`
}

// EChartsYAxis represents a list of y-axis definitions.
type EChartsYAxis struct {
	Data []EChartsYAxisData `json:"data"`
}

// UnmarshalJSON decodes y-axis options that may be a single object or an array.
func (ey *EChartsYAxis) UnmarshalJSON(data []byte) error {
	data = convertToArray(data)
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, &ey.Data)
}

// EChartsPadding represents padding values around a component.
type EChartsPadding struct {
	Box Box
}

// UnmarshalJSON decodes a padding array into a Box.
func (eb *EChartsPadding) UnmarshalJSON(data []byte) error {
	data = convertToArray(data)
	if len(data) == 0 {
		return nil
	}
	var arr []int
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	if len(arr) == 0 {
		return nil
	}
	switch len(arr) {
	case 1:
		eb.Box = NewBoxEqual(arr[0])
	case 2:
		eb.Box = NewBox(arr[1], arr[0], arr[1], arr[0])
	default:
		result := make([]int, 4)
		copy(result, arr)
		if len(arr) == 3 {
			result[3] = result[1]
		}
		// top, right, bottom, left
		eb.Box = NewBox(result[3], result[0], result[1], result[2])
	}
	return nil
}

// EChartsBox represents box dimensions with JSON tags for ECharts parsing.
type EChartsBox struct {
	Top    int  `json:"top"`
	Bottom int  `json:"bottom"`
	Left   int  `json:"left"`
	Right  int  `json:"right"`
	IsSet  bool `json:"isSet"`
}

// ToBox converts EChartsBox to Box.
func (eb EChartsBox) ToBox() Box {
	return Box{
		Top:    eb.Top,
		Bottom: eb.Bottom,
		Left:   eb.Left,
		Right:  eb.Right,
		IsSet:  eb.IsSet,
	}
}

// EChartsLabelOption configures data labels.
type EChartsLabelOption struct {
	Show      bool                  `json:"show"`
	Distance  int                   `json:"distance"`
	Color     string                `json:"color"`
	Formatter EChartsLabelFormatter `json:"formatter,omitempty"`
	Position  string                `json:"position,omitempty"`
}

const defaultLabelDistance = 5

// EChartsLabelFormatter holds a label formatter template string. Functions and
// rich-text objects are not expressible in JSON and are silently ignored.
type EChartsLabelFormatter struct {
	Template string
}

// UnmarshalJSON decodes a formatter when provided as a string, ignoring other forms.
func (f *EChartsLabelFormatter) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || data[0] != '"' {
		return nil
	}
	return json.Unmarshal(data, &f.Template)
}

// makeLabel builds the native label configuration, resolving the template
// formatter and approximating the position through a label offset.
func (el EChartsLabelOption) makeLabel(seriesName string, values []float64) SeriesLabel {
	label := SeriesLabel{
		Show: Ptr(el.Show),
		FontStyle: FontStyle{
			FontColor: ParseColor(el.Color),
		},
		Distance: el.Distance,
	}
	switch el.Position {
	case "bottom", "inside", "insideBottom":
		// the native label paints above the anchor by Distance, shift fully below
		label.Offset.Top = 2 * max(el.Distance, defaultLabelDistance)
	}
	if el.Formatter.Template != "" {
		label.LabelFormatter = makeTemplateFormatter(el.Formatter.Template, seriesName, values)
	}
	return label
}

// makeTemplateFormatter resolves ECharts label template tokens:
// {a} series name, {b} data name, {c} value, {d} percent of the series total.
// The library retains a single name per series, so {b} resolves to that name.
func makeTemplateFormatter(template, seriesName string, values []float64) SeriesLabelFormatter {
	var total float64
	for _, v := range values {
		if isValidExtent(v) {
			total += v
		}
	}
	return func(index int, name string, val float64) (string, *LabelStyle) {
		if !isValidExtent(val) {
			return "", nil
		}
		if name == "" {
			name = seriesName
		}
		label := strings.NewReplacer(
			"{a}", name,
			"{b}", name,
			"{c}", FormatValueHumanize(val, 2, false),
			"{d}", FormatValueHumanize(valuePercent(val, total), 2, false),
		).Replace(template)
		return strings.TrimSpace(label), nil
	}
}

// valuePercent returns the value share of the total in percent, 0 without a usable total.
func valuePercent(val, total float64) float64 {
	if total <= 0 {
		return 0
	}
	return val / total * 100
}

// EChartsLegend holds legend configuration from ECharts JSON.
type EChartsLegend struct {
	Show            *bool            `json:"show"`
	Data            []string         `json:"data"`
	Align           string           `json:"align"`
	Orient          string           `json:"orient"`
	Padding         EChartsPadding   `json:"padding,omitempty"`
	Left            EChartsPosition  `json:"left"`
	Top             EChartsPosition  `json:"top"`
	TextStyle       EChartsTextStyle `json:"textStyle"`
	BackgroundColor string           `json:"backgroundColor,omitempty"` // TODO - add support
	BorderColor     string           `json:"borderColor,omitempty"`
}

// EChartsMarkData represents mark lines or points in ECharts JSON.
type EChartsMarkData struct {
	Type string `json:"type"`
	// TODO - support position values below
	XAxis float64 `json:"xAxis,omitempty"`
	YAxis float64 `json:"yAxis,omitempty"`
}
type _EChartsMarkData EChartsMarkData

// UnmarshalJSON parses mark definitions provided as an object or array.
func (emd *EChartsMarkData) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return nil
	}
	data = convertToArray(data)
	var ds []*_EChartsMarkData
	if err := json.Unmarshal(data, &ds); err != nil {
		return err
	}
	for _, d := range ds {
		if d.Type != "" {
			emd.Type = d.Type
		}
	}
	return nil
}

// EChartsMarkPoint defines mark points for a series.
type EChartsMarkPoint struct {
	SymbolSize int               `json:"symbolSize"`
	Data       []EChartsMarkData `json:"data"`
}

// ToSeriesMarkPoint converts the mark point to the internal representation.
func (emp *EChartsMarkPoint) ToSeriesMarkPoint() SeriesMarkPoint {
	return SeriesMarkPoint{
		SymbolSize: emp.SymbolSize,
		Points: mapSlice(emp.Data, func(i EChartsMarkData) SeriesMark {
			return SeriesMark{Type: i.Type}
		}),
	}
}

// EChartsMarkLine defines mark lines for a series.
type EChartsMarkLine struct {
	Data []EChartsMarkData `json:"data"`
}

// ToSeriesMarkLine converts the mark line to the internal representation.
func (eml *EChartsMarkLine) ToSeriesMarkLine() SeriesMarkLine {
	return SeriesMarkLine{
		Lines: mapSlice(eml.Data, func(i EChartsMarkData) SeriesMark {
			return SeriesMark{Type: i.Type}
		}),
	}
}

// EChartsSeries holds data and styling for one chart series.
type EChartsSeries struct {
	Data       []EChartsSeriesData `json:"data"`
	Name       string              `json:"name"`
	Type       string              `json:"type"`
	Radius     EChartsRadius       `json:"radius"`
	Stack      string              `json:"stack,omitempty"`
	Symbol     string              `json:"symbol,omitempty"`
	SymbolSize *float64            `json:"symbolSize,omitempty"`
	LineStyle  struct {
		Width *float64 `json:"width,omitempty"`
	} `json:"lineStyle,omitempty"`
	AreaStyle      *EChartsAreaStyle `json:"areaStyle,omitempty"`
	BarWidth       string            `json:"barWidth,omitempty"`
	BarGap         string            `json:"barGap,omitempty"`
	BarCategoryGap string            `json:"barCategoryGap,omitempty"` // ignored, no library counterpart
	Smooth         EChartsSmooth     `json:"smooth,omitempty"`
	YAxisIndex     int               `json:"yAxisIndex"`
	ItemStyle      EChartStyle       `json:"itemStyle,omitempty"` // TODO - add support
	// label configuration
	Label     EChartsLabelOption `json:"label"`
	MarkPoint EChartsMarkPoint   `json:"markPoint"`
	MarkLine  EChartsMarkLine    `json:"markLine"`
	Max       *float64           `json:"max"` // TODO - add support
	Min       *float64           `json:"min"` // TODO - add support
}

// EChartsSeriesList is a list of EChartsSeries values.
type EChartsSeriesList []EChartsSeries

// ToSeriesList converts the ECharts series into native series, erroring when the
// configuration cannot be represented (notably mixed stacked and unstacked series).
func (esList EChartsSeriesList) ToSeriesList() (GenericSeriesList, error) {
	seriesList := make([]GenericSeries, 0, len(esList))
	var stackGroups []string
	var stackCount, unstackedCount int
	for _, item := range esList {
		if item.Stack != "" {
			if !slices.Contains(stackGroups, item.Stack) {
				stackGroups = append(stackGroups, item.Stack)
			}
			stackCount++
		} else if isStackableType(item.Type) {
			unstackedCount++
		}
		switch item.Type {
		case ChartTypePie, ChartTypeDoughnut:
			// each data item becomes its own series
			values := mapSlice(item.Data, func(dataItem EChartsSeriesData) float64 {
				return dataItem.Value.First()
			})
			for _, dataItem := range item.Data {
				label := item.Label.makeLabel(dataItem.Name, values)
				label.Show = Ptr(true) // circular chart labels are forced on
				seriesList = append(seriesList, GenericSeries{
					Type:   item.Type,
					Name:   dataItem.Name,
					Label:  label,
					Radius: item.Radius.Outer,
					Values: []float64{dataItem.Value.First()},
				})
			}
		case ChartTypeRadar, ChartTypeFunnel:
			// each data item becomes its own series
			for _, dataItem := range item.Data {
				seriesList = append(seriesList, GenericSeries{
					Name:   dataItem.Name,
					Type:   item.Type,
					Values: dataItem.Value.values,
					Label:  item.Label.makeLabel(item.Name, nil),
				})
			}
		default:
			values := mapSlice(item.Data, func(dataItem EChartsSeriesData) float64 {
				return dataItem.Value.First()
			})
			seriesList = append(seriesList, GenericSeries{
				Type:       item.Type,
				Values:     values,
				YAxisIndex: item.YAxisIndex,
				Label:      item.Label.makeLabel(item.Name, values),
				Name:       item.Name,
				MarkPoint:  item.MarkPoint.ToSeriesMarkPoint(),
				MarkLine:   item.MarkLine.ToSeriesMarkLine(),
			})
		}
	}
	if stackCount > 0 && unstackedCount > 0 {
		return nil, errors.New("stack must be defined on all series when any series uses a stack group")
	}
	if len(stackGroups) > 1 {
		return nil, errors.New("only a single stack group is supported, found: " + strings.Join(stackGroups, ", "))
	}
	return seriesList, nil
}

// isStackableType reports whether the chart type participates in stack series rendering.
func isStackableType(chartType string) bool {
	switch chartType {
	case ChartTypeLine, ChartTypeBar, ChartTypeHorizontalBar, ChartTypeScatter:
		return true
	}
	return false
}

// EChartsRadarIndicator maps radar indicator options from ECharts.
type EChartsRadarIndicator struct {
	Name string  `json:"name"`
	Max  float64 `json:"max"`
	Min  float64 `json:"min"`
}

// ToRadarIndicator converts to a RadarIndicator.
func (eri EChartsRadarIndicator) ToRadarIndicator() RadarIndicator {
	return RadarIndicator{
		Name: eri.Name,
		Max:  eri.Max,
		Min:  eri.Min,
	}
}

// EChartsTextStyle maps text style options from ECharts.
type EChartsTextStyle struct {
	Color      string  `json:"color"`
	FontFamily string  `json:"fontFamily"`
	FontSize   float64 `json:"fontSize"`
}

// ToFontStyle converts the text style to a FontStyle.
func (et *EChartsTextStyle) ToFontStyle() FontStyle {
	s := FontStyle{
		FontSize:  et.FontSize,
		FontColor: ParseColor(et.Color),
	}
	if et.FontFamily != "" {
		s.Font = GetFont(et.FontFamily)
	}
	return s
}

// EChartsOption mirrors a basic ECharts configuration.
type EChartsOption struct {
	Type       string         `json:"type"`
	Theme      string         `json:"theme"`
	FontFamily string         `json:"fontFamily"`
	Padding    EChartsPadding `json:"padding"`
	Box        EChartsBox     `json:"box"`
	Width      int            `json:"width"`
	Height     int            `json:"height"`
	Title      struct {
		Show            *bool            `json:"show,omitempty"`
		Text            string           `json:"text"`
		Subtext         string           `json:"subtext"`
		Left            EChartsPosition  `json:"left"`
		Top             EChartsPosition  `json:"top"`
		TextStyle       EChartsTextStyle `json:"textStyle"`
		SubtextStyle    EChartsTextStyle `json:"subtextStyle"`
		BackgroundColor string           `json:"backgroundColor,omitempty"` // TODO - add support
		BorderColor     string           `json:"borderColor,omitempty"`
	} `json:"title"`
	XAxis  EChartsXAxis  `json:"xAxis"`
	YAxis  EChartsYAxis  `json:"yAxis"`
	Legend EChartsLegend `json:"legend"`
	Radar  struct {
		Indicator []EChartsRadarIndicator `json:"indicator"`
	} `json:"radar"`
	Series          EChartsSeriesList `json:"series"`
	BackgroundColor string            `json:"backgroundColor,omitempty"`
	Children        []EChartsOption   `json:"children"`
}

// ToOption converts the ECharts options into a ChartOption, erroring when the
// configuration cannot be represented by the library.
func (eo *EChartsOption) ToOption() (ChartOption, error) {
	fontFamily := eo.FontFamily
	if len(fontFamily) == 0 {
		fontFamily = eo.Title.TextStyle.FontFamily
	}
	var fallbackFont *truetype.Font
	if fontFamily != "" {
		fallbackFont = GetFont(fontFamily)
	}
	theme := GetTheme(eo.Theme)
	backgroundColor := ParseColor(eo.BackgroundColor)
	if !backgroundColor.IsZero() {
		theme = theme.WithBackgroundColor(backgroundColor)
	}
	titleBorderColor := ParseColor(eo.Title.BorderColor)
	titleBorderWidth := 0.0
	if !titleBorderColor.IsZero() {
		theme = theme.WithTitleBorderColor(titleBorderColor)
		titleBorderWidth = defaultStrokeWidth
	}
	legendBorderColor := ParseColor(eo.Legend.BorderColor)
	legendBorderWidth := 0.0
	if !legendBorderColor.IsZero() {
		theme = theme.WithLegendBorderColor(legendBorderColor)
		legendBorderWidth = defaultStrokeWidth
	}
	titleTextStyle := eo.Title.TextStyle.ToFontStyle()
	titleSubtextStyle := eo.Title.SubtextStyle.ToFontStyle()
	legendTextStyle := eo.Legend.TextStyle.ToFontStyle()
	if fallbackFont != nil {
		if titleTextStyle.Font == nil {
			titleTextStyle.Font = fallbackFont
		}
		if titleSubtextStyle.Font == nil {
			titleSubtextStyle.Font = fallbackFont
		}
		if legendTextStyle.Font == nil {
			legendTextStyle.Font = fallbackFont
		}
	}
	seriesList, err := eo.Series.ToSeriesList()
	if err != nil {
		return ChartOption{}, err
	}
	o := ChartOption{
		OutputFormat: eo.Type,
		Theme:        theme,
		Title: TitleOption{
			Show:             eo.Title.Show,
			Text:             eo.Title.Text,
			Subtext:          eo.Title.Subtext,
			FontStyle:        titleTextStyle,
			SubtextFontStyle: titleSubtextStyle,
			Offset: OffsetStr{
				Left: string(eo.Title.Left),
				Top:  string(eo.Title.Top),
			},
			BorderWidth: titleBorderWidth,
		},
		Legend: LegendOption{
			Show:        eo.Legend.Show,
			FontStyle:   legendTextStyle,
			SeriesNames: eo.Legend.Data,
			Offset: OffsetStr{
				Left: string(eo.Legend.Left),
				Top:  string(eo.Legend.Top),
			},
			Align:       eo.Legend.Align,
			Vertical:    Ptr(strings.EqualFold(eo.Legend.Orient, "vertical")),
			Padding:     eo.Legend.Padding.Box,
			BorderWidth: legendBorderWidth,
		},
		RadarIndicators: mapSlice(eo.Radar.Indicator, EChartsRadarIndicator.ToRadarIndicator),
		Width:           eo.Width,
		Height:          eo.Height,
		Padding:         eo.Padding.Box,
		Box:             eo.Box.ToBox(),
		SeriesList:      seriesList,
	}
	eo.Series.applyChartFields(&o)
	if fallbackFont != nil {
		for i := range o.SeriesList {
			if o.SeriesList[i].Label.FontStyle.Font == nil {
				o.SeriesList[i].Label.FontStyle.Font = fallbackFont
			}
		}
	}
	isHorizontalChart := slices.ContainsFunc(eo.XAxis.Data, func(item EChartsXAxisData) bool {
		return item.Type == "value"
	})
	if isHorizontalChart {
		for index := range o.SeriesList {
			switch o.SeriesList[index].Type {
			case ChartTypeBar:
				o.SeriesList[index].Type = ChartTypeHorizontalBar
			case ChartTypeViolin:
				o.SeriesList[index].Type = ChartTypeHorizontalViolin
			}
		}
	}

	if len(eo.XAxis.Data) != 0 {
		xAxisData := eo.XAxis.Data[0]
		axisTheme := o.Theme
		axisLineColor := ParseColor(xAxisData.AxisLine.LineStyle.Color)
		if !axisLineColor.IsZero() {
			if xAxisData.AxisLine.LineStyle.Opacity != nil {
				axisLineColor = axisLineColor.WithAlpha(drawing.ColorChannelFromFloat(*xAxisData.AxisLine.LineStyle.Opacity))
			}
			axisTheme = o.Theme.WithXAxisColor(axisLineColor)
		}
		xLabelFontStyle := xAxisData.AxisLabel.makeFontStyle()
		if fallbackFont != nil && xLabelFontStyle.Font == nil {
			xLabelFontStyle.Font = fallbackFont
		}
		o.XAxis = XAxisOption{
			Theme:          axisTheme,
			BoundaryGap:    xAxisData.BoundaryGap,
			Labels:         xAxisData.Data,
			LabelCount:     xAxisData.SplitNumber,
			LabelFontStyle: xLabelFontStyle,
			Title:          xAxisData.Name,
			TitleFontStyle: xAxisData.NameTextStyle.makeFontStyle(),
		}
		if xAxisData.AxisLabel.Rotate != nil {
			o.XAxis.LabelRotation = DegreesToRadians(*xAxisData.AxisLabel.Rotate)
		}
		if xAxisData.AxisLabel.Margin != nil {
			o.XAxis.LabelOffset = OffsetInt{Top: *xAxisData.AxisLabel.Margin}
		}
		if o.XAxis.BoundaryGap == nil {
			// Ensure default ECharts behavior of centering labels and sets a "BoundaryGap"
			// https://echarts.apache.org/en/option.html#xAxis.boundaryGap
			o.XAxis.BoundaryGap = Ptr(true)
		}
	}
	yAxisOptions := make([]YAxisOption, len(eo.YAxis.Data))
	for index, item := range eo.YAxis.Data {
		axisTheme := o.Theme
		if axisLineColor := ParseColor(item.AxisLine.LineStyle.Color); !axisLineColor.IsZero() {
			if item.AxisLine.LineStyle.Opacity != nil {
				axisLineColor = axisLineColor.WithAlpha(drawing.ColorChannelFromFloat(*item.AxisLine.LineStyle.Opacity))
			}
			axisTheme = axisTheme.WithYAxisColor(axisLineColor).WithYAxisTextColor(axisLineColor)
		}
		var valFormatter ValueFormatter
		if item.AxisLabel.Formatter != "" {
			valFormatter = func(f float64) string {
				return strings.ReplaceAll(item.AxisLabel.Formatter, "{value}",
					FormatValueHumanize(f, 2, false))
			}
		}
		yLabelFontStyle := item.AxisLabel.makeFontStyle()
		if fallbackFont != nil && yLabelFontStyle.Font == nil {
			yLabelFontStyle.Font = fallbackFont
		}
		// ECharts shows split lines by default, the library hides them unless enabled
		splitLineShow := Ptr(true)
		if item.SplitLine.Show != nil {
			splitLineShow = item.SplitLine.Show
		}
		var axisPosition string
		if item.Position == PositionLeft || item.Position == PositionRight {
			axisPosition = item.Position
		}
		var labelOffset OffsetInt
		if item.AxisLabel.Margin != nil {
			// labels sit left of a left-positioned axis, right of a right-positioned one
			if axisPosition == PositionRight {
				labelOffset.Left = *item.AxisLabel.Margin
			} else {
				labelOffset.Left = -*item.AxisLabel.Margin
			}
		}
		var labelRotation float64
		if item.AxisLabel.Rotate != nil {
			labelRotation = DegreesToRadians(*item.AxisLabel.Rotate)
		}
		var labelSkipCount int
		if item.AxisLabel.Interval != nil {
			labelSkipCount = *item.AxisLabel.Interval
		}
		yAxisOptions[index] = YAxisOption{
			Min:            item.Min,
			Max:            item.Max,
			ValueFormatter: valFormatter,
			Theme:          axisTheme,
			Labels:         item.Data,
			LabelFontStyle: yLabelFontStyle,
			SpineLineShow:  item.AxisLine.Show,
			Title:          item.Name,
			TitleFontStyle: item.NameTextStyle.makeFontStyle(),
			Position:       axisPosition,
			LabelCount:     item.SplitNumber,
			LabelRotation:  labelRotation,
			LabelOffset:    labelOffset,
			LabelSkipCount: labelSkipCount,
			SplitLineShow:  splitLineShow,
		}
	}
	o.YAxis = yAxisOptions
	for _, child := range eo.Children {
		childOption, err := child.ToOption()
		if err != nil {
			return ChartOption{}, err
		}
		o.Children = append(o.Children, childOption)
	}
	return o, nil
}

// applyChartFields maps per-series fields onto the chart-level options;
// the first series defining each field wins.
func (esList EChartsSeriesList) applyChartFields(o *ChartOption) {
	var smoothSet bool
	for _, item := range esList {
		if symbolShape := makeSymbolShape(item.Symbol); symbolShape != "" && o.Symbol.Shape == "" {
			o.Symbol.Shape = symbolShape
		}
		if item.SymbolSize != nil && o.Symbol.Size == 0 {
			o.Symbol.Size = *item.SymbolSize
		}
		if item.LineStyle.Width != nil && *item.LineStyle.Width > 0 && o.LineStrokeWidth == 0 {
			o.LineStrokeWidth = *item.LineStyle.Width
		}
		if item.AreaStyle != nil && o.FillArea == nil {
			o.FillArea = Ptr(true)
			if item.AreaStyle.Opacity != nil {
				o.FillOpacity = drawing.ColorChannelFromFloat(*item.AreaStyle.Opacity)
			}
		}
		if ratio, ok := parsePercentRatio(item.BarWidth); ok && o.BarSize == 0 {
			o.BarSize = ratio
		}
		if ratio, ok := parsePercentRatio(item.BarGap); ok && o.BarMargin == nil {
			o.BarMargin = Ptr(ratio)
		}
		if item.Smooth.set && !smoothSet {
			o.StrokeSmoothingTension = item.Smooth.tension()
			smoothSet = true
		}
	}
	if slices.ContainsFunc(esList, func(item EChartsSeries) bool {
		return item.Stack != ""
	}) {
		o.StackSeries = Ptr(true) // ToSeriesList errors unless all series stack
	}
	for _, item := range esList {
		if item.Type == ChartTypeDoughnut && item.Radius.IsSet {
			o.Radius = item.Radius.Outer
			o.RadiusCenter = item.Radius.Inner
			break
		}
	}
}

// makeSymbolShape maps an ECharts symbol name to a library shape;
// unsupported shapes fall back to square.
func makeSymbolShape(symbol string) SymbolShape {
	switch symbol {
	case "":
		return ""
	case "circle":
		return SymbolCircle
	case "rect", "roundRect":
		return SymbolSquare
	case "diamond":
		return SymbolDiamond
	case "none":
		return SymbolNone
	default: // "triangle", "pin", "arrow" and unknown names
		return SymbolSquare
	}
}

// parsePercentRatio parses a percent string like "35%" into a 0-1 ratio.
// Non-percent values are rejected; absolute sizes have no slot-ratio equivalent.
func parsePercentRatio(value string) (float64, bool) {
	percentStr, ok := strings.CutSuffix(value, "%")
	if !ok {
		return 0, false
	}
	percent, err := strconv.ParseFloat(percentStr, 64)
	if err != nil {
		return 0, false
	}
	return percent / 100.0, true
}

func renderEcharts(options, outputType string) ([]byte, error) {
	o := EChartsOption{}
	if err := json.Unmarshal([]byte(options), &o); err != nil {
		return nil, err
	}
	opt, err := o.ToOption()
	if err != nil {
		return nil, err
	}
	opt.OutputFormat = outputType
	if p, err := Render(opt); err != nil {
		return nil, err
	} else {
		return p.Bytes()
	}
}

// RenderEChartsToPNG renders an ECharts option JSON string to PNG bytes.
func RenderEChartsToPNG(options string) ([]byte, error) {
	return renderEcharts(options, ChartOutputPNG)
}

// RenderEChartsToJPG renders an ECharts option JSON string to JPG bytes.
func RenderEChartsToJPG(options string) ([]byte, error) {
	return renderEcharts(options, ChartOutputJPG)
}

// RenderEChartsToSVG renders an ECharts option JSON string to SVG bytes.
func RenderEChartsToSVG(options string) ([]byte, error) {
	return renderEcharts(options, ChartOutputSVG)
}
