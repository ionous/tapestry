// Package doc builds web friendly documentation of the idl.
package doc

import (
	"cmp"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/lang/typeinfo"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

const SourceUrl = "https://pkg.go.dev/git.sr.ht/~ionous/tapestry@v0.24.4-1"
const typesFolder = "idl"
const slotFolder = "slot"
const baseUrl = "/api"

//go:embed static/*
var staticFS embed.FS

// hack
var allTypes []typeinfo.TypeSet

// outDir can point to a temp directory if need be
func Build(outDir string, idl []typeinfo.TypeSet) (err error) {
	allTypes = idl
	dc := docComments{allTypes}
	if tem, e := docTemplates(&dc); e != nil {
		err = e
	} else if e := os.MkdirAll(outDir, os.ModePerm); e != nil {
		err = e
	} else if e := CopyStaticFiles(outDir); e != nil {
		err = e
	} else {
		slots := make(SlotMap)
		for _, types := range idl {
			var m []Command
			for _, t := range types.Flow {
				if _, hidden := t.Markup["internal"]; !hidden {
					spec := BuildSpec(t)
					m = append(m, Command{types.Name, t, spec})
				}
			}
			sortCommands(m)
			splitBySlot(slots, types)

			// generate idl files:
			outFile := filepath.Join(outDir, typesFolder, types.Name)
			if e := Create(outFile, tem, map[string]any{
				"Name":     types.Name,
				"Types":    types,
				"Commands": m,
			}); e != nil {
				err = e
				return // early out
			}
		}
		// generate slot docs:
		for slot, part := range slots {
			outFile := filepath.Join(outDir, slotFolder, slot.Name)
			if cmt, e := compact.ExtractComment(slot.Markup); e != nil {
				err = e
				break
			} else {
				sortCommands(part.Commands)
				if e := Create(outFile, tem, map[string]any{
					"Name":     slot.Name,
					"Slot":     slot,
					"Commands": part.Commands,
					"Comment":  cmt,
				}); e != nil {
					err = e
					break
				}
			}
		}
	}
	return
}

func CopyStaticFiles(outDir string) (err error) {
	for _, fileName := range []string{"style.css", "custom.css"} {
		if b, e := fs.ReadFile(staticFS, "static/"+fileName); e != nil {
			err = e
		} else if fp, e := os.Create(filepath.Join(outDir, fileName)); e != nil {
			err = e
		} else {
			_, err = fp.Write(b)
		}
	}
	return
}

func Create(outPath string, tem *template.Template, data any) (err error) {
	os.MkdirAll(outPath, os.ModePerm)
	outFile := filepath.Join(outPath, "index.html")
	if fp, e := os.Create(outFile); e != nil {
		err = e
	} else if e := tem.ExecuteTemplate(fp, "page.tem", data); e != nil {
		err = fmt.Errorf("%w writing %s", e, outFile)
	}
	return
}

// - split everything into slots
func splitBySlot(out SlotMap, types typeinfo.TypeSet) {
	for _, t := range types.Flow {
		if _, hidden := t.Markup["internal"]; !hidden {
			for _, slot := range t.Slots {
				// create the command map
				m, ok := out[slot]
				if !ok {
					m = new(Commands)
					out[slot] = m
				}
				m.addFlow(types.Name, t)
			}
		}
	}
}

type SlotMap map[*typeinfo.Slot]*Commands

type Commands struct {
	Commands []Command
}

func sortCommands(cs []Command) {
	slices.SortFunc(cs, func(a, b Command) int {
		// we want specs ending in colons to be listed before ones with a new word.
		// ex. Numeral: before Numeral words:
		// colon is less than underscore ( but greater than space )
		x, y := strings.Replace(a.Spec, " ", "_", -1), strings.Replace(b.Spec, " ", "_", -1)
		return cmp.Compare(x, y)
	})
}
func (cat *Commands) addFlow(idlName string, t *typeinfo.Flow) {
	spec := BuildSpec(t)
	c := Command{idlName, t, spec}
	cat.Commands = append(cat.Commands, c)
}

type Command struct {
	idlName string // ex. "story"
	*typeinfo.Flow
	Spec string
}

// link to go package documentation
func (c Command) SourceLink() string {
	name := inflect.Pascal(c.Name)
	pkgPath := "dl/" + c.idlName
	return path.Join(SourceUrl, pkgPath+"#"+name)
}

func (c Command) Terms() (ret []typeinfo.Term) {
	// filter out hidden terms;
	for i, t := range c.Flow.Terms {
		if t.Private && ret == nil {
			ret = c.Flow.Terms[:i]
		}
		if !t.Private && ret != nil {
			ret = append(ret, t)
		}
	}
	if ret == nil {
		ret = c.Flow.Terms
	}
	return
}

// Build the document style signature for this flow
// it's different than the actual signature because,
// among other things, it includes markers for optional elements.
func BuildSpec(t *typeinfo.Flow) string {
	var str strings.Builder
	str.WriteString(inflect.Pascal(t.Lede))
	if len(t.Terms) == 0 {
		str.WriteString(":")
	} else {
		for i, t := range t.Terms {
			if t.Optional {
				str.WriteRune('[')
			}
			if len(t.Label) > 0 {
				if i == 0 {
					str.WriteRune(' ')
				}
				str.WriteString(inflect.Camelize(t.Label))
			}
			str.WriteRune(':')
			if t.Optional {
				str.WriteRune(']')
			}
		}
	}
	return str.String()
}

func TypeLink(t typeinfo.T) (ret string) {
	name := t.TypeName()
	switch t.(type) {
	case *typeinfo.Slot:
		ret = linkToSlot(name)
	default:
		if a, ok := findTypeSet(t); !ok {
			log.Panicf("couldnt find flow for %s", name)
		} else {
			ret = linkToType(a.Name, name)
		}
	}
	return
}

func findTypeSet(t typeinfo.T) (ret typeinfo.TypeSet, okay bool) {
Out:
	for _, types := range allTypes {
		switch t := t.(type) {
		case *typeinfo.Flow:
			if slices.Contains(types.Flow, t) {
				ret, okay = types, true
				break Out
			}
		case *typeinfo.Str:
			if slices.Contains(types.Str, t) {
				ret, okay = types, true
				break Out
			}
		case *typeinfo.Num:
			if slices.Contains(types.Num, t) {
				ret, okay = types, true
				break Out
			}
		case *typeinfo.Slot:
			if slices.Contains(types.Slot, t) {
				ret, okay = types, true
				break Out
			}
		}
	}
	return
}

// idl name is a .tells file without the extension.
func linkToType(idlName, typeName string) string {
	out := path.Join(baseUrl, typesFolder, idlName)
	if len(typeName) > 0 {
		out += "#" + inflect.Pascal(typeName)
	}
	return out
}

// slot name is something like "story_statement"
func linkToSlot(slotName string) string {
	return path.Join(baseUrl, slotFolder, slotName)
}
