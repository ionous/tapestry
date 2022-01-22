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
)

// fix: generate one file at a time, output one file at a time; see compact.go for example.
func main() {
	if ts, e := template.ParseFS(templates, "templates/*.tmpl"); e != nil {
		log.Fatal("parsing failed: ", e)
	} else if e := convert(ts, allBytes); e != nil {
		log.Fatal(e)
	}
}

//go:embed data/allTypes.jspec
var allBytes []byte

//go:embed templates/*
var templates embed.FS

func convert(ts *template.Template, b []byte) (err error) {
	var allTypes map[string]interface{}
	if e := json.Unmarshal(b, &allTypes); e != nil {
		err = e
	} else {
		groups := make(map[string][]*regen.Type)
		for k, _ := range allTypes {
			t := regen.NewType(regen.MapOf(k, allTypes))
			if u := t.Uses(); u != "group" {
				// the first group is the primary group.
				gn := t.AllGroups()[0]
				if len(gn) == 0 {
					panic(t.Name())
				}
				groups[gn] = append(groups[gn], t)
			}
		}
		for group, types := range groups {
			// go-lang templating doesnt allow us to control the whitespace of nested templates
			// it will always look bad, so just buffer it....
			log.Println("writing", group)
			gt := regen.NewType(regen.MapOf("_"+group, allTypes))

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
