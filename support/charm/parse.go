package charm

import (
	"github.com/ionous/errutil"
)

const eof = rune(-1)

// Parse is the main function of chart.
func Parse(first State, str string) (err error) {
	try := first
	for i, r := range str {
		if next := try.NewRune(r); next == nil {
			// no states left to parse remaining input
			err = EndpointError{str, i, try}
			break
		} else if n, ok := next.(errorState); ok {
			err = n.err
			break
		} else {
			try = next
		}
	}

	// after parsing the whole string
	// send an eof to flush any remaining data
	// ( ex. parsing a list might not know the list has finished )
	// fix: is this really needed?
	// states have "Get()" so shouldnt they be able to finish there?
	if err == nil && try != nil {
		if fini := try.NewRune(eof); fini != nil {
			if es, ok := fini.(errorState); ok {
				err = es.err
			} else {
				// and if we are passing eof, shouldnt the states check for it and return nil?
				// err = EndpointError{str, len(str), fini}
			}
		}
	}
	return
}

// ended before the whole input was parsed.
type EndpointError struct {
	str  string
	end  int
	last State
}

// index of the failure point in the input
func (e EndpointError) End() int {
	return e.end
}

func (e EndpointError) Error() (ret string) {
	return errutil.Sprintf("parsing `%s` ended in %T(%s) at index %d %q",
		e.str, e.last, e.last.StateName(), e.end, e.str[e.end:])
}
