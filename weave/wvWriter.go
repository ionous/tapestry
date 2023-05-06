package weave

import (
	"git.sr.ht/~ionous/tapestry/tables"
	"github.com/ionous/errutil"
)

// database/sql like interface
type Writer interface {
	Write(q string, args ...interface{}) error
}

// turn a function into a writer
type WrappedWriter func(q string, args ...interface{}) error

func (w WrappedWriter) Write(q string, args ...interface{}) (err error) {
	if e := w(q, args...); e != nil {
		err = errutil.New("writing", q, e)
	}
	return
}

// turn a transaction, etc into a writer.
func ExecWriter(db tables.Executer) WrappedWriter {
	return func(q string, args ...interface{}) (err error) {
		if _, e := db.Exec(q, args...); e != nil {
			err = errutil.New("exec", q, e)
		}
		return
	}
}
