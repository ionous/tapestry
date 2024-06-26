package dump

import (
	"database/sql"
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/qna/decoder"
	"git.sr.ht/~ionous/tapestry/qna/raw"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/tables"
)

func QueryPatterns(db *sql.DB, dec decoder.Decoder, scene string) (ret []raw.PatternData, err error) {
	if ps, e := QueryInnerPatterns(db, scene); e != nil {
		err = fmt.Errorf("%w while querying inner patterns", e)
	} else if e := QueryRules(db, dec, scene, ps); e != nil {
		err = fmt.Errorf("%w while querying rules", e)
	} else {
		ret = ps
	}
	return
}

func QueryInnerPatterns(db *sql.DB, scene string) (ret []raw.PatternData, err error) {
	if rows, e := db.Query(must("patterns"), scene); e != nil {
		err = e
	} else {
		var p raw.PatternData
		var labels string
		err = tables.ScanAll(rows, func() (_ error) {
			p.Labels = strings.Split(labels, ",")
			ret = append(ret, p)
			return
		}, &p.Id, &p.Pattern, &labels)
	}
	return
}

func QueryRules(db *sql.DB, dec decoder.Decoder, scene string, ps []raw.PatternData) (err error) {
	q := must("rules")
	for i, n := range ps {
		if rows, e := db.Query(q, scene, n.Id); e != nil {
			err = e
			break
		} else if vs, e := scanRules(rows, dec); e != nil {
			err = e
			break
		} else {
			ps[i].Rules = vs.rules
			ps[i].UpdateAll = vs.updateAll
		}
	}
	return
}

type ruleData struct {
	rules     []rt.Rule
	updateAll bool
}

func scanRules(rows *sql.Rows, dec decoder.Decoder) (ret ruleData, err error) {
	var rule rt.Rule
	var prog []byte
	err = tables.ScanAll(rows, func() (err error) {
		if exe, e := dec.DecodeProg(prog); e != nil {
			err = e
		} else {
			if rule.Updates {
				ret.updateAll = true
			}
			rule.Exe = exe
			ret.rules = append(ret.rules, rule)
		}
		return
	}, &rule.Name, &rule.Stop, &rule.Jump, &rule.Updates, &prog)
	return
}
