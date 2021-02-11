package pattern

import "git.sr.ht/~ionous/iffy/rt"

// call apply on each list rule ( in reverse order )
// last added rules wins over earlier rules, and "post fix" rules happen at the end.
// FIX: presort everything in the assembler?
func sortRules(run rt.Runtime, rules []*Rule) (ret []int, err error) {
	var pre, post int
	cnt := len(rules)
	a := make([]int, cnt)
	for i := cnt - 1; i >= 0; i-- {
		flags := rules[i].Flags
		// NOTE: previously we matched the rules, and then ran them
		// now: they are matched in order, so rules can affect each other.
		// and, if a pattern has a return value, the pattern stops only when a value gets set.
		//
		// apply the rule:
		// returns flags if the filters passed, -1 if they did not, error on any error.
		/*if flags, e := rules[i].ApplyRule(run); e != nil {
			err = e
			break
		} else */if flags >= 0 {
			if flags == Postfix {
				end := cnt - post - 1
				a[end], post = i, post+1
			} else {
				// FIX: add Replace as well
				a[pre], pre = i, pre+1
				/*if flags == Terminal {
					break
				}*/
			}
		}
	}
	if err == nil {
		if pre+post == cnt {
			ret = a
		} else {
			// keep all the prefixed items
			ret = a[:pre]
			// shift the post fixed items into the spot just after the prefixed items
			startOfPost := cnt - post
			for i := startOfPost; i < cnt; i++ {
				ret = append(ret, a[i])
			}
		}
	}
	return
}
