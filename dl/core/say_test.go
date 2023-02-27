package core

import (
	"os"

	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/rt/writer"
	"git.sr.ht/~ionous/tapestry/test/testutil"
)

var helloThereWorld = MakeActivity(
	&PrintText{Text: T("hello")},
	&PrintText{Text: T("there")},
	&PrintText{Text: T("world")},
)

func ExampleSpan() {
	var run sayTester
	run.SetWriter(os.Stdout)
	if e := safe.WriteText(&run, &SpanText{Does: helloThereWorld}); e != nil {
		panic(e)
	}
	// Output:
	// hello there world
}

func ExampleBracket() {
	var run sayTester
	run.SetWriter(os.Stdout)
	if e := safe.WriteText(&run, &BracketText{Does: helloThereWorld}); e != nil {
		panic(e)
	}
	// Output:
	// (hello there world)
}

func ExampleSlash() {
	var run sayTester
	run.SetWriter(os.Stdout)
	if e := safe.WriteText(&run, &SlashText{Does: helloThereWorld}); e != nil {
		panic(e)
	}
	// Output:
	// hello / there / world
}

func ExampleCommas() {
	var run sayTester
	run.SetWriter(os.Stdout)
	if e := safe.WriteText(&run, &CommaText{Does: helloThereWorld}); e != nil {
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
