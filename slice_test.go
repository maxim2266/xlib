package xlib

import (
	"slices"
	"testing"
)

func TestCloneSlice(t *testing.T) {
	// non-empty slice
	src := []int{1, 2, 3}
	clone := CloneSlice(src)

	if len(clone) != len(src) {
		t.Errorf("invalid length: %d instead of %d", len(clone), len(src))
		return
	}

	if cap(clone) != len(src) {
		t.Errorf("invalid capacity: %d instead of %d", cap(clone), len(src))
		return
	}

	if !slices.Equal(src, clone) {
		t.Errorf("data mismatch: %v", clone)
		return
	}

	// nil slice
	src = nil
	clone = CloneSlice(src)

	if clone != nil {
		t.Error("non-nil clone")
		return
	}

	// empty slice
	src = []int{}
	clone = CloneSlice(src)

	if clone == nil {
		t.Error("nil clone")
		return
	}

	if len(clone) != 0 {
		t.Errorf("invalid length: %d instead of 0", len(clone))
		return
	}

	if cap(clone) != 0 {
		t.Errorf("invalid capacity: %d instead of 0", cap(clone))
		return
	}
}
