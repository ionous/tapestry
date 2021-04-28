package core

import (
	"errors"

	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

// HasTrait - determine if a noun is currently in a particular state.
type HasTrait struct {
	Object rt.TextEval `if:"pb=obj"`
	Trait  rt.TextEval
}

// should be "When the target is publicly named"
func (*HasTrait) Compose() composer.Spec {
	return composer.Spec{
		Lede:  "get",
		Spec:  "{object:text_eval} has {trait:text_eval}",
		Group: "objects",
		Desc:  "Has Trait: Return true if the object is currently in the requested state.",
	}
}

func (op *HasTrait) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if obj := obj.String(); len(obj) == 0 {
		ret = g.False
	} else if trait, e := safe.GetText(run, op.Trait); e != nil {
		err = cmdError(op, e)
	} else {
		trait := trait.String()
		if p, e := run.GetField(obj, trait); e == nil {
			ret = p
		} else if errors.Is(e, g.UnknownField(obj, trait)) {
			ret = g.False
		} else {
			err = cmdError(op, e)
		}
	}
	return
}

// SetTrait a property value from an object by name.
// put obj:trait:
type SetTrait struct {
	Object rt.TextEval `if:"pb=obj"`
	Trait  rt.TextEval
}

// should be "When the target is publicly named"
func (*SetTrait) Compose() composer.Spec {
	return composer.Spec{
		Lede:  "put",
		Spec:  "set {object:text_eval} to {trait:text_eval}",
		Group: "objects",
		Desc:  "Set Trait: put an object into a particular state.",
	}
}

func (op *SetTrait) Execute(run rt.Runtime) (err error) {
	if obj, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if trait, e := safe.GetText(run, op.Trait); e != nil {
		err = cmdError(op, e)
	} else if e := run.SetField(obj.String(), trait.String(), g.BoolOf(true)); e != nil {
		err = cmdError(op, e)
	}
	return
}
