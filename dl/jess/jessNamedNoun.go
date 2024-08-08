package jess

import (
	"errors"
	"fmt"
)

func TryNamedNoun(q JessContext, in InputState,
	accept func(NamedNoun, ActualNoun, InputState),
	reject func(error)) {
	TryPronoun(q, in, func(pro Pronoun, noun ActualNoun, next InputState) {
		accept(NamedNoun{Pronoun: &pro}, noun, next)
	}, func(e error) {
		// differentiate between when matching a pronoun succeeded,
		// and when generating using that pronoun failed.
		var matchError FailedMatch
		if !errors.As(e, &matchError) {
			reject(e)
		} else {
			// matches a name, kind, and traits. ( "the animal called the cat" )
			TryInlineNoun(q, in, func(n InlineNoun, next InputState) {
				futureNoun := new(ActualNoun)
				q.SetTopic(futureNoun) // should this be after the generation?
				// create a noun with the matched data:
				GenerateNoun(q, n.Name, n.GetKind(), n.GetTraits(), func(an ActualNoun) {
					*futureNoun = an
					accept(NamedNoun{InlineNoun: &n}, an, next)
				}, reject)
			}, func(e error) {
				var matchError FailedMatch
				if !errors.As(e, &matchError) {
					reject(e)
				} else {
					// if all the earlier matches failed, try as an implicit noun.
					TryImplicitNoun(q, in, func(name Name, an ActualNoun, next InputState) {
						accept(NamedNoun{Name: &name}, an, next)
					}, reject)
				}
			})
		}
	})
}

// a phrase implies a noun exists, but no particular kind has been assigned.
func TryImplicitNoun(q JessContext, in InputState,
	accept func(Name, ActualNoun, InputState),
	reject func(error),
) {
	// tricksy: don't generate an implicit noun if it would conflict with a kind
	// ( its probably actually some other phrase )
	TryKind(q, in, func(kind Kind, rest InputState) {
		// tbd: allow this if there is an existing noun?
		// ( ordering to get the exact kind might be hard re: fallbacks )
		reject(fmt.Errorf("the phrase implies a noun, but there's already a kind of that name %s", kind.actualKind.Name))
	}, func(error) {
		var n Name
		if next := in; !n.Match(q, &next) {
			reject(FailedMatch{"expected some sort of name", next})
		} else {
			futureNoun := new(ActualNoun)
			q.SetTopic(futureNoun) // should this be after the generation?
			GenerateImplicitNoun(q, n, func(an ActualNoun) {
				*futureNoun = an
				accept(n, an, next)
			}, reject)
		}
	})
}
