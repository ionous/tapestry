package text

import (
	"os"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/rt/writer"
	"git.sr.ht/~ionous/tapestry/test/testutil"
)

func ExamplePrintNumDigits() {
	var run sayTester
	run.SetWriter(os.Stdout)
	if e := safe.WriteText(&run, &PrintNumDigits{Num: literal.F(213)}); e != nil {
		panic(e)
	}
	// Output:
	// 213
}

func ExamplePrintNumWords() {
	var run sayTester
	run.SetWriter(os.Stdout)
	if e := safe.WriteText(&run, &PrintNumWords{Num: literal.F(213)}); e != nil {
		panic(e)
	}
	// Output:
	// two hundred thirteen
}

type baseRuntime struct {
	testutil.PanicRuntime
}
type sayTester struct {
	baseRuntime
	writer.Sink
}
