package rt

import (
	"git.sr.ht/~ionous/iffy/affine"
	g "git.sr.ht/~ionous/iffy/rt/generic"
	"git.sr.ht/~ionous/iffy/rt/writer"
)

// Scope - establishes local variables.
// it abides by the rules of the matching g.Value methods: set copies.
type Scope interface {
	// return g.Unknown if the named field/variable isnt found.
	FieldByName(field string) (g.Value, error)
	SetFieldByName(field string, val g.Value) error
}

// VariableStack - a pool of record like name-value pairs.
// if a variable isnt found in the most recently pushed scope
// the next most recently pushed scope gets checked and so on.
type VariableStack interface {
	// swap the current scope for the passed scope
	// when init is true, try to initialize all the the values in the passed scope.
	// only init can return an error.
	ReplaceScope(scope Scope, init bool) (ret Scope, err error)
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
	ActivateDomain(name string) (ret string, err error)
	// record manipulation
	GetKindByName(name string) (*g.Kind, error)
	// return the runtime rules matching the passed pattern and target
	GetRules(pattern, target string, pflags *Flags) ([]Rule, error)
	// run the named pattern; add can be blank for execute style patterns.
	Call(name string, aff affine.Affinity, args []Arg) (g.Value, error)
	// trigger the named event, passing the objects to visit: target first, root-most last.
	Send(name string, up []string, arg []Arg) (g.Value, error)
	// returns a list, even for one-to-one relationships
	RelativesOf(a, relation string) (g.Value, error)
	// returns a list, even for one-to-one relationships
	ReciprocalsOf(b, relation string) (g.Value, error)
	// establish a new relation between nouns and and b
	RelateTo(a, b, relation string) error
	// modifies the behavior of Get/SetField meta.Variable
	VariableStack
	// various runtime objects (ex. nouns, kinds, etc. ) store data addressed by name.
	// the objects and their fields depend on implementation and context.
	// see package meta for a variety of common objects.
	GetField(meta, field string) (g.Value, error)
	// store, or at least attempt to store, the passed value at the named field in the named object.
	// it may return an error if the value is not of a compatible type,
	// if its considered to be read-only, or if there is no predeclared value of that name.
	SetField(meta, field string, value g.Value) error
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
