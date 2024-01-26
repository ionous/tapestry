package play

import (
	"database/sql"
	"encoding/json"

	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/lang/decode"
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
			var dir grammar.Directive
			var msg map[string]any
			if e := json.Unmarshal(prog, &msg); e != nil {
				err = e
			} else {
				var dec decode.Decoder
				if e := dec.
					Signatures(grammar.Z_Types.Signatures, assign.Z_Types.Signatures, literal.Z_Types.Signatures).
					Customize(literal.CustomDecoder).
					Decode(&dir, msg); e != nil {
					err = e
				} else {
					x := dir.MakeScanners()
					xs = append(xs, x)
				}
			}
			return
		}, &name, &prog,
	); e != nil {
		err = e
	} else {
		ret = &parser.AnyOf{Match: xs}
	}
	return
}
