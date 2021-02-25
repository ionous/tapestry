package story

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/ephemera"
	"git.sr.ht/~ionous/iffy/tables"
)

// ImportPhrase - action generates pattern ephemera for now.
func (op *ActionDecl) ImportPhrase(k *Importer) (err error) {
	// op.Name
	if n, e := op.Name.NewName(k); e != nil {
		err = e
	} else {
		actionType := k.NewName("execute", tables.NAMED_TYPE, op.At.String())
		actor := k.NewName("actor", tables.NAMED_PARAMETER, op.At.String())
		actorKind := k.NewName("actor", tables.NAMED_KIND, op.At.String())

		// need to declare the action itself
		k.NewPatternDecl(n, n, actionType, "")

		// the first parameter is always "actor"
		k.NewPatternDecl(n, actor, actorKind, affine.Object.String())

		// then the other parameters...
		err = op.ActionParams.Opt.(interface {
			ImportAction(*Importer, ephemera.Named) error
		}).ImportAction(k, n)
	}
	return
}

func (op *CommonAction) ImportAction(k *Importer, n ephemera.Named) (err error) {
	if kind, e := op.Kind.NewName(k); e != nil {
		err = e
	} else {
		noun := k.NewName(actionNoun, tables.NAMED_PARAMETER, op.At.String())
		k.NewPatternDecl(n, noun, kind, affine.Object.String())
	}
	return
}

func (op *ActionContext) ImportContext(k *Importer, n ephemera.Named) (err error) {
	if kind, e := op.Kind.NewName(k); e != nil {
		err = e
	} else {
		otherNoun := k.NewName(actionOtherNoun, tables.NAMED_PARAMETER, op.At.String())
		k.NewPatternDecl(n, otherNoun, kind, affine.Object.String())
	}
	return
}

const actionNoun = "noun"
const actionOtherNoun = "other_noun"

func (op *PairedAction) ImportAction(k *Importer, n ephemera.Named) (err error) {
	// inform calls the two objects "noun" and "second noun"
	if kind, e := op.Kinds.FixPlurals(k); e != nil {
		err = e
	} else {
		for _, name := range []string{actionNoun, actionOtherNoun} {
			noun := k.NewName(name, tables.NAMED_PARAMETER, op.At.String())
			k.NewPatternDecl(n, noun, kind, affine.Object.String())
		}
	}
	return
}
func (op *AbstractAction) ImportAction(k *Importer, n ephemera.Named) (err error) {
	// no extra parameters
	return
}
