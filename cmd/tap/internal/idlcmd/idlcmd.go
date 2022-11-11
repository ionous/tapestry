// Copyright (C) 2022 - Simon Travis. All rights reserved.
// Use of this source code is governed by the Hippocratic 2.1
// license that can be found in the LICENSE file.

package idlcmd

import (
	"context"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/idl"
)

var CmdIdl = &base.Command{
	Run:       runGenerate,
	UsageLine: "tap idlb [file.db]", //  [-run regexp] [-n] [-v] [-x] [build flags] [file.go... | packages]
	Short:     "generate a language database",
	Long: `
Idlb generates an sqlite database containing the tapestry command language.
Currently, language extensions must be built into tapestry itself
( because they require an implementation that tapestry can execute )
so there are no options for this command except the output filename.
	`,
}

var specs = idl.Specs

func runGenerate(ctx context.Context, cmd *base.Command, args []string) {

}
