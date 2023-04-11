package story

import (
	"git.sr.ht/~ionous/tapestry/imp"
	"git.sr.ht/~ionous/tapestry/imp/assert"
	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

// Execute - called by the macro runtime during weave.
func (op *KindOfRelation) Execute(macro rt.Runtime) error {
	return imp.StoryStatement(macro, op)
}

func (op *KindOfRelation) PostImport(k *imp.Importer) error {
	return op.Cardinality.DefineRelation(k, op.Relation.String())
}

type DefineRelation interface {
	DefineRelation(assert.Assertions, string) error
}

func (op *RelationCardinality) DefineRelation(k assert.Assertions, rel string) (err error) {
	if c, ok := op.Value.(DefineRelation); !ok {
		err = ImportError(op, errutil.Fmt("%w for %T", UnhandledSwap, op.Value))
	} else {
		err = c.DefineRelation(k, rel)
	}
	return
}

func (op *OneToOne) DefineRelation(k assert.Assertions, rel string) error {
	return k.AssertRelation(rel, op.Kind.String(), op.OtherKind.String(), false, false)
}
func (op *OneToMany) DefineRelation(k assert.Assertions, rel string) error {
	return k.AssertRelation(rel, op.Kind.String(), op.Kinds.String(), false, true)
}
func (op *ManyToOne) DefineRelation(k assert.Assertions, rel string) error {
	return k.AssertRelation(rel, op.Kind.String(), op.Kinds.String(), true, false)
}
func (op *ManyToMany) DefineRelation(k assert.Assertions, rel string) error {
	return k.AssertRelation(rel, op.Kinds.String(), op.OtherKinds.String(), true, true)
}
