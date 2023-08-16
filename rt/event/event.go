package event

import (
	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
)

// name of the (implicit) actor variable (that should be) defined by all events.
const Actor = "actor"

// name of the event object variable
// the fields of the object are listed as Field constants
const Object = "event" // n

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
	ReportPhase
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

	// setting this stops the rest of the current flow ( ex. bubbling )
	// clearing this stops lets other handlers on this target finish first.
	Interupt // interrupt

	// setting this interrupts all further processing;
	// clearing this allows the currently flow to finish first.
	Cancel // cancel
	//
	NumFields = iota
)

func (f Field) Index() int {
	return int(f)
}

func (f Field) Affine() (ret affine.Affinity) {
	if f == Interupt || f == Cancel {
		ret = affine.Bool
	} else {
		ret = affine.Text
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

// friendly string for the phase; this isn' the same as the prefix
func (p Phase) String() (ret string) {
	if p == TargetPhase {
		ret = "at event"
	} else {
		ret = p.prefix() + "event"
	}
	return
}

// event phases are represented by patterns following a certain naming convention:
// return the  string ( including separator ) used to differentiate the patterns of this phase.
func (p Phase) Prefix() string {
	return p.prefix()
}

// event phases are represented by separate patterns following a certain naming convention;
// each phase has a unique prefix followed by the name of the action;
// return that full, unique name:  ex. "begin traveling"
func (p Phase) PatternName(action string) string {
	return p.prefix() + action
}

// event phases are represented by patterns of different sub-types.
// ( used to determine when the call to a pattern needs to trigger a full-fledged event )
func (p Phase) PatternKind() (ret kindsOf.Kinds) {
	// fix these are reserved; can fix after all the old event handling is gone.
	if p == TargetPhase {
		ret = kindsOf.Action
	} else {
		ret = kindsOf.Event
	}
	return
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

var prefixes = []string{"before ", "", "after ", "report "}

func _() {
	var assert [1]struct{}
	_ = assert[int(NumPhases)-len(prefixes)]
}
