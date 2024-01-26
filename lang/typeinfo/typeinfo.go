// Package typeinfo describes tapestry autogenerated structures.
// It follows closely from tapestry typespecs:
// The typespecs describes the user specification ( in .ifspec, or .tells files )
// while the typeinfo describes the golang structures derived from those specs.
// unless otherwise specified, strings are underscore separated and lowercase.
package typeinfo

// implemented by every auto-generated command
type Inspector interface {
	// returns the typeinfo for the instance
	// and whether it repeats.
	Inspect() (T, bool)
}

// implemented by auto-generated flows
type FlowInspector interface {
	Inspector
	GetMarkup(ensure bool) map[string]any
}

// package listing of type data
type TypeSet struct {
	Name       string
	Slot       []*Slot
	Flow       []*Flow
	Str        []*Str
	Num        []*Num
	Signatures map[uint64]any
}

// marker interface implemented by each kind of typeinfo:
// Flow, Slot, Str, and Num.
type T interface {
	TypeInfo() T
	TypeName() string
	TypeMarkup() map[string]any
}

type Flow struct {
	Name   string         // unique name for this type
	Lede   string         // the compact format leading text
	Slots  []*Slot        // interfaces that a command implements
	Terms  []Term         // terms of the command
	Markup map[string]any // metadata shared by all instances of this type
}

// designates Flow as typeinfo; returns itself
func (t *Flow) TypeInfo() T {
	return t
}

func (t *Flow) TypeName() string {
	return t.Name
}

func (t *Flow) TypeMarkup() map[string]any {
	return t.Markup
}

// a member of a Flow.
type Term struct {
	Name     string // go lang name; unique within its flow.
	Label    string // the compact format signature
	Private  bool   // a member that only exists in memory; never serialized
	Optional bool   // true when the term can be omitted for instances of the flow.
	Repeats  bool   // true when the term can have multiple values ( all of the same type )
	Type     T      // a pointer to Flow, Slot, Str, or Num; or, nil if private
}

func (t Term) IsAnonymous() bool {
	return t.Label == "_"
}

type Slot struct {
	Name   string         // unique name for this type
	Markup map[string]any // metadata shared by all instances of this type
}

// designates Slot as typeinfo; returns itself
func (t *Slot) TypeInfo() T {
	return t
}

func (t *Slot) TypeName() string {
	return t.Name
}

func (t *Slot) TypeMarkup() map[string]any {
	return t.Markup
}

type Str struct {
	Name    string         // unique name for this type
	Options []string       // for enumerations; for plain strings, this is nil.
	Markup  map[string]any // metadata shared by all instances of this type
}

// designates Str as typeinfo; returns itself
func (t *Str) TypeInfo() T {
	return t
}

func (t *Str) TypeName() string {
	return t.Name
}

func (t *Str) TypeMarkup() map[string]any {
	return t.Markup
}

type Num struct {
	Name   string         // unique name for this type
	Markup map[string]any // metadata shared by all instances of this type
}

// designates Num as typeinfo; returns itself
func (t *Num) TypeInfo() T {
	return t
}

func (t *Num) TypeName() string {
	return t.Name
}

func (t *Num) TypeMarkup() map[string]any {
	return t.Markup
}
