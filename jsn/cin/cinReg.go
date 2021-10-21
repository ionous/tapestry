package cin

import (
	"hash/fnv"
	"io"
	r "reflect"
)

func findType(sig uint64, regs []map[uint64]interface{}) (ret interface{}, okay bool) {
	for _, reg := range regs {
		if a, ok := reg[sig]; ok {
			ret, okay = a, ok
			break
		}
	}
	return
}

func newFromType(cmd interface{}) interface{} {
	return r.New(r.TypeOf(cmd).Elem()).Interface()
}

func hash(k string) uint64 {
	w := fnv.New64a()
	io.WriteString(w, k)
	return w.Sum64()
}
