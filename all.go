package tapestry

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/dl/rel"
	"git.sr.ht/~ionous/tapestry/dl/render"
	"git.sr.ht/~ionous/tapestry/dl/story"
	"git.sr.ht/~ionous/tapestry/rt"
)

var AllSlats = [][]composer.Composer{
	assign.Slats,
	core.Slats,
	debug.Slats,
	grammar.Slats,
	literal.Slats,
	list.Slats,
	prim.Slats,
	rel.Slats,
	render.Slats,
	story.Slats,
}

var AllSignatures = []map[uint64]interface{}{
	assign.Signatures,
	core.Signatures,
	debug.Signatures,
	grammar.Signatures,
	literal.Signatures,
	list.Signatures,
	prim.Signatures,
	rel.Signatures,
	render.Signatures,
	rt.Signatures,
	story.Signatures,
}

var reg composer.TypeRegistry

func Registry() composer.TypeRegistry {
	if reg == nil {
		for _, slats := range AllSlats {
			reg.RegisterTypes(slats)
		}
	}
	return reg
}
