package decode

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/lang/compact"
)

func ValueError(msg string, value any) error {
	return fmt.Errorf("%s %T(%v)", msg, value, value)
}

func ParamError(msg compact.Message, p compact.Param, e error) error {
	return fmt.Errorf("%q(@%s:) %w", msg.Key, p.Label, e)
}
