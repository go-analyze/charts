package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValuesValues(t *testing.T) {
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
	assert.Len(t, values, 7)
	assert.Equal(t, float64(10), values[0])
	assert.Equal(t, float64(9), values[1])
	assert.Equal(t, float64(8), values[2])
	assert.Equal(t, float64(7), values[3])
	assert.Equal(t, float64(6), values[4])
	assert.Equal(t, float64(5), values[5])
	assert.Equal(t, float64(2), values[6])
}

func TestValuesValuesNormalized(t *testing.T) {
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
	assert.Len(t, values, 7)
	assert.Equal(t, 0.2127, values[0])
	assert.Equal(t, 0.0425, values[6])
}

func TestValuesNormalize(t *testing.T) {
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
	assert.Len(t, values, 7)
	assert.Equal(t, 0.2127, values[0].Value)
	assert.Equal(t, 0.0425, values[6].Value)
}
