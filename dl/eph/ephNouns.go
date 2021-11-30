package eph

import (
	"github.com/ionous/errutil"
)

type NounError struct {
	Noun string
	Err  error
}

func (n NounError) Error() string {
	return errutil.Sprintf("%v for noun %q", n.Err, n.Noun)
}
func (n NounError) Unwrap() error {
	return n.Err
}

// name of a noun to assembly info
func (el *EphNouns) Phase() Phase { return NounPhase }

// noun, kind
func (el *EphNouns) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if name, ok := UniformString(el.Noun); !ok {
		err = InvalidString(el.Noun)
	} else if k, ok := UniformString(el.Kind); !ok {
		err = InvalidString(el.Kind)
	} else {
		noun := d.EnsureNoun(name, at)
		// we can only add requirements to the noun in the same domain that it was declared
		// if in a different domain: the nouns have to match up
		if noun.domain == d {
			noun.AddRequirement(k)
		} else if ok, e := noun.HasAncestor(k); e != nil {
			err = e
		} else if !ok {
			err = NounError{name, errutil.Fmt("can't redefine parent as %q", k)}
		} else {
			e := errutil.New("duplicate noun definition at", at)
			LogWarning(NounError{name, e})
		}
	}
	return
}
