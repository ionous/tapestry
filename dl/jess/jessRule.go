package jess

func (op *TimedRule) Generate(ctx Context) error {
	panic("x")
}

func (op *TimedRule) Match(q Query, input *InputState) (okay bool) {
	return false
}

func (op *RulePrefix) Match(q Query, input *InputState) (okay bool) {
	return false
}

func (op *RuleSuffix) Match(q Query, input *InputState) (okay bool) {
	return false
}

func (op *ShortRuleName) Match(q Query, input *InputState) (okay bool) {
	return false
}

func (op *LongRuleName) Match(q Query, input *InputState) (okay bool) {
	return false
}
