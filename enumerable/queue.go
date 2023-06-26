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

// For now, increasing the size one at a time is likely optimal
// for the only useage of the queue type.  We may wish to change
// this at somepoint however.
const growthSize int = 1

type queue[T any] struct {
	// The values slice of this queue.
	//
	// Note: queue is implementated as a dynamically sized ring buffer, the zero index
	// is not nessecarily the next/current value. Also note that values are not explicitly
	// removed from this slice, which values are still 'in' the queue is tracked by index.
	values []T

	// The index of the current value.
	currentIndex int

	// The index of the last value added to the queue.
	lastSetIndex int

	// Will be true if values[0] has been set.
	zeroIndexSet bool

	// Will be true a value has been attempted to be read from an empty queue.
	waitingForWrite bool
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
	index := q.lastSetIndex + 1

	if index >= len(q.values) {
		if len(q.values) == 0 {
			q.values = make([]T, growthSize)
			q.currentIndex = -1
		} else if q.zeroIndexSet {
			// If the zero index is occupied, we cannot loop back to it here
			// and instead need to grow the values slice.
			newValues := make([]T, len(q.values)+growthSize)
			copy(newValues, q.values[:index])
			q.values = newValues
		} else {
			index = 0
			if q.currentIndex >= len(q.values) {
				q.currentIndex = -1
			}
		}
	} else if index == q.currentIndex {
		// If the write index has caught up to the read index
		// the new value needs to be written between the two
		// e.g: [3,4,here,1,2]
		// Note: The last value read should not be overwritten, as `Value`
		// may be called multiple times on it after a single `Next` call.
		newValues := make([]T, len(q.values)+growthSize)
		copy(newValues, q.values[:index])
		copy(newValues[index+growthSize:], q.values[index:])
		q.values = newValues
		// Shift the current read index to reflect its new location.
		q.currentIndex += growthSize
	}

	if index == 0 {
		q.zeroIndexSet = true
	}

	q.values[index] = value
	q.lastSetIndex = index

	return nil
}

func (q *queue[T]) Next() (bool, error) {
	// If the previous index was the zero-index the value is consumed (implicitly), so we update
	// the flag here.
	if q.currentIndex == 0 {
		q.zeroIndexSet = false
	}

	nextIndex := q.currentIndex + 1
	var hasValue bool
	if nextIndex >= len(q.values) {
		if q.zeroIndexSet {
			// Circle back to the beginning
			nextIndex = 0
			hasValue = true
		} else {
			hasValue = false
			if q.currentIndex == len(q.values) {
				// If we have reached the end of the values slice, and the previous
				// index was already out of bounds, we should avoid growing it further.
				nextIndex = q.currentIndex
			}
		}
	} else {
		// If the previous read index was the last index written to then the value has been
		// consumed and we have reached the edge of the ring: [v2, v3,^we are here, , v1]
		hasValue = q.currentIndex != q.lastSetIndex
	}

	q.currentIndex = nextIndex
	q.waitingForWrite = !hasValue
	return hasValue, nil
}

func (q *queue[T]) Value() (T, error) {
	// The read index might be out of bounds at this point (either outside the slice, or the ring)
	// and we should not return a value here if that is the case.
	if q.waitingForWrite {
		var zero T
		return zero, nil
	}
	return q.values[q.currentIndex], nil
}

func (q *queue[T]) Reset() {
	q.values = []T{}
	q.currentIndex = -1
	q.lastSetIndex = -1
	q.zeroIndexSet = false
	q.waitingForWrite = false
}

func (q *queue[T]) Size() int {
	return len(q.values)
}
