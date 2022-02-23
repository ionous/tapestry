package core

import (
	"os"

	"git.sr.ht/~ionous/tapestry/rt/safe"
)

func ExamplePrintNum() {
	var run sayTester
	run.SetWriter(os.Stdout)
	if e := safe.WriteText(&run, &PrintNum{Num: F(213)}); e != nil {
		panic(e)
	}
	// Output:
	// 213
}

func ExamplePrintNumWord() {
	var run sayTester
	run.SetWriter(os.Stdout)
	if e := safe.WriteText(&run, &PrintNumWord{Num: F(213)}); e != nil {
		panic(e)
	}
	// Output:
	// two hundred thirteen
}
