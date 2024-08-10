package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// match an existing noun
func TryExistingNoun(q JessContext, in InputState,
	accept func(ExistingNoun, InputState),
	reject func(error),
) {
	q.Try(After(weaver.NounPhase), func(weaver.Weaves, rt.Runtime) {
		var noun ExistingNoun
		if next := in; !noun.Match(q, &next) {
			reject(FailedMatch{"no such noun", in})
		} else {
			accept(noun, next)
		}
	}, reject)
}

// valid after match
func (op *ExistingNoun) GetKind() string {
	return op.actualNoun.Kind
}

// valid after match ( because it already exists )
func (op *ExistingNoun) BuildPropertyNoun(ctx BuildContext) (string, error) {
	return op.actualNoun.Name, nil
}

// valid after match ( because it already exists )
func (op *ExistingNoun) GetActualNoun() ActualNoun {
	return op.actualNoun
}

// so that nouns can be used as the *value* of a property
func (op *ExistingNoun) Assignment() rt.Assignment {
	return text(op.actualNoun.Name, op.actualNoun.Kind)
}

// fix: used for old phrase matching; but doesn't make a lot of sense for existing nouns.
func (op *ExistingNoun) BuildNouns(_ JessContext, w weaver.Weaves, _ rt.Runtime, props NounProperties) (ret []DesiredNoun, err error) {
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

func (op *ExistingNoun) Match(q JessContext, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.matchNoun(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *ExistingNoun) matchNoun(q JessContext, input *InputState) (okay bool) {
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
