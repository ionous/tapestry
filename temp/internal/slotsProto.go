package internal

import (
	"text/template"
)

// Create a new template and parse the letter into it.
var SlotsProto = template.Must(template.New("slotsProto").Parse(slotsProto))

//
const slotsProto = `
{{ range $m := . }}

message {{$m.Name}} {
{{- if $m.Sigs }}
  oneof eval {
{{- range $s := $m.Sigs }}
		{{ printf $m.Format $s.Type $s.Numbered $s.Crc}};
{{- end }}
	}
{{- end -}}
{{- if $m.Desc }}
	option (if).desc = "{{$m.Desc}}";
{{- end }}
}
{{- end }}
`
