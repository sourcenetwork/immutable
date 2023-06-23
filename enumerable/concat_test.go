package enumerable

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConcatYieldsNothingGivenEmpty(t *testing.T) {
	concat := Concat[int]()

	hasNext, err := concat.Next()
	require.NoError(t, err)
	require.False(t, hasNext)
}

func TestConcatYieldsItemsFromSource(t *testing.T) {
	v1 := 1
	v2 := 2
	v3 := 3
	source1 := New([]int{v1, v2, v3})

	concat := Concat(source1)

	hasNext, err := concat.Next()
	require.NoError(t, err)
	require.True(t, hasNext)

	r1, err := concat.Value()
	require.NoError(t, err)
	require.Equal(t, v1, r1)

	hasNext, err = concat.Next()
	require.NoError(t, err)
	require.True(t, hasNext)

	r2, err := concat.Value()
	require.NoError(t, err)
	require.Equal(t, v2, r2)

	hasNext, err = concat.Next()
	require.NoError(t, err)
	require.True(t, hasNext)

	r3, err := concat.Value()
	require.NoError(t, err)
	require.Equal(t, v3, r3)

	hasNext, err = concat.Next()
	require.NoError(t, err)
	require.False(t, hasNext)
}

func TestConcatYieldsItemsFromSourceInOrder(t *testing.T) {
	v1 := 1
	v2 := 2
	v3 := 3
	v4 := 4
	v5 := 5
	v6 := 6
	source1 := NewQueue[int]()
	var s1 Enumerable[int] = source1
	source2 := New([]int{v1, v2, v3})
	source3 := New([]int{v4, v5})

	concat := Concat(s1, source2, source3)

	// Start yielding *before* source1 has any items
	hasNext, err := concat.Next()
	require.NoError(t, err)
	require.True(t, hasNext)

	r1, err := concat.Value()
	require.NoError(t, err)
	require.Equal(t, v1, r1)

	// Put an item into source1
	err = source1.Put(v6)
	require.NoError(t, err)

	// Assert that the yielding of items from source2 is
	// not interupted by source1 recieving items
	hasNext, err = concat.Next()
	require.NoError(t, err)
	require.True(t, hasNext)

	r2, err := concat.Value()
	require.NoError(t, err)
	require.Equal(t, v2, r2)

	hasNext, err = concat.Next()
	require.NoError(t, err)
	require.True(t, hasNext)

	r3, err := concat.Value()
	require.NoError(t, err)
	require.Equal(t, v3, r3)

	hasNext, err = concat.Next()
	require.NoError(t, err)
	require.True(t, hasNext)

	// Assert that source3's items are yielded after
	// source2's
	r4, err := concat.Value()
	require.NoError(t, err)
	require.Equal(t, v4, r4)

	hasNext, err = concat.Next()
	require.NoError(t, err)
	require.True(t, hasNext)

	r5, err := concat.Value()
	require.NoError(t, err)
	require.Equal(t, v5, r5)

	// Then assert that source1's items are yielded
	// as the concat circles back round
	hasNext, err = concat.Next()
	require.NoError(t, err)
	require.True(t, hasNext)

	r6, err := concat.Value()
	require.NoError(t, err)
	require.Equal(t, v6, r6)

	hasNext, err = concat.Next()
	require.NoError(t, err)
	require.False(t, hasNext)
}
