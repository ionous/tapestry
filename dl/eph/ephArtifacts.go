package eph

import (
	"errors"

	"github.com/ionous/errutil"
)

// key to value and source of declaration.
type Artifacts map[string]Definition

type Definition struct {
	at, value string
}

type Conflict struct {
	Key    string
	Reason ReasonForConflict
	Was    Definition
	Value  string
}

func (n *Conflict) Error() string {
	return errutil.Sprintf("%s %s was %q at %s now %q", n.Reason.String(), n.Key, n.Was.value, n.Was.at, n.Value)
}

func newConflict(key string, why ReasonForConflict, was Definition, newval string) *Conflict {
	return &Conflict{key, why, was, newval}
}

// fix? consider moving domain error to catalog processing internals ( and removing explicit external use )
type DomainError struct {
	Domain string
	Err    error
}

func (n DomainError) Error() string {
	return errutil.Sprintf("%v in domain %q", n.Err, n.Domain)
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

func (defs Artifacts) Merge(from Artifacts, allowDupes bool) (err error) {
	for k, def := range from {
		var conflict *Conflict
		if e := defs.CheckConflict(k, def.value); e == nil {
			defs[k] = def // store if there was no conflict
		} else if allowDupes && errors.As(e, &conflict) && conflict.Reason == Duplicated {
			LogWarning(e) // warn if it was a duplicated definition
		} else {
			err = errutil.Append(err, e)
		}
	}
	return
}

func (defs Artifacts) CheckConflict(key, value string) (err error) {
	if def, ok := defs[key]; ok {
		var why ReasonForConflict
		if def.value == value {
			why = Duplicated // if its duplicated, the previous entry would have checked for redefined
		} else {
			why = Redefined
		}
		err = newConflict(key, why, def, value)
	}
	return
}
