package grok

import "github.com/ionous/errutil"

// after finding is/are, check for a macro;
// handle either 1. a pure "is" statement, or 2. an "is-verb" statement.
// 1.  (lhs: The bottle) is (rhs: closed.)
// 2a. (lhs: A device) is (rhs: in the lobby.)
// 2b. (lhs: The coffin) is (rhs: a closed container in the lobby.)
// 2c. (lhs: The coffin) is (rhs: closed in the lobby.)  <-- not my favorite.
//
// tbd: parse the sources when looping over the words ( in the caller? )
func beingPhrase(out *Results, lhs, rhs []Word) (err error) {
	// first, scan for leading traits on the rhs
	// ex. [is] ( rhs: fixed in place .... in the lobby )
	if rightLede, e := parseTraitSet(rhs); e != nil {
		err = e
	} else {
		// try to find a macro after the traits:
		afterRightLede := rhs[rightLede.wordCount:]
		if macro, ok := known.FindMacro(afterRightLede); !ok {
			// case 1. doesn't have a macro:
			if e := genNouns(&out.Sources, lhs, AllowMany|AllowAnonymous); e != nil {
				err = errutil.New("parsing subjects", e)
			}
		} else {
			// case 2: found a macro:
			out.Macro = macro
			postMacro := afterRightLede[macro.Width():]
			var lhsFlag, rhsFlag genFlag
			switch macro.Type() {
			case OneToMany, ManyToOne:
				lhsFlag, rhsFlag = AllowMany|OnlyNamed, OnlyOne|AllowAnonymous
			case ManyToMany:
				lhsFlag, rhsFlag = AllowMany|OnlyNamed, AllowMany|OnlyNamed
			}

			// [lhs: The coffin is] (rhs: (pre: a closed container) *in* (post: the antechamber.))
			if e := genNouns(&out.Sources, lhs, lhsFlag); e != nil {
				err = errutil.New("parsing subject", e)
			} else {
				// fix? this branching isnt satisfying
				if macro.Type() > ManyToOne {
					if e := genNouns(&out.Targets, postMacro, rhsFlag); e != nil {
						err = errutil.New("parsing target", e)
					}
				} else {
					// inform specifically denies these right leading traits in this case:
					// [The box is] (right lede: a closed container) kind of (post traits: closed container).
					if rightLede.wordCount > 0 {
						err = makeWordError(rhs[0], "some unexpected kind of properties")
					} else if postMacroTraits, e := parseTraitSet(postMacro); e != nil {
						err = e
					} else {
						postMacroTraits.applyTraits(out.Sources)
					}
				}
			}
		}
		if err == nil {
			rightLede.applyTraits(out.Sources)
		}
	}
	return
}
