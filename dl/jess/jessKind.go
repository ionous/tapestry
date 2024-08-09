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
	if k := op.actualKind.BaseKind; len(ks) > 0 && !slices.Contains(ks, k) {
		err = fmt.Errorf("matched an unexpected kind %q", k)
	} else {
		ret = op.actualKind.Name
	}
	return
}

func (op *Kind) Match(q JessContext, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Article) || true) &&
		op.matchKind(q, &next) {
		*input, okay = next, true
	}
	return
}

// for use in properties
func (op *Kind) Assignment() rt.Assignment {
	return text(op.actualKind.Name, "") // tbd: should these be typed? ex. as "kinds" or something?
}

func (op *Kind) matchKind(q JessContext, input *InputState) (okay bool) {
	var k kindsOf.Kinds
	if m, width := q.FindKind(input.Words(), &k); width > 0 && filterKind(q, k) {
		op.actualKind = ActualKind{m, k}
		op.Matched = input.Cut(width)
		*input, okay = input.Skip(width), true
	}
	return
}

// if no specific filter is set, then all kinds can match;
// otherwise one of the specific kinds must match.
func filterKind(q JessContext, k kindsOf.Kinds) (okay bool) {
	aspects, kinds := matchKindsOfAspects(q), matchKindsOfKinds(q)
	if !aspects && !kinds {
		okay = true
	} else {
		okay = (aspects && k == kindsOf.Aspect) ||
			(kinds && k == kindsOf.Kind)
	}
	return
}

// anonymous kinds: "the supporter"
func (op *Kind) BuildNouns(q JessContext, w weaver.Weaves, run rt.Runtime, props NounProperties) (ret []DesiredNoun, err error) {
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
	Name     string // as opposed to just what matched
	BaseKind kindsOf.Kinds
}

// search ancestry for an existing kind
func TryKind(q JessContext, in InputState,
	accept func(Kind, InputState),
	reject func(error)) {
	q.Try(After(weaver.AncestryPhase), func(weaver.Weaves, rt.Runtime) {
		var kind Kind
		if !kind.Match(q, &in) {
			reject(FailedMatch{"couldn't find matching kind", in})
		} else {
			accept(kind, in)
		}
		return
	}, reject)
}
