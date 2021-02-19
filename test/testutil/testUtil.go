package testutil

import (
	"fmt"
	r "reflect"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/export/tag"
	"git.sr.ht/~ionous/iffy/lang"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"github.com/ionous/errutil"
)

type Kinds struct {
	Kinds  KindMap
	Fields FieldMap
}

type KindMap map[string]*g.Kind
type FieldMap map[string][]g.Field

// register kinds from a struct using reflection
func (ks *Kinds) AddKinds(is ...interface{}) {
	for _, el := range is {
		ks.Fields = kindsForType(ks.Fields, r.TypeOf(el).Elem())
	}
}

func (ks *Kinds) New(name string, valuePairs ...interface{}) *g.Record {
	v := ks.Kind(name).NewRecord()
	if len(valuePairs) > 0 {
		if e := SetRecord(v, valuePairs...); e != nil {
			panic(e)
		}
	}
	return v
}

func (ks *Kinds) Kind(name string) (ret *g.Kind) {
	if k, e := ks.GetKindByName(name); e != nil {
		panic(e)
	} else {
		ret = k
	}
	return
}

func (ks *Kinds) GetKindByName(name string) (ret *g.Kind, err error) {
	if k, ok := ks.Kinds[name]; ok {
		ret = k // we created the kind already
	} else if fs, ok := ks.Fields[name]; !ok {
		err = errutil.New("unknown kind", name)
	} else {
		if ks.Kinds == nil {
			ks.Kinds = make(KindMap)
		}
		// create the kind from the stored fields
		k := g.NewKind(ks, name, fs)
		ks.Kinds[name] = k
		ret = k
	}
	return
}

// generate kinds from a struct using reflection
func kindsForType(kinds FieldMap, t r.Type) FieldMap {
	type stringer interface{ String() string }
	rstringer := r.TypeOf((*stringer)(nil)).Elem()
	if kinds == nil {
		kinds = make(FieldMap)
	}

	var fields []g.Field
	for i, cnt := 0, t.NumField(); i < cnt; i++ {
		f := t.Field(i)
		fieldType := f.Type
		var a affine.Affinity
		var t string
		switch k := fieldType.Kind(); k {
		default:
			panic(errutil.Sprint("unknown kind", k))
		case r.Bool:
			tags := tag.ReadTag(f.Tag)
			if _, ok := tags.Find("bool"); ok {
				a, t = affine.Bool, k.String()
			} else {
				a, t = affine.Text, "aspect"
				n := lang.Underscore(f.Name)
				// the name of the aspect is the name of the field
				kinds[n] = []g.Field{
					// false first.
					{Name: "not_" + n, Affinity: affine.Bool, Type: "trait"},
					{Name: "is_" + n, Affinity: affine.Bool, Type: "trait"},
				}
			}

		case r.String:
			a, t = affine.Text, k.String()

		case r.Struct:
			a, t = affine.Record, nameOfType(fieldType)
			kinds = kindsForType(kinds, fieldType)

		case r.Slice:
			elType := fieldType.Elem()
			switch k := elType.Kind(); k {
			case r.String:
				a, t = affine.TextList, k.String()
			case r.Float64:
				a, t = affine.NumList, k.String()
			case r.Struct:
				a, t = affine.RecordList, nameOfType(elType)
				kinds = kindsForType(kinds, elType)

			default:
				panic(errutil.Sprint("unknown slice", elType.String()))
			}

		case r.Float64:
			a, t = affine.Number, k.String()

		case r.Int:
			a, t = affine.Text, "aspect"
			if !fieldType.Implements(rstringer) {
				panic("unknown enum")
			}
			x := r.New(fieldType).Elem()
			var traits []g.Field
			for j := int64(0); j < 25; j++ {
				x.SetInt(j)
				trait := x.Interface().(stringer).String()
				end := fmt.Sprintf("%s(%d)", fieldType.Name(), j)
				if trait == end {
					break
				}
				name := lang.Underscore(trait)
				traits = append(traits, g.Field{Name: name, Affinity: affine.Bool, Type: "trait"})
			}
			aspect := nameOfType(fieldType)
			kinds[aspect] = traits
		}
		name := lang.Underscore(f.Name)
		fields = append(fields, g.Field{Name: name, Affinity: a, Type: t})

	}
	name := nameOfType(t)
	kinds[name] = fields
	return kinds
}

func nameOfType(t r.Type) string {
	return lang.Underscore(t.Name())
}
