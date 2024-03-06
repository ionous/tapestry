package jess

import (
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
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
		plural := next.Cut(index)          // cut up to the index of "are"
		one := inflect.Singularize(plural) // fix! should use the db
		span, _ := match.MakeSpan(one)     // fix! should find kind without span
		var ks kindsOf.Kinds
		if k, w := q.FindKind(span, &ks); w == index && ks == kindsOf.Aspect {
			// fix: clean this up some.
			op.Aspect.ActualKind = ActualKind{k, ks}
			op.Aspect.Matched = span.String()
			//
			next := next.Skip(w)         // skip the kind
			op.Are.Matched = next.Cut(1) // cut the word are
			next = next.Skip(1)          // move past are
			if op.PlainNames.Match(AddContext(q, PlainNameMatching), &next) {
				*input, okay = next, true
			}
		}
	}
	return
}

func (op *AspectsAreTraits) Generate(rar Registrar) (err error) {
	if aspect, e := op.Aspect.Validate(kindsOf.Aspect); e != nil {
		err = e
	} else {
		var names []string
		for it := op.PlainNames.GetNames(); it.HasNext(); {
			n := it.GetNext()
			names = append(names, n.Name.GetNormalizedName())
		}
		err = rar.AddTraits(aspect, names)
	}
	return
}
