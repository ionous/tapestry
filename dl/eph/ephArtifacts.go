package eph

import (
	"errors"

	"github.com/ionous/errutil"
)

type CategoryKey struct {
	Category string
	Key      string // arbitrary key dependent on the category
}

// key to value and source of declaration.
type Artifacts map[CategoryKey]Definition

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

type ArtifactFinder interface {
	GetArtifacts(collection string) (*Artifacts, bool)
}

// walks the properly cased named domain's dependencies ( non-recursively ) to find
// whether the new key,value pair contradicts or duplicates any existing value.
func CheckConflicts(collections []Dependency, as ArtifactFinder, key CategoryKey, at, value string) (err error) {
	var pf *Artifacts
	for i, cnt := 0, len(collections)-1; i <= cnt; i++ {
		n := collections[cnt-i].Name() // fix?
		if art, ok := as.GetArtifacts(n); !ok {
			err = errutil.New("%q unknown when checking for conflicts ", n)
			break
		} else {
			if pf == nil {
				pf = art
			}
			if e := art.CheckConflict(key, value); e != nil {
				err = DomainError{n, e} // reuse the domain error for now... ( collection name is often domain name )
				break
			}
		}
	}
	if err == nil && pf != nil {
		if *pf == nil {
			*pf = make(Artifacts)
		}
		(*pf)[key] = Definition{at: at, value: value}
	}
	return
}

func (defs Artifacts) Merge(from Artifacts) (err error) {
	for k, def := range from {
		var conflict *Conflict
		if e := defs.CheckConflict(k, def.value); e == nil {
			defs[k] = def // store if there was no conflict
		} else if errors.As(e, &conflict) && conflict.Reason == Duplicated {
			LogWarning(e) // warn if it was a duplicated definition
		} else {
			err = errutil.Append(err, e)
		}
	}
	return
}

func (defs Artifacts) CheckConflict(key CategoryKey, value string) (err error) {
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
