package testpat

import (
	"git.sr.ht/~ionous/tapestry/rt"
	"golang.org/x/exp/slices"
)

// Map - a simple helper to provide patterns w/o a db.
type Map map[string]*Pattern

func (m Map) GetRules(patname, target string) (ret []rt.Rule, err error) {
	if len(target) > 0 {
		patname += "." + target
	}
	if pat, ok := m[patname]; ok {
		ret = append(ret, pat.Rules...)
		slices.Reverse(ret)
	}
	return
}
