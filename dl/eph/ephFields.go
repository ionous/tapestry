package eph

import (
	"errors"
	"strings"

	"git.sr.ht/~ionous/iffy/affine"
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

func (op *EphFields) Phase() Phase { return FieldPhase }

// add some fields to a kind.
// see also: EphAspects which generates traits and adds them to a custom aspect kind.
func (op *EphFields) Assemble(c *Catalog, d *Domain, at string) (err error) {
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
}

func MakeUniformField(op EphParams) (ret UniformField, err error) {
	if name, ok := UniformString(op.Name); !ok {
		err = InvalidString(op.Name)
	} else if aff, ok := composer.FindChoice(&op.Affinity, op.Affinity.Str); !ok && len(op.Affinity.Str) > 0 {
		err = errutil.New("unknown affinity", aff)
	} else if class, ok := UniformString(op.Class); !ok && len(op.Class) > 0 {
		err = InvalidString(op.Class)
	} else {
		ret = UniformField{name, aff, class}
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
			name: op.name, affinity: op.affinity, class: op.class, at: at,
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
