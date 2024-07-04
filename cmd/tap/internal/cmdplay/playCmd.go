package cmdplay

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/dl/debug"
	"git.sr.ht/~ionous/tapestry/qna"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/meta"
	"git.sr.ht/~ionous/tapestry/support/player"
)

// called by tap.go
func playCmd(ctx context.Context, _ *base.Command, args []string) (err error) {
	if len(args) != 1 {
		err = fmt.Errorf("%w expected a scene name", base.UsageError)
	} else if lvl, ok := debug.MakeLoggingLevel(cfg.logLevel); !ok {
		log.Println("Unknown logging level", cfg.logLevel)
		log.Println("Expected one of:", strings.Join(debug.Zt_LoggingLevel.Options, ", "))
		err = fmt.Errorf("%w expected a valid logging level", base.UsageError)
	} else if saveDir, e := initSaveDir(cfg.saveDir); e != nil {
		// could make this a warning i suppose; then just not save...
		err = fmt.Errorf("couldn't initialize save directory because %s", e)
	} else {
		scene := args[0]
		debug.LogLevel = lvl
		opts := qna.NewOptions()
		opts.SetOption(meta.PrintResponseNames, rt.BoolOf(cfg.responses))
		opts.SetOption(meta.SaveDir, rt.StringOf(saveDir))
		err = player.PlayWithOptions(cfg.inFile, cfg.testString, scene, opts)
	}
	return
}

func initSaveDir(dir string) (ret string, err error) {
	if path, e := filepath.Abs(dir); e != nil {
		err = e
	} else if e := os.MkdirAll(path, os.ModePerm); e != nil {
		err = e
	} else {
		ret = path
	}
	return
}

// description of the play command; used by tap.go
var CmdPlay = &base.Command{
	Run:       playCmd,
	Flag:      buildFlags(),
	UsageLine: `tap play [-in dbpath] [-out savedir] "name of story"`,
	Short:     "play an existing story",
	Long: `Run a scene within a previously built story database.

Using '-test' can run the list of specified commands as if a player had typed them one by one.
`,
}

// filled with the user's choices as described by buildFlags()
var cfg = struct {
	inFile, saveDir, testString string
	json, responses             bool
	logLevel                    string
}{}

// returns a command line parsing object
func buildFlags() (ret flag.FlagSet) {
	var inFile string
	var saveDir string
	if home, e := os.UserHomeDir(); e == nil {
		inFile = filepath.Join(home, "Documents", "Tapestry", "build", "play.db")
		saveDir = filepath.Join(home, "Documents", "Tapestry", "save")
	}

	levels := strings.Join(debug.Zt_LoggingLevel.Options, ", ")
	ret.StringVar(&cfg.inFile, "in", inFile, "input file name (sqlite3)")
	ret.StringVar(&cfg.saveDir, "out", saveDir, " directory for save files")
	ret.StringVar(&cfg.testString, "test", "", "optional list of commands to run (non-interactive)")
	ret.BoolVar(&cfg.json, "json", false, "expect input/output in json (default is plain text)")
	ret.BoolVar(&cfg.responses, "responses", false, "print response names instead of values")
	ret.StringVar(&cfg.logLevel, "log", debug.C_LoggingLevel_Info.String(), levels)
	return
}
