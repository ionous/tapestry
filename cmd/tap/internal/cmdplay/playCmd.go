package cmdplay

import (
	"context"
	"log"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/qna"
	g "git.sr.ht/~ionous/tapestry/rt/generic"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/support/player"
	"github.com/ionous/errutil"
)

func goPlay(ctx context.Context, cmd *base.Command, args []string) (err error) {
	if lvl, ok := debug.MakeLoggingLevel(cfg.logLevel); !ok {
		err = errutil.New("Unknown log level", cfg.logLevel)
	} else {
		debug.LogLevel = lvl
		opts := qna.NewOptions()
		opts.SetOption(meta.PrintResponseNames, g.BoolOf(cfg.responses))
		if cnt, e := player.PlayWithOptions(cfg.inFile, cfg.testString, cfg.scene, opts); e != nil {
			errutil.PrintErrors(e, func(s string) { log.Println(s) })
			if errutil.Panic {
				log.Panic("mismatched")
			}
		} else {
			log.Println("done", cnt, cfg.inFile)
		}
	}
	return
}

var CmdPlay = &base.Command{
	Run:       goPlay,
	Flag:      buildFlags(),
	UsageLine: "tap play [-in dbpath] [-scene name]",
	Short:     "play a story",
	Long: `Run a scene within a previously built story database.

Using '-test' can run the list of specified commands as if a player had typed them one by one.
`,
}
