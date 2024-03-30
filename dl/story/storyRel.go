package story

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// Execute - called by the macro runtime during weave.
func (op *DefineRelation) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefineRelation) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weaver.MappingPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if rel, e := safe.GetText(run, op.Relation); e != nil {
			err = e
		} else {
			rel := rel.String()
			switch op.Cardinality {
			case C_RelationCardinality_OneToOne:
				err = op.addOneToOne(w, run, rel)
			case C_RelationCardinality_OneToMany:
				err = op.addOneToMany(w, run, rel)
			case C_RelationCardinality_ManyToOne:
				err = op.addManyToOne(w, run, rel)
			case C_RelationCardinality_ManyToMany:
				err = op.addManyToMany(w, run, rel)
			}
		}
		return
	})
}

func (op *DefineRelation) addOneToOne(w weaver.Weaves, run rt.Runtime, rel string) (err error) {
	if a, e := safe.GetText(run, op.Kind); e != nil {
		err = e
	} else if b, e := safe.GetText(run, op.OtherKind); e != nil {
		err = e
	} else {
		err = addRelation(w, rel, a.String(), b.String(), false, false)
	}
	return
}

func (op *DefineRelation) addOneToMany(w weaver.Weaves, run rt.Runtime, rel string) (err error) {
	if a, e := safe.GetText(run, op.Kind); e != nil {
		err = e
	} else if b, e := safe.GetText(run, op.OtherKind); e != nil {
		err = e
	} else {
		err = addRelation(w, rel, a.String(), b.String(), false, true)
	}
	return
}

func (op *DefineRelation) addManyToOne(w weaver.Weaves, run rt.Runtime, rel string) (err error) {
	if a, e := safe.GetText(run, op.Kind); e != nil {
		err = e
	} else if b, e := safe.GetText(run, op.OtherKind); e != nil {
		err = e
	} else {
		err = addRelation(w, rel, a.String(), b.String(), true, false)
	}
	return
}

func (op *DefineRelation) addManyToMany(w weaver.Weaves, run rt.Runtime, rel string) (err error) {
	if a, e := safe.GetText(run, op.Kind); e != nil {
		err = e
	} else if b, e := safe.GetText(run, op.OtherKind); e != nil {
		err = e
	} else {
		err = addRelation(w, rel, a.String(), b.String(), true, true)
	}
	return
}

func addRelation(w weaver.Weaves, rel, a, b string, amany, bmany bool) error {
	rn, an, bn := inflect.Normalize(rel), inflect.Normalize(a), inflect.Normalize(b)
	return w.AddRelation(rn, an, bn, amany, bmany)
}
