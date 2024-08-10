package jess

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

func (op *PropertyPronoun) BuildPropertyNoun(ctx BuildContext) (string, error) {
	return op.Pronoun.actualNoun.Name, nil
}

// PropertyNoun interface, only valid after a match.
func (op *PropertyPronoun) GetKind() string {
	return op.Pronoun.actualNoun.Kind
}

// PropertyNoun interface, only valid after build
func (op *PropertyPronoun) GetActualNoun() ActualNoun {
	return op.Pronoun.actualNoun
}

// when the query context has "MatchPronouns",
// try to match the *use* of a pronoun ( ex. "it" ).
func (op *Pronoun) Match(q JessContext, input *InputState) (okay bool) {
	// fix: ideally replace with "promise" so that it can return error okay
	if matchPronouns(q) {
		if width := input.MatchWord(keywords.It); width > 0 {
			if q.CurrentPhrase().UsePronoun() {
				op.Matched = input.Cut(width)
				*input = input.Skip(width)
				okay = true
			}
		}
	}
	return
}

// names are often potential nouns;
// this helper generates them as such.
func (op *Pronoun) BuildNouns(q JessContext, w weaver.Weaves, run rt.Runtime, props NounProperties) (ret []DesiredNoun, err error) {
	if an := q.GetTopic(); !an.IsValid() {
		err = fmt.Errorf("couldn't find topic of pronoun")
	} else {
		// duplicates Noun.BuildNouns:
		if e := writeKinds(w, an.Name, props.Kinds); e != nil {
			err = e
		} else {
			k := an.Kind
			if len(props.Kinds) > 0 {
				k = props.Kinds[0]
			}
			ret = []DesiredNoun{{Noun: an.Name, Traits: props.Traits, CreatedKind: k}}
		}
	}
	return
}

// search for a pronoun and the noun that it represents.
// accept fires in the value phase;
// reject can fire asap if it doesn't look like a pronoun.
func TryPronoun(q JessContext, in InputState,
	accept func(Pronoun, InputState),
	reject func(error),
) {
	// the name scan differentiates "the it girl" from "she is ..."
	// right now pronouns are all one word long
	if w := nameScan(in.words); w != 1 {
		reject(FailedMatch{"couldn't find a pronoun", in})
	} else if in.words[0].Hash() != keywords.It { // ( and always "it")
		reject(FailedMatch{"word isn't a known pronoun", in.Slice(w)})
	} else {
		RequestPronoun(q, func(an ActualNoun) {
			accept(Pronoun{
				Matched:    in.words,
				actualNoun: an,
			}, in.Skip(w))
		}, reject)
	}
}

// property pronouns can be implied: "The description is".
func TryImpliedPronoun(q JessContext,
	accept func(PropertyPronoun),
	reject func(error),
) {
	RequestPronoun(q, func(an ActualNoun) {
		accept(PropertyPronoun{Pronoun: Pronoun{actualNoun: an}})
	}, reject)
}

// determine the topic of the sentence based on an earlier definition.
func RequestPronoun(q JessContext,
	accept func(ActualNoun),
	reject func(error),
) {
	// fix: a direct dependency on previous phrase -- whatever its phase -- would be better.
	// (ex. add a notifier or a something to a topical sentence)
	q.Try(After(weaver.FallbackPhase), func(weaver.Weaves, rt.Runtime) {
		if !q.CurrentPhrase().UsePronoun() {
			reject(errors.New("sentence describes a particular noun"))
		} else if an := q.GetTopic(); !an.IsValid() {
			e := fmt.Errorf("couldn't find topic of pronoun")
			reject(e)
		} else {
			accept(an)
		}
	}, reject)

}
