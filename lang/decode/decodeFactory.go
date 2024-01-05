package decode

import (
	"hash/fnv"
	"io"
	r "reflect"
)

// hash to pointer to a nil value of the type
// ex. 15485098871275255450: (*Comparison)(nil)
type SignatureTable []map[uint64]any

// return the nil pointer to the type.
func (fac SignatureTable) FindType(hash uint64) (ret any, okay bool) {
	for _, reg := range fac {
		if a, ok := reg[hash]; ok {
			ret, okay = a, ok
			break
		}
	}
	return
}

// create a new instance of the type and return its pointer
func (fac SignatureTable) Create(slot, fullsig string) (ret r.Value, okay bool) {
	hash := Hash(slot, fullsig)
	if cmdPtr, ok := fac.FindType(hash); ok {
		cmdType := r.TypeOf(cmdPtr).Elem()
		ret = r.New(cmdType)
		okay = true
	}
	return
}

// Hash helper for generating signatures lookups.
// ex: bool_eval=Not:
// the sig includes all of the colon separated labels used to indicate a given type.
// the slot specifies the typename a command fits into.
func Hash(slot, fullsig string) uint64 {
	w := fnv.New64a()
	if len(slot) > 0 {
		io.WriteString(w, slot)
		w.Write([]byte{'='})
	}
	io.WriteString(w, fullsig)
	return w.Sum64()
}
