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
	var buf strings.Builder
	if n, e := getShuttle(); e != nil {
		err = e
	} else if e := n.Post(&buf, endpoint, msg); e != nil {
		err = e
	} else {
		ret = buf.String()
	}
	return
}

// go build -o taplib.so -buildmode=c-shared taplib.go
func main() {}

// fix: defer ctx.Close()
func getShuttle() (ret shuttle.Shuttle, err error) {
	if setup {
		ret = savedCtx
	} else {
		var inFile string
		if home, e := os.UserHomeDir(); e == nil {
			inFile = filepath.Join(home, "Documents", "Tapestry", "build", "play.db")
		}
		if n, e := shuttle.NewShuttle(inFile, qna.NewOptions()); e != nil {
			err = e
		} else {
			setup = true
			savedCtx = n
			ret = n
		}
	}
	return
}

var savedCtx shuttle.Shuttle
var setup bool
var lastResult unsafe.Pointer
