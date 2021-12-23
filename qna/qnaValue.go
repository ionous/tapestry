package qna

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/jsn/cin"
	"git.sr.ht/~ionous/iffy/rt"
)

func decodeAssignment(a affine.Affinity, prog []byte, signatures []map[uint64]interface{}) (ret rt.Assignment, err error) {
	if e := cin.Decode(rt.Assignment_Slot{&ret}, prog, signatures); e != nil {
		err = e
	}
	return
}
