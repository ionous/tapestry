package jess

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/rules"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

func (op *TimedRule) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	/**/ op.RulePrefix.Match(q, &next) &&
		(op.matchSomeone(q, &next) || true) &&
		op.Pattern.Match(q, &next) &&
		(Optional(q, &next, &op.RuleTarget) || true) &&
		(Optional(q, &next, &op.RuleSuffix) || true) &&
		(op.matchName(q, &next) || true) &&
		op.SubAssignment.Match(&next) {
		*input, okay = next, true
	}
	return
}

func (op *TimedRule) matchSomeone(_ Query, input *InputState) (okay bool) {
	if width := input.MatchWord(keywords.Someone); width > 0 {
		op.Someone = true
		*input, okay = input.Skip(width), true
	}
	return
}

// parenthetical name
func (op *TimedRule) matchName(q Query, input *InputState) (okay bool) {
	if val, ok := input.GetNext(match.Parenthetical); ok {
		// fix: it might make more sense to handle open/close parens separately
		// ... there also might be something clever here with rule_name matched
		// ex. caching until Generate and then validating there are close parens
		// no terminals, etc. in Generate....
		if words, e := match.Tokenize(val.String()); e != nil {
			// on error tokeninzing, set an empty rule name and error in Generate
			op.RuleName = new(RuleName)
			*input, okay = input.Skip(1), true
		} else {
			words := InputState(words)
			if Optional(q, &words, &op.RuleName) {
				*input, okay = input.Skip(1), true
			}
		}
	}
	return
}

// goal: schedule the rule
func (op *TimedRule) Generate(ctx Context) (err error) {
	if label, e := GetRuleName(op.RuleName); e != nil {
		err = e
	} else if pat, e := op.Pattern.Validate(kindsOf.Pattern, kindsOf.Action); e != nil {
		err = e
	} else if op.Someone && op.Pattern.ActualKind.BaseKind != kindsOf.Action {
		// fix: story had an additional test CanFilterActor. hrm.
		err = errors.New("can only filter actor actions")
	} else if exe, ok := op.SubAssignment.GetExe(); !ok {
		err = errors.New("rule expected a list of statements to execute")
	} else {
		err = ctx.Schedule(weaver.VerbPhase, func(w weaver.Weaves, run rt.Runtime) (err error) {
			if ks, e := run.GetField(meta.KindAncestry, pat); e != nil {
				err = e
			} else {
				n := rules.RuleName{
					Pattern: ks.Strings(),
					Label:   label,
					Prefix:  GetRulePrefix(op.RulePrefix),
					Suffix:  GetRuleSuffix(op.RuleSuffix),
				}
				var filters []rt.BoolEval

				// add custom filters ( if any )
				if tgt := op.RuleTarget; tgt != nil {
					if n := tgt.Noun; n != nil {
						filters = rules.AddNounFilter(n.ActualNoun.Name, filters)
					} else if k := tgt.Kind; k != nil {
						filters = rules.AddKindFilter(k.ActualKind.Name, filters)
					} else {
						panic("unhandled target match")
					}
				}

				// add other action filters
				if op.Pattern.ActualKind.BaseKind == kindsOf.Action {
					if !op.Someone {
						// all actions are triggered by actors;
						// then we automatically filter for the player
						filters = rules.AddPlayerFilter(filters)
					}
					// by default: all event handlers are filtered to the innermost target.
					filters = rules.AddEventFilters(filters)
				}

				if rule, e := n.GetRuleInfo(); e != nil {
					err = e
				} else {
					err = rule.WeaveRule(w, nil, exe)
				}
			}
			return
		})
	}
	return
}

// ----

func (op *SubAssignment) Match(input *InputState) (okay bool) {
	if a, ok := input.GetNext(match.Tell); ok {
		// for now panic to catch setup errors;
		// if there is a good reason to fail, could test for success.
		op.Assignment = a.Value.(rt.Assignment)
		*input, okay = input.Skip(1), true
	}
	return
}

func (op *SubAssignment) GetExe() (ret []rt.Execute, okay bool) {
	if a, ok := op.Assignment.(*assign.FromExe); ok {
		ret, okay = a.Exe, true
	}
	return
}

// ----

func (op *RuleTarget) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	Optional(q, &next, &op.Kind) ||
		Optional(q, &next, &op.Noun) {
		*input, okay = next, true
	}
	return
}

// ----

func GetRulePrefix(op RulePrefix) rules.Prefix {
	return rules.Prefix(op.PrefixValue)
}

func (op *RulePrefix) Match(q Query, input *InputState) (okay bool) {
	if idx, width := prefixes.FindPrefixIndex(input.Words()); width > 0 {
		op.PrefixValue = PrefixValue(idx)
		*input, okay = input.Skip(width), true
	}
	return
}

// ----

func GetRuleSuffix(op *RuleSuffix) (ret rules.Suffix) {
	if op != nil {
		ret = rules.Suffix(op.SuffixValue)
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

// ------

func GetRuleName(op *RuleName) (ret string, err error) {
	if op != nil {
		if str, w := match.Normalize(op.Matched); w != len(op.Matched) || w == 0 {
			err = fmt.Errorf("couldn't determine rule name")
		} else {
			ret = str
		}
	}
	return
}

// assumes we are inside the parens
func (op *RuleName) Match(q Query, input *InputState) (okay bool) {
	next := *input
	// "this is"
	if m, width := ruleNamePrefix.FindPrefix(next); m != nil {
		op.Prefix = true
		next = next[width:]
	}
	// "the"
	if m, width := match.FindCommonArticles(next); width > 0 {
		op.Article = m.String()
		next = next[width:]
	}
	// ... "rule"; tbd, add a LastIndexOf?
	if cnt := len(next); cnt > 0 && next[cnt-1].Hash() == keywords.Rule {
		op.Suffix = true
		next = next[:cnt-1]
	}
	// and this is the actual rule name:
	if len(next) > 0 {
		op.Matched = next.Words()
		okay = true
	}
	return
}

// ------

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
