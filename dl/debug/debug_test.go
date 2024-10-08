package debug

import (
	"log"
	"strings"
	"testing"

	"git.sr.ht/~ionous/tapestry/dl/call"
	"git.sr.ht/~ionous/tapestry/dl/literal"
)

func TestLog(t *testing.T) {
	w := log.Writer()
	defer log.SetOutput(w)
	var b strings.Builder
	log.SetOutput(&b)
	//
	lo := LogValue{
		LogLevel: C_LoggingLevel_Error,
		Value:    &call.FromText{Value: literal.T("hello")},
	}
	if e := lo.Execute(nil); e != nil {
		t.Fatal(e)
	} else if got := b.String(); !strings.HasSuffix(got, " #### error \"hello\"\n") {
		t.Fatalf("got %q", got)
	}
}
