package chartdraw

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuffer(t *testing.T) {
	t.Parallel()

	buffer := NewValueBuffer()

	buffer.Enqueue(1)
	assert.Equal(t, 1, buffer.Len())
	assert.Equal(t, float64(1), buffer.Peek())
	assert.Equal(t, float64(1), buffer.PeekBack())

	buffer.Enqueue(2)
	assert.Equal(t, 2, buffer.Len())
	assert.Equal(t, float64(1), buffer.Peek())
	assert.Equal(t, float64(2), buffer.PeekBack())

	buffer.Enqueue(3)
	assert.Equal(t, 3, buffer.Len())
	assert.Equal(t, float64(1), buffer.Peek())
	assert.Equal(t, float64(3), buffer.PeekBack())

	buffer.Enqueue(4)
	assert.Equal(t, 4, buffer.Len())
	assert.Equal(t, float64(1), buffer.Peek())
	assert.Equal(t, float64(4), buffer.PeekBack())

	buffer.Enqueue(5)
	assert.Equal(t, 5, buffer.Len())
	assert.Equal(t, float64(1), buffer.Peek())
	assert.Equal(t, float64(5), buffer.PeekBack())

	buffer.Enqueue(6)
	assert.Equal(t, 6, buffer.Len())
	assert.Equal(t, float64(1), buffer.Peek())
	assert.Equal(t, float64(6), buffer.PeekBack())

	buffer.Enqueue(7)
	assert.Equal(t, 7, buffer.Len())
	assert.Equal(t, float64(1), buffer.Peek())
	assert.Equal(t, float64(7), buffer.PeekBack())

	buffer.Enqueue(8)
	assert.Equal(t, 8, buffer.Len())
	assert.Equal(t, float64(1), buffer.Peek())
	assert.Equal(t, float64(8), buffer.PeekBack())

	value := buffer.Dequeue()
	assert.Equal(t, float64(1), value)
	assert.Equal(t, 7, buffer.Len())
	assert.Equal(t, float64(2), buffer.Peek())
	assert.Equal(t, float64(8), buffer.PeekBack())

	value = buffer.Dequeue()
	assert.Equal(t, float64(2), value)
	assert.Equal(t, 6, buffer.Len())
	assert.Equal(t, float64(3), buffer.Peek())
	assert.Equal(t, float64(8), buffer.PeekBack())

	value = buffer.Dequeue()
	assert.Equal(t, float64(3), value)
	assert.Equal(t, 5, buffer.Len())
	assert.Equal(t, float64(4), buffer.Peek())
	assert.Equal(t, float64(8), buffer.PeekBack())

	value = buffer.Dequeue()
	assert.Equal(t, float64(4), value)
	assert.Equal(t, 4, buffer.Len())
	assert.Equal(t, float64(5), buffer.Peek())
	assert.Equal(t, float64(8), buffer.PeekBack())

	value = buffer.Dequeue()
	assert.Equal(t, float64(5), value)
	assert.Equal(t, 3, buffer.Len())
	assert.Equal(t, float64(6), buffer.Peek())
	assert.Equal(t, float64(8), buffer.PeekBack())

	value = buffer.Dequeue()
	assert.Equal(t, float64(6), value)
	assert.Equal(t, 2, buffer.Len())
	assert.Equal(t, float64(7), buffer.Peek())
	assert.Equal(t, float64(8), buffer.PeekBack())

	value = buffer.Dequeue()
	assert.Equal(t, float64(7), value)
	assert.Equal(t, 1, buffer.Len())
	assert.Equal(t, float64(8), buffer.Peek())
	assert.Equal(t, float64(8), buffer.PeekBack())

	value = buffer.Dequeue()
	assert.Equal(t, float64(8), value)
	assert.Equal(t, 0, buffer.Len())
	assert.Zero(t, buffer.Peek())
	assert.Zero(t, buffer.PeekBack())
}

func TestBufferClear(t *testing.T) {
	t.Parallel()

	buffer := NewValueBuffer()
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)
	buffer.Enqueue(1)

	assert.Equal(t, 8, buffer.Len())

	buffer.Clear()
	assert.Equal(t, 0, buffer.Len())
	assert.Zero(t, buffer.Peek())
	assert.Zero(t, buffer.PeekBack())
}

func TestBufferArray(t *testing.T) {
	t.Parallel()

	buffer := NewValueBuffer()
	buffer.Enqueue(1)
	buffer.Enqueue(2)
	buffer.Enqueue(3)
	buffer.Enqueue(4)
	buffer.Enqueue(5)

	contents := buffer.Array()
	require.Len(t, contents, 5)
	assert.Equal(t, float64(1), contents[0])
	assert.Equal(t, float64(2), contents[1])
	assert.Equal(t, float64(3), contents[2])
	assert.Equal(t, float64(4), contents[3])
	assert.Equal(t, float64(5), contents[4])
}

func TestBufferEach(t *testing.T) {
	t.Parallel()

	buffer := NewValueBuffer()

	for x := 1; x < 17; x++ {
		buffer.Enqueue(float64(x))
	}

	called := 0
	buffer.Each(func(_ int, v float64) {
		if v == float64(called+1) {
			called++
		}
	})

	assert.Equal(t, 16, called)
}

func TestNewBuffer(t *testing.T) {
	t.Parallel()

	empty := NewValueBuffer()
	assert.NotNil(t, empty)
	assert.Zero(t, empty.Len())
	assert.Equal(t, bufferDefaultCapacity, empty.Capacity())
	assert.Zero(t, empty.Peek())
	assert.Zero(t, empty.PeekBack())
}

func TestNewBufferWithValues(t *testing.T) {
	t.Parallel()

	values := NewValueBuffer(1, 2, 3, 4, 5)
	assert.NotNil(t, values)
	assert.Equal(t, 5, values.Len())
	assert.Equal(t, float64(1), values.Peek())
	assert.Equal(t, float64(5), values.PeekBack())
}

func TestBufferGrowth(t *testing.T) {
	t.Parallel()

	values := NewValueBuffer(1, 2, 3, 4, 5)
	for i := 0; i < 1<<10; i++ {
		values.Enqueue(float64(i))
	}

	assert.Equal(t, float64(1<<10-1), values.PeekBack())
}
