package shared

import (
	"hash/fnv"
)

func Hash(value string) string {
	h := fnv.New32a()
	h.Write([]byte(value))
	return string(h.Sum32())
}
