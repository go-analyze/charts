package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-analyze/charts/chartdraw/matrix"
)

func TestLabelFormatPie(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "a: 12%",
		labelFormatPie("a", "", nil, 10, 0.12))

	assert.Equal(t, "a: f",
		labelFormatPie("a", "{b}: {c}", func(f float64) string {
			return "f"
		}, 10, 0.12))
}

func TestLabelFormatFunnel(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "a(12%)", labelFormatFunnel("a", "", nil, 10, 0.12))

	assert.Equal(t, "b(f, 25%)",
		labelFormatFunnel("b", "{b}({c}, {d})", func(f float64) string {
			return "f"
		}, 20, 0.25))
}

func TestLabelFormatter(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "10",
		labelFormatValue([]string{"a", "b"}, "", nil, 0, 10, 0.12))

	assert.Equal(t, "f f 12%",
		labelFormatValue([]string{"a", "b"}, "{c} {c} {d}",
			func(f float64) string {
				return "f"
			},
			0, 10, 0.12))

	assert.Equal(t, "Name: a, Value: 10, Percent: 12%",
		labelFormatPie("a", "Name: {b}, Value: {c}, Percent: {d}", nil, 10, 0.12))
}

func TestSeriesLabelFormatter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		label    SeriesLabel
		values   []labelValue
		expected []string
	}{
		{
			name: "label_formatter_with_index_based_highlighting",
			label: SeriesLabel{
				Show: Ptr(true),
				LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
					if index == 1 { // highlight only second value
						return "⭐ " + name + ": " + strconv.FormatFloat(val, 'f', 1, 64), nil
					}
					return "", nil // hide other labels
				},
			},
			values: []labelValue{
				{index: 0, value: 10.5},
				{index: 1, value: 20.7},
				{index: 2, value: 30.2},
			},
			expected: []string{"⭐ series1: 20.7"},
		},
		{
			name: "label_formatter_returns_empty_string",
			label: SeriesLabel{
				Show: Ptr(true),
				LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
					return "", nil // always return empty
				},
			},
			values: []labelValue{
				{index: 0, value: 10.5},
			},
			expected: []string{}, // should result in no rendered text
		},
		{
			name: "label_formatter_with_custom_format",
			label: SeriesLabel{
				Show: Ptr(true),
				LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
					return "Item " + strconv.Itoa(index) + ": $" + strconv.FormatFloat(val, 'f', 2, 64), nil
				},
			},
			values: []labelValue{
				{index: 0, value: 100.50},
				{index: 1, value: 200.75},
			},
			expected: []string{"Item 0: $100.50", "Item 1: $200.75"},
		},
		{
			name: "fallback_to_value_formatter_when_label_formatter_is_nil",
			label: SeriesLabel{
				Show: Ptr(true),
				ValueFormatter: func(f float64) string {
					return "Value: " + strconv.FormatFloat(f, 'f', 1, 64)
				},
			},
			values: []labelValue{
				{index: 0, value: 42.3},
			},
			expected: []string{"Value: 42.3"},
		},
		{
			name: "fallback_to_format_template_when_label_formatter_is_nil",
			label: SeriesLabel{
				Show:           Ptr(true),
				FormatTemplate: "Series {b} = {c}",
			},
			values: []labelValue{
				{index: 0, value: 42.3},
			},
			expected: []string{"Series series0 = 42.3"},
		},
		{
			name: "label_formatter_takes_precedence_over_format_template",
			label: SeriesLabel{
				Show:           Ptr(true),
				FormatTemplate: "Template: {c}",
				LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
					return "Formatter: " + strconv.FormatFloat(val, 'f', 0, 64), nil
				},
			},
			values: []labelValue{
				{index: 0, value: 42.3},
			},
			expected: []string{"Formatter: 42"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        400,
				Height:       300,
			}, PainterThemeOption(GetTheme(ThemeLight)))

			seriesNames := make([]string, len(tt.values))
			for i := range seriesNames {
				seriesNames[i] = "series" + strconv.Itoa(i)
			}

			painter := newSeriesLabelPainter(p, seriesNames, tt.label, GetTheme(ThemeLight))

			// Add all values to the painter
			for _, value := range tt.values {
				painter.Add(value)
			}

			// Filter out empty text entries for comparison
			var renderedText []string
			for _, v := range painter.values {
				if v.text != "" {
					renderedText = append(renderedText, v.text)
				}
			}

			// Handle empty slice vs nil slice case
			if len(tt.expected) == 0 && len(renderedText) == 0 {
				// Both empty, this is a pass
			} else {
				// Check that the correct text was generated
				assert.Equal(t, tt.expected, renderedText)
			}
		})
	}

	t.Run("empty_series_names", func(t *testing.T) {
		p := NewPainter(PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        400,
			Height:       300,
		}, PainterThemeOption(GetTheme(ThemeLight)))

		label := SeriesLabel{
			Show: Ptr(true),
			LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
				return "Index: " + strconv.Itoa(index) + ", Name: '" + name + "'", nil
			},
		}

		painter := newSeriesLabelPainter(p, []string{}, label, GetTheme(ThemeLight))
		painter.Add(labelValue{index: 0, value: 42.5})

		// Should handle empty series names gracefully
		assert.Len(t, painter.values, 1)
		assert.Equal(t, "Index: 0, Name: ''", painter.values[0].text)
	})

	t.Run("out_of_bounds_series_index", func(t *testing.T) {
		p := NewPainter(PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        400,
			Height:       300,
		}, PainterThemeOption(GetTheme(ThemeLight)))

		label := SeriesLabel{
			Show: Ptr(true),
			LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
				return "Index: " + strconv.Itoa(index) + ", Name: '" + name + "'", nil
			},
		}

		painter := newSeriesLabelPainter(p, []string{"Series0"}, label, GetTheme(ThemeLight))
		painter.Add(labelValue{index: 5, value: 42.5}) // Index beyond series names

		// Should handle out of bounds gracefully with empty name
		assert.Len(t, painter.values, 1)
		assert.Equal(t, "Index: 5, Name: ''", painter.values[0].text)
	})

	t.Run("show_flag_false", func(t *testing.T) {
		p := NewPainter(PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        400,
			Height:       300,
		}, PainterThemeOption(GetTheme(ThemeLight)))

		label := SeriesLabel{
			Show: Ptr(false), // Explicitly hide labels
			LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
				return "This should not appear", nil
			},
		}

		painter := newSeriesLabelPainter(p, []string{"Test"}, label, GetTheme(ThemeLight))
		painter.Add(labelValue{index: 0, value: 42.5})

		// Should not add any values when Show is false
		assert.Empty(t, painter.values)
	})

	t.Run("nil_label_formatter_fallback", func(t *testing.T) {
		p := NewPainter(PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        400,
			Height:       300,
		}, PainterThemeOption(GetTheme(ThemeLight)))

		label := SeriesLabel{
			Show: Ptr(true),
			// No LabelFormatter, should fall back to defaultValueFormatter
		}

		painter := newSeriesLabelPainter(p, []string{"Test"}, label, GetTheme(ThemeLight))
		painter.Add(labelValue{index: 0, value: 123.456})

		// Should use defaultValueFormatter
		assert.Len(t, painter.values, 1)
		assert.Equal(t, "123.46", painter.values[0].text) // defaultValueFormatter truncates to 2 decimals
	})
}

func TestDrawLabelWithBackground(t *testing.T) {
	t.Parallel()

	t.Run("empty_text", func(t *testing.T) {
		p := NewPainter(PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        400,
			Height:       400,
		})

		fontStyle := FontStyle{FontColor: ColorBlack, FontSize: 12}
		drawLabelWithBackground(p, "", 50, 50, 0, fontStyle, ColorRed, 5, ColorTransparent, 0)

		data, err := p.Bytes()
		require.NoError(t, err)
		assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 400\"></svg>", data)
	})

	t.Run("transparent_background", func(t *testing.T) {
		p := NewPainter(PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        400,
			Height:       400,
		})

		fontStyle := FontStyle{FontColor: ColorBlack, FontSize: 12}
		drawLabelWithBackground(p, "test", 50, 50, 0, fontStyle, ColorTransparent, 5, ColorTransparent, 0)

		data, err := p.Bytes()
		require.NoError(t, err)
		assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 400\"><text x=\"50\" y=\"50\" style=\"stroke:none;fill:black;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">test</text></svg>", data)
	})

	t.Run("solid_background", func(t *testing.T) {
		p := NewPainter(PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        400,
			Height:       400,
		})

		fontStyle := FontStyle{FontColor: ColorBlack, FontSize: 12}
		drawLabelWithBackground(p, "test", 50, 50, 0, fontStyle, ColorBlue, 0, ColorTransparent, 0)

		data, err := p.Bytes()
		require.NoError(t, err)
		assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 400\"><path  d=\"M 46 30\nL 81 30\nL 81 54\nL 46 54\nL 46 30\" style=\"stroke:none;fill:blue\"/><text x=\"50\" y=\"50\" style=\"stroke:none;fill:black;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">test</text></svg>", data)
	})

	t.Run("rounded_background", func(t *testing.T) {
		p := NewPainter(PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        400,
			Height:       400,
		})

		fontStyle := FontStyle{FontColor: ColorBlack, FontSize: 12}
		drawLabelWithBackground(p, "test", 50, 50, 0, fontStyle, ColorGreen, 5, ColorTransparent, 0)

		data, err := p.Bytes()
		require.NoError(t, err)
		assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 400\"><path  d=\"M 51 30\nL 76 30\nL 76 30\nA 5 5 90.00 0 1 81 35\nL 81 49\nL 81 49\nA 5 5 90.00 0 1 76 54\nL 51 54\nL 51 54\nA 5 5 90.00 0 1 46 49\nL 46 35\nL 46 35\nA 5 5 90.00 0 1 51 30\nZ\" style=\"stroke:none;fill:green\"/><text x=\"50\" y=\"50\" style=\"stroke:none;fill:black;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">test</text></svg>", data)
	})

	t.Run("negative_corner_radius", func(t *testing.T) {
		p := NewPainter(PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        400,
			Height:       400,
		})

		fontStyle := FontStyle{FontColor: ColorBlack, FontSize: 12}
		// Negative corner radius should be treated as 0 (square corners)
		drawLabelWithBackground(p, "test", 50, 50, 0, fontStyle, ColorRed, -5, ColorTransparent, 0)

		data, err := p.Bytes()
		require.NoError(t, err)
		assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 400\"><path  d=\"M 46 30\nL 81 30\nL 81 54\nL 46 54\nL 46 30\" style=\"stroke:none;fill:red\"/><text x=\"50\" y=\"50\" style=\"stroke:none;fill:black;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">test</text></svg>", data)
	})
}

func TestLabelStyleOverrides(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                    string
		label                   SeriesLabel
		values                  []labelValue
		expectedText            []string
		expectedFontSize        []float64
		expectedFontColor       []Color
		expectedDistance        []int
		expectedOffset          []OffsetInt
		expectedBackgroundColor []Color
		expectedCornerRadius    []int
	}{
		{
			name: "nil_override_uses_base_styles",
			label: SeriesLabel{
				Show: Ptr(true),
				LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
					return "base style", nil
				},
				FontStyle: FontStyle{FontSize: 12, FontColor: ColorRed},
				Distance:  10,
			},
			values: []labelValue{
				{index: 0, value: 10.5, offset: OffsetInt{Left: 1, Top: 2}},
			},
			expectedText:            []string{"base style"},
			expectedFontSize:        []float64{12},
			expectedFontColor:       []Color{ColorRed},
			expectedDistance:        []int{10},
			expectedOffset:          []OffsetInt{{Left: 1, Top: 2}},
			expectedBackgroundColor: []Color{{}}, // Zero color value
			expectedCornerRadius:    []int{0},
		},
		{
			name: "override_some_fields",
			label: SeriesLabel{
				Show: Ptr(true),
				LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
					return "style override", &LabelStyle{
						FontStyle: FontStyle{FontSize: 16}, // only override size
					}
				},
				FontStyle: FontStyle{FontSize: 12, FontColor: ColorRed},
			},
			values: []labelValue{
				{index: 0, value: 10.5},
			},
			expectedText:            []string{"style override"},
			expectedFontSize:        []float64{16},     // overridden
			expectedFontColor:       []Color{ColorRed}, // from base
			expectedBackgroundColor: []Color{{}},       // Zero color value
			expectedCornerRadius:    []int{0},
		},
		{
			name: "empty_text_hides_label",
			label: SeriesLabel{
				Show: Ptr(true),
				LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
					return "", &LabelStyle{
						FontStyle: FontStyle{FontSize: 20},
					} // style ignored when text is empty
				},
			},
			values: []labelValue{
				{index: 0, value: 10.5},
			},
			expectedText:            []string{}, // no label rendered
			expectedFontSize:        []float64{},
			expectedFontColor:       []Color{},
			expectedBackgroundColor: []Color{},
			expectedCornerRadius:    []int{},
		},
		{
			name: "per_point_color_changes",
			label: SeriesLabel{
				Show: Ptr(true),
				LabelFormatter: func(index int, name string, val float64) (string, *LabelStyle) {
					if index == 0 {
						return "red", &LabelStyle{
							FontStyle: FontStyle{FontColor: ColorRed},
						}
					}
					return "blue", &LabelStyle{
						FontStyle: FontStyle{FontColor: ColorBlue},
					}
				},
			},
			values: []labelValue{
				{index: 0, value: 10.5},
				{index: 1, value: 20.5},
			},
			expectedText:            []string{"red", "blue"},
			expectedFontColor:       []Color{ColorRed, ColorBlue},
			expectedBackgroundColor: []Color{{}, {}}, // Zero color values
			expectedCornerRadius:    []int{0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPainter(PainterOptions{
				OutputFormat: ChartOutputSVG,
				Width:        800,
				Height:       800,
			}, PainterThemeOption(GetTheme(ThemeLight)))

			seriesNames := make([]string, len(tt.values))
			for i := range seriesNames {
				seriesNames[i] = "series" + strconv.Itoa(i)
			}

			painter := newSeriesLabelPainter(p, seriesNames, tt.label, GetTheme(ThemeLight))

			// Add all values to the painter
			for _, value := range tt.values {
				painter.Add(value)
			}

			// Filter out empty text entries for most checks (they shouldn't be rendered)
			renderedValues := make([]labelRenderValue, 0, len(painter.values))
			actualRenderedText := make([]string, 0, len(painter.values))
			for _, v := range painter.values {
				if v.text != "" {
					renderedValues = append(renderedValues, v)
					actualRenderedText = append(actualRenderedText, v.text)
				}
			}

			// Check text (handle empty slice vs nil slice case)
			if len(tt.expectedText) == 0 && len(actualRenderedText) == 0 {
				// Both empty, this is a pass
			} else {
				assert.Equal(t, tt.expectedText, actualRenderedText)
			}

			// Check font size for rendered values
			if len(tt.expectedFontSize) > 0 {
				assert.Equal(t, len(tt.expectedFontSize), len(renderedValues))
				for i, expected := range tt.expectedFontSize {
					if i < len(renderedValues) {
						assert.InDelta(t, expected, renderedValues[i].fontStyle.FontSize, matrix.DefaultEpsilon)
					}
				}
			}

			// Check font color for rendered values
			if len(tt.expectedFontColor) > 0 {
				assert.Equal(t, len(tt.expectedFontColor), len(renderedValues))
				for i, expected := range tt.expectedFontColor {
					if i < len(renderedValues) {
						assert.Equal(t, expected, renderedValues[i].fontStyle.FontColor)
					}
				}
			}

			// Check background color for rendered values
			if len(tt.expectedBackgroundColor) > 0 {
				assert.Equal(t, len(tt.expectedBackgroundColor), len(renderedValues))
				for i, expected := range tt.expectedBackgroundColor {
					if i < len(renderedValues) {
						assert.Equal(t, expected, renderedValues[i].backgroundColor)
					}
				}
			}

			// Check corner radius for rendered values
			if len(tt.expectedCornerRadius) > 0 {
				assert.Equal(t, len(tt.expectedCornerRadius), len(renderedValues))
				for i, expected := range tt.expectedCornerRadius {
					if i < len(renderedValues) {
						assert.Equal(t, expected, renderedValues[i].cornerRadius)
					}
				}
			}
		})
	}
}

func TestLabelFormatterThresholdMin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		threshold float64
		testCases []struct {
			value    float64
			expected string
		}
	}{
		{
			name:      "threshold_100",
			threshold: 100,
			testCases: []struct {
				value    float64
				expected string
			}{
				{50, ""},
				{100, "100"},
				{150, "150"},
				{200.5, "200.5"},
			},
		},
		{
			name:      "threshold_0",
			threshold: 0,
			testCases: []struct {
				value    float64
				expected string
			}{
				{-10, ""},
				{0, "0"},
				{0.1, "0.1"},
				{50, "50"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := LabelFormatterThresholdMin(tt.threshold)

			for _, tc := range tt.testCases {
				text, style := formatter(0, "test", tc.value)
				assert.Equal(t, tc.expected, text)
				assert.Nil(t, style)
			}
		})
	}

	// Boundary conditions
	t.Run("boundary_conditions", func(t *testing.T) {
		t.Run("large_number", func(t *testing.T) {
			formatter := LabelFormatterThresholdMin(1e6)

			text, style := formatter(0, "test", 999999)
			assert.Equal(t, "", text)
			assert.Nil(t, style)

			text, style = formatter(0, "test", 1000000)
			assert.Equal(t, "1M", text) // defaultValueFormatter uses humanize for large numbers
			assert.Nil(t, style)
		})
	})
}

func TestLabelFormatterThresholdMax(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		threshold float64
		testCases []struct {
			value    float64
			expected string
		}
	}{
		{
			name:      "threshold_100",
			threshold: 100,
			testCases: []struct {
				value    float64
				expected string
			}{
				{50, "50"},
				{100, "100"},
				{150, ""},
				{200.5, ""},
			},
		},
		{
			name:      "threshold_0",
			threshold: 0,
			testCases: []struct {
				value    float64
				expected string
			}{
				{-10, "-10"},
				{0, "0"},
				{10, ""},
				{100, ""},
			},
		},
		{
			name:      "threshold_negative",
			threshold: -50,
			testCases: []struct {
				value    float64
				expected string
			}{
				{-100, "-100"},
				{-50, "-50"},
				{-25, ""},
				{0, ""},
				{50, ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := LabelFormatterThresholdMax(tt.threshold)
			for _, tc := range tt.testCases {
				text, style := formatter(0, "", tc.value)
				assert.Equal(t, tc.expected, text)
				assert.Nil(t, style)
			}
		})
	}
}

func TestLabelFormatterTopN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		values    []float64
		n         int
		testCases []struct {
			index    int
			value    float64
			expected string
		}
	}{
		{
			name:   "top_2_of_5",
			values: []float64{100, 200, 50, 300, 75},
			n:      2,
			testCases: []struct {
				index    int
				value    float64
				expected string
			}{
				{0, 100, ""},
				{1, 200, "200"}, // 2nd highest
				{2, 50, ""},
				{3, 300, "300"}, // highest
				{4, 75, ""},
			},
		},
		{
			name:   "top_1_of_3",
			values: []float64{10, 30, 20},
			n:      1,
			testCases: []struct {
				index    int
				value    float64
				expected string
			}{
				{0, 10, ""},
				{1, 30, "30"}, // highest
				{2, 20, ""},
			},
		},
		{
			name:   "top_5_of_3_all_shown",
			values: []float64{10, 30, 20},
			n:      5,
			testCases: []struct {
				index    int
				value    float64
				expected string
			}{
				{0, 10, "10"}, // all shown when n > series length
				{1, 30, "30"},
				{2, 20, "20"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := LabelFormatterTopN(tt.values, tt.n)

			for _, tc := range tt.testCases {
				text, style := formatter(tc.index, "test", tc.value)
				assert.Equal(t, tc.expected, text)
				assert.Nil(t, style)
			}
		})
	}

	edgeTests := []struct {
		name          string
		values        []float64
		n             int
		testValue     float64
		expectedText  string
		shouldShowAll bool // whether it should behave like LabelFormatterValueShort
	}{
		{"negative_n", []float64{10, 20, 30}, -1, 20.0, "", false},
		{"zero_n", []float64{10, 20, 30}, 0, 20.0, "", false},
		{"n_greater_than_length", []float64{10, 20}, 5, 20.0, "20", true},
		{"n_equal_to_length", []float64{10, 20}, 2, 20, "20", true},
		{"single_value", []float64{42}, 1, 42.0, "42", true},
	}

	for _, tt := range edgeTests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := LabelFormatterTopN(tt.values, tt.n)

			text, style := formatter(0, "test label", tt.testValue)

			assert.Equal(t, tt.expectedText, text)
			assert.Nil(t, style) // TopN doesn't add style, only controls visibility
		})
	}
}

func TestLabelFormatterGradientGreenRed(t *testing.T) {
	t.Parallel()

	values := []float64{10, 30, 20}

	formatter := LabelFormatterGradientGreenRed(values)

	tests := []struct {
		name          string
		value         float64
		expected      string
		expectedColor Color
	}{
		{"min_value", 10, "10", ColorGreen},      // minimum gets green
		{"mid_value", 20, "20", ColorYellowAlt1}, // middle gets yellowalt1
		{"max_value", 30, "30", ColorRed},        // maximum gets red
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, style := formatter(0, "test series label", tt.value)
			assert.Equal(t, tt.expected, text)
			assert.NotNil(t, style)
			assert.Equal(t, tt.expectedColor, style.FontStyle.FontColor)
		})
	}
}

func TestLabelFormatterGradientRedGreen(t *testing.T) {
	t.Parallel()

	values := []float64{10, 30, 20}

	formatter := LabelFormatterGradientRedGreen(values)

	tests := []struct {
		name          string
		value         float64
		expected      string
		expectedColor Color
	}{
		{"min_value", 10, "10", ColorRed},        // minimum gets red
		{"mid_value", 20, "20", ColorYellowAlt1}, // middle gets yellowalt1
		{"max_value", 30, "30", ColorGreen},      // maximum gets green
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, style := formatter(0, "test series label", tt.value)
			assert.Equal(t, tt.expected, text)
			assert.NotNil(t, style)
			assert.Equal(t, tt.expectedColor, style.FontStyle.FontColor)
		})
	}
}

func TestLabelFormatterGradientColor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		values    []float64
		colors    []Color
		testCases []struct {
			value         float64
			expectedText  string
			expectedColor Color
		}
	}{
		{
			name:   "blue_to_red_gradient_two_colors",
			values: []float64{0, 50, 100},
			colors: []Color{ColorBlue, ColorRed},
			testCases: []struct {
				value         float64
				expectedText  string
				expectedColor Color
			}{
				{0, "0", ColorBlue},    // minimum gets low color
				{100, "100", ColorRed}, // maximum gets high color
			},
		},
		{
			name:   "green_yellow_red_gradient_three_colors",
			values: []float64{10, 20, 30},
			colors: []Color{ColorGreen, ColorYellowAlt1, ColorRed},
			testCases: []struct {
				value         float64
				expectedText  string
				expectedColor Color
			}{
				{10, "10", ColorGreen},      // minimum gets green
				{20, "20", ColorYellowAlt1}, // middle gets yellowalt1
				{30, "30", ColorRed},        // maximum gets red
			},
		},
		{
			name:   "same_values_all_first_color",
			values: []float64{50, 50, 50},
			colors: []Color{ColorGreen, ColorYellowAlt1},
			testCases: []struct {
				value         float64
				expectedText  string
				expectedColor Color
			}{
				{50, "50", ColorGreen}, // all same values get first color
			},
		},
		{
			name:   "single_color",
			values: []float64{10, 20, 30},
			colors: []Color{ColorBlue},
			testCases: []struct {
				value         float64
				expectedText  string
				expectedColor Color
			}{
				{10, "10", ColorBlue}, // all values get the single color
				{20, "20", ColorBlue},
				{30, "30", ColorBlue},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := LabelFormatterGradientColor(tt.values, tt.colors...)

			for _, tc := range tt.testCases {
				text, style := formatter(0, "test", tc.value)
				assert.Equal(t, tc.expectedText, text)
				assert.NotNil(t, style)
				assert.Equal(t, tc.expectedColor, style.FontStyle.FontColor)
			}
		})
	}

	edgeTests := []struct {
		name            string
		values          []float64
		colors          []Color
		testValue       float64
		expectedText    string
		expectedColor   Color
		shouldHaveStyle bool
	}{
		{
			name:            "empty_values",
			values:          []float64{},
			colors:          []Color{ColorRed, ColorBlue},
			testValue:       50,
			expectedText:    "50",
			expectedColor:   ColorRed, // should default to first color
			shouldHaveStyle: true,
		},
		{
			name:            "no_colors",
			values:          []float64{10, 20, 30},
			colors:          []Color{},
			testValue:       20,
			expectedText:    "20",
			expectedColor:   ColorBlack, // fallback to black
			shouldHaveStyle: true,
		},
		{
			name:            "nil_values",
			values:          nil,
			colors:          []Color{ColorGreen},
			testValue:       100,
			expectedText:    "100",
			expectedColor:   ColorGreen,
			shouldHaveStyle: true,
		},
		{
			name:            "single_value",
			values:          []float64{42},
			colors:          []Color{ColorYellow, ColorRed},
			testValue:       42,
			expectedText:    "42",
			expectedColor:   ColorYellow, // all same values get first color
			shouldHaveStyle: true,
		},
	}

	for _, tt := range edgeTests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := LabelFormatterGradientColor(tt.values, tt.colors...)
			text, style := formatter(0, "test", tt.testValue)

			assert.Equal(t, tt.expectedText, text)
			if tt.shouldHaveStyle {
				assert.NotNil(t, style)
				assert.Equal(t, tt.expectedColor, style.FontStyle.FontColor)
			} else {
				assert.Nil(t, style)
			}
		})
	}
}

func TestLabelFormatterValueShort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    float64
		expected string
	}{
		{"integer", 100, "100"},
		{"decimal", 123.45, "123.45"},
		{"small_decimal", 0.123, "0.12"},
		{"zero", 0, "0"},
		{"negative", -50.67, "-50.67"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, style := LabelFormatterValueShort(0, "test", tt.value)
			assert.Equal(t, tt.expected, text)
			assert.Nil(t, style)
		})
	}
}

func TestLabelFormatterNameShortValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		seriesName string
		value      float64
		expected   string
	}{
		{"basic", "Sales", 100, "Sales: 100"},
		{"with_decimal", "Revenue", 123.45, "Revenue: 123.45"},
		{"empty_name", "", 50, ": 50"},
		{"special_chars", "Q1-Sales", 999.99, "Q1-Sales: 999.99"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, style := LabelFormatterNameShortValue(0, tt.seriesName, tt.value)
			assert.Equal(t, tt.expected, text)
			assert.Nil(t, style)
		})
	}
}

func TestInterpolateMultipleColors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		colors   []Color
		factor   float64
		expected Color
	}{
		{
			name:     "no_colors_fallback",
			colors:   []Color{},
			factor:   0.5,
			expected: ColorBlack,
		},
		{
			name:     "single_color",
			colors:   []Color{ColorBlue},
			factor:   0.5,
			expected: ColorBlue,
		},
		{
			name:     "two_colors_start",
			colors:   []Color{ColorRed, ColorBlue},
			factor:   0.0,
			expected: ColorRed,
		},
		{
			name:     "two_colors_end",
			colors:   []Color{ColorRed, ColorBlue},
			factor:   1.0,
			expected: ColorBlue,
		},
		{
			name:     "three_colors_start",
			colors:   []Color{ColorGreen, ColorYellow, ColorRed},
			factor:   0.0,
			expected: ColorGreen,
		},
		{
			name:     "three_colors_middle",
			colors:   []Color{ColorGreen, ColorYellow, ColorRed},
			factor:   0.5,
			expected: ColorYellow,
		},
		{
			name:     "three_colors_end",
			colors:   []Color{ColorGreen, ColorYellow, ColorRed},
			factor:   1.0,
			expected: ColorRed,
		},
		{
			name:     "factor_beyond_range",
			colors:   []Color{ColorGreen, ColorRed},
			factor:   1.5,
			expected: ColorRed, // should clamp to end
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := interpolateMultipleColors(tt.colors, tt.factor)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestInterpolateColor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		color1   Color
		color2   Color
		factor   float64
		expected Color
	}{
		{
			name:     "factor_negative",
			color1:   ColorRed,
			color2:   ColorBlue,
			factor:   -0.5,
			expected: ColorRed, // Should clamp to start color
		},
		{
			name:     "factor_greater_than_one",
			color1:   ColorRed,
			color2:   ColorBlue,
			factor:   1.5,
			expected: ColorBlue, // Should clamp to end color
		},
		{
			name:     "factor_exactly_zero",
			color1:   ColorGreen,
			color2:   ColorYellow,
			factor:   0.0,
			expected: ColorGreen,
		},
		{
			name:     "factor_exactly_one",
			color1:   ColorGreen,
			color2:   ColorYellow,
			factor:   1.0,
			expected: ColorYellow,
		},
		{
			name:     "same_colors",
			color1:   ColorBlue,
			color2:   ColorBlue,
			factor:   0.5,
			expected: ColorBlue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := interpolateColor(tt.color1, tt.color2, tt.factor)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDrawLabelWithBackgroundBorders(t *testing.T) {
	t.Parallel()

	t.Run("border_with_background", func(t *testing.T) {
		p := NewPainter(PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        400,
			Height:       400,
		})

		fontStyle := FontStyle{FontColor: ColorBlack, FontSize: 12}
		drawLabelWithBackground(p, "test", 50, 50, 0, fontStyle, ColorBlue, 0, ColorRed, 2)

		data, err := p.Bytes()
		require.NoError(t, err)
		assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 400\"><path  d=\"M 46 30\nL 81 30\nL 81 54\nL 46 54\nL 46 30\" style=\"stroke-width:2;stroke:red;fill:blue\"/><text x=\"50\" y=\"50\" style=\"stroke:none;fill:black;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">test</text></svg>", data)
	})

	t.Run("border_without_background", func(t *testing.T) {
		p := NewPainter(PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        400,
			Height:       400,
		})

		fontStyle := FontStyle{FontColor: ColorBlack, FontSize: 12}
		drawLabelWithBackground(p, "test", 50, 50, 0, fontStyle, ColorTransparent, 0, ColorGreen, 1.5)

		data, err := p.Bytes()
		require.NoError(t, err)
		assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 400\"><path  d=\"M 46 30\nL 81 30\nL 81 54\nL 46 54\nL 46 30\" style=\"stroke-width:1.5;stroke:green;fill:none\"/><text x=\"50\" y=\"50\" style=\"stroke:none;fill:black;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">test</text></svg>", data)
	})

	t.Run("no_border_no_background", func(t *testing.T) {
		p := NewPainter(PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        400,
			Height:       400,
		})

		fontStyle := FontStyle{FontColor: ColorBlack, FontSize: 12}
		drawLabelWithBackground(p, "test", 50, 50, 0, fontStyle, ColorTransparent, 0, ColorTransparent, 0)

		data, err := p.Bytes()
		require.NoError(t, err)
		assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 400\"><text x=\"50\" y=\"50\" style=\"stroke:none;fill:black;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">test</text></svg>", data)
	})

	t.Run("rounded_corners_with_border", func(t *testing.T) {
		p := NewPainter(PainterOptions{
			OutputFormat: ChartOutputSVG,
			Width:        400,
			Height:       400,
		})

		fontStyle := FontStyle{FontColor: ColorBlack, FontSize: 12}
		drawLabelWithBackground(p, "test", 50, 50, 0, fontStyle, ColorYellow, 8, ColorBlue, 3)

		data, err := p.Bytes()
		require.NoError(t, err)
		assertEqualSVG(t, "<svg xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" viewBox=\"0 0 400 400\"><path  d=\"M 54 30\nL 73 30\nL 73 30\nA 8 8 90.00 0 1 81 38\nL 81 46\nL 81 46\nA 8 8 90.00 0 1 73 54\nL 54 54\nL 54 54\nA 8 8 90.00 0 1 46 46\nL 46 38\nL 46 38\nA 8 8 90.00 0 1 54 30\nZ\" style=\"stroke-width:3;stroke:blue;fill:yellow\"/><text x=\"50\" y=\"50\" style=\"stroke:none;fill:black;font-size:15.3px;font-family:'Roboto Medium',sans-serif\">test</text></svg>", data)
	})
}
