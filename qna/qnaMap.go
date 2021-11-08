package qna

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/jsn/cin"
	"git.sr.ht/~ionous/iffy/rt"
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
func decodeValue(a affine.Affinity, i interface{}, signatures []map[uint64]interface{}) (ret rt.Assignment, err error) {
	if prog, ok := i.([]byte); !ok {
		ret = staticValue{a, i}
	} else {
		// the encoder needs a wrapper structure for writing interfaces
		// the writer ( ex. assembly WritePattern ) uses the same Fragment container
		var meta rt.Fragment
		if e := cin.Decode(&meta, prog, signatures); e != nil {
			err = e
		} else {
			ret = meta.Init
		}
	}
	return
}
