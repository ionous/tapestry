package doc

import (
	"embed"
	"go/doc/comment"
	"html/template"
	"log"
	"strings"

	"git.sr.ht/~ionous/tapestry/lang/compact"
	"git.sr.ht/~ionous/tapestry/support/inflect"
)

//go:embed templates/*
var temFS embed.FS

func docTemplates() (*template.Template, error) {
	funcMap := template.FuncMap{
		"Pascal":     inflect.Pascal,
		"Capitalize": inflect.Capitalize,
		"Title":      inflect.Titlecase,
		"Camel":      inflect.Camelize,
		"Lines":      extractLines,
		"GoComment": func(str []string) template.HTML {
			// use go's document parser to parse the document header
			// ( since its also used for the package comment, might as well.
			text := strings.Join(str, "\n")
			var p comment.Parser
			doc := p.Parse(text)
			var pr comment.Printer
			return template.HTML(pr.HTML(doc))
		},
		"TypeLink": TypeLink,
	}
	return template.New("").Funcs(funcMap).ParseFS(temFS, "templates/*.tem")
}

func normalizeLines(str any) (ret []string) {
	switch str := str.(type) {
	case string:
		ret = []string{str}
	case []any:
		if a, ok := compact.SliceStrings(str); !ok {
			log.Panicf("could parse string slice  %#v", str)
		} else {
			ret = a
		}
	case nil:
		// ?
	default:
		log.Panicf("unknown string type %T", str)
	}
	return
}

func extractLines(str any) (ret string) {
	switch str := str.(type) {
	case string:
		ret = str
	case []any:
		ret, _ = compact.JoinLines(str)
	case nil:
		// ?
	default:
		log.Panicf("unknown string type %T", str)
	}
	return
}
