// Copyright (C) 2022 - Simon Travis. All rights reserved.
// Use of this source code is governed by the Hippocratic 2.1
// license that can be found in the LICENSE file.
package idlcmd

import (
	"context"
	"database/sql"
	"os"
	"path"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
	"git.sr.ht/~ionous/tapestry/idl"
	"git.sr.ht/~ionous/tapestry/tables"
	idlrow "git.sr.ht/~ionous/tapestry/tables/idl"
	"github.com/ionous/errutil"
)

var CmdIdl = &base.Command{
	Run:       runGenerate,
	UsageLine: "tap idlb [file.db]", //  [-run regexp] [-n] [-v] [-x] [build flags] [file.go... | packages]
	Short:     "generate a language database",
	Long: `
Idlb generates an sqlite database containing the tapestry command language.
Currently, language extensions must be built into tapestry itself
( because they require an implementation that tapestry can execute )
so there are no options for this command except the output filename.
	`,
}

var specs = idl.Specs

// tbd: might consider implementing this in an idlb sub-package of idl
// so that we can build this for dynamic lookups
func runGenerate(ctx context.Context, cmd *base.Command, args []string) (err error) {
	var queue writeCache
	if len(args) != 1 {
		err = base.UsageError{Cmd: cmd, Cause: errutil.New("expected a filename")}
	} else {
		outFile := args[0]
		// see also: assembleCat
		if outFile, e := filepath.Abs(outFile); e != nil {
			err = e
		} else if e := os.Remove(outFile); e != nil && !os.IsNotExist(e) {
			err = errutil.New("couldn't clean output file", outFile, e)
		} else {
			// 0755 -> readable by all but only writable by the user
			// 0700 -> read/writable by user
			// 0777 -> ModePerm ... read/writable by all
			os.MkdirAll(path.Dir(outFile), os.ModePerm)
			if db, e := sql.Open(tables.DefaultDriver, outFile); e != nil {
				err = errutil.New("couldn't open db", outFile, e)
			} else {
				defer db.Close()
				if ts, e := rs.FromSpecs(specs); e != nil {
					err = e // returns the .ifspecs as a map of spec.TypeSpec
				} else if e := generateIdb(&queue, ts); e != nil {
					err = e
				} else if e := tables.CreateIdl(db); e != nil {
					err = e
				} else if tx, e := db.Begin(); e != nil {
					err = errutil.New("couldnt create transaction", e)
				} else {
					for _, q := range queue {
						if _, e := tx.Exec(q.tgt, q.args...); e != nil {
							tx.Rollback()
							err = errutil.New("couldnt write to", q.tgt, e)
							break
						}
					}
					if err == nil {
						if e := tx.Commit(); e != nil {
							err = errutil.New("couldnt commit", e)
						}
					}
				}
			}
		}
	}
	return
}

// a terrible way to optimize database writes
type cachedWrite struct {
	tgt  string
	args []interface{}
}
type writeCache []cachedWrite

func (q *writeCache) Write(tgt string, args ...interface{}) (err error) {
	(*q) = append(*q, cachedWrite{tgt, args})
	return
}

func generateIdb(w writer, ts rs.TypeSpecs) (err error) {

	// sorted list of keys ( for some stability of row ids when regenerating the database
	for _, key := range ts.Keys() {
		t := ts.Types[key]
		name := t.Name
		pack := t.Groups[0]
		uses := strings.ToLower(t.Spec.Choice[1:]) // $FLOW -> flow
		var closed bool
		switch uses := t.Spec.Value.(type) {
		case *spec.StrSpec:
			closed = uses.Exclusively
		case *spec.NumSpec:
			closed = uses.Exclusively
		}
		if e := w.Write(idlrow.Op, name, pack, uses, closed); e != nil {
			err = e
			break
		}
	}
	return
}

// database/sql like interface
type writer interface {
	Write(q string, args ...interface{}) error
}
