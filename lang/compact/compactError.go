package compact

import (
	"errors"

	"github.com/ionous/errutil"
)

type Unhandled string

func (u Unhandled) Error() string {
	return errutil.Sprint("unhandled", string(u))
}

func IsUnhandled(e error) bool {
	var u Unhandled
	return errors.As(e, &u)
}
