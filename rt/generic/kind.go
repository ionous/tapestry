package generic

import (
	"slices"

	"git.sr.ht/~ionous/tapestry/affine"
)

// Kinds database
// this isnt used by package generic, but its a common enough interface for tests and the runtime
type Kinds interface {
	GetKindByName(n string) (*Kind, error)
}

//
type Kind struct {
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

// a record without a named kind
func NewAnonymousRecord(fields []Field) *Record {
	return newKind("", nil, fields, nil).NewRecord()
}

// NewKind -
func NewKind(name string, parent *Kind, fields []Field) *Kind {
	return newKind(name, parent, fields, nil)
}

// when kinds are created via this method,
// the traits of any aspect fields act like separate boolean fields;
// without it, only the aspect text field itself exists.
// ex. if the list of fields contains a "colour" aspect with traits "red", "blue", "green"
// then the returned kind will respond to "colour", "red", "blue", and "green";
// with NewKind() it would respond only to "colour", the r,b,g fields wouldn't exist.
// its a bit of leaky abstraction because boolean traits are used only by objects.
func NewKindWithTraits(name string, parent *Kind, fields []Field, aspects []Aspect) *Kind {
	return newKind(name, parent, fields, aspects)
}

func newKind(name string, parent *Kind, fields []Field, aspects []Aspect) *Kind {
	if parent != nil { // fix? field lists are stored "flat" to simplify copy, record, etc.
		// have to copy or they can share memory, and bad things happen with other kinds.
		if len(parent.fields) > 0 {
			fields = append(append([]Field(nil), parent.fields...), fields...)
		}
		if len(parent.aspects) > 0 {
			aspects = append(append([]Aspect(nil), parent.aspects...), aspects...)
		}
	}
	return &Kind{name: name, parent: parent, fields: fields, aspects: aspects}
}

func (k *Kind) NewRecord() *Record {
	// we make a bunch of nil value placeholders which we fill by caching on demand.
	rec := &Record{kind: k, values: make([]Value, len(k.fields))}
	// set the default values for aspects?
	// alt: determine it on GetIndexedValue as per other defaults
	// for _, a := range k.aspects {
	// 	i := k.FieldIndex(a.Name)
	// 	rec.values[i] = StringFrom(a.Traits[0], a.Name)
	// }
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
