package cmdgenerate

import (
	"database/sql"
	"log"

	"strconv"

	"git.sr.ht/~ionous/tapestry/tables/idl"
)

type modelWriter struct {
	db *sql.DB
	tx *sql.Tx
}

func (m modelWriter) Close() (err error) {
	if m.db != nil {
		if e := m.tx.Commit(); e != nil {
			log.Println("couldnt commit", e)
		}
		err = m.db.Close()
	}
	return
}

func (m modelWriter) Write(q string, args ...interface{}) (err error) {
	var out string
	if sel, ok := idWriter[q]; ok {
		out = sel
	} else {
		out = q
	}
	_, err = m.tx.Exec(out, args...)
	return
}

// selects from idl_<key> where <key>=?<arg>
func opId(a int) string {
	return `(select rowid from idl_op where name=?` + strconv.Itoa(a) + `)`
}

// rewrite some tables to use ids
// the key of the table is the original, simplified insert statement
// the value is a more complex statement usually involving selects
var idWriter = map[string]string{
	idl.Op: idl.Op,
	idl.Markup: `insert into idl_markup( op, key, value ) values (` +
		opId(1) + // op name
		`, ?2` + // key
		`, ?3` + // value
		`)`,
	idl.Sig: `insert into idl_sig( op, slot, hash, body ) values (` +
		opId(1) + // op name
		`, ` + opId(2) + // op name
		`, ?3` + // hash
		`, ?4` + // body
		`)`,
	idl.Enum: `insert into idl_enum( op,  value ) values (` +
		opId(1) + // op name
		`, ?2` + // value
		`)`,
	idl.Term: `insert into idl_term(
		op, name, label, type, private, optional, repeats
	) values (` +
		opId(1) + // parent flow -> an op reference
		`, ?2` + // name
		`, ?3, ` + // label
		opId(4) + // type -> an op reference
		`, ?5` + // private bool
		`, ?6` + // optional bool
		`, ?7` + // repeats bool
		`)`,
}
