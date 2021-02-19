package pattern_test

import (
	"fmt"

	"git.sr.ht/~ionous/iffy/ephemera/debug"
	"git.sr.ht/~ionous/iffy/test/testutil"
)

// ExampleSayMe converts numbers to text
// http://learnyouahaskell.com/syntax-in-functions
func ExampleSayMe() {
	// rules are run in reverse order.
	var kinds testutil.Kinds
	kinds.AddKinds((*debug.SayMe)(nil))
	run := testutil.Runtime{
		Kinds: &kinds,
		PatternMap: testutil.PatternMap{
			"say_me": &debug.SayPattern,
		}}
	// say 4 numbers
	for i := 1; i <= 4; i++ {
		fmt.Printf(`say_me %d = "`, i)
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
