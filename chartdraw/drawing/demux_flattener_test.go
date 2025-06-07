package drawing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/image/math/fixed"
)

type recordAdder struct{ starts, adds []fixed.Point26_6 }

func (r *recordAdder) Start(p fixed.Point26_6)      { r.starts = append(r.starts, p) }
func (r *recordAdder) Add1(p fixed.Point26_6)       { r.adds = append(r.adds, p) }
func (r *recordAdder) Add2(b, c fixed.Point26_6)    {}
func (r *recordAdder) Add3(b, c, d fixed.Point26_6) {}

func TestDemuxFlattener(t *testing.T) {
	t.Parallel()

	r1 := &recordFlattener{}
	r2 := &recordFlattener{}
	d := DemuxFlattener{Flatteners: []Flattener{r1, r2}}
	d.MoveTo(1, 2)
	d.LineTo(3, 4)
	d.End()
	assert.Equal(t, r1.moves, r2.moves)
	assert.Equal(t, []string{"M1.0,2.0", "L3.0,4.0"}, r1.moves)
}

func TestFtLineBuilder(t *testing.T) {
	t.Parallel()

	ad := &recordAdder{}
	ft := FtLineBuilder{Adder: ad}
	ft.MoveTo(1, 1)
	ft.LineTo(2, 3)
	ft.End()
	if assert.Len(t, ad.starts, 1) {
		assert.Equal(t, fixed.Int26_6(64), ad.starts[0].X)
		assert.Equal(t, fixed.Int26_6(64), ad.starts[0].Y)
	}
	if assert.Len(t, ad.adds, 1) {
		assert.Equal(t, fixed.Int26_6(128), ad.adds[0].X)
		assert.Equal(t, fixed.Int26_6(192), ad.adds[0].Y)
	}
}
