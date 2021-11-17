package eph

import (
	"github.com/ionous/errutil"
	"github.com/ionous/inflect"
)

type Resolution int

const (
	Unresolved Resolution = iota
	Processing
	Resolved
	Errored
)

type Domain struct {
	name, at     string
	originalName string
	inflect      inflect.Ruleset
	deps         DomainList
	eph          EphList
	status       Resolution
	resolved     DomainList
	parents      DomainList
	defines      defines
	err          error
}

type AllDomains map[string]*Domain

// walks the domain's dependencies ( non-recursively ) to find
// whether the new key,value pair contradicts or duplicates any existing value.
func (d *Domain) CheckConflicts(cat, at, key, value string) (err error) {
	if deps, e := d.Resolve(); e != nil {
		err = e
	} else {
		fullKey := cat + " " + key
		if e := checkConflict(d, fullKey, value); e != nil {
			err = e
		} else {
			for _, dep := range deps {
				if e := checkConflict(dep, fullKey, value); e != nil {
					err = e
					break
				}
			}
			if d.defines == nil {
				d.defines = make(defines)
			}
			// store the result: rather than checking through all the resolved domains the next time
			// ( if there is a next time )
			d.defines[fullKey] = Definition{
				at:    at,
				value: value,
				err:   err, // might be null and that's cool.
			}
		}
	}
	return
}

func checkConflict(d *Domain, key, value string) (err error) {
	if def, ok := d.defines[key]; ok {
		if def.err != nil {
			err = def.err
		} else {
			var why ReasonForConflict
			if def.value == value {
				why = Duplicated
			} else {
				why = Redefined
			}
			err = &Conflict{why, d.name, def}
		}
	}
	return
}

func (d *Domain) Resolve() (ret DomainList, err error) {
	if e := d.resolve(); e != nil {
		err = e
	} else {
		ret = d.resolved
	}
	return
}

// Recursively determine the domain's dependency list;
// calling the passed function for each newly resolved dependency
func (d *Domain) resolve() (err error) {
	switch d.status {
	case Resolved:
		// ignore things that are already resolved

	case Processing:
		d.status, d.err = Errored, errutil.New("Circular reference detected:", d.name)
		err = d.err

	case Unresolved:
		d.status = Processing
		var res UniqueNames
		var deps, parents DomainList
		for _, dep := range d.deps {
			// recurse
			if e := dep.resolve(); e != nil {
				d.status, d.err = Errored, errutil.New(e, "->", d.name)
				err = d.err
				break
			} else {
				for _, sub := range dep.resolved {
					if res.AddName(sub.name) {
						deps = append(deps, sub)
					}
				}
				if res.AddName(dep.name) {
					parents = append(parents, dep)
				}
			}
		}
		if err == nil {
			d.status = Resolved
			d.parents = parents
			d.resolved = append(parents, deps...)
		}
	default:
		if e := d.err; e != nil {
			err = e
		} else {
			err = errutil.New("Unknown error processing", d.name)
		}
	}
	return
}
