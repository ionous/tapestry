// Builds a story database from story files.
package cmdweave

import (
	"context"
	"errors"
	"flag"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/cmdcheck"
	"git.sr.ht/~ionous/tapestry/content"
	"git.sr.ht/~ionous/tapestry/qna"
)

func runWeave(ctx context.Context, cmd *base.Command, args []string) (err error) {
	log.Println("reading:", weaveFlags.stories)
	log.Println("writing:", weaveFlags.outFile)

	stories := os.DirFS(weaveFlags.stories)
	shared := os.DirFS(weaveFlags.shared) // use the embedded stdlib if the shared folder doesn't exist
	if _, e := fs.Stat(shared, "."); errors.Is(e, fs.ErrNotExist) {
		shared = content.Shared
	}
	if outFile, e := filepath.Abs(weaveFlags.outFile); e != nil {
		err = e
	} else if e := WeavePaths(outFile, shared, stories); e != nil {
		err = e
	} else if weaveFlags.checkAll || len(weaveFlags.checkOne) > 0 {
		opt := qna.NewOptions()
		if cnt, e := cmdcheck.CheckFile(outFile, weaveFlags.checkOne, opt); e != nil {
			err = e
		} else {
			log.Println("Checked", cnt, outFile)
		}
	}
	return
}

var CmdWeave = &base.Command{
	Run:       runWeave,
	Flag:      buildFlags(),
	UsageLine: "tap weave [-in path] [-out path]",
	Short:     "compile a story",
	Long: `Turns story files into produces a playable database.

Using '-check' or '-run=<name>' can run all unit tests, or a specific one.
The weave command provides an option to locate the shared Tapestry libraries. 
If that location doesn't exist, weave uses a set of built-in libraries.
`,
}

// collection of local flags
var weaveFlags = struct {
	stories, shared, outFile string
	checkAll                 bool
	checkOne                 string
}{}

func buildFlags() (flags flag.FlagSet) {
	var shared string
	var stories string
	var outPath string
	if home, e := os.UserHomeDir(); e == nil {
		base := filepath.Join(home, "Documents", "Tapestry")
		shared = filepath.Join(base, "shared")
		stories = filepath.Join(base, "stories")
		outPath = filepath.Join(base, "build", "play.db")
	}
	flags.StringVar(&weaveFlags.stories, "in", stories, `input file or directory name.`)
	flags.StringVar(&weaveFlags.shared, "lib", shared, `folder containing the standard library.`)
	flags.StringVar(&weaveFlags.outFile, "out", outPath, "optional output filename (sqlite3)")
	flags.BoolVar(&weaveFlags.checkAll, "check", false, "run check after importing?")
	flags.StringVar(&weaveFlags.checkOne, "run", "", "run check on a particular test")
	return
}
