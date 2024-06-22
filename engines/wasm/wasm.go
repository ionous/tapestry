// Experimental wasm version
//
// use `npm run wasm-build` to build
// use `npm run wasm-serve` to serve
//
// https://go.dev/wiki/WebAssembly
// https://codeberg.org/meta/gzipped
// can strip symbols with -ldflags='-s -w'
// can compress and .gz files will be served by http-server
package main

import (
	"fmt"
	"strings"
	"syscall/js"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/frame"
	"git.sr.ht/~ionous/tapestry/parser"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/decode"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/support/play"
)

func main() {
	fmt.Println("hello")
	registerCallbacks(PlayGame())
	// stop from exiting?
	// might make more sense to go routine and wait on a process channel
	c := make(chan struct{})
	<-c
	fmt.Println("goodbye")
}

// register a js global object called "tapestry" containing everything needed.
// https://tutorialedge.net/golang/go-webassembly-tutorial/#registering-functions
func registerCallbacks(c *frame.Shuttle) {
	// expects a named endpoint, and data for that endpoint.
	// data is a string containing serialized json; returns in the same format.
	post := func(this js.Value, i []js.Value) (ret any, err error) {
		if cnt := len(i); cnt != 2 {
			err = fmt.Errorf("received %d args", cnt)
		} else if endpoint, e := getString(i[0]); e != nil {
			err = e
		} else if msg, e := getString(i[1]); e != nil {
			err = e
		} else {
			println("processing", endpoint, len(i))
			var buf strings.Builder
			if e := c.Post(&buf, endpoint, []byte(msg)); e != nil {
				err = e
			} else {
				ret = buf.String()
			}
		}
		return
	}
	// register the tapestry object
	js.Global().Set("tapestry", map[string]any{
		"post": promisify(post),
	})
}

type jsErrFunc func(js.Value, []js.Value) (any, error)

// https://stackoverflow.com/questions/67437284/how-to-throw-js-error-from-go-web-assembly
func promisify(call jsErrFunc) js.Func {
	jsError := js.Global().Get("Error")
	jsPromise := js.Global().Get("Promise")
	// this is the initiating call from javascript:
	// the js event loop is paused during handling.
	return js.FuncOf(func(those js.Value, args []js.Value) any {
		// if we were to re-enter javascript or if we were long running
		// this would have to be inside the promise and a go-routine
		// https://pkg.go.dev/syscall/js#FuncOf
		res, e := call(those, args)

		// we return a function of two callbacks that javascript can call
		return jsPromise.New(js.FuncOf(func(this js.Value, resolveReject []js.Value) any {
			resolve, reject := resolveReject[0], resolveReject[1]
			if e != nil {
				err := jsError.New(e.Error())
				reject.Invoke(err)
			} else {
				resolve.Invoke(js.ValueOf(res))
			}
			return nil // we return nothing: resolve/reject got our data
		}))
	})
}

func getString(data js.Value) (ret string, err error) {
	if data.Type() != js.TypeString {
		err = fmt.Errorf("not a string")
	} else {
		ret = data.String()
	}
	return
}

func PlayGame() *frame.Shuttle {
	return PlayWithOptions(qna.NewOptions())
}

// aka. makeShuttle
func PlayWithOptions(opts qna.Options) (ret *frame.Shuttle) {
	var grammar parser.AlwaysError
	var q query.QueryNone
	//
	decoder := decode.NewDecoder(tapestry.AllSignatures)
	run := qna.NewRuntimeOptions(q, decoder, opts)
	survey := play.MakeDefaultSurveyor(run)
	pt := play.NewPlaytime(run, survey, grammar)
	play.CaptureInput(pt)
	return frame.NewShuttle(pt, decoder)
}
