package eph

import (
	"git.sr.ht/~ionous/tapestry/imp/assert"
	"github.com/ionous/errutil"
)

type EphDefinition struct {
	Path  []string
	Value string
}

func (op *EphDefinition) Phase() assert.Phase { return assert.RefPhase }

func (op *EphDefinition) Weave(k assert.Assertions) (err error) {
	path := append(op.Path, op.Value)
	return k.AssertDefinition(path...)
}

func (ctx *Context) AssertDefinition(path ...string) (err error) {
	d, at := ctx.d, ctx.at
	if end := len(path) - 1; end <= 0 {
		err = errutil.New("path too short", path)
	} else {
		path, value := path[:end], path[end]
		err = d.AddDefinition(MakeKey(path...), at, value)
	}
	return

}

func (op *EphBeginDomain) Phase() assert.Phase { return assert.DomainPhase }

func (op *EphBeginDomain) Weave(k assert.Assertions) (err error) {
	return k.BeginDomain(op.Name, op.Requires)
}

func (op *EphEndDomain) Phase() assert.Phase { return assert.DomainPhase }

func (op *EphEndDomain) Weave(k assert.Assertions) (err error) {
	return k.EndDomain()
}
