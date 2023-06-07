// package kindsOf defines a handful of built-in names for various special kinds
package kindsOf

// default kinds
type Kinds int

//go:generate stringer -type=Kinds -linecomment
const None Kinds = 0 //

const (
	Aspect   Kinds = 1 << iota // aspect
	Kind                       // kind
	Macro                      // macro
	Pattern                    // pattern
	Record                     // record
	Relation                   // relation
	Response                   // response
	// subtypes of pattern:
	Action Kinds = Pattern | 1<<iota //action
	Event                            //event
)

var DefaultKinds = []Kinds{
	Aspect, Kind, Macro, Pattern, Record, Relation, Response, Action, Event,
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
