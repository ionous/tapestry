package weave

import (
	"errors"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
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
			k := dep.Leaf().(*ScopedKind)
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
	Name, Type string
	Affinity   affine.Affinity
	Initially  assign.Assignment
	At         string
}

// func (ep *Params) Unify(at string) (UniformField, error) {
// 	return MakeUniformField(ep.Affinity, ep.Name, ep.Class, at)
// }

// normalize the values of the field
func MakeUniformField(aff affine.Affinity, fieldName, fieldClass, at string) (ret UniformField, err error) {
	if name, ok := UniformString(fieldName); !ok {
		err = InvalidString(fieldName)
	} else if class, ok := UniformString(fieldClass); !ok && len(fieldClass) > 0 {
		err = InvalidString(fieldClass)
	} else {
		// shortcut: if we specify a field name for a record and no class, we'll expect the class to be the name.
		if len(class) == 0 && isRecordAffinity(aff) {
			class = name
		}
		ret = UniformField{Name: name, Affinity: aff, Type: class, At: at}
	}
	return
}

// if there's an initial value, make sure it works with our field
func (uf *UniformField) setAssignment(init assign.Assignment) (err error) {
	if init != nil {
		// fix? some statements have unknown affinity ( statements that pivot )
		if initAff := assign.GetAffinity(init); len(initAff) > 0 && initAff != uf.Affinity {
			err = errutil.Fmt("mismatched affinity of initial value (a %s) for field %q (a %s)", initAff, uf.Name, uf.Affinity)
		} else {
			uf.Initially = init
		}
	}
	return
}

func (uf *UniformField) assembleField(kind *ScopedKind) (err error) {
	if cls, classOk := kind.domain.GetPluralKind(uf.Type); !classOk && len(uf.Type) > 0 {
		err = KindError{kind.name, errutil.Fmt("unknown class %q for field %q", uf.Type, uf.Name)}
	} else if aff := uf.Affinity; classOk && !isClassAffinity(aff) {
		err = KindError{kind.name, errutil.Fmt("unexpected for field %q of class %q", uf.Name, uf.Type)}
	} else {
		var clsName string
		if classOk {
			clsName = cls.name
		}
		// checks for conflicts, allows duplicates.
		var conflict *Conflict
		if e := kind.AddField(&fieldDef{
			at:        uf.At,
			name:      uf.Name, // fieldName; already "uniform"
			affinity:  aff.String(),
			class:     clsName,
			initially: uf.Initially,
		}); errors.As(e, &conflict) && conflict.Reason == Duplicated {
			// handle conflicting inits...
			// -- AddField needs refactoring to put this in there easily.
			for i, was := range kind.fields {
				if was.name == uf.Name {
					hadInit := was.initially != nil
					wantsInit := uf.Initially != nil
					switch {
					case wantsInit && !hadInit:
						was.initially = uf.Initially // use the init
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
			if isAspect && uf.Name == clsName && aff == affine.Text {
				err = kind.AddField(&cls.aspects[0])
			}
		}
	}
	return
}

// if there is a class specified, only certain affinities are allowed.
func isRecordAffinity(a affine.Affinity) (okay bool) {
	switch a {
	case affine.Record, affine.RecordList:
		okay = true
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
