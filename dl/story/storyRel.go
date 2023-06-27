package story

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/weave"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *DefineRelation) Execute(macro rt.Runtime) error {
	return Weave(macro, op)
}

func (op *DefineRelation) Weave(cat *weave.Catalog) error {
	return cat.Schedule(assert.RequireAncestry, func(w *weave.Weaver) (err error) {
		if rel, e := safe.GetText(w, op.Relation); e != nil {
			err = e
		} else {
			err = op.Cardinality.DefineRelation(w, cat, rel.String())
		}
		return
	})
}

func (op *RelationCardinality) DefineRelation(run rt.Runtime, k assert.Assertions, rel string) (err error) {
	type RelationDefiner interface {
		DefineRelation(rt.Runtime, assert.Assertions, string) error
	}
	if c, ok := op.Value.(RelationDefiner); !ok {
		err = ImportError(op, errutil.Fmt("%w for %T", UnhandledSwap, op.Value))
	} else {
		err = c.DefineRelation(run, k, rel)
	}
	return
}

func (op *OneToOne) DefineRelation(run rt.Runtime, k assert.Assertions, rel string) (err error) {
	if a, e := safe.GetText(run, op.Kind); e != nil {
		err = e
	} else if b, e := safe.GetText(run, op.OtherKind); e != nil {
		err = e
	} else {
		err = k.AssertRelation(rel, a.String(), b.String(), false, false)
	}
	return
}

func (op *OneToMany) DefineRelation(run rt.Runtime, k assert.Assertions, rel string) (err error) {
	if a, e := safe.GetText(run, op.Kind); e != nil {
		err = e
	} else if b, e := safe.GetText(run, op.Kinds); e != nil {
		err = e
	} else {
		err = k.AssertRelation(rel, a.String(), b.String(), false, true)
	}
	return
}

func (op *ManyToOne) DefineRelation(run rt.Runtime, k assert.Assertions, rel string) (err error) {
	if a, e := safe.GetText(run, op.Kind); e != nil {
		err = e
	} else if b, e := safe.GetText(run, op.Kinds); e != nil {
		err = e
	} else {
		err = k.AssertRelation(rel, a.String(), b.String(), true, false)
	}
	return
}

func (op *ManyToMany) DefineRelation(run rt.Runtime, k assert.Assertions, rel string) (err error) {
	if a, e := safe.GetText(run, op.Kinds); e != nil {
		err = e
	} else if b, e := safe.GetText(run, op.OtherKinds); e != nil {
		err = e
	} else {
		err = k.AssertRelation(rel, a.String(), b.String(), true, true)
	}
	return
}
