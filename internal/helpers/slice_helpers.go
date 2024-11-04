package helpers

// Splice replaces elements from index `start` to `start + count` in `slice`
// with the elements in `elementsToAdd`. It returns a new slice with the modifications
// and the removed elements.
func Splice[T any](slice []T, start, count int, elementsToAdd []T) ([]T, []T) {
	// Ensure `start` and `count` do not exceed slice bounds
	if start < 0 || start >= len(slice) {
		return slice, nil // or handle error
	}
	if start+count > len(slice) {
		count = len(slice) - start // trim count if it goes out of bounds
	}

	// Capture the removed part
	removed := make([]T, count)
	copy(removed, slice[start:start+count])

	// Perform the splice operation to create the new slice
	newSlice := append(slice[:start], append(elementsToAdd, slice[start+count:]...)...)

	return removed, newSlice
}
