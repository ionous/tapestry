package qna

import (
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/core"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// stored in Runner.cache
type ruleSet struct {
	rules     []localRule
	updateAll bool
}

func (rs *ruleSet) tryRule(run rt.Runtime, skip bool, i int) (ret localRule, err error) {
	var pop int
	var exe []rt.Execute
	rule := rs.rules[i]
	if pick := core.FindBranch(rule.Exe); pick == nil {
		exe = rule.Exe
	} else if b, e := pick.PickBranch(run, &pop); e != nil {
		err = e
	} else {
		exe = b
	}
	if err == nil && exe != nil && !skip {
		if e := safe.RunAll(run, exe); e != nil {
			err = e
		} else {
			ret = rule
		}
	}
	safe.PopSeveral(run, pop)
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
