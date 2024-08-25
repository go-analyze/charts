package matrix

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	m := New(10, 5)
	rows, cols := m.Size()
	assert.Equal(t, 10, rows)
	assert.Equal(t, 5, cols)
	assert.Zero(t, m.Get(0, 0))
	assert.Zero(t, m.Get(9, 4))
}

func TestNewWithValues(t *testing.T) {
	m := New(5, 2, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	rows, cols := m.Size()
	assert.Equal(t, 5, rows)
	assert.Equal(t, 2, cols)
	assert.Equal(t, float64(1), m.Get(0, 0))
	assert.Equal(t, float64(10), m.Get(4, 1))
}

func TestIdentitiy(t *testing.T) {
	id := Identity(5)
	rows, cols := id.Size()
	assert.Equal(t, 5, rows)
	assert.Equal(t, 5, cols)
	assert.Equal(t, float64(1), id.Get(0, 0))
	assert.Equal(t, float64(1), id.Get(1, 1))
	assert.Equal(t, float64(1), id.Get(2, 2))
	assert.Equal(t, float64(1), id.Get(3, 3))
	assert.Equal(t, float64(1), id.Get(4, 4))
	assert.Equal(t, float64(0), id.Get(0, 1))
	assert.Equal(t, float64(0), id.Get(1, 0))
	assert.Equal(t, float64(0), id.Get(4, 0))
	assert.Equal(t, float64(0), id.Get(0, 4))
}

func TestNewFromArrays(t *testing.T) {
	m := NewFromArrays([][]float64{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
	})
	assert.NotNil(t, m)

	rows, cols := m.Size()
	assert.Equal(t, 2, rows)
	assert.Equal(t, 4, cols)
}

func TestOnes(t *testing.T) {
	ones := Ones(5, 10)
	rows, cols := ones.Size()
	assert.Equal(t, 5, rows)
	assert.Equal(t, 10, cols)

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			assert.Equal(t, float64(1), ones.Get(row, col))
		}
	}
}

func TestMatrixEpsilon(t *testing.T) {
	ones := Ones(2, 2)
	ones = ones.WithEpsilon(0.001)
	assert.Equal(t, 0.001, ones.Epsilon())
}

func TestMatrixArrays(t *testing.T) {
	m := NewFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
	})

	assert.NotNil(t, m)

	arrays := m.Arrays()

	assert.Equal(t, arrays, [][]float64{
		{1, 2, 3},
		{4, 5, 6},
	})
}

func TestMatrixIsSquare(t *testing.T) {
	assert.False(t, NewFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
	}).IsSquare())

	assert.False(t, NewFromArrays([][]float64{
		{1, 2},
		{3, 4},
		{5, 6},
	}).IsSquare())

	assert.True(t, NewFromArrays([][]float64{
		{1, 2},
		{3, 4},
	}).IsSquare())
}

func TestMatrixIsSymmetric(t *testing.T) {
	assert.False(t, NewFromArrays([][]float64{
		{1, 2, 3},
		{2, 1, 2},
	}).IsSymmetric())

	assert.False(t, NewFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}).IsSymmetric())

	assert.True(t, NewFromArrays([][]float64{
		{1, 2, 3},
		{2, 1, 2},
		{3, 2, 1},
	}).IsSymmetric())

}

func TestMatrixGet(t *testing.T) {
	m := NewFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})

	assert.Equal(t, float64(1), m.Get(0, 0))
	assert.Equal(t, float64(2), m.Get(0, 1))
	assert.Equal(t, float64(3), m.Get(0, 2))
	assert.Equal(t, float64(4), m.Get(1, 0))
	assert.Equal(t, float64(5), m.Get(1, 1))
	assert.Equal(t, float64(6), m.Get(1, 2))
	assert.Equal(t, float64(7), m.Get(2, 0))
	assert.Equal(t, float64(8), m.Get(2, 1))
	assert.Equal(t, float64(9), m.Get(2, 2))
}

func TestMatrixSet(t *testing.T) {
	m := NewFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})

	m.Set(1, 1, 99)
	assert.Equal(t, float64(99), m.Get(1, 1))
}

func TestMatrixCol(t *testing.T) {
	m := NewFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})

	assert.Equal(t, Vector([]float64{1, 4, 7}), m.Col(0))
	assert.Equal(t, Vector([]float64{2, 5, 8}), m.Col(1))
	assert.Equal(t, Vector([]float64{3, 6, 9}), m.Col(2))
}

func TestMatrixRow(t *testing.T) {
	m := NewFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})

	assert.Equal(t, Vector([]float64{1, 2, 3}), m.Row(0))
	assert.Equal(t, Vector([]float64{4, 5, 6}), m.Row(1))
	assert.Equal(t, Vector([]float64{7, 8, 9}), m.Row(2))
}

func TestMatrixSwapRows(t *testing.T) {
	m := NewFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})

	m.SwapRows(0, 1)

	assert.Equal(t, Vector([]float64{4, 5, 6}), m.Row(0))
	assert.Equal(t, Vector([]float64{1, 2, 3}), m.Row(1))
	assert.Equal(t, Vector([]float64{7, 8, 9}), m.Row(2))
}

func TestMatrixCopy(t *testing.T) {
	m := NewFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})

	m2 := m.Copy()
	assert.False(t, m == m2)
	assert.True(t, m.Equals(m2))
}

func TestMatrixDiagonalVector(t *testing.T) {
	m := NewFromArrays([][]float64{
		{1, 4, 7},
		{4, 2, 8},
		{7, 8, 3},
	})

	diag := m.DiagonalVector()
	assert.Equal(t, Vector([]float64{1, 2, 3}), diag)
}

func TestMatrixDiagonalVectorLandscape(t *testing.T) {
	m := NewFromArrays([][]float64{
		{1, 4, 7, 99},
		{4, 2, 8, 99},
	})

	diag := m.DiagonalVector()
	assert.Equal(t, Vector([]float64{1, 2}), diag)
}

func TestMatrixDiagonalVectorPortrait(t *testing.T) {
	m := NewFromArrays([][]float64{
		{1, 4},
		{4, 2},
		{99, 99},
	})

	diag := m.DiagonalVector()
	assert.Equal(t, Vector([]float64{1, 2}), diag)
}

func TestMatrixDiagonal(t *testing.T) {
	m := NewFromArrays([][]float64{
		{1, 4, 7},
		{4, 2, 8},
		{7, 8, 3},
	})

	m2 := NewFromArrays([][]float64{
		{1, 0, 0},
		{0, 2, 0},
		{0, 0, 3},
	})

	assert.True(t, m.Diagonal().Equals(m2))
}

func TestMatrixEquals(t *testing.T) {
	m := NewFromArrays([][]float64{
		{1, 4, 7},
		{4, 2, 8},
		{7, 8, 3},
	})

	assert.False(t, m.Equals(nil))
	var nilMatrix *Matrix
	assert.True(t, nilMatrix.Equals(nil))
	assert.False(t, m.Equals(New(1, 1)))
	assert.False(t, m.Equals(New(3, 3)))
	assert.True(t, m.Equals(New(3, 3, 1, 4, 7, 4, 2, 8, 7, 8, 3)))
}

func TestMatrixL(t *testing.T) {
	m := NewFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})

	l := m.L()
	assert.True(t, l.Equals(New(3, 3, 1, 2, 3, 0, 5, 6, 0, 0, 9)))
}

func TestMatrixU(t *testing.T) {
	m := NewFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})

	u := m.U()
	assert.True(t, u.Equals(New(3, 3, 0, 0, 0, 4, 0, 0, 7, 8, 0)))
}

func TestMatrixString(t *testing.T) {
	m := NewFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})

	assert.Equal(t, "1 2 3 \n4 5 6 \n7 8 9 \n", m.String())
}

func TestMatrixLU(t *testing.T) {
	m := NewFromArrays([][]float64{
		{1, 3, 5},
		{2, 4, 7},
		{1, 1, 0},
	})

	l, u, p := m.LU()
	assert.NotNil(t, l)
	assert.NotNil(t, u)
	assert.NotNil(t, p)
}

func TestMatrixQR(t *testing.T) {
	m := NewFromArrays([][]float64{
		{12, -51, 4},
		{6, 167, -68},
		{-4, 24, -41},
	})

	q, r := m.QR()
	assert.NotNil(t, q)
	assert.NotNil(t, r)
}

func TestMatrixTranspose(t *testing.T) {
	m := NewFromArrays([][]float64{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{10, 11, 12},
	})

	m2 := m.Transpose()

	rows, cols := m2.Size()
	assert.Equal(t, 3, rows)
	assert.Equal(t, 4, cols)

	assert.Equal(t, float64(1), m2.Get(0, 0))
	assert.Equal(t, float64(10), m2.Get(0, 3))
	assert.Equal(t, float64(3), m2.Get(2, 0))
}
