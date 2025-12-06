package charts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-analyze/bulk"
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
	if len(data) == 0 {
		return nil
	}
	if regexp.MustCompile(`^\d+`).Match(data) {
		data = []byte(fmt.Sprintf(`"%s"`, string(data)))
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
	return json.Unmarshal(data, &value.values)
}

// First returns the first value or 0 when empty.
func (value *EChartsSeriesDataValue) First() float64 {
	if len(value.values) == 0 {
		return 0
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

var numericRep = regexp.MustCompile(`^[-+]?[0-9]+(?:\.[0-9]+)?$`)

// UnmarshalJSON parses a series data item that may be a number or object.
func (es *EChartsSeriesData) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return nil
	}
	if numericRep.Match(data) {
		v, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return err
		}
		es.Value = EChartsSeriesDataValue{
			values: []float64{
				v,
			},
		}
		return nil
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
	BoundaryGap *bool            `json:"boundaryGap,omitempty"`
	SplitNumber int              `json:"splitNumber,omitempty"`
	AxisLabel   EChartsAxisLabel `json:"axisLabel,omitempty"`
	AxisLine    EChartsAxisLine  `json:"axisLine,omitempty"`
	Data        []string         `json:"data"`
	Type        string           `json:"type"`
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
	Formatter string `json:"formatter,omitempty"`
	Show      *bool  `json:"show,omitempty"`
	Color     string `json:"color,omitempty"`
	FontSize  *int   `json:"fontSize,omitempty"`
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

// EChartsYAxisData holds a single y-axis configuration block.
type EChartsYAxisData struct {
	Min       *float64         `json:"min,omitempty"`
	Max       *float64         `json:"max,omitempty"`
	AxisLabel EChartsAxisLabel `json:"axisLabel,omitempty"`
	AxisLine  EChartsAxisLine  `json:"axisLine,omitempty"`
	Data      []string         `json:"data"`
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
	arr := make([]int, 0)
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
	Show     bool   `json:"show"`
	Distance int    `json:"distance"`
	Color    string `json:"color"`
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
	ds := make([]*_EChartsMarkData, 0)
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
		Points: bulk.SliceTransform(func(i EChartsMarkData) SeriesMark {
			return SeriesMark{Type: i.Type}
		}, emp.Data),
	}
}

// EChartsMarkLine defines mark lines for a series.
type EChartsMarkLine struct {
	Data []EChartsMarkData `json:"data"`
}

// ToSeriesMarkLine converts the mark line to the internal representation.
func (eml *EChartsMarkLine) ToSeriesMarkLine() SeriesMarkLine {
	return SeriesMarkLine{
		Lines: bulk.SliceTransform(func(i EChartsMarkData) SeriesMark {
			return SeriesMark{Type: i.Type}
		}, eml.Data),
	}
}

// EChartsSeries holds data and styling for one chart series.
type EChartsSeries struct {
	Data       []EChartsSeriesData `json:"data"`
	Name       string              `json:"name"`
	Type       string              `json:"type"`
	Radius     string              `json:"radius"`
	YAxisIndex int                 `json:"yAxisIndex"`
	ItemStyle  EChartStyle         `json:"itemStyle,omitempty"` // TODO - add support
	// label configuration
	Label     EChartsLabelOption `json:"label"`
	MarkPoint EChartsMarkPoint   `json:"markPoint"`
	MarkLine  EChartsMarkLine    `json:"markLine"`
	Max       *float64           `json:"max"` // TODO - add support
	Min       *float64           `json:"min"` // TODO - add support
}

// EChartsSeriesList is a list of EChartsSeries values.
type EChartsSeriesList []EChartsSeries

func (esList EChartsSeriesList) ToSeriesList() GenericSeriesList {
	seriesList := make([]GenericSeries, 0, len(esList))
	for _, item := range esList {
		// if pie, each sub-recommendation generates a series
		if item.Type == ChartTypePie {
			for _, dataItem := range item.Data {
				seriesList = append(seriesList, GenericSeries{
					Type: item.Type,
					Name: dataItem.Name,
					Label: SeriesLabel{
						Show: Ptr(true),
					},
					Radius: item.Radius,
					Values: []float64{dataItem.Value.First()},
				})
			}
			continue
		}
		if item.Type == ChartTypeRadar ||
			item.Type == ChartTypeFunnel {
			for _, dataItem := range item.Data {
				seriesList = append(seriesList, GenericSeries{
					Name:   dataItem.Name,
					Type:   item.Type,
					Values: dataItem.Value.values,
					Label: SeriesLabel{
						FontStyle: FontStyle{
							FontColor: ParseColor(item.Label.Color),
						},
						Show:     Ptr(item.Label.Show),
						Distance: item.Label.Distance,
					},
				})
			}
			continue
		}
		seriesList = append(seriesList, GenericSeries{
			Type: item.Type,
			Values: bulk.SliceTransform(func(dataItem EChartsSeriesData) float64 {
				return dataItem.Value.First()
			}, item.Data),
			YAxisIndex: item.YAxisIndex,
			Label: SeriesLabel{
				FontStyle: FontStyle{
					FontColor: ParseColor(item.Label.Color),
				},
				Show:     Ptr(item.Label.Show),
				Distance: item.Label.Distance,
			},
			Name:      item.Name,
			MarkPoint: item.MarkPoint.ToSeriesMarkPoint(),
			MarkLine:  item.MarkLine.ToSeriesMarkLine(),
		})
	}
	return seriesList
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

// ToOption converts the ECharts options into a ChartOption.
func (eo *EChartsOption) ToOption() ChartOption {
	fontFamily := eo.FontFamily
	if len(fontFamily) == 0 {
		fontFamily = eo.Title.TextStyle.FontFamily
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
	o := ChartOption{
		OutputFormat: eo.Type,
		Font:         GetFont(fontFamily),
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
		RadarIndicators: bulk.SliceTransform(EChartsRadarIndicator.ToRadarIndicator, eo.Radar.Indicator),
		Width:           eo.Width,
		Height:          eo.Height,
		Padding:         eo.Padding.Box,
		Box:             eo.Box.ToBox(),
		SeriesList:      eo.Series.ToSeriesList(),
	}
	isHorizontalChart := false
	for _, item := range eo.XAxis.Data {
		if item.Type == "value" {
			isHorizontalChart = true
			break
		}
	}
	if isHorizontalChart {
		for index := range o.SeriesList {
			series := o.SeriesList[index]
			if series.Type == ChartTypeBar {
				o.SeriesList[index].Type = ChartTypeHorizontalBar
			}
		}
	}

	if len(eo.XAxis.Data) != 0 {
		xAxisData := eo.XAxis.Data[0]
		axisTheme := o.Theme
		axisLineColor := ParseColor(xAxisData.AxisLine.LineStyle.Color)
		if !axisLineColor.IsZero() {
			if xAxisData.AxisLine.LineStyle.Opacity != nil {
				axisLineColor = axisLineColor.WithAlpha(uint8(255 * *xAxisData.AxisLine.LineStyle.Opacity))
			}
			axisTheme = o.Theme.WithXAxisColor(axisLineColor)
		}
		o.XAxis = XAxisOption{
			Theme:       axisTheme,
			BoundaryGap: xAxisData.BoundaryGap,
			Labels:      xAxisData.Data,
			LabelCount:  xAxisData.SplitNumber,
			FontStyle:   xAxisData.AxisLabel.makeFontStyle(),
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
				axisLineColor = axisLineColor.WithAlpha(uint8(255 * *item.AxisLine.LineStyle.Opacity))
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
		yAxisOptions[index] = YAxisOption{
			Min:            item.Min,
			Max:            item.Max,
			ValueFormatter: valFormatter,
			Theme:          axisTheme,
			Labels:         item.Data,
			LabelFontStyle: item.AxisLabel.makeFontStyle(),
			SpineLineShow:  item.AxisLine.Show,
		}
	}
	o.YAxis = yAxisOptions
	o.Children = bulk.SliceTransform(func(child EChartsOption) ChartOption {
		return child.ToOption()
	}, eo.Children)
	return o
}

func renderEcharts(options, outputType string) ([]byte, error) {
	o := EChartsOption{}
	if err := json.Unmarshal([]byte(options), &o); err != nil {
		return nil, err
	}
	opt := o.ToOption()
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
