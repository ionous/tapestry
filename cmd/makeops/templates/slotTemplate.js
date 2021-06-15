// slotTemplate.js
'use strict';
module.exports = `
{{~#with type~}}
const Type_{{Pascal name}} = "{{name}}";

{{#if ../marshal}}
func {{Pascal name}}_Detailed_Marshal(n jsonexp.Context, ptr *{{Pascal name}}) (ret []byte,err error) {
  var b []byte
  if slat := *ptr; slat != nil {
    b, err = slat.(jsonexp.DetailedMarshaler).MarshalDetailed(n)
  }
  if err == nil {
    ret, err = json.Marshal(jsonexp.Node{
      Type:  {{ScopeOf name}}Type_{{Pascal name}},
      Value: b,
    })
  }
  return
}

func {{Pascal name}}_Detailed_Unmarshal(n jsonexp.Context, b []byte, out *{{Pascal name}}) (err error) {
  if ptr, e := jsonexp.UnmarshalDetailedSlot(n, b); e != nil {
    err =e
  } else if store, ok := ptr.({{Pascal name}}); !ok && ptr != nil {
    err = errutil.Fmt("couldnt store %T into %s", ptr, Type_{{Pascal name}})
  } else {
    (*out) = store
  }
  return
}
{{/if}}
{{/with}}
{{#if repeats}}{{>repeat this.type}}{{/if}}
`;
