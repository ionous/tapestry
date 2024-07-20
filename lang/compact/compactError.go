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
		ret = fmt.Sprintf("error at %s: %s", e.src, e.err)
	} else if ofs := e.src.Line; ofs > 0 {
		ret = fmt.Sprintf("error on line %d: %s", ofs, e.err)
	} else {
		ret = e.err.Error()
	}
	return
}

func MessageError(msg Message, err error) SourceError {
	return SourceError{MakeSource(msg.Markup), err}
}
