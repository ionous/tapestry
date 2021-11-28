package eph

import (
	"strings"

	"github.com/ionous/errutil"
)

const (
	AspectKinds = "aspect"
	RecordKinds = "record"
)

type KindError struct {
	Kind string
	Err  error
}

func (n KindError) Error() string {
	return errutil.Sprintf("%v in kind %q", n.Err, n.Kind)
}
func (n KindError) Unwrap() error {
	return n.Err
}

// name of a kind to assembly info
func (el *EphKinds) Phase() Phase { return AncestryPhase }

func (el *EphKinds) Assemble(c *Catalog, d *Domain, at string) (err error) {
	if singleKind, e := c.Singularize(d.name, strings.TrimSpace(el.Kinds)); e != nil {
		err = e
	} else if newKind, ok := UniformString(singleKind); !ok {
		err = InvalidString(el.Kinds)
	} else {
		kid := d.EnsureKind(newKind, at)
		if len(el.From) > 0 {
			if parentKind, ok := UniformString(el.From); !ok {
				err = InvalidString(el.From)
			} else {
				kid.AddRequirement(parentKind)
			}
		}
	}
	return
}
