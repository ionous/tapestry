package eph

import (
	"errors"
	"strings"

	"git.sr.ht/~ionous/iffy/dl/composer"
	"github.com/ionous/errutil"
)

//  traverse the domains and then kinds in a reasonable order
func (c *Catalog) WriteFields(w Writer) (err error) {
	if deps, e := c.ResolveKinds(); e != nil {
		err = e
	} else {
		for _, dep := range deps {
			k := dep.Leaf().(*ScopedKind)
			d := k.domain // for simplicity, fields exist at the  scope of the kind: regardless of the scope of the field's declaration.
			for _, f := range k.fields {
				if e := f.Write(&partialWriter{w: w, fields: []interface{}{d.Name(), k.Name()}}); e != nil {
					err = e
					break
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
	if singleKind, e := d.Singularize(strings.TrimSpace(el.Kinds)); e != nil {
		err = e
	} else if newKind, ok := UniformString(singleKind); !ok {
		err = InvalidString(el.Kinds)
	} else if kind, ok := d.GetKind(newKind); !ok {
		err = errutil.New("unknown kind", newKind)
	} else if param, e := MakeUniformField(EphParams{Affinity: el.Affinity, Name: el.Name, Class: el.Class}); e != nil {
		err = e
	} else if e := param.AssembleField(kind, at); e != nil {
		err = e // hrm.
	}
	return
}

type UniformField struct {
	name, affinity, class string
}

func MakeUniformField(el EphParams) (ret UniformField, err error) {
	if name, ok := UniformString(el.Name); !ok {
		err = InvalidString(el.Name)
	} else if aff, ok := composer.FindChoice(&el.Affinity, el.Affinity.Str); !ok && len(el.Affinity.Str) > 0 {
		err = errutil.New("unknown affinity", aff)
	} else if class, ok := UniformString(el.Class); !ok && len(el.Class) > 0 {
		err = InvalidString(el.Class)
	} else {
		ret = UniformField{name, aff, class}
	}
	return
}

func (el *UniformField) AssembleField(kind *ScopedKind, at string) (err error) {
	if cls, ok := kind.domain.GetKind(el.class); !ok && len(el.class) > 0 {
		err = KindError{kind.name, errutil.Fmt("unknown class %q for field %q", el.class, el.name)}
	} else {
		// checks for conflicts, allows duplicates.
		var conflict *Conflict
		if e := kind.AddField(&fieldDef{
			name: el.name, affinity: el.affinity, class: el.class, at: at,
		}); errors.As(e, &conflict) && conflict.Reason == Duplicated {
			LogWarning(e) // warn if it was a duplicated definition
		} else if e != nil {
			err = e // some other error
		} else if cls != nil && cls.HasParent(KindsOfAspect) && len(cls.aspects) > 0 {
			// if the field is a kind of aspect, then we not only add the aspect as a field
			// we add the set of traits as well
			err = kind.AddField(&cls.aspects[0])
		}
	}
	return
}
