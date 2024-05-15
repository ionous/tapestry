package rt

import (
	"io"

	"git.sr.ht/~ionous/tapestry/affine"
)

// Scope - establishes a pool of local variables.
// Scopes are usually accessed via the runtime Set/GetField.
type Scope interface {
	// return g.Unknown if the named field/variable doesn't exist.
	FieldByName(field string) (Value, error)
	// set without copying
	SetFieldByName(field string, val Value) error
}

// Type database
type Kinds interface {
	GetKindByName(n string) (*Kind, error)
}

// Runtime environment for an in-progress game.
type Runtime interface {
	// ActivateDomain - objects are grouped into potentially hierarchical "domains".
	// de/activating makes those groups hidden/visible to the runtime.
	// Domain hierarchy is defined at assembly time.
	ActivateDomain(name string) error
	// GetKindByName - type description
	GetKindByName(name string) (*Kind, error)
	// Call - run the pattern defined by the passed record.
	// passes the expected return because patterns can be called in ambiguous context ( ex. template expressions )
	Call(name string, expectedReturn affine.Affinity, keys []string, vals []Value) (Value, error)
	// RelativesOf - returns a list, even for one-to-one relationships.
	RelativesOf(a, relation string) (Value, error)
	// ReciprocalsOf - returns a list, even for one-to-one relationships.
	ReciprocalsOf(b, relation string) (Value, error)
	// RelateTo - establish a new relation between nouns and and b.
	RelateTo(a, b, relation string) error
	// PushScope - modifies the behavior of Get/SetField meta.Variable
	// by adding a set of variables to the current namespace.
	// ex. loops add iterators
	PushScope(Scope)
	// PopScope - remove the most recently added set of variables from the internal stack.
	PopScope()
	// GetField - various runtime objects (ex. nouns, kinds, etc.) store data addressed by name.
	// the objects and their fields depend on implementation and context.
	// see package meta for a variety of common objects.
	GetField(object, field string) (Value, error)
	// SetField - store, or at least attempt to store, a *copy* of the value into the field of the named object.
	// it can return an error:
	// if the value is not of a compatible type,
	// if the field is considered to read-only,
	// or if there is no object or field of the indicated names.
	SetField(object, field string, value Value) error
	// PluralOf - turn single words into their plural variants, and vice-versa.
	// each plural word maps to a unique singular.
	// for example, if the singular of "people" is "person", it cant also be "personage".
	PluralOf(single string) string
	// note: one plural word can map to multiple single words.
	// in that case the returned singular word is arbitrary ( if theoretically consistent )
	// for example, "person" can have the plural "people" or "persons" and this could return either.
	SingularOf(plural string) string
	// Random - return a pseudo-random number.
	Random(inclusiveMin, exclusiveMax int) int
	// Writer - Return the built-in writer, or the current override.
	Writer() io.Writer
	// SetWriter - Override the current writer.
	SetWriter(io.Writer) (prev io.Writer)
}
