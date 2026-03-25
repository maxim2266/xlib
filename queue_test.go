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
