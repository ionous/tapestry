package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

func (op *Nouns) Match(q Query, input *InputState) (okay bool) {
	return op.match(q, input, false)
}

// doesn't match anonymous leading nouns
func (op *Nouns) LimitedMatch(q Query, input *InputState) (okay bool) {
	return op.match(q, input, true)
}

func (op *Nouns) match(q Query, input *InputState, skipAnon bool) (okay bool) {
	next := *input
	if Optional(q, &next, &op.CountedNoun) {
		okay = true
	} else if Optional(q, &next, &op.KindCalled) {
		okay = true
	} else if Optional(q, &next, &op.Kind) {
		// match the kind. if it succeeds we're done;
		// but only return true if we wanted it to succeed
		okay = !skipAnon
	} else if Optional(q, &next, &op.Name) {
		okay = true
	}
	if okay { // as long as one succeeded, try matching additional nouns too...
		Optional(q, &next, &op.AdditionalNouns)
		*input = next
	}
	return
}

func (op *AdditionalNouns) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.CommaAnd.Match(q, &next) &&
		op.Nouns.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

// return the match of this, without any additional nouns
// panics if there wasn't actually a match
func (op *Nouns) GetName(traits, kinds []Matched) (ret grok.Name) {
	if n := op.CountedNoun; n != nil {
		ret = n.GetName(traits, kinds)
	} else if n := op.KindCalled; n != nil {
		ret = n.GetName(traits, kinds)
	} else if n := op.Kind; n != nil {
		ret = n.GetName(traits, kinds)
	} else if n := op.Name; n != nil {
		ret = n.GetName(traits, kinds)
	} else {
		panic("well that was unexpected")
	}
	return
}

func (op *Nouns) GetNames(traits, kinds []Matched) (ret []grok.Name) {
	for t := *op; ; {
		n := t.GetName(traits, kinds)
		ret = append(ret, n)
		// next name:
		if next := t.AdditionalNouns; next == nil {
			break
		} else {
			t = next.Nouns
		}
	}
	return
}
