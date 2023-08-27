package composer

import (
	"github.com/ionous/errutil"
)

type CommandError struct {
	Cmd Composer
	Ctx string
}

func (e CommandError) Error() string {
	name := SpecName(e.Cmd)
	var padding string
	if len(e.Ctx) > 0 {
		padding = " "
	}
	return errutil.Sprintf("# %s%s%s", name, padding, e.Ctx)
}
