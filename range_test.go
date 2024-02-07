package charts

import (
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
			name:             "PadMaxOnly",
			expectedMinValue: 0.0,
			expectedMaxValue: 10.5,
			minValue:         0.0,
			maxValue:         10.0,
			labelCount:       10,
		},
		{
			name:             "PadMinToZero",
			expectedMinValue: 0.0,
			expectedMaxValue: 21.0,
			minValue:         1.0,
			maxValue:         20.0,
			labelCount:       10,
		},
		{
			name:             "PadNegativeMinPositiveMax",
			expectedMinValue: -5.0,
			expectedMaxValue: 12.0,
			minValue:         -3.0,
			maxValue:         10.0,
			labelCount:       10,
		},
		{
			name:             "PadNegativeMinNegativeMax",
			expectedMinValue: -20.0,
			expectedMaxValue: -9.0,
			minValue:         -20.0,
			maxValue:         -10.0,
			labelCount:       10,
		},
		{
			name:             "PadPositiveMinPositiveMax",
			expectedMinValue: 100.0,
			expectedMaxValue: 214.0,
			minValue:         100.0,
			maxValue:         200.0,
			labelCount:       20,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			min, max := padRange(tc.labelCount, tc.minValue, tc.maxValue, 1.0, 1.0)

			assert.Equal(t, tc.expectedMinValue, min, "Unexpected value rounding %v", tc.minValue)
			assert.Equal(t, tc.expectedMaxValue, max, "Unexpected value rounding %v", tc.maxValue)
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
			name:          "OriginalZeroSub",
			expectedValue: 0.0,
			value:         0.0,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           false,
		},
		{
			name:          "OriginalZeroAdd",
			expectedValue: 0.0,
			value:         0.0,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           true,
		},
		{
			name:          "RoundFractionSub",
			expectedValue: -2.0,
			value:         -1.2,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           false,
		},
		{
			name:          "RoundFractionAdd",
			expectedValue: 2.0,
			value:         1.2,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           true,
		},
		{
			name:          "RoundVeryCloseToZeroSub",
			expectedValue: -1.0,
			value:         -0.01,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           false,
		},
		{
			name:          "RoundVeryCloseToZeroAdd",
			expectedValue: 0.0,
			value:         -0.01,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           true,
		},
		{
			name:          "OriginalLargeSub",
			expectedValue: 1337,
			value:         1337,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           false,
		},
		{
			name:          "OriginalLargeAdd",
			expectedValue: 1337,
			value:         1337,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           true,
		},
		{
			name:          "RoundThousandLargeSub",
			expectedValue: 1000,
			value:         1337,
			minMultiplier: 0.0,
			maxMultiplier: 1000.0,
			add:           false,
		},
		{
			name:          "RoundThousandLargeAdd",
			expectedValue: 2000,
			value:         1337,
			minMultiplier: 0.0,
			maxMultiplier: 1000.0,
			add:           true,
		},
		{
			name:          "RoundHundredLargeSub",
			expectedValue: 1300,
			value:         1337,
			minMultiplier: 0.0,
			maxMultiplier: 100.0,
			add:           false,
		},
		{
			name:          "RoundHundredLargeAdd",
			expectedValue: 1400,
			value:         1337,
			minMultiplier: 0.0,
			maxMultiplier: 100.0,
			add:           true,
		},
		{
			name:          "RoundNegativeSmallSub",
			expectedValue: -1.0,
			value:         -0.5,
			minMultiplier: 0.0,
			maxMultiplier: 2.0,
			add:           false,
		},
		{
			name:          "RoundHalfwayPointSub",
			expectedValue: 100.0,
			value:         150.0,
			minMultiplier: 0.0,
			maxMultiplier: 100.0,
			add:           false,
		},
		{
			name:          "RoundHalfwayPointAdd",
			expectedValue: 200.0,
			value:         150.0,
			minMultiplier: 0.0,
			maxMultiplier: 100.0,
			add:           true,
		},
		{
			name:          "RoundThousandsNegativeLargeSub",
			expectedValue: -2000.0,
			value:         -1337.0,
			minMultiplier: 0.0,
			maxMultiplier: 1000.0,
			add:           false,
		},
		{
			name:          "RoundHundredsNegativeLargeSub",
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

			assert.Equal(t, tc.expectedValue, val, "Unexpected value rounding %v", tc.value)
		})
	}
}
