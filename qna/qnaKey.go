package qna

import (
	"hash/fnv"
	"io"
)

type keyType struct {
	val  uint64
	strs []string
}

func makeKey(strs ...string) uint64 {
	w := fnv.New64a()
	for _, str := range strs {
		io.WriteString(w, str)
	}
	return w.Sum64()
}
