package tapestry

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/game"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/dl/rel"
	"git.sr.ht/~ionous/tapestry/dl/render"
)

var AllSignatures = []map[uint64]any{
	assign.Z_Types.Signatures,
	core.Z_Types.Signatures,
	debug.Z_Types.Signatures,
	grammar.Z_Types.Signatures,
	literal.Z_Types.Signatures,
	list.Z_Types.Signatures,
	prim.Z_Types.Signatures,
	rel.Z_Types.Signatures,
	render.Z_Types.Signatures,
	game.Z_Types.Signatures,
}
