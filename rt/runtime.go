package rt

import (
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/writer"
)

// Scope - implements a portion of generic.Value
// while usually used with record under the hood --
// it fields the Value interface to show it plays by the same aliasing rules:
// set generates a copy.
type Scope interface {
	// should return g.UnknownVariable or g.UnknownField
	FieldByName(field string) (g.Value, error)
	SetFieldByName(field string, val g.Value) error
}

// VariableStack - a pool of record like name-value pairs.
// if a variable isnt found in the most recently pushed scope
// the next most recently pushed scope gets checked and so on.
type VariableStack interface {
	// completely replace the current record lookup
	// ex. patterns have opaque namespaces
	ReplaceScope(Scope) (prev Scope)
	// add a set of variables to the internal stack.
	// ex. loops add to the current namespace.
	PushScope(Scope)
	// remove the most recently added set of variables from the internal stack.
	PopScope()
}

// Runtime environment for an in-progress game.
type Runtime interface {
	// objects are grouped into potentially hierarchical "domains"
	// de/activating makes those groups hidden/visible to the runtime.
	// Domain hierarchy is defined at assembly time.
	ActivateDomain(name string, enable bool)
	// find a function, test, or pattern addressed by name
	// pv should be a pointer to a concrete type.
	GetEvalByName(name string, pv interface{}) error
	// record manipulation
	GetKindByName(name string) (*g.Kind, error)
	//
	GetRules(name string, pflags *Flags) ([]Rule, error)
	//
	RelateTo(a, b, relation string) error
	RelativesOf(a, relation string) ([]string, error)
	ReciprocalsOf(b, relation string) ([]string, error)
	// modifies the behavior of Get/SetField object.Variable
	VariableStack
	// various runtime objects (ex. nouns, kinds, etc. ) store data addressed by name.
	// the objects and their fields depend on implementation and context.
	// see package object for a variety of common objects.
	GetField(object, field string) (g.Value, error)
	// store, or at least attempt to store, the passed value at the named field in the named object.
	// it may return an error if the value is not of a compatible type,
	// if its considered to be read-only, or if there is no predeclared value of that name.
	SetField(object, field string, value g.Value) error
	//
	// turn single words into their plural variants, and vice-versa
	PluralOf(single string) string
	SingularOf(plural string) string
	// return a pseudo-random number
	Random(inclusiveMin, exclusiveMax int) int
	// Return the built-in writer, or the current override.
	Writer() writer.Output
	// Override the current writer
	SetWriter(writer.Output) (prev writer.Output)
}
