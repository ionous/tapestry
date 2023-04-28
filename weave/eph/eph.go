package eph

import "git.sr.ht/~ionous/tapestry/imp/assert"

// implemented by individual commands
type Ephemera interface {
	Weave(assert.Assertions) error
}

type PrintOnce string

func (p *PrintOnce) PrintOnce() {
	if *p != "" {
		println(*p)
		*p = ""
	}
}
