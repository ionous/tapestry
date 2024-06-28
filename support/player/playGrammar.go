package player

import (
	"database/sql"
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/parser"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/tables"
)

func MakeGrammar(db *sql.DB) (ret []parser.Scanner, err error) {
	var dec decode.Decoder
	dec.Signatures(grammar.Z_Types.Signatures, call.Z_Types.Signatures, literal.Z_Types.Signatures).
		Customize(literal.CustomDecoder)
	q := (*query.QueryDecoder)(&dec)
	return ReadGrammar(db, q)
}

// fix: domains: rebuild on domain changes, or:
// add a special "AllOf" that does a db query / cache implicitly?
// add scanners which check the database domain?
func ReadGrammar(db *sql.DB, dec qdb.CommandDecoder) (ret []parser.Scanner, err error) {
	if rows, e := db.Query(
		`select name, prog  
		from mdl_grammar
		order by rowid`); e != nil {
		err = e
	} else {
		ret, err = ScanGrammar(rows, dec)
	}
	return
}

func ScanGrammar(rows *sql.Rows, dec qdb.CommandDecoder) (ret []parser.Scanner, err error) {
	var name string
	var prog []byte
	err = tables.ScanAll(rows, func() (err error) {
		var dir grammar.Directive
		if e := dec.DecodeValue(&dir, prog); e != nil {
			err = fmt.Errorf("%w decoding directive %q", e, name)
		} else {
			ret = append(ret, dir.MakeScanners())
		}
		return
	}, &name, &prog)
	return
}
