package core

import (
	"fmt"

	"git.sr.ht/~ionous/iffy/rt/safe"
	"git.sr.ht/~ionous/iffy/rt/writer"
)

func ExamplePrintNum() {
	var run sayTester
	run.SetWriter(writer.NewStdout())
	if e := safe.WriteText(&run, &PrintNum{&Number{213}}); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// 213
}

func ExamplePrintNumWord() {
	var run sayTester
	run.SetWriter(writer.NewStdout())
	if e := safe.WriteText(&run, &PrintNumWord{&Number{213}}); e != nil {
		fmt.Println("Error:", e)
	}
	// Output:
	// two hundred thirteen
}
