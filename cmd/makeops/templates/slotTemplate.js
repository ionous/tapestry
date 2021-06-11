// slotTemplate.js
'use strict';
module.exports = `
{{~#with type~}}
const Type_{{Pascal name}} = "{{name}}";

{{#if ../marshal}}
func {{Pascal name}}_Detailed_Marshal(n jsonexp.Context, ptr *{{Pascal name}}) (ret []byte, err error) {
  var b []byte
  if slot:= *ptr; slot != nil {
    b, err = slot.(jsonexp.DetailedMarshaler).MarshalDetailed(n)
  }
  if err == nil {
    ret, err= json.Marshal(jsonexp.Node{
      Type:  {{ScopeOf name}}Type_{{Pascal name}},
      Value: b,
    })
  }
  return
}

func {{Pascal name}}_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *{{Pascal name}}) (err error) {
  var msg jsonexp.Node
  if e := json.Unmarshal(b, &msg); e != nil {
    err= e
  } else if ptr, e := n.NewType(msg.Type); e != nil {
    err= e
  } else if e := ptr.UnmarshalDetailed(n, msg.Value); e != nil {
    err = e
  } else {
     (*out) = ptr.({{Pascal name}})
  }
  return
}
{{/if}}
{{/with}}
{{#if marshal}}
{{>repeat spec=this}}
{{/if}}
`;
