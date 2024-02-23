package jess

import (
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/support/match"
)

// the colors are....
// ( see also KindsOf )
func (op *AspectsAreTraits) Match(q Query, input *InputState) (okay bool) {
	next := *input
	Optional(q, &next, &op.Aspect.Article)
	// aspects are stored as *singular*
	// ideally, fix. but the use of those aspects in kinds as fields is also singular
	// and matching of the field to the field's type is used as a filter to detect aspects
	if index := scanUntil(next.Words(), keywords.Are); index > 0 {
		plural := next.Cut(index)                   // cut up to the index of "are"
		one := inflect.Singularize(plural.String()) // fix! should use the db
		span, _ := match.MakeSpan(one)              // fix! should find kind without span
		if k, w := q.FindKind(span); w == index {
			op.Aspect.Matched = matchedString{k, w}
			//
			next := next.Skip(w)         // skip the kind
			op.Are.Matched = next.Cut(1) // cut the word are
			next = next.Skip(1)          // move past are
			if op.Names.Match(AddContext(q, PlainNameMatching), &next) {
				*input, okay = next, true
			}
		}
	}
	return
}

func (op *AspectsAreTraits) Generate(rar Registrar) (err error) {
	var names []string
	for it := op.Names.Iterate(); it.HasNext(); {
		n := it.GetNext()
		names = append(names, inflect.Normalize(n.String()))
	}
	return rar.AddTraits(op.Aspect.String(), names)
}
