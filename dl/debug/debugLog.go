package debug

import (
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/affine"
	"git.sr.ht/~ionous/tapestry/dl/assign"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/generic"
	"github.com/ionous/errutil"
	"github.com/kr/pretty"
)

// LogLevel controls how much debugging to print
// The default level ( empty string ) means log everything but notes,
// otherwise it logs only at the named level and higher.
var LogLevel LoggingLevel

func (op *DebugLog) Execute(run rt.Runtime) (err error) {
	// fix? at this time we cant guarantee a lack of side-effects
	// so we always eval even if we don't print.
	if v, e := assign.GetValue(run, op.Value); e != nil {
		err = CmdError(op, e)
	} else {
		var i interface{}
		switch a := v.Affinity(); a {
		case affine.Bool:
			i = v.Bool()
		case affine.Number:
			i = v.Float()
		case affine.NumList:
			i = v.Floats()
		case affine.Text:
			i = v.String()
		case affine.TextList:
			i = v.Strings()
		case affine.Record:
			i = pretty.Sprint(generic.RecordToValue(v.Record()))
		case affine.RecordList:
			i = pretty.Sprint(generic.RecordsToValue(v.Records()))
		default:
			e := errutil.New("unknown affinity", a)
			err = CmdError(op, e)
		}
		global := LogLevel.Index()
		level := op.LogLevel.Index()
		if err == nil && ((global >= 0 && level >= global) || (global < 0 && level != 0)) {
			if level < 0 {
				level = 0
			}
			txt := op.LogLevel.Compose().Strings[level]
			header := strings.Repeat("#", 1+level)
			log.Println(header, txt, i)
		}
	}
	return
}

func (lvl LoggingLevel) Index() (ret int) {
	// FIX: who comes up with this stuff?
	if str := lvl.String(); len(str) == 0 {
		ret = -1
	} else {
		spec := lvl.Compose()
		_, ret = spec.IndexOfChoice(str)
	}
	return
}
