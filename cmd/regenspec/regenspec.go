// Package main for 'regenspec'.
// from the big json Type, generate .if style commands which are capable of generating that Type.
package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/template"

	regen "git.sr.ht/~ionous/tapestry/cmd/regenspec/internal"
	"git.sr.ht/~ionous/tapestry/web/js"
)

// fix: generate one file at a time, output one file at a time; see compact.go for example.
func main() {
	if ts, e := template.ParseFS(templates, "templates/*.tmpl"); e != nil {
		log.Fatal("parsing failed: ", e)
	} else if e := convert(ts, allBytes); e != nil {
		log.Fatal(e)
	}
}

// generated manually via makeops
// fix: in an ideal world this be using the original typespec format not the munged json output by makeops
// the munged output is a lot simpler. this is a step to maybe getting that right.
// ( though, ultimately the authority of the format should probably be written in the format,
//   so maybe the original typespec format was only useful as a stepping stone. )

//go:embed data/allTypes.jspec
var allBytes []byte

//go:embed templates/*
var templates embed.FS

func convert(ts *template.Template, b []byte) (err error) {
	var allTypes js.MapSlice // use MapSlice to keep the file specified order
	if e := json.Unmarshal(b, &allTypes); e != nil {
		err = e
	} else {
		// read all the types into their groups
		// note: the descriptions of the groups themselves start with an underscore
		groups := make(map[string][]*regen.Type)
		for _, n := range allTypes {
			var m map[string]interface{}
			if e := json.Unmarshal(n.Msg, &m); e != nil {
				err = e
				break
			} else {
				t := regen.NewType(m)
				// get the name of the group from the type
				var gn string
				if u := t.Uses(); u == "group" {
					gn = n.Key
					if gn[0] != '_' {
						panic(t.Name())
					}
				} else {
					// the first group is the primary group.
					gn = t.AllGroups()[0]
					if len(gn) == 0 {
						panic(t.Name())
					}
				}
				// add to the group
				groups[gn] = append(groups[gn], t)
			}
		}
		for group, types := range groups {
			if group[0] == '_' {
				continue // skip groups
			}
			// go-lang templating doesnt allow us to control the whitespace of nested templates
			// it will always look bad, so just buffer it....
			log.Println("writing", group)
			gt := groups["_"+group][0]

			var buf bytes.Buffer
			if e := ts.Execute(&buf, map[string]interface{}{
				"Type":  gt,
				"Group": group,
				"Types": types,
			}); e != nil {
				err = e
				break
			}
			// then reformat in json.... (requires a bytes.Buffer dst)
			b := buf.Bytes()
			var indented bytes.Buffer
			if e := json.Indent(&indented, b, "", "  "); e != nil {
				log.Println("couldnt format", group)
				indented.Reset()
				indented.Write(b)
			}
			// finally write it
			outFile := fmt.Sprintf("../../idl/%s.ifspecs", group)
			if fp, e := os.Create(outFile); e != nil {
				err = e
				break
			} else {
				indented.WriteTo(fp)
				fp.Close()
			}
		}
	}
	return
}
