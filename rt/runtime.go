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

// VariableStack - a pool of record like name-value pairs.
// if a variable isnt found in the most recently pushed scope
// the next most recently pushed scope gets checked and so on.
type VariableStack interface {
	// ReplaceScope - swap the current scope for the passed scope
	// when init is true, try to initialize all the the values in the passed scope.
	// only init can return an error.
	ReplaceScope(scope Scope, init bool) (ret Scope, err error)
	// PushScope - add a set of variables to the internal stack.
	// ex. loops add to the current namespace.
	PushScope(Scope)
	// PopScope - remove the most recently added set of variables from the internal stack.
	PopScope()
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
	Call(pat *g.Record, expectedReturn affine.Affinity) (g.Value, error)
	// Send - trigger the named event, passing the objects to visit: target first, root-most last.
	Send(pat *g.Record, up []string) (g.Value, error)
	// RelativesOf - returns a list, even for one-to-one relationships.
	RelativesOf(a, relation string) (g.Value, error)
	// ReciprocalsOf - returns a list, even for one-to-one relationships.
	ReciprocalsOf(b, relation string) (g.Value, error)
	// RelateTo - establish a new relation between nouns and and b.
	RelateTo(a, b, relation string) error
	// VariableStack - modifies the behavior of Get/SetField meta.Variable.
	VariableStack
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
	PluralOf(single string) string
	SingularOf(plural string) string
	// Random - return a pseudo-random number.
	Random(inclusiveMin, exclusiveMax int) int
	// Writer - Return the built-in writer, or the current override.
	Writer() io.Writer
	// SetWriter - Override the current writer.
	SetWriter(io.Writer) (prev io.Writer)
}
