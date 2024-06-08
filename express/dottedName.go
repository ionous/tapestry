package express

import (
	"git.sr.ht/~ionous/tapestry/dl/render"
	"git.sr.ht/~ionous/tapestry/rt"
)

// single dotted name supporting the printing of an object's name or
// a request for an object of a particular name.
// the full set of possibilities are:
//
//	{x}                  - print name of an object.
//	{kind_of: x}          - call a function with a local variable or object.
//	{print_plural_name: x} - call a user pattern with a local variable or object.
//
// where x can be, for instance:
//
//	{.target}   - a local variable named target
//	{.lantern}  - a variable or object named lantern
//	{.Lantern}  - an object named lantern
type dotName string

// when dotted names are used as arguments to concrete functions
//
//	ex. {numAsWords: .count}, {printPluralName: .target}
//
// we cant know the type of .count ( or even if it's a local or object )
// to know whether its a local variable or object,
//   - express would need a name stack during compilation
//
// to know what *type* its trying to render ( ex. text name of an object, or int from a local )
//   - express would need access to the known functions and patterns
func (on dotName) getNamedValue() *render.UnknownDot {
	return &render.UnknownDot{Name: T(string(on))}
}

// when dotted names are used directly:
//
//	ex {.target} or {.Lantern} or {.text}
//
// NameDot will attempt to read from the name as a variable,
// and if that fails, it will attempt to render the name as an object.
func (on dotName) getPrintedName() rt.TextEval {
	// the object.NameDot function itself handles the capitalization check
	// one thing missing here: if the text in a variable is not already an id
	// this will just print the text.
	return &render.RenderName{Name: string(on)}
}
