package xlib

type Queue[T any] struct {
	wi, ri int
	buff   []T
}

func MakeQueue[T any]() *Queue[T] {
	return &Queue[T]{
		buff: make([]T, 16),
	}
}

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

func (q *Queue[T]) Pop() (v T, ok bool) {
	if ok = !q.IsEmpty(); ok {
		v = q.buff[q.ri]
		q.ri = (q.ri + 1) & q.mask()
	}

	return
}

func (q *Queue[T]) IsEmpty() bool {
	return q.wi == q.ri
}

func (q *Queue[T]) mask() int {
	return len(q.buff) - 1
}
