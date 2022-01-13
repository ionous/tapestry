// package kindsOf defines a handful of built-in names for various special kinds
package kindsOf

// default kinds
type Kinds int

//go:generate stringer -type=Kinds -linecomment
const None Kinds = 0 //

const (
	Aspect   Kinds = 1 << iota // aspects
	Kind                       // kinds
	Pattern                    // patterns
	Record                     // records
	Relation                   // relations
	Response                   // responses
	// subtypes of pattern:
	Action Kinds = Pattern | 1<<iota //actions
	Event                            //events
)

var DefaultKinds = []Kinds{
	Aspect, Kind, Pattern, Record, Relation, Response, Action, Event,
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
