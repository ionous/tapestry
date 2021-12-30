package qna

import (
	"database/sql"
	"strings"

	"github.com/ionous/errutil"

	"git.sr.ht/~ionous/iffy"
	"git.sr.ht/~ionous/iffy/affine"
	"git.sr.ht/~ionous/iffy/dl/literal"
	"git.sr.ht/~ionous/iffy/dl/story"
	"git.sr.ht/~ionous/iffy/rt"
	"git.sr.ht/~ionous/iffy/tables"
)

// CheckAll tests stored in the passed db.
// It logs the results of running the tests, and only returns error on critical errors.
func CheckAll(db *sql.DB, actuallyJustThisOne string, options Options, signatures []map[uint64]interface{}) (ret int, err error) {
	if tests, e := readTests(db, actuallyJustThisOne, signatures); e != nil {
		err = e
	} else if len(tests) == 0 {
		err = errutil.New("no matching checks found")
	} else {
		for _, t := range tests {
			// fix? its currently necessary to activate a global domain, rather than jump straight into the check domain.
			// something about pair activation goes a bit wonky: multiple pairs can become active at once.
			tables.Must(db, `delete from run_domain; delete from run_pair`)
			run := NewRuntimeOptions(db, options, iffy.AllSignatures)
			if _, e := run.ActivateDomain("entire_game"); e != nil {
				err = e
			} else {
				if e := t.RunTest(run); e != nil {
					err = errutil.Append(err, e)
				}
				ret++
			}
		}
	}
	return
}

func readTests(db *sql.DB, actuallyJustThisOne string, signatures []map[uint64]interface{}) (ret []CheckOutput, err error) {
	var name, domain string
	var aff affine.Affinity
	var prog []byte
	var value []byte
	//
	if len(actuallyJustThisOne) > 0 {
		actuallyJustThisOne += ";"
	}
	// read all the matching tests from the db.
	// ( cant dynamically query them b/c it interferes with db writes; ex. ActivateDomain )
	err = tables.QueryAll(db,
		`select mc.name, md.domain, mc.value, mc.affinity, mc.prog
		from mdl_check mc
		join mdl_domain md
			on (mc.domain=md.rowid) 
		order by mc.domain, mc.name`,
		func() (err error) {
			if len(actuallyJustThisOne) == 0 || strings.Contains(actuallyJustThisOne, name+";") {
				var act rt.Execute
				if e := story.Decode(rt.Execute_Slot{&act}, prog, signatures); e != nil {
					err = e
				} else if v, e := literal.ReadLiteral(aff, "", value); e != nil {
					err = e
				} else if l, ok := v.(*literal.TextValue); !ok {
					err = errutil.New("can only handle text values right now")
				} else {
					ret = append(ret, CheckOutput{
						Name:   name,
						Domain: domain,
						Expect: l.String(),
						Test:   act,
					})
				}
			}
			return
		}, &name, &domain, &value, &aff, &prog)
	return
}
