package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAxisRender(t *testing.T) {
	t.Parallel()

	dayLabels := []string{
		"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun",
	}
	letterLabels := []string{"A", "B", "C", "D", "E", "F", "G"}
	fs := NewFontStyleWithSize(18)
	fs.FontColor = ColorGreen
	axisTheme := GetDefaultTheme().
		WithXAxisColor(ColorBlue).
		WithYAxisColor(ColorBlue).
		WithAxisSplitLineColor(ColorGray)

	tests := []struct {
		name          string
		optionFactory func(p *Painter) axisOption
	}{
		{
			name: "x-axis",
			optionFactory: func(p *Painter) axisOption {
				opt := XAxisOption{
					BoundaryGap:    Ptr(true),
					LabelFontStyle: fs,
				}
				return opt.prep(axisTheme, false).toAxisOption(newTestRangeForLabels(dayLabels, 0, fs))
			},
		},
		{
			name: "x-axis_rotation45",
			optionFactory: func(p *Painter) axisOption {
				opt := XAxisOption{
					BoundaryGap:    Ptr(true),
					LabelFontStyle: fs,
				}
				return opt.prep(axisTheme, false).toAxisOption(newTestRangeForLabels(dayLabels, DegreesToRadians(45), fs))
			},
		},
		{
			name: "x-axis_rotation90",
			optionFactory: func(p *Painter) axisOption {
				opt := XAxisOption{
					Labels:         dayLabels,
					BoundaryGap:    Ptr(true),
					LabelFontStyle: fs,
				}
				return opt.prep(axisTheme, false).toAxisOption(newTestRangeForLabels(dayLabels, DegreesToRadians(90), fs))
			},
		},
		{
			name: "x-axis_splitline",
			optionFactory: func(p *Painter) axisOption {
				return axisOption{
					theme:         axisTheme,
					aRange:        newTestRangeForLabels(letterLabels, 0, fs),
					splitLineShow: Ptr(true),
				}
			},
		},
		{
			name: "y-axis_left",
			optionFactory: func(p *Painter) axisOption {
				opt := YAxisOption{
					Position:       PositionLeft,
					isCategoryAxis: true,
				}
				return opt.prep(axisTheme, true).toAxisOption(newTestRangeForLabels(dayLabels, 0, fs))
			},
		},
		{
			name: "y-axis_right",
			optionFactory: func(p *Painter) axisOption {
				return axisOption{
					theme:         axisTheme,
					aRange:        newTestRangeForLabels(dayLabels, 0, fs),
					position:      PositionRight,
					boundaryGap:   Ptr(false),
					splitLineShow: Ptr(true),
				}
			},
		},
		{
			name: "reduced_label_count",
			optionFactory: func(p *Painter) axisOption {
				aRange := newTestRangeForLabels(letterLabels, 0, fs)
				aRange.labelCount -= 2
				return axisOption{
					theme:         axisTheme,
					aRange:        aRange,
					splitLineShow: Ptr(false),
				}
			},
		},
		{
			name: "custom_font",
			optionFactory: func(p *Painter) axisOption {
				fs := FontStyle{
					FontSize:  40.0,
					FontColor: ColorBlue,
				}
				return axisOption{
					theme:  axisTheme,
					aRange: newTestRangeForLabels(letterLabels, 0, fs),
				}
			},
		},
		{
			name: "boundary_gap_disable",
			optionFactory: func(p *Painter) axisOption {
				return axisOption{
					theme:       axisTheme,
					aRange:      newTestRangeForLabels(letterLabels, 0, fs),
					boundaryGap: Ptr(false),
				}
			},
		},
		{
			name: "boundary_gap_enable",
			optionFactory: func(p *Painter) axisOption {
				return axisOption{
					theme:       axisTheme,
					aRange:      newTestRangeForLabels(letterLabels, 0, fs),
					boundaryGap: Ptr(true),
				}
			},
		},
		{
			name: "dense_category_data",
			optionFactory: func(p *Painter) axisOption {
				const count = 1000
				labelLen := len(strconv.Itoa(count))
				labels := make([]string, count)
				tsl := testSeriesList{}
				for i := range labels {
					label := strconv.Itoa(i + 1)
					for len(label) < labelLen {
						label = "0" + label
					}
					labels[i] = label
					tsl = append(tsl, testSeries{values: []float64{float64(i)}})
				}
				return axisOption{
					theme: axisTheme,
					aRange: calculateCategoryAxisRange(p, p.Width(), false, false, labels,
						0, 0, 0, tsl, 0, fs),
					boundaryGap: Ptr(true),
				}
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			}, PainterPaddingOption(NewBoxEqual(50)))

			opt := tt.optionFactory(p)
			_, err := newAxisPainter(p, opt).Render()
			require.NoError(t, err)
			data, err := p.Bytes()
			require.NoError(t, err)
			assertTestdataSVG(t, data)
		})
	}
}

func TestBottomXAxis(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		makeOption func() XAxisOption
		makeValue  bool
	}{
		{
			name: "basic",
			makeOption: func() XAxisOption {
				return XAxisOption{
					Labels: []string{"a", "b", "c", "d"},
				}
			},
		},
		{
			name: "value",
			makeOption: func() XAxisOption {
				return XAxisOption{
					Labels: []string{"a", "b", "c", "d"},
				}
			},
			makeValue: true,
		},
		{
			name: "title",
			makeOption: func() XAxisOption {
				return XAxisOption{
					Title:  "Title",
					Labels: []string{"a", "b", "c", "d"},
				}
			},
		},
		{
			name: "boundary_gap_disabled",
			makeOption: func() XAxisOption {
				return XAxisOption{
					Labels:      []string{"a", "b", "c", "d"},
					BoundaryGap: Ptr(false),
				}
			},
		},
		{
			name: "font_style",
			makeOption: func() XAxisOption {
				return XAxisOption{
					Labels:         []string{"abc", "def", "ghi"},
					LabelFontStyle: NewFontStyleWithSize(20),
					LabelRotation:  DegreesToRadians(90),
				}
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			theme := GetTheme(ThemeLight)
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			}, PainterThemeOption(theme), PainterPaddingOption(NewBoxEqual(100)))

			xAxisOpt := tt.makeOption()
			xAxisOpt = *xAxisOpt.prep(theme, false)
			aRange := newTestRangeForLabels(xAxisOpt.Labels, xAxisOpt.LabelRotation,
				fillFontStyleDefaults(xAxisOpt.LabelFontStyle, defaultFontSize, theme.GetXAxisTextColor()))
			aRange.isCategory = !tt.makeValue
			_, err := newAxisPainter(p, xAxisOpt.toAxisOption(aRange)).Render()
			require.NoError(t, err)
			data, err := p.Bytes()
			require.NoError(t, err)
			assertTestdataSVG(t, data)
		})
	}
}

func TestYAxis(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		makeOption       func() *YAxisOption
		makeCategory     bool
		disableSplitLine bool
	}{
		{
			name: "basic",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Position: PositionLeft, // typically defaulted in defaultRender
					Labels:   []string{"a", "b", "c", "d"},
				}
			},
		},
		{
			name: "category",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Position:       PositionLeft, // typically defaulted in defaultRender
					Labels:         []string{"a", "b", "c", "d"},
					isCategoryAxis: true,
				}
			},
			makeCategory: true,
		},
		{
			name: "font_style_with_rotation",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Position:       PositionLeft, // typically defaulted in defaultRender
					Labels:         []string{"a", "b", "c"},
					LabelFontStyle: NewFontStyleWithSize(20),
					LabelRotation:  DegreesToRadians(270),
				}
			},
		},
		{
			name: "lines",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Position:      PositionLeft, // typically defaulted in defaultRender
					Labels:        []string{"a", "b", "c", "d"},
					SplitLineShow: Ptr(true),
					SpineLineShow: Ptr(true),
				}
			},
		},
		{
			name: "title_left",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Position: PositionLeft, // typically defaulted in defaultRender
					Title:    "title",
					Labels:   []string{"a", "b", "c", "d"},
				}
			},
		},
		{
			name: "title_right",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Position: PositionRight,
					Title:    "title",
					Labels:   []string{"a", "b", "c", "d"},
				}
			},
		},
		{
			name: "split_line_disable",
			makeOption: func() *YAxisOption {
				return &YAxisOption{
					Position: PositionRight,
					Labels:   []string{"a", "b", "c", "d"},
				}
			},
			disableSplitLine: true,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i)+"-"+tt.name, func(t *testing.T) {
			theme := GetTheme(ThemeLight)
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        600,
				Height:       400,
			}, PainterThemeOption(theme),
				PainterPaddingOption(NewBoxEqual(100)))

			yAxisOpt := tt.makeOption()
			yAxisOpt = yAxisOpt.prep(theme, true)
			aRange := newTestRangeForLabels(yAxisOpt.Labels, yAxisOpt.LabelRotation,
				fillFontStyleDefaults(yAxisOpt.LabelFontStyle, defaultFontSize, theme.GetYAxisTextColor()))
			aRange.isCategory = tt.makeCategory
			axisOpt := yAxisOpt.toAxisOption(aRange)
			if tt.disableSplitLine {
				axisOpt.splitLineShow = Ptr(false)
			}
			_, err := newAxisPainter(p, axisOpt).Render()
			require.NoError(t, err)
			data, err := p.Bytes()
			require.NoError(t, err)
			assertTestdataSVG(t, data)
		})
	}
}
