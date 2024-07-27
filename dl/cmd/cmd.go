package cmd

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

type cmdError struct {
	Cmd typeinfo.Instance
	Ctx string
	Err error
}

func (my cmdError) Unwrap() error {
	return my.Err
}

func (e cmdError) Error() string {
	name := "<nil>"
	if e.Cmd != nil {
		t := e.Cmd.TypeInfo()
		name = t.TypeName()
	}
	var padding string
	if len(e.Ctx) > 0 {
		padding = " "
	}
	return fmt.Sprintf("%s%s%s %s", name, padding, e.Ctx, e.Err)
}

// add error context... if there is an error
func Error(op typeinfo.Instance, e error) (err error) {
	if e != nil {
		err = cmdError{Cmd: op, Err: e}
	}
	return
}

func ErrorCtx(op typeinfo.Instance, ctx string, e error) (err error) {
	if e != nil {
		err = cmdError{Cmd: op, Ctx: ctx, Err: e}
	}
	return
}
