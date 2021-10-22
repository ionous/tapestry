package rt

// Flags tweak the ordering of rules.
// Prefix rules get a chance to run before Infix rules, Infix run before Postfix.
// If a prefix rule decides to end the pattern, nothing else runs.
// If an infix rule decides to end the pattern, the postfix rules still trigger.
// The postfix rules run until one decides to end the pattern.
type Flags int

// Rule triggers a series of statements when its filters are satisfied.
// ( for backwards compatibility it doesnt directly aggregate Handler )
type Rule struct {
	Name     string
	RawFlags Flags
	Filter   BoolEval
	Execute
}

// Handler executes a statement its filter passes
type Handler struct {
	Filter BoolEval
	Exe    Execute
}

type NoResult struct{}

func (e NoResult) Error() string { return "no result" }

func (e NoResult) Is(target error) bool { return target == e }

func (e NoResult) NoPanic() {}

//go:generate stringer -type=Flags
const (
	Prefix  Flags = (1 << iota) // all prefix rules get sorted towards the front of the list
	Infix                       // keeps the rule at the same relative location
	Postfix                     // all postfix rules get sorted towards the end of the list
	After                       //
	Filter                      // internal flag to find (and always update) counters in rules
)

func (l Rule) Flags() (ret Flags) {
	if l.RawFlags == 0 {
		ret = Infix
	} else {
		ret = l.RawFlags
	}
	return
}

// Phase - return a semi-opaque integer, the absolute value of which can be sorted to get phase order.
func (f Flags) Phase() Phase {
	v := Phase(f.Ordinal())
	if f&Filter != 0 {
		v = -v
	}
	return v
}

// Ordinal - return the sort order of the flag.
func (f Flags) Ordinal() (ret int) {
	switch f & ^Filter {
	case Prefix:
		ret = 1
	case Infix:
		ret = 2
	case Postfix:
		ret = 3
	case After:
		ret = 4
	}
	return
}

// Phase is used for database storage:
// negative values are used to indicate the rule wants its filter to always be updated.
type Phase int // phase

const (
	FirstPhase Phase = 1
	LastPhase        = 4
	NumPhases        = LastPhase - FirstPhase + 1
)

//
func MakeFlags(p Phase) (ret Flags) {
	if p == 0 {
		ret = Infix
	} else {
		var flags Flags
		if p < 0 {
			p = -p
			flags = Filter
		}
		ret = flags | 1<<(p-1)
	}
	return
}
