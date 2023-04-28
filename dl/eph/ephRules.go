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

func (rs *Rulesets) AppendRule(el *EphRules, tgt string, part int, at string) (err error) {
	slice := rt.Execute_Slice(el.Exe)
	if filter, e := marshalout(el.Filter); e != nil {
		err = e
	} else if prog, e := marshalout(&slice); e != nil {
		err = e
	} else {
		p := &rs.partitions[part]
		p.els = append(p.els, ephRules{
			Target: tgt, Filter: filter, Prog: prog, Touch: len(el.Touch.String()) > 0,
		})
		p.at = append(p.at, at)
	}
	return
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
func (op *EphRules) Assemble(ctx *Context) (err error) {
	d, at := ctx.d, ctx.at
	if name, ok := UniformString(op.PatternName); !ok {
		err = InvalidString(op.PatternName)
	} else if k, ok := d.GetKind(name); !ok || !k.HasAncestor(kindsOf.Pattern) {
		err = errutil.Fmt("unknown or invalid pattern %q", op.PatternName)
	} else if part, ok := op.When.GetPartition(); !ok {
		err = errutil.Fmt("couldn't compute flags for %q for pattern %q", op.When.Str, op.PatternName)
	} else if tgt, ok := op.getTargetName(d); !ok {
		err = errutil.Fmt("unknown or invalid target %q for pattern %q", op.Target, op.PatternName)
	} else {
		if d.rules == nil {
			d.rules = make(map[string]Rulesets)
		}
		rules := d.rules[name]
		if e := rules.AppendRule(op, tgt, part, at); e != nil {
			err = e
		} else {
			d.rules[name] = rules
		}
	}
	return
}

func (op *EphRules) getTargetName(d *Domain) (ret string, okay bool) {
	if tgt := op.Target; len(tgt) == 0 {
		okay = true
	} else if k, ok := d.GetKind(op.Target); ok {
		ret, okay = k.name, true
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
