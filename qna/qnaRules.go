package qna

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/pattern"
	"github.com/ionous/errutil"
)

// get the rules from the cache, or build them and add them to the cache
func (run *Runner) getRules(pat string) (ret pattern.RuleSet, err error) {
	if c, e := run.values.cache(func() (any, error) {
		return run.buildRules(pat)
	}, "rules", pat); e != nil {
		err = e
	} else {
		ret = c.(pattern.RuleSet)
	}
	return
}

// build the rules for the passed pat,tgt pair
func (run *Runner) buildRules(pat string) (ret pattern.RuleSet, err error) {
	if els, e := run.query.RulesFor(pat); e != nil {
		err = e
	} else {
		var rs pattern.RuleSet
		for _, el := range els {
			if exe, e := run.decode.DecodeProg(el.Prog); e != nil {
				err = errutil.New("decoding prog", pat, el.Name, e)
				break
			} else {
				rs.AddRule(rt.Rule{
					Name:    el.Name,
					Stop:    el.Stop,
					Jump:    rt.Jump(el.Jump),
					Updates: el.Updates,
					Exe:     exe,
				})
			}
		}
		if err == nil {
			ret = rs
		}
	}
	return
}
