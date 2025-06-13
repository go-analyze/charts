package drawing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/image/math/fixed"
)

type mockAdder struct {
	starts []fixed.Point26_6
	adds   []fixed.Point26_6
}

func (m *mockAdder) Start(p fixed.Point26_6)      { m.starts = append(m.starts, p) }
func (m *mockAdder) Add1(p fixed.Point26_6)       { m.adds = append(m.adds, p) }
func (m *mockAdder) Add2(b, c fixed.Point26_6)    {}
func (m *mockAdder) Add3(b, c, d fixed.Point26_6) {}

func TestFtLineBuilderMoveToLineTo(t *testing.T) {
	t.Parallel()

	ad := &mockAdder{}
	ft := FtLineBuilder{Adder: ad}
	ft.MoveTo(1, 2)
	ft.LineTo(3, 4)
	ft.End()

	if assert.Len(t, ad.starts, 1) {
		assert.Equal(t, fixed.Int26_6(64), ad.starts[0].X)
		assert.Equal(t, fixed.Int26_6(128), ad.starts[0].Y)
	}
	if assert.Len(t, ad.adds, 1) {
		assert.Equal(t, fixed.Int26_6(192), ad.adds[0].X)
		assert.Equal(t, fixed.Int26_6(256), ad.adds[0].Y)
	}
}
