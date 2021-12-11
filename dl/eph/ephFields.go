package eph

import (
	"errors"
	"strings"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
	"github.com/ionous/errutil"
)

// eph fields handles the parameter declarations that exist as part of the kinds
type ephFields struct {
	Kinds string
	EphParams
}

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

func (c *Catalog) WriteLocals(w Writer) (err error) {
	if deps, e := c.ResolveKinds(); e != nil {
		err = e
	} else {
		for _, dep := range deps {
			if k := dep.Leaf().(*ScopedKind); k.HasAncestor(KindsOfPattern) {
				for _, fd := range k.fields {
					if init := fd.initially; init != nil {
						if value, e := marshalout(init); e != nil {
							err = e
							break
						} else if e := w.Write(mdl_local, k.domain.name, k.name, fd.name, value); e != nil {
							err = e
							break
						}
					}
				}
			}
		}
	}
	return
}

func (op *ephFields) Phase() Phase { return FieldPhase }

// add some fields to a kind.
// see also: EphAspects which generates traits and adds them to a custom aspect kind.
func (op *ephFields) Assemble(c *Catalog, d *Domain, at string) (err error) {
	// hooray parameter validation
	// note: the kinds must exist ( and are resolved if they do ) already ( re: phased processing )
	if singleKind, e := d.Singularize(strings.TrimSpace(op.Kinds)); e != nil {
		err = e
	} else if newKind, ok := UniformString(singleKind); !ok {
		err = InvalidString(op.Kinds)
	} else if kind, ok := d.GetKind(newKind); !ok {
		err = errutil.New("unknown kind", newKind)
	} else if param, e := MakeUniformField(EphParams{Affinity: op.Affinity, Name: op.Name, Class: op.Class}); e != nil {
		err = e
	} else {
		err = param.AssembleField(kind, at)
	}
	return
}

type UniformField struct {
	name, affinity, class string
	initially             rt.Assignment
}

func MakeUniformField(op EphParams) (ret UniformField, err error) {
	if name, ok := UniformString(op.Name); !ok {
		err = InvalidString(op.Name)
	} else if aff, ok := composer.FindChoice(&op.Affinity, op.Affinity.Str); !ok && len(op.Affinity.Str) > 0 {
		err = errutil.New("unknown affinity", aff)
	} else if class, ok := UniformString(op.Class); !ok && len(op.Class) > 0 {
		err = InvalidString(op.Class)
	} else {
		// if there's an initial value, make sure it works with our field
		if init := op.Initially; init != nil {
			if initAff := init.Affinity(); initAff.String() != aff {
				err = errutil.Fmt("mismatched affinity of initial value (a %s) for field %q (a %s)", initAff, op.Name, aff)
			}
		}
		if err == nil {
			ret = UniformField{name, aff, class, op.Initially}
		}
	}
	return
}

func (op *UniformField) AssembleField(kind *ScopedKind, at string) (err error) {
	if cls, ok := kind.domain.GetKind(op.class); !ok && len(op.class) > 0 {
		err = KindError{kind.name, errutil.Fmt("unknown class %q for field %q", op.class, op.name)}
	} else if ok && (op.affinity != affine.Text.String() && op.affinity != affine.TextList.String()) {
		err = KindError{kind.name, errutil.Fmt("text affinity expected for field %q of class %q", op.name, op.class)}
	} else {

		// checks for conflicts, allows duplicates.
		var conflict *Conflict
		if e := kind.AddField(&fieldDef{
			at:        at,
			name:      op.name,
			affinity:  op.affinity,
			class:     op.class,
			initially: op.initially,
		}); errors.As(e, &conflict) && conflict.Reason == Duplicated {
			LogWarning(e) // warn if it was a duplicated definition
		} else if e != nil {
			err = e // some other error
		} else {
			// if the field is a kind of aspect
			isAspect := cls != nil && cls.HasParent(KindsOfAspect) && len(cls.aspects) > 0
			// when the name of the field is the same as the name of the aspect
			// that is our special "acts as trait" field, so add the set of traits.
			if isAspect && op.name == op.class && op.affinity == affine.Text.String() {
				err = kind.AddField(&cls.aspects[0])
			}
		}
	}
	return
}
