package xlib

import (
	"slices"
	"testing"
)

func TestMapIt(t *testing.T) {
	src := []int{1, 2, 3, 4}
	res := slices.Collect(MapIt(slices.Values(src), func(x int) int { return 2 * x }))

	if !slices.Equal(res, []int{2, 4, 6, 8}) {
		t.Fatalf("invalid result: %v", res)
	}
}

func TestFilterIt(t *testing.T) {
	src := []int{1, 2, 3, 4}
	res := slices.Collect(FilterIt(slices.Values(src), func(x int) bool { return x&1 != 0 }))

	if !slices.Equal(res, []int{1, 3}) {
		t.Fatalf("invalid result: %v", res)
	}
}

func TestTakeWhileIt(t *testing.T) {
	src := []int{1, 2, 3, 4}
	res := slices.Collect(TakeWhileIt(slices.Values(src), func(x int) bool { return x < 3 }))

	if !slices.Equal(res, []int{1, 2}) {
		t.Fatalf("invalid result: %v", res)
	}
}

func TestPipeIt(t *testing.T) {
	const N = 1000

	src := make([]int, N)

	for i := range N {
		src[i] = i
	}

	res := slices.Collect(PipeIt(slices.Values(src)))

	if !slices.Equal(res, src) {
		t.Fatal("(1) arrays don't match")
	}

	// early exit
	res = slices.Collect(TakeWhileIt(PipeIt(slices.Values(src)), func(x int) bool { return x < N/2 }))

	if !slices.Equal(res, src[:N/2]) {
		t.Fatal("(2) arrays don't match")
	}
}
