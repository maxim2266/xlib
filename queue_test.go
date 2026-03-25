package xlib

import (
	"log"
	"testing"
)

func TestQueueSimple(t *testing.T) {
	const N = 50

	q := MakeQueue[int](0)

	for n := 1; n <= N; n++ {
		for i := range n {
			q.Push(i)
		}

		for i := range n {
			x, ok := q.Pop()

			if !ok {
				log.Fatalf("(%d, %d) empty queue", i, n)
			}

			if x != i {
				log.Fatalf("(%d, %d) unexpected value %d", i, n, x)
			}
		}

		if !q.IsEmpty() {
			log.Fatalf("%d: non-empty queue", n)
		}
	}
}
