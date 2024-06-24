package cmddump

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/support/dump"
	"git.sr.ht/~ionous/tapestry/support/files"
	"git.sr.ht/~ionous/tapestry/tables"
)

func runDump(ctx context.Context, cmd *base.Command, args []string) (err error) {
	if inFile, e := filepath.Abs(dumpFlags.inFile); e != nil {
		flag.Usage()
		log.Fatal(e)
	} else if outFile, e := filepath.Abs(dumpFlags.outFile); e != nil {
		flag.Usage()
		log.Fatal(e)
	} else if db, e := tables.OpenModel(inFile); e != nil {
		err = e
	} else {
		defer db.Close()
		if data, e := dump.DumpAll(db, args[0]); e != nil {
			err = e
		} else {
			err = files.SaveJson(outFile, data, true)
		}
	}
	return
}

var CmdDump = &base.Command{
	Run:       runDump,
	Flag:      buildFlags(),
	UsageLine: `tap dump [-in dbpath] [-out file] "name of story"`,
	Short:     "generate raw scene dump",
	Long:      `Write the story db into a non-sql data format.`,
}

// collection of local flags
var dumpFlags = struct {
	inFile, outFile string
}{}

func buildFlags() (fs flag.FlagSet) {
	var inFile string
	var outFile string
	if home, e := os.UserHomeDir(); e == nil {
		inFile = filepath.Join(home, "Documents", "Tapestry", "build", "play.db")
		outFile = filepath.Join(home, "Documents", "Tapestry", "build", "play.json")
	}

	fs.StringVar(&dumpFlags.inFile, "in", inFile, "input file name (sqlite3)")
	fs.StringVar(&dumpFlags.outFile, "out", outFile, "output file name (json)")
	// fs.StringVar(&cfg.testString, "test", "", "optional list of commands to run (non-interactive)")
	// fs.BoolVar(&cfg.json, "json", false, "expect input/output in json (default is plain text)")
	// fs.BoolVar(&cfg.responses, "responses", false, "print response names instead of values")
	// fs.StringVar(&cfg.logLevel, "log", debug.C_LoggingLevel_Info.String(), levels)

	return
}
