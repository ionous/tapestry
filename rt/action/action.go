package action

import "git.sr.ht/~ionous/tapestry/rt/kindsOf"

// events raised for an action
type Event int

// in call order
const (
	BeforeEvent Event = iota
	AtEvent
	AfterEvent
	ReportEvent
	NumEvents
	FirstEvent = BeforeEvent
)

// fields for the events of actions
type Field int

//go:generate stringer -type=Field -linecomment
const (
	Noun      Field = iota // noun
	OtherNoun              // other noun
	Actor                  // actor
	Target                 // target
	Interupt               // interrupt event
	Cancel                 // cancel event
	NumFields = int(Cancel)
)

func (af Field) Index() int {
	return int(af)
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

// the opposite of bubbles is capture; and vice-versa.
func (f Flow) Reverse() Flow {
	return ^f & 1
}

// the specific name of an event depends on the action
// ex. "begin traveling"
func (evt Event) Name(action string) string {
	return evt.prefix() + action
}

func (evt Event) Prefix() string {
	return evt.prefix()
}

// the base class for the event
// ( used to determine when the call to a pattern needs to trigger a full-fledged event )
func (evt Event) Kind() (ret kindsOf.Kinds) {
	// fix these are reserved; can fix after all the old event handling is gone.
	if evt == AtEvent {
		ret = kindsOf.Event
	} else {
		ret = kindsOf.Action
	}
	return
}

// the order in which objects are visited when events are raised in the scene.
func (evt Event) Flow() (ret Flow) {
	switch evt {
	case BeforeEvent:
		ret = Captures
	case AfterEvent:
		ret = Bubbles
	}
	return
}

// friendly string for the event
func (evt Event) String() (ret string) {
	if evt == AtEvent {
		ret = "at event"
	} else {
		ret = evt.prefix() + "event"
	}
	return
}

func (evt Event) prefix() string {
	return prefixes[evt]
}

var prefixes = []string{"before ", "", "after ", "report "}

func _() {
	var assert [1]struct{}
	_ = assert[int(NumEvents)-len(prefixes)]
}
