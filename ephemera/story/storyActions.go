package story

import (
	"git.sr.ht/~ionous/iffy/affine"

	"git.sr.ht/~ionous/iffy/ephemera"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/tables"
)

// ImportPhrase - action generates pattern ephemera for now.
func (op *ActionDecl) ImportPhrase(k *Importer) (err error) {
	if _, e := op.makePattern(k, op.Action.Str, "agent", "actions"); e != nil {
		err = e
	} else if evt, e := op.makePattern(k, op.Event.Str, "actor", "events"); e != nil {
		err = e
	} else {
		// return success
		retName := k.NewName("success", tables.NAMED_RETURN, op.At.String())
		retType := k.NewName("bool_eval", tables.NAMED_TYPE, op.At.String())
		k.NewPatternDecl(evt, retName, retType, "")
	}
	return
}

func (op *ActionDecl) makePattern(k *Importer, name, kind, group string) (ret ephemera.Named, err error) {
	// declare the pattern
	n := k.NewName(lang.Breakcase(name), tables.NAMED_PATTERN, op.At.String())

	// need to declare the group itself at least once
	groupName := k.NewName(group, tables.NAMED_TYPE, op.At.String())
	k.NewPatternDecl(n, n, groupName, "")

	// the first parameter is always "agent"
	paramName := k.NewName(kind, tables.NAMED_PARAMETER, op.At.String())
	paramType := k.NewName(kind, tables.NAMED_KIND, op.At.String())
	k.NewPatternDecl(n, paramName, paramType, affine.Object.String())

	// then the other parameters...
	return n, op.ActionParams.Opt.(actionImporter).ImportAction(k, n)
}

type actionImporter interface {
	ImportAction(*Importer, ephemera.Named) error
}

func (op *CommonAction) ImportAction(k *Importer, n ephemera.Named) (err error) {
	if kind, e := NewSingularKind(k, op.Kind); e != nil {
		err = e
	} else {
		noun := k.NewName(actionNoun, tables.NAMED_PARAMETER, op.At.String())
		k.NewPatternDecl(n, noun, kind, affine.Object.String())
	}
	return
}

func (op *ActionContext) ImportContext(k *Importer, n ephemera.Named) (err error) {
	if kind, e := NewSingularKind(k, op.Kind); e != nil {
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
	if kind, e := FixSingular(k, op.Kinds); e != nil {
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
