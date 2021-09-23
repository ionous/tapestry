'use strict';
module.exports = `
func {{Pascal name}}_Override_Optional_Marshal(n jsn.Marshaler, val *{{OverrideOf name}}) {
  var zero {{OverrideOf name}}
  if *val != zero {
    {{Pascal name}}_Override_Marshal(n, val)
  }
}

func {{Pascal name}}_Override_Marshal(n jsn.Marshaler, val *{{OverrideOf name}}) {
{{#unless (IsBool name)}}
  {{Pascal name}}_Marshal(n, &{{Pascal name}}{*val})
{{else}}
  var out string
  if *val {
    out = Bool_True
  } else {
    out = Bool_False
  }
  {{Pascal name}}_Marshal(n, &{{Pascal name}}{out})
{{/unless}}
}
`;
