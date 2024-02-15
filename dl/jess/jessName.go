package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

func (op *Name) GetName(traits, kinds []Matched) grok.Name {
	return grok.Name{
		Article: ReduceArticle(op.Article),
		Span:    op.Matched.(Span),
		Traits:  traits,
		Kinds:   kinds,
	}
}

func (op *Name) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.matchName(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *Name) matchName(q Query, input *InputState) (okay bool) {
	if width := keywordScan(input.Words()); width > 0 {
		op.Matched, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}

func (op *Names) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Name.Match(q, &next) &&
		(Optional(q, &next, &op.AdditionalNames) || true) {
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

func (op *Names) Reduce(traits, kinds []Matched) (ret []grok.Name) {
	for t := *op; ; {
		n := t.Name.GetName(traits, kinds)
		ret = append(ret, n)
		// next name:
		if next := t.AdditionalNames; next == nil {
			break
		} else {
			t = next.Names
		}
	}
	return
}

// returns index of an "important" keyword
// or the end of the string if none found
func keywordScan(ws []grok.Word) (ret int) {
	ret = len(ws) // provisionally the whole thing.
Loop:
	for i, w := range ws {
		switch w.Hash() {
		case grok.Keyword.And,
			grok.Keyword.Are,
			grok.Keyword.Comma,
			grok.Keyword.Is:
			ret = i
			break Loop
		}
	}
	return
}

// similar to keyword scan; but only breaks on is/are.
func beScan(ws []grok.Word) (ret int) {
	ret = len(ws) // provisionally the whole thing.
Loop:
	for i, w := range ws {
		switch w.Hash() {
		case grok.Keyword.Are,
			grok.Keyword.Is:
			ret = i
			break Loop
		}
	}
	return
}
