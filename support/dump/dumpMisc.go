package dump

import (
	"database/sql"
	"fmt"

	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/qna/raw"
	"git.sr.ht/~ionous/tapestry/tables"
)

// fix? this is a little different than the way the normal player works
// this includes kinds
func QueryPlurals(db *sql.DB, scene string) (ret []raw.Plural, err error) {
	var p raw.Plural
	if rows, e := db.Query(must("plurals"), scene); e != nil {
		err = fmt.Errorf("%w while querying plurals", e)
	} else {
		err = tables.ScanAll(rows, func() (_ error) {
			ret = append(ret, p)
			return
		}, &p.One, &p.Many)
	}
	return
}

// directives produce scanners so that scanners can live separate from tapestry/commands
func QueryGrammar(db *sql.DB, dec qdb.CommandDecoder, scene string) (ret []grammar.Directive, err error) {
	if rows, e := db.Query(must("grammar"), scene); e != nil {
		err = fmt.Errorf("%w while querying plurals", e)
	} else {
		ret, err = scanGrammar(rows, dec)
	}
	return
}

func scanGrammar(rows *sql.Rows, dec qdb.CommandDecoder) (ret []grammar.Directive, err error) {
	var name string
	var prog []byte
	err = tables.ScanAll(rows, func() (err error) {
		var dir grammar.Directive
		if e := dec.DecodeValue(&dir, prog); e != nil {
			err = fmt.Errorf("%w decoding directive %q", e, name)
		} else {
			ret = append(ret, dir)
		}
		return
	}, &name, &prog)
	return
}
