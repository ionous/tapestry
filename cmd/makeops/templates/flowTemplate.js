// runTemplate.js
'use strict';
module.exports = `
{{~#with type~}}
// {{Pascal name}} {{desc.short}}
{{#if with.slots}}
// User implements:{{#each with.slots}} {{Pascal this}}{{#unless @last}},{{/unless}}{{/each}}.
{{/if}}
type {{Pascal name}} struct {
{{#each params}}
  {{Pascal key}} {{TypeOf this}} \`if:"
  {{~#if internal}}internal{{else}}label={{tag}}{{/if}}
  {{~#if optional}},optional{{/if}}\
  {{~#if (Unboxed type)}},type={{type}}{{/if}}"\`
{{/each}}
}

{{>spec~}}
{{~#each params}}{{#unless internal}}
const {{Pascal ../name}}_Field_{{Pascal key}} = "{{key}}";
{{~/unless}}{{~/each}}
{{#if ../marshal}}
{{>sig}}
{{#unless (NoHelpers name)}}
{{>repeat name=(Pascal name) el=(Pascal name)}}
{{/unless}}

func {{Pascal name}}_Optional_Marshal(m jsn.Marshaler, pv **{{Pascal name}}) (err error) {
  if enc := m.IsEncoding(); enc && *pv != nil {
    err = {{Pascal name}}_Marshal(m, *pv)
  } else if !enc {
    var v {{Pascal name}}
    if err = {{Pascal name}}_Marshal(m, &v); err == nil {
      *pv = &v
    }
  }
  return
}

func {{Pascal name}}_Marshal(m jsn.Marshaler, val *{{Pascal name}}) (err error) {
{{#if (IsCustom name)}}
  if fn, ok := m.CustomizedMarshal({{Pascal name}}_Type); ok {
    err = fn(m, val)
  } else {
    err = {{Pascal name}}_DefaultMarshal(m, val)
  }
  return
}
func {{Pascal name}}_DefaultMarshal(m jsn.Marshaler, val *{{Pascal name}}) (err error) {
{{/if}}
{{#if (IsPositioned this)}}
  m.SetCursor(val.At.Offset)
{{/if}}
  if err = m.MarshalBlock(jsn.MakeFlow(
{{~#if (LedeName this)}}"{{LedeName this}}"{{else}}{{Pascal name}}_Type{{/if
}}, {{Pascal name}}_Type, val)); err == nil {
{{~#each params}}{{#unless (IsInternal label)}}
    e{{@index}} := m.MarshalKey("{{sel}}", {{Pascal ../name}}_Field_{{Pascal key}})
    if e{{@index}} == nil {
      e{{@index}} = {{ScopeOf type}}{{Pascal type~}}
      {{#if (Unboxed type)}}_Unboxed{{/if~}}
      {{#if optional}}_Optional{{/if~}}
      {{#if repeats}}_Repeats{{/if~}}
      _Marshal(m, &val.{{Pascal key}})
    }
    if e{{@index}} != nil && e{{@index}} != jsn.Missing {
      m.Error(errutil.New(e{{@index}}, "in flow at", {{Pascal ../name}}_Field_{{Pascal key}}))
    }
{{~/unless}}{{/each}}
    m.EndBlock()
  }
  return
}
{{/if}}
{{/with}}
`;
