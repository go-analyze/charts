package charts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
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

type EChartsPosition string

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

type EChartStyle struct {
	Color   string   `json:"color"`
	Opacity *float64 `json:"opacity,omitempty"`
}

type EChartsSeriesDataValue struct {
	values []float64
}

func (value *EChartsSeriesDataValue) UnmarshalJSON(data []byte) error {
	data = convertToArray(data)
	return json.Unmarshal(data, &value.values)
}

func (value *EChartsSeriesDataValue) First() float64 {
	if len(value.values) == 0 {
		return 0
	}
	return value.values[0]
}

type EChartsSeriesData struct {
	Value     EChartsSeriesDataValue `json:"value"`
	Name      string                 `json:"name"`
	ItemStyle EChartStyle            `json:"itemStyle,omitempty"` // TODO - add support
}
type _EChartsSeriesData EChartsSeriesData

var numericRep = regexp.MustCompile(`^[-+]?[0-9]+(?:\.[0-9]+)?$`)

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

type EChartsXAxisData struct {
	BoundaryGap *bool            `json:"boundaryGap,omitempty"`
	SplitNumber int              `json:"splitNumber,omitempty"`
	AxisLabel   EChartsAxisLabel `json:"axisLabel,omitempty"`
	AxisLine    EChartsAxisLine  `json:"axisLine,omitempty"`
	Data        []string         `json:"data"`
	Type        string           `json:"type"`
}

type EChartsAxisLine struct {
	Show      *bool `json:"show,omitempty"`
	LineStyle struct {
		Color   string   `json:"color,omitempty"`
		Opacity *float64 `json:"opacity,omitempty"`
		Width   *int     `json:"width,omitempty"` // TODO - add support
	} `json:"lineStyle,omitempty"`
}

type EChartsXAxis struct {
	Data []EChartsXAxisData
}

func (ex *EChartsXAxis) UnmarshalJSON(data []byte) error {
	data = convertToArray(data)
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, &ex.Data)
}

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

type EChartsYAxisData struct {
	Min       *float64         `json:"min,omitempty"`
	Max       *float64         `json:"max,omitempty"`
	AxisLabel EChartsAxisLabel `json:"axisLabel,omitempty"`
	AxisLine  EChartsAxisLine  `json:"axisLine,omitempty"`
	Data      []string         `json:"data"`
}

type EChartsYAxis struct {
	Data []EChartsYAxisData `json:"data"`
}

func (ey *EChartsYAxis) UnmarshalJSON(data []byte) error {
	data = convertToArray(data)
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, &ey.Data)
}

type EChartsPadding struct {
	Box Box
}

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

type EChartsLabelOption struct {
	Show     bool   `json:"show"`
	Distance int    `json:"distance"`
	Color    string `json:"color"`
}

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

type EChartsMarkData struct {
	Type string `json:"type"`
	// TODO - support position values below
	XAxis float64 `json:"xAxis,omitempty"`
	YAxis float64 `json:"yAxis,omitempty"`
}
type _EChartsMarkData EChartsMarkData

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

type EChartsMarkPoint struct {
	SymbolSize int               `json:"symbolSize"`
	Data       []EChartsMarkData `json:"data"`
}

func (emp *EChartsMarkPoint) ToSeriesMarkPoint() SeriesMarkPoint {
	return SeriesMarkPoint{
		SymbolSize: emp.SymbolSize,
		Points: sliceConversion(emp.Data, func(i EChartsMarkData) SeriesMark {
			return SeriesMark{Type: i.Type}
		}),
	}
}

type EChartsMarkLine struct {
	Data []EChartsMarkData `json:"data"`
}

func (eml *EChartsMarkLine) ToSeriesMarkLine() SeriesMarkLine {
	return SeriesMarkLine{
		Lines: sliceConversion(eml.Data, func(i EChartsMarkData) SeriesMark {
			return SeriesMark{Type: i.Type}
		}),
	}
}

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
			Values: sliceConversion(item.Data, func(dataItem EChartsSeriesData) float64 {
				return dataItem.Value.First()
			}),
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

type EChartsTextStyle struct {
	Color      string  `json:"color"`
	FontFamily string  `json:"fontFamily"`
	FontSize   float64 `json:"fontSize"`
}

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

type EChartsOption struct {
	Type       string         `json:"type"`
	Theme      string         `json:"theme"`
	FontFamily string         `json:"fontFamily"`
	Padding    EChartsPadding `json:"padding"`
	Box        Box            `json:"box"`
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
		Indicator []RadarIndicator `json:"indicator"`
	} `json:"radar"`
	Series          EChartsSeriesList `json:"series"`
	BackgroundColor string            `json:"backgroundColor,omitempty"`
	Children        []EChartsOption   `json:"children"`
}

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
		RadarIndicators: eo.Radar.Indicator,
		Width:           eo.Width,
		Height:          eo.Height,
		Padding:         eo.Padding.Box,
		Box:             eo.Box,
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
	o.Children = sliceConversion(eo.Children, func(child EChartsOption) ChartOption {
		return child.ToOption()
	})
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

func RenderEChartsToPNG(options string) ([]byte, error) {
	return renderEcharts(options, ChartOutputPNG)
}

func RenderEChartsToJPG(options string) ([]byte, error) {
	return renderEcharts(options, ChartOutputJPG)
}

func RenderEChartsToSVG(options string) ([]byte, error) {
	return renderEcharts(options, ChartOutputSVG)
}
