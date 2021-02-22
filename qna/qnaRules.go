package qna

import (
	"database/sql"
	r "reflect"

	"git.sr.ht/~ionous/iffy/lang"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/tables"
)

type qnaRules struct {
	rules    map[string]ruleSet
	rulesFor *sql.Stmt
}

type ruleSet struct {
	err   error
	rules []rt.Rule
	flags rt.Flags
}

func (q *qnaRules) GetRules(pattern, target string, pflags *rt.Flags) (ret []rt.Rule, err error) {
	key := lang.Breakcase(pattern)
	if x, ok := q.rules[key]; ok {
		if e := x.err; e != nil {
			err = e
		} else {
			ret = x.rules
			if pflags != nil {
				*pflags = x.flags
			}
		}
	} else {
		var rulen string
		var phase int
		var handler rt.Handler
		hval := r.ValueOf(&handler).Elem()
		// var
		if rows, e := q.rulesFor.Query(key); e != nil {
			err = e
		} else if e := tables.ScanAll(rows, func() (err error) {
			flags := rt.MakeFlags(phase)
			x.rules = append(x.rules, rt.Rule{
				Name:    rulen,
				Filter:  handler.Filter,
				Execute: handler.Exe,
				Flags:   flags,
			})
			x.flags |= flags
			handler = rt.Handler{} // gob doesnt write nil values
			return
		}, &rulen, &phase, &tables.GobScanner{hval}); e != nil {
			err = e
		}
		if err != nil {
			x = ruleSet{err: err}
		} else {
			ret = x.rules
			if pflags != nil {
				*pflags = x.flags
			}
		}
		if q.rules == nil {
			q.rules = make(map[string]ruleSet)
		}
		q.rules[key] = x
	}
	return
}
