package charm

import (
	"github.com/ionous/errutil"
)

const Eof = rune(-1)

// Parse sends each rune of string to the passed state chart,
// and returns the last state which ran.

func Parse(str string, first State) (err error) {
	_, err = innerParse(first, str)
	return
}

func innerParse(first State, str string) (ret State, err error) {
	try := first
	for i, r := range str {
		if next := try.NewRune(r); next == nil {
			// no states left to parse remaining input
			err = EndpointError{str, i, try}
			break
		} else if es, ok := next.(Terminal); ok {
			err = errutil.Append(es.err, EndpointError{str, i, try})
			break
		} else {
			try = next
		}
	}
	if err == nil {
		ret = try
	}
	return
}

// ParseEof sends each rune of string to the passed state chart;
// after its done with the string, it sends an eof(-1) to flush any remaining data.
// see also Parse() which does not send the eof.
func ParseEof(first State, str string) (err error) {
	if last, e := innerParse(first, str); e != nil {
		err = e
	} else if last != nil {
		if fini := last.NewRune(Eof); fini != nil {
			if es, ok := fini.(Terminal); ok && es.err != nil {
				err = errutil.Fmt("%s handling eof for %q", es.err, str)
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
	rest := e.str[e.end:]
	return errutil.Sprintf("%q ended at %q, index %d, in %q",
		StateName(e.last), rest, e.end, e.str)
}
