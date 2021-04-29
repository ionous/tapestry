package internal

//
// const slatsCap = `using Go = import "/go.capnp";
// using  X = import "../options.capnp";
// {{- range .Deps }}
// using {{.|title}} = import "../{{.}}/{{.}}Slots.capnp";
// {{- end }}

// $Go.package("{{.Package}}");
// $Go.import("git.sr.ht/~ionous/dl/{{.Package}}");
// {{ range $m := .Slats }}
// struct {{$m.Name}} $X.label("{{$m.Camel}}"){{ if $m.Desc
// }} $X.desc("{{$m.Desc}}") {{ end }}
// {
// {{- range $i, $p := $m.Places }}
//   {{printf "%-12s" $p.Camel }} @{{printf "%-2d" $i}} :{{$p.CapType
// }}{{if $p.Internal}} $X.internal{{end
// }}{{if $p.Optional}} $X.optional{{end
// }}{{if $p.CapLabel}} $X.label("{{$p.CapLabel}}"){{end
// }}{{if $p.Pool}} $X.pool("{{$p.Pool}}"){{end
// }};
// {{- end }}
// }
// {{- end }}
// `
