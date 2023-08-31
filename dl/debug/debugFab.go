package debug

import (
	"log"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"github.com/ionous/errutil"
)

// used for fabricate to trigger a forward step of the parser
var Stepper func(words string) error

func (op *Fabricate) Execute(run rt.Runtime) (err error) {
	if Stepper == nil {
		err = errutil.New("no parser set for the fabricator")
	} else if words, e := safe.GetText(run, op.Text); e != nil {
		err = e
	} else {
		words := words.String()
		if LogLevel.Str <= LoggingLevel_Debug {
			log.Println("> ", words)
		}
		if len(words) > 0 {
			err = Stepper(words)
		}
	}
	return
}
