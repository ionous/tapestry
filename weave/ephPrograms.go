package weave

import (
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

// mdl.Prog, lede, "Directive", str
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
	for _, def := range d.defs {
		if vs := def.key.vals; vs[0] == "prog" {
			name := vs[1]
			if e := w.Write(mdl.Grammar, d.name, name, def.value, def.at); e != nil {
				err = errutil.Append(err, e)
			}
		}
	}
	return
}

// jump/skip/hop	{"Directive:scans:":[["jump","skip","hop"],[{"As:":"jumping"}]]}
func (cat *Catalog) AssertGrammar(opName string, opDirective *grammar.Directive) error {
	return cat.Schedule(assert.DirectivePhase, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		// fix: definitions probably need to be smarter.
		if str, e := marshalout(opDirective); e != nil {
			err = e
		} else {
			err = d.AddDefinition(MakeKey("prog", opName), at, str)
		}
		return
	})
}
