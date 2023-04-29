// Copyright (C) 2022 - Simon Travis. All rights reserved.
// Use of this source code is governed by the Hippocratic 2.1
// license that can be found in the LICENSE file.
package cmdidlb

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/dl/spec"
	"git.sr.ht/~ionous/tapestry/dl/spec/rs"
	"git.sr.ht/~ionous/tapestry/idl"
	"git.sr.ht/~ionous/tapestry/support/distill"
	"git.sr.ht/~ionous/tapestry/tables"
	idlrow "git.sr.ht/~ionous/tapestry/tables/idl"
	"github.com/ionous/errutil"
)

var CmdIdl = &base.Command{
	Run:       runGenerate,
	UsageLine: "tap idlb [file.db]", //  [-run regexp] [-n] [-v] [-x] [build flags] [file.go... | packages]
	Short:     "language database generator",
	Long: `
Idlb generates an sqlite database containing the tapestry command language.
Language extensions must be built into tapestry itself
( they require an implementation that tapestry can execute )
so there are no options for this command except the output filename.`,
}

// all specs
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
			log.Println("generating", outFile)
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
				} else {
					generateIdb(&queue, ts)
					if e := tables.CreateIdl(db); e != nil {
						err = e
					} else if tx, e := db.Begin(); e != nil {
						err = errutil.New("couldnt create transaction", e)
					} else {
						w := newSpecWriter(func(q string, args ...interface{}) (err error) {
							if _, e := tx.Exec(q, args...); e != nil {
								err = e
							}
							return
						})
						if e := queue.Flush(w); e != nil {
							tx.Rollback()
							err = errutil.New("couldnt write", e)
						} else if e := tx.Commit(); e != nil {
							err = errutil.New("couldnt commit", e)
						}
					}
				}
			}
		}
	}
	return
}

// doesn't return an error; assumes writer will handle this is internally in some implementation dependent way.
func generateIdb(w writer, ts rs.TypeSpecs) {
	ds := distill.MakeRegistry(ts.Types)
	// fix: it makes no sense to have to build the registry from the types
	// and then add them as a separate step
	for _, t := range ts.Types {
		ds.AddType(t)
	}
	ds.Sort() // sorted list of keys ( for some stability of row ids when regenerating )
	for _, key := range ds.Types {
		var t *spec.TypeSpec = ts.Types[key]
		var label string
		name := t.Name
		pack := t.Groups[0]
		uses := strings.ToLower(t.Spec.Choice[1:]) // $FLOW -> flow
		switch uses := t.Spec.Value.(type) {
		case *spec.FlowSpec:
			if len(uses.Name) > 0 {
				label = uses.Name
			} else {
				label = name
			}
			label = distill.Pascal(label)
			writeTerms(w, ts, t, uses)
		case *spec.StrSpec:
			label = closedOrOpenLabel(uses.Exclusively)
			writeEnumStrings(w, ts, t, uses)
		case *spec.NumSpec:
			label = closedOrOpenLabel(uses.Exclusively)
			writeEnumNumbers(w, ts, t, uses)
		case *spec.SwapSpec:
			if len(uses.Name) > 0 {
				label = uses.Name
			} else {
				label = name
			}
			label = distill.Pascal(label)
			writeChoices(w, ts, t, uses)
		}
		w.Write(idlrow.Op, name, pack, uses, label)
	}
	for _, sig := range ds.Sigs {
		w.Write(idlrow.Sig, sig.Type, sig.Slot, strconv.FormatUint(sig.Hash, 16), sig.Body())
	}
	return
}

func closedOrOpenLabel(exclusively bool) (ret string) {
	if exclusively {
		ret = "$CLOSED"
	} else {
		ret = "$OPEN"
	}
	return
}

func writeChoices(w writer, ts rs.TypeSpecs, t *spec.TypeSpec, swap *spec.SwapSpec) {
	for _, opt := range swap.Between {
		w.Write(idlrow.Swap, t.Name, opt.Key(), opt.Value(), opt.TypeName())
	}
}
func writeEnumStrings(w writer, ts rs.TypeSpecs, t *spec.TypeSpec, str *spec.StrSpec) {
	for _, opt := range str.Uses {
		w.Write(idlrow.Enum, t.Name, opt.Key(), opt.Value())
	}
}
func writeEnumNumbers(w writer, ts rs.TypeSpecs, t *spec.TypeSpec, num *spec.NumSpec) {
	if len(num.Uses) > 0 {
		panic("not implemented")
	}
}
func writeTerms(w writer, ts rs.TypeSpecs, t *spec.TypeSpec, uses *spec.FlowSpec) {
	for _, f := range uses.Terms {
		w.Write(idlrow.Term, t.Name, f.Field(), f.QuietLabel(), f.TypeName(), f.Private, f.Optional, f.Repeats)
		// 	Markup   map[string]any
	}
}

// database/sql like interface
type writer interface {
	Write(q string, args ...interface{}) error
}
