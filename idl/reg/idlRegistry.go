package reg

import (
	"github.com/ionous/errutil"
	capnp "zombiezen.com/go/capnproto2"
)

func Unpack(wrapperType uint64, ptr capnp.Ptr) (ret interface{}, err error) {
	// look up the set of implementations for the passed wrapper "interface"
	if es, ok := registry[wrapperType]; !ok {
		err = errutil.Fmt("unknown eval %0xd", wrapperType)
	} else if s := ptr.Struct(); !s.IsValid() {
		err = errutil.New("invalid value", wrapperType)
	} else {
		// all wrappers have the same format
		// the first u16 is the "which", the first pointer is the unioned value
		which := s.Uint16(0)
		if fn, ok := es[which]; !ok {
			err = errutil.Fmt("unknown value %0xd %0xd", wrapperType, which)
		} else if p, e := s.Ptr(0); e != nil {
			err = e
		} else {
			// call to the user code to cast to the proper implementation
			ret, err = fn(p.Struct())
		}
	}
	return
}

func RegisterMap(wrapper uint64, evals EvalMap) {
	es := registry[wrapper]
	for which, eval := range evals {
		es[which] = eval
	}
	registry[wrapper] = es
}

// turn a struct into a method
type EvalToPointer func(capnp.Struct) (interface{}, error)

// impl which -> function
type EvalMap map[uint16]EvalToPointer

// typeid -> eval map
var registry map[uint64]EvalMap
