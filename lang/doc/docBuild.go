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

const ext = ".html"
const SourceUrl = "https://pkg.go.dev/git.sr.ht/~ionous/tapestry@v0.24.4-1"
const typesFolder = "type"
const slotsFolder = "slot"

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
					m = append(m, Command{Spec: spec, Flow: t})
				}
			}
			sortCommands(m)
			splitBySlot(slots, types)

			// generate idl files:
			os.Mkdir(filepath.Join(outDir, typesFolder), os.ModePerm)
			outFile := filepath.Join(outDir, typesFolder, types.Name+ext)
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
		os.Mkdir(filepath.Join(outDir, slotsFolder), os.ModePerm)
		for slot, part := range slots {
			outFile := filepath.Join(outDir, slotsFolder, slot.Name+ext)
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

func Create(outFile string, tem *template.Template, data any) (err error) {
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
				// fix:
				pkgName, pkgPath := types.Name, "dl/"+types.Name
				m.addFlow(t, pkgName, pkgPath)
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
func (cat *Commands) addFlow(t *typeinfo.Flow, pkgName, pkgPath string) {
	spec := BuildSpec(t)
	c := Command{t, pkgName, pkgPath, spec}
	cat.Commands = append(cat.Commands, c)
}

type Command struct {
	*typeinfo.Flow
	pkgName, pkgPath string
	Spec             string
}

func (c Command) SourceLink() string {
	// fix: need the package name
	name := inflect.Pascal(c.Name)
	return path.Join(SourceUrl, c.pkgPath+"#"+name)
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
	pascal := inflect.Pascal(name)
	switch t.(type) {
	default:
		log.Panicf("unknown type %T", t)
	case *typeinfo.Slot:
		ret = join("../"+slotsFolder+"/", name, ext)

	case *typeinfo.Flow:
		var idl string
	Flow:
		for _, t := range allTypes {
			for _, el := range t.Flow {
				if el.Name == name {
					idl = t.Name
					break Flow
				}
			}
		}
		if len(idl) == 0 {
			log.Panicf("couldnt find flow for %s", name)
		}
		ret = join("../"+typesFolder+"/", idl, ext, "#", pascal)

	case *typeinfo.Str:
		var idl string
	Str:
		for _, t := range allTypes {
			for _, el := range t.Str {
				if el.Name == name {
					idl = t.Name
					break Str
				}
			}
		}
		if len(idl) == 0 {
			log.Panicf("couldnt find str for %s", name)
		}
		ret = join("../"+typesFolder+"/", idl, ext, "#", pascal)

	case *typeinfo.Num:
		var idl string
	Num:
		for _, t := range allTypes {
			for _, el := range t.Num {
				if el.Name == name {
					idl = t.Name
					break Num
				}
			}
		}
		if len(idl) == 0 {
			log.Panicf("couldnt find num for %q", name)
		}
		ret = join("../"+typesFolder+"/", idl, ext, "#", pascal)

	}
	return
}

func join(str ...string) string {
	return strings.Join(str, "")
}
