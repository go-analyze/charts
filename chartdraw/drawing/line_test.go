package drawing

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBresenhamDiagonal(t *testing.T) {
	t.Parallel()

	img := image.NewRGBA(image.Rect(0, 0, 5, 5))
	Bresenham(img, color.White, 0, 0, 4, 4)

	for i := 0; i <= 4; i++ {
		r, g, b, a := img.At(i, i).RGBA()
		assert.Equal(t, uint32(0xffff), r)
		assert.Equal(t, uint32(0xffff), g)
		assert.Equal(t, uint32(0xffff), b)
		assert.Equal(t, uint32(0xffff), a)
	}

	_, _, _, a := img.At(0, 1).RGBA()
	assert.Equal(t, uint32(0), a)
}

func TestPolylineBresenham(t *testing.T) {
	t.Parallel()

	img := image.NewRGBA(image.Rect(0, 0, 5, 5))
	PolylineBresenham(img, color.White, 0, 0, 2, 0, 2, 2)

	expected := [][2]int{{0, 0}, {1, 0}, {2, 0}, {2, 1}, {2, 2}}
	for _, p := range expected {
		_, _, _, a := img.At(p[0], p[1]).RGBA()
		assert.Equal(t, uint32(0xffff), a)
	}

	_, _, _, a := img.At(1, 1).RGBA()
	assert.Equal(t, uint32(0), a)
}
