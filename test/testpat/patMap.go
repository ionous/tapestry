package testpat

import (
	"git.sr.ht/~ionous/tapestry/rt/pattern"
)

// Map - a simple helper to provide patterns w/o a db.
type Map map[string]*Pattern

func (m Map) GetRules(patname string) (ret pattern.RuleSet, err error) {
	if pat, ok := m[patname]; ok {
		for i := len(pat.Rules) - 1; i >= 0; i-- {
			ret.AddRule(pat.Rules[i])
		}
	}
	return
}
