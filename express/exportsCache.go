package express

import (
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/dl/core"
)

// maps cmd spec name to a nil pointer to the cmd type in question
// ( nil ptr can be used for reflecting on types )
type nameCache struct {
	els map[string]interface{}
}

// the singleton cache of core commands
var coreCache nameCache

func (k *nameCache) get(n string) (ret interface{}, okay bool) {
	if len(k.els) == 0 {
		els := make(map[string]interface{})
		for _, v := range core.Slats {
			n := composer.SpecName(v) // ex. all_true
			els[n] = v
		}
		k.els = els
	}
	ret, okay = k.els[n]
	return
}
