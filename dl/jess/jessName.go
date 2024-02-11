package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

func (op *Name) Match(q Query, input *InputState) (okay bool) {
	if width := keywordScan(input.Words()); width > 0 {
		op.Matched, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}

func (op *Names) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	Optionally(q, &next, &op.Article) &&
		op.Name.Match(q, &next) &&
		Optionally(q, &next, &op.AdditionalNames) {
		*input, okay = next, true
	}
	return
}

func (op *AdditionalNames) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.CommaAnd.Match(q, &next) &&
		op.Names.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *Names) Reduce(traits, kinds []Matched) []grok.Name {
	var out []grok.Name
	for t := *op; ; {
		out = append(out, grok.Name{
			Article: ReduceArticle(t.Article),
			Span:    t.Name.Matched.(Span),
			Traits:  traits,
			Kinds:   kinds,
		})
		// next name:
		if next := t.AdditionalNames; next == nil {
			break
		} else {
			t = next.Names
		}
	}
	return out
}

// returns index of an "important" keyword
// or the end of the string if none found
func keywordScan(ws []grok.Word) (ret int) {
	ret = len(ws) // provisionally the whole thing.
	for i, w := range ws {
		switch w.Hash() {
		case grok.Keyword.And,
			grok.Keyword.Are,
			grok.Keyword.Comma,
			grok.Keyword.Is:
			ret = i
			break
		}
	}
	return
}
