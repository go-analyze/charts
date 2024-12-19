package chartdraw

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeValueFormatterWithFormat(t *testing.T) {
	t.Parallel()

	d := time.Now()
	df := TimeToFloat64(d)

	s := formatTime(d, DefaultDateFormat)
	sf := formatTime(df, DefaultDateFormat)
	assert.Equal(t, s, sf)

	sd := TimeValueFormatter(d)
	sdf := TimeValueFormatter(df)
	assert.Equal(t, s, sd)
	assert.Equal(t, s, sdf)
}

func TestFloatValueFormatter(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		input interface{}
	}{
		{
			name:  "basic_float",
			input: 1234.00,
		},
		{
			name:  "float32",
			input: float32(1234.00),
		},
		{
			name:  "int",
			input: 1234,
		},
		{
			name:  "int64",
			input: int64(1234),
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i)+"-"+tc.name, func(t *testing.T) {
			assert.Equal(t, "1234.00", FloatValueFormatter(tc.input))
		})
	}
}

func TestFloatValueFormatterWithFormat(t *testing.T) {
	t.Parallel()

	v := 123.456
	sv := FloatValueFormatterWithFormat(v, "%.3f")
	assert.Equal(t, "123.456", sv)
	assert.Equal(t, "123.000", FloatValueFormatterWithFormat(123, "%.3f"))
}

func TestExponentialValueFormatter(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "1.23e+02", ExponentialValueFormatter(123.456))
	assert.Equal(t, "1.24e+07", ExponentialValueFormatter(12421243.424))
	assert.Equal(t, "4.50e-01", ExponentialValueFormatter(0.45))
}
