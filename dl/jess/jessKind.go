package jess

import (
	"fmt"
	"slices"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// returns the real ( generally plural ) name of the kind
func (op *Kind) Validate(ks ...kindsOf.Kinds) (ret string, err error) {
	if k := op.ActualKind.base; len(ks) > 0 && !slices.Contains(ks, k) {
		err = fmt.Errorf("matched an unexpected kind %q", k)
	} else {
		ret = op.ActualKind.name
	}
	return
}

func (op *Kind) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.matchKind(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *Kind) matchKind(q Query, input *InputState) (okay bool) {
	var k kindsOf.Kinds
	if m, width := q.FindKind(input.Words(), &k); width > 0 && filterKind(q, k) {
		op.ActualKind = ActualKind{m, k}
		op.Matched, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}

func filterKind(q Query, k kindsOf.Kinds) (okay bool) {
	if matchKindsOfAspects(q) {
		okay = k == kindsOf.Aspect
	} else if matchKindsOfKinds(q) {
		okay = k == kindsOf.Kind
	} else {
		okay = true
	}
	return
}

// anonymous kinds: "the supporter"
func (op *Kind) BuildNouns(q Query, w weaver.Weaves, run rt.Runtime, props NounProperties) (ret []DesiredNoun, err error) {
	if plural, e := op.Validate(kindsOf.Kind); e != nil {
		err = e
	} else {
		singular := run.SingularOf(plural)
		if n, e := buildAnon(w, plural, singular, props); e != nil {
			err = e
		} else {
			ret = []DesiredNoun{n}
		}
	}
	return
}

type ActualKind struct {
	name string // as opposed to just what matched
	base kindsOf.Kinds
}
