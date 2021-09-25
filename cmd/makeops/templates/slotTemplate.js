// slotTemplate.js
'use strict';
module.exports = `
{{~#with type~}}
const {{Pascal name}}_Type = "{{name}}";
{{#if ../marshal}}
var {{Pascal name}}_Optional_Marshal = {{Pascal name}}_Marshal{{Custom name}}

func {{Pascal name}}_Marshal{{Custom name}}(n jsn.Marshaler, ptr *{{Pascal name}}) {
  if slat := *ptr; slat != nil {
    slat.(jsn.Marshalee).Marshal(n)
  }
  return
}
{{/if}}
{{/with}}
{{>repeat name=(Pascal type.name) el=(Pascal type.name)}}
`;
