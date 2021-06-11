// swapTemplate.js
'use strict';
module.exports = `
{{~#with type~}}
// {{Pascal name}} swaps between various options
type {{Pascal name}} struct {
{{~#if (IsPositioned this)}}
  At  reader.Position \`if:"internal"\`{{/if}}
  Opt interface{}
}

{{#each with.params}}
const {{Pascal ../name}}_{{Pascal @key}}= "{{@key}}";
{{/each}}

{{>spec spec=this}}

func (op* {{Pascal name}}) GetChoice() (ret string, okay bool) {
  switch op.Opt.(type) {
  case nil:
    okay = true
{{#each with.params}}
  case *{{TypeOf this}}:
    ret, okay = {{Pascal ../name}}_{{Pascal @key}}, true
{{/each}}
  }
  return
}

{{~#if ../marshal}}
{{>sig sig=this}}

func {{Pascal name}}_Detailed_Marshal(n jsonexp.Context, val *{{Pascal name}}) (ret []byte, err error) {
  if pick, ok := val.GetChoice(); !ok {
    err = errutil.Fmt("unknown choice %T in %s", val.Opt, Type_{{Pascal name}})
  } else if slat := val.Opt; len(pick) > 0 {
    if b, e := slat.(jsonexp.DetailedMarshaler).MarshalDetailed(n); e!= nil {
      err =  errutil.New("marshaling", Type_{{Pascal name}}, e)
    } else {
      ret, err= json.Marshal(
      jsonexp.Flow{
{{~#if (IsPositioned this)}}
        Id: val.At.Offset,{{/if}}
        Type: Type_{{Pascal name}},
        Fields: jsonexp.Fields{
          pick: b,
        },
      })
    }
  }
  return
}

func {{Pascal name}}_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *{{Pascal name}}) (err error) {
  var msg jsonexp.Flow
  if e := json.Unmarshal(b, &msg); e != nil {
    err =  errutil.New("unmarshaling msg", Type_{{Pascal name}}, e)
  } else {
    var ptr jsonexp.DetailedMarshaler
    var raw json.RawMessage
    for k, v := range msg.Fields {
      switch k {
    {{#each with.params}}
      case {{Pascal ../name}}_{{Pascal @key}}:
        ptr = new({{TypeOf this}})
    {{/each}}
      default:
        err = errutil.New("unknown choice", k, n.Source(), msg.Id)
      }
      raw= v
      break
    }
    if ptr == nil {
      err= errutil.New("missing choice", n.Source(), msg.Id)
    } else if err == nil {
      if e := ptr.UnmarshalDetailed(n, raw); e != nil {
        err =  errutil.New("unmarshaling contents", Type_{{Pascal name}}, e)
      } else {
        out.Opt = ptr
{{#if (IsPositioned this)}}
        out.At = reader.Position{Source:n.Source(), Offset: msg.Id}
{{/if}}
      }
    }
  }
  return
}
{{/if}}
{{/with}}
{{#if marshal}}
{{>repeat spec=this}}
{{/if}}
`;
