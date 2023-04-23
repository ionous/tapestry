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

func (op *EphDefinition) Assemble(c *Catalog, d *Domain, at string) (err error) {
	return d.AddDefinition(MakeKey(op.Path...), at, op.Value)
}
