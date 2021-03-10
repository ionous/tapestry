package evt

import (
	"strings"

	"git.sr.ht/~ionous/iffy/object"
	"git.sr.ht/~ionous/iffy/rt"
)

// 4 phases: capture, target, bubble, after
// with a list of rules in each set.
type Rulesets [4]Ruleset

// hold parallel arrays
// [ separated this way to preserve existing ApplyRules which needs a rule slice ]
type Ruleset []Target

type Target struct {
	Noun, Kind string
	rt.Rule
}

func (rs *Ruleset) Add(noun, kind string, rule rt.Rule) {
	(*rs) = append((*rs), Target{noun, kind, rule})
}

func (rs *Ruleset) Append(other Ruleset) {
	(*rs) = append((*rs), other...)
}

// note: scope has to be established before BuildPath gets called
func BuildPath(run rt.Runtime, event string, up []string, allFlags *rt.Flags) (ret Rulesets, err error) {
	// loop over all targets
	for _, target := range up {
		tgt := target
		if ks, e := AncestryOf(run, tgt); e != nil {
			err = e
		} else {
			// process the target, and its kinds
			for cls := ""; ; cls = tgt {
				var tgtFlags rt.Flags // get all of the rules for this target (sorted by increasing phase)
				if rules, e := run.GetRules(event, tgt, &tgtFlags); e != nil {
					err = e
					break
				} else if cnt := len(rules); cnt > 0 {
					var j int                // index into all of the rules
					for p := 0; p < 4; p++ { // separate them by phase
						flags := rt.Flags(1 << p)
						if flags&tgtFlags != 0 {
							var set Ruleset
							for ; (j < cnt) && (rules[j].Flags&flags) != 0; j++ {
								set.Add(target, cls, rules[j])
							}
							ret[p].Append(set)
						}
					}
					if allFlags != nil {
						*allFlags |= tgtFlags
					}
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
// root is to the right: ex. props,things,objects,kinds
func AncestryOf(run rt.Runtime, noun string) (ret []string, err error) {
	if kinds, e := run.GetField(object.Kinds, noun); e != nil {
		err = e
	} else {
		// fix? maybe kinds itself should be returning text list
		ret = strings.FieldsFunc(kinds.String(), func(b rune) bool { return b == ',' })
	}
	return
}
