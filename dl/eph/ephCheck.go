package eph

import (
	"git.sr.ht/~ionous/iffy/dl/literal"
	"github.com/ionous/errutil"
)

type asmChecks map[string]*asmCheck

type asmCheck struct {
	name   string
	domain *Domain
	expect interface{}
	prog   string
	at     string
}

// sometimes test are disabled in the script by renaming bits of them
// we dont consider than error -- though possible we should warn about it.
func (c *asmCheck) isValidCheck() bool {
	return c.expect != nil && len(c.prog) > 0
}

func (c *asmCheck) setExpectation(v literal.LiteralValue) (err error) {
	if v != nil && c.expect != nil {
		err = errutil.New("check %q cant have multiple expectations", c.name)
	} else if expect, e := encodeLiteral(v); e != nil {
		err = e
	} else if expect != nil {
		c.expect = expect
	}
	return
}

func (c *asmCheck) setProg(cmd interface{}) (err error) {
	if cmd != nil && len(c.prog) > 0 {
		err = errutil.New("check %q cant have multiple programs to check", c.name)
	} else if prog, e := marshalout(cmd); e != nil {
		err = e
	} else if len(prog) > 0 {
		c.prog = prog
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
