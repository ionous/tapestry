package story

import (
	"git.sr.ht/~ionous/iffy/rt"
)

// generates rules based on guards
type programHook interface {
	SlotType() string
	// gob requirement: return a pointer to the interface
	CmdPtr() interface{}
	// create a "pattern rule"
	// each rule returns its own kind of value -- so there's currently no common interface
	NewRule(guard rt.BoolEval, flags rt.Flags) (string, interface{})
}

type executeSlot struct{ cmd rt.Execute }

func (b *executeSlot) SlotType() string {
	return "execute"
}
func (b *executeSlot) NewRule(guard rt.BoolEval, flags rt.Flags) (string, interface{}) {
	return "rule", &rt.Rule{Filter: guard, Execute: b.cmd, Flags: flags}
}
func (b *executeSlot) CmdPtr() interface{} {
	return &b.cmd
}
