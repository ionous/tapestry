package internal

import (
	"encoding/json"
	"io"
	"strings"
	"text/template"

	"git.sr.ht/~ionous/iffy/lang"
)

type Entry struct {
	name, content string
	temp          *template.Template
}

func (x *Entry) Must(w io.Writer, data interface{}) {
	if x.temp == nil {
		x.temp = template.Must(template.New(x.name).Funcs(funcMap).Parse(x.content))
	}
	if x.temp != nil {
		template.Must(x.temp, x.temp.Execute(w, data))
	}
}

var Templates = struct {
	Ext  string
	Pack Entry
}{
	".ifspec",
	Entry{name: "packSpec", content: packSpec},
}

var funcMap = template.FuncMap{
	// The name "title" is what the function will be called in the template text.
	"title": strings.Title,
	"under": lang.Underscore,
	"esc":   escape,
}

func escape(i string) string {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	// Trim the beginning and trailing " character
	return string(b[1 : len(b)-1])
}
