package decoder

import (
	"fmt"
	"slices"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
)

type KindBuilder interface {
	KindOfAncestors(string) ([]string, error)
	FieldsOf(string) ([]rt.Field, error)
}

func BuildKind(kb KindBuilder, k string) (ret rt.Kind, err error) {
	if path, e := kb.KindOfAncestors(k); e != nil {
		err = fmt.Errorf("%w while building kind %q", e, k)
	} else if len(path) == 0 {
		err = fmt.Errorf("invalid kind %q", k)
	} else {
		k := path[0] // use the returned name in favor of the given name (ex. plurals)
		if fields, e := kb.FieldsOf(k); e != nil {
			err = fmt.Errorf("%w while building kind %q", e, k)
		} else if a, e := buildObjectAspects(kb, path, fields); e != nil {
			err = e
		} else {
			ret = rt.Kind{Path: path, Fields: fields, Aspects: a}
		}
	}
	return
}

// fix? currently only allow traits for objects. hrm.
func buildObjectAspects(kb KindBuilder, path []string, fields []rt.Field) (ret []rt.Aspect, err error) {
	if objectLike := path[len(path)-1] == kindsOf.Kind.String(); objectLike {
		ret, err = BuildAspects(kb, fields)
	}
	return
}

// fix: merge with the rtAsepcts version. :/
func BuildAspects(kb KindBuilder, fields []rt.Field) (ret []rt.Aspect, err error) {
	for _, ft := range fields {
		// tbd? currently a field with the same name and type is an aspect;
		// using string "aspects" might be better...
		// as there would be fewer false positives ( ex. a field of actor called actor )
		// although, it's nice the type is consistently the most derived kind...
		// ( ie. "illumination" is more specific than "aspects" )
		// and some of the db queries would have to change too
		if ft.Affinity == affine.Text && ft.Name == ft.Type {
			if states, e := BuildStates(kb, ft.Type); e != nil {
				err = fmt.Errorf("%w for field %q", e, ft.Name)
				break
			} else if len(states) > 0 {
				a := rt.Aspect{Name: ft.Type, Traits: states}
				ret = append(ret, a)
			}
		}
	}
	return
}

// return a list of state names if the passed kind represents a state set.
// an empty list if it doesnt
// and error if there was some problem determining what is what.
func BuildStates(kb KindBuilder, k string) (ret []string, err error) {
	if path, e := kb.KindOfAncestors(k); e != nil {
		err = fmt.Errorf("%w while building states of %q", e, k)
	} else if slices.Contains(path, k) {
		if fields, e := kb.FieldsOf(k); e != nil {
			err = fmt.Errorf("%w while building states of %q", e, k)
		} else {
			ts := make([]string, len(fields))
			for i, t := range fields {
				ts[i] = t.Name
			}
			ret = ts
		}
	}
	return
}
