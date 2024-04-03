package pattern

import (
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
)

//go:generate stringer -type=Category
type Category int

const (
	invalid     Category = iota
	Initializes          // some sort of record
	Calls                // some sort of callable pattern
	Sends                // some sort of pattern action
	Listens
)

func Categorize(k *g.Kind) (ret Category) {
	switch path := g.Ancestry(k); len(path) {
	case 2: // <record>(0) name; <pattern>(0) name
		switch base := path[0]; base {
		case kindsOf.Record.String():
			ret = Initializes
		case kindsOf.Pattern.String():
			ret = Calls
		}
	case 3: // <pattern>(1) <action> name
		switch base := path[1]; base {
		case kindsOf.Action.String():
			ret = Sends
		}
	case 4: // <pattern> <action>(1) action_name event_name
		switch base := path[1]; base {
		case kindsOf.Action.String():
			ret = Listens
		}
	}
	return
}
