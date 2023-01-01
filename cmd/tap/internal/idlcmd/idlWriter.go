package idlcmd

import (
	"strconv"

	"git.sr.ht/~ionous/tapestry/dl/eph"
	"git.sr.ht/~ionous/tapestry/tables/idl"
)

// return an object implementing the eph catalog writer
// remapping the simple table definitions from mdl.go ( package tables/mdl )
// to ones that can look up and store ids.
// ex. so a caller can Write(mdl.Check, raw args) and have the args remapped to ids.
func NewSpecWriter(fn writerFn) eph.Writer {
	return modelWriter{fn}
}

type writerFn func(q string, args ...interface{}) error
type modelWriter struct{ fn writerFn }

func (m modelWriter) Write(q string, args ...interface{}) (err error) {
	var out string
	if sel, ok := idWriter[q]; ok {
		out = sel
	} else {
		out = q
	}
	return m.fn(out, args...)
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
	idl.Sig: `insert into idl_sig( op, slot, hash, signature ) values (` +
		opId(1) + // op name
		`, ` + opId(2) + // op name
		`, ?3` + // hash
		`, ?4` + // signature
		`)`,
	idl.Enum: `insert into idl_enum( op, label, value ) values (` +
		opId(1) + // op name
		`, ?2` + // label
		`, ?3` + // value
		`)`,
	idl.Swap: `insert into idl_swap( op, label, type ) values (` +
		opId(1) + // op name
		`, ?2, ` + // label
		opId(3) + // type int -> an op reference
		`)`,
	idl.Term: `insert into idl_term(
		op, field, label, type, private, optional, repeats
	) values (` +
		opId(1) + // op name
		`, ?2` + // field
		`, ?3, ` + // label
		opId(4) + // type int -> an op reference
		`, ?5` + // private bool
		`, ?6` + // optional bool
		`, ?7` + // repeats bool
		`)`,
}
