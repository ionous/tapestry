package core

import (
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/rt/writer"
	"git.sr.ht/~ionous/tapestry/test/testutil"
)

var helloThereWorld = MakeActivity(
	&Say{Text: T("hello")},
	&Say{Text: T("there")},
	&Say{Text: T("world")},
)

func ExampleSpan() {
	var run sayTester
	run.SetWriter(writer.NewStdout())
	if e := safe.WriteText(&run, &SpanText{helloThereWorld}); e != nil {
		panic(e)
	}
	// Output:
	// hello there world
}

func ExampleBracket() {
	var run sayTester
	run.SetWriter(writer.NewStdout())
	if e := safe.WriteText(&run, &BracketText{helloThereWorld}); e != nil {
		panic(e)
	}
	// Output:
	// (hello there world)
}

func ExampleSlash() {
	var run sayTester
	run.SetWriter(writer.NewStdout())
	if e := safe.WriteText(&run, &SlashText{helloThereWorld}); e != nil {
		panic(e)
	}
	// Output:
	// hello / there / world
}

func ExampleCommas() {
	var run sayTester
	run.SetWriter(writer.NewStdout())
	if e := safe.WriteText(&run, &CommaText{helloThereWorld}); e != nil {
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
