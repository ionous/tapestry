package testutil

import (
	"fmt"
	"log"
	r "reflect"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

type Kinds struct {
	Kinds   KindMap // stores completed kinds
	Builder KindBuilder
}

type KindMap map[string]*rt.Kind
type FieldMap map[string]*[]rt.Field // kind name to fields; pointers are used to block recursion
type ParentMap map[string]string

type KindBuilder struct {
	Aspects []rt.Aspect
	Fields  FieldMap
	Parents ParentMap
}

// register kinds from a struct using reflection
// note: this doesnt actually create the kinds....
// it preps them for use -- they are created on Get()
func (ks *Kinds) AddKinds(is ...any) {
	for _, el := range is {
		ks.Builder.addType(ks, r.TypeOf(el).Elem())
	}
}

func (ks *Kinds) AddKind(el any) *rt.Kind {
	t := r.TypeOf(el).Elem()
	ks.Builder.addType(ks, t)
	return ks.Kind(nameOfType(t))
}

func (ks *Kinds) NewRecord(name string, valuePairs ...any) *rt.Record {
	v := rt.NewRecord(ks.Kind(name))
	if len(valuePairs) > 0 {
		if e := SetRecord(v, valuePairs...); e != nil {
			panic(e)
		}
	}
	return v
}

// return a kind ( declared via AddKinds ) panic if it doesn't exist.
func (ks *Kinds) Kind(name string) (ret *rt.Kind) {
	if k, e := ks.GetKindByName(name); e != nil {
		panic(e)
	} else {
		ret = k
	}
	return
}

// return a kind ( declared via AddKinds ) or error if it doesnt exist.
func (ks *Kinds) GetKindByName(name string) (ret *rt.Kind, err error) {
	if k, ok := ks.Kinds[name]; ok {
		ret = k // we created the kind already
	} else {
		if name == kindsOf.Aspect.String() {
			// special base for aspects
			ret = &rt.Kind{Path: []string{name}}
		} else {
			if k, e := ks.makeAspect(name); e != nil {
				err = e
			} else if k != nil {
				ret = k
			} else if fs, ok := ks.Builder.Fields[name]; !ok {
				err = fmt.Errorf("unknown kind %q", name)
			} else {
				if k, e := ks.makeKind(name, fs); e != nil {
					err = e
				} else {
					ret = k
				}
			}
		}
		// cache the return
		if ret != nil {
			if ks.Kinds == nil {
				ks.Kinds = make(KindMap)
			}
			ks.Kinds[name] = ret
		}
	}
	return
}

func (ks *Kinds) makeKind(name string, pfs *[]rt.Field) (ret *rt.Kind, err error) {
	var parents []string
	if p, ok := ks.Builder.Parents[name]; ok {
		if k, e := ks.GetKindByName(p); e != nil {
			err = e
		} else {
			parents = k.Ancestors()
		}
	}
	if err == nil {
		fields := *pfs
		path := append([]string{name}, parents...)
		ret = &rt.Kind{Path: path, Fields: fields, Aspects: rt.MakeAspects(ks, fields)}
	}
	return
}

func (ks *Kinds) makeAspect(name string) (ret *rt.Kind, err error) {
	for _, a := range ks.Builder.Aspects {
		if a.Name == name {
			if parent, e := ks.GetKindByName(kindsOf.Aspect.String()); e != nil {
				err = e
				break
			} else {
				// create the kind from the stored fields
				path := append([]string{name}, parent.Ancestors()...)
				ret = &rt.Kind{Path: path, Fields: MakeFieldsFromTraits(a.Traits)}
				break
			}
		}
	}
	return
}

func MakeFieldsFromTraits(ts []string) (ret []rt.Field) {
	for _, t := range ts {
		f := rt.Field{Name: t, Affinity: affine.Text, Type: t}
		ret = append(ret, f)
	}
	return
}

// generate fields from type t using reflection
func (kb *KindBuilder) addType(ks *Kinds, t r.Type) {
	type stringer interface{ String() string }
	rstringer := r.TypeOf((*stringer)(nil)).Elem()
	if kb.Fields == nil {
		kb.Fields = make(FieldMap)
		kb.Parents = make(ParentMap)
	}

	// already built?
	name := nameOfType(t)
	if kb.Fields[name] != nil {
		return
	}
	pfields := new([]rt.Field) // pointers are used to block recursion
	kb.Fields[name] = pfields

	for i, cnt := 0, t.NumField(); i < cnt; i++ {
		f := t.Field(i)
		fieldType := f.Type
		var b struct {
			Aff  affine.Affinity
			Type string //  comma separated path
		}
		switch k := fieldType.Kind(); k {
		default:
			log.Panic("unknown kind", k)

		case r.Bool:
			b.Aff = affine.Bool

		case r.String:
			// note: text type indicates kind, not golang type
			b.Aff, b.Type = affine.Text, ""

		case r.Ptr:
			fieldType = fieldType.Elem()
			b.Type = nameOfType(fieldType)
			b.Aff = affine.Record
			kb.addType(ks, fieldType)

		case r.Struct:
			if !f.Anonymous {
				b.Type = nameOfType(fieldType)
				b.Aff = affine.Record
				// recurse to add the type
				kb.addType(ks, fieldType)

			} else if len((*pfields)) > 0 {
				panic("anonymous structs are used for hierarchy and should be the first member")
			} else {
				parentName := nameOfType(fieldType)
				kb.Parents[name] = parentName

				// note: this doesnt set affinity, and so doesnt get added as a field
				// todo: document what's happening here.
				//kb.Fields[name] = pfields
				// parent := ks.Kind(name)
				// b.Type = parent.Name()
			}

		case r.Slice:
			elType := fieldType.Elem()
			switch k := elType.Kind(); k {
			case r.String:
				// note: text type indicates kind, not golang type
				b.Aff, b.Type = affine.TextList, ""
			case r.Float64:
				b.Aff, b.Type = affine.NumList, k.String()
			case r.Struct:
				b.Aff, b.Type = affine.RecordList, nameOfType(elType)
				kb.addType(ks, elType)

			default:
				log.Panic("unknown slice", elType.String())
			}

		case r.Float64:
			b.Aff, b.Type = affine.Num, k.String()

		case r.Int: // enumeration
			aspect := nameOfType(fieldType)
			b.Aff, b.Type = affine.Text, aspect
			if !fieldType.Implements(rstringer) {
				log.Panic("unknown enum", aspect)
			}
			x := r.New(fieldType).Elem()
			var traits []string
			for j := int64(0); j < 25; j++ {
				x.SetInt(j)
				trait := x.Interface().(stringer).String()
				end := fmt.Sprintf("%s(%d)", fieldType.Name(), j)
				if trait == end {
					break
				}
				name := inflect.MixedCaseToSpaces(trait)
				traits = append(traits, name)
			}
			kb.Aspects = append(kb.Aspects, rt.Aspect{
				Name:   aspect,
				Traits: traits,
			})
		}
		if len(b.Aff) > 0 {
			name := inflect.MixedCaseToSpaces(f.Name)
			(*pfields) = append((*pfields), rt.Field{Name: name, Affinity: b.Aff, Type: b.Type})
		}
	}
}

func nameOfType(t r.Type) string {
	return inflect.MixedCaseToSpaces(t.Name())
}
