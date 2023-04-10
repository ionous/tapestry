package eph

import (
	"git.sr.ht/~ionous/tapestry/imp/assert"
	"sort"

	"git.sr.ht/~ionous/tapestry/tables/mdl"
)

func (c *Catalog) WriteChecks(w Writer) (err error) {
	if ds, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
		for _, deps := range ds { // list of dependencies
			d := deps.Leaf().(*Domain) // panics if it fails
			names := make([]string, 0, len(d.checks))
			for k := range d.checks {
				// re: "expectation" changes, no longer requiring a test statement/ test output for the check to be valid.
				// maybe at runtime we can say "this didn't have any expectations", or "didn't have any program data".
				// if isValid := len(t.expectAff) > 0 && len(t.prog) > 0; isValid {
				// 	s := errutil.Sprintf("check %q at %s is missing an expectation or program", t.name, t.at)
				// 	LogWarning(errutil.Error(s)) // sprint the error to avoid triggering errutil.Panic
				// } else {
				names = append(names, k)
				// }
			}
			sort.Strings(names)
			for _, n := range names {
				t := d.checks[n]
				if e := w.Write(mdl.Check, d.name, t.name, t.expectVal, t.expectAff, t.prog, t.at); e != nil {
					err = e
					break
				}
			}
		}
	}
	return
}

// ensures that a domain exists for the named check
func (op *EphChecks) Phase() assert.Phase { return assert.DomainPhase }

func (op *EphChecks) Assemble(c *Catalog, d *Domain, at string) (err error) {
	// fix. todo: this isnt very well thought out right now --
	// what if a check is part of a story scene? shouldnt it have access to those objects?
	// if checks always establish their own domain, why do they have a duplicate name?
	// there are some checks that have their own scenes, some that have expressions to run, some that have things to check.
	// and others that have just one of those things. ( are there expectations that can simply verify the existence of objects in a model? )
	// etc.
	if name, ok := UniformString(op.Name); !ok {
		err = InvalidString(op.Name)
	} else if d, e := c.EnsureDomain(name, at); e != nil {
		err = e
	} else {
		// uses directive phase just to be near the end somewhere...
		err = d.AddEphemera(at, PhaseFunction{assert.DirectivePhase,
			func(c *Catalog, d *Domain, at string) (err error) {
				check := d.EnsureCheck(name, at)
				if e := check.setExpectation(op.Expect); e != nil {
					err = e
				} else if e := check.setProg(op.Exe); e != nil {
					err = e
				}
				return
			}})
	}
	return
}
