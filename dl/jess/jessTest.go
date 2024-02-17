package jess

import (
	"errors"
	"fmt"

	"git.sr.ht/~ionous/tapestry/support/grok"
)

// exists for testing the matching of trait sets
// not used directly during normal matching.
func MatchTraits(q Query, in InputState) (ret grok.TraitSet, err error) {
	var traits Traits
	if next := in; //
	!traits.Match(q, &next) {
		err = errors.New("failed to match traits")
	} else {
		// after the traits, there might be a kind
		var kind *Kind
		Optional(q, &next, &kind)
		if cnt := len(next); cnt != 0 {
			err = fmt.Errorf("partially matched %d traits", len(in)-cnt)
		} else {
			var k grok.Matched
			if kind != nil {
				k = kind.Matched
			}
			out := grok.TraitSet{Kind: k}
			for ts := traits.GetTraits(); ts.HasNext(); {
				t := ts.GetNext()
				out.Traits = append(out.Traits, t.Matched)
			}
			ret = out
		}
	}
	return
}
