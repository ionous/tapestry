// package kindsOf defines a handful of built-in names for various special kinds
package kindsOf

// default kinds
type Kinds int

//go:generate stringer -type=Kinds -linecomment
const None Kinds = 0 //

// note: iota increases even when not specified so every kind has a unique bit;
// re-specifying it combines with earlier kinds to indicate sub-types.
const (
	Aspect   Kinds = 1 << iota         // aspects
	Kind                               // kinds
	Pattern                            // patterns
	Record                             // records
	Relation                           // relations
	Response                           // responses
	Macro    Kinds = Pattern | 1<<iota // macros
	Action   Kinds = Pattern | 1<<iota // actions
)

var DefaultKinds = []Kinds{
	Aspect, Kind, Pattern, Record, Relation, Response, Action, Macro,
}

func (k Kinds) Parent() (ret Kinds) {
	for _, p := range DefaultKinds {
		if k != p && (k&p == p) {
			ret = p
			break
		}
	}
	return
}
