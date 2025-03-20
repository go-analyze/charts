package charts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLabelFormatPie(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "a: 12%",
		labelFormatPie([]string{"a", "b"}, "", nil, 0, 10, 0.12))

	assert.Equal(t, "b: 25%",
		labelFormatPie([]string{"a", "b"}, "", nil, 1, 20, 0.25))

	assert.Equal(t, "a: f",
		labelFormatPie([]string{"a", "b"}, "{b}: {c}", func(f float64) string {
			return "f"
		}, 0, 10, 0.12))
}

func TestLabelFormatFunnel(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "a(12%)",
		labelFormatFunnel([]string{"a", "b"}, "", nil, 0, 10, 0.12))

	assert.Equal(t, "b(25%)",
		labelFormatFunnel([]string{"a", "b"}, "", nil, 1, 20, 0.25))

	assert.Equal(t, "b(f, 25%)",
		labelFormatFunnel([]string{"a", "b"}, "{b}({c}, {d})", func(f float64) string {
			return "f"
		}, 1, 20, 0.25))
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
		labelFormatPie([]string{"a", "b"}, "Name: {b}, Value: {c}, Percent: {d}", nil,
			0, 10, 0.12))

	assert.Equal(t, "Name: b, Value: 20, Percent: 25%",
		labelFormatPie([]string{"a", "b"}, "Name: {b}, Value: {c}, Percent: {d}", nil,
			1, 20, 0.25))

	assert.Equal(t, "Empty Series '' 20",
		labelFormatPie([]string{}, "Empty Series '{b}' {c}", nil,
			1, 20, 0.25))
}
