package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

func (op *Name) GetNormalizedName() string {
	return inflect.Normalize(op.Matched.String())
}

// names are often potential nouns;
// this helper generates them as such.
func (op *Name) BuildNouns(q Query, w weaver.Weaves, run rt.Runtime, props NounProperties) (ret []DesiredNoun, err error) {
	if n, e := op.buildNoun(q, w, props); e != nil {
		err = e
	} else {
		ret = []DesiredNoun{n}
	}
	return
}

func (op *Name) buildNoun(q Query, w weaver.Weaves, props NounProperties) (ret DesiredNoun, err error) {
	nounSpan := op.Matched.(match.Span)
	if noun, created, e := ensureNoun(q, w, nounSpan, &props); e != nil {
		err = e
	} else if e := writeKinds(w, noun, props.Kinds); e != nil {
		err = e
	} else {
		n := DesiredNoun{Noun: noun, Traits: props.Traits}
		if created {
			n.appendArticle(op.Article)
		}
		ret = n
	}
	return
}

func (op *Name) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.matchName(&next) {
		*input, okay = next, true
	}
	return
}

func (op *Name) matchName(input *InputState) (okay bool) {
	if width := keywordScan(input.Words()); width > 0 {
		op.Matched = input.CutSpan(width)
		*input, okay = input.Skip(width), true
	}
	return
}

// returns index of an "important" keyword
// or the end of the string if none found.
// inform also has troubles with names like "the has been."
func keywordScan(ws []match.Word) (ret int) {
	ret = len(ws) // provisionally the whole thing.
Loop:
	for i, w := range ws {
		switch w.Hash() {
		case //
			keywords.And,
			keywords.Are,
			keywords.Comma,
			keywords.Has,
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

// returns the index of the matching word in the span
// -1 if not found
func scanUntil(span []match.Word, hashes ...uint64) (ret int) {
	ret = -1
Loop:
	for i, w := range span {
		m := w.Hash()
		for _, h := range hashes {
			if h == m {
				ret = i
				break Loop
			}
		}
	}
	return
}
