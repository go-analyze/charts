package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValuesValues(t *testing.T) {
	t.Parallel()

	vs := []Value{
		{Value: 10, Label: "Blue"},
		{Value: 9, Label: "Green"},
		{Value: 8, Label: "Gray"},
		{Value: 7, Label: "Orange"},
		{Value: 6, Label: "HEANG"},
		{Value: 5, Label: "??"},
		{Value: 2, Label: "!!"},
	}

	values := Values(vs).Values()
	require.Len(t, values, 7)
	assert.InDelta(t, float64(10), values[0], 0)
	assert.InDelta(t, float64(9), values[1], 0)
	assert.InDelta(t, float64(8), values[2], 0)
	assert.InDelta(t, float64(7), values[3], 0)
	assert.InDelta(t, float64(6), values[4], 0)
	assert.InDelta(t, float64(5), values[5], 0)
	assert.InDelta(t, float64(2), values[6], 0)
}

func TestValuesValuesNormalized(t *testing.T) {
	t.Parallel()

	vs := []Value{
		{Value: 10, Label: "Blue"},
		{Value: 9, Label: "Green"},
		{Value: 8, Label: "Gray"},
		{Value: 7, Label: "Orange"},
		{Value: 6, Label: "HEANG"},
		{Value: 5, Label: "??"},
		{Value: 2, Label: "!!"},
	}

	values := Values(vs).ValuesNormalized()
	require.Len(t, values, 7)
	assert.InDelta(t, 0.2127, values[0], 0)
	assert.InDelta(t, 0.0425, values[6], 0)
}

func TestValuesNormalize(t *testing.T) {
	t.Parallel()

	vs := []Value{
		{Value: 10, Label: "Blue"},
		{Value: 9, Label: "Green"},
		{Value: 8, Label: "Gray"},
		{Value: 7, Label: "Orange"},
		{Value: 6, Label: "HEANG"},
		{Value: 5, Label: "??"},
		{Value: 2, Label: "!!"},
	}

	values := Values(vs).Normalize()
	require.Len(t, values, 7)
	assert.InDelta(t, 0.2127, values[0].Value, 0)
	assert.InDelta(t, 0.0425, values[6].Value, 0)
}
