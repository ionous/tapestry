package jess

import "git.sr.ht/~ionous/tapestry/support/match"

type MatchedName interface {
	GetName(traits, kinds []Matched) resultName
	String() string
}

func (op *Name) String() string {
	return op.Matched.String()
}

func (op *Name) GetName(traits, kinds []Matched) resultName {
	return resultName{
		Article: reduceArticle(op.Article),
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

// returns index of an "important" keyword
// or the end of the string if none found
func keywordScan(ws []match.Word) (ret int) {
	ret = len(ws) // provisionally the whole thing.
Loop:
	for i, w := range ws {
		switch w.Hash() {
		case keywords.And,
			keywords.Are,
			keywords.Comma,
			keywords.Is:
			ret = i
			break Loop
		}
	}
	return
}

// similar to keyword scan; but only breaks on is/are.
func beScan(ws []match.Word) (ret int) {
	ret = len(ws) // provisionally the whole thing.
Loop:
	for i, w := range ws {
		switch w.Hash() {
		case keywords.Are,
			keywords.Is:
			ret = i
			break Loop
		}
	}
	return
}
