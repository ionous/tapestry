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

func (op *NamedTrait) Match(q Query, input *InputState) (okay bool) {
	if m, width := q.FindTrait(*input); width > 0 {
		op.Matched, *input, okay = m, input.Skip(width), true
	}
	return
}

// its interesting that we dont have to store anything else
// all the trait info is in this... even additional traits.
func (op *Traits) Match(q Query, input *InputState) (okay bool) {
	if next := *input; Optionally(q, &next, &op.Article) &&
		op.NamedTrait.Match(q, &next) &&
		Optionally(q, &next, &op.AdditionalTraits) {
		*input, okay = next, true
	}
	return
}

func (op *Traits) GetTraits() []Matched {
	var out []Matched
	for t := *op; ; {
		out = append(out, t.NamedTrait.Matched)
		if next := t.AdditionalTraits; next == nil {
			break
		} else {
			t = next.Traits
		}
	}
	return out
}

func (op *AdditionalTraits) Match(q Query, input *InputState) (okay bool) {
	if next := *input; Optionally(q, &next, &op.CommaAnd) &&
		op.Traits.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

// a generic method to read an optional element.
// always returns true.
// optional elements are implemented by pointer types
// the pointers are required to implement Interpreter
func Optionally[M any,
	IM interface {
		// the syntax for this feels very strange
		// the method takes an 'out' pointer ( so go can determine the type by inference )
		// the "interface" is a reused keyword, signifying a type constraint
		// indicating we want pointers to M to also be (implement) Interpreter.
		// *phew*
		*M
		Interpreter
	}](q Query, input *InputState, out **M) bool {
	var v M
	if next := *input; IM(&v).Match(q, &next) {
		*out = &v
		*input = next
	}
	return true
}
