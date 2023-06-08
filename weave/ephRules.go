package weave

import (
	"sort"

	"git.sr.ht/~ionous/tapestry/weave/assert"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/ionous/errutil"
)

func (c *Catalog) WriteRules(m mdl.Modeler) (err error) {
	if ds, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
	Done:
		// the rules for latter domains are more important than earlier domains so we put them first.
		// ( although sorting by the app, using materialized path or runtime domain sorting might be better )
		for rev := len(ds) - 1; rev >= 0; rev-- {
			d := ds[rev]
			// rules are stored by pattern name
			// we sort to give some consistency, although the order shouldnt really matter.
			names := make([]string, 0, len(d.rules))
			for k := range d.rules {
				names = append(names, k)
			}
			sort.Strings(names)
			for _, patternName := range names {
				rs := d.rules[patternName]
				for j, p := range rs.partitions {
					// write individual rules in reverse order: the last specified is the most important
					for i := len(p.els) - 1; i >= 0; i-- {
						el, at := p.els[i], p.at[i]
						//
						flags := j + int(rt.FirstPhase)
						if el.Touch {
							flags = -flags // marker for rules that need to always run (ex. counters "every third try" )
						}
						if e := m.Rule(d.name, patternName, el.Target, flags, el.Filter, el.Prog, at); e != nil {
							err = e
							break Done
						}
					}
				}
			}
		}
	}
	return
}

type Rulesets struct {
	partitions [rt.NumPhases]Partition
}

type Partition struct {
	els []ephRules
	at  []string
}

type ephRules struct {
	Target string
	Filter rt.BoolEval
	Prog   []rt.Execute
	Touch  bool
}

func fromTiming(t assert.EventTiming) (ret int, always bool) {
	if always = t&assert.RunAlways != 0; always {
		t ^= assert.RunAlways
	}
	switch t {
	case assert.Before:
		ret = 0
	case assert.During:
		ret = 1
	case assert.After:
		ret = 2
	case assert.Later:
		ret = 3
	}
	return
}

// validate that the pattern for the rule exists then add the rule to the *current* domain
// ( rules are de/activated based on domain, they can be part some child of the domain where the pattern was defined. )
func (cat *Catalog) AssertRule(opPatternName string, opTarget string, opGuard rt.BoolEval, opFlags assert.EventTiming, do []rt.Execute) error {
	return cat.Schedule(assert.RulePhase, func(ctx *Weaver) (err error) {
		d, at := ctx.d, ctx.at
		part, always := fromTiming(opFlags)
		if name, ok := UniformString(opPatternName); !ok {
			err = InvalidString(opPatternName)
		} else if tgt, ok := UniformString(opTarget); len(opTarget) > 0 && !ok {
			err = errutil.Fmt("unknown or invalid target %q for pattern %q", opTarget, opPatternName)
		} else {
			if d.rules == nil {
				d.rules = make(map[string]Rulesets)
			}
			rules := d.rules[name]
			p := &rules.partitions[part]
			p.els = append(p.els, ephRules{
				Target: tgt, Filter: opGuard, Prog: do, Touch: always,
			})
			p.at = append(p.at, at)
			d.rules[name] = rules
		}
		return
	})
}
