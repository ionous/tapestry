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
