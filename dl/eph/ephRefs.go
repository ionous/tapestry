package eph

func (op *EphRefs) Phase() Phase { return FieldPhase }

func (op *EphRefs) Assemble(c *Catalog, d *Domain, at string) (err error) {
	refsNotImplemented.PrintOnce()
	return
}

var refsNotImplemented PrintOnce = "refs not implemented"

type PrintOnce string

func (p *PrintOnce) PrintOnce() {
	if *p != "" {
		println(*p)
		*p = ""
	}
}
