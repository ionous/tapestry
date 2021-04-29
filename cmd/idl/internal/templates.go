package internal

import (
	"io"
	"strings"
	"text/template"
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

var Proto = struct {
	Ext                  string
	Header, Slots, Slats Entry
}{
	".proto",
	Entry{name: "header", content: headerProto},
	Entry{name: "slots", content: slotsProto},
	Entry{name: "slats", content: slatsProto},
}

var Cap = struct {
	Ext       string
	Pack, All Entry
}{
	".capnp",
	Entry{name: "pack", content: packCap},
	Entry{name: "all", content: allCap},
}

var funcMap = template.FuncMap{
	// The name "title" is what the function will be called in the template text.
	"title": strings.Title,
}
