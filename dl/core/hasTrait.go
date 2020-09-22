package core

import (
	"github.com/ionous/iffy/dl/composer"
	"github.com/ionous/iffy/rt"
)

// HasTrait a property value from an object by name.
type HasTrait struct {
	Obj   ObjectRef
	Trait rt.TextEval
}

// should be "When the target is publicly named"
func (*HasTrait) Compose() composer.Spec {
	return composer.Spec{
		Name:  "has_trait",
		Spec:  "{object%obj:object_ref} is {trait:text_eval}",
		Group: "objects",
		Desc:  "Has Trait: Return true if noun is currently in the requested state.",
	}
}

func (op *HasTrait) GetBool(run rt.Runtime) (ret bool, err error) {
	if id, e := GetObjectRef(run, op.Obj); e != nil {
		err = cmdError(op, e)
	} else if trait, e := rt.GetText(run, op.Trait); e != nil {
		err = cmdError(op, e)
	} else if p, e := run.GetField(id, trait); e != nil {
		err = cmdError(op, e)
	} else if ok, e := p.GetBool(run); e != nil {
		err = cmdError(op, e)
	} else {
		ret = ok
	}
	return
}
