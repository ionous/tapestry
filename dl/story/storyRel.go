package story

import (
	"log"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

func (op *DefineRelation) Weave(cat *weave.Catalog) error {
	return cat.ScheduleCmd(op, weaver.LanguagePhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
		if rel, e := safe.GetText(run, op.RelationName); e != nil {
			err = e
		} else if a, e := safe.GetText(run, op.KindName); e != nil {
			err = e
		} else if b, e := safe.GetText(run, op.OtherKindName); e != nil {
			err = e
		} else {
			amany, bmany := op.GetOneMany()
			rn, an, bn := inflect.Normalize(rel.String()), inflect.Normalize(a.String()), inflect.Normalize(b.String())
			// this can spin a bit until all members of the relation are known
			// having it in ancestry ensures that the kinds will exist for the property phase
			return cat.ScheduleCmd(op, weaver.AncestryPhase, func(w weaver.Weaves, run rt.Runtime) error {
				return w.AddRelation(rn, an, bn, amany, bmany)
			})
		}
		return
	})
}

func (op *DefineRelation) GetOneMany() (lhs, rhs bool) {
	switch op.Cardinality {
	case C_RelationCardinality_OneToOne:
		lhs, rhs = false, false
	case C_RelationCardinality_OneToMany:
		lhs, rhs = false, true
	case C_RelationCardinality_ManyToOne:
		lhs, rhs = true, false
	case C_RelationCardinality_ManyToMany:
		lhs, rhs = true, true
	default:
		// how did this load?
		log.Panicf("unexpected cardinality %s", op.Cardinality)
	}
	return
}
