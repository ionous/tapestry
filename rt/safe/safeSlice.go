package safe

import (
	"log"

	"git.sr.ht/~ionous/tapestry/rt"
)

func SliceIt(size int, next func(int) (rt.Value, error)) *sliceIt {
	return &sliceIt{size: size, next: next}
}

// panics if the value isnt a list
func ListIt(v rt.Value) (ret *sliceIt) {
	return SliceIt(v.Len(), func(i int) (ret rt.Value, _ error) {
		ret = v.Index(i)
		return
	})
}

type sliceIt struct {
	at   int
	size int
	next func(int) (rt.Value, error)
}

func (it *sliceIt) Remaining() int {
	return it.size - it.at
}

func (it *sliceIt) HasNext() bool {
	return it.at < it.size
}

func (it *sliceIt) GetNext() (ret rt.Value, err error) {
	if !it.HasNext() {
		log.Panicf("iterator out of range %d > %d", it.at, it.size)
	} else if v, e := it.next(it.at); e != nil {
		err = e
	} else {
		ret = v
		it.at++
	}
	return
}
