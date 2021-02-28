package testpat

import (
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/rt/pattern"
)

// Map - a simple helper for testing to provide patterns w/o a db.
type Map map[string]*Pattern

func (m Map) GetRules(patname, target string, pflags *rt.Flags) (ret []rt.Rule, err error) {
	if len(target) > 0 {
		patname += "." + target
	}
	if pat, ok := m[patname]; ok {
		inds, allFlags := pattern.SortRules(pat.Rules)
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
