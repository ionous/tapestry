// Builds a story database from story files.
package cmdweave

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/cmdcheck"
	"git.sr.ht/~ionous/tapestry/content"
	"git.sr.ht/~ionous/tapestry/qna"
)

func runWeave(ctx context.Context, cmd *base.Command, args []string) (err error) {
	var paths []NamedFS
	if pathExists(cfg.shared) {
		fsys := os.DirFS(cfg.shared)
		paths = append(paths, NamedFS{"shared", fsys})
		log.Println("reading:", cfg.shared)
	} else {
		fsys := content.Shared
		paths = append(paths, NamedFS{"embedded share", fsys}) // use the embedded stdlib if the shared folder doesn't exist
		log.Println("using shared static content")
	}
	if pathExists(cfg.stories) {
		log.Println("using local story content from", cfg.stories)
		fsys := os.DirFS(cfg.stories)
		paths = append(paths, NamedFS{"stories", fsys})
	}

	if outFile, e := filepath.Abs(cfg.outFile); e != nil {
		err = e
	} else {
		log.Println("writing:", outFile)
		if e := WeavePaths(outFile, paths...); e != nil {
			err = e
		} else if cfg.checkAll || len(cfg.checkOne) > 0 {
			opt := qna.NewOptions()
			if cnt, e := cmdcheck.CheckFile(outFile, cfg.checkOne, opt); e != nil {
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

// filled with the user's choices as described by buildFlags()
var cfg = struct {
	stories, shared, outFile string
	checkAll                 bool
	checkOne                 string
}{}

func buildFlags() (ret flag.FlagSet) {
	var shared string
	var stories string
	var outPath string
	if home, e := os.UserHomeDir(); e == nil {
		base := filepath.Join(home, "Documents", "Tapestry")
		shared = filepath.Join(base, "shared")
		stories = filepath.Join(base, "stories")
		outPath = filepath.Join(base, "build", "play.db")
	}
	ret.StringVar(&cfg.stories, "in", stories, `input file or directory name.`)
	ret.StringVar(&cfg.shared, "lib", shared, `folder containing the standard library.`)
	ret.StringVar(&cfg.outFile, "out", outPath, "optional output filename (sqlite3)")
	ret.BoolVar(&cfg.checkAll, "check", false, "run check after importing?")
	ret.StringVar(&cfg.checkOne, "run", "", "run check on a particular test")
	return
}
