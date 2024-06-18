package format

import (
	"os"

	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/rt/writer"
)

func ExamplePrintNum() {
	var run sayTester
	run.SetWriter(os.Stdout)
	if e := safe.WriteText(&run, &PrintNum{Num: literal.F(213)}); e != nil {
		panic(e)
	}
	// Output:
	// 213
}

func ExamplePrintCount() {
	var run sayTester
	run.SetWriter(os.Stdout)
	if e := safe.WriteText(&run, &PrintCount{Num: literal.F(213)}); e != nil {
		panic(e)
	}
	// Output:
	// two hundred thirteen
}

type sayTester struct {
	baseRuntime
	writer.Sink
}
