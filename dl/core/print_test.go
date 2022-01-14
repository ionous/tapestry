package core

import (
	"git.sr.ht/~ionous/tapestry/rt/safe"
	"git.sr.ht/~ionous/tapestry/rt/writer"
)

func ExamplePrintNum() {
	var run sayTester
	run.SetWriter(writer.NewStdout())
	if e := safe.WriteText(&run, &PrintNum{Num: F(213)}); e != nil {
		panic(e)
	}
	// Output:
	// 213
}

func ExamplePrintNumWord() {
	var run sayTester
	run.SetWriter(writer.NewStdout())
	if e := safe.WriteText(&run, &PrintNumWord{Num: F(213)}); e != nil {
		panic(e)
	}
	// Output:
	// two hundred thirteen
}
