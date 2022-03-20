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

		if e := gomake.WriteSpecs(idl.Specs, func(groupName string, b []byte) {
			path := path
			if groupName == "rt" && dlLike {
				path = path[:len(path)-2]
			}
			path = filepath.Join(path, groupName, groupName+"_lang.go_")
			if fp, e := os.Create(path); e != nil {
				log.Fatal(e)
			} else {
				fp.Write(b)
				fp.Close()
			}
		}); e != nil {
			log.Fatal(e)
		}
	}
}
