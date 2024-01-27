package play

import "git.sr.ht/~ionous/tapestry/lang/typeinfo"

// marker interface for play message types
type PlayMessage interface {
	PlayMessage()
	typeinfo.Instance
}

func (*PlayLog) PlayMessage()  {}
func (*PlayOut) PlayMessage()  {}
func (*PlayMode) PlayMessage() {}
