package story

import (
	"git.sr.ht/~ionous/iffy/dl/eph"
)

// ImportPhrase - action generates pattern ephemera for now.
func (op *ActionDecl) ImportPhrase(k *Importer) (err error) {
	extra := op.ActionParams.Value.(actionImporter).GetExtraParams()
	op.makePattern(k, op.Action.Str, "agent", "actions", extra, nil)
	op.makePattern(k, op.Event.Str, "actor", "events", extra, &eph.EphParams{
		Name:     "success",
		Affinity: eph.Affinity{eph.Affinity_Bool},
	})
	return
}

func (op *ActionDecl) makePattern(k *Importer, name, tgt, sub string, extra []eph.EphParams, res *eph.EphParams) {
	// pattern subtype -- maybe if we really need this an optional parameter of patterns?
	k.Write(&eph.EphKinds{
		Kinds: name,
		From:  sub,
	})
	// the first parameter is always "agent" of type "agent"
	ps := []eph.EphParams{{
		Name:     tgt,
		Affinity: eph.Affinity{eph.Affinity_Text},
		Class:    tgt,
	}}
	k.Write(&eph.EphPatterns{
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