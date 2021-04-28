package internal

import (
	"io"
	"text/template"
)

type Entry struct {
	temp *template.Template
}

func (x *Entry) Must(w io.Writer, data interface{}) {
	if x.temp != nil {
		template.Must(x.temp, x.temp.Execute(w, data))
	}
}

var Proto = struct {
	Ext                  string
	Header, Slots, Slats Entry
}{
	".proto",
	Entry{HeaderProto},
	Entry{SlotsProto},
	Entry{SlatsProto},
}

var Cap = struct {
	Ext                  string
	Header, Slots, Slats Entry
}{
	".capnp",
	Entry{HeaderCap},
	Entry{SlotsCap},
	Entry{SlatsCap},
}
