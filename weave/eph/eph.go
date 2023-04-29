package eph

import "git.sr.ht/~ionous/tapestry/weave/assert"

// implemented by individual commands
type Ephemera interface {
	Assert(assert.Assertions) error
}

type PrintOnce string

func (p *PrintOnce) PrintOnce() {
	if *p != "" {
		println(*p)
		*p = ""
	}
}
