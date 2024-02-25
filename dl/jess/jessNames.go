package jess

// all members of Names implement this so that they can be handled generically
type MatchedName interface {
	GetName(traits, kinds []string) resultName
	String() string
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

// checks Query flags for PlainNameMatching
func (op *Names) Match(q Query, input *InputState) (okay bool) {
	if next := *input; //
	(matchKinds(q) &&
		(Optional(q, &next, &op.CountedName) ||
			Optional(q, &next, &op.KindCalled) ||
			Optional(q, &next, &op.Kind))) ||
		Optional(q, &next, &op.Name) {
		// as long as one succeeded, try matching additional nouns too...
		// FIX: as far as i can tell, inform only allows "kind called" at the front of a list of names
		// maybe it'd be better to match and reject when used.
		if !Optional(q, &next, &op.AdditionalNames) || op.AdditionalNames.Names.KindCalled == nil {
			*input, okay = next, true
		}
	}
	return
}

func (op *Names) Pick() (ret MatchedName) {
	if n := op.CountedName; n != nil {
		ret = n
	} else if n := op.KindCalled; n != nil {
		ret = n
	} else if n := op.Kind; n != nil {
		ret = n
	} else if n := op.Name; n != nil {
		ret = n
	} else {
		panic("well that was unexpected")
	}
	return
}

func (op *Names) GetName(traits, kinds []string) (ret resultName) {
	return op.Pick().GetName(traits, kinds)
}

// return the match of this, without any additional nouns
// panics if there wasn't actually a match
func (op *Names) String() (ret string) {
	return op.Pick().String()
}

func (op *Names) GetNames(traits, kinds []string) (ret []resultName) {
	for t := *op; ; {
		n := t.GetName(traits, kinds)
		ret = append(ret, n)
		// next name:
		if next := t.AdditionalNames; next == nil {
			break
		} else {
			t = next.Names
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
