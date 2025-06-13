package drawing

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/image/draw"
)

func TestDrawImageTransform(t *testing.T) {
	t.Parallel()

	src := image.NewRGBA(image.Rect(0, 0, 1, 1))
	src.Set(0, 0, color.White)
	dst := image.NewRGBA(image.Rect(0, 0, 3, 3))
	DrawImage(src, dst, NewTranslationMatrix(1, 1), draw.Over, LinearFilter)
	_, _, _, a := dst.At(1, 1).RGBA()
	assert.Equal(t, uint32(0xffff), a)
}

func TestDrawImageScale(t *testing.T) {
	t.Parallel()

	src := image.NewRGBA(image.Rect(0, 0, 1, 1))
	src.Set(0, 0, color.White)
	dst := image.NewRGBA(image.Rect(0, 0, 2, 2))
	DrawImage(src, dst, NewScaleMatrix(2, 2), draw.Over, LinearFilter)
	_, _, _, a := dst.At(1, 1).RGBA()
	assert.Equal(t, uint32(0xffff), a)
}

func TestDrawImageFilters(t *testing.T) {
	t.Parallel()

	src := image.NewRGBA(image.Rect(0, 0, 1, 1))
	src.Set(0, 0, color.White)

	dst1 := image.NewRGBA(image.Rect(0, 0, 1, 1))
	DrawImage(src, dst1, NewIdentityMatrix(), draw.Over, BilinearFilter)
	_, _, _, a := dst1.At(0, 0).RGBA()
	assert.Equal(t, uint32(0xffff), a)

	dst2 := image.NewRGBA(image.Rect(0, 0, 1, 1))
	DrawImage(src, dst2, NewIdentityMatrix(), draw.Over, BicubicFilter)
	_, _, _, a = dst2.At(0, 0).RGBA()
	assert.Equal(t, uint32(0xffff), a)
}
