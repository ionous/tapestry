package testutil

import (
	"fmt"
	r "reflect"
	"strings"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/lang"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/test/tag"
	"github.com/ionous/errutil"
)

type Kinds struct {
	Kinds  KindMap
	Fields FieldMap
}

type KindMap map[string]*g.Kind
type FieldMap map[string][]g.Field

// register kinds from a struct using reflection
// note: this doesnt actually create the kinds....
// it preps them for use -- they are created on Get()
func (ks *Kinds) AddKinds(is ...interface{}) {
	for _, el := range is {
		ks.Fields = kindsForType(ks, ks.Fields, r.TypeOf(el).Elem())
	}
}

func (ks *Kinds) NewRecord(name string, valuePairs ...interface{}) *g.Record {
	v := ks.Kind(name).NewRecord()
	if len(valuePairs) > 0 {
		if e := SetRecord(v, valuePairs...); e != nil {
			panic(e)
		}
	}
	return v
}

// return a kind ( declared via AddKinds ) panic if it doesn't exist.
func (ks *Kinds) Kind(name string) (ret *g.Kind) {
	if k, e := ks.GetKindByName(name); e != nil {
		panic(e)
	} else {
		ret = k
	}
	return
}

// return a kind ( declared via AddKinds ) or error if it doesnt exist.
func (ks *Kinds) GetKindByName(name string) (ret *g.Kind, err error) {
	if k, ok := ks.Kinds[name]; ok {
		ret = k // we created the kind already
	} else if fs, ok := ks.Fields[name]; !ok {
		err = errutil.New("unknown kind", name)
	} else {
		if ks.Kinds == nil {
			ks.Kinds = make(KindMap)
		}
		// magic parent field for fake objects
		if len(fs) > 0 && len(fs[0].Affinity) == 0 {
			name = name + "," + fs[0].Type
		}
		// create the kind from the stored fields
		k := g.NewKind(ks, name, fs)
		ks.Kinds[name] = k
		ret = k
	}
	return
}

// generate fields from type t using reflection
func kindsForType(ks *Kinds, kinds FieldMap, t r.Type) FieldMap {
	type stringer interface{ String() string }
	rstringer := r.TypeOf((*stringer)(nil)).Elem()
	if kinds == nil {
		kinds = make(FieldMap)
	}

	var path string
	var fields []g.Field
	for i, cnt := 0, t.NumField(); i < cnt; i++ {
		f := t.Field(i)
		fieldType := f.Type
		var pending struct {
			Aff  affine.Affinity
			Type string
		}
		switch k := fieldType.Kind(); k {
		default:
			panic(errutil.Sprint("unknown kind", k))
		case r.Bool:
			tags := tag.ReadTag(f.Tag)
			if _, ok := tags.Find("bool"); ok {
				pending.Aff, pending.Type = affine.Bool, k.String()
			} else {
				// the name of the aspect is the name of the field and its class
				n := lang.Underscore(f.Name)
				pending.Aff, pending.Type = affine.Text, n
				kinds[n] = []g.Field{
					// false first.
					{Name: "not_" + n, Affinity: affine.Bool /*, Type: "trait"*/},
					{Name: "is_" + n, Affinity: affine.Bool /*, Type: "trait"*/},
				}
			}

		case r.String:
			pending.Aff, pending.Type = affine.Text, k.String()

		case r.Struct:
			pending.Type = nameOfType(fieldType)
			if f.Anonymous {
				if len(fields) > 0 {
					panic("anonymous structs are used for hierarchy and should be the first member")
				}
				parent := ks.Kind(pending.Type)
				pending.Type = strings.Join(parent.Path(), ",")

			} else {
				pending.Aff = affine.Record
				kinds = kindsForType(ks, kinds, fieldType)
			}

		case r.Slice:
			elType := fieldType.Elem()
			switch k := elType.Kind(); k {
			case r.String:
				pending.Aff, pending.Type = affine.TextList, k.String()
			case r.Float64:
				pending.Aff, pending.Type = affine.NumList, k.String()
			case r.Struct:
				pending.Aff, pending.Type = affine.RecordList, nameOfType(elType)
				kinds = kindsForType(ks, kinds, elType)

			default:
				panic(errutil.Sprint("unknown slice", elType.String()))
			}

		case r.Float64:
			pending.Aff, pending.Type = affine.Number, k.String()

		case r.Int:
			aspect := nameOfType(fieldType)
			pending.Aff, pending.Type = affine.Text, aspect
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
				traits = append(traits, g.Field{Name: name, Affinity: affine.Bool /*, Type: "trait"*/})
			}
			kinds[aspect] = traits
		}
		if len(pending.Aff) > 0 || len(path) > 0 {
			name := lang.Underscore(f.Name)
			fields = append(fields, g.Field{Name: name, Affinity: pending.Aff, Type: pending.Type})
		}
	}
	name := nameOfType(t)
	kinds[name] = fields
	return kinds
}

func nameOfType(t r.Type) string {
	return lang.Underscore(t.Name())
}
