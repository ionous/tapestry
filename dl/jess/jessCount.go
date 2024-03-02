package jess

import "git.sr.ht/~ionous/tapestry/rt/kindsOf"

func (op *CountedKind) Match(q Query, input *InputState) (okay bool) {
	if start := *input; //
	Optional(q, &start, &op.Article) || true {
		if next := start; //
		op.MatchingNumber.Match(q, &next) &&
			op.Kind.Match(q, &next) {
			op.Matched = start.Cut(start.Len() - next.Len())
			*input, okay = next, true
		}
	}
	return
}

// for CountedNoun's private field
type CountedText = string

func (op *CountedKind) String() string {
	return op.Matched
}

func (op *CountedKind) BuildNoun(traits, kinds []string) (ret DesiredNoun, err error) {
	if kind, e := op.Kind.Validate(kindsOf.Kind); e != nil {
		err = e
	} else {
		ret = DesiredNoun{
			Count:  int(op.MatchingNumber.Number),
			Traits: traits,
			Kinds:  append(kinds, kind),
			// no name, anonymous and counted.
		}
	}
	return
}
