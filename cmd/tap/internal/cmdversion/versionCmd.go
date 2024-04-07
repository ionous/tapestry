package cmdversion

import (
	"context"
	"flag"
	"fmt"
	"log"
	"runtime/debug"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
)

func GetVersion(details bool) (ret string) {
	if b, ok := debug.ReadBuildInfo(); !ok {
		ret = "tap version: ???" // only should happen if not built with modules
	} else {
		if details {
			ret = b.String()
		} else {
			if m := b.Main.Version; len(m) > 0 {
				ret = m
			} else {
				ret = "tap version: unknown" //provisionally
				for _, d := range b.Deps {
					if d != nil && d.Path == "git.sr.ht/~ionous/tapestry" {
						ret = d.Version
						break
					}
				}
			}
			ret = fmt.Sprintf("tap version: %s (%s)", ret, b.GoVersion)
		}
	}
	return
}

func run(ctx context.Context, cmd *base.Command, args []string) (err error) {
	log.Println(GetVersion(cfg.all))
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
