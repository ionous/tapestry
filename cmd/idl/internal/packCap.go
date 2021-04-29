package internal

// Create a new template and parse the letter into it.

const packCap = `@0x{{.Hash}};
using Go = import "/go.capnp";
using  X = import "../options.capnp";
{{- range .Deps }}
using {{.|title}} = import "../{{.}}/{{.}}.capnp";
{{- end }}

$Go.package("{{.Name}}");
$Go.import("git.sr.ht/~ionous/dl/{{.Name}}");
{{ range .Slots }}
struct {{.Name}} { eval @0:AnyPointer; }
{{- end }}
{{ range $m := .Slats }}
struct {{$m.Name}} $X.label("{{$m.Camel|title}}"){{ 
if $m.Group }} $X.group("{{$m.Group}}"){{ end }}{{ 
if $m.Desc }} $X.desc("{{$m.Desc}}"){{ end }} {
{{- range $i, $p := $m.Places }}
  {{printf "%-12s" $p.Camel }} @{{printf "%-2d" $i}} :{{$p.CapType
}}{{if $p.Internal}} $X.internal{{end
}}{{if $p.Optional}} $X.optional{{end
}}{{if $p.CapLabel}} $X.label("{{$p.CapLabel}}"){{end
}}{{if $p.Pool}} $X.pool("{{$p.Pool}}"){{end
}};
{{- end }}
}
{{- end }}
`
