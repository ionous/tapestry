package eph

import (
	"errors"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/composer"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/kindsOf"
	"git.sr.ht/~ionous/iffy/tables/mdl"
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
	Out:
		for _, dep := range deps {
			k := dep.Leaf().(*ScopedKind)
			d := k.domain // for simplicity, fields exist at the scope of the kind: regardless of the scope of the field's declaration.
			p := &partialWriter{w: w, fields: []interface{}{d.Name(), k.Name()}}
			if len(k.fields) > 0 {
				// note: fields might include an aspect field which is capable of storing the active trait ( the trait name text )
				for _, f := range k.fields {
					if e := f.Write(p); e != nil {
						err = e
						break Out
					}
				}
			} else if len(k.aspects) == 1 {
				// if there are no explicit fields; we might be a kind of aspect
				// and all we have are the traits for our aspect.
				if e := k.aspects[0].Write(p); e != nil {
					err = e
					break Out
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
			if k := dep.Leaf().(*ScopedKind); k.HasAncestor(kindsOf.Pattern) {
				for _, fd := range k.fields {
					if init := fd.initially; init != nil {
						if value, e := marshalout(init); e != nil {
							err = e
							break
						} else if e := w.Write(mdl.Assign, k.domain.name, k.name, fd.name, value); e != nil {
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
	// note: the kinds must exist ( and are resolved if they do ) already ( re: phased processing )
	if newKind, ok := UniformString(op.Kinds); !ok {
		err = InvalidString(op.Kinds)
	} else if kind, ok := d.GetPluralKind(newKind); !ok {
		err = errutil.New("unknown kind", newKind)
	} else if param, e := MakeUniformField(op.Affinity, op.Name, op.Class); e != nil {
		err = e
	} else {
		err = param.assembleField(kind, at)
	}
	return
}

type UniformField struct {
	name, affinity, class string
	initially             rt.Assignment
}

// normalize the values of the field
func MakeUniformField(fieldAffinity Affinity, fieldName, fieldClass string) (ret UniformField, err error) {
	if name, ok := UniformString(fieldName); !ok {
		err = InvalidString(fieldName)
	} else if aff, ok := composer.FindChoice(&fieldAffinity, fieldAffinity.Str); !ok && len(fieldAffinity.Str) > 0 {
		err = errutil.New("unknown affinity", aff)
	} else if class, ok := UniformString(fieldClass); !ok && len(fieldClass) > 0 {
		err = InvalidString(fieldClass)
	} else {
		ret = UniformField{name: name, affinity: aff, class: class}
	}
	return
}

// if there's an initial value, make sure it works with our field
func (uf *UniformField) setAssignment(init rt.Assignment) (err error) {
	if init != nil {
		// fix? some statements have unknown affinity ( statements that pivot )
		if initAff := init.Affinity(); len(initAff) > 0 && initAff.String() != uf.affinity {
			err = errutil.Fmt("mismatched affinity of initial value (a %s) for field %q (a %s)", initAff, uf.name, uf.affinity)
		} else {
			uf.initially = init
		}
	}
	return
}

func (uf *UniformField) assembleField(kind *ScopedKind, at string) (err error) {
	if cls, classOk := kind.domain.GetPluralKind(uf.class); !classOk && len(uf.class) > 0 {
		err = KindError{kind.name, errutil.Fmt("unknown class %q for field %q", uf.class, uf.name)}
	} else if aff := affine.Affinity(uf.affinity); classOk && !isClassAffinity(aff) {
		err = KindError{kind.name, errutil.Fmt("unexpected for field %q of class %q", uf.name, uf.class)}
	} else {
		var clsName string
		if classOk {
			clsName = cls.name
		}
		// checks for conflicts, allows duplicates.
		var conflict *Conflict
		if e := kind.AddField(&fieldDef{
			at:        at,
			name:      uf.name, // fieldName; already "uniform"
			affinity:  aff.String(),
			class:     clsName,
			initially: uf.initially,
		}); errors.As(e, &conflict) && conflict.Reason == Duplicated {
			LogWarning(e) // warn if it was a duplicated definition
		} else if e != nil {
			err = e // some other error
		} else {
			// if the field is a kind of aspect
			isAspect := cls != nil && cls.HasParent(kindsOf.Aspect) && len(cls.aspects) > 0
			// when the name of the field is the same as the name of the aspect
			// that is our special "acts as trait" field, so add the set of traits.
			if isAspect && uf.name == clsName && aff == affine.Text {
				err = kind.AddField(&cls.aspects[0])
			}
		}
	}
	return
}

// if there is a class specified, only certain affinities are allowed.
func isClassAffinity(a affine.Affinity) (okay bool) {
	switch a {
	case "", affine.Text, affine.TextList, affine.Record, affine.RecordList:
		okay = true
	}
	return
}
