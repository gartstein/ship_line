// Package utils provides common utility functions.
package utils

// CopyMap creates a deep copy of the given map.
// It returns a new map containing all key-value pairs from the original.
func CopyMap[K comparable, V any](original map[K]V) map[K]V {
	newMap := make(map[K]V, len(original))
	for k, v := range original {
		newMap[k] = v
	}
	return newMap
}

// Ptr returns a pointer to the given value.
func Ptr[T any](value T) *T {
	return &value
}
