package debug

import (
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
)

// LogLevel controls how much debugging to print
// The default level ( empty string ) means log everything but notes,
// otherwise it logs only at the named level and higher.
var LogLevel LoggingLevel

func (op *DebugLog) Execute(run rt.Runtime) (err error) {
	// fix? at this time we cant guarantee a lack of side-effects
	// so we always eval even if we don't print.
	if v, e := safe.GetAssignment(run, op.Value); e != nil {
		err = CmdError(op, e)
	} else {
		global := LogLevel
		level := op.LogLevel
		if err == nil && ((global >= 0 && level >= global) || (global < 0 && level != 0)) {
			if level < 0 {
				level = 0
			}
			txt := level.String()
			header := strings.Repeat("#", 1+int(level))
			log.Println(header, txt, Stringify(v))
		}
	}
	return
}
