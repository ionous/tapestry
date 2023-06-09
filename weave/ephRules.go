package weave

import (
	"git.sr.ht/~ionous/tapestry/weave/assert"

	"git.sr.ht/~ionous/tapestry/rt"
	"github.com/ionous/errutil"
)

// validate that the pattern for the rule exists then add the rule to the *current* domain
// ( rules are de/activated based on domain, they can be part some child of the domain where the pattern was defined. )
func (cat *Catalog) AssertRule(pattern string, target string, filter rt.BoolEval, flags assert.EventTiming, prog []rt.Execute) error {
	return cat.Schedule(assert.RulePhase, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		if name, ok := UniformString(pattern); !ok {
			err = InvalidString(pattern)
		} else if tgt, ok := UniformString(target); len(target) > 0 && !ok {
			err = errutil.Fmt("unknown or invalid target %q for pattern %q", target, pattern)
		} else {
			flags := fromTiming(flags)
			err = cat.writer.Rule(d.name, name, tgt, flags, filter, prog, at)
		}
		return
	})
}

func fromTiming(timing assert.EventTiming) int {
	var part int
	always := timing&assert.RunAlways != 0
	if always {
		timing ^= assert.RunAlways
	}
	switch timing {
	case assert.Before:
		part = 0
	case assert.During:
		part = 1
	case assert.After:
		part = 2
	case assert.Later:
		part = 3
	}
	flags := part + int(rt.FirstPhase)
	if always {
		flags = -flags // marker for rules that need to always run (ex. counters "every third try" )
	}
	return flags
}
