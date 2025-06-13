package charts

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNullValue(t *testing.T) {
	t.Parallel()

	assert.InDelta(t, math.MaxFloat64, GetNullValue(), 0.0)
}
