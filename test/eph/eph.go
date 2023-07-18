package eph

import (
	"git.sr.ht/~ionous/tapestry/weave"
)

// implemented by individual commands
type Ephemera interface {
	Assert(*weave.Catalog) error
}

type PrintOnce string

func (p *PrintOnce) PrintOnce() {
	if *p != "" {
		println(*p)
		*p = ""
	}
}
