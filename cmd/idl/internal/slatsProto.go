package internal

//
const slatsProto = `
{{ range $m := . }}

message {{$m.Name}} {
{{- range $p := $m.Places }}
  {{ $p.ProtoQualifier }} {{ printf $m.Format $p.ProtoType $p.Arg $p.Index}};
{{- end -}}
{{- if $m.Group }}

	option (if).group = "{{$m.Group}}";
{{- end }}
{{- if $m.Desc }}
	option (if).desc = "{{$m.Desc}}";
{{- end }}
}
{{- end }}
`
