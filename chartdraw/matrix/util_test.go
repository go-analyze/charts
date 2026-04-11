package matrix

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestF64s(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		value  float64
		expect string
	}{
		{"Zero", 0, "0"},
		{"Positive", 1.23, "1.23"},
		{"Integer", 3.0, "3"},
		{"Negative", -0.5, "-0.5"},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, f64s(tc.value))
		})
	}
}

func TestRoundToEpsilon(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   float64
		epsilon float64
	}{
		{"Zero", 0, 0},
		{"Positive", 1.234, 0.1},
		{"Negative", -2.5, 0.001},
	}

	for i, tc := range tests {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			r := roundToEpsilon(tc.value, tc.epsilon)
			assert.InDelta(t, tc.value, r, 0)
		})
	}
}
