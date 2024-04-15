package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

func (op *Noun) BuildNouns(_ Query, w weaver.Weaves, _ rt.Runtime, props NounProperties) (ret []DesiredNoun, err error) {
	n := op.ActualNoun.Name
	if e := writeKinds(w, n, props.Kinds); e != nil {
		err = e
	} else {
		ret = []DesiredNoun{{Noun: n, Traits: props.Traits}}
	}
	return
}

func (op *Noun) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.matchNoun(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *Noun) matchNoun(q Query, input *InputState) (okay bool) {
	if cnt := keywordScan(input.Words()); cnt > 0 {
		sub := input.Cut(cnt)
		// fix? it'd be nice if the mapping of "you" to "self" was handled by script;
		// or even not necessary at all.
		if width := 1; len(sub) == width && sub[0].Hash() == keywords.You {
			op.ActualNoun = ActualNoun{Name: PlayerSelf, Kind: Actors}
			op.Matched = input.Cut(width)
			*input, okay = input.Skip(width), true
		} else {
			// match the subsection normally:
			var kind string
			if m, width := q.FindNoun(sub, &kind); width > 0 {
				op.ActualNoun = ActualNoun{Name: m, Kind: kind}
				op.Matched = input.Cut(width)
				*input, okay = input.Skip(width), true
			}
		}
	}
	return
}

// the noun that matched ( as opposed to the name that matched )
type ActualNoun struct {
	Name string
	Kind string
}
