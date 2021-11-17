package eph

import (
	"log"
	"sort"
	"strings"

	"git.sr.ht/~ionous/iffy/lang"
	"github.com/ionous/errutil"
)

// implemented by individual commands
type Ephemera interface {
	Catalog(c *Catalog, d *Domain, at string) error
	// Phase?
}

// Catalog - receives ephemera from the importer.
type Catalog struct {
	Writer     // not a huge fan of this here.... hrm...
	domains    map[string]*Domain
	processing DomainStack
	conflicts  DomainConflicts
}

// helper to make the catalog compatible with the DependencyFinder ( for domains )
type catDependencyFinder Catalog

func (c *catDependencyFinder) GetDependencies(name string) (ret *Dependencies, okay bool) {
	if d, ok := c.domains[name]; ok {
		ret, okay = &d.deps, true
	}
	return
}

func (c *Catalog) Warn(e error) {
	log.Println(e) // for now good enough
}

func (c *Catalog) CheckConflicts(name, cat, at, key, value string) (err error) {
	if c.conflicts == nil {
		c.conflicts = make(DomainConflicts)
	}
	return c.conflicts.CheckConflicts(name, (*catDependencyFinder)(c), cat, at, key, value)
}

// primarily for testing: return list of all the domains that the passed named domain requires
func (c *Catalog) GetDependentDomains(name string) (ret []string, err error) {
	n := lang.Underscore(name)
	if dep, e := c.getDependentDomains(n); e != nil {
		err = e
	} else {
		ret = dep.GetFullTree(true)
	}
	return
}

func (c *Catalog) getDependentDomains(n string) (ret ResolvedDependencies, err error) {
	if d, ok := c.domains[n]; !ok {
		err = errutil.New("unknown domains have no dependencies")
	} else {
		ret, err = d.deps.Resolve(d.name, (*catDependencyFinder)(c))
	}
	return
}

// return the named domain ( creating it if necessary )
func (c *Catalog) GetDomain(name string) (ret *Domain) {
	n := lang.Underscore(name)
	if d, ok := c.domains[n]; ok {
		ret = d
	} else {
		d = &Domain{name: n}
		if c.domains == nil {
			c.domains = map[string]*Domain{n: d}
		} else {
			c.domains[n] = d
		}
		ret = d
	}
	return
}

// creates domains, suspends all other ephemera until the domains are resolved.
func (c *Catalog) AddEphemera(ephAt EphAt) (err error) {
	if d, ok := c.processing.Top(); !ok {
		err = errutil.New("no domain")
	} else {
		switch el := ephAt.Eph.(type) {
		case *EphEndDomain, *EphBeginDomain:
			err = el.Catalog(c, d, ephAt.At)
		default:
			d.eph.All = append(d.eph.All, ephAt)
		}
	}
	return
}

// domains should be in "most" core to least order
// each line should have all the dependencies it needs
func (c *Catalog) WriteDomains(fullTree bool) (err error) {
	sorted := make([]string, 0, len(c.domains))
	for _, d := range c.domains {
		if len(d.at) == 0 {
			err = errutil.Append(err, errutil.New("domain never declared", d.name))
		} else {
			sorted = append(sorted, d.name)
		}
	}
	if err == nil {
		// we *try* as much as possible to keep the order stableish
		sort.Strings(sorted)
		for _, n := range sorted {
			if deps, e := c.getDependentDomains(n); e != nil {
				err = errutil.Append(err, e)
			} else if e := c.Write(mdl_domain, n, strings.Join(deps.GetFullTree(fullTree), ",")); e != nil {
				err = errutil.Append(err, errutil.New("domain", n, "couldn't write", e))
			}
		}
	}
	return
}
