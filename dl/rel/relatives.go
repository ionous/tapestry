package rel

import (
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/safe"
	"github.com/ionous/errutil"
)

type RelativeOf struct {
	Via    Relation    `if:"selector"`
	Object rt.TextEval `if:"selector=of"`
}

type RelativesOf struct {
	Via    Relation    `if:"selector"`
	Object rt.TextEval `if:"selector=of"`
}

type ReciprocalOf struct {
	Via    Relation    `if:"selector"`
	Object rt.TextEval `if:"selector=of"`
}

type ReciprocalsOf struct {
	Via    Relation    `if:"selector"`
	Object rt.TextEval `if:"selector=of"`
}

func (*RelativeOf) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "relative", Role: composer.Function},
		Group:  "relations",
		Desc:   "RelativeOf: Returns the relative of a noun (ex. the target of a one-to-one relation.)",
	}
}
func (*RelativesOf) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "relatives", Role: composer.Function},
		Group:  "relations",
		Desc:   "RelativesOf: Returns the relatives of a noun as a list of names (ex. the targets of one-to-many relation).",
	}
}

func (*ReciprocalOf) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "reciprocal", Role: composer.Function},
		Group:  "relations",
		Desc:   "ReciprocalOf: Returns the implied relative of a noun (ex. the source in a one-to-many relation.)",
	}
}

func (*ReciprocalsOf) Compose() composer.Spec {
	return composer.Spec{
		Fluent: &composer.Fluid{Name: "reciprocals", Role: composer.Function},
		Group:  "relations",
		Desc:   "ReciprocalsOf: Returns the implied relative of a noun (ex. the sources of a many-to-many relation.)",
	}
}

func (op *RelativeOf) GetText(run rt.Runtime) (ret g.Value, err error) {
	if a, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if a := a.String(); len(a) == 0 {
		ret = g.Empty
	} else {
		noun, rel := a, op.Via.String()
		if vs, e := run.ReciprocalsOf(noun, rel); e != nil {
			err = cmdError(op, e)
		} else if cnt := len(vs); cnt > 1 {
			e := errutil.New("expected at most one relative for", noun, "in", rel)
			err = cmdError(op, e)
		} else {
			var rel string
			if cnt != 0 {
				rel = vs[0]
			}
			ret = g.StringOf(rel)
		}
	}
	return
}

func (op *RelativesOf) GetTextList(run rt.Runtime) (ret g.Value, err error) {
	if a, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if a := a.String(); len(a) == 0 {
		ret = g.StringsOf(nil)
	} else if vs, e := run.RelativesOf(a, op.Via.String()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.StringsOf(vs)
	}
	return
}

func (op *ReciprocalOf) GetText(run rt.Runtime) (ret g.Value, err error) {
	if a, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if a := a.String(); len(a) == 0 {
		ret = g.Empty
	} else {
		noun, rel := a, op.Via.String()
		if vs, e := run.ReciprocalsOf(noun, rel); e != nil {
			err = cmdError(op, e)
		} else if cnt := len(vs); cnt > 1 {
			e := errutil.New("expected at most one reciprocal for", noun, "in", rel)
			err = cmdError(op, e)
		} else {
			var rel string
			if cnt != 0 { // no relative is considered okay.
				rel = vs[0]
			}
			ret = g.StringOf(rel)
		}
	}
	return
}

func (op *ReciprocalsOf) GetTextList(run rt.Runtime) (ret g.Value, err error) {
	if a, e := safe.ObjectText(run, op.Object); e != nil {
		err = cmdError(op, e)
	} else if a := a.String(); len(a) == 0 {
		ret = g.StringsOf(nil)
	} else if vs, e := run.ReciprocalsOf(a, op.Via.String()); e != nil {
		err = cmdError(op, e)
	} else {
		ret = g.StringsOf(vs)
	}
	return
}
