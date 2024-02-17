package jess

import (
	"git.sr.ht/~ionous/tapestry/support/grok"
)

func (op *VerbLinks) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Verb.Match(q, &next) &&
		op.Names.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.OtherNames.Match(q, &next) {
		q.note("matched VerbLinks")
		*input, okay = next, true
	}
	return
}

func (op *VerbLinks) GetResults() (ret grok.Results, _ error) {
	ret = makeResult(
		op.Verb.Macro, !op.Verb.Macro.Reversed,
		op.Names.GetNames(nil, nil),
		op.OtherNames.GetNames(nil, nil),
	)
	return
}

func (op *LinksVerb) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Names.Match(q, &next) &&
		!op.Names.HasAnonymousKind() &&
		op.Are.Match(q, &next) &&
		op.Verb.Match(q, &next) &&
		op.OtherNames.Match(q, &next) {
		q.note("matched LinksVerb")
		*input, okay = next, true
	}
	return
}

func (op *LinksVerb) GetResults() (ret grok.Results, _ error) {
	ret = makeResult(
		op.Verb.Macro,
		op.Verb.Macro.Reversed,
		op.Names.GetNames(nil, nil),
		op.OtherNames.GetNames(nil, nil),
	)
	return
}

func (op *LinksAdjectives) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Names.Match(q, &next) &&
		!op.Names.HasAnonymousKind() &&
		op.Are.Match(q, &next) &&
		op.Adjectives.Match(q, &next) {
		q.note("matched LinksAdjectives")
		Optional(q, &next, &op.VerbPhrase)
		*input, okay = next, true
	}
	return
}

func (op *LinksAdjectives) GetOtherNames() (ret Iterator) {
	if v := op.VerbPhrase; v != nil {
		ret = v.Names.Iterate()
	}
	return
}

func (op *LinksAdjectives) GetMacro() (ret grok.Macro) {
	if v := op.VerbPhrase; v != nil {
		ret = v.Verb.Macro
	}
	return
}

func (op *LinksAdjectives) IsReversed() (okay bool) {
	if v := op.VerbPhrase; v != nil {
		okay = v.Verb.Macro.Reversed
	}
	return
}

func (op *LinksAdjectives) GetResults() (ret grok.Results, _ error) {
	var b []grok.Name
	for it := op.GetOtherNames(); it.HasNext(); {
		n := it.GetNext()
		b = append(b, n.GetName(nil, nil))
	}
	a := op.Names.GetNames(op.Adjectives.Reduce())
	ret = makeResult(
		op.GetMacro(), op.IsReversed(),
		a, b)

	return
}
