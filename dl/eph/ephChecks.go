package eph

func (op *EphChecks) Phase() Phase { return FieldPhase }

func (op *EphChecks) Assemble(c *Catalog, d *Domain, at string) (err error) {
	println("checks not implemented")
	return
}
