package cmdversion

import (
	"context"
	"flag"
	"log"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/support/files"
)

func run(ctx context.Context, cmd *base.Command, args []string) (err error) {
	version := files.GetVersion(cfg.all)
	log.Println("tap version:", version)
	return
}

var CmdVersion = &base.Command{
	Run:       run,
	Flag:      buildFlags(),
	UsageLine: `tap version`,
	Short:     "print Tapestry version",
	Long:      `Version prints the information about the tap tool recorded at the time it was built.`,
}

// returns b command line parsing object
func buildFlags() (fs flag.FlagSet) {
	fs.BoolVar(&cfg.all, "all", false, "print detailed information about the binary")
	return
}

// filled with the user's choices as described by buildFlags()
var cfg = struct {
	all bool
}{}
