// Package doc builds web friendly documentation of the idl.
package doc

import (
	"html/template"
	"io"
	"os"
	"path/filepath"

	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
)

const ext = ".html"

// outDir can point to a temp directory if need be
func Build(outDir string, idl []typeinfo.TypeSet) (err error) {
	if e := os.MkdirAll(outDir, os.ModePerm); e != nil {
		err = e
	} else if tmp, e := docTemplates(); e != nil {
		err = e
	} else {

		for _, types := range idl {
			outFile := filepath.Join(outDir, types.Name+ext)
			if fp, e := os.Create(outFile); e != nil {
				err = e
				break
			} else if e := Write(fp, tmp, types); e != nil {
				err = e
				break
			}
		}

	}
	return
}

func Write(w io.Writer, tmp *template.Template, types typeinfo.TypeSet) (err error) {
	// 1. create groups on a first name basis.
	// 2. group by commands underneath that -- maybe with the go name there separating them.
	// - and underneath each ( or run-inline to the command ) a comment, and then the return types.
	// 2. create inline parameters with optional elements in brackets and types
	// in between the colons
	// alt: make it more tell like, and list the components under in a sequence
	// 3. separate section -- probably at the top because its shorter anyway; containing the slots
	// 4. one for the types ( if any ) -- str, num ( maybe lsted first )

	// so, maybe a "slotHeader", "typesHeader" sub templates?
	err = tmp.ExecuteTemplate(w, "page.tem", types)
	return
}

// 1. generate one page per .idl file
// alt: per idl, generate "categories" for the listed slots
// then sort all commands from all the other idls into those categories
// ( so that "story" commands, "rt" commands, etc. aren't siloed into separate locations
// 2. generate an alphabetical index of all commands
