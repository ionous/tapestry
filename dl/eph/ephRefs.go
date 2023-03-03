package eph

func (op *EphRefs) Phase() Phase { return RefPhase }

func (op *EphRefs) Assemble(c *Catalog, d *Domain, at string) (err error) {
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
