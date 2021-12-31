package parser

import (
	"git.sr.ht/~ionous/tapestry/parser/ident"
)

type Context interface {
	// return true if the passed word is a plural word.
	IsPlural(word string) bool
	// the range of the player's known universe.
	// string names are defined by the "Focus" parts of a Scanner grammar.
	// for example, maybe "held" for objects held by the player.
	// the empty string is used as the default range when no focus has been declared.
	GetPlayerBounds(string) (Bounds, error)
	// ex. take from a container
	GetObjectBounds(ident.Id) (Bounds, error)
}

// Bounds encapsulates some set of objects.
// Searches visits every object in the set defined by the bounds.
// note: we use a visitor to support map traversal without copying keys if need be.
type Bounds func(NounVisitor) bool

// If the visitor function returns true, the search terminates and returns true;
// otherwise it returns false.
type NounVisitor func(NounInstance) bool

// NounInstance - allows parser to ask questions about a particular object.
// fix? it might be nicer for callers if these methods were part of context
//
type NounInstance interface {
	// Id for the noun. Returned via ResultList.Objects() on a successful match.
	Id() ident.Id
	// does the passed plural string apply to this object?
	// low-bar would be to return the same result as class,
	// better might be looking at plural printed name.
	HasPlural(string) bool
	// does the passed name apply to this object?
	HasName(string) bool
	// does the noun satisfy the passed named class?
	HasClass(string) bool
	// does the noun have the passed name attribute?
	HasAttribute(string) bool
}
