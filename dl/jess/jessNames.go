package jess

// all members of Names implement this so that they can be handled generically
type MatchedName interface {
	BuildNoun(traits, kinds []string) (DesiredNoun, error)
}

func (op *AdditionalNames) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	op.CommaAnd.Match(q, &next) &&
		op.Names.Match(q, &next) {
		*input, okay = next, true
	}
	return
}

// some callers want to fail matching on anonymous leading kinds
// tbd: would it be better to match, and error on generation?
// ( ie. to produce a message )
func (op *Names) HasAnonymousKind() bool {
	return op.Kind != nil
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
		if !Optional(q, &next, &op.AdditionalNames) || op.AdditionalNames.Names.KindCalled == nil {
			*input, okay = next, true
		}
	}
	return
}

func (op *Names) GetMatchedName() (ret MatchedName) {
	if m := op.CountedKind; m != nil {
		ret = m
	} else if m := op.KindCalled; m != nil {
		ret = m
	} else if m := op.Kind; m != nil {
		ret = m
	} else if m := op.Name; m != nil {
		ret = m
	} else if m := op.Noun; m != nil {
		ret = m
	} else {
		panic("well that was unexpected")
	}
	return
}

func (op *Names) BuildNouns(traits, kinds []string) (ret []DesiredNoun, err error) {
	for t := *op; ; {
		m := t.GetMatchedName()
		if m, e := m.BuildNoun(traits, kinds); e != nil {
			err = e
			break
		} else {
			ret = append(ret, m)
			// next name:
			if next := t.AdditionalNames; next == nil {
				break
			} else {
				t = next.Names
			}
		}
	}
	return
}

// unwind the tree of additional names
func (op *Names) Iterate() Iterator {
	return Iterator{op}
}

// unwind the traits ( if any ) of the names
func (op *Names) GetTraits() (ret Traitor) {
	if c := op.KindCalled; c != nil {
		ret = c.GetTraits()
	}
	return
}

type Iterator struct {
	next *Names
}

func (it Iterator) HasNext() bool {
	return it.next != nil
}

func (it *Iterator) GetNext() (ret *Names) {
	var next *Names
	if more := it.next.AdditionalNames; more != nil {
		next = &more.Names
	}
	ret, it.next = it.next, next
	return
}
