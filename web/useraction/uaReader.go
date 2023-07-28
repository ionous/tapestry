package useraction

import (
	"net/http"
	"time"

	"git.sr.ht/~ionous/tapestry/web"
)

// http endpoints get results from readers.
type Reader struct {
	next    chan Pending
	lock    chan struct{}
	results []Pending
}

type Pending struct {
	tok Token
	val string
	err error // not an http error
}

// result channel could be small, probably size 1
// result: make(chan result, 3),
func MakeReader(res chan Pending) Reader {
	return Reader{
		lock: make(chan struct{}, 1),
		next: res,
	}
}

// can be called by multiple http response handlers
// the result's token will match the passed token
func (r *Reader) Read(token Token) (string, error) {
	return r.readSince(time.Now(), Timeout, token)
}

func (r *Reader) readSince(start time.Time, wait time.Duration, token Token) (ret string, err error) {
	// because of the shared result pool, we only want one reader at a time.
	// lock by attempting to write into the reader's internal channel:
	// https://stackoverflow.com/questions/54488284/attempting-to-acquire-a-lock-with-a-deadline-in-golang
	select {
	case r.lock <- struct{}{}:
		if a, ok := r.pull(token); ok {
			ret = a
		} else {
			ret, err = r.poll(start, wait, token)
		}
		<-r.lock // unlock: read what we wrote to allow a new writer.

	case <-time.After(wait):
		err = web.Status(http.StatusRequestTimeout)
	}
	return
}

// search for a dequeued request matching the passed token.
// ( only one reader pull() happens at a time. )
func (r *Reader) pull(token Token) (ret string, okay bool) {
	if token.Valid() {
		now := time.Now()
		out, at := r.results, 0
		for _, el := range r.results {
			// only one will match; but keep looping to clear out expired
			if el.tok == token {
				ret, okay = el.val, true
			} else if t := timeWhen(el.tok.Time); t.Add(Expiry).After(now) {
				out[at] = el // keep this entry; might copy over itself; that's fine.
				at++
			}
		}
		if at < len(out) {
			out = out[:at] // slice down if need be.
		}
	}
	return
}

// wait for a request matching the passed token.
// ( only one reader poll() happens at a time. )
func (r *Reader) poll(start time.Time, wait time.Duration, token Token) (ret string, err error) {
Loop:
	for err == nil {
		elapsed := time.Now().Sub(start)

		select {
		case <-time.After(wait - elapsed):
			err = web.Status(http.StatusRequestTimeout)

		// read from the os action
		case res, okay := <-r.next:
			if !okay {
				// channel closed by dispatcher?
				err = web.Status(http.StatusServiceUnavailable)
			} else if res.tok != token {
				// not our result, so save it:
				// ( could maybe happen if the user had multiple browser windows open )
				r.results = append(r.results, res) // append our new result
				// doesn't break, doesnt error
				// tries again.
			} else {
				if res.err != nil {
					err = res.err
				} else {
					ret = res.val
					break Loop
				}
			}
		}
	}
	return
}
