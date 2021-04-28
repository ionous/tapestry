package internal

import (
	"text/template"
)

// Create a new template and parse the letter into it.
var SlotsCap = template.Must(template.New("slotsCap").Parse(slotsCap))

//
const slotsCap = `
struct Pos {
	source @0 :Text; 
	offset @1 :Text;
}

{{- range $m := . }}

struct {{$m.Name}} { eval @0:AnyPointer; }
struct {{$m.Name}}Impl {{ if $m.Desc -}}
	$X.desc("{{$m.Desc}}")
{{- end }} {
{{- if $m.Sigs }}
  union {
{{- range $i, $s := $m.Sigs }}
	{{ printf "%-30s" $s.Camel }} @{{printf "%-3d" $i}} :{{.Type}} $X.label("{{$s.Raw}}");
{{- end }}
  }
{{- end }}
}
{{- end -}}
`
