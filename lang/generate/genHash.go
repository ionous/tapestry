package generate

import (
	"hash/fnv"
	"io"
)

type Hashed struct {
	Value  uint64
	String string
}

// Hash helper for generating signatures lookups
// the slat includes all of the colon separated labels used to indicate a given type.
// the slot specifes the typename the slat fits into.
func Hash(slat, slot string) Hashed {
	str := slat
	if len(slot) > 0 {
		str = slot + "=" + slat
	}
	w := fnv.New64a()
	io.WriteString(w, str)
	return Hashed{Value: w.Sum64(), String: str}
}
