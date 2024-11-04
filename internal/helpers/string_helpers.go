// string_helpers.go
package helpers

// splitIntoPairs splits a string into pairs of two characters.
func SplitIntoPairs(input string) []string {
	var result []string
	for i := 0; i < len(input); i += 2 {
		if i+2 <= len(input) {
			result = append(result, input[i:i+2])
		}
	}
	return result
}
