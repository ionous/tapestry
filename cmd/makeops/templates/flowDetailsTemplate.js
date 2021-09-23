// flowDetailsTemplate.js
'use strict';
module.exports = `
func {{Pascal name}}_Marshal(n jsn.Marshaler, val *{{Pascal name}}) {
{{#if (IsPositioned this)}}
  n.SetCursor(val.At.Offset)
{{/if}}
  n.MapValues({{#if (LedeName this)~}}"{{LedeName this}}"{{~else~}}{{Pascal name}}_Type{{~/if~}}, {{Pascal name}}_Type)
{{~#each (ParamsOf this)}}{{#unless (IsInternal label)}}
  n.MapKey("{{SelectorOf @key this @index}}", {{Pascal ../name}}_Field_{{Pascal @key}})
  /* */ {{ScopeOf type}}{{Pascal type}}{{#if (OverrideOf type)}}_Override{{/if}}{{#if repeats}}_Repeats{{else if optional}}_Optional{{/if}}_Marshal(n, &val.{{Pascal @key}})
{{~/unless}}{{/each}}
  n.EndValues()
  return
}
`
