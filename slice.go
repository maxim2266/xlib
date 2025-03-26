package xlib

// Clone returns a shallow copy of the slice.
// Unlike [slices.Clone], the resulting slice has the capacity equal to the number
// of elements in the source slice.
func CloneSlice[S ~[]T, T any](src S) (res S) {
	if src != nil {
		res = make([]T, len(src))
		copy(res, src)
	}

	return
}
