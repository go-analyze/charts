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
			Data: []float64{1.0},
		},
	}, NewSeriesListBar([][]float64{
		{1},
	}))
}

func TestSeriesLists(t *testing.T) {
	t.Parallel()

	seriesList := NewSeriesListBar([][]float64{
		{1, 2},
		{10},
		{1, 2, 3, 4, 5, 6, 7, 8, 9},
	})

	assert.Len(t, seriesList.Filter(ChartTypeBar), 3)
	assert.Empty(t, seriesList.Filter(ChartTypeLine))

	min, max, maxSum := seriesList.getMinMaxSumMax(0, true)
	assert.InDelta(t, float64(12), maxSum, 0)
	assert.InDelta(t, float64(10), max, 0)
	assert.InDelta(t, float64(1), min, 0)
}

func TestSumSeries(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		var sl SeriesList
		result := sl.SumSeries()

		assert.Equal(t, "", result.Type)
		assert.Empty(t, result.Data)
	})

	t.Run("single", func(t *testing.T) {
		sl := SeriesList{
			{
				Type:       ChartTypeLine,
				Data:       []float64{1.5, 2.5},
				Name:       "SingleLine",
				YAxisIndex: 1,
				Radius:     "50%",
			},
		}

		result := sl.SumSeries()

		assert.Equal(t, sl[0], result)
	})

	t.Run("multiple", func(t *testing.T) {
		sl := NewSeriesListLine([][]float64{
			{1, 2, 3},
			{4, 5, 6},
		})

		result := sl.SumSeries()

		assert.Equal(t, ChartTypeLine, result.Type)
		assert.Equal(t, []float64{5, 7, 9}, result.Data)
	})

	t.Run("unequal_data_length", func(t *testing.T) {
		sl := NewSeriesListLine([][]float64{
			{1, 2},
			{3, 4, 5},
		})

		result := sl.SumSeries()

		assert.Equal(t, ChartTypeLine, result.Type)
		assert.Equal(t, []float64{4, 6, 5}, result.Data)
	})
}

func TestSeriesSummary(t *testing.T) {
	t.Parallel()

	seriesList := NewSeriesListLine([][]float64{
		{10},
		{1, 2},
		{1, 2, 3},
		{1, 2, 3, 4},
		{3, 7, 11, 13},
		{GetNullValue()},
		{10, GetNullValue()},
		{1, GetNullValue(), 2},
	})

	t.Run("empty_series", func(t *testing.T) {
		assert.Equal(t, populationSummary{
			MaxIndex: -1,
			MinIndex: -1,
		}, (&Series{}).Summary())
	})
	t.Run("one_value", func(t *testing.T) {
		assert.Equal(t, populationSummary{
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
		assert.Equal(t, populationSummary{
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
		assert.Equal(t, populationSummary{
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
		assert.Equal(t, populationSummary{
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
		assert.Equal(t, populationSummary{
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
	t.Run("null_only", func(t *testing.T) {
		assert.Equal(t, populationSummary{
			MaxIndex: -1,
			MinIndex: -1,
		}, seriesList[5].Summary())
	})
	t.Run("value_null", func(t *testing.T) {
		assert.Equal(t, populationSummary{
			Max:               10,
			MaxIndex:          0,
			Min:               10,
			MinIndex:          0,
			Average:           10,
			Median:            10,
			StandardDeviation: 0.0,
			Skewness:          0.0,
			Kurtosis:          0.0,
		}, seriesList[6].Summary())
	})
	t.Run("value_null_value", func(t *testing.T) {
		assert.Equal(t, populationSummary{
			Max:               2,
			MaxIndex:          2,
			Min:               1,
			MinIndex:          0,
			Average:           1.5,
			Median:            1.5,
			StandardDeviation: 0.5,
			Skewness:          0.0,
			Kurtosis:          1.0,
		}, seriesList[7].Summary())
	})
}

func TestFormatter(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "a: 12%",
		labelFormatPie([]string{"a", "b"}, "", 0, 10, 0.12))

	assert.Equal(t, "10",
		labelFormatValue([]string{"a", "b"}, "", 0, 10, 0.12))
}
