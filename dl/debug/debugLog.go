package debug

import (
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/dl/cmd"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// LogLevel controls how much debugging to print
// logs only at the named level and higher.
var LogLevel LoggingLevel = C_LoggingLevel_Debug

func (op *LogValue) Execute(run rt.Runtime) (err error) {
	// fix? currently, weave can't guarantee a lack of side-effects;
	// so this always evals even if it doesn't print.
	if v, e := safe.GetAssignment(run, op.Value); e != nil {
		err = cmd.Error(op, e)
	} else {
		printLog(op.LogLevel, Stringify(v))
	}
	return
}

func (op *Note) Execute(run rt.Runtime) (err error) {
	// fix? currently, weave can't guarantee a lack of side-effects;
	// so this always evals even if it doesn't print.
	if v, e := safe.GetText(run, op.Text); e != nil {
		err = cmd.Error(op, e)
	} else {
		printLog(op.LogLevel, v.String())
	}
	return
}

func printLog(level LoggingLevel, out string) {
	global := LogLevel
	if (global >= 0 && level >= global) || (global < 0 && level != 0) {
		if level < 0 {
			level = 0
		}
		txt := level.String()
		header := strings.Repeat("#", 1+int(level))
		log.Println(header, txt, out)
	}
}
