package xlib

import "math/bits"

// Queue is an unbounded FIFO data container with [Push] and [Pop] operations.
type Queue[T any] struct {
	wi, ri int
	buff   []T
}

// QueueOf constructs a new queue of type T and at least of the given size.
func QueueOf[T any](size int) *Queue[T] {
	switch {
	case size > 1024*1024*1024:
		panic("required queue size is too large")
	case size < 16:
		size = 16
	default:
		size = 1 << bits.Len(uint(size-1))
	}

	return &Queue[T]{
		buff: make([]T, size),
	}
}

// QueueFromSlice constructs a new queue initialised with the content of the given slice.
func QueueFromSlice[S ~[]T, T any](src S) *Queue[T] {
	q := QueueOf[T](len(src) + 1)

	for i, v := range src {
		q.buff[i] = v
	}

	q.wi = len(src)
	return q
}

// Push adds an item to the queue.
func (q *Queue[T]) Push(v T) {
	q.buff[q.wi] = v

	m := q.mask()

	if q.wi = (q.wi + 1) & m; q.wi == q.ri {
		// queue is full, resize
		b := make([]T, 2*len(q.buff))

		// copy queue content
		b[0] = q.buff[q.ri]

		for i, j := (q.ri+1)&m, 1; i != q.wi; i, j = (i+1)&m, j+1 {
			b[j] = q.buff[i]
		}

		// new queue
		*q = Queue[T]{
			wi:   len(q.buff),
			buff: b,
		}
	}
}

// Pop removes an item from the queue, returning the item and a flag indicating that
// the queue was not empty before the call.
func (q *Queue[T]) Pop() (v T, ok bool) {
	if ok = !q.IsEmpty(); ok {
		v, q.buff[q.ri] = q.buff[q.ri], v
		q.ri = (q.ri + 1) & q.mask()
	}

	return
}

// IsEmpty returns true if the queue is empty.
func (q *Queue[T]) IsEmpty() bool {
	return q.wi == q.ri
}

func (q *Queue[T]) mask() int {
	return len(q.buff) - 1
}
