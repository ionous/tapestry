package qna

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// stored in Runner.cache
type ruleSet struct {
	rules     []localRule
	updateAll bool
	skipRun   bool
}

func (rs *ruleSet) applyRule(run rt.Runtime, i int) (done bool, err error) {
	rule := rs.rules[i]
	if ok, e := safe.GetOptionalBool(run, rule.Filter, true); e != nil {
		err = e
	} else if ok.Bool() && !rs.skipRun {
		if e := safe.RunAll(run, rule.Exe); e != nil {
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
		var rules []localRule
		var updateAll bool
		for _, el := range els {

			if prog, e := run.decode.DecodeProg(el.Prog); e != nil {
				err = errutil.Append(err, errutil.New("decoding prog", pat, tgt, el.Name, e))
			} else {
				rules = append(rules, localRule{
					Name: el.Name,
					Prog: prog,
				})
				updateAll = updateAll || prog.Updates
			}
		}
		if err == nil {
			ret = ruleSet{rules: rules, updateAll: updateAll}
		}
	}
	return
}

// deserialized version of query rule
// tbd: hand deserialize interface to query so only one return type is needed?
type localRule struct {
	Name string
	assign.Prog
}
