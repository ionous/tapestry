package event

import (
	"git.sr.ht/~ionous/tapestry/affine"
)

// name of the optional event target parameter
const Actor = "actor"

// name of kind which incorporates the player
const Actors = "actors"

// name of the event object variable
// the fields of the object are listed as Field constants
const Object = "event"

// the individual events raised for an action
// they are represented collectively by a "pattern set"
// each pattern in the set corresponding to a phase
// using naming convention and pattern sub-types to identify which is which.
type Phase int

// phases are listed in call order
const (
	BeforePhase Phase = iota
	TargetPhase
	AfterPhase
	//
	NumPhases = iota
)

// fields for the event object
type Field int

//go:generate stringer -type=Field -linecomment
const (
	// name of the action
	Name Field = iota // name

	// target of the event ( used to determine event flow )
	Target // target

	// the object currently handling the event
	CurrentTarget // current target

	// a string containing the current cancellation status
	Status // status

	//
	NumFields = iota
)

// traits controlling event processing
// higher numbers indicate a higher severity.
type CancellationStatus int

//go:generate stringer -type=CancellationStatus -linecomment
const (
	// keep processing other handlers and other flows.
	ContinueNormally CancellationStatus = iota // continue normally

	// stops the current flow after other handlers finish.
	InterruptLater // interrupt later

	// stops the current flow without letting other handlers run.
	InterruptNow // interrupt now

	// stops all flows after the current flow.
	CancelLater // cancel later

	// stops all further processing.
	CancelNow // cancel now

	//
	NumStatus = iota
)

func (f Field) Index() int {
	return int(f)
}

func (f Field) Affine() (ret affine.Affinity) {
	return affine.Text
}

func (f Field) Type() (ret string) {
	if f == Status {
		ret = f.String()
	}
	return
}

// predefined pattern for finding event targets
const CapturePattern = "capture"

// behavior of a given event type
type Flow int

//go:generate stringer -type=Flow
const (
	// a direct call to a pattern
	Targets Flow = iota
	// patterns are raised from root object of the scene down to the targeted object.
	Captures
	// patterns are raised from the targeted object up to the root of the scene.
	Bubbles
)

// friendly string for the phase
func (p Phase) String() (ret string) {
	return p.prefix() + "event"
}

// event phases are represented by separate patterns following a certain naming convention;
// each phase has a unique prefix followed by the name of the action;
// return that full, unique name:  ex. "begin traveling"
func (p Phase) PatternName(action string) string {
	return p.prefix() + action
}

// determines the order to visit objects for a particular phase.
func (p Phase) Flow() (ret Flow) {
	switch p {
	case BeforePhase:
		ret = Captures
	case AfterPhase:
		ret = Bubbles
	default:
		ret = Targets
	}
	return
}

// raw prefix including separator
func (p Phase) prefix() string {
	return prefixes[p]
}

var prefixes = []string{"before ", "", "after "}

func _() {
	var assert [1]struct{}
	_ = assert[int(NumPhases)-len(prefixes)]
}
