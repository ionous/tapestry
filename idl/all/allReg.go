package all

import (
	core "git.sr.ht/~ionous/iffy/idl/core"
	"git.sr.ht/~ionous/iffy/idl/reg"
	capnp "zombiezen.com/go/capnproto2"
)

func init() {
	reg.Register(BoolEval_TypeID, BoolEvalImpl_Which_always, func(s capnp.Struct) interface{} {
		return &core.Always{Struct: s}
	})
}
