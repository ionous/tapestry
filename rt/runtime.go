package rt

import (
	"io"

	"git.sr.ht/~ionous/tapestry/affine"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
)

// Scope - establishes a pool of local variables.
// Scopes are usually accessed via the runtime Set/GetField.
type Scope interface {
	// note: return g.Unknown if the named field/variable isnt found.
	FieldByName(field string) (g.Value, error)
	// note: does not usually copy ( since Runtime does )
	SetFieldByName(field string, val g.Value) error
	// errors if the field doesnt exist
	SetFieldDirty(field string) error
}

// Runtime environment for an in-progress game.
type Runtime interface {
	// ActivateDomain - objects are grouped into potentially hierarchical "domains".
	// de/activating makes those groups hidden/visible to the runtime.
	// Domain hierarchy is defined at assembly time.
	ActivateDomain(name string) (ret string, err error)
	// GetKindByName  -record manipulation.
	GetKindByName(name string) (*g.Kind, error)
	// GetRules - return the runtime rules matching the passed pattern and target.
	GetRules(pattern, target string, pflags *Flags) ([]Rule, error)
	// Call - run the pattern defined by the passed record.
	// passes the expected return because patterns can be called in ambiguous context ( ex. template expressions )
	Call(name string, expectedReturn affine.Affinity, keys []string, vals []g.Value) (g.Value, error)
	// Send - trigger the named event, passing the objects to visit: target first, root-most last.
	Send(pat *g.Record, up []string) (g.Value, error)
	// RelativesOf - returns a list, even for one-to-one relationships.
	RelativesOf(a, relation string) (g.Value, error)
	// ReciprocalsOf - returns a list, even for one-to-one relationships.
	ReciprocalsOf(b, relation string) (g.Value, error)
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
	GetField(object, field string) (g.Value, error)
	// SetField - store, or at least attempt to store, a *copy* of the value into the field of the named object.
	// it can return an error:
	// if the value is not of a compatible type,
	// if the field is considered to read-only,
	// or if there is no object or field of the indicated names.
	SetField(object, field string, value g.Value) error
	// PluralOf - turn single words into their plural variants, and vice-versa.
	// each plural word maps to a unique singular.
	// for example, if the singular of "people" is "person", it cant also be "personage".
	PluralOf(single string) string
	// note: one plural word can map to multiple single words.
	// in that case the returned singular word is arbitrary ( if theoretically consistent )
	// for example, "person" can have the plural "people" or "persons" and this could return either.
	SingularOf(plural string) string
	// returns true if found; the word if not
	OppositeOf(string) string
	// Random - return a pseudo-random number.
	Random(inclusiveMin, exclusiveMax int) int
	// Writer - Return the built-in writer, or the current override.
	Writer() io.Writer
	// SetWriter - Override the current writer.
	SetWriter(io.Writer) (prev io.Writer)
}
