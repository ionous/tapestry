package testpat

import (
	"git.sr.ht/~ionous/tapestry/rt"
)

// Map - a simple helper for testing to provide patterns w/o a db.
type Map map[string]*Pattern

func (m Map) GetRules(patname, target string, pflags *rt.Flags) (ret []rt.Rule, err error) {
	if len(target) > 0 {
		patname += "." + target
	}
	if pat, ok := m[patname]; ok {
		inds, allFlags := SortRules(pat.Rules)
		ret = make([]rt.Rule, len(inds))
		for i, j := range inds {
			ret[i] = pat.Rules[j]
		}
		if pflags != nil {
			*pflags = allFlags
		}
	}
	return
}

// rules are sorted first, and then matched so rules can affect each other.
// except for some tests -- the sorting happens at assembly time.
func SortRules(rules []rt.Rule) (ret []int, retFlags rt.Flags) {
	var ls [rt.LastPhase + 1][]int
	cnt := len(rules)
	for i := cnt - 1; i >= 0; i-- {
		rule := rules[i]
		if flags := rule.Flags(); flags != -1 {
			ofs := flags.Ordinal()
			ls[ofs] = append(ls[ofs], i)
			retFlags |= flags
		}
	}
	if infixOnly := rt.Infix.Ordinal(); retFlags.Ordinal() == infixOnly {
		ret = ls[infixOnly] // this is the most common
	} else {
		for _, els := range ls {
			ret = append(ret, els...)
		}
	}
	return
}
