package core

import (
	"fmt"

	"git.sr.ht/~ionous/iffy/rt/safe"
	"git.sr.ht/~ionous/iffy/rt/writer"
	"git.sr.ht/~ionous/iffy/test/testutil"
)

var helloThereWorld = MakeActivity(
	&Say{&Text{"hello"}},
	&Say{&Text{"there"}},
	&Say{&Text{"world"}},
)

func ExampleSpan() {
	var run sayTester
	run.SetWriter(writer.NewStdout())
	if e := safe.WriteText(&run, &Span{helloThereWorld}); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// hello there world
}

func ExampleBracket() {
	var run sayTester
	run.SetWriter(writer.NewStdout())
	if e := safe.WriteText(&run, &Bracket{helloThereWorld}); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// (hello there world)
}

func ExampleSlash() {
	var run sayTester
	run.SetWriter(writer.NewStdout())
	if e := safe.WriteText(&run, &Slash{helloThereWorld}); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// hello / there / world
}

func ExampleCommas() {
	var run sayTester
	run.SetWriter(writer.NewStdout())
	if e := safe.WriteText(&run, &Commas{helloThereWorld}); e != nil {
		fmt.Println("Error:", e)
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
