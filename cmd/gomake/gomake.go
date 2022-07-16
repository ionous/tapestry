// Generate the golang dl from .ifspec(s)
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	gomake "git.sr.ht/~ionous/tapestry/cmd/gomake/internal"
	"git.sr.ht/~ionous/tapestry/idl"
	"github.com/ionous/errutil"
	"golang.org/x/tools/imports"
)

func main() {
	var out string
	flag.StringVar(&out, "out", "", "output directory")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	flag.Parse()
	if len(out) == 0 {
		flag.Usage()
	} else if path, e := filepath.Abs(out); e != nil {
		flag.Usage()
		log.Fatal(e)
	} else {
		dlLike := strings.HasSuffix(path, string(filepath.Separator)+"dl")
		if e := gomake.WriteSpecs(idl.Specs, func(groupName string, b []byte) (err error) {
			path := path
			if groupName == "rt" && dlLike {
				path = path[:len(path)-2]
			}
			base := filepath.Join(path, groupName)
			os.MkdirAll(base, 0700)
			path = filepath.Join(base, groupName+"_lang.go")
			// uses goimports to add imports and format the source
			if formatted, e := imports.Process(path, b, nil); e != nil {
				fmt.Println(string(b))
				err = errutil.New(e, "while formatting", groupName)
			} else {
				b = formatted
			}
			if fp, e := os.Create(path); e != nil {
				err = e // writing errors take precedence over the formatting error
			} else {
				println("writing", path)
				fp.Write(b)
				fp.Close()
			}
			return
		}); e != nil {
			log.Fatal(e)
		}
	}
}

const Description = //
`gomake generates go language serialization code for reading and writing .if files.
Currently, the compiler embeds the .ifspecs descriptions into the gomake executable.
`
const Example = "go run gomake.go -out ../../dl"

func init() {
	flag.Usage = func() {
		println(Description)
		flag.PrintDefaults()
		println("\nex.", Example)
	}
}
