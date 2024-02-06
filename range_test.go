package charts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
