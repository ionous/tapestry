package reg

import (
	"github.com/ionous/errutil"
	capnp "zombiezen.com/go/capnproto2"
)

func Unpack(wrapperType uint64, ptr capnp.Ptr) (ret interface{}, err error) {
	// look up the set of implementations for the passed wrapper "interface"
	if es, ok := registry[wrapperType]; !ok {
		err = errutil.Fmt("unknown eval %#x", wrapperType)
	} else if s := ptr.Struct(); !s.IsValid() {
		err = errutil.New("invalid value", wrapperType)
	} else {
		// all wrappers have the same format
		// the first u16 is the "which", the first pointer is the unioned value
		which := s.Uint16(0)
		if fn, ok := es[which]; !ok {
			err = errutil.Fmt("unknown value %#x %#x", wrapperType, which)
		} else if p, e := s.Ptr(0); e != nil {
			err = e
		} else {
			// call to the user code to cast to the proper implementation
			ret = fn(p.Struct())
		}
	}
	return
}

func RegisterMap(wrapper uint64, evals EvalMap) {
	// fix handling for tests, etc.
	if registry == nil {
		registry = make(Registry)
	}
	es := registry[wrapper]
	if es == nil {
		es = make(EvalMap)
	}
	for which, eval := range evals {
		es[which] = eval
	}
	registry[wrapper] = es
}

// turn a struct into a method
type EvalToPointer func(capnp.Struct) interface{}

// impl which -> function
type EvalMap map[uint16]EvalToPointer

// typeid -> eval map
type Registry map[uint64]EvalMap

var registry Registry
