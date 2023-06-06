package enumerable

// Queue is an extention of the enumerable interface allowing individual
// items to be added into the enumerable.
//
// Added items will be yielded in a FIFO order.  Items may be added after
// enumeration has begun.
type Queue[T any] interface {
	Enumerable[T]
	// Put adds an item to the queue.
	Put(T) error
	// Size returns the current length of the backing array.
	//
	// This may include empty space where yield items previously resided.
	// Useful for testing and debugging.
	Size() int
}

type queue[T any] struct {
	values       []T
	currentIndex int
	lastSetIndex int
	zeroIndexSet bool
}

var _ Queue[any] = (*queue[any])(nil)

// NewQueue creates an empty FIFO queue.
//
// It is implemented using a dynamically sized ring-buffer.
func NewQueue[T any]() Queue[T] {
	return &queue[T]{
		values:       []T{},
		currentIndex: -1,
		lastSetIndex: -1,
	}
}

func (q *queue[T]) Put(value T) error {
	var index int
	if !q.zeroIndexSet {
		// If the zero-index is empty, we should use it - circling
		// the ring buffer back to the beginning.
		index = 0
		q.zeroIndexSet = true
	} else {
		index = q.lastSetIndex + 1
	}

	if index >= len(q.values) {
		// For now, increasing the size one at a time is likely optimal
		// for the only useage of the queue type.  We may wish to change
		// this at somepoint however.
		newValues := make([]T, len(q.values)+1)
		copy(newValues, q.values)
		q.values = newValues
	}

	q.values[index] = value
	q.lastSetIndex = index

	return nil
}

func (q *queue[T]) Next() (bool, error) {
	if q.currentIndex >= q.lastSetIndex && !q.zeroIndexSet {
		// We have escaped the value-window and have no next value.
		return false, nil
	}

	nextIndex := q.currentIndex + 1
	if nextIndex == len(q.values) {
		// Circle back to the beginning
		nextIndex = 0
	}

	// If the next index is the zero-index the value is consumed (implicitly), so we update
	// the flag here.
	// Note: This may also be zero if this is the first Next call following either intialization
	// or a reset, it cannot be moved inside the len(q.values) if-block above.
	if nextIndex == 0 {
		q.zeroIndexSet = false
	}

	q.currentIndex = nextIndex
	return true, nil
}

func (q *queue[T]) Value() (T, error) {
	return q.values[q.currentIndex], nil
}

func (q *queue[T]) Reset() {
	q.values = []T{}
	q.currentIndex = -1
	q.lastSetIndex = -1
	q.zeroIndexSet = false
}

func (q *queue[T]) Size() int {
	return len(q.values)
}
