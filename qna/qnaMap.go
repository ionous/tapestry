package qna

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/tables"
)

// hold potential values.
// the meaning of a snapshot changes per value type.
// ex. snapshots from evals are unique instances,
// while snapshots of the same primitive list share the same memory.
// the interface mirrors safe.GetAssignedValue.
type valueMap map[keyType]rt.Assignment

// record a permanent error
func (m *valueMap) storeError(key keyType, err error) rt.Assignment {
	val := errorValue{err}
	(*m)[key] = val
	return val
}

// where i is from the db
func decodeValue(a affine.Affinity, i interface{}) (ret rt.Assignment, err error) {
	if prog, ok := i.([]byte); !ok {
		ret = staticValue{a, i}
	} else {
		// gob needs a wrapper structure (really a field?) for writing interfaces
		// the writer ( ex. assembly WritePattern ) has a matching anonymous struct
		meta := struct{ Init rt.Assignment }{}
		if e := tables.DecodeGob(prog, &meta); e != nil {
			err = e
		} else {
			ret = meta.Init
		}
	}
	return
}
