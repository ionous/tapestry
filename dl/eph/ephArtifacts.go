package eph

import (
	"errors"

	"github.com/ionous/errutil"
)

type DomainArtifacts map[string]Definitions

type Definitions map[string]Definition

type Definition struct {
	at, value string
	err       error
}

type Conflict struct {
	Reason ReasonForConflict
	Def    Definition
}

func (n *Conflict) Error() string {
	return errutil.Sprintf("%s definition at %s value %q", n.Reason.String(), n.Def.at, n.Def.value)
}

type DomainError struct {
	Domain string
	Err    error
}

func (n DomainError) Error() string {
	return errutil.Sprint(n.Err, "in domain", n.Domain)
}
func (n DomainError) Unwrap() error {
	return n.Err
}

type ReasonForConflict int

//go:generate stringer -type=ReasonForConflict
const (
	Redefined ReasonForConflict = iota
	Duplicated
)

// walks the properly cased named domain's dependencies ( non-recursively ) to find
// whether the new key,value pair contradicts or duplicates any existing value.
func (dc DomainArtifacts) CheckConflict(n string, l DependencyFinder, cat, at, key, value string) (err error) {
	fullKey := cat + " " + key
	if e := dc.checkConflict(n, fullKey, value); e != nil {
		err = e
	} else if deps, e := GetResolvedDependencies(n, l); e != nil {
		err = e
	} else {
		// ensure the definition maps exist:
		defs := dc[n]
		if defs == nil {
			defs = make(Definitions)
			dc[n] = defs
		}
		// slow: check for conflicts in all of our ancestors
		for _, depName := range deps.Ancestors(true) {
			if e := dc.checkConflict(depName, fullKey, value); e != nil {
				err = e
				break
			}
		}
		// store the new definition if there was no conflict with existing defs
		if err == nil {
			defs[fullKey] = Definition{at: at, value: value}
		}
	}
	return
}

// was anything stored before?
func (dc DomainArtifacts) checkConflict(n, key, value string) (err error) {
	if def, ok := dc[n]; ok {
		if e := def.checkConflict(key, value); e != nil {
			err = DomainError{n, e}
		}
	}
	return
}

func (defs Definitions) Merge(from Definitions, warn func(e error)) (err error) {
	for k, def := range from {
		var conflict *Conflict
		if e := defs.checkConflict(k, def.value); e == nil {
			defs[k] = def
		} else if errors.As(e, &conflict) && conflict.Reason == Duplicated {
			warn(e)
		} else {
			err = errutil.Append(err, e)
		}
	}
	return
}

func (defs Definitions) checkConflict(key, value string) (err error) {
	if def, ok := defs[key]; ok {
		if def.err != nil {
			err = def.err
		} else {
			var why ReasonForConflict
			if def.value == value {
				why = Duplicated // if its duplicated, the previous entry would have checked for redefined
			} else {
				why = Redefined
			}
			// store the error for next time (if any)
			e := &Conflict{why, def}
			def.err, err = e, e
		}
	}
	return
}
