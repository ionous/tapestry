package story

import (
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
)

// PostImport - action generates pattern ephemera for now.
func (op *ActionDecl) PostImport(k *imp.Importer) (err error) {
	extra := op.ActionParams.Value.(actionImporter).GetExtraParams()
	// the extra "pattern kind" is informational only; the different pattern types dont/shouldnt affect anything.
	op.makePattern(k, op.Action.Str, "agent", kindsOf.Action, extra, nil)
	op.makePattern(k, op.Event.Str, "actor", kindsOf.Event, extra, &eph.EphParams{
		Name:     "success",
		Affinity: eph.Affinity{eph.Affinity_Bool},
	})
	return
}

func (op *ActionDecl) makePattern(k *imp.Importer, name, tgt string, sub kindsOf.Kinds, extra []eph.EphParams, res *eph.EphParams) {
	// pattern subtype -- maybe if we really need this an optional parameter of patterns?
	k.WriteEphemera(&eph.EphKinds{
		Kinds: name,
		From:  sub.String(),
	})
	// the first parameter is always "agent" of type "agent"
	ps := []eph.EphParams{{
		Name:     tgt,
		Affinity: eph.Affinity{eph.Affinity_Text},
		Class:    tgt,
	}}
	k.WriteEphemera(&eph.EphPatterns{
		Name:   name,
		Params: append(ps, extra...),
		Result: res,
	})
}

const actionNoun = "noun"
const actionOtherNoun = "other_noun"

func (op *CommonAction) GetExtraParams() []eph.EphParams {
	return []eph.EphParams{{
		Name:     actionNoun,
		Affinity: eph.Affinity{eph.Affinity_Text},
		Class:    op.Kind.Str,
	}}
}

func (op *PairedAction) GetExtraParams() (ret []eph.EphParams) {
	return []eph.EphParams{{
		Name:     actionNoun,
		Affinity: eph.Affinity{eph.Affinity_Text},
		Class:    op.Kinds.Str,
	}, {
		Name:     actionOtherNoun,
		Affinity: eph.Affinity{eph.Affinity_Text},
		Class:    op.Kinds.Str,
	}}
}

func (op *AbstractAction) GetExtraParams() (ret []eph.EphParams) {
	// no extra parameters
	return
}

// then the other parameters...
type actionImporter interface {
	GetExtraParams() []eph.EphParams
}
