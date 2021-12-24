package testutil

import (
	"fmt"
	r "reflect"

	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/lang"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/kindsOf"
	"git.sr.ht/~ionous/iffy/test/tag"
	"github.com/ionous/errutil"
)

type Kinds struct {
	Kinds   KindMap // stores completed kinds
	Builder KindBuilder
}

type KindMap map[string]*g.Kind
type FieldMap map[string][]g.Field // kind name to fields
type AspectMap map[string]bool

type KindBuilder struct {
	Aspects AspectMap
	Fields  FieldMap
}

// register kinds from a struct using reflection
// note: this doesnt actually create the kinds....
// it preps them for use -- they are created on Get()
func (ks *Kinds) AddKinds(is ...interface{}) {
	for _, el := range is {
		ks.Builder.addType(ks, r.TypeOf(el).Elem())
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
	} else {
		b := ks.Builder
		if fs, ok := b.Fields[name]; !ok {
			err = errutil.New("unknown kind", name)
		} else {
			if ks.Kinds == nil {
				ks.Kinds = make(KindMap)
			}
			var path []string
			// magic parent field for fake objects
			if len(fs) > 0 && len(fs[0].Affinity) == 0 {
				path = append(path, fs[0].Type)
			}
			// magic field for aspects
			if b.Aspects[name] {
				path = append(path, kindsOf.Aspect.String())
			}
			// create the kind from the stored fields
			k := g.NewKind(ks, name, path, fs)
			ks.Kinds[name] = k
			ret = k
		}
	}
	return
}

// generate fields from type t using reflection
func (ft *KindBuilder) addType(ks *Kinds, t r.Type) {
	type stringer interface{ String() string }
	rstringer := r.TypeOf((*stringer)(nil)).Elem()
	if ft.Fields == nil {
		ft.Fields = make(FieldMap)
		ft.Aspects = make(AspectMap)
	}

	var path string
	var fields []g.Field
	for i, cnt := 0, t.NumField(); i < cnt; i++ {
		f := t.Field(i)
		fieldType := f.Type
		var b struct {
			Aff  affine.Affinity
			Type string //  comma separated path
		}
		switch k := fieldType.Kind(); k {
		default:
			panic(errutil.Sprint("unknown kind", k))
		case r.Bool:
			tags := tag.ReadTag(f.Tag)
			if _, ok := tags.Find("bool"); ok {
				b.Aff, b.Type = affine.Bool, k.String()
			} else {
				// the name of the aspect is the name of the field and its class
				aspect := lang.Underscore(f.Name)
				b.Aff, b.Type = affine.Text, aspect
				ft.Fields[aspect] = []g.Field{
					// false first.
					{Name: "not_" + aspect, Affinity: affine.Bool /*, Type: "trait"*/},
					{Name: "is_" + aspect, Affinity: affine.Bool /*, Type: "trait"*/},
				}
				ft.Aspects[aspect] = true
			}

		case r.String:
			b.Aff, b.Type = affine.Text, k.String()

		case r.Struct:
			b.Type = nameOfType(fieldType)
			if f.Anonymous {
				if len(fields) > 0 {
					panic("anonymous structs are used for hierarchy and should be the first member")
				}
				parent := ks.Kind(b.Type)
				b.Type = parent.Name()

			} else {
				b.Aff = affine.Record
				ft.addType(ks, fieldType)
			}

		case r.Slice:
			elType := fieldType.Elem()
			switch k := elType.Kind(); k {
			case r.String:
				b.Aff, b.Type = affine.TextList, k.String()
			case r.Float64:
				b.Aff, b.Type = affine.NumList, k.String()
			case r.Struct:
				b.Aff, b.Type = affine.RecordList, nameOfType(elType)
				ft.addType(ks, elType)

			default:
				panic(errutil.Sprint("unknown slice", elType.String()))
			}

		case r.Float64:
			b.Aff, b.Type = affine.Number, k.String()

		case r.Int:
			aspect := nameOfType(fieldType)
			b.Aff, b.Type = affine.Text, aspect
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
			ft.Fields[aspect] = traits
			ft.Aspects[aspect] = true
		}
		if len(b.Aff) > 0 || len(path) > 0 {
			name := lang.Underscore(f.Name)
			fields = append(fields, g.Field{Name: name, Affinity: b.Aff, Type: b.Type})
		}
	}
	name := nameOfType(t)
	ft.Fields[name] = fields
}

func nameOfType(t r.Type) string {
	return lang.Underscore(t.Name())
}
