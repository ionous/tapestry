package action

//go:generate stringer -type=Field -linecomment
type Field int

// fields for actions
// the specified order is assumed by story grammar and the runtime
const (
	Noun          Field = iota // noun
	OtherNoun                  // other noun
	Actor                      // actor
	Target                     // target
	CurrentTarget              // current target
	Interupt                   // interrupt event
	Cancel                     // cancel event
)

func (af Field) Index() int {
	return int(af)
}

// predefined pattern for finding event targets
const CapturePattern = "event capture"

// in call order
func EventNames(name string) []string {
	return []string{"before " + name, name, "after " + name, "report " + name}
}
