package cmdversion

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
)

func GetVersion() (ret string) {
	ret = "unknown"
	if a, ok := debug.ReadBuildInfo(); ok {
		if m := a.Main.Version; len(m) > 0 {
			ret = m
		} else {
			for _, d := range a.Deps {
				if d != nil && d.Path == "git.sr.ht/~ionous/tapestry" {
					ret = d.Version
					break
				}
			}
		}
		ret = fmt.Sprintf("%s (%s)", ret, a.GoVersion)
	}
	return
}

func run(ctx context.Context, cmd *base.Command, args []string) (err error) {
	log.Println("tap version", GetVersion())
	return
}

var CmdVersion = &base.Command{
	Run:       run,
	UsageLine: `tap version`,
	Short:     "print Tapestry version",
	Long:      `Version prints the information about the tap tool recorded at the time it was built.`,
}
