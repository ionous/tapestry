package eph

// mdl_prog, lede, "Directive", str
// fix? each phase seems to be getting its own writer.... should that be formalized?
func (c *Catalog) WriteDirectives(w Writer) (err error) {
	if deps, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
		for _, dep := range deps {
			d := dep.Leaf().(*Domain)
			if e := d.WriteDirectives(w); e != nil {
				err = e
				break
			}
		}
	}
	return
}

func (d *Domain) WriteDirectives(w Writer) (err error) {
	defs := d.phases[DirectivePhase].defs
	for k, def := range defs {
		if e := w.Write(mdl_prog, k, "Directive", def.value, def.at); e != nil {
			err = e
			break
		}
	}
	return
}

// we give it its own phase so we can keep its definitions separated out.
func (el *EphDirectives) Phase() Phase { return DirectivePhase }

// jump/skip/hop	{"Directive:scans:":[["jump","skip","hop"],[{"As:":"jumping"}]]}
func (el *EphDirectives) Assemble(c *Catalog, d *Domain, at string) (err error) {
	// fix? the original code decoded the grammar here and invented the lede from it
	// that would require us doing some work on the compact reader to separate story dependencies
	// so skip force the sender to do the decoding for now.
	return d.AddDefinition(el.Lede, at, el.Prog)
}

// func decodeGrammar(op grammar.GrammarMaker, prog string) error {
// 	return cin.Decode(op, prog, iffy.AllSignatures)
// }
