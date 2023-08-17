package qna

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// stored in Runner.cache
type ruleSet struct {
	rules     []rt.Rule
	updateAll bool
	skipRun   bool
}

func (rs *ruleSet) applyRule(run rt.Runtime, i int) (done bool, err error) {
	rule := rs.rules[i]
	if ok, e := safe.GetOptionalBool(run, rule.Filter, true); e != nil {
		err = e
	} else if ok.Bool() && !rs.skipRun {
		if e := safe.RunAll(run, rule.Execute); e != nil {
			err = e
		} else if rule.Terminates {
			if !rs.updateAll {
				done = true
			}
			rs.skipRun = true
		}
	}
	return
}

// get the rules from the cache, or build them and add them to the cache
func (run *Runner) getRules(pat, tgt string) (ret ruleSet, err error) {
	if c, e := run.values.cache(func() (any, error) {
		return run.buildRules(pat, tgt)
	}, "rules", pat, tgt); e != nil {
		err = e
	} else {
		ret = c.(ruleSet)
	}
	return
}

// build the rules for the passed pat,tgt pair
func (run *Runner) buildRules(pat, tgt string) (ret ruleSet, err error) {
	if els, e := run.query.RulesFor(pat, tgt); e != nil {
		err = e
	} else {
		var rules []rt.Rule
		var updateAll bool
		for _, el := range els {
			if filter, e := run.decode.DecodeFilter(el.Filter); e != nil {
				err = errutil.Append(err, errutil.New("decoding filter", pat, tgt, el.Id, e))
			} else if prog, e := run.decode.DecodeProg(el.Prog); e != nil {
				err = errutil.Append(err, errutil.New("decoding prog", pat, tgt, el.Id, e))
			} else {
				rules = append(rules, rt.Rule{
					Name:       el.Id,
					Filter:     filter,
					Execute:    prog,
					Updates:    el.Updates,
					Terminates: el.Terminates,
				})
				updateAll = updateAll || el.Updates
			}
		}
		if err == nil {
			ret = ruleSet{rules: rules, updateAll: updateAll}
		}
	}
	return
}
