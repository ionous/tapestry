package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

func (op *Noun) BuildNouns(_ JessContext, w weaver.Weaves, _ rt.Runtime, props NounProperties) (ret []DesiredNoun, err error) {
	n := op.actualNoun.Name
	if e := writeKinds(w, n, props.Kinds); e != nil {
		err = e
	} else {
		var k string
		if len(props.Kinds) > 0 {
			k = props.Kinds[0]
		}
		ret = []DesiredNoun{{Noun: n, Traits: props.Traits, CreatedKind: k}}
	}
	return
}

func (op *Noun) Match(q JessContext, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.matchNoun(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *Noun) matchNoun(q JessContext, input *InputState) (okay bool) {
	if cnt := nameScan(input.Words()); cnt > 0 {
		var kind string
		sub := input.Cut(cnt)
		if m, width := q.FindNoun(sub, &kind); width > 0 {
			op.actualNoun = ActualNoun{Name: m, Kind: kind}
			op.Matched = input.Cut(width)
			*input, okay = input.Skip(width), true
		}
	}
	return
}

// the noun that matched ( as opposed to the name that matched )
type ActualNoun struct {
	Name string
	Kind string
}

func (an *ActualNoun) GetActualNoun() ActualNoun {
	return *an
}

func (an ActualNoun) IsValid() bool {
	return len(an.Name) > 0
}
