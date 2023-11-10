package charm

import (
	"io"
	"strings"

	"github.com/ionous/errutil"
)

const Eof = rune(-1)

// Parse sends each rune of string to the passed state chart,
// and returns the last state which ran.
func Parse(str string, first State) (err error) {
	_, err = innerParse(first, strings.NewReader(str))
	return
}

func Read(in io.RuneReader, first State) (err error) {
	_, err = innerParse(first, in)
	return
}

func innerParse(first State, in io.RuneReader) (ret State, err error) {
	try := first
	var buf strings.Builder // fix: some sort of sliding window for error would be nice
	for i := 0; ; i++ {
		if r, _, e := in.ReadRune(); e != nil {
			if e != io.EOF {
				err = errutil.Append(e, EndpointError{buf.String(), i, try, e.Error()})
			}
			break
		} else {
			if next := try.NewRune(r); next == nil {
				// no states left to parse remaining input
				err = EndpointError{buf.String(), i, try, "unknown state"}
				break
			} else if es, ok := next.(Terminal); ok {
				err = EndpointError{buf.String(), i, try, es.err.Error()}
				break
			} else {
				try = next
			}
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
	if last, e := innerParse(first, strings.NewReader(str)); e != nil {
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
	str    string
	end    int
	last   State
	reason string
}

// index of the failure point in the input
func (e EndpointError) End() int {
	return e.end
}

func (e EndpointError) Error() (ret string) {
	return errutil.Sprintf("%s %q ended at index %d in %q",
		e.reason, StateName(e.last), e.end, e.str)
}
