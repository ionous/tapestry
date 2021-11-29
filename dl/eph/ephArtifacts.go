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
		err = &Conflict{why, def}
	}
	return
}
