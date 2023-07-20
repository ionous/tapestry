package story

import (
	"git.sr.ht/~ionous/tapestry/lang"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *DefineRelation) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefineRelation) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weave.RequireAncestry, func(w *weave.Weaver) (err error) {
		if rel, e := safe.GetText(w, op.Relation); e != nil {
			err = e
		} else {
			err = op.Cardinality.DefineRelation(w, rel.String())
		}
		return
	})
}

func (op *RelationCardinality) DefineRelation(w *weave.Weaver, rel string) (err error) {
	type RelationDefiner interface {
		addRelation(*weave.Weaver, string) error
	}
	if c, ok := op.Value.(RelationDefiner); !ok {
		err = ImportError(op, errutil.Fmt("%w for %T", UnhandledSwap, op.Value))
	} else {
		err = c.addRelation(w, rel)
	}
	return
}

func (op *OneToOne) addRelation(w *weave.Weaver, rel string) (err error) {
	if a, e := safe.GetText(w, op.Kind); e != nil {
		err = e
	} else if b, e := safe.GetText(w, op.OtherKind); e != nil {
		err = e
	} else {
		err = addRelation(w.Pin(), rel, a.String(), b.String(), false, false)
	}
	return
}

func (op *OneToMany) addRelation(w *weave.Weaver, rel string) (err error) {
	if a, e := safe.GetText(w, op.Kind); e != nil {
		err = e
	} else if b, e := safe.GetText(w, op.Kinds); e != nil {
		err = e
	} else {
		err = addRelation(w.Pin(), rel, a.String(), b.String(), false, true)
	}
	return
}

func (op *ManyToOne) addRelation(w *weave.Weaver, rel string) (err error) {
	if a, e := safe.GetText(w, op.Kind); e != nil {
		err = e
	} else if b, e := safe.GetText(w, op.Kinds); e != nil {
		err = e
	} else {
		err = addRelation(w.Pin(), rel, a.String(), b.String(), true, false)
	}
	return
}

func (op *ManyToMany) addRelation(w *weave.Weaver, rel string) (err error) {
	if a, e := safe.GetText(w, op.Kinds); e != nil {
		err = e
	} else if b, e := safe.GetText(w, op.OtherKinds); e != nil {
		err = e
	} else {
		err = addRelation(w.Pin(), rel, a.String(), b.String(), true, true)
	}
	return
}

func addRelation(pen *mdl.Pen, rel, a, b string, amany, bmany bool) error {
	rn, an, bn := lang.Normalize(rel), lang.Normalize(a), lang.Normalize(b)
	return pen.AddRelation(rn, an, bn, amany, bmany)
}
