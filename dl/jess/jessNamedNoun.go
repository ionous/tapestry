package jess

func TryNamedNoun(q JessContext, in InputState,
	accept func(NamedNoun, ActualNoun, InputState),
	reject func(error)) {
	TryPronoun(q, in, func(pro Pronoun, noun ActualNoun, next InputState) {
		// accept fires in the value phase;
		// reject can fire asap if it doesn't look like a pronoun.
		accept(NamedNoun{Pronoun: &pro}, noun, next)
	}, func(error) {
		// matches a name, kind, and traits
		TryInlineNoun(q, in, func(n InlineNoun, next InputState) {
			// create a noun with the matched data:
			GenerateNoun(q, n.Name, n.GetKind(), n.GetTraits(), func(noun ActualNoun) {
				// matched!
				accept(NamedNoun{InlineNoun: &n}, noun, next)
			}, reject)
		}, func(error) {
			var n Name
			if next := in; !n.Match(q, &next) {
				reject(FailedMatch{"TryNamedNoun expected a name", next})
			} else {
				// matched! now (in fallbacks) make sure the noun exists
				GenerateImplicitNoun(q, n, func(noun ActualNoun) {
					accept(NamedNoun{Name: &n}, noun, next)
				}, reject)
			}
		})
	})
	return
}
