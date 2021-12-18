// MIT License

// Copyright (c) 2021 Tree Xie

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package charts

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/wcharczuk/go-chart/v2"
)

type EChartStyle struct {
	Color string `json:"color"`
}
type ECharsSeriesData struct {
	Value     float64     `json:"value"`
	Name      string      `json:"name"`
	ItemStyle EChartStyle `json:"itemStyle"`
}
type _ECharsSeriesData ECharsSeriesData

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

func (es *ECharsSeriesData) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return nil
	}
	if regexp.MustCompile(`^\d+`).Match(data) {
		v, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return err
		}
		es.Value = v
		return nil
	}
	v := _ECharsSeriesData{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	es.Name = v.Name
	es.Value = v.Value
	es.ItemStyle = v.ItemStyle
	return nil
}

type EChartsPadding struct {
	box chart.Box
}

func (ep *EChartsPadding) UnmarshalJSON(data []byte) error {
	data = convertToArray(data)
	if len(data) == 0 {
		return nil
	}
	arr := make([]int, 0)
	err := json.Unmarshal(data, &arr)
	if err != nil {
		return err
	}
	if len(arr) == 0 {
		return nil
	}
	switch len(arr) {
	case 1:
		ep.box = chart.Box{
			Left:   arr[0],
			Top:    arr[0],
			Bottom: arr[0],
			Right:  arr[0],
		}
	case 2:
		ep.box = chart.Box{
			Top:    arr[0],
			Bottom: arr[0],
			Left:   arr[1],
			Right:  arr[1],
		}
	default:
		result := make([]int, 4)
		copy(result, arr)
		// 上右下左
		ep.box = chart.Box{
			Top:    arr[0],
			Right:  arr[1],
			Bottom: arr[2],
			Left:   arr[3],
		}
	}
	return nil
}

type EChartsYAxis struct {
	Data []struct {
		Min       *float64 `json:"min"`
		Max       *float64 `json:"max"`
		Interval  int      `json:"interval"`
		AxisLabel struct {
			Formatter string `json:"formatter"`
		} `json:"axisLabel"`
	} `json:"data"`
}

func (ey *EChartsYAxis) UnmarshalJSON(data []byte) error {
	data = convertToArray(data)
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, &ey.Data)
}

type EChartsXAxis struct {
	Data []struct {
		Type        string   `json:"type"`
		BoundaryGap *bool    `json:"boundaryGap"`
		SplitNumber int      `json:"splitNumber"`
		Data        []string `json:"data"`
	}
}

func (ex *EChartsXAxis) UnmarshalJSON(data []byte) error {
	data = convertToArray(data)
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, &ex.Data)
}

type ECharsOptions struct {
	Theme string `json:"theme"`
	Title struct {
		Text string `json:"text"`
		// 暂不支持(go-chart默认title只能居中)
		TextAlign string `json:"textAlign"`
		TextStyle struct {
			Color string `json:"color"`
			// TODO 字体支持
			FontFamily string  `json:"fontFamily"`
			FontSize   float64 `json:"fontSize"`
			Height     float64 `json:"height"`
		} `json:"textStyle"`
	} `json:"title"`
	XAxis  EChartsXAxis `json:"xAxis"`
	YAxis  EChartsYAxis `json:"yAxis"`
	Legend struct {
		Data    []string       `json:"data"`
		Align   string         `json:"align"`
		Padding EChartsPadding `json:"padding"`
	} `json:"legend"`
	Series []struct {
		Data       []ECharsSeriesData `json:"data"`
		Type       string             `json:"type"`
		YAxisIndex int                `json:"yAxisIndex"`
	} `json:"series"`
}

func convertEChartsSeries(e *ECharsOptions) ([]Series, chart.TickPosition) {
	tickPosition := chart.TickPositionUnset

	if len(e.Series) == 0 {
		return nil, tickPosition
	}
	seriesType := e.Series[0].Type
	if seriesType == SeriesPie {
		series := make([]Series, len(e.Series[0].Data))
		for index, item := range e.Series[0].Data {
			series[index] = Series{
				Data: []SeriesData{
					{
						Value: item.Value,
					},
				},
				Type: seriesType,
				Name: item.Name,
			}
		}
		return series, tickPosition
	}
	series := make([]Series, len(e.Series))
	for index, item := range e.Series {
		// bar默认tick居中
		if item.Type == SeriesBar {
			tickPosition = chart.TickPositionBetweenTicks
		}
		data := make([]SeriesData, len(item.Data))
		for j, itemData := range item.Data {
			sd := SeriesData{
				Value: itemData.Value,
			}
			if itemData.ItemStyle.Color != "" {
				c := parseColor(itemData.ItemStyle.Color)
				sd.Style.FillColor = c
				sd.Style.StrokeColor = c
			}
			data[j] = sd
		}
		series[index] = Series{
			YAxisIndex: item.YAxisIndex,
			Data:       data,
			Type:       item.Type,
		}
	}
	return series, tickPosition
}

func (e *ECharsOptions) ToOptions() Options {
	o := Options{
		Theme: e.Theme,
	}
	titleTextStyle := e.Title.TextStyle
	o.Title = Title{
		Text: e.Title.Text,
		Style: chart.Style{
			FontColor: parseColor(titleTextStyle.Color),
			FontSize:  titleTextStyle.FontSize,
		},
	}

	if titleTextStyle.FontSize != 0 && titleTextStyle.Height > titleTextStyle.FontSize {
		padding := int(titleTextStyle.Height-titleTextStyle.FontSize) / 2
		o.Title.Style.Padding.Top = padding
		o.Title.Style.Padding.Bottom = padding
	}

	boundaryGap := false
	if len(e.XAxis.Data) != 0 {
		xAxis := e.XAxis.Data[0]
		o.XAxis = XAxis{
			Type:        xAxis.Type,
			Data:        xAxis.Data,
			SplitNumber: xAxis.SplitNumber,
		}
		if xAxis.BoundaryGap == nil || *xAxis.BoundaryGap {
			boundaryGap = true
		}
	}

	o.Legend = Legend{
		Data:    e.Legend.Data,
		Align:   e.Legend.Align,
		Padding: e.Legend.Padding.box,
	}
	if len(e.YAxis.Data) != 0 {
		yAxisOptions := make([]*YAxisOption, len(e.YAxis.Data))
		for index, item := range e.YAxis.Data {
			opt := &YAxisOption{
				Max: item.Max,
				Min: item.Min,
			}
			template := item.AxisLabel.Formatter
			if template != "" {
				opt.Formater = func(v interface{}) string {
					str := defaultFloatFormater(v)
					return strings.ReplaceAll(template, "{value}", str)
				}
			}
			yAxisOptions[index] = opt
		}
		o.YAxisOptions = yAxisOptions
	}

	series, tickPosition := convertEChartsSeries(e)

	o.Series = series

	if boundaryGap {
		tickPosition = chart.TickPositionBetweenTicks
	}
	o.TickPosition = tickPosition
	return o
}

func ParseECharsOptions(options string) (Options, error) {
	e := ECharsOptions{}
	err := json.Unmarshal([]byte(options), &e)
	if err != nil {
		return Options{}, err
	}

	return e.ToOptions(), nil
}
