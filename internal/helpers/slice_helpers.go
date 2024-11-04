package helpers

// Splice replaces elements from index `start` to `start + count` in `slice`
// with the elements in `elementsToAdd`. It returns a new slice with the modifications
// and the removed elements.
func Splice[T any](slice []T, start, count int, elementsToAdd []T) ([]T, []T) {
	// Check if `start` is within valid bounds; if not, return the original slice and nil.
	// This prevents negative indexing and ensures `start` is not beyond the slice length.
	if start < 0 || start >= len(slice) {
		return slice, nil // Invalid `start` index; no elements removed or modified.
	}

	// Adjust `count` if the range `start` to `start + count` goes out of slice bounds.
	// This ensures we don't exceed the slice length when removing elements.
	if start+count > len(slice) {
		count = len(slice) - start // Set `count` to remove elements only up to the end of the slice.
	}

	// Capture the removed part
	removed := make([]T, count)
	copy(removed, slice[start:start+count])

	// Perform the splice operation to create the new slice
	newSlice := append(slice[:start], append(elementsToAdd, slice[start+count:]...)...)

	return removed, newSlice
}
