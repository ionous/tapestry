package pattern

import (
	r "reflect"

	"git.sr.ht/~ionous/iffy/rt"
	"github.com/ionous/errutil"
)

// Map - a simple helper for testing to provide patterns w/o a db.
type Map map[string]*Pattern

// skip assembling the pattern from the db
// we just want to test we can invoke a pattern successfully.
// pv is a pointer to a pattern instance, and we copy its contents in.
func (m Map) GetEvalByName(name string, pv interface{}) (err error) {
	if patternPtr, ok := m[name]; !ok {
		err = errutil.New("unknown pattern", name)
	} else {
		stored := r.ValueOf(patternPtr).Elem()
		outVal := r.ValueOf(pv).Elem()
		outVal.Set(stored)
	}
	return
}

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
