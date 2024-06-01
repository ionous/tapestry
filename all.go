package tapestry

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/dl/format"
	"git.sr.ht/~ionous/tapestry/dl/game"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/list"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/dl/logic"
	"git.sr.ht/~ionous/tapestry/dl/math"
	"git.sr.ht/~ionous/tapestry/dl/object"
	"git.sr.ht/~ionous/tapestry/dl/prim"
	"git.sr.ht/~ionous/tapestry/dl/rel"
	"git.sr.ht/~ionous/tapestry/dl/render"
	"git.sr.ht/~ionous/tapestry/dl/text"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// all important runtime signatures
var AllSignatures = []map[uint64]typeinfo.Instance{
	assign.Z_Types.Signatures,
	call.Z_Types.Signatures,
	debug.Z_Types.Signatures,
	// frame.Z_Types.Signatures,
	// game.Z_Types.Signatures,
	grammar.Z_Types.Signatures,
	// jess.Z_Types.Signatures,
	list.Z_Types.Signatures,
	literal.Z_Types.Signatures,
	logic.Z_Types.Signatures,
	math.Z_Types.Signatures,
	object.Z_Types.Signatures,
	// play.Z_Types.Signatures,
	prim.Z_Types.Signatures,
	format.Z_Types.Signatures,
	rel.Z_Types.Signatures,
	render.Z_Types.Signatures,
	game.Z_Types.Signatures,
	// rtti.Z_Types.Signatures,
	// story.Z_Types.Signatures,
	// testdl.Z_Types.Signatures,
	text.Z_Types.Signatures,
}
