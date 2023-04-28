package eph

import (
	"sort"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/imp/assert"
	"git.sr.ht/~ionous/tapestry/rt"

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

func (op *EphChecks) Weave(k assert.Assertions) (err error) {
	return k.AssertCheck(op.Name, op.Exe, op.Expect)
}

func (ctx *Context) AssertCheck(opName string, opExe []rt.Execute, opExpect literal.LiteralValue) (err error) {
	c, at := ctx.c, ctx.at
	// fix. todo: this isnt very well thought out right now --
	// what if a check is part of a story scene? shouldnt it have access to those objects?
	// if checks always establish their own domain, why do they have a duplicate name?
	// there are some checks that have their own scenes, some that have expressions to run, some that have things to check.
	// and others that have just one of those things. ( are there expectations that can simply verify the existence of objects in a model? )
	// etc.
	if name, ok := UniformString(opName); !ok {
		err = InvalidString(opName)
	} else if d, e := c.EnsureDomain(name, at); e != nil {
		err = e
	} else {
		// uses directive phase just to be near the end somewhere...
		err = d.QueueEphemera(at, PhaseFunction{assert.DirectivePhase,
			func(assert.World, assert.Assertions) (err error) {
				check := d.EnsureCheck(name, at)
				if e := check.setExpectation(opExpect); e != nil {
					err = e
				} else if e := check.setProg(opExe); e != nil {
					err = e
				}
				return
			}})
	}
	return
}
