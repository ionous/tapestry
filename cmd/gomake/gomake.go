// Generate the golang dl from .ifspec(s)
package main

import (
	"flag"
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
	flag.StringVar(&out, "out", "", "optional output directory")
	flag.BoolVar(&errutil.Panic, "panic", false, "panic on error?")
	if len(out) == 0 {
		out = "../../dl"
	}
	if path, e := filepath.Abs(out); e != nil {
		log.Fatal(e)
	} else {
		dlLike := strings.HasSuffix(path, string(filepath.Separator)+"dl")
		if e := gomake.WriteSpecs(idl.Specs, func(groupName string, b []byte) (err error) {
			path := path
			if groupName == "rt" && dlLike {
				path = path[:len(path)-2]
			}
			path = filepath.Join(path, groupName, groupName+"_lang.go")
			// use https://pkg.go.dev/golang.org/x/tools/cmd/goimports to add imports and format the source
			if formatted, e := imports.Process(path, b, nil); e != nil {
				err = errutil.New(e, "while formatting", groupName)
			} else {
				b = formatted
			}
			if fp, e := os.Create(path /*temp*/ + "_"); e != nil {
				err = e // writing errors take precedence over the formatting error
			} else {
				fp.Write(b)
				fp.Close()
			}
			return
		}); e != nil {
			log.Fatal(e)
		}
	}
}
