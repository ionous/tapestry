package eph

func (op *EphRefs) Phase() Phase { return FieldPhase }

func (op *EphRefs) Assemble(c *Catalog, d *Domain, at string) (err error) {
	println("refs not implemented")
	return
}
