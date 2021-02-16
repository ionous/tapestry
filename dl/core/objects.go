package core

import (
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

type ObjectExists struct {
	Object rt.TextEval `if:"selector=valid,placeholder=object"`
}

// NameOf returns the full name of an object as declared by the author.
// It doesnt change over the course of play. To change the name use the "printed name" property.
type NameOf struct {
	Object rt.TextEval `if:"selector"`
}

// KindOf returns the class of an object.
type KindOf struct {
	Object rt.TextEval
}

// IsKindOf is less about caring, and more about sharing;
// it returns true when the object is compatible with the named kind.
type IsKindOf struct {
	Object rt.TextEval `if:"selector"`
	Kind   string      `if:"selector=is"`
}

type IsExactKindOf struct {
	Object rt.TextEval `if:"selector"`
	Kind   string      `if:"selector=isExactly"`
}

// KindsOf returns all kinds of the specified type
type KindsOf struct {
	Kind string `if:"selector"`
}

func (*ObjectExists) Compose() composer.Spec {
	return composer.Spec{
		Group:  "objects",
		Desc:   "Object Exists: Returns whether there is a object of the specified name.",
		Fluent: &composer.Fluid{Name: "is", Role: composer.Function},
	}
}

func (op *ObjectExists) GetBool(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.ObjectFromText(run, op.Object); obj == nil {
		ret = g.False
	} else if e == nil {
		ret = g.True
	} else if _, isUnknown := e.(g.Unknown); isUnknown {
		ret = g.False
	} else {
		err = e
	}
	return
}

func (*NameOf) Compose() composer.Spec {
	return composer.Spec{
		Group:  "objects",
		Fluent: &composer.Fluid{Name: "nameOf", Role: composer.Function},
		Desc:   "Name Of: Full name of the object.",
	}
}

func (op *NameOf) GetText(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.ObjectFromText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if obj == nil {
		ret = g.Empty // fix: or, should it be "nothing"
	} else if v, e := safe.Unpack(obj, object.Name, affine.Text); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (*KindOf) Compose() composer.Spec {
	return composer.Spec{
		Group: "objects",
		Desc:  "Kind Of: Friendly name of the object's kind.",
		Spec:  "kind of {object:text_eval}",
	}
}

func (op *KindOf) GetText(run rt.Runtime) (ret g.Value, err error) {
	if obj, e := safe.ObjectFromText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if obj == nil {
		ret = g.Empty
	} else if v, e := safe.Unpack(obj, object.Kind, affine.Text); e != nil {
		err = cmdError(op, e)
	} else {
		ret = v
	}
	return
}

func (*IsKindOf) Compose() composer.Spec {
	return composer.Spec{
		Group:  "objects",
		Fluent: &composer.Fluid{Name: "kindOf", Role: composer.Function},
		Desc:   "Is Kind Of: True if the object is compatible with the named kind.",
	}
}

func (op *IsKindOf) GetBool(run rt.Runtime) (ret g.Value, err error) {
	kind := lang.Breakcase(op.Kind)
	if obj, e := safe.ObjectFromText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else {
		ok := safe.Compatible(obj, kind, false)
		ret = g.BoolOf(ok)
	}
	return
}

func (*IsExactKindOf) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "kindOf", Role: composer.Function},
		Group:  "objects",
		Desc:   "Is Kind Of: True if the object is compatible with the named kind.",
	}
}

func (op *IsExactKindOf) GetBool(run rt.Runtime) (ret g.Value, err error) {
	kind := lang.Breakcase(op.Kind)
	if obj, e := safe.ObjectFromText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else {
		ok := safe.Compatible(obj, kind, true)
		ret = g.BoolOf(ok)
	}
	return
}

func (*KindsOf) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "kindsOf", Role: composer.Function},
		Group:  "objects",
		Desc:   "Kinds Of: A list of compatible kinds.",
	}
}

func (op *KindsOf) GetTextList(run rt.Runtime) (g.Value, error) {
	kind := lang.Breakcase(op.Kind) // fix: assembly time.
	return run.GetField(object.Nouns, kind)
}
