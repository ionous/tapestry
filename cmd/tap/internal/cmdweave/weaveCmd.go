// Transforms detailed format json files to their compact format, and back again.
// Relies on the file extension ".ifx" being used for detailed format files, and ".if" for compact files.
package cmdweave

import (
	"context"
	"log"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/support/asm"
)

func runWeave(ctx context.Context, cmd *base.Command, args []string) (err error) {
	log.Println("reading:", weaveFlags.srcPath)
	log.Println("writing:", weaveFlags.outFile)
	if e := asm.AssembleFolder(weaveFlags.srcPath, weaveFlags.outFile); e != nil {
		err = e
	} else if weaveFlags.checkAll || len(weaveFlags.checkOne) > 0 {
		if cnt, e := asm.CheckOutput(weaveFlags.outFile, weaveFlags.checkOne); e != nil {
			err = e
		} else {
			log.Println("Checked", cnt, weaveFlags.outFile)
		}
	}
	return
}

var CmdWeave = &base.Command{
	Run:       runWeave,
	Flag:      buildFlags(),
	UsageLine: "tap weave [-in path] [-out path]",
	Short:     "make a playable story",
	Long: `Turns story files into produces a playable database.

Using '-check' or '-run=<name>' can run all unit tests, or a specific one.
`,
}
