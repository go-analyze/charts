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

func TestFilterSeriesListCandlestick(t *testing.T) {
	t.Parallel()

	generic := GenericSeriesList{
		{Values: []float64{1, 2}, Type: ChartTypeCandlestick},
		{Values: []float64{3, 4}, Type: ChartTypeLine},
	}

	filtered := filterSeriesList[CandlestickSeriesList](generic, ChartTypeCandlestick)
	assert.Len(t, filtered, 1)
	require.Len(t, filtered[0].Data, 2)
	assert.InDelta(t, 1.0, filtered[0].Data[0].Open, 0)
	assert.InDelta(t, 2.0, filtered[0].Data[1].Close, 0)
}

func TestGetSeriesMinMaxSumMaxEmpty(t *testing.T) {
	t.Parallel()

	empty := NewSeriesListLine([][]float64{{}})
	min, max, sum := getSeriesMinMaxSumMax(empty, 0, true)
	assert.InDelta(t, 0.0, min, 0)
	assert.InDelta(t, 0.0, max, 0)
	assert.InDelta(t, 0.0, sum, 0)

	nullVals := NewSeriesListLine([][]float64{{GetNullValue(), GetNullValue()}})
	min, max, sum = getSeriesMinMaxSumMax(nullVals, 0, true)
	assert.InDelta(t, 0.0, min, 0)
	assert.InDelta(t, 0.0, max, 0)
	assert.InDelta(t, 0.0, sum, 0)
}

func TestExpandSingleValueScatterSeries(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		values []float64
	}{
		{
			name:   "empty",
			values: []float64{},
		},
		{
			name:   "single",
			values: []float64{42},
		},
		{
			name:   "multiple",
			values: []float64{1, 2, 3},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := expandSingleValueScatterSeries(tc.values)

			assert.Equal(t, len(tc.values), len(got))
			for i, v := range tc.values {
				require.Len(t, got[i], 1)
				assert.InDelta(t, v, got[i][0], 0.0)
			}
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
			expected: []float64{4.0},
		},
		{
			name: "multiple",
			values: [][]float64{
				{1, 2, 3},
				{4, 5, 6},
			},
			expected: []float64{6, 15},
		},
		{
			name: "unequal_data_length",
			values: [][]float64{
				{1, 2},
				{3, 4, 5},
			},
			expected: []float64{3, 12},
		},
		{
			name: "null_values",
			values: [][]float64{
				{GetNullValue(), 2, 3},
				{4, GetNullValue(), 6},
			},
			expected: []float64{5, 10},
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

func TestSumSeriesValues(t *testing.T) {
	t.Parallel()

	type summableSeries interface {
		SumSeriesValues() []float64
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
				result := series.SumSeriesValues()

				assert.Equal(t, tc.expected, result)
			})
		}
	}
}

func TestSumSeriesDataAndMaxCount(t *testing.T) {
	t.Parallel()

	seriesList := LineSeriesList{
		{Values: []float64{1, 2}, YAxisIndex: 0},
		{Values: []float64{3, 4, 5}, YAxisIndex: 1},
		{Values: []float64{6}, YAxisIndex: 0},
	}

	assert.Equal(t, []float64{10, 6, 5}, sumSeriesData(seriesList, -1))
	assert.Equal(t, []float64{7, 2, 0}, sumSeriesData(seriesList, 0))
	assert.Equal(t, []float64{3, 4, 5}, sumSeriesData(seriesList, 1))
	assert.Equal(t, 3, getSeriesMaxDataCount(seriesList))
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
			MaxFirstIndex:     0,
			MaxIndex:          0,
			Min:               10,
			MinFirstIndex:     0,
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
			MaxFirstIndex:     1,
			Min:               1,
			MinFirstIndex:     0,
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
			MaxFirstIndex:     2,
			MaxIndex:          2,
			Min:               1,
			MinFirstIndex:     0,
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
			MaxFirstIndex:     3,
			MaxIndex:          3,
			Min:               1,
			MinFirstIndex:     0,
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
			MaxFirstIndex:     3,
			MaxIndex:          3,
			Min:               3,
			MinFirstIndex:     0,
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
			MaxFirstIndex:     0,
			MaxIndex:          0,
			Min:               10,
			MinFirstIndex:     0,
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
			MaxFirstIndex:     2,
			MaxIndex:          2,
			Min:               1,
			MinFirstIndex:     0,
			MinIndex:          0,
			Average:           1.5,
			Median:            1.5,
			StandardDeviation: 0.5,
			Skewness:          0.0,
			Kurtosis:          1.0,
		}, seriesList[7].Summary())
	})
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
		_ = seriesList.SumSeriesValues()
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

func BenchmarkSeriesMarkListSplitGlobal(b *testing.B) {
	pure := NewSeriesMarkList(SeriesMarkTypeMax, SeriesMarkTypeMin, SeriesMarkTypeAverage)
	mixed := NewSeriesMarkList(SeriesMarkTypeMax, SeriesMarkTypeMin, SeriesMarkTypeAverage)
	mixed[1].Global = true

	for i := 0; i < b.N; i++ {
		_, _ = pure.splitGlobal()
		_, _ = mixed.splitGlobal()
	}
}

func BenchmarkSeriesMarkListFilterGlobal(b *testing.B) {
	pure := NewSeriesMarkList(SeriesMarkTypeMax, SeriesMarkTypeMin, SeriesMarkTypeAverage)
	mixed := NewSeriesMarkList(SeriesMarkTypeMax, SeriesMarkTypeMin, SeriesMarkTypeAverage)
	mixed[1].Global = true

	for i := 0; i < b.N; i++ {
		_ = pure.filterGlobal(false)
		_ = mixed.filterGlobal(false)
	}
}

func TestCandlestickGenericBidirectionalConversion(t *testing.T) {
	t.Parallel()

	// Create original candlestick data with varied OHLC values
	originalData := []OHLCData{
		{Open: 100, High: 110, Low: 95, Close: 105},
		{Open: 105, High: 115, Low: 98, Close: 108},
		{Open: 108, High: 120, Low: 102, Close: 112},
		// Include invalid data to test null value handling
		{Open: GetNullValue(), High: GetNullValue(), Low: GetNullValue(), Close: GetNullValue()},
		{Open: GetNullValue(), High: 100, Low: 80, Close: GetNullValue()},
	}

	original := CandlestickSeries{
		Data:       originalData,
		YAxisIndex: 1,
		Label:      SeriesLabel{Show: Ptr(true)},
		Name:       "Test Series",
	}

	// Convert CandlestickSeries to GenericSeries
	genericList := CandlestickSeriesList{original}.ToGenericSeriesList()
	require.Len(t, genericList, 1)
	generic := genericList[0]

	// Verify the generic series has correct structure
	assert.Equal(t, ChartTypeCandlestick, generic.Type)
	assert.Equal(t, 1, generic.YAxisIndex)
	assert.Equal(t, "Test Series", generic.Name)
	assert.True(t, *generic.Label.Show)

	// Verify OHLC encoding: 4 values per candlestick
	expectedLen := len(originalData) * 4
	require.Len(t, generic.Values, expectedLen)

	// Check encoded values
	assert.InDelta(t, 100.0, generic.Values[0], 0) // First candle Open
	assert.InDelta(t, 110.0, generic.Values[1], 0) // First candle High
	assert.InDelta(t, 95.0, generic.Values[2], 0)  // First candle Low
	assert.InDelta(t, 105.0, generic.Values[3], 0) // First candle Close

	assert.InDelta(t, 105.0, generic.Values[4], 0) // Second candle Open
	assert.InDelta(t, 115.0, generic.Values[5], 0) // Second candle High
	assert.InDelta(t, 98.0, generic.Values[6], 0)  // Second candle Low
	assert.InDelta(t, 108.0, generic.Values[7], 0) // Second candle Close

	// Check null values are preserved
	assert.InDelta(t, GetNullValue(), generic.Values[12], 0) // Fourth candle Open
	assert.InDelta(t, GetNullValue(), generic.Values[13], 0) // Fourth candle High
	assert.InDelta(t, GetNullValue(), generic.Values[14], 0) // Fourth candle Low
	assert.InDelta(t, GetNullValue(), generic.Values[15], 0) // Fourth candle Close

	// Convert GenericSeries back to CandlestickSeries
	reconstructed := filterSeriesList[CandlestickSeriesList](GenericSeriesList{generic}, ChartTypeCandlestick)
	require.Len(t, reconstructed, 1)
	result := reconstructed[0]

	// Verify the result is equivalent to the original
	assert.Equal(t, original.YAxisIndex, result.YAxisIndex)
	assert.Equal(t, original.Name, result.Name)
	assert.Equal(t, original.Label.Show, result.Label.Show)
	require.Len(t, result.Data, len(original.Data))

	// Check each OHLC data point
	for i, expectedOHLC := range original.Data {
		actualOHLC := result.Data[i]
		assert.InDelta(t, expectedOHLC.Open, actualOHLC.Open, 0)
		assert.InDelta(t, expectedOHLC.High, actualOHLC.High, 0)
		assert.InDelta(t, expectedOHLC.Low, actualOHLC.Low, 0)
		assert.InDelta(t, expectedOHLC.Close, actualOHLC.Close, 0)
	}
}

func TestCandlestickGenericConversionFallback(t *testing.T) {
	t.Parallel()

	// Test fallback behavior when GenericSeries doesn't have OHLC-encoded data
	generic := GenericSeries{
		Values:     []float64{100, 105, 110}, // Not divisible by 4
		YAxisIndex: 0,
		Name:       "Test Fallback",
		Type:       ChartTypeCandlestick,
	}

	// Convert to CandlestickSeries
	reconstructed := filterSeriesList[CandlestickSeriesList](GenericSeriesList{generic}, ChartTypeCandlestick)
	require.Len(t, reconstructed, 1)
	result := reconstructed[0]

	// Should create flat OHLC data where O=H=L=C
	require.Len(t, result.Data, 3)
	for i, expectedValue := range generic.Values {
		ohlc := result.Data[i]
		assert.InDelta(t, expectedValue, ohlc.Open, 0)
		assert.InDelta(t, expectedValue, ohlc.High, 0)
		assert.InDelta(t, expectedValue, ohlc.Low, 0)
		assert.InDelta(t, expectedValue, ohlc.Close, 0)
	}
}
