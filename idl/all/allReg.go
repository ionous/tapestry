package all

import (
	core "git.sr.ht/~ionous/iffy/idl/core"
	"git.sr.ht/~ionous/iffy/idl/reg"
	"git.sr.ht/~ionous/iffy/idl/rtx"
	capnp "zombiezen.com/go/capnproto2"
)

func RegisterTypes() {
	// fix? or maybe this should be an array of pairs and validate they are all in-order
	// fix? i dont like that are having to register all the signatures, id rather register all the types
	// could we look up the union from the type map maybe?
	reg.RegisterMap(rtx.BoolEval_TypeID, reg.EvalMap{
		// BoolEvalImpl_Which_allOf           BoolEvalImpl_Which = 0
		// BoolEvalImpl_Which_always          BoolEvalImpl_Which = 1
		// BoolEvalImpl_Which_anyOf           BoolEvalImpl_Which = 2
		// BoolEvalImpl_Which_bool            BoolEvalImpl_Which = 3
		// BoolEvalImpl_Which_cmpIsNum        BoolEvalImpl_Which = 4
		// BoolEvalImpl_Which_cmpIsTxt        BoolEvalImpl_Which = 5
		// BoolEvalImpl_Which_containsPart    BoolEvalImpl_Which = 6
		// BoolEvalImpl_Which_countOfTrigger  BoolEvalImpl_Which = 7
		// BoolEvalImpl_Which_determineArgs   BoolEvalImpl_Which = 8
		// BoolEvalImpl_Which_during          BoolEvalImpl_Which = 9
		// BoolEvalImpl_Which_findList        BoolEvalImpl_Which = 10
		// BoolEvalImpl_Which_getObjTrait     BoolEvalImpl_Which = 11
		// BoolEvalImpl_Which_getFrom         BoolEvalImpl_Which = 12
		// BoolEvalImpl_Which_hasDominion     BoolEvalImpl_Which = 13
		// BoolEvalImpl_Which_isEmpty         BoolEvalImpl_Which = 14
		// BoolEvalImpl_Which_isValid         BoolEvalImpl_Which = 15
		// BoolEvalImpl_Which_kindOfIs        BoolEvalImpl_Which = 16
		// BoolEvalImpl_Which_kindOfIsExactly BoolEvalImpl_Which = 17
		// BoolEvalImpl_Which_matchesTo       BoolEvalImpl_Which = 18
		// BoolEvalImpl_Which_not             BoolEvalImpl_Which = 19
		// BoolEvalImpl_Which_renderArgs      BoolEvalImpl_Which = 20
		// BoolEvalImpl_Which_renderRefFlags  BoolEvalImpl_Which = 21
		// BoolEvalImpl_Which_sendToArgs      BoolEvalImpl_Which = 22
		// BoolEvalImpl_Which_var             BoolEvalImpl_Which = 23
		uint16(BoolEvalImpl_Which_allOf): func(s capnp.Struct) interface{} {
			return &core.AllTrue{Struct: s}
		},
		uint16(BoolEvalImpl_Which_always): func(s capnp.Struct) interface{} {
			return &core.Always{Struct: s}
		},
		uint16(BoolEvalImpl_Which_anyOf): func(s capnp.Struct) interface{} {
			return &core.AnyTrue{Struct: s}
		},
	})
}
