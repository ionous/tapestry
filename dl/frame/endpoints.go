package frame

// Destinations for POST
// the same endpoint may respond in different ways depending on the game State(s)
//
//go:generate stringer -type=Endpoint -linecomment
type Endpoint int

// maybe also "$weave", "$shutdown" ...
// note: instead of sending "input" we send a fabricate command.....
const (
	Restart Endpoint = iota // restart
	Query                   // query
)
