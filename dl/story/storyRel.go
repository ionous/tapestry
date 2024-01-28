package story

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/mdl"
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
			rel := rel.String()
			switch op.Cardinality {
			case C_RelationCardinality_OneToOne:
				err = op.addOneToOne(w, rel)
			case C_RelationCardinality_OneToMany:
				err = op.addOneToMany(w, rel)
			case C_RelationCardinality_ManyToOne:
				err = op.addManyToOne(w, rel)
			case C_RelationCardinality_ManyToMany:
				err = op.addManyToMany(w, rel)
			}
		}
		return
	})
}

func (op *DefineRelation) addOneToOne(w *weave.Weaver, rel string) (err error) {
	if a, e := safe.GetText(w, op.Kind); e != nil {
		err = e
	} else if b, e := safe.GetText(w, op.OtherKind); e != nil {
		err = e
	} else {
		err = addRelation(w.Pin(), rel, a.String(), b.String(), false, false)
	}
	return
}

func (op *DefineRelation) addOneToMany(w *weave.Weaver, rel string) (err error) {
	if a, e := safe.GetText(w, op.Kind); e != nil {
		err = e
	} else if b, e := safe.GetText(w, op.OtherKind); e != nil {
		err = e
	} else {
		err = addRelation(w.Pin(), rel, a.String(), b.String(), false, true)
	}
	return
}

func (op *DefineRelation) addManyToOne(w *weave.Weaver, rel string) (err error) {
	if a, e := safe.GetText(w, op.Kind); e != nil {
		err = e
	} else if b, e := safe.GetText(w, op.OtherKind); e != nil {
		err = e
	} else {
		err = addRelation(w.Pin(), rel, a.String(), b.String(), true, false)
	}
	return
}

func (op *DefineRelation) addManyToMany(w *weave.Weaver, rel string) (err error) {
	if a, e := safe.GetText(w, op.Kind); e != nil {
		err = e
	} else if b, e := safe.GetText(w, op.OtherKind); e != nil {
		err = e
	} else {
		err = addRelation(w.Pin(), rel, a.String(), b.String(), true, true)
	}
	return
}

func addRelation(pen *mdl.Pen, rel, a, b string, amany, bmany bool) error {
	rn, an, bn := inflect.Normalize(rel), inflect.Normalize(a), inflect.Normalize(b)
	return pen.AddRelation(rn, an, bn, amany, bmany)
}
