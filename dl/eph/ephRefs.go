package eph

import "git.sr.ht/~ionous/tapestry/imp/assert"

func (op *EphRefs) Phase() assert.Phase { return assert.RefPhase }

func (op *EphRefs) Weave(k assert.Assertions) (err error) {
	// not implemented.
	return
}

func (op *EphRefs) Assemble(ctx *Context) (err error) {
	refsNotImplemented.PrintOnce()
	return
}

// refs imply some fact about the world that will be defined elsewhere.
// assembly would verify that the referenced thing really exists
var refsNotImplemented PrintOnce = "refs not implemented"

type PrintOnce string

func (p *PrintOnce) PrintOnce() {
	if *p != "" {
		println(*p)
		*p = ""
	}
}
