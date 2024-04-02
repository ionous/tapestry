package cmdxform

import (
	"context"
	"flag"
	"fmt"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"github.com/reiver/go-porterstemmer"
)

func run(ctx context.Context, cmd *base.Command, args []string) (_ error) {
	for _, word := range args {
		if cfg.stem {
			fmt.Println("stem:", porterstemmer.StemString(word))
		}
		if cfg.plural {
			fmt.Println("plural:", inflect.Pluralize(word))
		}
		if cfg.single {
			fmt.Println("singular:", inflect.Singularize(word))
		}
	}

	return
}

var CmdXform = &base.Command{
	Run:       run,
	UsageLine: "tap xform [-stem] [-plural] [-single] word",
	Flag:      buildFlags(),
	Short:     "Transform English words",
	Long:      `Transform English words.`,
}

var cfg = struct {
	stem, plural, single bool
}{}

func buildFlags() (fs flag.FlagSet) {
	fs.BoolVar(&cfg.stem, "stem", false, "Print suffix of a word")
	fs.BoolVar(&cfg.plural, "plural", false, "Print plural of a word")
	fs.BoolVar(&cfg.single, "single", false, "Print singular of a word")
	return
}
