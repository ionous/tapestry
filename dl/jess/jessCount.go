package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

func (op *CountedKind) Match(q JessContext, input *InputState) (okay bool) {
	if start := *input; //
	Optional(q, &start, &op.Article) || true {
		if next := start; //
		op.MatchingNum.Match(q, &next) &&
			op.Kind.Match(q, &next) {
			op.Matched = start.Cut(start.Len() - next.Len())
			*input, okay = next, true
		}
	}
	return
}

// generates n initial instances (and their aliases, cause why not.)
// delays the desired traits and additional kinds
// ( tbd if that makes sense or not )
func (op *CountedKind) BuildNouns(q JessContext, w weaver.Weaves, run rt.Runtime, props NounProperties) (ret []DesiredNoun, err error) {
	if plural, e := op.Kind.Validate(kindsOf.Kind); e != nil {
		err = e
	} else {
		if cnt := int(op.MatchingNum.Value); cnt > 0 {
			singular := run.SingularOf(plural)
			ret = make([]DesiredNoun, cnt)
			for i := 0; i < cnt; i++ {
				if n, e := buildAnon(w, plural, singular, props); e != nil {
					err = e
					break
				} else {
					ret[i] = n
				}
			}
		}
	}
	return
}
