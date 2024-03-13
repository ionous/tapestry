package jess

// some callers want to fail matching on anonymous leading kinds
// tbd: would it be better to match, and error on generation?
// ( ie. to produce a message )
func (op *Names) HasAnonymousKind() bool {
	return op.Kind != nil
}

// unwind the tree of additional names
func (op *Names) GetNames() Iterator {
	return Iterator{op}
}

// unwind the traits ( if any ) of the names
func (op *Names) GetTraits() (ret Traitor) {
	if c := op.KindCalled; c != nil {
		ret = c.GetTraits()
	}
	return
}

// checks Query flags to control matching
func (op *Names) Match(q Query, input *InputState) (okay bool) {
	matchNouns := matchNouns(q)
	matchKinds := matchKinds(q)
	if next := *input; ( //
	matchNouns &&
		// "the bottle"
		Optional(q, &next, &op.Noun)) || ( //
	matchKinds &&
		// "5 containers",
		Optional(q, &next, &op.CountedKind) ||
		// "the container called the bottle"
		Optional(q, &next, &op.KindCalled) ||
		// "the container"
		Optional(q, &next, &op.Kind)) || ( //
	// "the unknown name"
	Optional(q, &next, &op.Name)) {
		// as long as one succeeded, try matching additional names too...
		// inform seems to only allow "kind called" at the front of a list of names...
		// fix? but maybe it'd be better to match and reject when used incorrectly.
		// [ an advantage is -- we could register a whole side ( lhs/rhs ) to nouns at once ]
		if !Optional(q, &next, &op.AdditionalNames) || op.AdditionalNames.Names.KindCalled == nil {
			*input, okay = next, true
		}
	}
	return
}

// implements NounBuilder by calling BuildNouns on all matched names
func (op Names) BuildNouns(q Query, rar *Context, ts, ks []string) (ret []DesiredNoun, err error) {
	for n := op.GetNames(); n.HasNext(); {
		at := n.GetNext()
		if ns, e := buildNounsFrom(q, rar, ts, ks,
			ref(at.CountedKind),
			ref(at.KindCalled),
			ref(at.Kind),
			ref(at.Name),
			ref(at.Noun), //
		); e != nil {
			err = e
			break
		} else {
			ret = append(ret, ns...)
		}
	}
	return
}
