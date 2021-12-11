package story

import (
	"git.sr.ht/~ionous/iffy/dl/eph"
	"github.com/ionous/errutil"
)

// ex. colors are a kind of value
func (op *KindsOfAspect) ImportPhrase(k *Importer) (err error) {
	k.Write(&eph.EphAspects{Aspects: op.Aspect.Str})
	return
}

// ex. "cats are a kind of animal"
func (op *KindsOfKind) ImportPhrase(k *Importer) (err error) {
	k.Write(&eph.EphKinds{Kinds: op.PluralKinds.Str, From: op.SingularKind.Str})
	return
}

// ex. cats have some text called breed.
// ex. horses have an aspect called speed.
func (op *KindsPossessProperties) ImportPhrase(k *Importer) (err error) {
	for _, n := range op.PropertyDecl {
		if e := n.ImportProperty(k, op.PluralKinds.Str); e != nil {
			err = errutil.Append(err, e)
		}
	}
	return
}
