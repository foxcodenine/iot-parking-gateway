package helpers

// Splice function to remove and optionally add elements to a slice.
func Splice[T any](slice []T, start, count int, elementsToAdd []T) ([]T, []T) {
	// Ensure `start` is within valid bounds; if not, return the original slice and an empty slice.
	if start < 0 || start > len(slice) {
		return []T{}, slice
	}

	// If `count` exceeds the bounds, adjust it to avoid out-of-range issues.
	if start+count > len(slice) {
		count = len(slice) - start
	}

	// Capture the removed elements.
	removed := make([]T, count)
	copy(removed, slice[start:start+count])

	// Create the new slice with removed elements and added elements, if any.
	newSlice := append(slice[:start], append(elementsToAdd, slice[start+count:]...)...)

	return removed, newSlice
}

// contains checks if a slice contains a particular element.
func Contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
