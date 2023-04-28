package eph

import (
	"sort"

	"git.sr.ht/~ionous/tapestry/imp/assert"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/tables/mdl"
	"github.com/ionous/errutil"
)

func (c *Catalog) WriteRules(w Writer) (err error) {
	if ds, e := c.ResolveDomains(); e != nil {
		err = e
	} else {
	Done:
		// the rules for latter domains are more important than earlier domains so we put them first.
		// ( although sorting by the app, using materialized path or runtime domain sorting might be better )
		for rev := len(ds) - 1; rev >= 0; rev-- {
			deps := ds[rev]
			d := deps.Leaf().(*Domain)
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
						if e := w.Write(mdl.Rule, d.name, patternName, el.Target, flags, el.Filter, el.Prog, at); e != nil {
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
	Filter string
	Prog   string
	Touch  bool
}

// rules are assembled after kinds and their fields...
func (op *EphRules) Phase() assert.Phase { return assert.PatternPhase }

func (op *EphRules) Weave(k assert.Assertions) (err error) {
	flags := toTiming(op.When, op.Touch)
	return k.AssertRule(op.PatternName, op.Target, op.Filter, flags, op.Exe)
}

// validate that the pattern for the rule exists then add the rule to the *current* domain
// ( rules are de/activated based on domain, they can be part some child of the domain where the pattern was defined. )
func (ctx *Context) AssertRule(opPatternName string, opTarget string, opGuard rt.BoolEval, opFlags assert.EventTiming, do []rt.Execute) (err error) {
	d, at := ctx.d, ctx.at
	when, always := fromTiming(opFlags)
	if name, ok := UniformString(opPatternName); !ok {
		err = InvalidString(opPatternName)
	} else if k, ok := d.GetKind(name); !ok || !k.HasAncestor(kindsOf.Pattern) {
		err = errutil.Fmt("unknown or invalid pattern %q(%s)", opPatternName, name)
	} else if part, ok := when.GetPartition(); !ok {
		err = errutil.Fmt("couldn't compute flags for %q for pattern %q", when.Str, opPatternName)
	} else if tgt, ok := getTargetName(d, opTarget); !ok {
		err = errutil.Fmt("unknown or invalid target %q for pattern %q", opTarget, opPatternName)
	} else {
		if d.rules == nil {
			d.rules = make(map[string]Rulesets)
		}
		rules := d.rules[name]
		slice := rt.Execute_Slice(do)
		if filter, e := marshalout(opGuard); e != nil {
			err = e
		} else if prog, e := marshalout(&slice); e != nil {
			err = e
		} else {
			p := &rules.partitions[part]
			p.els = append(p.els, ephRules{
				Target: tgt, Filter: filter, Prog: prog, Touch: len(always.String()) > 0,
			})
			p.at = append(p.at, at)
			d.rules[name] = rules
		}
	}
	return
}

func getTargetName(d *Domain, opTarget string) (ret string, okay bool) {
	if tgt := opTarget; len(tgt) == 0 {
		okay = true
	} else if x, ok := UniformString(opTarget); ok {
		if k, ok := d.GetKind(x); ok {
			ret, okay = k.name, true
		}
	}
	return
}

func (op *EphTiming) GetPartition() (ret int, okay bool) {
	// probably shouldnt but just rely on declared order for now.
	spec := op.Compose()
	if _, i := spec.IndexOfChoice(op.Str); i >= 0 {
		ret, okay = i, true
	}
	return
}
