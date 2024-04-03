// package kindsOf defines a handful of built-in base types for tapestry kinds.
// a kind is a sharable set of author defined fields used by the runtime.
package kindsOf

// default kinds
type Kinds int

//go:generate stringer -type=Kinds -linecomment
const None Kinds = 0 //

// note: iota increases even when not specified so every kind has a unique bit;
// re-specifying it combines with earlier kinds to indicate sub-types.
// ( Kind is an abstract or concrete noun defined by some story. )
const (
	// a set of exclusive boolean fields which apply to a noun;
	// a kind supports many aspects, a noun can only have one field of each aspect.
	Aspect Kinds = 1 << iota // aspects

	// a base for all abstract and concrete nouns:
	// nouns are globally accessible objects identified by unique names.
	Kind // kinds

	// used for user defined functions.
	// the fields of a pattern define the parameters and return values.
	Pattern // patterns

	// a base for any structured data type held by variables and kinds.
	Record // records

	// describe how nouns pair with other nouns.
	// the fields of relation define the number and kind of nouns on each side of the pairing.
	Relation // relations

	// a single kind, the fields of which define stock phrases relayed to the player.
	// ex. "can't act in the dark" -> "It is pitch dark, and you can't see a thing."
	// ( patterns have parameters and return values; responses don't. )
	Response // responses

	// a special kind of pattern used for handling events.
	Action Kinds = Pattern | 1<<iota // actions
)

var DefaultKinds = []Kinds{
	Aspect, Kind, Pattern, Record, Relation, Response, Action,
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

// given a string, return the default kind that it exactly matches (if any)
func FindDefaultKind(str string) (ret Kinds) {
	for _, k := range DefaultKinds {
		if str == k.String() {
			ret = k
			break
		}
	}
	return
}
