package internal

import (
	"text/template"
)

// Create a new template and parse the letter into it.
var SlatsCap = template.Must(template.New("slatsCap").Parse(slatsCap))

//
const slatsCap = `
{{ range $m := . }}

struct {{$m.Name}} $X.label("{{$m.Camel}}"){{- if $m.Group
}} $X.group("{{$m.Group}}") {{ end 
}}{{ if $m.Desc 
}} $X.desc("{{$m.Desc}}") {{ end 
}} {
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
