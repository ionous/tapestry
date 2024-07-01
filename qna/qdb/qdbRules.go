package qdb

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/tables"
)

func (q *Query) getRules(pat string) (ret query.RuleSet, err error) {
	if rows, e := q.rulesFor.Query(pat); e != nil {
		err = e
	} else {
		var rs query.RuleSet
		var rule struct {
			Name    string
			Stop    bool
			Jump    rt.Jump
			Updates bool
			Prog    []byte
		}
		if e := tables.ScanAll(rows, func() (err error) {
			if prog, e := q.dec.DecodeProg(rule.Prog); e != nil {
				err = fmt.Errorf("%w decoding %q rule %v", e, pat, rule.Name)
			} else {
				rs.Rules = append(rs.Rules, rt.Rule{
					Name:    rule.Name,
					Stop:    rule.Stop,
					Jump:    rule.Jump,
					Updates: rule.Updates,
					Exe:     prog,
				})
			}
			if rule.Updates {
				rs.UpdateAll = true
			}
			return
		}, &rule.Name, &rule.Stop, &rule.Jump, &rule.Updates, &rule.Prog); e != nil {
			err = e
		} else {
			ret = rs
		}
	}
	return
}
