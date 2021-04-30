package core

import (
	capnp "zombiezen.com/go/capnproto2"
)

type stringer interface{ String() string }

func cmdError(s stringer, cmd capnp.Struct, e error) error {
	return cmdErrorCtx(s, cmd, "", e)
}

func cmdErrorCtx(s stringer, cmd capnp.Struct, ctx string, e error) error {
	// avoid triggering errutil panics for break statements
	// if _, ok := err.(DoInterrupt); !ok {
	// e := &composer.CommandError{Cmd: op, Ctx: ctx}
	// err := errutil.Append(err, e)
	// }
	return e
}
