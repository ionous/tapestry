package jess

import "git.sr.ht/~ionous/tapestry/rt/kindsOf"

// called can have its own kind, its own specific article, and its name is flagged as "exact"
// ( where regular names are treated as potential aliases of existing names. )
func (op *KindCalled) GetName(traits, kinds []string) (ret resultName, err error) {
	if kind, e := op.GetKind(); e != nil {
		err = e
	} else {
		for ts := op.GetTraits(); ts.HasNext(); {
			t := ts.GetNext()
			traits = append(traits, t.String())
		}
		ret = resultName{
			// ignores the article of the kind,
			// in favor of the article closest to the named noun
			Article: reduceArticle(op.CalledName.Article),
			Matched: op.CalledName.Matched,
			Exact:   true,
			Traits:  traits,
			Kinds:   append(kinds, kind),
		}
	}
	return
}

func (op *KindCalled) GetKind() (string, error) {
	return op.Kind.Validate(kindsOf.Kind)
}

func (op *KindCalled) GetTraits() (ret Traitor) {
	if op.Traits != nil {
		ret = op.Traits.GetTraits()
	}
	return
}

func (op *KindCalled) String() string {
	return op.CalledName.Matched
}

func (op *KindCalled) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(Optional(q, &next, &op.Traits) || true) &&
		op.Kind.Match(q, &next) &&
		op.CalledName.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

func (op *CalledName) String() string {
	return op.Matched
}

func (op *CalledName) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.Called.Match(q, &next, keywords.Called) &&
		(Optional(q, &next, &op.Article) || true) &&
		op.matchName(q, &next) {
		*input, okay = next, true
	}
	return
}

// match the words after "called" ending with either "is/are" or the end of the string.
// this means that something like "a container called the bottle and a woman called the genie"
// generates a single object with a long, strange name.
func (op *CalledName) matchName(q Query, input *InputState) (okay bool) {
	if width := beScan(input.Words()); width > 0 {
		op.Matched, *input, okay = input.Cut(width), input.Skip(width), true
	}
	return
}
