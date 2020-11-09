package rt

import (
	g "github.com/ionous/iffy/rt/generic"
	"github.com/ionous/iffy/rt/writer"
)

// Scope provides access to a pool of named objects.
// Various runtime objects (ex. nouns, kinds, etc. ) store data addressed by name.
// The objects and their fields depend on implementation and context.
// See package object for a variety of common objects.
type Scope interface {
	GetField(object, field string) (g.Value, error)
	// Store, or at least attempt to store, the passed value at the named field in the named object.
	// It may return an error if the value is not of a compatible type,
	// if its considered to be read-only, or if there is no predeclared value of that name.
	SetField(object, field string, value g.Value) error
}

type Notary interface {
	// create a new field set from the named kind.
	// an implementation might only support the creation of some kinds.
	Make(kind string) (g.Value, error)
	// copy the contents of the passed value into a new value.
	Copy(g.Value) (g.Value, error)
	// move the contents of the passed value into a new value.
	// the fields of the old value are left at defaults
	// currently implemented as SetField with a null value
	// Move(g.Value) (g.Value, error)
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
	Notary
	// the runtime behaves as stack of scopes.
	// if a variable isnt found in the most recently pushed scope
	// the next most recently pushed scope will be checked and so on.
	Scope
	// add a set of variables to the internal stack.
	PushScope(Scope)
	// remove the most recently added set of variables from the internal stack.
	PopScope()
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

// WritersBlock applies a writer to the runtime for the duration of fn.
// If the writer also implements io.Closer and fn is free of errors,
// w.Close gets called and its result returned.
func WritersBlock(run Runtime, w writer.Output, fn func() error) (err error) {
	was := run.SetWriter(w)
	e := fn()
	run.SetWriter(was)
	if e != nil {
		err = e
	} else {
		err = writer.Close(w)
	}
	return
}
