package internal

import (
	"database/sql"

	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/jsn/cin"
	"git.sr.ht/~ionous/tapestry/parser"
	"git.sr.ht/~ionous/tapestry/tables"
)

// fix: domains: rebuild on domain changes, or:
// add a special "AllOf" that does a db query / cache implicitly?
// add scanners which check the database domain?
func MakeGrammar(db *sql.DB) (ret parser.Scanner, err error) {
	var xs []parser.Scanner
	var name string
	var prog []byte
	if e := tables.QueryAll(db,
		`select name, prog  
		from mdl_grammar
		order by rowid`,
		func() (err error) {
			var d grammar.Directive
			if e := cin.Decode(&d, prog, cin.Signatures{grammar.Signatures}); e != nil {
				err = e
			} else {
				x := d.MakeScanners()
				xs = append(xs, x)
			}
			return
		}, &name, &prog,
	); e != nil {
		err = e
	} else {
		ret = &parser.AnyOf{xs}
	}
	return
}
