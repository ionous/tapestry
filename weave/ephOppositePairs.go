package weave

import "github.com/ionous/errutil"

// a one-to-one mapping of a word and its opposite
type OppositePair struct{ one, other string }
type OppositePairs []OppositePair

func (ps *OppositePairs) AddPair(a, b string) (err error) {
	if a == b || len(a) == 0 || len(b) == 0 {
		err = errutil.Fmt("can't make opposites of %q and %q", a, b)
	} else {
		switch p, n := ps.matches(a, b); {
		case n == 0:
			*ps = append(*ps, OppositePair{a, b})
		case n < 0:
			err = errutil.Fmt("trying to add a pair of opposites (%q, %q) where one or the other already exists as (%q, %q)",
				a, b, p.one, p.other)
		default:
			// already have the exact match
		}
	}
	return
}

func (ps OppositePairs) FindOpposite(word string) (ret string, okay bool) {
	for _, p := range ps {
		if p.one == word {
			ret, okay = p.other, true
		} else if p.other == word {
			ret, okay = p.one, true
		}
	}
	return
}

func (ps OppositePairs) matches(a, b string) (retPair OppositePair, retMatch int) {
	for _, p := range ps {
		if n := p.match(OppositePair{a, b}); n != 0 {
			retPair, retMatch = p, n
			break
		} else if n := p.match(OppositePair{b, a}); n != 0 {
			retPair, retMatch = p, n
			break
		}
	}
	return
}

// tests if a matches b
// 0 no match. 1 exact match. -1 mismatch.
func (a OppositePair) match(b OppositePair) (ret int) {
	one, other := a.one == b.one, a.other == b.other
	if one && other {
		ret = 1
	} else if one || other {
		ret = -1
	}
	return
}
