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
	name, at      string
	originalName  string
	inflect       inflect.Ruleset
	deps          DomainList
	eph           EphList
	status        Resolution
	resolvedNames UniqueNames
	resolved      DomainList
	defines       defines
	err           error
}

type AllDomains map[string]*Domain

func (d *Domain) Resolved() []string {
	return []string(d.resolvedNames)
}

// returns:
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
	if e := d.resolveCb(nil); e != nil {
		err = e
	} else {
		ret = d.resolved
	}
	return
}

// Recursively determine the domain's dependency list
func (d *Domain) resolveCb(newlyResolved func(*Domain) error) (err error) {
	switch d.status {
	case Resolved:
		// ignore things that are already resolved

	case Processing:
		d.status, d.err = Errored, errutil.New("Circular reference detected:", d.name)
		err = d.err

	case Unresolved:
		d.status = Processing
		var res UniqueNames
		var deps DomainList
		for _, dep := range d.deps {
			if e := dep.resolveCb(newlyResolved); e != nil {
				d.status, d.err = Errored, errutil.New(e, "->", d.name)
				err = d.err
				break
			} else {
				if res.AddName(dep.name) {
					deps = append(deps, dep)
				}
				for _, sub := range dep.resolved {
					if res.AddName(sub.name) {
						deps = append(deps, sub)
					}
				}
			}
		}
		if err == nil {
			d.status = Resolved
			d.resolvedNames = res
			d.resolved = deps
			//
			if newlyResolved != nil {
				if e := newlyResolved(d); e != nil {
					d.status, d.err = Errored, e
					err = d.err
				}
			}
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
