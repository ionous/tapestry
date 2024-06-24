package dump

import (
	"database/sql"
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/tables"
)

func QueryPatterns(db *sql.DB, scene string) (ret []PatternData, err error) {
	if ps, e := QueryInnerPatterns(db, scene); e != nil {
		err = fmt.Errorf("%w while querying inner patterns", e)
	} else if e := QueryRules(db, scene, ps); e != nil {
		err = fmt.Errorf("%w while querying rules", e)
	} else {
		ret = ps
	}
	return
}

func QueryInnerPatterns(db *sql.DB, scene string) (ret []PatternData, err error) {
	var p PatternData
	var labels string
	if rows, e := db.Query(must("patterns"), scene); e != nil {
		err = e
	} else {
		err = tables.ScanAll(rows, func() (_ error) {
			p.Labels = strings.Split(labels, ",")
			ret = append(ret, p)
			return
		}, &p.Id, &p.Pattern, &labels, &p.Result)
	}
	return
}

func QueryRules(db *sql.DB, scene string, ps []PatternData) (err error) {
	q := must("rules")
	for i, n := range ps {
		if rows, e := db.Query(q, scene, n.Id); e != nil {
			err = e
		} else if vs, e := queryRules(rows); e != nil {
			err = e
		} else {
			ps[i].Rules = vs
		}
	}
	return
}

func queryRules(rows *sql.Rows) (ret []query.RuleData, err error) {
	var rule query.RuleData
	err = tables.ScanAll(rows, func() (_ error) {
		ret = append(ret, rule)
		return
	}, &rule.Name, &rule.Stop, &rule.Jump, &rule.Updates, &rule.Prog)
	return
}
