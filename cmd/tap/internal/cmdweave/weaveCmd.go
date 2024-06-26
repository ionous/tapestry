// Builds a story database from story files.
package cmdweave

import (
	"context"
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
	var paths []fs.FS
	if pathExists(weaveFlags.shared) {
		paths = append(paths, os.DirFS(weaveFlags.shared))
		log.Println("reading:", weaveFlags.shared)
	} else {
		paths = append(paths, content.Shared) // use the embedded stdlib if the shared folder doesn't exist
		log.Println("using shared static content")
	}
	if pathExists(weaveFlags.stories) {
		log.Println("using local story content from", weaveFlags.stories)
		paths = append(paths, os.DirFS(weaveFlags.stories))
	}

	if outFile, e := filepath.Abs(weaveFlags.outFile); e != nil {
		err = e
	} else {
		log.Println("writing:", outFile)
		if e := WeavePaths(outFile, paths...); e != nil {
			err = e
		} else if weaveFlags.checkAll || len(weaveFlags.checkOne) > 0 {
			opt := qna.NewOptions()
			if cnt, e := cmdcheck.CheckFile(outFile, weaveFlags.checkOne, opt); e != nil {
				err = e
			} else {
				log.Println("Checked", cnt, outFile)
			}
		}
	}
	return
}

func pathExists(str string) bool {
	stat, e := os.Stat(str)
	return e == nil && stat.IsDir()
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
