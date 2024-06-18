package format_test

import (
	"os"

	"git.sr.ht/~ionous/tapestry/dl/format"
	"git.sr.ht/~ionous/tapestry/dl/literal"
	"git.sr.ht/~ionous/tapestry/rt"
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/rt/writer"
	"git.sr.ht/~ionous/tapestry/test/testutil"
)

var helloThereWorld = []rt.Execute{
	&format.PrintText{Text: literal.T("hello")},
	&format.PrintText{Text: literal.T("there")},
	&format.PrintText{Text: literal.T("world")},
}

func ExamplePrintWords() {
	var run sayTester
	run.SetWriter(os.Stdout)
	if e := safe.WriteText(&run, &format.PrintWords{Exe: helloThereWorld}); e != nil {
		panic(e)
	}
	// Output:
	// hello there world
}

func ExamplePrintBrackets() {
	var run sayTester
	run.SetWriter(os.Stdout)
	if e := safe.WriteText(&run, &format.PrintBrackets{Exe: helloThereWorld}); e != nil {
		panic(e)
	}
	// Output:
	// (hello there world)
}

func ExampleSlash() {
	var run sayTester
	run.SetWriter(os.Stdout)
	if e := safe.WriteText(&run,
		&format.PrintWords{
			Separator: literal.T(" / "),
			Exe:       helloThereWorld,
		}); e != nil {
		panic(e)
	}
	// Output:
	// hello / there / world
}

func ExamplePrintCommas() {
	var run sayTester
	run.SetWriter(os.Stdout)
	if e := safe.WriteText(&run, &format.PrintCommas{Exe: helloThereWorld}); e != nil {
		panic(e)
	}
	// Output:
	// hello, there, and world
}

type baseRuntime struct {
	testutil.PanicRuntime
}
type sayTester struct {
	baseRuntime
	writer.Sink
}
