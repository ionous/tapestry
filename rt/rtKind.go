package rt

import (
	"slices"

	"git.sr.ht/~ionous/tapestry/affine"
)

// type system for tapestry objects and records.
type Kind struct {
	path      []string // the name is first
	fields    []Field
	aspects   []Aspect
	lastField int // cache of last accessed field
}

// kinds are composed of fields.
type Field struct {
	Name     string          // name of the field, unique within the kind
	Affinity affine.Affinity // one of the built in data types
	Type     string          // subtype; ex. kind for text types
	Init     Assignment      // default initialization
}

// the traits of any aspects behave like separate boolean fields;
// ex. if the list of fields contains a "colour" aspect with traits "red", "blue", "green"
// then the returned kind will respond to "colour", "red", "blue", and "green".
type Aspect struct {
	Name   string // matches the field name
	Traits []string
}

// NewKind -
// path indicates ancestry, the name of this kind at the start, parent kinds to the right.
// fields are expected to be the accumulated fields of all ancestors;
// parent fields to the *left*, more derived fields to the right. ( the reverse of path. )
// descendant kinds can have different initial assignments for the same named field.
func NewKind(path []string, fields []Field, aspects []Aspect) *Kind {
	return &Kind{path: path, fields: fields, aspects: aspects}
}

// semi unique identifier for this kind
func (k *Kind) Name() (ret string) {
	if len(k.path) > 0 {
		ret = k.path[0]
	}
	return
}

// ancestors of this kind, the name of this kind at the head of the list
// can be nil if the kind is "anonymous"
func (k *Kind) Path() []string {
	return k.path
}

// does this kind have the named ancestor?
// ( this is a shortcut
func (k *Kind) Implements(name string) bool {
	return slices.Contains(k.path, name)
}

// number of fields contained by this ( and all parent kinds )
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
