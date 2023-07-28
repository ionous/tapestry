// Package useraction: handle long running os level actions initiated from the browser.
package useraction

import (
	"log"
	"net/http"
	"time"

	"git.sr.ht/~ionous/tapestry/web"
)

// uaRequest to start a new named action.
// Launched from a go routine, assumed to be blocking so that there's only dispatched action running at a time.
// For now: no parameters, and only a single string return ( ex. a filename )
type Dispatch func(action string) (string, error)

// returns a Dispatcher that can be used by multiple http handlers if need be.
// rus a go routine internally
func Start(dispatch Dispatch) *Dispatcher {
	results, requests := make(chan Pending), make(chan uaRequest)

	// dispatch function:
	go func() {
		for {
			if req, ok := <-requests; !ok || !req.token.Valid() {
				log.Println("dispatch closed")
				close(results)
				break // closed
			} else {
				v, e := dispatch(req.action)
				results <- Pending{
					tok: req.token,
					val: v,
					err: e,
				}
			}
		}
	}()

	return &Dispatcher{
		reader:   MakeReader(results),
		requests: requests,
	}
}

// approximate amount of time an http request should wait to fulfill an action.
// note: chrome and firefox are both 300 seconds ( 5 minutes )
const Timeout = 120 * time.Second

// max time between an initial request and the attempt to get its results.
// needed so that unrequested results dont take up memory; shorter seems better.
const Expiry = 1 * time.Second

type Dispatcher struct {
	reader   Reader
	requests chan uaRequest
}

// only one of token or value gets set;
// indicated by whether Token.Valid() is true
type Result struct {
	Token Token
	Value string
}

type uaRequest struct {
	action string
	token  Token
}

func makeRequest(action string) uaRequest {
	return uaRequest{
		action: action,
		token:  MakeToken(),
	}
}

// this will block if a dispatch callback is in progress
// not much we can do about that here i dont think
// func (m *Dispatcher) Close() {
// 	m.requests <- uaRequest{}
// }

// given a token returned by post, see if the result is ready.
// err might be an http Status, or some error from the dispatcher.
// StatusuaRequestTimeout means the caller should try again.
func (m *Dispatcher) Get(token Token) (ret string, err error) {
	return m.reader.Read(token)
}

// request that the dispatch function specified in Start() process the requested action as soon as it can;
// returns a valid Token if the action is still pending or in progress after some timeout.
// err might be an http Status, or some error from the dispatcher.
func (m *Dispatcher) Post(action string) (ret Result, err error) {
	start := time.Now()
Hack:
	if token, e := m.post(makeRequest(action)); e != nil {
		firstTime := err == nil
		err = e
		if firstTime && e == web.Status(http.StatusTooManyRequests) {
			// when debugging, the dispatch can get stuck trying to write results
			// because all of the actual readers are dead
			m.reader.readSince(start, 0, Token{})
			goto Hack
		}
	} else {
		// tries to read on post, but if it times out that's okay.
		// returns the value on success, and the token on "still pending".
		if res, e := m.reader.readSince(start, Timeout, token); e == nil {
			ret.Value = res
		} else if e == web.Status(http.StatusRequestTimeout) {
			ret.Token = token
		} else {
			err = e
		}
	}
	return
}

// attempt to post to the
func (w *Dispatcher) post(req uaRequest) (ret Token, err error) {
	select {
	case w.requests <- req:
		ret = req.token

	default:
		// to handle closing properly, fail immediately if a request ( another post )
		// is already in progress. tbd: would it be worthwhile to have a timeout instead?
		// ( ex. so that two browser windows could *try* to open a file dialog,
		// and the second wouldnt wait till the first completes )
		err = web.Status(http.StatusTooManyRequests)
	}
	return
}
