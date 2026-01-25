package charts

import (
	"testing"

	"github.com/golang/freetype/truetype"
	"github.com/stretchr/testify/assert"

	"github.com/go-analyze/charts/chartdraw/matrix"
)

func TestFillFontStyleDefaults(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		input        FontStyle
		defaultSize  float64
		defaultColor Color
		fontOptions  []*truetype.Font
		expected     FontStyle
	}{
		{
			name:         "empty_font_style_gets_all_defaults",
			input:        FontStyle{},
			defaultSize:  12,
			defaultColor: ColorRed,
			expected: FontStyle{
				FontSize:  12,
				FontColor: ColorRed,
				Font:      getPreferredFont(),
			},
		},
		{
			name: "partial_font_style_keeps_existing_values",
			input: FontStyle{
				FontSize:  16,
				FontColor: ColorBlue,
			},
			defaultSize:  12,
			defaultColor: ColorRed,
			expected: FontStyle{
				FontSize:  16,        // keeps existing
				FontColor: ColorBlue, // keeps existing
				Font:      getPreferredFont(),
			},
		},
		{
			name: "zero_font_size_gets_default",
			input: FontStyle{
				FontSize:  0,
				FontColor: ColorGreen,
			},
			defaultSize:  14,
			defaultColor: ColorRed,
			expected: FontStyle{
				FontSize:  14,         // uses default
				FontColor: ColorGreen, // keeps existing
				Font:      getPreferredFont(),
			},
		},
		{
			name: "zero_color_gets_default",
			input: FontStyle{
				FontSize:  18,
				FontColor: Color{}, // zero color
			},
			defaultSize:  12,
			defaultColor: ColorYellow,
			expected: FontStyle{
				FontSize:  18,          // keeps existing
				FontColor: ColorYellow, // uses default
				Font:      getPreferredFont(),
			},
		},
		{
			name: "existing_font_is_preserved",
			input: FontStyle{
				Font: &truetype.Font{}, // some font
			},
			defaultSize:  12,
			defaultColor: ColorRed,
			expected: FontStyle{
				FontSize:  12,
				FontColor: ColorRed,
				Font:      &truetype.Font{}, // keeps existing font
			},
		},
		{
			name: "all_values_set_no_defaults_applied",
			input: FontStyle{
				FontSize:  20,
				FontColor: ColorPurple,
				Font:      &truetype.Font{},
			},
			defaultSize:  12,
			defaultColor: ColorRed,
			expected: FontStyle{
				FontSize:  20,               // keeps existing
				FontColor: ColorPurple,      // keeps existing
				Font:      &truetype.Font{}, // keeps existing
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fillFontStyleDefaults(tt.input, tt.defaultSize, tt.defaultColor, tt.fontOptions...)
			assert.InDelta(t, tt.expected.FontSize, result.FontSize, matrix.DefaultEpsilon)
			assert.Equal(t, tt.expected.FontColor, result.FontColor)
			if tt.expected.Font != nil && result.Font != nil {
				// Both have fonts, compare if they're the same instance
				if tt.input.Font != nil {
					// If input had a font, it should be preserved
					assert.Equal(t, tt.input.Font, result.Font)
				} else {
					// If input had no font, should get default font
					assert.NotNil(t, result.Font)
				}
			} else if tt.expected.Font == nil && result.Font == nil {
				// Both should be nil
			} else {
				t.Errorf("Font mismatch: expected %v, got %v", tt.expected.Font, result.Font)
			}
		})
	}
}

func TestMergeFontStyles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		primary  FontStyle
		defaults []FontStyle
		expected FontStyle
	}{
		{
			name:    "empty_primary_uses_first_available_defaults",
			primary: FontStyle{},
			defaults: []FontStyle{
				{FontSize: 12, FontColor: ColorRed, Font: &truetype.Font{}},
			},
			expected: FontStyle{
				FontSize:  12,
				FontColor: ColorRed,
				Font:      &truetype.Font{},
			},
		},
		{
			name: "primary_values_take_precedence",
			primary: FontStyle{
				FontSize:  16,
				FontColor: ColorBlue,
			},
			defaults: []FontStyle{
				{FontSize: 12, FontColor: ColorRed, Font: &truetype.Font{}},
			},
			expected: FontStyle{
				FontSize:  16,               // from primary
				FontColor: ColorBlue,        // from primary
				Font:      &truetype.Font{}, // from default
			},
		},
		{
			name:    "uses_first_available_value_from_multiple_defaults",
			primary: FontStyle{},
			defaults: []FontStyle{
				{},                                    // empty first default
				{FontSize: 14},                        // has size
				{FontSize: 18, FontColor: ColorGreen}, // has size and color
			},
			expected: FontStyle{
				FontSize:  14,         // from second default (first with size)
				FontColor: ColorGreen, // from third default (first with color)
				Font:      nil,        // no defaults have font
			},
		},
		{
			name: "partial_primary_fills_from_defaults",
			primary: FontStyle{
				FontSize: 20, // has size
				// missing color and font
			},
			defaults: []FontStyle{
				{FontColor: ColorPurple}, // has color
				{Font: &truetype.Font{}}, // has font
			},
			expected: FontStyle{
				FontSize:  20,               // from primary
				FontColor: ColorPurple,      // from first default
				Font:      &truetype.Font{}, // from second default
			},
		},
		{
			name: "zero_values_in_primary_use_defaults",
			primary: FontStyle{
				FontSize:  0,       // zero value
				FontColor: Color{}, // zero color
				Font:      nil,     // nil font
			},
			defaults: []FontStyle{
				{FontSize: 15, FontColor: ColorYellow, Font: &truetype.Font{}},
			},
			expected: FontStyle{
				FontSize:  15,
				FontColor: ColorYellow,
				Font:      &truetype.Font{},
			},
		},
		{
			name: "no_defaults_primary_unchanged",
			primary: FontStyle{
				FontSize:  24,
				FontColor: ColorBlack,
			},
			defaults: []FontStyle{},
			expected: FontStyle{
				FontSize:  24,
				FontColor: ColorBlack,
				Font:      nil,
			},
		},
		{
			name: "complex_precedence_scenario",
			primary: FontStyle{
				FontColor: ColorRed, // only color set
			},
			defaults: []FontStyle{
				{FontSize: 10},                         // first default has size
				{},                                     // empty default
				{FontSize: 12, Font: &truetype.Font{}}, // third has size and font
				{FontSize: 14, FontColor: ColorBlue, Font: &truetype.Font{}}, // fourth has all
			},
			expected: FontStyle{
				FontSize:  10,               // from first default with size
				FontColor: ColorRed,         // from primary
				Font:      &truetype.Font{}, // from third default with font
			},
		},
		{
			name: "all_defaults_empty_primary_unchanged",
			primary: FontStyle{
				FontSize:  18,
				FontColor: ColorGreen,
				Font:      &truetype.Font{},
			},
			defaults: []FontStyle{
				{}, {}, {}, // all empty
			},
			expected: FontStyle{
				FontSize:  18,
				FontColor: ColorGreen,
				Font:      &truetype.Font{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mergeFontStyles(tt.primary, tt.defaults...)
			assert.InDelta(t, tt.expected.FontSize, result.FontSize, matrix.DefaultEpsilon)
			assert.Equal(t, tt.expected.FontColor, result.FontColor)

			// Handle font comparison carefully since their pointers
			if tt.expected.Font != nil && result.Font != nil {
				// Both should be non-nil, check if they match expected source
				if tt.primary.Font != nil {
					assert.Equal(t, tt.primary.Font, result.Font, "should preserve primary font")
				} else {
					// Should come from defaults, just verify it's not nil
					assert.NotNil(t, result.Font, "should get font from defaults")
				}
			} else if tt.expected.Font == nil {
				assert.Nil(t, result.Font, "font should be nil")
			}
		})
	}
}

func TestMergeFontStylesEdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("nil_font_handling", func(t *testing.T) {
		primary := FontStyle{}
		defaults := []FontStyle{
			{Font: nil},              // explicit nil
			{Font: &truetype.Font{}}, // has font
		}

		result := mergeFontStyles(primary, defaults...)
		assert.NotNil(t, result.Font, "should get font from second default")
	})

	t.Run("zero_vs_non_zero_color", func(t *testing.T) {
		primary := FontStyle{
			FontColor: Color{R: 0, G: 0, B: 0, A: 0}, // transparent/zero
		}
		defaults := []FontStyle{
			{FontColor: ColorRed},
		}

		result := mergeFontStyles(primary, defaults...)
		assert.Equal(t, ColorRed, result.FontColor, "zero color should be replaced")
	})

	t.Run("partial_alpha_color", func(t *testing.T) {
		partialColor := Color{R: 100, G: 100, B: 100, A: 0} // has RGB but no alpha
		primary := FontStyle{
			FontColor: partialColor,
		}
		defaults := []FontStyle{
			{FontColor: ColorBlue},
		}

		result := mergeFontStyles(primary, defaults...)
		// Should keep the primary color even with zero alpha since IsZero() checks all components
		if partialColor.IsZero() {
			assert.Equal(t, ColorBlue, result.FontColor)
		} else {
			assert.Equal(t, partialColor, result.FontColor)
		}
	})
}
