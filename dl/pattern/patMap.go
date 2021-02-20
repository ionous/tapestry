package pattern

import (
	"git.sr.ht/~ionous/iffy/rt"
	"github.com/ionous/errutil"
)

// Map - a simple helper for testing to provide patterns w/o a db.
type Map map[string]*Pattern

func (m Map) GetRules(name string, pflags *rt.Flags) (ret []rt.Rule, err error) {
	if pat, ok := m[name]; !ok {
		err = errutil.New("unknown pattern", name)
	} else {
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
