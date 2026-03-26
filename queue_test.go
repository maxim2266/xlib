package xlib

import (
	"math/bits"
	"testing"
)

func TestQueueSimple(t *testing.T) {
	const N = 50

	q := QueueOf[int](0)

	for n := 1; n <= N; n++ {
		for i := range n {
			q.Push(i)
		}

		for i := range n {
			x, ok := q.Pop()

			if !ok {
				t.Fatalf("(%d, %d) empty queue", i, n)
			}

			if x != i {
				t.Fatalf("(%d, %d) unexpected value %d", i, n, x)
			}
		}

		if !q.IsEmpty() {
			t.Fatalf("%d: non-empty queue", n)
		}
	}
}

func TestQueueFromSlice(t *testing.T) {
	const N = 32

	src := make([]int, N)

	for i := range len(src) {
		src[i] = i
	}

	q := QueueFromSlice(src)

	if len(q.buff) != 1<<bits.Len(N) {
		t.Fatalf("unexpected queue size %d", len(q.buff))
	}

	for i := range len(src) {
		x, ok := q.Pop()

		if !ok {
			t.Fatalf("(%d) empty queue", i)
		}

		if x != i {
			t.Fatalf("(%d) unexpected value %d", i, x)
		}
	}

	if !q.IsEmpty() {
		t.Fatal("non-empty queue")
	}
}

func TestQueueIter(t *testing.T) {
	q := QueueFromArgs(1, 2, 3, 4)
	i := 0

	for v := range q.All {
		if i++; v != i {
			if v == -1 {
				break
			}

			t.Fatalf("unexpected value: %d instead of %d", v, i)
		}

		q.Push(-i)
	}

	i = -1

	for v := range q.All {
		if i--; v != i {
			t.Fatalf("unexpected value: %d instead of %d", v, i)
		}
	}

	if i != -4 {
		t.Fatalf("strange last value %d", i)
	}
}
