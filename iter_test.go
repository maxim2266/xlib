package xlib

import (
	"slices"
	"testing"
)

func TestMapIt(t *testing.T) {
	src := []int{1, 2, 3, 4}
	res := slices.Collect(MapIt(slices.Values(src), func(x int) int { return 2 * x }))

	if !slices.Equal(res, []int{2, 4, 6, 8}) {
		t.Errorf("invalid result: %v", res)
		return
	}
}

func TestFilterIt(t *testing.T) {
	src := []int{1, 2, 3, 4}
	res := slices.Collect(FilterIt(slices.Values(src), func(x int) bool { return x&1 != 0 }))

	if !slices.Equal(res, []int{1, 3}) {
		t.Errorf("invalid result: %v", res)
		return
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
		t.Error("(1) arrays don't match")
		return
	}

	// early exit
	res = res[:0]

	for v := range PipeIt(slices.Values(src)) {
		if v >= N/2 {
			break
		}

		res = append(res, v)
	}

	if !slices.Equal(res, src[:N/2]) {
		t.Error("(2) arrays don't match")
		return
	}
}
