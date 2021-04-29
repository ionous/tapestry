package internal

//
const allCap = ` @0x838375eaedd19910;
using Go = import "/go.capnp";
using  X = import "options.capnp";
{{- range .Deps }}
using {{.|title}} = import "{{.}}/{{.}}.capnp";
{{- end }}

$Go.package("dl");
$Go.import("git.sr.ht/~ionous/dl");
{{ range $m := .Slots }}
struct {{$m.Name}}Impl {{ if $m.Desc -}}
	$X.desc("{{$m.Desc}}")
{{- end }} {
{{- if $m.Sigs }}
  union {
{{- range $i, $s := $m.Sigs }}
	{{ printf "%-30s" $s.Camel }} @{{printf "%-3d" $i}} :{{$s.Package|title}}.{{$s.Type}} $X.label("{{$s.Raw}}");
{{- end }}
  }
{{- end }}
}
{{- end -}}
`
