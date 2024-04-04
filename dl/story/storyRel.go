package story

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

func (op *DefineRelation) Weave(cat *weave.Catalog) error {
	return cat.Schedule(weaver.LanguagePhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if rel, e := safe.GetText(run, op.Relation); e != nil {
			err = e
		} else {
			rel := rel.String()
			switch op.Cardinality {
			case C_RelationCardinality_OneToOne:
				err = op.addOneToOne(cat, run, rel)
			case C_RelationCardinality_OneToMany:
				err = op.addOneToMany(cat, run, rel)
			case C_RelationCardinality_ManyToOne:
				err = op.addManyToOne(cat, run, rel)
			case C_RelationCardinality_ManyToMany:
				err = op.addManyToMany(cat, run, rel)
			}
		}
		return
	})
}

func (op *DefineRelation) addOneToOne(cat *weave.Catalog, run rt.Runtime, rel string) (err error) {
	if a, e := safe.GetText(run, op.Kind); e != nil {
		err = e
	} else if b, e := safe.GetText(run, op.OtherKind); e != nil {
		err = e
	} else {
		err = addRelation(cat, rel, a.String(), b.String(), false, false)
	}
	return
}

func (op *DefineRelation) addOneToMany(cat *weave.Catalog, run rt.Runtime, rel string) (err error) {
	if a, e := safe.GetText(run, op.Kind); e != nil {
		err = e
	} else if b, e := safe.GetText(run, op.OtherKind); e != nil {
		err = e
	} else {
		err = addRelation(cat, rel, a.String(), b.String(), false, true)
	}
	return
}

func (op *DefineRelation) addManyToOne(cat *weave.Catalog, run rt.Runtime, rel string) (err error) {
	if a, e := safe.GetText(run, op.Kind); e != nil {
		err = e
	} else if b, e := safe.GetText(run, op.OtherKind); e != nil {
		err = e
	} else {
		err = addRelation(cat, rel, a.String(), b.String(), true, false)
	}
	return
}

func (op *DefineRelation) addManyToMany(cat *weave.Catalog, run rt.Runtime, rel string) (err error) {
	if a, e := safe.GetText(run, op.Kind); e != nil {
		err = e
	} else if b, e := safe.GetText(run, op.OtherKind); e != nil {
		err = e
	} else {
		err = addRelation(cat, rel, a.String(), b.String(), true, true)
	}
	return
}

func addRelation(cat *weave.Catalog, rel, a, b string, amany, bmany bool) error {
	rn, an, bn := inflect.Normalize(rel), inflect.Normalize(a), inflect.Normalize(b)
	// this can spin a bit until all members of the relation are known
	// having it in ancestry ensures that the kinds will exist for the property phase
	return cat.Schedule(weaver.AncestryPhase, func(w weaver.Weaves, run rt.Runtime) error {
		return w.AddRelation(rn, an, bn, amany, bmany)
	})
}
