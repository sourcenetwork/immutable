package enumerable

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueueYieldsNothingGivenEmpty(t *testing.T) {
	queue := NewQueue[int]()

	hasNext, err := queue.Next()
	require.NoError(t, err)
	require.False(t, hasNext)
}

func TestQueueYieldsSingleItemGivenValueAdded(t *testing.T) {
	v1 := 1
	queue := NewQueue[int]()

	queue.Put(v1)

	hasNext, err := queue.Next()
	require.NoError(t, err)
	require.True(t, hasNext)

	r1, err := queue.Value()
	require.NoError(t, err)
	require.Equal(t, v1, r1)

	hasNext, err = queue.Next()
	require.NoError(t, err)
	require.False(t, hasNext)
}

func TestQueueValueReturnsSameItemEachTime(t *testing.T) {
	v1 := 1
	queue := NewQueue[int]()

	queue.Put(v1)

	hasNext, err := queue.Next()
	require.NoError(t, err)
	require.True(t, hasNext)

	r1, err := queue.Value()
	require.NoError(t, err)
	require.Equal(t, v1, r1)

	// Calling Value multiple times without progressing the enumeration
	// via next should keep returning the same value
	r1, err = queue.Value()
	require.NoError(t, err)
	require.Equal(t, v1, r1)

	hasNext, err = queue.Next()
	require.NoError(t, err)
	require.False(t, hasNext)
}

func TestQueueYields100ItemsGiven100ItemsAddedPreviously(t *testing.T) {
	numberOfItems := 100
	queue := NewQueue[int]()

	for i := 1; i <= numberOfItems; i++ {
		queue.Put(i)
	}

	for i := 1; i <= numberOfItems; i++ {
		hasNext, err := queue.Next()
		require.NoError(t, err)
		require.True(t, hasNext)

		r1, err := queue.Value()
		require.NoError(t, err)
		require.Equal(t, i, r1)
	}
}

func TestQueueYieldsItemsGiven100ItemsReadAsAdded(t *testing.T) {
	numberOfItems := 100
	queue := NewQueue[int]()

	for i := 1; i <= numberOfItems; i++ {
		queue.Put(i)

		hasNext, err := queue.Next()
		require.NoError(t, err)
		require.True(t, hasNext)

		r1, err := queue.Value()
		require.NoError(t, err)
		require.Equal(t, i, r1)
	}
}

func TestQueueYieldsItemsGiven100ItemsReadInPairsAsAdded(t *testing.T) {
	numberOfItems := 100
	queue := NewQueue[int]()

	for i := 1; i <= numberOfItems; i = i + 2 {
		vi1 := i
		vi2 := i + 1
		queue.Put(vi1)
		queue.Put(vi2)

		hasNext, err := queue.Next()
		require.NoError(t, err)
		require.True(t, hasNext)

		ri1, err := queue.Value()
		require.NoError(t, err)
		require.Equal(t, vi1, ri1)

		hasNext, err = queue.Next()
		require.NoError(t, err)
		require.True(t, hasNext)

		ri2, err := queue.Value()
		require.NoError(t, err)
		require.Equal(t, vi2, ri2)
	}
}

func TestQueueYieldsNothinGivenReset(t *testing.T) {
	queue := NewQueue[int]()

	queue.Reset()

	hasNext, err := queue.Next()
	require.NoError(t, err)
	require.False(t, hasNext)
}

func TestQueueYieldsNothinGivenResetAfterValueAdded(t *testing.T) {
	v1 := 1
	queue := NewQueue[int]()

	queue.Put(v1)
	queue.Reset()

	hasNext, err := queue.Next()
	require.NoError(t, err)
	require.False(t, hasNext)
}

func TestQueueYieldsSingleItemGivenValueAddedAfterReset(t *testing.T) {
	v1 := 1
	v2 := 2
	queue := NewQueue[int]()

	queue.Put(v1)
	queue.Reset()
	queue.Put(v2)

	hasNext, err := queue.Next()
	require.NoError(t, err)
	require.True(t, hasNext)

	r1, err := queue.Value()
	require.NoError(t, err)
	require.Equal(t, v2, r1)

	hasNext, err = queue.Next()
	require.NoError(t, err)
	require.False(t, hasNext)
}

func TestQueueYieldsItemCorrectlyGivenCircle(t *testing.T) {
	v1 := 1
	v2 := 2
	v3 := 3
	v4 := 4
	queue := NewQueue[int]()

	queue.Put(v1)
	queue.Put(v2)
	queue.Put(v3)
	// [1, 2, 3]

	hasNext, err := queue.Next()
	require.NoError(t, err)
	require.True(t, hasNext)

	r1, err := queue.Value()
	// [, 2, 3]
	require.NoError(t, err)
	require.Equal(t, v1, r1)

	hasNext, err = queue.Next()
	require.NoError(t, err)
	require.True(t, hasNext)

	r2, err := queue.Value()
	// [, , 3]
	require.NoError(t, err)
	require.Equal(t, v2, r2)

	queue.Put(v4)
	// [4, , 3]

	hasNext, err = queue.Next()
	require.NoError(t, err)
	require.True(t, hasNext)

	r3, err := queue.Value()
	// [4, ,]
	require.NoError(t, err)
	require.Equal(t, v3, r3)

	hasNext, err = queue.Next()
	require.NoError(t, err)
	require.True(t, hasNext)

	r4, err := queue.Value()
	// [, ,]
	require.NoError(t, err)
	require.Equal(t, v4, r4)

	hasNext, err = queue.Next()
	require.NoError(t, err)
	require.False(t, hasNext)

	size := queue.Size()
	require.Equal(t, 3, size)
}

func TestQueueYieldsItemAddedAfterFullEnumeration(t *testing.T) {
	v1 := 1
	v2 := 2
	v3 := 3
	queue := NewQueue[int]()

	queue.Put(v1)
	queue.Put(v2)

	hasNext, err := queue.Next()
	require.NoError(t, err)
	require.True(t, hasNext)

	hasNext, err = queue.Next()
	require.NoError(t, err)
	require.True(t, hasNext)

	hasNext, err = queue.Next()
	require.NoError(t, err)
	require.False(t, hasNext)

	queue.Put(v3)

	hasNext, err = queue.Next()
	require.NoError(t, err)

	require.True(t, hasNext)

	r3, err := queue.Value()
	require.NoError(t, err)

	require.Equal(t, v3, r3)
}
