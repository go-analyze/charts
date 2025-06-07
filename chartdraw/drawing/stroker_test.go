package drawing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLineStrokerLine(t *testing.T) {
	t.Parallel()

	rec := &recordFlattenerEnd{}
	ls := NewLineStroker(rec)
	ls.MoveTo(0, 0)
	ls.LineTo(2, 0)
	ls.End()

	expect := []string{"M0.0,-0.5", "L2.0,-0.5", "L2.0,0.5", "L0.0,0.5", "L0.0,-0.5", "E"}
	assert.Equal(t, expect, rec.moves)
}
