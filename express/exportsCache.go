package express

import (
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/dl/math"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/dl/printer"
	"git.sr.ht/~ionous/tapestry/dl/text"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// maps cmd spec name to a nil pointer to the cmd type in question
// ( nil ptr can be used for reflecting on types )
type nameCache struct {
	els map[string]any
}

// the singleton cache of core commands
var coreCache nameCache

// returns the null pointer
func (k *nameCache) get(n string) (ret any, okay bool) {
	if len(k.els) == 0 {
		m := make(map[string]any)
		// tbd: add the pointers to flow, or a list to the generated type list?
		// ( see also: unblock.MakeBlockCreator )
		addSig(m, logic.Z_Types.Signatures)
		addSig(m, math.Z_Types.Signatures)
		addSig(m, object.Z_Types.Signatures)
		addSig(m, text.Z_Types.Signatures)
		addSig(m, printer.Z_Types.Signatures)
		k.els = m
	}
	ret, okay = k.els[n]
	return
}

func addSig(m map[string]any, sig map[uint64]typeinfo.Instance) {
	for _, ptr := range sig {
		t := ptr.TypeInfo()
		m[t.TypeName()] = ptr // note: multiple signatures can generate the same type
	}
}
