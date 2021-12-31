package testpat_test

import (
	"fmt"

	"git.sr.ht/~ionous/tapestry/test/debug"
	"git.sr.ht/~ionous/tapestry/test/testpat"
	"git.sr.ht/~ionous/tapestry/test/testutil"
)

// ExampleSayMe converts numbers to text.
// note: this relies almost entirely on debug functionality,
// so it's just an example not a test of runtime behavior.
// see also: http://learnyouahaskell.com/syntax-in-functions
func ExampleSayMe() {
	// rules are run in reverse order.
	var kinds testutil.Kinds
	kinds.AddKinds((*debug.SayMe)(nil))
	run := testpat.Runtime{
		testpat.Map{
			"say_me": &debug.SayPattern,
		}, testutil.Runtime{
			Kinds: &kinds,
		},
	}
	// say 4 numbers
	for i := 1; i <= 4; i++ {
		fmt.Printf(`say_me %d = "`, i)
		// creates a core.CallPattern and runs it through the debug runtime
		if e := debug.DetermineSay(i).Execute(&run); e != nil {
			fmt.Println("Error:", e)
		}
		fmt.Println(`"`)
	}

	// Output:
	// say_me 1 = "One!"
	// say_me 2 = "Two!"
	// say_me 3 = "Three!"
	// say_me 4 = "Not between 1 and 3."
}
