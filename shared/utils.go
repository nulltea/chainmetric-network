package shared

import (
	"hash/fnv"
)

func Hash(value string) string {
	h := fnv.New32a()
	h.Write([]byte(value))
	return string(h.Sum32())
}

func ContainsString(value string, values []string) bool {
	contains := false

	if values == nil {
		return false
	}

	for i := range values {
		if values[i] == value {
			contains = true
			break
		}
	}

	return contains
}
