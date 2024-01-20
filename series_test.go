package charts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSeriesListDataFromValues(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(SeriesList{
		{
			Type: ChartTypeBar,
			Data: []SeriesData{
				{
					Value: 1.0,
				},
			},
		},
	}, NewSeriesListDataFromValues([][]float64{
		{
			1,
		},
	}, ChartTypeBar))
}

func TestSeriesLists(t *testing.T) {
	assert := assert.New(t)
	seriesList := NewSeriesListDataFromValues([][]float64{
		{
			1,
			2,
		},
		{
			10,
		},
	}, ChartTypeBar)

	assert.Equal(2, len(seriesList.Filter(ChartTypeBar)))
	assert.Equal(0, len(seriesList.Filter(ChartTypeLine)))

	max, min := seriesList.GetMaxMin(0)
	assert.Equal(float64(10), max)
	assert.Equal(float64(1), min)

	assert.Equal(seriesSummary{
		MaxIndex:     1,
		MaxValue:     2,
		MinIndex:     0,
		MinValue:     1,
		AverageValue: 1.5,
	}, seriesList[0].Summary())
}

func TestFormatter(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("a: 12%", NewPieLabelFormatter([]string{
		"a",
		"b",
	}, "")(0, 10, 0.12))

	assert.Equal("10", NewValueLabelFormatter([]string{
		"a",
		"b",
	}, "")(0, 10, 0.12))
}
