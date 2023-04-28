package weave

import "git.sr.ht/~ionous/tapestry/tables"

// database/sql like interface
type Writer interface {
	Write(q string, args ...interface{}) error
}

// turn a function into a writer
type WrappedWriter func(q string, args ...interface{}) error

func (w WrappedWriter) Write(q string, args ...interface{}) error {
	return w(q, args...)
}

// turn a transaction, etc into a writer.
func ExecWriter(db tables.Executer) WrappedWriter {
	return func(q string, args ...interface{}) (err error) {
		if _, e := db.Exec(q, args...); e != nil {
			err = e
		}
		return
	}
}
