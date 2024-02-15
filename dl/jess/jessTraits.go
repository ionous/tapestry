package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

func (op *Article) Match(q Query, input *InputState) (okay bool) {
	if m, width := q.FindArticle(*input); width > 0 {
		op.Matched, *input, okay = m, input.Skip(width), true
	}
	return
}

func ReduceArticle(op *Article) (ret grok.Article) {
	if op != nil {
		ret = grok.Article{
			Matched: op.Matched,
		}
	}
	return
}

func (op *Trait) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.matchTrait(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *Trait) matchTrait(q Query, input *InputState) (okay bool) {
	if m, width := q.FindTrait(*input); width > 0 {
		op.Matched, *input, okay = m, input.Skip(width), true
	}
	return
}

// its interesting that we dont have to store anything else
// all the trait info is in this... even additional traits.
func (op *Traits) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Trait.Match(q, &next) {
		Optional(q, &next, &op.AdditionalTraits)
		*input, okay = next, true
	}
	return
}

func (op *Traits) GetTraits() []Matched {
	var out []Matched
	for t := *op; ; {
		out = append(out, t.Trait.Matched)
		if next := t.AdditionalTraits; next == nil {
			break
		} else {
			t = next.Traits
		}
	}
	return out
}

func (op *AdditionalTraits) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.CommaAnd) || true) &&
		op.Traits.Match(q, &next) {
		*input, okay = next, true
	}
	return
}
