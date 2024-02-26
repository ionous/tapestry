package jess

import "git.sr.ht/~ionous/tapestry/rt/kindsOf"

func (op *CountedName) Match(q Query, input *InputState) (okay bool) {
	if start := *input; //
	Optional(q, &start, &op.Article) || true {
		if next := start; //
		op.MatchingNumber.Match(q, &next) &&
			op.Kind.Match(q, &next) {
			op.Matched = start.Cut(len(start) - len(next))
			*input, okay = next, true
		}
	}
	return
}

func (op *CountedName) String() string {
	return op.Matched.String()
}

func (op *CountedName) GetName(traits, kinds []string) (ret resultName, err error) {
	if kind, e := op.Kind.Validate(kindsOf.Kind); e != nil {
		err = e
	} else {
		ret = resultName{
			Article: articleResult{
				Count: int(op.MatchingNumber.Number),
			},
			Traits: traits,
			Kinds:  append(kinds, kind),
			// no name, anonymous and counted.
		}
	}
	return
}
