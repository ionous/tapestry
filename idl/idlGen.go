// +build ignore
package dl

//go:generate capnp compile -I ../../../../zombiezen.com/go/capnproto2/std -ogo:core core.capnp
//go:generate capnp compile -I ../../../../zombiezen.com/go/capnproto2/std -ogo:debug debug.capnp
//go:generate capnp compile -I ../../../../zombiezen.com/go/capnproto2/std -ogo:grammar grammar.capnp
//go:generate capnp compile -I ../../../../zombiezen.com/go/capnproto2/std -ogo:list list.capnp
//go:generate capnp compile -I ../../../../zombiezen.com/go/capnproto2/std -ogo:reader reader.capnp
//go:generate capnp compile -I ../../../../zombiezen.com/go/capnproto2/std -ogo:rel rel.capnp
//go:generate capnp compile -I ../../../../zombiezen.com/go/capnproto2/std -ogo:render render.capnp
//go:generate capnp compile -I ../../../../zombiezen.com/go/capnproto2/std -ogo:rtx rtx.capnp
//go:generate capnp compile -I ../../../../zombiezen.com/go/capnproto2/std -ogo:all allCmds.capnp
