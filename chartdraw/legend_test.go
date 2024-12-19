package chartdraw

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLegend(t *testing.T) {
	t.Parallel()

	graph := Chart{
		Series: []Series{
			ContinuousSeries{
				Name:    "A test series",
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},
		},
	}

	//note we have to do this as a separate step because we need a reference to graph
	graph.Elements = []Renderable{
		Legend(&graph),
	}
	buf := bytes.NewBuffer([]byte{})
	require.NoError(t, graph.Render(PNG, buf))
	assert.NotZero(t, buf.Len())
}
