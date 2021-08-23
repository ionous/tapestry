// flowDetails.js
'use strict';
// Flow gets the name capitalized, plus all of the parameters colon separated as a single key object.
// probably need a sigbuilder b/c of the optional parameters.
// explicit lede over name... but where does that even come from?
// specTemplate
module.exports = `
func {{Pascal name}}_Compact_Optional_Marshal(n jsonexp.Context, val **{{Pascal name}}) (ret []byte,err error) {
  if *val != nil {
    ret, err= {{Pascal name}}_Compact_Marshal(n, *val)
  }
  return
}
func {{Pascal name}}_Compact_Marshal(n jsonexp.Context, val *{{Pascal name}}) (ret []byte, err error) {
  var sig jsonexp.Sig
  var fields []json.RawMessage
  sig.WriteLede({{Pascal name}}_Lede)
{{~#each (ParamsOf this)}}{{#unless (IsInternal label)}}
  if b, e := {{ScopeOf type}}{{Pascal type}}{{#if (OverrideOf type)}}_Override{{/if}}_Compact{{#if repeats}}_Repeats{{else if optional}}_Optional{{/if}}_Marshal(n, &val.{{Pascal @key}}); e != nil {
    err = errutil.Append(err, e)
  } else{{#if optional}} if len(b) > 0{{/if}} {
    sig.WriteLabel("{{SelectorOf @key this @index}}")
    fields = append(fields, b)
  }
{{/unless}}{{/each}}
  if err == nil {
    ret, err = json.Marshal(map[string]interface{}{
{{~#if (IsPositioned this)}}{{#if val.At.Offset}}
    "id": val.At.Offset,{{/if}}{{/if}}
      sig.String(): fields,
    })
  }
  return
}

func {{Pascal name}}_Compact_Optional_Unmarshal(n jsonexp.Context, b []byte, out **{{Pascal name}}) (err error) {
  if len(b) > 0 {
    var val {{Pascal name}}
    if e:= {{Pascal name}}_Compact_Unmarshal(n, b, &val); e!= nil {
      err = e
    } else {
      *out = &val
    }
  }
  return
}
func {{Pascal name}}_Compact_Unmarshal(n jsonexp.Context, b []byte, out *{{Pascal name}}) (err error) {
  var msg jsonexp.Flow
  if e := json.Unmarshal(b, &msg); e != nil {
    err = errutil.New({{Pascal name}}_Type, "-", e)
  } {{#each (ParamsOf this)}}{{#unless (IsInternal label)~}}
  else if e := {{ScopeOf type}}{{Pascal type}}{{#if (OverrideOf type)}}_Override{{/if}}_Compact{{#if repeats}}_Repeats{{else if optional}}_Optional{{/if}}_Unmarshal(n, msg.Fields[{{Pascal ../name}}_Field_{{Pascal @key}}], &out.{{Pascal @key}}); e != nil {
    err = errutil.New({{Pascal ../name}}_Type + "." + {{Pascal ../name}}_Field_{{Pascal @key}}, "-", e)
  } {{/unless}}{{/each~}}
{{~#if (IsPositioned this)~}}
  else {
      out.At = reader.Position{Source: n.Source(), Offset: msg.Id}
  }{{/if}}
  return
}
`
