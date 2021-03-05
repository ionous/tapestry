package core

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

// HasTrait a property value from an object by name.
type HasTrait struct {
	Object rt.TextEval
	Trait  rt.TextEval
}

// should be "When the target is publicly named"
func (*HasTrait) Compose() composer.Spec {
	return composer.Spec{
		Spec:  "{object:text_eval} has {trait:text_eval}",
		Group: "objects",
		Desc:  "Has Trait: Return true if noun is currently in the requested state.",
	}
}

func (op *HasTrait) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if obj := obj.String(); len(obj) == 0 {
		ret = g.False
	} else if trait, e := safe.GetText(run, op.Trait); e != nil {
		err = cmdError(op, e)
	} else if p, e := run.GetField(obj, trait.String()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = p
	}
	return
}
