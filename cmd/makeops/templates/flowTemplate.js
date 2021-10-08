// runTemplate.js
'use strict';
module.exports = `
{{~#with type~}}
// {{Pascal name}} {{desc.short}}
type {{Pascal name}} struct {
{{#each (ParamsOf this)}}
  {{Pascal @key}} {{TypeOf this}} \`if:"
  {{~#if (IsInternal label)}}internal{{else}}label={{LabelOf @key this @index}}{{/if}}
  {{~#if optional}},optional{{/if}}\
  {{~#if (Unboxed type)}},type={{type}}{{/if}}"\`
{{/each}}
}

{{>spec~}}
{{~#each (ParamsOf this)}}{{#unless (IsInternal label)}}
const {{Pascal ../name}}_Field_{{Pascal @key}} = "{{@key}}";
{{~/unless}}{{~/each}}
{{#if ../marshal}}
{{>sig}}
{{#unless (NoHelpers name)}}
{{>repeat name=(Pascal name) el=(Pascal name)}}
{{/unless}}

func {{Pascal name}}_Optional_Marshal(n jsn.Marshaler, pv **{{Pascal name}}) {
  if enc := n.IsEncoding(); enc && *pv != nil {
    {{Pascal name}}_Marshal(n, *pv)
  } else if !enc {
    var v {{Pascal name}}
    if {{Pascal name}}_Marshal(n, &v) {
      *pv = &v
    }
  }
}

func {{Pascal name}}_Marshal{{Custom name}}(n jsn.Marshaler, val *{{Pascal name}}) (okay bool) {
{{#if (IsPositioned this)}}
  n.SetCursor(val.At.Offset)
{{/if}}
  if okay = n.Map
    {{~#if (IsLiteral group)}}Literal{{else}}Values{{/if~}}
    ({{#if (LedeName this)~}}"{{LedeName this}}"{{~else~}}{{Pascal name}}_Type{{~/if~}}, {{Pascal name}}_Type); okay {
{{~#each (ParamsOf this)}}{{#unless (IsInternal label)}}
    if n.MapKey("{{SelectorOf @key this @index}}", {{Pascal ../name}}_Field_{{Pascal @key}}) {
      {{ScopeOf type}}{{Pascal type}}
      {{~#if (Unboxed type)}}_Unboxed{{/if}}
      {{~#if repeats}}_Repeats
      {{~else if optional}}_Optional
      {{~/if}}_Marshal(n, &val.{{Pascal @key}})
    }
{{~/unless}}{{/each}}
    n.EndValues()
  }
  return
}
{{/if}}
{{/with}}
`;
