package generic

import (
	"log"

	"github.com/ionous/iffy/affine"
)

// we bake it down for faster, easier indexed access.
type Kind struct {
	kinds   Kinds
	name    string
	fields  []Field
	traits  []trait
	lastOne int // one-based index of last field
}

type Field struct {
	Name     string
	Affinity affine.Affinity
	Type     string // ex. record name, "aspect", "trait", "float64", ...
}

// aspects are a specific kind of record where every field is a boolean trait
func NewKind(kinds Kinds, name string, fields []Field) *Kind {
	return &Kind{kinds: kinds, name: name, fields: fields}
}

// fix: temp till all kinds are moved to assembly
func (k *Kind) IsStaleKind(kinds Kinds) bool {
	return kinds != k.kinds
}

func (k *Kind) NewRecord() *Record {
	return &Record{kind: k, values: make([]interface{}, len(k.fields))}
}

func (k *Kind) Name() string {
	return k.name
}

func (k *Kind) NumField() int {
	return len(k.fields)
}

func (k *Kind) Field(i int) Field {
	k.lastOne = i + 1
	return k.fields[i]
}

// searches for the field which handles the passed field;
// for traits, it returns the index of its associated aspect.
// returns -1 if no matching field was found
func (k *Kind) FieldIndex(n string) (ret int) {
	if prev := k.lastOne - 1; prev >= 0 && k.fields[prev].Name == n {
		ret = prev
	} else {
		k.ensureTraits()
		if aspect := findAspect(n, k.traits); len(aspect) > 0 {
			ret = k.fieldIndex(aspect)
		} else {
			ret = k.fieldIndex(n)
		}
		k.lastOne = ret + 1
	}
	return
}

func (k *Kind) fieldIndex(field string) (ret int) {
	ret = -1 // provisionally
	for i, f := range k.fields {
		if f.Name == field {
			ret = i
			break
		}
	}
	return
}

func (k *Kind) ensureTraits() {
	if k.traits == nil {
		var ts []trait
		for _, ft := range k.fields {
			if ft.Type == "aspect" {
				if aspect, e := k.kinds.GetKindByName(ft.Name); e != nil {
					log.Println("unknown aspect", ft.Name)
				} else {
					ts = makeTraits(aspect, ts)
				}
			}
		}
		if len(ts) == 0 {
			ts = make([]trait, 0, 0)
		} else {
			sortTraits(ts)
		}
		k.traits = ts
	}
}
