package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/kindsOf"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"git.sr.ht/~ionous/tapestry/support/match"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// runs in the AncestryPhase; but requires that the kind is known already.
// so "The colors are a kind of aspect. The colors are black and blue." is fine;
// but reversing those two sentences will fail.
func (op *AspectsAreTraits) Phase() weaver.Phase {
	// needs to be before PropertyPhase so properties can find the aspect w/o spinning.
	return weaver.AncestryPhase
}

// the colors are....
// ( see also KindsOf )
func (op *AspectsAreTraits) MatchLine(q Query, line InputState) (ret InputState, okay bool) {
	next := line //
	Optional(q, &next, &op.Aspect.Article)
	// aspects are stored as *singular*
	// ideally, fix. but the use of those aspects in kinds as fields is also singular
	// and matching of the field to the field's type is used as a filter to detect aspects
	if index := scanUntil(next.words, keywords.Are); index > 0 {
		// cut up to the index of "are"
		org := next.Cut(index)
		plural, width := match.Normalize(org)
		if width == index {
			one := inflect.Singularize(plural)                  // fix! should use the db
			if span, e := match.TokenizeString(one); e == nil { // fix! should find kind without span
				var ks kindsOf.Kinds
				if k, w := q.FindKind(span, &ks); w == index && ks == kindsOf.Aspect {
					// fix: clean this up some.
					op.Aspect.actualKind = ActualKind{k, ks}
					op.Aspect.Matched = org
					//
					next := next.Skip(w)         // skip the kind
					op.Are.Matched = next.Cut(1) // cut the word are
					next = next.Skip(1)          // move past are
					if op.PlainNames.Match(AddContext(q, PlainNameMatching), &next) {
						ret, okay = next, true
					}
				}
			}
		}
	}
	return
}

func (op *AspectsAreTraits) Weave(w weaver.Weaves, _ rt.Runtime) (err error) {
	if aspect, e := op.Aspect.Validate(kindsOf.Aspect); e != nil {
		err = e
	} else {
		var names []string
		for it := &op.PlainNames; it != nil; it = it.Next() {
			if name, e := match.NormalizeAll(it.Name.Matched); e != nil {
				err = e
				break
			} else {
				names = append(names, name)
			}
		}
		if err == nil {
			err = w.AddAspectTraits(aspect, names)
		}
	}
	return
}
