package iffy

import (
	"encoding/gob"

	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/dl/core"
	"git.sr.ht/~ionous/iffy/dl/debug"
	"git.sr.ht/~ionous/iffy/dl/grammar"
	"git.sr.ht/~ionous/iffy/dl/list"
	"git.sr.ht/~ionous/iffy/dl/reader"
	"git.sr.ht/~ionous/iffy/dl/rel"
	"git.sr.ht/~ionous/iffy/dl/render"
	"git.sr.ht/~ionous/iffy/dl/value"
	"git.sr.ht/~ionous/iffy/rt"
)

var AllSlots = [][]interface{}{
	core.Slots,
	grammar.Slots,
	list.Slots,
	// story.Slots, // re:story importer, etc.
	rt.Slots,
}
var AllSlats = [][]composer.Composer{
	core.Slats,
	debug.Slats,
	grammar.Slats,
	list.Slats,
	reader.Slats,
	rel.Slats,
	render.Slats,
	value.Slats,
	// story.Slats,
}
var AllSignatures = []map[uint64]interface{}{
	core.Signatures,
	debug.Signatures,
	grammar.Signatures,
	list.Signatures,
	reader.Signatures,
	rel.Signatures,
	render.Signatures,
	value.Signatures,
}

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
