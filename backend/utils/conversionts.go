package utils

import "strconv"

// ConvertMapKeys converts a map's int keys to string keys.
func ConvertMapKeys[V any](m map[int]V) map[string]V {
	result := make(map[string]V, len(m))
	for k, v := range m {
		result[strconv.Itoa(k)] = v
	}
	return result
}
