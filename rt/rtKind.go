package rt

import (
	"slices"

	"git.sr.ht/~ionous/tapestry/affine"
)

// used for tapestry objects and records
// a kind is a collection of fields
// with optional inheritance.
type Kind struct {
	// indicates ancestry, the name of this kind at the start, parent kinds to the right.
	// can be nil if anonymous
	Path []string
	// holds the accumulated fields of all ancestors.
	// parent fields to the *left*, more derived fields to the right. ( the reverse of path. )
	// descendant kinds can have different initial assignments for the same named field.
	Fields []Field `json:",omitempty"`
	// a subset of fields for those containing traits
	// see MakeAspects.
	Aspects []Aspect `json:",omitempty"`
	// a cache of last accessed field
	lastField int
}

// member of a kind.
type Field struct {
	Name     string          // name of the field, unique within the kind
	Affinity affine.Affinity // one of the built in data types
	Type     string          `json:",omitempty"` // subtype; ex. kind for text types
	Init     Assignment      `json:",omitempty"` // default initialization
}

// the traits of any aspects behave like separate boolean fields;
// ex. if the list of fields contains a "colour" aspect with traits "red", "blue", "green"
// then the returned kind will respond to "colour", "red", "blue", and "green".
type Aspect struct {
	Name   string // matches the field name
	Traits []string
}

// semi unique identifier for this kind
// returns the empty string for anonymous kinds
func (k *Kind) Name() (ret string) {
	if len(k.Path) > 0 {
		ret = k.Path[0]
	}
	return
}

// ancestors of this kind, the name of this kind at the head of the list.
// the slice isn't a copy; it should be treated as read-only.
// nil for anonymous kinds.
func (k *Kind) Ancestors() []string {
	return k.Path
}

// true if the kind has the named ancestor.
// ( this is a shortcut for testing the passed name in Path() )
func (k *Kind) Implements(name string) bool {
	return slices.Contains(k.Path, name)
}

// number of fields contained by this ( and all parent kinds )
func (k *Kind) FieldCount() int {
	return len(k.Fields)
}

// return a description of an indexed field.
// panics if out of range
func (k *Kind) Field(i int) (ret Field) {
	return k.Fields[i]
}

// returns the index of the matching field;
// for traits, that's the index of its aspect.
// returns -1 if no matching field was found.
func (k *Kind) FieldIndex(n string) (ret int) {
	if prev := k.lastField; prev >= 0 && prev < len(k.Fields) && k.Fields[prev].Name == n {
		ret = prev
	} else {
		ret = -1 // provisionally
		if i := k.FindAspectByTrait(n); i >= 0 {
			n = k.Aspects[i].Name
		}
		if i := findField(n, k.Fields); i >= 0 {
			ret = i
		}
		k.lastField = ret
	}
	return
}

func (k *Kind) Aspect(i int) (ret Aspect) {
	return k.Aspects[i]
}

func (k *Kind) AspectIndex(aspect string) (ret int) {
	return slices.IndexFunc(k.Aspects, func(a Aspect) bool {
		return a.Name == aspect
	})
}

// return the index of the aspect containing the passed trait
func (k *Kind) FindAspectByTrait(trait string) int {
	return slices.IndexFunc(k.Aspects, func(a Aspect) bool {
		return slices.Contains(a.Traits, trait)
	})
}

func findField(field string, fields []Field) (ret int) {
	return slices.IndexFunc(fields, func(f Field) bool {
		return f.Name == field
	})
}
