package jess

import (
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// "the <inline kind: container> called the box..."
func TryInlineNoun(q JessContext, in InputState,
	accept func(InlineNoun, InputState),
	reject func(error)) {
	var called Called
	var name Name
	// {lhs: words} 'called' {rhs: name}
	if lhs, rhs, ok := called.Split(in); !ok {
		reject(FailedMatch{"expected a kind called something", in})
	} else if !name.Match(AddContext(q, CheckIndefiniteArticles), &rhs) {
		reject(FailedMatch{"expected a name", rhs})
	} else {
		TryInlineKind(q, lhs, func(kind InlineKind, afterKind InputState) {
			if afterKind.Len() > 0 {
				reject(FailedMatch{"unexpected words after kind", afterKind})
			} else {
				// needs to generate the noun with its kind,
				// any specified traits, article properties,
				// and aliases
				accept(InlineNoun{
					InlineKind: kind,
					Called:     called,
					Name:       name,
				}, rhs)
			}
		}, reject)
	}
}

// "the closed transparent container..."
// expects to consume all of the input
// sends itself to the done function.
func TryInlineKind(q JessContext, in InputState,
	accept func(InlineKind, InputState),
	reject func(error)) {
	// a, the, etc.
	TryArticle(q, in, func(article *Article, next InputState) {
		// matches existing traits ( trait names are globally unique )
		TryTraits(q, next, func(traits *Traits, next InputState) {
			// matches an existing kind
			TryKind(q, next, func(kind Kind, rest InputState) {
				// done!
				accept(InlineKind{
					Article: article,
					Traits:  traits,
					Kind:    kind,
				}, rest)
			}, reject)
		}, reject)
	}, reject)
}

func (op *InlineNoun) BuildPropertyNoun(ctx BuildContext) (ret string, err error) {
	if an, e := generateNoun(ctx, ctx, op.Name, op.GetKind(), op.GetTraits()); e != nil {
		err = e
	} else {
		op.actualNoun = an
		ret = an.Name
	}
	return
}

// valid after generation
func (op *InlineNoun) GetActualNoun() ActualNoun {
	return op.actualNoun
}

func (op *InlineNoun) GetKind() string {
	return op.InlineKind.GetKind()
}

func (op *InlineNoun) GetTraits() []string {
	return op.InlineKind.GetTraits()
}

func (op *InlineKind) GetKind() string {
	return op.Kind.actualKind.Name
}

func (op *InlineKind) GetTraits() []string {
	return ReduceTraits(op.Traits)
}

func generateNoun(q Query, w weaver.Weaves, name Name, kind string, traits []string) (ret ActualNoun, err error) {
	// ick.
	if n, e := name.BuildPropertyNoun(q, w, NounProperties{
		Kinds:  []string{kind},
		Traits: traits,
	}); e != nil {
		err = e
	} else if e := n.writeNounValues(w); e != nil {
		err = e
	} else {
		ret = ActualNoun{Name: n.Noun, Kind: kind}
	}
	return
}
