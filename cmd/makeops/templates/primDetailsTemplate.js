// primDetailsTemplate.js
// note: partial templates seem to COPY arguments,
// so they lose .. relative hierarchy.
'use strict';
module.exports = `
func {{Pascal name}}_Optional_Marshal(n jsn.Marshaler, val *{{Pascal name}}) {
  var zero {{Pascal name}}
  if val.{{Pascal uses}} != zero.{{Pascal uses}} {
     {{Pascal name}}_Marshal(n, val)
  }
}

func {{Pascal name}}_Marshal(n jsn.Marshaler, val *{{Pascal name}}) {
{{~#if (IsPositioned this)}}
  n.SetCursor(val.At.Offset)
{{/if}}
  n.WriteValue({{Pascal name}}_Type, val.{{Pascal uses}})
}
`;
