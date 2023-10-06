package shuttle

// Destinations for POST
// ( different State(s) may respond to in different ways to the same endpoint
//  not very RESTful... fine for now )
//go:generate stringer -type=Endpoint -linecomment
type Endpoint int

// maybe also "$weave", "$shutdown" ...
// note: instead of sending "input" we send a fabricate command.....
const (
	Restart Endpoint = iota // restart
	Query                   // query
)
