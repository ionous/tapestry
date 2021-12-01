package eph

import (
	"errors"
	"strings"

	"git.sr.ht/~ionous/iffy/dl/composer"
	"github.com/ionous/errutil"
)

//  traverse the domains and then kinds in a reasonable order
func (cat *Catalog) WriteFields(w Writer) (err error) {
	if ds, e := cat.ResolveDomains(); e != nil {
		err = e
	} else {
		for _, dep := range ds {
			d := dep.Leaf().(*Domain)
			if ks, e := d.ResolveKinds(); e != nil {
				err = e
				break
			} else {
				for _, kep := range ks {
					k := kep.Leaf().(*ScopedKind)
					for _, f := range k.fields {
						f.Write(&partialFields{w: w, fields: []interface{}{d.Name(), k.Name()}})
					}
				}
			}
		}
	}
	return
}

func (el *EphFields) Phase() Phase { return FieldPhase }

// add some fields to a kind.
// see also: EphAspects which generates traits and adds them to a custom aspect kind.
func (el *EphFields) Assemble(c *Catalog, d *Domain, at string) (err error) {
	// hooray parameter validation
	// note: the kinds must exist ( and are resolved if they do ) already ( re: phased processing )
	if singleKind, e := c.Singularize(d.name, strings.TrimSpace(el.Kinds)); e != nil {
		err = e
	} else if newKind, ok := UniformString(singleKind); !ok {
		err = InvalidString(el.Kinds)
	} else if kind, ok := d.GetKind(newKind); !ok {
		err = errutil.New("unknown kind", newKind)
	} else if name, ok := UniformString(el.Name); !ok {
		err = InvalidString(el.Name)
	} else if aff, ok := composer.FindChoice(&el.Affinity, el.Affinity.Str); !ok && len(el.Affinity.Str) > 0 {
		err = errutil.New("unknown affinity", aff)
	} else if class, ok := UniformString(el.Class); !ok && len(el.Class) > 0 {
		err = DomainError{d.name, KindError{kind.name, InvalidString(el.Class)}}
	} else if _, ok := d.GetKind(class); !ok && len(class) > 0 {
		err = DomainError{d.name, KindError{kind.name, errutil.New("unknown field class", class)}}
	} else {
		// checks for conflicts, allows duplicates.
		var conflict *Conflict
		if e := kind.AddField(&fieldDef{
			name: name, affinity: aff, class: class, at: at,
		}); errors.As(e, &conflict) && conflict.Reason == Duplicated {
			LogWarning(e) // warn if it was a duplicated definition
		} else {
			err = e // some other error ( or nil )
		}
	}
	return
}
