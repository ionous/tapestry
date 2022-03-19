// Generate the golang dl from .ifspec(s)
package main

import (
	"log"
	"os"

	gomake "git.sr.ht/~ionous/tapestry/cmd/gomake/internal"
	"git.sr.ht/~ionous/tapestry/idl"
)

func main() {
	if e := gomake.WriteSpecs(os.Stdout, idl.Specs); e != nil {
		log.Fatal(e)
	}
}
