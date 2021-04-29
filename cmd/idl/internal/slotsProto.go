package internal

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
