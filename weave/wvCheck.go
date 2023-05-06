package weave

import (
	"sort"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/ionous/errutil"
)

type asmChecks map[string]*asmCheck

type asmCheck struct {
	// sometimes test are disabled in the script by renaming bits of them
	// we dont consider than error -- though possible we should warn about it.
	name      string
	domain    *Domain
	expectVal interface{}
	expectAff affine.Affinity
	prog      string
	at        string
}

func (c *asmCheck) setExpectation(v literal.LiteralValue) (err error) {
	if v != nil {
		if len(c.expectAff) > 0 {
			err = errutil.Fmt("check %q cant have multiple expectations", c.name)
		} else if res, e := marshalout(v); e != nil {
			err = e
		} else {
			c.expectVal, c.expectAff = res, v.Affinity()
		}
	}
	return
}

func (c *asmCheck) setProg(exe []rt.Execute) (err error) {
	if len(exe) > 0 {
		if len(c.prog) > 0 {
			err = errutil.Fmt("check %q cant have multiple programs to check", c.name)
		} else {
			slice := rt.Execute_Slice(exe)
			if str, e := marshalout(&slice); e != nil {
				err = e
			} else if len(str) > 0 {
				c.prog = str
			}
		}
	}
	return
}

// return the uniformly named domain ( if it exists )
func (d *Domain) GetCheck(name string) (ret *asmCheck, okay bool) {
	if e := VisitTree(d, func(dep Dependency) (err error) {
		if n, ok := dep.(*Domain).checks[name]; ok {
			ret, okay, err = n, true, Visited
		}
		return
	}); e != nil && e != Visited {
		LogWarning(e)
	}
	return
}

// return the uniformly named domain ( creating it in this domain if necessary )
func (d *Domain) EnsureCheck(name, at string) (ret *asmCheck) {
	if k, ok := d.GetCheck(name); ok {
		ret = k
	} else {
		k = &asmCheck{name: name, at: at, domain: d}
		if d.checks == nil {
			d.checks = map[string]*asmCheck{name: k}
		} else {
			d.checks[name] = k
		}
		ret = k
	}
	return
}

func (c *Catalog) WriteChecks(w Writer) (err error) {
	if ds, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
		for _, d := range ds { // list of dependencies
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
