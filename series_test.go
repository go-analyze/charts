package charts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSeriesListDataFromValues(t *testing.T) {
	t.Parallel()

	assert.Equal(t, SeriesList{
		{
			Type: ChartTypeBar,
			Data: []float64{
				1.0,
			},
		},
	}, NewSeriesListDataFromValues([][]float64{
		{
			1,
		},
	}, ChartTypeBar))
}

func TestSeriesLists(t *testing.T) {
	t.Parallel()

	seriesList := NewSeriesListDataFromValues([][]float64{
		{1, 2},
		{10},
	}, ChartTypeBar)

	assert.Equal(t, 2, len(seriesList.Filter(ChartTypeBar)))
	assert.Equal(t, 0, len(seriesList.Filter(ChartTypeLine)))

	min, max := seriesList.GetMinMax(0)
	assert.Equal(t, float64(10), max)
	assert.Equal(t, float64(1), min)
}

func TestSeriesSummary(t *testing.T) {
	t.Parallel()

	seriesList := NewSeriesListDataFromValues([][]float64{
		{10},
		{1, 2},
		{1, 2, 3},
		{1, 2, 3, 4},
		{3, 7, 11, 13},
	}, ChartTypeLine)

	t.Run("empty_series", func(t *testing.T) {
		assert.Equal(t, seriesSummary{
			MaxIndex: -1,
			MinIndex: -1,
		}, (&Series{}).Summary())
	})
	t.Run("one_value", func(t *testing.T) {
		assert.Equal(t, seriesSummary{
			Max:               10,
			MaxIndex:          0,
			Min:               10,
			MinIndex:          0,
			Average:           10,
			Median:            10,
			StandardDeviation: 0.0,
			Skewness:          0.0,
			Kurtosis:          0.0,
		}, seriesList[0].Summary())
	})
	t.Run("two_values", func(t *testing.T) {
		assert.Equal(t, seriesSummary{
			Max:               2,
			MaxIndex:          1,
			Min:               1,
			MinIndex:          0,
			Average:           1.5,
			Median:            1.5,
			StandardDeviation: 0.5,
			Skewness:          0.0,
			Kurtosis:          1.0,
		}, seriesList[1].Summary())
	})
	t.Run("three_values", func(t *testing.T) {
		assert.Equal(t, seriesSummary{
			Max:               3,
			MaxIndex:          2,
			Min:               1,
			MinIndex:          0,
			Average:           2,
			Median:            2,
			StandardDeviation: 0.8164965809277263,
			Skewness:          0.0,
			Kurtosis:          1.4999999999999987,
		}, seriesList[2].Summary())
	})
	t.Run("four_values", func(t *testing.T) {
		assert.Equal(t, seriesSummary{
			Max:               4,
			MaxIndex:          3,
			Min:               1,
			MinIndex:          0,
			Average:           2.5,
			Median:            2.5,
			StandardDeviation: 1.118033988749895,
			Skewness:          0.0,
			Kurtosis:          1.64,
		}, seriesList[3].Summary())
	})
	t.Run("prime_values", func(t *testing.T) {
		assert.Equal(t, seriesSummary{
			Max:               13,
			MaxIndex:          3,
			Min:               3,
			MinIndex:          0,
			Average:           8.5,
			Median:            9,
			StandardDeviation: 3.840572873934304,
			Skewness:          -0.2780305556539629,
			Kurtosis:          1.5733984487216317,
		}, seriesList[4].Summary())
	})
}

func TestFormatter(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "a: 12%", labelFormatPie([]string{
		"a",
		"b",
	}, "", 0, 10, 0.12))

	assert.Equal(t, "10", labelFormatValue([]string{
		"a",
		"b",
	}, "", 0, 10, 0.12))
}
