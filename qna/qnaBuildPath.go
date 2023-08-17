package qna

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
)

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
func BuildPath(run rt.Runtime, event string, up []string) (ret Ruleset, err error) {
	// loop over all targets
	for _, target := range up {
		tgt := target
		if ks, e := ancestryOf(run, tgt); e != nil {
			err = e
		} else {
			// process the target, and its kinds
			for cls := ""; ; cls = tgt {
				if rules, e := run.GetRules(event, tgt); e != nil {
					err = e
					break
				} else {
					for i := len(rules) - 1; i >= 0; i-- {
						ret.Add(target, cls, rules[i])
					}
				}
				// now ancestors: root last.
				if cnt := len(ks); cnt == 0 {
					break
				} else {
					tgt, ks = ks[cnt-1], ks[:cnt-1]
				}
			}
		}
	}
	return
}

// return the ancestors of the passed noun as a slice of strings
// root is at the start. ex. kinds,objects,things,props
func ancestryOf(run rt.Runtime, noun string) (ret []string, err error) {
	if kinds, e := run.GetField(meta.ObjectKinds, noun); e != nil {
		err = e
	} else {
		ret = kinds.Strings()
	}
	return
}
