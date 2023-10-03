package rules

// exports for testing
type UpdateTracker = updateTracker

// exports for testing
var Ranks = ranks

// exports for testing
func (n RuleName) RuleForPattern(ks Kinds) (ret RuleInfo, err error) {
	return n.ruleForPattern(ks)
}
