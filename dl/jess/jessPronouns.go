package jess

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

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
		// tbd: does it make sense to wait here ( and not, say, in the build? )
		q.TryTopic(func(an ActualNoun) {
			accept(Pronoun{
				Matched: in.words,
				topic:   Topic{noun: an},
			}, in.Skip(w))
		}, reject)
	}
}

// property pronouns can be implied: "The description is".
func TryImpliedPronoun(q JessContext,
	accept func(Pronoun),
	reject func(error),
) {
	q.TryTopic(func(an ActualNoun) {
		accept(Pronoun{
			topic: Topic{noun: an},
		})
	}, reject)
}

// PropertyNoun interface, only valid after a match.
// ( the kind of the sentence topic is needed to match properties )
func (op *Pronoun) GetKind() string {
	return op.topic.noun.Kind
}

// PropertyNoun interface, only valid after a match.
func (op *Pronoun) BuildPropertyNoun(ctx BuildContext) (ActualNoun, error) {
	return op.topic.Resolve()
}

// for non-promise matching code pathss
// when the query context has "MatchPronouns",
// try to match the *use* of a pronoun ( ex. "it" ).
func (op *Pronoun) Match(q JessContext, input *InputState) (okay bool) {
	// fix: ideally replace with "promise" so that it can return error okay
	if matchPronouns(q) {
		if width := input.MatchWord(keywords.It); width > 0 {
			// request the topic of the former sentence
			q.TryTopic(op.topic.acceptedTopic, op.topic.rejectedTopic)
			op.Matched = input.Cut(width)
			*input = input.Skip(width)
			okay = true
		}
	}
	return
}

// names are often potential nouns;
// this helper generates them as such.
func (op *Pronoun) BuildNouns(q JessContext, w weaver.Weaves, run rt.Runtime, props NounProperties) (ret []DesiredNoun, err error) {
	if an, e := op.topic.Resolve(); e != nil {
		err = e
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
