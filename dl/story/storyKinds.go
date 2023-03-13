package story

import (
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// ex. "cats are a kind of animal"
func (op *DefineKinds) PostImport(k *imp.Importer) (err error) {
	// FIX: macro runtime
	if kinds, e := safe.GetTextList(nil, op.Kinds); e != nil {
		err = e
	} else if ancestor, e := safe.GetText(nil, op.Ancestor); e != nil {
		err = e
	} else {
		for _, kind := range kinds.Strings() {
			k.WriteEphemera(&eph.EphKinds{Kind: kind, Ancestor: ancestor.String()})
		}
	}
	return
}

// ex. cats have some text called breed.
// ex. horses have an aspect called speed.
func (op *DefineFields) PostImport(k *imp.Importer) (err error) {
	if kind, e := safe.GetText(nil, op.Kind); e != nil {
		err = e
	} else if len(op.Fields) > 0 {
		var ps []eph.EphParams
		for _, el := range op.Fields {
			// bool fields become implicit aspects
			// ( vs. bool pattern vars which stay bools -- see reduceProps )
			if p, ok := el.GetParam(); ok && p.Affinity.Str != eph.Affinity_Bool {
				ps = append(ps, p)
			} else if ok {
				// first: add the aspect
				aspect := p.Name
				traits := []string{"not_" + aspect, "is_" + aspect}
				k.WriteEphemera(&eph.EphAspects{Aspects: aspect, Traits: traits})
				// second: add the field that uses the aspect....
				// fix: future: it'd be nicer to support single trait kinds
				// not_aspect would instead be: Not{IsTrait{PositiveName}}
				ps = append(ps, eph.AspectParam(aspect))
			}
		}

		k.WriteEphemera(&eph.EphKinds{Kind: kind.String(), Contain: ps})
	}
	return
}
