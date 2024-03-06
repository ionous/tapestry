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

// generates n initial instances (and their aliases, cause why not.)
// delays the desired traits and additional kinds
// ( tbd if that makes sense or not )
func (op *CountedKind) BuildNouns(q Query, rar Registrar, ts, ks []string) (ret []DesiredNoun, err error) {
	if plural, e := op.Kind.Validate(kindsOf.Kind); e != nil {
		err = e
	} else {
		if cnt := int(op.MatchingNumber.Number); cnt > 0 {
			singular := rar.GetSingular(plural)
			ret = make([]DesiredNoun, cnt)
			for i := 0; i < cnt; i++ {
				if n, e := buildAnon(rar, plural, singular, ts, ks); e != nil {
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
