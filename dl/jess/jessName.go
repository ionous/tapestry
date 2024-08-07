package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// the exact name specified by the player
func (op *Name) GetNormalizedName() (string, error) {
	return match.NormalizeAll(op.Matched)
}

// names are often potential nouns;
// this helper treats them as such,
// ensuring that they exist in the database.
func (op *Name) BuildNouns(q Query, w weaver.Weaves, run rt.Runtime, props NounProperties) (ret []DesiredNoun, err error) {
	if n, e := op.BuildNoun(q, w, props); e != nil {
		err = e
	} else {
		ret = []DesiredNoun{n}
	}
	return
}

// note: this doesn't apply values ( nor articles )
// it stores them until writeNounValues ( left until after the kind of the noun is firmly established )
func (op *Name) BuildNoun(q Query, w weaver.Weaves, props NounProperties) (ret DesiredNoun, err error) {
	if noun, kind, created, e := ensureNoun(q, w, op.Matched, &props); e != nil {
		err = e
	} else if e := writeKinds(w, noun, props.Kinds); e != nil {
		err = e
	} else {
		n := DesiredNoun{Noun: noun, Traits: props.Traits, CreatedKind: kind}
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
	if width := nameScan(input.Words()); width > 0 {
		op.Matched = input.Cut(width)
		*input, okay = input.Skip(width), true
	}
	return
}

// returns index of an "important" keyword
// or the end of the string if none found.
// inform also has troubles with names like "the has been."
func nameScan(ts []match.TokenValue) (ret int) {
	if cnt := len(ts); cnt > 0 {
		if ts[0].Token == match.Quoted {
			ret = 1
		} else if i := scanUntil(ts, nameSeparators...); i < 0 {
			ret = cnt
		} else {
			ret = i
		}
	}
	return
}

var nameSeparators = []uint64{
	keywords.And,
	keywords.Are,
	keywords.Comma,
	keywords.Has,
	keywords.Have,
	keywords.Is,
}

// returns the index of the matching word in the span
// -1 if not found
func scanUntil(ts []match.TokenValue, hashes ...uint64) (ret int) {
	ret = -1
Loop:
	for i, tv := range ts {
		m := tv.Hash()
		for _, h := range hashes {
			if h == m {
				ret = i
				break Loop
			}
		}
	}
	return
}
