package internal

import (
	"database/sql"

	"git.sr.ht/~ionous/iffy/dl/grammar"
	"git.sr.ht/~ionous/iffy/parser"
	"git.sr.ht/~ionous/iffy/tables"
)

// re: domains: rebuild on domain changes?
// add a special "AllOf" that does a db query / cache implicitly?
// add scanners which check the database domain?
func MakeGrammar(db *sql.DB) (ret parser.Scanner, err error) {
	var xs []parser.Scanner
	var prog []byte
	if e := tables.QueryAll(db,
		`select bytes from mdl_prog where type='Grammar' order by rowid`,
		func() (err error) {
			var gram grammar.Grammar
			if e := tables.DecodeGob(prog, &gram); e != nil {
				err = e
			} else {
				x := gram.Scanner.MakeScanner()
				xs = append(xs, x)
			}
			return
		}, &prog,
	); e != nil {
		err = e
	} else {
		ret = &parser.AnyOf{xs}
	}
	return
}
