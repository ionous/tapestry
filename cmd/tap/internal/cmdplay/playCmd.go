package cmdplay

import (
	"context"
	"log"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/support/play"
	"github.com/ionous/errutil"
)

func goPlay(ctx context.Context, cmd *base.Command, args []string) (err error) {
	if cnt, e := play.PlayGame(playFlags.inFile, playFlags.testString, playFlags.scene, playFlags.json); e != nil {
		errutil.PrintErrors(e, func(s string) { log.Println(s) })
		if errutil.Panic {
			log.Panic("mismatched")
		}
	} else {
		log.Println("done", cnt, playFlags.inFile)
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
