package play

type PlayMessage interface {
	PlayMessage()
}

func (*PlayLog) PlayMessage()  {}
func (*PlayOut) PlayMessage()  {}
func (*PlayMode) PlayMessage() {}
