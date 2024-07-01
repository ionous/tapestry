//go:build js && wasm

// Experimental wasm version
// http://localhost:5173/
//
// use `npm run build-wasm` to build
// use `npm run dev` to serve
//
// https://go.dev/wiki/WebAssembly
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"
	"syscall/js"

	"git.sr.ht/~ionous/tapestry"
	"git.sr.ht/~ionous/tapestry/dl/frame"
	"git.sr.ht/~ionous/tapestry/parser"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/qna/query"
	"git.sr.ht/~ionous/tapestry/qna/raw"
	"git.sr.ht/~ionous/tapestry/support/play"
)

func main() {
	fmt.Println("hello")
	tapestry.Register(gob.Register)
	registerCallbacks()
	// stop from exiting?
	// might make more sense to go routine and wait on a process channel
	c := make(chan struct{})
	<-c
	fmt.Println("goodbye")
}

// register a js global object called "tapestry" containing everything needed.
// https://tutorialedge.net/golang/go-webassembly-tutorial/#registering-functions
func registerCallbacks() {
	var shuttle *frame.Shuttle
	var tapestry js.Value
	// expects a named endpoint, and data for that endpoint.
	// data is a string containing serialized json; returns in the same format.
	post := func(this js.Value, i []js.Value) (ret any, err error) {
		if shuttle == nil {
			err = fmt.Errorf("game hasn't been loaded, or has finished.")
		} else if cnt := len(i); cnt != 2 {
			err = fmt.Errorf("post received %d args", cnt)
		} else if endpoint, e := getString(i[0]); e != nil {
			err = e
		} else if msg, e := getString(i[1]); e != nil {
			err = e
		} else {
			println("processing", endpoint, len(i))
			var buf strings.Builder
			if e := shuttle.Post(&buf, endpoint, []byte(msg)); e == nil {
				ret = buf.String()
			} else {
				err = e
				// shuttle = nil?
			}
		}
		return
	}
	// expects binary data
	// data is a string containing serialized json; returns in the same format.
	play := func(this js.Value, i []js.Value) (ret any, err error) {
		if cnt := len(i); cnt != 1 {
			err = fmt.Errorf("play received %d args", cnt)
		} else {
			v := i[0]
			b := make([]byte, v.Length())
			js.CopyBytesToGo(b, v)
			if storyData, e := LoadStory(b); e != nil {
				err = e
			} else {
				name := storyData.Scenes[len(storyData.Scenes)-1]
				fmt.Printf("loaded %s...\n", name)
				tapestry.Set("story", name)
				shuttle = PlayGame(storyData)
			}
		}
		return
	}

	// register the tapestry object
	tapestry = js.ValueOf(map[string]any{
		"post": promisify(post),
		"play": promisify(play),
		// the story name is a little loopy
		// its needed for some of the calls (ex. restart)
		// but it comes from the gob
		"story": nil,
	})
	js.Global().Set("tapestry", tapestry)
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

func getString(v js.Value) (ret string, err error) {
	if v.Type() != js.TypeString {
		err = fmt.Errorf("not a string")
	} else {
		ret = v.String()
	}
	return
}

func PlayGame(storyData raw.Data) *frame.Shuttle {
	return PlayWithOptions(storyData, qna.NewOptions())
}

// aka. makeShuttle
func PlayWithOptions(storyData raw.Data, opts qna.Options) (ret *frame.Shuttle) {
	q := raw.MakeQuery(&storyData)
	scan := make([]parser.Scanner, len(q.Grammar))
	for i, d := range q.Grammar {
		scan[i] = d.MakeScanners()
	}
	//
	run := qna.NewRuntimeOptions(q, opts)
	survey := play.MakeDefaultSurveyor(run) // fix: this should live in play.NewPlaytime; maybe an options?
	pt := play.NewPlaytime(run, survey, scan)
	play.CaptureInput(pt) // fix: can shuttle capture input?
	decoder := query.NewDecoder(tapestry.AllSignatures)
	return frame.NewShuttle(pt, decoder)
}

// read the binary data
func LoadStory(b []byte) (ret raw.Data, err error) {
	dec := gob.NewDecoder(bytes.NewReader(b))
	err = dec.Decode(&ret)
	return
}
