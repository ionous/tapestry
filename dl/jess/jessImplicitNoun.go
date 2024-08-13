package jess

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/weave/weaver"
)

// a phrase implies a noun exists, but no particular kind has been assigned.
func TryImplicitNoun(q JessContext, in InputState,
	accept func(ImplicitNoun, InputState),
	reject func(error),
) {
	// its not clear to me when exactly this should happen
	q.Try(weaver.FallbackPhase, func(weaver.Weaves, rt.Runtime) {
		// tricksy: don't match an implicit noun if it would conflict with a kind
		// ( its probably actually some other phrase )
		TryKind(q, in, func(kind Kind, rest InputState) {
			reject(fmt.Errorf("the phrase implies a noun %q, but there's already a kind of that name", kind.actualKind.Name))
		}, func(error) {
			var n Name
			if next := in; !n.Match(q, &next) {
				reject(FailedMatch{"expected some sort of name", next})
			} else {
				accept(ImplicitNoun{Name: n}, next)
			}
		})
	}, reject)
}

// implicit nouns can only use properties of things.
// there is however an earlier phase check for Noun.
func (op *ImplicitNoun) GetKind() string {
	return Things
}

func (op *ImplicitNoun) BuildPropertyNoun(ctx BuildContext) (ret ActualNoun, err error) {
	// fix: for backwards compatibility with tests, this first creates the noun as "object"
	// and then generates it as Things. i dont remember why the placeholder was necessary.
	// the test output will list the name before the kind when created this way.
	if an, created, e := ensureNoun(ctx, ctx, op.Name.Matched, nil); e != nil {
		err = e
	} else if !created {
		ret = an
	} else if an, e := generateNoun(ctx, ctx, op.Name, Things, nil); e != nil {
		err = e
	} else {
		ret = an
	}
	return
}
