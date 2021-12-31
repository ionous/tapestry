package generic

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
)

// we bake it down for faster, easier indexed access.
type Kind struct {
	kinds   Kinds
	name    string // keeping name *and* path makes debugging easier
	path    []string
	fields  []Field // currently, stored flat. fix? save some space and search through the hierarchy?
	traits  []trait
	lastOne int // one-based index of last field
}

type Field struct {
	Name     string
	Affinity affine.Affinity
	Type     string // ex. kind for text types ( future? "aspect", "trait", "float64", ... )
}

// see also eph.AspectParm
func (ft *Field) isAspectLike() bool {
	return ft.Affinity == affine.Text && ft.Name == ft.Type
}

// path is a list of ancestors with the root at the start.
func NewKind(kinds Kinds, name string, path []string, fields []Field) *Kind {
	fullpath := append(path, name)
	return &Kind{kinds: kinds, name: name, path: fullpath, fields: fields}
}

func NewAnonymousRecord(kinds Kinds, fields []Field) *Record {
	return NewKind(kinds, "", nil, fields).NewRecord()
}

func (k *Kind) NewRecord() *Record {
	// we make a bunch of nil value placeholders which we fill by caching on demand.
	return &Record{kind: k, values: make([]Value, len(k.fields))}
}

// Ancestor list, root towards the start; the name of this kind at the end.
func (k *Kind) Path() (ret []string) {
	ret = append(ret, k.path...) // copies the slice
	return
}

func (k *Kind) Implements(i string) (ret bool) {
	for _, p := range k.path {
		if i == p {
			ret = true
			break
		}
	}
	return
}

func (k *Kind) Name() (ret string) {
	return k.name
}

func (k *Kind) NumField() int {
	return len(k.fields)
}

// 0 indexed
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
			if ft.isAspectLike() {
				// if this fails, we are likely to return an error through GetIndexedField at some point so...
				if aspect, e := k.kinds.GetKindByName(ft.Type); e == nil {
					if aok := aspect.Implements(kindsOf.Aspect.String()); aok {
						ts = makeTraits(aspect, ts)
					}
				}
			}
		}
		if len(ts) == 0 {
			ts = make([]trait, 0)
		} else {
			sortTraits(ts)
		}
		k.traits = ts
	}
}
