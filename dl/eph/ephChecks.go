package eph

import (
	"sort"

	"github.com/ionous/errutil"
)

func (c *Catalog) WriteChecks(w Writer) (err error) {
	if ds, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
		for _, deps := range ds { // list of dependencies
			d := deps.Leaf().(*Domain) // panics if it fails
			names := make([]string, 0, len(d.checks))
			for k, t := range d.checks {
				if !t.isValidCheck() {
					s := errutil.Sprintf("check %q at %s is missing an expectation or program", t.name, t.at)
					LogWarning(errutil.Error(s)) // sprint the error to avoid triggering errutil.Panic
				} else {
					names = append(names, k)
				}
			}
			sort.Strings(names)
			for _, n := range names {
				t := d.checks[n]
				if e := w.Write(mdl_check, d.name, t.name, t.expect, t.prog, t.at); e != nil {
					err = e
					break
				}
			}
		}
	}
	return
}

// uses directive phase just to be near the end somewhere...
func (op *EphChecks) Phase() Phase { return DirectivePhase }

func (op *EphChecks) Assemble(c *Catalog, d *Domain, at string) (err error) {
	// writes an expectation;
	// not much to verify right now.
	if name, ok := UniformString(op.Name); !ok {
		err = InvalidString(op.Name)
	} else {
		check := d.EnsureCheck(name, at)
		if check.domain != d {
			err = errutil.New("can't extend check %q from another domain", name)
		} else if e := check.setExpectation(op.Expect); e != nil {
			err = e
		} else if e := check.setProg(op.Exe); e != nil {
			err = e
		}
	}
	return
}
