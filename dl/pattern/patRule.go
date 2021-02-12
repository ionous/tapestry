package pattern

import (
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/safe"
)

// Flags tweak the ordering of rules.
// Prefix rules get a chance to run before Infix rules, Infix run before Postfix.
// If a prefix rule decides to end the pattern, nothing else runs.
// If an infix rule decides to end the pattern, the postfix rules still trigger.
// The postfix rules run until one decides to end the pattern.
type Flags int

//go:generate stringer -type=Flags
const (
	Infix   Flags = (1 << iota) // keeps the rule at the same relative location
	Prefix                      // all prefix rules get sorted towards the front of the list
	Postfix                     // all postfix rules get sorted towards the end of the list
)

// Rule triggers a series of statements when its filters are satisfied.
type Rule struct {
	Filter rt.BoolEval
	Flags
	rt.Execute
}

func (my *Rule) GetFlags() (ret Flags) {
	ret = my.Flags
	if ret == 0 {
		ret = Infix
	}
	return
}
func (my *Rule) ApplyRule(run rt.Runtime, allow Flags) (okay Flags, err error) {
	if flags := my.GetFlags(); allow&flags != 0 {
		if ok, e := safe.GetOptionalBool(run, my.Filter, true); e != nil {
			err = e
		} else if ok.Bool() { // the rule returns false if it didnt apply
			if e := safe.Run(run, my.Execute); e != nil {
				err = e
			} else {
				okay = flags
			}
		}
	}
	return
}
