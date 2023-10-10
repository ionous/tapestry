package main

// #include <stdlib.h>
import "C"

import (
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"unsafe"

	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/support/shuttle"
	"github.com/ionous/errutil"
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
			err = errutil.Fmt("Recovered %s:\n%s", r, debug.Stack())
		}
	}()
	var buf strings.Builder
	if n, e := getShuttle(); e != nil {
		err = e
	} else if e := n.Post(&buf, endpoint, []byte(msg)); e != nil {
		err = e
	} else {
		ret = buf.String()
	}
	return
}

// go build -o taplib.so -buildmode=c-shared taplib.go
func main() {}

// fix: defer ctx.Close()
func getShuttle() (ret *shuttle.Shuttle, err error) {
	if savedCtx != nil {
		ret = savedCtx
	} else {
		var inFile string
		if home, e := os.UserHomeDir(); e == nil {
			inFile = filepath.Join(home, "Documents", "Tapestry", "build", "play.db")
		}
		if n, e := shuttle.NewShuttle(inFile, qna.NewOptions()); e != nil {
			err = e
		} else {
			savedCtx = &n
			ret = savedCtx
		}
	}
	return
}

var savedCtx *shuttle.Shuttle
var lastResult unsafe.Pointer
