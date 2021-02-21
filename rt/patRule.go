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
	Name string
	Flags
	Filter BoolEval
	Execute
}

// Handler executes a statement its filter passes
type Handler struct {
	Filter BoolEval
	Exe    Execute
}

func (my *Rule) GetFlags() (ret Flags) {
	ret = my.Flags
	if ret == 0 {
		ret = Infix
	}
	return
}

//go:generate stringer -type=Flags
const (
	Prefix  Flags = (1 << iota) // all prefix rules get sorted towards the front of the list
	Infix                       // keeps the rule at the same relative location
	Postfix                     // all postfix rules get sorted towards the end of the list
	After
)

func (f Flags) Ordinal() (ret int) {
	switch f {
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

func MakeFlags(i int) (ret Flags) {
	if i > 0 {
		ret = 1 << (i - 1)
	}
	return
}
