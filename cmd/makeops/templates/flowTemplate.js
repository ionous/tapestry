// runTemplate.js
'use strict';
module.exports = `
{{#with type~}}
// {{Pascal name}}
{{~#if (IsString desc.short)}} {{{desc.short}}}
{{else~}}
{{#each desc.short}}
{{#if @index}}//{{/if}} {{{this}}}
{{else}}

{{/each}}
{{/if}}
{{#if with.slots}}
// User implements:{{#each with.slots}} {{Pascal this}}{{#unless @last}},{{/unless}}{{/each}}.
{{/if}}
type {{Pascal name}} struct {
{{#each params}}{{#unless embedded}}
  {{~#unless expanded}}{{Pascal key}}{{/unless}} {{TypeOf this}} \`if:"
  {{~#if internal}}internal{{else}}label={{tag}}{{/if}}
  {{~#if optional}},optional{{/if}}\
  {{~#if (Unboxed type)}},type={{type}}{{/if}}"\`
{{/unless}}{{/each}}
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

type {{Pascal name}}_Flow struct { ptr* {{Pascal name}} }

func (n {{Pascal name}}_Flow) GetType() string { return {{Pascal name}}_Type }
func (n {{Pascal name}}_Flow) GetLede() string { return {{#if (LedeName this)}}"{{LedeName this}}"{{else}}{{Pascal name}}_Type{{/if}} }
func (n {{Pascal name}}_Flow) GetFlow() interface{} { return n.ptr }
func (n {{Pascal name}}_Flow) SetFlow(i interface{}) (okay bool) {
  if ptr, ok := i.(*{{Pascal name}}); ok {
    *n.ptr, okay = *ptr, true
  }
  return
}

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
{{#if (IsPositioned this)}}
  m.SetCursor(val.At.Offset)
{{/if}}
  if err = m.MarshalBlock({{Pascal name}}_Flow{val}); err == nil {
{{~#each params}}{{#unless internal}}{{#unless expanded}}
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
{{~/unless}}{{/unless}}{{/each}}
    m.EndBlock()
  }
  return
}
{{/if}}
{{/with}}
`;
