package weave

import (
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"git.sr.ht/~ionous/tapestry/weave/assert"
	"github.com/ionous/errutil"
)

// note: values are written per *noun* not per domain....
func (cat *Catalog) AssertNounValue(opNoun, opField string, opPath []string, opValue literal.LiteralValue) error {
	return cat.Schedule(assert.ValuePhase, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		if noun, ok := UniformString(opNoun); !ok {
			err = InvalidString(opNoun)
		} else if field, ok := UniformString(opField); !ok {
			err = InvalidString(opField)
		} else if path, e := UniformStrings(opPath); e != nil {
			err = e
		} else if noun, e := d.GetClosestNoun(noun); e != nil {
			err = e
		} else if n, ok := cat.domainNouns[domainNoun{noun.domain, noun.name}]; !ok {
			err = errutil.Fmt("unexpected noun %q in domain %q", noun.name, noun.domain)
		} else if rv, e := n.recordValues(at); e != nil {
			err = e
		} else if value := opValue; value == nil {
			err = errutil.New("null value", opNoun, opField)
		} else {
			return rv.writeValue(noun.name, at, field, path, value)
		}
		return
	})
}

// fix: at some point it'd be nice to write values as they are generated
// the basic idea i think would be to write each field AND sub-record path individually
// and, on write, do a test to ensure the path is meaningful,
// and that no "directory value" value exists for any sub path
// ex. "a.b.c" is okay, so long as there's no record stored at "a.b" directly.
// the runtime would change the way it reconstitutes values to handle all that.
func (cat *Catalog) WriteValues(m *mdl.Modeler) (err error) {
Loop:
	for _, n := range cat.domainNouns {
		if rv := n.localRecord; rv.isValid() {
			for _, fv := range rv.rec.Fields {
				if e := m.Value(n.domain.name, n.name, fv.Field, fv.Value, rv.at); e != nil {
					err = e
					break Loop
				}
			}
		}
	}
	return
}
