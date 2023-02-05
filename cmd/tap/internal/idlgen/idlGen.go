package idlgen

import (
	"context"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
)

var CmdIdl = &base.Command{
	Run:       runCodeGen,
	UsageLine: "tap idlgen [-dl] [-out]",
	Short:     "generate golang serializers from .ifspecs",
	Long: `
Idlb generates an sqlite database containing the tapestry command language.
Language extensions must be built into tapestry itself
( they require an implementation that tapestry can execute )
so there are no options for this command except the output filename.
	`,
}

// FIX: where do keyword specs come from?
func runCodeGen(ctx context.Context, cmd *base.Command, args []string) (err error) {
}
