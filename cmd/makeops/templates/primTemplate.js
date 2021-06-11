// strTemplate.js
'use strict';
module.exports = `
{{~#with type~}}
// {{Pascal name}} requires a user-specified string.
type {{Pascal name}} struct {
{{~#if (IsPositioned this)}}
  At  reader.Position \`if:"internal"\`{{/if}}
  {{Pascal uses}} {{#if (Uses name 'num')}}float64{{else}}string{{/if}}
}
{{#if (Uses name 'str')}}
func (op *{{Pascal name}}) String()(ret string) {
  return op.Str
}
{{/if}}

{{#each (Choices @this)}}
const {{Pascal ../name}}_{{Pascal this.token}}= "{{this.token}}";
{{/each~}}

{{>spec spec=this}}
{{#if ../marshal}}
{{>sig sig=this}}

{{~#if (OverrideOf name)}}

func {{Pascal name}}_Detailed_Override_Marshal(n jsonexp.Context, val *{{OverrideOf name}}) ([]byte, error) {
{{#unless (IsBool name)}}
  return {{Pascal name}}_Detailed_Marshal(n, &{{Pascal name}}{*val})
{{else}}
  var out string
  if *val {
    out = Bool_True
  } else {
    out = Bool_False
  }
  return {{Pascal name}}_Detailed_Marshal(n, &{{Pascal name}}{out})
{{/unless}}
}

func {{Pascal name}}_Detailed_Override_Unmarshal(n jsonexp.Context, b []byte, out *{{OverrideOf name}}) (err error) {
  var msg {{Pascal name}}
  if e:= {{Pascal name}}_Detailed_Unmarshal(n, b, &msg); e!= nil {
    err = errutil.New("unmarshaling", Type_{{Pascal name}}, e)
  } else {
{{#unless (IsBool name)}}
    *out= msg.{{Pascal uses}}
{{else}}
    *out = msg.{{Pascal uses}} == Bool_True
{{/unless}}
  }
  return
}
{{#if ../repeats}}
func {{Pascal name}}_Detailed_Override_Repeats_Marshal(n jsonexp.Context, vals *[]{{OverrideOf name}}) (ret []byte, err error) {
  var msgs []json.RawMessage
  msgs= make([]json.RawMessage, len(*vals))
  for i, el:= range *vals {
    if b, e:= {{Pascal name}}_Detailed_Override_Marshal(n, &el); e!= nil {
      err = errutil.New("marshaling", Type_{{Pascal name}}, "at", i, e)
      break
    } else {
      msgs[i]= b
    }
  }
  if err== nil {
    ret, err= json.Marshal(msgs)
  }
  return
}

func {{Pascal name}}_Detailed_Override_Repeats_Unmarshal(n jsonexp.Context, b []byte, out *[]{{OverrideOf name}}) (err error) {
  var msgs []json.RawMessage
  if e:= json.Unmarshal(b, &msgs); e!= nil  {
    err = errutil.New("unmarshaling", Type_{{Pascal name}}, e)
  } else {
    vals:= make([]{{OverrideOf name}}, len(msgs))
    for i, msg:= range msgs {
      if e:= {{Pascal name}}_Detailed_Override_Unmarshal(n, msg, &vals[i]); e!= nil {
        err = errutil.New("unmarshaling", Type_{{Pascal name}}, "at", i, e)
        break
      }
    }
    if err== nil {
      *out= vals
    }
  }
  return
}
{{/if}}
{{/if~}}

func {{Pascal name}}_Detailed_Marshal(n jsonexp.Context, val *{{Pascal name}}) ([]byte, error) {
  return json.Marshal(jsonexp.{{Pascal uses}}{
{{#if (IsPositioned this)}}
    Id: val.At.Offset,
{{/if}}
    Type:  Type_{{Pascal name}},
    Value: val.{{Pascal uses}},
  })
}

func {{Pascal name}}_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *{{Pascal name}}) (err error) {
  var msg jsonexp.{{Pascal uses}}
  if e:= json.Unmarshal(b, &msg); e!= nil {
    err =  errutil.New("unmarshaling", Type_{{Pascal name}}, e)
  } else {
{{#if (IsPositioned this)~}}
    out.At = reader.Position{ Source:n.Source(), Offset: msg.Id }
{{/if}}
    out.{{Pascal uses}}= msg.Value
  }
  return
}
{{/if}}
{{/with}}
{{#if marshal}}
{{>repeat spec=this}}
{{/if}}
`;
