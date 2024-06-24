package main

// #include <stdlib.h>
import "C"

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"unsafe"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/frame"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/qdb"
	"git.sr.ht/~ionous/tapestry/support/play"
	"git.sr.ht/~ionous/tapestry/support/player"
	"git.sr.ht/~ionous/tapestry/tables"
)

//export Post
func Post(endpoint, msg string) (ret *C.char) {
	res, e := post(endpoint, msg)
	if e != nil {
		res = e.Error()
		println("error", res)
	}
	if lastResult != nil {
		C.free(lastResult)
	}
	ret = C.CString(res)
	lastResult = unsafe.Pointer(ret)
	return
}

func post(endpoint, msg string) (ret string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Recovered %s:\n%s", r, debug.Stack())
		}
	}()
	var buf strings.Builder
	if c, e := getShuttle(); e != nil {
		err = e
	} else if e := c.Post(&buf, endpoint, []byte(msg)); e != nil {
		err = e
	} else {
		ret = buf.String()
	}
	return
}

// go build -o taplib.so -buildmode=c-shared taplib.go
func main() {}

// fix: defer ctx.Close()
func getShuttle() (ret *frame.Shuttle, err error) {
	if savedCtx != nil {
		ret = savedCtx
	} else {
		var inFile string
		if home, e := os.UserHomeDir(); e == nil {
			inFile = filepath.Join(home, "Documents", "Tapestry", "build", "play.db")
		}
		if db, e := tables.CreateRunTime(inFile); e != nil {
			err = e
		} else if c, e := makeShuttle(db, qna.NewOptions()); e != nil {
			err = e
		} else {
			savedCtx = c
			ret = savedCtx
		}
	}
	return
}

func makeShuttle(db *sql.DB, opts qna.Options) (ret *frame.Shuttle, err error) {
	decoder := decode.NewDecoder(tapestry.AllSignatures)
	if grammar, e := player.MakeGrammar(db); e != nil {
		err = e
	} else if q, e := qdb.NewQueries(db); e != nil {
		err = e
	} else {
		run := qna.NewRuntimeOptions(q, decoder, opts)
		survey := play.MakeDefaultSurveyor(run)
		pt := play.NewPlaytime(run, survey, grammar)
		play.CaptureInput(pt)
		ret = frame.NewShuttle(pt, decoder)
	}
	return
}

var savedCtx *frame.Shuttle
var lastResult unsafe.Pointer
