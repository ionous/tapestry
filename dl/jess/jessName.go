package jess

import "git.sr.ht/~ionous/tapestry/support/grok"

func (op *Name) GetName(art grok.Article, traits, kinds []Matched) grok.Name {
	return grok.Name{
		Article: art,
		Span:    op.Matched.(Span),
		Traits:  traits,
		Kinds:   kinds,
	}
}

func (op *Name) Match(q Query, input *InputState) (okay bool) {
	if width := keywordScan(input.Words()); width > 0 {
		op.Matched, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}

func (op *Names) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	Optionally(q, &next, &op.Article) &&
		TryNameCalled(q, &next, &op.NameCalled) &&
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
		// fix? move t.Article into the sub-phrases?
		// name, and named_kind, named_trait would become:
		// TheName, TheKind, TheTrait maybe (article would move out of kinds, names, traits )
		a := ReduceArticle(op.Article)
		n := t.NameCalled.GetName(a, traits, kinds)
		out = append(out, n)
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
