package qna

import (
	"git.sr.ht/~ionous/tapestry/rt/pattern"
)

func (run *Runner) getRules(name string) (ret pattern.RuleSet, err error) {
	if rs, e := run.query.RulesFor(name); e != nil {
		err = e
	} else {
		ret = pattern.MakeRules(rs.Rules, rs.UpdateAll)
	}
	return
}
