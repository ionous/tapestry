package jess

import (
	"log"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

// allows partial matches; test that there's no input left to verify a complete match.
func (op *MatchingPhrases) Match(q Query, input *InputState) (ret Generator, okay bool) {
	type matcher interface {
		Generator
		Interpreter
		typeinfo.Instance
	}
	var best InputState
	var bestMatch matcher
	// fix? could change to reflect ( or expand type info ) to walk generically
	for _, m := range []matcher{
		// understand "..." as .....
		&op.Understand,
		// names are "a kind of"/"kinds of" [traits] kind:any.
		&op.KindsOf,
		// kind:objects are "usually" traits.
		&op.KindsAreTraits,
		// kinds:records|objects "have" a ["list of"] number|text|records|objects|aspects ["called a" ...]
		&op.KindsHaveProperties,
		// kinds:objects ("can be"|"are either") new_trait [or new_trait...]
		&op.KindsAreEither,
		// kind:aspects are names
		&op.AspectsAreTraits,
		&op.VerbNamesAreNames,
		&op.NamesVerbNames,
		&op.NamesAreLikeVerbs,
		&op.PropertyNounValue,
		&op.NounPropertyValue,
		&op.MapLocations,
		&op.MapDirections,
		&op.MapConnections,
	} {
		if next := *input; //
		m.Match(q, &next) /* && len(next) == 0 */ {
			if !okay || next.Len() < best.Len() {
				best = next
				bestMatch, okay = m, true
				if best.Len() == 0 {
					break
				}
			}
		}
	}
	if okay {
		ret, *input = bestMatch, best
		if useLogging(q) {
			log.Println("matched", bestMatch.TypeInfo().TypeName())
		}
	}
	return
}

func (op *VerbNamesAreNames) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Verb.Match(q, &next) &&
		op.Names.Match(q, &next) &&
		op.Are.Match(q, &next) &&
		op.OtherNames.Match(q, &next) {
		// q.note("matched VerbNamesAreNames")
		*input, okay = next, true
	}
	return
}

func (op *VerbNamesAreNames) Generate(rar Registrar) (err error) {
	if res, e := op.compile(); e != nil {
		err = e
	} else {
		err = applyResults(rar, res)
	}
	return
}

func (op *VerbNamesAreNames) compile() (ret localResults, err error) {
	if lhs, e := op.Names.BuildNouns(nil, nil); e != nil {
		err = e
	} else if rhs, e := op.OtherNames.BuildNouns(nil, nil); e != nil {
		err = e
	} else {
		ret = makeResult(op.Verb.Macro, !op.Verb.Macro.Reversed, lhs, rhs)
	}
	return
}

func (op *NamesVerbNames) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Names.Match(q, &next) &&
		!op.Names.HasAnonymousKind() &&
		op.Are.Match(q, &next) &&
		op.Verb.Match(q, &next) &&
		op.OtherNames.Match(q, &next) {
		// q.note("matched NamesVerbNames")
		*input, okay = next, true
	}
	return
}

func (op *NamesVerbNames) Generate(rar Registrar) (err error) {
	if res, e := op.compile(); e != nil {
		err = e
	} else {
		err = applyResults(rar, res)
	}
	return
}

func (op *NamesVerbNames) compile() (ret localResults, err error) {
	if lhs, e := op.Names.BuildNouns(nil, nil); e != nil {
		err = e
	} else if rhs, e := op.OtherNames.BuildNouns(nil, nil); e != nil {
		err = e
	} else {
		ret = makeResult(op.Verb.Macro, op.Verb.Macro.Reversed, lhs, rhs)
	}
	return
}

func (op *NamesAreLikeVerbs) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Names.Match(q, &next) &&
		!op.Names.HasAnonymousKind() &&
		op.Are.Match(q, &next) &&
		op.Adjectives.Match(q, &next) {
		Optional(q, &next, &op.VerbPhrase)
		*input, okay = next, true
	}
	return
}

func (op *NamesAreLikeVerbs) BuildOtherNouns() (ret []DesiredNoun, err error) {
	if v := op.VerbPhrase; v != nil {
		ret, err = op.VerbPhrase.BuildNouns(nil, nil)
	}
	return
}

func (op *NamesAreLikeVerbs) GetMacro() (ret Macro) {
	if v := op.VerbPhrase; v != nil {
		ret = v.Verb.Macro
	}
	return
}

func (op *NamesAreLikeVerbs) IsReversed() (okay bool) {
	if v := op.VerbPhrase; v != nil {
		okay = v.Verb.Macro.Reversed
	}
	return
}

func (op *NamesAreLikeVerbs) Generate(rar Registrar) (err error) {
	if res, e := op.compile(); e != nil {
		err = e
	} else {
		err = applyResults(rar, res)
	}
	return
}

func (op *NamesAreLikeVerbs) compile() (ret localResults, err error) {
	if rhs, e := op.BuildOtherNouns(); e != nil {
		err = e
	} else if ts, ks, e := op.Adjectives.Reduce(); e != nil {
		err = e
	} else if lhs, e := op.Names.BuildNouns(ts, ks); e != nil {
		err = e
	} else {
		ret = makeResult(op.GetMacro(), op.IsReversed(), lhs, rhs)
	}
	return
}
