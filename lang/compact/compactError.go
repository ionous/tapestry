package compact

import (
	"errors"
	"fmt"
)

type Unhandled string

func (u Unhandled) Error() string {
	return fmt.Sprintf("unhandled %s", string(u))
}

func IsUnhandled(e error) bool {
	var u Unhandled
	return errors.As(e, &u)
}

type SourceError struct {
	src Source
	err error
}

func (e SourceError) Source() Source {
	return e.src
}

func (e SourceError) Unwrap() error {
	return e.err
}

func (e SourceError) Error() (ret string) {
	if len(e.src.File) > 0 {
		ret = fmt.Sprintf("error at %s: %s", e.src.ErrorString(), e.err)
	} else if ofs := e.src.Line; ofs > 0 {
		ret = fmt.Sprintf("error near line %d: %s", ofs+1, e.err)
	} else {
		ret = e.err.Error()
	}
	return
}

func MessageError(msg Message, e error) error {
	return MakeSourceError(MakeSource(msg.Markup), e)
}

func MakeSourceError(src Source, e error) (err error) {
	// avoid stacking multiple line number errors:
	// keep the inner most
	var prev SourceError
	if errors.As(e, &prev) {
		err = prev
	} else {
		err = SourceError{src, e}
	}
	return
}
