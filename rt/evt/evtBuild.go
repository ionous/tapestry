package evt

import (
	"strings"

	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
)

// 4 phases: capture, target, bubble, after.
type Rules [4][]rt.Rule

// note: scope has to be established before BuildPath gets called
func BuildPath(run rt.Runtime, event string, up []string, allFlags *rt.Flags) (ret Rules, err error) {
	// loop over all targets
	for i, tgt := range up {
		if ks, e := AncestryOf(run, tgt); e != nil {
			err = e
		} else {
			// process the target, and its kinds
			for {
				var tgtFlags rt.Flags
				if rules, e := run.GetRules(event, tgt, &tgtFlags); e != nil {
					err = e
					break
				} else if cnt := len(rules); cnt > 0 {
					var j int
					for p := 0; p < 4; p++ {
						flags := rt.Flags(1 << p)
						if flags&tgtFlags != 0 {
							var set []rt.Rule
							for ; (j < cnt) && (rules[j].Flags&flags) != 0; j++ {
								// for now, skip adding the questionable rows
								// fix? perhaps GetRules() could exclude them based on flags.
								if i == 0 || p != 1 {
									set = append(set, rules[j])
								}
							}
							ret[p] = append(ret[p], set...)
						}
					}
				}
				if allFlags != nil {
					*allFlags |= tgtFlags
				}
				// now its class
				if len(ks) == 0 {
					break
				} else {
					tgt, ks = ks[0], ks[1:]
				}
			}
		}
	}
	return
}

// return the ancestors of the passed noun as a slice of strings
// root is to the right
// ex. maybe props,things,objects,kinds
func AncestryOf(run rt.Runtime, noun string) (ret []string, err error) {
	if kinds, e := run.GetField(object.Kinds, noun); e != nil {
		err = e
	} else {
		// fix? maybe kinds should be returning text list
		ret = strings.FieldsFunc(kinds.String(), func(b rune) bool { return b == ',' })
	}
	return
}
