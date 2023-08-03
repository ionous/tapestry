package mdl

import (
	"log"
	"strings"

	"github.com/ionous/errutil"
)

var LogWarning = func(e error) {
	log.Println("Warning:", e) // for now good enough
}

type Warnings []error

// override the global warning function
// returns a defer-able function which:
// 1. restores the warning function; and,
// 2. raises a Fatal error if there are any unhandled warnings.
func (w *Warnings) Catch(fatal func(args ...any)) func() {
	was := LogWarning
	LogWarning = func(e error) {
		(*w) = append((*w), e)
	}
	return func() {
		if len(*w) > 0 {
			fatal("unhandled warnings", *w)
		}
		LogWarning = was
	}
}

// return the warnings as a raw list, clear all stored errors.
func (w *Warnings) All() (ret []error) {
	ret, (*w) = (*w), nil
	return ret
}

// remove and return the first warning,
// error if there is none, or if warning doesnt start with the passed prefix.
func (w *Warnings) Expect(prefix string) (err error) {
	if e := w.pop(); e == nil {
		err = errutil.Fmt("expected %q, received nothing.", prefix)
	} else if str := e.Error(); !strings.HasPrefix(str, prefix) {
		err = errutil.Fmt("expected %q, received %q.", prefix, str)
	}
	return
}

func (w *Warnings) pop() (err error) {
	if cnt := len(*w); cnt > 0 {
		err, (*w) = (*w)[0], (*w)[1:]
	}
	return
}
