package express

import (
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// maps cmd spec name to a nil pointer to the cmd type in question
// ( nil ptr can be used for reflecting on types )
type nameCache struct {
	els map[string]interface{}
}

// the singleton cache of core commands
var coreCache nameCache

// returns the null pointer
func (k *nameCache) get(n string) (ret interface{}, okay bool) {
	if len(k.els) == 0 {
		m := make(map[string]interface{})
		// tbd: add the pointers to flow, or a list to the generated type list?
		// ( see also: unblock.MakeBlockCreator )
		for _, ptr := range core.Z_Types.Signatures {
			i, _ := ptr.(typeinfo.Inspector).Inspect()
			m[i.TypeName()] = ptr // note: multiple signatures can generate the same type
		}
		k.els = m
	}
	ret, okay = k.els[n]
	return
}
