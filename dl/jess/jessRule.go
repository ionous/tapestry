package jess

import (
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/rules"
)

func (op *TimedRule) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	/**/ op.RulePrefix.Match(q, &next) &&
		op.Pattern.Match(q, &next) &&
		(Optional(q, &next, &op.Target) || true) &&
		(Optional(q, &next, &op.RuleSuffix) || true) &&
		(op.matchName(&next) || true) &&
		op.SubAssignment.Match(&next) {
		*input, okay = next, true
	}
	return

}

// parenthetical name
func (op *TimedRule) matchName(input *InputState) (okay bool) {
	// if width := input.MatchWord(flexParens); width > 0 {
	// 	ws := input
	// 	op.RuleName = ws[0].String()
	// 	*input, okay = input.Skip(width), true
	// }
	return false
	return
}

// goal: schedule the rule
func (op *TimedRule) Generate(ctx Context) (err error) {
	if pat, e := op.Pattern.Validate(kindsOf.Pattern); e != nil {
		err = e
	} else {
		_ = pat
	}
	panic("not implemented")
}

func (op *RulePrefix) Match(q Query, input *InputState) (okay bool) {
	if idx, width := prefixes.FindPrefixIndex(input.Words()); width > 0 {
		op.PrefixValue = PrefixValue(idx)
		*input, okay = input.Skip(width), true
	}
	return
}

func (op *RuleSuffix) Match(q Query, input *InputState) (okay bool) {
	if idx, width := suffixes.FindPrefixIndex(input.Words()); width > 0 {
		op.SuffixValue = SuffixValue(idx)
		*input, okay = input.Skip(width), true
	}
	return
}

// assumes we are inside the parens
func (op *RuleName) MatchWords(next []match.TokenValue) (okay bool) {
	// // "this is"
	// if m, width := ruleNamePrefix.FindPrefix(next); m != nil {
	// 	op.Prefix = true
	// 	next = next[width:]
	// }
	// // "the"
	// if m, width := match.FindCommonArticles(next); width > 0 {
	// 	op.Article = m.String()
	// 	next = next[width:]
	// }
	// // ... "rule"
	// if cnt := len(next); cnt > 0 && next[cnt-1].Hash() == keywords.Rule {
	// 	op.Suffix = true
	// 	next = next[:cnt-1]
	// }
	// // and this is the rule name:
	// if len(next) > 0 {
	// 	op.Matched = match.JoinWords(next)
	// 	okay = true
	// }
	return
}

var flexParens = match.Hash("()")
var ruleNamePrefix = match.PanicSpans("this is")
var prefixes = make(match.SpanList, rules.NumPrefixes)
var suffixes = make(match.SpanList, rules.NumSuffixes)

type PrefixValue rules.Prefix
type SuffixValue rules.Suffix

func init() {
	for i := 0; i < rules.NumPrefixes; i++ {
		n := rules.Prefix(i)
		prefixes[i] = match.PanicSpan(n.String())
	}
	for i := 0; i < rules.NumSuffixes; i++ {
		n := rules.Suffix(i)
		suffixes[i] = match.PanicSpan(n.String())
	}
}
