package pattern

import (
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

// meta.KindAncestry
func Categorize(path []string) (ret Category) {
	switch len(path) {
	case 2: //  [] name, (record|pattern)
		switch path[len(path)-1] {
		case kindsOf.Record.String():
			ret = Initializes
		case kindsOf.Pattern.String():
			ret = Calls
		}
	case 3: // [] name, action, pattern
		switch path[len(path)-2] {
		case kindsOf.Action.String():
			ret = Sends
		}
	case 4: // [] name, action, pattern
		switch path[len(path)-2] {
		case kindsOf.Action.String():
			ret = Listens
		}
	}
	return
}
