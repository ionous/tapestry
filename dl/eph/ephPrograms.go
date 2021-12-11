package eph

import (
	"strings"

	"github.com/ionous/errutil"
)

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
		if i := strings.Index(def.value, ":"); i < 0 {
			e := errutil.New("badly formatted program", def.value)
			err = errutil.Append(err, e)
		} else {
			typeName, prog := def.value[:i], def.value[i+1:]
			//
			if e := w.Write(mdl_prog, k, typeName, prog, def.at); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}

// we give it its own phase so we can keep its definitions separated out.
func (op *EphDirectives) Phase() Phase { return DirectivePhase }

// jump/skip/hop	{"Directive:scans:":[["jump","skip","hop"],[{"As:":"jumping"}]]}
func (op *EphDirectives) Assemble(c *Catalog, d *Domain, at string) (err error) {
	prog := op.Type + ":" + op.Prog // fix: definitions probably need to be smarter.
	return d.AddDefinition(op.Name, at, prog)
}

// fix? the original code decoded the grammar here and invented the lede from it
// that would require us doing some work on the compact reader to separate story dependencies
// so skip force the sender to do the decoding for now.
// func decodeGrammar(op grammar.GrammarMaker, prog string) error {
// 	return cin.Decode(op, prog, iffy.AllSignatures)
// }
