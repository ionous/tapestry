package eph

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/composer"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/ionous/errutil"
)

// write the fields of each kind in kind order
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

// write the field initializers in kind order
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

// after queuing up all the fields, assemble them; parent kinds first.
var FieldActions = PhaseAction{
	Do: func(d *Domain) (err error) {
		if deps, e := d.ResolveKinds(); e != nil {
			err = e
		} else {
			for _, dep := range deps {
				k := dep.Leaf().(*ScopedKind)
				for _, p := range k.pendingFields {
					if e := p.assembleField(k); e != nil {
						err = e
						break
					}
				}
			}
		}
		return
	},
}

type UniformField struct {
	name, affinity, class string
	initially             rt.Assignment
	at                    string
}

func (ep *EphParams) unify(at string) (UniformField, error) {
	return MakeUniformField(ep.Affinity, ep.Name, ep.Class, at)
}

// normalize the values of the field
func MakeUniformField(fieldAffinity Affinity, fieldName, fieldClass, at string) (ret UniformField, err error) {
	if name, ok := UniformString(fieldName); !ok {
		err = InvalidString(fieldName)
	} else if aff, ok := composer.FindChoice(&fieldAffinity, fieldAffinity.Str); !ok && len(fieldAffinity.Str) > 0 {
		err = errutil.New("unknown affinity", aff)
	} else if class, ok := UniformString(fieldClass); !ok && len(fieldClass) > 0 {
		err = InvalidString(fieldClass)
	} else {
		ret = UniformField{name: name, affinity: aff, class: class, at: at}
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

func (uf *UniformField) assembleField(kind *ScopedKind) (err error) {
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
			at:        uf.at,
			name:      uf.name, // fieldName; already "uniform"
			affinity:  aff.String(),
			class:     clsName,
			initially: uf.initially,
		}); errors.As(e, &conflict) && conflict.Reason == Duplicated {
			// handle conflicting inits...
			// -- AddField needs refactoring to put this in there easily.
			for i, was := range kind.fields {
				if was.name == uf.name {
					hadInit := was.initially != nil
					wantsInit := uf.initially != nil
					switch {
					case wantsInit && !hadInit:
						was.initially = uf.initially // use the init
						kind.fields[i] = was         // update the list of structs
					case wantsInit && hadInit:
						conflict.Reason = Redefined
						err = conflict // really should wrap this up, but really should fix AddFields
					case !wantsInit && !hadInit:
						LogWarning(e)
					}
					break // out of loop
				}
			}
		} else if e != nil {
			err = e // some other error
		} else {
			// if the field is a kind of aspect
			isAspect := cls != nil && cls.HasParent(kindsOf.Aspect) && len(cls.aspects) > 0
			// when the name of the field is the same as the name of the aspect
			// that is our special "acts as trait" field, so add the set of traits ( to check for conflicts )
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
