package jess

import "errors"

func TryNamedNoun(q JessContext, in InputState,
	accept func(NamedNoun, ActualNoun, InputState),
	reject func(error)) {
	TryPronoun(q, in, func(pro Pronoun, noun ActualNoun, next InputState) {
		// accept fires in the value phase;
		// reject can fire asap if it doesn't look like a pronoun.
		accept(NamedNoun{Pronoun: &pro}, noun, next)
	}, func(e error) {
		// differentiate between when matching a pronoun succeeded,
		// and when generating using that pronoun failed.
		var matchError FailedMatch
		if !errors.As(e, &matchError) {
			reject(e)
		} else {
			// matches a name, kind, and traits
			TryInlineNoun(q, in, func(n InlineNoun, next InputState) {
				// create a noun with the matched data:
				GenerateNoun(q, n.Name, n.GetKind(), n.GetTraits(), func(noun ActualNoun) {
					// matched!
					accept(NamedNoun{InlineNoun: &n}, noun, next)
				}, reject)
			}, func(e error) {
				var matchError FailedMatch
				if !errors.As(e, &matchError) {
					reject(e)
				} else {
					var n Name
					if next := in; !n.Match(q, &next) {
						reject(FailedMatch{"expected some sort of name", next})
					} else {
						// matched! now (in fallbacks) make sure the noun exists
						GenerateImplicitNoun(q, n, func(noun ActualNoun) {
							accept(NamedNoun{Name: &n}, noun, next)
						}, reject)
					}
				}
			})
		}
	})
	return
}
