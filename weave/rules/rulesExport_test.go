package rules

import g "git.sr.ht/~ionous/tapestry/rt/generic"

// exports for testing
type UpdateTracker = updateTracker

// exports for testing
var Ranks = ranks

// exports for testing
func (n RuleName) RuleForPattern(ks g.Kinds) (ret RuleInfo, err error) {
	return n.ruleForPattern(ks)
}
