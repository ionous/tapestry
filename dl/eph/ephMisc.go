package eph

import "git.sr.ht/~ionous/tapestry/imp/assert"

type EphDefinition struct {
	Path  []string
	Value string
}

func (op *EphDefinition) Phase() assert.Phase { return assert.RefPhase }

func (op *EphDefinition) Weave(k assert.Assertions) (err error) {
	path := append(op.Path, op.Value)
	return k.AssertDefinition(path...)
}

func (op *EphDefinition) Assemble(ctx *Context) (err error) {
	d, at := ctx.d, ctx.at
	return d.AddDefinition(MakeKey(op.Path...), at, op.Value)
}

func (op *EphBeginDomain) Phase() assert.Phase { return assert.DomainPhase }

func (op *EphBeginDomain) Assemble(*Context) error {
	panic("what should happen here?")
}

func (op *EphBeginDomain) Weave(k assert.Assertions) (err error) {
	return k.BeginDomain(op.Name, op.Requires)
}

func (op *EphEndDomain) Phase() assert.Phase { return assert.DomainPhase }

func (op *EphEndDomain) Assemble(*Context) error {
	panic("what should happen here?")
}

func (op *EphEndDomain) Weave(k assert.Assertions) (err error) {
	return k.EndDomain()
}
