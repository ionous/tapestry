package play

// expected play time patterns
const (
	RequestingPlayerInput = "requesting player input"
	StartGame             = "start game"
	PassTime              = "pass time"
	RunningAnAction       = "running an action"
)

// expected parameter names
const (
	Actor      = "actor"
	Action     = "action"
	FirstNoun  = "first noun"
	SecondNoun = "second noun"
)

// convention for actions which don't advance time
const OutOfWorldPrefix = "request "
