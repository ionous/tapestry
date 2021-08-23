// swapDetailsTemplate.js
'use strict';
module.exports = `
func {{Pascal name}}_Detailed_Marshal(n jsonexp.Context, val *{{Pascal name}}) (ret []byte,err error) {
  if pick, ok := val.GetChoice(); !ok {
    err = errutil.Fmt("unknown choice %T in %s", val.Opt, {{Pascal name}}_Type)
  } else if slat := val.Opt; len(pick) > 0 {
    if b, e := slat.(jsonexp.DetailedMarshaler).MarshalDetailed(n); e != nil {
      err =  errutil.New({{Pascal name}}_Type, "-", e)
    } else {
      ret, err = json.Marshal(
      jsonexp.Flow{
{{~#if (IsPositioned this)}}
        Id: val.At.Offset,{{/if}}
        Type: {{Pascal name}}_Type,
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
    err =  errutil.New("value of", {{Pascal name}}_Type, "-", e)
  } else {
    var ptr jsonexp.DetailedMarshaler
    var raw json.RawMessage
    for k, v := range msg.Fields {
      switch k {
    {{#each with.params}}
      case {{Pascal ../name}}_{{Pascal @key}}_Opt:
        ptr = new({{TypeOf this}})
    {{/each}}
      default:
        err = errutil.New("unknown choice", k, n.Source(), msg.Id)
      }
      raw= v
      break
    }
    if ptr == nil {
      err = errutil.New("missing choice", n.Source(), msg.Id)
    } else if err == nil {
      if e := ptr.UnmarshalDetailed(n, raw); e != nil {
        err =  errutil.New("contents of", {{Pascal name}}_Type, "-", e)
      } else {
        out.Opt = ptr
{{#if (IsPositioned this)}}
        out.At = reader.Position{Source: n.Source(), Offset: msg.Id}
{{/if}}
      }
    }
  }
  return
}
`;
