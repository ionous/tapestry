package player

import (
	"database/sql"
	"encoding/json"

	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/grammar"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/lang/decode"
	"git.sr.ht/~ionous/tapestry/parser"
	"git.sr.ht/~ionous/tapestry/qna/raw"
	"git.sr.ht/~ionous/tapestry/tables"
)

func MakeGrammar(db *sql.DB) (ret parser.Scanner, err error) {
	var dec decode.Decoder
	return ReadGrammar(db, dec.
		Signatures(grammar.Z_Types.Signatures, call.Z_Types.Signatures, literal.Z_Types.Signatures).
		Customize(literal.CustomDecoder))
}

// fix: domains: rebuild on domain changes, or:
// add a special "AllOf" that does a db query / cache implicitly?
// add scanners which check the database domain?
func ReadGrammar(db *sql.DB, dec *decode.Decoder) (ret parser.Scanner, err error) {
	g := grammarBuilder{dec: dec}
	var name string
	var prog []byte
	if e := tables.QueryAll(db,
		`select name, prog  
		from mdl_grammar
		order by rowid`,
		func() error {
			return g.addBytes(name, prog)
		}, &name, &prog,
	); e != nil {
		err = e
	} else {
		ret = g.getScanner()
	}
	return
}

func readRawGrammar(dec *decode.Decoder, els []raw.Grammar) (ret parser.Scanner, err error) {
	g := grammarBuilder{dec: dec}
	for _, el := range els {
		if e := g.addBytes(el.Name, el.Prog); e != nil {
			err = e
			break
		}
	}
	if err == nil {
		ret = g.getScanner()
	}
	return
}

type grammarBuilder struct {
	dec   *decode.Decoder
	match []parser.Scanner
}

func (g *grammarBuilder) getScanner() parser.Scanner {
	return &parser.AnyOf{Match: g.match}
}

func (g *grammarBuilder) addBytes(name string, prog []byte) (err error) {
	var dir grammar.Directive
	var msg map[string]any
	if e := json.Unmarshal(prog, &msg); e != nil {
		err = e
	} else if e := g.dec.Decode(&dir, msg); e != nil {
		err = e
	} else {
		x := dir.MakeScanners()
		g.match = append(g.match, x)
	}
	return
}
