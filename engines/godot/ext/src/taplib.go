package main

// #include <stdlib.h>
import "C"

import (
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/support/shuttle"
)

//export Post
func Post(msg string) (ret *C.char) {
	res, e := post(msg)
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

func post(msg string) (ret string, err error) {
	var buf strings.Builder
	if ctx, e := getContext(); e != nil {
		err = e
	} else if strOrMap, e := shuttle.Decode(strings.NewReader(msg)); e != nil {
		err = e
	} else {
		println("got str or map")
		if next, e := shuttle.Post(&buf, ctx, state, strOrMap); e != nil {
			err = e
		} else {
			ret = buf.String()
			if len(next.Name) > 0 {
				state = next
			}
		}
	}
	return
}

// go build -o taplib.so -buildmode=c-shared taplib.go
func main() {}

func getContext() (ret shuttle.Context, err error) {
	if setup {
		ret = savedCtx
	} else {
		var inFile string
		if home, e := os.UserHomeDir(); e == nil {
			inFile = filepath.Join(home, "Documents", "Tapestry", "build", "play.db")
		}
		opts := qna.NewOptions()
		shuttle.NewContext(inFile, opts)
		// defer ctx.Close()
		if n, e := shuttle.NewContext(inFile, opts); e != nil {
			err = e
		} else {
			setup = true
			savedCtx = n
			ret = n
		}
	}
	return
}

var savedCtx shuttle.Context
var setup bool
var state shuttle.State
var lastResult unsafe.Pointer
