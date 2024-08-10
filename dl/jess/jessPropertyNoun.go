package jess

import (
	"errors"
)

type PropertyNoun interface {
	GetKind() string // to pick the names of properties from phrases.
	BuildPropertyNoun(BuildContext) (string, error)
}

func TryPropertyNoun(q JessContext, in InputState,
	accept func(PropertyNoun, InputState),
	reject func(error)) {
	// because property matching needs the full kind
	// this has to wait until the previous phrases has been built.
	TryPronoun(q, in, func(pn Pronoun, next InputState) {
		prop := PropertyPronoun{Pronoun: pn}
		accept(&prop, next)
	}, func(e error) {
		// differentiate between when matching a pronoun succeeded,
		// and when generating using that pronoun failed.
		var matchError FailedMatch
		if !errors.As(e, &matchError) {
			reject(e)
		} else {
			// matches a name, kind, and traits. ( "the animal called the cat" )
			TryInlineNoun(q, in, func(n InlineNoun, next InputState) {
				q.CurrentPhrase().SetTopic(&n)
				accept(&n, next)
			}, func(e error) {
				var matchError FailedMatch
				if !errors.As(e, &matchError) {
					reject(e)
				} else {
					// see if there is already such a noun.
					TryExistingNoun(q, in, func(noun ExistingNoun, next InputState) {
						q.CurrentPhrase().SetTopic(&noun)
						accept(&noun, next)
					}, func(error) {
						// if all the earlier matches failed, generate an implicit noun.
						TryImplicitNoun(q, in, func(noun ImplicitNoun, next InputState) {
							q.CurrentPhrase().SetTopic(&noun)
							accept(&noun, next)
						}, reject)
					})
				}
			})
		}
	})
}
