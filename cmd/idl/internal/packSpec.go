package internal

// object {
// 	string group;
// 	object {
// 			object {
// 					string flow; # the short spec
// 					# ex. "cmd label%arg#type"
// 					string desc;
// 					union {
// 						array { string };
// 						string;
// 					} slots;
// 			} string name;
// 	} types;
// }

const packSpec = `{
  "slots": { {{- range $i, $m := .Slots -}}{{if gt $i 0 }},{{end}}
    "{{ $m.Name | under }}": "{{$m.Desc|esc}}"
   {{- end}}
  },
  "types": { {{- range $i, $m := .Slats -}}{{if gt $i 0 }},{{end}}
    "{{ $m.Name | under }}": {
      "flow": "{{$m.FlowText}}"{{if $m.Slots}},
      "slots": {{$m.SlotText}}{{end}}{{if $m.Desc}},
      "desc": "{{$m.Desc|esc}}"{{end}}{{if $m.Group}},
      "group": "{{.Group}}"{{end}}
    }
  {{- end}}
  }
}
`
