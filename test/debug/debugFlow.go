package debug

import (
	r "reflect"

	"git.sr.ht/~ionous/iffy/lang"
)

// Flow creates a debug flow to simulate ifspec generation
type Flow struct{ op interface{} }

// MakeFlow indicates the start of a set of key-value pairs.
// Unlike the other block types, the block itself is not mutable -- only its values.
func MakeFlow(op interface{}) Flow { return Flow{op} }

func (n Flow) GetType() string {
	return r.TypeOf(n.op).Elem().Name()
}

func (n Flow) GetLede() string {
	return lang.Breakcase(n.GetType())
}

func (n Flow) GetFlow() interface{} {
	return n.op
}

func (n Flow) SetFlow(i interface{}) (okay bool) {
	panic("not implemented... probably use reflection")
	return false
}