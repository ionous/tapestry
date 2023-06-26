package testweave

import (
	"log"
	"testing"

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
func (w *Warnings) Catch(t *testing.T) func() {
	was := LogWarning
	LogWarning = func(e error) {
		(*w) = append((*w), e)
	}
	return func() {
		if len(*w) > 0 {
			t.Fatal("unhandled warnings", *w)
		}
		LogWarning = was
	}
}

// return the warnings as a raw list, clear all stored errors.
func (w *Warnings) All() (ret []error) {
	ret, (*w) = (*w), nil
	return ret
}

// remove and return the first warning, or error if there are none left.
func (w *Warnings) Shift() (err error) {
	if cnt := len(*w); cnt == 0 {
		err = errutil.New("out of warnings")
	} else {
		err, (*w) = (*w)[0], (*w)[1:]
	}
	return
}
