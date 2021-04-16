package iffy

import (
	"encoding/gob"

	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/debug"
	"git.sr.ht/~ionous/iffy/dl/grammar"
	"git.sr.ht/~ionous/iffy/dl/list"
	"git.sr.ht/~ionous/iffy/dl/rel"
	"git.sr.ht/~ionous/iffy/dl/render"
	"git.sr.ht/~ionous/iffy/rt"
)

var AllSlots = [][]composer.Slot{rt.Slots, core.Slots, list.Slots, grammar.Slots}
var AllSlats = [][]composer.Composer{core.Slats, debug.Slats, render.Slats, list.Slats, rel.Slats, grammar.Slats}

func RegisterGobs() {
	registerGob()
}

// where should this live?
func init() {
	registerGob()
}

var registeredGob = false

func registerGob() {
	if !registeredGob {
		for _, slats := range AllSlats {
			for _, cmd := range slats {
				gob.Register(cmd)
			}
		}
		registeredGob = true
	}
}
