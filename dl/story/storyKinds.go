package story

import (
	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/imp"
)

// ex. colors are a kind of value
func (op *KindsOfAspect) PostImport(k *imp.Importer) (err error) {
	// fix: is this even useful? see EphAspects.Assemble which has to work around the empty traits list.
	k.WriteEphemera(&eph.EphAspects{Aspects: op.Aspect.Str})
	return
}

// ex. "cats are a kind of animal"
func (op *KindsOfKind) PostImport(k *imp.Importer) (err error) {
	k.WriteEphemera(&eph.EphKinds{Kinds: op.PluralKinds.Str, From: op.SingularKind.Str})
	return
}

// ex. cats have some text called breed.
// ex. horses have an aspect called speed.
func (op *KindsHaveProperties) PostImport(k *imp.Importer) (err error) {
	if len(op.Props) > 0 {
		var ps []eph.EphParams
		for _, el := range op.Props {
			// bool fields become implicit aspects
			// ( vs. bool pattern vars which stay bools -- see reduceProps )
			if p := el.GetParam(); p.Affinity.Str != eph.Affinity_Bool {
				ps = append(ps, p)
			} else {
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
		k.WriteEphemera(&eph.EphKinds{Kinds: op.PluralKinds.Str, Contain: ps})
	}
	return
}
