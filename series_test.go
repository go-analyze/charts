package charts

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSeriesLists(t *testing.T) {
	t.Parallel()

	values := [][]float64{
		{1, 2},
		{10},
		{1, 2, 3, 4, 5, 6, 7, 8, 9},
	}
	for i, tc := range []string{ChartTypeLine, ChartTypeBar, ChartTypeHorizontalBar} {
		t.Run(strconv.Itoa(i)+"-"+tc, func(t *testing.T) {
			var seriesList seriesList
			switch tc { // switch case to ensure chart type and generic type match expectations
			case ChartTypeLine:
				seriesList = NewSeriesListLine(values)
				assert.Len(t, filterSeriesList[LineSeriesList](seriesList, ChartTypeLine), 3)
				assert.Empty(t, filterSeriesList[BarSeriesList](seriesList, ChartTypeBar))
			case ChartTypeBar:
				seriesList = NewSeriesListBar(values)
				assert.Len(t, filterSeriesList[BarSeriesList](seriesList, ChartTypeBar), 3)
				assert.Empty(t, filterSeriesList[LineSeriesList](seriesList, ChartTypeLine))
			case ChartTypeHorizontalBar:
				seriesList = NewSeriesListHorizontalBar(values)
				assert.Len(t, filterSeriesList[HorizontalBarSeriesList](seriesList, ChartTypeHorizontalBar), 3)
				assert.Empty(t, filterSeriesList[LineSeriesList](seriesList, ChartTypeLine))
			default:
				require.Fail(t, "Need to implement chart type test")
			}

			min, max, maxSum := getSeriesMinMaxSumMax(seriesList, 0, true)
			assert.InDelta(t, float64(12), maxSum, 0)
			assert.InDelta(t, float64(10), max, 0)
			assert.InDelta(t, float64(1), min, 0)
		})
	}
}

func TestSumSeries(t *testing.T) {
	t.Parallel()

	type summableSeries interface {
		SumSeries() []float64
	}
	testTypes := []struct {
		name       string
		seriesFact func([][]float64) summableSeries
	}{
		{
			name: "line",
			seriesFact: func(values [][]float64) summableSeries {
				return NewSeriesListLine(values)
			},
		},
		{
			name: "bar",
			seriesFact: func(values [][]float64) summableSeries {
				return NewSeriesListBar(values)
			},
		},
		{
			name: "horizontal_bar",
			seriesFact: func(values [][]float64) summableSeries {
				return NewSeriesListHorizontalBar(values)
			},
		},
	}
	tests := []struct {
		name     string
		values   [][]float64
		expected []float64
	}{
		{
			name:     "empty",
			values:   [][]float64{},
			expected: []float64{},
		},
		{
			name:     "single",
			values:   [][]float64{{1.5, 2.5}},
			expected: []float64{1.5, 2.5},
		},
		{
			name: "multiple",
			values: [][]float64{
				{1, 2, 3},
				{4, 5, 6},
			},
			expected: []float64{5, 7, 9},
		},
		{
			name: "unequal_data_length",
			values: [][]float64{
				{1, 2},
				{3, 4, 5},
			},
			expected: []float64{4, 6, 5},
		},
		{
			name: "null_values",
			values: [][]float64{
				{GetNullValue(), 2, 3},
				{4, GetNullValue(), 6},
			},
			expected: []float64{4, 2, 9},
		},
	}

	for _, typeCase := range testTypes {
		for _, tc := range tests {
			t.Run(typeCase.name+"-"+tc.name, func(t *testing.T) {
				series := typeCase.seriesFact(tc.values)
				result := series.SumSeries()

				assert.Equal(t, tc.expected, result)
			})
		}
	}
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
		}, summarizePopulationData(nil))
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

func BenchmarkGetSeriesYAxisCount(b *testing.B) { // benchmark used to evaluate methods for iterating the series
	nameCount := 100
	seriesList := make(LineSeriesList, nameCount)
	for i := 0; i < nameCount; i++ {
		seriesList[i] = LineSeries{}
	}

	for i := 0; i < b.N; i++ {
		_ = getSeriesYAxisCount(seriesList)
	}
}

func BenchmarkGetSeriesMinMaxSumMax(b *testing.B) { // benchmark used to evaluate methods for iterating the series
	seriesCount := 100
	seriesSize := 100
	seriesList := make(LineSeriesList, seriesCount)
	for i := 0; i < seriesCount; i++ {
		data := make([]float64, seriesSize)
		for si := 0; si < seriesSize; si++ {
			if si+1%10 == 0 {
				data[si] = GetNullValue()
			} else {
				data[si] = float64(si)
			}
		}
		seriesList[i] = LineSeries{
			Values: data,
		}
	}

	for i := 0; i < b.N; i++ {
		_, _, _ = getSeriesMinMaxSumMax(seriesList, 0, true)
	}
}

func BenchmarkSumSeries(b *testing.B) { // benchmark used to evaluate methods for iterating the series
	seriesCount := 100
	seriesSize := 100
	seriesList := make(LineSeriesList, seriesCount)
	for i := 0; i < seriesCount; i++ {
		seriesList[i] = LineSeries{
			Values: make([]float64, seriesSize),
		}
	}

	for i := 0; i < b.N; i++ {
		_ = seriesList.SumSeries()
	}
}

func BenchmarkSeriesNames(b *testing.B) { // benchmark used to evaluate methods for iterating the series
	nameCount := 100
	seriesList := make(LineSeriesList, nameCount)
	for i := 0; i < nameCount; i++ {
		seriesList[i] = LineSeries{
			Name: strconv.Itoa(i),
		}
	}

	for i := 0; i < b.N; i++ {
		_ = seriesList.names()
	}
}

func BenchmarkGetSeriesMaxDataCount(b *testing.B) { // benchmark used to evaluate methods for iterating the series
	seriesCount := 100
	seriesList := make(LineSeriesList, seriesCount)
	for i := 0; i < seriesCount; i++ {
		seriesList[i] = LineSeries{
			Values: make([]float64, i),
		}
	}

	for i := 0; i < b.N; i++ {
		_ = getSeriesMaxDataCount(seriesList)
	}
}
