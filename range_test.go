package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPadRange(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		expectedMinValue float64
		expectedMaxValue float64
		minValue         float64
		maxValue         float64
		labelCount       int
	}{
		{
			name:             "pad_max_only",
			expectedMinValue: 0.0,
			expectedMaxValue: 10.5,
			minValue:         0.0,
			maxValue:         10.0,
			labelCount:       10,
		},
		{
			name:             "pad_min_to_zero",
			expectedMinValue: 0.0,
			expectedMaxValue: 21.0,
			minValue:         1.0,
			maxValue:         20.0,
			labelCount:       10,
		},
		{
			name:             "pad_negative_min_positive_max",
			expectedMinValue: -5.0,
			expectedMaxValue: 12.0,
			minValue:         -3.0,
			maxValue:         10.0,
			labelCount:       10,
		},
		{
			name:             "pad_negative_min_negative_max",
			expectedMinValue: -20.0,
			expectedMaxValue: -9.0,
			minValue:         -20.0,
			maxValue:         -10.0,
			labelCount:       10,
		},
		{
			name:             "pad_positive_min_positive_max",
			expectedMinValue: 100.0,
			expectedMaxValue: 214.0,
			minValue:         100.0,
			maxValue:         200.0,
			labelCount:       20,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			min, max := padRange(tc.labelCount, tc.minValue, tc.maxValue, 1.0, 1.0)

			assert.InDelta(t, tc.expectedMinValue, min, 0, "Unexpected value rounding %v", tc.minValue)
			assert.InDelta(t, tc.expectedMaxValue, max, 0, "Unexpected value rounding %v", tc.maxValue)
		})
	}
}

func TestFriendlyRound(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		expectedValue float64
		value         float64
		minMultiplier float64
		maxMultiplier float64
		add           bool
	}{
		{
			name:          "original_zero_sub",
			expectedValue: 0.0,
			value:         0.0,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           false,
		},
		{
			name:          "original_zero_add",
			expectedValue: 0.0,
			value:         0.0,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           true,
		},
		{
			name:          "round_fraction_sub",
			expectedValue: -2.0,
			value:         -1.2,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           false,
		},
		{
			name:          "round_fraction_add",
			expectedValue: 2.0,
			value:         1.2,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           true,
		},
		{
			name:          "round_very_close_to_zero_sub",
			expectedValue: -1.0,
			value:         -0.01,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           false,
		},
		{
			name:          "round_very_close_to_zero_add",
			expectedValue: 0.0,
			value:         -0.01,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           true,
		},
		{
			name:          "original_large_sub",
			expectedValue: 1337,
			value:         1337,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           false,
		},
		{
			name:          "original_large_add",
			expectedValue: 1337,
			value:         1337,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           true,
		},
		{
			name:          "round_thousand_large_sub",
			expectedValue: 1000,
			value:         1337,
			minMultiplier: 0.0,
			maxMultiplier: 1000.0,
			add:           false,
		},
		{
			name:          "round_thousand_large_add",
			expectedValue: 2000,
			value:         1337,
			minMultiplier: 0.0,
			maxMultiplier: 1000.0,
			add:           true,
		},
		{
			name:          "round_hundred_large_sub",
			expectedValue: 1300,
			value:         1337,
			minMultiplier: 0.0,
			maxMultiplier: 100.0,
			add:           false,
		},
		{
			name:          "round_hundred_large_add",
			expectedValue: 1400,
			value:         1337,
			minMultiplier: 0.0,
			maxMultiplier: 100.0,
			add:           true,
		},
		{
			name:          "round_negative_small_sub",
			expectedValue: -1.0,
			value:         -0.5,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           false,
		},
		{
			name:          "round_halfway_point_sub",
			expectedValue: 100.0,
			value:         150.0,
			minMultiplier: 0.0,
			maxMultiplier: 100.0,
			add:           false,
		},
		{
			name:          "round_halfway_point_add",
			expectedValue: 200.0,
			value:         150.0,
			minMultiplier: 0.0,
			maxMultiplier: 100.0,
			add:           true,
		},
		{
			name:          "round_thousands_negative_large_sub",
			expectedValue: -2000.0,
			value:         -1337.0,
			minMultiplier: 0.0,
			maxMultiplier: 1000.0,
			add:           false,
		},
		{
			name:          "round_hundreds_negative_large_sub",
			expectedValue: -1400.0,
			value:         -1337.0,
			minMultiplier: 0.0,
			maxMultiplier: 100.0,
			add:           false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, _ := friendlyRound(tc.value, 1.0, 0.0,
				tc.minMultiplier, tc.maxMultiplier, tc.add)

			assert.InDelta(t, tc.expectedValue, val, 0, "Unexpected value rounding %v", tc.value)
		})
	}
}
