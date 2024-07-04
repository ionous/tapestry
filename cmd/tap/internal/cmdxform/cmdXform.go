package cmdxform

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"git.sr.ht/~ionous/tapestry/cmd/tap/internal/base"
	"git.sr.ht/~ionous/tapestry/support/inflect"
	"github.com/reiver/go-porterstemmer"
)

func run(ctx context.Context, cmd *base.Command, args []string) (err error) {
	if len(args) == 0 {
		err = fmt.Errorf("%w expected at least one word to transform", base.UsageError)
	} else if !hasAnyFlags() {
		err = fmt.Errorf("%w expected at least one transform option", base.UsageError)
	} else {
		var fns = []func(word string) string{
			porterstemmer.StemString,
			inflect.Pluralize,
			inflect.Singularize,
			inflect.Normalize,
		}
		//
		var x [1]struct{} // static assert that the lengths are the same
		var _ = x[numOptions-len(fns)]

		for _, word := range args {
			for i, ok := range optionFlags {
				if ok {
					opt := options(i).String()
					str := fns[i](word)
					fmt.Println(opt+":", str)
				}
			}
		}
	}
	return
}

var CmdXform = &base.Command{
	Run:       run,
	UsageLine: buildUsage(),
	Flag:      buildFlags(),
	Short:     "transform english words",
	Long:      `Transform English words in various helpful ways.`,
}

//go:generate stringer -type=options
type options int

const (
	stem options = iota
	plural
	singular
	normal
	numOptions = iota
)

// set if
var optionFlags [numOptions]bool

// text written to the user for "tap help xform"
func buildUsage() (ret string) {
	var str strings.Builder
	str.WriteString("tap xform")
	for i := 0; i < numOptions; i++ {
		opt := options(i).String()
		str.WriteString(" [-")
		str.WriteString(opt)
		str.WriteRune(']')
	}
	str.WriteString(` "word"`)
	return str.String()
}

func buildFlags() (ret flag.FlagSet) {
	for i := 0; i < numOptions; i++ {
		opt := options(i).String()
		ret.BoolVar(&optionFlags[i], opt, false, fmt.Sprintf("Print %s of a word", opt))
	}
	return
}

func hasAnyFlags() (okay bool) {
	for _, ok := range optionFlags {
		if ok {
			okay = true
			break
		}
	}
	return
}
