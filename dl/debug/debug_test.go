package debug

import (
	"log"
	"strings"
	"testing"

	"git.sr.ht/~ionous/iffy/dl/core"
)

func TestLog(t *testing.T) {
	w := log.Writer()
	defer log.SetOutput(w)
	var b strings.Builder
	log.SetOutput(&b)
	//
	lo := DebugLog{
		LogLevel: LoggingLevel{Str: LoggingLevel_Error},
		Value:    &core.FromText{&core.TextValue{"hello"}}}
	if e := lo.Execute(nil); e != nil {
		t.Fatal(e)
	} else if got := b.String(); !strings.HasSuffix(got, " ###### error hello\n") {
		t.Fatalf("got %q", got)
	}
}
