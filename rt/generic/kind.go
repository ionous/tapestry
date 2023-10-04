package generic

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"golang.org/x/exp/slices"
)

//
type Kind struct {
	kinds     Kinds
	name      string // keeping name *and* path makes debugging easier
	parent    *Kind
	fields    []Field
	aspects   []Aspect
	lastField int // cache of last accessed field
}

type Field struct {
	Name     string
	Affinity affine.Affinity
	Type     string // ex. kind for text types
}

type Aspect struct {
	Name   string // matches the field name
	Traits []string
}

// NewKind -
func NewKind(kinds Kinds, name string, parent *Kind, fields []Field) *Kind {
	return newKind(kinds, name, parent, fields, nil)
}

// a record without a named kind
func NewAnonymousRecord(kinds Kinds, fields []Field) *Record {
	return newKind(kinds, "", nil, fields, nil).NewRecord()
}

func newKind(kinds Kinds, name string, parent *Kind, fields []Field, aspects []Aspect) *Kind {
	if parent != nil { // fix? field lists are stored "flat" to simplify copy, record, etc.
		// have to copy or they can share memory, and bad things happen with other kinds.
		fields = append(append([]Field(nil), parent.fields...), fields...)
		aspects = append(append([]Aspect(nil), parent.aspects...), aspects...)
	}
	return &Kind{kinds: kinds, name: name, parent: parent, fields: fields, aspects: aspects}
}

func (k *Kind) NewRecord() *Record {
	// we make a bunch of nil value placeholders which we fill by caching on demand.
	rec := &Record{kind: k, values: make([]Value, len(k.fields))}
	// set the default values for aspects?
	// alt: determine it on GetIndexedValue as per other defaults
	for _, a := range k.aspects {
		i := k.FieldIndex(a.Name)
		rec.values[i] = StringFrom(a.Traits[0], a.Name)
	}
	return rec
}

func Base(k *Kind) string {
	for ; k.parent != nil; k = k.parent {
	}
	return k.name
}

// Ancestor list, root towards the start; the name of this kind at the end.
func Path(k *Kind) (ret []string) {
	for ; k != nil; k = k.parent {
		ret = append(ret, k.name)
	}
	slices.Reverse(ret)
	return
}

func (k *Kind) Parent() (ret *Kind) {
	return k.parent
}

func (k *Kind) Implements(name string) (okay bool) {
	for ; k != nil; k = k.parent {
		if k.name == name {
			okay = true
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

// panics if out of range
func (k *Kind) Field(i int) (ret Field) {
	return k.fields[i]
}

// searches for the field which handles the passed field;
// for traits, it returns the index of its associated aspect.
// returns -1 if no matching field was found
func (k *Kind) FieldIndex(n string) (ret int) {
	if prev := k.lastField; prev >= 0 && prev < len(k.fields) && k.fields[prev].Name == n {
		ret = prev
	} else {
		ret = -1 // provisionally
		if i := findTrait(n, k.aspects); i >= 0 {
			n = k.aspects[i].Name
		}
		if i := findField(n, k.fields); i >= 0 {
			ret = i
		}
		k.lastField = ret
	}
	return
}

func findField(field string, fields []Field) (ret int) {
	ret = -1 // provisionally
	for i, f := range fields {
		if f.Name == field {
			ret = i
			break
		}
	}
	return
}

// find aspect from trait name in a sorted list of traits
func findTrait(trait string, aspects []Aspect) (ret int) {
	ret = -1 // provisionally
	for i, a := range aspects {
		for _, t := range a.Traits {
			if trait == t {
				ret = i
				break
			}
		}
	}
	return
}
